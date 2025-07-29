package vanta

import (
	"context"
	"fmt"
	"strings"

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

	// Validate configuration
	if err := validateConfig(vantaConfig); err != nil {
		plugin.Logger(ctx).Error("vanta.getClient", "config_validation_error", err)
		return nil, err
	}

	var options []restapi.Option
	options = append(options, restapi.WithScopes(restapi.ScopeAllRead))

	// If OAuth credentials are provided, use them
	if vantaConfig.ClientID != nil && vantaConfig.ClientSecret != nil {
		options = append(options, restapi.WithOAuthCredentials(*vantaConfig.ClientID, *vantaConfig.ClientSecret))
	} else if vantaConfig.AccessToken != nil {
		// If access token is provided, use it
		options = append(options, restapi.WithToken(*vantaConfig.AccessToken))
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

// validateConfig validates the Vanta configuration and returns appropriate errors
func validateConfig(config vantaConfig) error {
	// Check if OAuth credentials are provided
	hasOAuthCredentials := config.ClientID != nil && config.ClientSecret != nil
	hasPartialOAuthCredentials := (config.ClientID != nil && config.ClientSecret == nil) || (config.ClientID == nil && config.ClientSecret != nil)
	hasAccessToken := config.AccessToken != nil

	// Validate that OAuth credentials are complete if provided
	if hasPartialOAuthCredentials {
		if config.ClientID != nil && config.ClientSecret == nil {
			return fmt.Errorf("invalid configuration: client_secret is required when client_id is provided")
		}
		if config.ClientID == nil && config.ClientSecret != nil {
			return fmt.Errorf("invalid configuration: client_id is required when client_secret is provided")
		}
	}

	// Validate that at least one authentication method is provided
	if !hasOAuthCredentials && !hasAccessToken {
		return fmt.Errorf("authentication required: provide either OAuth credentials (client_id and client_secret) or access_token in connection config")
	}

	// Validate credential formats
	if hasOAuthCredentials {
		if err := validateCredentialFormat("client_id", *config.ClientID, "vci_"); err != nil {
			return err
		}
		if err := validateCredentialFormat("client_secret", *config.ClientSecret, "vcs_"); err != nil {
			return err
		}
	}

	if hasAccessToken {
		if err := validateCredentialFormat("access_token", *config.AccessToken, "vat_"); err != nil {
			return err
		}
	}
	return nil
}

// validateCredentialFormat validates that credentials follow the expected format
func validateCredentialFormat(credType, credential, expectedPrefix string) error {
	if credential == "" {
		return fmt.Errorf("invalid configuration: %s cannot be empty", credType)
	}

	if !strings.HasPrefix(credential, expectedPrefix) {
		return fmt.Errorf("invalid configuration: %s should start with '%s'", credType, expectedPrefix)
	}

	// Basic length validation (Vanta credentials are typically longer than just the prefix)
	if len(credential) <= len(expectedPrefix)+5 {
		return fmt.Errorf("invalid configuration: %s appears to be too short or invalid", credType)
	}

	return nil
}
