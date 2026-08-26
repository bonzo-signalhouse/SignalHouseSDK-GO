package signalhouse

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/url"
)

// MessagesService handles message-related API calls.
//
// Date-range parameters (StartDate, EndDate) on list and analytics methods
// accept ISO-8601 date or timestamp strings. The server normalizes them to
// full UTC-day boundaries: StartDate -> start-of-UTC-day, EndDate -> end-of-
// UTC-day. The end day is always included; hourly resolution is not supported.
// Business-number filters and SMS/MMS senders accept digits-only 5-6 digit
// Short Codes or 10+ digit long numbers. Recipient numbers remain 10+ digits.
type MessagesService struct {
	client *Client
}

// GetMessagesParams are the parameters for the GetMessages method.
type GetMessagesParams struct {
	ID                   string
	CampaignID           string
	BrandID              string
	SubgroupID           string
	GroupID              string
	PhoneNumber          string
	SenderPhoneNumber    string
	RecipientPhoneNumber string
	Status               string
	Direction            string
	// MessageType accepts a single value (e.g. "SMS") or a comma-separated list (e.g. "SMS,MMS").
	// Allowed values: SMS, MMS, RCS, WHATSAPP, VIBER, P2P. Empty defaults to all non-P2P types.
	MessageType string
	// Channel filters by "tenDLC", "tollFree", "shortCode", or "p2p" (single value or comma-separated list).
	// Older messages without a stored channel count as tenDLC.
	Channel string
	Carrier string
	// SentimentLabel filters by sentiment bucket: "positive", "neutral", or "negative" (single value or
	// comma-separated list). Applied as the sentimentScore range the label threshold produces, so unscored
	// messages match no bucket and outbound messages (never scored) are excluded.
	SentimentLabel string
	StartDate      string
	EndDate        string
	// SortField is one of "createdAt", "segmentCount", "status", or "sentimentScore".
	SortField string
	SortOrder string
	Page      string
	Limit     string
}

// GetMessages returns a list of messages with optional filters.
//
// Inbound messages in the response also carry sentimentScore (-100..100, null when unscored),
// sentimentLabel ("positive", "neutral" or "negative", null when sentimentScore is null) and
// sentimentScoredAt. Outbound messages are never scored.
func (s *MessagesService) GetMessages(ctx context.Context, params GetMessagesParams, opts *RequestOptions) (*Response, error) {
	qs := buildQueryString(map[string]string{
		"id":                   params.ID,
		"campaignId":           params.CampaignID,
		"brandId":              params.BrandID,
		"subgroupId":           params.SubgroupID,
		"groupId":              params.GroupID,
		"phoneNumber":          params.PhoneNumber,
		"senderPhoneNumber":    params.SenderPhoneNumber,
		"recipientPhoneNumber": params.RecipientPhoneNumber,
		"status":               params.Status,
		"direction":            params.Direction,
		"messageType":          params.MessageType,
		"channel":              params.Channel,
		"carrier":              params.Carrier,
		"sentimentLabel":       params.SentimentLabel,
		"startDate":            params.StartDate,
		"endDate":              params.EndDate,
		"sortField":            params.SortField,
		"sortOrder":            params.SortOrder,
		"page":                 params.Page,
		"limit":                params.Limit,
	})
	return s.client.doJSON(ctx, "GET", "/message"+qs, nil, opts)
}

// GetAnalyticsParams are the parameters for the GetAnalytics method.
//
// SubgroupID, BrandID, CampaignID, PhoneNumber and Carrier each accept a single value or a
// comma-separated list of values. GroupID accepts a single value only.
type GetAnalyticsParams struct {
	GroupID     string
	SubgroupID  string
	BrandID     string
	CampaignID  string
	PhoneNumber string
	Carrier     string
	StartDate   string
	EndDate     string
	// Channel filters by "tenDLC", "tollFree", "shortCode", or "p2p" (single value or comma-separated list).
	// A tenDLC selection also includes older messages with no channel, and p2p is matched by carrier.
	// MessageType scopes the SENTIMENT figures to "SMS", "MMS" or "P2P" (single value or
	// comma-separated list). Empty means all scored types. It does NOT narrow the message counts,
	// which are always returned split per type and have no filterable type dimension.
	MessageType string
	Channel     string
}

// GetAnalytics returns aggregated analytics for messages.
//
// The response also includes inbound-message sentiment for the requested scope: sentimentScore
// (volume-weighted average of every scored inbound message in range, -100..100, null when nothing
// was scored), sentimentLabel, sentimentScoredCount, sentimentPositive, sentimentNeutral and
// sentimentNegative. sentimentScore is an average, so it must never be summed across responses.
func (s *MessagesService) GetAnalytics(ctx context.Context, params GetAnalyticsParams, opts *RequestOptions) (*Response, error) {
	qs := buildQueryString(map[string]string{
		"groupId":     params.GroupID,
		"subgroupId":  params.SubgroupID,
		"brandId":     params.BrandID,
		"campaignId":  params.CampaignID,
		"phoneNumber": params.PhoneNumber,
		"carrier":     params.Carrier,
		"startDate":   params.StartDate,
		"endDate":     params.EndDate,
		"channel":     params.Channel,
		"messageType": params.MessageType,
	})
	return s.client.doJSON(ctx, "GET", "/message/analytics"+qs, nil, opts)
}

// GetAnalyticsDetailParams are the parameters for the GetAnalyticsDetail method.
//
// SubgroupID, BrandID, CampaignID, PhoneNumber and Carrier each accept a single value or a
// comma-separated list of values. GroupID accepts a single value only.
type GetAnalyticsDetailParams struct {
	GroupID     string
	SubgroupID  string
	BrandID     string
	CampaignID  string
	PhoneNumber string
	Carrier     string
	StartDate   string
	EndDate     string
	// Channel filters by "tenDLC", "tollFree", "shortCode", or "p2p" (single value or comma-separated list).
	// A tenDLC selection also includes older messages with no channel, and p2p is a real channel here.
	// MessageType scopes the SENTIMENT figures to "SMS", "MMS" or "P2P" (single value or
	// comma-separated list). Empty means all scored types. It does NOT narrow the message counts,
	// which are always returned split per type and have no filterable type dimension.
	MessageType string
	Channel     string
}

// GetAnalyticsDetail returns detailed analytics snapshot records for charting and aggregation.
//
// cardTotals and every byDate/byPhoneNumber/byCarrier row also carry that row's own inbound-message
// sentiment: sentimentScore (volume-weighted average, -100..100, null when nothing in that row's
// scope was scored), sentimentLabel, sentimentScoredCount, sentimentPositive, sentimentNeutral and
// sentimentNegative. The fields are always present; sentimentScore is an average and must never be
// summed across rows or responses.
func (s *MessagesService) GetAnalyticsDetail(ctx context.Context, params GetAnalyticsDetailParams, opts *RequestOptions) (*Response, error) {
	qs := buildQueryString(map[string]string{
		"groupId":     params.GroupID,
		"subgroupId":  params.SubgroupID,
		"brandId":     params.BrandID,
		"campaignId":  params.CampaignID,
		"phoneNumber": params.PhoneNumber,
		"carrier":     params.Carrier,
		"startDate":   params.StartDate,
		"endDate":     params.EndDate,
		"channel":     params.Channel,
		"messageType": params.MessageType,
	})
	return s.client.doJSON(ctx, "GET", "/message/analytics/detail"+qs, nil, opts)
}

// GetAnalyticsFilterOptionsParams are the parameters for the GetAnalyticsFilterOptions method.
type GetAnalyticsFilterOptionsParams struct {
	GroupID string
}

// GetAnalyticsFilterOptions returns filter dropdown options (subgroups, brands, campaigns, phone numbers)
// for the Analytics page, scoped to a single group. Sourced from ClickHouse — only items with at least
// one message are returned.
func (s *MessagesService) GetAnalyticsFilterOptions(ctx context.Context, params GetAnalyticsFilterOptionsParams, opts *RequestOptions) (*Response, error) {
	qs := buildQueryString(map[string]string{
		"groupId": params.GroupID,
	})
	return s.client.doJSON(ctx, "GET", "/message/analytics/filter-options"+qs, nil, opts)
}

// GetAnalyticsBySubgroupParams are the parameters for the GetAnalyticsBySubgroup method.
// Limit is capped server-side at 50 per page.
// Channel: "both" (default when empty, all activity) or "tenDLC" / "tollFree" / "shortCode" / "p2p" (single value
// or comma-separated list) — scopes ORDER BY + row inclusion by the stored channel column so a
// single-channel caller doesn't get pages dominated by other channels; toll-free is kept separate
// from 10DLC.
type GetAnalyticsBySubgroupParams struct {
	GroupID     string
	SubgroupID  string
	BrandID     string
	CampaignID  string
	PhoneNumber string
	Carrier     string
	StartDate   string
	EndDate     string
	Page        string
	Limit       string
	Channel     string
	// MessageType scopes ONLY the sentiment figures (SMS, MMS or P2P; comma-separated for several).
	// Omit for all scored types. It does not narrow the message counts, which are always returned
	// split per type and have no filterable type dimension.
	MessageType string
}

// GetAnalyticsBySubgroup returns a paginated breakdown of message metrics grouped by subgroup.
// Each row carries the full set of sms/mms/p2p metric columns plus smsOptOuts/mmsOptOuts (sourced
// from the DNC analytics MV — P2P has no opt-outs) so callers can apply channel toggles client-side.
// Each row also carries that subgroup's inbound-message sentiment: sentimentScore (volume-weighted,
// -100..100, null when nothing was scored), sentimentLabel, sentimentScoredCount and the
// positive/neutral/negative counts, plus the same six per message type under smsSentiment*,
// mmsSentiment* and p2pSentiment*. The score is derived server-side; combine counts across rows,
// never the scores.
func (s *MessagesService) GetAnalyticsBySubgroup(ctx context.Context, params GetAnalyticsBySubgroupParams, opts *RequestOptions) (*Response, error) {
	qs := buildQueryString(map[string]string{
		"groupId":     params.GroupID,
		"subgroupId":  params.SubgroupID,
		"brandId":     params.BrandID,
		"campaignId":  params.CampaignID,
		"phoneNumber": params.PhoneNumber,
		"carrier":     params.Carrier,
		"startDate":   params.StartDate,
		"endDate":     params.EndDate,
		"page":        params.Page,
		"limit":       params.Limit,
		"channel":     params.Channel,
		"messageType": params.MessageType,
	})
	return s.client.doJSON(ctx, "GET", "/message/analytics/by-subgroup"+qs, nil, opts)
}

// GetAnalyticsByErrorCodeParams are the parameters for the GetAnalyticsByErrorCode method.
// Limit is capped server-side at 50 per page.
// Channel: "both" (default when empty, every code) or "tenDLC" / "tollFree" / "shortCode" / "p2p" (single value
// or comma-separated list) — scopes ORDER BY + totalCount by the stored channel column; when one
// channel is selected, totalErrors reflects only that channel.
type GetAnalyticsByErrorCodeParams struct {
	GroupID     string
	SubgroupID  string
	BrandID     string
	CampaignID  string
	PhoneNumber string
	Carrier     string
	StartDate   string
	EndDate     string
	Page        string
	Limit       string
	Channel     string
}

// GetAnalyticsByErrorCode returns a paginated breakdown of failed messages grouped by error code.
// Each row contains per-channel (sms/mms/p2p) error counts plus an enriched description.
func (s *MessagesService) GetAnalyticsByErrorCode(ctx context.Context, params GetAnalyticsByErrorCodeParams, opts *RequestOptions) (*Response, error) {
	qs := buildQueryString(map[string]string{
		"groupId":     params.GroupID,
		"subgroupId":  params.SubgroupID,
		"brandId":     params.BrandID,
		"campaignId":  params.CampaignID,
		"phoneNumber": params.PhoneNumber,
		"carrier":     params.Carrier,
		"startDate":   params.StartDate,
		"endDate":     params.EndDate,
		"page":        params.Page,
		"limit":       params.Limit,
		"channel":     params.Channel,
	})
	return s.client.doJSON(ctx, "GET", "/message/analytics/by-error-code"+qs, nil, opts)
}

// GetDncAnalyticsParams are the parameters for the GetDncAnalytics method.
type GetDncAnalyticsParams struct {
	GroupID     string
	SubgroupID  string
	BrandID     string
	CampaignID  string
	PhoneNumber string
	Carrier     string
	StartDate   string
	EndDate     string
	// Channel filters by "tenDLC", "tollFree", or "shortCode" (single value or comma-separated list). Opt-outs are
	// A2P-only, so "p2p" applies no filter, and a tenDLC selection also includes older opt-outs with
	// no channel.
	Channel string
}

// GetDncAnalytics returns aggregated DNC (Do Not Contact) opt-out analytics.
func (s *MessagesService) GetDncAnalytics(ctx context.Context, params GetDncAnalyticsParams, opts *RequestOptions) (*Response, error) {
	qs := buildQueryString(map[string]string{
		"groupId":     params.GroupID,
		"subgroupId":  params.SubgroupID,
		"brandId":     params.BrandID,
		"campaignId":  params.CampaignID,
		"phoneNumber": params.PhoneNumber,
		"carrier":     params.Carrier,
		"startDate":   params.StartDate,
		"endDate":     params.EndDate,
		"channel":     params.Channel,
	})
	return s.client.doJSON(ctx, "GET", "/message/dnc/analytics"+qs, nil, opts)
}

// GetDncRecordsParams are the parameters for the GetDncRecords method.
type GetDncRecordsParams struct {
	GroupID     string
	SubgroupID  string
	BrandID     string
	CampaignID  string
	PhoneNumber string
	Carrier     string
	StartDate   string
	EndDate     string
	Page        string
	Limit       string
	SortField   string
	SortOrder   string
	// Channel filters by "tenDLC", "tollFree", or "shortCode" (single value or comma-separated list). Opt-outs are
	// A2P-only, so "p2p" applies no filter, and a tenDLC selection also includes older opt-outs with
	// no channel.
	Channel string
}

// GetDncRecords returns paginated Do Not Call records.
func (s *MessagesService) GetDncRecords(ctx context.Context, params GetDncRecordsParams, opts *RequestOptions) (*Response, error) {
	qs := buildQueryString(map[string]string{
		"groupId":     params.GroupID,
		"subgroupId":  params.SubgroupID,
		"brandId":     params.BrandID,
		"campaignId":  params.CampaignID,
		"phoneNumber": params.PhoneNumber,
		"carrier":     params.Carrier,
		"startDate":   params.StartDate,
		"endDate":     params.EndDate,
		"page":        params.Page,
		"limit":       params.Limit,
		"sortField":   params.SortField,
		"sortOrder":   params.SortOrder,
		"channel":     params.Channel,
	})
	return s.client.doJSON(ctx, "GET", "/message/dnc"+qs, nil, opts)
}

// SendSMSParams are the parameters for the SendSMS method.
type SendSMSParams struct {
	SenderPhoneNumber                 string   `json:"senderPhoneNumber"`
	RecipientPhoneNumbers             []string `json:"recipientPhoneNumber"`
	MessageBody                       string   `json:"messageBody"`
	StatusCallbackURL                 string   `json:"statusCallbackUrl,omitempty"`
	EnableShortlink                   bool     `json:"enableShortlink"`
	FilterLandlinesAndInactiveNumbers bool     `json:"filterLandlinesAndInactiveNumbers,omitempty"`
}

// SendSMS sends an SMS message to one or more recipients. SenderPhoneNumber accepts a digits-only
// 5-6 digit Short Code or 10+ digit long number. A Short Code sender permits exactly one recipient;
// the API rejects Short Code sends with multiple recipients.
//
// The response reports the send outcome via requestedRecipientCount, enqueuedCount, and failedCount.
// enqueuedCount is the number of recipients accepted for delivery and is 0 when the send was fully
// blocked (by the moderation filter or a campaign opt-out); the request still succeeds in that case,
// so treat enqueuedCount, not the HTTP status, as the signal that a message was accepted. Blocked
// recipients are returned in insertedMessages with status "FAILED"; recipients on the global Do Not
// Contact list are reported in dncBlockedNumbers with no message record. Each FAILED message carries
// errorCode, a stable machine-readable cause (campaign opt-out is "OUT"; null on success or on
// failures with no assigned code) that is also included on the FAILED status-callback payload.
func (s *MessagesService) SendSMS(ctx context.Context, params SendSMSParams, opts *RequestOptions) (*Response, error) {
	if err := require(map[string]string{
		"senderPhoneNumber": params.SenderPhoneNumber,
		"messageBody":       params.MessageBody,
	}); err != nil {
		return nil, err
	}
	if len(params.RecipientPhoneNumbers) == 0 {
		return nil, &ValidationError{Field: "recipientPhoneNumbers"}
	}
	return s.client.doJSON(ctx, "POST", "/message/sms", params, opts)
}

// SendP2PParams are the parameters for the SendP2P method.
type SendP2PParams struct {
	RecipientPhoneNumbers []string `json:"recipientPhoneNumber"`
	MessageBody           string   `json:"messageBody"`
	StatusCallbackURL     string   `json:"statusCallbackUrl,omitempty"`
	// GroupID is optional. Explicit group to send from; defaults to the group in your auth context.
	GroupID string `json:"groupId,omitempty"`
	// SubgroupID is optional. Explicit subgroup to attribute the send (and its inbound reply) to. P2P has
	// no dedicated sender number, so pass this to record the send against a specific subgroup; it also
	// resolves the group. Defaults to the group's oldest subgroup.
	SubgroupID string `json:"subgroupId,omitempty"`
	// UseSignalHouseShortlinks is optional. When false, SignalHouse applies no
	// link-shortening or text-spin and your pre-spun links/content are sent verbatim
	// (bring-your-own spinner). Defaults to true. Pointer so the literal false is sent
	// while a nil pointer omits the field and preserves the server default.
	UseSignalHouseShortlinks *bool `json:"useSignalHouseShortlinks,omitempty"`
}

// SendP2P sends a P2P message via Rogue Mobile SMPP.
func (s *MessagesService) SendP2P(ctx context.Context, params SendP2PParams, opts *RequestOptions) (*Response, error) {
	if err := require(map[string]string{
		"messageBody": params.MessageBody,
	}); err != nil {
		return nil, err
	}
	if len(params.RecipientPhoneNumbers) == 0 {
		return nil, &ValidationError{Field: "recipientPhoneNumbers"}
	}
	return s.client.doJSON(ctx, "POST", "/message/p2p", params, opts)
}

// GetP2PBatchParams are the parameters for the GetP2PBatch method.
type GetP2PBatchParams struct {
	BatchID string
}

// GetP2PBatch returns a P2P batch by its ID.
func (s *MessagesService) GetP2PBatch(ctx context.Context, params GetP2PBatchParams, opts *RequestOptions) (*Response, error) {
	if err := require(map[string]string{"batchId": params.BatchID}); err != nil {
		return nil, err
	}
	path := fmt.Sprintf("/message/p2p/batches/%s", url.PathEscape(params.BatchID))
	return s.client.doJSON(ctx, "GET", path, nil, opts)
}

// GetDeliveryReportsParams are the parameters for the GetDeliveryReports method.
type GetDeliveryReportsParams struct {
	MessageIDs []string `json:"messageIds"`
}

// GetDeliveryReports fetches the latest upstream carrier delivery report for a set of messages by ID,
// useful for reconciling messages stuck in a non-terminal state. SMS and MMS messages are looked up
// against the carrier; other message types are returned with their stored status.
func (s *MessagesService) GetDeliveryReports(ctx context.Context, params GetDeliveryReportsParams, opts *RequestOptions) (*Response, error) {
	if len(params.MessageIDs) == 0 {
		return nil, &ValidationError{Field: "messageIds"}
	}
	return s.client.doJSON(ctx, "POST", "/message/delivery-report", params, opts)
}

// SendMMSParams are the parameters for the SendMMS method.
type SendMMSParams struct {
	SenderPhoneNumber                 string        `json:"-"`
	RecipientPhoneNumbers             []string      `json:"-"`
	MessageBody                       string        `json:"-"`
	MediaURLs                         []string      `json:"-"`
	StatusCallbackURL                 string        `json:"-"`
	EnableShortlink                   bool          `json:"-"`
	EnableCompression                 bool          `json:"-"`
	FilterLandlinesAndInactiveNumbers bool          `json:"-"`
	Images                            []*FileUpload `json:"-"`
}

// SendMMS sends an MMS message with optional media attachments. SenderPhoneNumber accepts a
// digits-only 5-6 digit Short Code or 10+ digit long number. A Short Code sender permits exactly
// one recipient; the API rejects Short Code sends with multiple recipients.
//
// The response reports the send outcome via requestedRecipientCount, enqueuedCount, and failedCount.
// enqueuedCount is the number of recipients accepted for delivery and is 0 when the send was fully
// blocked (by the moderation filter or a campaign opt-out); the request still succeeds in that case,
// so treat enqueuedCount, not the HTTP status, as the signal that a message was accepted. Blocked
// recipients are returned in insertedMessages with status "FAILED"; recipients on the global Do Not
// Contact list are reported in dncBlockedNumbers with no message record. Each FAILED message carries
// errorCode, a stable machine-readable cause (campaign opt-out is "OUT"; null on success or on
// failures with no assigned code) that is also included on the FAILED status-callback payload.
func (s *MessagesService) SendMMS(ctx context.Context, params SendMMSParams, opts *RequestOptions) (*Response, error) {
	if err := require(map[string]string{
		"senderPhoneNumber": params.SenderPhoneNumber,
	}); err != nil {
		return nil, err
	}
	if len(params.RecipientPhoneNumbers) == 0 {
		return nil, &ValidationError{Field: "recipientPhoneNumbers"}
	}

	writeForm := func(w *multipart.Writer) error {
		return writeMMSForm(w, params.SenderPhoneNumber, params.RecipientPhoneNumbers,
			params.MessageBody, params.MediaURLs, params.StatusCallbackURL,
			params.EnableShortlink, params.EnableCompression, params.FilterLandlinesAndInactiveNumbers, params.Images)
	}

	return s.client.doMultipart(ctx, "POST", "/message/mms", writeForm, opts)
}

// SendGroupMessageParams are the parameters for the SendGroupMessage method.
type SendGroupMessageParams struct {
	SenderPhoneNumber                 string        `json:"-"`
	RecipientPhoneNumbers             []string      `json:"-"`
	MessageBody                       string        `json:"-"`
	MediaURLs                         []string      `json:"-"`
	StatusCallbackURL                 string        `json:"-"`
	EnableShortlink                   bool          `json:"-"`
	EnableCompression                 bool          `json:"-"`
	FilterLandlinesAndInactiveNumbers bool          `json:"-"`
	Images                            []*FileUpload `json:"-"`
}

// SendGroupMessage sends a group MMS message with optional media attachments. Short Codes do not
// support group messaging; SenderPhoneNumber must be a long number.
//
// The response reports the send outcome via requestedRecipientCount, enqueuedCount, and failedCount.
// enqueuedCount is the number of recipients accepted for delivery and is 0 when the send was fully
// blocked (by the moderation filter or a campaign opt-out); the request still succeeds in that case,
// so treat enqueuedCount, not the HTTP status, as the signal that a message was accepted. Blocked
// recipients are returned in insertedMessages with status "FAILED"; recipients on the global Do Not
// Contact list are reported in dncBlockedNumbers with no message record. Each FAILED message carries
// errorCode, a stable machine-readable cause (campaign opt-out is "OUT"; null on success or on
// failures with no assigned code) that is also included on the FAILED status-callback payload.
func (s *MessagesService) SendGroupMessage(ctx context.Context, params SendGroupMessageParams, opts *RequestOptions) (*Response, error) {
	if err := require(map[string]string{
		"senderPhoneNumber": params.SenderPhoneNumber,
	}); err != nil {
		return nil, err
	}
	if len(params.RecipientPhoneNumbers) == 0 {
		return nil, &ValidationError{Field: "recipientPhoneNumbers"}
	}

	writeForm := func(w *multipart.Writer) error {
		return writeMMSForm(w, params.SenderPhoneNumber, params.RecipientPhoneNumbers,
			params.MessageBody, params.MediaURLs, params.StatusCallbackURL,
			params.EnableShortlink, params.EnableCompression, params.FilterLandlinesAndInactiveNumbers, params.Images)
	}

	return s.client.doMultipart(ctx, "POST", "/message/groupMessage", writeForm, opts)
}

// CarrierIdLookupParams are the parameters for the CarrierIdLookup method.
type CarrierIdLookupParams struct {
	PhoneNumber string `json:"phoneNumber"`
}

// CarrierIdLookup looks up the carrier for a phone number. Always performs a fresh lookup (not cached). Billed per request.
func (s *MessagesService) CarrierIdLookup(ctx context.Context, params CarrierIdLookupParams, opts *RequestOptions) (*Response, error) {
	if err := require(map[string]string{
		"phoneNumber": params.PhoneNumber,
	}); err != nil {
		return nil, err
	}
	return s.client.doJSON(ctx, "POST", "/message/carrier-id", params, opts)
}

// writeMMSForm writes MMS/GroupMessage form fields to a multipart writer.
func writeMMSForm(w *multipart.Writer, senderPhone string, recipientPhones []string,
	messageBody string, mediaURLs []string, statusCallbackURL string,
	enableShortlink, enableCompression, filterLandlinesAndInactiveNumbers bool, images []*FileUpload) error {

	// Add images
	for _, img := range images {
		if img != nil {
			part, err := w.CreateFormFile("images", img.Filename)
			if err != nil {
				return err
			}
			if _, err := io.Copy(part, img.Reader); err != nil {
				return err
			}
		}
	}

	// Add form fields
	if err := w.WriteField("senderPhoneNumber", senderPhone); err != nil {
		return err
	}

	recipientJSON, err := json.Marshal(recipientPhones)
	if err != nil {
		return err
	}
	if err := w.WriteField("recipientPhoneNumber", string(recipientJSON)); err != nil {
		return err
	}

	if err := w.WriteField("messageBody", messageBody); err != nil {
		return err
	}

	if len(mediaURLs) > 0 {
		urlsJSON, err := json.Marshal(mediaURLs)
		if err != nil {
			return err
		}
		if err := w.WriteField("mediaUrls", string(urlsJSON)); err != nil {
			return err
		}
	}

	if statusCallbackURL != "" {
		if err := w.WriteField("statusCallbackUrl", statusCallbackURL); err != nil {
			return err
		}
	}

	if err := w.WriteField("enableShortlink", fmt.Sprintf("%t", enableShortlink)); err != nil {
		return err
	}

	if err := w.WriteField("enableCompression", fmt.Sprintf("%t", enableCompression)); err != nil {
		return err
	}

	if filterLandlinesAndInactiveNumbers {
		if err := w.WriteField("filterLandlinesAndInactiveNumbers", "true"); err != nil {
			return err
		}
	}

	return nil
}
