package controller

import (
	"fmt"
	"slices"
	"sort"
	"strings"

	"pentagi/cmd/installer/files"
)

func GetEmbeddedLLMConfigsPath(files files.Files) []string {
	providersConfigsPath := make([]string, 0)
	if confFiles, err := files.List(EmbeddedLLMConfigsPath); err == nil {
		for _, confFile := range confFiles {
			confPath := DefaultLLMConfigsPath + strings.TrimPrefix(confFile, EmbeddedLLMConfigsPath+"/")
			providersConfigsPath = append(providersConfigsPath, confPath)
		}
		sort.Strings(providersConfigsPath)
	}

	return providersConfigsPath
}

// GetLLMProviders returns configured LLM providers.
func (c *controller) GetLLMProviders() map[string]*LLMProviderConfig {
	return map[string]*LLMProviderConfig{
		"openai":    c.GetLLMProviderConfig("openai"),
		"anthropic": c.GetLLMProviderConfig("anthropic"),
		"gemini":    c.GetLLMProviderConfig("gemini"),
		"bedrock":   c.GetLLMProviderConfig("bedrock"),
		"ollama":    c.GetLLMProviderConfig("ollama"),
		"deepseek":  c.GetLLMProviderConfig("deepseek"),
		"glm":       c.GetLLMProviderConfig("glm"),
		"kimi":      c.GetLLMProviderConfig("kimi"),
		"qwen":      c.GetLLMProviderConfig("qwen"),
		"custom":    c.GetLLMProviderConfig("custom"),
	}
}

// GetLLMProviderConfig returns the current LLM provider configuration.
func (c *controller) GetLLMProviderConfig(providerID string) *LLMProviderConfig {
	providersConfigsPath := GetEmbeddedLLMConfigsPath(c.files)
	providerConfig := &LLMProviderConfig{
		Name:                   "Unknown",
		EmbeddedLLMConfigsPath: providersConfigsPath,
	}

	switch providerID {
	case "openai":
		providerConfig.Name = "OpenAI"
		providerConfig.APIKey, _ = c.GetVar("OPEN_AI_KEY")
		providerConfig.BaseURL, _ = c.GetVar("OPEN_AI_SERVER_URL")
		providerConfig.Configured = providerConfig.APIKey.Value != ""

	case "anthropic":
		providerConfig.Name = "Anthropic"
		providerConfig.APIKey, _ = c.GetVar("ANTHROPIC_API_KEY")
		providerConfig.BaseURL, _ = c.GetVar("ANTHROPIC_SERVER_URL")
		providerConfig.Configured = providerConfig.APIKey.Value != ""

	case "gemini":
		providerConfig.Name = "Google Gemini"
		providerConfig.APIKey, _ = c.GetVar("GEMINI_API_KEY")
		providerConfig.BaseURL, _ = c.GetVar("GEMINI_SERVER_URL")
		providerConfig.Configured = providerConfig.APIKey.Value != ""

	case "bedrock":
		providerConfig.Name = "AWS Bedrock"
		providerConfig.Region, _ = c.GetVar("BEDROCK_REGION")
		providerConfig.DefaultAuth, _ = c.GetVar("BEDROCK_DEFAULT_AUTH")
		providerConfig.BearerToken, _ = c.GetVar("BEDROCK_BEARER_TOKEN")
		providerConfig.AccessKey, _ = c.GetVar("BEDROCK_ACCESS_KEY_ID")
		providerConfig.SecretKey, _ = c.GetVar("BEDROCK_SECRET_ACCESS_KEY")
		providerConfig.SessionToken, _ = c.GetVar("BEDROCK_SESSION_TOKEN")
		providerConfig.BaseURL, _ = c.GetVar("BEDROCK_SERVER_URL")
		providerConfig.Configured = providerConfig.DefaultAuth.Value == "true" ||
			providerConfig.BearerToken.Value != "" ||
			(providerConfig.AccessKey.Value != "" && providerConfig.SecretKey.Value != "")

	case "ollama":
		providerConfig.Name = "Ollama"
		providerConfig.BaseURL, _ = c.GetVar("OLLAMA_SERVER_URL")
		providerConfig.APIKey, _ = c.GetVar("OLLAMA_SERVER_API_KEY")
		providerConfig.ConfigPath, _ = c.GetVar("OLLAMA_SERVER_CONFIG_PATH")
		providerConfig.HostConfigPath, _ = c.GetVar("PENTAGI_OLLAMA_SERVER_CONFIG_PATH")
		if slices.Contains(providersConfigsPath, providerConfig.ConfigPath.Value) {
			providerConfig.HostConfigPath.Value = providerConfig.ConfigPath.Value
		}
		providerConfig.Model, _ = c.GetVar("OLLAMA_SERVER_MODEL")
		providerConfig.PullTimeout, _ = c.GetVar("OLLAMA_SERVER_PULL_MODELS_TIMEOUT")
		providerConfig.PullEnabled, _ = c.GetVar("OLLAMA_SERVER_PULL_MODELS_ENABLED")
		providerConfig.LoadModelsEnabled, _ = c.GetVar("OLLAMA_SERVER_LOAD_MODELS_ENABLED")
		providerConfig.Configured = providerConfig.BaseURL.Value != ""

	case "deepseek":
		providerConfig.Name = "DeepSeek"
		providerConfig.APIKey, _ = c.GetVar("DEEPSEEK_API_KEY")
		providerConfig.BaseURL, _ = c.GetVar("DEEPSEEK_SERVER_URL")
		providerConfig.ProviderName, _ = c.GetVar("DEEPSEEK_PROVIDER")
		providerConfig.Configured = providerConfig.APIKey.Value != ""

	case "glm":
		providerConfig.Name = "GLM"
		providerConfig.APIKey, _ = c.GetVar("GLM_API_KEY")
		providerConfig.BaseURL, _ = c.GetVar("GLM_SERVER_URL")
		providerConfig.ProviderName, _ = c.GetVar("GLM_PROVIDER")
		providerConfig.Configured = providerConfig.APIKey.Value != ""

	case "kimi":
		providerConfig.Name = "Kimi"
		providerConfig.APIKey, _ = c.GetVar("KIMI_API_KEY")
		providerConfig.BaseURL, _ = c.GetVar("KIMI_SERVER_URL")
		providerConfig.ProviderName, _ = c.GetVar("KIMI_PROVIDER")
		providerConfig.Configured = providerConfig.APIKey.Value != ""

	case "qwen":
		providerConfig.Name = "Qwen"
		providerConfig.APIKey, _ = c.GetVar("QWEN_API_KEY")
		providerConfig.BaseURL, _ = c.GetVar("QWEN_SERVER_URL")
		providerConfig.ProviderName, _ = c.GetVar("QWEN_PROVIDER")
		providerConfig.Configured = providerConfig.APIKey.Value != ""

	case "custom":
		providerConfig.Name = "Custom"
		providerConfig.BaseURL, _ = c.GetVar("LLM_SERVER_URL")
		providerConfig.APIKey, _ = c.GetVar("LLM_SERVER_KEY")
		providerConfig.Model, _ = c.GetVar("LLM_SERVER_MODEL")
		providerConfig.ConfigPath, _ = c.GetVar("LLM_SERVER_CONFIG_PATH")
		providerConfig.HostConfigPath, _ = c.GetVar("PENTAGI_LLM_SERVER_CONFIG_PATH")
		if slices.Contains(providersConfigsPath, providerConfig.ConfigPath.Value) {
			providerConfig.HostConfigPath.Value = providerConfig.ConfigPath.Value
		}
		providerConfig.LegacyReasoning, _ = c.GetVar("LLM_SERVER_LEGACY_REASONING")
		providerConfig.PreserveReasoning, _ = c.GetVar("LLM_SERVER_PRESERVE_REASONING")
		providerConfig.ProviderName, _ = c.GetVar("LLM_SERVER_PROVIDER")
		providerConfig.Configured = providerConfig.BaseURL.Value != "" && providerConfig.APIKey.Value != "" &&
			(providerConfig.Model.Value != "" || providerConfig.ConfigPath.Value != "")
	}

	return providerConfig
}

// UpdateLLMProviderConfig updates a specific LLM provider configuration.
func (c *controller) UpdateLLMProviderConfig(providerID string, config *LLMProviderConfig) error {
	switch providerID {
	case "openai", "anthropic", "gemini":
		if err := c.SetVar(config.APIKey.Name, config.APIKey.Value); err != nil {
			return fmt.Errorf("failed to set %s: %w", config.APIKey.Name, err)
		}
		if err := c.SetVar(config.BaseURL.Name, config.BaseURL.Value); err != nil {
			return fmt.Errorf("failed to set %s: %w", config.BaseURL.Name, err)
		}

	case "bedrock":
		if err := c.SetVar(config.Region.Name, config.Region.Value); err != nil {
			return fmt.Errorf("failed to set %s: %w", config.Region.Name, err)
		}
		if err := c.SetVar(config.DefaultAuth.Name, config.DefaultAuth.Value); err != nil {
			return fmt.Errorf("failed to set %s: %w", config.DefaultAuth.Name, err)
		}
		if err := c.SetVar(config.BearerToken.Name, config.BearerToken.Value); err != nil {
			return fmt.Errorf("failed to set %s: %w", config.BearerToken.Name, err)
		}
		if err := c.SetVar(config.AccessKey.Name, config.AccessKey.Value); err != nil {
			return fmt.Errorf("failed to set %s: %w", config.AccessKey.Name, err)
		}
		if err := c.SetVar(config.SecretKey.Name, config.SecretKey.Value); err != nil {
			return fmt.Errorf("failed to set %s: %w", config.SecretKey.Name, err)
		}
		if err := c.SetVar(config.SessionToken.Name, config.SessionToken.Value); err != nil {
			return fmt.Errorf("failed to set %s: %w", config.SessionToken.Name, err)
		}
		if err := c.SetVar(config.BaseURL.Name, config.BaseURL.Value); err != nil {
			return fmt.Errorf("failed to set %s: %w", config.BaseURL.Name, err)
		}

	case "ollama":
		if err := c.SetVar(config.BaseURL.Name, config.BaseURL.Value); err != nil {
			return fmt.Errorf("failed to set %s: %w", config.BaseURL.Name, err)
		}
		if err := c.SetVar(config.APIKey.Name, config.APIKey.Value); err != nil {
			return fmt.Errorf("failed to set %s: %w", config.APIKey.Name, err)
		}
		if err := c.SetVar(config.Model.Name, config.Model.Value); err != nil {
			return fmt.Errorf("failed to set %s: %w", config.Model.Name, err)
		}
		if err := c.SetVar(config.PullTimeout.Name, config.PullTimeout.Value); err != nil {
			return fmt.Errorf("failed to set %s: %w", config.PullTimeout.Name, err)
		}
		if err := c.SetVar(config.PullEnabled.Name, config.PullEnabled.Value); err != nil {
			return fmt.Errorf("failed to set %s: %w", config.PullEnabled.Name, err)
		}
		if err := c.SetVar(config.LoadModelsEnabled.Name, config.LoadModelsEnabled.Value); err != nil {
			return fmt.Errorf("failed to set %s: %w", config.LoadModelsEnabled.Name, err)
		}

		containerPath, hostPath := "", ""
		if config.HostConfigPath.Value != "" {
			if slices.Contains(config.EmbeddedLLMConfigsPath, config.HostConfigPath.Value) {
				containerPath = config.HostConfigPath.Value
			} else {
				containerPath = DefaultOllamaConfigsPath
				hostPath = config.HostConfigPath.Value
			}
		}

		if err := c.SetVar(config.ConfigPath.Name, containerPath); err != nil {
			return fmt.Errorf("failed to set %s: %w", config.ConfigPath.Name, err)
		}
		if err := c.SetVar(config.HostConfigPath.Name, hostPath); err != nil {
			return fmt.Errorf("failed to set %s: %w", config.HostConfigPath.Name, err)
		}

	case "deepseek", "glm", "kimi", "qwen":
		if err := c.SetVar(config.APIKey.Name, config.APIKey.Value); err != nil {
			return fmt.Errorf("failed to set %s: %w", config.APIKey.Name, err)
		}
		if err := c.SetVar(config.BaseURL.Name, config.BaseURL.Value); err != nil {
			return fmt.Errorf("failed to set %s: %w", config.BaseURL.Name, err)
		}
		if err := c.SetVar(config.ProviderName.Name, config.ProviderName.Value); err != nil {
			return fmt.Errorf("failed to set %s: %w", config.ProviderName.Name, err)
		}

	case "custom":
		if err := c.SetVar(config.BaseURL.Name, config.BaseURL.Value); err != nil {
			return fmt.Errorf("failed to set %s: %w", config.BaseURL.Name, err)
		}
		if err := c.SetVar(config.APIKey.Name, config.APIKey.Value); err != nil {
			return fmt.Errorf("failed to set %s: %w", config.APIKey.Name, err)
		}
		if err := c.SetVar(config.Model.Name, config.Model.Value); err != nil {
			return fmt.Errorf("failed to set %s: %w", config.Model.Name, err)
		}
		if err := c.SetVar(config.LegacyReasoning.Name, config.LegacyReasoning.Value); err != nil {
			return fmt.Errorf("failed to set %s: %w", config.LegacyReasoning.Name, err)
		}
		if err := c.SetVar(config.PreserveReasoning.Name, config.PreserveReasoning.Value); err != nil {
			return fmt.Errorf("failed to set %s: %w", config.PreserveReasoning.Name, err)
		}
		if err := c.SetVar(config.ProviderName.Name, config.ProviderName.Value); err != nil {
			return fmt.Errorf("failed to set %s: %w", config.ProviderName.Name, err)
		}

		containerPath, hostPath := "", ""
		if config.HostConfigPath.Value != "" {
			if slices.Contains(config.EmbeddedLLMConfigsPath, config.HostConfigPath.Value) {
				containerPath = config.HostConfigPath.Value
			} else {
				containerPath = DefaultCustomConfigsPath
				hostPath = config.HostConfigPath.Value
			}
		}

		if err := c.SetVar(config.ConfigPath.Name, containerPath); err != nil {
			return fmt.Errorf("failed to set %s: %w", config.ConfigPath.Name, err)
		}
		if err := c.SetVar(config.HostConfigPath.Name, hostPath); err != nil {
			return fmt.Errorf("failed to set %s: %w", config.HostConfigPath.Name, err)
		}
	}

	return nil
}

// ResetLLMProviderConfig resets a specific LLM provider configuration.
func (c *controller) ResetLLMProviderConfig(providerID string) map[string]*LLMProviderConfig {
	var vars []string
	switch providerID {
	case "openai":
		vars = []string{"OPEN_AI_KEY", "OPEN_AI_SERVER_URL"}
	case "anthropic":
		vars = []string{"ANTHROPIC_API_KEY", "ANTHROPIC_SERVER_URL"}
	case "gemini":
		vars = []string{"GEMINI_API_KEY", "GEMINI_SERVER_URL"}
	case "bedrock":
		vars = []string{
			"BEDROCK_DEFAULT_AUTH", "BEDROCK_BEARER_TOKEN",
			"BEDROCK_ACCESS_KEY_ID", "BEDROCK_SECRET_ACCESS_KEY", "BEDROCK_SESSION_TOKEN",
			"BEDROCK_REGION", "BEDROCK_SERVER_URL",
		}
	case "ollama":
		vars = []string{
			"OLLAMA_SERVER_URL",
			"OLLAMA_SERVER_API_KEY",
			"OLLAMA_SERVER_MODEL",
			"OLLAMA_SERVER_CONFIG_PATH",
			"OLLAMA_SERVER_PULL_MODELS_TIMEOUT",
			"OLLAMA_SERVER_PULL_MODELS_ENABLED",
			"OLLAMA_SERVER_LOAD_MODELS_ENABLED",
			"PENTAGI_OLLAMA_SERVER_CONFIG_PATH",
		}
	case "deepseek":
		vars = []string{"DEEPSEEK_API_KEY", "DEEPSEEK_SERVER_URL", "DEEPSEEK_PROVIDER"}
	case "glm":
		vars = []string{"GLM_API_KEY", "GLM_SERVER_URL", "GLM_PROVIDER"}
	case "kimi":
		vars = []string{"KIMI_API_KEY", "KIMI_SERVER_URL", "KIMI_PROVIDER"}
	case "qwen":
		vars = []string{"QWEN_API_KEY", "QWEN_SERVER_URL", "QWEN_PROVIDER"}
	case "custom":
		vars = []string{
			"LLM_SERVER_URL", "LLM_SERVER_KEY", "LLM_SERVER_MODEL",
			"LLM_SERVER_CONFIG_PATH", "LLM_SERVER_LEGACY_REASONING",
			"LLM_SERVER_PRESERVE_REASONING", "LLM_SERVER_PROVIDER",
			"PENTAGI_LLM_SERVER_CONFIG_PATH",
		}
	}

	if len(vars) != 0 {
		if err := c.ResetVars(vars); err != nil {
			return nil
		}
	}

	return c.GetLLMProviders()
}
