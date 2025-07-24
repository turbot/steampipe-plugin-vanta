package restapi

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/turbot/steampipe-plugin-vanta/restapi/model"
)

// ListComputers retrieves a paginated list of computers from Vanta
func (c *RestClient) ListComputers(ctx context.Context, options *model.ListComputersOptions) (*model.ListComputersOutput, error) {
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

	resp, err := c.makeRequest(ctx, "GET", "/v1/monitored-computers", params)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}

	respBodyBytes, err := c.readResponseBody(resp)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var result *model.ListComputersOutput
	if err = json.Unmarshal(respBodyBytes, &result); err != nil {
		return nil, fmt.Errorf("failed to JSON-decode response body: %w", err)
	}

	return result, nil
}

// GetComputerByID retrieves a specific computer by its ID
func (c *RestClient) GetComputerByID(ctx context.Context, id string) (*model.Computer, error) {
	if id == "" {
		return nil, fmt.Errorf("computer ID cannot be empty")
	}

	resp, err := c.makeRequest(ctx, "GET", fmt.Sprintf("/v1/monitored-computers/%s", id), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}

	respBodyBytes, err := c.readResponseBody(resp)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var computer *model.Computer
	if err = json.Unmarshal(respBodyBytes, &computer); err != nil {
		return nil, fmt.Errorf("failed to JSON-decode response body: %w", err)
	}

	return computer, nil
}
