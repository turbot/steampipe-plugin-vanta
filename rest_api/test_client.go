package rest_api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"

	"github.com/turbot/steampipe-plugin-vanta/v2/rest_api/model"
)

// ListTests retrieves a paginated list of comprehensive tests from Vanta
func (c *RestClient) ListTests(ctx context.Context, options *model.ListTestsOptions) (*model.TestResults, error) {
	// Build URL with query parameters
	params := url.Values{}

	if options != nil {
		if options.PageSize > 0 {
			params.Set("pageSize", strconv.Itoa(options.PageSize))
		}
		if options.PageCursor != "" {
			params.Set("pageCursor", options.PageCursor)
		}
		if options.StatusFilter != "" {
			params.Set("statusFilter", options.StatusFilter)
		}
		if options.FrameworkFilter != "" {
			params.Set("frameworkFilter", options.FrameworkFilter)
		}
		if options.IntegrationFilter != "" {
			params.Set("integrationFilter", options.IntegrationFilter)
		}
		if options.ControlFilter != "" {
			params.Set("controlFilter", options.ControlFilter)
		}
		if options.OwnerFilter != "" {
			params.Set("ownerFilter", options.OwnerFilter)
		}
		if options.CategoryFilter != "" {
			params.Set("categoryFilter", options.CategoryFilter)
		}
		if options.IsInRollout != nil {
			params.Set("isInRollout", strconv.FormatBool(*options.IsInRollout))
		}
	}

	resp, err := c.makeRequest(ctx, "GET", "/v1/tests", params)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}

	respBodyBytes, err := c.readResponseBody(resp)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var result *model.TestResults
	if err = json.Unmarshal(respBodyBytes, &result); err != nil {
		return nil, fmt.Errorf("failed to JSON-decode response body: %w", err)
	}

	return result, nil
}

// GetTestByID retrieves a specific comprehensive test by its ID
func (c *RestClient) GetTestByID(ctx context.Context, id string) (*model.Test, error) {
	if id == "" {
		return nil, fmt.Errorf("test ID cannot be empty")
	}

	resp, err := c.makeRequest(ctx, "GET", fmt.Sprintf("/v1/tests/%s", id), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}

	respBodyBytes, err := c.readResponseBody(resp)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var test *model.Test
	if err = json.Unmarshal(respBodyBytes, &test); err != nil {
		return nil, fmt.Errorf("failed to JSON-decode response body: %w", err)
	}

	return test, nil
}
