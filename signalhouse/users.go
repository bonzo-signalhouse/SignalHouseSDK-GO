package signalhouse

import (
	"context"
	"fmt"
	"net/url"
)

// UsersService handles user-related API calls.
type UsersService struct {
	client      *Client
	enableAdmin bool

	// Admin contains admin-only user methods. Only populated when EnableAdmin is true.
	Admin *UsersAdminService
}

// UsersAdminService contains admin-only user methods.
type UsersAdminService struct {
	client *Client
}

// GetUsersParams are the parameters for the GetUsers method.
type GetUsersParams struct {
	ID       string
	GroupID  string
	UserType string
}

// GetUsers returns a list of users with optional filters.
func (s *UsersService) GetUsers(ctx context.Context, params GetUsersParams, opts *RequestOptions) (*Response, error) {
	qs := buildQueryString(map[string]string{
		"id":       params.ID,
		"groupId":  params.GroupID,
		"userType": params.UserType,
	})
	return s.client.doJSON(ctx, "GET", "/user"+qs, nil, opts)
}

// CreateUserData represents the data for creating a user.
type CreateUserData struct {
	GroupID          string `json:"groupId"`
	Role             string `json:"role"`
	Password         string `json:"password"`
	Email            string `json:"email"`
	FirstName        string `json:"firstName,omitempty"`
	LastName         string `json:"lastName,omitempty"`
	CompanyName      string `json:"companyName,omitempty"`
	Phone            string `json:"phone,omitempty"`
	ProfileImageLink string `json:"profileImageLink,omitempty"`
	AddressLine1     string `json:"addressLine1,omitempty"`
	AddressLine2     string `json:"addressLine2,omitempty"`
	City             string `json:"city,omitempty"`
	State            string `json:"state,omitempty"`
	Country          string `json:"country,omitempty"`
	PostalCode       string `json:"postalCode,omitempty"`
	SignupSource     string `json:"signupSource,omitempty"`
}

// CreateUserParams are the parameters for the CreateUser method.
type CreateUserParams struct {
	Data CreateUserData
}

// CreateUser creates a new user.
func (s *UsersService) CreateUser(ctx context.Context, params CreateUserParams, opts *RequestOptions) (*Response, error) {
	return s.client.doJSON(ctx, "POST", "/user", params.Data, opts)
}

// CreateServiceUserData represents the data for creating a service user.
type CreateServiceUserData struct {
	GroupID string `json:"groupId"`
	Name    string `json:"name"`
	Role    string `json:"role"`
}

// CreateServiceUserParams are the parameters for the CreateServiceUser method.
type CreateServiceUserParams struct {
	Data CreateServiceUserData
}

// CreateServiceUser creates a new service user (API key).
func (s *UsersService) CreateServiceUser(ctx context.Context, params CreateServiceUserParams, opts *RequestOptions) (*Response, error) {
	return s.client.doJSON(ctx, "POST", "/user/serviceuser", params.Data, opts)
}

// UpdateUserData represents the data for updating a user.
type UpdateUserData struct {
	Role             string `json:"role,omitempty"`
	Email            string `json:"email,omitempty"`
	FirstName        string `json:"firstName,omitempty"`
	LastName         string `json:"lastName,omitempty"`
	CompanyName      string `json:"companyName,omitempty"`
	Phone            string `json:"phone,omitempty"`
	ProfileImageLink string `json:"profileImageLink,omitempty"`
	AddressLine1     string `json:"addressLine1,omitempty"`
	AddressLine2     string `json:"addressLine2,omitempty"`
	City             string `json:"city,omitempty"`
	State            string `json:"state,omitempty"`
	Country          string `json:"country,omitempty"`
	PostalCode       string `json:"postalCode,omitempty"`
	Status           string `json:"status,omitempty"`
}

// UpdateUserParams are the parameters for the UpdateUser method.
type UpdateUserParams struct {
	ID   string
	Data UpdateUserData
}

// UpdateUser updates a user's information.
func (s *UsersService) UpdateUser(ctx context.Context, params UpdateUserParams, opts *RequestOptions) (*Response, error) {
	if err := require(map[string]string{"id": params.ID}); err != nil {
		return nil, err
	}
	path := fmt.Sprintf("/user/%s", url.PathEscape(params.ID))
	return s.client.doJSON(ctx, "PUT", path, params.Data, opts)
}

// DeleteUserParams are the parameters for the DeleteUser method.
type DeleteUserParams struct {
	ID string
}

// DeleteUser deletes (marks as inactive) a user.
func (s *UsersService) DeleteUser(ctx context.Context, params DeleteUserParams, opts *RequestOptions) (*Response, error) {
	if err := require(map[string]string{"id": params.ID}); err != nil {
		return nil, err
	}
	path := fmt.Sprintf("/user/%s", url.PathEscape(params.ID))
	return s.client.doJSON(ctx, "DELETE", path, nil, opts)
}

// GetNotificationPreferencesParams are the parameters for the GetNotificationPreferences method.
type GetNotificationPreferencesParams struct {
	ID string
}

// GetNotificationPreferences returns notification preferences for a user.
func (s *UsersService) GetNotificationPreferences(ctx context.Context, params GetNotificationPreferencesParams, opts *RequestOptions) (*Response, error) {
	if err := require(map[string]string{"id": params.ID}); err != nil {
		return nil, err
	}
	path := fmt.Sprintf("/user/%s/notification-preferences", url.PathEscape(params.ID))
	return s.client.doJSON(ctx, "GET", path, nil, opts)
}

// NotificationPreference represents a single notification preference entry.
type NotificationPreference struct {
	Name  string `json:"name"`
	Web   bool   `json:"web"`
	Email bool   `json:"email"`
}

// UpdateNotificationPreferencesParams are the parameters for the UpdateNotificationPreferences method.
type UpdateNotificationPreferencesParams struct {
	ID                      string
	NotificationPreferences []NotificationPreference `json:"notificationPreferences"`
}

// UpdateNotificationPreferences updates notification preferences for a user.
func (s *UsersService) UpdateNotificationPreferences(ctx context.Context, params UpdateNotificationPreferencesParams, opts *RequestOptions) (*Response, error) {
	if err := require(map[string]string{"id": params.ID}); err != nil {
		return nil, err
	}
	if len(params.NotificationPreferences) == 0 {
		return nil, &ValidationError{Field: "notificationPreferences"}
	}
	path := fmt.Sprintf("/user/%s/notification-preferences", url.PathEscape(params.ID))
	body := map[string]interface{}{
		"notificationPreferences": params.NotificationPreferences,
	}
	return s.client.doJSON(ctx, "PUT", path, body, opts)
}

// AdminGetUsersParams are the parameters for the admin GetUsers method.
type AdminGetUsersParams struct {
	Email    string
	UserType string
}

// GetUsers returns a list of all users. Admin only.
func (s *UsersAdminService) GetUsers(ctx context.Context, params AdminGetUsersParams, opts *RequestOptions) (*Response, error) {
	qs := buildQueryString(map[string]string{
		"email":    params.Email,
		"userType": params.UserType,
	})
	return s.client.doJSON(ctx, "GET", "/user"+qs, nil, opts)
}

// AdminCreateInternalUserParams are the parameters for the admin CreateInternalUser method.
type AdminCreateInternalUserParams struct {
	Data map[string]interface{}
}

// CreateInternalUser creates a new internal user. Admin only.
func (s *UsersAdminService) CreateInternalUser(ctx context.Context, params AdminCreateInternalUserParams, opts *RequestOptions) (*Response, error) {
	if err := requireNonNil(map[string]interface{}{"data": params.Data}); err != nil {
		return nil, err
	}
	return s.client.doJSON(ctx, "POST", "/user/signalhouse", params.Data, opts)
}
