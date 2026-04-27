package checker

import (
	"context"
	"errors"
	"io"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"syscall"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/api/types/volume"
	"github.com/docker/docker/client"
)

type DockerVersion struct {
	Version string
	Valid   bool
}

type ImageInfo struct {
	Name string
	Tag  string
	Hash string
}

type DockerErrorType string

const (
	DockerErrorNone         DockerErrorType = ""
	DockerErrorNotInstalled DockerErrorType = "not_installed"
	DockerErrorNotRunning   DockerErrorType = "not_running"
	DockerErrorAPIError     DockerErrorType = "api_error"
	DockerErrorPermission   DockerErrorType = "permission"
)

func createDockerClient(host, certPath string, tlsVerify bool) (*client.Client, error) {
	opts := []client.Opt{
		client.WithAPIVersionNegotiation(),
	}

	if host != "" {
		opts = append(opts, client.WithHost(host))
	}

	if tlsVerify && certPath != "" {
		opts = append(opts, client.WithTLSClientConfig(
			filepath.Join(certPath, "ca.pem"),
			filepath.Join(certPath, "cert.pem"),
			filepath.Join(certPath, "key.pem"),
		))
	}

	return client.NewClientWithOpts(opts...)
}

func createDockerClientFromEnv(ctx context.Context) (*client.Client, DockerErrorType) {
	_, err := exec.LookPath("docker")
	if err != nil {
		return nil, DockerErrorNotInstalled
	}

	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, DockerErrorAPIError
	}

	_, err = cli.Ping(ctx)
	if err != nil {
		cli.Close()
		if strings.Contains(err.Error(), "Cannot connect to the Docker daemon") ||
			strings.Contains(err.Error(), "Is the docker daemon running") ||
			strings.Contains(err.Error(), "connection refused") ||
			strings.Contains(err.Error(), "no such host") ||
			strings.Contains(err.Error(), "dial unix") {
			return nil, DockerErrorNotRunning
		}
		if strings.Contains(err.Error(), "permission denied") ||
			strings.Contains(err.Error(), "Got permission denied") {
			return nil, DockerErrorPermission
		}
		return nil, DockerErrorAPIError
	}

	return cli, DockerErrorNone
}

func checkDockerVersion(ctx context.Context, cli *client.Client) DockerVersion {
	version, err := cli.ServerVersion(ctx)
	if err != nil {
		return DockerVersion{Version: "", Valid: false}
	}

	versionStr := version.Version
	valid := checkVersionCompatibility(versionStr, "20.0.0")

	return DockerVersion{Version: versionStr, Valid: valid}
}

func checkDockerCliVersion() DockerVersion {
	_, err := exec.LookPath("docker")
	if err != nil {
		return DockerVersion{Version: "", Valid: false}
	}

	cmd := exec.Command("docker", "version", "--format", "{{.Client.Version}}")
	output, err := cmd.Output()
	if err != nil && len(output) == 0 {
		return DockerVersion{Version: "", Valid: false}
	}

	versionStr := extractVersionFromOutput(string(output))
	valid := checkVersionCompatibility(versionStr, "20.0.0")

	return DockerVersion{Version: versionStr, Valid: valid}
}

func checkDockerComposeVersion() DockerVersion {
	cmd := exec.Command("docker", "compose", "version")
	output, err := cmd.Output()
	if err != nil {
		cmd = exec.Command("docker-compose", "--version")
		output, err = cmd.Output()
		if err != nil {
			return DockerVersion{Version: "", Valid: false}
		}
	}

	versionStr := extractVersionFromOutput(string(output))
	valid := checkVersionCompatibility(versionStr, "1.25.0")

	return DockerVersion{Version: versionStr, Valid: valid}
}

func extractVersionFromOutput(output string) string {
	re := regexp.MustCompile(`v?(\d+\.\d+\.\d+)`)
	matches := re.FindStringSubmatch(output)
	if len(matches) > 1 {
		return matches[1]
	}
	return ""
}

func checkVersionCompatibility(version, minVersion string) bool {
	if version == "" || minVersion == "" {
		return false
	}

	versionParts := strings.Split(version, ".")
	minVersionParts := strings.Split(minVersion, ".")

	for i := 0; i < len(versionParts) && i < len(minVersionParts); i++ {
		v, err1 := strconv.Atoi(versionParts[i])
		minV, err2 := strconv.Atoi(minVersionParts[i])

		if err1 != nil || err2 != nil {
			return false
		}

		if v > minV {
			return true
		}
		if v < minV {
			return false
		}
	}

	return len(versionParts) >= len(minVersionParts)
}

func checkContainerExists(ctx context.Context, cli *client.Client, name string) (exists, running bool) {
	containers, err := cli.ContainerList(ctx, container.ListOptions{All: true})
	if err != nil {
		return false, false
	}

	for _, cont := range containers {
		for _, containerName := range cont.Names {
			if strings.TrimPrefix(containerName, "/") == name {
				return true, cont.State == "running"
			}
		}
	}

	return false, false
}

func checkVolumesExist(ctx context.Context, cli *client.Client, volumeNames []string) bool {
	if cli == nil || len(volumeNames) == 0 {
		return false
	}

	volumes, err := cli.VolumeList(ctx, volume.ListOptions{})
	if err != nil {
		return false
	}

	existingVolumes := make([]string, 0, len(volumes.Volumes))
	for _, vol := range volumes.Volumes {
		existingVolumes = append(existingVolumes, vol.Name)
	}

	for _, volumeName := range volumeNames {
		for _, existingVolume := range existingVolumes {
			if existingVolume == volumeName || strings.HasSuffix(existingVolume, "_"+volumeName) {
				return true
			}
		}
	}

	return false
}

func getContainerImageInfo(ctx context.Context, cli *client.Client, containerName string) *ImageInfo {
	containers, err := cli.ContainerList(ctx, container.ListOptions{All: true})
	if err != nil {
		return nil
	}

	for _, cont := range containers {
		for _, name := range cont.Names {
			if strings.TrimPrefix(name, "/") == containerName {
				return parseImageRef(cont.Image, cont.ImageID)
			}
		}
	}

	return nil
}

func checkImageExists(ctx context.Context, cli *client.Client, imageName string) bool {
	imageInfo := getImageInfo(ctx, cli, imageName)
	return imageInfo != nil && imageInfo.Hash != ""
}

func getImageInfo(ctx context.Context, cli *client.Client, imageName string) *ImageInfo {
	if cli == nil {
		return nil
	}

	images, err := cli.ImageList(ctx, image.ListOptions{})
	if err != nil {
		return nil
	}

	imageInfo := parseImageRef(imageName, "")
	if imageInfo == nil {
		return nil
	}

	fullImageName := imageInfo.Name + ":" + imageInfo.Tag
	for _, img := range images {
		for _, tag := range img.RepoTags {
			if tag == imageName || tag == fullImageName {
				imageInfo.Hash = img.ID
				break
			}
		}
	}

	return imageInfo
}

func parseImageRef(imageRef, imageID string) *ImageInfo {
	if imageRef == "" {
		return nil
	}

	info := &ImageInfo{
		Hash: imageID,
	}

	originalRef := imageRef

	if strings.Contains(imageRef, "@") {
		parts := strings.SplitN(imageRef, "@", 2)
		imageRef = parts[0]
		if len(parts) > 1 {
			info.Hash = parts[1]
		}
	}

	var name string
	if strings.Contains(imageRef, "/") {
		nameParts := strings.Split(imageRef, "/")
		if len(nameParts) >= 2 {
			if strings.Contains(nameParts[0], ".") || strings.Contains(nameParts[0], ":") {
				if len(nameParts) > 2 {
					name = strings.Join(nameParts[1:], "/")
				} else {
					name = nameParts[1]
				}
			} else {
				name = imageRef
			}
		} else {
			name = imageRef
		}
	} else {
		name = imageRef
	}

	if strings.Contains(name, ":") {
		parts := strings.SplitN(name, ":", 2)
		info.Name = parts[0]
		if len(parts) > 1 && parts[1] != "" {
			if !strings.Contains(parts[1], ".") {
				info.Tag = parts[1]
			} else {
				info.Name = name
				info.Tag = "latest"
			}
		}
	} else {
		info.Name = name
		info.Tag = "latest"
	}

	if info.Tag == "" {
		info.Tag = "latest"
	}

	if info.Name == "" {
		info.Name = originalRef
		info.Tag = "latest"
	}

	return info
}

func checkDockerPullConnectivity(ctx context.Context, dockerClient, workerClient *client.Client) bool {
	if dockerClient != nil {
		if checkSingleDockerPull(ctx, dockerClient, DefaultImage) {
			return true
		}
	}

	if workerClient != nil {
		if checkSingleDockerPull(ctx, workerClient, DefaultImage) {
			return true
		}
	}

	return false
}

func checkSingleDockerPull(ctx context.Context, cli *client.Client, imageName string) bool {
	reader, err := cli.ImagePull(ctx, imageName, image.PullOptions{})
	if err != nil {
		return false
	}
	defer reader.Close()

	buf := make([]byte, 1024)
	_, err = reader.Read(buf)
	return err == nil || errors.Is(err, io.EOF) || errors.Is(err, syscall.EIO)
}
