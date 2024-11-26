package api

import (
	"context"
	"fmt"

	"github.com/machinebox/graphql"
	errorsHandler "github.com/turbot/steampipe-plugin-vanta/errors"
)

type HrUserMap struct {
	ApolloId  string `json:"apolloId"`
	EndDate   string `json:"endDate"`
	IsActive  bool   `json:"isActive"`
	JobTitle  string `json:"jobTitle"`
	Service   string `json:"service"`
	StartDate string `json:"startDate"`
	UniqueId  string `json:"uniqueId"`
}

type UserRoleMap struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type UserStatusMap struct {
	CompletionDate string `json:"completionDate"`
	DueDate        string `json:"dueDate"`
	Status         string `json:"status"`
}

// Data about a user within Vanta
type User struct {
	Id                          string        `json:"id"`
	DisplayName                 string        `json:"displayName"`
	CreatedAt                   string        `json:"createdAt"`
	Email                       string        `json:"email"`
	PermissionLevel             string        `json:"permissionLevel"`
	EmploymentStatus            string        `json:"employmentStatus"`
	TaskStatus                  string        `json:"taskStatus"`
	StartDate                   string        `json:"startDate"`
	EndDate                     string        `json:"endDate"`
	FamilyName                  string        `json:"familyName"`
	GivenName                   string        `json:"givenName"`
	IsActive                    bool          `json:"isActive"`
	IsNotHuman                  bool          `json:"isNotHuman"`
	IsFromScan                  bool          `json:"isFromScan"`
	NeedsEmployeeDigestReminder bool          `json:"needsEmployeeDigestReminder"`
	HrUser                      HrUserMap     `json:"hrUser"`
	Role                        UserRoleMap   `json:"role"`
	TaskStatusInfo              UserStatusMap `json:"taskStatusInfo"`
	OrganizationName            string        `json:"-"`
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
	Name     string   `json:"name"`
	UserList UserList `json:"people"`
}

// ListUsersResponse is returned by ListUsers on success
type ListUsersResponse struct {
	Organization UserQueryOrganization `json:"organization"`
}

// Filter object to filter the API response.
type UserFilters struct {
	// Filter using user employment status.
	// Supported values: CURRENTLY_EMPLOYED, PREVIOUSLY_EMPLOYED, NOT_PEOPLE, INACTIVE_EMPLOYEE, UPCOMING_EMPLOYEE
	EmploymentStatusFilter string

	// Filter using the security task status.
	// Supported values: SECURITY_TASKS_DUE_SOON, SECURITY_TASKS_OVERDUE, SECURITY_TASKS_COMPLETE, OFFBOARDING_DUE_SOON, OFFBOARDING_OVERDUE, OFFBOARDING_COMPLETE
	TaskStatusFilters []string
}

type ListUsersRequestConfiguration struct {
	// The maximum number of results to return in a single call. To retrieve the
	// remaining results, make another call with the returned EndCursor value.
	//
	// Maximum limit is 100.
	Limit int

	// When paginating forwards, the cursor to continue.
	EndCursor string

	// The filter object used to filter the API response.
	Filters *UserFilters
}

// Define the query
const (
	queryUserList = `
query fetchUsersForPeoplePage($first: Int!, $after: String, $sortParams: sortParams!, $filters: UserFilters!, $utcOffset: Int! = 0) {
  organization {
    id
    name
    people(first: $first, after: $after, sortParams: $sortParams, filters: $filters, utcOffset: $utcOffset) {
      totalCount
      pageInfo {
        endCursor
        hasNextPage
      }
      edges {
        node {
          id
          ...peoplePageFields
        }
      }
    }
  }
}

fragment peoplePageFields on User {
  id
  displayName
  email
  isActive
  familyName
  givenName
  createdAt
  isFromScan
  isNotHuman
  needsEmployeeDigestReminder
  startDate
  endDate
  hrUser {
    apolloId
    endDate
    jobTitle
    service
    startDate
    uniqueId
    isActive
  }
  ...peoplePageRoleFields
  ...PeoplePageUserStatuses
}

fragment peoplePageRoleFields on User {
  id
  role {
    id
    name
  }
}

fragment PeoplePageUserStatuses on User {
  id
  employmentStatus
  taskStatus
  taskStatusInfo {
    status
    dueDate
    completionDate
  }
}	
`
)

// ListUsers returns a paginated list of currently active users
//
// @param ctx context for configuration
//
// @param client the API client
//
// @param options the API parameters
func ListUsers(
	ctx context.Context,
	client *Client,
	options *ListUsersRequestConfiguration,
) (*ListUsersResponse, error) {

	// Make a request
	req := graphql.NewRequest(queryUserList)

	// Default vars
	req.Var("sortParams", map[string]interface{}{
		"field":     "name",
		"direction": 1,
	})

	// Define the query filters
	filters := map[string]interface{}{
		"includeNonHumanUsers": true,
		"includeRemovedUsers":  true,
	}

	if options.Filters != nil {
		if options.Filters.EmploymentStatusFilter != "" {
			filters["employmentStatusFilter"] = options.Filters.EmploymentStatusFilter
		}

		if len(options.Filters.TaskStatusFilters) > 0 {
			filters["taskStatusFilters"] = options.Filters.TaskStatusFilters
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
	var data ListUsersResponse

	if err := client.Graphql.Run(ctx, req, &data); err != nil {
		err = errorsHandler.BuildErrorMessage(err)
		return nil, err
	}

	return &data, err
}
