package vanta

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
	"github.com/turbot/steampipe-plugin-vanta/restapi/model"
)

//// TABLE DEFINITION

func tableVantaEvidence(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "vanta_evidence",
		Description: "Vanta Evidence",
		List: &plugin.ListConfig{
			Hydrate: listVantaEvidences,
		},
		Get: &plugin.GetConfig{
			Hydrate:    getVantaEvidence,
			KeyColumns: plugin.SingleColumn("id"),
		},
		Columns: []*plugin.Column{
			// Available columns from REST API (Vendor data)
			{Name: "id", Type: proto.ColumnType_STRING, Transform: transform.FromField("ID"), Description: "A unique identifier of the vendor."},
			{Name: "name", Type: proto.ColumnType_STRING, Transform: transform.FromField("Name"), Description: "The name of the vendor."},
			{Name: "website_url", Type: proto.ColumnType_STRING, Transform: transform.FromField("WebsiteURL"), Description: "The website URL of the vendor."},
			{Name: "account_manager_name", Type: proto.ColumnType_STRING, Transform: transform.FromField("AccountManagerName"), Description: "The name of the account manager."},
			{Name: "account_manager_email", Type: proto.ColumnType_STRING, Transform: transform.FromField("AccountManagerEmail"), Description: "The email of the account manager."},
			{Name: "services_provided", Type: proto.ColumnType_STRING, Transform: transform.FromField("ServicesProvided"), Description: "Description of services provided by the vendor."},
			{Name: "additional_notes", Type: proto.ColumnType_STRING, Transform: transform.FromField("AdditionalNotes"), Description: "Additional notes about the vendor."},
			{Name: "security_owner_user_id", Type: proto.ColumnType_STRING, Transform: transform.FromField("SecurityOwnerUserID"), Description: "The user ID of the security owner."},
			{Name: "business_owner_user_id", Type: proto.ColumnType_STRING, Transform: transform.FromField("BusinessOwnerUserID"), Description: "The user ID of the business owner."},
			{Name: "contract_start_date", Type: proto.ColumnType_TIMESTAMP, Transform: transform.FromField("ContractStartDate"), Description: "The contract start date."},
			{Name: "contract_renewal_date", Type: proto.ColumnType_TIMESTAMP, Transform: transform.FromField("ContractRenewalDate"), Description: "The contract renewal date."},
			{Name: "contract_termination_date", Type: proto.ColumnType_TIMESTAMP, Transform: transform.FromField("ContractTerminationDate"), Description: "The contract termination date."},
			{Name: "next_security_review_due_date", Type: proto.ColumnType_TIMESTAMP, Transform: transform.FromField("NextSecurityReviewDueDate"), Description: "The next security review due date."},
			{Name: "last_security_review_completion_date", Type: proto.ColumnType_TIMESTAMP, Transform: transform.FromField("LastSecurityReviewCompletionDate"), Description: "The last security review completion date."},
			{Name: "is_visible_to_auditors", Type: proto.ColumnType_BOOL, Transform: transform.FromField("IsVisibleToAuditors"), Description: "If true, the vendor is visible to auditors."},
			{Name: "is_risk_auto_scored", Type: proto.ColumnType_BOOL, Transform: transform.FromField("IsRiskAutoScored"), Description: "If true, the vendor risk is auto-scored."},
			{Name: "category", Type: proto.ColumnType_JSON, Transform: transform.FromField("Category"), Description: "The category information of the vendor."},
			{Name: "auth_details", Type: proto.ColumnType_JSON, Transform: transform.FromField("AuthDetails"), Description: "Authentication details for the vendor."},
			{Name: "risk_attribute_ids", Type: proto.ColumnType_JSON, Transform: transform.FromField("RiskAttributeIDs"), Description: "List of risk attribute IDs."},
			{Name: "status", Type: proto.ColumnType_STRING, Transform: transform.FromField("Status"), Description: "The status of the vendor."},
			{Name: "inherent_risk_level", Type: proto.ColumnType_STRING, Transform: transform.FromField("InherentRiskLevel"), Description: "The inherent risk level of the vendor."},
			{Name: "residual_risk_level", Type: proto.ColumnType_STRING, Transform: transform.FromField("ResidualRiskLevel"), Description: "The residual risk level of the vendor."},
			{Name: "vendor_headquarters", Type: proto.ColumnType_STRING, Transform: transform.FromField("VendorHeadquarters"), Description: "The headquarters location of the vendor."},
			{Name: "contract_amount", Type: proto.ColumnType_DOUBLE, Transform: transform.FromField("ContractAmount"), Description: "The contract amount."},
			{Name: "custom_fields", Type: proto.ColumnType_JSON, Transform: transform.FromField("CustomFields"), Description: "Custom fields for the vendor."},

			// Derived columns from nested data
			{Name: "category_display_name", Type: proto.ColumnType_STRING, Transform: transform.From(getCategoryDisplayName), Description: "The display name of the vendor category."},

			// Deprecated columns (original evidence columns not available in vendor REST API)
			{Name: "title", Type: proto.ColumnType_STRING, Description: "[DEPRECATED] The title of the document."},
			{Name: "evidence_request_id", Type: proto.ColumnType_STRING, Description: "[DEPRECATED] A unique identifier for this evidence request."},
			{Name: "description", Type: proto.ColumnType_STRING, Description: "[DEPRECATED] A human-readable description of the evidence requested."},
			{Name: "uid", Type: proto.ColumnType_STRING, Description: "[DEPRECATED] An identifier that is unique across all of Vanta."},
			{Name: "app_upload_enabled", Type: proto.ColumnType_BOOL, Description: "[DEPRECATED] If true, applications are allowed to upload documents on behalf of customers for this evidence request."},
			{Name: "restricted", Type: proto.ColumnType_BOOL, Description: "[DEPRECATED] If true, access to the contents of the evidence documents is restricted."},
			{Name: "dismissed_status", Type: proto.ColumnType_JSON, Description: "[DEPRECATED] Information about the dismissed status of the evidence request."},
			{Name: "renewal_metadata", Type: proto.ColumnType_JSON, Description: "[DEPRECATED] Information on the renewal cadence of the evidence request."},
			{Name: "organization_name", Type: proto.ColumnType_STRING, Description: "[DEPRECATED] The name of the organization."},
		},
	}
}

//// LIST FUNCTION

func listVantaEvidences(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// Create REST client
	client, err := getClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("vanta_evidence.listVantaEvidences", "connection_error", err)
		return nil, err
	}

	// Default page limit
	pageLimit := 100

	// Adjust page limit if query limit is smaller
	if d.QueryContext.Limit != nil && int(*d.QueryContext.Limit) < pageLimit {
		pageLimit = int(*d.QueryContext.Limit)
	}

	options := &model.ListVendorsOptions{
		Limit:  pageLimit,
		Cursor: "",
	}

	for {
		result, err := client.ListVendors(ctx, options)
		if err != nil {
			plugin.Logger(ctx).Error("vanta_evidence.listVantaEvidences", "api_error", err)
			return nil, err
		}

		for _, vendor := range result.Results.Data {
			// Stream the raw Vendor object
			d.StreamListItem(ctx, vendor)

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

func getVantaEvidence(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	id := d.EqualsQualString("id")
	if id == "" {
		return nil, nil
	}

	// Create REST client
	client, err := getClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("vanta_evidence.getVantaEvidence", "connection_error", err)
		return nil, err
	}

	vendor, err := client.GetVendorByID(ctx, id)
	if err != nil {
		plugin.Logger(ctx).Error("vanta_evidence.getVantaEvidence", "api_error", err)
		return nil, err
	}

	if vendor == nil {
		return nil, nil
	}

	return vendor, nil
}

//// TRANSFORM FUNCTIONS

// getCategoryDisplayName extracts the category display name from the vendor object
func getCategoryDisplayName(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	item := d.HydrateItem
	vendor, ok := item.(*model.Vendor)
	if !ok {
		return nil, nil
	}

	if vendor.Category != nil {
		return vendor.Category.DisplayName, nil
	}
	return nil, nil
}
