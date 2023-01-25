package vanta

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-vanta/api"
)

//// TABLE DEFINITION

func tableVantaUser(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "vanta_user",
		Description: "Vanta User",
		List: &plugin.ListConfig{
			Hydrate: listVantaUsers,
			KeyColumns: plugin.KeyColumnSlice{
				{Name: "employment_status", Require: plugin.Optional},
				{Name: "task_status", Require: plugin.Optional},
			},
		},
		Columns: []*plugin.Column{
			{Name: "display_name", Type: proto.ColumnType_STRING, Description: "The display name of the user."},
			{Name: "id", Type: proto.ColumnType_STRING, Description: "A unique identifier of the user."},
			{Name: "email", Type: proto.ColumnType_STRING, Description: "The email of the user."},
			{Name: "permission_level", Type: proto.ColumnType_STRING, Description: "The permission level of the user."},
			{Name: "employment_status", Type: proto.ColumnType_STRING, Description: "The current employment status of the user."},
			{Name: "created_at", Type: proto.ColumnType_TIMESTAMP, Description: "The time when the user was created."},
			{Name: "task_status", Type: proto.ColumnType_STRING, Description: "The security task status of the user."},
			{Name: "start_date", Type: proto.ColumnType_TIMESTAMP, Description: "The on-boarding time of the user."},
			{Name: "end_date", Type: proto.ColumnType_TIMESTAMP, Description: "The off-boarding time of the user."},
			{Name: "family_name", Type: proto.ColumnType_STRING, Description: "The family name of the user."},
			{Name: "given_name", Type: proto.ColumnType_STRING, Description: "The given name of the user."},
			{Name: "is_active", Type: proto.ColumnType_BOOL, Description: "If true, the user is active."},
			{Name: "is_not_human", Type: proto.ColumnType_BOOL, Description: "If true, the resource is not a human."},
			{Name: "is_from_scan", Type: proto.ColumnType_BOOL, Description: "If true, the user was discovered by the security scan."},
			{Name: "needs_employee_digest_reminder", Type: proto.ColumnType_BOOL, Description: "If true, user will get an email digest of their incomplete security tasks."},
			{Name: "hr_user", Type: proto.ColumnType_JSON, Description: "Specifies the embedded HR information of the user."},
			{Name: "role", Type: proto.ColumnType_JSON, Description: "Specifies the role information the user is member of."},
			{Name: "task_status_info", Type: proto.ColumnType_JSON, Description: "Specifies the security task information of the user."},
			{Name: "organization_name", Type: proto.ColumnType_STRING, Description: "The name of the organization."},
		},
	}
}

//// LIST FUNCTION

func listVantaUsers(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// Create client
	conn, err := getVantaAppClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("vanta_user.listVantaUsers", "connection_error", err)
		return nil, err
	}

	options := &api.ListUsersRequestConfiguration{}

	// Default set to 100.
	// This is the maximum number of items can be requested
	pageLimit := 100

	// Adjust page limit, if less than default value
	limit := d.QueryContext.Limit
	if limit != nil && int(*limit) < pageLimit {
		pageLimit = int(*limit)
	}
	options.Limit = pageLimit

	// Additional filters
	filters := &api.UserFilters{}
	if d.EqualsQualString("employment_status") != "" {
		filters.EmploymentStatusFilter = d.EqualsQualString("employment_status")
	}

	if d.EqualsQualString("task_status") != "" {
		filters.TaskStatusFilters = []string{d.EqualsQualString("task_status")}
	}
	options.Filters = filters

	for {
		query, err := api.ListUsers(context.Background(), conn, options)
		if err != nil {
			plugin.Logger(ctx).Error("vanta_user.listVantaUsers", "query_error", err)
			return nil, err
		}

		for _, e := range query.Organization.UserList.Edges {
			user := e.User
			user.OrganizationName = query.Organization.Name
			d.StreamListItem(ctx, user)

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}

		// Return if all resources are processed
		if !query.Organization.UserList.PageInfo.HasNextPage {
			break
		}

		// Else set the next page cursor
		options.EndCursor = query.Organization.UserList.PageInfo.EndCursor
	}

	return nil, nil
}
