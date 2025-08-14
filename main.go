package main

import (
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-vanta/v2/vanta"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{PluginFunc: vanta.Plugin})
}
