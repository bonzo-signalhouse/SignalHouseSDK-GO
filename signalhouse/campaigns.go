package signalhouse

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"net/url"
)

// CampaignsService handles campaign-related API calls.
type CampaignsService struct {
	client      *Client
	enableAdmin bool

	// Admin contains admin-only campaign methods. Only populated when EnableAdmin is true.
	Admin *CampaignsAdminService
}

// CampaignsAdminService contains admin-only campaign methods.
type CampaignsAdminService struct {
	client *Client
}

// GetCampaignsParams are the parameters for the GetCampaigns method.
type GetCampaignsParams struct {
	ID         string
	BrandID    string
	SubgroupID string
	GroupID    string
	Page       string
	Limit      string
	Status     string
	// RegistrationType optionally filters by "TEN_DLC" or "TOLL_FREE".
	RegistrationType string
}

// GetCampaigns returns a list of campaigns with optional filters.
//
// Every record carries a `sentiment` object -- trailing 7-day and 30-day inbound reply sentiment,
// refreshed periodically from analytics rather than computed per request: `{ sevenDay, thirtyDay,
// updatedAt }`, each window being `{ score, label, scoredCount, positive, neutral, negative }`. The
// object is always present and fully populated, so it can be read without a presence check.
// `score` is a volume-weighted average in -100..100 and is null when nothing was scored in the
// window, which is not the same as a score of 0 -- null means nobody replied, 0 means replies
// averaged neutral. `updatedAt` is null until the first rollup writes the record, and it tracks when
// the figures last CHANGED, not when the job last ran.
func (s *CampaignsService) GetCampaigns(ctx context.Context, params GetCampaignsParams, opts *RequestOptions) (*Response, error) {
	qs := buildQueryString(map[string]string{
		"id":               params.ID,
		"brandId":          params.BrandID,
		"subgroupId":       params.SubgroupID,
		"groupId":          params.GroupID,
		"page":             params.Page,
		"limit":            params.Limit,
		"status":           params.Status,
		"registrationType": params.RegistrationType,
	})
	return s.client.doJSON(ctx, "GET", "/campaign"+qs, nil, opts)
}

// GetCampaignHealthParams are the parameters for the GetCampaignHealth method.
type GetCampaignHealthParams struct {
	CampaignID     string
	IncludeNumbers *bool
}

// GetCampaignHealth returns aggregated campaign health (7-day and 30-day windows) for a campaign.
func (s *CampaignsService) GetCampaignHealth(ctx context.Context, params GetCampaignHealthParams, opts *RequestOptions) (*Response, error) {
	if err := require(map[string]string{"campaignId": params.CampaignID}); err != nil {
		return nil, err
	}

	include := ""
	if params.IncludeNumbers != nil {
		if *params.IncludeNumbers {
			include = "true"
		} else {
			include = "false"
		}
	}

	qs := buildQueryString(map[string]string{
		"campaignId":     params.CampaignID,
		"includeNumbers": include,
	})
	return s.client.doJSON(ctx, "GET", "/campaign/health"+qs, nil, opts)
}

// CreateCampaignData represents the data for creating a campaign.
type CreateCampaignData struct {
	BrandID                string   `json:"brandId"`
	Usecase                string   `json:"usecase"`
	SubUsecases            []string `json:"subUsecases,omitempty"`
	Description            string   `json:"description"`
	EmbeddedLink           *bool    `json:"embeddedLink,omitempty"`
	EmbeddedPhone          *bool    `json:"embeddedPhone,omitempty"`
	TermsAndConditions     *bool    `json:"termsAndConditions,omitempty"`
	NumberPool             *bool    `json:"numberPool,omitempty"`
	AgeGated               *bool    `json:"ageGated,omitempty"`
	DirectLending          *bool    `json:"directLending,omitempty"`
	SubscriberOptIn        *bool    `json:"subscriberOptIn,omitempty"`
	SubscriberOptOut       *bool    `json:"subscriberOptOut,omitempty"`
	SubscriberHelp         *bool    `json:"subscriberHelp,omitempty"`
	Sample1                string   `json:"sample1,omitempty"`
	Sample2                string   `json:"sample2,omitempty"`
	Sample3                string   `json:"sample3,omitempty"`
	Sample4                string   `json:"sample4,omitempty"`
	Sample5                string   `json:"sample5,omitempty"`
	MessageFlow            string   `json:"messageFlow"`
	HelpMessage            string   `json:"helpMessage,omitempty"`
	AutoRenewal            *bool    `json:"autoRenewal,omitempty"`
	Tag                    string   `json:"tag,omitempty"`
	OptinKeywords          string   `json:"optinKeywords,omitempty"`
	OptoutKeywords         string   `json:"optoutKeywords,omitempty"`
	HelpKeywords           string   `json:"helpKeywords,omitempty"`
	OptinMessage           string   `json:"optinMessage,omitempty"`
	OptoutMessage          string   `json:"optoutMessage,omitempty"`
	PrivacyPolicyLink      string   `json:"privacyPolicyLink"`
	TermsAndConditionsLink string   `json:"termsAndConditionsLink"`
	EmbeddedLinkSample     string   `json:"embeddedLinkSample,omitempty"`
	PhoneNumbers           []string `json:"phoneNumbers"`
}

// CreateCampaignParams are the parameters for the CreateCampaign method.
type CreateCampaignParams struct {
	CampaignData CreateCampaignData
}

// CreateCampaign creates a new campaign.
func (s *CampaignsService) CreateCampaign(ctx context.Context, params CreateCampaignParams, opts *RequestOptions) (*Response, error) {
	return s.client.doJSON(ctx, "POST", "/campaign", params.CampaignData, opts)
}

// TollFreeOptIn describes a single Toll-Free opt-in channel.
type TollFreeOptIn struct {
	// CallToAction is the opt-in call-to-action text (20-1024 characters when provided).
	CallToAction string `json:"callToAction,omitempty"`
	// URL is the opt-in web URL (web channel only).
	URL string `json:"url,omitempty"`
	// Keywords are the opt-in keywords (keyword channel only).
	Keywords []string `json:"keywords,omitempty"`
}

// TollFreeOptIns groups the per-channel opt-in details. Each channel is optional.
type TollFreeOptIns struct {
	Verbal                   *TollFreeOptIn `json:"verbal,omitempty"`
	Web                      *TollFreeOptIn `json:"web,omitempty"`
	Keyword                  *TollFreeOptIn `json:"keyword,omitempty"`
	InteractiveVoiceResponse *TollFreeOptIn `json:"interactiveVoiceResponse,omitempty"`
}

// TollFreeCampaignDetails is the Toll-Free campaign sub-object (campaignData.tollFree). On create,
// UseCase, MessageVolume, ProgramSummary, ExampleMessage, CustomerCareEmail, and OptInImageURLs
// (at least one) are required; OptIns and MultiNumberReason are optional. When reused on
// UpdateCampaignData every field is optional (only populated fields are sent).
type TollFreeCampaignDetails struct {
	// UseCase is the Toll-Free use case (e.g. TWO_FA, GENERAL_MARKETING, HEALTHCARE, MIXED).
	UseCase string `json:"useCase,omitempty"`
	// MessageVolume is the monthly volume tier: one of TEN, HUNDRED, THOUSAND, TEN_THOUSAND,
	// HUNDRED_THOUSAND, TWO_HUNDRED_FIFTY_THOUSAND, FIVE_HUNDRED_THOUSAND, SEVEN_HUNDRED_FIFTY_THOUSAND,
	// ONE_MILLION, FIVE_MILLION, TEN_MILLION_PLUS.
	MessageVolume string `json:"messageVolume,omitempty"`
	// ProgramSummary is a 40-500 character summary of the messaging program.
	ProgramSummary string `json:"programSummary,omitempty"`
	// ExampleMessage is stored verbatim; the server appends no STOP/HELP text. Holds one sample per use
	// case and is capped at 1000 characters, matching the carrier limit.
	ExampleMessage    string   `json:"exampleMessage,omitempty"`
	CustomerCareEmail string   `json:"customerCareEmail,omitempty"`
	OptInImageURLs    []string `json:"optInImageURLs,omitempty"`
	OptIns            *TollFreeOptIns `json:"optIns,omitempty"`
	// MultiNumberReason is required when more than one number is assigned; internal to Signal House.
	MultiNumberReason string `json:"multiNumberReason,omitempty"`
	// UseCases holds the sender's real use-case picks. UseCase above carries the single value the carrier
	// accepts (MIXED when there is more than one); these are what a Signal House landing page renders a
	// consent checkbox per. Never forwarded to the carrier.
	UseCases []string `json:"useCases,omitempty"`
	// Channels are the channels consent covers: SMS, MMS, VOICE. Never forwarded to the carrier.
	Channels []string `json:"channels,omitempty"`
	// LandingID is the Signal House landing page whose consent boxes this campaign's call-to-action
	// describes. Must be an active page owned by the same brand and group.
	LandingID string `json:"landingId,omitempty"`
}

// CreateTollFreeCampaignData represents the data for creating a Toll-Free (TFN) campaign.
// registrationType is forced to TOLL_FREE server-side — do not set it here.
type CreateTollFreeCampaignData struct {
	BrandID                string                  `json:"brandId"`
	TollFree               TollFreeCampaignDetails `json:"tollFree"`
	PrivacyPolicyLink      string                  `json:"privacyPolicyLink,omitempty"`
	TermsAndConditionsLink string                  `json:"termsAndConditionsLink,omitempty"`
	Tag                    string                  `json:"tag,omitempty"`
	AutoRenewal            *bool                   `json:"autoRenewal,omitempty"`
	// PhoneNumbers must list 1-5 Toll-Free numbers, locked to the campaign once assigned.
	PhoneNumbers []string `json:"phoneNumbers"`
}

// CreateTollFreeCampaignParams are the parameters for the CreateTollFreeCampaign method.
type CreateTollFreeCampaignParams struct {
	CampaignData CreateTollFreeCampaignData
}

// CreateTollFreeCampaign creates a new Toll-Free (TFN) campaign and submits it for Signal House review.
func (s *CampaignsService) CreateTollFreeCampaign(ctx context.Context, params CreateTollFreeCampaignParams, opts *RequestOptions) (*Response, error) {
	return s.client.doJSON(ctx, "POST", "/campaign/toll-free", params.CampaignData, opts)
}

// ShortCodeCampaignDetails is the Short Code campaign sub-object (campaignData.shortCode).
type ShortCodeCampaignDetails struct {
	// UseCases are the selected use cases (at least 1; ALERTS, CHAT, DONATION, FRAUD_ALERTS,
	// PROMOTIONAL_MARKETING, PUBLIC_SERVICE_ANNOUNCEMENT, SWEEPSTAKES_CONTEST,
	// TWO_FACTOR_AUTHENTICATION, VOTING_POLLING, POLITICAL, EMERGENCY_NOTIFICATIONS,
	// DELIVERY_NOTIFICATIONS, SECURITY_ALERTS).
	UseCases []string `json:"useCases"`
	// OptInMethods are the selected opt-in methods (at least 1; IVR_POS, WEB, TV, PRINT, KEYWORD,
	// VERBAL, RADIO).
	OptInMethods []string `json:"optInMethods"`
	// OptInMethodDescriptions holds a generated-but-editable description per selected opt-in
	// method, keyed by the method value (required for every selected method).
	OptInMethodDescriptions map[string]string `json:"optInMethodDescriptions"`
	// MessageFrequency is "RECURRING" or "SINGLE".
	MessageFrequency string `json:"messageFrequency"`
	// PricingTier is "STANDARD" or "FTEU".
	PricingTier string `json:"pricingTier"`
	// AdultContent flags SHAFT content; when true, DoubleOptInMessage is required.
	AdultContent bool `json:"adultContent,omitempty"`
	// DoubleOptInMessage is required when AdultContent is true (20-1024 characters).
	DoubleOptInMessage string `json:"doubleOptInMessage,omitempty"`
	// ProgramSummary is a generated-but-editable campaign/program summary (40-2048 characters).
	ProgramSummary string `json:"programSummary"`
	// OptInConfirmationMessage is the opt-in confirmation message (20-1024 characters).
	OptInConfirmationMessage string `json:"optInConfirmationMessage"`
	// OptInURL is the optional public HTTPS URL for the live opt-in experience. Signal House uses the
	// brand's OptInLink when it is empty, then captures its screenshot asynchronously, retrying up to three times.
	OptInURL string `json:"optInUrl"`
}

// CreateShortCodeCampaignData represents the data for creating a Short Code campaign.
// registrationType is forced to SHORT_CODE server-side — do not set it here.
type CreateShortCodeCampaignData struct {
	BrandID                string                   `json:"brandId"`
	PrivacyPolicyLink      string                   `json:"privacyPolicyLink"`
	TermsAndConditionsLink string                   `json:"termsAndConditionsLink"`
	OptinMessage           string                   `json:"optinMessage"`
	OptoutMessage          string                   `json:"optoutMessage"`
	HelpMessage            string                   `json:"helpMessage"`
	Sample1                string                   `json:"sample1"`
	Sample2                string                   `json:"sample2"`
	Sample3                string                   `json:"sample3"`
	AutoRenewal            *bool                    `json:"autoRenewal,omitempty"`
	Tag                    string                   `json:"tag,omitempty"`
	ShortCode              ShortCodeCampaignDetails `json:"shortCode"`
}

// CreateShortCodeCampaignParams are the parameters for the CreateShortCodeCampaign method.
type CreateShortCodeCampaignParams struct {
	CampaignData CreateShortCodeCampaignData
}

// CreateShortCodeCampaign creates a new Short Code campaign and submits it for Signal House review
// (SHGHL-2225). Requires an approved (VERIFIED) Short Code brand. Request or register the campaign's
// Short Code separately through NumbersService.RequestShortCodeAcquisition after creation.
func (s *CampaignsService) CreateShortCodeCampaign(ctx context.Context, params CreateShortCodeCampaignParams, opts *RequestOptions) (*Response, error) {
	writeForm := func(w *multipart.Writer) error {
		return marshalToField(w, "campaignData", params.CampaignData)
	}
	return s.client.doMultipart(ctx, "POST", "/campaign/short-code", writeForm, opts)
}

// UpdateShortCodeCampaignData represents the editable fields for a Short Code campaign update.
// Only permitted while the campaign is in Signal House Review or Signal House Rejected status; the
// number source cannot be changed once submitted.
type UpdateShortCodeCampaignData struct {
	PrivacyPolicyLink      string `json:"privacyPolicyLink,omitempty"`
	TermsAndConditionsLink string `json:"termsAndConditionsLink,omitempty"`
	OptinMessage           string `json:"optinMessage,omitempty"`
	OptoutMessage          string `json:"optoutMessage,omitempty"`
	HelpMessage            string `json:"helpMessage,omitempty"`
	Sample1                string `json:"sample1,omitempty"`
	Sample2                string `json:"sample2,omitempty"`
	Sample3                string `json:"sample3,omitempty"`
	Tag                    string `json:"tag,omitempty"`
	// ShortCode carries the editable Short Code fields (all optional; useCases/optInMethods/etc.).
	ShortCode *ShortCodeCampaignDetails `json:"shortCode,omitempty"`
}

// UpdateShortCodeCampaignParams are the parameters for the UpdateShortCodeCampaign method.
type UpdateShortCodeCampaignParams struct {
	CampaignID string
	UpdateData UpdateShortCodeCampaignData
}

// UpdateShortCodeCampaign updates a Short Code campaign's editable fields (SHGHL-2225).
func (s *CampaignsService) UpdateShortCodeCampaign(ctx context.Context, params UpdateShortCodeCampaignParams, opts *RequestOptions) (*Response, error) {
	if err := require(map[string]string{"campaignId": params.CampaignID}); err != nil {
		return nil, err
	}
	path := fmt.Sprintf("/campaign/short-code/%s", url.PathEscape(params.CampaignID))
	return s.client.doJSON(ctx, "PUT", path, params.UpdateData, opts)
}

// CancelShortCodeCampaignParams are the parameters for the CancelShortCodeCampaign method.
type CancelShortCodeCampaignParams struct {
	CampaignID string
}

// CancelShortCodeCampaign cancels a Short Code campaign (SHGHL-2225). Customers may cancel only
// while the campaign is in Signal House Review or Signal House Rejected status; Signal House staff
// may cancel from any non-terminal status. Persists as "EXPIRED" (displayed as "Cancelled"). A real
// Registry lease (external or an already-fulfilled Signal House request) is never auto-released —
// see SHGHL-2228 for offboarding.
func (s *CampaignsService) CancelShortCodeCampaign(ctx context.Context, params CancelShortCodeCampaignParams, opts *RequestOptions) (*Response, error) {
	if err := require(map[string]string{"campaignId": params.CampaignID}); err != nil {
		return nil, err
	}
	path := fmt.Sprintf("/campaign/short-code/%s/cancel", url.PathEscape(params.CampaignID))
	return s.client.doJSON(ctx, "POST", path, nil, opts)
}

// ReadCampaignArtifactParams are the parameters for the ReadCampaignArtifact method.
type ReadCampaignArtifactParams struct {
	ArtifactID string
}

// ReadCampaignArtifact downloads a private Short Code external lease receipt. Opt-in screenshots are
// public at the shortCode.screenshotUrl returned on the campaign.
// The artifact bytes are raw file content, not JSON — this client's standard JSON envelope does not
// apply; callers needing the raw bytes should issue their own authenticated HTTP GET to
// "/campaign/short-code/artifact/{artifactId}" instead.
func (s *CampaignsService) ReadCampaignArtifact(ctx context.Context, params ReadCampaignArtifactParams, opts *RequestOptions) (*Response, error) {
	if err := require(map[string]string{"artifactId": params.ArtifactID}); err != nil {
		return nil, err
	}
	path := fmt.Sprintf("/campaign/short-code/artifact/%s", url.PathEscape(params.ArtifactID))
	return s.client.doJSON(ctx, "GET", path, nil, opts)
}

// UploadOptInImageParams are the parameters for the UploadOptInImage method.
type UploadOptInImageParams struct {
	File *FileUpload
}

// UploadOptInImage uploads an opt-in proof image (multipart) and returns the hosted image's id and
// url, suitable for a Toll-Free campaign's optInImageURLs.
func (s *CampaignsService) UploadOptInImage(ctx context.Context, params UploadOptInImageParams, opts *RequestOptions) (*Response, error) {
	if params.File == nil {
		return nil, &ValidationError{Field: "file"}
	}
	writeForm := func(w *multipart.Writer) error {
		part, err := w.CreateFormFile("file", params.File.Filename)
		if err != nil {
			return err
		}
		if _, err := io.Copy(part, params.File.Reader); err != nil {
			return err
		}
		return nil
	}
	return s.client.doMultipart(ctx, "POST", "/campaign/opt-in-image", writeForm, opts)
}

// CaptureOptInImageFromLandingParams are the parameters for the CaptureOptInImageFromLanding method.
type CaptureOptInImageFromLandingParams struct {
	BrandID string
	// LandingID is the specific page to capture. Required for Toll-Free brands, which have one page per campaign — omitted, the request is refused rather than guessing. A 10DLC brand has a single page, so it may be omitted there.
	LandingID string
}

// CaptureOptInImageFromLanding auto-captures an opt-in proof image from a brand's generated landing
// page and returns the hosted image's id and url.
func (s *CampaignsService) CaptureOptInImageFromLanding(ctx context.Context, params CaptureOptInImageFromLandingParams, opts *RequestOptions) (*Response, error) {
	if err := require(map[string]string{"brandId": params.BrandID}); err != nil {
		return nil, err
	}
	body := map[string]string{"brandId": params.BrandID}
	if params.LandingID != "" {
		body["landingId"] = params.LandingID
	}
	return s.client.doJSON(ctx, "POST", "/campaign/opt-in-image/from-landing", body, opts)
}

// UpdateCampaignData represents the data for updating a campaign.
type UpdateCampaignData struct {
	Usecase                string   `json:"usecase,omitempty"`
	SubUsecases            []string `json:"subUsecases,omitempty"`
	Description            string   `json:"description,omitempty"`
	EmbeddedLink           *bool    `json:"embeddedLink,omitempty"`
	EmbeddedPhone          *bool    `json:"embeddedPhone,omitempty"`
	TermsAndConditions     *bool    `json:"termsAndConditions,omitempty"`
	NumberPool             *bool    `json:"numberPool,omitempty"`
	AgeGated               *bool    `json:"ageGated,omitempty"`
	DirectLending          *bool    `json:"directLending,omitempty"`
	SubscriberOptIn        *bool    `json:"subscriberOptIn,omitempty"`
	SubscriberOptOut       *bool    `json:"subscriberOptOut,omitempty"`
	SubscriberHelp         *bool    `json:"subscriberHelp,omitempty"`
	Sample1                string   `json:"sample1,omitempty"`
	Sample2                string   `json:"sample2,omitempty"`
	Sample3                string   `json:"sample3,omitempty"`
	Sample4                string   `json:"sample4,omitempty"`
	Sample5                string   `json:"sample5,omitempty"`
	MessageFlow            string   `json:"messageFlow,omitempty"`
	HelpMessage            string   `json:"helpMessage,omitempty"`
	AutoRenewal            *bool    `json:"autoRenewal,omitempty"`
	Tag                    string   `json:"tag,omitempty"`
	OptinKeywords          string   `json:"optinKeywords,omitempty"`
	OptoutKeywords         string   `json:"optoutKeywords,omitempty"`
	HelpKeywords           string   `json:"helpKeywords,omitempty"`
	OptinMessage           string   `json:"optinMessage,omitempty"`
	OptoutMessage          string   `json:"optoutMessage,omitempty"`
	PrivacyPolicyLink      string   `json:"privacyPolicyLink,omitempty"`
	TermsAndConditionsLink string   `json:"termsAndConditionsLink,omitempty"`
	EmbeddedLinkSample     string   `json:"embeddedLinkSample,omitempty"`
	// TollFree carries the editable Toll-Free fields when updating a Toll-Free campaign (all
	// optional; phoneNumbers cannot be changed — Toll-Free numbers are locked to their campaign).
	TollFree *TollFreeCampaignDetails `json:"tollFree,omitempty"`
}

// UpdateCampaignParams are the parameters for the UpdateCampaign method.
type UpdateCampaignParams struct {
	CampaignID   string
	CampaignData UpdateCampaignData
}

// UpdateCampaign updates an existing campaign.
func (s *CampaignsService) UpdateCampaign(ctx context.Context, params UpdateCampaignParams, opts *RequestOptions) (*Response, error) {
	if err := require(map[string]string{"campaignId": params.CampaignID}); err != nil {
		return nil, err
	}
	path := fmt.Sprintf("/campaign/%s", url.PathEscape(params.CampaignID))
	return s.client.doJSON(ctx, "PUT", path, params.CampaignData, opts)
}

// DeleteCampaignParams are the parameters for the DeleteCampaign method.
type DeleteCampaignParams struct {
	CampaignID string
}

// DeleteCampaign deletes (marks as EXPIRED) a campaign.
func (s *CampaignsService) DeleteCampaign(ctx context.Context, params DeleteCampaignParams, opts *RequestOptions) (*Response, error) {
	if err := require(map[string]string{"campaignId": params.CampaignID}); err != nil {
		return nil, err
	}
	path := fmt.Sprintf("/campaign/%s", url.PathEscape(params.CampaignID))
	return s.client.doJSON(ctx, "DELETE", path, nil, opts)
}

// AppealDcaRejectionParams are the parameters for the AppealDcaRejection method.
type AppealDcaRejectionParams struct {
	CampaignID string
}

// AppealDcaRejection appeals a DCA-rejected campaign.
func (s *CampaignsService) AppealDcaRejection(ctx context.Context, params AppealDcaRejectionParams, opts *RequestOptions) (*Response, error) {
	if err := require(map[string]string{"campaignId": params.CampaignID}); err != nil {
		return nil, err
	}
	path := fmt.Sprintf("/campaign/appealDcaRejection/%s", url.PathEscape(params.CampaignID))
	return s.client.doJSON(ctx, "POST", path, nil, opts)
}

// NudgeDcaForCampaignParams are the parameters for the NudgeDcaForCampaign method.
type NudgeDcaForCampaignParams struct {
	CampaignID string
}

// NudgeDcaForCampaign nudges a connectivity partner to prioritize review of a campaign in PENDING_DCA_APPROVAL status.
func (s *CampaignsService) NudgeDcaForCampaign(ctx context.Context, params NudgeDcaForCampaignParams, opts *RequestOptions) (*Response, error) {
	if err := require(map[string]string{"campaignId": params.CampaignID}); err != nil {
		return nil, err
	}
	path := fmt.Sprintf("/campaign/nudge/%s", url.PathEscape(params.CampaignID))
	return s.client.doJSON(ctx, "POST", path, nil, opts)
}

// ApproveCampaignParams are the parameters for the ApproveCampaign admin method.
type ApproveCampaignParams struct {
	CampaignID string
}

// ApproveCampaign approves a campaign that is pending approval. Admin only.
func (s *CampaignsAdminService) ApproveCampaign(ctx context.Context, params ApproveCampaignParams, opts *RequestOptions) (*Response, error) {
	if err := require(map[string]string{"campaignId": params.CampaignID}); err != nil {
		return nil, err
	}
	path := fmt.Sprintf("/campaign/approve/%s", url.PathEscape(params.CampaignID))
	return s.client.doJSON(ctx, "POST", path, nil, opts)
}

// RejectCampaignParams are the parameters for the RejectCampaign admin method.
type RejectCampaignParams struct {
	CampaignID      string
	RejectionReason string // Optional rejection reason (10-256 characters)
}

// RejectCampaign rejects a campaign that is pending approval. Admin only.
func (s *CampaignsAdminService) RejectCampaign(ctx context.Context, params RejectCampaignParams, opts *RequestOptions) (*Response, error) {
	if err := require(map[string]string{"campaignId": params.CampaignID}); err != nil {
		return nil, err
	}
	path := fmt.Sprintf("/campaign/reject/%s", url.PathEscape(params.CampaignID))
	var body interface{}
	if params.RejectionReason != "" {
		body = map[string]string{"rejectionReason": params.RejectionReason}
	}
	return s.client.doJSON(ctx, "POST", path, body, opts)
}

// UpdateShortCodeCampaignStatusParams are the parameters for the UpdateShortCodeCampaignStatus admin method.
type UpdateShortCodeCampaignStatusParams struct {
	CampaignID string
	// Status is the target status: "PENDING_REVIEW", "REJECTED", "PENDING_CREATION",
	// "DCA_REJECTED", or "ACTIVE". Use SubmitShortCodeCampaignToCarrier for PENDING_DCA_APPROVAL.
	Status string
	// RejectionReason is required for "REJECTED"/"DCA_REJECTED" (10-1024 characters).
	RejectionReason string
}

// UpdateShortCodeCampaignStatus transitions a Short Code campaign's review status (SHGHL-2225).
// Admin only.
func (s *CampaignsAdminService) UpdateShortCodeCampaignStatus(ctx context.Context, params UpdateShortCodeCampaignStatusParams, opts *RequestOptions) (*Response, error) {
	if err := require(map[string]string{"campaignId": params.CampaignID, "status": params.Status}); err != nil {
		return nil, err
	}
	path := fmt.Sprintf("/campaign/short-code/%s/status", url.PathEscape(params.CampaignID))
	body := map[string]interface{}{
		"status":          params.Status,
		"rejectionReason": params.RejectionReason,
	}
	return s.client.doJSON(ctx, "PUT", path, body, opts)
}

// FulfillShortCodeCampaignParams are the parameters for the FulfillShortCodeCampaign admin method.
type FulfillShortCodeCampaignParams struct { CampaignID string; ActualCode string `json:"actualCode"`; InternalNotes string `json:"internalNotes,omitempty"` }

// FulfillShortCodeCampaign records the issued Short Code for a campaign-bound Signal House request. Admin only.
func (s *CampaignsAdminService) FulfillShortCodeCampaign(ctx context.Context, params FulfillShortCodeCampaignParams, opts *RequestOptions) (*Response, error) {
	if err := require(map[string]string{"campaignId": params.CampaignID, "actualCode": params.ActualCode}); err != nil { return nil, err }
	path := fmt.Sprintf("/campaign/short-code/%s/fulfill", url.PathEscape(params.CampaignID))
	return s.client.doJSON(ctx, "POST", path, map[string]string{"actualCode": params.ActualCode, "internalNotes": params.InternalNotes}, opts)
}

// SubmitShortCodeCampaignToCarrierParams are the parameters for the SubmitShortCodeCampaignToCarrier admin method.
type SubmitShortCodeCampaignToCarrierParams struct { CampaignID string }

// SubmitShortCodeCampaignToCarrier submits an internally approved Short Code campaign to carrier review. Admin only.
func (s *CampaignsAdminService) SubmitShortCodeCampaignToCarrier(ctx context.Context, params SubmitShortCodeCampaignToCarrierParams, opts *RequestOptions) (*Response, error) {
	if err := require(map[string]string{"campaignId": params.CampaignID}); err != nil { return nil, err }
	path := fmt.Sprintf("/campaign/short-code/%s/submit-to-carrier", url.PathEscape(params.CampaignID))
	return s.client.doJSON(ctx, "POST", path, nil, opts)
}
