package providers

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/json"
	"fmt"
	"math"
	"math/big"
	"strings"
	"sync"
	"sync/atomic"
	"time"
	"unicode"

	"pentagi/pkg/config"
	"pentagi/pkg/csum"
	"pentagi/pkg/database"
	"pentagi/pkg/docker"
	"pentagi/pkg/graphiti"
	obs "pentagi/pkg/observability"
	"pentagi/pkg/providers/custom"
	"pentagi/pkg/providers/embeddings"
	"pentagi/pkg/providers/ollama"
	"pentagi/pkg/providers/openai"
	"pentagi/pkg/providers/pconfig"
	"pentagi/pkg/providers/provider"
	"pentagi/pkg/providers/tester"
	"pentagi/pkg/templates"
	"pentagi/pkg/tools"

	"github.com/sirupsen/logrus"
)

const deltaCallCounter = 10000

const defaultTestParallelWorkersNumber = 16

const pentestDockerImage = "vxcontrol/kali-linux"

type ProviderController interface {
	NewFlowProvider(
		ctx context.Context,
		prvname provider.ProviderName,
		prompter templates.Prompter,
		executor tools.FlowToolsExecutor,
		flowID, userID int64,
		askUser bool,
		input string,
	) (FlowProvider, error)
	LoadFlowProvider(
		ctx context.Context,
		prvname provider.ProviderName,
		prompter templates.Prompter,
		executor tools.FlowToolsExecutor,
		flowID, userID int64,
		askUser bool,
		image, language, title, tcIDTemplate string,
	) (FlowProvider, error)
	NewAssistantProvider(
		ctx context.Context,
		prvname provider.ProviderName,
		prompter templates.Prompter,
		executor tools.FlowToolsExecutor,
		assistantID, flowID, userID int64,
		image, input string,
		streamCb StreamMessageHandler,
	) (AssistantProvider, error)
	LoadAssistantProvider(
		ctx context.Context,
		prvname provider.ProviderName,
		prompter templates.Prompter,
		executor tools.FlowToolsExecutor,
		assistantID, flowID, userID int64,
		image, language, title, tcIDTemplate string,
		streamCb StreamMessageHandler,
	) (AssistantProvider, error)

	Embedder() embeddings.Embedder
	GraphitiClient() *graphiti.Client
	DefaultProviders() provider.Providers
	DefaultProvidersConfig() provider.ProvidersConfig
	GetProvider(
		ctx context.Context,
		prvname provider.ProviderName,
		userID int64,
	) (provider.Provider, error)
	GetProviders(
		ctx context.Context,
		userID int64,
	) (provider.Providers, error)

	NewProvider(prv database.Provider) (provider.Provider, error)
	CreateProvider(
		ctx context.Context,
		userID int64,
		prvname provider.ProviderName,
		prvtype provider.ProviderType,
		config *pconfig.ProviderConfig,
	) (database.Provider, error)
	UpdateProvider(
		ctx context.Context,
		userID int64,
		prvID int64,
		prvname provider.ProviderName,
		config *pconfig.ProviderConfig,
	) (database.Provider, error)
	DeleteProvider(
		ctx context.Context,
		userID int64,
		prvID int64,
	) (database.Provider, error)

	TestAgent(
		ctx context.Context,
		prvtype provider.ProviderType,
		agentType pconfig.ProviderOptionsType,
		config *pconfig.AgentConfig,
	) (tester.AgentTestResults, error)
	TestProvider(
		ctx context.Context,
		prvtype provider.ProviderType,
		config *pconfig.ProviderConfig,
	) (tester.ProviderTestResults, error)
}

type providerController struct {
	db             database.Querier
	cfg            *config.Config
	docker         docker.DockerClient
	publicIP       string
	dockerNetwork  string
	embedder       embeddings.Embedder
	graphitiClient *graphiti.Client

	startCallNumber *atomic.Int64

	defaultDockerImageForPentest string

	summarizerAgent     csum.Summarizer
	summarizerAssistant csum.Summarizer

	defaultConfigs provider.ProvidersConfig

	provider.Providers
}

func NewProviderController(
	cfg *config.Config,
	db database.Querier,
	docker docker.DockerClient,
) (ProviderController, error) {
	if cfg == nil {
		return nil, fmt.Errorf("config is required")
	}

	embedder, err := embeddings.New(cfg)
	if err != nil {
		logrus.WithError(err).Errorf("failed to create embedder '%s'", cfg.EmbeddingProvider)
	}

	providers := make(provider.Providers)
	defaultConfigs := make(provider.ProvidersConfig)

	if config, err := openai.DefaultProviderConfig(); err != nil {
		return nil, fmt.Errorf("failed to create openai provider config: %w", err)
	} else {
		defaultConfigs[provider.ProviderOpenAI] = config
	}



	if config, err := ollama.DefaultProviderConfig(cfg); err != nil {
		return nil, fmt.Errorf("failed to create ollama provider config: %w", err)
	} else {
		defaultConfigs[provider.ProviderOllama] = config
	}

	if config, err := custom.DefaultProviderConfig(cfg); err != nil {
		return nil, fmt.Errorf("failed to create custom provider config: %w", err)
	} else {
		defaultConfigs[provider.ProviderCustom] = config
	}



	if cfg.OpenAIKey != "" {
		p, err := openai.New(cfg, provider.DefaultProviderNameOpenAI, defaultConfigs[provider.ProviderOpenAI])
		if err != nil {
			return nil, fmt.Errorf("failed to create openai provider: %w", err)
		}

		providers[provider.DefaultProviderNameOpenAI] = p
	}



	if cfg.OllamaServerURL != "" {
		p, err := ollama.New(cfg, provider.DefaultProviderNameOllama, defaultConfigs[provider.ProviderOllama])
		if err != nil {
			return nil, fmt.Errorf("failed to create ollama provider: %w", err)
		}
		providers[provider.DefaultProviderNameOllama] = p
	}

	if cfg.LLMServerURL != "" && (cfg.LLMServerModel != "" || cfg.LLMServerConfig != "") {
		p, err := custom.New(cfg, provider.DefaultProviderNameCustom, defaultConfigs[provider.ProviderCustom])
		if err != nil {
			return nil, fmt.Errorf("failed to create custom provider: %w", err)
		}

		providers[provider.DefaultProviderNameCustom] = p
	}



	summarizerAgent := csum.NewSummarizer(csum.SummarizerConfig{
		PreserveLast:   cfg.SummarizerPreserveLast,
		UseQA:          cfg.SummarizerUseQA,
		SummHumanInQA:  cfg.SummarizerSumHumanInQA,
		LastSecBytes:   cfg.SummarizerLastSecBytes,
		MaxBPBytes:     cfg.SummarizerMaxBPBytes,
		MaxQASections:  cfg.SummarizerMaxQASections,
		MaxQABytes:     cfg.SummarizerMaxQABytes,
		KeepQASections: cfg.SummarizerKeepQASections,
	})

	summarizerAssistant := csum.NewSummarizer(csum.SummarizerConfig{
		PreserveLast:   cfg.AssistantSummarizerPreserveLast,
		UseQA:          true,
		SummHumanInQA:  false,
		LastSecBytes:   cfg.AssistantSummarizerLastSecBytes,
		MaxBPBytes:     cfg.AssistantSummarizerMaxBPBytes,
		MaxQASections:  cfg.AssistantSummarizerMaxQASections,
		MaxQABytes:     cfg.AssistantSummarizerMaxQABytes,
		KeepQASections: cfg.AssistantSummarizerKeepQASections,
	})

	graphitiClient, err := graphiti.NewClient(
		cfg.GraphitiURL,
		time.Duration(cfg.GraphitiTimeout)*time.Second,
		cfg.GraphitiEnabled && cfg.GraphitiURL != "",
	)
	if err != nil {
		logrus.WithError(err).Warn("failed to initialize graphiti client, continuing without it")
		graphitiClient = &graphiti.Client{}
	}

	return &providerController{
		db:             db,
		cfg:            cfg,
		docker:         docker,
		publicIP:       cfg.DockerPublicIP,
		dockerNetwork:  cfg.DockerNetwork,
		embedder:       embedder,
		graphitiClient: graphitiClient,

		startCallNumber: newAtomicInt64(0), // 0 means to make it random

		defaultDockerImageForPentest: cfg.DockerDefaultImageForPentest,

		summarizerAgent:     summarizerAgent,
		summarizerAssistant: summarizerAssistant,

		defaultConfigs: defaultConfigs,

		Providers: providers,
	}, nil
}

func (pc *providerController) NewFlowProvider(
	ctx context.Context,
	prvname provider.ProviderName,
	prompter templates.Prompter,
	executor tools.FlowToolsExecutor,
	flowID, userID int64,
	askUser bool,
	input string,
) (FlowProvider, error) {
	ctx, span := obs.Observer.NewSpan(ctx, obs.SpanKindInternal, "providers.NewFlowProvider")
	defer span.End()

	prv, err := pc.GetProvider(ctx, prvname, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get provider: %w", err)
	}

	var image string
	if isPentestTask(input) {
		image = pc.defaultDockerImageForPentest
	} else {
		image = pc.docker.GetDefaultImage()
	}
	image = pc.normalizeFlowImage(image)

	language := detectLanguage(input)

	title := generateTitleHeuristic(input)

	tcIDTemplate, err := prv.GetToolCallIDTemplate(ctx, prompter)
	if err != nil {
		return nil, fmt.Errorf("failed to determine tool call ID template: %w", err)
	}

	fp := &flowProvider{
		db:              pc.db,
		mx:              &sync.RWMutex{},
		embedder:        pc.embedder,
		graphitiClient:  pc.graphitiClient,
		flowID:          flowID,
		publicIP:        pc.publicIP,
		dockerNetwork:   pc.dockerNetwork,
		callCounter:     newAtomicInt64(pc.startCallNumber.Add(deltaCallCounter)),
		image:           image,
		title:           title,
		language:        language,
		askUser:         askUser,
		planning:        pc.cfg.AgentPlanningStepEnabled,
		tcIDTemplate:    tcIDTemplate,
		prompter:        prompter,
		executor:        executor,
		summarizer:      pc.summarizerAgent,
		Provider:        prv,
		maxGACallsLimit: pc.cfg.MaxGeneralAgentToolCalls,
		maxLACallsLimit: pc.cfg.MaxLimitedAgentToolCalls,
		buildMonitor: func() *executionMonitor {
			return &executionMonitor{
				enabled:        pc.cfg.ExecutionMonitorEnabled,
				sameThreshold:  pc.cfg.ExecutionMonitorSameToolLimit,
				totalThreshold: pc.cfg.ExecutionMonitorTotalToolLimit,
			}
		},
	}

	return fp, nil
}

func (pc *providerController) LoadFlowProvider(
	ctx context.Context,
	prvname provider.ProviderName,
	prompter templates.Prompter,
	executor tools.FlowToolsExecutor,
	flowID, userID int64,
	askUser bool,
	image, language, title, tcIDTemplate string,
) (FlowProvider, error) {
	ctx, span := obs.Observer.NewSpan(ctx, obs.SpanKindInternal, "providers.LoadFlowProvider")
	defer span.End()

	prv, err := pc.GetProvider(ctx, prvname, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get provider: %w", err)
	}

	// Warm the in-memory tool-call-ID cache from the DB-stored template so that
	// subsequent NewFlowProvider calls skip the LLM-based sample-collection phase.
	provider.WarmToolCallIDCache(prv.Type(), tcIDTemplate)

	image = pc.normalizeFlowImage(image)

	fp := &flowProvider{
		db:              pc.db,
		mx:              &sync.RWMutex{},
		embedder:        pc.embedder,
		graphitiClient:  pc.graphitiClient,
		flowID:          flowID,
		publicIP:        pc.publicIP,
		dockerNetwork:   pc.dockerNetwork,
		callCounter:     newAtomicInt64(pc.startCallNumber.Add(deltaCallCounter)),
		image:           image,
		title:           title,
		language:        language,
		askUser:         askUser,
		planning:        pc.cfg.AgentPlanningStepEnabled,
		tcIDTemplate:    tcIDTemplate,
		prompter:        prompter,
		executor:        executor,
		summarizer:      pc.summarizerAgent,
		Provider:        prv,
		maxGACallsLimit: pc.cfg.MaxGeneralAgentToolCalls,
		maxLACallsLimit: pc.cfg.MaxLimitedAgentToolCalls,
		buildMonitor: func() *executionMonitor {
			return &executionMonitor{
				enabled:        pc.cfg.ExecutionMonitorEnabled,
				sameThreshold:  pc.cfg.ExecutionMonitorSameToolLimit,
				totalThreshold: pc.cfg.ExecutionMonitorTotalToolLimit,
			}
		},
	}

	return fp, nil
}

func (pc *providerController) normalizeFlowImage(image string) string {
	normalizedImage := strings.ToLower(strings.TrimSpace(image))
	defaultImage := strings.ToLower(strings.TrimSpace(pc.docker.GetDefaultImage()))
	pentestImage := strings.ToLower(strings.TrimSpace(pc.defaultDockerImageForPentest))

	if pentestImage == "" {
		return normalizedImage
	}

	if normalizedImage == "" || normalizedImage == defaultImage {
		return pentestImage
	}

	return normalizedImage
}

// detectLanguage identifies the natural language of the input text using
// Unicode script ranges, avoiding an LLM round-trip for a trivial task.
func detectLanguage(input string) string {
	counts := make(map[string]int)
	for _, r := range input {
		switch {
		case unicode.Is(unicode.Han, r):
			counts["Chinese"]++
		case unicode.Is(unicode.Hiragana, r) || unicode.Is(unicode.Katakana, r):
			counts["Japanese"]++
		case unicode.Is(unicode.Hangul, r):
			counts["Korean"]++
		case unicode.Is(unicode.Cyrillic, r):
			counts["Russian"]++
		case unicode.Is(unicode.Arabic, r):
			counts["Arabic"]++
		case unicode.Is(unicode.Hebrew, r):
			counts["Hebrew"]++
		case unicode.Is(unicode.Thai, r):
			counts["Thai"]++
		case unicode.Is(unicode.Devanagari, r):
			counts["Hindi"]++
		case unicode.Is(unicode.Greek, r):
			counts["Greek"]++
		}
	}
	best, bestN := "English", 0
	for lang, n := range counts {
		if n > bestN {
			best, bestN = lang, n
		}
	}
	return best
}

// isPentestTask identifies if the input is a penetration testing task using keyword heuristics,
// avoiding an LLM round-trip for image selection.
func isPentestTask(input string) bool {
	lowerInput := strings.ToLower(input)
	pentestKeywords := []string{
		"pentest", "penetration", "exploit", "vulnerability",
		"attack", "security test", "security audit", "red team",
		"hacking", "breach", "payload", "reverse shell", "webshell",
		"sql inject", "xss", "rce", "privilege escalat", "lateral mov",
		"enumerat", "recon", "footprint", "scan", "bruteforce",
	}
	for _, keyword := range pentestKeywords {
		if strings.Contains(lowerInput, keyword) {
			return true
		}
	}
	return false
}

// generateTitleHeuristic extracts a quick title from input without an LLM call.
// Extracts the first sentence (up to 80 chars) or the first 80 chars,
// trimming common articles and providing a reasonable default.
func generateTitleHeuristic(input string) string {
	input = strings.TrimSpace(input)
	if input == "" {
		return "New Task"
	}

	// Find first sentence (up to period, newline, or 80 chars)
	title := input
	if idx := strings.IndexAny(title, ".\n"); idx > 0 && idx < 80 {
		title = title[:idx]
	} else if len(title) > 80 {
		title = title[:80]
		// Trim back to last word boundary
		if idx := strings.LastIndex(title, " "); idx > 20 {
			title = title[:idx]
		}
	}
	title = strings.TrimSpace(title)

	// Remove common articles and prefixes for brevity
	title = strings.TrimPrefix(strings.TrimPrefix(strings.TrimPrefix(title, "Can you "), "Please "), "I want to ")

	return strings.TrimSpace(title)
}

func (pc *providerController) Embedder() embeddings.Embedder {
	return pc.embedder
}

func (pc *providerController) GraphitiClient() *graphiti.Client {
	return pc.graphitiClient
}

func (pc *providerController) NewAssistantProvider(
	ctx context.Context,
	prvname provider.ProviderName,
	prompter templates.Prompter,
	executor tools.FlowToolsExecutor,
	assistantID, flowID, userID int64,
	image, input string,
	streamCb StreamMessageHandler,
) (AssistantProvider, error) {
	ctx, span := obs.Observer.NewSpan(ctx, obs.SpanKindInternal, "providers.NewAssistantProvider")
	defer span.End()

	prv, err := pc.GetProvider(ctx, prvname, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get provider: %w", err)
	}

	language := detectLanguage(input)

	title := generateTitleHeuristic(input)

	tcIDTemplate, err := prv.GetToolCallIDTemplate(ctx, prompter)
	if err != nil {
		return nil, fmt.Errorf("failed to determine tool call ID template: %w", err)
	}

	ap := &assistantProvider{
		id:         assistantID,
		summarizer: pc.summarizerAssistant,
		fp: flowProvider{
			db:              pc.db,
			mx:              &sync.RWMutex{},
			embedder:        pc.embedder,
			graphitiClient:  pc.graphitiClient,
			flowID:          flowID,
			publicIP:        pc.publicIP,
			dockerNetwork:   pc.dockerNetwork,
			callCounter:     newAtomicInt64(pc.startCallNumber.Add(deltaCallCounter)),
			image:           image,
			title:           title,
			language:        language,
			tcIDTemplate:    tcIDTemplate,
			prompter:        prompter,
			executor:        executor,
			streamCb:        streamCb,
			summarizer:      pc.summarizerAgent,
			Provider:        prv,
			maxGACallsLimit: pc.cfg.MaxGeneralAgentToolCalls,
			maxLACallsLimit: pc.cfg.MaxLimitedAgentToolCalls,
			buildMonitor: func() *executionMonitor {
				return &executionMonitor{
					enabled:        pc.cfg.ExecutionMonitorEnabled,
					sameThreshold:  pc.cfg.ExecutionMonitorSameToolLimit,
					totalThreshold: pc.cfg.ExecutionMonitorTotalToolLimit,
				}
			},
		},
	}

	return ap, nil
}

func (pc *providerController) LoadAssistantProvider(
	ctx context.Context,
	prvname provider.ProviderName,
	prompter templates.Prompter,
	executor tools.FlowToolsExecutor,
	assistantID, flowID, userID int64,
	image, language, title, tcIDTemplate string,
	streamCb StreamMessageHandler,
) (AssistantProvider, error) {
	ctx, span := obs.Observer.NewSpan(ctx, obs.SpanKindInternal, "providers.LoadAssistantProvider")
	defer span.End()

	prv, err := pc.GetProvider(ctx, prvname, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get provider: %w", err)
	}

	// Warm cache from DB-stored template (same rationale as LoadFlowProvider).
	provider.WarmToolCallIDCache(prv.Type(), tcIDTemplate)

	ap := &assistantProvider{
		id:         assistantID,
		summarizer: pc.summarizerAssistant,
		fp: flowProvider{
			db:              pc.db,
			mx:              &sync.RWMutex{},
			embedder:        pc.embedder,
			graphitiClient:  pc.graphitiClient,
			flowID:          flowID,
			publicIP:        pc.publicIP,
			dockerNetwork:   pc.dockerNetwork,
			callCounter:     newAtomicInt64(pc.startCallNumber.Add(deltaCallCounter)),
			image:           image,
			title:           title,
			language:        language,
			tcIDTemplate:    tcIDTemplate,
			prompter:        prompter,
			executor:        executor,
			streamCb:        streamCb,
			summarizer:      pc.summarizerAgent,
			Provider:        prv,
			maxGACallsLimit: pc.cfg.MaxGeneralAgentToolCalls,
			maxLACallsLimit: pc.cfg.MaxLimitedAgentToolCalls,
			buildMonitor: func() *executionMonitor {
				return &executionMonitor{
					enabled:        pc.cfg.ExecutionMonitorEnabled,
					sameThreshold:  pc.cfg.ExecutionMonitorSameToolLimit,
					totalThreshold: pc.cfg.ExecutionMonitorTotalToolLimit,
				}
			},
		},
	}

	return ap, nil
}

func (pc *providerController) DefaultProviders() provider.Providers {
	return pc.Providers
}

func (pc *providerController) DefaultProvidersConfig() provider.ProvidersConfig {
	return pc.defaultConfigs
}

func (pc *providerController) GetProvider(
	ctx context.Context,
	prvname provider.ProviderName,
	userID int64,
) (provider.Provider, error) {
	// Lookup user defined providers first so they take precedence over built-in providers
	prv, err := pc.db.GetUserProviderByName(ctx, database.GetUserProviderByNameParams{
		Name:   string(prvname),
		UserID: userID,
	})
	if err != nil && err != sql.ErrNoRows {
		return nil, fmt.Errorf("failed to get provider '%s' from database: %w", prvname, err)
	}
	if err == nil {
		return pc.NewProvider(prv)
	}

	// Fall back to built-in default providers
	switch prvname {
	case provider.DefaultProviderNameOpenAI:
		return pc.Providers.Get(provider.DefaultProviderNameOpenAI)
	case provider.DefaultProviderNameOllama:
		return pc.Providers.Get(provider.DefaultProviderNameOllama)
	case provider.DefaultProviderNameCustom:
		return pc.Providers.Get(provider.DefaultProviderNameCustom)
	}

	return nil, fmt.Errorf("provider '%s' not found", prvname)
}

func (pc *providerController) GetProviders(
	ctx context.Context,
	userID int64,
) (provider.Providers, error) {
	providersMap := make(provider.Providers, len(pc.Providers))

	// Copy default providers
	for prvname, prv := range pc.Providers {
		providersMap[prvname] = prv
	}

	// Copy user providers
	providers, err := pc.db.GetUserProviders(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user providers: %w", err)
	}

	for _, prv := range providers {
		p, err := pc.NewProvider(prv)
		if err != nil {
			return nil, fmt.Errorf("failed to build provider: %w", err)
		}
		providersMap[provider.ProviderName(prv.Name)] = p
	}

	return providersMap, nil
}

func (pc *providerController) NewProvider(prv database.Provider) (provider.Provider, error) {
	if len(prv.Config) == 0 {
		prv.Config = []byte(pconfig.EmptyProviderConfigRaw)
	}

	// Check if the provider type is available via check default one
	providerName := provider.ProviderName(prv.Name)
	providerType := provider.ProviderType(prv.Type)
	if !pc.ListTypes().Contains(providerType) {
		return nil, fmt.Errorf("provider type '%s' is not available", prv.Type)
	}

	switch providerType {
	case provider.ProviderOpenAI:
		openaiConfig, err := openai.BuildProviderConfig(prv.Config)
		if err != nil {
			return nil, fmt.Errorf("failed to build openai provider config: %w", err)
		}
		return openai.New(pc.cfg, providerName, openaiConfig)
	case provider.ProviderOllama:
		ollamaConfig, err := ollama.BuildProviderConfig(pc.cfg, prv.Config)
		if err != nil {
			return nil, fmt.Errorf("failed to build ollama provider config: %w", err)
		}
		return ollama.New(pc.cfg, providerName, ollamaConfig)
	case provider.ProviderCustom:
		customConfig, err := custom.BuildProviderConfig(pc.cfg, prv.Config)
		if err != nil {
			return nil, fmt.Errorf("failed to build custom provider config: %w", err)
		}
		return custom.New(pc.cfg, providerName, customConfig)
	default:
		return nil, fmt.Errorf("unknown provider type: %s", prv.Type)
	}
}

func (pc *providerController) CreateProvider(
	ctx context.Context,
	userID int64,
	prvname provider.ProviderName,
	prvtype provider.ProviderType,
	config *pconfig.ProviderConfig,
) (database.Provider, error) {
	ctx, span := obs.Observer.NewSpan(ctx, obs.SpanKindInternal, "providers.CreateProvider")
	defer span.End()

	var (
		err    error
		result database.Provider
	)

	if config, err = pc.patchProviderConfig(prvtype, config); err != nil {
		return result, fmt.Errorf("failed to patch provider config: %w", err)
	}

	rawConfig, err := json.Marshal(config)
	if err != nil {
		return result, fmt.Errorf("failed to marshal provider config: %w", err)
	}

	result, err = pc.db.CreateProvider(ctx, database.CreateProviderParams{
		UserID: userID,
		Type:   database.ProviderType(prvtype),
		Name:   string(prvname),
		Config: rawConfig,
	})
	if err != nil {
		return result, fmt.Errorf("failed to create provider: %w", err)
	}

	return result, nil
}

func (pc *providerController) UpdateProvider(
	ctx context.Context,
	userID int64,
	prvID int64,
	prvname provider.ProviderName,
	config *pconfig.ProviderConfig,
) (database.Provider, error) {
	ctx, span := obs.Observer.NewSpan(ctx, obs.SpanKindInternal, "providers.UpdateProvider")
	defer span.End()

	var (
		err    error
		result database.Provider
	)

	prv, err := pc.db.GetUserProvider(ctx, database.GetUserProviderParams{
		ID:     prvID,
		UserID: userID,
	})
	if err != nil {
		return result, fmt.Errorf("failed to get provider: %w", err)
	}
	prvtype := provider.ProviderType(prv.Type)

	if config, err = pc.patchProviderConfig(prvtype, config); err != nil {
		return result, fmt.Errorf("failed to patch provider config: %w", err)
	}

	rawConfig, err := json.Marshal(config)
	if err != nil {
		return result, fmt.Errorf("failed to marshal provider config: %w", err)
	}

	result, err = pc.db.UpdateUserProvider(ctx, database.UpdateUserProviderParams{
		ID:     prvID,
		UserID: userID,
		Name:   string(prvname),
		Config: rawConfig,
	})
	if err != nil {
		return result, fmt.Errorf("failed to update provider: %w", err)
	}

	return result, nil
}

func (pc *providerController) DeleteProvider(
	ctx context.Context,
	userID int64,
	prvID int64,
) (database.Provider, error) {
	ctx, span := obs.Observer.NewSpan(ctx, obs.SpanKindInternal, "providers.DeleteProvider")
	defer span.End()

	result, err := pc.db.DeleteUserProvider(ctx, database.DeleteUserProviderParams{
		ID:     prvID,
		UserID: userID,
	})
	if err != nil {
		return result, fmt.Errorf("failed to delete provider: %w", err)
	}

	return result, nil
}

func (pc *providerController) TestAgent(
	ctx context.Context,
	prvtype provider.ProviderType,
	agentType pconfig.ProviderOptionsType,
	config *pconfig.AgentConfig,
) (tester.AgentTestResults, error) {
	ctx, span := obs.Observer.NewSpan(ctx, obs.SpanKindInternal, "providers.TestAgent")
	defer span.End()

	var result tester.AgentTestResults

	// Create provider config with single agent configuration
	testConfig := &pconfig.ProviderConfig{}

	// Set the agent config to the appropriate field based on agent type
	switch agentType {
	case pconfig.OptionsTypeSimple:
		testConfig.Simple = config
	case pconfig.OptionsTypeSimpleJSON:
		testConfig.SimpleJSON = config
	case pconfig.OptionsTypePrimaryAgent:
		testConfig.PrimaryAgent = config
	case pconfig.OptionsTypeAssistant:
		testConfig.Assistant = config
	case pconfig.OptionsTypeGenerator:
		testConfig.Generator = config
	case pconfig.OptionsTypeRefiner:
		testConfig.Refiner = config
	case pconfig.OptionsTypeAdviser:
		testConfig.Adviser = config
	case pconfig.OptionsTypeReflector:
		testConfig.Reflector = config
	case pconfig.OptionsTypeSearcher:
		testConfig.Searcher = config
	case pconfig.OptionsTypeEnricher:
		testConfig.Enricher = config
	case pconfig.OptionsTypeCoder:
		testConfig.Coder = config
	case pconfig.OptionsTypeInstaller:
		testConfig.Installer = config
	case pconfig.OptionsTypePentester:
		testConfig.Pentester = config
	default:
		return result, fmt.Errorf("unsupported agent type: %s", agentType)
	}

	// Patch with defaults
	patchedConfig, err := pc.patchProviderConfig(prvtype, testConfig)
	if err != nil {
		return result, fmt.Errorf("failed to patch provider config: %w", err)
	}

	// Create temporary provider for testing using existing provider logic
	providerName := provider.ProviderName("test-provider")
	tempProvider, err := pc.buildProviderFromConfig(prvtype, providerName, patchedConfig)
	if err != nil {
		return result, fmt.Errorf("failed to create provider for testing: %w", err)
	}

	// Run tests for specific agent type only
	testWorkers := pc.getTestParallelWorkers(prvtype)
	results, err := tester.TestProvider(
		ctx,
		tempProvider,
		tester.WithAgentTypes(agentType),
		tester.WithVerbose(false),
		tester.WithParallelWorkers(testWorkers),
	)
	if err != nil {
		return result, fmt.Errorf("failed to test agent: %w", err)
	}

	// Extract results for the specific agent type
	switch agentType {
	case pconfig.OptionsTypeSimple:
		result = results.Simple
	case pconfig.OptionsTypeSimpleJSON:
		result = results.SimpleJSON
	case pconfig.OptionsTypePrimaryAgent:
		result = results.PrimaryAgent
	case pconfig.OptionsTypeAssistant:
		result = results.Assistant
	case pconfig.OptionsTypeGenerator:
		result = results.Generator
	case pconfig.OptionsTypeRefiner:
		result = results.Refiner
	case pconfig.OptionsTypeAdviser:
		result = results.Adviser
	case pconfig.OptionsTypeReflector:
		result = results.Reflector
	case pconfig.OptionsTypeSearcher:
		result = results.Searcher
	case pconfig.OptionsTypeEnricher:
		result = results.Enricher
	case pconfig.OptionsTypeCoder:
		result = results.Coder
	case pconfig.OptionsTypeInstaller:
		result = results.Installer
	case pconfig.OptionsTypePentester:
		result = results.Pentester
	default:
		return result, fmt.Errorf("unexpected agent type: %s", agentType)
	}

	return result, nil
}

func (pc *providerController) TestProvider(
	ctx context.Context,
	prvtype provider.ProviderType,
	config *pconfig.ProviderConfig,
) (tester.ProviderTestResults, error) {
	ctx, span := obs.Observer.NewSpan(ctx, obs.SpanKindInternal, "providers.TestProvider")
	defer span.End()

	var results tester.ProviderTestResults

	// Patch config with defaults
	patchedConfig, err := pc.patchProviderConfig(prvtype, config)
	if err != nil {
		return results, fmt.Errorf("failed to patch provider config: %w", err)
	}

	// Create provider for testing
	providerName := provider.ProviderName("test-provider")
	testProvider, err := pc.buildProviderFromConfig(prvtype, providerName, patchedConfig)
	if err != nil {
		return results, fmt.Errorf("failed to create provider for testing: %w", err)
	}

	// Run full provider testing
	testWorkers := pc.getTestParallelWorkers(prvtype)
	results, err = tester.TestProvider(
		ctx,
		testProvider,
		tester.WithVerbose(false),
		tester.WithParallelWorkers(testWorkers),
	)
	if err != nil {
		return results, fmt.Errorf("failed to test provider: %w", err)
	}

	return results, nil
}

func (pc *providerController) getTestParallelWorkers(prvtype provider.ProviderType) int {
	if prvtype != provider.ProviderCustom {
		return defaultTestParallelWorkersNumber
	}

	if pc.cfg == nil || pc.cfg.LLMServerTestParallelWorkers <= 0 {
		return 1
	}

	return pc.cfg.LLMServerTestParallelWorkers
}

func (pc *providerController) patchProviderConfig(
	prvtype provider.ProviderType,
	config *pconfig.ProviderConfig,
) (*pconfig.ProviderConfig, error) {
	var (
		defaultCfg *pconfig.ProviderConfig
		ok         bool
	)

	if defaultCfg, ok = pc.defaultConfigs[prvtype]; !ok {
		return nil, fmt.Errorf("default provider config not found for type: %s", prvtype.String())
	}

	if config == nil {
		return defaultCfg, nil
	}

	if config.Simple == nil {
		config.Simple = defaultCfg.Simple
	}
	if config.SimpleJSON == nil {
		config.SimpleJSON = defaultCfg.SimpleJSON
	}
	if config.PrimaryAgent == nil {
		config.PrimaryAgent = defaultCfg.PrimaryAgent
	}
	if config.Assistant == nil {
		config.Assistant = defaultCfg.Assistant
	}
	if config.Generator == nil {
		config.Generator = defaultCfg.Generator
	}
	if config.Refiner == nil {
		config.Refiner = defaultCfg.Refiner
	}
	if config.Adviser == nil {
		config.Adviser = defaultCfg.Adviser
	}
	if config.Reflector == nil {
		config.Reflector = defaultCfg.Reflector
	}
	if config.Searcher == nil {
		config.Searcher = defaultCfg.Searcher
	}
	if config.Enricher == nil {
		config.Enricher = defaultCfg.Enricher
	}
	if config.Coder == nil {
		config.Coder = defaultCfg.Coder
	}
	if config.Installer == nil {
		config.Installer = defaultCfg.Installer
	}
	if config.Pentester == nil {
		config.Pentester = defaultCfg.Pentester
	}

	config.SetDefaultOptions(defaultCfg.GetDefaultOptions())

	return config, nil
}

func (pc *providerController) buildProviderFromConfig(
	prvtype provider.ProviderType,
	prvname provider.ProviderName,
	config *pconfig.ProviderConfig,
) (provider.Provider, error) {
	switch prvtype {
	case provider.ProviderOpenAI:
		return openai.New(pc.cfg, prvname, config)
	case provider.ProviderCustom:
		return custom.New(pc.cfg, prvname, config)
	case provider.ProviderOllama:
		return ollama.New(pc.cfg, prvname, config)
	default:
		return nil, fmt.Errorf("unknown provider type: %s", prvtype)
	}
}

func newAtomicInt64(seed int64) *atomic.Int64 {
	var number atomic.Int64

	if seed == 0 {
		bigID, err := rand.Int(rand.Reader, big.NewInt(math.MaxInt64))
		if err != nil {
			return &number
		}
		seed = bigID.Int64()
	}

	number.Store(seed)
	return &number
}
