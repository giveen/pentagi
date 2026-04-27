package controller

import (
	"fmt"

	"pentagi/cmd/installer/checker"
)

// GetLangfuseConfig returns the current Langfuse configuration.
func (c *controller) GetLangfuseConfig() *LangfuseConfig {
	vars, _ := c.GetVars([]string{
		"LANGFUSE_LISTEN_IP",
		"LANGFUSE_LISTEN_PORT",
		"LANGFUSE_BASE_URL",
		"LANGFUSE_PROJECT_ID",
		"LANGFUSE_PUBLIC_KEY",
		"LANGFUSE_SECRET_KEY",
		"LANGFUSE_INIT_USER_EMAIL",
		"LANGFUSE_INIT_USER_PASSWORD",
		"LANGFUSE_INIT_USER_NAME",
		"LANGFUSE_EE_LICENSE_KEY",
	})

	if v := vars["LANGFUSE_LISTEN_IP"]; v.Default == "" {
		v.Default = "127.0.0.1"
		vars["LANGFUSE_LISTEN_IP"] = v
	}
	if v := vars["LANGFUSE_LISTEN_PORT"]; v.Default == "" {
		v.Default = "4000"
		vars["LANGFUSE_LISTEN_PORT"] = v
	}

	var deploymentType string
	baseURL := vars["LANGFUSE_BASE_URL"]
	projectID := vars["LANGFUSE_PROJECT_ID"]
	publicKey := vars["LANGFUSE_PUBLIC_KEY"]
	secretKey := vars["LANGFUSE_SECRET_KEY"]
	adminEmail := vars["LANGFUSE_INIT_USER_EMAIL"]
	adminPassword := vars["LANGFUSE_INIT_USER_PASSWORD"]
	adminName := vars["LANGFUSE_INIT_USER_NAME"]
	licenseKey := vars["LANGFUSE_EE_LICENSE_KEY"]

	switch baseURL.Value {
	case "":
		deploymentType = "disabled"
	case checker.DefaultLangfuseEndpoint:
		deploymentType = "embedded"
		if projectID.Value == "" && !projectID.IsChanged {
			if initProjectID, ok := c.GetVar("LANGFUSE_INIT_PROJECT_ID"); ok {
				projectID.Value = initProjectID.Value
				projectID.IsChanged = true
			}
		}
		if publicKey.Value == "" && !publicKey.IsChanged {
			if initPublicKey, ok := c.GetVar("LANGFUSE_INIT_PROJECT_PUBLIC_KEY"); ok {
				publicKey.Value = initPublicKey.Value
				publicKey.IsChanged = true
			}
		}
		if secretKey.Value == "" && !secretKey.IsChanged {
			if initSecretKey, ok := c.GetVar("LANGFUSE_INIT_PROJECT_SECRET_KEY"); ok {
				secretKey.Value = initSecretKey.Value
				secretKey.IsChanged = true
			}
		}
	default:
		deploymentType = "external"
	}

	return &LangfuseConfig{
		DeploymentType: deploymentType,
		ListenIP:       vars["LANGFUSE_LISTEN_IP"],
		ListenPort:     vars["LANGFUSE_LISTEN_PORT"],
		BaseURL:        baseURL,
		ProjectID:      projectID,
		PublicKey:      publicKey,
		SecretKey:      secretKey,
		AdminEmail:     adminEmail,
		AdminPassword:  adminPassword,
		AdminName:      adminName,
		Installed:      c.checker.LangfuseInstalled,
		LicenseKey:     licenseKey,
	}
}

// UpdateLangfuseConfig updates Langfuse configuration with proper endpoint handling.
func (c *controller) UpdateLangfuseConfig(config *LangfuseConfig) error {
	if config == nil {
		return fmt.Errorf("config cannot be nil")
	}

	switch config.DeploymentType {
	case "embedded":
		config.BaseURL.Value = checker.DefaultLangfuseEndpoint

		if err := c.SetVar("LANGFUSE_LISTEN_IP", config.ListenIP.Value); err != nil {
			return fmt.Errorf("failed to set LANGFUSE_LISTEN_IP: %w", err)
		}
		if err := c.SetVar("LANGFUSE_LISTEN_PORT", config.ListenPort.Value); err != nil {
			return fmt.Errorf("failed to set LANGFUSE_LISTEN_PORT: %w", err)
		}
		if err := c.SetVar("LANGFUSE_EE_LICENSE_KEY", config.LicenseKey.Value); err != nil {
			return fmt.Errorf("failed to set LANGFUSE_EE_LICENSE_KEY: %w", err)
		}

		if !config.Installed {
			if err := c.SetVar("LANGFUSE_INIT_PROJECT_ID", config.ProjectID.Value); err != nil {
				return fmt.Errorf("failed to set LANGFUSE_INIT_PROJECT_ID: %w", err)
			}
			if err := c.SetVar("LANGFUSE_INIT_PROJECT_PUBLIC_KEY", config.PublicKey.Value); err != nil {
				return fmt.Errorf("failed to set LANGFUSE_INIT_PROJECT_PUBLIC_KEY: %w", err)
			}
			if err := c.SetVar("LANGFUSE_INIT_PROJECT_SECRET_KEY", config.SecretKey.Value); err != nil {
				return fmt.Errorf("failed to set LANGFUSE_INIT_PROJECT_SECRET_KEY: %w", err)
			}
			if err := c.SetVar("LANGFUSE_INIT_USER_EMAIL", config.AdminEmail.Value); err != nil {
				return fmt.Errorf("failed to set LANGFUSE_INIT_USER_EMAIL: %w", err)
			}
			if err := c.SetVar("LANGFUSE_INIT_USER_NAME", config.AdminName.Value); err != nil {
				return fmt.Errorf("failed to set LANGFUSE_INIT_USER_NAME: %w", err)
			}
			if err := c.SetVar("LANGFUSE_INIT_USER_PASSWORD", config.AdminPassword.Value); err != nil {
				return fmt.Errorf("failed to set LANGFUSE_INIT_USER_PASSWORD: %w", err)
			}
		}

	case "external":
		// use provided endpoint
	case "disabled":
		config.BaseURL.Value = ""
	}

	if err := c.SetVar("LANGFUSE_BASE_URL", config.BaseURL.Value); err != nil {
		return fmt.Errorf("failed to set LANGFUSE_BASE_URL: %w", err)
	}
	if err := c.SetVar("LANGFUSE_PROJECT_ID", config.ProjectID.Value); err != nil {
		return fmt.Errorf("failed to set LANGFUSE_PROJECT_ID: %w", err)
	}
	if err := c.SetVar("LANGFUSE_PUBLIC_KEY", config.PublicKey.Value); err != nil {
		return fmt.Errorf("failed to set LANGFUSE_PUBLIC_KEY: %w", err)
	}
	if err := c.SetVar("LANGFUSE_SECRET_KEY", config.SecretKey.Value); err != nil {
		return fmt.Errorf("failed to set LANGFUSE_SECRET_KEY: %w", err)
	}

	return nil
}

func (c *controller) ResetLangfuseConfig() *LangfuseConfig {
	vars := []string{
		"LANGFUSE_BASE_URL",
		"LANGFUSE_PROJECT_ID",
		"LANGFUSE_PUBLIC_KEY",
		"LANGFUSE_SECRET_KEY",
		"LANGFUSE_LISTEN_IP",
		"LANGFUSE_LISTEN_PORT",
		"LANGFUSE_EE_LICENSE_KEY",
	}

	if !c.checker.LangfuseInstalled {
		vars = append(vars,
			"LANGFUSE_INIT_USER_EMAIL",
			"LANGFUSE_INIT_USER_NAME",
			"LANGFUSE_INIT_USER_PASSWORD",
			"LANGFUSE_INIT_PROJECT_ID",
			"LANGFUSE_INIT_PROJECT_PUBLIC_KEY",
			"LANGFUSE_INIT_PROJECT_SECRET_KEY",
		)
	}

	if err := c.ResetVars(vars); err != nil {
		return nil
	}

	return c.GetLangfuseConfig()
}

// GetGraphitiConfig returns the current Graphiti configuration.
func (c *controller) GetGraphitiConfig() *GraphitiConfig {
	vars, _ := c.GetVars([]string{
		"GRAPHITI_URL",
		"GRAPHITI_TIMEOUT",
		"GRAPHITI_MODEL_NAME",
		"NEO4J_USER",
		"NEO4J_PASSWORD",
		"NEO4J_DATABASE",
		"NEO4J_URI",
	})

	if v := vars["GRAPHITI_TIMEOUT"]; v.Default == "" {
		v.Default = "30"
		vars["GRAPHITI_TIMEOUT"] = v
	}
	if v := vars["GRAPHITI_MODEL_NAME"]; v.Default == "" {
		v.Default = "gpt-5-mini"
		vars["GRAPHITI_MODEL_NAME"] = v
	}
	if v := vars["NEO4J_USER"]; v.Default == "" {
		v.Default = "neo4j"
		vars["NEO4J_USER"] = v
	}
	if v := vars["NEO4J_PASSWORD"]; v.Default == "" {
		v.Default = "devpassword"
		vars["NEO4J_PASSWORD"] = v
	}
	if v := vars["NEO4J_DATABASE"]; v.Default == "" {
		v.Default = "neo4j"
		vars["NEO4J_DATABASE"] = v
	}
	if v := vars["NEO4J_URI"]; v.Default == "" {
		v.Default = "bolt://neo4j:7687"
		vars["NEO4J_URI"] = v
	}

	graphitiURL := vars["GRAPHITI_URL"]
	graphitiEnabled, _ := c.GetVar("GRAPHITI_ENABLED")

	deploymentType := "external"
	if graphitiEnabled.Value != "true" || graphitiURL.Value == "" {
		deploymentType = "disabled"
	} else if graphitiURL.Value == checker.DefaultGraphitiEndpoint {
		deploymentType = "embedded"
	}

	return &GraphitiConfig{
		DeploymentType: deploymentType,
		GraphitiURL:    graphitiURL,
		Timeout:        vars["GRAPHITI_TIMEOUT"],
		ModelName:      vars["GRAPHITI_MODEL_NAME"],
		Neo4jUser:      vars["NEO4J_USER"],
		Neo4jPassword:  vars["NEO4J_PASSWORD"],
		Neo4jDatabase:  vars["NEO4J_DATABASE"],
		Neo4jURI:       vars["NEO4J_URI"],
		Installed:      c.checker.GraphitiInstalled,
	}
}

// UpdateGraphitiConfig updates Graphiti configuration.
func (c *controller) UpdateGraphitiConfig(config *GraphitiConfig) error {
	if config == nil {
		return fmt.Errorf("config cannot be nil")
	}

	switch config.DeploymentType {
	case "embedded":
		config.GraphitiURL.Value = checker.DefaultGraphitiEndpoint
		if err := c.SetVar("GRAPHITI_ENABLED", "true"); err != nil {
			return fmt.Errorf("failed to set GRAPHITI_ENABLED: %w", err)
		}
		if err := c.SetVar("GRAPHITI_TIMEOUT", config.Timeout.Value); err != nil {
			return fmt.Errorf("failed to set GRAPHITI_TIMEOUT: %w", err)
		}
		if err := c.SetVar("GRAPHITI_MODEL_NAME", config.ModelName.Value); err != nil {
			return fmt.Errorf("failed to set GRAPHITI_MODEL_NAME: %w", err)
		}
		if err := c.SetVar("NEO4J_USER", config.Neo4jUser.Value); err != nil {
			return fmt.Errorf("failed to set NEO4J_USER: %w", err)
		}
		if err := c.SetVar("NEO4J_PASSWORD", config.Neo4jPassword.Value); err != nil {
			return fmt.Errorf("failed to set NEO4J_PASSWORD: %w", err)
		}
		if err := c.SetVar("NEO4J_DATABASE", config.Neo4jDatabase.Value); err != nil {
			return fmt.Errorf("failed to set NEO4J_DATABASE: %w", err)
		}
	case "external":
		if err := c.SetVar("GRAPHITI_ENABLED", "true"); err != nil {
			return fmt.Errorf("failed to set GRAPHITI_ENABLED: %w", err)
		}
		if err := c.SetVar("GRAPHITI_TIMEOUT", config.Timeout.Value); err != nil {
			return fmt.Errorf("failed to set GRAPHITI_TIMEOUT: %w", err)
		}
	case "disabled":
		if err := c.SetVar("GRAPHITI_ENABLED", "false"); err != nil {
			return fmt.Errorf("failed to set GRAPHITI_ENABLED: %w", err)
		}
		config.GraphitiURL.Value = ""
	}

	if err := c.SetVar("GRAPHITI_URL", config.GraphitiURL.Value); err != nil {
		return fmt.Errorf("failed to set GRAPHITI_URL: %w", err)
	}

	return nil
}

func (c *controller) ResetGraphitiConfig() *GraphitiConfig {
	vars := []string{
		"GRAPHITI_ENABLED",
		"GRAPHITI_URL",
		"GRAPHITI_TIMEOUT",
		"GRAPHITI_MODEL_NAME",
		"NEO4J_USER",
		"NEO4J_PASSWORD",
		"NEO4J_DATABASE",
		"NEO4J_URI",
	}

	if err := c.ResetVars(vars); err != nil {
		return nil
	}

	return c.GetGraphitiConfig()
}

// GetObservabilityConfig returns the current observability configuration.
func (c *controller) GetObservabilityConfig() *ObservabilityConfig {
	vars, _ := c.GetVars([]string{
		"OTEL_HOST",
		"GRAFANA_LISTEN_IP",
		"GRAFANA_LISTEN_PORT",
		"OTEL_GRPC_LISTEN_IP",
		"OTEL_GRPC_LISTEN_PORT",
		"OTEL_HTTP_LISTEN_IP",
		"OTEL_HTTP_LISTEN_PORT",
	})

	defaults := map[string]string{
		"GRAFANA_LISTEN_IP":     "127.0.0.1",
		"GRAFANA_LISTEN_PORT":   "3000",
		"OTEL_GRPC_LISTEN_IP":   "127.0.0.1",
		"OTEL_GRPC_LISTEN_PORT": "8148",
		"OTEL_HTTP_LISTEN_IP":   "127.0.0.1",
		"OTEL_HTTP_LISTEN_PORT": "4318",
	}
	for k, def := range defaults {
		if v := vars[k]; v.Default == "" {
			v.Default = def
			vars[k] = v
		}
	}

	otelHost := vars["OTEL_HOST"]
	deploymentType := "external"
	switch otelHost.Value {
	case "":
		deploymentType = "disabled"
	case checker.DefaultObservabilityEndpoint:
		deploymentType = "embedded"
	}

	return &ObservabilityConfig{
		DeploymentType:     deploymentType,
		OTelHost:           otelHost,
		GrafanaListenIP:    vars["GRAFANA_LISTEN_IP"],
		GrafanaListenPort:  vars["GRAFANA_LISTEN_PORT"],
		OTelGrpcListenIP:   vars["OTEL_GRPC_LISTEN_IP"],
		OTelGrpcListenPort: vars["OTEL_GRPC_LISTEN_PORT"],
		OTelHttpListenIP:   vars["OTEL_HTTP_LISTEN_IP"],
		OTelHttpListenPort: vars["OTEL_HTTP_LISTEN_PORT"],
	}
}

// UpdateObservabilityConfig updates the observability configuration.
func (c *controller) UpdateObservabilityConfig(config *ObservabilityConfig) error {
	if config == nil {
		return fmt.Errorf("config cannot be nil")
	}

	langfuseOtelEnvVar, _ := c.GetVar("LANGFUSE_OTEL_EXPORTER_OTLP_ENDPOINT")

	switch config.DeploymentType {
	case "embedded":
		config.OTelHost.Value = checker.DefaultObservabilityEndpoint
		langfuseOtelEnvVar.Value = checker.DefaultLangfuseOtelEndpoint
		updates := map[string]string{
			"GRAFANA_LISTEN_IP":     config.GrafanaListenIP.Value,
			"GRAFANA_LISTEN_PORT":   config.GrafanaListenPort.Value,
			"OTEL_GRPC_LISTEN_IP":   config.OTelGrpcListenIP.Value,
			"OTEL_GRPC_LISTEN_PORT": config.OTelGrpcListenPort.Value,
			"OTEL_HTTP_LISTEN_IP":   config.OTelHttpListenIP.Value,
			"OTEL_HTTP_LISTEN_PORT": config.OTelHttpListenPort.Value,
		}
		if err := c.SetVars(updates); err != nil {
			return fmt.Errorf("failed to set embedded listen vars: %w", err)
		}
	case "external":
		langfuseOtelEnvVar.Value = ""
	case "disabled":
		config.OTelHost.Value = ""
		langfuseOtelEnvVar.Value = ""
	}

	if err := c.SetVar(langfuseOtelEnvVar.Name, langfuseOtelEnvVar.Value); err != nil {
		return fmt.Errorf("failed to set LANGFUSE_OTEL_EXPORTER_OTLP_ENDPOINT: %w", err)
	}
	if err := c.SetVar(config.OTelHost.Name, config.OTelHost.Value); err != nil {
		return fmt.Errorf("failed to set OTEL_HOST: %w", err)
	}

	return nil
}

func (c *controller) ResetObservabilityConfig() *ObservabilityConfig {
	vars := []string{
		"OTEL_HOST",
		"GRAFANA_LISTEN_IP",
		"GRAFANA_LISTEN_PORT",
		"OTEL_GRPC_LISTEN_IP",
		"OTEL_GRPC_LISTEN_PORT",
		"OTEL_HTTP_LISTEN_IP",
		"OTEL_HTTP_LISTEN_PORT",
		"LANGFUSE_OTEL_EXPORTER_OTLP_ENDPOINT",
	}

	if err := c.ResetVars(vars); err != nil {
		return nil
	}

	return c.GetObservabilityConfig()
}
