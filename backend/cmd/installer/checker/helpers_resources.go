package checker

import (
	"bufio"
	"context"
	"os"
	"os/exec"
	"regexp"
	"runtime"
	"strconv"
	"strings"
)

func checkCPUResources() bool {
	return runtime.NumCPU() >= 2
}

func determineComponentNeeds(c *CheckResult) (needsForPentagi, needsForGraphiti, needsForLangfuse, needsForObservability bool) {
	needsForPentagi = !c.PentagiRunning
	needsForGraphiti = c.GraphitiConnected && !c.GraphitiExternal && !c.GraphitiRunning
	needsForLangfuse = c.LangfuseConnected && !c.LangfuseExternal && !c.LangfuseRunning
	needsForObservability = c.ObservabilityConnected && !c.ObservabilityExternal && !c.ObservabilityRunning
	return
}

func calculateRequiredMemoryGB(needsForPentagi, needsForGraphiti, needsForLangfuse, needsForObservability bool) float64 {
	requiredGB := MinFreeMemGB
	if needsForPentagi {
		requiredGB += MinFreeMemGBForPentagi
	}
	if needsForGraphiti {
		requiredGB += MinFreeMemGBForGraphiti
	}
	if needsForLangfuse {
		requiredGB += MinFreeMemGBForLangfuse
	}
	if needsForObservability {
		requiredGB += MinFreeMemGBForObservability
	}
	return requiredGB
}

func checkMemoryResources(needsForPentagi, needsForGraphiti, needsForLangfuse, needsForObservability bool) bool {
	if !needsForPentagi && !needsForGraphiti && !needsForLangfuse && !needsForObservability {
		return true
	}

	requiredGB := calculateRequiredMemoryGB(needsForPentagi, needsForGraphiti, needsForLangfuse, needsForObservability)

	switch runtime.GOOS {
	case "linux":
		return checkLinuxMemory(requiredGB)
	case "darwin":
		return checkDarwinMemory(requiredGB)
	default:
		return true
	}
}

func getAvailableMemoryGB() float64 {
	switch runtime.GOOS {
	case "linux":
		return getLinuxAvailableMemoryGB()
	case "darwin":
		return getDarwinAvailableMemoryGB()
	default:
		return 0.0
	}
}

func getLinuxAvailableMemoryGB() float64 {
	file, err := os.Open("/proc/meminfo")
	if err != nil {
		return 0.0
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var memFree, memAvailable int64

	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "MemAvailable:") {
			fields := strings.Fields(line)
			if len(fields) >= 2 {
				if val, err := strconv.ParseInt(fields[1], 10, 64); err == nil {
					memAvailable = val * 1024
				}
			}
			break
		}
		if strings.HasPrefix(line, "MemFree:") {
			fields := strings.Fields(line)
			if len(fields) >= 2 {
				if val, err := strconv.ParseInt(fields[1], 10, 64); err == nil {
					memFree = val * 1024
				}
			}
		}
	}

	availableMemGB := float64(memAvailable) / (1024 * 1024 * 1024)
	if availableMemGB > 0 {
		return availableMemGB
	}

	return float64(memFree) / (1024 * 1024 * 1024)
}

func getDarwinAvailableMemoryGB() float64 {
	cmd := exec.Command("vm_stat")
	output, err := cmd.Output()
	if err != nil {
		return 0.0
	}

	lines := strings.Split(string(output), "\n")
	var pageSize, freePages, inactivePages, purgeablePages int64 = 4096, 0, 0, 0

	for _, line := range lines {
		if strings.Contains(line, "page size of") {
			re := regexp.MustCompile(`(\d+) bytes`)
			if matches := re.FindStringSubmatch(line); len(matches) > 1 {
				if val, err := strconv.ParseInt(matches[1], 10, 64); err == nil {
					pageSize = val
				}
			}
		}
		if strings.HasPrefix(line, "Pages free:") {
			re := regexp.MustCompile(`(\d+)`)
			if matches := re.FindStringSubmatch(line); len(matches) > 1 {
				if val, err := strconv.ParseInt(matches[1], 10, 64); err == nil {
					freePages = val
				}
			}
		}
		if strings.HasPrefix(line, "Pages inactive:") {
			re := regexp.MustCompile(`(\d+)`)
			if matches := re.FindStringSubmatch(line); len(matches) > 1 {
				if val, err := strconv.ParseInt(matches[1], 10, 64); err == nil {
					inactivePages = val
				}
			}
		}
		if strings.HasPrefix(line, "Pages purgeable:") {
			re := regexp.MustCompile(`(\d+)`)
			if matches := re.FindStringSubmatch(line); len(matches) > 1 {
				if val, err := strconv.ParseInt(matches[1], 10, 64); err == nil {
					purgeablePages = val
				}
			}
		}
	}

	availablePages := freePages + inactivePages + purgeablePages
	return float64(availablePages*pageSize) / (1024 * 1024 * 1024)
}

func checkLinuxMemory(requiredGB float64) bool {
	file, err := os.Open("/proc/meminfo")
	if err != nil {
		return true
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var memFree, memAvailable int64

	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "MemAvailable:") {
			fields := strings.Fields(line)
			if len(fields) >= 2 {
				if val, err := strconv.ParseInt(fields[1], 10, 64); err == nil {
					memAvailable = val * 1024
				}
			}
			break
		}
		if strings.HasPrefix(line, "MemFree:") {
			fields := strings.Fields(line)
			if len(fields) >= 2 {
				if val, err := strconv.ParseInt(fields[1], 10, 64); err == nil {
					memFree = val * 1024
				}
			}
		}
	}

	availableMemGB := float64(memAvailable) / (1024 * 1024 * 1024)
	if availableMemGB > 0 {
		return availableMemGB >= requiredGB
	}

	freeMemGB := float64(memFree) / (1024 * 1024 * 1024)
	return freeMemGB >= requiredGB
}

func checkDarwinMemory(requiredGB float64) bool {
	cmd := exec.Command("vm_stat")
	output, err := cmd.Output()
	if err != nil {
		return true
	}

	lines := strings.Split(string(output), "\n")
	var pageSize, freePages, inactivePages, purgeablePages int64 = 4096, 0, 0, 0

	for _, line := range lines {
		if strings.Contains(line, "page size of") {
			re := regexp.MustCompile(`(\d+) bytes`)
			if matches := re.FindStringSubmatch(line); len(matches) > 1 {
				if val, err := strconv.ParseInt(matches[1], 10, 64); err == nil {
					pageSize = val
				}
			}
		}
		if strings.HasPrefix(line, "Pages free:") {
			re := regexp.MustCompile(`(\d+)`)
			if matches := re.FindStringSubmatch(line); len(matches) > 1 {
				if val, err := strconv.ParseInt(matches[1], 10, 64); err == nil {
					freePages = val
				}
			}
		}
		if strings.HasPrefix(line, "Pages inactive:") {
			re := regexp.MustCompile(`(\d+)`)
			if matches := re.FindStringSubmatch(line); len(matches) > 1 {
				if val, err := strconv.ParseInt(matches[1], 10, 64); err == nil {
					inactivePages = val
				}
			}
		}
		if strings.HasPrefix(line, "Pages purgeable:") {
			re := regexp.MustCompile(`(\d+)`)
			if matches := re.FindStringSubmatch(line); len(matches) > 1 {
				if val, err := strconv.ParseInt(matches[1], 10, 64); err == nil {
					purgeablePages = val
				}
			}
		}
	}

	availablePages := freePages + inactivePages + purgeablePages
	availableMemGB := float64(availablePages*pageSize) / (1024 * 1024 * 1024)
	return availableMemGB >= requiredGB
}

func calculateRequiredDiskGB(workerImageExists bool, localComponents int) float64 {
	if !workerImageExists {
		return MinFreeDiskGBForWorkerImages
	} else if localComponents > 0 {
		return MinFreeDiskGBForComponents + float64(localComponents)*MinFreeDiskGBPerComponents
	}
	return MinFreeDiskGB
}

func countLocalComponentsToInstall(
	pentagiInstalled,
	graphitiConnected, graphitiExternal, graphitiInstalled,
	langfuseConnected, langfuseExternal, langfuseInstalled,
	obsConnected, obsExternal, obsInstalled bool,
) int {
	localComponents := 0
	if !pentagiInstalled {
		localComponents++
	}
	if graphitiConnected && !graphitiExternal && !graphitiInstalled {
		localComponents++
	}
	if langfuseConnected && !langfuseExternal && !langfuseInstalled {
		localComponents++
	}
	if obsConnected && !obsExternal && !obsInstalled {
		localComponents++
	}
	return localComponents
}

func checkDiskSpaceWithContext(
	ctx context.Context,
	workerImageExists, pentagiInstalled,
	graphitiConnected, graphitiExternal, graphitiInstalled,
	langfuseConnected, langfuseExternal, langfuseInstalled,
	obsConnected, obsExternal, obsInstalled bool,
) bool {
	localComponents := countLocalComponentsToInstall(
		pentagiInstalled,
		graphitiConnected, graphitiExternal, graphitiInstalled,
		langfuseConnected, langfuseExternal, langfuseInstalled,
		obsConnected, obsExternal, obsInstalled,
	)

	requiredGB := calculateRequiredDiskGB(workerImageExists, localComponents)

	switch runtime.GOOS {
	case "linux":
		return checkLinuxDiskSpace(ctx, requiredGB)
	case "darwin":
		return checkDarwinDiskSpace(ctx, requiredGB)
	default:
		return true
	}
}

func getAvailableDiskGB(ctx context.Context) float64 {
	switch runtime.GOOS {
	case "linux":
		return getLinuxAvailableDiskGB(ctx)
	case "darwin":
		return getDarwinAvailableDiskGB(ctx)
	default:
		return 0.0
	}
}

func getLinuxAvailableDiskGB(ctx context.Context) float64 {
	cmd := exec.CommandContext(ctx, "df", "-BG", ".")
	output, err := cmd.Output()
	if err != nil {
		return 0.0
	}

	lines := strings.Split(string(output), "\n")
	if len(lines) < 2 {
		return 0.0
	}

	fields := strings.Fields(lines[1])
	if len(fields) < 4 {
		return 0.0
	}

	availableStr := strings.TrimSuffix(fields[3], "G")
	if available, err := strconv.ParseFloat(availableStr, 64); err == nil {
		return available
	}

	return 0.0
}

func getDarwinAvailableDiskGB(ctx context.Context) float64 {
	cmd := exec.CommandContext(ctx, "df", "-g", ".")
	output, err := cmd.Output()
	if err != nil {
		return 0.0
	}

	lines := strings.Split(string(output), "\n")
	if len(lines) < 2 {
		return 0.0
	}

	fields := strings.Fields(lines[1])
	if len(fields) < 4 {
		return 0.0
	}

	if available, err := strconv.ParseFloat(fields[3], 64); err == nil {
		return available
	}

	return 0.0
}

func checkLinuxDiskSpace(ctx context.Context, requiredGB float64) bool {
	cmd := exec.CommandContext(ctx, "df", "-BG", ".")
	output, err := cmd.Output()
	if err != nil {
		return true
	}

	lines := strings.Split(string(output), "\n")
	if len(lines) < 2 {
		return true
	}

	fields := strings.Fields(lines[1])
	if len(fields) < 4 {
		return true
	}

	availableStr := strings.TrimSuffix(fields[3], "G")
	if available, err := strconv.ParseFloat(availableStr, 64); err == nil {
		return available >= requiredGB
	}

	return true
}

func checkDarwinDiskSpace(ctx context.Context, requiredGB float64) bool {
	cmd := exec.CommandContext(ctx, "df", "-g", ".")
	output, err := cmd.Output()
	if err != nil {
		return true
	}

	lines := strings.Split(string(output), "\n")
	if len(lines) < 2 {
		return true
	}

	fields := strings.Fields(lines[1])
	if len(fields) < 4 {
		return true
	}

	if available, err := strconv.ParseFloat(fields[3], 64); err == nil {
		return available >= requiredGB
	}

	return true
}
