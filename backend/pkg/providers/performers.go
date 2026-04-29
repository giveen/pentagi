package providers

import (
"context"
"encoding/json"
"fmt"

"pentagi/pkg/database"
"pentagi/pkg/providers/pconfig"
"pentagi/pkg/tools"

"github.com/sirupsen/logrus"
)

func (fp *flowProvider) performCoder(
	ctx context.Context,
	taskID, subtaskID *int64,
	systemCoderTmpl, userCoderTmpl, question string,
) (string, error) {
	var (
		codeResult   tools.CodeResult
		optAgentType = pconfig.OptionsTypeCoder
		msgChainType = database.MsgchainTypeCoder
	)

	adviser, err := fp.GetAskAdviceHandler(ctx, taskID, subtaskID)
	if err != nil {
		return "", fmt.Errorf("failed to get adviser handler: %w", err)
	}

	installer, err := fp.GetInstallerHandler(ctx, taskID, subtaskID)
	if err != nil {
		return "", fmt.Errorf("failed to get installer handler: %w", err)
	}

	memorist, err := fp.GetMemoristHandler(ctx, taskID, subtaskID)
	if err != nil {
		return "", fmt.Errorf("failed to get memorist handler: %w", err)
	}

	searcher, err := fp.GetSubtaskSearcherHandler(ctx, taskID, subtaskID)
	if err != nil {
		return "", fmt.Errorf("failed to get searcher handler: %w", err)
	}

	ctx = tools.PutAgentContext(ctx, msgChainType)
	cfg := tools.CoderExecutorConfig{
		TaskID:    taskID,
		SubtaskID: subtaskID,
		Adviser:   adviser,
		Installer: installer,
		Memorist:  memorist,
		Searcher:  searcher,
		CodeResult: func(ctx context.Context, name string, args json.RawMessage) (string, error) {
			err := json.Unmarshal(args, &codeResult)
			if err != nil {
				return "", fmt.Errorf("failed to unmarshal result: %w", err)
			}
			return "code result successfully processed", nil
		},
		Summarizer: fp.GetSummarizeResultHandler(taskID, subtaskID),
	}
	executor, err := fp.executor.GetCoderExecutor(cfg)
	if err != nil {
		return "", fmt.Errorf("failed to get coder executor: %w", err)
	}

	if fp.planning {
		userCoderTmplWithPlan, err := fp.performPlanner(
			ctx, taskID, subtaskID, optAgentType, executor, userCoderTmpl, question,
		)
		if err != nil {
			logrus.WithContext(ctx).WithError(err).Warn("failed to get task plan from planner, proceeding without plan")
		} else {
			userCoderTmpl = userCoderTmplWithPlan
		}
	}

	return fp.performAndLogTextResult(
		ctx,
		taskID,
		subtaskID,
		optAgentType,
		msgChainType,
		systemCoderTmpl,
		userCoderTmpl,
		question,
		executor,
		"coder",
		func() string { return codeResult.Result },
	)
}

func (fp *flowProvider) performInstaller(
	ctx context.Context,
	taskID, subtaskID *int64,
	systemInstallerTmpl, userInstallerTmpl, question string,
) (string, error) {
	var (
		maintenanceResult tools.MaintenanceResult
		optAgentType      = pconfig.OptionsTypeInstaller
		msgChainType      = database.MsgchainTypeInstaller
	)

	adviser, err := fp.GetAskAdviceHandler(ctx, taskID, subtaskID)
	if err != nil {
		return "", fmt.Errorf("failed to get adviser handler: %w", err)
	}

	memorist, err := fp.GetMemoristHandler(ctx, taskID, subtaskID)
	if err != nil {
		return "", fmt.Errorf("failed to get memorist handler: %w", err)
	}

	searcher, err := fp.GetSubtaskSearcherHandler(ctx, taskID, subtaskID)
	if err != nil {
		return "", fmt.Errorf("failed to get searcher handler: %w", err)
	}

	ctx = tools.PutAgentContext(ctx, msgChainType)
	cfg := tools.InstallerExecutorConfig{
		TaskID:    taskID,
		SubtaskID: subtaskID,
		Adviser:   adviser,
		Memorist:  memorist,
		Searcher:  searcher,
		MaintenanceResult: func(ctx context.Context, name string, args json.RawMessage) (string, error) {
			err := json.Unmarshal(args, &maintenanceResult)
			if err != nil {
				return "", fmt.Errorf("failed to unmarshal result: %w", err)
			}
			return "maintenance result successfully processed", nil
		},
		Summarizer: fp.GetSummarizeResultHandler(taskID, subtaskID),
	}
	executor, err := fp.executor.GetInstallerExecutor(cfg)
	if err != nil {
		return "", fmt.Errorf("failed to get installer executor: %w", err)
	}

	if fp.planning {
		userInstallerTmplWithPlan, err := fp.performPlanner(
			ctx, taskID, subtaskID, optAgentType, executor, userInstallerTmpl, question,
		)
		if err != nil {
			logrus.WithContext(ctx).WithError(err).Warn("failed to get task plan from planner, proceeding without plan")
		} else {
			userInstallerTmpl = userInstallerTmplWithPlan
		}
	}

	return fp.performAndLogTextResult(
		ctx,
		taskID,
		subtaskID,
		optAgentType,
		msgChainType,
		systemInstallerTmpl,
		userInstallerTmpl,
		question,
		executor,
		"installer",
		func() string { return maintenanceResult.Result },
	)
}

func (fp *flowProvider) performMemorist(
	ctx context.Context,
	taskID, subtaskID *int64,
	systemMemoristTmpl, userMemoristTmpl, question string,
) (string, error) {
	var (
		memoristResult tools.MemoristResult
		optAgentType   = pconfig.OptionsTypeSearcher
		msgChainType   = database.MsgchainTypeMemorist
	)

	ctx = tools.PutAgentContext(ctx, msgChainType)
	cfg := tools.MemoristExecutorConfig{
		TaskID:    taskID,
		SubtaskID: subtaskID,
		SearchResult: func(ctx context.Context, name string, args json.RawMessage) (string, error) {
			err := json.Unmarshal(args, &memoristResult)
			if err != nil {
				return "", fmt.Errorf("failed to unmarshal result: %w", err)
			}
			return "memorist result successfully processed", nil
		},
		Summarizer: fp.GetSummarizeResultHandler(taskID, subtaskID),
	}
	executor, err := fp.executor.GetMemoristExecutor(cfg)
	if err != nil {
		return "", fmt.Errorf("failed to get memorist executor: %w", err)
	}

	return fp.performAndLogTextResult(
		ctx,
		taskID,
		subtaskID,
		optAgentType,
		msgChainType,
		systemMemoristTmpl,
		userMemoristTmpl,
		question,
		executor,
		"memorist",
		func() string { return memoristResult.Result },
	)
}

func (fp *flowProvider) performPentester(
	ctx context.Context,
	taskID, subtaskID *int64,
	systemPentesterTmpl, userPentesterTmpl, question string,
) (string, error) {
	var (
		hackResult   tools.HackResult
		optAgentType = pconfig.OptionsTypePentester
		msgChainType = database.MsgchainTypePentester
	)

	adviser, err := fp.GetAskAdviceHandler(ctx, taskID, subtaskID)
	if err != nil {
		return "", fmt.Errorf("failed to get adviser handler: %w", err)
	}

	coder, err := fp.GetCoderHandler(ctx, taskID, subtaskID)
	if err != nil {
		return "", fmt.Errorf("failed to get coder handler: %w", err)
	}

	installer, err := fp.GetInstallerHandler(ctx, taskID, subtaskID)
	if err != nil {
		return "", fmt.Errorf("failed to get installer handler: %w", err)
	}

	memorist, err := fp.GetMemoristHandler(ctx, taskID, subtaskID)
	if err != nil {
		return "", fmt.Errorf("failed to get memorist handler: %w", err)
	}

	searcher, err := fp.GetSubtaskSearcherHandler(ctx, taskID, subtaskID)
	if err != nil {
		return "", fmt.Errorf("failed to get searcher handler: %w", err)
	}

	ctx = tools.PutAgentContext(ctx, msgChainType)
	cfg := tools.PentesterExecutorConfig{
		TaskID:    taskID,
		SubtaskID: subtaskID,
		Adviser:   adviser,
		Coder:     coder,
		Installer: installer,
		Memorist:  memorist,
		Searcher:  searcher,
		HackResult: func(ctx context.Context, name string, args json.RawMessage) (string, error) {
			err := json.Unmarshal(args, &hackResult)
			if err != nil {
				return "", fmt.Errorf("failed to unmarshal result: %w", err)
			}
			return "hack result successfully processed", nil
		},
		Summarizer: fp.GetSummarizeResultHandler(taskID, subtaskID),
	}
	executor, err := fp.executor.GetPentesterExecutor(cfg)
	if err != nil {
		return "", fmt.Errorf("failed to get pentester executor: %w", err)
	}

	if fp.planning {
		userPentesterTmplWithPlan, err := fp.performPlanner(
			ctx, taskID, subtaskID, optAgentType, executor, userPentesterTmpl, question,
		)
		if err != nil {
			logrus.WithContext(ctx).WithError(err).Warn("failed to get task plan from planner, proceeding without plan")
		} else {
			userPentesterTmpl = userPentesterTmplWithPlan
		}
	}

	return fp.performAndLogTextResult(
		ctx,
		taskID,
		subtaskID,
		optAgentType,
		msgChainType,
		systemPentesterTmpl,
		userPentesterTmpl,
		question,
		executor,
		"pentester",
		func() string { return hackResult.Result },
	)
}

func (fp *flowProvider) performSearcher(
	ctx context.Context,
	taskID, subtaskID *int64,
	systemSearcherTmpl, userSearcherTmpl, question string,
) (string, error) {
	var (
		searchResult tools.SearchResult
		optAgentType = pconfig.OptionsTypeSearcher
		msgChainType = database.MsgchainTypeSearcher
	)

	memorist, err := fp.GetMemoristHandler(ctx, taskID, subtaskID)
	if err != nil {
		return "", fmt.Errorf("failed to get memorist handler: %w", err)
	}

	ctx = tools.PutAgentContext(ctx, msgChainType)
	cfg := tools.SearcherExecutorConfig{
		TaskID:    taskID,
		SubtaskID: subtaskID,
		Memorist:  memorist,
		SearchResult: func(ctx context.Context, name string, args json.RawMessage) (string, error) {
			err := json.Unmarshal(args, &searchResult)
			if err != nil {
				return "", fmt.Errorf("failed to unmarshal result: %w", err)
			}
			return "search result successfully processed", nil
		},
		Summarizer: fp.GetSummarizeResultHandler(taskID, subtaskID),
	}
	executor, err := fp.executor.GetSearcherExecutor(cfg)
	if err != nil {
		return "", fmt.Errorf("failed to get searcher executor: %w", err)
	}

	return fp.performAndLogTextResult(
		ctx,
		taskID,
		subtaskID,
		optAgentType,
		msgChainType,
		systemSearcherTmpl,
		userSearcherTmpl,
		question,
		executor,
		"searcher",
		func() string { return searchResult.Result },
	)
}

func (fp *flowProvider) performEnricher(
	ctx context.Context,
	taskID, subtaskID *int64,
	systemEnricherTmpl, userEnricherTmpl, question string,
) (string, error) {
	var (
		enricherResult tools.EnricherResult
		optAgentType   = pconfig.OptionsTypeEnricher
		msgChainType   = database.MsgchainTypeEnricher
	)

	ctx = tools.PutAgentContext(ctx, msgChainType)
	cfg := tools.EnricherExecutorConfig{
		TaskID:    taskID,
		SubtaskID: subtaskID,
		EnricherResult: func(ctx context.Context, name string, args json.RawMessage) (string, error) {
			err := json.Unmarshal(args, &enricherResult)
			if err != nil {
				return "", fmt.Errorf("failed to unmarshal result: %w", err)
			}
			return "enrich result successfully processed", nil
		},
		Summarizer: fp.GetSummarizeResultHandler(taskID, subtaskID),
	}
	executor, err := fp.executor.GetEnricherExecutor(cfg)
	if err != nil {
		return "", fmt.Errorf("failed to get enricher executor: %w", err)
	}

	return fp.performAndLogTextResult(
		ctx,
		taskID,
		subtaskID,
		optAgentType,
		msgChainType,
		systemEnricherTmpl,
		userEnricherTmpl,
		question,
		executor,
		"enricher",
		func() string { return enricherResult.Result },
	)
}

func (fp *flowProvider) performAndLogTextResult(
	ctx context.Context,
	taskID, subtaskID *int64,
	optAgentType pconfig.ProviderOptionsType,
	msgChainType database.MsgchainType,
	systemTmpl, userTmpl, question string,
	executor tools.ContextToolsExecutor,
	agentName string,
	resultProvider func() string,
) (string, error) {
	msgChainID, chain, err := fp.restoreChain(
		ctx, taskID, subtaskID, optAgentType, msgChainType, systemTmpl, userTmpl,
	)
	if err != nil {
		return "", fmt.Errorf("failed to restore chain: %w", err)
	}

	err = fp.performAgentChain(ctx, optAgentType, msgChainID, taskID, subtaskID, chain, executor, fp.summarizer)
	if err != nil {
		return "", fmt.Errorf("failed to get task %s result: %w", agentName, err)
	}

	result := resultProvider()
	if agentCtx, ok := tools.GetAgentContext(ctx); ok {
		fp.putAgentLog(
			ctx,
			agentCtx.ParentAgentType,
			agentCtx.CurrentAgentType,
			question,
			result,
			taskID,
			subtaskID,
		)
	}

	return result, nil
}

