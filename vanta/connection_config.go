package vanta

import (
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

type vantaConfig struct {
	ClientID     *string `hcl:"client_id"`
	ClientSecret *string `hcl:"client_secret"`
	AccessToken  *string `hcl:"access_token"`
	RefreshToken *string `hcl:"refresh_token"`

	ApiToken  *string `hcl:"api_token"`
	SessionId *string `hcl:"session_id"` // This is the connect.sid cookie from a logged in Vanta browser session. Required to access tables that are using the deprecated https://app.vanta.com/graphql endpoint
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
