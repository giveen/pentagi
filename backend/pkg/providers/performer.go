package providers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"pentagi/pkg/csum"
	"pentagi/pkg/database"
	obs "pentagi/pkg/observability"
	"pentagi/pkg/observability/langfuse"
	"pentagi/pkg/providers/autotuner"
	"pentagi/pkg/providers/pconfig"
	"pentagi/pkg/tools"

	"github.com/sirupsen/logrus"
	"github.com/vxcontrol/langchaingo/llms"
	"github.com/vxcontrol/langchaingo/llms/reasoning"
	"github.com/vxcontrol/langchaingo/llms/streaming"
)

const (
	maxRetriesToCallSimpleChain    = 3
	maxRetriesToCallAgentChain     = 2
	maxRetriesToCallFunction       = 3
	maxReflectorCallsPerChain      = 3
	maxGeneralAgentChainIterations = 100
	maxLimitedAgentChainIterations = 20
	maxAgentShutdownIterations     = 3
	maxSoftDetectionsBeforeAbort   = 4
	delayBetweenRetries            = 2 * time.Second
)

var nonRepairableToolCallErrors = []string{
	"failed to update toolcall result:",
}

type callResult struct {
	streamID  int64
	funcCalls []llms.ToolCall
	info      map[string]any
	thinking  *reasoning.ContentReasoning
	content   string
}

func (fp *flowProvider) performAgentChain(
	ctx context.Context,
	optAgentType pconfig.ProviderOptionsType,
	chainID int64,
	taskID, subtaskID *int64,
	chain []llms.MessageContent,
	executor tools.ContextToolsExecutor,
	summarizer csum.Summarizer,
) error {
	ctx, span := obs.Observer.NewSpan(ctx, obs.SpanKindInternal, "providers.flowProvider.performAgentChain")
	defer span.End()

	var (
		wantToStop        bool
		monitor           = fp.buildMonitor()
		detector          = &repeatingDetector{}
		summarizerHandler = fp.GetSummarizeResultHandler(taskID, subtaskID)
	)

	logger := logrus.WithContext(ctx).WithFields(enrichLogrusFields(fp.flowID, taskID, subtaskID, logrus.Fields{
		"provider":     fp.Type(),
		"agent":        optAgentType,
		"msg_chain_id": chainID,
	}))

	// Track execution time for duration calculation
	lastUpdateTime := time.Now()
	rollLastUpdateTime := func() float64 {
		durationDelta := time.Since(lastUpdateTime).Seconds()
		lastUpdateTime = time.Now()
		return durationDelta
	}

	executionContext, err := fp.getExecutionContext(ctx, taskID, subtaskID)
	if err != nil {
		logger.WithError(err).Error("failed to get execution context")
		return fmt.Errorf("failed to get execution context: %w", err)
	}

	groupID := fmt.Sprintf("flow-%d", fp.flowID)
	toolTypeMapping := tools.GetToolTypeMapping()

	var maxCallsLimit int
	switch optAgentType {
	case pconfig.OptionsTypeAssistant, pconfig.OptionsTypePrimaryAgent,
		pconfig.OptionsTypePentester, pconfig.OptionsTypeCoder, pconfig.OptionsTypeInstaller:
		if fp.maxGACallsLimit <= 0 {
			maxCallsLimit = maxGeneralAgentChainIterations
		} else {
			maxCallsLimit = max(fp.maxGACallsLimit, maxAgentShutdownIterations*2)
		}
	default:
		if fp.maxLACallsLimit <= 0 {
			maxCallsLimit = maxLimitedAgentChainIterations
		} else {
			maxCallsLimit = max(fp.maxLACallsLimit, maxAgentShutdownIterations*2)
		}
	}

	for iteration := 0; ; iteration++ {
		if iteration >= maxCallsLimit {
			msg := fmt.Sprintf("agent chain exceeded maximum iterations (%d)", maxCallsLimit)
			logger.WithField("iteration", iteration).Error(msg)
			return errors.New(msg)
		}

		var result *callResult
		if iteration >= maxCallsLimit-maxAgentShutdownIterations {
			logger.WithFields(logrus.Fields{
				"iteration": iteration,
				"limit":     maxCallsLimit,
			}).Warn("max tool calls limit will be reached soon, invoking reflector for graceful termination")

			// Format reflector message for graceful termination
			result = &callResult{
				content: fmt.Sprintf(
					"I can’t continue this multi-turn chain because I’m too close to the AI agent iteration limit (%d).",
					maxCallsLimit,
				),
			}
		} else {
			result, err = fp.callWithRetries(ctx, optAgentType, chainID, taskID, subtaskID, chain, executor, executionContext)
			if err != nil {
				logger.WithError(err).Error("failed to call agent chain")
				return err
			}

			if err := fp.updateMsgChainUsage(ctx, chainID, optAgentType, result.info, rollLastUpdateTime()); err != nil {
				logger.WithError(err).Error("failed to update msg chain usage")
				return err
			}
		}

		if len(result.funcCalls) == 0 {
			if optAgentType == pconfig.OptionsTypeAssistant {
				fp.storeAgentResponseToGraphiti(ctx, groupID, optAgentType, result, taskID, subtaskID, chainID)
				return fp.processAssistantResult(ctx, logger, chainID, chain, result, summarizer, summarizerHandler, rollLastUpdateTime())
			} else {
				// Build AI message with reasoning for reflector (universal pattern)
				reflectorMsg := llms.MessageContent{Role: llms.ChatMessageTypeAI}
				if result.content != "" || !result.thinking.IsEmpty() {
					reflectorMsg.Parts = append(reflectorMsg.Parts, llms.TextPartWithReasoning(result.content, result.thinking))
				}
				result, err = fp.performReflector(
					ctx, optAgentType, chainID, taskID, subtaskID,
					append(chain, reflectorMsg), executor,
					fp.getLastHumanMessage(chain), result.content, executionContext, 1)
				if err != nil {
					fields := logrus.Fields{}
					if result != nil {
						fields["content"] = result.content[:min(1000, len(result.content))]
						if !result.thinking.IsEmpty() {
							fields["thinking"] = result.thinking.Content[:min(1000, len(result.thinking.Content))]
						}
						fields["execution"] = executionContext[:min(1000, len(executionContext))]
					}
					logger.WithError(err).WithFields(fields).Error("failed to perform reflector")
					return err
				}
			}
		}

		fp.storeAgentResponseToGraphiti(ctx, groupID, optAgentType, result, taskID, subtaskID, chainID)

		msg := llms.MessageContent{Role: llms.ChatMessageTypeAI}
		// Universal pattern: preserve content with or without reasoning (works for all providers thanks to deduplication)
		if result.content != "" || !result.thinking.IsEmpty() {
			msg.Parts = append(msg.Parts, llms.TextPartWithReasoning(result.content, result.thinking))
		}
		for _, toolCall := range result.funcCalls {
			msg.Parts = append(msg.Parts, toolCall)
		}
		chain = append(chain, msg)

		if err := fp.updateMsgChain(ctx, optAgentType, chainID, chain, rollLastUpdateTime()); err != nil {
			logger.WithError(err).Error("failed to update msg chain")
			return err
		}

		for idx, toolCall := range result.funcCalls {
			if toolCall.FunctionCall == nil {
				continue
			}

			funcName := toolCall.FunctionCall.Name
			response, err := fp.execToolCall(
				ctx, optAgentType, chainID, idx, result, monitor, detector, executor, taskID, subtaskID, chain,
			)

			if toolTypeMapping[funcName] != tools.AgentToolType {
				fp.storeToolExecutionToGraphiti(
					ctx, groupID, optAgentType, toolCall, response, err, executor, taskID, subtaskID, chainID,
				)
			}

			if err != nil {
				logger.WithError(err).WithFields(logrus.Fields{
					"func_name": funcName,
					"func_args": toolCall.FunctionCall.Arguments,
				}).Error("failed to exec tool call")
				return err
			}

			chain = append(chain, llms.MessageContent{
				Role: llms.ChatMessageTypeTool,
				Parts: []llms.ContentPart{
					llms.ToolCallResponse{
						ToolCallID: toolCall.ID,
						Name:       funcName,
						Content:    response,
					},
				},
			})
			if err := fp.updateMsgChain(ctx, optAgentType, chainID, chain, rollLastUpdateTime()); err != nil {
				logger.WithError(err).Error("failed to update msg chain")
				return err
			}

			if executor.IsBarrierFunction(funcName) {
				wantToStop = true
			}
		}

		if wantToStop {
			return nil
		}

		if summarizer != nil {
			// it returns the same chain state if error occurs
			chain, err = summarizer.SummarizeChain(ctx, summarizerHandler, chain, fp.tcIDTemplate)
			if err != nil {
				// log swallowed error
				_, observation := obs.Observer.NewObservation(ctx)
				observation.Event(
					langfuse.WithEventName("chain summarization error swallowed"),
					langfuse.WithEventInput(chain),
					langfuse.WithEventStatus(err.Error()),
					langfuse.WithEventLevel(langfuse.ObservationLevelWarning),
					langfuse.WithEventMetadata(langfuse.Metadata{
						"tc_id_template": fp.tcIDTemplate,
						"msg_chain_id":   chainID,
						"error":          err.Error(),
					}),
				)
				logger.WithError(err).Warn("failed to summarize chain")
			} else if err := fp.updateMsgChain(ctx, optAgentType, chainID, chain, rollLastUpdateTime()); err != nil {
				logger.WithError(err).Error("failed to update msg chain")
				return err
			}
		}
	}
}

func (fp *flowProvider) execToolCall(
	ctx context.Context,
	optAgentType pconfig.ProviderOptionsType,
	chainID int64,
	toolCallIDx int,
	result *callResult,
	monitor *executionMonitor,
	detector *repeatingDetector,
	executor tools.ContextToolsExecutor,
	taskID, subtaskID *int64,
	chain []llms.MessageContent,
) (string, error) {
	var (
		streamID int64
		thinking string
	)

	// use streamID and thinking only for first tool call to minimize content
	if toolCallIDx == 0 {
		streamID = result.streamID
		if !result.thinking.IsEmpty() {
			thinking = result.thinking.Content
		}
	}

	toolCall := result.funcCalls[toolCallIDx]
	if toolCall.FunctionCall == nil {
		return "", fmt.Errorf("tool call function call is nil")
	}

	funcName := toolCall.FunctionCall.Name
	funcArgs := json.RawMessage(toolCall.FunctionCall.Arguments)

	logger := logrus.WithContext(ctx).WithFields(enrichLogrusFields(fp.flowID, taskID, subtaskID, logrus.Fields{
		"agent":        fp.Type(),
		"func_name":    funcName,
		"func_args":    string(funcArgs)[:min(1000, len(funcArgs))],
		"tool_call_id": toolCall.ID,
		"msg_chain_id": chainID,
	}))

	if detector.detect(toolCall) {
		if len(detector.funcCalls) >= RepeatingToolCallThreshold+maxSoftDetectionsBeforeAbort {
			errMsg := fmt.Sprintf("tool '%s' repeated %d times consecutively, aborting chain", funcName, len(detector.funcCalls))
			logger.WithField("repeat_count", len(detector.funcCalls)).Error(errMsg)
			return "", errors.New(errMsg)
		}

		response := fmt.Sprintf("tool call '%s' is repeating, please try another tool", funcName)

		_, observation := obs.Observer.NewObservation(ctx)
		observation.Event(
			langfuse.WithEventName("repeating tool call detected"),
			langfuse.WithEventInput(funcArgs),
			langfuse.WithEventMetadata(map[string]any{
				"tool_call_id": toolCall.ID,
				"tool_name":    funcName,
				"msg_chain_id": chainID,
			}),
			langfuse.WithEventStatus("failed"),
			langfuse.WithEventLevel(langfuse.ObservationLevelError),
			langfuse.WithEventOutput(response),
		)
		logger.Warn("failed to exec function: tool call is repeating")

		return response, nil
	}

	var (
		err      error
		response string
	)

	for idx := 0; idx <= maxRetriesToCallFunction; idx++ {
		if idx == maxRetriesToCallFunction {
			err = fmt.Errorf("reached max retries to call function: %w", err)
			logger.WithError(err).Error("failed to exec function")
			return "", fmt.Errorf("failed to exec function '%s': %w", funcName, err)
		}

		response, err = executor.Execute(ctx, streamID, toolCall.ID, funcName, funcName, thinking, funcArgs)
		if err != nil {
			if errors.Is(err, context.Canceled) {
				return "", err
			}

			if !shouldRepairToolCallArgs(err) {
				logger.WithError(err).Warn("failed to exec function with non-repairable error")
				return "", fmt.Errorf("failed to exec function '%s': %w", funcName, err)
			}

			logger.WithError(err).Warn("failed to exec function")

			funcExecErr := err
			funcSchema, err := executor.GetToolSchema(funcName)
			if err != nil {
				logger.WithError(err).Error("failed to get tool schema")
				return "", fmt.Errorf("failed to get tool schema: %w", err)
			}

			funcArgs, err = fp.fixToolCallArgs(ctx, funcName, funcArgs, funcSchema, funcExecErr)
			if err != nil {
				logger.WithError(err).Error("failed to fix tool call args")
				return "", fmt.Errorf("failed to fix tool call args: %w", err)
			}
		} else {
			break
		}
	}

	if monitor.shouldInvokeMentor(toolCall) && executor.IsFunctionExists(tools.AdviceToolName) {
		logger.WithFields(logrus.Fields{
			"same_tool_count":  monitor.sameToolCount,
			"total_call_count": monitor.totalCallCount,
		}).Debug("execution monitor threshold reached, invoking mentor for progress review")

		mentorResponse, err := fp.performMentor(
			ctx, optAgentType, chainID, taskID, subtaskID, chain, executor, toolCall, response,
		)
		if err != nil {
			logger.WithError(err).Warn("failed to invoke execution mentor, continuing with normal execution")
		} else {
			monitor.reset()
			response = formatEnhancedToolResponse(response, mentorResponse)
		}
	}

	return response, nil
}

func isJSONParseError(err error) bool {
	return err != nil && strings.Contains(err.Error(), "Failed to parse tool call arguments as JSON")
}

func shouldRepairToolCallArgs(err error) bool {
	if err == nil {
		return false
	}

	errMsg := strings.ToLower(err.Error())
	for _, marker := range nonRepairableToolCallErrors {
		if strings.Contains(errMsg, strings.ToLower(marker)) {
			return false
		}
	}

	return true
}

func (fp *flowProvider) callWithRetries(
	ctx context.Context,
	optAgentType pconfig.ProviderOptionsType,
	chainID int64,
	taskID, subtaskID *int64,
	chain []llms.MessageContent,
	executor tools.ContextToolsExecutor,
	executionContext string,
) (*callResult, error) {
	var (
		err     error
		errs    []error
		msgType = database.MsglogTypeAnswer
		resp    *llms.ContentResponse
		result  callResult
	)

	logger := logrus.WithContext(ctx).WithFields(enrichLogrusFields(fp.flowID, taskID, subtaskID, logrus.Fields{
		"agent":        fp.Type(),
		"msg_chain_id": chainID,
		"agent_type":   optAgentType,
	}))

	ticker := time.NewTicker(delayBetweenRetries)
	defer ticker.Stop()

	fillResult := func(resp *llms.ContentResponse) error {
		var stopReason string
		var parts []string

		if resp == nil || len(resp.Choices) == 0 {
			return fmt.Errorf("no choices in response")
		}

		for _, choice := range resp.Choices {
			if stopReason == "" {
				stopReason = choice.StopReason
			}

			if choice.GenerationInfo != nil {
				result.info = choice.GenerationInfo
			}

			// Extract reasoning for logging/analytics (provider-aware)
			if result.thinking.IsEmpty() {
				if !choice.Reasoning.IsEmpty() {
					result.thinking = choice.Reasoning
				} else if len(choice.ToolCalls) > 0 && !choice.ToolCalls[0].Reasoning.IsEmpty() {
					// Gemini puts reasoning in first tool call when tools are used
					result.thinking = choice.ToolCalls[0].Reasoning
				}
			}

			if strings.TrimSpace(choice.Content) != "" {
				parts = append(parts, choice.Content)
			}

			for _, toolCall := range choice.ToolCalls {
				if toolCall.FunctionCall == nil {
					continue
				}
				result.funcCalls = append(result.funcCalls, toolCall)
			}
		}

		result.content = strings.Join(parts, "\n")
		if strings.Trim(result.content, "' \"\n\r\t") == "" && len(result.funcCalls) == 0 {
			return fmt.Errorf("no content and tool calls in response: stop reason '%s'", stopReason)
		}

		return nil
	}

	for idx := 0; idx <= maxRetriesToCallAgentChain; idx++ {
		if idx == maxRetriesToCallAgentChain {
			reflectorResult, err := fp.performCallerReflector(
				ctx, optAgentType, chainID, taskID, subtaskID, chain, executor, executionContext, errs,
			)
			if err != nil {
				msg := fmt.Sprintf("failed to call agent chain: max retries reached, %d", idx)
				return nil, fmt.Errorf(msg+": %w", errors.Join(append(errs, err)...))
			}

			return reflectorResult, nil
		}

		var streamCb streaming.Callback
		if fp.streamCb != nil {
			result.streamID = fp.callCounter.Add(1)
			streamCb = func(ctx context.Context, chunk streaming.Chunk) error {
				switch chunk.Type {
				case streaming.ChunkTypeReasoning:
					if chunk.Reasoning.IsEmpty() {
						return nil
					}
					return fp.streamCb(ctx, &StreamMessageChunk{
						Type:     StreamMessageChunkTypeThinking,
						MsgType:  msgType,
						Thinking: chunk.Reasoning,
						StreamID: result.streamID,
					})
				case streaming.ChunkTypeText:
					return fp.streamCb(ctx, &StreamMessageChunk{
						Type:     StreamMessageChunkTypeContent,
						MsgType:  msgType,
						Content:  chunk.Content,
						StreamID: result.streamID,
					})
				case streaming.ChunkTypeToolCall:
					// skip tool call chunks (we don't need them for now)
				case streaming.ChunkTypeDone:
					return fp.streamCb(ctx, &StreamMessageChunk{
						Type:     StreamMessageChunkTypeFlush,
						MsgType:  msgType,
						StreamID: result.streamID,
					})
				}
				return nil
			}
		}

		// On retries after a JSON parse error, inject a targeted hint so the model
		// knows to use single-quotes instead of double-quotes in shell values.
		callChain := chain
		if idx > 0 && len(errs) > 0 && isJSONParseError(errs[len(errs)-1]) {
			callChain = append(make([]llms.MessageContent, 0, len(chain)+1), chain...)
			callChain = append(callChain, llms.TextParts(llms.ChatMessageTypeHuman,
				"SYSTEM HINT: Your previous tool call could not be decoded because it contained "+
					"a JSON string with unescaped double-quotes (e.g., FAKETIME=\"value with spaces\"). "+
					"Use single-quotes for any value that contains spaces or special characters "+
					"(e.g., FAKETIME='2026-04-26 10:00:00'). "+
					"Alternatively, supply environment variables via the 'env' map field instead of 'export' statements. "+
					"Please retry your tool call with corrected argument formatting."))
		}

		// Thread the attempt count so WrapGenerateContent can activate sober mode.
		callCtx := autotuner.WithAttemptCount(ctx, idx+1)
		resp, err = fp.CallWithTools(callCtx, optAgentType, callChain, executor.Tools(), streamCb)
		if err == nil {
			err = fillResult(resp)
		}
		if err == nil {
			break
		} else {
			errs = append(errs, err)
			logger.WithFields(logrus.Fields{
				"retry_iteration": idx,
				"error":           err.Error()[:min(200, len(err.Error()))],
			}).Warn("agent chain call failed, will retry")
		}

		ticker.Reset(delayBetweenRetries)
		select {
		case <-ticker.C:
		case <-ctx.Done():
			return nil, fmt.Errorf("context canceled while waiting for retry: %w", ctx.Err())
		}
	}

	if fp.streamCb != nil && result.streamID != 0 {
		fp.streamCb(ctx, &StreamMessageChunk{
			Type:     StreamMessageChunkTypeUpdate,
			MsgType:  msgType,
			Content:  result.content,
			Thinking: result.thinking,
			StreamID: result.streamID,
		})
		// don't update stream by ID if we got content separately from tool calls
		// because we stored thinking and content into standalone messages
		if len(result.funcCalls) > 0 && result.content != "" {
			result.streamID = 0
		}
	}

	return &result, nil
}

