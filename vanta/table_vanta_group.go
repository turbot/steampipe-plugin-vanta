package vanta

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
	"github.com/turbot/steampipe-plugin-vanta/restapi/model"
)

//// TABLE DEFINITION

func tableVantaGroup(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "vanta_group",
		Description: "Vanta Group",
		List: &plugin.ListConfig{
			Hydrate: listVantaGroups,
		},
		Get: &plugin.GetConfig{
			Hydrate:    getVantaGroup,
			KeyColumns: plugin.SingleColumn("id"),
		},
		Columns: []*plugin.Column{
			{Name: "name", Type: proto.ColumnType_STRING, Transform: transform.FromField("Name"), Description: "The name of the group."},
			{Name: "id", Type: proto.ColumnType_STRING, Transform: transform.FromField("ID"), Description: "A unique identifier of the group."},
			{Name: "creation_date", Type: proto.ColumnType_TIMESTAMP, Transform: transform.FromField("CreationDate"), Description: "The creation date of the group."},
			{Name: "checklist", Type: proto.ColumnType_JSON, Description: "[DEPRECATED] Describes the security requirements for the group."},
			{Name: "embedded_idp_group", Type: proto.ColumnType_JSON, Description: "[DEPRECATED] A list of embedded IDP group."},
			{Name: "organization_name", Type: proto.ColumnType_STRING, Description: "[DEPRECATED] The name of the organization."},
		},
	}
}

//// LIST FUNCTION

func listVantaGroups(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// Create REST client
	client, err := getClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("vanta_group.listVantaGroups", "connection_error", err)
		return nil, err
	}

	// Default page limit
	pageLimit := 100

	// Adjust page limit if query limit is smaller
	if d.QueryContext.Limit != nil && int(*d.QueryContext.Limit) < pageLimit {
		pageLimit = int(*d.QueryContext.Limit)
	}

	options := &model.ListGroupsOptions{
		Limit:  pageLimit,
		Cursor: "",
	}

	for {
		result, err := client.ListGroups(ctx, options)
		if err != nil {
			plugin.Logger(ctx).Error("vanta_group.listVantaGroups", "query_error", err)
			return nil, err
		}

		for _, group := range result.Results.Data {
			// Stream the raw GroupItem object
			d.StreamListItem(ctx, group)

			// Check if we should stop (limit reached or context cancelled)
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}

		// Check if there are more pages
		if !result.Results.PageInfo.HasNextPage {
			break
		}

		// Set cursor for next page
		options.Cursor = result.Results.PageInfo.EndCursor
	}

	return nil, nil
}

//// GET FUNCTION

func getVantaGroup(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	id := d.EqualsQualString("id")
	if id == "" {
		return nil, nil
	}

	// Create REST client
	client, err := getClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("vanta_group.getVantaGroup", "connection_error", err)
		return nil, err
	}

	group, err := client.GetGroupByID(ctx, id)
	if err != nil {
		plugin.Logger(ctx).Error("vanta_group.getVantaGroup", "query_error", err)
		return nil, err
	}

	if group == nil {
		return nil, nil
	}

	// Return the raw GroupItem object
	return group, nil
}
