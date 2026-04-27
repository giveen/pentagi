package checker

import (
	"bytes"
	"context"
	"encoding/json"
	"net"
	"net/http"
	"net/url"
	"time"

	"github.com/docker/docker/client"
)

type CheckUpdatesRequest struct {
	InstallerOsType         string  `json:"installer_os_type"`
	InstallerVersion        string  `json:"installer_version"`
	PentagiImageName        *string `json:"pentagi_image_name,omitempty"`
	PentagiImageTag         *string `json:"pentagi_image_tag,omitempty"`
	PentagiImageHash        *string `json:"pentagi_image_hash,omitempty"`
	WorkerImageName         *string `json:"worker_image_name,omitempty"`
	WorkerImageTag          *string `json:"worker_image_tag,omitempty"`
	WorkerImageHash         *string `json:"worker_image_hash,omitempty"`
	GraphitiConnected       bool    `json:"graphiti_connected"`
	GraphitiInstalled       bool    `json:"graphiti_installed"`
	GraphitiExternal        bool    `json:"graphiti_external"`
	GraphitiImageName       *string `json:"graphiti_image_name,omitempty"`
	GraphitiImageTag        *string `json:"graphiti_image_tag,omitempty"`
	GraphitiImageHash       *string `json:"graphiti_image_hash,omitempty"`
	Neo4jImageName          *string `json:"neo4j_image_name,omitempty"`
	Neo4jImageTag           *string `json:"neo4j_image_tag,omitempty"`
	Neo4jImageHash          *string `json:"neo4j_image_hash,omitempty"`
	LangfuseConnected       bool    `json:"langfuse_connected"`
	LangfuseInstalled       bool    `json:"langfuse_installed"`
	LangfuseExternal        bool    `json:"langfuse_external"`
	ObservabilityConnected  bool    `json:"observability_connected"`
	ObservabilityExternal   bool    `json:"observability_external"`
	ObservabilityInstalled  bool    `json:"observability_installed"`
	LangfuseWorkerImageName *string `json:"langfuse_worker_image_name,omitempty"`
	LangfuseWorkerImageTag  *string `json:"langfuse_worker_image_tag,omitempty"`
	LangfuseWorkerImageHash *string `json:"langfuse_worker_image_hash,omitempty"`
	LangfuseWebImageName    *string `json:"langfuse_web_image_name,omitempty"`
	LangfuseWebImageTag     *string `json:"langfuse_web_image_tag,omitempty"`
	LangfuseWebImageHash    *string `json:"langfuse_web_image_hash,omitempty"`
	GrafanaImageName        *string `json:"grafana_image_name,omitempty"`
	GrafanaImageTag         *string `json:"grafana_image_tag,omitempty"`
	GrafanaImageHash        *string `json:"grafana_image_hash,omitempty"`
	OpenTelemetryImageName  *string `json:"otel_image_name,omitempty"`
	OpenTelemetryImageTag   *string `json:"otel_image_tag,omitempty"`
	OpenTelemetryImageHash  *string `json:"otel_image_hash,omitempty"`
}

type CheckUpdatesResponse struct {
	InstallerIsUpToDate     bool `json:"installer_is_up_to_date"`
	PentagiIsUpToDate       bool `json:"pentagi_is_up_to_date"`
	GraphitiIsUpToDate      bool `json:"graphiti_is_up_to_date"`
	LangfuseIsUpToDate      bool `json:"langfuse_is_up_to_date"`
	ObservabilityIsUpToDate bool `json:"observability_is_up_to_date"`
	WorkerIsUpToDate        bool `json:"worker_is_up_to_date"`
}

func getNetworkFailures(ctx context.Context, proxyURL string, dockerClient, workerClient *client.Client) []string {
	var failures []string

	if !checkDNSResolution("docker.io") {
		failures = append(failures, "• DNS resolution failed for docker.io")
	}

	if !checkHTTPConnectivity(ctx, proxyURL) {
		failures = append(failures, "• Cannot reach external services via HTTPS")
	}

	if dockerClient != nil && workerClient != nil && !checkDockerPullConnectivity(ctx, dockerClient, workerClient) {
		failures = append(failures, "• Cannot pull Docker images from registry")
	}

	return failures
}

func checkUpdatesServer(
	ctx context.Context,
	serverURL, proxyURL string,
	request CheckUpdatesRequest,
) *CheckUpdatesResponse {
	jsonData, err := json.Marshal(request)
	if err != nil {
		return nil
	}

	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	if proxyURL != "" {
		if proxyURLParsed, err := url.Parse(proxyURL); err == nil {
			client.Transport = &http.Transport{
				Proxy: http.ProxyURL(proxyURLParsed),
			}
		}
	}

	fullURL := serverURL + UpdatesCheckEndpoint
	req, err := http.NewRequestWithContext(ctx, "POST", fullURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", UserAgent)

	resp, err := client.Do(req)
	if err != nil {
		return nil
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil
	}

	var response CheckUpdatesResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil
	}

	return &response
}

func checkDNSResolution(hostname string) bool {
	_, err := net.LookupHost(hostname)
	return err == nil
}

func checkHTTPConnectivity(ctx context.Context, proxyURL string) bool {
	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	if proxyURL != "" {
		if proxyURLParsed, err := url.Parse(proxyURL); err == nil {
			client.Transport = &http.Transport{
				Proxy: http.ProxyURL(proxyURLParsed),
			}
		}
	}

	req, err := http.NewRequestWithContext(ctx, "GET", "https://docker.io", nil)
	if err != nil {
		return false
	}

	resp, err := client.Do(req)
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	return resp.StatusCode < 500
}
