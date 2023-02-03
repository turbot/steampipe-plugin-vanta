package vanta

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-vanta/api"
)

//// TABLE DEFINITION

func tableVantaGroup(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "vanta_group",
		Description: "Vanta Group",
		List: &plugin.ListConfig{
			Hydrate: listVantaGroups,
		},
		Columns: []*plugin.Column{
			{Name: "name", Type: proto.ColumnType_STRING, Description: "The name of the group."},
			{Name: "id", Type: proto.ColumnType_STRING, Description: "A unique identifier of the group."},
			{Name: "checklist", Type: proto.ColumnType_JSON, Description: "Describes the security requirements for the group."},
			{Name: "embedded_idp_group", Type: proto.ColumnType_JSON, Description: "A list of embedded IDP group."},
			{Name: "organization_name", Type: proto.ColumnType_STRING, Description: "The name of the organization."},
		},
	}
}

//// LIST FUNCTION

func listVantaGroups(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// Create client
	conn, err := getVantaAppClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("vanta_group.listVantaGroups", "connection_error", err)
		return nil, err
	}

	// As of Jan 13, 2023, the query doesn't provide the paging information
	query, err := api.ListGroups(context.Background(), conn)
	if err != nil {
		plugin.Logger(ctx).Error("vanta_group.listVantaGroups", "query_error", err)
		return nil, err
	}

	for _, group := range query.Organization.Groups {
		group.OrganizationName = query.Organization.Name
		d.StreamListItem(ctx, group)

		// Context can be cancelled due to manual cancellation or the limit has been hit
		if d.RowsRemaining(ctx) == 0 {
			return nil, nil
		}
	}

	return nil, nil
}
