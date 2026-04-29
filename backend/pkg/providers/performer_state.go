package providers

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"pentagi/pkg/csum"
	"pentagi/pkg/database"
	obs "pentagi/pkg/observability"
	"pentagi/pkg/observability/langfuse"
	"pentagi/pkg/providers/pconfig"
	"pentagi/pkg/tools"

	"github.com/sirupsen/logrus"
	"github.com/vxcontrol/langchaingo/llms"
)

func (fp *flowProvider) processAssistantResult(
	ctx context.Context,
	logger *logrus.Entry,
	chainID int64,
	chain []llms.MessageContent,
	result *callResult,
	summarizer csum.Summarizer,
	summarizerHandler tools.SummarizeHandler,
	durationDelta float64,
) error {
	var err error

	processAssistantResultStartTime := time.Now()

	if fp.streamCb != nil {
		if result.streamID == 0 {
			result.streamID = fp.callCounter.Add(1)
		}
		err := fp.streamCb(ctx, &StreamMessageChunk{
			Type:     StreamMessageChunkTypeUpdate,
			MsgType:  database.MsglogTypeAnswer,
			Content:  result.content,
			Thinking: result.thinking,
			StreamID: result.streamID,
		})
		if err != nil {
			return fmt.Errorf("failed to stream assistant result: %w", err)
		}
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
		}
	}

	// Preserve reasoning for assistant responses using universal pattern
	msg := llms.MessageContent{Role: llms.ChatMessageTypeAI}
	if result.content != "" || !result.thinking.IsEmpty() {
		msg.Parts = append(msg.Parts, llms.TextPartWithReasoning(result.content, result.thinking))
	}
	chain = append(chain, msg)
	durationDelta += time.Since(processAssistantResultStartTime).Seconds()
	if err := fp.updateMsgChain(ctx, pconfig.OptionsTypeAssistant, chainID, chain, durationDelta); err != nil {
		return fmt.Errorf("failed to update msg chain: %w", err)
	}

	return nil
}

func (fp *flowProvider) updateMsgChain(
	ctx context.Context,
	optAgentType pconfig.ProviderOptionsType,
	chainID int64,
	chain []llms.MessageContent,
	durationDelta float64,
) error {
	chainBlob, err := json.Marshal(chain)
	if err != nil {
		return fmt.Errorf("failed to marshal msg chain: %w", err)
	}

	_, err = fp.db.UpdateMsgChain(ctx, database.UpdateMsgChainParams{
		Chain:           chainBlob,
		DurationSeconds: durationDelta,
		Model:           fp.Model(optAgentType),
		ModelProvider:   string(fp.Type()),
		ID:              chainID,
	})
	if err != nil {
		return fmt.Errorf("failed to update msg chain in DB: %w", err)
	}

	return nil
}

func (fp *flowProvider) updateMsgChainUsage(
	ctx context.Context,
	chainID int64,
	optAgentType pconfig.ProviderOptionsType,
	info map[string]any,
	durationDelta float64,
) error {
	usage := fp.GetUsage(info)
	if usage.IsZero() {
		return nil
	}

	price := fp.GetPriceInfo(optAgentType)
	if price != nil {
		usage.UpdateCost(price)
	}

	_, err := fp.db.UpdateMsgChainUsage(ctx, database.UpdateMsgChainUsageParams{
		UsageIn:         usage.Input,
		UsageOut:        usage.Output,
		UsageCacheIn:    usage.CacheRead,
		UsageCacheOut:   usage.CacheWrite,
		UsageCostIn:     usage.CostInput,
		UsageCostOut:    usage.CostOutput,
		DurationSeconds: durationDelta,
		ID:              chainID,
	})
	if err != nil {
		return fmt.Errorf("failed to update msg chain usage in DB: %w", err)
	}

	return nil
}
