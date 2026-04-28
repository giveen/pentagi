package custom

import (
	"context"
	"os"
	"strings"
	"sync"
	"time"

	"pentagi/pkg/config"
	"pentagi/pkg/providers/pconfig"
	"pentagi/pkg/providers/provider"
	"pentagi/pkg/system"
	"pentagi/pkg/templates"

	"github.com/vxcontrol/langchaingo/llms"
	"github.com/vxcontrol/langchaingo/llms/openai"
	"github.com/vxcontrol/langchaingo/llms/streaming"
)

func BuildProviderConfig(cfg *config.Config, configData []byte) (*pconfig.ProviderConfig, error) {
	defaultOptions := []llms.CallOption{
		llms.WithTemperature(1.0),
		llms.WithTopP(1.0),
		llms.WithN(1),
		llms.WithMaxTokens(16384),
	}

	if cfg.LLMServerModel != "" {
		defaultOptions = append(defaultOptions, llms.WithModel(cfg.LLMServerModel))
	}

	providerConfig, err := pconfig.LoadConfigData(configData, defaultOptions)
	if err != nil {
		return nil, err
	}

	return providerConfig, nil
}

func DefaultProviderConfig(cfg *config.Config) (*pconfig.ProviderConfig, error) {
	if cfg.LLMServerConfig == "" {
		return BuildProviderConfig(cfg, []byte(pconfig.EmptyProviderConfigRaw))
	}

	configData, err := os.ReadFile(cfg.LLMServerConfig)
	if err != nil {
		return nil, err
	}

	return BuildProviderConfig(cfg, configData)
}

type customProvider struct {
	llm            *openai.LLM
	model          string
	models         pconfig.ModelsConfig
	providerName   provider.ProviderName
	providerConfig *pconfig.ProviderConfig
	providerPrefix string

	limiter          chan struct{}
	modelSwitchDelay time.Duration
	modelStateMu     sync.Mutex
	lastModel        string
}

func New(
	cfg *config.Config,
	providerName provider.ProviderName,
	providerConfig *pconfig.ProviderConfig,
) (provider.Provider, error) {
	baseKey := strings.TrimSpace(providerConfig.APIKey)
	if baseKey == "" {
		baseKey = cfg.LLMServerKey
	}
	baseURL := strings.TrimSpace(providerConfig.APIURL)
	if baseURL == "" {
		baseURL = cfg.LLMServerURL
	}
	baseModel := cfg.LLMServerModel
	if baseModel == "" && providerConfig.Simple != nil {
		baseModel = providerConfig.Simple.Model
	}
	maxParallel := cfg.LLMServerMaxParallel
	if maxParallel <= 0 {
		maxParallel = 1
	}
	modelSwitchDelay := time.Duration(cfg.LLMServerModelSwitchDelay) * time.Millisecond
	httpClient, err := system.GetHTTPClient(cfg)
	if err != nil {
		return nil, err
	}

	opts := []openai.Option{
		openai.WithToken(baseKey),
		openai.WithModel(baseModel),
		openai.WithBaseURL(baseURL),
		openai.WithHTTPClient(httpClient),
	}
	if !cfg.LLMServerLegacyReasoning {
		opts = append(opts,
			openai.WithUsingReasoningMaxTokens(),
			openai.WithModernReasoningFormat(),
		)
	}
	if cfg.LLMServerPreserveReasoning {
		opts = append(opts,
			openai.WithPreserveReasoningContent(),
		)
	}
	client, err := openai.New(opts...)
	if err != nil {
		return nil, err
	}

	// Use centralized model loading with prefix filtering
	models, err := provider.LoadModelsFromHTTP(baseURL, baseKey, httpClient, cfg.LLMServerProvider)
	if err != nil {
		// If loading fails, fallback to empty models list
		models = pconfig.ModelsConfig{}
	}

	return &customProvider{
		llm:              client,
		model:            baseModel,
		models:           models,
		providerName:     providerName,
		providerConfig:   providerConfig,
		providerPrefix:   cfg.LLMServerProvider,
		limiter:          make(chan struct{}, maxParallel),
		modelSwitchDelay: modelSwitchDelay,
	}, nil
}

func (p *customProvider) acquireLimiter(ctx context.Context) (func(), error) {
	select {
	case p.limiter <- struct{}{}:
		return func() {
			<-p.limiter
		}, nil
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

func (p *customProvider) waitModelSwitch(ctx context.Context, model string) error {
	p.modelStateMu.Lock()
	shouldWait := p.modelSwitchDelay > 0 && p.lastModel != "" && p.lastModel != model
	p.modelStateMu.Unlock()

	if shouldWait {
		timer := time.NewTimer(p.modelSwitchDelay)
		defer timer.Stop()

		select {
		case <-timer.C:
		case <-ctx.Done():
			return ctx.Err()
		}
	}

	p.modelStateMu.Lock()
	p.lastModel = model
	p.modelStateMu.Unlock()

	return nil
}

func (p *customProvider) prepareCall(ctx context.Context, opt pconfig.ProviderOptionsType) (func(), error) {
	release, err := p.acquireLimiter(ctx)
	if err != nil {
		return nil, err
	}

	if err := p.waitModelSwitch(ctx, p.ModelWithPrefix(opt)); err != nil {
		release()
		return nil, err
	}

	return release, nil
}

func (p *customProvider) Type() provider.ProviderType {
	return provider.ProviderCustom
}

func (p *customProvider) Name() provider.ProviderName {
	return p.providerName
}

func (p *customProvider) GetRawConfig() []byte {
	return p.providerConfig.GetRawConfig()
}

func (p *customProvider) GetProviderConfig() *pconfig.ProviderConfig {
	return p.providerConfig
}

func (p *customProvider) GetPriceInfo(opt pconfig.ProviderOptionsType) *pconfig.PriceInfo {
	return p.providerConfig.GetPriceInfoForType(opt)
}

func (p *customProvider) GetModels() pconfig.ModelsConfig {
	return p.models
}

func (p *customProvider) Model(opt pconfig.ProviderOptionsType) string {
	model := p.model
	opts := llms.CallOptions{Model: &model}
	for _, option := range p.providerConfig.GetOptionsForType(opt) {
		option(&opts)
	}

	return opts.GetModel()
}

func (p *customProvider) ModelWithPrefix(opt pconfig.ProviderOptionsType) string {
	return provider.ApplyModelPrefix(p.Model(opt), p.providerPrefix)
}

func (p *customProvider) Call(
	ctx context.Context,
	opt pconfig.ProviderOptionsType,
	prompt string,
) (string, error) {
	release, err := p.prepareCall(ctx, opt)
	if err != nil {
		return "", err
	}
	defer release()

	return provider.WrapGenerateFromSinglePrompt(
		ctx, p, opt, p.llm, prompt,
		p.providerConfig.GetOptionsForType(opt)...,
	)
}

func (p *customProvider) CallEx(
	ctx context.Context,
	opt pconfig.ProviderOptionsType,
	chain []llms.MessageContent,
	streamCb streaming.Callback,
) (*llms.ContentResponse, error) {
	release, err := p.prepareCall(ctx, opt)
	if err != nil {
		return nil, err
	}
	defer release()

	return provider.WrapGenerateContent(
		ctx, p, opt, p.llm.GenerateContent, chain,
		append([]llms.CallOption{
			llms.WithStreamingFunc(streamCb),
		}, p.providerConfig.GetOptionsForType(opt)...)...,
	)
}

func (p *customProvider) CallWithTools(
	ctx context.Context,
	opt pconfig.ProviderOptionsType,
	chain []llms.MessageContent,
	tools []llms.Tool,
	streamCb streaming.Callback,
) (*llms.ContentResponse, error) {
	release, err := p.prepareCall(ctx, opt)
	if err != nil {
		return nil, err
	}
	defer release()

	return provider.WrapGenerateContent(
		ctx, p, opt, p.llm.GenerateContent, chain,
		append([]llms.CallOption{
			llms.WithTools(tools),
			llms.WithStreamingFunc(streamCb),
		}, p.providerConfig.GetOptionsForType(opt)...)...,
	)
}

func (p *customProvider) GetUsage(info map[string]any) pconfig.CallUsage {
	return pconfig.NewCallUsage(info)
}

func (p *customProvider) GetToolCallIDTemplate(ctx context.Context, prompter templates.Prompter) (string, error) {
	return provider.DetermineToolCallIDTemplate(ctx, p, pconfig.OptionsTypeSimple, prompter, "")
}
