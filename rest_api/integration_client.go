package rest_api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/turbot/steampipe-plugin-vanta/v2/rest_api/model"
)

// ListConnectedIntegrations retrieves a paginated list of integrations from Vanta
func (c *RestClient) ListConnectedIntegrations(ctx context.Context, options *model.ListIntegrationsOptions) (*model.ListIntegrationsOutput, error) {
	// Build URL with query parameters
	params := url.Values{}

	if options != nil {
		if options.Limit > 0 {
			params.Set("pageSize", fmt.Sprintf("%d", options.Limit))
		}
		if options.Cursor != "" {
			params.Set("pageCursor", options.Cursor)
		}
	}

	resp, err := c.makeRequest(ctx, "GET", "/v1/integrations", params)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}

	respBodyBytes, err := c.readResponseBody(resp)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var result *model.ListIntegrationsOutput
	if err = json.Unmarshal(respBodyBytes, &result); err != nil {
		return nil, fmt.Errorf("failed to JSON-decode response body: %w", err)
	}

	return result, nil
}

// GetIntegrationByID retrieves a specific integration by its ID
func (c *RestClient) GetIntegrationByID(ctx context.Context, id string) (*model.Integration, error) {
	if id == "" {
		return nil, fmt.Errorf("integration ID cannot be empty")
	}

	resp, err := c.makeRequest(ctx, "GET", fmt.Sprintf("/v1/integrations/%s", id), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}

	respBodyBytes, err := c.readResponseBody(resp)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var integration *model.Integration
	if err = json.Unmarshal(respBodyBytes, &integration); err != nil {
		return nil, fmt.Errorf("failed to JSON-decode response body: %w", err)
	}

	return integration, nil
}
