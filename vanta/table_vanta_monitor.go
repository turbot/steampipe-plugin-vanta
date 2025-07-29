package vanta

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
	"github.com/turbot/steampipe-plugin-vanta/restapi/model"
)

//// TABLE DEFINITION

func tableVantaMonitor(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "vanta_monitor",
		Description: "Vanta Monitor",
		List: &plugin.ListConfig{
			Hydrate: listVantaMonitors,
			KeyColumns: plugin.KeyColumnSlice{
				{Name: "status", Require: plugin.Optional},
				{Name: "category", Require: plugin.Optional},
			},
		},
		Get: &plugin.GetConfig{
			Hydrate:    getVantaMonitor,
			KeyColumns: plugin.SingleColumn("id"),
		},
		Columns: []*plugin.Column{
			// Available columns from REST API
			{Name: "id", Type: proto.ColumnType_STRING, Transform: transform.FromField("ID"), Description: "An internal Vanta generated ID of the test."},
			{Name: "name", Type: proto.ColumnType_STRING, Transform: transform.FromField("Name"), Description: "A human-readable name of the test."},
			{Name: "last_test_run_date", Type: proto.ColumnType_TIMESTAMP, Transform: transform.FromField("LastTestRunDate"), Description: "The date when the test was last run."},
			{Name: "latest_flip_date", Type: proto.ColumnType_TIMESTAMP, Transform: transform.FromField("LatestFlipDate"), Description: "The last time the test flipped to a passing or failing state."},
			{Name: "description", Type: proto.ColumnType_STRING, Transform: transform.FromField("Description"), Description: "A human-readable description of the test."},
			{Name: "failure_description", Type: proto.ColumnType_STRING, Transform: transform.FromField("FailureDescription"), Description: "Description of what failure means for this test."},
			{Name: "remediation_description", Type: proto.ColumnType_STRING, Transform: transform.FromField("RemediationDescription"), Description: "Description of how to remediate failures for this test."},
			{Name: "version", Type: proto.ColumnType_JSON, Transform: transform.FromField("Version"), Description: "Version information for the test."},
			{Name: "category", Type: proto.ColumnType_STRING, Transform: transform.FromField("Category"), Description: "A high-level categorization of the test."},
			{Name: "integrations", Type: proto.ColumnType_JSON, Transform: transform.FromField("Integrations"), Description: "List of integrations associated with this test."},
			{Name: "status", Type: proto.ColumnType_STRING, Transform: transform.FromField("Status"), Description: "Current status of the test."},
			{Name: "deactivated_status_info", Type: proto.ColumnType_JSON, Transform: transform.FromField("DeactivatedStatusInfo"), Description: "Information about deactivation status."},
			{Name: "remediation_status_info", Type: proto.ColumnType_JSON, Transform: transform.FromField("RemediationStatusInfo"), Description: "Specifies the remediation information."},
			{Name: "owner", Type: proto.ColumnType_JSON, Transform: transform.FromField("Owner"), Description: "Owner information for the test."},

			// Derived columns from nested data
			{Name: "owner_display_name", Type: proto.ColumnType_STRING, Transform: transform.From(getMonitorOwnerDisplayName), Description: "Display name of the test owner."},
			{Name: "owner_email", Type: proto.ColumnType_STRING, Transform: transform.From(getMonitorOwnerEmail), Description: "Email address of the test owner."},
			{Name: "is_deactivated", Type: proto.ColumnType_BOOL, Transform: transform.From(getMonitorIsDeactivated), Description: "Whether the test is deactivated."},
			{Name: "deactivated_reason", Type: proto.ColumnType_STRING, Transform: transform.From(getMonitorDeactivatedReason), Description: "Reason for deactivation if the test is deactivated."},
			{Name: "remediation_status", Type: proto.ColumnType_STRING, Transform: transform.From(getMonitorRemediationStatus), Description: "Status of remediation efforts."},
			{Name: "remediation_item_count", Type: proto.ColumnType_INT, Transform: transform.From(getMonitorRemediationItemCount), Description: "Number of items requiring remediation."},
			{Name: "version_major", Type: proto.ColumnType_INT, Transform: transform.From(getMonitorVersionMajor), Description: "Major version number of the test."},
			{Name: "version_minor", Type: proto.ColumnType_INT, Transform: transform.From(getMonitorVersionMinor), Description: "Minor version number of the test."},

			// Backward compatibility columns (mapped from REST API data)
			{Name: "test_id", Type: proto.ColumnType_STRING, Transform: transform.FromField("ID"), Description: "A unique identifier for this test (same as id)."},
			{Name: "latest_flip_time", Type: proto.ColumnType_TIMESTAMP, Transform: transform.FromField("LatestFlipDate"), Description: "The last time the test flipped to a passing or failing state."},
			{Name: "outcome", Type: proto.ColumnType_STRING, Transform: transform.From(getMonitorOutcome), Description: "Outcome of the test run (mapped from status)."},
			{Name: "compliance_status", Type: proto.ColumnType_STRING, Transform: transform.FromField("Status"), Description: "The compliance status of the test."},
			{Name: "services", Type: proto.ColumnType_JSON, Transform: transform.FromField("Integrations"), Description: "A list of services (mapped from integrations)."},
			{Name: "assignees", Type: proto.ColumnType_JSON, Transform: transform.From(getMonitorAssignees), Description: "A list of users assigned as owner for this test."},
			{Name: "disabled_status", Type: proto.ColumnType_JSON, Transform: transform.FromField("DeactivatedStatusInfo"), Description: "Metadata about whether this test is disabled."},
			{Name: "controls", Type: proto.ColumnType_JSON, Description: "[DEPRECATED] A list of controls being checked during the test."},
			{Name: "organization_name", Type: proto.ColumnType_STRING, Description: "[DEPRECATED] The name of the organization."},
			//{Name: "failing_resource_entities", Type: proto.ColumnType_JSON, Description: "[DEPRECATED] Failing resource entities (requires separate API call).", Hydrate: listVantaMonitorFailingResourceEntities, Transform: transform.FromValue()},
		},
	}
}

//// LIST FUNCTION

func listVantaMonitors(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// Create REST client
	client, err := getClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("vanta_monitor.listVantaMonitors", "connection_error", err)
		return nil, err
	}

	// Default page limit
	pageLimit := 100

	// Adjust page limit if query limit is smaller
	if d.QueryContext.Limit != nil && int(*d.QueryContext.Limit) < pageLimit {
		pageLimit = int(*d.QueryContext.Limit)
	}

	options := &model.ListMonitorsOptions{
		Limit:  pageLimit,
		Cursor: "",
	}

	for {
		result, err := client.ListMonitors(ctx, options)
		if err != nil {
			plugin.Logger(ctx).Error("vanta_monitor.listVantaMonitors", "query_error", err)
			return nil, err
		}

		for _, monitor := range result.Results.Data {
			// Apply optional filters
			if shouldFilterMonitor(d, monitor) {
				continue
			}

			// Stream the raw Monitor object
			d.StreamListItem(ctx, monitor)

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

func getVantaMonitor(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	id := d.EqualsQualString("id")
	if id == "" {
		return nil, nil
	}

	// Create REST client
	client, err := getClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("vanta_monitor.getVantaMonitor", "connection_error", err)
		return nil, err
	}

	monitor, err := client.GetMonitorByID(ctx, id)
	if err != nil {
		plugin.Logger(ctx).Error("vanta_monitor.getVantaMonitor", "query_error", err)
		return nil, err
	}

	if monitor == nil {
		return nil, nil
	}

	return monitor, nil
}

//// HELPER FUNCTIONS

// shouldFilterMonitor applies optional filters
func shouldFilterMonitor(d *plugin.QueryData, monitor *model.Monitor) bool {
	// Filter by status
	if statusFilter := d.EqualsQualString("status"); statusFilter != "" {
		if monitor.Status != statusFilter {
			return true
		}
	}

	// Filter by category
	if categoryFilter := d.EqualsQualString("category"); categoryFilter != "" {
		if monitor.Category != categoryFilter {
			return true
		}
	}

	return false
}

//// TRANSFORM FUNCTIONS

// getMonitorOwnerDisplayName extracts the owner display name
func getMonitorOwnerDisplayName(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	item := d.HydrateItem
	monitor, ok := item.(*model.Monitor)
	if !ok {
		return nil, nil
	}

	if monitor.Owner != nil {
		return monitor.Owner.DisplayName, nil
	}
	return nil, nil
}

// getMonitorOwnerEmail extracts the owner email
func getMonitorOwnerEmail(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	item := d.HydrateItem
	monitor, ok := item.(*model.Monitor)
	if !ok {
		return nil, nil
	}

	if monitor.Owner != nil {
		return monitor.Owner.EmailAddress, nil
	}
	return nil, nil
}

// getMonitorIsDeactivated extracts the deactivated status
func getMonitorIsDeactivated(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	item := d.HydrateItem
	monitor, ok := item.(*model.Monitor)
	if !ok {
		return nil, nil
	}

	if monitor.DeactivatedStatusInfo != nil {
		return monitor.DeactivatedStatusInfo.IsDeactivated, nil
	}
	return false, nil
}

// getMonitorDeactivatedReason extracts the deactivated reason
func getMonitorDeactivatedReason(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	item := d.HydrateItem
	monitor, ok := item.(*model.Monitor)
	if !ok {
		return nil, nil
	}

	if monitor.DeactivatedStatusInfo != nil && monitor.DeactivatedStatusInfo.DeactivatedReason != nil {
		return *monitor.DeactivatedStatusInfo.DeactivatedReason, nil
	}
	return nil, nil
}

// getMonitorRemediationStatus extracts the remediation status
func getMonitorRemediationStatus(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	item := d.HydrateItem
	monitor, ok := item.(*model.Monitor)
	if !ok {
		return nil, nil
	}

	if monitor.RemediationStatusInfo != nil {
		return monitor.RemediationStatusInfo.Status, nil
	}
	return nil, nil
}

// getMonitorRemediationItemCount extracts the remediation item count
func getMonitorRemediationItemCount(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	item := d.HydrateItem
	monitor, ok := item.(*model.Monitor)
	if !ok {
		return nil, nil
	}

	if monitor.RemediationStatusInfo != nil {
		return monitor.RemediationStatusInfo.ItemCount, nil
	}
	return nil, nil
}

// getMonitorVersionMajor extracts the major version
func getMonitorVersionMajor(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	item := d.HydrateItem
	monitor, ok := item.(*model.Monitor)
	if !ok {
		return nil, nil
	}

	if monitor.Version != nil {
		return monitor.Version.Major, nil
	}
	return nil, nil
}

// getMonitorVersionMinor extracts the minor version
func getMonitorVersionMinor(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	item := d.HydrateItem
	monitor, ok := item.(*model.Monitor)
	if !ok {
		return nil, nil
	}

	if monitor.Version != nil {
		return monitor.Version.Minor, nil
	}
	return nil, nil
}

// getMonitorOutcome maps status to outcome for backward compatibility
func getMonitorOutcome(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	item := d.HydrateItem
	monitor, ok := item.(*model.Monitor)
	if !ok {
		return nil, nil
	}

	// Map status values to legacy outcome values
	switch monitor.Status {
	case "NEEDS_ATTENTION":
		return "FAIL", nil
	case "DEACTIVATED":
		return "DISABLED", nil
	case "PASSING":
		return "PASS", nil
	default:
		return monitor.Status, nil
	}
}

// getMonitorAssignees creates assignees array for backward compatibility
func getMonitorAssignees(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	item := d.HydrateItem
	monitor, ok := item.(*model.Monitor)
	if !ok {
		return nil, nil
	}

	if monitor.Owner != nil {
		// Return as an array for backward compatibility
		return []interface{}{map[string]interface{}{
			"id":           monitor.Owner.ID,
			"emailAddress": monitor.Owner.EmailAddress,
			"displayName":  monitor.Owner.DisplayName,
		}}, nil
	}
	return []interface{}{}, nil
}

//// HYDRATE FUNCTIONS (Legacy support)

// func listVantaMonitorFailingResourceEntities(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
// 	// Get monitor ID
// 	var testId string

// 	// Try to get from REST API monitor object
// 	if monitor, ok := h.Item.(*model.Monitor); ok {
// 		testId = monitor.ID
// 	} else if data, ok := h.Item.(api.Monitor); ok {
// 		// Fallback for legacy GraphQL object
// 		testId = data.TestId
// 	} else {
// 		return []interface{}{}, nil
// 	}

// 	// Create client
// 	conn, err := getVantaAppClient(ctx, d)
// 	if err != nil {
// 		plugin.Logger(ctx).Error("vanta_monitor.listVantaMonitorFailingResourceEntities", "connection_error", err)
// 		return []interface{}{}, nil // Return empty instead of error for backward compatibility
// 	}

// 	options := &api.ListTestFailingResourceEntitiesRequestConfiguration{
// 		Limit:   100, // Default to maximum; e.g. 100
// 		TestIds: []string{testId},
// 	}

// 	var failingResourceEntities []api.Resource
// 	for {
// 		query, err := api.ListTestFailingResourceEntities(context.Background(), conn, options)
// 		if err != nil {
// 			plugin.Logger(ctx).Error("vanta_monitor.listVantaMonitorFailingResourceEntities", "query_error", err)
// 			return []interface{}{}, nil // Return empty instead of error for backward compatibility
// 		}

// 		for _, result := range query.Organization.CurrentTestResults {
// 			for _, e := range result.FailingResourceEntities.Edges {
// 				failingResourceEntities = append(failingResourceEntities, e.Node.Resource)
// 			}

// 			// Return if all resources are processed
// 			if !result.FailingResourceEntities.PageInfo.HasNextPage {
// 				return failingResourceEntities, nil
// 			}

// 			// Else set the next page cursor
// 			options.EndCursor = result.FailingResourceEntities.PageInfo.EndCursor
// 		}
// 	}
// }
