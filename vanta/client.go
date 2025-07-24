package vanta

import (
	"context"
	"fmt"

	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-vanta/restapi"
)

// CreateRestClient creates a new REST API client using our restapi package
func CreateRestClient(ctx context.Context, queryData *plugin.QueryData) (restapi.Vanta, error) {
	config := GetConfig(queryData.Connection)

	var options []restapi.Option
	options = append(options, restapi.WithScopes(restapi.ScopeAllRead))

	// If OAuth credentials are provided, use them
	if config.ClientID != nil && config.ClientSecret != nil {
		options = append(options, restapi.WithOAuthCredentials(*config.ClientID, *config.ClientSecret))
	} else if config.AccessToken != nil {
		// If access token is provided, use it
		options = append(options, restapi.WithToken(*config.AccessToken))
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

	return client, nil
}

// Legacy function - kept for backward compatibility with external SDK
func CreateClient(ctx context.Context, queryData *plugin.QueryData) (interface{}, error) {
	// For now, redirect to the REST client
	// This maintains the function signature while using our new implementation
	return CreateRestClient(ctx, queryData)
}
