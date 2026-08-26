package signalhouse

import (
	"context"
	"net/url"
)

// VoiceCallsService handles voice call API calls — outbound origination and
// call log queries. Wraps voice-backend's /voice/v1/calls surface.
//
// Origination modes:
//   - Single-leg (no trunk/profile/identity): Twilio-style direct outbound —
//     dial To from caller-ID From. From must be an account-owned number.
//   - SIP trunk (SipTrunkID): Route via a configured trunk; the trunk's
//     auth/ACL governs caller-ID rules.
//   - Two-leg-with-bridge (SipProfileID or ToIdentity): Ring the SDK-
//     registered endpoint first, then dial To and bridge. ToIdentity
//     resolves to the most recently-minted unexpired ephemeral identity
//     for that name.
type VoiceCallsService struct {
	client *Client
}

// CreateCallParams is the body for CreateCall.
type CreateCallParams struct {
	// To is the destination phone number in E.164 format (e.g. "+15551234567") or SIP URI.
	To string `json:"to"`
	// From is the caller ID to present on the outbound leg. Must be an
	// account-owned phone number when not using a SIP trunk.
	From string `json:"from"`
	// SipTrunkID routes the call via a configured SIP trunk. Mutually
	// exclusive with SipProfileID and ToIdentity for the trunk path.
	SipTrunkID string `json:"sip_trunk_id,omitempty"`
	// SipProfileID enables two-leg-with-bridge: ring this persistent SIP
	// profile first, then dial To and bridge. Mutually exclusive with ToIdentity.
	SipProfileID string `json:"sip_profile_id,omitempty"`
	// ToIdentity enables two-leg-with-bridge: ring the SDK-registered
	// identity by logical name (the value passed to Tokens.CreateToken with
	// Identity), then dial To and bridge. Mutually exclusive with SipProfileID.
	ToIdentity string `json:"to_identity,omitempty"`
	// AnswerURL is the webhook URL fetched on call answer; returns
	// call-control instructions.
	AnswerURL string `json:"answer_url,omitempty"`
	// AnswerMethod is the HTTP method for AnswerURL: "GET" or "POST". Defaults to POST.
	AnswerMethod string `json:"answer_method,omitempty"`
	// StatusCallback is the webhook URL for call status updates
	// (QUEUED, RINGING, ANSWERED, COMPLETED, FAILED).
	StatusCallback string `json:"status_callback,omitempty"`
	// StatusCallbackMethod is the HTTP method for StatusCallback: "GET" or "POST". Defaults to POST.
	StatusCallbackMethod string `json:"status_callback_method,omitempty"`
	// RecordingEnabled records both legs of the call. Defaults to false.
	RecordingEnabled *bool `json:"recording_enabled,omitempty"`
	// TranscriptionEnabled transcribes the recording. Defaults to false.
	TranscriptionEnabled *bool `json:"transcription_enabled,omitempty"`
	// IvrFlowID routes the call into an AI IVR flow on answer instead of
	// webhook control.
	IvrFlowID string `json:"ivr_flow_id,omitempty"`
	// Metadata is free-form metadata persisted on the call log.
	Metadata map[string]interface{} `json:"metadata,omitempty"`
}

// ListCallsParams holds query filters for ListCalls. All fields optional.
// Page and Limit are strings to match this SDK's query-param convention.
type ListCallsParams struct {
	// Page is the page number (1-indexed) as a string. Defaults server-side to "1".
	Page string
	// Limit is the page size (1-100) as a string. Defaults server-side to "20".
	Limit string
	// Status filters by call status: QUEUED, RINGING, IN_PROGRESS,
	// COMPLETED, FAILED, BUSY, NO_ANSWER, CANCELED.
	Status string
	// Direction filters by call direction: INBOUND or OUTBOUND.
	Direction string
	// From filters by caller number.
	From string
	// To filters by destination number.
	To string
	// DateFrom is an ISO-8601 lower bound on start_time.
	DateFrom string
	// DateTo is an ISO-8601 upper bound on start_time.
	DateTo string
}

type hangupCallBody struct {
	Reason string `json:"reason,omitempty"`
}

// CreateCall places an outbound call.
// POST /voice/v1/calls
func (s *VoiceCallsService) CreateCall(ctx context.Context, params CreateCallParams, opts *RequestOptions) (*Response, error) {
	if err := require(map[string]string{
		"to":   params.To,
		"from": params.From,
	}); err != nil {
		return nil, err
	}
	return s.client.doJSON(ctx, "POST", "/voice/v1/calls", params, opts)
}

// ListCalls returns paginated call logs for the current account.
// GET /voice/v1/calls
func (s *VoiceCallsService) ListCalls(ctx context.Context, params ListCallsParams, opts *RequestOptions) (*Response, error) {
	qs := buildQueryString(map[string]string{
		"page":      params.Page,
		"limit":     params.Limit,
		"status":    params.Status,
		"direction": params.Direction,
		"from":      params.From,
		"to":        params.To,
		"date_from": params.DateFrom,
		"date_to":   params.DateTo,
	})
	return s.client.doJSON(ctx, "GET", "/voice/v1/calls"+qs, nil, opts)
}

// GetCall returns a single call log by ID (the value returned as call_id
// from CreateCall).
// GET /voice/v1/calls/:id
func (s *VoiceCallsService) GetCall(ctx context.Context, id string, opts *RequestOptions) (*Response, error) {
	if err := require(map[string]string{"id": id}); err != nil {
		return nil, err
	}
	return s.client.doJSON(ctx, "GET", "/voice/v1/calls/"+url.PathEscape(id), nil, opts)
}

// HangupCall hangs up an in-progress call. Reason is optional and recorded
// on the call log; valid values are NORMAL, BUSY, NO_ANSWER, REJECTED.
// Defaults to NORMAL when reason is empty.
// POST /voice/v1/calls/:id/hangup
func (s *VoiceCallsService) HangupCall(ctx context.Context, id, reason string, opts *RequestOptions) (*Response, error) {
	if err := require(map[string]string{"id": id}); err != nil {
		return nil, err
	}
	body := hangupCallBody{Reason: reason}
	return s.client.doJSON(ctx, "POST", "/voice/v1/calls/"+url.PathEscape(id)+"/hangup", body, opts)
}
