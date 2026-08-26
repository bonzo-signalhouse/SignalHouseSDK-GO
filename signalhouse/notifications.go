package signalhouse

import (
	"context"
	"fmt"
	"net/url"
)

// NotificationsService handles notification-related API calls.
type NotificationsService struct {
	client *Client
}

// GetNotificationsParams are the parameters for the GetNotifications method.
type GetNotificationsParams struct {
	ID         string
	GroupID    string
	Page       string
	Limit      string
	Status     string
	EventTypes string
}

// GetNotifications returns notifications by ID or groupID with optional filters.
func (s *NotificationsService) GetNotifications(ctx context.Context, params GetNotificationsParams, opts *RequestOptions) (*Response, error) {
	if err := requireAny(map[string]string{"id": params.ID, "groupId": params.GroupID}); err != nil {
		return nil, err
	}
	qs := buildQueryString(map[string]string{
		"id":         params.ID,
		"groupId":    params.GroupID,
		"page":       params.Page,
		"limit":      params.Limit,
		"status":     params.Status,
		"eventTypes": params.EventTypes,
	})
	return s.client.doJSON(ctx, "GET", "/notification"+qs, nil, opts)
}

// UpdateNotificationStatusParams are the parameters for the UpdateNotificationStatus method.
type UpdateNotificationStatusParams struct {
	IDs    interface{} `json:"ids"`
	Status string      `json:"status"`
}

// UpdateNotificationStatus updates the status of one or more notifications.
func (s *NotificationsService) UpdateNotificationStatus(ctx context.Context, params UpdateNotificationStatusParams, opts *RequestOptions) (*Response, error) {
	if params.Status == "" {
		return nil, &ValidationError{Field: "status"}
	}
	if params.IDs == nil {
		return nil, &ValidationError{Field: "ids"}
	}
	return s.client.doJSON(ctx, "PUT", "/notification/status", params, opts)
}

// DeleteNotificationParams are the parameters for the DeleteNotification method.
type DeleteNotificationParams struct {
	ID string
}

// DeleteNotification deletes a notification by ID.
func (s *NotificationsService) DeleteNotification(ctx context.Context, params DeleteNotificationParams, opts *RequestOptions) (*Response, error) {
	if err := require(map[string]string{"id": params.ID}); err != nil {
		return nil, err
	}
	path := fmt.Sprintf("/notification/%s", url.PathEscape(params.ID))
	return s.client.doJSON(ctx, "DELETE", path, nil, opts)
}
