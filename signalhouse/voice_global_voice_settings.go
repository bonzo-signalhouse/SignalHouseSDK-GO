package signalhouse

import (
	"context"
)

// VoiceGlobalVoiceSettingsService handles Global Voice Settings API calls —
// account-wide voice defaults (accepted regions, max spend per minute, E911).
//
// SIP trunks and Programmable Voice Profiles override these. Wraps
// voice-backend's /voice/api/v1/global-voice-settings surface. Access via
// client.Voice.GlobalVoiceSettings.
type VoiceGlobalVoiceSettingsService struct {
	client *Client
}

// UpdateGlobalVoiceSettingsParams is the body for UpdateGlobalVoiceSettings.
// The update is partial — only the fields you provide are written; omitted
// fields are preserved.
type UpdateGlobalVoiceSettingsParams struct {
	// AcceptedRegions are accepted regions as country calling codes; empty = all allowed.
	AcceptedRegions []string `json:"acceptedRegions,omitempty"`
	// MaxSpendPerMinute is the baseline max spend per minute (dollars); null = no limit.
	MaxSpendPerMinute *float64 `json:"maxSpendPerMinute,omitempty"`
	// E911Enabled controls whether numbers may dial 911.
	E911Enabled *bool `json:"e911Enabled,omitempty"`
}

// GetGlobalVoiceSettings returns the current account's global voice settings.
// Returns data: null if the account has never configured them (defaults apply).
// GET /voice/api/v1/global-voice-settings
func (s *VoiceGlobalVoiceSettingsService) GetGlobalVoiceSettings(ctx context.Context, opts *RequestOptions) (*Response, error) {
	return s.client.doJSON(ctx, "GET", "/voice/api/v1/global-voice-settings", nil, opts)
}

// UpdateGlobalVoiceSettings upserts the account's global voice settings.
// Partial — only the fields you provide are written; omitted fields are preserved.
// PUT /voice/api/v1/global-voice-settings
func (s *VoiceGlobalVoiceSettingsService) UpdateGlobalVoiceSettings(ctx context.Context, params UpdateGlobalVoiceSettingsParams, opts *RequestOptions) (*Response, error) {
	return s.client.doJSON(ctx, "PUT", "/voice/api/v1/global-voice-settings", params, opts)
}
