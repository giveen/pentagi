package tools

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
	"sync"
	"time"

	obs "pentagi/pkg/observability"
	"pentagi/pkg/observability/langfuse"

	"github.com/sirupsen/logrus"
	"golang.org/x/net/publicsuffix"
)

const (
	httpClientDefaultTimeout = 30 * time.Second
	httpClientMaxTimeout     = 300 * time.Second
	httpClientDefaultMaxBody = 50 * 1024  // 50 KB
	httpClientAbsMaxBody     = 512 * 1024 // 512 KB hard ceiling
	httpClientMinMaxBody     = 512        // 512 B minimum
)

// httpClientTool performs structured HTTP requests with per-subtask session cookie jars.
type httpClientTool struct {
	flowID    int64
	taskID    *int64
	subtaskID *int64

	mu   sync.Mutex
	jars map[string]*cookiejar.Jar // keyed by subtask session key
}

// NewHttpClientTool creates a new httpclient tool instance.
func NewHttpClientTool(
	flowID int64, taskID, subtaskID *int64,
) Tool {
	return &httpClientTool{
		flowID:    flowID,
		taskID:    taskID,
		subtaskID: subtaskID,
		jars:      make(map[string]*cookiejar.Jar),
	}
}

// IsAvailable always returns true — no external dependencies required.
func (h *httpClientTool) IsAvailable() bool {
	return true
}

// sessionKey builds a stable string key from task/subtask IDs.
func (h *httpClientTool) sessionKey() string {
	var tid, sid int64
	if h.taskID != nil {
		tid = *h.taskID
	}
	if h.subtaskID != nil {
		sid = *h.subtaskID
	}
	return fmt.Sprintf("%d:%d", tid, sid)
}

// getOrCreateJar returns (or lazily creates) the cookie jar for the current subtask.
func (h *httpClientTool) getOrCreateJar() *cookiejar.Jar {
	key := h.sessionKey()
	h.mu.Lock()
	defer h.mu.Unlock()
	if jar, ok := h.jars[key]; ok {
		return jar
	}
	jar, _ := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
	h.jars[key] = jar
	return jar
}

// Handle processes an HTTP request tool call from an AI agent.
func (h *httpClientTool) Handle(ctx context.Context, name string, args json.RawMessage) (string, error) {
	var action HttpRequest
	ctx, observation := obs.Observer.NewObservation(ctx)
	logger := logrus.WithContext(ctx).WithFields(enrichLogrusFields(h.flowID, h.taskID, h.subtaskID, logrus.Fields{
		"tool": name,
		"args": string(args),
	}))

	if err := json.Unmarshal(args, &action); err != nil {
		logger.WithError(err).Error("failed to unmarshal httpclient action")
		return "", fmt.Errorf("failed to unmarshal %s action arguments: %w", name, err)
	}

	result, err := h.doRequest(ctx, &action)
	if err != nil {
		observation.Event(
			langfuse.WithEventName("httpclient request error swallowed"),
			langfuse.WithEventInput(map[string]any{
				"method": action.Method,
				"url":    action.URL,
			}),
			langfuse.WithEventStatus(err.Error()),
			langfuse.WithEventLevel(langfuse.ObservationLevelWarning),
			langfuse.WithEventMetadata(langfuse.Metadata{
				"tool_name": HttpClientToolName,
				"method":    action.Method,
				"url":       action.URL,
				"error":     err.Error(),
			}),
		)
		logger.WithError(err).Error("httpclient request failed")
		return fmt.Sprintf("httpclient request failed: %v", err), nil
	}

	return result, nil
}

// doRequest builds and executes the HTTP request, returning a structured text result.
func (h *httpClientTool) doRequest(ctx context.Context, action *HttpRequest) (string, error) {
	// --- validate / normalise inputs ---

	method := strings.ToUpper(strings.TrimSpace(action.Method))
	if method == "" {
		method = http.MethodGet
	}

	rawURL := strings.TrimSpace(action.URL)
	if rawURL == "" {
		return "", fmt.Errorf("url is required")
	}
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return "", fmt.Errorf("invalid url %q: %w", rawURL, err)
	}
	if parsedURL.Scheme != "http" && parsedURL.Scheme != "https" {
		return "", fmt.Errorf("unsupported url scheme %q (only http/https allowed)", parsedURL.Scheme)
	}

	// --- timeout ---
	timeoutSecs := int64(action.TimeoutSecs)
	if timeoutSecs <= 0 {
		timeoutSecs = int64(httpClientDefaultTimeout.Seconds())
	}
	if timeoutSecs > int64(httpClientMaxTimeout.Seconds()) {
		timeoutSecs = int64(httpClientMaxTimeout.Seconds())
	}
	timeout := time.Duration(timeoutSecs) * time.Second

	// --- max body ---
	maxBody := int64(action.MaxBodyBytes)
	if maxBody < httpClientMinMaxBody {
		maxBody = httpClientDefaultMaxBody
	}
	if maxBody > httpClientAbsMaxBody {
		maxBody = httpClientAbsMaxBody
	}

	// --- build body ---
	var bodyReader io.Reader
	var contentTypeOverride string

	switch {
	case strings.TrimSpace(action.BodyJSON) != "":
		bodyReader = strings.NewReader(action.BodyJSON)
		contentTypeOverride = "application/json"
	case len(action.BodyForm) > 0:
		form := url.Values{}
		for k, v := range action.BodyForm {
			form.Set(k, v)
		}
		bodyReader = strings.NewReader(form.Encode())
		contentTypeOverride = "application/x-www-form-urlencoded"
	case action.BodyRaw != "":
		bodyReader = strings.NewReader(action.BodyRaw)
	}

	req, err := http.NewRequestWithContext(ctx, method, rawURL, bodyReader)
	if err != nil {
		return "", fmt.Errorf("failed to build request: %w", err)
	}

	// --- headers ---
	for k, v := range action.Headers {
		req.Header.Set(k, v)
	}
	if contentTypeOverride != "" && req.Header.Get("Content-Type") == "" {
		req.Header.Set("Content-Type", contentTypeOverride)
	}

	// --- TLS config ---
	skipVerify := !bool(action.VerifyTLS)
	tlsCfg := &tls.Config{InsecureSkipVerify: skipVerify} // #nosec G402 — intentional for pentest context

	// --- transport / proxy ---
	transport := &http.Transport{
		TLSClientConfig: tlsCfg,
	}
	if action.ProxyURL != "" {
		proxyParsed, perr := url.Parse(action.ProxyURL)
		if perr != nil {
			return "", fmt.Errorf("invalid proxy_url %q: %w", action.ProxyURL, perr)
		}
		transport.Proxy = http.ProxyURL(proxyParsed)
	}

	// --- cookie jar ---
	jar := h.getOrCreateJar()

	client := &http.Client{
		Timeout:   timeout,
		Transport: transport,
	}

	if bool(action.SendCookies) {
		client.Jar = jar
	}
	if !bool(action.FollowRedirects) {
		client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		}
	}

	// --- execute ---
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("request error: %w", err)
	}
	defer resp.Body.Close()

	// --- save cookies if requested ---
	if bool(action.SaveCookies) {
		for _, c := range resp.Cookies() {
			jar.SetCookies(parsedURL, []*http.Cookie{c})
		}
	}

	// --- read body ---
	limitedBody := io.LimitReader(resp.Body, maxBody+1)
	bodyBytes, err := io.ReadAll(limitedBody)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %w", err)
	}
	truncated := false
	if int64(len(bodyBytes)) > maxBody {
		bodyBytes = bodyBytes[:maxBody]
		truncated = true
	}

	// --- format result ---
	var sb strings.Builder
	fmt.Fprintf(&sb, "Status: %d %s\n", resp.StatusCode, http.StatusText(resp.StatusCode))
	fmt.Fprintf(&sb, "URL: %s\n", resp.Request.URL.String())

	if len(resp.Header) > 0 {
		sb.WriteString("\nResponse Headers:\n")
		for k, vv := range resp.Header {
			fmt.Fprintf(&sb, "  %s: %s\n", k, strings.Join(vv, ", "))
		}
	}

	if bool(action.SaveCookies) || bool(action.SendCookies) {
		cookies := jar.Cookies(parsedURL)
		if len(cookies) > 0 {
			sb.WriteString("\nSession Cookies:\n")
			for _, c := range cookies {
				fmt.Fprintf(&sb, "  %s=%s\n", c.Name, c.Value)
			}
		}
	}

	if len(bodyBytes) > 0 {
		sb.WriteString("\nResponse Body:\n")
		sb.Write(bodyBytes)
	}

	if truncated {
		fmt.Fprintf(&sb, "\n\n[Response body truncated at %d bytes]", maxBody)
	}

	return sb.String(), nil
}
