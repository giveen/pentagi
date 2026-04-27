package checker

import (
	"context"
	"errors"

	"pentagi/pkg/version"
)

var (
	InstallerVersion = version.GetBinaryVersion()
	UserAgent        = "PentAGI-Installer/" + InstallerVersion
)

const (
	DockerComposeFile            = "docker-compose.yml"
	GraphitiComposeFile          = "docker-compose-graphiti.yml"
	LangfuseComposeFile          = "docker-compose-langfuse.yml"
	ObservabilityComposeFile     = "docker-compose-observability.yml"
	ExampleCustomConfigLLMFile   = "example.custom.provider.yml"
	ExampleOllamaConfigLLMFile   = "example.ollama.provider.yml"
	PentagiScriptFile            = "/usr/local/bin/pentagi"
	PentagiContainerName         = "pentagi"
	GraphitiContainerName        = "graphiti"
	Neo4jContainerName           = "neo4j"
	LangfuseWorkerContainerName  = "langfuse-worker"
	LangfuseWebContainerName     = "langfuse-web"
	GrafanaContainerName         = "grafana"
	OpenTelemetryContainerName   = "otel"
	DefaultImage                 = "debian:latest"
	DefaultImageForPentest       = "vxcontrol/kali-linux"
	DefaultGraphitiEndpoint      = "http://graphiti:8000"
	DefaultLangfuseEndpoint      = "http://langfuse-web:3000"
	DefaultObservabilityEndpoint = "otelcol:8148"
	DefaultLangfuseOtelEndpoint  = "http://otelcol:4318"
	DefaultUpdateServerEndpoint  = "https://update.pentagi.com"
	UpdatesCheckEndpoint         = "/api/v1/updates/check"
	MinFreeMemGB                 = 0.5
	MinFreeMemGBForPentagi       = 0.5
	MinFreeMemGBForGraphiti      = 2.0
	MinFreeMemGBForLangfuse      = 1.5
	MinFreeMemGBForObservability = 1.5
	MinFreeDiskGB                = 5.0
	MinFreeDiskGBForComponents   = 10.0
	MinFreeDiskGBPerComponents   = 2.0
	MinFreeDiskGBForWorkerImages = 25.0
)

var (
	ErrAppStateNotInitialized = errors.New("appState not initialized")
	ErrHandlerNotInitialized  = errors.New("handler not initialized")
)

type CheckResult struct {
	EnvFileExists           bool   `json:"env_file_exists" yaml:"env_file_exists"`
	DockerApiAccessible     bool   `json:"docker_api_accessible" yaml:"docker_api_accessible"`
	WorkerEnvApiAccessible  bool   `json:"worker_env_api_accessible" yaml:"worker_env_api_accessible"`
	WorkerImageExists       bool   `json:"worker_image_exists" yaml:"worker_image_exists"`
	DockerInstalled         bool   `json:"docker_installed" yaml:"docker_installed"`
	DockerComposeInstalled  bool   `json:"docker_compose_installed" yaml:"docker_compose_installed"`
	DockerVersion           string `json:"docker_version" yaml:"docker_version"`
	DockerVersionOK         bool   `json:"docker_version_ok" yaml:"docker_version_ok"`
	DockerComposeVersion    string `json:"docker_compose_version" yaml:"docker_compose_version"`
	DockerComposeVersionOK  bool   `json:"docker_compose_version_ok" yaml:"docker_compose_version_ok"`
	PentagiScriptInstalled  bool   `json:"pentagi_script_installed" yaml:"pentagi_script_installed"`
	PentagiExtracted        bool   `json:"pentagi_extracted" yaml:"pentagi_extracted"`
	PentagiInstalled        bool   `json:"pentagi_installed" yaml:"pentagi_installed"`
	PentagiRunning          bool   `json:"pentagi_running" yaml:"pentagi_running"`
	PentagiVolumesExist     bool   `json:"pentagi_volumes_exist" yaml:"pentagi_volumes_exist"`
	GraphitiConnected       bool   `json:"graphiti_connected" yaml:"graphiti_connected"`
	GraphitiExternal        bool   `json:"graphiti_external" yaml:"graphiti_external"`
	GraphitiExtracted       bool   `json:"graphiti_extracted" yaml:"graphiti_extracted"`
	GraphitiInstalled       bool   `json:"graphiti_installed" yaml:"graphiti_installed"`
	GraphitiRunning         bool   `json:"graphiti_running" yaml:"graphiti_running"`
	GraphitiVolumesExist    bool   `json:"graphiti_volumes_exist" yaml:"graphiti_volumes_exist"`
	LangfuseConnected       bool   `json:"langfuse_connected" yaml:"langfuse_connected"`
	LangfuseExternal        bool   `json:"langfuse_external" yaml:"langfuse_external"`
	LangfuseExtracted       bool   `json:"langfuse_extracted" yaml:"langfuse_extracted"`
	LangfuseInstalled       bool   `json:"langfuse_installed" yaml:"langfuse_installed"`
	LangfuseRunning         bool   `json:"langfuse_running" yaml:"langfuse_running"`
	LangfuseVolumesExist    bool   `json:"langfuse_volumes_exist" yaml:"langfuse_volumes_exist"`
	ObservabilityConnected  bool   `json:"observability_connected" yaml:"observability_connected"`
	ObservabilityExternal   bool   `json:"observability_external" yaml:"observability_external"`
	ObservabilityExtracted  bool   `json:"observability_extracted" yaml:"observability_extracted"`
	ObservabilityInstalled  bool   `json:"observability_installed" yaml:"observability_installed"`
	ObservabilityRunning    bool   `json:"observability_running" yaml:"observability_running"`
	SysNetworkOK            bool   `json:"sys_network_ok" yaml:"sys_network_ok"`
	SysCPUOK                bool   `json:"sys_cpu_ok" yaml:"sys_cpu_ok"`
	SysMemoryOK             bool   `json:"sys_memory_ok" yaml:"sys_memory_ok"`
	SysDiskFreeSpaceOK      bool   `json:"sys_disk_free_space_ok" yaml:"sys_disk_free_space_ok"`
	UpdateServerAccessible  bool   `json:"update_server_accessible" yaml:"update_server_accessible"`
	InstallerIsUpToDate     bool   `json:"installer_is_up_to_date" yaml:"installer_is_up_to_date"`
	PentagiIsUpToDate       bool   `json:"pentagi_is_up_to_date" yaml:"pentagi_is_up_to_date"`
	GraphitiIsUpToDate      bool   `json:"graphiti_is_up_to_date" yaml:"graphiti_is_up_to_date"`
	LangfuseIsUpToDate      bool   `json:"langfuse_is_up_to_date" yaml:"langfuse_is_up_to_date"`
	ObservabilityIsUpToDate bool   `json:"observability_is_up_to_date" yaml:"observability_is_up_to_date"`
	WorkerIsUpToDate        bool   `json:"worker_is_up_to_date" yaml:"worker_is_up_to_date"`

	// System resource details for UI display
	SysCPUCount        int             `json:"sys_cpu_count" yaml:"sys_cpu_count"`
	SysMemoryRequired  float64         `json:"sys_memory_required_gb" yaml:"sys_memory_required_gb"`
	SysMemoryAvailable float64         `json:"sys_memory_available_gb" yaml:"sys_memory_available_gb"`
	SysDiskRequired    float64         `json:"sys_disk_required_gb" yaml:"sys_disk_required_gb"`
	SysDiskAvailable   float64         `json:"sys_disk_available_gb" yaml:"sys_disk_available_gb"`
	SysNetworkFailures []string        `json:"sys_network_failures" yaml:"sys_network_failures"`
	DockerErrorType    DockerErrorType `json:"docker_error_type" yaml:"docker_error_type"`
	EnvDirWritable     bool            `json:"env_dir_writable" yaml:"env_dir_writable"`

	// handler controls how information is gathered. If nil, skip gathering
	handler CheckHandler
}

// CheckHandler defines how to gather information into a CheckResult
type CheckHandler interface {
	GatherAllInfo(ctx context.Context, c *CheckResult) error
	GatherDockerInfo(ctx context.Context, c *CheckResult) error
	GatherWorkerInfo(ctx context.Context, c *CheckResult) error
	GatherPentagiInfo(ctx context.Context, c *CheckResult) error
	GatherGraphitiInfo(ctx context.Context, c *CheckResult) error
	GatherLangfuseInfo(ctx context.Context, c *CheckResult) error
	GatherObservabilityInfo(ctx context.Context, c *CheckResult) error
	GatherSystemInfo(ctx context.Context, c *CheckResult) error
	GatherUpdatesInfo(ctx context.Context, c *CheckResult) error
}

// Delegating methods that preserve public API
func (c *CheckResult) GatherAllInfo(ctx context.Context) error {
	if c.handler == nil {
		return ErrHandlerNotInitialized
	}
	return c.handler.GatherAllInfo(ctx, c)
}

func (c *CheckResult) GatherDockerInfo(ctx context.Context) error {
	if c.handler == nil {
		return ErrHandlerNotInitialized
	}
	return c.handler.GatherDockerInfo(ctx, c)
}

func (c *CheckResult) GatherWorkerInfo(ctx context.Context) error {
	if c.handler == nil {
		return ErrHandlerNotInitialized
	}
	return c.handler.GatherWorkerInfo(ctx, c)
}

func (c *CheckResult) GatherPentagiInfo(ctx context.Context) error {
	if c.handler == nil {
		return ErrHandlerNotInitialized
	}
	return c.handler.GatherPentagiInfo(ctx, c)
}

func (c *CheckResult) GatherGraphitiInfo(ctx context.Context) error {
	if c.handler == nil {
		return ErrHandlerNotInitialized
	}
	return c.handler.GatherGraphitiInfo(ctx, c)
}

func (c *CheckResult) GatherLangfuseInfo(ctx context.Context) error {
	if c.handler == nil {
		return ErrHandlerNotInitialized
	}
	return c.handler.GatherLangfuseInfo(ctx, c)
}

func (c *CheckResult) GatherObservabilityInfo(ctx context.Context) error {
	if c.handler == nil {
		return ErrHandlerNotInitialized
	}
	return c.handler.GatherObservabilityInfo(ctx, c)
}

func (c *CheckResult) GatherSystemInfo(ctx context.Context) error {
	if c.handler == nil {
		return ErrHandlerNotInitialized
	}
	return c.handler.GatherSystemInfo(ctx, c)
}

func (c *CheckResult) GatherUpdatesInfo(ctx context.Context) error {
	if c.handler == nil {
		return ErrHandlerNotInitialized
	}
	return c.handler.GatherUpdatesInfo(ctx, c)
}

func GatherWithHandler(ctx context.Context, handler CheckHandler) (CheckResult, error) {
	if handler == nil {
		return CheckResult{}, ErrHandlerNotInitialized
	}

	c := CheckResult{
		handler: handler,
	}

	if err := handler.GatherAllInfo(ctx, &c); err != nil {
		return c, err
	}

	return c, nil
}
