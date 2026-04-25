package tools

import (
	"context"
	"fmt"
	"io"
	"strings"
	"time"

	"pentagi/pkg/docker"

	"github.com/docker/docker/api/types/container"
)

const cancelExecTimeout = 5 * time.Second

func cancelRunningExecCommands(ctx context.Context, dockerClient docker.DockerClient, containerID string) error {
	pids := execTracker.listPIDs(containerID)
	if len(pids) == 0 {
		return nil
	}

	strPIDs := make([]string, 0, len(pids))
	for _, pid := range pids {
		if pid <= 0 {
			continue
		}
		strPIDs = append(strPIDs, fmt.Sprintf("%d", pid))
	}
	if len(strPIDs) == 0 {
		return nil
	}

	killCmd := fmt.Sprintf(
		"for pid in %s; do kill -TERM \"$pid\" 2>/dev/null || true; done",
		strings.Join(strPIDs, " "),
	)

	execResp, err := dockerClient.ContainerExecCreate(ctx, containerID, container.ExecOptions{
		Cmd:          []string{"sh", "-c", killCmd},
		AttachStdout: true,
		AttachStderr: true,
		Tty:          false,
	})
	if err != nil {
		return fmt.Errorf("failed to create kill exec: %w", err)
	}

	execCtx, cancel := context.WithTimeout(ctx, cancelExecTimeout)
	defer cancel()

	attachResp, err := dockerClient.ContainerExecAttach(execCtx, execResp.ID, container.ExecAttachOptions{Tty: false})
	if err != nil {
		return fmt.Errorf("failed to attach kill exec: %w", err)
	}
	defer attachResp.Close()

	_, _ = io.Copy(io.Discard, attachResp.Reader)

	return nil
}
