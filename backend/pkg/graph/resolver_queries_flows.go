package graph

import (
"context"

"pentagi/pkg/database"
"pentagi/pkg/database/converter"
"pentagi/pkg/graph/model"
"pentagi/pkg/providers/autotuner"
"pentagi/pkg/providers/pconfig"

"github.com/sirupsen/logrus"
)

// Providers is the resolver for the providers field.
func (r *queryResolver) Providers(ctx context.Context) ([]*model.Provider, error) {
	uid, _, err := validatePermission(ctx, "providers.view")
	if err != nil {
		return nil, err
	}

	r.Logger.WithFields(logrus.Fields{
		"uid": uid,
	}).Debug("get providers")

	providers, err := r.ProvidersCtrl.GetProviders(ctx, uid)
	if err != nil {
		return nil, err
	}

	providersList := make([]*model.Provider, len(providers))
	for i, prvname := range providers.ListNames() {
		providersList[i] = &model.Provider{
			Name: string(prvname),
			Type: model.ProviderType(providers[prvname].Type()),
		}
	}

	return providersList, nil
}

// AutoTuneParams is the resolver for the autoTuneParams field.
func (r *queryResolver) AutoTuneParams(ctx context.Context, agentType model.AgentConfigType, modelName string) (*model.AutoTuneParams, error) {
	if _, _, err := validatePermission(ctx, "settings.view"); err != nil {
		return nil, err
	}

	role := pconfig.ProviderOptionsType(agentType)
	params, profile, family := autotuner.Get().GetParamsWithMeta(role, 1, modelName)

	return &model.AutoTuneParams{
		Temperature:       params.Temperature,
		TopP:              params.TopP,
		TopK:              int(params.TopK),
		MinP:              params.MinP,
		FrequencyPenalty:  params.FrequencyPenalty,
		PresencePenalty:   params.PresencePenalty,
		RepetitionPenalty: params.RepetitionPenalty,
		Profile:           profile,
		Family:            family,
	}, nil
}

// Assistants is the resolver for the assistants field.
func (r *queryResolver) Assistants(ctx context.Context, flowID int64) ([]*model.Assistant, error) {
	uid, err := validatePermissionWithFlowID(ctx, "assistants.view", flowID, r.DB)
	if err != nil {
		return nil, err
	}

	r.Logger.WithFields(logrus.Fields{
		"uid":  uid,
		"flow": flowID,
	}).Debug("get assistants")

	assistants, err := r.DB.GetFlowAssistants(ctx, flowID)
	if err != nil {
		return nil, err
	}

	return converter.ConvertAssistants(assistants), nil
}

// Flows is the resolver for the flows field.
func (r *queryResolver) Flows(ctx context.Context) ([]*model.Flow, error) {
	uid, admin, err := validatePermission(ctx, "flows.view")
	if err != nil {
		return nil, err
	}

	r.Logger.WithFields(logrus.Fields{
		"uid": uid,
	}).Debug("get flows")

	var (
		flows      []database.Flow
		containers []database.Container
	)

	if admin {
		flows, err = r.DB.GetFlows(ctx)
	} else {
		flows, err = r.DB.GetUserFlows(ctx, uid)
	}
	if err != nil {
		return nil, err
	}

	if _, admin, err = validatePermission(ctx, "containers.view"); err == nil {
		if admin {
			containers, err = r.DB.GetContainers(ctx)
		} else {
			containers, err = r.DB.GetUserContainers(ctx, uid)
		}
		if err != nil {
			return nil, err
		}
	}

	return converter.ConvertFlows(flows, containers), nil
}

// Flow is the resolver for the flow field.
func (r *queryResolver) Flow(ctx context.Context, flowID int64) (*model.Flow, error) {
	uid, err := validatePermissionWithFlowID(ctx, "flows.view", flowID, r.DB)
	if err != nil {
		return nil, err
	}

	r.Logger.WithFields(logrus.Fields{
		"uid":  uid,
		"flow": flowID,
	}).Debug("get flow")

	var (
		flow       database.Flow
		containers []database.Container
	)

	flow, err = r.DB.GetFlow(ctx, flowID)
	if err != nil {
		return nil, err
	}

	if _, _, err = validatePermission(ctx, "containers.view"); err == nil {
		containers, err = r.DB.GetFlowContainers(ctx, flowID)
		if err != nil {
			return nil, err
		}
	}

	return converter.ConvertFlow(flow, containers), nil
}

// Tasks is the resolver for the tasks field.
func (r *queryResolver) Tasks(ctx context.Context, flowID int64) ([]*model.Task, error) {
	uid, err := validatePermissionWithFlowID(ctx, "tasks.view", flowID, r.DB)
	if err != nil {
		return nil, err
	}

	r.Logger.WithFields(logrus.Fields{
		"uid":  uid,
		"flow": flowID,
	}).Debug("get tasks")

	tasks, err := r.DB.GetFlowTasks(ctx, flowID)
	if err != nil {
		return nil, err
	}

	var subtasks []database.Subtask
	if _, _, err = validatePermission(ctx, "subtasks.view"); err == nil {
		subtasks, err = r.DB.GetFlowSubtasks(ctx, flowID)
		if err != nil {
			return nil, err
		}
	}

	return converter.ConvertTasks(tasks, subtasks), nil
}

// Screenshots is the resolver for the screenshots field.
func (r *queryResolver) Screenshots(ctx context.Context, flowID int64) ([]*model.Screenshot, error) {
	uid, err := validatePermissionWithFlowID(ctx, "screenshots.view", flowID, r.DB)
	if err != nil {
		return nil, err
	}

	r.Logger.WithFields(logrus.Fields{
		"uid":  uid,
		"flow": flowID,
	}).Debug("get screenshots")

	screenshots, err := r.DB.GetFlowScreenshots(ctx, flowID)
	if err != nil {
		return nil, err
	}

	return converter.ConvertScreenshots(screenshots), nil
}

// TerminalLogs is the resolver for the terminalLogs field.
func (r *queryResolver) TerminalLogs(ctx context.Context, flowID int64) ([]*model.TerminalLog, error) {
	uid, err := validatePermissionWithFlowID(ctx, "termlogs.view", flowID, r.DB)
	if err != nil {
		return nil, err
	}

	r.Logger.WithFields(logrus.Fields{
		"uid":  uid,
		"flow": flowID,
	}).Debug("get term logs")

	logs, err := r.DB.GetFlowTermLogs(ctx, flowID)
	if err != nil {
		return nil, err
	}

	return converter.ConvertTerminalLogs(logs), nil
}

// MessageLogs is the resolver for the messageLogs field.
func (r *queryResolver) MessageLogs(ctx context.Context, flowID int64) ([]*model.MessageLog, error) {
	uid, err := validatePermissionWithFlowID(ctx, "msglogs.view", flowID, r.DB)
	if err != nil {
		return nil, err
	}

	r.Logger.WithFields(logrus.Fields{
		"uid":  uid,
		"flow": flowID,
	}).Debug("get msg logs")

	logs, err := r.DB.GetFlowMsgLogs(ctx, flowID)
	if err != nil {
		return nil, err
	}

	return converter.ConvertMessageLogs(logs), nil
}

// AgentLogs is the resolver for the agentLogs field.
func (r *queryResolver) AgentLogs(ctx context.Context, flowID int64) ([]*model.AgentLog, error) {
	uid, err := validatePermissionWithFlowID(ctx, "agentlogs.view", flowID, r.DB)
	if err != nil {
		return nil, err
	}

	r.Logger.WithFields(logrus.Fields{
		"uid":  uid,
		"flow": flowID,
	}).Debug("get agent logs")

	logs, err := r.DB.GetFlowAgentLogs(ctx, flowID)
	if err != nil {
		return nil, err
	}

	return converter.ConvertAgentLogs(logs), nil
}

// SearchLogs is the resolver for the searchLogs field.
func (r *queryResolver) SearchLogs(ctx context.Context, flowID int64) ([]*model.SearchLog, error) {
	uid, err := validatePermissionWithFlowID(ctx, "searchlogs.view", flowID, r.DB)
	if err != nil {
		return nil, err
	}

	r.Logger.WithFields(logrus.Fields{
		"uid":  uid,
		"flow": flowID,
	}).Debug("get search logs")

	logs, err := r.DB.GetFlowSearchLogs(ctx, flowID)
	if err != nil {
		return nil, err
	}

	return converter.ConvertSearchLogs(logs), nil
}

// VectorStoreLogs is the resolver for the vectorStoreLogs field.
func (r *queryResolver) VectorStoreLogs(ctx context.Context, flowID int64) ([]*model.VectorStoreLog, error) {
	uid, err := validatePermissionWithFlowID(ctx, "vecstorelogs.view", flowID, r.DB)
	if err != nil {
		return nil, err
	}

	r.Logger.WithFields(logrus.Fields{
		"uid":  uid,
		"flow": flowID,
	}).Debug("get vector store logs")

	logs, err := r.DB.GetFlowVectorStoreLogs(ctx, flowID)
	if err != nil {
		return nil, err
	}

	return converter.ConvertVectorStoreLogs(logs), nil
}

// AssistantLogs is the resolver for the assistantLogs field.
func (r *queryResolver) AssistantLogs(ctx context.Context, flowID int64, assistantID int64) ([]*model.AssistantLog, error) {
	uid, err := validatePermissionWithFlowID(ctx, "assistantlogs.view", flowID, r.DB)
	if err != nil {
		return nil, err
	}

	r.Logger.WithFields(logrus.Fields{
		"uid":       uid,
		"flow":      flowID,
		"assistant": assistantID,
	}).Debug("get assistant logs")

	logs, err := r.DB.GetFlowAssistantLogs(ctx, database.GetFlowAssistantLogsParams{
		FlowID:      flowID,
		AssistantID: assistantID,
	})
	if err != nil {
		return nil, err
	}

	return converter.ConvertAssistantLogs(logs), nil
}
