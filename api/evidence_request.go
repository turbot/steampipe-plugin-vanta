package api

import (
	"context"

	"github.com/machinebox/graphql"
	errorsHandler "github.com/turbot/steampipe-plugin-vanta/errors"
)

// Information about the dismissed status of the evidence request
type EvidenceRequestDismissedStatus struct {
	CreatedAt   string `json:"createdAt"`
	IsDismissed bool   `json:"isDismissed"`
	Reason      string `json:"reason"`
}

// Information on the renewal cadence of the evidence request
type EvidenceRequestRenewalMetadata struct {
	Cadence              string `json:"cadence"`
	CadenceLastUpdatedAt string `json:"cadenceLastUpdatedAt"`
	NextDate             string `json:"nextDate"`
}

// Data about an evidence request
type EvidenceRequest struct {
	AppUploadEnabled  bool                           `json:"appUploadEnabled"`
	Category          string                         `json:"category"`
	Description       string                         `json:"description"`
	DismissedStatus   EvidenceRequestDismissedStatus `json:"dismissedStatus"`
	EvidenceRequestId string                         `json:"evidenceRequestId"`
	OrganizationName  string                         `json:"-"`
	RenewalMetadata   EvidenceRequestRenewalMetadata `json:"renewalMetadata"`
	Restricted        bool                           `json:"restricted"`
	Title             string                         `json:"title"`
	Uid               string                         `json:"uid"`
}

// Relay-style edge for evidence request
type EvidenceRequestEdge struct {
	EvidenceRequest EvidenceRequest `json:"node"`
}

// Paginated list of evidence requests
type EvidenceRequestConnection struct {
	Edges      []EvidenceRequestEdge `json:"edges"`
	PageInfo   PageInfo              `json:"pageInfo"`
	TotalCount int                   `json:"totalCount"`
}

type EvidenceRequestQueryOrganization struct {
	Name             string                    `json:"name"`
	EvidenceRequests EvidenceRequestConnection `json:"evidenceRequests"`
}

// ListEvidenceRequestsResponse is returned by ListEvidenceRequests on success
type ListEvidenceRequestsResponse struct {
	Organization EvidenceRequestQueryOrganization `json:"organization"`
}

type ListEvidenceRequestsRequestConfiguration struct {
	// The maximum number of results to return in a single call. To retrieve the
	// remaining results, make another call with the returned EndCursor value.
	//
	// Maximum limit is 100.
	Limit int

	// When paginating forwards, the cursor to continue.
	EndCursor string
}

// Define the query
const (
	queryEvidenceRequestList = `
query ListEvidenceRequests($first: Int!, $after: String, $evidenceRequestIds: [String!]) {
  organization {
    name
    evidenceRequests(first: $first, after: $after, evidenceRequestIds: $evidenceRequestIds) {
      pageInfo {
        hasNextPage
        endCursor
      }
      totalCount
      edges {
        node {
          title
          category
          description
          appUploadEnabled
          evidenceRequestId
          restricted
          uid
          dismissedStatus {
            reason
            isDismissed
            createdAt
          }
          renewalMetadata {
            cadence
            cadenceLastUpdatedAt
            nextDate
          }
        }
      }
    }
  }
}
`
)

// ListEvidenceRequests returns a paginated list of evidence requests
//
// @param ctx context for configuration
//
// @param client the API client
//
// @param options the API parameters
func ListEvidenceRequests(
	ctx context.Context,
	client *Client,
	options *ListEvidenceRequestsRequestConfiguration,
) (*ListEvidenceRequestsResponse, error) {

	// Make a request
	req := graphql.NewRequest(queryEvidenceRequestList)

	// Check for options and set it
	if options.Limit > 0 {
		req.Var("first", options.Limit)
	}

	if options.EndCursor != "" {
		req.Var("after", options.EndCursor)
	}

	// set header fields
	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("Authorization", "token "+*client.Token)

	var err error
	var data ListEvidenceRequestsResponse

	if err := client.Graphql.Run(ctx, req, &data); err != nil {
		err = errorsHandler.BuildErrorMessage(err)
		return nil, err
	}

	return &data, err
}

// GetEvidenceRequest returns a specific evidence requests
//
// @param ctx context for configuration
//
// @param client the API client
//
// @param id unique identifier of the evidence request
func GetEvidenceRequest(
	ctx context.Context,
	client *Client,
	id string,
) (*ListEvidenceRequestsResponse, error) {

	// Make a request
	req := graphql.NewRequest(queryEvidenceRequestList)

	// Set the variables

	// expects single single object in the response body, since using filter
	req.Var("first", 1)

	if id != "" {
		req.Var("evidenceRequestIds", []string{id})
	}

	// set header fields
	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("Authorization", "token "+*client.Token)

	var err error
	var data ListEvidenceRequestsResponse

	if err := client.Graphql.Run(ctx, req, &data); err != nil {
		err = errorsHandler.BuildErrorMessage(err)
		return nil, err
	}

	return &data, err
}
