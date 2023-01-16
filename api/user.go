package api

import (
	"context"

	"github.com/machinebox/graphql"
	errorsHandler "github.com/turbot/steampipe-plugin-vanta/errors"
)

// Data about a user within Vanta
type User struct {
	CreatedAt   string `json:"createdAt"`
	DisplayName string `json:"displayName"`
	Email       string `json:"email"`
	Uid         string `json:"uid"`
}

// Relay-style edge for User
type UserEdge struct {
	User User `json:"node"`
}

// Paginated list of currently active users
type UserList struct {
	Edges      []UserEdge `json:"edges"`
	PageInfo   PageInfo   `json:"pageInfo"`
	TotalCount int        `json:"totalCount"`
}

type UserQueryOrganization struct {
	UserList UserList `json:"userList"`
}

type ListUsersResponse struct {
	Organization UserQueryOrganization `json:"organization"`
}

type ListUsersRequestConfiguration struct {
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
	queryUserList = `
query ListUsers($first: Int, $after: String) {
  organization {
    userList(first: $first, after: $after) {
      pageInfo {
        hasNextPage
        endCursor
      }
      totalCount
      edges {
        node {
          uid
          email
          displayName
          createdAt
        }
      }
    }
  }
}
`
)

// ListUsers returns a paginated list of currently active users
func ListUsers(
	ctx context.Context,
	client *Client,
	options *ListUsersRequestConfiguration,
) (*ListUsersResponse, error) {

	// Make a request
	req := graphql.NewRequest(queryUserList)

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
	var data ListUsersResponse

	if err := client.Graphql.Run(ctx, req, &data); err != nil {
		err = errorsHandler.BuildErrorMessage(err)
		return nil, err
	}

	return &data, err
}
