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
			Schema:      ConfigSchema,
		},
		// DefaultGetConfig: &plugin.GetConfig{
		// 	ShouldIgnoreError: errors.NotFoundError,
		// },
		DefaultTransform: transform.FromCamel().Transform(transform.NullIfZeroValue),
		TableMap: map[string]*plugin.Table{
			"vanta_monitor": tableVantaMonitor(ctx),
			"vanta_policy":  tableVantaPolicy(ctx),
			"vanta_user":    tableVantaUser(ctx),
		},
	}
	return p
}
