package vanta

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-vanta/api"
)

//// TABLE DEFINITION

func tableVantaPolicy(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "vanta_policy",
		Description: "Vanta Policy",
		List: &plugin.ListConfig{
			Hydrate: listVantaPolicies,
		},
		Columns: []*plugin.Column{
			{Name: "display_name", Type: proto.ColumnType_STRING, Description: "The display name of the policy."},
			{Name: "uid", Type: proto.ColumnType_STRING, Description: "A unique identifier of the policy."},
			{Name: "description", Type: proto.ColumnType_STRING, Description: "A human-readable description of the policy."},
			{Name: "policy_type", Type: proto.ColumnType_STRING, Description: "The type of the policy."},
			{Name: "url", Type: proto.ColumnType_STRING, Description: "The URL of the policy."},
			{Name: "approved_at", Type: proto.ColumnType_TIMESTAMP, Description: "The time when the policy was approved."},
			{Name: "created_at", Type: proto.ColumnType_TIMESTAMP, Description: "The time when the policy was created."},
			{Name: "updated_at", Type: proto.ColumnType_TIMESTAMP, Description: "The time when the policy was last modified."},
			{Name: "pre_signed_url", Type: proto.ColumnType_STRING, Description: "The pre-signed URL of the policy."},
			{Name: "approver", Type: proto.ColumnType_JSON, Description: "The Vanta user who approved the policy."},
			{Name: "uploader", Type: proto.ColumnType_JSON, Description: "The Vanta user that uploaded the document to Vanta."},
			{Name: "title", Type: proto.ColumnType_STRING, Description: "The title of the policy."},
			{Name: "organization_name", Type: proto.ColumnType_STRING, Description: "The name of the organization."},
		},
	}
}

//// LIST FUNCTION

func listVantaPolicies(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// Create client
	conn, err := getClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("vanta_policy.listVantaPolicies", "connection_error", err)
		return nil, err
	}

	// As of Jan 13, 2023, the query doesn't provide the paging information
	query, err := api.ListPolicies(context.Background(), conn)
	if err != nil {
		plugin.Logger(ctx).Error("vanta_policy.listVantaPolicies", "query_error", err)
		return nil, err
	}

	for _, policy := range query.Organization.Policies {
		policy.OrganizationName = query.Organization.Name
		d.StreamListItem(ctx, policy)

		// Context can be cancelled due to manual cancellation or the limit has been hit
		if d.RowsRemaining(ctx) == 0 {
			return nil, nil
		}
	}

	return nil, nil
}
