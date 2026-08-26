package signalhouse

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/url"
)

// BrandsService handles brand-related API calls.
//
// Brand lookup id: carrier Brand ID (B… or TFNB…), Mongo _id, or internal reference UUID.
// Read/create responses return brandId: null for pending brands — poll GET /brand?id=<Mongo _id>.
type BrandsService struct {
	client *Client
}

// GetBrandsParams are the parameters for the GetBrands method.
type GetBrandsParams struct {
	ID         string
	SubgroupID string
	GroupID    string
	Page       string
	Limit      string
	Status     string
	// RegistrationType optionally filters by "TEN_DLC", "TOLL_FREE", or "SHORT_CODE".
	RegistrationType string
}

// GetBrands returns a list of brands with optional filters. brandId may be null until carrier assignment.
func (s *BrandsService) GetBrands(ctx context.Context, params GetBrandsParams, opts *RequestOptions) (*Response, error) {
	qs := buildQueryString(map[string]string{
		"id":               params.ID,
		"subgroupId":       params.SubgroupID,
		"groupId":          params.GroupID,
		"page":             params.Page,
		"limit":            params.Limit,
		"status":           params.Status,
		"registrationType": params.RegistrationType,
	})
	return s.client.doJSON(ctx, "GET", "/brand"+qs, nil, opts)
}

// GetExternalVettingParams are the parameters for the GetExternalVetting method.
// BrandID is a brand lookup id (carrier id, Mongo _id, or internal reference).
type GetExternalVettingParams struct {
	BrandID string
}

// GetExternalVetting returns external vetting information for a brand.
func (s *BrandsService) GetExternalVetting(ctx context.Context, params GetExternalVettingParams, opts *RequestOptions) (*Response, error) {
	if err := require(map[string]string{"brandId": params.BrandID}); err != nil {
		return nil, err
	}
	path := fmt.Sprintf("/brand/externalvetting/%s", url.PathEscape(params.BrandID))
	return s.client.doJSON(ctx, "GET", path, nil, opts)
}

// CreateBrandData represents the data for creating a brand.
type CreateBrandData struct {
	SubgroupID           string   `json:"subgroupId"`
	EntityType           string   `json:"entityType"`
	DisplayName          string   `json:"displayName"`
	CompanyName          string   `json:"companyName"`
	EIN                  string   `json:"ein"`
	EINIssuingCountry    string   `json:"einIssuingCountry,omitempty"`
	AltBusinessID        string   `json:"altBusinessId,omitempty"`
	AltBusinessIDType    string   `json:"altBusinessIdType,omitempty"`
	FirstName            string   `json:"firstName,omitempty"`
	LastName             string   `json:"lastName,omitempty"`
	Phone                string   `json:"phone"`
	Street               string   `json:"street"`
	City                 string   `json:"city"`
	State                string   `json:"state"`
	PostalCode           string   `json:"postalCode"`
	Country              string   `json:"country"`
	Email                string   `json:"email"`
	MobilePhone          string   `json:"mobilePhone,omitempty"`
	StockSymbol          string   `json:"stockSymbol,omitempty"`
	StockExchange        string   `json:"stockExchange,omitempty"`
	Website              string   `json:"website,omitempty"`
	Vertical             string   `json:"vertical"`
	ReferenceID          string   `json:"referenceId,omitempty"`
	Tag                  []string `json:"tag,omitempty"`
	Mock                 *bool    `json:"mock,omitempty"`
	BusinessContactEmail string   `json:"businessContactEmail,omitempty"`
}

// CreateBrandParams are the parameters for the CreateBrand method.
type CreateBrandParams struct {
	BrandData CreateBrandData
}

// CreateBrand creates a new brand. Response brandId may be null until carrier assignment.
func (s *BrandsService) CreateBrand(ctx context.Context, params CreateBrandParams, opts *RequestOptions) (*Response, error) {
	return s.client.doJSON(ctx, "POST", "/brand", params.BrandData, opts)
}

// TollFreeBrandDetails is the Toll-Free brand sub-object (brandData.tollFree). On create,
// BusinessRegistrationType, LegalEntityType, TaxID, CountryCode, and SupportPhone are required;
// TaxIDIssuingCountry, BusinessDBA, and BusinessWebsite are optional. When reused on UpdateBrandData every field is
// optional (only populated fields are sent). SupportPhone must be E.164 with an optional leading
// "+" and no spaces, parentheses or hyphens (e.g. "17025550100").
type TollFreeBrandDetails struct {
	// BusinessRegistrationType is one of EIN, CBN, NEQ, PROVINCIAL_NUMBER, CRN, VAT, ACN, ABN, BRN,
	// SIREN, SIRET, NZBN, UST_IDNR, CIF, NIF, CNPJ, UID, OTHER.
	BusinessRegistrationType string `json:"businessRegistrationType,omitempty"`
	// LegalEntityType is one of PRIVATE_COMPANY, PUBLIC_COMPANY, NON_PROFIT, GOVERNMENT, SOLE_PROPRIETOR.
	LegalEntityType     string `json:"legalEntityType,omitempty"`
	TaxID               string `json:"taxId,omitempty"`
	CountryCode         string `json:"countryCode,omitempty"`
	SupportPhone        string `json:"supportPhone,omitempty"`
	TaxIDIssuingCountry string `json:"taxIdIssuingCountry,omitempty"`
	BusinessDBA         string `json:"businessDBA,omitempty"`
	BusinessWebsite     string `json:"businessWebsite,omitempty"`
}

// CreateTollFreeBrandData represents the data for creating a Toll-Free (TFN) brand.
// registrationType is forced to TOLL_FREE server-side — do not set it here.
type CreateTollFreeBrandData struct {
	SubgroupID             string               `json:"subgroupId"`
	DisplayName            string               `json:"displayName"`
	CompanyName            string               `json:"companyName"`
	Street                 string               `json:"street"`
	City                   string               `json:"city"`
	State                  string               `json:"state"`
	PostalCode             string               `json:"postalCode"`
	Email                  string               `json:"email"`
	FirstName              string               `json:"firstName"`
	LastName               string               `json:"lastName"`
	Website                string               `json:"website"`
	TollFree               TollFreeBrandDetails `json:"tollFree"`
	ReferenceID            string               `json:"referenceId,omitempty"`
	Tag                    []string             `json:"tag,omitempty"`
	Mock                   *bool                `json:"mock,omitempty"`
	OptInLink              string               `json:"optInLink,omitempty"`
	PrivacyPolicyLink      string               `json:"privacyPolicyLink,omitempty"`
	TermsAndConditionsLink string               `json:"termsAndConditionsLink,omitempty"`
	LandingID              string               `json:"landingId,omitempty"`
	Country                string               `json:"country,omitempty"`
}

// CreateTollFreeBrandParams are the parameters for the CreateTollFreeBrand method.
type CreateTollFreeBrandParams struct {
	BrandData CreateTollFreeBrandData
}

// CreateTollFreeBrand creates a new Toll-Free (TFN) brand. For Toll-Free, brandId is a
// TFNB-prefixed id when assigned; otherwise null — poll by Mongo _id.
func (s *BrandsService) CreateTollFreeBrand(ctx context.Context, params CreateTollFreeBrandParams, opts *RequestOptions) (*Response, error) {
	return s.client.doJSON(ctx, "POST", "/brand/toll-free", params.BrandData, opts)
}

// ShortCodeBrandDetails is the Short Code brand sub-object (brandData.shortCode).
type ShortCodeBrandDetails struct {
	BusinessRegistrationType string `json:"businessRegistrationType"`
	LegalEntityType          string `json:"legalEntityType"`
	TaxID                    string `json:"taxId"`
	CountryCode              string `json:"countryCode"`
	SupportPhone             string `json:"supportPhone"`
}

// CreateShortCodeBrandData represents the data for creating a local Short Code brand.
// registrationType is forced to SHORT_CODE server-side.
type CreateShortCodeBrandData struct {
	SubgroupID             string                `json:"subgroupId"`
	DisplayName            string                `json:"displayName"`
	CompanyName            string                `json:"companyName"`
	Street                 string                `json:"street"`
	City                   string                `json:"city"`
	State                  string                `json:"state"`
	PostalCode             string                `json:"postalCode"`
	Country                string                `json:"country"`
	Email                  string                `json:"email"`
	FirstName              string                `json:"firstName"`
	LastName               string                `json:"lastName"`
	Website                string                `json:"website"`
	OptInLink              string                `json:"optInLink,omitempty"`
	PrivacyPolicyLink      string                `json:"privacyPolicyLink,omitempty"`
	TermsAndConditionsLink string                `json:"termsAndConditionsLink,omitempty"`
	ShortCode              ShortCodeBrandDetails `json:"shortCode"`
}

// CreateShortCodeBrandParams are the parameters for CreateShortCodeBrand.
type CreateShortCodeBrandParams struct {
	BrandData CreateShortCodeBrandData
}

// CreateShortCodeBrand creates a local Short Code brand in PENDING_APPROVAL for manual review.
// It returns a synchronous 201-created record with a stable SCB-prefixed brand ID.
func (s *BrandsService) CreateShortCodeBrand(ctx context.Context, params CreateShortCodeBrandParams, opts *RequestOptions) (*Response, error) {
	return s.client.doJSON(ctx, "POST", "/brand/short-code", params.BrandData, opts)
}

// ApproveShortCodeBrandParams are the parameters for ApproveShortCodeBrand.
type ApproveShortCodeBrandParams struct {
	BrandID string
}

// ApproveShortCodeBrand approves a pending Short Code brand. signalhouse_admin only.
func (s *BrandsService) ApproveShortCodeBrand(ctx context.Context, params ApproveShortCodeBrandParams, opts *RequestOptions) (*Response, error) {
	if err := require(map[string]string{"brandId": params.BrandID}); err != nil {
		return nil, err
	}
	return s.client.doJSON(ctx, "POST", fmt.Sprintf("/brand/approve/%s", url.PathEscape(params.BrandID)), nil, opts)
}

// RejectShortCodeBrandParams are the parameters for RejectShortCodeBrand.
type RejectShortCodeBrandParams struct {
	BrandID                 string
	InternalRejectionReason string
}

// RejectShortCodeBrand rejects a pending Short Code brand with a customer-visible reason. signalhouse_admin only.
func (s *BrandsService) RejectShortCodeBrand(ctx context.Context, params RejectShortCodeBrandParams, opts *RequestOptions) (*Response, error) {
	if err := require(map[string]string{"brandId": params.BrandID, "internalRejectionReason": params.InternalRejectionReason}); err != nil {
		return nil, err
	}
	return s.client.doJSON(ctx, "POST", fmt.Sprintf("/brand/reject/%s", url.PathEscape(params.BrandID)), map[string]string{"internalRejectionReason": params.InternalRejectionReason}, opts)
}

// TransferBrandParams are the parameters for the TransferBrand method.
// BrandIDs are brand lookup ids (carrier id, Mongo _id, or internal reference).
type TransferBrandParams struct {
	SubgroupID string   `json:"-"`
	BrandIDs   []string `json:"brandIds"`
}

// TransferBrand transfers brands to a different subgroup.
func (s *BrandsService) TransferBrand(ctx context.Context, params TransferBrandParams, opts *RequestOptions) (*Response, error) {
	if err := require(map[string]string{"subgroupId": params.SubgroupID}); err != nil {
		return nil, err
	}
	if len(params.BrandIDs) == 0 {
		return nil, &ValidationError{Field: "brandIds"}
	}
	path := fmt.Sprintf("/brand/transfer/%s", url.PathEscape(params.SubgroupID))
	body := map[string][]string{"brandIds": params.BrandIDs}
	return s.client.doJSON(ctx, "POST", path, body, opts)
}

// CreateExternalVettingParams are the parameters for the CreateExternalVetting method.
type CreateExternalVettingParams struct {
	BrandID string
}

// CreateExternalVetting creates external vetting for a brand.
func (s *BrandsService) CreateExternalVetting(ctx context.Context, params CreateExternalVettingParams, opts *RequestOptions) (*Response, error) {
	if err := require(map[string]string{"brandId": params.BrandID}); err != nil {
		return nil, err
	}
	path := fmt.Sprintf("/brand/externalvetting/%s", url.PathEscape(params.BrandID))
	return s.client.doJSON(ctx, "POST", path, nil, opts)
}

// ImportExternalVettingParams are the parameters for the ImportExternalVetting method.
type ImportExternalVettingParams struct {
	BrandID           string
	VettingProviderID string
	VettingID         string
	// VettingToken is required by some providers (e.g. AEGIS); leave empty to omit it.
	VettingToken string
}

// ImportExternalVetting imports an existing external vetting record for a brand.
//
// Unlike CreateExternalVetting (which orders a new, billable vetting), this attaches a vetting the
// brand already completed directly with the provider, using the provider-issued vettingId and
// vettingToken. It is synchronous and not billable.
func (s *BrandsService) ImportExternalVetting(ctx context.Context, params ImportExternalVettingParams, opts *RequestOptions) (*Response, error) {
	if err := require(map[string]string{"brandId": params.BrandID, "vettingProviderId": params.VettingProviderID, "vettingId": params.VettingID}); err != nil {
		return nil, err
	}
	path := fmt.Sprintf("/brand/externalvetting/import/%s", url.PathEscape(params.BrandID))
	body := map[string]string{
		"vettingProviderId": params.VettingProviderID,
		"vettingId":         params.VettingID,
	}
	if params.VettingToken != "" {
		body["vettingToken"] = params.VettingToken
	}
	return s.client.doJSON(ctx, "POST", path, body, opts)
}

// UpdateBrandData represents the data for updating a brand.
type UpdateBrandData struct {
	CompanyName          string `json:"companyName,omitempty"`
	EIN                  string `json:"ein,omitempty"`
	EINIssuingCountry    string `json:"einIssuingCountry,omitempty"`
	EntityType           string `json:"entityType,omitempty"`
	FirstName            string `json:"firstName,omitempty"`
	LastName             string `json:"lastName,omitempty"`
	DisplayName          string `json:"displayName,omitempty"`
	Website              string `json:"website,omitempty"`
	Street               string `json:"street,omitempty"`
	City                 string `json:"city,omitempty"`
	State                string `json:"state,omitempty"`
	PostalCode           string `json:"postalCode,omitempty"`
	Country              string `json:"country,omitempty"`
	Email                string `json:"email,omitempty"`
	ReferenceID          string `json:"referenceId,omitempty"`
	Phone                string `json:"phone,omitempty"`
	Vertical             string `json:"vertical,omitempty"`
	StockSymbol          string `json:"stockSymbol,omitempty"`
	StockExchange        string `json:"stockExchange,omitempty"`
	AltBusinessID        string `json:"altBusinessId,omitempty"`
	AltBusinessIDType    string `json:"altBusinessIdType,omitempty"`
	BrandRelationship    string `json:"brandRelationship,omitempty"`
	BusinessContactEmail string `json:"businessContactEmail,omitempty"`
	// TollFree carries the editable Toll-Free fields when updating a Toll-Free brand (all optional;
	// subgroupId is immutable).
	TollFree *TollFreeBrandDetails `json:"tollFree,omitempty"`
}

// UpdateBrandParams are the parameters for the UpdateBrand method.
type UpdateBrandParams struct {
	BrandID   string
	BrandData UpdateBrandData
}

// UpdateBrand updates a brand's information.
func (s *BrandsService) UpdateBrand(ctx context.Context, params UpdateBrandParams, opts *RequestOptions) (*Response, error) {
	if err := require(map[string]string{"brandId": params.BrandID}); err != nil {
		return nil, err
	}
	path := fmt.Sprintf("/brand/%s", url.PathEscape(params.BrandID))
	return s.client.doJSON(ctx, "PUT", path, params.BrandData, opts)
}

// RevetBrandParams are the parameters for the RevetBrand method.
type RevetBrandParams struct {
	BrandID string
}

// RevetBrand revets a brand that is in UNVERIFIED status.
func (s *BrandsService) RevetBrand(ctx context.Context, params RevetBrandParams, opts *RequestOptions) (*Response, error) {
	if err := require(map[string]string{"brandId": params.BrandID}); err != nil {
		return nil, err
	}
	path := fmt.Sprintf("/brand/revet/%s", url.PathEscape(params.BrandID))
	return s.client.doJSON(ctx, "PUT", path, nil, opts)
}

// DeleteBrandParams are the parameters for the DeleteBrand method.
type DeleteBrandParams struct {
	BrandID string
}

// DeleteBrand deletes (marks as DELETED) a brand.
func (s *BrandsService) DeleteBrand(ctx context.Context, params DeleteBrandParams, opts *RequestOptions) (*Response, error) {
	if err := require(map[string]string{"brandId": params.BrandID}); err != nil {
		return nil, err
	}
	path := fmt.Sprintf("/brand/%s", url.PathEscape(params.BrandID))
	return s.client.doJSON(ctx, "DELETE", path, nil, opts)
}

// GetAppealHistoryParams are the parameters for the GetAppealHistory method.
type GetAppealHistoryParams struct {
	BrandID string
}

// GetAppealHistory returns the appeal history for a brand.
func (s *BrandsService) GetAppealHistory(ctx context.Context, params GetAppealHistoryParams, opts *RequestOptions) (*Response, error) {
	if err := require(map[string]string{"brandId": params.BrandID}); err != nil {
		return nil, err
	}
	path := fmt.Sprintf("/brand/appeal/%s", url.PathEscape(params.BrandID))
	return s.client.doJSON(ctx, "GET", path, nil, opts)
}

// SubmitAppealParams are the parameters for the SubmitAppeal method.
type SubmitAppealParams struct {
	BrandID          string        `json:"-"`
	AppealCategories []string      `json:"-"`
	Explanation      string        `json:"-"`
	File             *FileUpload   `json:"-"`
}

// SubmitAppeal submits an appeal for a brand with a file attachment.
func (s *BrandsService) SubmitAppeal(ctx context.Context, params SubmitAppealParams, opts *RequestOptions) (*Response, error) {
	if err := require(map[string]string{
		"brandId":     params.BrandID,
		"explanation": params.Explanation,
	}); err != nil {
		return nil, err
	}
	if len(params.AppealCategories) == 0 {
		return nil, &ValidationError{Field: "appealCategories"}
	}
	if params.File == nil {
		return nil, &ValidationError{Field: "file"}
	}

	path := fmt.Sprintf("/brand/appeal/%s", url.PathEscape(params.BrandID))

	writeForm := func(w *multipart.Writer) error {
		// Add the file
		part, err := w.CreateFormFile("file", params.File.Filename)
		if err != nil {
			return err
		}
		if _, err := io.Copy(part, params.File.Reader); err != nil {
			return err
		}

		// Add appealCategories as a JSON-stringified array
		categoriesJSON, err := json.Marshal(params.AppealCategories)
		if err != nil {
			return err
		}
		if err := w.WriteField("appealCategories", string(categoriesJSON)); err != nil {
			return err
		}

		// Add explanation
		if err := w.WriteField("explanation", params.Explanation); err != nil {
			return err
		}

		return nil
	}

	return s.client.doMultipart(ctx, "POST", path, writeForm, opts)
}
