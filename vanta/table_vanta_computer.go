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

func tableVantaComputer(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "vanta_computer",
		Description: "Vanta Computer",
		List: &plugin.ListConfig{
			Hydrate: listVantaComputers,
		},
		Columns: []*plugin.Column{
			{Name: "owner_name", Type: proto.ColumnType_STRING, Description: "The name of the workstation owner.", Transform: transform.FromField("Owner.DisplayName")},
			{Name: "id", Type: proto.ColumnType_STRING, Description: "A unique Vanta generated identifier of the computer.", Transform: transform.FromField("Data.Id")},
			{Name: "serial_number", Type: proto.ColumnType_STRING, Description: "The serial number of the workstation.", Transform: transform.FromField("Data.SerialNumber")},
			{Name: "agent_version", Type: proto.ColumnType_STRING, Description: "The Vanta agent version.", Transform: transform.FromField("Data.AgentVersion")},
			{Name: "os_version", Type: proto.ColumnType_STRING, Description: "The OS version of the workstation.", Transform: transform.FromField("Data.OsVersion")},
			{Name: "is_password_manager_installed", Type: proto.ColumnType_BOOL, Description: "If true, a password manager is installed in the workstation.", Transform: transform.FromField("Data.IsPasswordManagerInstalled")},
			{Name: "is_encrypted", Type: proto.ColumnType_BOOL, Description: "If true, the workstation's hard drive is encrypted.", Transform: transform.FromField("Data.IsEncrypted")},
			{Name: "has_screen_lock", Type: proto.ColumnType_BOOL, Description: "If true, the workstation has a screen lock configured.", Transform: transform.FromField("Data.HasScreenLock")},
			{Name: "hostname", Type: proto.ColumnType_STRING, Description: "The hostname of the workstation.", Transform: transform.FromField("Data.Hostname")},
			{Name: "host_identifier", Type: proto.ColumnType_STRING, Description: "The host identifier of the workstation.", Transform: transform.FromField("Data.HostIdentifier")},
			{Name: "last_ping", Type: proto.ColumnType_TIMESTAMP, Description: "The time when the workstation was last scanned by the Vanta agent.", Transform: transform.FromField("Data.LastPing").Transform(transform.NullIfZeroValue).Transform(transform.UnixToTimestamp)},
			{Name: "num_browser_extensions", Type: proto.ColumnType_INT, Description: "The number of browser extensions installed in the workstation.", Transform: transform.FromField("Data.NumBrowserExtensions")},
			{Name: "endpoint_applications", Type: proto.ColumnType_JSON, Description: "A list of applications installed on the device.", Hydrate: listVantaComputerApplications, Transform: transform.FromValue()},
			{Name: "installed_av_programs", Type: proto.ColumnType_JSON, Description: "A list of anti-virus programs installed in the workstation.", Transform: transform.FromField("Data.InstalledAvPrograms")},
			{Name: "installed_password_managers", Type: proto.ColumnType_JSON, Description: "A list of password managers installed in the workstation.", Transform: transform.FromField("Data.InstalledPasswordManagers")},
			{Name: "unsupported_reasons", Type: proto.ColumnType_JSON, Description: "Specifies the reason for unmonitored computers."},
			{Name: "owner_id", Type: proto.ColumnType_STRING, Description: "A unique identifier of the owner of the workstation.", Transform: transform.FromField("Owner.Id")},
			{Name: "organization_name", Type: proto.ColumnType_STRING, Description: "The name of the organization."},
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
			result.OrganizationName = query.Organization.Name
			result.Owner = api.WorkstationOwner{
				DisplayName: user.DisplayName,
				Id:          user.Id,
			}
			result.UnsupportedReasons = workstation.UnsupportedReasons
			d.StreamListItem(ctx, result)

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func listVantaComputerApplications(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	computer := h.Item.(api.Workstation)

	// If the computer is unmonitored/unsupported, the Vanta agent does not monitor
	// the application installed on that device.
	// Hence, the API returns an error if a request is being made for an unmonitored device.
	// Return nil, if the device is unmonitored / unsupported.
	var unsupportedReasons api.WorkstationUnsupportedReasons
	if computer.UnsupportedReasons != unsupportedReasons {
		return nil, nil
	}

	// Create client
	conn, err := getVantaAppClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("vanta_computer.listVantaComputerApplications", "connection_error", err)
		return nil, err
	}

	// As of Feb 7, 2023, the query doesn't provide the paging information
	query, err := api.ListEndpointApplications(context.Background(), conn, computer.Data.Id)
	if err != nil {
		plugin.Logger(ctx).Error("vanta_computer.listVantaComputerApplications", "query_error", err)
		return nil, err
	}

	var result []api.AppDetails
	for _, endpointIds := range query.Organization.OsqueryEndpointsByIds {
		
		for _, app := range endpointIds.Data.ApplicationData {

			// Remove all entries with empty data
			if strings.ReplaceAll(app.Name, " ", "") == "" {
				continue
			}
			result = append(result, app)
		}
	}

	return result, nil
}
