package controller

import (
	"fmt"
	"net/url"
	"strings"
)

// GetScraperConfig returns current scraper configuration.
func (c *controller) GetScraperConfig() *ScraperConfig {
	publicURL, _ := c.GetVar("SCRAPER_PUBLIC_URL")
	privateURL, _ := c.GetVar("SCRAPER_PRIVATE_URL")
	localUsername, _ := c.GetVar("LOCAL_SCRAPER_USERNAME")
	localPassword, _ := c.GetVar("LOCAL_SCRAPER_PASSWORD")
	maxSessions, _ := c.GetVar("LOCAL_SCRAPER_MAX_CONCURRENT_SESSIONS")

	config := &ScraperConfig{
		PublicURL:             publicURL,
		PrivateURL:            privateURL,
		LocalUsername:         localUsername,
		LocalPassword:         localPassword,
		MaxConcurrentSessions: maxSessions,
	}

	config.Mode = c.determineScraperMode(privateURL.Value, publicURL.Value)

	if config.Mode == "external" || config.Mode == "embedded" {
		config.PublicUsername, config.PublicPassword = c.extractCredentialsFromURL(publicURL.Value)
		config.PrivateUsername, config.PrivatePassword = c.extractCredentialsFromURL(privateURL.Value)
		config.PublicURL.Value = RemoveCredentialsFromURL(publicURL.Value)
		config.PrivateURL.Value = RemoveCredentialsFromURL(privateURL.Value)
	}

	return config
}

func (c *controller) determineScraperMode(privateURL, publicURL string) string {
	if privateURL == "" && publicURL == "" {
		return "disabled"
	}

	parsedURL, err := url.Parse(privateURL)
	if err != nil {
		return "external"
	}

	if parsedURL.Scheme == DefaultScraperSchema && parsedURL.Hostname() == DefaultScraperDomain {
		return "embedded"
	}

	return "external"
}

func (c *controller) extractCredentialsFromURL(urlStr string) (username, password string) {
	if urlStr == "" {
		return "", ""
	}

	parsedURL, err := url.Parse(urlStr)
	if err != nil || parsedURL.User == nil {
		return "", ""
	}

	username = parsedURL.User.Username()
	password, _ = parsedURL.User.Password()
	return username, password
}

// UpdateScraperConfig updates scraper configuration.
func (c *controller) UpdateScraperConfig(config *ScraperConfig) error {
	if config == nil {
		return fmt.Errorf("config cannot be nil")
	}

	switch config.Mode {
	case "disabled":
		if err := c.SetVar("SCRAPER_PUBLIC_URL", ""); err != nil {
			return fmt.Errorf("failed to clear SCRAPER_PUBLIC_URL: %w", err)
		}
		if err := c.SetVar("SCRAPER_PRIVATE_URL", ""); err != nil {
			return fmt.Errorf("failed to clear SCRAPER_PRIVATE_URL: %w", err)
		}
	case "external":
		privateURL := config.PrivateURL.Value
		if config.PrivateUsername != "" && config.PrivatePassword != "" {
			privateURL = c.addCredentialsToURL(config.PrivateURL.Value, config.PrivateUsername, config.PrivatePassword)
		}
		publicURL := config.PublicURL.Value
		if config.PublicUsername != "" && config.PublicPassword != "" {
			publicURL = c.addCredentialsToURL(config.PublicURL.Value, config.PublicUsername, config.PublicPassword)
		}
		if err := c.SetVar("SCRAPER_PUBLIC_URL", publicURL); err != nil {
			return fmt.Errorf("failed to set SCRAPER_PUBLIC_URL: %w", err)
		}
		if err := c.SetVar("SCRAPER_PRIVATE_URL", privateURL); err != nil {
			return fmt.Errorf("failed to set SCRAPER_PRIVATE_URL: %w", err)
		}
	case "embedded":
		privateURL := DefaultScraperBaseURL
		if config.PrivateUsername != "" && config.PrivatePassword != "" {
			privateURL = c.addCredentialsToURL(privateURL, config.PrivateUsername, config.PrivatePassword)
		}

		publicURL := config.PublicURL.Value
		if config.PublicUsername != "" && config.PublicPassword != "" {
			if publicURL == "" {
				publicURL = privateURL
			}
			publicURL = c.addCredentialsToURL(publicURL, config.PublicUsername, config.PublicPassword)
		}

		if err := c.SetVar("SCRAPER_PUBLIC_URL", publicURL); err != nil {
			return fmt.Errorf("failed to set SCRAPER_PUBLIC_URL: %w", err)
		}
		if err := c.SetVar("SCRAPER_PRIVATE_URL", privateURL); err != nil {
			return fmt.Errorf("failed to set SCRAPER_PRIVATE_URL: %w", err)
		}
		if err := c.SetVar("LOCAL_SCRAPER_USERNAME", config.PrivateUsername); err != nil {
			return fmt.Errorf("failed to set LOCAL_SCRAPER_USERNAME: %w", err)
		}
		if err := c.SetVar("LOCAL_SCRAPER_PASSWORD", config.PrivatePassword); err != nil {
			return fmt.Errorf("failed to set LOCAL_SCRAPER_PASSWORD: %w", err)
		}
		if err := c.SetVar("LOCAL_SCRAPER_MAX_CONCURRENT_SESSIONS", config.MaxConcurrentSessions.Value); err != nil {
			return fmt.Errorf("failed to set LOCAL_SCRAPER_MAX_CONCURRENT_SESSIONS: %w", err)
		}
	}

	return nil
}

func (c *controller) addCredentialsToURL(urlStr, username, password string) string {
	if username == "" || password == "" {
		return urlStr
	}
	if urlStr == "" {
		return ""
	}
	parsedURL, err := url.Parse(urlStr)
	if err != nil {
		return urlStr
	}
	parsedURL.User = url.UserPassword(username, password)
	return parsedURL.String()
}

func (c *controller) ResetScraperConfig() *ScraperConfig {
	vars := []string{
		"SCRAPER_PUBLIC_URL",
		"SCRAPER_PRIVATE_URL",
		"LOCAL_SCRAPER_USERNAME",
		"LOCAL_SCRAPER_PASSWORD",
		"LOCAL_SCRAPER_MAX_CONCURRENT_SESSIONS",
	}
	if err := c.ResetVars(vars); err != nil {
		return nil
	}
	return c.GetScraperConfig()
}

// GetSearchEnginesConfig returns current search engines configuration.
func (c *controller) GetSearchEnginesConfig() *SearchEnginesConfig {
	duckduckgoEnabled, _ := c.GetVar("DUCKDUCKGO_ENABLED")
	duckduckgoRegion, _ := c.GetVar("DUCKDUCKGO_REGION")
	duckduckgoSafeSearch, _ := c.GetVar("DUCKDUCKGO_SAFESEARCH")
	duckduckgoTimeRange, _ := c.GetVar("DUCKDUCKGO_TIME_RANGE")
	sploitusEnabled, _ := c.GetVar("SPLOITUS_ENABLED")
	perplexityAPIKey, _ := c.GetVar("PERPLEXITY_API_KEY")
	tavilyAPIKey, _ := c.GetVar("TAVILY_API_KEY")
	traversaalAPIKey, _ := c.GetVar("TRAVERSAAL_API_KEY")
	googleAPIKey, _ := c.GetVar("GOOGLE_API_KEY")
	googleCXKey, _ := c.GetVar("GOOGLE_CX_KEY")
	googleLRKey, _ := c.GetVar("GOOGLE_LR_KEY")
	perplexityModel, _ := c.GetVar("PERPLEXITY_MODEL")
	perplexityContextSize, _ := c.GetVar("PERPLEXITY_CONTEXT_SIZE")
	searxngURL, _ := c.GetVar("SEARXNG_URL")
	searxngCategories, _ := c.GetVar("SEARXNG_CATEGORIES")
	searxngLanguage, _ := c.GetVar("SEARXNG_LANGUAGE")
	searxngSafeSearch, _ := c.GetVar("SEARXNG_SAFESEARCH")
	searxngTimeRange, _ := c.GetVar("SEARXNG_TIME_RANGE")
	searxngTimeout, _ := c.GetVar("SEARXNG_TIMEOUT")

	config := &SearchEnginesConfig{
		DuckDuckGoEnabled:     duckduckgoEnabled,
		DuckDuckGoRegion:      duckduckgoRegion,
		DuckDuckGoSafeSearch:  duckduckgoSafeSearch,
		DuckDuckGoTimeRange:   duckduckgoTimeRange,
		SploitusEnabled:       sploitusEnabled,
		PerplexityAPIKey:      perplexityAPIKey,
		PerplexityModel:       perplexityModel,
		PerplexityContextSize: perplexityContextSize,
		TavilyAPIKey:          tavilyAPIKey,
		TraversaalAPIKey:      traversaalAPIKey,
		GoogleAPIKey:          googleAPIKey,
		GoogleCXKey:           googleCXKey,
		GoogleLRKey:           googleLRKey,
		SearxngURL:            searxngURL,
		SearxngCategories:     searxngCategories,
		SearxngLanguage:       searxngLanguage,
		SearxngSafeSearch:     searxngSafeSearch,
		SearxngTimeRange:      searxngTimeRange,
		SearxngTimeout:        searxngTimeout,
	}

	configuredCount := 0
	if duckduckgoEnabled.Value == "true" || (duckduckgoEnabled.Value == "" && duckduckgoEnabled.Default == "true") {
		configuredCount++
	}
	if sploitusEnabled.Value == "true" || (sploitusEnabled.Value == "" && sploitusEnabled.Default == "true") {
		configuredCount++
	}
	if perplexityAPIKey.Value != "" {
		configuredCount++
	}
	if tavilyAPIKey.Value != "" {
		configuredCount++
	}
	if traversaalAPIKey.Value != "" {
		configuredCount++
	}
	if googleAPIKey.Value != "" && googleCXKey.Value != "" {
		configuredCount++
	}
	if searxngURL.Value != "" {
		configuredCount++
	}
	config.ConfiguredCount = configuredCount

	return config
}

func (c *controller) UpdateSearchEnginesConfig(config *SearchEnginesConfig) error {
	if config == nil {
		return fmt.Errorf("config cannot be nil")
	}

	updates := map[string]string{
		"DUCKDUCKGO_ENABLED":      config.DuckDuckGoEnabled.Value,
		"DUCKDUCKGO_REGION":       config.DuckDuckGoRegion.Value,
		"DUCKDUCKGO_SAFESEARCH":   config.DuckDuckGoSafeSearch.Value,
		"DUCKDUCKGO_TIME_RANGE":   config.DuckDuckGoTimeRange.Value,
		"SPLOITUS_ENABLED":        config.SploitusEnabled.Value,
		"PERPLEXITY_API_KEY":      config.PerplexityAPIKey.Value,
		"PERPLEXITY_MODEL":        config.PerplexityModel.Value,
		"PERPLEXITY_CONTEXT_SIZE": config.PerplexityContextSize.Value,
		"TAVILY_API_KEY":          config.TavilyAPIKey.Value,
		"TRAVERSAAL_API_KEY":      config.TraversaalAPIKey.Value,
		"GOOGLE_API_KEY":          config.GoogleAPIKey.Value,
		"GOOGLE_CX_KEY":           config.GoogleCXKey.Value,
		"GOOGLE_LR_KEY":           config.GoogleLRKey.Value,
		"SEARXNG_URL":             config.SearxngURL.Value,
		"SEARXNG_CATEGORIES":      config.SearxngCategories.Value,
		"SEARXNG_LANGUAGE":        config.SearxngLanguage.Value,
		"SEARXNG_SAFESEARCH":      config.SearxngSafeSearch.Value,
		"SEARXNG_TIME_RANGE":      config.SearxngTimeRange.Value,
		"SEARXNG_TIMEOUT":         config.SearxngTimeout.Value,
	}
	if err := c.SetVars(updates); err != nil {
		return err
	}
	return nil
}

func (c *controller) ResetSearchEnginesConfig() *SearchEnginesConfig {
	vars := []string{
		"DUCKDUCKGO_ENABLED", "DUCKDUCKGO_REGION", "DUCKDUCKGO_SAFESEARCH", "DUCKDUCKGO_TIME_RANGE",
		"SPLOITUS_ENABLED", "PERPLEXITY_API_KEY", "PERPLEXITY_MODEL", "PERPLEXITY_CONTEXT_SIZE",
		"TAVILY_API_KEY", "TRAVERSAAL_API_KEY", "GOOGLE_API_KEY", "GOOGLE_CX_KEY", "GOOGLE_LR_KEY",
		"SEARXNG_URL", "SEARXNG_CATEGORIES", "SEARXNG_LANGUAGE", "SEARXNG_SAFESEARCH", "SEARXNG_TIME_RANGE", "SEARXNG_TIMEOUT",
	}
	if err := c.ResetVars(vars); err != nil {
		return nil
	}
	return c.GetSearchEnginesConfig()
}

func (c *controller) GetDockerConfig() *DockerConfig {
	vars, _ := c.GetVars([]string{
		"DOCKER_INSIDE", "DOCKER_NET_ADMIN", "DOCKER_SOCKET", "DOCKER_NETWORK", "DOCKER_PUBLIC_IP", "DOCKER_WORK_DIR",
		"DOCKER_DEFAULT_IMAGE", "DOCKER_DEFAULT_IMAGE_FOR_PENTEST", "DOCKER_HOST", "DOCKER_TLS_VERIFY", "PENTAGI_DOCKER_CERT_PATH",
	})

	config := &DockerConfig{
		DockerInside:                 vars["DOCKER_INSIDE"],
		DockerNetAdmin:               vars["DOCKER_NET_ADMIN"],
		DockerSocket:                 vars["DOCKER_SOCKET"],
		DockerNetwork:                vars["DOCKER_NETWORK"],
		DockerPublicIP:               vars["DOCKER_PUBLIC_IP"],
		DockerWorkDir:                vars["DOCKER_WORK_DIR"],
		DockerDefaultImage:           vars["DOCKER_DEFAULT_IMAGE"],
		DockerDefaultImageForPentest: vars["DOCKER_DEFAULT_IMAGE_FOR_PENTEST"],
		DockerHost:                   vars["DOCKER_HOST"],
		DockerTLSVerify:              vars["DOCKER_TLS_VERIFY"],
		HostDockerCertPath:           vars["PENTAGI_DOCKER_CERT_PATH"],
	}
	if config.DockerHost.Default == "" {
		config.DockerHost.Default = "unix:///var/run/docker.sock"
	}
	config.Configured = config.DockerInside.Value != "" || config.DockerDefaultImage.Value != "" || config.DockerDefaultImageForPentest.Value != ""
	return config
}

func (c *controller) UpdateDockerConfig(config *DockerConfig) error {
	updates := map[string]string{
		"DOCKER_INSIDE":                    config.DockerInside.Value,
		"DOCKER_NET_ADMIN":                 config.DockerNetAdmin.Value,
		"DOCKER_SOCKET":                    config.DockerSocket.Value,
		"DOCKER_NETWORK":                   config.DockerNetwork.Value,
		"DOCKER_PUBLIC_IP":                 config.DockerPublicIP.Value,
		"DOCKER_WORK_DIR":                  config.DockerWorkDir.Value,
		"DOCKER_DEFAULT_IMAGE":             config.DockerDefaultImage.Value,
		"DOCKER_DEFAULT_IMAGE_FOR_PENTEST": config.DockerDefaultImageForPentest.Value,
		"DOCKER_HOST":                      config.DockerHost.Value,
		"DOCKER_TLS_VERIFY":                config.DockerTLSVerify.Value,
		"PENTAGI_DOCKER_CERT_PATH":         config.HostDockerCertPath.Value,
	}

	dockerHost := config.DockerHost.Value
	if strings.HasPrefix(dockerHost, "unix://") && !config.DockerHost.IsDefault() {
		updates["PENTAGI_DOCKER_SOCKET"] = strings.TrimPrefix(dockerHost, "unix://")
	} else {
		updates["PENTAGI_DOCKER_SOCKET"] = ""
	}

	if config.HostDockerCertPath.Value != "" {
		updates["DOCKER_CERT_PATH"] = DefaultDockerCertPath
	} else {
		updates["DOCKER_CERT_PATH"] = ""
	}

	return c.SetVars(updates)
}

func (c *controller) ResetDockerConfig() *DockerConfig {
	vars := []string{
		"DOCKER_INSIDE", "DOCKER_NET_ADMIN", "DOCKER_SOCKET", "DOCKER_NETWORK", "DOCKER_PUBLIC_IP", "DOCKER_WORK_DIR",
		"DOCKER_DEFAULT_IMAGE", "DOCKER_DEFAULT_IMAGE_FOR_PENTEST", "DOCKER_HOST", "DOCKER_TLS_VERIFY", "DOCKER_CERT_PATH",
		"PENTAGI_DOCKER_SOCKET", "PENTAGI_DOCKER_CERT_PATH",
	}
	if err := c.ResetVars(vars); err != nil {
		return nil
	}
	return c.GetDockerConfig()
}

func (c *controller) GetServerSettingsConfig() *ServerSettingsConfig {
	vars, _ := c.GetVars([]string{
		"LICENSE_KEY", "PENTAGI_LISTEN_IP", "PENTAGI_LISTEN_PORT", "PUBLIC_URL", "CORS_ORIGINS",
		"COOKIE_SIGNING_SALT", "PROXY_URL", "HTTP_CLIENT_TIMEOUT", "EXTERNAL_SSL_CA_PATH", "EXTERNAL_SSL_INSECURE",
		"PENTAGI_SSL_DIR", "PENTAGI_DATA_DIR",
	})

	defaults := map[string]string{
		"LICENSE_KEY": "", "PENTAGI_LISTEN_IP": "127.0.0.1", "PENTAGI_LISTEN_PORT": "8443", "PUBLIC_URL": "https://localhost:8443",
		"CORS_ORIGINS": "https://localhost:8443", "PENTAGI_DATA_DIR": "pentagi-data", "PENTAGI_SSL_DIR": "pentagi-ssl",
		"HTTP_CLIENT_TIMEOUT": "600", "EXTERNAL_SSL_INSECURE": "false",
	}
	for varName, defaultValue := range defaults {
		if v := vars[varName]; v.Default == "" {
			v.Default = defaultValue
			vars[varName] = v
		}
	}

	cfg := &ServerSettingsConfig{
		LicenseKey:          vars["LICENSE_KEY"],
		ListenIP:            vars["PENTAGI_LISTEN_IP"],
		ListenPort:          vars["PENTAGI_LISTEN_PORT"],
		PublicURL:           vars["PUBLIC_URL"],
		CorsOrigins:         vars["CORS_ORIGINS"],
		CookieSigningSalt:   vars["COOKIE_SIGNING_SALT"],
		ProxyURL:            vars["PROXY_URL"],
		HTTPClientTimeout:   vars["HTTP_CLIENT_TIMEOUT"],
		ExternalSSLCAPath:   vars["EXTERNAL_SSL_CA_PATH"],
		ExternalSSLInsecure: vars["EXTERNAL_SSL_INSECURE"],
		SSLDir:              vars["PENTAGI_SSL_DIR"],
		DataDir:             vars["PENTAGI_DATA_DIR"],
	}

	if cfg.ProxyURL.Value != "" {
		user, pass := c.extractCredentialsFromURL(cfg.ProxyURL.Value)
		cfg.ProxyUsername = user
		cfg.ProxyPassword = pass
		cfg.ProxyURL.Value = RemoveCredentialsFromURL(cfg.ProxyURL.Value)
	}

	return cfg
}

func (c *controller) UpdateServerSettingsConfig(config *ServerSettingsConfig) error {
	if config == nil {
		return fmt.Errorf("config cannot be nil")
	}

	proxyURL := config.ProxyURL.Value
	if proxyURL != "" && config.ProxyUsername != "" && config.ProxyPassword != "" {
		proxyURL = c.addCredentialsToURL(proxyURL, config.ProxyUsername, config.ProxyPassword)
	}

	updates := map[string]string{
		"LICENSE_KEY": config.LicenseKey.Value, "PENTAGI_LISTEN_IP": config.ListenIP.Value, "PENTAGI_LISTEN_PORT": config.ListenPort.Value,
		"PUBLIC_URL": config.PublicURL.Value, "CORS_ORIGINS": config.CorsOrigins.Value, "COOKIE_SIGNING_SALT": config.CookieSigningSalt.Value,
		"PROXY_URL": proxyURL, "HTTP_CLIENT_TIMEOUT": config.HTTPClientTimeout.Value, "EXTERNAL_SSL_CA_PATH": config.ExternalSSLCAPath.Value,
		"EXTERNAL_SSL_INSECURE": config.ExternalSSLInsecure.Value, "PENTAGI_SSL_DIR": config.SSLDir.Value, "PENTAGI_DATA_DIR": config.DataDir.Value,
	}

	return c.SetVars(updates)
}

func (c *controller) ResetServerSettingsConfig() *ServerSettingsConfig {
	vars := []string{
		"LICENSE_KEY", "PENTAGI_LISTEN_IP", "PENTAGI_LISTEN_PORT", "PUBLIC_URL", "CORS_ORIGINS", "COOKIE_SIGNING_SALT",
		"PROXY_URL", "HTTP_CLIENT_TIMEOUT", "EXTERNAL_SSL_CA_PATH", "EXTERNAL_SSL_INSECURE", "PENTAGI_SSL_DIR", "PENTAGI_DATA_DIR",
	}
	if err := c.ResetVars(vars); err != nil {
		return nil
	}
	return c.GetServerSettingsConfig()
}
