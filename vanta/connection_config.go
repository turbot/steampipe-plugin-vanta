package vanta

import (
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

type vantaConfig struct {
	ApiToken  *string `hcl:"api_token"`
	SessionId *string `hcl:"session_id"`
}

// ConfigInstance returns an instance of a connection config struct
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
