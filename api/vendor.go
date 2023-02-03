package api

import (
	"context"
	"fmt"

	"github.com/machinebox/graphql"
	errorsHandler "github.com/turbot/steampipe-plugin-vanta/errors"
)

type VendorAssessmentDocument struct {
	Id string `json:"id"`
}

// Vendor owner object
type VendorOwner struct {
	DisplayName string `json:"displayName"`
	Id          string `json:"id"`
}

// Data about a vendor within Vanta
type Vendor struct {
	AssessmentDocuments             []VendorAssessmentDocument `json:"assessmentDocuments"`
	Id                              string                     `json:"id"`
	LatestSecurityReviewCompletedAt string                     `json:"latestSecurityReviewCompletedAt"`
	Name                            string                     `json:"name"`
	OrganizationName                string                     `json:"-"`
	Owner                           VendorOwner                `json:"owner"`
	RiskAttributes                  []string                   `json:"riskAttributes"`
	ServicesProvided                string                     `json:"servicesProvided"`
	Severity                        string                     `json:"severity"`
	SharesCreditCardData            bool                       `json:"sharesCreditCardData"`
	SubmittedVaqs                   []string                   `json:"submittedVAQs"`
	Url                             string                     `json:"url"`
	VendorCategory                  string                     `json:"vendorCategory"`
	VendorRiskLocked                bool                       `json:"vendorRiskLocked"`
}

// Relay-style edge for vendor
type VendorEdge struct {
	Vendor Vendor `json:"node"`
}

// Paginated list of vendors
type Vendors struct {
	Edges      []VendorEdge `json:"edges"`
	PageInfo   PageInfo     `json:"pageInfo"`
	TotalCount int          `json:"totalCount"`
}

type VendorQueryOrganization struct {
	Name    string  `json:"name"`
	Vendors Vendors `json:"vendors"`
}

// ListVendorsResponse is returned by ListVendors on success
type ListVendorsResponse struct {
	Organization VendorQueryOrganization `json:"organization"`
}

type VendorFilters struct {
	// Filter using vendor risk level.
	// Supported values: HIGH, MEDIUM, LOW
	SeverityFilter []string
}

type ListVendorsRequestConfiguration struct {
	// The maximum number of results to return in a single call. To retrieve the
	// remaining results, make another call with the returned EndCursor value.
	//
	// Maximum limit is 100.
	Limit int

	// When paginating forwards, the cursor to continue.
	EndCursor string

	Filters *VendorFilters
}

// Define the query
const (
	queryVendorList = `
query ListVendors($first: Int!, $after: String, $sortParams: sortParams!, $filters: VendorFilters!) {
  organization {
    name
    vendors(first: $first, after: $after, sortParams: $sortParams, filters: $filters) {
      totalCount
      pageInfo {
        endCursor
        hasNextPage
      }
      edges {
        node {
          id
          ...VenderTableData
        }
      }
    }
  }
}

fragment VenderTableData on Vendor {
  id
  assessmentDocuments {
    id
  }
  attestationOfComplianceDocuments {
    id
  }
  baaDocuments {
    id
  }
  name
  latestSecurityReviewCompletedAt
  securityReviewCompletedAt
  owner {
    id
    displayName
  }
  servicesProvided
  severity
  sharesCreditCardData
  sharesEPHI
  submittedVAQs {
    id
  }
  url
  riskAttributes {
    id
    riskLevel
    riskCategoryId
    riskName
  }
  vendorCategory
  vendorRiskLocked
}
`
)

// ListVendors returns a paginated list of vendors
//
// @param ctx context for configuration
//
// @param client the API client
//
// @param options the API parameters
func ListVendors(
	ctx context.Context,
	client *Client,
	options *ListVendorsRequestConfiguration,
) (*ListVendorsResponse, error) {

	// Make a request
	req := graphql.NewRequest(queryVendorList)

	// Default vars
	req.Var("sortParams", map[string]interface{}{
		"field":     "name",
		"direction": 1,
	})

	// Define the query filters
	filters := map[string]interface{}{}

	if options.Filters != nil {
		if len(options.Filters.SeverityFilter) > 0 {
			filters["severityFilter"] = options.Filters.SeverityFilter
		}
	}
	req.Var("filters", filters)

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
	var data ListVendorsResponse

	if err := client.Graphql.Run(ctx, req, &data); err != nil {
		err = errorsHandler.BuildErrorMessage(err)
		return nil, err
	}

	return &data, err
}
