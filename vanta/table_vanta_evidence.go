package vanta

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-vanta/api"
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
			KeyColumns: plugin.SingleColumn("evidence_request_id"),
		},
		Columns: []*plugin.Column{
			{Name: "title", Type: proto.ColumnType_STRING, Description: "The title of the document."},
			{Name: "evidence_request_id", Type: proto.ColumnType_STRING, Description: "A unique identifier for this evidence request."},
			{Name: "category", Type: proto.ColumnType_STRING, Description: "Specifies the category of the evidence request."},
			{Name: "description", Type: proto.ColumnType_STRING, Description: "A human-readable description of the evidence requested."},
			{Name: "uid", Type: proto.ColumnType_STRING, Description: "An identifier that is unique across all of Vanta."},
			{Name: "app_upload_enabled", Type: proto.ColumnType_BOOL, Description: "If true, applications are allowed to upload documents on behalf of customers for this evidence request."},
			{Name: "restricted", Type: proto.ColumnType_BOOL, Description: "If true, access to the contents of the evidence documents is restricted."},
			{Name: "dismissed_status", Type: proto.ColumnType_JSON, Description: "Information about the dismissed status of the evidence request."},
			{Name: "renewal_metadata", Type: proto.ColumnType_JSON, Description: "Information on the renewal cadence of the evidence request."},
			{Name: "organization_name", Type: proto.ColumnType_STRING, Description: "The name of the organization."},
		},
	}
}

//// LIST FUNCTION

func listVantaEvidences(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// Create client
	conn, err := getClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("vanta_evidence.listVantaEvidences", "connection_error", err)
		return nil, err
	}

	options := &api.ListEvidencesConfiguration{}

	// Default set to 100.
	// This is the maximum number of items can be requested
	pageLimit := 100

	// Adjust page limit, if less than default value
	limit := d.QueryContext.Limit
	if limit != nil && int(*limit) < pageLimit {
		pageLimit = int(*limit)
	}
	options.Limit = pageLimit

	for {
		query, err := api.ListEvidences(context.Background(), conn, options)
		if err != nil {
			plugin.Logger(ctx).Error("vanta_evidence.listVantaEvidences", "query_error", err)
			return nil, err
		}

		for _, e := range query.Organization.Evidences.Edges {
			user := e.Evidence
			user.OrganizationName = query.Organization.Name
			d.StreamListItem(ctx, user)

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}

		// Return if all resources are processed
		if !query.Organization.Evidences.PageInfo.HasNextPage {
			break
		}

		// Else set the next page cursor
		options.EndCursor = query.Organization.Evidences.PageInfo.EndCursor
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getVantaEvidence(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	id := d.EqualsQualString("evidence_request_id")

	// Return nil, if empty
	if id == "" {
		return nil, nil
	}

	// Create client
	conn, err := getClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("vanta_evidence.getVantaEvidence", "connection_error", err)
		return nil, err
	}

	query, err := api.GetEvidence(context.Background(), conn, id)
	if err != nil {
		plugin.Logger(ctx).Error("vanta_evidence.getVantaEvidence", "query_error", err)
		return nil, err
	}

	// Since GET uses the same LIST query with additional filter to extract the queried evidence request,
	// return the first element
	if len(query.Organization.Evidences.Edges) > 0 {
		result := query.Organization.Evidences.Edges[0].Evidence
		result.OrganizationName = query.Organization.Name

		return result, nil
	}

	return nil, nil
}
