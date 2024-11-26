package api

import (
	"context"
	"fmt"

	"github.com/machinebox/graphql"
	errorsHandler "github.com/turbot/steampipe-plugin-vanta/errors"
)

type IntegrationControl struct {
	ControlId string `json:"controlId"`
}

type Test struct {
	TestId   string               `json:"testId"`
	Controls []IntegrationControl `json:"controls"`
}

type IntegrationCredential struct {
	DisableCause       string      `json:"disableCause"`
	DisabledFetchKinds interface{} `json:"disabledFetchKinds"`
	EnabledProducts    interface{} `json:"enabledProducts"`
	ExternalAccountId  string      `json:"externalAccountId"`
	Id                 string      `json:"id"`
	IsDisabled         bool        `json:"isDisabled"`
	LastUpdated        string      `json:"lastUpdated"`
	Metadata           string      `json:"metadata"`
	Service            string      `json:"service"`
	WillExpire         string      `json:"willExpire"`
}

type Integration struct {
	ApplicationUrl        string                  `json:"applicationUrl"`
	Credentials           []IntegrationCredential `json:"credentials"`
	Description           string                  `json:"description"`
	DisplayName           string                  `json:"displayName"`
	Id                    string                  `json:"id"`
	InstallationUrl       string                  `json:"installationUrl"`
	IntegrationCategories []string                `json:"integrationCategories"`
	LogoSlugId            string                  `json:"logoSlugId"`
	OrganizationName      string                  `json:"-"`
	ScopableResources     []string                `json:"scopableResources"`
	ServiceCategories     []string                `json:"serviceCategories"`
	Tests                 []Test                  `json:"tests"`
}

// Relay-style edge for Vanta integration
type IntegrationEdge struct {
	Integration Integration `json:"node"`
}

// Paginated list of Vanta integrations
type Integrations struct {
	Edges      []IntegrationEdge `json:"edges"`
	PageInfo   PageInfo          `json:"pageInfo"`
	TotalCount int               `json:"totalCount"`
}

type IntegrationQueryOrganization struct {
	Name         string       `json:"name"`
	Integrations Integrations `json:"integrations"`
}

// ListIntegrationsResponse is returned by ListIntegrations on success
type ListIntegrationsResponse struct {
	Organization IntegrationQueryOrganization `json:"organization"`
}

type ListIntegrationsRequestConfiguration struct {
	// The maximum number of results to return in a single call. To retrieve the
	// remaining results, make another call with the returned EndCursor value.
	//
	// Maximum limit is 50.
	Limit int

	// When paginating forwards, the cursor to continue.
	EndCursor string

	// If true, only connected integrations will be listed.
	// Set false to list all the available integrations.
	// Default is true.
	OnlyConnected *bool
}

// Define the query
const (
	queryIntegrationList = `
query ListIntegrations($first: Int!, $after: String) {
  organization {
    name
    integrations(first: $first, after: $after) {
      totalCount
      pageInfo {
        endCursor
        hasNextPage
      }
      edges {
        node {
          id
          displayName
          description
          serviceCategories
          integrationCategories
          scopableResources
          tests {
            testId
            controls {
              controlId
            }
          }
          ... on FirstPartyIntegration {
            helpCenterArticleLink
            permissionsDescription
            additionalInformation
            credentials {
              id
              service
              isDisabled
              lastUpdated
              ... on FirstPartyIntegrationCredential {
                externalAccountId
                disableCause
                metadata
                willExpire {
                  expirationDate
                  customerFacingSummary
                  helpdeskArticleUrl
                }
                disabledFetchKinds {
                  kind
                }
                enabledProducts
              }
            }
          }
          ... on ThirdPartyIntegration {
            logoSlugId
            applicationUrl
            installationUrl
            credentials {
              id
              service
              isDisabled
              lastUpdated
              ... on ThirdPartyIntegrationCredential {
                sourceId
                dataLastReceivedAt
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

// ListIntegrations returns a paginated list of vendors
//
// @param ctx context for configuration
//
// @param client the API client
//
// @param options the API parameters
func ListIntegrations(
	ctx context.Context,
	client *Client,
	options *ListIntegrationsRequestConfiguration,
) (*ListIntegrationsResponse, error) {

	// Make a request
	req := graphql.NewRequest(queryIntegrationList)

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
	var data ListIntegrationsResponse

	if err := client.Graphql.Run(ctx, req, &data); err != nil {
		err = errorsHandler.BuildErrorMessage(err)
		return nil, err
	}

	return &data, err
}
