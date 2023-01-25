package api

import (
	"context"

	"github.com/machinebox/graphql"
	errorsHandler "github.com/turbot/steampipe-plugin-vanta/errors"
)

// Metadata about whether this test is disabled and by whom
type TestDisabledStatus struct {
	CreatedAt      string `json:"createdAt"`
	Disabled       bool   `json:"disabled"`
	DisabledReason string `json:"disabledReason"`
	Expiration     string `json:"expiration"`
	UpdatedAt      string `json:"updatedAt"`
}

type Monitor struct {
	Category                string                      `json:"category"`
	Description             string                      `json:"description"`
	DisabledStatus          TestDisabledStatus          `json:"disabledStatus"`
	FailMessage             string                      `json:"failMessage"`
	FailureDescription      string                      `json:"failureDescription"`
	LatestFlip              string                      `json:"latestFlip"`
	Name                    string                      `json:"name"`
	OrganizationName        string                      `json:"-"`
	Outcome                 string                      `json:"outcome"`
	Remediation             string                      `json:"remediation"`
	TestId                  string                      `json:"testId"`
	Timestamp               string                      `json:"timestamp"`
	FailingResourceEntities FailingResourceEntitiesEdge `json:"failingResourceEntities"`
}

type Resource struct {
	DisplayName string `json:"displayName"`
	Uid         string `json:"uid"`
	Type        string `json:"__typename"`
}

type FailingResourceEntityResource struct {
	Resource Resource `json:"resource"`
}

type FailingResourceEntityNode struct {
	Node FailingResourceEntityResource `json:"node"`
}

type FailingResourceEntitiesEdge struct {
	Edges []FailingResourceEntityNode `json:"edges"`
}

type MonitorQueryOrganization struct {
	Name    string    `json:"name"`
	Results []Monitor `json:"currentTestResults"`
}

// ListMonitorsResponse is returned by ListMonitors on success
type ListMonitorsResponse struct {
	Organization MonitorQueryOrganization `json:"organization"`
}

type ListMonitorsRequestConfiguration struct {
	// Optional outcome to filter by; if null, test results will not be filtered by outcome.
	// Supported values are: DISABLED, FAIL, IN_PROGRESS, INVALID, NA.
	Outcome string

	// Optional list of testIds to filter by; if null or empty, test results will not be filtered by testId
	TestIds []string
}

type ListFailingResourceEntitiesRequestConfiguration struct {
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
	queryTestResultList = `
query getTestResults($filter: TestResultsFilter, $first: Int!) {
  organization {
    name
    currentTestResults(filter: $filter) {
      name
      category
      outcome
      latestFlip
      timestamp
      failMessage
      testId
      remediation
      failMessage
      failureDescription
      disabledStatus {
        createdAt
        disabled
        disabledReason
        expiration
        updatedAt
      }
      failingResourceEntities(first: $first) {
        edges {
          node {
            resource {
              ... on User {
                displayName
                uid
                __typename
              }
              ... on Vendor {
                displayName
                uid
                __typename
              }
              ... on VantaAgentMonitoredComputer {
                displayName
                uid
                __typename
              }
              ... on Evidence {
                uid
                displayName
                __typename
              }
              ... on Policy {
                displayName
                uid
                __typename
              }
            }
          }
        }
      }
    }
  }
}
`
)

// ListMonitors returns a list of most recent test runs and metadata about them
//
// @param ctx context for configuration
//
// @param client the API client
//
// @param options the API parameters
func ListMonitors(
	ctx context.Context,
	client *Client,
	options *ListMonitorsRequestConfiguration,
) (*ListMonitorsResponse, error) {

	// Make a request
	req := graphql.NewRequest(queryTestResultList)

	// Check for options and set it
	filters := map[string]interface{}{}
	if options.Outcome != "" {
		filters["outcome"] = options.Outcome
	}
	if options.TestIds != nil || len(options.TestIds) > 0 {
		filters["testIds"] = options.TestIds
	}
	req.Var("filter", filters)

	// Set the limit for the failingResourceEntities
	// Set the max limit - 100
	req.Var("first", 100)

	// set header fields
	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("Authorization", "token "+*client.Token)

	var err error
	var data ListMonitorsResponse

	if err := client.Graphql.Run(ctx, req, &data); err != nil {
		err = errorsHandler.BuildErrorMessage(err)
		return nil, err
	}

	return &data, err
}
