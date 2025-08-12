package vanta

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
	"github.com/turbot/steampipe-plugin-vanta/rest_api/model"
)

//// TABLE DEFINITION

func tableVantaVendor(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "vanta_vendor",
		Description: "Vanta Vendor",
		List: &plugin.ListConfig{
			Hydrate: listVantaVendors,
			KeyColumns: plugin.KeyColumnSlice{
				{Name: "severity", Require: plugin.Optional},
				{Name: "inherent_risk_level", Require: plugin.Optional},
				{Name: "status", Require: plugin.Optional},
			},
		},
		Get: &plugin.GetConfig{
			Hydrate:    getVantaVendor,
			KeyColumns: plugin.SingleColumn("id"),
		},
		Columns: []*plugin.Column{
			// Available columns from REST API
			{Name: "id", Type: proto.ColumnType_STRING, Transform: transform.FromField("ID"), Description: "A unique identifier of the vendor."},
			{Name: "name", Type: proto.ColumnType_STRING, Description: "The display name of the vendor."},
			{Name: "website_url", Type: proto.ColumnType_STRING, Transform: transform.FromField("WebsiteURL"), Description: "The website URL of the vendor."},
			{Name: "account_manager_name", Type: proto.ColumnType_STRING, Description: "The name of the account manager."},
			{Name: "account_manager_email", Type: proto.ColumnType_STRING, Description: "The email of the account manager."},
			{Name: "services_provided", Type: proto.ColumnType_STRING, Description: "Description of services provided by the vendor."},
			{Name: "additional_notes", Type: proto.ColumnType_STRING, Description: "Additional notes about the vendor."},
			{Name: "security_owner_user_id", Type: proto.ColumnType_STRING, Transform: transform.FromField("SecurityOwnerUserID"), Description: "The user ID of the security owner."},
			{Name: "business_owner_user_id", Type: proto.ColumnType_STRING, Transform: transform.FromField("BusinessOwnerUserID"), Description: "The user ID of the business owner."},
			{Name: "contract_start_date", Type: proto.ColumnType_TIMESTAMP, Description: "The contract start date."},
			{Name: "contract_renewal_date", Type: proto.ColumnType_TIMESTAMP, Description: "The contract renewal date."},
			{Name: "contract_termination_date", Type: proto.ColumnType_TIMESTAMP, Description: "The contract termination date."},
			{Name: "next_security_review_due_date", Type: proto.ColumnType_TIMESTAMP, Description: "The next security review due date."},
			{Name: "last_security_review_completion_date", Type: proto.ColumnType_TIMESTAMP, Description: "The last security review completion date."},
			{Name: "is_visible_to_auditors", Type: proto.ColumnType_BOOL, Description: "If true, the vendor is visible to auditors."},
			{Name: "is_risk_auto_scored", Type: proto.ColumnType_BOOL, Description: "If true, the vendor risk is auto-scored."},
			{Name: "category", Type: proto.ColumnType_JSON, Description: "The category information of the vendor."},
			{Name: "auth_details", Type: proto.ColumnType_JSON, Description: "Authentication details for the vendor."},
			{Name: "risk_attribute_ids", Type: proto.ColumnType_JSON, Transform: transform.FromField("RiskAttributeIDs"), Description: "List of risk attribute IDs."},
			{Name: "status", Type: proto.ColumnType_STRING, Description: "The status of the vendor."},
			{Name: "inherent_risk_level", Type: proto.ColumnType_STRING, Description: "The inherent risk level of the vendor."},
			{Name: "residual_risk_level", Type: proto.ColumnType_STRING, Description: "The residual risk level of the vendor."},
			{Name: "vendor_headquarters", Type: proto.ColumnType_STRING, Description: "The headquarters location of the vendor."},
			{Name: "contract_amount", Type: proto.ColumnType_DOUBLE, Description: "The contract amount."},
			{Name: "custom_fields", Type: proto.ColumnType_JSON, Description: "Custom fields for the vendor."},

			// Derived columns from nested data
			{Name: "category_display_name", Type: proto.ColumnType_STRING, Transform: transform.From(getVendorCategoryDisplayName), Description: "The display name of the vendor category."},

			// Backward compatibility columns (derived from REST API data)
			{Name: "severity", Type: proto.ColumnType_STRING, Transform: transform.From(getSeverity), Description: "The risk level of the vendor (mapped from inherent_risk_level)."},
			{Name: "url", Type: proto.ColumnType_STRING, Transform: transform.FromField("WebsiteURL"), Description: "The URL of the vendor tool."},
			{Name: "latest_security_review_completed_at", Type: proto.ColumnType_TIMESTAMP, Transform: transform.FromField("LastSecurityReviewCompletionDate"), Description: "The time when the security assessment was last reviewed."},
		},
	}
}

//// LIST FUNCTION

func listVantaVendors(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// Create REST client
	client, err := getClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("vanta_vendor.listVantaVendors", "connection_error", err)
		return nil, err
	}

	maxLimit := int32(100)
	if d.QueryContext.Limit != nil {
		limit := int32(*d.QueryContext.Limit)
		if limit < maxLimit {
			maxLimit = limit
		}
	}

	options := &model.ListVendorsOptions{
		Limit:  int(maxLimit),
		Cursor: "",
	}

	for {
		result, err := client.ListVendors(ctx, options)
		if err != nil {
			plugin.Logger(ctx).Error("vanta_vendor.listVantaVendors", "api_error", err)
			return nil, err
		}

		for _, vendor := range result.Results.Data {
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

func getVantaVendor(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	id := d.EqualsQualString("id")
	if id == "" {
		return nil, nil
	}

	// Create REST client
	client, err := getClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("vanta_vendor.getVantaVendor", "connection_error", err)
		return nil, err
	}

	vendor, err := client.GetVendorByID(ctx, id)
	if err != nil {
		plugin.Logger(ctx).Error("vanta_vendor.getVantaVendor", "api_error", err)
		return nil, err
	}

	if vendor == nil {
		return nil, nil
	}

	return vendor, nil
}

//// TRANSFORM FUNCTIONS

// getVendorCategoryDisplayName extracts the category display name from the vendor object
func getVendorCategoryDisplayName(ctx context.Context, d *transform.TransformData) (interface{}, error) {
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

// getSeverity maps inherent_risk_level to severity for backward compatibility
func getSeverity(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	item := d.HydrateItem
	vendor, ok := item.(*model.Vendor)
	if !ok {
		return nil, nil
	}

	// Map inherent risk level to severity for backward compatibility
	return vendor.InherentRiskLevel, nil
}
