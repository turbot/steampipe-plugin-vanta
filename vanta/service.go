package vanta

import (
	"context"
	"fmt"
	"os"

	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-vanta/api"
)

// getClient:: returns vanta client after authentication
func getClient(ctx context.Context, d *plugin.QueryData) (*api.Client, error) {
	// Load connection from cache, which preserves throttling protection etc
	cacheKey := "vanta"
	if cachedData, ok := d.ConnectionManager.Cache.Get(cacheKey); ok {
		return cachedData.(*api.Client), nil
	}

	// Get the config
	vantaConfig := GetConfig(d.Connection)

	/*
		precedence of credentials:
		- Credentials set in config
		- VANTA_API_TOKEN env var
	*/
	var token string
	token = os.Getenv("VANTA_API_TOKEN")

	if vantaConfig.ApiToken != nil {
		token = *vantaConfig.ApiToken
	}

	// Return if no credential specified
	if token == "" {
		return nil, fmt.Errorf("api_token must be configured")
	}

	// Start with an empty Vanta config
	config := api.ClientConfig{ApiToken: vantaConfig.ApiToken}

	// Create the client
	client, err := api.CreateClient(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("error creating client: %s", err.Error())
	}

	// Save to cache
	d.ConnectionManager.Cache.Set(cacheKey, client)

	return client, nil
}

// getVantaAppClient:: returns vanta client using an endpoint https://app.vanta.com/graphql
//
// Vanta has 2 separate endpoint to access the resources
// The public endpoint - https://app.vanta.com/graphql; and the other is
// https://app.vanta.com/graphql, which is being used in the console.
//
// The public one required the users' personal access token to authenticate the request; whereas
// the endpoint https://app.vanta.com/graphql requires a session id which is created when the user logged in to the console
// and valid until the session is expired.
func getVantaAppClient(ctx context.Context, d *plugin.QueryData) (*api.Client, error) {
	// Load connection from cache, which preserves throttling protection etc
	cacheKey := "vanta-app"
	if cachedData, ok := d.ConnectionManager.Cache.Get(cacheKey); ok {
		return cachedData.(*api.Client), nil
	}

	// Get the config
	vantaConfig := GetConfig(d.Connection)

	var sessionId string
	if vantaConfig.SessionId != nil {
		sessionId = *vantaConfig.SessionId
	}

	// Return if no credential specified
	if sessionId == "" {
		return nil, fmt.Errorf("session_id must be configured")
	}

	// Start with an empty Vanta config
	config := api.ClientConfig{SessionId: &sessionId}

	// Create the client
	client, err := api.CreateAppClient(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("error creating client: %s", err.Error())
	}

	// Save to cache
	d.ConnectionManager.Cache.Set(cacheKey, client)

	return client, nil
}
