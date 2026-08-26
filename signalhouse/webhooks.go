package signalhouse

import (
	"context"
	"fmt"
	"net/url"
)

// WebhooksService handles webhook-related API calls.
type WebhooksService struct {
	client *Client
}

// GetWebhooksParams are the parameters for the GetWebhooks method.
type GetWebhooksParams struct {
	ID           string
	GroupID      string
	EndpointType string
	PhoneNumber  string
	Page         string
	Limit        string
}

// GetWebhooks returns a list of webhooks with optional filters.
func (s *WebhooksService) GetWebhooks(ctx context.Context, params GetWebhooksParams, opts *RequestOptions) (*Response, error) {
	qs := buildQueryString(map[string]string{
		"id":           params.ID,
		"groupId":      params.GroupID,
		"endpointType": params.EndpointType,
		"phoneNumber":  params.PhoneNumber,
		"page":         params.Page,
		"limit":        params.Limit,
	})
	return s.client.doJSON(ctx, "GET", "/webhook"+qs, nil, opts)
}

// CreateWebhookParams are the parameters for the CreateWebhook method.
type CreateWebhookParams struct {
	WebhookData map[string]interface{}
}

// CreateWebhook creates a new webhook.
func (s *WebhooksService) CreateWebhook(ctx context.Context, params CreateWebhookParams, opts *RequestOptions) (*Response, error) {
	if err := requireNonNil(map[string]interface{}{"webhookData": params.WebhookData}); err != nil {
		return nil, err
	}
	return s.client.doJSON(ctx, "POST", "/webhook", params.WebhookData, opts)
}

// UpdateWebhookParams are the parameters for the UpdateWebhook method.
type UpdateWebhookParams struct {
	ID         string
	UpdateData map[string]interface{}
}

// UpdateWebhook updates an existing webhook.
func (s *WebhooksService) UpdateWebhook(ctx context.Context, params UpdateWebhookParams, opts *RequestOptions) (*Response, error) {
	if err := require(map[string]string{"id": params.ID}); err != nil {
		return nil, err
	}
	if err := requireNonNil(map[string]interface{}{"updateData": params.UpdateData}); err != nil {
		return nil, err
	}
	path := fmt.Sprintf("/webhook/%s", url.PathEscape(params.ID))
	return s.client.doJSON(ctx, "PUT", path, params.UpdateData, opts)
}

// DeleteWebhookParams are the parameters for the DeleteWebhook method.
type DeleteWebhookParams struct {
	ID string
}

// DeleteWebhook deletes (marks as inactive) a webhook.
func (s *WebhooksService) DeleteWebhook(ctx context.Context, params DeleteWebhookParams, opts *RequestOptions) (*Response, error) {
	if err := require(map[string]string{"id": params.ID}); err != nil {
		return nil, err
	}
	path := fmt.Sprintf("/webhook/%s", url.PathEscape(params.ID))
	return s.client.doJSON(ctx, "DELETE", path, nil, opts)
}
