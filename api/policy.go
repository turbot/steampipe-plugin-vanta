package api

import (
	"context"
	"fmt"

	"github.com/machinebox/graphql"
	errorsHandler "github.com/turbot/steampipe-plugin-vanta/errors"
)

type PolicyUploadedDoc struct {
	Id     string `json:"id"`
	SlugId string `json:"slugId"`
	Url    string `json:"url"`
}

type PolicyAcceptanceControl struct {
	Id               string                                   `json:"id"`
	Name             string                                   `json:"name"`
	StandardSections []PolicyAcceptanceControlStandardSection `json:"standardSections"`
	Assignees        []PolicyAcceptanceControlAssignee        `json:"assignees"`
}

type PolicyAcceptanceControlStandardSection struct {
	Standard string   `json:"standard"`
	Sections []string `json:"sections"`
}

type PolicyAcceptanceControlAssignee struct {
	Id          string `json:"id"`
	DisplayName string `json:"displayName"`
}

// Data about a policy within Vanta
type Policy struct {
	ApprovedAt       string            `json:"approvedAt"`
	Approver         PolicyApprover    `json:"approver"`
	CreatedAt        string            `json:"createdAt"`
	DisplayName      string            `json:"displayName"`
	Id               string            `json:"id"`
	Metadata         PolicyDocStub     `json:"-"`
	NumUsers         int               `json:"numUsers"`
	NumUsersAccepted int               `json:"numUsersAccepted"`
	OrganizationName string            `json:"-"`
	PolicyType       string            `json:"policyType"`
	Source           string            `json:"source"`
	Title            string            `json:"title"`
	UpdatedAt        string            `json:"updatedAt"`
	UploadedDoc      PolicyUploadedDoc `json:"uploadedDoc"`
	Uploader         PolicyUploader    `json:"uploader"`
}

type PolicyApprover struct {
	Id          string `json:"id"`
	DisplayName string `json:"displayName"`
}

type PolicyUploader struct {
	Id          string `json:"id"`
	DisplayName string `json:"displayName"`
}

type PolicyDocStub struct {
	Description              string                    `json:"description"`
	PolicyType               string                    `json:"policyType"`
	Status                   string                    `json:"status"`
	EmployeeAcceptanceTestId string                    `json:"employeeAcceptanceTestId"`
	AcceptanceControls       []PolicyAcceptanceControl `json:"acceptanceControls"`
}

type PolicyQueryOrganization struct {
	Name           string          `json:"name"`
	Policies       []Policy        `json:"policies"`
	PolicyDocStubs []PolicyDocStub `json:"policyDocStubs"`
}

// ListPoliciesResponse is returned by ListPolicies on success
type ListPoliciesResponse struct {
	Organization PolicyQueryOrganization `json:"organization"`
}

// Define the query
const (
	queryPolicyList = `
query PoliciesAndPolicyDocStubs {
  organization {
    name
    policies(onlyApproved: false) {
      id
      ...PoliciesV2PolicyInfo
    }
    policyDocStubs(includeDisabled: false) {
      ...PolicyDocStubInfo
    }
  }
}

fragment PoliciesV2PolicyInfo on Policy {
  id
  title
  policyType
  createdAt
  updatedAt
  approvedAt
  approver {
    id
    displayName
  }
  approverName
  uploader {
    id
    displayName
  }
  numUsers
  numUsersAccepted
  source
  uploadedDoc {
    id
    slugId
    url
  }
}

fragment PolicyDocStubInfo on policyDocStub {
  description
  policyType
  status
  employeeAcceptanceTestId
  acceptanceControls {
    id
    name
    standardSections {
      standard
      sections
    }
    assignees {
      id
      displayName
    }
  }
}
`
)

// ListPolicies returns a list of the most recent policies of each type
//
// @param ctx context for configuration
//
// @param client the API client
func ListPolicies(
	ctx context.Context,
	client *Client,
) (*ListPoliciesResponse, error) {

	// Make a request
	req := graphql.NewRequest(queryPolicyList)

	// set header fields
	req.Header.Set("user-agent", "steampipe")
	req.Header.Set("Cookie", fmt.Sprintf("connect.sid=%s", *client.SessionId))
	req.Header.Set("x-csrf-token", "this_csrf_header_is_constant")
	req.Header.Set("Cache-Control", "no-cache")

	var err error
	var data ListPoliciesResponse

	if err := client.Graphql.Run(ctx, req, &data); err != nil {
		err = errorsHandler.BuildErrorMessage(err)
		return nil, err
	}

	return &data, err
}
