package graph

import (
"context"
"fmt"

"pentagi/pkg/database"
"pentagi/pkg/database/converter"
"pentagi/pkg/graph/model"

"github.com/sirupsen/logrus"
)

// UsageStatsTotal is the resolver for the usageStatsTotal field.
func (r *queryResolver) UsageStatsTotal(ctx context.Context) (*model.UsageStats, error) {
	uid, _, err := validatePermission(ctx, "usage.view")
	if err != nil {
		return nil, err
	}

	r.Logger.WithFields(logrus.Fields{
		"uid": uid,
	}).Debug("get total usage stats")

	stats, err := r.DB.GetUserTotalUsageStats(ctx, uid)
	if err != nil {
		return nil, err
	}

	return converter.ConvertUsageStats(stats), nil
}

// UsageStatsByPeriod is the resolver for the usageStatsByPeriod field.
func (r *queryResolver) UsageStatsByPeriod(ctx context.Context, period model.UsageStatsPeriod) ([]*model.DailyUsageStats, error) {
	uid, _, err := validatePermission(ctx, "usage.view")
	if err != nil {
		return nil, err
	}

	r.Logger.WithFields(logrus.Fields{
		"uid":    uid,
		"period": period,
	}).Debug("get usage stats by period")

	switch period {
	case model.UsageStatsPeriodWeek:
		stats, err := r.DB.GetUsageStatsByDayLastWeek(ctx, uid)
		if err != nil {
			return nil, err
		}
		return converter.ConvertDailyUsageStats(stats), nil
	case model.UsageStatsPeriodMonth:
		stats, err := r.DB.GetUsageStatsByDayLastMonth(ctx, uid)
		if err != nil {
			return nil, err
		}
		return converter.ConvertDailyUsageStatsMonth(stats), nil
	case model.UsageStatsPeriodQuarter:
		stats, err := r.DB.GetUsageStatsByDayLast3Months(ctx, uid)
		if err != nil {
			return nil, err
		}
		return converter.ConvertDailyUsageStatsQuarter(stats), nil
	default:
		return nil, fmt.Errorf("invalid period: %s", period)
	}
}

// UsageStatsByProvider is the resolver for the usageStatsByProvider field.
func (r *queryResolver) UsageStatsByProvider(ctx context.Context) ([]*model.ProviderUsageStats, error) {
	uid, _, err := validatePermission(ctx, "usage.view")
	if err != nil {
		return nil, err
	}

	r.Logger.WithFields(logrus.Fields{
		"uid": uid,
	}).Debug("get usage stats by provider")

	stats, err := r.DB.GetUsageStatsByProvider(ctx, uid)
	if err != nil {
		return nil, err
	}

	return converter.ConvertProviderUsageStats(stats), nil
}

// UsageStatsByModel is the resolver for the usageStatsByModel field.
func (r *queryResolver) UsageStatsByModel(ctx context.Context) ([]*model.ModelUsageStats, error) {
	uid, _, err := validatePermission(ctx, "usage.view")
	if err != nil {
		return nil, err
	}

	r.Logger.WithFields(logrus.Fields{
		"uid": uid,
	}).Debug("get usage stats by model")

	stats, err := r.DB.GetUsageStatsByModel(ctx, uid)
	if err != nil {
		return nil, err
	}

	return converter.ConvertModelUsageStats(stats), nil
}

// UsageStatsByAgentType is the resolver for the usageStatsByAgentType field.
func (r *queryResolver) UsageStatsByAgentType(ctx context.Context) ([]*model.AgentTypeUsageStats, error) {
	uid, _, err := validatePermission(ctx, "usage.view")
	if err != nil {
		return nil, err
	}

	r.Logger.WithFields(logrus.Fields{
		"uid": uid,
	}).Debug("get usage stats by agent type")

	stats, err := r.DB.GetUsageStatsByType(ctx, uid)
	if err != nil {
		return nil, err
	}

	return converter.ConvertAgentTypeUsageStats(stats), nil
}

// UsageStatsByFlow is the resolver for the usageStatsByFlow field.
func (r *queryResolver) UsageStatsByFlow(ctx context.Context, flowID int64) (*model.UsageStats, error) {
	uid, err := validatePermissionWithFlowID(ctx, "usage.view", flowID, r.DB)
	if err != nil {
		return nil, err
	}

	r.Logger.WithFields(logrus.Fields{
		"uid":  uid,
		"flow": flowID,
	}).Debug("get usage stats by flow")

	stats, err := r.DB.GetFlowUsageStats(ctx, flowID)
	if err != nil {
		return nil, err
	}

	return converter.ConvertUsageStats(stats), nil
}

// UsageStatsByAgentTypeForFlow is the resolver for the usageStatsByAgentTypeForFlow field.
func (r *queryResolver) UsageStatsByAgentTypeForFlow(ctx context.Context, flowID int64) ([]*model.AgentTypeUsageStats, error) {
	uid, err := validatePermissionWithFlowID(ctx, "usage.view", flowID, r.DB)
	if err != nil {
		return nil, err
	}

	r.Logger.WithFields(logrus.Fields{
		"uid":  uid,
		"flow": flowID,
	}).Debug("get usage stats by agent type for flow")

	stats, err := r.DB.GetUsageStatsByTypeForFlow(ctx, flowID)
	if err != nil {
		return nil, err
	}

	return converter.ConvertAgentTypeUsageStatsForFlow(stats), nil
}

// ToolcallsStatsTotal is the resolver for the toolcallsStatsTotal field.
func (r *queryResolver) ToolcallsStatsTotal(ctx context.Context) (*model.ToolcallsStats, error) {
	uid, _, err := validatePermission(ctx, "usage.view")
	if err != nil {
		return nil, err
	}

	r.Logger.WithFields(logrus.Fields{
		"uid": uid,
	}).Debug("get total toolcalls stats")

	stats, err := r.DB.GetUserTotalToolcallsStats(ctx, uid)
	if err != nil {
		return nil, err
	}

	return converter.ConvertToolcallsStats(stats), nil
}

// ToolcallsStatsByPeriod is the resolver for the toolcallsStatsByPeriod field.
func (r *queryResolver) ToolcallsStatsByPeriod(ctx context.Context, period model.UsageStatsPeriod) ([]*model.DailyToolcallsStats, error) {
	uid, _, err := validatePermission(ctx, "usage.view")
	if err != nil {
		return nil, err
	}

	r.Logger.WithFields(logrus.Fields{
		"uid":    uid,
		"period": period,
	}).Debug("get toolcalls stats by period")

	switch period {
	case model.UsageStatsPeriodWeek:
		stats, err := r.DB.GetToolcallsStatsByDayLastWeek(ctx, uid)
		if err != nil {
			return nil, err
		}
		return converter.ConvertDailyToolcallsStatsWeek(stats), nil
	case model.UsageStatsPeriodMonth:
		stats, err := r.DB.GetToolcallsStatsByDayLastMonth(ctx, uid)
		if err != nil {
			return nil, err
		}
		return converter.ConvertDailyToolcallsStatsMonth(stats), nil
	case model.UsageStatsPeriodQuarter:
		stats, err := r.DB.GetToolcallsStatsByDayLast3Months(ctx, uid)
		if err != nil {
			return nil, err
		}
		return converter.ConvertDailyToolcallsStatsQuarter(stats), nil
	default:
		return nil, fmt.Errorf("unsupported period: %s", period)
	}
}

// ToolcallsStatsByFunction is the resolver for the toolcallsStatsByFunction field.
func (r *queryResolver) ToolcallsStatsByFunction(ctx context.Context) ([]*model.FunctionToolcallsStats, error) {
	uid, _, err := validatePermission(ctx, "usage.view")
	if err != nil {
		return nil, err
	}

	r.Logger.WithFields(logrus.Fields{
		"uid": uid,
	}).Debug("get toolcalls stats by function")

	stats, err := r.DB.GetToolcallsStatsByFunction(ctx, uid)
	if err != nil {
		return nil, err
	}

	return converter.ConvertFunctionToolcallsStats(stats), nil
}

// ToolcallsStatsByFlow is the resolver for the toolcallsStatsByFlow field.
func (r *queryResolver) ToolcallsStatsByFlow(ctx context.Context, flowID int64) (*model.ToolcallsStats, error) {
	uid, err := validatePermissionWithFlowID(ctx, "usage.view", flowID, r.DB)
	if err != nil {
		return nil, err
	}

	r.Logger.WithFields(logrus.Fields{
		"uid":  uid,
		"flow": flowID,
	}).Debug("get toolcalls stats by flow")

	stats, err := r.DB.GetFlowToolcallsStats(ctx, flowID)
	if err != nil {
		return nil, err
	}

	return converter.ConvertToolcallsStats(stats), nil
}

// ToolcallsStatsByFunctionForFlow is the resolver for the toolcallsStatsByFunctionForFlow field.
func (r *queryResolver) ToolcallsStatsByFunctionForFlow(ctx context.Context, flowID int64) ([]*model.FunctionToolcallsStats, error) {
	uid, err := validatePermissionWithFlowID(ctx, "usage.view", flowID, r.DB)
	if err != nil {
		return nil, err
	}

	r.Logger.WithFields(logrus.Fields{
		"uid":  uid,
		"flow": flowID,
	}).Debug("get toolcalls stats by function for flow")

	stats, err := r.DB.GetToolcallsStatsByFunctionForFlow(ctx, flowID)
	if err != nil {
		return nil, err
	}

	return converter.ConvertFunctionToolcallsStatsForFlow(stats), nil
}

// FlowsStatsTotal is the resolver for the flowsStatsTotal field.
func (r *queryResolver) FlowsStatsTotal(ctx context.Context) (*model.FlowsStats, error) {
	uid, _, err := validatePermission(ctx, "usage.view")
	if err != nil {
		return nil, err
	}

	r.Logger.WithFields(logrus.Fields{
		"uid": uid,
	}).Debug("get total flows stats")

	stats, err := r.DB.GetUserTotalFlowsStats(ctx, uid)
	if err != nil {
		return nil, err
	}

	return converter.ConvertFlowsStats(stats), nil
}

// FlowsStatsByPeriod is the resolver for the flowsStatsByPeriod field.
func (r *queryResolver) FlowsStatsByPeriod(ctx context.Context, period model.UsageStatsPeriod) ([]*model.DailyFlowsStats, error) {
	uid, _, err := validatePermission(ctx, "usage.view")
	if err != nil {
		return nil, err
	}

	r.Logger.WithFields(logrus.Fields{
		"uid":    uid,
		"period": period,
	}).Debug("get flows stats by period")

	switch period {
	case model.UsageStatsPeriodWeek:
		stats, err := r.DB.GetFlowsStatsByDayLastWeek(ctx, uid)
		if err != nil {
			return nil, err
		}
		return converter.ConvertDailyFlowsStatsWeek(stats), nil
	case model.UsageStatsPeriodMonth:
		stats, err := r.DB.GetFlowsStatsByDayLastMonth(ctx, uid)
		if err != nil {
			return nil, err
		}
		return converter.ConvertDailyFlowsStatsMonth(stats), nil
	case model.UsageStatsPeriodQuarter:
		stats, err := r.DB.GetFlowsStatsByDayLast3Months(ctx, uid)
		if err != nil {
			return nil, err
		}
		return converter.ConvertDailyFlowsStatsQuarter(stats), nil
	default:
		return nil, fmt.Errorf("unsupported period: %s", period)
	}
}

// FlowStatsByFlow is the resolver for the flowStatsByFlow field.
func (r *queryResolver) FlowStatsByFlow(ctx context.Context, flowID int64) (*model.FlowStats, error) {
	uid, err := validatePermissionWithFlowID(ctx, "usage.view", flowID, r.DB)
	if err != nil {
		return nil, err
	}

	r.Logger.WithFields(logrus.Fields{
		"uid":  uid,
		"flow": flowID,
	}).Debug("get flow stats by flow")

	stats, err := r.DB.GetFlowStats(ctx, flowID)
	if err != nil {
		return nil, err
	}

	return converter.ConvertFlowStats(stats), nil
}

// FlowsExecutionStatsByPeriod is the resolver for the flowsExecutionStatsByPeriod field.
func (r *queryResolver) FlowsExecutionStatsByPeriod(ctx context.Context, period model.UsageStatsPeriod) ([]*model.FlowExecutionStats, error) {
	uid, _, err := validatePermission(ctx, "usage.view")
	if err != nil {
		return nil, err
	}

	r.Logger.WithFields(logrus.Fields{
		"uid":    uid,
		"period": period,
	}).Debug("get flows execution stats by period")

	// Step 1: Get flows for the period
	type flowInfo struct {
		ID    int64
		Title string
	}
	var flows []flowInfo

	switch period {
	case model.UsageStatsPeriodWeek:
		rows, err := r.DB.GetFlowsForPeriodLastWeek(ctx, uid)
		if err != nil {
			return nil, err
		}
		for _, row := range rows {
			flows = append(flows, flowInfo{ID: row.ID, Title: row.Title})
		}
	case model.UsageStatsPeriodMonth:
		rows, err := r.DB.GetFlowsForPeriodLastMonth(ctx, uid)
		if err != nil {
			return nil, err
		}
		for _, row := range rows {
			flows = append(flows, flowInfo{ID: row.ID, Title: row.Title})
		}
	case model.UsageStatsPeriodQuarter:
		rows, err := r.DB.GetFlowsForPeriodLast3Months(ctx, uid)
		if err != nil {
			return nil, err
		}
		for _, row := range rows {
			flows = append(flows, flowInfo{ID: row.ID, Title: row.Title})
		}
	default:
		return nil, fmt.Errorf("unsupported period: %s", period)
	}

	// Step 2: Build stats for each flow using analytics functions
	result := make([]*model.FlowExecutionStats, 0, len(flows))

	for _, flow := range flows {
		// Get raw data for this flow
		tasks, err := r.DB.GetTasksForFlow(ctx, flow.ID)
		if err != nil {
			return nil, err
		}

		// Collect task IDs
		taskIDs := make([]int64, len(tasks))
		for i, task := range tasks {
			taskIDs[i] = task.ID
		}

		// Get subtasks for all tasks
		var subtasks []database.GetSubtasksForTasksRow
		if len(taskIDs) > 0 {
			subtasks, err = r.DB.GetSubtasksForTasks(ctx, taskIDs)
			if err != nil {
				return nil, err
			}
		}

		// Get msgchains for the flow
		msgchains, err := r.DB.GetMsgchainsForFlow(ctx, flow.ID)
		if err != nil {
			return nil, err
		}

		// Get toolcalls for the flow
		toolcalls, err := r.DB.GetToolcallsForFlow(ctx, flow.ID)
		if err != nil {
			return nil, err
		}

		// Get assistants count for the flow
		assistantsCount, err := r.DB.GetAssistantsCountForFlow(ctx, flow.ID)
		if err != nil {
			return nil, err
		}

		// Build execution stats using analytics functions
		flowStats := converter.BuildFlowExecutionStats(flow.ID, flow.Title, tasks, subtasks, msgchains, toolcalls, int(assistantsCount))
		result = append(result, flowStats)
	}

	return result, nil
}
