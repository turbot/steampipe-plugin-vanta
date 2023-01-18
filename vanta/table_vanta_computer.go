package vanta

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
	"github.com/turbot/steampipe-plugin-vanta/api"
)

//// TABLE DEFINITION

func tableVantaComputer(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "vanta_computer",
		Description: "Vanta Computer",
		List: &plugin.ListConfig{
			Hydrate: listVantaComputers,
		},
		Columns: []*plugin.Column{
			{Name: "owner_name", Type: proto.ColumnType_STRING, Description: "", Transform: transform.FromField("Owner.DisplayName")},
			{Name: "serial_number", Type: proto.ColumnType_STRING, Description: "", Transform: transform.FromField("Data.SerialNumber")},
			{Name: "agent_version", Type: proto.ColumnType_STRING, Description: "", Transform: transform.FromField("Data.AgentVersion")},
			{Name: "os_version", Type: proto.ColumnType_STRING, Description: "", Transform: transform.FromField("Data.OsVersion")},

			{Name: "is_password_manager_installed", Type: proto.ColumnType_BOOL, Description: "", Transform: transform.FromField("Data.IsPasswordManagerInstalled")},
			{Name: "is_encrypted", Type: proto.ColumnType_BOOL, Description: "", Transform: transform.FromField("Data.IsEncrypted")},
			{Name: "has_screen_lock", Type: proto.ColumnType_BOOL, Description: "", Transform: transform.FromField("Data.HasScreenLock")},

			{Name: "hostname", Type: proto.ColumnType_STRING, Description: "", Transform: transform.FromField("Data.Hostname")},
			{Name: "host_identifier", Type: proto.ColumnType_STRING, Description: "", Transform: transform.FromField("Data.HostIdentifier")},
			{Name: "last_ping", Type: proto.ColumnType_STRING, Description: "", Transform: transform.FromField("Data.LastPing")},

			{Name: "num_browser_extensions", Type: proto.ColumnType_INT, Description: "", Transform: transform.FromField("Data.NumBrowserExtensions")},
			{Name: "installed_av_programs", Type: proto.ColumnType_JSON, Description: "", Transform: transform.FromField("Data.InstalledAvPrograms")},
			{Name: "installed_password_managers", Type: proto.ColumnType_JSON, Description: "", Transform: transform.FromField("Data.InstalledPasswordManagers")},
		},
	}
}

//// LIST FUNCTION

func listVantaComputers(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// Create client
	conn, err := getVantaAppClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("vanta_computer.listVantaComputers", "connection_error", err)
		return nil, err
	}

	// As of Jan 13, 2023, the query doesn't provide the paging information
	query, err := api.ListWorkstations(context.Background(), conn)
	if err != nil {
		plugin.Logger(ctx).Error("vanta_computer.listVantaComputers", "query_error", err)
		return nil, err
	}

	for _, user := range query.Organization.Users {
		for _, workstation := range user.Workstations {
			result := workstation
			result.Owner = api.WorkstationOwner{
				DisplayName: user.DisplayName,
				Email:       user.Email,
			}
			d.StreamListItem(ctx, result)
		}
	}

	return nil, nil
}
