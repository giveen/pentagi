package graph

import (
"context"
"database/sql"
"encoding/json"
"fmt"
"time"

"pentagi/pkg/database"
"pentagi/pkg/database/converter"
"pentagi/pkg/graph/model"
"pentagi/pkg/providers/openai"
"pentagi/pkg/providers/pconfig"
"pentagi/pkg/providers/provider"
"pentagi/pkg/templates"

"github.com/sirupsen/logrus"
)

// Settings is the resolver for the settings field.
func (r *queryResolver) Settings(ctx context.Context) (*model.Settings, error) {
	_, _, err := validatePermission(ctx, "settings.view")
	if err != nil {
		return nil, err
	}

	settings := &model.Settings{
		Debug:              r.Config.Debug,
		AskUser:            r.Config.AskUser,
		DockerInside:       r.Config.DockerInside,
		AssistantUseAgents: r.Config.AssistantUseAgents,
	}

	return settings, nil
}

// SettingsProviders is the resolver for the settingsProviders field.
func (r *queryResolver) SettingsProviders(ctx context.Context) (*model.ProvidersConfig, error) {
	uid, _, err := validatePermission(ctx, "settings.providers.view")
	if err != nil {
		return nil, err
	}

	r.Logger.WithFields(logrus.Fields{
		"uid": uid,
	}).Debug("get providers")

	config := model.ProvidersConfig{
		Enabled:     &model.ProvidersReadinessStatus{},
		Default:     &model.DefaultProvidersConfig{},
		Models:      &model.ProvidersModelsList{},
		UserDefined: make([]*model.ProviderConfig, 0),
	}

	now := time.Now()
	defaultProvidersConfig := r.ProvidersCtrl.DefaultProvidersConfig()
	for prvtype, pcfg := range defaultProvidersConfig {
		var apiURL *string
		var apiKey *string
		if pcfg != nil {
			if pcfg.APIURL != "" {
				apiURL = &pcfg.APIURL
			}
			if pcfg.APIKey != "" {
				apiKey = &pcfg.APIKey
			}
		}
		mpcfg := &model.ProviderConfig{
			Name:      string(prvtype),
			Type:      model.ProviderType(prvtype),
			APIURL:    apiURL,
			APIKey:    apiKey,
			Agents:    converter.ConvertProviderConfigToGqlModel(pcfg),
			CreatedAt: now,
			UpdatedAt: now,
		}

		switch prvtype {
		case provider.ProviderOpenAI:
			config.Default.Openai = mpcfg
			if models, err := openai.DefaultModels(); err == nil {
				config.Models.Openai = converter.ConvertModels(models)
			}
		case provider.ProviderOllama:
			config.Default.Ollama = mpcfg
		case provider.ProviderCustom:
			config.Default.Custom = mpcfg
		}
	}

	defaultProviders := r.ProvidersCtrl.DefaultProviders()
	for _, prvtype := range defaultProviders.ListTypes() {
		switch prvtype {
		case provider.ProviderOpenAI:
			config.Enabled.Openai = true
		case provider.ProviderOllama:
			config.Enabled.Ollama = true
			if p, ok := defaultProviders[provider.DefaultProviderNameOllama]; ok {
				config.Models.Ollama = converter.ConvertModels(p.GetModels())
			}
		case provider.ProviderCustom:
			config.Enabled.Custom = true
			if p, ok := defaultProviders[provider.DefaultProviderNameCustom]; ok {
				config.Models.Custom = converter.ConvertModels(p.GetModels())
			}
		}
	}

	providers, err := r.DB.GetUserProviders(ctx, uid)
	if err != nil {
		return nil, fmt.Errorf("failed to get user providers: %w", err)
	}

	for _, prv := range providers {
		var cfg pconfig.ProviderConfig

		if len(prv.Config) == 0 {
			prv.Config = []byte(pconfig.EmptyProviderConfigRaw)
		}
		if err := json.Unmarshal(prv.Config, &cfg); err != nil {
			r.Logger.WithError(err).Errorf("failed to unmarshal provider config: %s", prv.Config)
			continue
		}

		config.UserDefined = append(config.UserDefined, converter.ConvertProvider(prv, &cfg))
	}

	return &config, nil
}

// SettingsPrompts is the resolver for the settingsPrompts field.
func (r *queryResolver) SettingsPrompts(ctx context.Context) (*model.PromptsConfig, error) {
	uid, _, err := validatePermission(ctx, "settings.prompts.view")
	if err != nil {
		return nil, err
	}

	r.Logger.WithFields(logrus.Fields{
		"uid": uid,
	}).Debug("get prompts")

	prompts, err := r.DB.GetUserPrompts(ctx, uid)
	if err != nil {
		return nil, err
	}

	defaultPrompts, err := templates.GetDefaultPrompts()
	if err != nil {
		return nil, err
	}

	promptsConfig := model.PromptsConfig{
		Default:     converter.ConvertDefaultPrompts(defaultPrompts),
		UserDefined: make([]*model.UserPrompt, 0, len(prompts)),
	}

	for _, prompt := range prompts {
		promptsConfig.UserDefined = append(promptsConfig.UserDefined, converter.ConvertPrompt(prompt))
	}

	return &promptsConfig, nil
}

// SettingsUser is the resolver for the settingsUser field.
func (r *queryResolver) SettingsUser(ctx context.Context) (*model.UserPreferences, error) {
	uid, _, err := validatePermission(ctx, "settings.user.view")
	if err != nil {
		return nil, err
	}

	isUserSession, err := validateUserType(ctx, userSessionTypes...)
	if err != nil {
		return nil, err
	}

	if !isUserSession {
		return nil, fmt.Errorf("unauthorized: non-user session is not allowed to get user preferences")
	}

	r.Logger.WithFields(logrus.Fields{
		"uid": uid,
	}).Debug("get user preferences")

	prefs, err := r.DB.GetUserPreferencesByUserID(ctx, uid)
	if err != nil {
		if err == sql.ErrNoRows {
			return &model.UserPreferences{
				FavoriteFlows: []int64{},
			}, nil
		}
		return nil, fmt.Errorf("failed to get user preferences: %w", err)
	}

	return converter.ConvertUserPreferences(prefs), nil
}

// APIToken is the resolver for the apiToken field.
func (r *queryResolver) APIToken(ctx context.Context, tokenID string) (*model.APIToken, error) {
	uid, admin, err := validatePermission(ctx, "settings.tokens.view")
	if err != nil {
		return nil, err
	}

	isUserSession, err := validateUserType(ctx, userSessionTypes...)
	if err != nil {
		return nil, err
	}

	if !isUserSession {
		return nil, fmt.Errorf("unauthorized: non-user session is not allowed to get API tokens")
	}

	r.Logger.WithFields(logrus.Fields{
		"uid":     uid,
		"tokenID": tokenID,
	}).Debug("get api token")

	var token database.ApiToken

	if admin {
		token, err = r.DB.GetAPITokenByTokenID(ctx, tokenID)
	} else {
		token, err = r.DB.GetUserAPITokenByTokenID(ctx, database.GetUserAPITokenByTokenIDParams{
			TokenID: tokenID,
			UserID:  uid,
		})
	}
	if err != nil {
		return nil, fmt.Errorf("token not found: %w", err)
	}

	return converter.ConvertAPIToken(token), nil
}

// APITokens is the resolver for the apiTokens field.
func (r *queryResolver) APITokens(ctx context.Context) ([]*model.APIToken, error) {
	uid, admin, err := validatePermission(ctx, "settings.tokens.view")
	if err != nil {
		return nil, err
	}

	isUserSession, err := validateUserType(ctx, userSessionTypes...)
	if err != nil {
		return nil, err
	}

	if !isUserSession {
		return nil, fmt.Errorf("unauthorized: non-user session is not allowed to get API tokens")
	}

	r.Logger.WithFields(logrus.Fields{
		"uid": uid,
	}).Debug("get api tokens")

	var tokens []database.ApiToken

	if admin {
		tokens, err = r.DB.GetAPITokens(ctx)
	} else {
		tokens, err = r.DB.GetUserAPITokens(ctx, uid)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get tokens: %w", err)
	}

	return converter.ConvertAPITokens(tokens), nil
}

// FlowTemplate is the resolver for the flowTemplate field.
func (r *queryResolver) FlowTemplate(ctx context.Context, templateID int64) (*model.FlowTemplate, error) {
	uid, _, err := validatePermission(ctx, "templates.view")
	if err != nil {
		return nil, err
	}

	isUserSession, err := validateUserType(ctx, userSessionTypes...)
	if err != nil {
		return nil, err
	}

	if !isUserSession {
		return nil, fmt.Errorf("unauthorized: non-user session is not allowed to view templates")
	}

	r.Logger.WithFields(logrus.Fields{
		"uid":        uid,
		"templateID": templateID,
	}).Debug("get flow template")

	template, err := r.DB.GetFlowTemplate(ctx, database.GetFlowTemplateParams{
		ID:     templateID,
		UserID: uid,
	})
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("template not found")
		}
		return nil, fmt.Errorf("failed to get template: %w", err)
	}

	return converter.ConvertFlowTemplate(template), nil
}

// FlowTemplates is the resolver for the flowTemplates field.
func (r *queryResolver) FlowTemplates(ctx context.Context) ([]*model.FlowTemplate, error) {
	uid, _, err := validatePermission(ctx, "templates.view")
	if err != nil {
		return nil, err
	}

	isUserSession, err := validateUserType(ctx, userSessionTypes...)
	if err != nil {
		return nil, err
	}

	if !isUserSession {
		return nil, fmt.Errorf("unauthorized: non-user session is not allowed to view templates")
	}

	r.Logger.WithFields(logrus.Fields{
		"uid": uid,
	}).Debug("get flow templates")

	templates, err := r.DB.GetFlowTemplatesByUserID(ctx, uid)
	if err != nil {
		return nil, fmt.Errorf("failed to get templates: %w", err)
	}

	return converter.ConvertFlowTemplates(templates), nil
}

// McpServers is the resolver for the mcpServers field.
func (r *queryResolver) McpServers(ctx context.Context) ([]*model.McpServer, error) {
	_, _, err := validatePermission(ctx, "settings.mcp.view")
	if err != nil {
		return nil, err
	}

	servers, err := r.DB.GetMcpServers(ctx)
	if err != nil {
		return nil, err
	}
	return converter.ConvertMcpServers(servers), nil
}

// McpServer is the resolver for the mcpServer field.
func (r *queryResolver) McpServer(ctx context.Context, mcpServerID int64) (*model.McpServer, error) {
	_, _, err := validatePermission(ctx, "settings.mcp.view")
	if err != nil {
		return nil, err
	}

	srv, err := r.DB.GetMcpServer(ctx, mcpServerID)
	if err != nil {
		return nil, err
	}
	return converter.ConvertMcpServer(srv), nil
}
