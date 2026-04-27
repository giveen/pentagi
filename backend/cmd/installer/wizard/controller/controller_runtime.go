package controller

import "fmt"

import "pentagi/cmd/installer/loader"

func (c *controller) GetSummarizerConfig(summarizerType SummarizerType) *SummarizerConfig {
	var prefix string
	if summarizerType == SummarizerTypeAssistant {
		prefix = "ASSISTANT_SUMMARIZER_"
	} else {
		prefix = "SUMMARIZER_"
	}

	config := &SummarizerConfig{
		Type: summarizerType,
	}

	// Read variables directly from state (defaults handled by loader)
	config.PreserveLast, _ = c.GetVar(prefix + "PRESERVE_LAST")

	if summarizerType == SummarizerTypeGeneral {
		config.UseQA, _ = c.GetVar(prefix + "USE_QA")
		config.SumHumanInQA, _ = c.GetVar(prefix + "SUM_MSG_HUMAN_IN_QA")
	}

	// Size settings
	config.LastSecBytes, _ = c.GetVar(prefix + "LAST_SEC_BYTES")
	config.MaxBPBytes, _ = c.GetVar(prefix + "MAX_BP_BYTES")
	config.MaxQABytes, _ = c.GetVar(prefix + "MAX_QA_BYTES")

	// Count settings
	config.MaxQASections, _ = c.GetVar(prefix + "MAX_QA_SECTIONS")
	config.KeepQASections, _ = c.GetVar(prefix + "KEEP_QA_SECTIONS")

	return config
}

// UpdateSummarizerConfig updates summarizer configuration
func (c *controller) UpdateSummarizerConfig(config *SummarizerConfig) error {
	var prefix string
	if config.Type == SummarizerTypeAssistant {
		prefix = "ASSISTANT_SUMMARIZER_"
	} else {
		prefix = "SUMMARIZER_"
	}

	// Update boolean settings
	if err := c.SetVar(prefix+"PRESERVE_LAST", config.PreserveLast.Value); err != nil {
		return fmt.Errorf("failed to set %s: %w", prefix+"PRESERVE_LAST", err)
	}

	// General-specific boolean settings
	if config.Type == SummarizerTypeGeneral {
		if err := c.SetVar(prefix+"USE_QA", config.UseQA.Value); err != nil {
			return fmt.Errorf("failed to set %s: %w", prefix+"USE_QA", err)
		}
		if err := c.SetVar(prefix+"SUM_MSG_HUMAN_IN_QA", config.SumHumanInQA.Value); err != nil {
			return fmt.Errorf("failed to set %s: %w", prefix+"SUM_MSG_HUMAN_IN_QA", err)
		}
	}

	// Update size settings
	if err := c.SetVar(prefix+"LAST_SEC_BYTES", config.LastSecBytes.Value); err != nil {
		return fmt.Errorf("failed to set %s: %w", prefix+"LAST_SEC_BYTES", err)
	}
	if err := c.SetVar(prefix+"MAX_BP_BYTES", config.MaxBPBytes.Value); err != nil {
		return fmt.Errorf("failed to set %s: %w", prefix+"MAX_BP_BYTES", err)
	}
	if err := c.SetVar(prefix+"MAX_QA_BYTES", config.MaxQABytes.Value); err != nil {
		return fmt.Errorf("failed to set %s: %w", prefix+"MAX_QA_BYTES", err)
	}

	// Update count settings
	if err := c.SetVar(prefix+"MAX_QA_SECTIONS", config.MaxQASections.Value); err != nil {
		return fmt.Errorf("failed to set %s: %w", prefix+"MAX_QA_SECTIONS", err)
	}
	if err := c.SetVar(prefix+"KEEP_QA_SECTIONS", config.KeepQASections.Value); err != nil {
		return fmt.Errorf("failed to set %s: %w", prefix+"KEEP_QA_SECTIONS", err)
	}

	return nil
}

func (c *controller) ResetSummarizerConfig(summarizerType SummarizerType) *SummarizerConfig {
	var prefix string
	if summarizerType == SummarizerTypeAssistant {
		prefix = "ASSISTANT_SUMMARIZER_"
	} else {
		prefix = "SUMMARIZER_"
	}

	vars := []string{
		prefix + "PRESERVE_LAST",
		prefix + "LAST_SEC_BYTES",
		prefix + "MAX_BP_BYTES",
		prefix + "MAX_QA_BYTES",
		prefix + "MAX_QA_SECTIONS",
		prefix + "KEEP_QA_SECTIONS",
	}

	if summarizerType == SummarizerTypeGeneral {
		vars = append(vars,
			prefix+"USE_QA",
			prefix+"SUM_MSG_HUMAN_IN_QA",
		)
	}

	if err := c.ResetVars(vars); err != nil {
		return nil
	}

	return c.GetSummarizerConfig(summarizerType)
}

// EmbedderConfig represents embedder configuration settings
type EmbedderConfig struct {
	// direct form field mappings using loader.EnvVar
	// these fields directly correspond to environment variables and form inputs (not computed)
	Provider      loader.EnvVar // EMBEDDING_PROVIDER
	URL           loader.EnvVar // EMBEDDING_URL
	APIKey        loader.EnvVar // EMBEDDING_KEY
	Model         loader.EnvVar // EMBEDDING_MODEL
	BatchSize     loader.EnvVar // EMBEDDING_BATCH_SIZE
	StripNewLines loader.EnvVar // EMBEDDING_STRIP_NEW_LINES

	// computed fields (not directly mapped to env vars)
	Configured bool
	Installed  bool
}

// GetEmbedderConfig returns current embedder configuration
func (c *controller) GetEmbedderConfig() *EmbedderConfig {
	config := &EmbedderConfig{}
	config.Provider, _ = c.GetVar("EMBEDDING_PROVIDER")
	config.URL, _ = c.GetVar("EMBEDDING_URL")
	config.APIKey, _ = c.GetVar("EMBEDDING_KEY")
	config.Model, _ = c.GetVar("EMBEDDING_MODEL")
	config.BatchSize, _ = c.GetVar("EMBEDDING_BATCH_SIZE")
	config.StripNewLines, _ = c.GetVar("EMBEDDING_STRIP_NEW_LINES")
	config.Installed = c.checker.PentagiInstalled

	// Determine if configured based on provider requirements
	switch config.Provider.Value {
	case "openai", "":
		// For OpenAI, check if we have API key either in EMBEDDING_KEY or OPEN_AI_KEY
		openaiKey, _ := c.GetVar("OPEN_AI_KEY")
		config.Configured = config.APIKey.Value != "" || openaiKey.Value != ""
	case "ollama":
		// for Ollama, no API key required, but URL must be provided
		config.Configured = config.URL.Value != ""
	case "huggingface", "googleai":
		// These require API key
		config.Configured = config.APIKey.Value != ""
	default:
		// Others are configured if API key is present
		config.Configured = config.APIKey.Value != ""
	}

	return config
}

// UpdateEmbedderConfig updates embedder configuration
func (c *controller) UpdateEmbedderConfig(config *EmbedderConfig) error {
	if config == nil {
		return fmt.Errorf("config cannot be nil")
	}

	// Update environment variables
	if err := c.SetVar("EMBEDDING_PROVIDER", config.Provider.Value); err != nil {
		return fmt.Errorf("failed to set EMBEDDING_PROVIDER: %w", err)
	}
	if err := c.SetVar("EMBEDDING_URL", config.URL.Value); err != nil {
		return fmt.Errorf("failed to set EMBEDDING_URL: %w", err)
	}
	if err := c.SetVar("EMBEDDING_KEY", config.APIKey.Value); err != nil {
		return fmt.Errorf("failed to set EMBEDDING_KEY: %w", err)
	}
	if err := c.SetVar("EMBEDDING_MODEL", config.Model.Value); err != nil {
		return fmt.Errorf("failed to set EMBEDDING_MODEL: %w", err)
	}
	if err := c.SetVar("EMBEDDING_BATCH_SIZE", config.BatchSize.Value); err != nil {
		return fmt.Errorf("failed to set EMBEDDING_BATCH_SIZE: %w", err)
	}
	if err := c.SetVar("EMBEDDING_STRIP_NEW_LINES", config.StripNewLines.Value); err != nil {
		return fmt.Errorf("failed to set EMBEDDING_STRIP_NEW_LINES: %w", err)
	}

	return nil
}

func (c *controller) ResetEmbedderConfig() *EmbedderConfig {
	vars := []string{
		"EMBEDDING_PROVIDER",
		"EMBEDDING_URL",
		"EMBEDDING_KEY",
		"EMBEDDING_MODEL",
		"EMBEDDING_BATCH_SIZE",
		"EMBEDDING_STRIP_NEW_LINES",
	}

	if err := c.ResetVars(vars); err != nil {
		return nil
	}

	return c.GetEmbedderConfig()
}

// AIAgentsConfig represents extra AI agents configuration
type AIAgentsConfig struct {
	// direct form field mappings using loader.EnvVar
	// these fields directly correspond to environment variables and form inputs (not computed)
	HumanInTheLoop                 loader.EnvVar // ASK_USER
	AssistantUseAgents             loader.EnvVar // ASSISTANT_USE_AGENTS
	ExecutionMonitorEnabled        loader.EnvVar // EXECUTION_MONITOR_ENABLED
	ExecutionMonitorSameToolLimit  loader.EnvVar // EXECUTION_MONITOR_SAME_TOOL_LIMIT
	ExecutionMonitorTotalToolLimit loader.EnvVar // EXECUTION_MONITOR_TOTAL_TOOL_LIMIT
	MaxGeneralAgentToolCalls       loader.EnvVar // MAX_GENERAL_AGENT_TOOL_CALLS
	MaxLimitedAgentToolCalls       loader.EnvVar // MAX_LIMITED_AGENT_TOOL_CALLS
	AgentPlanningStepEnabled       loader.EnvVar // AGENT_PLANNING_STEP_ENABLED
}

func (c *controller) GetAIAgentsConfig() *AIAgentsConfig {
	config := &AIAgentsConfig{}

	config.HumanInTheLoop, _ = c.GetVar("ASK_USER")
	config.AssistantUseAgents, _ = c.GetVar("ASSISTANT_USE_AGENTS")
	config.ExecutionMonitorEnabled, _ = c.GetVar("EXECUTION_MONITOR_ENABLED")
	config.ExecutionMonitorSameToolLimit, _ = c.GetVar("EXECUTION_MONITOR_SAME_TOOL_LIMIT")
	config.ExecutionMonitorTotalToolLimit, _ = c.GetVar("EXECUTION_MONITOR_TOTAL_TOOL_LIMIT")
	config.MaxGeneralAgentToolCalls, _ = c.GetVar("MAX_GENERAL_AGENT_TOOL_CALLS")
	config.MaxLimitedAgentToolCalls, _ = c.GetVar("MAX_LIMITED_AGENT_TOOL_CALLS")
	config.AgentPlanningStepEnabled, _ = c.GetVar("AGENT_PLANNING_STEP_ENABLED")

	return config
}

func (c *controller) UpdateAIAgentsConfig(config *AIAgentsConfig) error {
	if config == nil {
		return fmt.Errorf("config cannot be nil")
	}

	if err := c.SetVar("ASK_USER", config.HumanInTheLoop.Value); err != nil {
		return fmt.Errorf("failed to set ASK_USER: %w", err)
	}
	if err := c.SetVar("ASSISTANT_USE_AGENTS", config.AssistantUseAgents.Value); err != nil {
		return fmt.Errorf("failed to set ASSISTANT_USE_AGENTS: %w", err)
	}
	if err := c.SetVar("EXECUTION_MONITOR_ENABLED", config.ExecutionMonitorEnabled.Value); err != nil {
		return fmt.Errorf("failed to set EXECUTION_MONITOR_ENABLED: %w", err)
	}
	if err := c.SetVar("EXECUTION_MONITOR_SAME_TOOL_LIMIT", config.ExecutionMonitorSameToolLimit.Value); err != nil {
		return fmt.Errorf("failed to set EXECUTION_MONITOR_SAME_TOOL_LIMIT: %w", err)
	}
	if err := c.SetVar("EXECUTION_MONITOR_TOTAL_TOOL_LIMIT", config.ExecutionMonitorTotalToolLimit.Value); err != nil {
		return fmt.Errorf("failed to set EXECUTION_MONITOR_TOTAL_TOOL_LIMIT: %w", err)
	}
	if err := c.SetVar("MAX_GENERAL_AGENT_TOOL_CALLS", config.MaxGeneralAgentToolCalls.Value); err != nil {
		return fmt.Errorf("failed to set MAX_GENERAL_AGENT_TOOL_CALLS: %w", err)
	}
	if err := c.SetVar("MAX_LIMITED_AGENT_TOOL_CALLS", config.MaxLimitedAgentToolCalls.Value); err != nil {
		return fmt.Errorf("failed to set MAX_LIMITED_AGENT_TOOL_CALLS: %w", err)
	}
	if err := c.SetVar("AGENT_PLANNING_STEP_ENABLED", config.AgentPlanningStepEnabled.Value); err != nil {
		return fmt.Errorf("failed to set AGENT_PLANNING_STEP_ENABLED: %w", err)
	}

	return nil
}

func (c *controller) ResetAIAgentsConfig() *AIAgentsConfig {
	vars := []string{
		"ASK_USER",
		"ASSISTANT_USE_AGENTS",
	}

	if err := c.ResetVars(vars); err != nil {
		return nil
	}

	return c.GetAIAgentsConfig()
}
