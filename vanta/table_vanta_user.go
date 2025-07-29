package vanta

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
	"github.com/turbot/steampipe-plugin-vanta/restapi/model"
)

//// TABLE DEFINITION

func tableVantaUser(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "vanta_user",
		Description: "Vanta User",
		List: &plugin.ListConfig{
			Hydrate: listVantaUsers,
			KeyColumns: plugin.KeyColumnSlice{
				{Name: "employment_status", Require: plugin.Optional},
			},
		},
		Get: &plugin.GetConfig{
			Hydrate:    getVantaUser,
			KeyColumns: plugin.SingleColumn("id"),
		},
		Columns: []*plugin.Column{
			{Name: "display_name", Type: proto.ColumnType_STRING, Transform: transform.FromField("Name.Display"), Description: "The display name of the user."},
			{Name: "id", Type: proto.ColumnType_STRING, Transform: transform.FromField("ID"), Description: "A unique identifier of the user."},
			{Name: "email", Type: proto.ColumnType_STRING, Transform: transform.FromField("EmailAddress"), Description: "The email of the user."},
			{Name: "employment_status", Type: proto.ColumnType_STRING, Transform: transform.FromField("Employment.Status"), Description: "The current employment status of the user."},
			{Name: "job_title", Type: proto.ColumnType_STRING, Transform: transform.FromField("Employment.JobTitle"), Description: "The job title of the user."},
			{Name: "task_status", Type: proto.ColumnType_STRING, Transform: transform.FromField("TasksSummary.Status"), Description: "The security task status of the user."},
			{Name: "start_date", Type: proto.ColumnType_TIMESTAMP, Transform: transform.FromField("Employment.StartDate"), Description: "The employment start date of the user."},
			{Name: "end_date", Type: proto.ColumnType_TIMESTAMP, Transform: transform.FromField("Employment.EndDate"), Description: "The employment end date of the user."},
			{Name: "family_name", Type: proto.ColumnType_STRING, Transform: transform.FromField("Name.Last"), Description: "The family name of the user."},
			{Name: "given_name", Type: proto.ColumnType_STRING, Transform: transform.FromField("Name.First"), Description: "The given name of the user."},
			{Name: "is_active", Type: proto.ColumnType_BOOL, Transform: transform.From(getIsActiveStatus), Description: "If true, the user is active."},
			{Name: "group_ids", Type: proto.ColumnType_JSON, Description: "List of group IDs the user belongs to."},
			{Name: "employment", Type: proto.ColumnType_JSON, Description: "Employment information including job title and dates."},
			{Name: "name", Type: proto.ColumnType_JSON, Description: "Name information including display, first, and last name."},
			{Name: "sources", Type: proto.ColumnType_JSON, Description: "Information about data sources for this user."},
			{Name: "tasks_summary", Type: proto.ColumnType_JSON, Description: "Summary of security task completion status."},
			{Name: "is_from_scan", Type: proto.ColumnType_BOOL, Description: "[DEPRECATED] If true, the user was discovered by the security scan."},
			{Name: "needs_employee_digest_reminder", Type: proto.ColumnType_BOOL, Description: "[DEPRECATED] If true, user will get an email digest of their incomplete security tasks."},
			{Name: "is_not_human", Type: proto.ColumnType_BOOL, Description: "[DEPRECATED] If true, the resource is not a human."},
		},
	}
}

//// LIST FUNCTION

func listVantaUsers(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// Create REST client
	client, err := getClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("vanta_user.listVantaUsers", "connection_error", err)
		return nil, err
	}

	maxLimit := int32(100)
	if d.QueryContext.Limit != nil {
		limit := int32(*d.QueryContext.Limit)
		if limit < maxLimit {
			maxLimit = limit
		}
	}

	// Check for employment status filter
	// employmentStatusFilter := d.EqualsQualString("employment_status")

	options := &model.ListPeopleOptions{
		Limit:  int(maxLimit),
		Cursor: "",
	}

	for {
		result, err := client.ListPeople(ctx, options)
		if err != nil {
			plugin.Logger(ctx).Error("vanta_user.listVantaUsers", "api_error", err)
			return nil, err
		}

		for _, person := range result.Results.Data {
			// Stream the raw Person object
			d.StreamListItem(ctx, person)

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

func getVantaUser(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	id := d.EqualsQualString("id")
	if id == "" {
		return nil, nil
	}

	// Create REST client
	client, err := getClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("vanta_user.getVantaUser", "connection_error", err)
		return nil, err
	}

	person, err := client.GetPersonByID(ctx, id)
	if err != nil {
		plugin.Logger(ctx).Error("vanta_user.getVantaUser", "api_error", err)
		return nil, err
	}

	if person == nil {
		return nil, nil
	}

	// Return the raw Person object
	return person, nil
}

//// HELPER FUNCTIONS

// getIsActiveStatus determines if a user is active based on employment status
func getIsActiveStatus(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	item := d.HydrateItem

	person, ok := item.(*model.Person)
	if !ok {
		plugin.Logger(ctx).Error("getIsActiveStatus", "casting_error", "HydrateItem is not *model.Person")
		return false, nil
	}

	if person.Employment != nil && person.Employment.Status != nil {
		return *person.Employment.Status == model.EmploymentStatusCurrent, nil
	}

	return false, nil
}
