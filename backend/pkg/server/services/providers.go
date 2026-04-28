package services

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"slices"
	"strings"
	"time"

	"pentagi/pkg/config"
	"pentagi/pkg/providers"
	"pentagi/pkg/providers/pconfig"
	"pentagi/pkg/providers/provider"
	"pentagi/pkg/server/logger"
	"pentagi/pkg/server/models"
	"pentagi/pkg/server/response"
	"pentagi/pkg/system"

	"github.com/gin-gonic/gin"
)

type ProviderService struct {
	providers providers.ProviderController
	cfg       *config.Config
}

func NewProviderService(providers providers.ProviderController, cfg *config.Config) *ProviderService {
	return &ProviderService{
		providers: providers,
		cfg:       cfg,
	}
}

// GetProviders is a function to return providers list
// @Summary Retrieve providers list
// @Tags Providers
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.successResp{data=models.ProviderInfo} "providers list received successful"
// @Failure 403 {object} response.errorResp "getting providers not permitted"
// @Router /providers/ [get]
func (s *ProviderService) GetProviders(c *gin.Context) {
	privs := c.GetStringSlice("prm")
	if !slices.Contains(privs, "providers.view") {
		logger.FromContext(c).Errorf("error filtering user role permissions: permission not found")
		response.Error(c, response.ErrNotPermitted, nil)
		return
	}

	providers, err := s.providers.GetProviders(c, int64(c.GetUint64("uid")))
	if err != nil {
		logger.FromContext(c).Errorf("error getting providers: %v", err)
		response.Error(c, response.ErrInternal, nil)
		return
	}

	providerInfos := make([]models.ProviderInfo, len(providers))
	for i, name := range providers.ListNames() {
		providerInfos[i] = models.ProviderInfo{
			Name: name.String(),
			Type: models.ProviderType(providers[name].Type()),
		}
	}

	response.Success(c, http.StatusOK, providerInfos)
}

// CheckHealth validates provider endpoint connectivity and embedding endpoint readiness.
// @Summary Check provider and embedding endpoint health
// @Tags Providers
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param json body models.ProviderHealthCheckRequest true "provider endpoint health check payload"
// @Success 200 {object} response.successResp{data=models.ProviderHealthCheckResponse} "provider endpoint health received successfully"
// @Failure 400 {object} response.errorResp "invalid provider health request"
// @Failure 403 {object} response.errorResp "checking provider health not permitted"
// @Failure 500 {object} response.errorResp "internal error on provider health check"
// @Router /providers/health [post]
func (s *ProviderService) CheckHealth(c *gin.Context) {
	privs := c.GetStringSlice("prm")
	if !slices.Contains(privs, "providers.view") {
		logger.FromContext(c).Errorf("error filtering user role permissions: permission not found")
		response.Error(c, response.ErrNotPermitted, nil)
		return
	}

	var req models.ProviderHealthCheckRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.FromContext(c).WithError(err).Errorf("error binding provider health payload")
		response.Error(c, response.ErrProvidersInvalidRequest, err)
		return
	}
	if err := req.Valid(); err != nil {
		logger.FromContext(c).WithError(err).Errorf("error validating provider health payload")
		response.Error(c, response.ErrProvidersInvalidData, err)
		return
	}

	httpClient, err := system.GetHTTPClient(s.cfg)
	if err != nil {
		logger.FromContext(c).WithError(err).Errorf("error creating http client for provider health check")
		response.Error(c, response.ErrInternal, err)
		return
	}

	providerURL, providerKey := s.resolveProviderEndpoint(req)
	providerHealth := s.checkProviderEndpoint(c, provider.ProviderType(req.Type), providerURL, providerKey, httpClient)
	embeddingProvider, embeddingHealth := s.checkEmbeddingEndpoint(
		c,
		req,
		provider.ProviderType(req.Type),
		providerURL,
		providerKey,
		httpClient,
	)

	resp := models.ProviderHealthCheckResponse{
		ProviderType:      req.Type,
		Provider:          providerHealth,
		EmbeddingProvider: embeddingProvider,
		Embedding:         embeddingHealth,
	}

	response.Success(c, http.StatusOK, resp)
}

func (s *ProviderService) resolveProviderEndpoint(req models.ProviderHealthCheckRequest) (string, string) {
	apiURL := strings.TrimSpace(ptrToString(req.APIURL))
	apiKey := strings.TrimSpace(ptrToString(req.APIKey))

	switch provider.ProviderType(req.Type) {
	case provider.ProviderOpenAI:
		if apiURL == "" {
			apiURL = strings.TrimSpace(s.cfg.OpenAIServerURL)
		}
		if apiKey == "" {
			apiKey = strings.TrimSpace(s.cfg.OpenAIKey)
		}
	case provider.ProviderOllama:
		if apiURL == "" {
			apiURL = strings.TrimSpace(s.cfg.OllamaServerURL)
		}
		if apiKey == "" {
			apiKey = strings.TrimSpace(s.cfg.OllamaServerAPIKey)
		}
	case provider.ProviderCustom:
		if apiURL == "" {
			apiURL = strings.TrimSpace(s.cfg.LLMServerURL)
		}
		if apiKey == "" {
			apiKey = strings.TrimSpace(s.cfg.LLMServerKey)
		}
	}

	return apiURL, apiKey
}

func (s *ProviderService) checkProviderEndpoint(
	c *gin.Context,
	providerType provider.ProviderType,
	apiURL, apiKey string,
	httpClient *http.Client,
) models.EndpointHealth {
	result := models.EndpointHealth{URL: apiURL}
	if apiURL == "" {
		result.Reachable = false
		result.Error = ptr("provider endpoint URL is empty")
		return result
	}

	switch providerType {
	case provider.ProviderOllama:
		statusCode, modelsCount, err := checkOllamaTags(apiURL, apiKey, httpClient)
		if statusCode != nil {
			result.StatusCode = statusCode
		}
		if modelsCount != nil {
			result.Models = modelsCount
		}
		if err != nil {
			result.Reachable = false
			errMsg := err.Error()
			result.Error = &errMsg
			return result
		}
		result.Reachable = true
		return result
	default:
		providerPrefix := ""
		if providerType == provider.ProviderCustom {
			providerPrefix = s.cfg.LLMServerProvider
		}

		models, err := provider.LoadModelsFromHTTP(apiURL, apiKey, httpClient, providerPrefix)
		if err != nil {
			result.Reachable = false
			errMsg := err.Error()
			result.Error = &errMsg
			logger.FromContext(c).WithError(err).Warn("provider models endpoint check failed")
			return result
		}
		count := len(models)
		result.Models = &count
		result.Reachable = true
		return result
	}
}

func (s *ProviderService) checkEmbeddingEndpoint(
	c *gin.Context,
	req models.ProviderHealthCheckRequest,
	providerType provider.ProviderType,
	providerURL, providerKey string,
	httpClient *http.Client,
) (string, models.EndpointHealth) {
	providerName := strings.TrimSpace(strings.ToLower(s.cfg.EmbeddingProvider))
	result := models.EndpointHealth{}
	embeddingModel := strings.TrimSpace(s.cfg.EmbeddingModel)
	if req.EmbeddingModel != nil && strings.TrimSpace(*req.EmbeddingModel) != "" {
		embeddingModel = strings.TrimSpace(*req.EmbeddingModel)
		providerName = strings.ToLower(string(providerType))
		result.URL = providerURL
	}

	if embeddingModel != "" {
		model := embeddingModel
		result.Model = &model
	}

	if providerName == "none" {
		result.Reachable = false
		result.Error = ptr("embedding provider is disabled")
		return providerName, result
	}

	if req.EmbeddingModel != nil && strings.TrimSpace(*req.EmbeddingModel) != "" {
		if providerURL == "" {
			result.Reachable = false
			result.Error = ptr("embedding endpoint URL is empty")
			return providerName, result
		}

		switch providerType {
		case provider.ProviderOllama:
			statusCode, modelsCount, modelNames, err := checkOllamaTagsWithNames(providerURL, providerKey, httpClient)
			if statusCode != nil {
				result.StatusCode = statusCode
			}
			if modelsCount != nil {
				result.Models = modelsCount
			}
			if err != nil {
				result.Reachable = false
				errMsg := err.Error()
				result.Error = &errMsg
				return providerName, result
			}
			if result.Model != nil && !containsModel(*result.Model, modelNames) {
				result.Reachable = false
				errMsg := fmt.Sprintf("embedding model '%s' was not found in provider model list", *result.Model)
				result.Error = &errMsg
				return providerName, result
			}
			result.Reachable = true
			return providerName, result
		default:
			providerPrefix := ""
			if providerType == provider.ProviderCustom {
				providerPrefix = s.cfg.LLMServerProvider
			}
			models, err := provider.LoadModelsFromHTTP(providerURL, providerKey, httpClient, providerPrefix)
			if err != nil {
				result.Reachable = false
				errMsg := err.Error()
				result.Error = &errMsg
				return providerName, result
			}
			count := len(models)
			result.Models = &count
			if result.Model != nil && !containsConfiguredModel(*result.Model, models) {
				result.Reachable = false
				errMsg := fmt.Sprintf("embedding model '%s' was not found in provider model list", *result.Model)
				result.Error = &errMsg
				return providerName, result
			}
			result.Reachable = true
			return providerName, result
		}
	}

	endpoint := strings.TrimSpace(s.cfg.EmbeddingURL)
	switch providerName {
	case "openai":
		if endpoint == "" {
			endpoint = strings.TrimSpace(s.cfg.OpenAIServerURL)
		}
		key := strings.TrimSpace(s.cfg.EmbeddingKey)
		if key == "" {
			key = strings.TrimSpace(s.cfg.OpenAIKey)
		}
		result.URL = endpoint
		if endpoint == "" {
			result.Reachable = false
			result.Error = ptr("embedding endpoint URL is empty")
			return providerName, result
		}
		models, err := provider.LoadModelsFromHTTP(endpoint, key, httpClient, "")
		if err != nil {
			result.Reachable = false
			errMsg := err.Error()
			result.Error = &errMsg
			logger.FromContext(c).WithError(err).Warn("embedding models endpoint check failed")
			return providerName, result
		}
		count := len(models)
		result.Models = &count
		result.Reachable = true
		return providerName, result
	case "ollama":
		if endpoint == "" {
			endpoint = strings.TrimSpace(s.cfg.OllamaServerURL)
		}
		key := strings.TrimSpace(s.cfg.EmbeddingKey)
		if key == "" {
			key = strings.TrimSpace(s.cfg.OllamaServerAPIKey)
		}
		result.URL = endpoint
		if endpoint == "" {
			result.Reachable = false
			result.Error = ptr("embedding endpoint URL is empty")
			return providerName, result
		}
		statusCode, modelsCount, modelNames, err := checkOllamaTagsWithNames(endpoint, key, httpClient)
		if statusCode != nil {
			result.StatusCode = statusCode
		}
		if modelsCount != nil {
			result.Models = modelsCount
		}
		if err != nil {
			result.Reachable = false
			errMsg := err.Error()
			result.Error = &errMsg
			return providerName, result
		}
		if result.Model != nil && !containsModel(*result.Model, modelNames) {
			result.Reachable = false
			errMsg := fmt.Sprintf("embedding model '%s' was not found in Ollama model list", *result.Model)
			result.Error = &errMsg
			return providerName, result
		}
		result.Reachable = true
		return providerName, result
	default:
		result.URL = endpoint
		if endpoint == "" {
			result.Reachable = false
			result.Error = ptr("embedding endpoint URL is empty")
			return providerName, result
		}
		statusCode, err := checkHTTPReachable(endpoint, strings.TrimSpace(s.cfg.EmbeddingKey), httpClient)
		if statusCode != nil {
			result.StatusCode = statusCode
		}
		if err != nil {
			result.Reachable = false
			errMsg := err.Error()
			result.Error = &errMsg
			return providerName, result
		}
		result.Reachable = true
		return providerName, result
	}
}

func checkHTTPReachable(endpoint, apiKey string, client *http.Client) (*int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, strings.TrimRight(endpoint, "/"), nil)
	if err != nil {
		return nil, err
	}
	if apiKey != "" {
		req.Header.Set("Authorization", "Bearer "+apiKey)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	statusCode := resp.StatusCode
	if statusCode >= 200 && statusCode < 500 {
		return &statusCode, nil
	}
	return &statusCode, fmt.Errorf("unexpected status code: %d", statusCode)
}

func checkOllamaTags(endpoint, apiKey string, client *http.Client) (*int, *int, error) {
	statusCode, modelsCount, _, err := checkOllamaTagsWithNames(endpoint, apiKey, client)
	return statusCode, modelsCount, err
}

func checkOllamaTagsWithNames(endpoint, apiKey string, client *http.Client) (*int, *int, []string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	tagsURL := strings.TrimRight(strings.TrimSpace(endpoint), "/") + "/api/tags"
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, tagsURL, nil)
	if err != nil {
		return nil, nil, nil, err
	}
	if apiKey != "" {
		req.Header.Set("Authorization", "Bearer "+apiKey)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, nil, nil, err
	}
	defer resp.Body.Close()

	status := resp.StatusCode
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return &status, nil, nil, err
	}
	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusBadRequest {
		return &status, nil, nil, fmt.Errorf("ollama tags check failed with status %d: %s", resp.StatusCode, strings.TrimSpace(string(body)))
	}

	var listResp struct {
		Models []struct {
			Name string `json:"name"`
		} `json:"models"`
	}
	if err := json.Unmarshal(body, &listResp); err != nil {
		return &status, nil, nil, err
	}

	count := len(listResp.Models)
	names := make([]string, 0, len(listResp.Models))
	for _, model := range listResp.Models {
		names = append(names, model.Name)
	}

	return &status, &count, names, nil
}

func containsModel(target string, names []string) bool {
	target = strings.TrimSpace(strings.ToLower(target))
	for _, name := range names {
		lowerName := strings.TrimSpace(strings.ToLower(name))
		if lowerName == target || strings.HasPrefix(lowerName, target+":") {
			return true
		}
	}

	return false
}

func containsConfiguredModel(target string, models pconfig.ModelsConfig) bool {
	names := make([]string, 0, len(models))
	for _, model := range models {
		names = append(names, model.Name)
	}

	return containsModel(target, names)
}

func ptrToString(value *string) string {
	if value == nil {
		return ""
	}

	return *value
}

func ptr(value string) *string {
	v := value
	return &v
}
