package signalhouse

import (
	"context"
	"fmt"
	"net/url"
)

// BillingService handles billing and wallet-related API calls.
//
// Date-range parameters (StartDate, EndDate) on list and history methods
// accept ISO-8601 date or timestamp strings. The server normalizes them to
// full UTC-day boundaries: StartDate -> start-of-UTC-day, EndDate -> end-of-
// UTC-day. The end day is always included; hourly resolution is not supported.
type BillingService struct {
	client *Client
}

// GetTransactionHistoryParams are the parameters for the GetTransactionHistory method.
type GetTransactionHistoryParams struct {
	GroupID         string
	SubgroupID      string
	BrandID         string
	CampaignID      string
	EntryType       string
	TransactionType string
	StartDate       string
	EndDate         string
}

// GetTransactionHistory returns transaction history with optional filters.
func (s *BillingService) GetTransactionHistory(ctx context.Context, params GetTransactionHistoryParams, opts *RequestOptions) (*Response, error) {
	qs := buildQueryString(map[string]string{
		"groupId":         params.GroupID,
		"subgroupId":      params.SubgroupID,
		"brandId":         params.BrandID,
		"campaignId":      params.CampaignID,
		"entryType":       params.EntryType,
		"transactionType": params.TransactionType,
		"startDate":       params.StartDate,
		"endDate":         params.EndDate,
	})
	return s.client.doJSON(ctx, "GET", "/billing/wallet/transactionHistory"+qs, nil, opts)
}

// GetPaymentHistoryParams are the parameters for the GetPaymentHistory method.
type GetPaymentHistoryParams struct {
	GroupID   string
	StartDate string
	EndDate   string
	Page      string
	Limit     string
	Timezone  string
}

// GetPaymentHistory returns payment history (wallet recharge transactions) with server-side pagination.
func (s *BillingService) GetPaymentHistory(ctx context.Context, params GetPaymentHistoryParams, opts *RequestOptions) (*Response, error) {
	qs := buildQueryString(map[string]string{
		"groupId":   params.GroupID,
		"startDate": params.StartDate,
		"endDate":   params.EndDate,
		"page":      params.Page,
		"limit":     params.Limit,
		"timezone":  params.Timezone,
	})
	return s.client.doJSON(ctx, "GET", "/billing/wallet/paymentHistory"+qs, nil, opts)
}

// GetWalletParams are the parameters for the GetWallet method.
type GetWalletParams struct {
	GroupID string
}

// GetWallet returns wallet information for a group.
func (s *BillingService) GetWallet(ctx context.Context, params GetWalletParams, opts *RequestOptions) (*Response, error) {
	if err := require(map[string]string{"groupId": params.GroupID}); err != nil {
		return nil, err
	}
	path := fmt.Sprintf("/billing/wallet/%s", url.PathEscape(params.GroupID))
	return s.client.doJSON(ctx, "GET", path, nil, opts)
}

// UpdateWalletParams are the parameters for the UpdateWallet method.
type UpdateWalletParams struct {
	GroupID                  string `json:"-"`
	AutoRechargeEnabled      *bool  `json:"autoRechargeEnabled,omitempty"`
	AutoRechargeThreshold    *int   `json:"autoRechargeThreshold,omitempty"`
	AutoRechargeToAmount     *int   `json:"autoRechargeToAmount,omitempty"`
	PrimaryPaymentMethodID   string `json:"primaryPaymentMethodId,omitempty"`
	SecondaryPaymentMethodID string `json:"secondaryPaymentMethodId,omitempty"`
}

// UpdateWallet updates wallet settings for a group.
func (s *BillingService) UpdateWallet(ctx context.Context, params UpdateWalletParams, opts *RequestOptions) (*Response, error) {
	if err := require(map[string]string{"groupId": params.GroupID}); err != nil {
		return nil, err
	}
	path := fmt.Sprintf("/billing/wallet/%s", url.PathEscape(params.GroupID))
	return s.client.doJSON(ctx, "PUT", path, params, opts)
}

// GetPaymentMethodsParams are the parameters for the GetPaymentMethods method.
type GetPaymentMethodsParams struct {
	GroupID string
}

// GetPaymentMethods returns payment methods for a group.
func (s *BillingService) GetPaymentMethods(ctx context.Context, params GetPaymentMethodsParams, opts *RequestOptions) (*Response, error) {
	if err := require(map[string]string{"groupId": params.GroupID}); err != nil {
		return nil, err
	}
	path := fmt.Sprintf("/billing/wallet/paymentMethods/%s", url.PathEscape(params.GroupID))
	return s.client.doJSON(ctx, "GET", path, nil, opts)
}

// AddFundsParams are the parameters for the AddFunds method.
type AddFundsParams struct {
	GroupID         string `json:"-"`
	Amount          int    `json:"amount"`
	PaymentMethodID string `json:"paymentMethodId,omitempty"`
}

// AddFunds adds funds to a group's wallet.
func (s *BillingService) AddFunds(ctx context.Context, params AddFundsParams, opts *RequestOptions) (*Response, error) {
	if err := require(map[string]string{"groupId": params.GroupID}); err != nil {
		return nil, err
	}
	if params.Amount == 0 {
		return nil, &ValidationError{Field: "amount"}
	}
	path := fmt.Sprintf("/billing/wallet/addfunds/%s", url.PathEscape(params.GroupID))
	return s.client.doJSON(ctx, "POST", path, params, opts)
}

// AddPaymentMethodParams are the parameters for the AddPaymentMethod method.
type AddPaymentMethodParams struct {
	GroupID         string `json:"-"`
	PaymentMethodID string `json:"paymentMethodId"`
}

// AddPaymentMethod adds a payment method to a group's wallet.
func (s *BillingService) AddPaymentMethod(ctx context.Context, params AddPaymentMethodParams, opts *RequestOptions) (*Response, error) {
	if err := require(map[string]string{"groupId": params.GroupID, "paymentMethodId": params.PaymentMethodID}); err != nil {
		return nil, err
	}
	path := fmt.Sprintf("/billing/wallet/paymentMethod/%s", url.PathEscape(params.GroupID))
	body := map[string]string{"paymentMethodId": params.PaymentMethodID}
	return s.client.doJSON(ctx, "POST", path, body, opts)
}

// RemovePaymentMethodParams are the parameters for the RemovePaymentMethod method.
type RemovePaymentMethodParams struct {
	GroupID         string
	PaymentMethodID string
}

// RemovePaymentMethod removes a payment method from a group's wallet.
func (s *BillingService) RemovePaymentMethod(ctx context.Context, params RemovePaymentMethodParams, opts *RequestOptions) (*Response, error) {
	if err := require(map[string]string{"groupId": params.GroupID, "paymentMethodId": params.PaymentMethodID}); err != nil {
		return nil, err
	}
	path := fmt.Sprintf("/billing/wallet/paymentMethod/%s/%s", url.PathEscape(params.GroupID), url.PathEscape(params.PaymentMethodID))
	return s.client.doJSON(ctx, "DELETE", path, nil, opts)
}

// GetInvoiceDetailsParams are the parameters for the GetInvoiceDetails method.
type GetInvoiceDetailsParams struct {
	GroupID    string
	SubgroupID string
	StartDate  string
	EndDate    string
}

// GetInvoiceDetails returns invoice details with optional filters.
func (s *BillingService) GetInvoiceDetails(ctx context.Context, params GetInvoiceDetailsParams, opts *RequestOptions) (*Response, error) {
	qs := buildQueryString(map[string]string{
		"groupId":    params.GroupID,
		"subgroupId": params.SubgroupID,
		"startDate":  params.StartDate,
		"endDate":    params.EndDate,
	})
	return s.client.doJSON(ctx, "GET", "/billing/invoiceDetails"+qs, nil, opts)
}

// GetFeesParams are the parameters for the GetFees method.
type GetFeesParams struct {
	GroupID string
}

// GetFees returns the fee schedule for a group.
func (s *BillingService) GetFees(ctx context.Context, params GetFeesParams, opts *RequestOptions) (*Response, error) {
	if err := require(map[string]string{"groupId": params.GroupID}); err != nil {
		return nil, err
	}
	path := fmt.Sprintf("/billing/fees/%s", url.PathEscape(params.GroupID))
	return s.client.doJSON(ctx, "GET", path, nil, opts)
}
