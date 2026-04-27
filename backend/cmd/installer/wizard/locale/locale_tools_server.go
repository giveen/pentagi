package locale

// Tools screen strings
const (
	ToolsTitle       = "Tools Configuration"
	ToolsDescription = "Enhance agent capabilities with additional tools and options"
	ToolsName        = "Tools"
	ToolsOverview    = `Configure additional tools and capabilities for AI agents.
Each tool can be enabled and configured according to your requirements.

Available settings:
• Human-in-the-loop - Enable user interaction during testing
• AI Agents Settings - Configure global behavior for AI agents
• Search Engines - Configure external search providers
• Scraper - Web content extraction and analysis
• Graphiti (beta) - Temporal knowledge graph for semantic memory
• Docker - Container environment configuration`
)

// Server Settings screen strings
const (
	ServerSettingsFormTitle       = "Server Settings"
	ServerSettingsFormDescription = "Configure PentAGI server network access and public routing"
	ServerSettingsFormName        = "Server Settings"
	ServerSettingsFormOverview    = `• Network binding - control which interface and port PentAGI listens on
• Public URL - external address and optional base path used in redirects
• CORS - allowed origins for browser access
• Proxy - HTTP/HTTPS proxy for outbound traffic to LLM/search providers
• SSL directory - custom certificates directory containing server.crt and server.key (PEM)
• Data directory - persistent storage for agent artifacts and flow workspaces`

	// Field labels and descriptions
	ServerSettingsLicenseKey     = "License Key"
	ServerSettingsLicenseKeyDesc = "PentAGI License Key in format of XXXX-XXXX-XXXX-XXXX"

	ServerSettingsHost     = "Server Host (Listen IP)"
	ServerSettingsHostDesc = "Bind address used by Docker port mapping (e.g., 0.0.0.0 to expose on all interfaces)"

	ServerSettingsPort     = "Server Port (Listen Port)"
	ServerSettingsPortDesc = "External TCP port exposed by Docker for PentAGI web UI"

	ServerSettingsPublicURL     = "Public URL"
	ServerSettingsPublicURLDesc = "Base public URL for redirects and links (supports base path, e.g., https://example.com/pentagi/)"

	ServerSettingsCORSOrigins     = "CORS Origins"
	ServerSettingsCORSOriginsDesc = "Comma-separated list of allowed origins (e.g., https://localhost:8443,https://localhost)"

	ServerSettingsProxyURL     = "HTTP/HTTPS Proxy"
	ServerSettingsProxyURLDesc = "Proxy for outbound requests to LLMs and external tools (not used for Docker API access)"

	ServerSettingsProxyUsername     = "Proxy Username"
	ServerSettingsProxyUsernameDesc = "Username for proxy authentication (optional)"
	ServerSettingsProxyPassword     = "Proxy Password"
	ServerSettingsProxyPasswordDesc = "Password for proxy authentication (optional)"

	ServerSettingsHTTPClientTimeout     = "HTTP Client Timeout"
	ServerSettingsHTTPClientTimeoutDesc = "Timeout in seconds for external API calls (LLM providers, search engines, etc.)"

	ServerSettingsExternalSSLCAPath     = "Custom CA Certificate Path"
	ServerSettingsExternalSSLCAPathDesc = "Path inside container to custom root CA cert (e.g., /opt/pentagi/ssl/ca-bundle.pem)"

	ServerSettingsExternalSSLInsecure     = "Skip SSL Verification"
	ServerSettingsExternalSSLInsecureDesc = "Disable SSL/TLS certificate validation (use only for testing with self-signed certs)"

	ServerSettingsSSLDir     = "SSL Directory"
	ServerSettingsSSLDirDesc = "Directory containing server.crt and server.key in PEM format (server.crt may include fullchain)"

	ServerSettingsDataDir     = "Data Directory"
	ServerSettingsDataDirDesc = "Directory for all agent-generated files; contains flow-N subdirectories used as /work in worker containers"

	ServerSettingsCookieSigningSalt     = "Cookie Signing Salt"
	ServerSettingsCookieSigningSaltDesc = "Secret used to sign cookies (keep private)"

	// Hints for fields overview
	ServerSettingsLicenseKeyHint          = "License Key"
	ServerSettingsHostHint                = "Listen IP"
	ServerSettingsPortHint                = "Listen Port"
	ServerSettingsPublicURLHint           = "Public URL"
	ServerSettingsCORSOriginsHint         = "CORS Origins"
	ServerSettingsProxyURLHint            = "Proxy URL"
	ServerSettingsProxyUsernameHint       = "Proxy Username"
	ServerSettingsProxyPasswordHint       = "Proxy Password"
	ServerSettingsHTTPClientTimeoutHint   = "HTTP Timeout"
	ServerSettingsExternalSSLCAPathHint   = "Custom CA Path"
	ServerSettingsExternalSSLInsecureHint = "Skip SSL Verification"
	ServerSettingsSSLDirHint              = "SSL Directory"
	ServerSettingsDataDirHint             = "Data Directory"

	// Help texts per-field
	ServerSettingsGeneralHelp = `PentAGI exposes its web UI via Docker with configurable host and port.

Public URL must reflect how users reach the server. If using a subpath (e.g., /pentagi/), include it here. CORS controls browser access from specified origins. Proxy affects outbound traffic to LLM/search providers and other external services used by Tools.

SSL directory allows providing custom certificates. When set, server will use server.crt and server.key from that directory. Data directory stores artifacts and working files for flows.`

	ServerSettingsLicenseKeyHelp = `PentAGI License Key in format of XXXX-XXXX-XXXX-XXXX. It's used to communicate with PentAGI Cloud API.`

	ServerSettingsHostHelp = `Bind address for published port in docker-compose mapping.

Examples:
• 127.0.0.1 — local-only access
• 0.0.0.0 — expose on all interfaces`

	ServerSettingsPortHelp = `External port for PentAGI UI. Must be available on the host. Example: 8443.`

	ServerSettingsPublicURLHelp = `Set the public base URL used in redirects and links.

Examples:
• http://localhost:8443
• https://example.com/
• https://example.com/pentagi/ (with base path)`

	ServerSettingsCORSOriginsHelp = `Comma-separated allowed origins for browser access.`

	ServerSettingsProxyURLHelp = `HTTP or HTTPS proxy for outbound requests to LLM providers and external tools. Not used for Docker API communication.`

	ServerSettingsHTTPClientTimeoutHelp = `Timeout in seconds for all external HTTP/HTTPS API calls including:
• LLM provider requests (OpenAI, Anthropic, Bedrock, etc.)
• Search engine queries (Google, Tavily, Perplexity, etc.)
• External tool integrations
• Embedding generation requests

Default: 600 seconds (10 minutes)
Setting to 0 disables timeout (not recommended in production)
Too low values may cause legitimate long-running requests to fail.`

	ServerSettingsExternalSSLCAPathHelp = `Path to custom CA certificate file (PEM format) inside the container.

Must point to /opt/pentagi/ssl/ directory, which is mounted from pentagi-ssl volume on the host.

Examples:
• /opt/pentagi/ssl/ca-bundle.pem
• /opt/pentagi/ssl/corporate-ca.pem

File can contain multiple root and intermediate certificates.`

	ServerSettingsExternalSSLInsecureHelp = `Disable SSL/TLS certificate validation for connections to LLM providers and external services.

⚠ WARNING: Use only for testing with self-signed certificates. Never enable in production.

When enabled, all certificate validation is bypassed, making connections vulnerable to man-in-the-middle attacks.`

	ServerSettingsSSLDirHelp = `Path to directory with server.crt and server.key in PEM format. server.crt may include fullchain. Overrides default generated certificate behavior.`

	ServerSettingsDataDirHelp = `Host directory for persistent data. PentAGI stores agent artifacts under flow-N subdirectories, which map to /work inside worker containers.`

	ServerSettingsCookieSigningSaltHelp = `Secret salt used to sign cookies. Keep it private.`
)

// Human-in-the-loop screen strings
const (
	// AI Agents Settings screen strings
	ToolsAIAgentsSettingsFormTitle       = "AI Agents Settings"
	ToolsAIAgentsSettingsFormDescription = "Configure global behavior for AI agents"
	ToolsAIAgentsSettingsFormName        = "AI Agents Settings"
	ToolsAIAgentsSettingsFormOverview    = `This section configures global behavior of AI agents across PentAGI.

Basic Settings:
• Enable User Interaction: allow agents to request user input when needed
• Use Multi-Agent Mode: enable assistant to orchestrate multiple specialized agents

Execution Monitoring (⚠️  BETA):
• Enable Execution Monitoring: automatic mentor supervision for pattern analysis
• Same Tool Call Threshold: consecutive identical tool calls before mentor review
• Total Tool Call Threshold: total tool calls before mentor review

Tool Call Limits:
• Max Tool Calls (General Agents): prevent runaway executions for Assistant, Primary Agent, Pentester, Coder, Installer
• Max Tool Calls (Limited Agents): prevent runaway executions for Searcher, Enricher, Memorist, etc.

Task Planning (⚠️  BETA):
• Enable Task Planning: generate structured execution plans for specialist agents

⚠️  BETA features are under active development. Enable for testing only.`

	// field labels and descriptions
	ToolsAIAgentsSettingHumanInTheLoop          = "Enable User Interaction"
	ToolsAIAgentsSettingHumanInTheLoopDesc      = "Allow agents to ask for user input when needed"
	ToolsAIAgentsSettingUseAgents               = "Use Multi-Agent Mode"
	ToolsAIAgentsSettingUseAgentsDesc           = "Enable assistant to orchestrate multiple specialized agents"
	ToolsAIAgentsSettingExecutionMonitor        = "Enable Execution Monitoring (beta)"
	ToolsAIAgentsSettingExecutionMonitorDesc    = "Automatically invoke mentor for execution pattern analysis"
	ToolsAIAgentsSettingSameToolLimit           = "Same Tool Call Threshold"
	ToolsAIAgentsSettingSameToolLimitDesc       = "Consecutive identical tool calls before mentor review"
	ToolsAIAgentsSettingTotalToolLimit          = "Total Tool Call Threshold"
	ToolsAIAgentsSettingTotalToolLimitDesc      = "Total tool calls before mentor review"
	ToolsAIAgentsSettingMaxGeneralToolCalls     = "Max Tool Calls (General Agents)"
	ToolsAIAgentsSettingMaxGeneralToolCallsDesc = "Maximum tool calls for Assistant, Primary Agent, Pentester, Coder, Installer"
	ToolsAIAgentsSettingMaxLimitedToolCalls     = "Max Tool Calls (Limited Agents)"
	ToolsAIAgentsSettingMaxLimitedToolCallsDesc = "Maximum tool calls for Searcher, Enricher, Memorist, etc."
	ToolsAIAgentsSettingTaskPlanning            = "Enable Task Planning (beta)"
	ToolsAIAgentsSettingTaskPlanningDesc        = "Generate structured execution plans for specialist agents"

	// help content
	ToolsAIAgentsSettingsHelp = `AI Agents Settings define how agents collaborate, interact with users, and handle execution control.

Basic Settings:
• Enable User Interaction: allow agents to ask for user input when needed
• Use Multi-Agent Mode: enable assistant to orchestrate specialized agents for complex tasks

Execution Monitoring (⚠️  BETA):
Automatically invokes adviser (mentor) to analyze execution patterns, detect loops, suggest alternative strategies, and prevent agents from fixating on single approach. Thresholds: consecutive identical calls (default: 5) and total calls (default: 10).

Task Planning (⚠️  BETA):
Generates 3-7 step execution plans before specialist agents begin work. Prevents scope creep and improves success rates. Works best when adviser uses enhanced configuration (stronger model or maximum reasoning mode).

Tool Call Limits (always active):
Hard limits prevent infinite loops: General agents default 100, Limited agents default 20. Works independently from beta features.

OPEN SOURCE MODELS < 32B (Qwen3.5-27B, DeepSeek-V3, Llama-3.1-70B):
✓ ENABLE both beta features - ESSENTIAL for quality results
✓ Testing shows 2x improvement in result quality vs. baseline
✓ Configure adviser with enhanced settings for best performance
✓ Ideal for air-gapped deployments with local LLM inference

Performance: 2-3x increase in tokens/time, 2x improvement in quality for models < 32B.

⚠️  BETA WARNING: Features under active development. Recommended for open source models < 32B despite beta status. For cloud APIs with larger models, keep disabled.

Note: Changes require service restart.`
)

// Search Engines screen strings
const (
	ToolsSearchEnginesFormTitle       = "Search Engines Configuration"
	ToolsSearchEnginesFormDescription = "Configure search engines for AI agents to gather intelligence during testing"
	ToolsSearchEnginesFormName        = "Search Engines"
	ToolsSearchEnginesFormOverview    = `Available search engines:
• DuckDuckGo - Free search engine (no API key required)
• Sploitus - Security exploits and vulnerabilities database (no API key required)
• Perplexity - AI-powered search with reasoning
• Tavily - Search API for AI applications
• Traversaal - Web scraping and search
• Google Search - Requires API key and Custom Search Engine ID
• Searxng - Internet metasearch engine

Get API keys from:
• Perplexity: https://www.perplexity.ai/
• Tavily: https://tavily.com/
• Traversaal: https://traversaal.ai/
• Google: https://developers.google.com/custom-search/v1/introduction`

	ToolsSearchEnginesDuckDuckGo               = "DuckDuckGo Search"
	ToolsSearchEnginesDuckDuckGoDesc           = "Enable DuckDuckGo search (no API key required)"
	ToolsSearchEnginesDuckDuckGoRegion         = "DuckDuckGo Region"
	ToolsSearchEnginesDuckDuckGoRegionDesc     = "DuckDuckGo region code (e.g., us-en, uk-en, cn-zh)"
	ToolsSearchEnginesDuckDuckGoSafeSearch     = "DuckDuckGo Safe Search"
	ToolsSearchEnginesDuckDuckGoSafeSearchDesc = "DuckDuckGo safe search (strict, moderate, off)"
	ToolsSearchEnginesDuckDuckGoTimeRange      = "DuckDuckGo Time Range"
	ToolsSearchEnginesDuckDuckGoTimeRangeDesc  = "DuckDuckGo time range (d: day, w: week, m: month, y: year)"
	ToolsSearchEnginesSploitus                 = "Sploitus Search"
	ToolsSearchEnginesSploitusDesc             = "Enable Sploitus search for exploits and vulnerabilities (no API key required)"
	ToolsSearchEnginesPerplexityKey            = "Perplexity API Key"
	ToolsSearchEnginesPerplexityKeyDesc        = "API key for Perplexity AI search"
	ToolsSearchEnginesTavilyKey                = "Tavily API Key"
	ToolsSearchEnginesTavilyKeyDesc            = "API key for Tavily search service"
	ToolsSearchEnginesTraversaalKey            = "Traversaal API Key"
	ToolsSearchEnginesTraversaalKeyDesc        = "API key for Traversaal web scraping"
	ToolsSearchEnginesGoogleKey                = "Google Search API Key"
	ToolsSearchEnginesGoogleKeyDesc            = "Google Custom Search API key"
	ToolsSearchEnginesGoogleCX                 = "Google Search Engine ID"
	ToolsSearchEnginesGoogleCXDesc             = "Google Custom Search Engine ID"
	ToolsSearchEnginesGoogleLR                 = "Google Language Restriction"
	ToolsSearchEnginesGoogleLRDesc             = "Google Search Engine language restriction (e.g., lang_en, lang_cn, etc.)"
	ToolsSearchEnginesSearxngURL               = "Searxng Search URL"
	ToolsSearchEnginesSearxngURLDesc           = "Searxng search engine URL"
	ToolsSearchEnginesSearxngCategories        = "Searxng Search Categories"
	ToolsSearchEnginesSearxngCategoriesDesc    = "Searxng search engine categories (e.g., general, it, web, news, technology, science, health, other)"
	ToolsSearchEnginesSearxngLanguage          = "Searxng Search Language"
	ToolsSearchEnginesSearxngLanguageDesc      = "Searxng search engine language (en, ch, fr, de, it, es, pt, ru, zh, empty for all languages)"
	ToolsSearchEnginesSearxngSafeSearch        = "Searxng Safe Search"
	ToolsSearchEnginesSearxngSafeSearchDesc    = "Searxng search engine safe search (0: off, 1: moderate, 2: strict)"
	ToolsSearchEnginesSearxngTimeRange         = "Searxng Time Range"
	ToolsSearchEnginesSearxngTimeRangeDesc     = "Searxng search engine time range (day, month, year)"
	ToolsSearchEnginesSearxngTimeout           = "Searxng Timeout"
	ToolsSearchEnginesSearxngTimeoutDesc       = "Searxng request timeout in seconds"
)

// Scraper screen strings
const (
	ToolsScraperFormTitle       = "Scraper Configuration"
	ToolsScraperFormDescription = "Configure web scraping service"
	ToolsScraperFormName        = "Scraper"
	ToolsScraperFormOverview    = `Web scraper service for content extraction and analysis using vxcontrol/scraper Docker image.

Modes:
• Embedded - Run local scraper container (recommended)
• External - Use external scraper services
• Disabled - No web scraping capabilities

Docker image: https://hub.docker.com/r/vxcontrol/scraper

The scraper supports:
• Public URL access for external links
• Private URL access for internal/local links
• Content extraction and analysis
• Multiple output formats`

	ToolsScraperModeTitle                 = "Scraper Mode"
	ToolsScraperModeDesc                  = "Select how the scraper service should operate"
	ToolsScraperEmbedded                  = "Embedded Container"
	ToolsScraperExternal                  = "External Service"
	ToolsScraperDisabled                  = "Disabled"
	ToolsScraperPublicURL                 = "Public Scraper URL"
	ToolsScraperPublicURLDesc             = "URL for scraping public/external websites. If empty, the same value as private URL will be used."
	ToolsScraperPublicURLEmbeddedDesc     = "URL for embedded scraper (optional override). If empty, the same value as private URL will be used."
	ToolsScraperPrivateURL                = "Private Scraper URL"
	ToolsScraperPrivateURLDesc            = "URL for scraping private/internal websites"
	ToolsScraperPublicUsername            = "Public URL Username"
	ToolsScraperPublicUsernameDesc        = "Username for public scraper access"
	ToolsScraperPublicPassword            = "Public URL Password"
	ToolsScraperPublicPasswordDesc        = "Password for public scraper access"
	ToolsScraperPrivateUsername           = "Private URL Username"
	ToolsScraperPrivateUsernameDesc       = "Username for private scraper access"
	ToolsScraperPrivatePassword           = "Private URL Password"
	ToolsScraperPrivatePasswordDesc       = "Password for private scraper access"
	ToolsScraperLocalUsername             = "Local URL Username"
	ToolsScraperLocalUsernameDesc         = "Username for embedded scraper service"
	ToolsScraperLocalPassword             = "Local URL Password"
	ToolsScraperLocalPasswordDesc         = "Password for embedded scraper service"
	ToolsScraperMaxConcurrentSessions     = "Max Concurrent Sessions"
	ToolsScraperMaxConcurrentSessionsDesc = "Maximum number of concurrent scraping sessions"
	ToolsScraperEmbeddedHelp              = "Embedded mode runs a local scraper container that can access both public and private resources. The default configuration uses https://someuser:somepass@scraper/."
	ToolsScraperExternalHelp              = "External mode uses separate scraper services. Configure different URLs for public and private access as needed."
	ToolsScraperDisabledHelp              = "Scraper is disabled. Web content extraction and analysis capabilities will not be available."
)

// Docker Environment screen strings
const (
	ToolsDockerFormTitle       = "Docker Environment Configuration"
	ToolsDockerFormDescription = "Configure Docker environment for worker containers"
	ToolsDockerFormName        = "Docker Environment"
	ToolsDockerFormOverview    = `• Worker Isolation - Containers provide security boundaries for tasks
• Network Capabilities - Enable privileged network operations for pentesting
• Container Management - Control how workers access Docker daemon
• Storage Configuration - Define workspace and artifact storage
• Image Selection - Set default images for different task types

Critical for penetration testing workflows requiring network scanning, custom tools, and secure task isolation.`

	// General help text
	ToolsDockerGeneralHelp = `Each AI agent task runs in an isolated Docker container with two ports (28000-32000 range) automatically allocated per flow. Worker containers are created on-demand from default images or agent-selected ones.

Basic setup requires enabling capabilities: Docker Access allows spawning additional containers for specialized tools, while Network Admin grants low-level network permissions essential for scanning tools like nmap.

Storage operates via Docker volumes by default, or host directories when Work Directory is specified. Connection settings control the Docker daemon location - local socket for standard setups, or remote TCP with TLS for distributed environments.

Default images serve as fallbacks: general tasks use standard images, while security testing defaults to pentesting-focused containers. Public IP enables reverse shell attacks by providing workers with a reachable address for target callbacks. Usually it's a local interface address of the host machine with Docker daemon running for the workers containers.

Configuration combines based on scenario: enable both capabilities for full pentesting, use Work Directory for persistent artifacts, or configure remote connection for isolated Docker environments.`

	// Container capabilities
	ToolsDockerInside       = "Docker Access"
	ToolsDockerInsideDesc   = "Allow workers to manage Docker containers"
	ToolsDockerNetAdmin     = "Network Admin"
	ToolsDockerNetAdminDesc = "Grant NET_ADMIN capability for network scanning tools like nmap"

	// Connection settings
	ToolsDockerSocket       = "Docker Socket"
	ToolsDockerSocketDesc   = "Path to Docker socket on host filesystem"
	ToolsDockerNetwork      = "Docker Network"
	ToolsDockerNetworkDesc  = "Custom network name for worker containers, or 'host' for direct host network access"
	ToolsDockerPublicIP     = "Public IP Address"
	ToolsDockerPublicIPDesc = "Public IP for reverse connections in OOB attacks"

	// Storage configuration
	ToolsDockerWorkDir     = "Work Directory"
	ToolsDockerWorkDirDesc = "Host directory for worker filesystems (default: Docker volumes)"

	// Default images
	ToolsDockerDefaultImage               = "Default Image"
	ToolsDockerDefaultImageDesc           = "Default Docker image for general tasks"
	ToolsDockerDefaultImageForPentest     = "Pentesting Image"
	ToolsDockerDefaultImageForPentestDesc = "Default Docker image for security testing tasks"

	// TLS connection settings (optional)
	ToolsDockerHost          = "Docker Host"
	ToolsDockerHostDesc      = "Docker daemon connection (unix:// or tcp://)"
	ToolsDockerTLSVerify     = "TLS Verification"
	ToolsDockerTLSVerifyDesc = "Enable TLS verification for Docker connection"
	ToolsDockerCertPath      = "TLS Certificates"
	ToolsDockerCertPathDesc  = "Directory containing ca.pem, cert.pem, key.pem files"

	// Help content for specific configurations
	ToolsDockerInsideHelp = `Docker Access enables workers to spawn additional containers for specialized tools and environments. Required when tasks need custom software not available in default images.

When enabled, workers can pull and run any Docker image, providing maximum flexibility for complex testing scenarios.`

	ToolsDockerNetAdminHelp = `Network Admin capability allows workers to perform low-level network operations essential for penetration testing.

Required for:
• Network scanning with nmap, masscan
• Custom packet crafting
• Network interface manipulation
• Raw socket operations

Critical for comprehensive security assessments.`

	ToolsDockerSocketHelp = `Docker Socket path defines how workers access the Docker daemon. Use only file path to the socket file. Used with Docker Access to enable container management.

For enhanced security, consider using docker-in-docker (DinD) instead of exposing the main Docker daemon directly to workers.
When using DinD, use the path to the Docker socket file of the DinD container which binded to the host filesystem.

Example: /var/run/docker.sock`

	ToolsDockerNetworkHelp = `Docker Network controls network isolation mode for worker containers:

Bridge Mode (custom network name):
• Isolated communication between containers
• Port forwarding from container to host
• Enhanced security boundaries
• Network-based monitoring and filtering
• Recommended for most use cases

Host Mode (value: 'host'):
• Direct access to host network interfaces
• No port forwarding - ports bind directly to host
• Required for raw packet manipulation
• Advanced network testing capabilities
• Lower isolation - use with caution

Examples:
• 'pentagi-network' - creates isolated bridge network
• 'host' - enables direct host network access

Security Note: Host network mode reduces container isolation. Only use when necessary for advanced penetration testing tasks requiring direct network stack access.`

	ToolsDockerPublicIPHelp = `Public IP Address enables out-of-band (OOB) attack techniques by providing workers with a reachable address for reverse connections.

Workers automatically receive two random ports (28000-32000 range) mapped to this IP for receiving callbacks from exploited targets.

By default agents will try to get public address from the services api.ipify.org, ipinfo.io/ip or ifconfig.me.`

	ToolsDockerWorkDirHelp = `Work Directory specifies host filesystem location for worker storage. When set, replaces default Docker volumes with host directory mounts.

Benefits:
• Persistent storage across restarts
• Direct file system access
• Easier artifact management
• Custom backup strategies

By default uses Docker dedicated volume per worker container.

Example: /path/to/workdir/`

	ToolsDockerDefaultImageHelp = `Default Image provides fallback for workers when task requirements don't specify a particular container image.

Should contain basic utilities and tools for general-purpose tasks. Default: debian:latest`

	ToolsDockerDefaultImageForPentestHelp = `Pentesting Image serves as default for security testing tasks. Should include comprehensive security tools and utilities.

Recommended images include Kali Linux, Parrot Security, or custom security-focused containers. Default: vxcontrol/kali-linux`

	ToolsDockerHostHelp = `Docker Host uses for start primary worker containers and overrides default Docker daemon connection. Supports Unix sockets and TCP connections.

Examples:
• unix:///var/run/docker.sock (local)
• tcp://docker-host:2376 (remote)

Enable TLS for remote connections.`

	ToolsDockerTLSVerifyHelp = `TLS Verification secures Docker daemon connections over TCP. Strongly recommended for remote Docker hosts.

Requires valid certificates in the specified certificate directory.`

	ToolsDockerCertPathHelp = `TLS Certificates directory must contain:
• ca.pem - Certificate Authority
• cert.pem - Client certificate
• key.pem - Private key

Required for secure remote Docker connections when using TLS to manage worker containers.

Example: /path/to/certs`
)

// Embedder form strings
const (
	EmbedderFormTitle       = "Embedder Configuration"
	EmbedderFormDescription = "Configure text vectorization for semantic search and knowledge storage"
	EmbedderFormName        = "Embedder"
	EmbedderFormOverview    = `Text embeddings convert documents into vectors for semantic search and knowledge storage.
Different providers offer various models with different capabilities and pricing.

Choose carefully as changing providers requires reindexing all stored data.`

	EmbedderFormProvider     = "Embedding Provider"
	EmbedderFormProviderDesc = "Select the provider for text vectorization. Embeddings are used for semantic search and knowledge storage."

	EmbedderFormURL     = "API Endpoint URL"
	EmbedderFormURLDesc = "Custom API endpoint (leave empty to use default)"

	EmbedderFormAPIKey     = "API Key"
	EmbedderFormAPIKeyDesc = "Authentication key for the provider (not required for Ollama)"

	EmbedderFormModel     = "Model Name"
	EmbedderFormModelDesc = "Specific embedding model to use (leave empty for provider default)"

	EmbedderFormBatchSize     = "Batch Size"
	EmbedderFormBatchSizeDesc = "Number of documents to process in a single batch (1-1000)"

	EmbedderFormStripNewLines     = "Strip New Lines"
	EmbedderFormStripNewLinesDesc = "Remove line breaks from text before embedding (true/false)"

	EmbedderFormHelpTitle   = "Embedding Configuration"
	EmbedderFormHelpContent = `Configure text vectorization for semantic search and knowledge storage.

If no specific embedding settings are configured, the system will use OpenAI embeddings with the API key from LLM Providers.

Change providers carefully - different embedders produce incompatible vectors requiring database reindexing.`

	EmbedderFormHelpOpenAI      = "OpenAI: Most reliable option with excellent quality. Requires API key from LLM Providers if not set here."
	EmbedderFormHelpOllama      = "Ollama: Local embeddings, no API key needed. Requires Ollama server running."
	EmbedderFormHelpHuggingFace = "HuggingFace: Open source models with API key required."
	EmbedderFormHelpGoogleAI    = "Google AI: Quality embeddings, requires API key."

	// Provider names and descriptions
	EmbedderProviderDefault         = "Default (OpenAI)"
	EmbedderProviderDefaultDesc     = "Use OpenAI embeddings with API key from LLM Providers configuration"
	EmbedderProviderOpenAI          = "OpenAI"
	EmbedderProviderOpenAIDesc      = "OpenAI text embeddings API (text-embedding-3-small, ada-002)"
	EmbedderProviderOllama          = "Ollama"
	EmbedderProviderOllamaDesc      = "Local Ollama server for open-source embedding models"
	EmbedderProviderMistral         = "Mistral"
	EmbedderProviderMistralDesc     = "Mistral AI embedding models"
	EmbedderProviderJina            = "Jina"
	EmbedderProviderJinaDesc        = "Jina AI embedding API"
	EmbedderProviderHuggingFace     = "HuggingFace"
	EmbedderProviderHuggingFaceDesc = "HuggingFace inference API for embedding models"
	EmbedderProviderGoogleAI        = "Google AI"
	EmbedderProviderGoogleAIDesc    = "Google AI embedding models (embedding-001)"
	EmbedderProviderVoyageAI        = "VoyageAI"
	EmbedderProviderVoyageAIDesc    = "VoyageAI embedding API"
	EmbedderProviderDisabled        = "Disabled"
	EmbedderProviderDisabledDesc    = "Disable embeddings functionality completely"

	// Provider-specific placeholders and help
	EmbedderURLPlaceholderOpenAI      = "https://api.openai.com/v1"
	EmbedderURLPlaceholderOllama      = "http://localhost:11434"
	EmbedderURLPlaceholderMistral     = "https://api.mistral.ai/v1"
	EmbedderURLPlaceholderJina        = "https://api.jina.ai/v1"
	EmbedderURLPlaceholderHuggingFace = "https://api-inference.huggingface.co"
	EmbedderURLPlaceholderGoogleAI    = "Not supported - uses default endpoint"
	EmbedderURLPlaceholderVoyageAI    = "Not supported - uses default endpoint"

	EmbedderAPIKeyPlaceholderOllama      = "Not required for local models"
	EmbedderAPIKeyPlaceholderMistral     = "Mistral API key"
	EmbedderAPIKeyPlaceholderJina        = "Jina API key"
	EmbedderAPIKeyPlaceholderHuggingFace = "HuggingFace API key"
	EmbedderAPIKeyPlaceholderGoogleAI    = "Google AI API key"
	EmbedderAPIKeyPlaceholderVoyageAI    = "VoyageAI API key"
	EmbedderAPIKeyPlaceholderDefault     = "API key for the provider"

	EmbedderModelPlaceholderOpenAI      = "text-embedding-3-small"
	EmbedderModelPlaceholderOllama      = "nomic-embed-text"
	EmbedderModelPlaceholderMistral     = "mistral-embed"
	EmbedderModelPlaceholderJina        = "jina-embeddings-v2-base-en"
	EmbedderModelPlaceholderHuggingFace = "sentence-transformers/all-MiniLM-L6-v2"
	EmbedderModelPlaceholderGoogleAI    = "gemini-embedding-001"
	EmbedderModelPlaceholderVoyageAI    = "voyage-2"
	EmbedderModelPlaceholderDefault     = "Model name"

	// Provider IDs for internal use
	EmbedderProviderIDDefault     = "default"
	EmbedderProviderIDOpenAI      = "openai"
	EmbedderProviderIDOllama      = "ollama"
	EmbedderProviderIDMistral     = "mistral"
	EmbedderProviderIDJina        = "jina"
	EmbedderProviderIDHuggingFace = "huggingface"
	EmbedderProviderIDGoogleAI    = "googleai"
	EmbedderProviderIDVoyageAI    = "voyageai"
	EmbedderProviderIDDisabled    = "none"

	EmbedderHelpGeneral = `Embeddings convert text into vectors for semantic search and knowledge storage. This enables PentAGI to understand meaning rather than just keywords, making search results more relevant and intelligent.

Key benefits:
• Find documents by meaning, not exact words
• Build a smart knowledge base from pentesting results
• Enable AI agents to locate relevant information quickly
• Support advanced reasoning with contextual data

Choose Ollama for completely local processing - your data never leaves your infrastructure. Other providers offer cloud-based processing with different model capabilities and pricing.

Configure carefully as changing providers requires rebuilding the entire knowledge base.`

	EmbedderHelpAttentionPrefix = "Important:"
	EmbedderHelpAttention       = `Different embedding providers create incompatible vectors. Changing providers or models will break existing semantic search.

You must flush or reindex your entire knowledge base using the etester utility:
• Run 'etester flush' to clear old embeddings
• Run 'etester reindex' to rebuild with new provider
• This process can take significant time for large datasets`

	EmbedderHelpAttentionSuffix = `Only change providers if absolutely necessary.`

	// Provider help texts
	EmbedderHelpDefault = `Default mode uses OpenAI embeddings with the API key configured in LLM Providers.

This is the recommended option for most users as it requires no additional configuration if you already have OpenAI set up.`

	EmbedderHelpOpenAI = `Direct OpenAI API access for embedding generation.

Get your API key from:
https://platform.openai.com/api-keys

Recommended models:
• text-embedding-3-small (cost-effective, 1536 dimensions)
• text-embedding-3-large (highest quality, 3072 dimensions)
• text-embedding-ada-002 (legacy, still supported)`

	EmbedderHelpOllama = `Local Ollama server for open-source embedding models.

Popular embedding models:
• nomic-embed-text (recommended, 768 dimensions)
• mxbai-embed-large (large model, 1024 dimensions)
• snowflake-arctic-embed (multilingual support)

Install Ollama from:
https://ollama.com/

Start with: ollama pull nomic-embed-text`

	EmbedderHelpMistral = `Mistral AI embedding models via API.

Get your API key from:
https://console.mistral.ai/

Uses Mistral's embedding model with fixed configuration.
No model selection required - uses the default embedding model.`

	EmbedderHelpJina = `Jina AI embedding API with specialized models.

Get your API key from:
https://jina.ai/

Recommended models:
• jina-embeddings-v2-base-en (general purpose, 768 dimensions)
• jina-embeddings-v2-small-en (lightweight, 512 dimensions)
• jina-embeddings-v2-base-code (code-specific embeddings)`

	EmbedderHelpHuggingFace = `HuggingFace Inference API for open-source models.

Get your API key from:
https://huggingface.co/settings/tokens

Popular models:
• sentence-transformers/all-MiniLM-L6-v2 (384 dimensions)
• sentence-transformers/all-mpnet-base-v2 (768 dimensions)
• intfloat/e5-large-v2 (1024 dimensions)`

	EmbedderHelpGoogleAI = `Google AI embedding models (Gemini).

Get your API key from:
https://aistudio.google.com/app/apikey

Available models:
• gemini-embedding-001 (latest model, 768 dimensions)
• text-embedding-004 (legacy Vertex AI model)

Uses Google's fixed endpoint - URL configuration not supported.`

	EmbedderHelpVoyageAI = `VoyageAI embedding API optimized for retrieval.

Get your API key from:
https://www.voyageai.com/

Recommended models:
• voyage-2 (general purpose, 1024 dimensions)
• voyage-large-2 (highest quality, 1536 dimensions)
• voyage-code-2 (code embeddings, 1536 dimensions)`

	EmbedderHelpDisabled = `Disables all embedding functionality.

This will:
• Disable semantic search capabilities
• Turn off knowledge storage vectorization
• Reduce memory and computational requirements

Only recommended if embeddings are not needed for your use case.`
)
