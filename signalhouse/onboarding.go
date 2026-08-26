package signalhouse

import (
	"context"
	"fmt"
	"net/url"
)

// OnboardingService handles onboarding-related API calls. All methods are
// staff-only and are exposed through the Admin sub-service only.
type OnboardingService struct {
	client      *Client
	enableAdmin bool

	// Admin contains admin-only onboarding methods. Only populated when EnableAdmin is true.
	Admin *OnboardingAdminService
}

// OnboardingAdminService contains admin-only onboarding methods.
type OnboardingAdminService struct {
	client *Client
}

// AdminGetOnboardingsParams are the parameters for the admin GetOnboardings method.
// Reserved for future pagination options; currently empty.
type AdminGetOnboardingsParams struct{}

// GetOnboardings returns all onboarding records (one per group). Admin only.
func (s *OnboardingAdminService) GetOnboardings(ctx context.Context, _ AdminGetOnboardingsParams, opts *RequestOptions) (*Response, error) {
	return s.client.doJSON(ctx, "GET", "/group/onboarding", nil, opts)
}

// AdminGetOnboardingParams are the parameters for the admin GetOnboarding method.
type AdminGetOnboardingParams struct {
	GroupID string
}

// GetOnboarding returns a single onboarding record by group ID. Admin only.
func (s *OnboardingAdminService) GetOnboarding(ctx context.Context, params AdminGetOnboardingParams, opts *RequestOptions) (*Response, error) {
	if err := require(map[string]string{"groupId": params.GroupID}); err != nil {
		return nil, err
	}
	path := fmt.Sprintf("/group/onboarding/%s", url.PathEscape(params.GroupID))
	return s.client.doJSON(ctx, "GET", path, nil, opts)
}
