package api

import (
	"context"
	"fmt"

	"github.com/machinebox/graphql"
	errorsHandler "github.com/turbot/steampipe-plugin-vanta/errors"
)

// Metadata about whether this test is disabled and by whom
type TestDisabledStatus struct {
	CreatedAt  string `json:"createdAt"`
	Disabled   bool   `json:"disabled"`
	Expiration string `json:"expiration"`
}

type TestRemediationStatus struct {
	Status                 string `json:"status"`
	SoonestRemediateByDate string `json:"soonestRemediateByDate"`
}

// Information about assignee
type TestAssignee struct {
	Id          string `json:"id"`
	DisplayName string `json:"displayName"`
}

type TestControlStandardInfo struct {
	Standard string `json:"standard"`
}

type TestControlStandardSection struct {
	StandardInfo TestControlStandardInfo `json:"standardInfo"`
}

type TestControlInfo struct {
	Id               string                       `json:"id"`
	Name             string                       `json:"name"`
	StandardSections []TestControlStandardSection `json:"standardSections"`
}

// Information about an individual test run
type Monitor struct {
	Assignees               []TestAssignee        `json:"assignees"`
	Category                string                `json:"category"`
	ComplianceStatus        string                `json:"complianceStatus"`
	Controls                []TestControlInfo     `json:"controls"`
	Description             string                `json:"description"`
	DisabledStatus          TestDisabledStatus    `json:"disabledStatus"`
	Id                      string                `json:"id"`
	LatestFlipTime          string                `json:"latestFlipTime"`
	Name                    string                `json:"name"`
	OrganizationName        string                `json:"-"`
	Outcome                 string                `json:"outcome"`
	RemediationStatus       TestRemediationStatus `json:"remediationStatus"`
	Services                []string              `json:"services"`
	TestId                  string                `json:"testId"`
	UseRemediationTimelines bool                  `json:"useRemediationTimelines"`
}

type MonitorQueryOrganization struct {
	Name    string    `json:"name"`
	Results []Monitor `json:"testResults"`
}

// ListMonitorsResponse is returned by ListMonitors on success
type ListMonitorsResponse struct {
	Organization MonitorQueryOrganization `json:"organization"`
}

type Resource struct {
	DisplayName string `json:"displayName"`
	Uid         string `json:"uid"`
	Type        string `json:"__typename"`
}

// Type for an entity that was tested
type FailingResourceEntityResource struct {
	Resource Resource `json:"resource"`
}

type FailingResourceEntityNode struct {
	Node FailingResourceEntityResource `json:"node"`
}

// Relay-style connection for FailingResourceEntity
type FailingResourceEntitiesEdge struct {
	Edges      []FailingResourceEntityNode `json:"edges"`
	PageInfo   PageInfo                    `json:"pageInfo"`
	TotalCount int                         `json:"totalCount"`
}

type CurrentTestResult struct {
	FailingResourceEntities FailingResourceEntitiesEdge `json:"failingResourceEntities"`
}

type FailingResourceEntityQueryOrganization struct {
	Name               string              `json:"name"`
	CurrentTestResults []CurrentTestResult `json:"currentTestResults"`
}

// ListTestFailingResourceEntitiesResponse is returned by ListTestFailingResourceEntities on success
type ListTestFailingResourceEntitiesResponse struct {
	Organization FailingResourceEntityQueryOrganization `json:"organization"`
}

type ListTestFailingResourceEntitiesRequestConfiguration struct {
	// Required
	//
	// A list of testIds to filter by.
	TestIds []string

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
query getTestsForTestsPage {
  organization {
    name
    testResults(includeInRollout: true, filterNA: true) {
      id
      ...TestsPageTestResult
    }
  }
}

fragment TestsPageTestResult on TestResult {
  name
  category
  outcome
  latestFlipTime
  testId
  id
  description
  complianceStatus
  useRemediationTimelines
  services
  controls {
    id
    name
    standardSections {
      standardInfo {
        standard
      }
    }
  }
  assignees {
    displayName
    id
    employmentStatus
  }
  disabledStatus {
    disabled
    createdAt
    expiration
  }
  remediationStatus {
    status
    soonestRemediateByDate
    itemCount
  }
}
`

	queryListFailingResourceEntities = `
query getTestResults($filter: TestResultsFilter, $first: Int!, $after: String) {
	organization {
		name
		currentTestResults(filter: $filter) {
			failingResourceEntities(first: $first, after: $after) {
				edges {
					node {
						resource {
							uid
							displayName
							__typename
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
func ListMonitors(
	ctx context.Context,
	client *Client,
) (*ListMonitorsResponse, error) {

	// Make a request
	req := graphql.NewRequest(queryTestResultList)

	// set header fields
	req.Header.Set("user-agent", "steampipe")
	req.Header.Set("Cookie", fmt.Sprintf("connect.sid=%s", *client.SessionId))
	req.Header.Set("x-csrf-token", "this_csrf_header_is_constant")
	req.Header.Set("Cache-Control", "no-cache")

	var err error
	var data ListMonitorsResponse

	if err := client.Graphql.Run(ctx, req, &data); err != nil {
		err = errorsHandler.BuildErrorMessage(err)
		return nil, err
	}

	return &data, err
}

// ListTestFailingResourceEntities returns a paginated list of the entities that were in a failing state when the test was run
//
// @param ctx context for configuration
//
// @param client the API client
//
// @param options the API parameters
func ListTestFailingResourceEntities(
	ctx context.Context,
	client *Client,
	options *ListTestFailingResourceEntitiesRequestConfiguration,
) (*ListTestFailingResourceEntitiesResponse, error) {

	// Make a request
	req := graphql.NewRequest(queryListFailingResourceEntities)

	// Check for options and set it
	filters := map[string]interface{}{}
	if options.TestIds != nil || len(options.TestIds) > 0 {
		filters["testIds"] = options.TestIds
	}
	req.Var("filter", filters)

	// Check for options and set it
	if options.Limit > 0 {
		req.Var("first", options.Limit)
	}

	if options.EndCursor != "" {
		req.Var("after", options.EndCursor)
	}

	// set header fields
	req.Header.Set("user-agent", "steampipe")
	req.Header.Set("Cookie", fmt.Sprintf("connect.sid=%s", *client.SessionId))
	req.Header.Set("x-csrf-token", "this_csrf_header_is_constant")
	req.Header.Set("Cache-Control", "no-cache")

	var err error
	var data ListTestFailingResourceEntitiesResponse

	if err := client.Graphql.Run(ctx, req, &data); err != nil {
		err = errorsHandler.BuildErrorMessage(err)
		return nil, err
	}

	return &data, err
}
