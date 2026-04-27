package tools

import "sync"

type trackedExec struct {
	containerID string
	pid         int
}

type runningExecTracker struct {
	mx   sync.Mutex
	exec map[string]trackedExec
}

var execTracker = &runningExecTracker{
	exec: make(map[string]trackedExec),
}

func (t *runningExecTracker) register(execID, containerID string, pid int) {
	if execID == "" || containerID == "" || pid <= 0 {
		return
	}

	t.mx.Lock()
	defer t.mx.Unlock()
	t.exec[execID] = trackedExec{
		containerID: containerID,
		pid:         pid,
	}
}

func (t *runningExecTracker) unregister(execID string) {
	if execID == "" {
		return
	}

	t.mx.Lock()
	defer t.mx.Unlock()
	delete(t.exec, execID)
}

func (t *runningExecTracker) listPIDs(containerID string) []int {
	if containerID == "" {
		return nil
	}

	t.mx.Lock()
	defer t.mx.Unlock()

	seen := make(map[int]struct{})
	pids := make([]int, 0)
	for _, e := range t.exec {
		if e.containerID != containerID || e.pid <= 0 {
			continue
		}
		if _, ok := seen[e.pid]; ok {
			continue
		}
		seen[e.pid] = struct{}{}
		pids = append(pids, e.pid)
	}

	return pids
}
