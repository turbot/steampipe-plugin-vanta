package vanta

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
	"github.com/turbot/steampipe-plugin-vanta/v2/rest_api/model"
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
			{Name: "name", Type: proto.ColumnType_STRING, Description: "The name of the group."},
			{Name: "id", Type: proto.ColumnType_STRING, Transform: transform.FromField("ID"), Description: "A unique identifier of the group."},
			{Name: "creation_date", Type: proto.ColumnType_TIMESTAMP, Description: "The creation date of the group."},
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

	maxLimit := int32(100)
	if d.QueryContext.Limit != nil {
		limit := int32(*d.QueryContext.Limit)
		if limit < maxLimit {
			maxLimit = limit
		}
	}

	options := &model.ListGroupsOptions{
		Limit:  int(maxLimit),
		Cursor: "",
	}

	for {
		result, err := client.ListGroups(ctx, options)
		if err != nil {
			plugin.Logger(ctx).Error("vanta_group.listVantaGroups", "api_error", err)
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
		plugin.Logger(ctx).Error("vanta_group.getVantaGroup", "api_error", err)
		return nil, err
	}

	if group == nil {
		return nil, nil
	}

	// Return the raw GroupItem object
	return group, nil
}
