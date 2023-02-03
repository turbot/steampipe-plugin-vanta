package api

import (
	"context"
	"fmt"

	"github.com/machinebox/graphql"
	errorsHandler "github.com/turbot/steampipe-plugin-vanta/errors"
)

// Group security requirements
type SecurityRequirementsMap struct {
	MustAcceptPolicies                bool `json:"mustAcceptPolicies"`
	MustBeBackgroundChecked           bool `json:"mustBeBackgroundChecked"`
	MustCompleteCcpaSecurityTraining  bool `json:"mustCompleteCcpaSecurityTraining"`
	MustCompleteGdprSecurityTraining  bool `json:"mustCompleteGdprSecurityTraining"`
	MustCompleteHipaaSecurityTraining bool `json:"mustCompleteHipaaSecurityTraining"`
	MustCompletePciSecurityTraining   bool `json:"mustCompletePciSecurityTraining"`
	MustCompleteSecurityTraining      bool `json:"mustCompleteSecurityTraining"`
	MustInstallLaptopMonitoring       bool `json:"mustInstallLaptopMonitoring"`
}

// Group checklist
type GroupChecklist struct {
	Id           string                  `json:"id"`
	Name         string                  `json:"name"`
	Requirements SecurityRequirementsMap `json:"securityRequirements"`
}

type EmbeddedIdpGroupMap struct {
	Id        string `json:"id"`
	Service   string `json:"service"`
	UpdatedAt string `json:"updatedAt"`
}

type Group struct {
	Checklist        GroupChecklist      `json:"checklist"`
	EmbeddedIdpGroup EmbeddedIdpGroupMap `json:"embeddedIdpGroup"`
	Id               string              `json:"id"`
	Name             string              `json:"name"`
	OrganizationName string              `json:"_"`
}

type GroupQueryOrganization struct {
	Groups []Group `json:"roles"`
	Name   string  `json:"name"`
}

// ListGroupsResponse is returned by ListGroups on success
type ListGroupsResponse struct {
	Organization GroupQueryOrganization `json:"organization"`
}

// Define the query
const (
	queryGroupList = `
query ListGroups {
  organization {
    name
    roles {
      id
      checklist {
        id
        ...ChecklistRequirementsForGroupsPage
      }
      name
      embeddedIdpGroup {
        id
        updatedAt
        service
      }
    }
  }
}

fragment ChecklistRequirementsForGroupsPage on IChecklist {
  name
  securityRequirements {
    ...SecurityRequirementsMap
  }
}

fragment SecurityRequirementsMap on securityRequirements {
  mustAcceptPolicies
  mustBeBackgroundChecked
  mustCompleteCcpaSecurityTraining
  mustCompleteGdprSecurityTraining
  mustCompleteHipaaSecurityTraining
  mustCompletePciSecurityTraining
  mustCompleteSecurityTraining
  mustInstallLaptopMonitoring
}
`
)

// ListGroups returns a list all Vanta groups
//
// @param ctx context for configuration
//
// @param client the API client
func ListGroups(
	ctx context.Context,
	client *Client,
) (*ListGroupsResponse, error) {

	// Make a request
	req := graphql.NewRequest(queryGroupList)

	// set header fields
	req.Header.Set("user-agent", "steampipe")
	req.Header.Set("Cookie", fmt.Sprintf("connect.sid=%s", *client.SessionId))
	req.Header.Set("x-csrf-token", "this_csrf_header_is_constant")
	req.Header.Set("Cache-Control", "no-cache")

	var err error
	var data ListGroupsResponse

	if err := client.Graphql.Run(ctx, req, &data); err != nil {
		err = errorsHandler.BuildErrorMessage(err)
		return nil, err
	}

	return &data, err
}
