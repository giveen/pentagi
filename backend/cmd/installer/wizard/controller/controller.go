package controller

import (
	"net/url"
	"slices"
	"strings"

	"pentagi/cmd/installer/checker"
	"pentagi/cmd/installer/files"
	"pentagi/cmd/installer/loader"
	"pentagi/cmd/installer/state"
)

const (
	EmbeddedLLMConfigsPath   = "providers-configs"
	DefaultDockerCertPath    = "/opt/pentagi/docker/ssl"
	DefaultCustomConfigsPath = "/opt/pentagi/conf/custom.provider.yml"
	DefaultOllamaConfigsPath = "/opt/pentagi/conf/ollama.provider.yml"
	DefaultLLMConfigsPath    = "/opt/pentagi/conf/"
	DefaultScraperBaseURL    = "https://scraper/"
	DefaultScraperDomain     = "scraper"
	DefaultScraperSchema     = "https"
)

type Controller interface {
	GetState() state.State
	GetChecker() *checker.CheckResult

	state.State

	LLMProviderConfigController
	LangfuseConfigController
	GraphitiConfigController
	ObservabilityConfigController
	SummarizerConfigController
	EmbedderConfigController
	AIAgentsConfigController
	ScraperConfigController
	SearchEnginesConfigController
	DockerConfigController
	ChangesConfigController
	ServerSettingsConfigController
}

type LLMProviderConfigController interface {
	GetLLMProviders() map[string]*LLMProviderConfig
	GetLLMProviderConfig(providerID string) *LLMProviderConfig
	UpdateLLMProviderConfig(providerID string, config *LLMProviderConfig) error
	ResetLLMProviderConfig(providerID string) map[string]*LLMProviderConfig
}

type LangfuseConfigController interface {
	GetLangfuseConfig() *LangfuseConfig
	UpdateLangfuseConfig(config *LangfuseConfig) error
	ResetLangfuseConfig() *LangfuseConfig
}

type GraphitiConfigController interface {
	GetGraphitiConfig() *GraphitiConfig
	UpdateGraphitiConfig(config *GraphitiConfig) error
	ResetGraphitiConfig() *GraphitiConfig
}

type ObservabilityConfigController interface {
	GetObservabilityConfig() *ObservabilityConfig
	UpdateObservabilityConfig(config *ObservabilityConfig) error
	ResetObservabilityConfig() *ObservabilityConfig
}

type SummarizerConfigController interface {
	GetSummarizerConfig(summarizerType SummarizerType) *SummarizerConfig
	UpdateSummarizerConfig(config *SummarizerConfig) error
	ResetSummarizerConfig(summarizerType SummarizerType) *SummarizerConfig
}

type EmbedderConfigController interface {
	GetEmbedderConfig() *EmbedderConfig
	UpdateEmbedderConfig(config *EmbedderConfig) error
	ResetEmbedderConfig() *EmbedderConfig
}

type AIAgentsConfigController interface {
	GetAIAgentsConfig() *AIAgentsConfig
	UpdateAIAgentsConfig(config *AIAgentsConfig) error
	ResetAIAgentsConfig() *AIAgentsConfig
}

type ScraperConfigController interface {
	GetScraperConfig() *ScraperConfig
	UpdateScraperConfig(config *ScraperConfig) error
	ResetScraperConfig() *ScraperConfig
}

type SearchEnginesConfigController interface {
	GetSearchEnginesConfig() *SearchEnginesConfig
	UpdateSearchEnginesConfig(config *SearchEnginesConfig) error
	ResetSearchEnginesConfig() *SearchEnginesConfig
}

type DockerConfigController interface {
	GetDockerConfig() *DockerConfig
	UpdateDockerConfig(config *DockerConfig) error
	ResetDockerConfig() *DockerConfig
}

type ChangesConfigController interface {
	GetApplyChangesConfig() *ApplyChangesConfig
}

type ServerSettingsConfigController interface {
	GetServerSettingsConfig() *ServerSettingsConfig
	UpdateServerSettingsConfig(config *ServerSettingsConfig) error
	ResetServerSettingsConfig() *ServerSettingsConfig
}

// controller bridges TUI models with the state package
type controller struct {
	files   files.Files
	checker checker.CheckResult
	state.State
}

func NewController(state state.State, files files.Files, checker checker.CheckResult) Controller {
	return &controller{
		files:   files,
		checker: checker,
		State:   state,
	}
}

// GetState returns the underlying state interface for processor integration
func (c *controller) GetState() state.State {
	return c.State
}

// GetChecker returns the checker result for processor integration
func (c *controller) GetChecker() *checker.CheckResult {
	return &c.checker
}

// LLMProviderConfig represents LLM provider configuration
type LLMProviderConfig struct {
	// dependent on the provider type
	Name string

	// direct form field mappings using loader.EnvVar
	// these fields directly correspond to environment variables and form inputs (not computed)
	BaseURL loader.EnvVar // OPEN_AI_SERVER_URL | ANTHROPIC_SERVER_URL | GEMINI_SERVER_URL | BEDROCK_SERVER_URL | OLLAMA_SERVER_URL | DEEPSEEK_SERVER_URL | GLM_SERVER_URL | KIMI_SERVER_URL | QWEN_SERVER_URL | LLM_SERVER_URL
	APIKey  loader.EnvVar // OPEN_AI_KEY | ANTHROPIC_API_KEY | GEMINI_API_KEY | LLM_SERVER_KEY | DEEPSEEK_API_KEY | GLM_API_KEY | KIMI_API_KEY | QWEN_API_KEY | OLLAMA_SERVER_API_KEY
	Model   loader.EnvVar // LLM_SERVER_MODEL
	// AWS Bedrock specific fields
	DefaultAuth  loader.EnvVar // BEDROCK_DEFAULT_AUTH
	BearerToken  loader.EnvVar // BEDROCK_BEARER_TOKEN
	AccessKey    loader.EnvVar // BEDROCK_ACCESS_KEY_ID
	SecretKey    loader.EnvVar // BEDROCK_SECRET_ACCESS_KEY
	SessionToken loader.EnvVar // BEDROCK_SESSION_TOKEN
	Region       loader.EnvVar // BEDROCK_REGION
	// Ollama and Custom specific fields
	ConfigPath        loader.EnvVar // OLLAMA_SERVER_CONFIG_PATH | LLM_SERVER_CONFIG_PATH
	HostConfigPath    loader.EnvVar // PENTAGI_OLLAMA_SERVER_CONFIG_PATH | PENTAGI_LLM_SERVER_CONFIG_PATH
	LegacyReasoning   loader.EnvVar // LLM_SERVER_LEGACY_REASONING
	PreserveReasoning loader.EnvVar // LLM_SERVER_PRESERVE_REASONING
	// Custom specific fields
	ProviderName loader.EnvVar // LLM_SERVER_PROVIDER | DEEPSEEK_PROVIDER | GLM_PROVIDER | KIMI_PROVIDER | QWEN_PROVIDER
	// Ollama specific fields
	PullTimeout       loader.EnvVar // OLLAMA_SERVER_PULL_MODELS_TIMEOUT
	PullEnabled       loader.EnvVar // OLLAMA_SERVER_PULL_MODELS_ENABLED
	LoadModelsEnabled loader.EnvVar // OLLAMA_SERVER_LOAD_MODELS_ENABLED

	// computed fields (not directly mapped to env vars)
	Configured bool

	// local path to the embedded LLM config files inside the container
	EmbeddedLLMConfigsPath []string
}

// LangfuseConfig represents Langfuse configuration
type LangfuseConfig struct {
	// deployment configuration
	DeploymentType string // "embedded" or "external" or "disabled"

	// embedded listen settings
	ListenIP   loader.EnvVar // LANGFUSE_LISTEN_IP
	ListenPort loader.EnvVar // LANGFUSE_LISTEN_PORT

	// integration settings (always required)
	BaseURL   loader.EnvVar // LANGFUSE_BASE_URL
	ProjectID loader.EnvVar // LANGFUSE_PROJECT_ID | LANGFUSE_INIT_PROJECT_ID
	PublicKey loader.EnvVar // LANGFUSE_PUBLIC_KEY | LANGFUSE_INIT_PROJECT_PUBLIC_KEY
	SecretKey loader.EnvVar // LANGFUSE_SECRET_KEY | LANGFUSE_INIT_PROJECT_SECRET_KEY

	// embedded instance settings (only for embedded mode)
	AdminEmail    loader.EnvVar // LANGFUSE_INIT_USER_EMAIL
	AdminPassword loader.EnvVar // LANGFUSE_INIT_USER_PASSWORD
	AdminName     loader.EnvVar // LANGFUSE_INIT_USER_NAME

	// enterprise license (optional for embedded mode)
	LicenseKey loader.EnvVar // LANGFUSE_EE_LICENSE_KEY

	// computed fields (not directly mapped to env vars)
	Installed bool
}

// GraphitiConfig represents Graphiti knowledge graph configuration
type GraphitiConfig struct {
	// deployment configuration
	DeploymentType string // "embedded" or "external" or "disabled"

	// integration settings (always)
	GraphitiURL loader.EnvVar // GRAPHITI_URL
	Timeout     loader.EnvVar // GRAPHITI_TIMEOUT
	ModelName   loader.EnvVar // GRAPHITI_MODEL_NAME

	// neo4j settings (embedded only)
	Neo4jUser     loader.EnvVar // NEO4J_USER
	Neo4jPassword loader.EnvVar // NEO4J_PASSWORD
	Neo4jDatabase loader.EnvVar // NEO4J_DATABASE
	Neo4jURI      loader.EnvVar // NEO4J_URI

	// computed fields (not directly mapped to env vars)
	Installed bool
}

// ObservabilityConfig represents observability configuration
type ObservabilityConfig struct {
	// deployment configuration
	DeploymentType string // "embedded" or "external" or "disabled"

	// embedded listen settings
	GrafanaListenIP    loader.EnvVar // GRAFANA_LISTEN_IP
	GrafanaListenPort  loader.EnvVar // GRAFANA_LISTEN_PORT
	OTelGrpcListenIP   loader.EnvVar // OTEL_GRPC_LISTEN_IP
	OTelGrpcListenPort loader.EnvVar // OTEL_GRPC_LISTEN_PORT
	OTelHttpListenIP   loader.EnvVar // OTEL_HTTP_LISTEN_IP
	OTelHttpListenPort loader.EnvVar // OTEL_HTTP_LISTEN_PORT

	// integration settings
	OTelHost loader.EnvVar // OTEL_HOST
}

type SummarizerType string

const (
	SummarizerTypeGeneral   SummarizerType = "general"
	SummarizerTypeAssistant SummarizerType = "assistant"
)

// SummarizerConfig represents summarizer configuration settings
type SummarizerConfig struct {
	// type identifier ("general" or "assistant")
	Type SummarizerType

	// common boolean settings
	PreserveLast loader.EnvVar // PREFIX + PRESERVE_LAST
	UseQA        loader.EnvVar // PREFIX + USE_QA
	SumHumanInQA loader.EnvVar // PREFIX + SUM_MSG_HUMAN_IN_QA

	// size settings (in bytes)
	LastSecBytes loader.EnvVar // PREFIX + LAST_SEC_BYTES
	MaxBPBytes   loader.EnvVar // PREFIX + MAX_BP_BYTES
	MaxQABytes   loader.EnvVar // PREFIX + MAX_QA_BYTES

	// count settings
	MaxQASections  loader.EnvVar // PREFIX + MAX_QA_SECTIONS
	KeepQASections loader.EnvVar // PREFIX + KEEP_QA_SECTIONS
}

// GetSummarizerConfig returns summarizer configuration for specified type
// ScraperConfig represents scraper configuration settings
type ScraperConfig struct {
	// direct form field mappings using loader.EnvVar
	// these fields directly correspond to environment variables and form inputs (not computed)
	PublicURL             loader.EnvVar // SCRAPER_PUBLIC_URL
	PrivateURL            loader.EnvVar // SCRAPER_PRIVATE_URL
	LocalUsername         loader.EnvVar // LOCAL_SCRAPER_USERNAME
	LocalPassword         loader.EnvVar // LOCAL_SCRAPER_PASSWORD
	MaxConcurrentSessions loader.EnvVar // LOCAL_SCRAPER_MAX_CONCURRENT_SESSIONS

	// computed fields (not directly mapped to env vars)
	// these are derived from the above EnvVar fields
	Mode string // "embedded", "external", "disabled" - computed from PrivateURL

	// parsed credentials for external mode (extracted from URLs)
	PublicUsername  string
	PublicPassword  string
	PrivateUsername string
	PrivatePassword string
}

// SearchEnginesConfig represents search engines configuration settings
type SearchEnginesConfig struct {
	// direct form field mappings using loader.EnvVar
	// these fields directly correspond to environment variables and form inputs (not computed)
	DuckDuckGoEnabled loader.EnvVar // DUCKDUCKGO_ENABLED
	SploitusEnabled   loader.EnvVar // SPLOITUS_ENABLED
	PerplexityAPIKey  loader.EnvVar // PERPLEXITY_API_KEY
	TavilyAPIKey      loader.EnvVar // TAVILY_API_KEY
	TraversaalAPIKey  loader.EnvVar // TRAVERSAAL_API_KEY
	GoogleAPIKey      loader.EnvVar // GOOGLE_API_KEY
	GoogleCXKey       loader.EnvVar // GOOGLE_CX_KEY
	GoogleLRKey       loader.EnvVar // GOOGLE_LR_KEY

	// duckduckgo extra settings
	DuckDuckGoRegion     loader.EnvVar // DUCKDUCKGO_REGION
	DuckDuckGoSafeSearch loader.EnvVar // DUCKDUCKGO_SAFESEARCH
	DuckDuckGoTimeRange  loader.EnvVar // DUCKDUCKGO_TIME_RANGE

	// perplexity extra settings
	PerplexityModel       loader.EnvVar // PERPLEXITY_MODEL
	PerplexityContextSize loader.EnvVar // PERPLEXITY_CONTEXT_SIZE

	// searxng extra settings
	SearxngURL        loader.EnvVar // SEARXNG_URL
	SearxngCategories loader.EnvVar // SEARXNG_CATEGORIES
	SearxngLanguage   loader.EnvVar // SEARXNG_LANGUAGE
	SearxngSafeSearch loader.EnvVar // SEARXNG_SAFESEARCH
	SearxngTimeRange  loader.EnvVar // SEARXNG_TIME_RANGE
	SearxngTimeout    loader.EnvVar // SEARXNG_TIMEOUT

	// computed fields (not directly mapped to env vars)
	ConfiguredCount int // number of configured engines
}

// DockerConfig represents Docker environment configuration
type DockerConfig struct {
	// direct form field mappings using loader.EnvVar
	// these fields directly correspond to environment variables and form inputs (not computed)
	DockerInside                 loader.EnvVar // DOCKER_INSIDE
	DockerNetAdmin               loader.EnvVar // DOCKER_NET_ADMIN
	DockerSocket                 loader.EnvVar // DOCKER_SOCKET
	DockerNetwork                loader.EnvVar // DOCKER_NETWORK
	DockerPublicIP               loader.EnvVar // DOCKER_PUBLIC_IP
	DockerWorkDir                loader.EnvVar // DOCKER_WORK_DIR
	DockerDefaultImage           loader.EnvVar // DOCKER_DEFAULT_IMAGE
	DockerDefaultImageForPentest loader.EnvVar // DOCKER_DEFAULT_IMAGE_FOR_PENTEST

	// TLS connection settings (optional)
	DockerHost         loader.EnvVar // DOCKER_HOST
	DockerTLSVerify    loader.EnvVar // DOCKER_TLS_VERIFY
	HostDockerCertPath loader.EnvVar // PENTAGI_DOCKER_CERT_PATH

	// computed fields (not directly mapped to env vars)
	Configured bool
}

// ServerSettingsConfig represents PentAGI server settings configuration
type ServerSettingsConfig struct {
	// direct form field mappings using loader.EnvVar
	LicenseKey          loader.EnvVar // LICENSE_KEY
	ListenIP            loader.EnvVar // PENTAGI_LISTEN_IP
	ListenPort          loader.EnvVar // PENTAGI_LISTEN_PORT
	PublicURL           loader.EnvVar // PUBLIC_URL
	CorsOrigins         loader.EnvVar // CORS_ORIGINS
	CookieSigningSalt   loader.EnvVar // COOKIE_SIGNING_SALT
	ProxyURL            loader.EnvVar // PROXY_URL
	HTTPClientTimeout   loader.EnvVar // HTTP_CLIENT_TIMEOUT
	ExternalSSLCAPath   loader.EnvVar // EXTERNAL_SSL_CA_PATH
	ExternalSSLInsecure loader.EnvVar // EXTERNAL_SSL_INSECURE
	SSLDir              loader.EnvVar // PENTAGI_SSL_DIR
	DataDir             loader.EnvVar // PENTAGI_DATA_DIR

	// parsed credentials for proxy server (extracted from URLs)
	ProxyUsername string
	ProxyPassword string
}

// ChangeInfo represents information about a single environment variable change
type ChangeInfo struct {
	Variable    string // environment variable name
	Description string // localized description
	NewValue    string // new value being set
	Masked      bool   // whether the value should be masked in display
}

// ApplyChangesConfig contains information about pending changes and installation status
type ApplyChangesConfig struct {
	// installation state
	IsInstalled bool // whether PentAGI is currently installed

	// deployment selections
	LangfuseEnabled      bool // whether Langfuse embedded deployment is selected
	ObservabilityEnabled bool // whether Observability embedded deployment is selected

	// changes information
	Changes      []ChangeInfo // list of pending environment variable changes
	ChangesCount int          // total number of changes
	HasCritical  bool         // whether there are critical changes requiring restart
	HasSecrets   bool         // whether there are secret/sensitive changes
}

// GetApplyChangesConfig returns the current apply changes configuration
func (c *controller) GetApplyChangesConfig() *ApplyChangesConfig {
	config := &ApplyChangesConfig{
		IsInstalled: c.checker.PentagiInstalled,
		Changes:     []ChangeInfo{},
	}

	// check deployment selections
	langfuseConfig := c.GetLangfuseConfig()
	config.LangfuseEnabled = langfuseConfig.DeploymentType == "embedded"

	observabilityConfig := c.GetObservabilityConfig()
	config.ObservabilityEnabled = observabilityConfig.DeploymentType == "embedded"

	// collect all changed variables
	allVars := c.GetAllVars()
	for varName, envVar := range allVars {
		if envVar.IsChanged {
			description := c.getVariableDescription(varName)
			masked := c.isVariableMasked(varName)
			value := envVar.Value
			if value == "" {
				value = "{EMPTY}"
			}

			config.Changes = append(config.Changes, ChangeInfo{
				Variable:    varName,
				Description: description,
				NewValue:    value,
				Masked:      masked,
			})

			// mark critical and secret changes
			if c.isCriticalVariable(varName) {
				config.HasCritical = true
			}
			if masked {
				config.HasSecrets = true
			}
		}
	}

	slices.SortFunc(config.Changes, func(a, b ChangeInfo) int {
		return strings.Compare(a.Description, b.Description)
	})

	config.ChangesCount = len(config.Changes)
	return config
}

// getVariableDescription returns a user-friendly description for an environment variable
// RemoveCredentialsFromURL removes credentials from URL - public method for form display
func RemoveCredentialsFromURL(urlStr string) string {
	if urlStr == "" {
		return urlStr
	}

	parsedURL, err := url.Parse(urlStr)
	if err != nil {
		return urlStr
	}

	parsedURL.User = nil

	return parsedURL.String()
}
