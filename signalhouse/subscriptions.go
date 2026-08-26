package signalhouse

import (
	"context"
	"fmt"
	"net/url"
)

// SubscriptionsService handles subscription-related API calls.
type SubscriptionsService struct {
	client      *Client
	enableAdmin bool

	// Admin contains admin-only subscription methods. Only populated when EnableAdmin is true.
	Admin *SubscriptionsAdminService
}

// SubscriptionsAdminService contains admin-only subscription methods.
type SubscriptionsAdminService struct {
	client *Client
}

// GetSubscriptionsParams are the parameters for the GetSubscriptions method.
type GetSubscriptionsParams struct {
	GroupID    string
	OnlyActive string
}

// GetSubscriptions returns subscription history for a group.
func (s *SubscriptionsService) GetSubscriptions(ctx context.Context, params GetSubscriptionsParams, opts *RequestOptions) (*Response, error) {
	qs := buildQueryString(map[string]string{
		"groupId":    params.GroupID,
		"onlyActive": params.OnlyActive,
	})
	return s.client.doJSON(ctx, "GET", "/subscription/user"+qs, nil, opts)
}

// SubscribeParams are the parameters for the Subscribe method.
type SubscribeParams struct {
	GroupID    string
	TemplateID string
}

// Subscribe subscribes a group to a template.
func (s *SubscriptionsService) Subscribe(ctx context.Context, params SubscribeParams, opts *RequestOptions) (*Response, error) {
	if err := require(map[string]string{"groupId": params.GroupID, "templateId": params.TemplateID}); err != nil {
		return nil, err
	}
	path := fmt.Sprintf("/subscription/user/subscribe/%s/%s",
		url.PathEscape(params.GroupID), url.PathEscape(params.TemplateID))
	return s.client.doJSON(ctx, "POST", path, map[string]interface{}{}, opts)
}

// UnsubscribeParams are the parameters for the Unsubscribe method.
type UnsubscribeParams struct {
	GroupID string
}

// Unsubscribe cancels a group's subscription.
func (s *SubscriptionsService) Unsubscribe(ctx context.Context, params UnsubscribeParams, opts *RequestOptions) (*Response, error) {
	if err := require(map[string]string{"groupId": params.GroupID}); err != nil {
		return nil, err
	}
	path := fmt.Sprintf("/subscription/user/cancel/%s", url.PathEscape(params.GroupID))
	return s.client.doJSON(ctx, "POST", path, nil, opts)
}

// AdminGetTemplatesParams are the parameters for the admin GetTemplates method.
type AdminGetTemplatesParams struct {
	TemplateID string
}

// GetTemplates returns subscription templates. Admin only.
func (s *SubscriptionsAdminService) GetTemplates(ctx context.Context, params AdminGetTemplatesParams, opts *RequestOptions) (*Response, error) {
	qs := buildQueryString(map[string]string{
		"templateId": params.TemplateID,
	})
	return s.client.doJSON(ctx, "GET", "/subscription/template"+qs, nil, opts)
}

// AdminCreateTemplateParams are the parameters for the admin CreateTemplate method.
type AdminCreateTemplateParams struct {
	TemplateData map[string]interface{}
}

// CreateTemplate creates a new subscription template. Admin only.
func (s *SubscriptionsAdminService) CreateTemplate(ctx context.Context, params AdminCreateTemplateParams, opts *RequestOptions) (*Response, error) {
	if err := requireNonNil(map[string]interface{}{"templateData": params.TemplateData}); err != nil {
		return nil, err
	}
	return s.client.doJSON(ctx, "POST", "/subscription/template", params.TemplateData, opts)
}

// AdminUpdateTemplateParams are the parameters for the admin UpdateTemplate method.
type AdminUpdateTemplateParams struct {
	TemplateID   string
	TemplateData map[string]interface{}
}

// UpdateTemplate updates a subscription template. Admin only.
func (s *SubscriptionsAdminService) UpdateTemplate(ctx context.Context, params AdminUpdateTemplateParams, opts *RequestOptions) (*Response, error) {
	if err := require(map[string]string{"templateId": params.TemplateID}); err != nil {
		return nil, err
	}
	if err := requireNonNil(map[string]interface{}{"templateData": params.TemplateData}); err != nil {
		return nil, err
	}
	path := fmt.Sprintf("/subscription/template/%s", url.PathEscape(params.TemplateID))
	return s.client.doJSON(ctx, "PUT", path, params.TemplateData, opts)
}

// AdminDeleteTemplateParams are the parameters for the admin DeleteTemplate method.
type AdminDeleteTemplateParams struct {
	TemplateID string
}

// DeleteTemplate deletes a subscription template. Admin only.
func (s *SubscriptionsAdminService) DeleteTemplate(ctx context.Context, params AdminDeleteTemplateParams, opts *RequestOptions) (*Response, error) {
	if err := require(map[string]string{"templateId": params.TemplateID}); err != nil {
		return nil, err
	}
	path := fmt.Sprintf("/subscription/template/%s", url.PathEscape(params.TemplateID))
	return s.client.doJSON(ctx, "DELETE", path, nil, opts)
}

// AdminCreateCustomSubscriptionParams are the parameters for the admin CreateCustomSubscription method.
type AdminCreateCustomSubscriptionParams struct {
	GroupID          string
	SubscriptionData map[string]interface{}
}

// CreateCustomSubscription creates a custom subscription for a group. Admin only.
func (s *SubscriptionsAdminService) CreateCustomSubscription(ctx context.Context, params AdminCreateCustomSubscriptionParams, opts *RequestOptions) (*Response, error) {
	if err := require(map[string]string{"groupId": params.GroupID}); err != nil {
		return nil, err
	}
	if err := requireNonNil(map[string]interface{}{"subscriptionData": params.SubscriptionData}); err != nil {
		return nil, err
	}
	path := fmt.Sprintf("/subscription/user/custom/%s", url.PathEscape(params.GroupID))
	return s.client.doJSON(ctx, "POST", path, params.SubscriptionData, opts)
}
