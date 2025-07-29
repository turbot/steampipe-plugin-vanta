package vanta

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

const pluginName = "steampipe-plugin-vanta"

// Plugin creates this (vanta) plugin
func Plugin(ctx context.Context) *plugin.Plugin {
	p := &plugin.Plugin{
		Name: pluginName,
		ConnectionConfigSchema: &plugin.ConnectionConfigSchema{
			NewInstance: ConfigInstance,
		},
		DefaultShouldIgnoreError: isNotFoundError([]string{"not found"}),
		DefaultTransform:         transform.FromCamel().Transform(transform.NullIfZeroValue),
		TableMap: map[string]*plugin.Table{
			"vanta_computer":    tableVantaComputer(ctx),
			"vanta_evidence":    tableVantaEvidence(ctx),
			"vanta_group":       tableVantaGroup(ctx),
			"vanta_integration": tableVantaIntegration(ctx),
			"vanta_monitor":     tableVantaMonitor(ctx),
			"vanta_policy":      tableVantaPolicy(ctx),
			"vanta_user":        tableVantaUser(ctx),
			"vanta_vendor":      tableVantaVendor(ctx),
		},
	}
	return p
}
