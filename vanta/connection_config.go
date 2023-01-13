package vanta

import (
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/schema"
)

type vantaConfig struct {
	ApiToken *string `cty:"api_token"`
}

var ConfigSchema = map[string]*schema.Attribute{
	"api_token": {
		Type: schema.TypeString,
	},
}

func ConfigInstance() interface{} {
	return &vantaConfig{}
}

// GetConfig :: retrieve and cast connection config from query data
func GetConfig(connection *plugin.Connection) vantaConfig {
	if connection == nil || connection.Config == nil {
		return vantaConfig{}
	}
	config, _ := connection.Config.(vantaConfig)
	return config
}
