package vanta

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-vanta/api"
)

//// TABLE DEFINITION

func tableVantaIntegration(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "vanta_integration",
		Description: "Vanta Integration",
		List: &plugin.ListConfig{
			Hydrate: listVantaIntegrations,
		},
		Columns: []*plugin.Column{
			{Name: "display_name", Type: proto.ColumnType_STRING, Description: "The display name of the integration."},
			{Name: "id", Type: proto.ColumnType_STRING, Description: "A unique identifier of the integration."},
			{Name: "description", Type: proto.ColumnType_STRING, Description: "A human-readable description of the integration."},
			{Name: "application_url", Type: proto.ColumnType_STRING, Description: "The URL of the application."},
			{Name: "installation_url", Type: proto.ColumnType_STRING, Description: "The installation URl of the integration."},
			{Name: "logo_slug_id", Type: proto.ColumnType_STRING, Description: "The slug of the logo used for the integration."},
			{Name: "credentials", Type: proto.ColumnType_JSON, Description: "The credential metadata of the integration."},
			{Name: "integration_categories", Type: proto.ColumnType_JSON, Description: "A list of integration categories."},
			{Name: "scopable_reource", Type: proto.ColumnType_JSON, Description: "A list of scopable resources."},
			{Name: "service_categories", Type: proto.ColumnType_JSON, Description: "A list of service categories."},
			{Name: "tests", Type: proto.ColumnType_JSON, Description: "A list of tests defined for monitoring the integrations."},
			{Name: "organization_name", Type: proto.ColumnType_STRING, Description: "The name of the organization."},
		},
	}
}

//// LIST FUNCTION

func listVantaIntegrations(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// Create client
	conn, err := getVantaAppClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("vanta_integration.listVantaIntegrations", "connection_error", err)
		return nil, err
	}

	options := &api.ListIntegrationsRequestConfiguration{}

	// Default set to 50.
	// This is the maximum number of items can be requested
	pageLimit := 50

	// Adjust page limit, if less than default value
	limit := d.QueryContext.Limit
	if limit != nil && int(*limit) < pageLimit {
		pageLimit = int(*limit)
	}
	options.Limit = pageLimit

	for {
		query, err := api.ListIntegrations(context.Background(), conn, options)
		if err != nil {
			plugin.Logger(ctx).Error("vanta_integration.listVantaIntegrations", "query_error", err)
			return nil, err
		}

		for _, e := range query.Organization.Integrations.Edges {
			integration := e.Integration
			integration.OrganizationName = query.Organization.Name
			d.StreamListItem(ctx, integration)

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}

		// Return if all resources are processed
		if !query.Organization.Integrations.PageInfo.HasNextPage {
			break
		}

		// Else set the next page cursor
		options.EndCursor = query.Organization.Integrations.PageInfo.EndCursor
	}

	return nil, nil
}
