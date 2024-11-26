package api

import (
	"context"
	"fmt"

	"github.com/machinebox/graphql"
	errorsHandler "github.com/turbot/steampipe-plugin-vanta/errors"
)

// VendorOwner represents the owner of a vendor.
type VendorOwner struct {
	Id          string `json:"id"`
	DisplayName string `json:"displayName"`
}

// VendorDocumentStatus represents the status of vendor agreements.
type VendorDocumentStatus struct {
	Baa struct {
		Required bool `json:"required"`
		Uploaded bool `json:"uploaded"`
	} `json:"baa"`
	Dpa struct {
		Required bool `json:"required"`
		Uploaded bool `json:"uploaded"`
	} `json:"dpa"`
}

// VendorRiskAttribute represents the risk attribute of a vendor.
type VendorRiskAttribute struct {
	Id             string `json:"id"`
	RiskName       string `json:"riskName"`
	RiskLevel      string `json:"riskLevel"`
	IsActive       bool   `json:"isActive"`
	Description    string `json:"description"`
	RiskCategoryId string `json:"riskCategoryId"`
}

// VendorRiskCategory represents a risk category for a vendor.
type VendorRiskCategory struct {
	RiskCategoryId   string                `json:"riskCategoryId"`
	RiskCategoryName string                `json:"riskCategoryName"`
	Active           bool                  `json:"active"`
	RiskAttributes   []VendorRiskAttribute `json:"riskAttributes"`
}

// VendorRiskProfile represents the risk profile of a vendor.
type VendorRiskProfile struct {
	RiskCategories []VendorRiskCategory `json:"riskCategories"`
}

// Vendor represents the vendor entity with associated data.
type Vendor struct {
	Id                   string `json:"id"`
	RiskLevel            string `json:"riskLevel"`
	CreatedAt            string `json:"createdAt"`
	LatestSecurityReview struct {
		CompletionDate string `json:"completionDate"`
	} `json:"latestSecurityReview"`
	Name                      string               `json:"name"`
	Url                       string               `json:"url"`
	VendorCategory            string               `json:"vendorCategory"`
	PublicTrustReportURL      string               `json:"publicTrustReportURL"`
	OrganizationName          string               `json:"organizationName"`
	VendorState               string               `json:"vendorState"`
	Owner                     VendorOwner          `json:"owner"`
	VendorRiskLocked          bool                 `json:"vendorRiskLocked"`
	RiskProfile               VendorRiskProfile    `json:"riskProfile"`
	VendorDocumentStatuses    VendorDocumentStatus `json:"vendorDocumentStatuses"`
	SecurityReviewStatus      string               `json:"securityReviewStatus"`
	ContractTerminationDate   string               `json:"contractTerminationDate"`
	LastArchivedByActor       VendorOwner          `json:"lastArchivedByActor"`
	SecurityReviewRenewalInfo struct {
		DueDate              string `json:"dueDate"`
		LatestCompletionDate string `json:"latestCompletionDate"`
	} `json:"securityReviewRenewalInfo"`
}

// VendorEdge represents a single edge in the vendors query result.
type VendorEdge struct {
	Vendor Vendor `json:"node"`
}

// Vendors represents the paginated list of vendors.
type Vendors struct {
	Edges      []VendorEdge `json:"edges"`
	PageInfo   PageInfo     `json:"pageInfo"`
	TotalCount int          `json:"totalCount"`
}

// VendorQueryOrganization represents the organization data in the query.
type VendorQueryOrganization struct {
	Name    string  `json:"name"`
	Vendors Vendors `json:"vendors"`
}

// ListVendorsResponse represents the full response of the ListVendors query.
type ListVendorsResponse struct {
	Organization VendorQueryOrganization `json:"organization"`
}

// VendorFilters defines filters for querying vendors.
type VendorFilters struct {
	SeverityFilter []string `json:"severityFilter"`
}

// ListVendorsRequestConfiguration defines request parameters for the ListVendors function.
type ListVendorsRequestConfiguration struct {
	Limit     int
	EndCursor string
	Filters   *VendorFilters
}

// GraphQL query with fragments
const queryVendorList = `
query fetchVendorTableData($first: Int!, $after: String, $sortParams: sortParams!, $filters: VendorFilters!) {
  organization {
    name
    vendors(first: $first, after: $after, sortParams: $sortParams, filters: $filters) {
      totalCount
      pageInfo {
        startCursor
        endCursor
        hasNextPage
        hasPreviousPage
      }
      edges {
        node {
          id
		  riskLevel
		  createdAt
		  latestSecurityReview {
			completionDate
		  }
          ...VendorTableContent_Vendor
        }
      }
    }
  }
}

fragment VendorTableContent_Vendor on Vendor {
  id
  name
  url
  vendorCategory
  publicTrustReportURL
  vendorState
  owner {
    id
    displayName
  }
  vendorRiskLocked
  riskProfile {
    riskCategories {
      riskCategoryId
      riskCategoryName
      active
      riskAttributes {
        id
        riskName
        riskLevel
        description
        riskCategoryId
      }
    }
  }
  vendorDocumentStatuses {
    baa {
      required
      uploaded
    }
    dpa {
      required
      uploaded
    }
  }
  securityReviewStatus
  securityReviewRenewalInfo {
    dueDate
    latestCompletionDate
  }
  contractTerminationDate
  lastArchivedByActor {
    displayName
  }
}
`

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
			filters["riskFilter"] = options.Filters.SeverityFilter
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
