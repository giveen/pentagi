package locale

// LLM Providers Screen constants
const (
	LLMProvidersTitle       = "LLM Providers Configuration"
	LLMProvidersDescription = "Configure Large Language Model providers for AI agents"
	LLMProvidersName        = "LLM Providers"
	LLMProvidersOverview    = `PentAGI uses specialized AI agents (researcher, developer, executor, pentester) that require different LLM capabilities for optimal penetration testing results.

Why multiple providers matter:
• Agent Specialization: Different agents benefit from models optimized for reasoning, coding, or analysis
• Cost Efficiency: Mix expensive reasoning models (o3, grok-4, claude-sonnet-4, gemini-2.5-pro) for complex tasks with cheaper models for simple operations
• Performance Optimization: Each provider excels in different areas - OpenAI for medium tasks, Anthropic for complex tasks, Gemini for saving costs

Provider Selection Guide:
• Cloud Production: OpenAI + Anthropic + Gemini for industry-leading performance and reliability
• Enterprise/Compliance: AWS Bedrock for SOC2, HIPAA, and access to multiple model families
• Privacy/On-premises: Ollama or vLLM with Llama 3.1, Qwen3, or other open models for complete data control

Ready-to-use configurations for OpenRouter, DeepInfra, vLLM, Ollama, and other providers are available in the /opt/pentagi/conf/ directory inside the container`
)

// LLM Provider titles and descriptions
const (
	LLMProviderOpenAI        = "OpenAI"
	LLMProviderAnthropic     = "Anthropic"
	LLMProviderGemini        = "Google Gemini"
	LLMProviderBedrock       = "AWS Bedrock"
	LLMProviderOllama        = "Ollama"
	LLMProviderDeepSeek      = "DeepSeek"
	LLMProviderGLM           = "GLM Zhipu AI"
	LLMProviderKimi          = "Kimi Moonshot AI"
	LLMProviderQwen          = "Qwen Alibaba Cloud"
	LLMProviderCustom        = "Custom"
	LLMProviderOpenAIDesc    = "Industry-leading GPT models with excellent general performance"
	LLMProviderAnthropicDesc = "Claude models with superior reasoning and safety features"
	LLMProviderGeminiDesc    = "Google's advanced multimodal models with broad knowledge"
	LLMProviderBedrockDesc   = "Enterprise AWS access to multiple foundation model providers"
	LLMProviderOllamaDesc    = "Local and cloud open-source models for privacy and flexibility"
	LLMProviderDeepSeekDesc  = "Advanced Chinese AI models with strong reasoning and multilingual capabilities"
	LLMProviderGLMDesc       = "Zhipu AI's GLM models for Chinese and English tasks"
	LLMProviderKimiDesc      = "Moonshot AI's long-context models for document analysis"
	LLMProviderQwenDesc      = "Alibaba Cloud's Qwen models for multilingual tasks"
	LLMProviderCustomDesc    = "Custom OpenAI-compatible endpoint for maximum flexibility"
)

// Provider-specific help text
const (
	LLMFormOpenAIHelp = `OpenAI delivers industry-leading models with cutting-edge reasoning capabilities perfect for sophisticated penetration testing.

Default PentAGI Models:
• o1, o4-mini: Advanced reasoning models for complex vulnerability analysis and strategic planning
• GPT-4.1, GPT-4.1-mini: Flagship models optimized for exploit development and code generation
• Automatic model selection based on agent type and task complexity

Key Advantages:
• Most advanced reasoning capabilities with step-by-step analysis (o-series models)
• Excellent coding abilities for custom exploit development and payload generation
• Reliable performance with consistent uptime and extensive API documentation
• Proven track record in security research and penetration testing scenarios

Best for: Production environments requiring cutting-edge AI capabilities, teams prioritizing performance over cost
Cost: Premium pricing, but optimized configurations balance cost with quality

Setup: Get your API key from https://platform.openai.com/api-keys`

	LLMFormAnthropicHelp = `Anthropic Claude models excel in safety-conscious penetration testing with superior reasoning and analytical capabilities.

Default PentAGI Models:
• Claude Sonnet-4: Premium reasoning model for complex security analysis and strategic vulnerability assessment
• Claude 3.5 Haiku: High-speed model optimized for rapid information gathering and simple parsing tasks
• Balanced cost-performance ratio across all security testing scenarios

Key Advantages:
• Exceptional safety and ethics focus - reduces harmful output while maintaining security testing effectiveness
• Superior reasoning for methodical vulnerability analysis and systematic penetration testing approaches
• Large context windows ideal for analyzing extensive codebases and configuration files
• Excellent at understanding complex security contexts and regulatory compliance requirements

Best for: Security teams prioritizing responsible testing practices, compliance-focused environments, detailed analysis
Cost: Mid-range pricing with excellent value for reasoning-heavy security workflows

Setup: Get your API key from https://console.anthropic.com/`

	LLMFormGeminiHelp = `Google Gemini combines multimodal capabilities with advanced reasoning, perfect for comprehensive security assessments.

Default PentAGI Models:
• Gemini 2.5 Pro: Advanced reasoning model for deep vulnerability analysis and complex exploit development
• Gemini 2.5 Flash: High-performance model balancing speed and intelligence for most security testing tasks
• Gemini 2.0 Flash Lite: Cost-effective model for rapid scanning and information gathering operations
• Reasoning capabilities with step-by-step analysis for thorough penetration testing

Key Advantages:
• Multimodal support enables analysis of screenshots, network diagrams, and security documentation
• Competitive pricing with generous rate limits for development and testing environments
• Large context windows (up to 2M tokens) for analyzing massive codebases and system configurations
• Strong performance in code analysis and vulnerability identification across multiple programming languages

Best for: Budget-conscious teams, development environments, scenarios requiring image/document analysis
Cost: Most cost-effective option among major cloud providers with excellent performance/price ratio

Setup: Get your API key from https://aistudio.google.com/app/apikey`

	LLMFormBedrockHelp = `AWS Bedrock provides enterprise-grade access to 20+ foundation models with multiple authentication methods and enhanced security.

Default PentAGI Models:
• Claude Sonnet-4.5 (via Bedrock): Premium reasoning model with AWS enterprise security and extended thinking capabilities
• OpenAI GPT OSS 120B: Strong reasoning model for scientific analysis and complex security tasks
• Claude Haiku-4.5, DeepSeek V3.2, Qwen3-32B: Efficient models for specific agent roles and cost optimization
• Access to Amazon Nova (multimodal), Mistral, Moonshot, and more through single unified interface

Authentication Methods (priority order):
1. Default AWS Auth (BEDROCK_DEFAULT_AUTH=true): Use AWS SDK credential chain - recommended for EC2/ECS/Lambda
2. Bearer Token (BEDROCK_BEARER_TOKEN): Token-based authentication for custom auth scenarios
3. Static Credentials (ACCESS_KEY + SECRET_KEY): Traditional IAM credentials for development and testing

Key Advantages:
• Enterprise compliance: SOC2, HIPAA, FedRAMP certifications with data residency and governance controls
• Multi-provider access: 20+ models from Anthropic, Amazon, OpenAI, Qwen, DeepSeek, Cohere, Mistral, Moonshot
• Flexible authentication: Three methods to suit different deployment scenarios and security requirements
• Enhanced security: VPC integration, CloudTrail logging, IAM controls, private endpoints for complete isolation
• Regional deployment: Deploy in preferred AWS regions for latency optimization and data sovereignty

Best for: Enterprise environments, regulated industries, teams requiring compliance controls and flexible authentication
Cost: Competitive pricing with provisioned throughput options, but new accounts have restrictive rate limits (2-20 req/min)
Important: Request quota increases through AWS Service Quotas console for production penetration testing workflows

Setup: Choose authentication method and configure credentials. Verify rate limits at https://docs.aws.amazon.com/bedrock/`

	LLMFormOllamaHelp = `Ollama supports two deployment scenarios for complete flexibility.

Scenario 1: Local Ollama Server (Self-Hosted)
• Run Ollama on your own hardware (8GB+ RAM recommended, GPU optional but beneficial)
• Complete data privacy - all processing happens locally
• Zero ongoing costs - only infrastructure
• No API key needed - authentication handled by network access
• Setup: Install from https://ollama.ai/ and configure OLLAMA_SERVER_URL=http://ollama-server:11434

Scenario 2: Ollama Cloud (Managed Service)
• Cloud-hosted models without local infrastructure requirements
• No hardware needed - models run on Ollama's infrastructure
• Pay-per-use pricing with free tier available
• API key required - generate at https://ollama.com/settings/keys
• Setup: Register at https://ollama.com, configure OLLAMA_SERVER_URL=https://ollama.com + OLLAMA_SERVER_API_KEY=your_key

Default PentAGI Models:
• Llama 3.1:8b, Qwen3:32b, and other open models
• Customizable - switch between 100+ available models
• Model auto-download and loading options for convenience

Key Advantages:
• Dual deployment options: Choose between privacy (local) and convenience (cloud)
• Cost flexibility: Zero ongoing costs for local, pay-per-use for cloud
• Extensive model library: Access to latest open-source models (Llama, Qwen, Mistral, Gemma, and more)
• Air-gapped support: Local deployment works in isolated networks

Best for: Privacy-focused teams (local), budget-conscious deployments (cloud), organizations with data sovereignty requirements
Setup options: Local installation from https://10.10.10.10:11434 or cloud registration at https://ollama.com`

	LLMFormDeepSeekHelp = `DeepSeek provides advanced AI models with strong reasoning capabilities and multilingual support.

Default PentAGI Models:
• DeepSeek-Chat: Flagship model for general-purpose tasks with strong coding and reasoning capabilities
• DeepSeek-Reasoner: Advanced reasoning model for complex security analysis
• Cost-effective pricing with competitive performance compared to leading models

Key Advantages:
• Strong coding and reasoning capabilities for security analysis and exploit development
• Multilingual support (Chinese and English) for international penetration testing scenarios
• Competitive pricing with excellent performance-to-cost ratio
• OpenAI-compatible API for seamless integration

LiteLLM Integration:
• Set Provider Name to 'deepseek' when using LiteLLM proxy
• Enables model prefix (e.g., deepseek/deepseek-chat) without modifying config.yml
• Optional for direct DeepSeek API usage

Best for: Teams requiring multilingual support, cost-conscious deployments, Chinese language security testing
Cost: Highly competitive pricing with strong performance characteristics

Setup: Get your API key from https://platform.deepseek.com/`

	LLMFormGLMHelp = `GLM from Zhipu AI provides advanced language models with strong NLP and reasoning capabilities developed by Tsinghua University.

Default PentAGI Models:
• GLM-4-Air: High performance general dialogue model optimized for regular tasks and tool calling
• GLM-4-Plus: Flagship model with strong reasoning and code generation capabilities
• GLM-Z1-Plus: Advanced reasoning model with deep analysis capabilities for security research

Key Advantages:
• Exceptional Chinese and English NLP capabilities
• Strong performance in multilingual security testing and analysis scenarios
• GLM-4 and GLM-Z1 model families with enhanced reasoning and coding
• OpenAI-compatible API for easy integration

Alternative API Endpoints:
• International: https://api.z.ai/api/paas/v4 (default)
• China: https://open.bigmodel.cn/api/paas/v4
• Coding-specific: https://api.z.ai/api/coding/paas/v4

LiteLLM Integration:
• Set Provider Name to 'zai' when using LiteLLM proxy
• Enables model prefix (e.g., zai/glm-4) without modifying config.yml
• Optional for direct GLM API usage

Best for: Chinese and English multilingual penetration testing, teams operating in Asian markets
Cost: Competitive pricing with good performance for multilingual tasks

Setup: Get your API key from https://open.bigmodel.cn/`

	LLMFormKimiHelp = `Kimi from Moonshot AI provides ultra-long context models perfect for analyzing extensive codebases and documentation.

Default PentAGI Models:
• Moonshot-v1-8k: Long-context model supporting up to 8K tokens for general dialogue
• Kimi-k2.5: Advanced model with strong reasoning and document understanding
• Optimized for processing large volumes of text and code

Key Advantages:
• Ultra-long context windows (up to 1M tokens) for comprehensive codebase analysis
• Strong Chinese and English language support for multilingual penetration testing
• Cost-effective for document-heavy security assessments and threat intelligence analysis
• Excellent at understanding complex system architectures and long-form technical documentation

Alternative API Endpoints:
• International: https://api.moonshot.ai/v1 (default)
• China: https://api.moonshot.cn/v1

LiteLLM Integration:
• Set Provider Name to 'moonshot' when using LiteLLM proxy
• Enables model prefix (e.g., moonshot/kimi-k2.5) without modifying config.yml
• Optional for direct Kimi API usage

Best for: Large codebase analysis, document-heavy assessments, teams needing extended context for security research
Cost: Competitive pricing with excellent value for long-context use cases

Setup: Get your API key from https://platform.moonshot.ai/`

	LLMFormQwenHelp = `Qwen from Alibaba Cloud Model Studio (DashScope) provides powerful multilingual models with multimodal capabilities.

Default PentAGI Models:
• Qwen-Turbo: Fastest lightweight model for high-frequency tasks and real-time response scenarios
• Qwen-Plus: Balanced performance model for general dialogue, code generation, and tool calling
• Qwen-Max: Flagship reasoning model with strong instruction following and complex task handling
• QwQ-Plus: Deep reasoning model with extended chain-of-thought for complex logic analysis

Key Advantages:
• Strong multilingual support (Chinese, English, and multiple other languages)
• Multimodal capabilities with Qwen-VL for visual security analysis
• Alibaba Cloud integration for enterprise deployments
• DashScope ecosystem with additional AI services and tools
• Qwen2.5, Qwen3, and QwQ model families with various sizes and specializations

Alternative API Endpoints:
• US: https://dashscope-us.aliyuncs.com/compatible-mode/v1 (default)
• Singapore: https://dashscope-intl.aliyuncs.com/compatible-mode/v1
• China: https://dashscope.aliyuncs.com/compatible-mode/v1

LiteLLM Integration:
• Set Provider Name to 'dashscope' when using LiteLLM proxy
• Enables model prefix (e.g., dashscope/qwen-plus) without modifying config.yml
• Optional for direct Qwen API usage

Best for: Teams operating in Asian markets, multilingual security testing, visual analysis with Qwen-VL, Alibaba Cloud ecosystem integration
Cost: Competitive pricing with flexible tiers for different use cases

Setup: Get your API key from https://dashscope.console.aliyun.com/`

	LLMFormCustomHelp = `Configure any OpenAI-compatible API endpoint for maximum flexibility and integration with existing infrastructure.

Ready-to-use Configurations:
• vLLM deployments: High-throughput on-premises inference with optimal GPU utilization
• OpenRouter: Access 200+ models from multiple providers through single API with competitive pricing
• DeepInfra: Serverless inference for popular open models with pay-per-use pricing
• Together AI, Groq, Fireworks: Alternative cloud providers with specialized performance optimizations
• LiteLLM Proxy: Universal gateway to 100+ providers with load balancing and unified interface (use LLM_SERVER_PROVIDER for model prefixing)
• Some reasoning models and LLM providers may require preserving reasoning content while using tool calls (LLM_SERVER_PRESERVE_REASONING=true)

Popular On-Premises Options:
• vLLM: Production-grade serving for Qwen, Llama, Mistral models with batching and GPU optimization
• LocalAI: OpenAI-compatible API wrapper for various local models and embedding services
• Text Generation WebUI: Community-favorite interface with extensive model support and fine-tuning capabilities
• Hugging Face TGI: Enterprise text generation inference with auto-scaling and monitoring

Key Advantages:
• Unlimited flexibility: Use any OpenAI-compatible endpoint or service
• Cost optimization: Choose providers with competitive pricing or deploy models on your own infrastructure
• Vendor independence: Avoid lock-in with ability to switch between providers and models seamlessly
• Custom fine-tuning: Deploy specialized models trained on your security testing scenarios

Best for: Teams with specific model requirements, cost optimization needs, or existing LLM infrastructure
LiteLLM Integration: Set LLM_SERVER_PROVIDER to match your provider name (e.g., "openrouter", "moonshot") to use the same config files with both direct API access and LiteLLM proxy
Examples available: Pre-configured setups for major providers in /opt/pentagi/conf/ directory inside the container`
)

// LLM Provider Form field labels and descriptions
const (
	LLMFormFieldBaseURL           = "Base URL"
	LLMFormFieldAPIKey            = "API Key"
	LLMFormFieldDefaultAuth       = "Use Default AWS Auth"
	LLMFormFieldBearerToken       = "Bearer Token"
	LLMFormFieldAccessKey         = "Access Key ID"
	LLMFormFieldSecretKey         = "Secret Access Key"
	LLMFormFieldSessionToken      = "Session Token"
	LLMFormFieldRegion            = "Region"
	LLMFormFieldModel             = "Model"
	LLMFormFieldConfigPath        = "Config Path"
	LLMFormFieldLegacyReasoning   = "Legacy Reasoning"
	LLMFormFieldPreserveReasoning = "Preserve Reasoning"
	LLMFormFieldProviderName      = "Provider Name"
	LLMFormFieldPullTimeout       = "Model Pull Timeout"
	LLMFormFieldPullEnabled       = "Auto-pull Models"
	LLMFormFieldLoadModelsEnabled = "Load Models from Server"
	LLMFormBaseURLDesc            = "API endpoint URL for the provider"
	LLMFormAPIKeyDesc             = "Your API key for authentication"
	LLMFormDefaultAuthDesc        = "Use AWS SDK default credential chain (environment, EC2 role, ~/.aws/credentials) - highest priority"
	LLMFormBearerTokenDesc        = "Bearer token for authentication - takes priority over static credentials"
	LLMFormAccessKeyDesc          = "AWS Access Key ID for static credentials authentication"
	LLMFormSecretKeyDesc          = "AWS Secret Access Key for static credentials authentication"
	LLMFormSessionTokenDesc       = "AWS Session Token for temporary credentials (optional, used with static credentials)"
	LLMFormRegionDesc             = "AWS region for Bedrock service"
	LLMFormModelDesc              = "Default model to use for this provider"
	LLMFormConfigPathDesc         = "Path to configuration file (optional)"
	LLMFormLegacyReasoningDesc    = "Enable legacy reasoning mode (true/false)"
	LLMFormPreserveReasoningDesc  = "Preserve reasoning content in multi-turn conversations (required by some providers)"
	LLMFormProviderNameDesc       = "Provider name prefix for model names (useful for LiteLLM proxy)"
	LLMFormPullTimeoutDesc        = "Timeout in seconds for downloading models (default: 600)"
	LLMFormPullEnabledDesc        = "Automatically download required models on startup"
	LLMFormLoadModelsEnabledDesc  = "Load available models list from Ollama server"
	LLMFormOllamaAPIKeyDesc       = "Ollama Cloud API key (optional, leave empty for local Ollama server)"
)

// LLM Provider Form status messages
const (
	LLMProviderFormTitle       = "LLM Provider %s Configuration"
	LLMProviderFormDescription = "Configure your Large Language Model provider settings"
	LLMProviderFormName        = "LLM Provider %s"
	LLMProviderFormOverview    = `Agent Role Assignment:
• Primary Agent & Pentester: Use reasoning models (o3, grok-4, claude-sonnet-4, gemini-2.5-pro) for complex vulnerability analysis
• Assistant & Adviser: Advanced models (o4-mini, claude-sonnet-4) for strategic planning and recommendations
• Coder & Installer: Precision models (gpt-4.1, claude-sonnet-4) for exploit development and system configuration
• Searcher & Enricher: Fast models (gpt-4.1-mini, claude-3.5-haiku, gemini-2.0-flash-lite) for information gathering
• Simple tasks: Lightweight models for JSON parsing and basic operations

Performance Considerations:
• Reasoning models provide step-by-step analysis but are slower and more expensive
• Standard models offer faster responses suitable for high-frequency agent interactions
• Each agent type uses provider-specific model configurations optimized for security testing workflows

Your configuration will determine which models each agent uses for different penetration testing scenarios.`
)

// Monitoring Screen
const (
	MonitoringTitle       = "Monitoring Configuration"
	MonitoringDescription = "Configure monitoring and observability platforms for comprehensive system insights"
	MonitoringName        = "Monitoring"
	MonitoringOverview    = `Comprehensive monitoring and observability for production-ready deployments.

Why monitoring matters:
• Track performance bottlenecks: Identify slow LLM calls, database queries, and system resources
• Debug issues faster: Detailed traces help diagnose problems across distributed components
• Optimize costs: Monitor token usage patterns and optimize expensive LLM interactions
• Production readiness: Essential for reliable operation in critical environments

Platform Options:
Langfuse: Specialized LLM observability with conversation tracking, prompt engineering insights, and cost analytics
Observability: Full-stack monitoring with metrics, traces, logs, and alerting for infrastructure and application health

Quick Setup:
• Development: Enable Langfuse for LLM insights only
• Production: Enable both platforms for comprehensive monitoring
• Cost-conscious: Use embedded modes to avoid external service fees`
)

// Langfuse Integration constants
const (
	MonitoringLangfuseFormTitle       = "Langfuse Configuration"
	MonitoringLangfuseFormDescription = "Configuration of Langfuse integration for LLM monitoring"
	MonitoringLangfuseFormName        = "Langfuse"
	MonitoringLangfuseFormOverview    = `Langfuse provides:
• Complete conversation tracking
• Model performance metrics
• Cost monitoring and optimization
• User behavior analytics
• Debug traces for AI interactions

Choose between embedded instance or external connection.`

	// Deployment types
	MonitoringLangfuseEmbedded = "Embedded Server"
	MonitoringLangfuseExternal = "External Server"
	MonitoringLangfuseDisabled = "Disabled"

	// Form fields
	MonitoringLangfuseDeploymentType     = "Deployment Type"
	MonitoringLangfuseDeploymentTypeDesc = "Select the deployment type for Langfuse"
	MonitoringLangfuseBaseURL            = "Server URL"
	MonitoringLangfuseBaseURLDesc        = "Address of the Langfuse server (e.g., https://cloud.langfuse.com)"
	MonitoringLangfuseProjectID          = "Project ID"
	MonitoringLangfuseProjectIDDesc      = "Project identifier in Langfuse"
	MonitoringLangfusePublicKey          = "Public Key"
	MonitoringLangfusePublicKeyDesc      = "Public API key for project access"
	MonitoringLangfuseSecretKey          = "Secret Key"
	MonitoringLangfuseSecretKeyDesc      = "Secret API key for project access"
	MonitoringLangfuseListenIP           = "Listen IP"
	MonitoringLangfuseListenIPDesc       = "Bind address used by Docker port mapping (e.g., 0.0.0.0 to expose on all interfaces)"
	MonitoringLangfuseListenPort         = "Listen Port"
	MonitoringLangfuseListenPortDesc     = "External TCP port exposed by Docker for Langfuse web UI"

	// Admin settings for embedded
	MonitoringLangfuseAdminEmail        = "Admin Email"
	MonitoringLangfuseAdminEmailDesc    = "Email for accessing the Langfuse admin panel"
	MonitoringLangfuseAdminPassword     = "Admin Password"
	MonitoringLangfuseAdminPasswordDesc = "Password for accessing the Langfuse admin panel"
	MonitoringLangfuseAdminName         = "Admin Username"
	MonitoringLangfuseAdminNameDesc     = "Administrator username in Langfuse"
	MonitoringLangfuseLicenseKey        = "Enterprise License Key"
	MonitoringLangfuseLicenseKeyDesc    = "Langfuse Enterprise license key (optional)"

	// Help text
	MonitoringLangfuseModeGuide    = "Choose deployment: Embedded (local control), External (cloud/existing), Disabled (no analytics)"
	MonitoringLangfuseEmbeddedHelp = `Embedded deploys complete Langfuse stack:
• PostgreSQL + ClickHouse databases
• MinIO S3 storage + Redis cache
• Full LLM conversation tracking
• Cost analysis and performance metrics
• Private data stays on your server

Resource requirements:
• ~2GB RAM, 5GB disk space minimum
• Additional storage for conversation logs
• Automatic setup and maintenance

Best for: Teams wanting data privacy, custom configurations, or no external dependencies. All analytics data stored locally with full administrative control.

Default admin access:
• Web UI: http://localhost:4000
• Login: admin@pentagi.com
• Password: password (change required)`
	MonitoringLangfuseExternalHelp = `External connects to cloud.langfuse.com or your existing Langfuse server:

• No local infrastructure needed
• Managed updates and maintenance
• Shared analytics across teams
• Enterprise features available
• Data stored on external provider

Setup requirements:
• Langfuse account and API keys
• Internet connectivity required
• Project ID and authentication keys

Best for: Teams using cloud services, wanting managed infrastructure, or integrating with existing Langfuse deployments across organizations.`
	MonitoringLangfuseDisabledHelp = `Langfuse is disabled. Without LLM observability you will not have:

• Conversation history tracking
• Token usage and cost analysis
• Model performance metrics
• Debug traces for AI interactions
• User behavior analytics
• Prompt engineering insights

Consider enabling for production use
to monitor AI agent performance and
optimize costs effectively.`
)

// Graphiti Integration constants
const (
	MonitoringGraphitiFormTitle       = "Graphiti Configuration (beta)"
	MonitoringGraphitiFormDescription = "Configuration of Graphiti knowledge graph integration"
	MonitoringGraphitiFormName        = "Graphiti (beta)"
	MonitoringGraphitiFormOverview    = `⚠️  BETA FEATURE: This functionality is currently under active development. Please monitor updates for improvements and stability fixes.

Graphiti provides temporal knowledge graph capabilities:
• Entity and relationship extraction
• Semantic memory for AI agents
• Temporal context tracking
• Knowledge reuse across flows

⚠️  REQUIREMENT: Graphiti requires configured OpenAI provider (LLM Providers → OpenAI) for entity extraction.

Choose between embedded instance or external connection.`

	// Deployment types
	MonitoringGraphitiEmbedded = "Embedded Stack"
	MonitoringGraphitiExternal = "External Service"
	MonitoringGraphitiDisabled = "Disabled"

	// Form fields
	MonitoringGraphitiDeploymentType     = "Deployment Type"
	MonitoringGraphitiDeploymentTypeDesc = "Select the deployment type for Graphiti"
	MonitoringGraphitiURL                = "Graphiti Server URL"
	MonitoringGraphitiURLDesc            = "Address of the Graphiti API server"
	MonitoringGraphitiTimeout            = "Request Timeout"
	MonitoringGraphitiTimeoutDesc        = "Timeout in seconds for Graphiti operations"
	MonitoringGraphitiModelName          = "Extraction Model"
	MonitoringGraphitiModelNameDesc      = "LLM model for entity extraction (uses OpenAI provider from LLM Providers configuration)"
	MonitoringGraphitiNeo4jUser          = "Neo4j Username"
	MonitoringGraphitiNeo4jUserDesc      = "Username for Neo4j database access"
	MonitoringGraphitiNeo4jPassword      = "Neo4j Password"
	MonitoringGraphitiNeo4jPasswordDesc  = "Password for Neo4j database access"
	MonitoringGraphitiNeo4jDatabase      = "Neo4j Database"
	MonitoringGraphitiNeo4jDatabaseDesc  = "Neo4j database name"

	// Help text
	MonitoringGraphitiModeGuide    = "Choose deployment: Embedded (local Neo4j), External (existing Graphiti), Disabled (no knowledge graph)"
	MonitoringGraphitiEmbeddedHelp = `⚠️  BETA: This feature is under active development. Monitor updates for improvements.

Embedded deploys complete Graphiti stack:
• Neo4j graph database
• Graphiti API service
• Automatic entity extraction from agent interactions
• Temporal relationship tracking
• Private knowledge graph on your server

Prerequisites:
• OpenAI provider must be configured (LLM Providers → OpenAI)
• OpenAI API key is used for entity extraction
• Configured model will be used for knowledge graph operations

Resource requirements:
• ~1.5GB RAM, 3GB disk space minimum
• Neo4j UI: http://localhost:7474
• Graphiti API: http://localhost:8000
• Automatic setup and maintenance

Best for: Teams wanting knowledge graph capabilities with full data control and privacy.`
	MonitoringGraphitiExternalHelp = `⚠️  BETA: This feature is under active development. Monitor updates for improvements.

External connects to your existing Graphiti server:

• No local infrastructure needed
• Managed updates and maintenance
• Shared knowledge graph across teams
• Data stored on external provider

Setup requirements:
• Graphiti server URL and access
• Network connectivity required
• External server must be configured with OpenAI API key
• Model and extraction settings configured on external server

Best for: Teams using existing Graphiti deployments or cloud services.`
	MonitoringGraphitiDisabledHelp = `Graphiti is disabled. You will not have:

• Temporal knowledge graph
• Entity and relationship extraction
• Semantic memory for AI agents
• Knowledge reuse across flows
• Advanced contextual search

Note: Graphiti is currently in beta.
Consider enabling for production use
to build a knowledge base from
penetration testing results.`
)

// Observability Integration constants
const (
	MonitoringObservabilityFormTitle       = "Observability Configuration"
	MonitoringObservabilityFormDescription = "Configuration of monitoring and observability stack"
	MonitoringObservabilityFormName        = "Observability"
	MonitoringObservabilityFormOverview    = `Observability stack includes:
• Grafana dashboards for visualization
• VictoriaMetrics for time-series data
• Jaeger for distributed tracing
• Loki for log aggregation
• OpenTelemetry for data collection

Monitor PentAGI performance and system health.`

	// Deployment types
	MonitoringObservabilityEmbedded = "Embedded Stack"
	MonitoringObservabilityExternal = "External Collector"
	MonitoringObservabilityDisabled = "Disabled"

	// Form fields
	MonitoringObservabilityDeploymentType     = "Deployment Type"
	MonitoringObservabilityDeploymentTypeDesc = "Select the deployment type for monitoring"
	MonitoringObservabilityOTelHost           = "OpenTelemetry Host"
	MonitoringObservabilityOTelHostDesc       = "Address of the external OpenTelemetry collector"

	// embedded listen fields
	MonitoringObservabilityGrafanaListenIP        = "Grafana Listen IP"
	MonitoringObservabilityGrafanaListenIPDesc    = "Bind address used by Docker port mapping (e.g., 0.0.0.0 to expose on all interfaces)"
	MonitoringObservabilityGrafanaListenPort      = "Grafana Listen Port"
	MonitoringObservabilityGrafanaListenPortDesc  = "External TCP port exposed by Docker for Grafana web UI"
	MonitoringObservabilityOTelGrpcListenIP       = "OTel gRPC Listen IP"
	MonitoringObservabilityOTelGrpcListenIPDesc   = "Bind address used by Docker port mapping (e.g., 0.0.0.0 to expose on all interfaces)"
	MonitoringObservabilityOTelGrpcListenPort     = "OTel gRPC Listen Port"
	MonitoringObservabilityOTelGrpcListenPortDesc = "External TCP port exposed by Docker for OTel gRPC receiver"
	MonitoringObservabilityOTelHttpListenIP       = "OTel HTTP Listen IP"
	MonitoringObservabilityOTelHttpListenIPDesc   = "Bind address used by Docker port mapping (e.g., 0.0.0.0 to expose on all interfaces)"
	MonitoringObservabilityOTelHttpListenPort     = "OTel HTTP Listen Port"
	MonitoringObservabilityOTelHttpListenPortDesc = "External TCP port exposed by Docker for OTel HTTP receiver"

	// Help text
	MonitoringObservabilityModeGuide    = "Choose monitoring: Embedded (full stack), External (existing infra), Disabled (no monitoring)"
	MonitoringObservabilityEmbeddedHelp = `Embedded deploys complete monitoring:
• Grafana dashboards and alerting
• VictoriaMetrics time-series database
• Jaeger distributed tracing UI
• Loki log aggregation system
• ClickHouse analytical database
• Node Exporter + cAdvisor metrics
• OpenTelemetry data collection

Auto-instrumented components with
pre-built dashboards for system health,
performance analysis, and debugging.

Resource requirements:
• ~1.5GB RAM, 3GB disk space minimum
• Grafana UI: http://localhost:3000
• Profiling: http://localhost:7777

Best for: Complete system visibility,
troubleshooting, and performance tuning.`
	MonitoringObservabilityExternalHelp = `External sends telemetry to your existing monitoring infrastructure:

• OTLP protocol over HTTP/2 (no TLS)
• Your collector must support:
  - OTLP HTTP receiver (port 4318)
  - OTLP gRPC receiver (port 8148)
  - tls: insecure: true setting
• Sends metrics, traces, and logs
• Compatible with enterprise platforms:
  Datadog, New Relic, Splunk, etc.

OTEL_HOST example:
your-collector:4318

Collector config requirement:
tls: insecure: true

Best for: Organizations with existing
monitoring infrastructure or centralized
observability platforms.`
	MonitoringObservabilityDisabledHelp = `Observability is disabled. You will not have:

• System performance monitoring
• Distributed request tracing
• Structured log aggregation
• Resource usage analytics
• Error tracking and alerting
• Performance bottleneck analysis

Consider enabling for production use
to monitor system health, debug issues,
and optimize performance effectively.`
)

// Summarizer Screen
const (
	SummarizerTitle       = "Summarizer Configuration"
	SummarizerDescription = "Enable conversation summarization to reduce LLM costs and improve context management"
	SummarizerName        = "Summarizer"
	SummarizerOverview    = `Optimize context usage, reduce LLM costs, and match your model capabilities.

When to adjust summarization:
• High token costs: Reduce context size (4K-12K vs 22K+ tokens)
• "Context too long" errors: Configure for your model's limits
• Poor conversation flow: Increase context retention for quality
• Different model types: Short-context vs long-context model tuning

General Summarization: Maximum cost control and precision tuning for research/analysis tasks
Assistant Summarization: Optimal conversation quality with intelligent context management for interactive sessions

Quick wins:
• Cost reduction: Use General, reduce Recent Sections to 1-2
• Context errors: Match limits to your model (8K/32K/128K)
• Quality priority: Use Assistant with increased limits`

	SummarizerTypeGeneralName = "General Summarization"
	SummarizerTypeGeneralDesc = "Global summarization settings for conversation context management"

	SummarizerTypeGeneralInfo = `Choose this for maximum cost control and short-context model compatibility.

Perfect when you need:
• Aggressive cost reduction: Fine-tune every parameter for minimal token usage
• Short-context models (8K-32K): Precise limits to avoid overflow errors
• Research/analysis tasks: Controlled compression without losing key data
• Custom QA handling: Full control over question-answer pair processing

Typical results:
• 40-70% cost reduction vs default settings
• 4K-12K token contexts (vs 22K+ in Assistant mode)
• Better performance on GPT-3.5, Claude Instant, smaller models
• Precise control over conversation memory vs fresh context balance

Best practices:
• Start with 1-2 Recent Sections for maximum savings
• Enable Size Management for automatic overflow protection
• Disable QA compression only for critical reasoning tasks`

	SummarizerTypeAssistantName = "Assistant Summarization"
	SummarizerTypeAssistantDesc = "Specialized summarization settings for AI assistant contexts"

	SummarizerTypeAssistantInfo = `Choose this for optimal conversation quality and dialogue continuity.

Perfect when you need:
• Extended reasoning chains: Maintain context for complex multi-step thinking
• High-quality conversations: Preserve dialogue flow and assistant personality
• Long-context models (64K+): Leverage full model capabilities efficiently
• Interactive sessions: Better memory of user preferences and conversation history

Typical results:
• 8K-40K token contexts with intelligent compression
• Superior conversation continuity vs manual settings
• Automatic context optimization for reasoning tasks
• Balanced cost vs quality (3x more context than General mode)

Best practices:
• Use default settings for most scenarios - they're pre-optimized
• Increase Recent Sections only for very complex tasks
• Monitor context usage - costs scale with token count
• Perfect for GPT-4, Claude, and other large context models`
)

// Summarizer Form Screen
const (
	SummarizerFormGeneralTitle   = "General Summarizer Configuration"
	SummarizerFormAssistantTitle = "Assistant Summarizer Configuration"
	SummarizerFormDescription    = "Configure %s Settings"

	// Field Labels and Descriptions
	SummarizerFormPreserveLast     = "Size Management"
	SummarizerFormPreserveLastDesc = "Controls last section compression. Enabled: sections fit LastSecBytes (smaller context). Disabled: sections grow freely (larger context)"

	SummarizerFormUseQA     = "QA Summarization"
	SummarizerFormUseQADesc = "Enables question-answer pair compression when total QA content exceeds MaxQABytes or MaxQASections limits"

	SummarizerFormSumHumanInQA     = "Compress User Messages"
	SummarizerFormSumHumanInQADesc = "Include user messages in QA compression. Disabled: preserves original user text (recommended for most cases)"

	SummarizerFormLastSecBytes     = "Section Size Limit"
	SummarizerFormLastSecBytesDesc = "Maximum bytes per recent section when Size Management enabled. Larger: more detail per section, higher token usage"

	SummarizerFormMaxBPBytes     = "Response Size Limit"
	SummarizerFormMaxBPBytesDesc = "Maximum bytes for individual AI responses before compression. Prevents single large responses from dominating context"

	SummarizerFormMaxQASections     = "QA Section Limit"
	SummarizerFormMaxQASectionsDesc = "Maximum question-answer sections before QA compression triggers. Works with MaxQABytes to control total QA memory"

	SummarizerFormMaxQABytes     = "Total QA Memory"
	SummarizerFormMaxQABytesDesc = "Maximum bytes for all QA sections combined. When exceeded (with MaxQASections), triggers QA compression to fit limit"

	SummarizerFormKeepQASections     = "Recent Sections"
	SummarizerFormKeepQASectionsDesc = "Number of most recent conversation sections preserved without compression. PRIMARY parameter affecting context size"

	// Enhanced Help Text - General (common principles)
	SummarizerFormGeneralHelp = `Context estimation: 4K-22K tokens (typical), up to 94K (maximum settings).

Key relationships:
• Recent Sections: Most critical - each +1 adds ~1.5-9K tokens
• Size Management OFF: 2-3x larger context (less compression)
• Section/Response Limits: Control individual component sizes
• QA Memory: Manages total conversation history when limits exceeded

Parameter interactions:
• QA compression activates when BOTH MaxQABytes AND MaxQASections exceeded
• Size Management disabled → sections can grow 2x larger than limits
• Response Limit prevents single large outputs from dominating context
• User message compression (SummHumanInQA) saves 5% but loses original phrasing

Reduce for smaller models:
• Recent Sections: 1-2 (vs 3+ default)
• Section Limit: 25-35KB (vs 50KB+)
• Disable Size Management for simple conversations

Common mistakes:
• Setting Recent Sections too high (main cause of context overflow)
• Enabling Size Management with very low Section Limits (over-compression)
• Mismatched QA limits (high bytes + low sections = ineffective)

Current algorithm compresses older content while preserving recent context quality.`

	// Enhanced Help Text - Assistant specific (interactive conversations)
	SummarizerFormAssistantHelp = `Optimized for interactive conversations requiring context continuity.

Default tuning (3 Recent Sections, 75KB limits):
• Typical range: 8K-40K tokens
• Good for: Extended dialogues, reasoning chains, context-dependent tasks
• Models: Works well with 32K+ context models

Adjustments by model type:
• Short context (≤16K): Recent Sections=1-2, Section Limit=45KB
• Long context (128K+): Can increase Recent Sections=5-7
• High-frequency chat: Reduce Recent Sections=2 for faster responses

Advanced tuning:
• QA Memory 200KB+ for document analysis conversations
• Response Limit 24-32KB for detailed technical responses
• Keep User Messages uncompressed (SummHumanInQA=false) for better context

Performance optimization:
• Each Recent Section ≈ 9-18KB in assistant mode
• Size Management reduces growth by ~20% but may lose detail
• QA compression triggers less often due to larger default limits

Size Management enabled by default - maintains conversation flow while preventing context overflow.
Monitor actual token usage and adjust Recent Sections first, then limits.`

	// Context size estimation
	SummarizerContextEstimatedSize    = "Estimated context size: %s\n%s"
	SummarizerContextTokenRange       = "~%s tokens"
	SummarizerContextTokenRangeMinMax = "~%s-%s tokens"
	SummarizerContextRequires256K     = "Requires 256K+ context model"
	SummarizerContextRequires128K     = "Requires 128K+ context model"
	SummarizerContextRequires64K      = "Requires 64K+ context model"
	SummarizerContextRequires32K      = "Requires 32K+ context model"
	SummarizerContextRequires16K      = "Requires 16K+ context model"
	SummarizerContextFitsIn8K         = "Fits in 8K+ context model"
)
