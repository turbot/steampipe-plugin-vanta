package restapi

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"

	"github.com/turbot/steampipe-plugin-vanta/restapi/model"
)

// ListMonitors retrieves a paginated list of monitors/tests from Vanta
func (c *RestClient) ListMonitors(ctx context.Context, options *model.ListMonitorsOptions) (*model.MonitorResults, error) {
	// Build URL with query parameters
	params := url.Values{}

	if options != nil {
		if options.Limit > 0 {
			params.Set("pageSize", strconv.Itoa(options.Limit))
		}
		if options.Cursor != "" {
			params.Set("pageCursor", options.Cursor)
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

	var result *model.MonitorResults
	if err = json.Unmarshal(respBodyBytes, &result); err != nil {
		return nil, fmt.Errorf("failed to JSON-decode response body: %w", err)
	}

	return result, nil
}

// GetMonitorByID retrieves a specific monitor/test by its ID
func (c *RestClient) GetMonitorByID(ctx context.Context, id string) (*model.Monitor, error) {
	if id == "" {
		return nil, fmt.Errorf("monitor ID cannot be empty")
	}

	resp, err := c.makeRequest(ctx, "GET", fmt.Sprintf("/v1/tests/%s", id), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}

	respBodyBytes, err := c.readResponseBody(resp)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var monitor *model.Monitor
	if err = json.Unmarshal(respBodyBytes, &monitor); err != nil {
		return nil, fmt.Errorf("failed to JSON-decode response body: %w", err)
	}

	return monitor, nil
}

// ListTestEntities retrieves a paginated list of entities for a specific test
func (c *RestClient) ListTestEntities(ctx context.Context, testID string, options *model.ListTestEntitiesOptions) (*model.TestEntitiesResults, error) {
	if testID == "" {
		return nil, fmt.Errorf("test ID cannot be empty")
	}

	// Build URL with query parameters
	params := url.Values{}

	if options != nil {
		if options.Limit > 0 {
			params.Set("pageSize", strconv.Itoa(options.Limit))
		}
		if options.Cursor != "" {
			params.Set("pageCursor", options.Cursor)
		}
		if options.EntityStatus != "" {
			params.Set("entityStatus", options.EntityStatus)
		}
	}

	resp, err := c.makeRequest(ctx, "GET", fmt.Sprintf("/v1/tests/%s/entities", testID), params)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}

	respBodyBytes, err := c.readResponseBody(resp)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var result *model.TestEntitiesResults
	if err = json.Unmarshal(respBodyBytes, &result); err != nil {
		return nil, fmt.Errorf("failed to JSON-decode response body: %w", err)
	}

	return result, nil
}
