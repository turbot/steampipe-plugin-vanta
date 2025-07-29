package vanta

import (
	"context"
	"fmt"

	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-vanta/restapi"
)

// getClient:: returns vanta client after authentication
func getClient(ctx context.Context, d *plugin.QueryData) (restapi.Vanta, error) {
	// Load connection from cache, which preserves throttling protection etc
	cacheKey := "vanta"
	if cachedData, ok := d.ConnectionManager.Cache.Get(cacheKey); ok {
		return cachedData.(restapi.Vanta), nil
	}

	// Get the config
	vantaConfig := GetConfig(d.Connection)

	var options []restapi.Option
	options = append(options, restapi.WithScopes(restapi.ScopeAllRead))

	// If OAuth credentials are provided, use them
	if vantaConfig.ClientID != nil && vantaConfig.ClientSecret != nil {
		options = append(options, restapi.WithOAuthCredentials(*vantaConfig.ClientID, *vantaConfig.ClientSecret))
	} else if vantaConfig.AccessToken != nil {
		// If access token is provided, use it
		options = append(options, restapi.WithToken(*vantaConfig.AccessToken))
	} else {
		// Fallback: this maintains backward compatibility with existing configurations
		plugin.Logger(ctx).Warn("vanta.CreateRestClient", "warning", "No OAuth credentials or access token found, REST client requires authentication")
		return nil, fmt.Errorf("authentication required: provide either client_id/client_secret or access_token in connection config")
	}

	client, err := restapi.New(ctx, options...)
	if err != nil {
		plugin.Logger(ctx).Error("vanta.CreateRestClient", "error", err)
		return nil, err
	}

	// Save to cache
	d.ConnectionManager.Cache.Set(cacheKey, client)

	return client, nil
}
