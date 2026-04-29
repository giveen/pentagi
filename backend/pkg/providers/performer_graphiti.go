package providers

import (
	"context"
	"fmt"
	"strings"
	"time"

	"pentagi/pkg/graphiti"
	obs "pentagi/pkg/observability"
	"pentagi/pkg/observability/langfuse"
	"pentagi/pkg/providers/pconfig"
	"pentagi/pkg/templates"
	"pentagi/pkg/tools"

	"github.com/sirupsen/logrus"
	"github.com/vxcontrol/langchaingo/llms"
)

// storeToGraphiti stores messages to Graphiti with timeout
func (fp *flowProvider) storeToGraphiti(
	ctx context.Context,
	observation langfuse.Observation,
	groupID string,
	messages []graphiti.Message,
) error {
	if fp.graphitiClient == nil || !fp.graphitiClient.IsEnabled() {
		return nil
	}

	storeCtx, cancel := context.WithTimeout(ctx, fp.graphitiClient.GetTimeout())
	defer cancel()

	err := fp.graphitiClient.AddMessages(storeCtx, graphiti.AddMessagesRequest{
		GroupID:  groupID,
		Messages: messages,
		Observation: &graphiti.Observation{
			ID:      observation.ID(),
			TraceID: observation.TraceID(),
			Time:    time.Now().UTC(),
		},
	})
	if err != nil {
		logrus.WithError(err).
			WithField("group_id", groupID).
			Warn("failed to store messages to graphiti")
	}

	return err
}

// storeAgentResponseToGraphiti stores agent response to Graphiti
func (fp *flowProvider) storeAgentResponseToGraphiti(
	ctx context.Context,
	groupID string,
	agentType pconfig.ProviderOptionsType,
	result *callResult,
	taskID, subtaskID *int64,
	chainID int64,
) {
	if fp.graphitiClient == nil || !fp.graphitiClient.IsEnabled() {
		return
	}

	if result.content == "" {
		return
	}

	tmpl, err := templates.ReadGraphitiTemplate("agent_response.tmpl")
	if err != nil {
		logrus.WithError(err).Warn("failed to read agent response template for graphiti")
		return
	}

	content, err := templates.RenderPrompt("agent_response", tmpl, map[string]any{
		"AgentType": string(agentType),
		"Response":  result.content,
		"TaskID":    taskID,
		"SubtaskID": subtaskID,
	})
	if err != nil {
		logrus.WithError(err).Warn("failed to render agent response template for graphiti")
		return
	}

	parts := []string{fmt.Sprintf("PentAGI %s agent execution in flow %d", agentType, fp.flowID)}
	if taskID != nil {
		parts = append(parts, fmt.Sprintf("task %d", *taskID))
	}
	if subtaskID != nil {
		parts = append(parts, fmt.Sprintf("subtask %d", *subtaskID))
	}
	sourceDescription := strings.Join(parts, ", ")

	messages := []graphiti.Message{
		{
			Content:           content,
			Author:            fmt.Sprintf("%s Agent", string(agentType)),
			Timestamp:         time.Now(),
			Name:              "agent_response",
			SourceDescription: sourceDescription,
		},
	}
	logrus.WithField("messages", messages).Debug("storing agent response to graphiti")

	ctx, observation := obs.Observer.NewObservation(ctx)
	storeEvaluator := observation.Evaluator(
		langfuse.WithEvaluatorName("store messages to graphiti"),
		langfuse.WithEvaluatorInput(messages),
		langfuse.WithEvaluatorMetadata(langfuse.Metadata{
			"group_id":     groupID,
			"agent_type":   agentType,
			"task_id":      taskID,
			"subtask_id":   subtaskID,
			"msg_chain_id": chainID,
		}),
	)

	ctx, observation = storeEvaluator.Observation(ctx)
	if err := fp.storeToGraphiti(ctx, observation, groupID, messages); err != nil {
		storeEvaluator.End(
			langfuse.WithEvaluatorStatus(err.Error()),
			langfuse.WithEvaluatorLevel(langfuse.ObservationLevelError),
		)
		return
	}

	storeEvaluator.End(
		langfuse.WithEvaluatorStatus("success"),
	)
}

// storeToolExecutionToGraphiti stores tool execution to Graphiti
func (fp *flowProvider) storeToolExecutionToGraphiti(
	ctx context.Context,
	groupID string,
	agentType pconfig.ProviderOptionsType,
	toolCall llms.ToolCall,
	response string,
	execErr error,
	executor tools.ContextToolsExecutor,
	taskID, subtaskID *int64,
	chainID int64,
) {
	if fp.graphitiClient == nil || !fp.graphitiClient.IsEnabled() {
		return
	}

	if toolCall.FunctionCall == nil {
		return
	}

	funcName := toolCall.FunctionCall.Name
	funcArgs := toolCall.FunctionCall.Arguments

	registryDefs := tools.GetRegistryDefinitions()
	toolDef, ok := registryDefs[funcName]
	description := ""
	if ok {
		description = toolDef.Description
	}

	isBarrier := executor.IsBarrierFunction(funcName)

	status := "success"
	if execErr != nil {
		status = "failure"
		response = fmt.Sprintf("Error: %s", execErr.Error())
	}

	toolExecTmpl, err := templates.ReadGraphitiTemplate("tool_execution.tmpl")
	if err != nil {
		logrus.WithError(err).Warn("failed to read tool execution template for graphiti")
		return
	}

	toolExecContent, err := templates.RenderPrompt("tool_execution", toolExecTmpl, map[string]any{
		"ToolName":    funcName,
		"Description": description,
		"IsBarrier":   isBarrier,
		"Arguments":   funcArgs,
		"AgentType":   string(agentType),
		"Status":      status,
		"Result":      response,
		"TaskID":      taskID,
		"SubtaskID":   subtaskID,
	})
	if err != nil {
		logrus.WithError(err).Warn("failed to render tool execution template for graphiti")
		return
	}

	parts := []string{fmt.Sprintf("PentAGI tool execution in flow %d", fp.flowID)}
	if taskID != nil {
		parts = append(parts, fmt.Sprintf("task %d", *taskID))
	}
	if subtaskID != nil {
		parts = append(parts, fmt.Sprintf("subtask %d", *subtaskID))
	}
	sourceDescription := strings.Join(parts, ", ")

	messages := []graphiti.Message{
		{
			Content:           toolExecContent,
			Author:            fmt.Sprintf("%s Agent", string(agentType)),
			Timestamp:         time.Now(),
			Name:              fmt.Sprintf("tool_execution_%s", funcName),
			SourceDescription: sourceDescription,
		},
	}

	ctx, observation := obs.Observer.NewObservation(ctx)
	storeEvaluator := observation.Evaluator(
		langfuse.WithEvaluatorName("store tool execution to graphiti"),
		langfuse.WithEvaluatorInput(messages),
		langfuse.WithEvaluatorMetadata(langfuse.Metadata{
			"group_id":     groupID,
			"agent_type":   agentType,
			"tool_name":    funcName,
			"tool_args":    funcArgs,
			"task_id":      taskID,
			"subtask_id":   subtaskID,
			"msg_chain_id": chainID,
		}),
	)

	ctx, observation = storeEvaluator.Observation(ctx)
	if err := fp.storeToGraphiti(ctx, observation, groupID, messages); err != nil {
		storeEvaluator.End(
			langfuse.WithEvaluatorStatus(err.Error()),
			langfuse.WithEvaluatorLevel(langfuse.ObservationLevelError),
		)
		return
	}

	storeEvaluator.End(
		langfuse.WithEvaluatorStatus("success"),
	)
}
