package vanta

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
	"github.com/turbot/steampipe-plugin-vanta/restapi/model"
)

//// TABLE DEFINITION

func tableVantaIntegration(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "vanta_integration",
		Description: "Vanta Integration",
		List: &plugin.ListConfig{
			Hydrate: listVantaIntegrations,
		},
		Get: &plugin.GetConfig{
			Hydrate:    getVantaIntegration,
			KeyColumns: plugin.SingleColumn("id"),
		},
		Columns: []*plugin.Column{
			{Name: "display_name", Type: proto.ColumnType_STRING, Description: "The display name of the integration."},
			{Name: "id", Type: proto.ColumnType_STRING, Transform: transform.FromField("IntegrationID"), Description: "A unique identifier of the integration."},
			{Name: "description", Type: proto.ColumnType_STRING, Description: "[DEPRECATED] A human-readable description of the integration."},
			{Name: "application_url", Type: proto.ColumnType_STRING, Description: "[DEPRECATED] The URL of the application."},
			{Name: "installation_url", Type: proto.ColumnType_STRING, Description: "[DEPRECATED] The installation URL of the integration."},
			{Name: "logo_slug_id", Type: proto.ColumnType_STRING, Description: "[DEPRECATED] The slug of the logo used for the integration."},
			{Name: "credentials", Type: proto.ColumnType_JSON, Description: "[DEPRECATED] The credential metadata of the integration."},
			{Name: "integration_categories", Type: proto.ColumnType_JSON, Description: "[DEPRECATED] A list of integration categories."},
			{Name: "scopable_resource", Type: proto.ColumnType_JSON, Transform: transform.FromField("ResourceKinds"), Description: "A list of scopable resources (resource kinds)."},
			{Name: "service_categories", Type: proto.ColumnType_JSON, Description: "[DEPRECATED] A list of service categories."},
			{Name: "tests", Type: proto.ColumnType_JSON, Description: "[DEPRECATED] A list of tests defined for monitoring the integrations."},
			{Name: "connections", Type: proto.ColumnType_JSON, Description: "A list of connections for this integration."},
			{Name: "organization_name", Type: proto.ColumnType_STRING, Description: "[DEPRECATED] The name of the organization."},
		},
	}
}

//// LIST FUNCTION

func listVantaIntegrations(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// Create REST client
	client, err := CreateRestClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("vanta_integration.listVantaIntegrations", "connection_error", err)
		return nil, err
	}

	// Default page limit (maximum allowed by the API)
	pageLimit := 50

	// Adjust page limit if query limit is smaller
	if d.QueryContext.Limit != nil && int(*d.QueryContext.Limit) < pageLimit {
		pageLimit = int(*d.QueryContext.Limit)
	}

	options := &model.ListIntegrationsOptions{
		Limit:  pageLimit,
		Cursor: "",
	}

	for {
		result, err := client.ListConnectedIntegrations(ctx, options)
		if err != nil {
			plugin.Logger(ctx).Error("vanta_integration.listVantaIntegrations", "query_error", err)
			return nil, err
		}

		for _, integration := range result.Results.Data {
			// Stream the raw Integration object
			d.StreamListItem(ctx, integration)

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

func getVantaIntegration(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	id := d.EqualsQualString("id")
	if id == "" {
		return nil, nil
	}

	// Create REST client
	client, err := CreateRestClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("vanta_integration.getVantaIntegration", "connection_error", err)
		return nil, err
	}

	integration, err := client.GetIntegrationByID(ctx, id)
	if err != nil {
		plugin.Logger(ctx).Error("vanta_integration.getVantaIntegration", "query_error", err)
		return nil, err
	}

	if integration == nil {
		return nil, nil
	}

	return integration, nil
}
