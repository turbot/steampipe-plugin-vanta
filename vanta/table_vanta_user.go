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
		},
		Columns: []*plugin.Column{
			{Name: "display_name", Type: proto.ColumnType_STRING, Description: "The display name of the user."},
			{Name: "uid", Type: proto.ColumnType_STRING, Description: "A unique identifier of the user."},
			{Name: "email", Type: proto.ColumnType_STRING, Description: "The email of the user."},
			{Name: "created_at", Type: proto.ColumnType_TIMESTAMP, Description: "The date and time when the user was created."},
		},
	}
}

//// LIST FUNCTION

func listVantaUsers(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// Create client
	conn, err := getClient(ctx, d)
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

	for {
		query, err := api.ListUsers(context.Background(), conn, options)
		if err != nil {
			plugin.Logger(ctx).Error("vanta_user.listVantaUsers", "query_error", err)
			return nil, err
		}

		for _, e := range query.Organization.UserList.Edges {
			user := e.User

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
