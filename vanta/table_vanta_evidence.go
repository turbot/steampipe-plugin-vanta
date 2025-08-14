package vanta

import (
	"context"
	"fmt"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
	"github.com/turbot/steampipe-plugin-vanta/v2/rest_api/model"
)

//// TABLE DEFINITION

func tableVantaEvidence(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "vanta_evidence",
		Description: "Vanta Evidence",
		List: &plugin.ListConfig{
			Hydrate: listVantaEvidences,
			KeyColumns: []*plugin.KeyColumn{
				{Name: "audit_id", Require: plugin.Required},
			},
		},
		Columns: []*plugin.Column{
			// Required parameters
			{Name: "audit_id", Type: proto.ColumnType_STRING, Transform: transform.FromQual("audit_id"), Description: "The audit ID (required parameter)."},

			// Evidence fields from API response
			{Name: "id", Type: proto.ColumnType_STRING, Transform: transform.FromField("ID"), Description: "Vanta internal reference to evidence."},
			{Name: "external_id", Type: proto.ColumnType_STRING, Transform: transform.FromField("ExternalID"), Description: "This is a static UUID to map Audit Firm controls to Vanta controls."},
			{Name: "status", Type: proto.ColumnType_STRING, Description: "Vanta internal statuses for audit evidence."},
			{Name: "name", Type: proto.ColumnType_STRING, Description: "Mutable name for evidence. Not guaranteed to be unique."},
			{Name: "deletion_date", Type: proto.ColumnType_TIMESTAMP, Description: "The date this Audit Evidence was deleted."},
			{Name: "creation_date", Type: proto.ColumnType_TIMESTAMP, Description: "The date this Audit Evidence was created."},
			{Name: "status_updated_date", Type: proto.ColumnType_TIMESTAMP, Description: "Point in time that status was last updated."},
			{Name: "test_status", Type: proto.ColumnType_STRING, Description: "The outcome of the automated test run, for Test-type evidence."},
			{Name: "evidence_type", Type: proto.ColumnType_STRING, Description: "The type of Audit Evidence."},
			{Name: "evidence_id", Type: proto.ColumnType_STRING, Transform: transform.FromField("EvidenceID"), Description: "Unique identifier for evidence."},
			{Name: "description", Type: proto.ColumnType_STRING, Description: "The description for the evidence. It will be set to null if the evidence is deleted."},
			{Name: "related_controls", Type: proto.ColumnType_JSON, Description: "The controls associated to this evidence."},

			// Derived columns from nested data
			{Name: "related_control_names", Type: proto.ColumnType_JSON, Transform: transform.From(getRelatedControlNames), Description: "Names of controls associated to this evidence."},
		},
	}
}

//// LIST FUNCTION

func listVantaEvidences(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// Get audit_id from required qualifier
	auditID := d.EqualsQualString("audit_id")
	if auditID == "" {
		return nil, fmt.Errorf("audit_id is required")
	}

	// Create REST client
	client, err := getClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("vanta_evidence.listVantaEvidences", "connection_error", err)
		return nil, err
	}

	maxLimit := int32(50)
	if d.QueryContext.Limit != nil {
		limit := int32(*d.QueryContext.Limit)
		if limit < maxLimit {
			maxLimit = limit
		}
	}

	options := &model.ListEvidenceOptions{
		AuditID: auditID,
		Limit:   int(maxLimit),
		Cursor:  "",
	}

	for {
		result, err := client.ListEvidence(ctx, auditID, options)
		if err != nil {
			plugin.Logger(ctx).Error("vanta_evidence.listVantaEvidences", "api_error", err)
			return nil, err
		}

		for _, evidence := range result.Results.Data {
			// Stream the evidence object
			d.StreamListItem(ctx, evidence)

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

//// TRANSFORM FUNCTIONS

// getRelatedControlNames extracts the control names from the evidence object
func getRelatedControlNames(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	item := d.HydrateItem
	evidence, ok := item.(*model.Evidence)
	if !ok {
		return nil, nil
	}

	var controlNames []string
	for _, control := range evidence.RelatedControls {
		controlNames = append(controlNames, control.Name)
	}

	return controlNames, nil
}
