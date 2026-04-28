package checker

import (
	"context"
	"fmt"
	"path/filepath"
	"runtime"
	"sync"

	"pentagi/cmd/installer/state"

	"github.com/docker/docker/client"
)

// defaultCheckHandler provides the existing implementation of gathering logic
type defaultCheckHandler struct {
	mx           *sync.Mutex
	appState     state.State
	dockerClient *client.Client
	workerClient *client.Client
}

func (h *defaultCheckHandler) GatherAllInfo(ctx context.Context, c *CheckResult) error {
	envPath := h.appState.GetEnvPath()
	c.EnvFileExists = checkFileExists(envPath) && checkFileIsReadable(envPath)
	if !c.EnvFileExists {
		return fmt.Errorf("environment file %s does not exist or is not readable", envPath)
	}

	// check write permissions to .env directory
	envDir := filepath.Dir(envPath)
	c.EnvDirWritable = checkDirIsWritable(envDir)

	if err := h.GatherDockerInfo(ctx, c); err != nil {
		return err
	}
	if err := h.GatherWorkerInfo(ctx, c); err != nil {
		return err
	}
	if err := h.GatherPentagiInfo(ctx, c); err != nil {
		return err
	}
	if err := h.GatherGraphitiInfo(ctx, c); err != nil {
		return err
	}
	if err := h.GatherLangfuseInfo(ctx, c); err != nil {
		return err
	}
	if err := h.GatherObservabilityInfo(ctx, c); err != nil {
		return err
	}
	if err := h.GatherSystemInfo(ctx, c); err != nil {
		return err
	}
	if err := h.GatherUpdatesInfo(ctx, c); err != nil {
		return err
	}

	return nil
}

func (h *defaultCheckHandler) GatherDockerInfo(ctx context.Context, c *CheckResult) error {
	h.mx.Lock()
	defer h.mx.Unlock()

	var cli *client.Client

	if cli, c.DockerErrorType = createDockerClientFromEnv(ctx); c.DockerErrorType != DockerErrorNone {
		c.DockerApiAccessible = false
		c.DockerInstalled = c.DockerErrorType != DockerErrorNotInstalled
		if c.DockerInstalled {
			version := checkDockerCliVersion()
			c.DockerVersion = version.Version
			c.DockerVersionOK = version.Valid
		}
	} else {
		h.dockerClient = cli
		c.DockerApiAccessible = true
		c.DockerInstalled = true

		version := checkDockerVersion(ctx, cli)
		c.DockerVersion = version.Version
		c.DockerVersionOK = version.Valid
	}

	composeVersion := checkDockerComposeVersion()
	c.DockerComposeInstalled = composeVersion.Version != ""
	c.DockerComposeVersion = composeVersion.Version
	c.DockerComposeVersionOK = composeVersion.Valid

	return nil
}

func (h *defaultCheckHandler) GatherWorkerInfo(ctx context.Context, c *CheckResult) error {
	h.mx.Lock()
	defer h.mx.Unlock()

	dockerHost := getEnvVar(h.appState, "DOCKER_HOST", "")
	dockerCertPath := getEnvVar(h.appState, "PENTAGI_DOCKER_CERT_PATH", "")
	dockerTLSVerify := getEnvVar(h.appState, "DOCKER_TLS_VERIFY", "") != ""

	cli, err := createDockerClient(dockerHost, dockerCertPath, dockerTLSVerify)
	if err != nil {
		// fallback to DOCKER_CERT_PATH for backward compatibility
		// this handles cases where migration failed or user manually edited .env
		// note: after migration, DOCKER_CERT_PATH contains container path, not host path
		dockerCertPath = getEnvVar(h.appState, "DOCKER_CERT_PATH", "")
		cli, err = createDockerClient(dockerHost, dockerCertPath, dockerTLSVerify)
		if err != nil {
			c.WorkerEnvApiAccessible = false
			c.WorkerImageExists = false
			return nil
		}
	}

	h.workerClient = cli
	c.WorkerEnvApiAccessible = true

	pentestImage := getEnvVar(h.appState, "DOCKER_DEFAULT_IMAGE_FOR_PENTEST", DefaultImageForPentest)
	c.WorkerImageExists = checkImageExists(ctx, cli, pentestImage)

	return nil
}

func (h *defaultCheckHandler) GatherPentagiInfo(ctx context.Context, c *CheckResult) error {
	h.mx.Lock()
	defer h.mx.Unlock()

	envDir := filepath.Dir(h.appState.GetEnvPath())
	dockerComposeFile := filepath.Join(envDir, DockerComposeFile)
	c.PentagiExtracted = checkFileExists(dockerComposeFile) &&
		checkFileExists(ExampleCustomConfigLLMFile) &&
		checkFileExists(ExampleOllamaConfigLLMFile)
	c.PentagiScriptInstalled = checkFileExists(PentagiScriptFile)

	if h.dockerClient != nil {
		exists, running := checkContainerExists(ctx, h.dockerClient, PentagiContainerName)
		c.PentagiInstalled = exists
		c.PentagiRunning = running

		// check if pentagi-related volumes exist (indicates previous installation)
		pentagiVolumes := []string{"pentagi-postgres-data", "pentagi-data", "pentagi-ssl", "scraper-ssl"}
		c.PentagiVolumesExist = checkVolumesExist(ctx, h.dockerClient, pentagiVolumes)
	}

	return nil
}

func (h *defaultCheckHandler) GatherGraphitiInfo(ctx context.Context, c *CheckResult) error {
	h.mx.Lock()
	defer h.mx.Unlock()

	graphitiEnabled := getEnvVar(h.appState, "GRAPHITI_ENABLED", "")
	graphitiURL := getEnvVar(h.appState, "GRAPHITI_URL", "")

	c.GraphitiConnected = graphitiEnabled == "true" && graphitiURL != ""
	c.GraphitiExternal = graphitiURL != DefaultGraphitiEndpoint

	envDir := filepath.Dir(h.appState.GetEnvPath())
	graphitiComposeFile := filepath.Join(envDir, GraphitiComposeFile)
	c.GraphitiExtracted = checkFileExists(graphitiComposeFile)

	if h.dockerClient != nil {
		graphitiExists, graphitiRunning := checkContainerExists(ctx, h.dockerClient, GraphitiContainerName)
		neo4jExists, neo4jRunning := checkContainerExists(ctx, h.dockerClient, Neo4jContainerName)

		c.GraphitiInstalled = graphitiExists && neo4jExists
		c.GraphitiRunning = graphitiRunning && neo4jRunning

		// check if graphiti-related volumes exist (indicates previous installation)
		graphitiVolumes := []string{"neo4j_data"}
		c.GraphitiVolumesExist = checkVolumesExist(ctx, h.dockerClient, graphitiVolumes)
	}

	return nil
}

func (h *defaultCheckHandler) GatherLangfuseInfo(ctx context.Context, c *CheckResult) error {
	h.mx.Lock()
	defer h.mx.Unlock()

	baseURL := getEnvVar(h.appState, "LANGFUSE_BASE_URL", "")
	projectID := getEnvVar(h.appState, "LANGFUSE_PROJECT_ID", "")
	publicKey := getEnvVar(h.appState, "LANGFUSE_PUBLIC_KEY", "")
	secretKey := getEnvVar(h.appState, "LANGFUSE_SECRET_KEY", "")

	c.LangfuseConnected = baseURL != "" && projectID != "" && publicKey != "" && secretKey != ""
	c.LangfuseExternal = baseURL != DefaultLangfuseEndpoint

	envDir := filepath.Dir(h.appState.GetEnvPath())
	langfuseFile := filepath.Join(envDir, LangfuseComposeFile)
	c.LangfuseExtracted = checkFileExists(langfuseFile)

	if h.dockerClient != nil {
		workerExists, workerRunning := checkContainerExists(ctx, h.dockerClient, LangfuseWorkerContainerName)
		webExists, webRunning := checkContainerExists(ctx, h.dockerClient, LangfuseWebContainerName)

		c.LangfuseInstalled = workerExists && webExists
		c.LangfuseRunning = workerRunning && webRunning

		// check if langfuse-related volumes exist (indicates previous installation)
		langfuseVolumes := []string{"langfuse-postgres-data", "langfuse-clickhouse-data", "langfuse-minio-data"}
		c.LangfuseVolumesExist = checkVolumesExist(ctx, h.dockerClient, langfuseVolumes)
	}

	return nil
}

func (h *defaultCheckHandler) GatherObservabilityInfo(ctx context.Context, c *CheckResult) error {
	h.mx.Lock()
	defer h.mx.Unlock()

	otelHost := getEnvVar(h.appState, "OTEL_HOST", "")
	c.ObservabilityConnected = otelHost != ""
	c.ObservabilityExternal = otelHost != DefaultObservabilityEndpoint

	envDir := filepath.Dir(h.appState.GetEnvPath())
	obsFile := filepath.Join(envDir, ObservabilityComposeFile)
	c.ObservabilityExtracted = checkFileExists(obsFile)

	if h.dockerClient != nil {
		exists, running := checkContainerExists(ctx, h.dockerClient, OpenTelemetryContainerName)
		c.ObservabilityInstalled = exists
		c.ObservabilityRunning = running
	}

	return nil
}

func (h *defaultCheckHandler) GatherSystemInfo(ctx context.Context, c *CheckResult) error {
	h.mx.Lock()
	defer h.mx.Unlock()

	// CPU check and count
	c.SysCPUCount = runtime.NumCPU()
	c.SysCPUOK = checkCPUResources()

	// memory check and calculations
	needsForPentagi, needsForGraphiti, needsForLangfuse, needsForObservability := determineComponentNeeds(c)

	// calculate required memory using shared function
	c.SysMemoryRequired = calculateRequiredMemoryGB(needsForPentagi, needsForGraphiti, needsForLangfuse, needsForObservability)

	// get available memory and check if sufficient
	c.SysMemoryAvailable = getAvailableMemoryGB()
	c.SysMemoryOK = checkMemoryResources(needsForPentagi, needsForGraphiti, needsForLangfuse, needsForObservability)

	// disk check and calculations
	localComponents := countLocalComponentsToInstall(
		c.PentagiInstalled,
		c.GraphitiConnected, c.GraphitiExternal, c.GraphitiInstalled,
		c.LangfuseConnected, c.LangfuseExternal, c.LangfuseInstalled,
		c.ObservabilityConnected, c.ObservabilityExternal, c.ObservabilityInstalled,
	)

	// calculate required disk space using shared function
	c.SysDiskRequired = calculateRequiredDiskGB(c.WorkerImageExists, localComponents)

	// get available disk space and check if sufficient
	c.SysDiskAvailable = getAvailableDiskGB(ctx)
	c.SysDiskFreeSpaceOK = checkDiskSpaceWithContext(
		ctx,
		c.WorkerImageExists,
		c.PentagiInstalled,
		c.GraphitiConnected,
		c.GraphitiExternal,
		c.GraphitiInstalled,
		c.LangfuseConnected,
		c.LangfuseExternal,
		c.LangfuseInstalled,
		c.ObservabilityConnected,
		c.ObservabilityExternal,
		c.ObservabilityInstalled,
	)

	// network check with proxy and docker clients
	proxyURL := getProxyURL(h.appState)
	c.SysNetworkFailures = getNetworkFailures(ctx, proxyURL, h.dockerClient, h.workerClient)
	c.SysNetworkOK = len(c.SysNetworkFailures) == 0

	return nil
}

func (h *defaultCheckHandler) GatherUpdatesInfo(ctx context.Context, c *CheckResult) error {
	h.mx.Lock()
	defer h.mx.Unlock()

	proxyURL := getProxyURL(h.appState)
	updateServerURL := getEnvVar(h.appState, "UPDATE_SERVER_URL", DefaultUpdateServerEndpoint)

	request := CheckUpdatesRequest{
		InstallerOsType:        runtime.GOOS,
		InstallerVersion:       InstallerVersion,
		GraphitiConnected:      c.GraphitiConnected,
		GraphitiExternal:       c.GraphitiExternal,
		GraphitiInstalled:      c.GraphitiInstalled,
		LangfuseConnected:      c.LangfuseConnected,
		LangfuseExternal:       c.LangfuseExternal,
		LangfuseInstalled:      c.LangfuseInstalled,
		ObservabilityConnected: c.ObservabilityConnected,
		ObservabilityExternal:  c.ObservabilityExternal,
		ObservabilityInstalled: c.ObservabilityInstalled,
	}

	// get PentAGI container image info
	if h.dockerClient != nil && c.PentagiInstalled {
		if imageInfo := getContainerImageInfo(ctx, h.dockerClient, PentagiContainerName); imageInfo != nil {
			request.PentagiImageName = &imageInfo.Name
			request.PentagiImageTag = &imageInfo.Tag
			request.PentagiImageHash = &imageInfo.Hash
		}
	}

	// get Worker image info from environment
	if h.workerClient != nil {
		defaultImage := getEnvVar(h.appState, "DOCKER_DEFAULT_IMAGE_FOR_PENTEST", DefaultImageForPentest)
		if imageInfo := getImageInfo(ctx, h.workerClient, defaultImage); imageInfo != nil {
			request.WorkerImageName = &imageInfo.Name
			request.WorkerImageTag = &imageInfo.Tag
			request.WorkerImageHash = &imageInfo.Hash
		}
	}

	// get Graphiti image info if installed locally
	if h.dockerClient != nil && c.GraphitiConnected && !c.GraphitiExternal && c.GraphitiInstalled {
		if graphitiInfo := getContainerImageInfo(ctx, h.dockerClient, GraphitiContainerName); graphitiInfo != nil {
			request.GraphitiImageName = &graphitiInfo.Name
			request.GraphitiImageTag = &graphitiInfo.Tag
			request.GraphitiImageHash = &graphitiInfo.Hash
		}
		if neo4jInfo := getContainerImageInfo(ctx, h.dockerClient, Neo4jContainerName); neo4jInfo != nil {
			request.Neo4jImageName = &neo4jInfo.Name
			request.Neo4jImageTag = &neo4jInfo.Tag
			request.Neo4jImageHash = &neo4jInfo.Hash
		}
	}

	// get Langfuse image info if installed locally
	if h.dockerClient != nil && c.LangfuseConnected && !c.LangfuseExternal && c.LangfuseInstalled {
		if workerInfo := getContainerImageInfo(ctx, h.dockerClient, LangfuseWorkerContainerName); workerInfo != nil {
			request.LangfuseWorkerImageName = &workerInfo.Name
			request.LangfuseWorkerImageTag = &workerInfo.Tag
			request.LangfuseWorkerImageHash = &workerInfo.Hash
		}
		if webInfo := getContainerImageInfo(ctx, h.dockerClient, LangfuseWebContainerName); webInfo != nil {
			request.LangfuseWebImageName = &webInfo.Name
			request.LangfuseWebImageTag = &webInfo.Tag
			request.LangfuseWebImageHash = &webInfo.Hash
		}
	}

	// get Grafana and OpenTelemetry image info if observability installed locally
	if h.dockerClient != nil && c.ObservabilityConnected && !c.ObservabilityExternal && c.ObservabilityInstalled {
		if grafanaInfo := getContainerImageInfo(ctx, h.dockerClient, GrafanaContainerName); grafanaInfo != nil {
			request.GrafanaImageName = &grafanaInfo.Name
			request.GrafanaImageTag = &grafanaInfo.Tag
			request.GrafanaImageHash = &grafanaInfo.Hash
		}
		if otelInfo := getContainerImageInfo(ctx, h.dockerClient, OpenTelemetryContainerName); otelInfo != nil {
			request.OpenTelemetryImageName = &otelInfo.Name
			request.OpenTelemetryImageTag = &otelInfo.Tag
			request.OpenTelemetryImageHash = &otelInfo.Hash
		}
	}

	response := checkUpdatesServer(ctx, updateServerURL, proxyURL, request)
	if response != nil {
		c.UpdateServerAccessible = true
		c.InstallerIsUpToDate = response.InstallerIsUpToDate
		c.PentagiIsUpToDate = response.PentagiIsUpToDate
		c.GraphitiIsUpToDate = response.GraphitiIsUpToDate
		c.LangfuseIsUpToDate = response.LangfuseIsUpToDate
		c.ObservabilityIsUpToDate = response.ObservabilityIsUpToDate
		c.WorkerIsUpToDate = response.WorkerIsUpToDate
	} else {
		c.UpdateServerAccessible = false
		c.InstallerIsUpToDate = false
		c.PentagiIsUpToDate = false
		c.GraphitiIsUpToDate = false
		c.LangfuseIsUpToDate = false
		c.ObservabilityIsUpToDate = false
	}

	return nil
}

// Gather builds a CheckResult using the default handler for the given appState.
func Gather(ctx context.Context, appState state.State) (CheckResult, error) {
	if appState == nil {
		return CheckResult{}, ErrAppStateNotInitialized
	}

	c := CheckResult{
		// default to the built-in handler
		handler: &defaultCheckHandler{
			mx:       &sync.Mutex{},
			appState: appState,
		},
	}

	if err := c.GatherAllInfo(ctx); err != nil {
		return c, err
	}

	return c, nil
}
