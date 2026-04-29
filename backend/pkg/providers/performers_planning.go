package providers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"pentagi/pkg/database"
	obs "pentagi/pkg/observability"
	"pentagi/pkg/providers/pconfig"
	"pentagi/pkg/templates"
	"pentagi/pkg/tools"

	"github.com/sirupsen/logrus"
	"github.com/vxcontrol/langchaingo/llms"
	"github.com/vxcontrol/langchaingo/llms/reasoning"
)

// performPlanner invokes adviser to create an execution plan for agent tasks
func (fp *flowProvider) performPlanner(
	ctx context.Context,
	taskID, subtaskID *int64,
	opt pconfig.ProviderOptionsType,
	executor tools.ContextToolsExecutor,
	userTmpl, question string,
) (string, error) {
	ctx, span := obs.Observer.NewSpan(ctx, obs.SpanKindInternal, "providers.flowProvider.performPlanner")
	defer span.End()

	toolCallID := templates.GenerateFromPattern(fp.tcIDTemplate, tools.AdviceToolName)
	logger := logrus.WithContext(ctx).WithFields(logrus.Fields{
		"task_id":      taskID,
		"subtask_id":   subtaskID,
		"agent_type":   string(opt),
		"tool_call_id": toolCallID,
	})

	logger.Debug("requesting task plan from adviser (planner)")

	// 1. Format Question for task planning
	strategicState := ""
	if taskID != nil {
		tasksInfo, infoErr := fp.getTasksInfo(ctx, *taskID)
		if infoErr != nil {
			logger.WithError(infoErr).Warn("failed to load tasks info for strategic state, continuing without it")
		} else {
			strategicState = fp.getTaskStrategicState(tasksInfo.Task, tasksInfo.Tasks, tasksInfo.Subtasks)
		}
	}

	planQuestionData := map[string]any{
		"AgentType":      string(opt),
		"TaskQuestion":   question,
		"StrategicState": strategicState,
	}

	planQuestion, err := fp.prompter.RenderTemplate(templates.PromptTypeQuestionTaskPlanner, planQuestionData)
	if err != nil {
		return "", fmt.Errorf("failed to render task planner question: %w", err)
	}

	// 2. Call adviser handler with custom observation name "planner"
	askAdvice := tools.AskAdvice{
		Question: planQuestion,
	}

	askAdviceJSON, err := json.Marshal(askAdvice)
	if err != nil {
		return "", fmt.Errorf("failed to marshal ask advice: %w", err)
	}

	logger.Debug("executing adviser handler for task planning")
	plan, err := executor.Execute(ctx, 0, toolCallID, tools.AdviceToolName, "planner", "", askAdviceJSON)
	if err != nil {
		return "", fmt.Errorf("failed to execute adviser handler: %w", err)
	}

	logger.WithField("plan_length", len(plan)).Debug("task plan created successfully")

	// Wrap original request with execution plan using template
	taskAssignment, err := fp.prompter.RenderTemplate(templates.PromptTypeTaskAssignmentWrapper, map[string]any{
		"OriginalRequest": userTmpl,
		"ExecutionPlan":   plan,
	})
	if err != nil {
		return "", fmt.Errorf("failed to render task assignment wrapper: %w", err)
	}

	return taskAssignment, nil
}

// performMentor invokes adviser to monitor agent execution progress
func (fp *flowProvider) performMentor(
	ctx context.Context,
	opt pconfig.ProviderOptionsType,
	chainID int64,
	taskID, subtaskID *int64,
	chain []llms.MessageContent,
	executor tools.ContextToolsExecutor,
	lastToolCall llms.ToolCall,
	lastToolResult string,
) (string, error) {
	ctx, span := obs.Observer.NewSpan(ctx, obs.SpanKindInternal, "providers.flowProvider.performMentor")
	defer span.End()

	if lastToolCall.FunctionCall == nil {
		return "", fmt.Errorf("last tool call function call is nil")
	}

	toolCallID := templates.GenerateFromPattern(fp.tcIDTemplate, tools.AdviceToolName)
	logger := logrus.WithContext(ctx).WithFields(logrus.Fields{
		"chain_id":       chainID,
		"task_id":        taskID,
		"subtask_id":     subtaskID,
		"last_tool_name": lastToolCall.FunctionCall.Name,
		"agent_type":     string(opt),
		"tool_call_id":   toolCallID,
	})

	logger.Debug("invoking execution adviser for progress monitoring (mentor)")

	// 1. Collect recent messages from chain
	recentMessages := getRecentMessages(chain)

	// 2. Extract all executed tool calls from chain
	executedToolCalls := extractToolCallsFromChain(chain)

	// 3. Get subtask description
	subtaskDesc := ""
	if subtaskID != nil {
		if subtask, err := fp.db.GetSubtask(ctx, *subtaskID); err == nil {
			subtaskDesc = subtask.Description
		}
	}

	// 4. Extract original agent prompt from chain
	agentPrompt := extractAgentPromptFromChain(chain)

	// 5. Format Question through new template
	questionData := map[string]any{
		"SubtaskDescription": subtaskDesc,
		"AgentType":          string(opt),
		"AgentPrompt":        agentPrompt,
		"RecentMessages":     recentMessages,
		"ExecutedToolCalls":  executedToolCalls,
		"LastToolName":       lastToolCall.FunctionCall.Name,
		"LastToolArgs":       formatToolCallArguments(lastToolCall.FunctionCall.Arguments),
		"LastToolResult":     cutString(lastToolResult, 4096),
	}

	question, err := fp.prompter.RenderTemplate(templates.PromptTypeQuestionExecutionMonitor, questionData)
	if err != nil {
		return "", fmt.Errorf("failed to render execution monitor question: %w", err)
	}

	// 6. Call adviser handler with custom observation name "mentor"
	askAdvice := tools.AskAdvice{
		Question: question,
	}

	askAdviceJSON, err := json.Marshal(askAdvice)
	if err != nil {
		return "", fmt.Errorf("failed to marshal ask advice: %w", err)
	}

	logger.Debug("executing adviser handler for execution monitoring")
	result, err := executor.Execute(ctx, 0, toolCallID, tools.AdviceToolName, "mentor", "", askAdviceJSON)
	if err != nil {
		return "", fmt.Errorf("failed to execute adviser handler: %w", err)
	}

	logger.WithField("result_length", len(result)).Debug("execution mentor completed successfully")
	return result, nil
}

func (fp *flowProvider) performSimpleChain(
	ctx context.Context,
	taskID, subtaskID *int64,
	opt pconfig.ProviderOptionsType,
	msgChainType database.MsgchainType,
	systemTmpl, userTmpl string,
) (string, error) {
	var (
		resp *llms.ContentResponse
		err  error
	)

	startTime := time.Now()

	chain := []llms.MessageContent{
		llms.TextParts(llms.ChatMessageTypeSystem, systemTmpl),
		llms.TextParts(llms.ChatMessageTypeHuman, userTmpl),
	}

	for idx := 0; idx <= maxRetriesToCallSimpleChain; idx++ {
		if idx == maxRetriesToCallSimpleChain {
			return "", fmt.Errorf("failed to call simple chain: %w", err)
		}

		resp, err = fp.CallEx(ctx, opt, chain, nil)
		if err == nil {
			break
		} else {
			if errors.Is(err, context.Canceled) {
				return "", err
			}

			select {
			case <-ctx.Done():
				return "", ctx.Err()
			case <-time.After(time.Second * 5):
			default:
			}
		}
	}

	if len(resp.Choices) == 0 {
		return "", fmt.Errorf("no choices in response")
	}

	var parts []string
	var usage pconfig.CallUsage
	var reasoning *reasoning.ContentReasoning
	for _, choice := range resp.Choices {
		parts = append(parts, choice.Content)
		usage.Merge(fp.GetUsage(choice.GenerationInfo))
		// Preserve reasoning from first choice for simple chains (safe for all providers)
		if reasoning == nil && !choice.Reasoning.IsEmpty() {
			reasoning = choice.Reasoning
		}
	}

	// Update cost based on price info
	usage.UpdateCost(fp.GetPriceInfo(opt))

	// Universal pattern for simple chains - preserve reasoning if present
	msg := llms.MessageContent{Role: llms.ChatMessageTypeAI}
	content := strings.Join(parts, "\n")
	if content != "" || reasoning != nil {
		msg.Parts = append(msg.Parts, llms.TextPartWithReasoning(content, reasoning))
	}
	chain = append(chain, msg)

	chainBlob, err := json.Marshal(chain)
	if err != nil {
		return "", fmt.Errorf("failed to marshal summarizer msg chain: %w", err)
	}

	_, err = fp.db.CreateMsgChain(ctx, database.CreateMsgChainParams{
		Type:            msgChainType,
		Model:           fp.Model(opt),
		ModelProvider:   string(fp.Type()),
		UsageIn:         usage.Input,
		UsageOut:        usage.Output,
		UsageCacheIn:    usage.CacheRead,
		UsageCacheOut:   usage.CacheWrite,
		UsageCostIn:     usage.CostInput,
		UsageCostOut:    usage.CostOutput,
		DurationSeconds: time.Since(startTime).Seconds(),
		Chain:           chainBlob,
		FlowID:          fp.flowID,
		TaskID:          database.Int64ToNullInt64(taskID),
		SubtaskID:       database.Int64ToNullInt64(subtaskID),
	})

	return strings.Join(parts, "\n\n"), nil
}
