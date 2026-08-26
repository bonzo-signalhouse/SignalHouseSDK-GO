package signalhouse

import (
	"context"
	"fmt"
	"net/url"
)

// AuthService handles authentication-related API calls.
type AuthService struct {
	client *Client
}

// LoginParams are the parameters for the Login method.
type LoginParams struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Login authenticates a user with email and password.
func (s *AuthService) Login(ctx context.Context, params LoginParams, opts *RequestOptions) (*Response, error) {
	return s.client.doJSON(ctx, "POST", "/auth", params, opts)
}

// ResetPasswordParams are the parameters for the ResetPassword method.
type ResetPasswordParams struct {
	UserID      string `json:"-"`
	NewPassword string `json:"newPassword"`
}

// ResetPassword resets a user's password.
func (s *AuthService) ResetPassword(ctx context.Context, params ResetPasswordParams, opts *RequestOptions) (*Response, error) {
	if err := require(map[string]string{"userId": params.UserID, "newPassword": params.NewPassword}); err != nil {
		return nil, err
	}
	path := fmt.Sprintf("/auth/resetpassword/%s", url.PathEscape(params.UserID))
	body := map[string]string{"newPassword": params.NewPassword}
	return s.client.doJSON(ctx, "PUT", path, body, opts)
}

// ForgotPasswordParams are the parameters for the ForgotPassword method.
type ForgotPasswordParams struct {
	Email string `json:"email"`
}

// ForgotPassword requests a password-reset link be emailed to the given address (public).
// Always succeeds regardless of whether an account exists for the email, so it cannot be
// used to enumerate registered emails.
func (s *AuthService) ForgotPassword(ctx context.Context, params ForgotPasswordParams, opts *RequestOptions) (*Response, error) {
	if err := require(map[string]string{"email": params.Email}); err != nil {
		return nil, err
	}
	return s.client.doJSON(ctx, "POST", "/auth/forgot-password", params, opts)
}

// ResetPasswordWithTokenParams are the parameters for the ResetPasswordWithToken method.
type ResetPasswordWithTokenParams struct {
	Token    string `json:"token"`
	Password string `json:"password"`
}

// ResetPasswordWithToken resets a password using the single-use token from a reset email (public).
func (s *AuthService) ResetPasswordWithToken(ctx context.Context, params ResetPasswordWithTokenParams, opts *RequestOptions) (*Response, error) {
	if err := require(map[string]string{"token": params.Token, "password": params.Password}); err != nil {
		return nil, err
	}
	return s.client.doJSON(ctx, "POST", "/auth/reset-password", params, opts)
}

// GetAuthHistoryParams are the parameters for the GetAuthHistory method.
type GetAuthHistoryParams struct {
	GroupID string
	UserID  string
	Page    string
	Limit   string
}

// GetAuthHistory returns login history for a group or user.
func (s *AuthService) GetAuthHistory(ctx context.Context, params GetAuthHistoryParams, opts *RequestOptions) (*Response, error) {
	if err := requireAny(map[string]string{"groupId": params.GroupID, "userId": params.UserID}); err != nil {
		return nil, err
	}
	qs := buildQueryString(map[string]string{
		"groupId": params.GroupID,
		"userId":  params.UserID,
		"page":    params.Page,
		"limit":   params.Limit,
	})
	return s.client.doJSON(ctx, "GET", "/auth/history"+qs, nil, opts)
}

// LogoutAll logs out all other users in the caller's active group.
func (s *AuthService) LogoutAll(ctx context.Context, opts *RequestOptions) (*Response, error) {
	return s.client.doJSON(ctx, "POST", "/auth/logout-all", nil, opts)
}

// GetGroupID returns the Group ID associated with the caller's JWT (their active group).
func (s *AuthService) GetGroupID(ctx context.Context, opts *RequestOptions) (*Response, error) {
	return s.client.doJSON(ctx, "GET", "/auth/group-id", nil, opts)
}

// RequestExternalLinkTokenParams are the parameters for the RequestExternalLinkToken method.
type RequestExternalLinkTokenParams struct {
	Product string `json:"product"`
}

// RequestExternalLinkToken mints a single-use, short-lived external-link token for the
// authenticated caller. The token is handed to the GHL/Shopify backend so it can link the
// caller's existing V2 group to its tenant.
func (s *AuthService) RequestExternalLinkToken(ctx context.Context, params RequestExternalLinkTokenParams, opts *RequestOptions) (*Response, error) {
	if err := require(map[string]string{"product": params.Product}); err != nil {
		return nil, err
	}
	return s.client.doJSON(ctx, "POST", "/auth/external-link-token", params, opts)
}
