// Package signalhouse provides a Go client for the SignalHouse API.
//
// Usage:
//
//	client, err := signalhouse.NewClient(signalhouse.Config{
//	    APIKey:  "your-api-key",
//	    BaseURL: "https://api.signalhouse.com",
//	})
//	resp, err := client.Auth.Login(ctx, signalhouse.LoginParams{
//	    Email:    "user@example.com",
//	    Password: "secret",
//	})
package signalhouse

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"strings"
)

// Config holds the configuration for the SignalHouse SDK client.
type Config struct {
	// APIKey is the API key for authenticating requests.
	APIKey string
	// BaseURL is the base URL for the SignalHouse API (e.g., "https://api.signalhouse.com").
	BaseURL string
	// EnableAdmin enables admin-only methods on domain services.
	EnableAdmin bool
	// HTTPClient is an optional custom HTTP client. If nil, http.DefaultClient is used.
	HTTPClient *http.Client
}

// Response is the standardized response returned by all SDK methods.
type Response struct {
	Success bool            `json:"success"`
	Data    json.RawMessage `json:"data,omitempty"`
	Status  int             `json:"status"`
	Error   string          `json:"error,omitempty"`
}

// RequestOptions holds per-request overrides.
type RequestOptions struct {
	// Token overrides the default API key for this request.
	Token string
	// Headers are additional headers to include in the request.
	Headers map[string]string
}

// Client is the main SignalHouse SDK client. Access API domains via its exported fields.
type Client struct {
	baseURL     string
	apiKey      string
	enableAdmin bool
	httpClient  *http.Client

	Auth          *AuthService
	Billing       *BillingService
	Brands        *BrandsService
	Campaigns     *CampaignsService
	Groups        *GroupsService
	Landings      *LandingsService
	Messages      *MessagesService
	Notifications *NotificationsService
	Numbers       *NumbersService
	Onboarding    *OnboardingService
	Shortlinks    *ShortlinksService
	Subgroups     *SubgroupsService
	Subscriptions *SubscriptionsService
	Tickets       *TicketsService
	Users         *UsersService
	Voice         *VoiceService
	Webhooks      *WebhooksService
}

// NewClient creates a new SignalHouse SDK client with the provided configuration.
func NewClient(cfg Config) (*Client, error) {
	if cfg.APIKey == "" {
		return nil, fmt.Errorf("API key is required to initialize SignalHouse SDK")
	}
	if cfg.BaseURL == "" {
		return nil, fmt.Errorf("base URL is required to initialize SignalHouse SDK")
	}

	httpClient := cfg.HTTPClient
	if httpClient == nil {
		httpClient = http.DefaultClient
	}

	c := &Client{
		baseURL:     strings.TrimRight(cfg.BaseURL, "/"),
		apiKey:      cfg.APIKey,
		enableAdmin: cfg.EnableAdmin,
		httpClient:  httpClient,
	}

	c.Auth = &AuthService{client: c}
	c.Billing = &BillingService{client: c}
	c.Brands = &BrandsService{client: c}
	c.Campaigns = &CampaignsService{client: c, enableAdmin: cfg.EnableAdmin}
	c.Groups = &GroupsService{client: c, enableAdmin: cfg.EnableAdmin}
	c.Landings = &LandingsService{client: c}
	c.Messages = &MessagesService{client: c}
	c.Notifications = &NotificationsService{client: c}
	c.Numbers = &NumbersService{client: c, enableAdmin: cfg.EnableAdmin}
	c.Onboarding = &OnboardingService{client: c, enableAdmin: cfg.EnableAdmin}
	c.Shortlinks = &ShortlinksService{client: c}
	c.Subgroups = &SubgroupsService{client: c}
	c.Subscriptions = &SubscriptionsService{client: c, enableAdmin: cfg.EnableAdmin}
	c.Tickets = &TicketsService{client: c, enableAdmin: cfg.EnableAdmin}
	c.Users = &UsersService{client: c, enableAdmin: cfg.EnableAdmin}
	c.Voice = &VoiceService{client: c}
	c.Voice.SipTrunks = &VoiceSipTrunksService{client: c}
	c.Voice.SipProfiles = &VoiceSipProfilesService{client: c}
	c.Voice.ProgrammableVoiceProfiles = &VoiceProgrammableVoiceProfilesService{client: c}
	c.Voice.Calls = &VoiceCallsService{client: c}
	c.Voice.CallLogs = &VoiceCallLogsService{client: c}
	c.Voice.Analytics = &VoiceAnalyticsService{client: c}
	c.Voice.GlobalVoiceSettings = &VoiceGlobalVoiceSettingsService{client: c}
	c.Voice.Tokens = &VoiceTokensService{client: c}
	c.Webhooks = &WebhooksService{client: c}

	// Initialize admin sub-services when enabled
	if cfg.EnableAdmin {
		c.Campaigns.Admin = &CampaignsAdminService{client: c}
		c.Groups.Admin = &GroupsAdminService{client: c}
		c.Subscriptions.Admin = &SubscriptionsAdminService{client: c}
		c.Users.Admin = &UsersAdminService{client: c}
		c.Numbers.Admin = &NumbersAdminService{client: c}
		c.Onboarding.Admin = &OnboardingAdminService{client: c}
		c.Tickets.Admin = &TicketsAdminService{client: c}
	}

	return c, nil
}

// require validates that all provided fields are non-empty.
func require(fields map[string]string) error {
	for name, value := range fields {
		if value == "" {
			return &ValidationError{Field: name}
		}
	}
	return nil
}

// requireAny validates that at least one value is non-empty.
func requireAny(fields map[string]string) error {
	for _, value := range fields {
		if value != "" {
			return nil
		}
	}
	names := make([]string, 0, len(fields))
	for name := range fields {
		names = append(names, name)
	}
	return &ValidationError{Field: strings.Join(names, " or ")}
}

// requireNonNil validates that values are non-nil.
func requireNonNil(fields map[string]interface{}) error {
	for name, value := range fields {
		if value == nil {
			return &ValidationError{Field: name}
		}
	}
	return nil
}

// buildQueryString builds a URL query string from key-value pairs, omitting empty values.
func buildQueryString(params map[string]string) string {
	q := url.Values{}
	for k, v := range params {
		if v != "" {
			q.Set(k, v)
		}
	}
	encoded := q.Encode()
	if encoded == "" {
		return ""
	}
	return "?" + encoded
}

// doJSON performs a JSON API request and returns the standardized Response.
func (c *Client) doJSON(ctx context.Context, method, path string, body interface{}, opts *RequestOptions) (*Response, error) {
	var bodyReader io.Reader
	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
		bodyReader = bytes.NewReader(data)
	}

	req, err := http.NewRequestWithContext(ctx, method, c.baseURL+path, bodyReader)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	c.setAuth(req, opts)

	return c.executeRequest(req)
}

// doMultipart performs a multipart form-data request and returns the standardized Response.
func (c *Client) doMultipart(ctx context.Context, method, path string, writeForm func(w *multipart.Writer) error, opts *RequestOptions) (*Response, error) {
	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)

	if err := writeForm(writer); err != nil {
		return nil, fmt.Errorf("failed to write multipart form: %w", err)
	}
	if err := writer.Close(); err != nil {
		return nil, fmt.Errorf("failed to close multipart writer: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, method, c.baseURL+path, &buf)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())
	c.setAuth(req, opts)

	return c.executeRequest(req)
}

// setAuth sets the Authorization header on a request.
func (c *Client) setAuth(req *http.Request, opts *RequestOptions) {
	token := c.apiKey
	if opts != nil {
		if opts.Token != "" {
			token = opts.Token
		}
		for k, v := range opts.Headers {
			req.Header.Set(k, v)
		}
	}
	req.Header.Set("Authorization", "Bearer "+token)
}

// executeRequest sends the request and parses the response into a standardized Response.
func (c *Client) executeRequest(req *http.Request) (*Response, error) {
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		return &Response{
			Success: true,
			Data:    json.RawMessage(respBody),
			Status:  resp.StatusCode,
		}, nil
	}

	// Try to extract error message from response body
	errMsg := string(respBody)
	var errResp struct {
		Message string `json:"message"`
	}
	if json.Unmarshal(respBody, &errResp) == nil && errResp.Message != "" {
		errMsg = errResp.Message
	}

	return &Response{
		Success: false,
		Status:  resp.StatusCode,
		Error:   errMsg,
	}, nil
}
