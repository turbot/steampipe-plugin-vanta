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
	Category           string             `json:"category"`
	Description        string             `json:"description"`
	DisabledStatus     TestDisabledStatus `json:"disabledStatus"`
	FailMessage        string             `json:"failMessage"`
	FailureDescription string             `json:"failureDescription"`
	LatestFlip         string             `json:"latestFlip"`
	Name               string             `json:"name"`
	Outcome            string             `json:"outcome"`
	Remediation        string             `json:"remediation"`
	TestId             string             `json:"testId"`
	Timestamp          string             `json:"timestamp"`
}

type MonitorQueryOrganization struct {
	Results []Monitor `json:"currentTestResults"`
}

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

// Define the query
const (
	queryTestResultList = `
query getTestResults($filter: TestResultsFilter) {
  organization {
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
    }
  }
}
`
)

// ListMonitors returns a list of most recent test runs and metadata about them.
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
