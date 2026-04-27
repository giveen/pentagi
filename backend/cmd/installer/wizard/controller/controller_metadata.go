package controller

import "pentagi/cmd/installer/wizard/locale"

func (c *controller) getVariableDescription(varName string) string {
	// map of environment variable name -> description
	envVarDescriptions := map[string]string{
		"OPEN_AI_KEY":                       locale.EnvDesc_OPEN_AI_KEY,
		"OPEN_AI_SERVER_URL":                locale.EnvDesc_OPEN_AI_SERVER_URL,
		"ANTHROPIC_API_KEY":                 locale.EnvDesc_ANTHROPIC_API_KEY,
		"ANTHROPIC_SERVER_URL":              locale.EnvDesc_ANTHROPIC_SERVER_URL,
		"GEMINI_API_KEY":                    locale.EnvDesc_GEMINI_API_KEY,
		"GEMINI_SERVER_URL":                 locale.EnvDesc_GEMINI_SERVER_URL,
		"BEDROCK_DEFAULT_AUTH":              locale.EnvDesc_BEDROCK_DEFAULT_AUTH,
		"BEDROCK_BEARER_TOKEN":              locale.EnvDesc_BEDROCK_BEARER_TOKEN,
		"BEDROCK_ACCESS_KEY_ID":             locale.EnvDesc_BEDROCK_ACCESS_KEY_ID,
		"BEDROCK_SECRET_ACCESS_KEY":         locale.EnvDesc_BEDROCK_SECRET_ACCESS_KEY,
		"BEDROCK_SESSION_TOKEN":             locale.EnvDesc_BEDROCK_SESSION_TOKEN,
		"BEDROCK_REGION":                    locale.EnvDesc_BEDROCK_REGION,
		"BEDROCK_SERVER_URL":                locale.EnvDesc_BEDROCK_SERVER_URL,
		"OLLAMA_SERVER_URL":                 locale.EnvDesc_OLLAMA_SERVER_URL,
		"OLLAMA_SERVER_API_KEY":             locale.EnvDesc_OLLAMA_SERVER_API_KEY,
		"OLLAMA_SERVER_MODEL":               locale.EnvDesc_OLLAMA_SERVER_MODEL,
		"OLLAMA_SERVER_CONFIG_PATH":         locale.EnvDesc_OLLAMA_SERVER_CONFIG_PATH,
		"OLLAMA_SERVER_PULL_MODELS_TIMEOUT": locale.EnvDesc_OLLAMA_SERVER_PULL_MODELS_TIMEOUT,
		"OLLAMA_SERVER_PULL_MODELS_ENABLED": locale.EnvDesc_OLLAMA_SERVER_PULL_MODELS_ENABLED,
		"OLLAMA_SERVER_LOAD_MODELS_ENABLED": locale.EnvDesc_OLLAMA_SERVER_LOAD_MODELS_ENABLED,
		"DEEPSEEK_API_KEY":                  locale.EnvDesc_DEEPSEEK_API_KEY,
		"DEEPSEEK_SERVER_URL":               locale.EnvDesc_DEEPSEEK_SERVER_URL,
		"DEEPSEEK_PROVIDER":                 locale.EnvDesc_DEEPSEEK_PROVIDER,
		"GLM_API_KEY":                       locale.EnvDesc_GLM_API_KEY,
		"GLM_SERVER_URL":                    locale.EnvDesc_GLM_SERVER_URL,
		"GLM_PROVIDER":                      locale.EnvDesc_GLM_PROVIDER,
		"KIMI_API_KEY":                      locale.EnvDesc_KIMI_API_KEY,
		"KIMI_SERVER_URL":                   locale.EnvDesc_KIMI_SERVER_URL,
		"KIMI_PROVIDER":                     locale.EnvDesc_KIMI_PROVIDER,
		"QWEN_API_KEY":                      locale.EnvDesc_QWEN_API_KEY,
		"QWEN_SERVER_URL":                   locale.EnvDesc_QWEN_SERVER_URL,
		"QWEN_PROVIDER":                     locale.EnvDesc_QWEN_PROVIDER,
		"LLM_SERVER_URL":                    locale.EnvDesc_LLM_SERVER_URL,
		"LLM_SERVER_KEY":                    locale.EnvDesc_LLM_SERVER_KEY,
		"LLM_SERVER_MODEL":                  locale.EnvDesc_LLM_SERVER_MODEL,
		"LLM_SERVER_CONFIG_PATH":            locale.EnvDesc_LLM_SERVER_CONFIG_PATH,
		"LLM_SERVER_LEGACY_REASONING":       locale.EnvDesc_LLM_SERVER_LEGACY_REASONING,
		"LLM_SERVER_PRESERVE_REASONING":     locale.EnvDesc_LLM_SERVER_PRESERVE_REASONING,
		"LLM_SERVER_PROVIDER":               locale.EnvDesc_LLM_SERVER_PROVIDER,

		"LANGFUSE_LISTEN_IP":   locale.EnvDesc_LANGFUSE_LISTEN_IP,
		"LANGFUSE_LISTEN_PORT": locale.EnvDesc_LANGFUSE_LISTEN_PORT,
		"LANGFUSE_BASE_URL":    locale.EnvDesc_LANGFUSE_BASE_URL,
		"LANGFUSE_PROJECT_ID":  locale.EnvDesc_LANGFUSE_PROJECT_ID,
		"LANGFUSE_PUBLIC_KEY":  locale.EnvDesc_LANGFUSE_PUBLIC_KEY,
		"LANGFUSE_SECRET_KEY":  locale.EnvDesc_LANGFUSE_SECRET_KEY,

		// langfuse init variables
		"LANGFUSE_INIT_PROJECT_ID":         locale.EnvDesc_LANGFUSE_INIT_PROJECT_ID,
		"LANGFUSE_INIT_PROJECT_PUBLIC_KEY": locale.EnvDesc_LANGFUSE_INIT_PROJECT_PUBLIC_KEY,
		"LANGFUSE_INIT_PROJECT_SECRET_KEY": locale.EnvDesc_LANGFUSE_INIT_PROJECT_SECRET_KEY,
		"LANGFUSE_INIT_USER_EMAIL":         locale.EnvDesc_LANGFUSE_INIT_USER_EMAIL,
		"LANGFUSE_INIT_USER_NAME":          locale.EnvDesc_LANGFUSE_INIT_USER_NAME,
		"LANGFUSE_INIT_USER_PASSWORD":      locale.EnvDesc_LANGFUSE_INIT_USER_PASSWORD,

		"LANGFUSE_OTEL_EXPORTER_OTLP_ENDPOINT": locale.EnvDesc_LANGFUSE_OTEL_EXPORTER_OTLP_ENDPOINT,

		"GRAFANA_LISTEN_IP":     locale.EnvDesc_GRAFANA_LISTEN_IP,
		"GRAFANA_LISTEN_PORT":   locale.EnvDesc_GRAFANA_LISTEN_PORT,
		"OTEL_GRPC_LISTEN_IP":   locale.EnvDesc_OTEL_GRPC_LISTEN_IP,
		"OTEL_GRPC_LISTEN_PORT": locale.EnvDesc_OTEL_GRPC_LISTEN_PORT,
		"OTEL_HTTP_LISTEN_IP":   locale.EnvDesc_OTEL_HTTP_LISTEN_IP,
		"OTEL_HTTP_LISTEN_PORT": locale.EnvDesc_OTEL_HTTP_LISTEN_PORT,
		"OTEL_HOST":             locale.EnvDesc_OTEL_HOST,

		"SUMMARIZER_PRESERVE_LAST":       locale.EnvDesc_SUMMARIZER_PRESERVE_LAST,
		"SUMMARIZER_USE_QA":              locale.EnvDesc_SUMMARIZER_USE_QA,
		"SUMMARIZER_SUM_MSG_HUMAN_IN_QA": locale.EnvDesc_SUMMARIZER_SUM_MSG_HUMAN_IN_QA,
		"SUMMARIZER_LAST_SEC_BYTES":      locale.EnvDesc_SUMMARIZER_LAST_SEC_BYTES,
		"SUMMARIZER_MAX_BP_BYTES":        locale.EnvDesc_SUMMARIZER_MAX_BP_BYTES,
		"SUMMARIZER_MAX_QA_BYTES":        locale.EnvDesc_SUMMARIZER_MAX_QA_BYTES,
		"SUMMARIZER_MAX_QA_SECTIONS":     locale.EnvDesc_SUMMARIZER_MAX_QA_SECTIONS,
		"SUMMARIZER_KEEP_QA_SECTIONS":    locale.EnvDesc_SUMMARIZER_KEEP_QA_SECTIONS,

		"ASSISTANT_SUMMARIZER_PRESERVE_LAST":    locale.EnvDesc_ASSISTANT_SUMMARIZER_PRESERVE_LAST,
		"ASSISTANT_SUMMARIZER_LAST_SEC_BYTES":   locale.EnvDesc_ASSISTANT_SUMMARIZER_LAST_SEC_BYTES,
		"ASSISTANT_SUMMARIZER_MAX_BP_BYTES":     locale.EnvDesc_ASSISTANT_SUMMARIZER_MAX_BP_BYTES,
		"ASSISTANT_SUMMARIZER_MAX_QA_BYTES":     locale.EnvDesc_ASSISTANT_SUMMARIZER_MAX_QA_BYTES,
		"ASSISTANT_SUMMARIZER_MAX_QA_SECTIONS":  locale.EnvDesc_ASSISTANT_SUMMARIZER_MAX_QA_SECTIONS,
		"ASSISTANT_SUMMARIZER_KEEP_QA_SECTIONS": locale.EnvDesc_ASSISTANT_SUMMARIZER_KEEP_QA_SECTIONS,

		"EMBEDDING_PROVIDER":        locale.EnvDesc_EMBEDDING_PROVIDER,
		"EMBEDDING_URL":             locale.EnvDesc_EMBEDDING_URL,
		"EMBEDDING_KEY":             locale.EnvDesc_EMBEDDING_KEY,
		"EMBEDDING_MODEL":           locale.EnvDesc_EMBEDDING_MODEL,
		"EMBEDDING_BATCH_SIZE":      locale.EnvDesc_EMBEDDING_BATCH_SIZE,
		"EMBEDDING_STRIP_NEW_LINES": locale.EnvDesc_EMBEDDING_STRIP_NEW_LINES,

		"ASK_USER": locale.EnvDesc_ASK_USER,

		"ASSISTANT_USE_AGENTS": locale.EnvDesc_ASSISTANT_USE_AGENTS,

		"EXECUTION_MONITOR_ENABLED":          locale.EnvDesc_EXECUTION_MONITOR_ENABLED,
		"EXECUTION_MONITOR_SAME_TOOL_LIMIT":  locale.EnvDesc_EXECUTION_MONITOR_SAME_TOOL_LIMIT,
		"EXECUTION_MONITOR_TOTAL_TOOL_LIMIT": locale.EnvDesc_EXECUTION_MONITOR_TOTAL_TOOL_LIMIT,
		"MAX_GENERAL_AGENT_TOOL_CALLS":       locale.EnvDesc_MAX_GENERAL_AGENT_TOOL_CALLS,
		"MAX_LIMITED_AGENT_TOOL_CALLS":       locale.EnvDesc_MAX_LIMITED_AGENT_TOOL_CALLS,
		"AGENT_PLANNING_STEP_ENABLED":        locale.EnvDesc_AGENT_PLANNING_STEP_ENABLED,

		"SCRAPER_PUBLIC_URL":                    locale.EnvDesc_SCRAPER_PUBLIC_URL,
		"SCRAPER_PRIVATE_URL":                   locale.EnvDesc_SCRAPER_PRIVATE_URL,
		"LOCAL_SCRAPER_USERNAME":                locale.EnvDesc_LOCAL_SCRAPER_USERNAME,
		"LOCAL_SCRAPER_PASSWORD":                locale.EnvDesc_LOCAL_SCRAPER_PASSWORD,
		"LOCAL_SCRAPER_MAX_CONCURRENT_SESSIONS": locale.EnvDesc_LOCAL_SCRAPER_MAX_CONCURRENT_SESSIONS,

		"DUCKDUCKGO_ENABLED":    locale.EnvDesc_DUCKDUCKGO_ENABLED,
		"DUCKDUCKGO_REGION":     locale.EnvDesc_DUCKDUCKGO_REGION,
		"DUCKDUCKGO_SAFESEARCH": locale.EnvDesc_DUCKDUCKGO_SAFESEARCH,
		"DUCKDUCKGO_TIME_RANGE": locale.EnvDesc_DUCKDUCKGO_TIME_RANGE,
		"SPLOITUS_ENABLED":      locale.EnvDesc_SPLOITUS_ENABLED,
		"PERPLEXITY_API_KEY":    locale.EnvDesc_PERPLEXITY_API_KEY,
		"TAVILY_API_KEY":        locale.EnvDesc_TAVILY_API_KEY,
		"TRAVERSAAL_API_KEY":    locale.EnvDesc_TRAVERSAAL_API_KEY,
		"GOOGLE_API_KEY":        locale.EnvDesc_GOOGLE_API_KEY,
		"GOOGLE_CX_KEY":         locale.EnvDesc_GOOGLE_CX_KEY,
		"GOOGLE_LR_KEY":         locale.EnvDesc_GOOGLE_LR_KEY,

		"PERPLEXITY_MODEL":        locale.EnvDesc_PERPLEXITY_MODEL,
		"PERPLEXITY_CONTEXT_SIZE": locale.EnvDesc_PERPLEXITY_CONTEXT_SIZE,

		"SEARXNG_URL":        locale.EnvDesc_SEARXNG_URL,
		"SEARXNG_CATEGORIES": locale.EnvDesc_SEARXNG_CATEGORIES,
		"SEARXNG_LANGUAGE":   locale.EnvDesc_SEARXNG_LANGUAGE,
		"SEARXNG_SAFESEARCH": locale.EnvDesc_SEARXNG_SAFESEARCH,
		"SEARXNG_TIME_RANGE": locale.EnvDesc_SEARXNG_TIME_RANGE,
		"SEARXNG_TIMEOUT":    locale.EnvDesc_SEARXNG_TIMEOUT,

		"DOCKER_INSIDE":                    locale.EnvDesc_DOCKER_INSIDE,
		"DOCKER_NET_ADMIN":                 locale.EnvDesc_DOCKER_NET_ADMIN,
		"DOCKER_SOCKET":                    locale.EnvDesc_DOCKER_SOCKET,
		"DOCKER_NETWORK":                   locale.EnvDesc_DOCKER_NETWORK,
		"DOCKER_PUBLIC_IP":                 locale.EnvDesc_DOCKER_PUBLIC_IP,
		"DOCKER_WORK_DIR":                  locale.EnvDesc_DOCKER_WORK_DIR,
		"DOCKER_DEFAULT_IMAGE":             locale.EnvDesc_DOCKER_DEFAULT_IMAGE,
		"DOCKER_DEFAULT_IMAGE_FOR_PENTEST": locale.EnvDesc_DOCKER_DEFAULT_IMAGE_FOR_PENTEST,
		"DOCKER_HOST":                      locale.EnvDesc_DOCKER_HOST,
		"DOCKER_TLS_VERIFY":                locale.EnvDesc_DOCKER_TLS_VERIFY,
		"DOCKER_CERT_PATH":                 locale.EnvDesc_DOCKER_CERT_PATH,

		"LICENSE_KEY":                       locale.EnvDesc_LICENSE_KEY,
		"PENTAGI_LISTEN_IP":                 locale.EnvDesc_PENTAGI_LISTEN_IP,
		"PENTAGI_LISTEN_PORT":               locale.EnvDesc_PENTAGI_LISTEN_PORT,
		"PUBLIC_URL":                        locale.EnvDesc_PUBLIC_URL,
		"CORS_ORIGINS":                      locale.EnvDesc_CORS_ORIGINS,
		"COOKIE_SIGNING_SALT":               locale.EnvDesc_COOKIE_SIGNING_SALT,
		"PROXY_URL":                         locale.EnvDesc_PROXY_URL,
		"EXTERNAL_SSL_CA_PATH":              locale.EnvDesc_EXTERNAL_SSL_CA_PATH,
		"EXTERNAL_SSL_INSECURE":             locale.EnvDesc_EXTERNAL_SSL_INSECURE,
		"PENTAGI_SSL_DIR":                   locale.EnvDesc_PENTAGI_SSL_DIR,
		"PENTAGI_DATA_DIR":                  locale.EnvDesc_PENTAGI_DATA_DIR,
		"PENTAGI_DOCKER_SOCKET":             locale.EnvDesc_PENTAGI_DOCKER_SOCKET,
		"PENTAGI_DOCKER_CERT_PATH":          locale.EnvDesc_PENTAGI_DOCKER_CERT_PATH,
		"PENTAGI_LLM_SERVER_CONFIG_PATH":    locale.EnvDesc_PENTAGI_LLM_SERVER_CONFIG_PATH,
		"PENTAGI_OLLAMA_SERVER_CONFIG_PATH": locale.EnvDesc_PENTAGI_OLLAMA_SERVER_CONFIG_PATH,

		"STATIC_DIR":     locale.EnvDesc_STATIC_DIR,
		"STATIC_URL":     locale.EnvDesc_STATIC_URL,
		"SERVER_PORT":    locale.EnvDesc_SERVER_PORT,
		"SERVER_HOST":    locale.EnvDesc_SERVER_HOST,
		"SERVER_SSL_CRT": locale.EnvDesc_SERVER_SSL_CRT,
		"SERVER_SSL_KEY": locale.EnvDesc_SERVER_SSL_KEY,
		"SERVER_USE_SSL": locale.EnvDesc_SERVER_USE_SSL,

		"OAUTH_GOOGLE_CLIENT_ID":     locale.EnvDesc_OAUTH_GOOGLE_CLIENT_ID,
		"OAUTH_GOOGLE_CLIENT_SECRET": locale.EnvDesc_OAUTH_GOOGLE_CLIENT_SECRET,
		"OAUTH_GITHUB_CLIENT_ID":     locale.EnvDesc_OAUTH_GITHUB_CLIENT_ID,
		"OAUTH_GITHUB_CLIENT_SECRET": locale.EnvDesc_OAUTH_GITHUB_CLIENT_SECRET,

		"LANGFUSE_EE_LICENSE_KEY": locale.EnvDesc_LANGFUSE_EE_LICENSE_KEY,

		"GRAPHITI_URL":        locale.EnvDesc_GRAPHITI_URL,
		"GRAPHITI_TIMEOUT":    locale.EnvDesc_GRAPHITI_TIMEOUT,
		"GRAPHITI_MODEL_NAME": locale.EnvDesc_GRAPHITI_MODEL_NAME,
		"NEO4J_USER":          locale.EnvDesc_NEO4J_USER,
		"NEO4J_DATABASE":      locale.EnvDesc_NEO4J_DATABASE,

		"PENTAGI_POSTGRES_PASSWORD": locale.EnvDesc_PENTAGI_POSTGRES_PASSWORD,
		"NEO4J_PASSWORD":            locale.EnvDesc_NEO4J_PASSWORD,
	}
	if desc, ok := envVarDescriptions[varName]; ok {
		return desc
	}
	return varName
}

// maskedVariables contains environment variable names that should be masked in display
var maskedVariables = map[string]bool{
	// API keys and Secrets
	"OPEN_AI_KEY":               true,
	"ANTHROPIC_API_KEY":         true,
	"GEMINI_API_KEY":            true,
	"BEDROCK_BEARER_TOKEN":      true,
	"BEDROCK_ACCESS_KEY_ID":     true,
	"BEDROCK_SECRET_ACCESS_KEY": true,
	"BEDROCK_SESSION_TOKEN":     true,
	"OLLAMA_SERVER_API_KEY":     true,
	"DEEPSEEK_API_KEY":          true,
	"GLM_API_KEY":               true,
	"KIMI_API_KEY":              true,
	"QWEN_API_KEY":              true,
	"LLM_SERVER_KEY":            true,
	"LANGFUSE_PUBLIC_KEY":       true,
	"LANGFUSE_SECRET_KEY":       true,
	"EMBEDDING_KEY":             true,
	"LOCAL_SCRAPER_PASSWORD":    true,
	"PERPLEXITY_API_KEY":        true,
	"TAVILY_API_KEY":            true,
	"TRAVERSAAL_API_KEY":        true,
	"GOOGLE_API_KEY":            true,
	"GOOGLE_CX_KEY":             true,

	// oauth client secrets
	"OAUTH_GOOGLE_CLIENT_SECRET": true,
	"OAUTH_GITHUB_CLIENT_SECRET": true,

	// urls can embed credentials; mask to avoid leaking secrets
	"PROXY_URL":           true,
	"SCRAPER_PUBLIC_URL":  true,
	"SCRAPER_PRIVATE_URL": true,

	// langfuse init secrets
	"LANGFUSE_INIT_PROJECT_PUBLIC_KEY": true,
	"LANGFUSE_INIT_PROJECT_SECRET_KEY": true,
	"LANGFUSE_INIT_USER_PASSWORD":      true,

	// langfuse license key
	"LANGFUSE_EE_LICENSE_KEY": true,

	// postgres password for pentagi service (pgvector binds on localhost)
	"PENTAGI_POSTGRES_PASSWORD": true,

	// neo4j password for graphiti service (neo4j binds on localhost)
	"NEO4J_PASSWORD": true,

	// langfuse stack secrets (compose-managed)
	"LANGFUSE_SALT":                      true,
	"LANGFUSE_ENCRYPTION_KEY":            true,
	"LANGFUSE_NEXTAUTH_SECRET":           true,
	"LANGFUSE_CLICKHOUSE_PASSWORD":       true,
	"LANGFUSE_S3_ACCESS_KEY_ID":          true,
	"LANGFUSE_S3_SECRET_ACCESS_KEY":      true,
	"LANGFUSE_REDIS_AUTH":                true,
	"LANGFUSE_AUTH_CUSTOM_CLIENT_SECRET": true,

	// server settings
	"COOKIE_SIGNING_SALT": true,
}

// isVariableMasked returns true if the variable should be masked in display
func (c *controller) isVariableMasked(varName string) bool {
	return maskedVariables[varName]
}

// criticalVariables contains environment variable names that require service restart
var criticalVariables = map[string]bool{
	// LLM Provider changes
	"OPEN_AI_KEY":                       true,
	"OPEN_AI_SERVER_URL":                true,
	"ANTHROPIC_API_KEY":                 true,
	"ANTHROPIC_SERVER_URL":              true,
	"GEMINI_API_KEY":                    true,
	"GEMINI_SERVER_URL":                 true,
	"BEDROCK_DEFAULT_AUTH":              true,
	"BEDROCK_BEARER_TOKEN":              true,
	"BEDROCK_ACCESS_KEY_ID":             true,
	"BEDROCK_SECRET_ACCESS_KEY":         true,
	"BEDROCK_SESSION_TOKEN":             true,
	"BEDROCK_REGION":                    true,
	"OLLAMA_SERVER_URL":                 true,
	"OLLAMA_SERVER_API_KEY":             true,
	"OLLAMA_SERVER_MODEL":               true,
	"OLLAMA_SERVER_CONFIG_PATH":         true,
	"OLLAMA_SERVER_PULL_MODELS_TIMEOUT": true,
	"OLLAMA_SERVER_PULL_MODELS_ENABLED": true,
	"OLLAMA_SERVER_LOAD_MODELS_ENABLED": true,
	"DEEPSEEK_API_KEY":                  true,
	"DEEPSEEK_SERVER_URL":               true,
	"DEEPSEEK_PROVIDER":                 true,
	"GLM_API_KEY":                       true,
	"GLM_SERVER_URL":                    true,
	"GLM_PROVIDER":                      true,
	"KIMI_API_KEY":                      true,
	"KIMI_SERVER_URL":                   true,
	"KIMI_PROVIDER":                     true,
	"QWEN_API_KEY":                      true,
	"QWEN_SERVER_URL":                   true,
	"QWEN_PROVIDER":                     true,
	"LLM_SERVER_URL":                    true,
	"LLM_SERVER_KEY":                    true,
	"LLM_SERVER_MODEL":                  true,
	"LLM_SERVER_CONFIG_PATH":            true,
	"LLM_SERVER_LEGACY_REASONING":       true,
	"LLM_SERVER_PRESERVE_REASONING":     true,
	"LLM_SERVER_PROVIDER":               true,

	// tools changes
	"DUCKDUCKGO_ENABLED":      true,
	"DUCKDUCKGO_REGION":       true,
	"DUCKDUCKGO_SAFESEARCH":   true,
	"DUCKDUCKGO_TIME_RANGE":   true,
	"SPLOITUS_ENABLED":        true,
	"PERPLEXITY_API_KEY":      true,
	"PERPLEXITY_MODEL":        true,
	"PERPLEXITY_CONTEXT_SIZE": true,
	"TAVILY_API_KEY":          true,
	"TRAVERSAAL_API_KEY":      true,
	"GOOGLE_API_KEY":          true,
	"GOOGLE_CX_KEY":           true,
	"GOOGLE_LR_KEY":           true,
	"SEARXNG_URL":             true,
	"SEARXNG_CATEGORIES":      true,
	"SEARXNG_LANGUAGE":        true,
	"SEARXNG_SAFESEARCH":      true,
	"SEARXNG_TIME_RANGE":      true,
	"SEARXNG_TIMEOUT":         true,

	// mounting custom LLM server config into pentagi container changes volume mapping
	"PENTAGI_LLM_SERVER_CONFIG_PATH":    true,
	"PENTAGI_OLLAMA_SERVER_CONFIG_PATH": true,

	// Embedding provider changes
	"EMBEDDING_PROVIDER":        true,
	"EMBEDDING_URL":             true,
	"EMBEDDING_KEY":             true,
	"EMBEDDING_MODEL":           true,
	"EMBEDDING_BATCH_SIZE":      true,
	"EMBEDDING_STRIP_NEW_LINES": true,

	// Docker configuration changes
	"DOCKER_INSIDE":                    true,
	"DOCKER_NET_ADMIN":                 true,
	"DOCKER_SOCKET":                    true,
	"DOCKER_NETWORK":                   true,
	"DOCKER_PUBLIC_IP":                 true,
	"DOCKER_DEFAULT_IMAGE":             true,
	"DOCKER_DEFAULT_IMAGE_FOR_PENTEST": true,
	"DOCKER_HOST":                      true,
	"DOCKER_TLS_VERIFY":                true,
	"DOCKER_CERT_PATH":                 true,
	"PENTAGI_DOCKER_SOCKET":            true,

	// observability changes
	"OTEL_HOST": true,

	// graphiti changes
	"GRAPHITI_URL":        true,
	"GRAPHITI_TIMEOUT":    true,
	"GRAPHITI_MODEL_NAME": true,

	// server settings changes
	"ASK_USER":                           true,
	"EXECUTION_MONITOR_ENABLED":          true,
	"EXECUTION_MONITOR_SAME_TOOL_LIMIT":  true,
	"EXECUTION_MONITOR_TOTAL_TOOL_LIMIT": true,
	"MAX_GENERAL_AGENT_TOOL_CALLS":       true,
	"MAX_LIMITED_AGENT_TOOL_CALLS":       true,
	"AGENT_PLANNING_STEP_ENABLED":        true,

	"LICENSE_KEY":           true,
	"PENTAGI_LISTEN_IP":     true,
	"PENTAGI_LISTEN_PORT":   true,
	"PUBLIC_URL":            true,
	"CORS_ORIGINS":          true,
	"COOKIE_SIGNING_SALT":   true,
	"PROXY_URL":             true,
	"EXTERNAL_SSL_CA_PATH":  true,
	"EXTERNAL_SSL_INSECURE": true,
	"STATIC_DIR":            true,
	"STATIC_URL":            true,
	"SERVER_PORT":           true,
	"SERVER_HOST":           true,
	"SERVER_SSL_CRT":        true,
	"SERVER_SSL_KEY":        true,
	"SERVER_USE_SSL":        true,
	"PENTAGI_SSL_DIR":       true,
	"PENTAGI_DATA_DIR":      true,

	// scraper settings
	"SCRAPER_PUBLIC_URL":  true,
	"SCRAPER_PRIVATE_URL": true,

	// oauth settings
	"OAUTH_GOOGLE_CLIENT_ID":     true,
	"OAUTH_GOOGLE_CLIENT_SECRET": true,
	"OAUTH_GITHUB_CLIENT_ID":     true,
	"OAUTH_GITHUB_CLIENT_SECRET": true,

	// langfuse integration settings passed to pentagi
	"LANGFUSE_BASE_URL":   true,
	"LANGFUSE_PROJECT_ID": true,
	"LANGFUSE_PUBLIC_KEY": true,
	"LANGFUSE_SECRET_KEY": true,

	// summarizer settings (general)
	"SUMMARIZER_PRESERVE_LAST":       true,
	"SUMMARIZER_USE_QA":              true,
	"SUMMARIZER_SUM_MSG_HUMAN_IN_QA": true,
	"SUMMARIZER_LAST_SEC_BYTES":      true,
	"SUMMARIZER_MAX_BP_BYTES":        true,
	"SUMMARIZER_MAX_QA_SECTIONS":     true,
	"SUMMARIZER_MAX_QA_BYTES":        true,
	"SUMMARIZER_KEEP_QA_SECTIONS":    true,

	// assistant-level settings
	"ASSISTANT_USE_AGENTS":                  true,
	"ASSISTANT_SUMMARIZER_PRESERVE_LAST":    true,
	"ASSISTANT_SUMMARIZER_LAST_SEC_BYTES":   true,
	"ASSISTANT_SUMMARIZER_MAX_BP_BYTES":     true,
	"ASSISTANT_SUMMARIZER_MAX_QA_SECTIONS":  true,
	"ASSISTANT_SUMMARIZER_MAX_QA_BYTES":     true,
	"ASSISTANT_SUMMARIZER_KEEP_QA_SECTIONS": true,
}

// isCriticalVariable returns true if changing this variable requires service restart
func (c *controller) isCriticalVariable(varName string) bool {
	return criticalVariables[varName]
}
