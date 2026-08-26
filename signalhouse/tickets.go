package signalhouse

import (
	"context"
)

// TicketsService handles customer-ticket (Jira) API calls. All methods are
// staff-only and are exposed through the Admin sub-service only.
type TicketsService struct {
	client      *Client
	enableAdmin bool

	// Admin contains admin-only ticket methods. Only populated when EnableAdmin is true.
	Admin *TicketsAdminService
}

// TicketsAdminService contains admin-only customer ticket methods.
type TicketsAdminService struct {
	client *Client
}

// GetJiraMetadata returns Jira form metadata for the admin ticket form. Admin only.
func (s *TicketsAdminService) GetJiraMetadata(ctx context.Context, opts *RequestOptions) (*Response, error) {
	return s.client.doJSON(ctx, "GET", "/admin/ticket/jira/metadata", nil, opts)
}

// SearchJiraEpicsParams are the parameters for SearchJiraEpics.
type SearchJiraEpicsParams struct {
	Query string
}

// SearchJiraEpics searches open SHGHL epics. Admin only.
func (s *TicketsAdminService) SearchJiraEpics(ctx context.Context, params SearchJiraEpicsParams, opts *RequestOptions) (*Response, error) {
	qs := buildQueryString(map[string]string{"query": params.Query})
	return s.client.doJSON(ctx, "GET", "/admin/ticket/jira/epics"+qs, nil, opts)
}

// SearchJiraParentIssuesParams are the parameters for SearchJiraParentIssues.
type SearchJiraParentIssuesParams struct {
	Query string
}

// SearchJiraParentIssues searches open non-Epic, non-Subtask issues for the Subtask parent picker. Admin only.
func (s *TicketsAdminService) SearchJiraParentIssues(ctx context.Context, params SearchJiraParentIssuesParams, opts *RequestOptions) (*Response, error) {
	qs := buildQueryString(map[string]string{"query": params.Query})
	return s.client.doJSON(ctx, "GET", "/admin/ticket/jira/parent-issues"+qs, nil, opts)
}

// SearchJiraAssigneesParams are the parameters for SearchJiraAssignees.
type SearchJiraAssigneesParams struct {
	Query string
}

// SearchJiraAssignees searches assignable Jira users. Admin only.
func (s *TicketsAdminService) SearchJiraAssignees(ctx context.Context, params SearchJiraAssigneesParams, opts *RequestOptions) (*Response, error) {
	qs := buildQueryString(map[string]string{"query": params.Query})
	return s.client.doJSON(ctx, "GET", "/admin/ticket/jira/assignees"+qs, nil, opts)
}

// SearchJiraLabelsParams are the parameters for SearchJiraLabels.
type SearchJiraLabelsParams struct {
	Query string
}

// SearchJiraLabels suggests Jira labels matching a query. Admin only.
func (s *TicketsAdminService) SearchJiraLabels(ctx context.Context, params SearchJiraLabelsParams, opts *RequestOptions) (*Response, error) {
	qs := buildQueryString(map[string]string{"query": params.Query})
	return s.client.doJSON(ctx, "GET", "/admin/ticket/jira/labels"+qs, nil, opts)
}

// CreateJiraTicketData is the body for creating a customer Jira ticket.
type CreateJiraTicketData struct {
	GroupID            string   `json:"groupId"`
	ContactEmail       string   `json:"contactEmail,omitempty"`
	SubscriptionPlan   string   `json:"subscriptionPlan,omitempty"`
	IssueTypeID        string   `json:"issueTypeId"`
	PriorityID         string   `json:"priorityId"`
	Summary            string   `json:"summary"`
	Description        string   `json:"description,omitempty"`
	AssigneeAccountID  string   `json:"assigneeAccountId,omitempty"`
	ParentEpicID       string   `json:"parentEpicId,omitempty"`
	ParentEpicKey      string   `json:"parentEpicKey,omitempty"`
	SprintID           *int     `json:"sprintId,omitempty"`
	Labels             []string `json:"labels,omitempty"`
}

// CreateJiraTicketParams are the parameters for CreateJiraTicket.
type CreateJiraTicketParams struct {
	Data CreateJiraTicketData
}

// CreateJiraTicket creates a Jira ticket from the admin customer ticket form. Admin only.
func (s *TicketsAdminService) CreateJiraTicket(ctx context.Context, params CreateJiraTicketParams, opts *RequestOptions) (*Response, error) {
	if err := require(map[string]string{
		"groupId":     params.Data.GroupID,
		"issueTypeId": params.Data.IssueTypeID,
		"priorityId":  params.Data.PriorityID,
		"summary":     params.Data.Summary,
	}); err != nil {
		return nil, err
	}
	return s.client.doJSON(ctx, "POST", "/admin/ticket/jira", params.Data, opts)
}
