package signalhouse

import (
	"context"
	"net/url"
)

// VoiceSipTrunksService handles SIP trunk API calls.
type VoiceSipTrunksService struct {
	client *Client
}

// SipTrunkHost is an entry in destinationHosts / sourceHosts. Port is only
// used on destination hosts.
type SipTrunkHost struct {
	Host string `json:"host"`
	Port string `json:"port,omitempty"`
}

// CreateSipTrunkParams is the body for CreateSipTrunk.
type CreateSipTrunkParams struct {
	Name              string         `json:"name"`
	Region            string         `json:"region"`
	ConnectionType    string         `json:"connectionType"` // "IP_AUTH" or "REGISTRATION"
	AllowedIPs        []string       `json:"allowedIps,omitempty"`
	Transport         string         `json:"transport,omitempty"` // "tcp" | "udp" | "tls"
	MaxSpendPerMinute *float64       `json:"maxSpendPerMinute,omitempty"`
	DestinationHosts  []SipTrunkHost `json:"destinationHosts,omitempty"`
	SourceHosts       []SipTrunkHost `json:"sourceHosts,omitempty"`
}

// UpdateSipTrunkParams is the body for UpdateSipTrunk. All fields optional.
type UpdateSipTrunkParams struct {
	Name              string         `json:"name,omitempty"`
	Region            string         `json:"region,omitempty"`
	ConnectionType    string         `json:"connectionType,omitempty"`
	AllowedIPs        []string       `json:"allowedIps,omitempty"`
	Transport         string         `json:"transport,omitempty"`
	MaxSpendPerMinute *float64       `json:"maxSpendPerMinute,omitempty"`
	DestinationHosts  []SipTrunkHost `json:"destinationHosts,omitempty"`
	SourceHosts       []SipTrunkHost `json:"sourceHosts,omitempty"`
}

// ListSipTrunks returns all SIP trunks for the current account.
// GET /voice/sip-trunks
func (s *VoiceSipTrunksService) ListSipTrunks(ctx context.Context, opts *RequestOptions) (*Response, error) {
	return s.client.doJSON(ctx, "GET", "/voice/sip-trunks", nil, opts)
}

// GetSipTrunk returns a single SIP trunk by ID.
// GET /voice/sip-trunks/:id
func (s *VoiceSipTrunksService) GetSipTrunk(ctx context.Context, id string, opts *RequestOptions) (*Response, error) {
	if err := require(map[string]string{"id": id}); err != nil {
		return nil, err
	}
	return s.client.doJSON(ctx, "GET", "/voice/sip-trunks/"+url.PathEscape(id), nil, opts)
}

// CreateSipTrunk creates a new SIP trunk.
// Registration trunks return a one-time-visible password on the response.
// POST /voice/sip-trunks
func (s *VoiceSipTrunksService) CreateSipTrunk(ctx context.Context, params CreateSipTrunkParams, opts *RequestOptions) (*Response, error) {
	if err := require(map[string]string{
		"name":           params.Name,
		"region":         params.Region,
		"connectionType": params.ConnectionType,
	}); err != nil {
		return nil, err
	}
	return s.client.doJSON(ctx, "POST", "/voice/sip-trunks", params, opts)
}

// UpdateSipTrunk updates an existing SIP trunk.
// PATCH /voice/sip-trunks/:id
func (s *VoiceSipTrunksService) UpdateSipTrunk(ctx context.Context, id string, params UpdateSipTrunkParams, opts *RequestOptions) (*Response, error) {
	if err := require(map[string]string{"id": id}); err != nil {
		return nil, err
	}
	return s.client.doJSON(ctx, "PATCH", "/voice/sip-trunks/"+url.PathEscape(id), params, opts)
}

// DeleteSipTrunk deletes a SIP trunk.
// DELETE /voice/sip-trunks/:id
func (s *VoiceSipTrunksService) DeleteSipTrunk(ctx context.Context, id string, opts *RequestOptions) (*Response, error) {
	if err := require(map[string]string{"id": id}); err != nil {
		return nil, err
	}
	return s.client.doJSON(ctx, "DELETE", "/voice/sip-trunks/"+url.PathEscape(id), nil, opts)
}

// ToggleSipTrunkActive flips a trunk's active/inactive status.
// POST /voice/sip-trunks/:id/toggle-active
func (s *VoiceSipTrunksService) ToggleSipTrunkActive(ctx context.Context, id string, opts *RequestOptions) (*Response, error) {
	if err := require(map[string]string{"id": id}); err != nil {
		return nil, err
	}
	return s.client.doJSON(ctx, "POST", "/voice/sip-trunks/"+url.PathEscape(id)+"/toggle-active", nil, opts)
}

// RegenerateSipTrunkPassword rotates the SIP password on a REGISTRATION-type
// trunk. The old password stops working immediately. Returns { password: "..." }.
// POST /voice/sip-trunks/:id/regenerate-password
func (s *VoiceSipTrunksService) RegenerateSipTrunkPassword(ctx context.Context, id string, opts *RequestOptions) (*Response, error) {
	if err := require(map[string]string{"id": id}); err != nil {
		return nil, err
	}
	return s.client.doJSON(ctx, "POST", "/voice/sip-trunks/"+url.PathEscape(id)+"/regenerate-password", nil, opts)
}

// GetSipPops returns available SIP Points of Presence.
// GET /voice/sip-trunks/pops
func (s *VoiceSipTrunksService) GetSipPops(ctx context.Context, opts *RequestOptions) (*Response, error) {
	return s.client.doJSON(ctx, "GET", "/voice/sip-trunks/pops", nil, opts)
}

// sipTrunkNumbersParams is the body for AssignNumbers / UnassignNumbers.
type sipTrunkNumbersParams struct {
	PhoneNumbers []string `json:"phoneNumbers"`
}

// AssignNumbers attaches one or more phone numbers (by E.164) to a SIP trunk,
// routing each DID down the trunk. Idempotent; numbers already configured on a
// SIP profile are skipped and returned in `skipped` with a reason.
// Returns { assigned: [...], skipped: [...] }.
// POST /voice/sip-trunks/:id/assign-numbers
func (s *VoiceSipTrunksService) AssignNumbers(ctx context.Context, id string, phoneNumbers []string, opts *RequestOptions) (*Response, error) {
	if err := require(map[string]string{"id": id}); err != nil {
		return nil, err
	}
	if len(phoneNumbers) == 0 {
		return nil, &ValidationError{Field: "phoneNumbers"}
	}
	return s.client.doJSON(ctx, "POST", "/voice/sip-trunks/"+url.PathEscape(id)+"/assign-numbers", sipTrunkNumbersParams{PhoneNumbers: phoneNumbers}, opts)
}

// UnassignNumbers detaches one or more phone numbers (by E.164) from a SIP
// trunk, removing each number's SIP_TRUNK inbound route when it points at this
// trunk. Returns { unassigned: [...] }.
// POST /voice/sip-trunks/:id/unassign-numbers
func (s *VoiceSipTrunksService) UnassignNumbers(ctx context.Context, id string, phoneNumbers []string, opts *RequestOptions) (*Response, error) {
	if err := require(map[string]string{"id": id}); err != nil {
		return nil, err
	}
	if len(phoneNumbers) == 0 {
		return nil, &ValidationError{Field: "phoneNumbers"}
	}
	return s.client.doJSON(ctx, "POST", "/voice/sip-trunks/"+url.PathEscape(id)+"/unassign-numbers", sipTrunkNumbersParams{PhoneNumbers: phoneNumbers}, opts)
}
