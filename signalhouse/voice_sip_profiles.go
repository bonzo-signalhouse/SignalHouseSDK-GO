package signalhouse

import (
	"context"
	"net/url"
)

// VoiceSipProfilesService handles SIP profile (endpoint) API calls.
//
// A SIP profile represents a single registerable UA — a desk phone, a
// softphone, a SIP-capable app — that registers to Signal House with a
// username + password. Distinct from a SIP trunk, which is a peer-to-peer
// link to a PBX or carrier.
type VoiceSipProfilesService struct {
	client *Client
}

// CreateSipProfileParams is the body for CreateSipProfile.
type CreateSipProfileParams struct {
	Name                 string `json:"name"`
	RecordingAllowed     *bool  `json:"recordingAllowed,omitempty"`
	TranscriptionEnabled *bool  `json:"transcriptionEnabled,omitempty"`
	SentimentFlag        *bool  `json:"sentimentFlag,omitempty"`
}

// UpdateSipProfileParams is the body for UpdateSipProfile. All fields optional.
// To rotate the password, pass a new Password value — there is no dedicated
// regenerate endpoint for profiles.
type UpdateSipProfileParams struct {
	Name                 string `json:"name,omitempty"`
	SIPUsername          string `json:"sipUsername,omitempty"`
	Password             string `json:"password,omitempty"`
	Transport            string `json:"transport,omitempty"` // "UDP" | "TCP" | "TLS"
	RecordingAllowed     *bool  `json:"recordingAllowed,omitempty"`
	TranscriptionEnabled *bool  `json:"transcriptionEnabled,omitempty"`
	SentimentFlag        *bool  `json:"sentimentFlag,omitempty"`
}

// ListSipProfiles returns all SIP profiles for the current account.
// GET /voice/sip-profiles
func (s *VoiceSipProfilesService) ListSipProfiles(ctx context.Context, opts *RequestOptions) (*Response, error) {
	return s.client.doJSON(ctx, "GET", "/voice/sip-profiles", nil, opts)
}

// GetSipProfile returns a single SIP profile by ID. The response does NOT
// include the password — call GetSipProfilePassword to retrieve it.
// GET /voice/sip-profiles/:id
func (s *VoiceSipProfilesService) GetSipProfile(ctx context.Context, id string, opts *RequestOptions) (*Response, error) {
	if err := require(map[string]string{"id": id}); err != nil {
		return nil, err
	}
	return s.client.doJSON(ctx, "GET", "/voice/sip-profiles/"+url.PathEscape(id), nil, opts)
}

// GetSipProfilePassword fetches the SIP password for this profile. Returns
// { password: "..." }.
// GET /voice/sip-profiles/:id/password
func (s *VoiceSipProfilesService) GetSipProfilePassword(ctx context.Context, id string, opts *RequestOptions) (*Response, error) {
	if err := require(map[string]string{"id": id}); err != nil {
		return nil, err
	}
	return s.client.doJSON(ctx, "GET", "/voice/sip-profiles/"+url.PathEscape(id)+"/password", nil, opts)
}

// GetSipProfileTransports lists valid SIP transports (UDP/TCP/TLS) with
// their address and port for the active region.
// GET /voice/sip-profiles/transports
func (s *VoiceSipProfilesService) GetSipProfileTransports(ctx context.Context, opts *RequestOptions) (*Response, error) {
	return s.client.doJSON(ctx, "GET", "/voice/sip-profiles/transports", nil, opts)
}

// CreateSipProfile creates a new SIP profile. The server generates the SIP
// username and password; the password is returned at the TOP LEVEL of the
// create response (not nested inside sipProfile) and is the only chance to
// surface it without an explicit GetSipProfilePassword call.
// POST /voice/sip-profiles
func (s *VoiceSipProfilesService) CreateSipProfile(ctx context.Context, params CreateSipProfileParams, opts *RequestOptions) (*Response, error) {
	if err := require(map[string]string{"name": params.Name}); err != nil {
		return nil, err
	}
	return s.client.doJSON(ctx, "POST", "/voice/sip-profiles", params, opts)
}

// UpdateSipProfile updates an existing SIP profile.
// PATCH /voice/sip-profiles/:id
func (s *VoiceSipProfilesService) UpdateSipProfile(ctx context.Context, id string, params UpdateSipProfileParams, opts *RequestOptions) (*Response, error) {
	if err := require(map[string]string{"id": id}); err != nil {
		return nil, err
	}
	return s.client.doJSON(ctx, "PATCH", "/voice/sip-profiles/"+url.PathEscape(id), params, opts)
}

// DeleteSipProfile deletes a SIP profile. Unassigns any linked numbers as a
// side effect.
// DELETE /voice/sip-profiles/:id
func (s *VoiceSipProfilesService) DeleteSipProfile(ctx context.Context, id string, opts *RequestOptions) (*Response, error) {
	if err := require(map[string]string{"id": id}); err != nil {
		return nil, err
	}
	return s.client.doJSON(ctx, "DELETE", "/voice/sip-profiles/"+url.PathEscape(id), nil, opts)
}

type sipProfileAssignNumberBody struct {
	PhoneNumberID string `json:"phoneNumberId"`
}

// AssignNumberToSipProfile assigns a phone number to this SIP profile (routes
// inbound calls on that number to the endpoint). phoneNumberID may be the
// number's E.164 (preferred; the backend detects it) or a legacy voice-backend
// phone-number UUID.
// POST /voice/sip-profiles/:id/assign-number
func (s *VoiceSipProfilesService) AssignNumberToSipProfile(ctx context.Context, id, phoneNumberID string, opts *RequestOptions) (*Response, error) {
	if err := require(map[string]string{"id": id, "phoneNumberId": phoneNumberID}); err != nil {
		return nil, err
	}
	body := sipProfileAssignNumberBody{PhoneNumberID: phoneNumberID}
	return s.client.doJSON(ctx, "POST", "/voice/sip-profiles/"+url.PathEscape(id)+"/assign-number", body, opts)
}

// UnassignNumberFromSipProfile removes a phone number from this SIP profile.
// phoneNumberID may be the number's E.164 (preferred; the backend detects it)
// or a legacy voice-backend phone-number UUID.
// POST /voice/sip-profiles/:id/unassign-number
func (s *VoiceSipProfilesService) UnassignNumberFromSipProfile(ctx context.Context, id, phoneNumberID string, opts *RequestOptions) (*Response, error) {
	if err := require(map[string]string{"id": id, "phoneNumberId": phoneNumberID}); err != nil {
		return nil, err
	}
	body := sipProfileAssignNumberBody{PhoneNumberID: phoneNumberID}
	return s.client.doJSON(ctx, "POST", "/voice/sip-profiles/"+url.PathEscape(id)+"/unassign-number", body, opts)
}
