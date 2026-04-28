package locale

// Common status and UI strings
const (
	// Common status and UI strings
	UIStatistics       = "Statistics"
	UIStatus           = "Status: "
	UIMode             = "Mode: "
	UINoConfigSelected = "No configuration selected"
	UILoading          = "Loading..."
	UINotImplemented   = "Not implemented yet"
	UIUnsavedChanges   = "Unsaved changes"
	UIConfigSaved      = "Configuration saved"

	// Status labels
	StatusEnabled       = "Enabled"
	StatusDisabled      = "Disabled"
	StatusConfigured    = "Configured"
	StatusNotConfigured = "Not configured"
	StatusEmbedded      = "Embedded"
	StatusExternal      = "External"

	// Success/Warning messages
	MessageSearchEnginesNone       = "⚠ No search engines configured"
	MessageSearchEnginesConfigured = "✓ %d search engines configured"
	MessageDockerConfigured        = "✓ Docker environment configured"
	MessageDockerNotConfigured     = "⚠ Docker environment not configured"
)

// Legend constants
const (
	LegendConfigured    = "✓ Configured"
	LegendNotConfigured = "✗ Not configured"
)

// Common Navigation Actions (always available)
const (
	NavBack       = "Esc: Back"
	NavExit       = "Ctrl+Q: Exit"
	NavUpDown     = "↑/↓: Scroll/Select"
	NavLeftRight  = "←/→: Move"
	NavPgUpPgDown = "PgUp/PgDn: Page"
	NavHomeEnd    = "Home/End: Start/End"
	NavEnter      = "Enter: Continue"
	NavYn         = "Y/N: Accept/Reject"
	NavCtrlC      = "Ctrl+C: Cancel"
	NavCtrlS      = "Ctrl+S: Save"
	NavCtrlR      = "Ctrl+R: Reset"
	NavCtrlH      = "Ctrl+H: Show/Hide"
	NavTab        = "Tab: Complete"
	NavSeparator  = " • "
)

// Welcome Screen constants
const (
	// Form interface implementation
	WelcomeFormTitle       = "Welcome to PentAGI"
	WelcomeFormDescription = "PentAGI is an autonomous penetration testing platform that leverages AI technologies to perform comprehensive security assessments."
	WelcomeFormName        = "Welcome"
	WelcomeFormOverview    = `System checks verify:
• Environment configuration file presence
• Docker API accessibility and version compatibility
• Worker environment readiness
• System resources (CPU, memory, disk space)
• Network connectivity for external dependencies

Once all checks pass, proceed through the configuration wizard to set up LLM providers, monitoring, and security tools.

The installer guides you through each component setup with recommendations for different deployment scenarios.`

	// Configuration status messages
	WelcomeConfigurationFailed = "⚠ Failed checks: %s"
	WelcomeConfigurationPassed = "✓ All system checks passed"

	// Workflow steps
	WelcomeWorkflowTitle = "Installation Workflow:"
	WelcomeWorkflowStep1 = "1. Accept End User License Agreement"
	WelcomeWorkflowStep2 = "2. Configure LLM providers (OpenAI, Anthropic, etc.)"
	WelcomeWorkflowStep3 = "3. Set up integrations (Langfuse, Observability)"
	WelcomeWorkflowStep4 = "4. Configure security settings"
	WelcomeWorkflowStep5 = "5. Deploy and start PentAGI services"
	WelcomeSystemReady   = "✓ System ready - Press Enter to continue"
)

// Troubleshooting on welcome screen constants
const (
	TroubleshootTitle = "System Requirements Not Met"

	// Environment file issues
	TroubleshootEnvFileTitle = "Environment Configuration Missing"
	TroubleshootEnvFileDesc  = "The .env file is required for PentAGI configuration but was not found or is not readable."
	TroubleshootEnvFileFix   = `To fix:
1. Copy .env.example to .env in your installation directory
2. Edit .env and configure at least one LLM provider API key
3. Ensure the file has read permissions (chmod 644 .env)

Quick fix:
cp .env.example .env && chmod 644 .env`

	// Write permissions
	TroubleshootWritePermTitle = "Write Permissions Required"
	TroubleshootWritePermDesc  = "The installer needs write access to the configuration directory to save settings and deploy services."
	TroubleshootWritePermFix   = `To fix:
1. Check directory permissions: ls -la
2. Grant write access: chmod 755 .
3. Or run installer from a writable location
4. Ensure sufficient disk space is available`

	// Docker not installed
	TroubleshootDockerNotInstalledTitle = "Docker Not Installed"
	TroubleshootDockerNotInstalledDesc  = "Docker is not installed on this system. PentAGI requires Docker to run containers."
	TroubleshootDockerNotInstalledFix   = `To fix:
1. Install Docker Desktop: https://docs.docker.com/get-docker/
2. For Linux: Follow distribution-specific instructions
3. Verify installation: docker --version
4. Ensure docker command is in your PATH`

	// Docker not running
	TroubleshootDockerNotRunningTitle = "Docker Daemon Not Running"
	TroubleshootDockerNotRunningDesc  = "Docker is installed but the daemon is not running. The Docker service must be active."
	TroubleshootDockerNotRunningFix   = `To fix:
1. Start Docker Desktop (Windows/Mac)
2. Linux: sudo systemctl start docker
3. Check status: docker ps
4. If using DOCKER_HOST, verify the remote daemon is accessible`

	// Docker permission issues
	TroubleshootDockerPermissionTitle = "Docker Permission Denied"
	TroubleshootDockerPermissionDesc  = "Your user account lacks permission to access Docker. This is common on Linux systems."
	TroubleshootDockerPermissionFix   = `To fix:
1. Add user to docker group: sudo usermod -aG docker $USER
2. Log out and back in for changes to take effect
3. Or run with sudo (not recommended for production)
4. Verify: docker ps (should work without sudo)`

	// Generic Docker API issues
	TroubleshootDockerAPITitle = "Docker API Connection Failed"
	TroubleshootDockerAPIDesc  = "Cannot establish connection to Docker API. This may be due to configuration or network issues."
	TroubleshootDockerAPIFix   = `To fix:
1. Check DOCKER_HOST environment variable
2. Verify Docker is running: docker version
3. For remote Docker: ensure network connectivity
4. Check firewall settings if using TCP connection
5. Try: export DOCKER_HOST=unix:///var/run/docker.sock`

	// Docker version issues
	TroubleshootDockerVersionTitle = "Docker Version Too Old"
	TroubleshootDockerVersionDesc  = "Your Docker version is incompatible. PentAGI requires Docker 20.0.0 or newer."
	TroubleshootDockerVersionFix   = `To fix:
1. Update Docker to version 20.0.0 or newer
2. Visit https://docs.docker.com/engine/install/

Current version: %s
Required: 20.0.0+`

	// Docker Compose issues
	TroubleshootComposeTitle = "Docker Compose Not Found"
	TroubleshootComposeDesc  = "Docker Compose is required but not installed or not in PATH."
	TroubleshootComposeFix   = `To fix:
1. Install Docker Desktop (includes Compose) or
2. Install standalone: https://docs.docker.com/compose/install/

Verify installation: docker compose version`

	// Docker Compose version issues
	TroubleshootComposeVersionTitle = "Docker Compose Version Too Old"
	TroubleshootComposeVersionDesc  = "Your Docker Compose version is incompatible. PentAGI requires Docker Compose 1.25.0 or newer."
	TroubleshootComposeVersionFix   = `Current version: %s
Required: 1.25.0+

To fix:
1. Update Docker Desktop to latest version
2. Or install newer Docker Compose:
   https://docs.docker.com/compose/install/`

	// Worker environment issues
	TroubleshootWorkerTitle = "Worker Docker Environment Not Accessible"
	TroubleshootWorkerDesc  = "Cannot connect to the Docker environment for worker containers. This may be a remote or local Docker setup issue."
	TroubleshootWorkerFix   = `To fix:
1. For remote Docker, set env vars before installer:
   export DOCKER_HOST=tcp://remote:2376
   export DOCKER_CERT_PATH=/path/to/certs
   export DOCKER_TLS_VERIFY=1
2. Verify connection: docker -H $DOCKER_HOST ps
3. For local Docker, leave these vars unset
4. Check firewall allows Docker port (2375/2376)
5. Ensure certificates are valid if using TLS`

	// CPU issues
	TroubleshootCPUTitle = "Insufficient CPU Cores"
	TroubleshootCPUDesc  = "PentAGI requires at least 2 CPU cores for proper operation."
	TroubleshootCPUFix   = `Your system has %d CPU core(s), but 2+ are required.

For virtual machines:
1. Increase CPU allocation in VM settings
2. Ensure host has sufficient resources

Docker Desktop users:
Settings → Resources → CPUs: Set to 2 or more`

	// Memory issues
	TroubleshootMemoryTitle = "Insufficient Memory"
	TroubleshootMemoryDesc  = "Not enough free memory for selected components."
	TroubleshootMemoryFix   = `Memory requirements:
• Base system: 0.5 GB
• PentAGI core: +0.5 GB
• Langfuse (if enabled): +1.5 GB
• Observability (if enabled): +1.5 GB

Total needed: %.1f GB
Available: %.1f GB

To fix:
1. Close unnecessary applications
2. Increase Docker memory limit
3. Disable optional components (Langfuse/Observability)`

	// Disk space issues
	TroubleshootDiskTitle = "Insufficient Disk Space"
	TroubleshootDiskDesc  = "Not enough free disk space for installation and operation."
	TroubleshootDiskFix   = `Disk requirements:
• Base installation: 5 GB minimum
• With components: 10 GB + 2 GB per component
• Worker images: 25 GB (includes 6GB+ Kali image)

Required: %.1f GB
Available: %.1f GB

To fix:
1. Free up disk space
2. Use external storage for Docker
3. Prune unused Docker resources:
   docker system prune -a`

	// Network issues
	TroubleshootNetworkTitle = "Network Connectivity Failed"
	TroubleshootNetworkDesc  = "Cannot reach required external services. This prevents downloading Docker images and updates."
	TroubleshootNetworkFix   = `Failed checks:
%s

To fix:
1. Verify internet connection: ping docker.io
2. Check DNS resolution: nslookup docker.io
3. If behind proxy, set before running installer:
   export HTTP_PROXY=http://proxy:port
   export HTTPS_PROXY=http://proxy:port
4. For persistent proxy, add to .env:
   PROXY_URL=http://proxy:port
5. Check firewall allows outbound HTTPS (port 443)
6. Try alternative DNS servers if DNS fails`

	// Generic hint at the bottom
	TroubleshootFixHint = "\nResolve the issues above and run the installer again."

	// Network failure messages (used in checker/helpers.go)
	NetworkFailureDNS        = "• DNS resolution failed for docker.io"
	NetworkFailureHTTPS      = "• Cannot reach external services via HTTPS"
	NetworkFailureDockerPull = "• Cannot pull Docker images from registry"
)

// System Checks constants
const (
	ChecksTitle               = "System Checks"
	ChecksWarningFailed       = "⚠ Some checks failed"
	CheckEnvironmentFile      = "Environment file"
	CheckWritePermissions     = "Write permissions"
	CheckDockerAPI            = "Docker API"
	CheckDockerVersion        = "Docker version"
	CheckDockerCompose        = "Docker Compose"
	CheckDockerComposeVersion = "Docker Compose version"
	CheckWorkerEnvironment    = "Worker environment"
	CheckSystemResources      = "System resources"
	CheckNetworkConnectivity  = "Network connectivity"
)

// EULA Screen constants
const (
	// Form interface implementation
	EULAFormDescription = "Legal terms and conditions for PentAGI usage"
	EULAFormName        = "EULA"
	EULAFormOverview    = `Review and accept the End User License Agreement to proceed with PentAGI installation.

The EULA contains:
• Software license terms and usage rights
• Limitation of liability and warranties
• Data collection and privacy policies
• Compliance requirements and restrictions
• Support and maintenance terms

You must scroll through the entire document and accept the terms to continue with the installation process.

Use arrow keys, page up/down, or home/end keys to navigate through the document.`

	// Error and status messages
	EULAErrorLoadingTitle     = "# Error Loading EULA\n\nFailed to load EULA: %v"
	EULAContentFallback       = "# EULA Content\n\n%s\n\n---\n\n*Note: Markdown rendering failed: %v*"
	EULAConfigurationRead     = "✓ EULA reviewed"
	EULAConfigurationAccepted = "✓ EULA accepted"
	EULAConfigurationPending  = "⚠ EULA not reviewed"
	EULALoading               = "Loading EULA..."
	EULAProgress              = "Progress: %d%%"
	EULAProgressComplete      = " • Complete"
)

// Main Menu Screen constants
const (
	MainMenuTitle       = "PentAGI Configuration"
	MainMenuDescription = "Configure all PentAGI components and settings"
	MainMenuName        = "Main Menu"
	MainMenuOverview    = `Welcome to PentAGI Configuration Center.

Configure essential components:
• LLM Providers - AI language models for autonomous testing
• Monitoring - Observability and analytics platforms
• Tools - Additional capabilities for enhanced testing
• System Settings - Environment and deployment options

Navigate through each section to complete your PentAGI setup.`

	MenuTitle        = "Configuration Menu"
	MenuSystemStatus = "System Status"
)

// Main Menu Status Labels (not used)
const (
	MainMenuStatusPentagiRunning     = "PentAGI is already running"
	MainMenuStatusPentagiNotRunning  = "Ready to start PentAGI services"
	MainMenuStatusUpToDate           = "PentAGI is up to date"
	MainMenuStatusUpdatesAvailable   = "Updates are available"
	MainMenuStatusReadyToStart       = "Ready to start"
	MainMenuStatusAllServicesRunning = "All services are running"
	MainMenuStatusNoUpdatesAvailable = "No updates available"
)
