package vanta

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
	"github.com/turbot/steampipe-plugin-vanta/api"
)

//// TABLE DEFINITION

func tableVantaMonitor(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "vanta_monitor",
		Description: "Vanta Monitor",
		List: &plugin.ListConfig{
			Hydrate: listVantaMonitors,
		},
		Columns: []*plugin.Column{
			{Name: "name", Type: proto.ColumnType_STRING, Description: "A human-readable name of the test."},
			{Name: "id", Type: proto.ColumnType_STRING, Description: "A unique identifier of the test."},
			{Name: "category", Type: proto.ColumnType_STRING, Description: "A high-level categorization of the test."},
			{Name: "outcome", Type: proto.ColumnType_STRING, Description: "Outcome of the test run. Possible values are: 'PASS', 'DISABLED', 'FAIL', 'IN_PROGRESS', 'INVALID' and 'NA'."},
			{Name: "test_id", Type: proto.ColumnType_STRING, Description: "A unique identifier for this test."},
			{Name: "latest_flip_time", Type: proto.ColumnType_TIMESTAMP, Description: "The last time the test flipped to a passing or failing state."},
			{Name: "compliance_status", Type: proto.ColumnType_STRING, Description: "The compliance status of the test."},
			{Name: "description", Type: proto.ColumnType_STRING, Description: "A human-readable description of the test."},
			{Name: "services", Type: proto.ColumnType_JSON, Description: "A list of services."},
			{Name: "assignees", Type: proto.ColumnType_JSON, Description: "A list of users assigned as owner for this test."},
			{Name: "disabled_status", Type: proto.ColumnType_JSON, Description: "Metadata about whether this test is disabled."},
			{Name: "remediation_status", Type: proto.ColumnType_JSON, Description: "Specifies the remediation information."},
			{Name: "controls", Type: proto.ColumnType_JSON, Description: "A list of controls being checked during the test."},
			{Name: "organization_name", Type: proto.ColumnType_STRING, Description: "The name of the organization."},
			{Name: "failing_resource_entities", Type: proto.ColumnType_JSON, Description: "The name of the organization.", Hydrate: listVantaMonitorFailingResourceEntities, Transform: transform.FromValue()},
		},
	}
}

//// LIST FUNCTION

func listVantaMonitors(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// Create client
	conn, err := getVantaAppClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("vanta_monitor.listVantaMonitors", "connection_error", err)
		return nil, err
	}

	// As of Jan 13, 2023, the query doesn't provide the paging information
	query, err := api.ListMonitors(context.Background(), conn)
	if err != nil {
		plugin.Logger(ctx).Error("vanta_monitor.listVantaMonitors", "query_error", err)
		return nil, err
	}

	for _, result := range query.Organization.Results {
		result.OrganizationName = query.Organization.Name
		d.StreamListItem(ctx, result)

		// Context can be cancelled due to manual cancellation or the limit has been hit
		if d.RowsRemaining(ctx) == 0 {
			return nil, nil
		}
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func listVantaMonitorFailingResourceEntities(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	data := h.Item.(api.Monitor)
	testId := data.TestId

	// Create client
	conn, err := getVantaAppClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("vanta_monitor.listVantaMonitorFailingResourceEntities", "connection_error", err)
		return nil, err
	}

	options := &api.ListTestFailingResourceEntitiesRequestConfiguration{
		Limit:   100, // Default to maximum; e.g. 100
		TestIds: []string{testId},
	}

	var failingResourceEntities []api.Resource
	for {
		query, err := api.ListTestFailingResourceEntities(context.Background(), conn, options)
		if err != nil {
			plugin.Logger(ctx).Error("vanta_monitor.listVantaMonitorFailingResourceEntities", "query_error", err)
			return nil, err
		}

		for _, result := range query.Organization.CurrentTestResults {
			for _, e := range result.FailingResourceEntities.Edges {
				failingResourceEntities = append(failingResourceEntities, e.Node.Resource)
			}

			// Return if all resources are processed
			if !result.FailingResourceEntities.PageInfo.HasNextPage {
				return failingResourceEntities, nil
			}

			// Else set the next page cursor
			options.EndCursor = result.FailingResourceEntities.PageInfo.EndCursor
		}
	}
}
