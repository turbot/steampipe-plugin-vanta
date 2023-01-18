package vanta

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-vanta/api"
)

//// TABLE DEFINITION

func tableVantaMonitor(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "vanta_monitor",
		Description: "Vanta Monitor",
		List: &plugin.ListConfig{
			Hydrate: listVantaMonitors,
			KeyColumns: plugin.KeyColumnSlice{
				{Name: "outcome", Require: plugin.Optional},
				{Name: "test_id", Require: plugin.Optional},
			},
		},
		Columns: []*plugin.Column{
			{Name: "name", Type: proto.ColumnType_STRING, Description: "A human-readable name of the test."},
			{Name: "category", Type: proto.ColumnType_STRING, Description: "A high-level categorization of the test."},
			{Name: "outcome", Type: proto.ColumnType_STRING, Description: "Outcome of the test run. Possible values are: 'PASS', 'DISABLED', 'FAIL', 'IN_PROGRESS', 'INVALID' and 'NA'."},
			{Name: "latest_flip", Type: proto.ColumnType_TIMESTAMP, Description: "The last time the test flipped to a passing or failing state."},
			{Name: "timestamp", Type: proto.ColumnType_TIMESTAMP, Description: "The time when the test was run."},
			{Name: "test_id", Type: proto.ColumnType_STRING, Description: "A unique identifier for this test."},
			{Name: "description", Type: proto.ColumnType_STRING, Description: "A human-readable description of the test."},
			{Name: "remediation", Type: proto.ColumnType_STRING, Description: "Instructions for how to fix this test if it's failing."},
			{Name: "fail_message", Type: proto.ColumnType_STRING, Description: "Describes the reason of a failed test."},
			{Name: "failure_description", Type: proto.ColumnType_STRING, Description: "Description under which the conditions the test would fail."},
			{Name: "disabled_status", Type: proto.ColumnType_JSON, Description: "Metadata about whether this test is disabled and by whom."},
			{Name: "organization_name", Type: proto.ColumnType_STRING, Description: "The name of the organization."},
		},
	}
}

//// LIST FUNCTION

func listVantaMonitors(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// Create client
	conn, err := getClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("vanta_monitor.listVantaMonitors", "connection_error", err)
		return nil, err
	}

	// API parameters
	option := &api.ListMonitorsRequestConfiguration{}

	// Check for additional filters
	if d.EqualsQualString("outcome") != "" {
		option.Outcome = d.EqualsQualString("outcome")
	}

	if d.EqualsQualString("test_id") != "" {
		option.TestIds = []string{d.EqualsQualString("test_id")}
	}

	// As of Jan 13, 2023, the query doesn't provide the paging information
	query, err := api.ListMonitors(context.Background(), conn, option)
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
