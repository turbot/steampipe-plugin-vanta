package vanta

import (
	"context"

	"github.com/turbot/go-kit/helpers"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
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
			{Name: "title", Type: proto.ColumnType_STRING, Description: "The title of the policy."},
			{Name: "id", Type: proto.ColumnType_STRING, Description: "A unique identifier of the policy."},
			{Name: "description", Type: proto.ColumnType_STRING, Description: "A human-readable description of the policy.", Transform: transform.FromField("Metadata.Description")},
			{Name: "policy_type", Type: proto.ColumnType_STRING, Description: "The type of the policy."},
			{Name: "status", Type: proto.ColumnType_STRING, Description: "The current status of the policy.", Transform: transform.FromField("Metadata.Status")},
			{Name: "approved_at", Type: proto.ColumnType_TIMESTAMP, Description: "The time when the policy was approved."},
			{Name: "created_at", Type: proto.ColumnType_TIMESTAMP, Description: "The time when the policy was created."},
			{Name: "updated_at", Type: proto.ColumnType_TIMESTAMP, Description: "The time when the policy was last modified."},
			{Name: "employee_acceptance_test_id", Type: proto.ColumnType_STRING, Description: "The test Id of the control that runs against employees policy acceptance.", Transform: transform.FromField("Metadata.EmployeeAcceptanceTestId")},
			{Name: "num_users", Type: proto.ColumnType_INT, Description: "The number of users assigned with the policy."},
			{Name: "num_users_accepted", Type: proto.ColumnType_STRING, Description: "The number of user accepted the policy."},
			{Name: "source", Type: proto.ColumnType_STRING, Description: "The source of the policy."},
			{Name: "acceptance_controls", Type: proto.ColumnType_JSON, Description: "Specifies the acceptance controls.", Transform: transform.FromField("Metadata.AcceptanceControls")},
			{Name: "approver", Type: proto.ColumnType_JSON, Description: "The Vanta user who approved the policy."},
			{Name: "standards", Type: proto.ColumnType_JSON, Description: "A list of policy standards.", Hydrate: getVantaPolicyStandards, Transform: transform.FromValue()},
			{Name: "uploaded_doc", Type: proto.ColumnType_JSON, Description: "Specifies the docs uploaded for the policy."},
			{Name: "uploader", Type: proto.ColumnType_JSON, Description: "The Vanta user that uploaded the document to Vanta."},
			{Name: "organization_name", Type: proto.ColumnType_STRING, Description: "The name of the organization."},
		},
	}
}

//// LIST FUNCTION

func listVantaPolicies(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// Create client
	conn, err := getVantaAppClient(ctx, d)
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

	// ListPolicies returns two sets of policy responses;
	// first one is Policies, contains the details of the policy; and
	// the other is PolicyDocStubs, contains the data related to policy standards and controls.
	for _, policy := range query.Organization.Policies {
		for _, policyDocStub := range query.Organization.PolicyDocStubs {
			if policy.PolicyType == policyDocStub.PolicyType {
				policy.Metadata = policyDocStub
				policy.OrganizationName = query.Organization.Name

				// Stream result
				d.StreamListItem(ctx, policy)

				// Context can be cancelled due to manual cancellation or the limit has been hit
				if d.RowsRemaining(ctx) == 0 {
					return nil, nil
				}
			}
		}
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getVantaPolicyStandards(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	data := h.Item.(api.Policy)

	var standards []string
	for _, i := range data.Metadata.AcceptanceControls {
		for _, standardSection := range i.StandardSections {
			if !helpers.StringSliceContains(standards, standardSection.Standard) {
				standards = append(standards, standardSection.Standard)
			}
		}
	}

	return standards, nil
}
