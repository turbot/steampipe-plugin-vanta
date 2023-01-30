package api

import (
	"context"

	"github.com/machinebox/graphql"
	errorsHandler "github.com/turbot/steampipe-plugin-vanta/errors"
)

// Information about the dismissed status of the evidence
type EvidenceDismissedStatus struct {
	CreatedAt   string `json:"createdAt"`
	IsDismissed bool   `json:"isDismissed"`
	Reason      string `json:"reason"`
}

// Information on the renewal cadence of the evidence
type EvidenceRenewalMetadata struct {
	Cadence              string `json:"cadence"`
	CadenceLastUpdatedAt string `json:"cadenceLastUpdatedAt"`
	NextDate             string `json:"nextDate"`
}

// Data about an evidence
type Evidence struct {
	AppUploadEnabled  bool                    `json:"appUploadEnabled"`
	Category          string                  `json:"category"`
	Description       string                  `json:"description"`
	DismissedStatus   EvidenceDismissedStatus `json:"dismissedStatus"`
	EvidenceRequestId string                  `json:"evidenceRequestId"`
	OrganizationName  string                  `json:"-"`
	RenewalMetadata   EvidenceRenewalMetadata `json:"renewalMetadata"`
	Restricted        bool                    `json:"restricted"`
	Title             string                  `json:"title"`
	Uid               string                  `json:"uid"`
}

// Relay-style edge for evidence
type EvidenceEdge struct {
	Evidence Evidence `json:"node"`
}

// Paginated list of evidences
type EvidenceConnection struct {
	Edges      []EvidenceEdge `json:"edges"`
	PageInfo   PageInfo       `json:"pageInfo"`
	TotalCount int            `json:"totalCount"`
}

type EvidenceQueryOrganization struct {
	Name      string             `json:"name"`
	Evidences EvidenceConnection `json:"evidenceRequests"`
}

// ListEvidencesResponse is returned by ListEvidences on success
type ListEvidencesResponse struct {
	Organization EvidenceQueryOrganization `json:"organization"`
}

type ListEvidencesConfiguration struct {
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
	queryEvidenceList = `
query ListEvidences($first: Int!, $after: String, $evidenceRequestIds: [String!]) {
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

// ListEvidences returns a paginated list of evidences
//
// @param ctx context for configuration
//
// @param client the API client
//
// @param options the API parameters
func ListEvidences(
	ctx context.Context,
	client *Client,
	options *ListEvidencesConfiguration,
) (*ListEvidencesResponse, error) {

	// Make a request
	req := graphql.NewRequest(queryEvidenceList)

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
	var data ListEvidencesResponse

	if err := client.Graphql.Run(ctx, req, &data); err != nil {
		err = errorsHandler.BuildErrorMessage(err)
		return nil, err
	}

	return &data, err
}

// GetEvidence returns a specific evidence
//
// @param ctx context for configuration
//
// @param client the API client
//
// @param id unique identifier of the evidence
func GetEvidence(
	ctx context.Context,
	client *Client,
	id string,
) (*ListEvidencesResponse, error) {

	// Make a request
	req := graphql.NewRequest(queryEvidenceList)

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
	var data ListEvidencesResponse

	if err := client.Graphql.Run(ctx, req, &data); err != nil {
		err = errorsHandler.BuildErrorMessage(err)
		return nil, err
	}

	return &data, err
}
