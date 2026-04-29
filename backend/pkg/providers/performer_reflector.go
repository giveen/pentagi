package providers

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"strings"

	"pentagi/pkg/cast"
	"pentagi/pkg/database"
	obs "pentagi/pkg/observability"
	"pentagi/pkg/observability/langfuse"
	"pentagi/pkg/providers/pconfig"
	"pentagi/pkg/templates"
	"pentagi/pkg/tools"

	"github.com/sirupsen/logrus"
	"github.com/vxcontrol/langchaingo/llms"
)

func (fp *flowProvider) performReflector(
	ctx context.Context,
	optOriginType pconfig.ProviderOptionsType,
	chainID int64,
	taskID, subtaskID *int64,
	chain []llms.MessageContent,
	executor tools.ContextToolsExecutor,
	humanMessage, content, executionContext string,
	iteration int,
) (*callResult, error) {
	ctx, span := obs.Observer.NewSpan(ctx, obs.SpanKindInternal, "providers.flowProvider.performReflector")
	defer span.End()

	var (
		optAgentType = pconfig.OptionsTypeReflector
		msgChainType = database.MsgchainTypeReflector
	)

	logger := logrus.WithContext(ctx).WithFields(enrichLogrusFields(fp.flowID, taskID, subtaskID, logrus.Fields{
		"provider":     fp.Type(),
		"agent":        optAgentType,
		"origin":       optOriginType,
		"msg_chain_id": chainID,
		"iteration":    iteration,
	}))

	if iteration > maxReflectorCallsPerChain {
		msg := "reflector called too many times"
		_, observation := obs.Observer.NewObservation(ctx)
		observation.Event(
			langfuse.WithEventName("reflector limit calls reached"),
			langfuse.WithEventInput(content),
			langfuse.WithEventStatus("failed"),
			langfuse.WithEventLevel(langfuse.ObservationLevelError),
			langfuse.WithEventOutput(msg),
			langfuse.WithEventMetadata(map[string]any{
				"iteration": iteration,
			}),
		)
		logger.WithField("content", content[:min(1000, len(content))]).Warn(msg)
		return nil, errors.New(msg)
	}

	logger.WithField("content", content[:min(1000, len(content))]).Warn("got message instead of tool call")

	reflectorContext := map[string]map[string]any{
		"user": {
			"Message":          content,
			"BarrierToolNames": executor.GetBarrierToolNames(),
		},
		"system": {
			"BarrierTools":     executor.GetBarrierTools(),
			"CurrentTime":      getCurrentTime(),
			"ExecutionContext": executionContext,
		},
	}

	if humanMessage != "" {
		reflectorContext["system"]["Request"] = humanMessage
	}

	ctx, observation := obs.Observer.NewObservation(ctx)
	reflectorAgent := observation.Agent(
		langfuse.WithAgentName("reflector"),
		langfuse.WithAgentInput(content),
		langfuse.WithAgentMetadata(langfuse.Metadata{
			"user_context":   reflectorContext["user"],
			"system_context": reflectorContext["system"],
		}),
	)
	ctx, observation = reflectorAgent.Observation(ctx)

	reflectorEvaluator := observation.Evaluator(
		langfuse.WithEvaluatorName("render reflector agent prompts"),
		langfuse.WithEvaluatorInput(reflectorContext),
		langfuse.WithEvaluatorMetadata(langfuse.Metadata{
			"user_context":   reflectorContext["user"],
			"system_context": reflectorContext["system"],
			"lang":           fp.language,
		}),
	)

	userReflectorTmpl, err := fp.prompter.RenderTemplate(templates.PromptTypeQuestionReflector, reflectorContext["user"])
	if err != nil {
		msg := "failed to get user reflector template"
		return nil, wrapErrorEndEvaluatorSpan(ctx, reflectorEvaluator, msg, err)
	}

	systemReflectorTmpl, err := fp.prompter.RenderTemplate(templates.PromptTypeReflector, reflectorContext["system"])
	if err != nil {
		msg := "failed to get system reflector template"
		return nil, wrapErrorEndEvaluatorSpan(ctx, reflectorEvaluator, msg, err)
	}

	reflectorEvaluator.End(
		langfuse.WithEvaluatorOutput(map[string]any{
			"user_template":   userReflectorTmpl,
			"system_template": systemReflectorTmpl,
		}),
		langfuse.WithEvaluatorStatus("success"),
		langfuse.WithEvaluatorLevel(langfuse.ObservationLevelDebug),
	)

	advice, err := fp.performSimpleChain(ctx, taskID, subtaskID, optAgentType,
		msgChainType, systemReflectorTmpl, userReflectorTmpl)
	if err != nil {
		advice = ToolPlaceholder
	}

	opts := []langfuse.AgentOption{
		langfuse.WithAgentStatus("failed"),
		langfuse.WithAgentOutput(advice),
		langfuse.WithAgentLevel(langfuse.ObservationLevelWarning),
	}
	defer func() {
		reflectorAgent.End(opts...)
	}()

	chain = append(chain, llms.TextParts(llms.ChatMessageTypeHuman, advice))
	result, err := fp.callWithRetries(ctx, optOriginType, chainID, taskID, subtaskID, chain, executor, executionContext)
	if err != nil {
		logger.WithError(err).Error("failed to call agent chain by reflector")
		opts = append(opts,
			langfuse.WithAgentStatus(err.Error()),
			langfuse.WithAgentLevel(langfuse.ObservationLevelError),
		)
		return nil, err
	}

	// don't update duration delta for reflector because it's already included in the performAgentChain
	if err := fp.updateMsgChainUsage(ctx, chainID, optAgentType, result.info, 0); err != nil {
		logger.WithError(err).Error("failed to update msg chain usage")
		opts = append(opts,
			langfuse.WithAgentStatus(err.Error()),
			langfuse.WithAgentLevel(langfuse.ObservationLevelError),
		)
		return nil, err
	}

	// preserve reasoning in reflector response using universal pattern
	reflectorMsg := llms.MessageContent{Role: llms.ChatMessageTypeAI}
	if result.content != "" || !result.thinking.IsEmpty() {
		reflectorMsg.Parts = append(reflectorMsg.Parts, llms.TextPartWithReasoning(result.content, result.thinking))
	}
	chain = append(chain, reflectorMsg)
	if len(result.funcCalls) == 0 {
		// Check if we are already in a reflector retry cycle to prevent infinite recursion.
		// This blocks recursive performReflector calls after caller reflector was invoked.
		if isReflectorRetry(ctx) {
			logger.Error("reflector recursion detected: cannot recursively call reflector after caller reflector")
			return nil, errors.New("reflector recursion detected: LLM returned no tool calls after reflector advice")
		}

		return fp.performReflector(ctx, optOriginType, chainID, taskID, subtaskID, chain, executor,
			humanMessage, result.content, executionContext, iteration+1)
	}

	opts = append(opts, langfuse.WithAgentStatus("success"))
	return result, nil
}

func (fp *flowProvider) performCallerReflector(
	ctx context.Context,
	optAgentType pconfig.ProviderOptionsType,
	chainID int64,
	taskID, subtaskID *int64,
	chain []llms.MessageContent,
	executor tools.ContextToolsExecutor,
	executionContext string,
	errs []error,
) (*callResult, error) {
	ctx, span := obs.Observer.NewSpan(ctx, obs.SpanKindInternal, "providers.flowProvider.performCallerReflector")
	defer span.End()

	logger := logrus.WithContext(ctx).WithFields(enrichLogrusFields(fp.flowID, taskID, subtaskID, logrus.Fields{
		"provider":     fp.Type(),
		"agent":        optAgentType,
		"msg_chain_id": chainID,
		"errors_count": len(errs),
	})).WithError(errors.Join(errs...))

	// Check if we are already in a reflector retry cycle to prevent infinite recursion.
	// This blocks repeated calls to performCallerReflector after reflector advice failed.
	if isReflectorRetry(ctx) {
		logger.Error("reflector recursion detected: caller reflector already invoked in this chain")
		return nil, errors.New("reflector recursion detected: cannot invoke caller reflector again after reflector advice failed")
	}

	// Mark context to prevent any further reflector recursion.
	// This flag will be checked in:
	// 1. performCallerReflector (here) - if reflector advice fails again
	// 2. performReflector - before recursive call when no tool calls returned
	ctx = markReflectorRetry(ctx)
	logger = logger.WithContext(ctx)

	logger.Warn("max retries reached, invoking caller reflector for guidance")

	reflectorContent := fmt.Sprintf(
		"I'm having trouble generating a proper tool call response. "+
			"I've attempted %d times but each attempt failed with errors:\n\n%s\n\n"+
			"I'm not sure how to proceed correctly. Should I try a different approach, "+
			"or should I use one of the barrier tools to report this issue?",
		len(errs), errors.Join(errs...).Error(),
	)

	reflectorResult, err := fp.performReflector(
		ctx, optAgentType, chainID, taskID, subtaskID, chain, executor,
		fp.getLastHumanMessage(chain), reflectorContent, executionContext, 1,
	)
	if err == nil {
		return reflectorResult, nil
	}

	return nil, fmt.Errorf("failed to perform caller reflector: %w", err)
}

func (fp *flowProvider) getLastHumanMessage(chain []llms.MessageContent) string {
	ast, err := cast.NewChainAST(chain, true)
	if err != nil {
		return ""
	}

	slices.Reverse(ast.Sections)
	for _, section := range ast.Sections {
		if section.Header.HumanMessage != nil {
			var hparts []string
			for _, part := range section.Header.HumanMessage.Parts {
				if text, ok := part.(llms.TextContent); ok {
					hparts = append(hparts, text.Text)
				}
			}
			return strings.Join(hparts, "\n")
		}
	}

	return ""
}
