package signalhouse

import (
	"context"
	"fmt"
	"net/url"
)

// SubgroupsService handles subgroup-related API calls.
type SubgroupsService struct {
	client *Client
}

// GetSubgroupsParams are the parameters for the GetSubgroups method.
type GetSubgroupsParams struct {
	ID      string
	GroupID string
	Page    string
	Limit   string
}

// GetSubgroups returns a list of subgroups with optional filters.
func (s *SubgroupsService) GetSubgroups(ctx context.Context, params GetSubgroupsParams, opts *RequestOptions) (*Response, error) {
	qs := buildQueryString(map[string]string{
		"id":      params.ID,
		"groupId": params.GroupID,
		"page":    params.Page,
		"limit":   params.Limit,
	})
	return s.client.doJSON(ctx, "GET", "/subgroup"+qs, nil, opts)
}

// GetSubgroupCountsParams are the optional filters for the GetSubgroupCounts method.
type GetSubgroupCountsParams struct {
	GroupID    string
	SubgroupID string
	Limit      string
	Offset     string
}

// GetSubgroupCounts returns resource counts grouped by subgroup ID. Only signalhouse_admin can query
// all groups; every other allowed role is restricted to its JWT group scope.
func (s *SubgroupsService) GetSubgroupCounts(ctx context.Context, params GetSubgroupCountsParams, opts *RequestOptions) (*Response, error) {
	qs := buildQueryString(map[string]string{
		"groupId":    params.GroupID,
		"subgroupId": params.SubgroupID,
		"limit":      params.Limit,
		"offset":     params.Offset,
	})
	return s.client.doJSON(ctx, "GET", "/subgroup/counts"+qs, nil, opts)
}

// CreateSubgroupData represents the data for creating a subgroup.
type CreateSubgroupData struct {
	GroupID             string `json:"groupId"`
	SubgroupName        string `json:"subgroupName"`
	ContactFirstName    string `json:"contactFirstName,omitempty"`
	ContactLastName     string `json:"contactLastName,omitempty"`
	ContactMiddleName   string `json:"contactMiddleName,omitempty"`
	SubgroupCompanyName string `json:"subgroupCompanyName,omitempty"`
	AddressLine1        string `json:"addressLine1,omitempty"`
	AddressLine2        string `json:"addressLine2,omitempty"`
	City                string `json:"city,omitempty"`
	State               string `json:"state,omitempty"`
	Country             string `json:"country,omitempty"`
	PostalCode          string `json:"postalCode,omitempty"`
	Phone               string `json:"phone,omitempty"`
	// CarrierIDFamily and CarrierIDRegion require signalhouse_api, signalhouse_admin,
	// signalhouse_user, api, admin, developer, or billing role when set.
	CarrierIDFamily     string `json:"carrierIdFamily,omitempty"`
	CarrierIDRegion     string `json:"carrierIdRegion,omitempty"`
}

// CreateSubgroupParams are the parameters for the CreateSubgroup method.
type CreateSubgroupParams struct {
	SubgroupData CreateSubgroupData
}

// CreateSubgroup creates a new subgroup.
func (s *SubgroupsService) CreateSubgroup(ctx context.Context, params CreateSubgroupParams, opts *RequestOptions) (*Response, error) {
	return s.client.doJSON(ctx, "POST", "/subgroup", params.SubgroupData, opts)
}

// UpdateSubgroupData represents the data for updating a subgroup.
type UpdateSubgroupData struct {
	SubgroupName        string `json:"subgroupName,omitempty"`
	ContactFirstName    string `json:"contactFirstName,omitempty"`
	ContactLastName     string `json:"contactLastName,omitempty"`
	ContactMiddleName   string `json:"contactMiddleName,omitempty"`
	SubgroupCompanyName string `json:"subgroupCompanyName,omitempty"`
	AddressLine1        string `json:"addressLine1,omitempty"`
	AddressLine2        string `json:"addressLine2,omitempty"`
	City                string `json:"city,omitempty"`
	State               string `json:"state,omitempty"`
	Country             string `json:"country,omitempty"`
	PostalCode          string `json:"postalCode,omitempty"`
	Phone               string `json:"phone,omitempty"`
	Status              string `json:"status,omitempty"`
	// CarrierIDFamily and CarrierIDRegion require signalhouse_api, signalhouse_admin,
	// signalhouse_user, api, admin, developer, or billing role when set.
	CarrierIDFamily     string `json:"carrierIdFamily,omitempty"`
	CarrierIDRegion     string `json:"carrierIdRegion,omitempty"`
}

// UpdateSubgroupParams are the parameters for the UpdateSubgroup method.
type UpdateSubgroupParams struct {
	ID         string
	UpdateData UpdateSubgroupData
}

// UpdateSubgroup updates an existing subgroup.
func (s *SubgroupsService) UpdateSubgroup(ctx context.Context, params UpdateSubgroupParams, opts *RequestOptions) (*Response, error) {
	if err := require(map[string]string{"id": params.ID}); err != nil {
		return nil, err
	}
	path := fmt.Sprintf("/subgroup/%s", url.PathEscape(params.ID))
	return s.client.doJSON(ctx, "PUT", path, params.UpdateData, opts)
}

// DeleteSubgroupParams are the parameters for the DeleteSubgroup method.
type DeleteSubgroupParams struct {
	ID string
}

// DeleteSubgroup deletes (marks as inactive) a subgroup.
func (s *SubgroupsService) DeleteSubgroup(ctx context.Context, params DeleteSubgroupParams, opts *RequestOptions) (*Response, error) {
	if err := require(map[string]string{"id": params.ID}); err != nil {
		return nil, err
	}
	path := fmt.Sprintf("/subgroup/%s", url.PathEscape(params.ID))
	return s.client.doJSON(ctx, "DELETE", path, nil, opts)
}
