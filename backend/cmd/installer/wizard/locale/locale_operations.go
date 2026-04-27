package locale

// Development and Mock Screen constants
const (
	MockScreenTitle       = "Development Screen"
	MockScreenDescription = "This screen is under development"
)

// Apply Changes screen constants
const (
	ApplyChangesFormTitle       = "Apply Configuration Changes"
	ApplyChangesFormName        = "Apply Changes"
	ApplyChangesFormDescription = "Review and apply your configuration changes"

	// Apply Changes overview and help
	ApplyChangesFormOverview = `This screen allows you to review all pending configuration changes and apply them to your PentAGI installation.

When you apply changes, the system will:
• Save all modified environment variables to the .env file
• Restart affected services with the new configuration
• Install additional components if needed`

	// Apply Changes status messages
	ApplyChangesNotStarted     = "Configuration changes are ready to be applied"
	ApplyChangesInProgress     = "Applying configuration changes...\n"
	ApplyChangesCompleted      = "Configuration changes have been successfully applied\n"
	ApplyChangesFailed         = "Failed to perform configuration changes"
	ApplyChangesResetCompleted = "Configuration changes have been successfully reset\n"

	ApplyChangesTerminalIsNotInitialized = "Terminal is not initialized"

	// Apply Changes instructions
	ApplyChangesInstructions = `Press Enter to begin applying the configuration changes.`

	ApplyChangesNoChanges = "No configuration changes are pending"

	// Apply Changes installation status
	ApplyChangesInstallNotFound = `PentAGI is not currently installed on this system.

The following actions will be performed:
• Docker environment setup and validation
• Creation of docker-compose.yml file
• Installation and startup of PentAGI core services`

	ApplyChangesInstallFoundLangfuse      = `• Installation of Langfuse observability stack (docker-compose-langfuse.yml)`
	ApplyChangesInstallFoundObservability = `• Installation of comprehensive observability stack with Grafana, VictoriaMetrics, and Jaeger (docker-compose-observability.yml)`

	ApplyChangesUpdateFound = `PentAGI is currently installed on this system.

The following actions will be performed:
• Update environment variables in .env file
• Recreate and restart affected Docker containers
• Apply new configuration to running services`

	// Apply Changes warnings and notes
	ApplyChangesWarningCritical = "⚠️  Critical changes detected - services will be restarted"
	ApplyChangesWarningSecrets  = "🔒 Secret values detected - they will be securely stored"
	ApplyChangesNoteBackup      = "💾 Current configuration will be backed up before changes"
	ApplyChangesNoteTime        = "⏱️  This process may take less than a minute depending on selected components"

	// Apply Changes progress messages
	ApplyChangesStageValidation = "Validating environment and dependencies..."
	ApplyChangesStageBackup     = "Creating configuration backup..."
	ApplyChangesStageEnvFile    = "Updating environment file..."
	ApplyChangesStageCompose    = "Generating Docker Compose files..."
	ApplyChangesStageDocker     = "Managing Docker containers..."
	ApplyChangesStageServices   = "Starting services..."
	ApplyChangesStageComplete   = "Configuration changes applied successfully"

	// Apply Changes change list headers
	ApplyChangesChangesTitle  = "Pending Configuration Changes"
	ApplyChangesChangesCount  = "Total changes: %d"
	ApplyChangesChangesMasked = "(hidden for security)"
	ApplyChangesChangesEmpty  = "No changes to apply"

	// Apply Changes help content
	ApplyChangesHelpTitle   = "Applying Configuration Changes"
	ApplyChangesHelpContent = `Be sure to check the current configuration before applying changes.`
)

// apply changes integrity prompt
const (
	ApplyChangesIntegrityPromptTitle   = "File integrity check"
	ApplyChangesIntegrityPromptMessage = "Out-of-date files were detected.\nDo you want to update them to the latest version?"
	ApplyChangesIntegrityOutdatedList  = "Out-of-date files:\n%s\nConfirm update? (y/n)"
	ApplyChangesIntegrityChecking      = "Collecting file integrity information..."
	ApplyChangesIntegrityNoOutdated    = "No out-of-date files found. Proceeding with apply."
)

// Maintenance Screen constants
const (
	MaintenanceTitle       = "System Maintenance"
	MaintenanceDescription = "Manage PentAGI services and perform maintenance operations"
	MaintenanceName        = "Maintenance"
	MaintenanceOverview    = `Perform system maintenance operations for PentAGI.

Available operations depend on the current system state and will only be shown when applicable.

Operations include:
• Service lifecycle management (Start/Stop/Restart)
• Component updates and downloads
• System reset and cleanup
• Container and image management

Each operation will provide real-time status updates and confirmation when required.`

	// Maintenance menu items
	MaintenanceStartPentagi            = "Start PentAGI"
	MaintenanceStartPentagiDesc        = "Start all configured PentAGI services"
	MaintenanceStopPentagi             = "Stop PentAGI"
	MaintenanceStopPentagiDesc         = "Stop all running PentAGI services"
	MaintenanceRestartPentagi          = "Restart PentAGI"
	MaintenanceRestartPentagiDesc      = "Restart all PentAGI services"
	MaintenanceDownloadWorkerImage     = "Download Worker Image"
	MaintenanceDownloadWorkerImageDesc = "Download pentesting container image for worker tasks"
	MaintenanceUpdateWorkerImage       = "Update Worker Image"
	MaintenanceUpdateWorkerImageDesc   = "Update pentesting container image to latest version"
	MaintenanceUpdatePentagi           = "Update PentAGI"
	MaintenanceUpdatePentagiDesc       = "Update PentAGI to the latest version"
	MaintenanceUpdateInstaller         = "Update Installer"
	MaintenanceUpdateInstallerDesc     = "Update this installer to the latest version"
	MaintenanceFactoryReset            = "Factory Reset"
	MaintenanceFactoryResetDesc        = "Reset PentAGI to factory defaults"
	MaintenanceRemovePentagi           = "Remove PentAGI"
	MaintenanceRemovePentagiDesc       = "Remove PentAGI containers but keep data"
	MaintenancePurgePentagi            = "Purge PentAGI"
	MaintenancePurgePentagiDesc        = "Completely remove PentAGI including all data"
	MaintenanceResetPassword           = "Reset Admin Password"
	MaintenanceResetPasswordDesc       = "Reset the administrator password for PentAGI"
)

// Reset Password Screen constants
const (
	ResetPasswordFormTitle       = "Reset Admin Password"
	ResetPasswordFormDescription = "Reset the administrator password for PentAGI"
	ResetPasswordFormName        = "Reset Password"
	ResetPasswordFormOverview    = `Reset the password for the default administrator account (admin@pentagi.com).

This operation requires PentAGI to be running and will update the password in the PostgreSQL database.

Enter your new password twice to confirm and press Enter to apply the change.

Password requirements:
• Minimum 5 characters
• Both password fields must match`

	// Form fields
	ResetPasswordNewPassword         = "New Password"
	ResetPasswordNewPasswordDesc     = "Enter the new administrator password"
	ResetPasswordConfirmPassword     = "Confirm Password"
	ResetPasswordConfirmPasswordDesc = "Re-enter the new password to confirm"

	// Status messages
	ResetPasswordNotAvailable = "PentAGI must be running to reset password"
	ResetPasswordAvailable    = "Password reset is available"
	ResetPasswordInProgress   = "Resetting password..."
	ResetPasswordSuccess      = "Password has been successfully reset"
	ResetPasswordErrorPrefix  = "Error: "

	// Validation errors
	ResetPasswordErrorEmptyPassword = "Password cannot be empty"
	ResetPasswordErrorShortPassword = "Password must be at least 5 characters long"
	ResetPasswordErrorMismatch      = "Passwords do not match"

	// Help content
	ResetPasswordHelpContent = `Reset the administrator password for accessing PentAGI.

This operation:
• Updates the password for admin@pentagi.com account
• Sets the user status to 'active'
• Requires PentAGI database to be accessible
• Does not affect other user accounts

The password change takes effect immediately after successful completion.

Enter the same password in both fields and press Enter to confirm the change.`
)

// Processor Operation Form constants
const (
	// Dynamic title templates
	ProcessorOperationFormTitle       = "%s"
	ProcessorOperationFormDescription = "Execute %s operation"
	ProcessorOperationFormName        = "%s"

	// Common status messages
	ProcessorOperationNotStarted = "Ready to execute %s operation"
	ProcessorOperationInProgress = "Executing %s operation...\n"
	ProcessorOperationCompleted  = "%s operation completed successfully\n"
	ProcessorOperationFailed     = "Failed to execute %s operation"

	// Confirmation messages
	ProcessorOperationConfirmation = "Are you sure you want to %s?"
	ProcessorOperationPressEnter   = "Press Enter to %s"
	ProcessorOperationPressYN      = "Press Y to confirm, N to cancel"
	// Short notice without hotkeys (for static help panel)
	ProcessorOperationRequiresConfirmationShort = "This operation requires confirmation"
	// Additional terminal messages
	ProcessorOperationCancelled = "Operation cancelled"
	ProcessorOperationUnknown   = "Unknown operation: %s"

	// Operation specific messages
	ProcessorOperationStarting    = "Starting services..."
	ProcessorOperationStopping    = "Stopping services..."
	ProcessorOperationRestarting  = "Restarting services..."
	ProcessorOperationDownloading = "Downloading images..."
	ProcessorOperationUpdating    = "Updating components..."
	ProcessorOperationResetting   = "Resetting to factory defaults..."
	ProcessorOperationRemoving    = "Removing containers..."
	ProcessorOperationPurging     = "Purging all data..."
	ProcessorOperationInstalling  = "Installing PentAGI services..."

	// Help text templates
	ProcessorOperationHelpTitle           = "%s Operation"
	ProcessorOperationHelpContent         = "This operation will %s."
	ProcessorOperationHelpContentDownload = "This operation will download %s components."
	ProcessorOperationHelpContentUpdate   = "This operation will update %s components."
	// Generic title/description/builders for dynamic operations
	OperationTitleInstallPentagi    = "Install PentAGI"
	OperationDescInstallPentagi     = "Install and configure PentAGI services"
	OperationTitleDownload          = "Download %s"
	OperationDescDownloadComponents = "Download %s components"
	OperationTitleUpdate            = "Update %s"
	OperationDescUpdateToLatest     = "Update %s to latest version"
	OperationTitleExecute           = "Execute %s"
	OperationDescExecuteOn          = "Execute %s on %s"
	OperationProgressExecuting      = "Executing %s..."

	// Terminal not initialized
	ProcessorOperationTerminalNotInitialized = "Terminal is not initialized"
)

// Operation-specific help texts
const (
	ProcessorHelpInstallPentagi = `This will:
• Deploy Docker containers for selected services
• Configure networking and volumes
• Start all enabled services
• Set up monitoring if configured

Installation will use your current configuration settings.`

	ProcessorHelpStartPentagi = `This will:
• Core PentAGI API and web interface
• Configured Langfuse analytics (if enabled)
• Observability stack (if enabled)

Services will be started in the correct dependency order.`

	ProcessorHelpStopPentagi = `This will:
• Gracefully shutdown containers
• Preserve all data and configurations
• Network connections will be closed

You can restart services later without losing any data.`

	ProcessorHelpRestartPentagi = `This will:
• Stop running containers
• Apply any configuration changes
• Start services with fresh state

Useful after configuration updates or to resolve issues.`

	ProcessorHelpDownloadWorkerImage = `This large image (6GB+) contains:
• Kali Linux tools and utilities
• Security testing frameworks
• Network analysis software

Required for pentesting operations.`

	ProcessorHelpUpdateWorkerImage = `This will:
• Pull the latest pentesting image
• Update security tools and frameworks
• Preserve existing worker containers

Note: This is a large download (6GB+).`

	ProcessorHelpUpdatePentagi = `This will:
• Download latest container images
• Perform rolling update of services
• Preserve all data and configurations

Services will be briefly unavailable during update.`

	ProcessorHelpUpdateInstaller = `This will:
• Download the latest installer binary
• Replace the current installer
• Exit for manual restart

You'll need to restart the installer after update.`

	ProcessorHelpFactoryReset = `⚠️  WARNING: This operation will:
• Remove all containers and networks
• Delete all configuration files
• Clear stored data and volumes
• Restore default settings

This action cannot be undone!`

	ProcessorHelpRemovePentagi = `This will:
• Stop and remove all containers
• Remove Docker networks
• Preserve volumes and data
• Keep configuration files

You can reinstall later without losing data.`

	ProcessorHelpPurgePentagi = `⚠️  WARNING: This will permanently delete:
• All containers and images
• All data volumes
• All configuration files
• All stored results

This action cannot be undone!`
)

// environment variable descriptions (centralized)
const (
	EnvDesc_OPEN_AI_KEY                       = "OpenAI API Key"
	EnvDesc_OPEN_AI_SERVER_URL                = "OpenAI Server URL"
	EnvDesc_ANTHROPIC_API_KEY                 = "Anthropic API Key"
	EnvDesc_ANTHROPIC_SERVER_URL              = "Anthropic Server URL"
	EnvDesc_GEMINI_API_KEY                    = "Google Gemini API Key"
	EnvDesc_GEMINI_SERVER_URL                 = "Gemini Server URL"
	EnvDesc_BEDROCK_DEFAULT_AUTH              = "AWS Bedrock Use Default Credential Chain"
	EnvDesc_BEDROCK_BEARER_TOKEN              = "AWS Bedrock Bearer Token"
	EnvDesc_BEDROCK_ACCESS_KEY_ID             = "AWS Bedrock Access Key ID"
	EnvDesc_BEDROCK_SECRET_ACCESS_KEY         = "AWS Bedrock Secret Access Key"
	EnvDesc_BEDROCK_SESSION_TOKEN             = "AWS Bedrock Session Token"
	EnvDesc_BEDROCK_REGION                    = "AWS Bedrock Region"
	EnvDesc_BEDROCK_SERVER_URL                = "AWS Bedrock Custom Endpoint URL"
	EnvDesc_OLLAMA_SERVER_URL                 = "Ollama Server URL"
	EnvDesc_OLLAMA_SERVER_API_KEY             = "Ollama Server API Key (Cloud)"
	EnvDesc_OLLAMA_SERVER_MODEL               = "Ollama Default Model"
	EnvDesc_OLLAMA_SERVER_CONFIG_PATH         = "Ollama Container Config Path"
	EnvDesc_OLLAMA_SERVER_PULL_MODELS_TIMEOUT = "Ollama Model Pull Timeout"
	EnvDesc_OLLAMA_SERVER_PULL_MODELS_ENABLED = "Ollama Auto-pull Models"
	EnvDesc_OLLAMA_SERVER_LOAD_MODELS_ENABLED = "Ollama Load Models List"
	EnvDesc_DEEPSEEK_API_KEY                  = "DeepSeek API Key"
	EnvDesc_DEEPSEEK_SERVER_URL               = "DeepSeek Server URL"
	EnvDesc_DEEPSEEK_PROVIDER                 = "DeepSeek Provider Name Prefix (for LiteLLM, e.g., 'deepseek')"
	EnvDesc_GLM_API_KEY                       = "GLM API Key"
	EnvDesc_GLM_SERVER_URL                    = "GLM Server URL"
	EnvDesc_GLM_PROVIDER                      = "GLM Provider Name Prefix (for LiteLLM, e.g., 'zai')"
	EnvDesc_KIMI_API_KEY                      = "Kimi API Key"
	EnvDesc_KIMI_SERVER_URL                   = "Kimi Server URL"
	EnvDesc_KIMI_PROVIDER                     = "Kimi Provider Name Prefix (for LiteLLM, e.g., 'moonshot')"
	EnvDesc_QWEN_API_KEY                      = "Qwen API Key"
	EnvDesc_QWEN_SERVER_URL                   = "Qwen Server URL"
	EnvDesc_QWEN_PROVIDER                     = "Qwen Provider Name Prefix (for LiteLLM, e.g., 'dashscope')"
	EnvDesc_LLM_SERVER_URL                    = "Custom LLM Server URL"
	EnvDesc_LLM_SERVER_KEY                    = "Custom LLM API Key"
	EnvDesc_LLM_SERVER_MODEL                  = "Custom LLM Model"
	EnvDesc_LLM_SERVER_CONFIG_PATH            = "Custom LLM Container Config Path"
	EnvDesc_LLM_SERVER_LEGACY_REASONING       = "Custom LLM Legacy Reasoning"
	EnvDesc_LLM_SERVER_PRESERVE_REASONING     = "Custom LLM Preserve Reasoning Content"
	EnvDesc_LLM_SERVER_PROVIDER               = "Custom LLM Provider Name"

	EnvDesc_LANGFUSE_LISTEN_IP   = "Langfuse Listen IP"
	EnvDesc_LANGFUSE_LISTEN_PORT = "Langfuse Listen Port"
	EnvDesc_LANGFUSE_BASE_URL    = "Langfuse Base URL"
	EnvDesc_LANGFUSE_PROJECT_ID  = "Langfuse Project ID"
	EnvDesc_LANGFUSE_PUBLIC_KEY  = "Langfuse Public Key"
	EnvDesc_LANGFUSE_SECRET_KEY  = "Langfuse Secret Key"

	// langfuse init variables
	EnvDesc_LANGFUSE_INIT_PROJECT_ID         = "Langfuse Init Project ID"
	EnvDesc_LANGFUSE_INIT_PROJECT_PUBLIC_KEY = "Langfuse Init Project Public Key"
	EnvDesc_LANGFUSE_INIT_PROJECT_SECRET_KEY = "Langfuse Init Project Secret Key"
	EnvDesc_LANGFUSE_INIT_USER_EMAIL         = "Langfuse Init User Email"
	EnvDesc_LANGFUSE_INIT_USER_NAME          = "Langfuse Init User Name"
	EnvDesc_LANGFUSE_INIT_USER_PASSWORD      = "Langfuse Init User Password"

	EnvDesc_LANGFUSE_OTEL_EXPORTER_OTLP_ENDPOINT = "Langfuse OTLP endpoint for OpenTelemetry exporter"

	EnvDesc_GRAFANA_LISTEN_IP     = "Grafana Listen IP"
	EnvDesc_GRAFANA_LISTEN_PORT   = "Grafana Listen Port"
	EnvDesc_OTEL_GRPC_LISTEN_IP   = "OTel gRPC Listen IP"
	EnvDesc_OTEL_GRPC_LISTEN_PORT = "OTel gRPC Listen Port"
	EnvDesc_OTEL_HTTP_LISTEN_IP   = "OTel HTTP Listen IP"
	EnvDesc_OTEL_HTTP_LISTEN_PORT = "OTel HTTP Listen Port"
	EnvDesc_OTEL_HOST             = "OpenTelemetry Host"

	EnvDesc_SUMMARIZER_PRESERVE_LAST       = "Summarizer Preserve Last"
	EnvDesc_SUMMARIZER_USE_QA              = "Summarizer Use QA"
	EnvDesc_SUMMARIZER_SUM_MSG_HUMAN_IN_QA = "Summarizer Human in QA"
	EnvDesc_SUMMARIZER_LAST_SEC_BYTES      = "Summarizer Last Section Bytes"
	EnvDesc_SUMMARIZER_MAX_BP_BYTES        = "Summarizer Max BP Bytes"
	EnvDesc_SUMMARIZER_MAX_QA_BYTES        = "Summarizer Max QA Bytes"
	EnvDesc_SUMMARIZER_MAX_QA_SECTIONS     = "Summarizer Max QA Sections"
	EnvDesc_SUMMARIZER_KEEP_QA_SECTIONS    = "Summarizer Keep QA Sections"

	EnvDesc_ASSISTANT_SUMMARIZER_PRESERVE_LAST    = "Assistant Summarizer Preserve Last"
	EnvDesc_ASSISTANT_SUMMARIZER_LAST_SEC_BYTES   = "Assistant Summarizer Last Section Bytes"
	EnvDesc_ASSISTANT_SUMMARIZER_MAX_BP_BYTES     = "Assistant Summarizer Max BP Bytes"
	EnvDesc_ASSISTANT_SUMMARIZER_MAX_QA_BYTES     = "Assistant Summarizer Max QA Bytes"
	EnvDesc_ASSISTANT_SUMMARIZER_MAX_QA_SECTIONS  = "Assistant Summarizer Max QA Sections"
	EnvDesc_ASSISTANT_SUMMARIZER_KEEP_QA_SECTIONS = "Assistant Summarizer Keep QA Sections"

	EnvDesc_EMBEDDING_PROVIDER        = "Embedding Provider"
	EnvDesc_EMBEDDING_URL             = "Embedding URL"
	EnvDesc_EMBEDDING_KEY             = "Embedding API Key"
	EnvDesc_EMBEDDING_MODEL           = "Embedding Model"
	EnvDesc_EMBEDDING_BATCH_SIZE      = "Embedding Batch Size"
	EnvDesc_EMBEDDING_STRIP_NEW_LINES = "Embedding Strip New Lines"

	EnvDesc_ASK_USER = "Human-in-the-loop"

	EnvDesc_ASSISTANT_USE_AGENTS = "Enable multi-agent mode for assistant"

	EnvDesc_EXECUTION_MONITOR_ENABLED          = "Enable Execution Monitoring (beta)"
	EnvDesc_EXECUTION_MONITOR_SAME_TOOL_LIMIT  = "Same Tool Call Threshold"
	EnvDesc_EXECUTION_MONITOR_TOTAL_TOOL_LIMIT = "Total Tool Call Threshold"
	EnvDesc_MAX_GENERAL_AGENT_TOOL_CALLS       = "Max Tool Calls for General Agents"
	EnvDesc_MAX_LIMITED_AGENT_TOOL_CALLS       = "Max Tool Calls for Limited Agents"
	EnvDesc_AGENT_PLANNING_STEP_ENABLED        = "Enable Task Planning (beta)"

	EnvDesc_SCRAPER_PUBLIC_URL                    = "Scraper Public URL"
	EnvDesc_SCRAPER_PRIVATE_URL                   = "Scraper Private URL"
	EnvDesc_LOCAL_SCRAPER_USERNAME                = "Local Scraper Username"
	EnvDesc_LOCAL_SCRAPER_PASSWORD                = "Local Scraper Password"
	EnvDesc_LOCAL_SCRAPER_MAX_CONCURRENT_SESSIONS = "Scraper Max Concurrent Sessions"

	EnvDesc_DUCKDUCKGO_ENABLED    = "DuckDuckGo Search"
	EnvDesc_DUCKDUCKGO_REGION     = "DuckDuckGo Region"
	EnvDesc_DUCKDUCKGO_SAFESEARCH = "DuckDuckGo Safe Search"
	EnvDesc_DUCKDUCKGO_TIME_RANGE = "DuckDuckGo Time Range"
	EnvDesc_SPLOITUS_ENABLED      = "Sploitus Search"
	EnvDesc_PERPLEXITY_API_KEY    = "Perplexity API Key"
	EnvDesc_TAVILY_API_KEY        = "Tavily API Key"
	EnvDesc_TRAVERSAAL_API_KEY    = "Traversaal API Key"
	EnvDesc_GOOGLE_API_KEY        = "Google Search API Key"
	EnvDesc_GOOGLE_CX_KEY         = "Google Search CX Key"
	EnvDesc_GOOGLE_LR_KEY         = "Google Search LR Key"

	EnvDesc_DOCKER_INSIDE                    = "Docker Inside Container"
	EnvDesc_DOCKER_NET_ADMIN                 = "Docker Network Admin"
	EnvDesc_DOCKER_SOCKET                    = "Docker Socket Path"
	EnvDesc_DOCKER_NETWORK                   = "Docker Network"
	EnvDesc_DOCKER_PUBLIC_IP                 = "Docker Public IP"
	EnvDesc_DOCKER_WORK_DIR                  = "Docker Work Directory"
	EnvDesc_DOCKER_DEFAULT_IMAGE             = "Docker Default Image"
	EnvDesc_DOCKER_DEFAULT_IMAGE_FOR_PENTEST = "Docker Pentest Image"
	EnvDesc_DOCKER_HOST                      = "Docker Host"
	EnvDesc_DOCKER_TLS_VERIFY                = "Docker TLS Verify"
	EnvDesc_DOCKER_CERT_PATH                 = "Docker Certificate Path"

	EnvDesc_LICENSE_KEY                       = "PentAGI License Key"
	EnvDesc_PENTAGI_LISTEN_IP                 = "PentAGI Server Host"
	EnvDesc_PENTAGI_LISTEN_PORT               = "PentAGI Server Port"
	EnvDesc_PUBLIC_URL                        = "PentAGI Public URL"
	EnvDesc_CORS_ORIGINS                      = "PentAGI CORS Origins"
	EnvDesc_COOKIE_SIGNING_SALT               = "PentAGI Cookie Signing Salt"
	EnvDesc_PROXY_URL                         = "HTTP/HTTPS Proxy URL"
	EnvDesc_HTTP_CLIENT_TIMEOUT               = "HTTP Client Timeout (seconds)"
	EnvDesc_EXTERNAL_SSL_CA_PATH              = "Custom CA Certificate Path"
	EnvDesc_EXTERNAL_SSL_INSECURE             = "Skip SSL Verification"
	EnvDesc_PENTAGI_SSL_DIR                   = "PentAGI SSL Directory"
	EnvDesc_PENTAGI_DATA_DIR                  = "PentAGI Data Directory"
	EnvDesc_PENTAGI_DOCKER_SOCKET             = "Mount Docker Socket Path"
	EnvDesc_PENTAGI_DOCKER_CERT_PATH          = "Mount Docker Certificate Path"
	EnvDesc_PENTAGI_LLM_SERVER_CONFIG_PATH    = "Custom LLM Host Config Path"
	EnvDesc_PENTAGI_OLLAMA_SERVER_CONFIG_PATH = "Ollama Host Config Path"

	EnvDesc_STATIC_DIR     = "Frontend Static Directory"
	EnvDesc_STATIC_URL     = "Frontend Static URL"
	EnvDesc_SERVER_PORT    = "Backend Server Port"
	EnvDesc_SERVER_HOST    = "Backend Server Host"
	EnvDesc_SERVER_SSL_CRT = "Backend Server SSL Certificate Path"
	EnvDesc_SERVER_SSL_KEY = "Backend Server SSL Key Path"
	EnvDesc_SERVER_USE_SSL = "Backend Server Use SSL"

	EnvDesc_PERPLEXITY_MODEL        = "Perplexity Model"
	EnvDesc_PERPLEXITY_CONTEXT_SIZE = "Perplexity Context Size"

	EnvDesc_SEARXNG_URL        = "Searxng Search URL"
	EnvDesc_SEARXNG_CATEGORIES = "Searxng Search Categories"
	EnvDesc_SEARXNG_LANGUAGE   = "Searxng Search Language"
	EnvDesc_SEARXNG_SAFESEARCH = "Searxng Safe Search"
	EnvDesc_SEARXNG_TIME_RANGE = "Searxng Time Range"
	EnvDesc_SEARXNG_TIMEOUT    = "Searxng Timeout"

	EnvDesc_OAUTH_GOOGLE_CLIENT_ID     = "OAuth Google Client ID"
	EnvDesc_OAUTH_GOOGLE_CLIENT_SECRET = "OAuth Google Client Secret"
	EnvDesc_OAUTH_GITHUB_CLIENT_ID     = "OAuth GitHub Client ID"
	EnvDesc_OAUTH_GITHUB_CLIENT_SECRET = "OAuth GitHub Client Secret"

	EnvDesc_LANGFUSE_EE_LICENSE_KEY   = "Langfuse Enterprise License Key"
	EnvDesc_PENTAGI_POSTGRES_PASSWORD = "PentAGI PostgreSQL Password"

	EnvDesc_GRAPHITI_URL        = "Graphiti Server URL"
	EnvDesc_GRAPHITI_TIMEOUT    = "Graphiti Request Timeout"
	EnvDesc_GRAPHITI_MODEL_NAME = "Graphiti Extraction Model"
	EnvDesc_NEO4J_USER          = "Neo4j Username"
	EnvDesc_NEO4J_DATABASE      = "Neo4j Database Name"
	EnvDesc_NEO4J_PASSWORD      = "Neo4j Database Password"
)

// dynamic, contextual sections used in processor operation forms
const (
	// section headers
	ProcessorSectionCurrentState = "Current state"
	ProcessorSectionPlanned      = "Planned actions"
	ProcessorSectionEffects      = "Effects"

	// component labels
	ProcessorComponentPentagi       = "PentAGI"
	ProcessorComponentLangfuse      = "Langfuse"
	ProcessorComponentObservability = "Observability"

	ProcessorComponentWorkerImage           = "worker image"
	ProcessorComponentComposeStacks         = "compose stacks"
	ProcessorComponentDefaultFiles          = "default files"
	ProcessorItemComposeFiles               = "compose files"
	ProcessorItemComposeStacksImagesVolumes = "compose stacks, images, volumes"

	// common states
	ProcessorStateInstalled = "installed"
	ProcessorStateMissing   = "not installed"
	ProcessorStateRunning   = "running"
	ProcessorStateStopped   = "stopped"
	ProcessorStateEmbedded  = "embedded"
	ProcessorStateExternal  = "external"
	ProcessorStateConnected = "connected"
	ProcessorStateDisabled  = "disabled"
	ProcessorStateUnknown   = "unknown"

	// planned action bullet prefixes
	PlannedWillStart    = "will start:"
	PlannedWillStop     = "will stop:"
	PlannedWillRestart  = "will restart:"
	PlannedWillUpdate   = "will update:"
	PlannedWillSkip     = "will skip:"
	PlannedWillRemove   = "will remove:"
	PlannedWillPurge    = "will purge:"
	PlannedWillDownload = "will download:"
	PlannedWillRestore  = "will restore:"

	// effect notes per operation (concise and practical)
	EffectsStart           = "PentAGI web UI becomes available. Background services are brought online in the required order."
	EffectsStop            = "Web UI becomes unavailable. In-progress flows pause safely. When you start PentAGI again, flows resume automatically. A small portion of the current agent step may be lost."
	EffectsRestart         = "Services stop and start again with a clean state. Brief downtime is expected. Flows resume automatically afterwards."
	EffectsUpdateAll       = "Images are pulled and services are recreated where needed. External or disabled components are skipped. Temporary downtime is expected."
	EffectsDownloadWorker  = "Running worker containers are not touched. New flows will use the downloaded image. To switch an existing flow to the new image, finish the flow and start a new task or create a new assistant."
	EffectsUpdateWorker    = "Pulls latest worker image. Running worker containers keep using the old image; new containers will use the updated one."
	EffectsUpdateInstaller = "The installer binary will be updated and the app will exit. Start the installer again to continue."
	EffectsFactoryReset    = "Removes containers, volumes and networks, restores default .env and embedded files. Produces a clean baseline. This action cannot be undone."
	EffectsRemove          = "Stops and removes containers but keeps volumes and images. Data is preserved. Web UI becomes unavailable until you start again."
	EffectsPurge           = "Complete cleanup: containers, images, volumes and configuration files are deleted. Irreversible."
	EffectsInstall         = "Required files are created and services are started. External components are detected and skipped."
)
