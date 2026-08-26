package signalhouse

import (
	"context"
	"net/url"
)

// VoiceCallLogsService handles Call Logs API calls — the account-scoped voice
// call history. Wraps voice-backend's /voice/api/v1/call-logs surface.
//
// Distinct from VoiceCallsService (client.Voice.Calls), which is the
// Twilio-compatible programmatic call-control surface: this is the richer log
// read surface with cost, recording, and voice-config attribution used by the
// portal Call Logs view. Access via client.Voice.CallLogs.
type VoiceCallLogsService struct {
	client *Client
}

// ListCallLogsParams holds query filters for ListCallLogs. All fields optional.
// Page and Limit are strings to match this SDK's query-param convention.
type ListCallLogsParams struct {
	// Page is the page number (1-indexed) as a string. Defaults server-side to "1".
	Page string
	// Limit is the page size (1-1000) as a string. Defaults server-side to "20".
	Limit string
	// Direction filters by call direction: INBOUND, OUTBOUND, or INTERNAL.
	Direction string
	// Status filters by call status: QUEUED, RINGING, IN_PROGRESS,
	// COMPLETED, FAILED, BUSY, NO_ANSWER, CANCELED.
	Status string
	// CallSource filters by how the call was handled: SIP_TRUNK,
	// PROGRAMMABLE_VOICE, AI_VOICE, WEBRTC, VOICEMAIL, CONFERENCE, PSTN.
	CallSource string
	// SubgroupID filters by owning subgroup.
	SubgroupID string
	// SipTrunkID filters by SIP trunk id.
	SipTrunkID string
	// From filters by caller number (contains).
	From string
	// To filters by called number (contains).
	To string
	// Q is a free-text search over callId / from / to.
	Q string
	// Sentiment filters by sentiment: positive, negative, or neutral.
	Sentiment string
	// Voicemail, when "true", returns only voicemails.
	Voicemail string
	// DateFrom is an ISO-8601 lower bound on start time.
	DateFrom string
	// DateTo is an ISO-8601 upper bound on start time.
	DateTo string
	// Sort is the sort field (default startTime).
	Sort string
	// Order is the sort order: asc or desc (default desc).
	Order string
}

// ListCallLogs returns paginated call logs for the current account.
// Returns { data: [...call logs], meta: { total, page, limit, totalPages } }.
// GET /voice/api/v1/call-logs
func (s *VoiceCallLogsService) ListCallLogs(ctx context.Context, params ListCallLogsParams, opts *RequestOptions) (*Response, error) {
	qs := buildQueryString(map[string]string{
		"page":       params.Page,
		"limit":      params.Limit,
		"direction":  params.Direction,
		"status":     params.Status,
		"callSource": params.CallSource,
		"subgroupId": params.SubgroupID,
		"sipTrunkId": params.SipTrunkID,
		"from":       params.From,
		"to":         params.To,
		"q":          params.Q,
		"sentiment":  params.Sentiment,
		"voicemail":  params.Voicemail,
		"dateFrom":   params.DateFrom,
		"dateTo":     params.DateTo,
		"sort":       params.Sort,
		"order":      params.Order,
	})
	return s.client.doJSON(ctx, "GET", "/voice/api/v1/call-logs"+qs, nil, opts)
}

// GetCallLog returns a single call log by its UUID or callId.
// Returns { data: { ...call log } }.
// GET /voice/api/v1/call-logs/:id
func (s *VoiceCallLogsService) GetCallLog(ctx context.Context, id string, opts *RequestOptions) (*Response, error) {
	if err := require(map[string]string{"id": id}); err != nil {
		return nil, err
	}
	return s.client.doJSON(ctx, "GET", "/voice/api/v1/call-logs/"+url.PathEscape(id), nil, opts)
}

// GetCallLogRecording returns a short-TTL presigned playback URL for a call's
// recording. Account-scoped — only the owning account can presign. 404 if the
// call has no recording. Returns { data: { callId, url, expiresIn } }.
// GET /voice/api/v1/call-logs/:id/recording
func (s *VoiceCallLogsService) GetCallLogRecording(ctx context.Context, id string, opts *RequestOptions) (*Response, error) {
	if err := require(map[string]string{"id": id}); err != nil {
		return nil, err
	}
	return s.client.doJSON(ctx, "GET", "/voice/api/v1/call-logs/"+url.PathEscape(id)+"/recording", nil, opts)
}

type markVoicemailReadBody struct {
	Read bool `json:"read"`
}

// MarkVoicemailRead marks a voicemail (a call log with a recording) read or
// unread. Account-scoped — only the owning account can update. Pass read=true
// to mark read, read=false to mark unread.
// Returns { data: { id, voicemailReadAt } }.
// PATCH /voice/api/v1/call-logs/:id/voicemail-read
func (s *VoiceCallLogsService) MarkVoicemailRead(ctx context.Context, id string, read bool, opts *RequestOptions) (*Response, error) {
	if err := require(map[string]string{"id": id}); err != nil {
		return nil, err
	}
	body := markVoicemailReadBody{Read: read}
	return s.client.doJSON(ctx, "PATCH", "/voice/api/v1/call-logs/"+url.PathEscape(id)+"/voicemail-read", body, opts)
}
