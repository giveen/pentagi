package graph

import (
"context"
"database/sql"
"encoding/json"
"fmt"
"strings"
"time"

"pentagi/pkg/database"
"pentagi/pkg/database/converter"
"pentagi/pkg/graph/model"
mcpstdio "pentagi/pkg/mcp/stdio"
"pentagi/pkg/providers/pconfig"
"pentagi/pkg/providers/provider"
"pentagi/pkg/server/auth"
"pentagi/pkg/templates"
"pentagi/pkg/templates/validator"

"github.com/sirupsen/logrus"
)

// TestAgent is the resolver for the testAgent field.
func (r *mutationResolver) TestAgent(ctx context.Context, typeArg model.ProviderType, agentType model.AgentConfigType, agent model.AgentConfig) (*model.AgentTestResult, error) {
	uid, _, err := validatePermission(ctx, "settings.providers.view")
	if err != nil {
		return nil, err
	}

	r.Logger.WithFields(logrus.Fields{
		"uid":  uid,
		"type": typeArg.String(),
	}).Debug("test agent")

	cfg := converter.ConvertAgentConfigFromGqlModel(&agent)
	prvtype := provider.ProviderType(typeArg)
	atype := pconfig.ProviderOptionsType(agentType)
	result, err := r.ProvidersCtrl.TestAgent(ctx, prvtype, atype, cfg)
	if err != nil {
		return nil, err
	}

	return converter.ConvertTestResults(result), nil
}

// TestProvider is the resolver for the testProvider field.
func (r *mutationResolver) TestProvider(ctx context.Context, typeArg model.ProviderType, agents model.AgentsConfig, apiURL *string, apiKey *string) (*model.ProviderTestResult, error) {
	uid, _, err := validatePermission(ctx, "settings.providers.view")
	if err != nil {
		return nil, err
	}

	r.Logger.WithFields(logrus.Fields{
		"uid":  uid,
		"type": typeArg.String(),
	}).Debug("test provider")

	cfg := converter.ConvertAgentsConfigFromGqlModel(&agents)
	if apiURL != nil {
		cfg.APIURL = *apiURL
	}
	if apiKey != nil {
		cfg.APIKey = *apiKey
	}
	prvtype := provider.ProviderType(typeArg)
	result, err := r.ProvidersCtrl.TestProvider(ctx, prvtype, cfg)
	if err != nil {
		return nil, err
	}

	return converter.ConvertProviderTestResults(result), nil
}

// CreateProvider is the resolver for the createProvider field.
func (r *mutationResolver) CreateProvider(ctx context.Context, name string, typeArg model.ProviderType, agents model.AgentsConfig, apiURL *string, apiKey *string) (*model.ProviderConfig, error) {
	uid, _, err := validatePermission(ctx, "settings.providers.edit")
	if err != nil {
		return nil, err
	}

	r.Logger.WithFields(logrus.Fields{
		"uid":  uid,
		"name": name,
		"type": typeArg.String(),
	}).Debug("create provider")

	cfg := converter.ConvertAgentsConfigFromGqlModel(&agents)
	if apiURL != nil {
		cfg.APIURL = *apiURL
	}
	if apiKey != nil {
		cfg.APIKey = *apiKey
	}
	prvname, prvtype := provider.ProviderName(name), provider.ProviderType(typeArg)
	prv, err := r.ProvidersCtrl.CreateProvider(ctx, uid, prvname, prvtype, cfg)
	if err != nil {
		return nil, err
	}

	r.Subscriptions.NewFlowPublisher(uid, 0).ProviderCreated(ctx, prv, cfg)

	return converter.ConvertProvider(prv, cfg), nil
}

// UpdateProvider is the resolver for the updateProvider field.
func (r *mutationResolver) UpdateProvider(ctx context.Context, providerID int64, name string, agents model.AgentsConfig, apiURL *string, apiKey *string) (*model.ProviderConfig, error) {
	uid, _, err := validatePermission(ctx, "settings.providers.edit")
	if err != nil {
		return nil, err
	}

	r.Logger.WithFields(logrus.Fields{
		"uid":      uid,
		"provider": providerID,
		"name":     name,
	}).Debug("update provider")

	cfg := converter.ConvertAgentsConfigFromGqlModel(&agents)
	if apiURL != nil {
		cfg.APIURL = *apiURL
	}
	if apiKey != nil {
		cfg.APIKey = *apiKey
	}
	prvname := provider.ProviderName(name)
	prv, err := r.ProvidersCtrl.UpdateProvider(ctx, uid, providerID, prvname, cfg)
	if err != nil {
		return nil, err
	}

	r.Subscriptions.NewFlowPublisher(uid, 0).ProviderUpdated(ctx, prv, cfg)

	return converter.ConvertProvider(prv, cfg), nil
}

// DeleteProvider is the resolver for the deleteProvider field.
func (r *mutationResolver) DeleteProvider(ctx context.Context, providerID int64) (model.ResultType, error) {
	uid, _, err := validatePermission(ctx, "settings.providers.edit")
	if err != nil {
		return model.ResultTypeError, err
	}

	r.Logger.WithFields(logrus.Fields{
		"uid":      uid,
		"provider": providerID,
	}).Debug("delete provider")

	prv, err := r.ProvidersCtrl.DeleteProvider(ctx, uid, providerID)
	if err != nil {
		return model.ResultTypeError, err
	}

	var cfg pconfig.ProviderConfig
	if err := json.Unmarshal(prv.Config, &cfg); err != nil {
		return model.ResultTypeError, err
	}

	r.Subscriptions.NewFlowPublisher(uid, 0).ProviderDeleted(ctx, prv, &cfg)

	return model.ResultTypeSuccess, nil
}

// ValidatePrompt is the resolver for the validatePrompt field.
func (r *mutationResolver) ValidatePrompt(ctx context.Context, typeArg model.PromptType, template string) (*model.PromptValidationResult, error) {
	uid, _, err := validatePermission(ctx, "settings.prompts.edit")
	if err != nil {
		return nil, err
	}

	r.Logger.WithFields(logrus.Fields{
		"uid":      uid,
		"type":     typeArg.String(),
		"template": template[:min(len(template), 1000)],
	}).Debug("validate prompt")

	var (
		result    model.ResultType = model.ResultTypeSuccess
		errorType *model.PromptValidationErrorType
		message   *string
		line      *int
		details   *string
	)

	if err := validator.ValidatePrompt(templates.PromptType(typeArg), template); err != nil {
		result = model.ResultTypeError
		errType := model.PromptValidationErrorTypeUnknownType
		if err, ok := err.(*validator.ValidationError); ok {
			switch err.Type {
			case validator.ErrorTypeSyntax:
				errType = model.PromptValidationErrorTypeSyntaxError
			case validator.ErrorTypeUnauthorizedVar:
				errType = model.PromptValidationErrorTypeUnauthorizedVariable
			case validator.ErrorTypeRenderingFailed:
				errType = model.PromptValidationErrorTypeRenderingFailed
			case validator.ErrorTypeEmptyTemplate:
				errType = model.PromptValidationErrorTypeEmptyTemplate
			case validator.ErrorTypeVariableTypeMismatch:
				errType = model.PromptValidationErrorTypeVariableTypeMismatch
			}
			if err.Message != "" {
				message = &err.Message
			}
			if err.Line > 0 {
				line = &err.Line
			}
			if err.Details != "" {
				details = &err.Details
			}
		}
		errorType = &errType
	}

	return &model.PromptValidationResult{
		Result:    result,
		ErrorType: errorType,
		Message:   message,
		Line:      line,
		Details:   details,
	}, nil
}

// CreatePrompt is the resolver for the createPrompt field.
func (r *mutationResolver) CreatePrompt(ctx context.Context, typeArg model.PromptType, template string) (*model.UserPrompt, error) {
	uid, _, err := validatePermission(ctx, "settings.prompts.edit")
	if err != nil {
		return nil, err
	}

	r.Logger.WithFields(logrus.Fields{
		"uid":      uid,
		"type":     typeArg.String(),
		"template": template[:min(len(template), 1000)],
	}).Debug("create prompt")

	if err := validator.ValidatePrompt(templates.PromptType(typeArg), template); err != nil {
		return nil, err
	}

	prompt, err := r.DB.CreateUserPrompt(ctx, database.CreateUserPromptParams{
		UserID: uid,
		Type:   database.PromptType(typeArg),
		Prompt: template,
	})
	if err != nil {
		return nil, err
	}

	return converter.ConvertPrompt(prompt), nil
}

// UpdatePrompt is the resolver for the updatePrompt field.
func (r *mutationResolver) UpdatePrompt(ctx context.Context, promptID int64, template string) (*model.UserPrompt, error) {
	uid, _, err := validatePermission(ctx, "settings.prompts.edit")
	if err != nil {
		return nil, err
	}

	r.Logger.WithFields(logrus.Fields{
		"uid":      uid,
		"prompt":   promptID,
		"template": template[:min(len(template), 1000)],
	}).Debug("update prompt")

	prompt, err := r.DB.GetUserPrompt(ctx, database.GetUserPromptParams{
		ID:     promptID,
		UserID: uid,
	})
	if err != nil {
		return nil, err
	}

	if err := validator.ValidatePrompt(templates.PromptType(prompt.Type), template); err != nil {
		return nil, err
	}

	prompt, err = r.DB.UpdateUserPrompt(ctx, database.UpdateUserPromptParams{
		ID:     promptID,
		Prompt: template,
		UserID: uid,
	})
	if err != nil {
		return nil, err
	}

	return converter.ConvertPrompt(prompt), nil
}

// DeletePrompt is the resolver for the deletePrompt field.
func (r *mutationResolver) DeletePrompt(ctx context.Context, promptID int64) (model.ResultType, error) {
	uid, _, err := validatePermission(ctx, "settings.prompts.edit")
	if err != nil {
		return model.ResultTypeError, err
	}

	r.Logger.WithFields(logrus.Fields{
		"uid":    uid,
		"prompt": promptID,
	}).Debug("delete prompt")

	err = r.DB.DeleteUserPrompt(ctx, database.DeleteUserPromptParams{
		ID:     promptID,
		UserID: uid,
	})
	if err != nil {
		return model.ResultTypeError, err
	}

	return model.ResultTypeSuccess, nil
}

// CreateAPIToken is the resolver for the createAPIToken field.
func (r *mutationResolver) CreateAPIToken(ctx context.Context, input model.CreateAPITokenInput) (*model.APITokenWithSecret, error) {
	uid, _, err := validatePermission(ctx, "settings.tokens.create")
	if err != nil {
		return nil, err
	}

	isUserSession, err := validateUserType(ctx, userSessionTypes...)
	if err != nil {
		return nil, err
	}

	if !isUserSession {
		return nil, fmt.Errorf("unauthorized: non-user session is not allowed to create API tokens")
	}

	if r.Config.CookieSigningSalt == "" || r.Config.CookieSigningSalt == "salt" {
		return nil, fmt.Errorf("token creation is disabled with default salt")
	}

	if input.TTL < 60 || input.TTL > 94608000 {
		return nil, fmt.Errorf("invalid TTL: must be between 60 and 94608000 seconds")
	}

	r.Logger.WithFields(logrus.Fields{
		"uid":  uid,
		"name": input.Name,
		"ttl":  input.TTL,
	}).Debug("create api token")

	user, err := r.DB.GetUser(ctx, uid)
	if err != nil {
		return nil, err
	}

	tokenID, err := auth.GenerateTokenID()
	if err != nil {
		return nil, fmt.Errorf("failed to generate token ID: %w", err)
	}

	claims := auth.MakeAPITokenClaims(tokenID, user.Hash, uint64(uid), uint64(user.RoleID), uint64(input.TTL))

	tokenString, err := auth.MakeAPIToken(r.Config.CookieSigningSalt, claims)
	if err != nil {
		return nil, fmt.Errorf("failed to create token: %w", err)
	}

	var nameStr sql.NullString
	if input.Name != nil && *input.Name != "" {
		nameStr = sql.NullString{String: *input.Name, Valid: true}
	}

	apiToken, err := r.DB.CreateAPIToken(ctx, database.CreateAPITokenParams{
		TokenID: tokenID,
		UserID:  uid,
		RoleID:  user.RoleID,
		Name:    nameStr,
		Ttl:     int64(input.TTL),
		Status:  database.TokenStatusActive,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create token in database: %w", err)
	}

	tokenWithSecret := database.APITokenWithSecret{
		ApiToken: apiToken,
		Token:    tokenString,
	}

	r.TokenCache.Invalidate(tokenID)
	r.TokenCache.InvalidateUser(uint64(uid))

	r.Subscriptions.NewFlowPublisher(uid, 0).APITokenCreated(ctx, tokenWithSecret)

	return converter.ConvertAPITokenWithSecret(tokenWithSecret), nil
}

// UpdateAPIToken is the resolver for the updateAPIToken field.
func (r *mutationResolver) UpdateAPIToken(ctx context.Context, tokenID string, input model.UpdateAPITokenInput) (*model.APIToken, error) {
	uid, _, err := validatePermission(ctx, "settings.tokens.edit")
	if err != nil {
		return nil, err
	}

	isUserSession, err := validateUserType(ctx, userSessionTypes...)
	if err != nil {
		return nil, err
	}

	if !isUserSession {
		return nil, fmt.Errorf("unauthorized: non-user session is not allowed to update API tokens")
	}

	r.Logger.WithFields(logrus.Fields{
		"uid":     uid,
		"tokenID": tokenID,
	}).Debug("update api token")

	token, err := r.DB.GetUserAPITokenByTokenID(ctx, database.GetUserAPITokenByTokenIDParams{
		TokenID: tokenID,
		UserID:  uid,
	})
	if err != nil {
		return nil, fmt.Errorf("token not found: %w", err)
	}

	var nameStr sql.NullString
	if input.Name != nil {
		if *input.Name != "" {
			nameStr = sql.NullString{String: *input.Name, Valid: true}
		}
	} else {
		nameStr = token.Name
	}

	status := token.Status
	if input.Status != nil {
		switch s := *input.Status; s {
		case model.TokenStatusActive:
			status = database.TokenStatusActive
		case model.TokenStatusRevoked:
			status = database.TokenStatusRevoked
		default:
			return nil, fmt.Errorf("invalid token status: %s", s.String())
		}
	}

	updatedToken, err := r.DB.UpdateUserAPIToken(ctx, database.UpdateUserAPITokenParams{
		ID:     token.ID,
		UserID: uid,
		Name:   nameStr,
		Status: status,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to update token: %w", err)
	}

	if input.Status != nil {
		r.TokenCache.Invalidate(tokenID)
		r.TokenCache.InvalidateUser(uint64(uid))
	}

	r.Subscriptions.NewFlowPublisher(uid, 0).APITokenUpdated(ctx, updatedToken)

	return converter.ConvertAPIToken(updatedToken), nil
}

// DeleteAPIToken is the resolver for the deleteAPIToken field.
func (r *mutationResolver) DeleteAPIToken(ctx context.Context, tokenID string) (bool, error) {
	uid, _, err := validatePermission(ctx, "settings.tokens.delete")
	if err != nil {
		return false, err
	}

	isUserSession, err := validateUserType(ctx, userSessionTypes...)
	if err != nil {
		return false, err
	}

	if !isUserSession {
		return false, fmt.Errorf("unauthorized: non-user session is not allowed to delete API tokens")
	}

	r.Logger.WithFields(logrus.Fields{
		"uid":     uid,
		"tokenID": tokenID,
	}).Debug("delete api token")

	token, err := r.DB.DeleteUserAPITokenByTokenID(ctx, database.DeleteUserAPITokenByTokenIDParams{
		TokenID: tokenID,
		UserID:  uid,
	})
	if err != nil {
		return false, fmt.Errorf("failed to delete token: %w", err)
	}

	r.TokenCache.Invalidate(tokenID)
	r.TokenCache.InvalidateUser(uint64(uid))

	r.Subscriptions.NewFlowPublisher(uid, 0).APITokenDeleted(ctx, token)

	return true, nil
}

// AddFavoriteFlow is the resolver for the addFavoriteFlow field.
func (r *mutationResolver) AddFavoriteFlow(ctx context.Context, flowID int64) (model.ResultType, error) {
	_, err := validatePermissionWithFlowID(ctx, "flows.view", flowID, r.DB)
	if err != nil {
		return model.ResultTypeError, err
	}

	uid, _, err := validatePermission(ctx, "settings.user.edit")
	if err != nil {
		return model.ResultTypeError, err
	}

	isUserSession, err := validateUserType(ctx, userSessionTypes...)
	if err != nil {
		return model.ResultTypeError, err
	}

	if !isUserSession {
		return model.ResultTypeError, fmt.Errorf("unauthorized: non-user session is not allowed to manage favorites")
	}

	r.Logger.WithFields(logrus.Fields{
		"uid":    uid,
		"flowID": flowID,
	}).Debug("add favorite flow")

	flow, err := r.DB.GetFlow(ctx, flowID)
	if err != nil {
		return model.ResultTypeError, fmt.Errorf("flow not found: %w", err)
	}

	if flow.UserID != uid {
		return model.ResultTypeError, fmt.Errorf("unauthorized: cannot favorite other user's flow")
	}

	prefs, err := r.DB.AddFavoriteFlow(ctx, database.AddFavoriteFlowParams{
		UserID: uid,
		FlowID: flowID,
	})
	if err != nil {
		return model.ResultTypeError, fmt.Errorf("failed to add favorite flow: %w", err)
	}

	r.Subscriptions.NewFlowPublisher(uid, 0).SettingsUserUpdated(ctx, prefs)

	return model.ResultTypeSuccess, nil
}

// DeleteFavoriteFlow is the resolver for the deleteFavoriteFlow field.
func (r *mutationResolver) DeleteFavoriteFlow(ctx context.Context, flowID int64) (model.ResultType, error) {
	_, err := validatePermissionWithFlowID(ctx, "flows.view", flowID, r.DB)
	if err != nil {
		return model.ResultTypeError, err
	}

	uid, err := validatePermissionWithFlowID(ctx, "settings.user.edit", flowID, r.DB)
	if err != nil {
		return model.ResultTypeError, err
	}

	isUserSession, err := validateUserType(ctx, userSessionTypes...)
	if err != nil {
		return model.ResultTypeError, err
	}

	if !isUserSession {
		return model.ResultTypeError, fmt.Errorf("unauthorized: non-user session is not allowed to manage favorites")
	}

	r.Logger.WithFields(logrus.Fields{
		"uid":    uid,
		"flowID": flowID,
	}).Debug("delete favorite flow")

	prefs, err := r.DB.DeleteFavoriteFlow(ctx, database.DeleteFavoriteFlowParams{
		FlowID: flowID,
		UserID: uid,
	})
	if err != nil {
		return model.ResultTypeError, fmt.Errorf("failed to delete favorite flow: %w", err)
	}

	r.Subscriptions.NewFlowPublisher(uid, 0).SettingsUserUpdated(ctx, prefs)

	return model.ResultTypeSuccess, nil
}

// CreateFlowTemplate is the resolver for the createFlowTemplate field.
func (r *mutationResolver) CreateFlowTemplate(ctx context.Context, input model.CreateFlowTemplateInput) (*model.FlowTemplate, error) {
	uid, _, err := validatePermission(ctx, "templates.create")
	if err != nil {
		return nil, err
	}

	isUserSession, err := validateUserType(ctx, userSessionTypes...)
	if err != nil {
		return nil, err
	}

	if !isUserSession {
		return nil, fmt.Errorf("unauthorized: non-user session is not allowed to create templates")
	}

	r.Logger.WithFields(logrus.Fields{
		"uid":   uid,
		"title": input.Title,
	}).Debug("create flow template")

	template, err := r.DB.CreateFlowTemplate(ctx, database.CreateFlowTemplateParams{
		UserID: uid,
		Title:  input.Title,
		Text:   input.Text,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create template: %w", err)
	}

	r.Subscriptions.NewFlowPublisher(uid, 0).FlowTemplateCreated(ctx, template)

	return converter.ConvertFlowTemplate(template), nil
}

// UpdateFlowTemplate is the resolver for the updateFlowTemplate field.
func (r *mutationResolver) UpdateFlowTemplate(ctx context.Context, templateID int64, input model.UpdateFlowTemplateInput) (*model.FlowTemplate, error) {
	uid, _, err := validatePermission(ctx, "templates.edit")
	if err != nil {
		return nil, err
	}

	isUserSession, err := validateUserType(ctx, userSessionTypes...)
	if err != nil {
		return nil, err
	}

	if !isUserSession {
		return nil, fmt.Errorf("unauthorized: non-user session is not allowed to update templates")
	}

	r.Logger.WithFields(logrus.Fields{
		"uid":        uid,
		"templateID": templateID,
	}).Debug("update flow template")

	_, err = r.DB.GetFlowTemplate(ctx, database.GetFlowTemplateParams{
		ID:     templateID,
		UserID: uid,
	})
	if err != nil {
		return nil, fmt.Errorf("template not found: %w", err)
	}

	template, err := r.DB.UpdateFlowTemplate(ctx, database.UpdateFlowTemplateParams{
		ID:     templateID,
		UserID: uid,
		Title:  input.Title,
		Text:   input.Text,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to update template: %w", err)
	}

	r.Subscriptions.NewFlowPublisher(uid, 0).FlowTemplateUpdated(ctx, template)

	return converter.ConvertFlowTemplate(template), nil
}

// DeleteFlowTemplate is the resolver for the deleteFlowTemplate field.
func (r *mutationResolver) DeleteFlowTemplate(ctx context.Context, templateID int64) (model.ResultType, error) {
	uid, _, err := validatePermission(ctx, "templates.delete")
	if err != nil {
		return model.ResultTypeError, err
	}

	isUserSession, err := validateUserType(ctx, userSessionTypes...)
	if err != nil {
		return model.ResultTypeError, err
	}

	if !isUserSession {
		return model.ResultTypeError, fmt.Errorf("unauthorized: non-user session is not allowed to delete templates")
	}

	r.Logger.WithFields(logrus.Fields{
		"uid":        uid,
		"templateID": templateID,
	}).Debug("delete flow template")

	template, err := r.DB.GetFlowTemplate(ctx, database.GetFlowTemplateParams{
		ID:     templateID,
		UserID: uid,
	})
	if err != nil {
		return model.ResultTypeError, fmt.Errorf("template not found: %w", err)
	}

	err = r.DB.DeleteFlowTemplate(ctx, database.DeleteFlowTemplateParams{
		ID:     templateID,
		UserID: uid,
	})
	if err != nil {
		return model.ResultTypeError, fmt.Errorf("failed to delete template: %w", err)
	}

	r.Subscriptions.NewFlowPublisher(uid, 0).FlowTemplateDeleted(ctx, template)

	return model.ResultTypeSuccess, nil
}

// CreateMcpServer is the resolver for the createMcpServer field.
func (r *mutationResolver) CreateMcpServer(ctx context.Context, input model.CreateMcpServerInput) (*model.McpServer, error) {
	_, _, err := validatePermission(ctx, "settings.mcp.create")
	if err != nil {
		return nil, err
	}

	// build DB params
	var stdioEnv json.RawMessage
	var sseHeaders json.RawMessage
	var tools json.RawMessage

	if input.Stdio != nil {
		if b, err := json.Marshal(input.Stdio.Env); err == nil {
			stdioEnv = b
		}
	}
	if input.Sse != nil {
		if b, err := json.Marshal(input.Sse.Headers); err == nil {
			sseHeaders = b
		}
	}
	if input.Tools != nil {
		if b, err := json.Marshal(input.Tools); err == nil {
			tools = b
		}
	}

	params := database.CreateMcpServerParams{
		Name:      input.Name,
		Transport: database.McpTransport(input.Transport),
		StdioCommand: func() sql.NullString {
			if input.Stdio != nil {
				return sql.NullString{String: input.Stdio.Command, Valid: true}
			}
			return sql.NullString{}
		}(),
		StdioArgs: func() sql.NullString {
			if input.Stdio != nil && input.Stdio.Args != nil {
				return sql.NullString{String: *input.Stdio.Args, Valid: true}
			}
			return sql.NullString{}
		}(),
		StdioEnv: stdioEnv,
		SseUrl: func() sql.NullString {
			if input.Sse != nil {
				return sql.NullString{String: input.Sse.URL, Valid: true}
			}
			return sql.NullString{}
		}(),
		SseHeaders: sseHeaders,
		Tools:      tools,
	}

	srv, err := r.DB.CreateMcpServer(ctx, params)
	if err != nil {
		return nil, err
	}

	// convert to GraphQL model
	return converter.ConvertMcpServer(srv), nil
}

// UpdateMcpServer is the resolver for the updateMcpServer field.
func (r *mutationResolver) UpdateMcpServer(ctx context.Context, mcpServerID int64, input model.UpdateMcpServerInput) (*model.McpServer, error) {
	_, _, err := validatePermission(ctx, "settings.mcp.edit")
	if err != nil {
		return nil, err
	}

	// fetch existing server to fill missing fields
	existing, err := r.DB.GetMcpServer(ctx, mcpServerID)
	if err != nil {
		return nil, err
	}

	// compute fields: if input value provided use it, otherwise keep existing
	name := existing.Name
	if input.Name != nil {
		name = *input.Name
	}

	var transport database.McpTransport = existing.Transport
	if input.Transport != nil {
		transport = database.McpTransport(*input.Transport)
	}

	var stdioCmd sql.NullString = existing.StdioCommand
	var stdioArgs sql.NullString = existing.StdioArgs
	var stdioEnv json.RawMessage = existing.StdioEnv
	if input.Stdio != nil {
		stdioCmd = sql.NullString{String: input.Stdio.Command, Valid: true}
		if input.Stdio.Args != nil {
			stdioArgs = sql.NullString{String: *input.Stdio.Args, Valid: true}
		} else {
			stdioArgs = sql.NullString{}
		}
		if b, err := json.Marshal(input.Stdio.Env); err == nil {
			stdioEnv = b
		}
	}

	var sseUrl sql.NullString = existing.SseUrl
	var sseHeaders json.RawMessage = existing.SseHeaders
	if input.Sse != nil {
		sseUrl = sql.NullString{String: input.Sse.URL, Valid: true}
		if b, err := json.Marshal(input.Sse.Headers); err == nil {
			sseHeaders = b
		}
	}

	var tools json.RawMessage = existing.Tools
	if input.Tools != nil {
		if b, err := json.Marshal(input.Tools); err == nil {
			tools = b
		}
	}

	params := database.UpdateMcpServerParams{
		Name:         name,
		Transport:    transport,
		StdioCommand: stdioCmd,
		StdioArgs:    stdioArgs,
		StdioEnv:     stdioEnv,
		SseUrl:       sseUrl,
		SseHeaders:   sseHeaders,
		Tools:        tools,
		ID:           mcpServerID,
	}

	srv, err := r.DB.UpdateMcpServer(ctx, params)
	if err != nil {
		return nil, err
	}
	return converter.ConvertMcpServer(srv), nil
}

// DeleteMcpServer is the resolver for the deleteMcpServer field.
func (r *mutationResolver) DeleteMcpServer(ctx context.Context, mcpServerID int64) (model.ResultType, error) {
	_, _, err := validatePermission(ctx, "settings.mcp.delete")
	if err != nil {
		return model.ResultTypeError, err
	}

	if err := r.DB.DeleteMcpServer(ctx, mcpServerID); err != nil {
		return model.ResultTypeError, err
	}

	return model.ResultTypeSuccess, nil
}

// TestMcpServer is the resolver for the testMcpServer field.
func (r *mutationResolver) TestMcpServer(ctx context.Context, mcpServerID int64) (model.ResultType, error) {
	_, _, err := validatePermission(ctx, "settings.mcp.view")
	if err != nil {
		return model.ResultTypeError, err
	}

	// Retrieve server config
	srv, err := r.DB.GetMcpServer(ctx, mcpServerID)
	if err != nil {
		return model.ResultTypeError, err
	}

	switch srv.Transport {
	case database.McpTransportStdio:
		if !srv.StdioCommand.Valid {
			return model.ResultTypeError, fmt.Errorf("stdio transport configured without command")
		}
		cmd := srv.StdioCommand.String
		var args []string
		if srv.StdioArgs.Valid {
			args = strings.Fields(srv.StdioArgs.String)
		}

		envMap := map[string]string{}
		if len(srv.StdioEnv) > 0 {
			var env []struct {
				Key   string `json:"key"`
				Value string `json:"value"`
			}
			if err := json.Unmarshal(srv.StdioEnv, &env); err == nil {
				for _, e := range env {
					envMap[e.Key] = e.Value
				}
			}
		}

		// run a short test ping
		timeout := 10 * time.Second
		if err := mcpstdio.Test(ctx, cmd, args, envMap, timeout); err != nil {
			return model.ResultTypeError, fmt.Errorf("stdio connector test failed: %w", err)
		}
		return model.ResultTypeSuccess, nil
	default:
		// Other transports not implemented yet
		return model.ResultTypeError, fmt.Errorf("unsupported transport: %s", srv.Transport)
	}
}
