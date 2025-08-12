package rest_api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/turbot/steampipe-plugin-vanta/rest_api/model"
)

// ListVendors retrieves a paginated list of vendors from Vanta
func (c *RestClient) ListVendors(ctx context.Context, options *model.ListVendorsOptions) (*model.ListVendorsOutput, error) {
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

	resp, err := c.makeRequest(ctx, "GET", "/v1/vendors", params)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}

	respBodyBytes, err := c.readResponseBody(resp)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var result *model.ListVendorsOutput
	if err = json.Unmarshal(respBodyBytes, &result); err != nil {
		return nil, fmt.Errorf("failed to JSON-decode response body: %w", err)
	}

	return result, nil
}

// GetVendorByID retrieves a specific vendor by its ID
func (c *RestClient) GetVendorByID(ctx context.Context, id string) (*model.Vendor, error) {
	if id == "" {
		return nil, fmt.Errorf("vendor ID cannot be empty")
	}

	resp, err := c.makeRequest(ctx, "GET", fmt.Sprintf("/v1/vendors/%s", id), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}

	respBodyBytes, err := c.readResponseBody(resp)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var vendor *model.Vendor
	if err = json.Unmarshal(respBodyBytes, &vendor); err != nil {
		return nil, fmt.Errorf("failed to JSON-decode response body: %w", err)
	}

	return vendor, nil
}
