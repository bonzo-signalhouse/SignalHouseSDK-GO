package signalhouse

import (
	"context"
	"net/url"
)

// VoiceProgrammableVoiceProfilesService handles Programmable Voice Profile API calls.
//
// A Programmable Voice Profile groups a set of the account's numbers (across one
// or more subgroups) under a single inbound call-handling decision routed through
// Signal House's own voice servers — the same idea as a SIP trunk, but without a
// customer SBC. A number can belong to at most one profile and is mutually
// exclusive with SIP trunks.
type VoiceProgrammableVoiceProfilesService struct {
	client *Client
}

// CreateProgrammableVoiceProfileParams is the body for CreateProgrammableVoiceProfile.
// Name and SubgroupIDs are required.
type CreateProgrammableVoiceProfileParams struct {
	Name                string   `json:"name"`
	SubgroupIDs         []string `json:"subgroupIds"`
	Region              string   `json:"region,omitempty"`
	RouteAction         string   `json:"routeAction,omitempty"` // "FORWARD" | "WEBRTC" | "SIP_TRUNK" | "SIP_PROFILE" (default "FORWARD"). Targets: FORWARD → ForwardToE164, SIP_TRUNK → RouteSipTrunkID, SIP_PROFILE → RouteSipProfileID.
	ForwardToE164       string   `json:"forwardToE164,omitempty"`
	ForwardAfterSeconds *int     `json:"forwardAfterSeconds,omitempty"`
	RouteSipTrunkID     string   `json:"routeSipTrunkId,omitempty"`   // SIP trunk id, required when RouteAction is "SIP_TRUNK".
	RouteSipProfileID   string   `json:"routeSipProfileId,omitempty"` // SIP endpoint id, required when RouteAction is "SIP_PROFILE".
	RecordingEnabled    *bool    `json:"recordingEnabled,omitempty"`
	Enabled             *bool    `json:"enabled,omitempty"`
}

// UpdateProgrammableVoiceProfileParams is the body for UpdateProgrammableVoiceProfile.
// All fields optional (partial update).
type UpdateProgrammableVoiceProfileParams struct {
	Name                string   `json:"name,omitempty"`
	SubgroupIDs         []string `json:"subgroupIds,omitempty"`
	Region              string   `json:"region,omitempty"`
	RouteAction         string   `json:"routeAction,omitempty"` // "FORWARD" | "WEBRTC" | "SIP_TRUNK" | "SIP_PROFILE". Targets: FORWARD → ForwardToE164, SIP_TRUNK → RouteSipTrunkID, SIP_PROFILE → RouteSipProfileID.
	ForwardToE164       string   `json:"forwardToE164,omitempty"`
	ForwardAfterSeconds *int     `json:"forwardAfterSeconds,omitempty"`
	RouteSipTrunkID     string   `json:"routeSipTrunkId,omitempty"`   // SIP trunk id, required when RouteAction is "SIP_TRUNK".
	RouteSipProfileID   string   `json:"routeSipProfileId,omitempty"` // SIP endpoint id, required when RouteAction is "SIP_PROFILE".
	RecordingEnabled    *bool    `json:"recordingEnabled,omitempty"`
	Enabled             *bool    `json:"enabled,omitempty"`
}

// ListProgrammableVoiceProfiles returns all Programmable Voice Profiles for the
// current account. Returns { profiles: [...] }.
// GET /voice/api/v1/programmable-voice-profiles
func (s *VoiceProgrammableVoiceProfilesService) ListProgrammableVoiceProfiles(ctx context.Context, opts *RequestOptions) (*Response, error) {
	return s.client.doJSON(ctx, "GET", "/voice/api/v1/programmable-voice-profiles", nil, opts)
}

// GetProgrammableVoiceProfile returns a single Programmable Voice Profile by ID.
// Returns { profile }.
// GET /voice/api/v1/programmable-voice-profiles/:id
func (s *VoiceProgrammableVoiceProfilesService) GetProgrammableVoiceProfile(ctx context.Context, id string, opts *RequestOptions) (*Response, error) {
	if err := require(map[string]string{"id": id}); err != nil {
		return nil, err
	}
	return s.client.doJSON(ctx, "GET", "/voice/api/v1/programmable-voice-profiles/"+url.PathEscape(id), nil, opts)
}

// CreateProgrammableVoiceProfile creates a new Programmable Voice Profile.
// Returns { profile }.
// POST /voice/api/v1/programmable-voice-profiles
func (s *VoiceProgrammableVoiceProfilesService) CreateProgrammableVoiceProfile(ctx context.Context, params CreateProgrammableVoiceProfileParams, opts *RequestOptions) (*Response, error) {
	if err := require(map[string]string{"name": params.Name}); err != nil {
		return nil, err
	}
	return s.client.doJSON(ctx, "POST", "/voice/api/v1/programmable-voice-profiles", params, opts)
}

// UpdateProgrammableVoiceProfile updates an existing Programmable Voice Profile
// (partial). Returns { profile }.
// PATCH /voice/api/v1/programmable-voice-profiles/:id
func (s *VoiceProgrammableVoiceProfilesService) UpdateProgrammableVoiceProfile(ctx context.Context, id string, params UpdateProgrammableVoiceProfileParams, opts *RequestOptions) (*Response, error) {
	if err := require(map[string]string{"id": id}); err != nil {
		return nil, err
	}
	return s.client.doJSON(ctx, "PATCH", "/voice/api/v1/programmable-voice-profiles/"+url.PathEscape(id), params, opts)
}

// ToggleProgrammableVoiceProfileActive toggles a profile active/inactive.
// Returns { profile }.
// POST /voice/api/v1/programmable-voice-profiles/:id/toggle-active
func (s *VoiceProgrammableVoiceProfilesService) ToggleProgrammableVoiceProfileActive(ctx context.Context, id string, opts *RequestOptions) (*Response, error) {
	if err := require(map[string]string{"id": id}); err != nil {
		return nil, err
	}
	return s.client.doJSON(ctx, "POST", "/voice/api/v1/programmable-voice-profiles/"+url.PathEscape(id)+"/toggle-active", nil, opts)
}

// DeleteProgrammableVoiceProfile deletes a Programmable Voice Profile. Its number
// assignments cascade away, so those numbers fall back to subgroup/global routing.
// Returns { message }.
// DELETE /voice/api/v1/programmable-voice-profiles/:id
func (s *VoiceProgrammableVoiceProfilesService) DeleteProgrammableVoiceProfile(ctx context.Context, id string, opts *RequestOptions) (*Response, error) {
	if err := require(map[string]string{"id": id}); err != nil {
		return nil, err
	}
	return s.client.doJSON(ctx, "DELETE", "/voice/api/v1/programmable-voice-profiles/"+url.PathEscape(id), nil, opts)
}

type programmableVoiceProfileNumberBody struct {
	E164 string `json:"e164"`
}

// AssignNumberToProgrammableVoiceProfile assigns a phone number (by E.164) to
// this profile. Rejected if the number is already configured on a SIP
// trunk/endpoint or another profile — it must be de-configured there first
// (mutual exclusivity). Returns { profile }.
// POST /voice/api/v1/programmable-voice-profiles/:id/assign-number
func (s *VoiceProgrammableVoiceProfilesService) AssignNumberToProgrammableVoiceProfile(ctx context.Context, id, e164 string, opts *RequestOptions) (*Response, error) {
	if err := require(map[string]string{"id": id, "e164": e164}); err != nil {
		return nil, err
	}
	body := programmableVoiceProfileNumberBody{E164: e164}
	return s.client.doJSON(ctx, "POST", "/voice/api/v1/programmable-voice-profiles/"+url.PathEscape(id)+"/assign-number", body, opts)
}

// UnassignNumberFromProgrammableVoiceProfile unassigns a phone number (by E.164)
// from this profile. Returns { profile }.
// POST /voice/api/v1/programmable-voice-profiles/:id/unassign-number
func (s *VoiceProgrammableVoiceProfilesService) UnassignNumberFromProgrammableVoiceProfile(ctx context.Context, id, e164 string, opts *RequestOptions) (*Response, error) {
	if err := require(map[string]string{"id": id, "e164": e164}); err != nil {
		return nil, err
	}
	body := programmableVoiceProfileNumberBody{E164: e164}
	return s.client.doJSON(ctx, "POST", "/voice/api/v1/programmable-voice-profiles/"+url.PathEscape(id)+"/unassign-number", body, opts)
}
