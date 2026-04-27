package checker

// IsReadyToContinue returns true when all required baseline conditions are met.
func (c *CheckResult) IsReadyToContinue() bool {
	return c.EnvFileExists &&
		c.EnvDirWritable &&
		c.DockerApiAccessible &&
		c.WorkerEnvApiAccessible &&
		c.DockerComposeInstalled &&
		c.DockerVersionOK &&
		c.DockerComposeVersionOK &&
		c.SysNetworkOK &&
		c.SysCPUOK &&
		c.SysMemoryOK &&
		c.SysDiskFreeSpaceOK
}

// availability helpers for installer operations
// these functions centralize complex visibility/availability logic for UI

// CanStartAll returns true when at least one embedded stack is installed and not running
func (c *CheckResult) CanStartAll() bool {
	if c.PentagiInstalled && !c.PentagiRunning {
		return true
	}
	if c.GraphitiConnected && !c.GraphitiExternal && c.GraphitiInstalled && !c.GraphitiRunning {
		return true
	}
	if c.LangfuseConnected && !c.LangfuseExternal && c.LangfuseInstalled && !c.LangfuseRunning {
		return true
	}
	if c.ObservabilityConnected && !c.ObservabilityExternal && c.ObservabilityInstalled && !c.ObservabilityRunning {
		return true
	}
	return false
}

// CanStopAll returns true when any compose stack is running
func (c *CheckResult) CanStopAll() bool {
	return c.PentagiRunning || c.GraphitiRunning || c.LangfuseRunning || c.ObservabilityRunning
}

// CanRestartAll mirrors stop logic (requires running services)
func (c *CheckResult) CanRestartAll() bool { return c.CanStopAll() }

// CanDownloadWorker returns true when worker image is missing
func (c *CheckResult) CanDownloadWorker() bool { return !c.WorkerImageExists }

// CanUpdateWorker returns true when worker image exists but is not up to date
func (c *CheckResult) CanUpdateWorker() bool { return c.WorkerImageExists && !c.WorkerIsUpToDate }

// CanUpdateAll returns true when any installed stack has updates available
func (c *CheckResult) CanUpdateAll() bool {
	if c.PentagiInstalled && !c.PentagiIsUpToDate {
		return true
	}
	if c.GraphitiInstalled && !c.GraphitiIsUpToDate {
		return true
	}
	if c.LangfuseInstalled && !c.LangfuseIsUpToDate {
		return true
	}
	if c.ObservabilityInstalled && !c.ObservabilityIsUpToDate {
		return true
	}
	return false
}

// CanUpdateInstaller returns true when installer update is available and update server accessible
func (c *CheckResult) CanUpdateInstaller() bool {
	return !c.InstallerIsUpToDate && c.UpdateServerAccessible
}

// CanFactoryReset returns true when any compose stack is installed
func (c *CheckResult) CanFactoryReset() bool {
	return c.PentagiInstalled || c.GraphitiInstalled || c.LangfuseInstalled || c.ObservabilityInstalled
}

// CanRemoveAll returns true when any compose stack is installed
func (c *CheckResult) CanRemoveAll() bool { return c.CanFactoryReset() }

// CanPurgeAll returns true when any compose stack is installed
func (c *CheckResult) CanPurgeAll() bool { return c.CanFactoryReset() }

// CanResetPassword returns true when PentAGI is running
func (c *CheckResult) CanResetPassword() bool { return c.PentagiRunning }

// CanInstallAll returns true when main stack is not installed yet
func (c *CheckResult) CanInstallAll() bool { return !c.PentagiInstalled }
