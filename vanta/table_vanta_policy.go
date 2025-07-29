package vanta

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
	"github.com/turbot/steampipe-plugin-vanta/restapi/model"
)

//// TABLE DEFINITION

func tableVantaPolicy(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "vanta_policy",
		Description: "Vanta Policy",
		List: &plugin.ListConfig{
			Hydrate: listVantaPolicies,
		},
		Get: &plugin.GetConfig{
			Hydrate:    getVantaPolicy,
			KeyColumns: plugin.SingleColumn("id"),
		},
		Columns: []*plugin.Column{
			{Name: "title", Type: proto.ColumnType_STRING, Transform: transform.FromField("Name"), Description: "The title of the policy."},
			{Name: "id", Type: proto.ColumnType_STRING, Transform: transform.FromField("ID"), Description: "A unique identifier of the policy."},
			{Name: "description", Type: proto.ColumnType_STRING, Transform: transform.FromField("Description"), Description: "A human-readable description of the policy."},
			{Name: "policy_type", Type: proto.ColumnType_STRING, Transform: transform.FromConstant(""), Description: "[DEPRECATED] The type of the policy."},
			{Name: "status", Type: proto.ColumnType_STRING, Transform: transform.FromField("Status"), Description: "The current status of the policy."},
			{Name: "approved_at", Type: proto.ColumnType_TIMESTAMP, Transform: transform.FromField("ApprovedAtDate"), Description: "The time when the policy was approved."},
			{Name: "created_at", Type: proto.ColumnType_TIMESTAMP, Transform: transform.FromConstant(nil), Description: "[DEPRECATED] The time when the policy was created."},
			{Name: "updated_at", Type: proto.ColumnType_TIMESTAMP, Transform: transform.FromConstant(nil), Description: "[DEPRECATED] The time when the policy was last modified."},
			{Name: "employee_acceptance_test_id", Type: proto.ColumnType_STRING, Transform: transform.FromConstant(""), Description: "[DEPRECATED] The test Id of the control that runs against employees policy acceptance."},
			{Name: "num_users", Type: proto.ColumnType_INT, Transform: transform.FromConstant(0), Description: "[DEPRECATED] The number of users assigned with the policy."},
			{Name: "num_users_accepted", Type: proto.ColumnType_STRING, Transform: transform.FromConstant(""), Description: "[DEPRECATED] The number of user accepted the policy."},
			{Name: "source", Type: proto.ColumnType_STRING, Transform: transform.FromConstant(""), Description: "[DEPRECATED] The source of the policy."},
			{Name: "acceptance_controls", Type: proto.ColumnType_JSON, Transform: transform.FromConstant(nil), Description: "[DEPRECATED] Specifies the acceptance controls."},
			{Name: "approver", Type: proto.ColumnType_JSON, Transform: transform.FromConstant(nil), Description: "[DEPRECATED] The Vanta user who approved the policy."},
			{Name: "standards", Type: proto.ColumnType_JSON, Transform: transform.FromConstant([]string{}), Description: "[DEPRECATED] A list of policy standards."},
			{Name: "uploaded_doc", Type: proto.ColumnType_JSON, Transform: transform.FromConstant(nil), Description: "[DEPRECATED] Specifies the docs uploaded for the policy."},
			{Name: "uploader", Type: proto.ColumnType_JSON, Transform: transform.FromConstant(nil), Description: "[DEPRECATED] The Vanta user that uploaded the document to Vanta."},
			{Name: "organization_name", Type: proto.ColumnType_STRING, Transform: transform.FromConstant(""), Description: "[DEPRECATED] The name of the organization."},
			{Name: "latest_version_status", Type: proto.ColumnType_STRING, Transform: transform.FromField("LatestVersion.Status"), Description: "The status of the latest version of the policy."},
		},
	}
}

//// LIST FUNCTION

func listVantaPolicies(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// Create REST client
	client, err := getClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("vanta_policy.listVantaPolicies", "connection_error", err)
		return nil, err
	}

	// Default page limit
	pageLimit := 100

	// Adjust page limit if query limit is smaller
	if d.QueryContext.Limit != nil && int(*d.QueryContext.Limit) < pageLimit {
		pageLimit = int(*d.QueryContext.Limit)
	}

	options := &model.ListPoliciesOptions{
		Limit:  pageLimit,
		Cursor: "",
	}

	for {
		result, err := client.ListPolicies(ctx, options)
		if err != nil {
			plugin.Logger(ctx).Error("vanta_policy.listVantaPolicies", "query_error", err)
			return nil, err
		}

		for _, policy := range result.Results.Data {
			// Stream the raw PolicyItem object
			d.StreamListItem(ctx, policy)

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

func getVantaPolicy(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	id := d.EqualsQualString("id")
	if id == "" {
		return nil, nil
	}

	// Create REST client
	client, err := getClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("vanta_policy.getVantaPolicy", "connection_error", err)
		return nil, err
	}

	policy, err := client.GetPolicyByID(ctx, id)
	if err != nil {
		plugin.Logger(ctx).Error("vanta_policy.getVantaPolicy", "query_error", err)
		return nil, err
	}

	if policy == nil {
		return nil, nil
	}

	// Return the raw PolicyItem object
	return policy, nil
}
