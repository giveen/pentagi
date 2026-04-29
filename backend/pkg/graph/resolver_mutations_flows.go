package graph

import (
"context"
"errors"
"fmt"

"pentagi/pkg/controller"
"pentagi/pkg/database"
"pentagi/pkg/database/converter"
"pentagi/pkg/graph/model"
"pentagi/pkg/providers/provider"

"github.com/sirupsen/logrus"
)

// CreateFlow is the resolver for the createFlow field.
func (r *mutationResolver) CreateFlow(ctx context.Context, modelProvider string, input string) (*model.Flow, error) {
	uid, _, err := validatePermission(ctx, "flows.create")
	if err != nil {
		return nil, err
	}

	r.Logger.WithFields(logrus.Fields{
		"uid":      uid,
		"provider": modelProvider,
		"input":    input[:min(len(input), 1000)],
	}).Debug("create flow")

	if modelProvider == "" {
		return nil, fmt.Errorf("model provider is required")
	}

	if input == "" {
		return nil, fmt.Errorf("user input is required")
	}

	prvname := provider.ProviderName(modelProvider)
	prv, err := r.ProvidersCtrl.GetProvider(ctx, prvname, uid)
	if err != nil {
		return nil, err
	}
	prvtype := prv.Type()

	fw, err := r.Controller.CreateFlow(ctx, uid, input, prvname, prvtype, nil)
	if err != nil {
		return nil, err
	}

	flow, err := r.DB.GetFlow(ctx, fw.GetFlowID())
	if err != nil {
		return nil, err
	}

	var containers []database.Container
	if _, _, err = validatePermission(ctx, "containers.view"); err == nil {
		containers, err = r.DB.GetFlowContainers(ctx, fw.GetFlowID())
		if err != nil {
			return nil, err
		}
	}

	return converter.ConvertFlow(flow, containers), nil
}

// PutUserInput is the resolver for the putUserInput field.
func (r *mutationResolver) PutUserInput(ctx context.Context, flowID int64, input string, modelProvider *string) (model.ResultType, error) {
	uid, err := validatePermissionWithFlowID(ctx, "flows.edit", flowID, r.DB)
	if err != nil {
		return model.ResultTypeError, err
	}

	fields := logrus.Fields{
		"uid":  uid,
		"flow": flowID,
	}
	if modelProvider != nil {
		if *modelProvider == "" {
			fields["modelProvider"] = "empty"
		} else {
			fields["modelProvider"] = *modelProvider
		}
	} else {
		fields["modelProvider"] = "unknown"
	}

	r.Logger.WithFields(fields).Debug("put user input")

	fw, err := r.Controller.GetFlow(ctx, flowID)
	if err != nil {
		return model.ResultTypeError, err
	}

	var prv provider.Provider
	if modelProvider != nil && *modelProvider != "" {
		name := provider.ProviderName(*modelProvider)
		prv, err = r.ProvidersCtrl.GetProvider(ctx, name, uid)
		if err != nil {
			return model.ResultTypeError, fmt.Errorf("failed to get provider '%s': %w", *modelProvider, err)
		}
	}

	if err := fw.PutInput(ctx, input, prv); err != nil {
		return model.ResultTypeError, err
	}

	return model.ResultTypeSuccess, nil
}

// StopFlow is the resolver for the stopFlow field.
func (r *mutationResolver) StopFlow(ctx context.Context, flowID int64) (model.ResultType, error) {
	uid, err := validatePermissionWithFlowID(ctx, "flows.edit", flowID, r.DB)
	if err != nil {
		return model.ResultTypeError, err
	}

	r.Logger.WithFields(logrus.Fields{
		"uid":  uid,
		"flow": flowID,
	}).Debug("stop flow")

	if err := r.Controller.StopFlow(ctx, flowID); err != nil {
		return model.ResultTypeError, err
	}

	return model.ResultTypeSuccess, nil
}

// FinishFlow is the resolver for the finishFlow field.
func (r *mutationResolver) FinishFlow(ctx context.Context, flowID int64) (model.ResultType, error) {
	uid, err := validatePermissionWithFlowID(ctx, "flows.edit", flowID, r.DB)
	if err != nil {
		return model.ResultTypeError, err
	}

	r.Logger.WithFields(logrus.Fields{
		"uid":  uid,
		"flow": flowID,
	}).Debug("finish flow")

	err = r.Controller.FinishFlow(ctx, flowID)
	if err != nil {
		return model.ResultTypeError, err
	}

	return model.ResultTypeSuccess, nil
}

// DeleteFlow is the resolver for the deleteFlow field.
func (r *mutationResolver) DeleteFlow(ctx context.Context, flowID int64) (model.ResultType, error) {
	uid, err := validatePermissionWithFlowID(ctx, "flows.delete", flowID, r.DB)
	if err != nil {
		return model.ResultTypeError, err
	}

	r.Logger.WithFields(logrus.Fields{
		"uid":  uid,
		"flow": flowID,
	}).Debug("delete flow")

	if fw, err := r.Controller.GetFlow(ctx, flowID); err == nil {
		if err := fw.Finish(ctx); err != nil {
			return model.ResultTypeError, err
		}
	} else if !errors.Is(err, controller.ErrFlowNotFound) {
		return model.ResultTypeError, err
	}

	flow, err := r.DB.GetFlow(ctx, flowID)
	if err != nil {
		return model.ResultTypeError, err
	}

	containers, err := r.DB.GetFlowContainers(ctx, flow.ID)
	if err != nil {
		return model.ResultTypeError, err
	}

	if _, err := r.DB.DeleteFlow(ctx, flow.ID); err != nil {
		return model.ResultTypeError, err
	}

	publisher := r.Subscriptions.NewFlowPublisher(flow.UserID, flow.ID)
	publisher.FlowUpdated(ctx, flow, containers)
	publisher.FlowDeleted(ctx, flow, containers)

	return model.ResultTypeSuccess, nil
}

// RenameFlow is the resolver for the renameFlow field.
func (r *mutationResolver) RenameFlow(ctx context.Context, flowID int64, title string) (model.ResultType, error) {
	uid, err := validatePermissionWithFlowID(ctx, "flows.edit", flowID, r.DB)
	if err != nil {
		return model.ResultTypeError, err
	}

	r.Logger.WithFields(logrus.Fields{
		"uid":   uid,
		"flow":  flowID,
		"title": title,
	}).Debug("rename flow")

	err = r.Controller.RenameFlow(ctx, flowID, title)
	if errors.Is(err, controller.ErrFlowNotFound) {
		// if flow worker not found, update flow title in DB and notify about it
		flow, err := r.DB.UpdateFlowTitle(ctx, database.UpdateFlowTitleParams{
			ID:    flowID,
			Title: title,
		})
		if err != nil {
			return model.ResultTypeError, err
		}

		containers, err := r.DB.GetFlowContainers(ctx, flow.ID)
		if err != nil {
			return model.ResultTypeError, err
		}

		publisher := r.Subscriptions.NewFlowPublisher(flow.UserID, flow.ID)
		publisher.FlowUpdated(ctx, flow, containers)
	} else if err != nil {
		return model.ResultTypeError, err
	}

	return model.ResultTypeSuccess, nil
}

// CreateAssistant is the resolver for the createAssistant field.
func (r *mutationResolver) CreateAssistant(ctx context.Context, flowID int64, modelProvider string, input string, useAgents bool) (*model.FlowAssistant, error) {
	var (
		err error
		uid int64
	)

	if flowID == 0 {
		uid, _, err = validatePermission(ctx, "assistants.create")
		if err != nil {
			return nil, err
		}
		uid, _, err = validatePermission(ctx, "flows.create")
		if err != nil {
			return nil, err
		}
	} else {
		uid, err = validatePermissionWithFlowID(ctx, "assistants.create", flowID, r.DB)
		if err != nil {
			return nil, err
		}
	}

	r.Logger.WithFields(logrus.Fields{
		"uid":      uid,
		"flow":     flowID,
		"provider": modelProvider,
		"input":    input,
	}).Debug("create assistant")

	if modelProvider == "" {
		return nil, fmt.Errorf("model provider is required")
	}

	if input == "" {
		return nil, fmt.Errorf("user input is required")
	}

	prvname := provider.ProviderName(modelProvider)
	prv, err := r.ProvidersCtrl.GetProvider(ctx, prvname, uid)
	if err != nil {
		return nil, err
	}
	prvtype := prv.Type()

	aw, err := r.Controller.CreateAssistant(ctx, uid, flowID, input, useAgents, prvname, prvtype, nil)
	if err != nil {
		return nil, err
	}

	assistant, err := r.DB.GetAssistant(ctx, aw.GetAssistantID())
	if err != nil {
		return nil, err
	}

	flow, err := r.DB.GetFlow(ctx, assistant.FlowID)
	if err != nil {
		return nil, err
	}

	containers, err := r.DB.GetFlowContainers(ctx, assistant.FlowID)
	if err != nil {
		return nil, err
	}

	return converter.ConvertFlowAssistant(flow, containers, assistant), nil
}

// CallAssistant is the resolver for the callAssistant field.
func (r *mutationResolver) CallAssistant(ctx context.Context, flowID int64, assistantID int64, input string, useAgents bool) (model.ResultType, error) {
	uid, err := validatePermissionWithFlowID(ctx, "assistants.edit", flowID, r.DB)
	if err != nil {
		return model.ResultTypeError, err
	}

	r.Logger.WithFields(logrus.Fields{
		"uid":       uid,
		"flow":      flowID,
		"assistant": assistantID,
	}).Debug("call assistant")

	fw, err := r.Controller.GetFlow(ctx, flowID)
	if err != nil {
		return model.ResultTypeError, err
	}

	aw, err := fw.GetAssistant(ctx, assistantID)
	if err != nil {
		return model.ResultTypeError, err
	}

	if err := aw.PutInput(ctx, input, useAgents); err != nil {
		return model.ResultTypeError, err
	}

	return model.ResultTypeSuccess, nil
}

// StopAssistant is the resolver for the stopAssistant field.
func (r *mutationResolver) StopAssistant(ctx context.Context, flowID int64, assistantID int64) (*model.Assistant, error) {
	uid, err := validatePermissionWithFlowID(ctx, "assistants.edit", flowID, r.DB)
	if err != nil {
		return nil, err
	}

	r.Logger.WithFields(logrus.Fields{
		"uid":       uid,
		"flow":      flowID,
		"assistant": assistantID,
	}).Debug("stop assistant")

	fw, err := r.Controller.GetFlow(ctx, flowID)
	if err != nil {
		return nil, err
	}

	aw, err := fw.GetAssistant(ctx, assistantID)
	if err != nil {
		return nil, err
	}

	if err := aw.Stop(ctx); err != nil {
		return nil, err
	}

	assistant, err := r.DB.GetFlowAssistant(ctx, database.GetFlowAssistantParams{
		ID:     assistantID,
		FlowID: flowID,
	})
	if err != nil {
		return nil, err
	}

	r.Subscriptions.NewFlowPublisher(fw.GetUserID(), flowID).AssistantUpdated(ctx, assistant)

	return converter.ConvertAssistant(assistant), nil
}

// DeleteAssistant is the resolver for the deleteAssistant field.
func (r *mutationResolver) DeleteAssistant(ctx context.Context, flowID int64, assistantID int64) (model.ResultType, error) {
	uid, err := validatePermissionWithFlowID(ctx, "assistants.delete", flowID, r.DB)
	if err != nil {
		return model.ResultTypeError, err
	}

	r.Logger.WithFields(logrus.Fields{
		"uid":       uid,
		"flow":      flowID,
		"assistant": assistantID,
	}).Debug("delete assistant")

	fw, err := r.Controller.GetFlow(ctx, flowID)
	if err != nil {
		return model.ResultTypeError, err
	}

	assistant, err := r.DB.GetFlowAssistant(ctx, database.GetFlowAssistantParams{
		ID:     assistantID,
		FlowID: flowID,
	})
	if err != nil {
		return model.ResultTypeError, err
	}

	if err := fw.DeleteAssistant(ctx, assistantID); err != nil {
		return model.ResultTypeError, err
	}

	r.Subscriptions.NewFlowPublisher(fw.GetUserID(), flowID).AssistantDeleted(ctx, assistant)

	return model.ResultTypeSuccess, nil
}
