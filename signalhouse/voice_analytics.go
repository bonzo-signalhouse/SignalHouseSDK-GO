package signalhouse

import (
	"context"
)

// VoiceAnalyticsService handles Voice Analytics API calls — aggregated call
// metrics for the portal Voice Analytics view. Wraps voice-backend's
// GET /voice/stats/voice-analytics surface.
//
// Distinct from VoiceCallLogsService (client.Voice.CallLogs), which is the raw
// call history read surface: this returns per-direction summary tiles (both
// directions) plus byDate / byNumber / byCarrier status breakdowns for charts.
// Access via client.Voice.Analytics.
type VoiceAnalyticsService struct {
	client *Client
}

// VoiceAnalyticsParams holds query filters for Get. All fields optional and
// omitted from the request when empty. Summary tiles honor every filter except
// Direction; breakdowns honor all filters.
type VoiceAnalyticsParams struct {
	// DateFrom is an ISO-8601 lower bound on call start time.
	DateFrom string
	// DateTo is an ISO-8601 upper bound on call start time.
	DateTo string
	// GroupBy is the bucket size for the byDate breakdown: day, week, or month
	// (default day).
	GroupBy string
	// Direction restricts the breakdowns to one direction: INBOUND or OUTBOUND.
	// Summary tiles always cover both directions.
	Direction string
	// CallSource is a comma-separated CallLog.callSource list (SIP_TRUNK,
	// PROGRAMMABLE_VOICE, AI_VOICE, WEBRTC, VOICEMAIL, CONFERENCE, PSTN).
	CallSource string
	// SubgroupID filters by owning subgroup.
	SubgroupID string
	// Number is a substring match against either leg (from OR to).
	Number string
	// Carrier filters by carrier.
	Carrier string
}

// VoiceAnalyticsSummary is the per-direction summary tile block. Nullable
// metrics are pointers and omitted when the server returns null.
type VoiceAnalyticsSummary struct {
	TotalCalls      int `json:"totalCalls"`
	SuccessfulCalls int `json:"successfulCalls"`
	// SuccessRate is a percentage (0–100).
	SuccessRate    float64 `json:"successRate"`
	AvgCallsPerDay float64 `json:"avgCallsPerDay"`
	// AvgCallDurationSeconds is null when there are no answered calls.
	AvgCallDurationSeconds *float64 `json:"avgCallDurationSeconds,omitempty"`
	TotalCallMinutes       float64  `json:"totalCallMinutes"`
	// AnswerRate is a percentage (ASR / answer seizure rate).
	AnswerRate float64 `json:"answerRate"`
	// AbandonmentRate is a percentage (CANCELED before answer).
	AbandonmentRate float64 `json:"abandonmentRate"`
	FailedCalls     int     `json:"failedCalls"`
	// FailureRate is a percentage.
	FailureRate float64 `json:"failureRate"`
	// PctUnder6s is the percentage of answered calls under 6 seconds.
	PctUnder6s float64 `json:"pctUnder6s"`
	// AvgRingDurationSeconds is the average answer_time − start_time; null when
	// unavailable.
	AvgRingDurationSeconds *float64 `json:"avgRingDurationSeconds,omitempty"`
	// SentimentScore is the positive-sentiment percentage; null when unscored.
	SentimentScore *float64 `json:"sentimentScore,omitempty"`
	// SentimentLabel is Positive, Negative, or Neutral; null when unscored.
	SentimentLabel *string `json:"sentimentLabel,omitempty"`
}

// VoiceAnalyticsBucket is one row of a byDate / byNumber / byCarrier breakdown.
type VoiceAnalyticsBucket struct {
	// Key is the date / number / carrier this bucket aggregates.
	Key                  string  `json:"key"`
	Total                int     `json:"total"`
	Completed            int     `json:"completed"`
	Failed               int     `json:"failed"`
	NoAnswer             int     `json:"noAnswer"`
	Busy                 int     `json:"busy"`
	Canceled             int     `json:"canceled"`
	InProgress           int     `json:"inProgress"`
	TotalDurationSeconds float64 `json:"totalDurationSeconds"`
	// AvgSentimentScore is null when no calls in the bucket are scored.
	AvgSentimentScore *float64 `json:"avgSentimentScore,omitempty"`
}

// VoiceAnalyticsRange is the resolved date window the analytics cover.
type VoiceAnalyticsRange struct {
	DateFrom string `json:"dateFrom"`
	DateTo   string `json:"dateTo"`
	Days     int    `json:"days"`
}

// VoiceAnalyticsResult is the typed shape returned by GET
// /voice/stats/voice-analytics. Unmarshal a Response body into this to read the
// analytics with types.
type VoiceAnalyticsResult struct {
	Success   bool                   `json:"success"`
	Range     VoiceAnalyticsRange    `json:"range"`
	Outbound  VoiceAnalyticsSummary  `json:"outbound"`
	Inbound   VoiceAnalyticsSummary  `json:"inbound"`
	ByDate    []VoiceAnalyticsBucket `json:"byDate"`
	ByNumber  []VoiceAnalyticsBucket `json:"byNumber"`
	ByCarrier []VoiceAnalyticsBucket `json:"byCarrier"`
	ByChannel []VoiceAnalyticsBucket `json:"byChannel"`
}

// Get returns aggregated voice analytics for the current account. Returns
// { success, range, outbound, inbound, byDate, byNumber, byCarrier, byChannel }; the body
// mirrors VoiceAnalyticsResult.
// GET /voice/stats/voice-analytics
func (s *VoiceAnalyticsService) Get(ctx context.Context, params VoiceAnalyticsParams, opts *RequestOptions) (*Response, error) {
	qs := buildQueryString(map[string]string{
		"dateFrom":   params.DateFrom,
		"dateTo":     params.DateTo,
		"groupBy":    params.GroupBy,
		"direction":  params.Direction,
		"callSource": params.CallSource,
		"subgroupId": params.SubgroupID,
		"number":     params.Number,
		"carrier":    params.Carrier,
	})
	return s.client.doJSON(ctx, "GET", "/voice/stats/voice-analytics"+qs, nil, opts)
}
