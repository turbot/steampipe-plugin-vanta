package vanta

import (
	"context"
	"strings"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
	"github.com/turbot/steampipe-plugin-vanta/api"
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
			},
		},
		Columns: []*plugin.Column{
			{Name: "name", Type: proto.ColumnType_STRING, Description: "The display name of the vendor."},
			{Name: "id", Type: proto.ColumnType_STRING, Description: "A unique identifier of the vendor."},
			{Name: "severity", Type: proto.ColumnType_STRING, Transform: transform.FromField("RiskLevel"), Description: "The risk level of the vendor."},
			{Name: "url", Type: proto.ColumnType_STRING, Description: "The URL of the vendor tool."},
			{Name: "latest_security_review_completed_at", Type: proto.ColumnType_TIMESTAMP, Transform: transform.FromField("LatestSecurityReview.CompletionDate"), Description: "The time when the security assessment was last reviewed."},
			{Name: "vendor_risk_locked", Type: proto.ColumnType_BOOL, Description: "If true, the vendor risk level is locked."},
			{Name: "owner", Type: proto.ColumnType_JSON, Description: "The owner of the vendor."},
			{Name: "risk_profile", Type: proto.ColumnType_JSON, Description: "Specifies the risk profile of the vendor."},
			{Name: "organization_name", Type: proto.ColumnType_STRING, Description: "The name of the organization."},
		},
	}
}

//// LIST FUNCTION

func listVantaVendors(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// Create client
	conn, err := getVantaAppClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("vanta_vendor.listVantaVendors", "connection_error", err)
		return nil, err
	}

	options := &api.ListVendorsRequestConfiguration{}

	// Default set to 100.
	// This is the maximum number of items can be requested
	pageLimit := 100

	// Adjust page limit, if less than default value
	limit := d.QueryContext.Limit
	if limit != nil && int(*limit) < pageLimit {
		pageLimit = int(*limit)
	}
	options.Limit = pageLimit

	// Additional filters
	filters := &api.VendorFilters{}
	if d.EqualsQualString("severity") != "" {
		severity := d.EqualsQualString("severity")
		filters.SeverityFilter = []string{strings.ToUpper(severity)}
	}
	options.Filters = filters

	for {
		query, err := api.ListVendors(context.Background(), conn, options)
		if err != nil {
			plugin.Logger(ctx).Error("vanta_vendor.listVantaVendors", "query_error", err)
			return nil, err
		}

		for _, e := range query.Organization.Vendors.Edges {
			vendor := e.Vendor
			vendor.OrganizationName = query.Organization.Name
			d.StreamListItem(ctx, vendor)

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}

		// Return if all resources are processed
		if !query.Organization.Vendors.PageInfo.HasNextPage {
			break
		}

		// Else set the next page cursor
		options.EndCursor = query.Organization.Vendors.PageInfo.EndCursor
	}

	return nil, nil
}
