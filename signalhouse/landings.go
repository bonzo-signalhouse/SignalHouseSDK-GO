package signalhouse

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/url"
)

// LandingsService handles landing page-related API calls.
type LandingsService struct {
	client *Client
}

// GetLandingsParams are the parameters for the GetLandings method.
type GetLandingsParams struct {
	LandingID string
}

// GetLandings returns details of a landing page by its ID.
func (s *LandingsService) GetLandings(ctx context.Context, params GetLandingsParams, opts *RequestOptions) (*Response, error) {
	if err := require(map[string]string{"landingId": params.LandingID}); err != nil {
		return nil, err
	}
	path := fmt.Sprintf("/landing/%s", url.PathEscape(params.LandingID))
	return s.client.doJSON(ctx, "GET", path, nil, opts)
}

// GetPublicLandingParams are the parameters for the GetPublicLanding method.
type GetPublicLandingParams struct {
	LandingID string
}

// GetPublicLanding returns the public view of a landing page by its ID. This endpoint requires no authentication.
func (s *LandingsService) GetPublicLanding(ctx context.Context, params GetPublicLandingParams, opts *RequestOptions) (*Response, error) {
	if err := require(map[string]string{"landingId": params.LandingID}); err != nil {
		return nil, err
	}
	path := fmt.Sprintf("/landing/public/%s", url.PathEscape(params.LandingID))
	return s.client.doJSON(ctx, "GET", path, nil, opts)
}

// GetLandingByBrandIDParams are the parameters for the GetLandingByBrandID method.
type GetLandingByBrandIDParams struct {
	BrandID string
}

// GetLandingByBrandID returns a landing page by its associated brand ID.
func (s *LandingsService) GetLandingByBrandID(ctx context.Context, params GetLandingByBrandIDParams, opts *RequestOptions) (*Response, error) {
	if err := require(map[string]string{"brandId": params.BrandID}); err != nil {
		return nil, err
	}
	path := fmt.Sprintf("/landing/brand/%s", url.PathEscape(params.BrandID))
	return s.client.doJSON(ctx, "GET", path, nil, opts)
}

// FileUpload represents a file to be uploaded in a multipart request.
type FileUpload struct {
	// Filename is the name of the file.
	Filename string
	// Reader provides the file content.
	Reader io.Reader
}

// TollFreeConsentBox is one consent checkbox rendered on a Toll-Free landing page. The page renders these
// values and the call-to-action submitted to the carrier repeats them, so the two cannot diverge.
type TollFreeConsentBox struct {
	UseCase string `json:"useCase"`
	Channel string `json:"channel"`
	Text    string `json:"text"`
}

// TollFreeLandingData is the Toll-Free specific content of a landing page. Required when
// RegistrationType is "TOLL_FREE", ignored otherwise.
type TollFreeLandingData struct {
	UseCases          []string             `json:"useCases"`
	Channels          []string             `json:"channels"`
	ConsentDisclosure string               `json:"consentDisclosure"`
	ConsentBoxes      []TollFreeConsentBox `json:"consentBoxes"`
}

// CreateLandingData represents the data for creating a landing page.
type CreateLandingData struct {
	BrandID                  string `json:"brandId"`
	Description              string `json:"description"`
	PrimaryBackgroundColor   string `json:"primaryBackgroundColor"`
	SecondaryBackgroundColor string `json:"secondaryBackgroundColor"`
	PrimaryTextColor         string `json:"primaryTextColor"`
	SecondaryTextColor       string `json:"secondaryTextColor"`
	// RegistrationType selects the page type: "TEN_DLC" (default) or "TOLL_FREE".
	RegistrationType string `json:"registrationType,omitempty"`
	// TollFree carries the consent content for a Toll-Free page; required when RegistrationType is "TOLL_FREE".
	TollFree *TollFreeLandingData `json:"tollFree,omitempty"`
	// UseBrandTemplate builds the page from the brand's landing page template, taking its colors, logo and
	// webhook URL from there. Toll-Free only; no logo file is required when set. Description is still sent by
	// the caller, since the About paragraph names the services this campaign's use cases cover.
	UseBrandTemplate bool `json:"useBrandTemplate,omitempty"`
}

// CreateLandingParams are the parameters for the CreateLanding method.
type CreateLandingParams struct {
	LandingData CreateLandingData
	File        *FileUpload
}

// CreateLanding creates a new landing page with optional logo file upload.
func (s *LandingsService) CreateLanding(ctx context.Context, params CreateLandingParams, opts *RequestOptions) (*Response, error) {
	writeForm := func(w *multipart.Writer) error {
		if params.File != nil {
			part, err := w.CreateFormFile("file", params.File.Filename)
			if err != nil {
				return err
			}
			if _, err := io.Copy(part, params.File.Reader); err != nil {
				return err
			}
		}
		return writeLandingFields(w, params.LandingData)
	}
	return s.client.doMultipart(ctx, "POST", "/landing", writeForm, opts)
}

// UpdateLandingData represents the data for updating a landing page.
type UpdateLandingData struct {
	Description              string `json:"description,omitempty"`
	PrimaryBackgroundColor   string `json:"primaryBackgroundColor,omitempty"`
	SecondaryBackgroundColor string `json:"secondaryBackgroundColor,omitempty"`
	PrimaryTextColor         string `json:"primaryTextColor,omitempty"`
	SecondaryTextColor       string `json:"secondaryTextColor,omitempty"`
	// RegistrationType selects the page type: "TEN_DLC" or "TOLL_FREE".
	RegistrationType string `json:"registrationType,omitempty"`
	// TollFree replaces the whole Toll-Free consent subdocument when set.
	TollFree *TollFreeLandingData `json:"tollFree,omitempty"`
	// UseBrandTemplate re-copies the brand's landing page template into this page rather than creating a
	// second one. Toll-Free only.
	UseBrandTemplate bool `json:"useBrandTemplate,omitempty"`
}

// UpdateLandingParams are the parameters for the UpdateLanding method.
type UpdateLandingParams struct {
	LandingID   string
	LandingData UpdateLandingData
	File        *FileUpload
}

// UpdateLanding updates an existing landing page with optional logo file upload.
func (s *LandingsService) UpdateLanding(ctx context.Context, params UpdateLandingParams, opts *RequestOptions) (*Response, error) {
	if err := require(map[string]string{"landingId": params.LandingID}); err != nil {
		return nil, err
	}
	path := fmt.Sprintf("/landing/%s", url.PathEscape(params.LandingID))
	writeForm := func(w *multipart.Writer) error {
		if params.File != nil {
			part, err := w.CreateFormFile("file", params.File.Filename)
			if err != nil {
				return err
			}
			if _, err := io.Copy(part, params.File.Reader); err != nil {
				return err
			}
		}
		return writeUpdateLandingFields(w, params.LandingData)
	}
	return s.client.doMultipart(ctx, "PUT", path, writeForm, opts)
}

// LandingTemplateData is the reusable look every page a Toll-Free brand generates snapshots.
type LandingTemplateData struct {
	// UserDescription is the customer's own sentence about their business. It is NOT the finished About
	// paragraph: the services clause is added per campaign when a page is generated.
	UserDescription          string `json:"userDescription"`
	PrimaryBackgroundColor   string `json:"primaryBackgroundColor,omitempty"`
	SecondaryBackgroundColor string `json:"secondaryBackgroundColor,omitempty"`
	PrimaryTextColor         string `json:"primaryTextColor,omitempty"`
	SecondaryTextColor       string `json:"secondaryTextColor,omitempty"`
	// WebhookURL is where generated pages post opt-in submissions.
	WebhookURL string `json:"webhookUrl,omitempty"`
}

// GetLandingTemplateParams are the parameters for the GetLandingTemplate method.
type GetLandingTemplateParams struct {
	BrandID string
}

// GetLandingTemplate returns a brand's Toll-Free landing page template, or an empty list when it has none.
func (s *LandingsService) GetLandingTemplate(ctx context.Context, params GetLandingTemplateParams, opts *RequestOptions) (*Response, error) {
	if err := require(map[string]string{"brandId": params.BrandID}); err != nil {
		return nil, err
	}
	path := fmt.Sprintf("/landing/template/%s", url.PathEscape(params.BrandID))
	return s.client.doJSON(ctx, "GET", path, nil, opts)
}

// UpsertLandingTemplateParams are the parameters for the UpsertLandingTemplate method.
type UpsertLandingTemplateParams struct {
	BrandID      string
	TemplateData LandingTemplateData
	// File is required the first time only; omit it on a later save to keep the stored logo.
	File *FileUpload
}

// UpsertLandingTemplate creates or replaces a brand's Toll-Free landing page template. Editing a template
// never alters a page that already exists, so a URL already given to a carrier keeps rendering what was
// declared.
func (s *LandingsService) UpsertLandingTemplate(ctx context.Context, params UpsertLandingTemplateParams, opts *RequestOptions) (*Response, error) {
	if err := require(map[string]string{"brandId": params.BrandID}); err != nil {
		return nil, err
	}
	path := fmt.Sprintf("/landing/template/%s", url.PathEscape(params.BrandID))
	writeForm := func(w *multipart.Writer) error {
		if params.File != nil {
			part, err := w.CreateFormFile("file", params.File.Filename)
			if err != nil {
				return err
			}
			if _, err := io.Copy(part, params.File.Reader); err != nil {
				return err
			}
		}
		return writeLandingTemplateFields(w, params.TemplateData)
	}
	return s.client.doMultipart(ctx, "PUT", path, writeForm, opts)
}

// writeLandingTemplateFields writes LandingTemplateData fields to a multipart writer.
func writeLandingTemplateFields(w *multipart.Writer, data LandingTemplateData) error {
	fields := map[string]string{
		"userDescription":          data.UserDescription,
		"primaryBackgroundColor":   data.PrimaryBackgroundColor,
		"secondaryBackgroundColor": data.SecondaryBackgroundColor,
		"primaryTextColor":         data.PrimaryTextColor,
		"secondaryTextColor":       data.SecondaryTextColor,
		"webhookUrl":               data.WebhookURL,
	}
	for k, v := range fields {
		if v != "" {
			if err := w.WriteField(k, v); err != nil {
				return err
			}
		}
	}
	return nil
}

// DeleteLandingParams are the parameters for the DeleteLanding method.
type DeleteLandingParams struct {
	LandingID string
}

// DeleteLanding deletes a landing page.
func (s *LandingsService) DeleteLanding(ctx context.Context, params DeleteLandingParams, opts *RequestOptions) (*Response, error) {
	if err := require(map[string]string{"landingId": params.LandingID}); err != nil {
		return nil, err
	}
	path := fmt.Sprintf("/landing/%s", url.PathEscape(params.LandingID))
	return s.client.doJSON(ctx, "DELETE", path, nil, opts)
}

// writeLandingFields writes CreateLandingData fields to a multipart writer.
func writeLandingFields(w *multipart.Writer, data CreateLandingData) error {
	fields := map[string]string{
		"brandId":                  data.BrandID,
		"description":              data.Description,
		"primaryBackgroundColor":   data.PrimaryBackgroundColor,
		"secondaryBackgroundColor": data.SecondaryBackgroundColor,
		"primaryTextColor":         data.PrimaryTextColor,
		"secondaryTextColor":       data.SecondaryTextColor,
		"registrationType":         data.RegistrationType,
	}
	for k, v := range fields {
		if v != "" {
			if err := w.WriteField(k, v); err != nil {
				return err
			}
		}
	}
	if data.TollFree != nil {
		if err := marshalToField(w, "tollFree", data.TollFree); err != nil {
			return err
		}
	}
	if data.UseBrandTemplate {
		if err := w.WriteField("useBrandTemplate", "true"); err != nil {
			return err
		}
	}
	return nil
}

// writeUpdateLandingFields writes UpdateLandingData fields to a multipart writer.
func writeUpdateLandingFields(w *multipart.Writer, data UpdateLandingData) error {
	fields := map[string]string{
		"description":              data.Description,
		"primaryBackgroundColor":   data.PrimaryBackgroundColor,
		"secondaryBackgroundColor": data.SecondaryBackgroundColor,
		"primaryTextColor":         data.PrimaryTextColor,
		"secondaryTextColor":       data.SecondaryTextColor,
		"registrationType":         data.RegistrationType,
	}
	for k, v := range fields {
		if v != "" {
			if err := w.WriteField(k, v); err != nil {
				return err
			}
		}
	}
	if data.TollFree != nil {
		if err := marshalToField(w, "tollFree", data.TollFree); err != nil {
			return err
		}
	}
	if data.UseBrandTemplate {
		if err := w.WriteField("useBrandTemplate", "true"); err != nil {
			return err
		}
	}
	return nil
}

// marshalToField is a helper to marshal a value to JSON and write it as a form field.
func marshalToField(w *multipart.Writer, key string, value interface{}) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return w.WriteField(key, string(data))
}
