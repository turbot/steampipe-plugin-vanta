package vanta

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
	"github.com/turbot/steampipe-plugin-vanta/restapi/model"
)

//// TABLE DEFINITION

func tableVantaComputer(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "vanta_computer",
		Description: "Vanta Computer",
		List: &plugin.ListConfig{
			Hydrate: listVantaComputers,
		},
		Get: &plugin.GetConfig{
			Hydrate:    getVantaComputer,
			KeyColumns: plugin.SingleColumn("id"),
		},
		Columns: []*plugin.Column{
			// Available columns from REST API
			{Name: "id", Type: proto.ColumnType_STRING, Transform: transform.FromField("ID"), Description: "A unique identifier of the computer."},
			{Name: "integration_id", Type: proto.ColumnType_STRING, Transform: transform.FromField("IntegrationID"), Description: "The integration ID associated with this computer."},
			{Name: "last_check_date", Type: proto.ColumnType_TIMESTAMP, Transform: transform.FromField("LastCheckDate"), Description: "The last time this computer was checked."},
			{Name: "operating_system", Type: proto.ColumnType_JSON, Transform: transform.FromField("OperatingSystem"), Description: "Operating system information including type and version."},
			{Name: "owner", Type: proto.ColumnType_JSON, Transform: transform.FromField("Owner"), Description: "Owner information including ID, email, and display name."},
			{Name: "serial_number", Type: proto.ColumnType_STRING, Transform: transform.FromField("SerialNumber"), Description: "The serial number of the computer."},
			{Name: "udid", Type: proto.ColumnType_STRING, Transform: transform.FromField("UDID"), Description: "The unique device identifier."},
			{Name: "screenlock", Type: proto.ColumnType_JSON, Transform: transform.FromField("Screenlock"), Description: "Screenlock security check results."},
			{Name: "disk_encryption", Type: proto.ColumnType_JSON, Transform: transform.FromField("DiskEncryption"), Description: "Disk encryption security check results."},
			{Name: "password_manager", Type: proto.ColumnType_JSON, Transform: transform.FromField("PasswordManager"), Description: "Password manager security check results."},
			{Name: "antivirus_installation", Type: proto.ColumnType_JSON, Transform: transform.FromField("AntivirusInstallation"), Description: "Antivirus installation security check results."},

			// Derived columns (available via transform from REST API)
			{Name: "owner_name", Type: proto.ColumnType_STRING, Transform: transform.From(getOwnerName), Description: "The name of the workstation owner."},
			{Name: "owner_id", Type: proto.ColumnType_STRING, Transform: transform.From(getOwnerID), Description: "A unique identifier of the owner of the workstation."},
			{Name: "os_version", Type: proto.ColumnType_STRING, Transform: transform.From(getOSVersion), Description: "The OS version of the workstation."},
			{Name: "has_screen_lock", Type: proto.ColumnType_BOOL, Transform: transform.From(getScreenlockStatus), Description: "If true, the workstation has a screen lock configured."},
			{Name: "is_encrypted", Type: proto.ColumnType_BOOL, Transform: transform.From(getDiskEncryptionStatus), Description: "If true, the workstation's hard drive is encrypted."},
			{Name: "is_password_manager_installed", Type: proto.ColumnType_BOOL, Transform: transform.From(getPasswordManagerStatus), Description: "If true, a password manager is installed in the workstation."},

			// Deprecated columns (not available in REST API)
			{Name: "agent_version", Type: proto.ColumnType_STRING, Description: "[DEPRECATED] The Vanta agent version."},
			{Name: "hostname", Type: proto.ColumnType_STRING, Description: "[DEPRECATED] The hostname of the workstation."},
			{Name: "host_identifier", Type: proto.ColumnType_STRING, Description: "[DEPRECATED] The host identifier of the workstation."},
			{Name: "last_ping", Type: proto.ColumnType_TIMESTAMP, Description: "[DEPRECATED] The time when the workstation was last scanned by the Vanta agent."},
			{Name: "num_browser_extensions", Type: proto.ColumnType_INT, Description: "[DEPRECATED] The number of browser extensions installed in the workstation."},
			{Name: "endpoint_applications", Type: proto.ColumnType_JSON, Description: "[DEPRECATED] A list of applications installed on the device."},
			{Name: "installed_av_programs", Type: proto.ColumnType_JSON, Description: "[DEPRECATED] A list of anti-virus programs installed in the workstation."},
			{Name: "installed_password_managers", Type: proto.ColumnType_JSON, Description: "[DEPRECATED] A list of password managers installed in the workstation."},
			{Name: "unsupported_reasons", Type: proto.ColumnType_JSON, Description: "[DEPRECATED] Specifies the reason for unmonitored computers."},
			{Name: "organization_name", Type: proto.ColumnType_STRING, Description: "[DEPRECATED] The name of the organization."},
		},
	}
}

//// LIST FUNCTION

func listVantaComputers(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// Create REST client
	client, err := getClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("vanta_computer.listVantaComputers", "connection_error", err)
		return nil, err
	}

	// Default page limit
	pageLimit := 100

	// Adjust page limit if query limit is smaller
	if d.QueryContext.Limit != nil && int(*d.QueryContext.Limit) < pageLimit {
		pageLimit = int(*d.QueryContext.Limit)
	}

	options := &model.ListComputersOptions{
		Limit:  pageLimit,
		Cursor: "",
	}

	for {
		result, err := client.ListComputers(ctx, options)
		if err != nil {
			plugin.Logger(ctx).Error("vanta_computer.listVantaComputers", "query_error", err)
			return nil, err
		}

		for _, computer := range result.Results.Data {
			// Stream the raw Computer object
			d.StreamListItem(ctx, computer)

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

func getVantaComputer(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	id := d.EqualsQualString("id")
	if id == "" {
		return nil, nil
	}

	// Create REST client
	client, err := getClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("vanta_computer.getVantaComputer", "connection_error", err)
		return nil, err
	}

	computer, err := client.GetComputerByID(ctx, id)
	if err != nil {
		plugin.Logger(ctx).Error("vanta_computer.getVantaComputer", "query_error", err)
		return nil, err
	}

	if computer == nil {
		return nil, nil
	}

	return computer, nil
}

//// TRANSFORM FUNCTIONS

// getOwnerName extracts the owner display name from the computer object
func getOwnerName(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	item := d.HydrateItem
	computer, ok := item.(*model.Computer)
	if !ok {
		return nil, nil
	}

	if computer.Owner != nil {
		return computer.Owner.DisplayName, nil
	}
	return nil, nil
}

// getOwnerID extracts the owner ID from the computer object
func getOwnerID(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	item := d.HydrateItem
	computer, ok := item.(*model.Computer)
	if !ok {
		return nil, nil
	}

	if computer.Owner != nil {
		return computer.Owner.ID, nil
	}
	return nil, nil
}

// getOSVersion extracts the operating system version from the computer object
func getOSVersion(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	item := d.HydrateItem
	computer, ok := item.(*model.Computer)
	if !ok {
		return nil, nil
	}

	if computer.OperatingSystem != nil {
		return computer.OperatingSystem.Version, nil
	}
	return nil, nil
}

// getScreenlockStatus determines if a computer has screen lock based on security check outcome
func getScreenlockStatus(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	item := d.HydrateItem
	computer, ok := item.(*model.Computer)
	if !ok {
		return false, nil
	}

	if computer.Screenlock != nil && computer.Screenlock.Outcome == "PASS" {
		return true, nil
	}
	return false, nil
}

// getDiskEncryptionStatus determines if a computer has disk encryption based on security check outcome
func getDiskEncryptionStatus(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	item := d.HydrateItem
	computer, ok := item.(*model.Computer)
	if !ok {
		return false, nil
	}

	if computer.DiskEncryption != nil && computer.DiskEncryption.Outcome == "PASS" {
		return true, nil
	}
	return false, nil
}

// getPasswordManagerStatus determines if a computer has password manager based on security check outcome
func getPasswordManagerStatus(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	item := d.HydrateItem
	computer, ok := item.(*model.Computer)
	if !ok {
		return false, nil
	}

	if computer.PasswordManager != nil && computer.PasswordManager.Outcome == "PASS" {
		return true, nil
	}
	return false, nil
}
