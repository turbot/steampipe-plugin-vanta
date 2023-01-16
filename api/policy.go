package api

import (
	"context"

	"github.com/machinebox/graphql"
	errorsHandler "github.com/turbot/steampipe-plugin-vanta/errors"
)

// Data about a policy within Vanta
type Policy struct {
	ApprovedAt    string `json:"approvedAt"`
	Approver     User   `json:"approver"`
	CreatedAt    string `json:"createdAt"`
	Description  string `json:"description"`
	DisplayName  string `json:"displayName"`
	PolicyType   string `json:"policyType"`
	PreSignedUrl string `json:"preSignedURL"`
	Title        string `json:"title"`
	Uid          string `json:"uid"`
	UpdatedAt    string `json:"updatedAt"`
	Uploader     User   `json:"uploader"`
	Url          string `json:"url"`
}

type PolicyQueryOrganization struct {
	Policies []Policy `json:"policies"`
}

type ListPoliciesResponse struct {
	Organization PolicyQueryOrganization `json:"organization"`
}

// Define the query
const (
	queryPolicyList = `
query ListPolicies {
  organization {
    policies {
      displayName
      title
      description
      policyType
      url
      createdAt
      updatedAt
      approvedAt
      uid
      preSignedURL
      approver {
        createdAt
        displayName
        email
        uid
      }
      uploader {
        createdAt
        displayName
        email
        uid
      }
    }
  }
}
`
)

// ListPolicies returns a list of the most recent policies of each type
func ListPolicies(
	ctx context.Context,
	client *Client,
) (*ListPoliciesResponse, error) {

	// Make a request
	req := graphql.NewRequest(queryPolicyList)

	// set header fields
	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("Authorization", "token "+*client.Token)

	var err error
	var data ListPoliciesResponse

	if err := client.Graphql.Run(ctx, req, &data); err != nil {
		err = errorsHandler.BuildErrorMessage(err)
		return nil, err
	}

	return &data, err
}
