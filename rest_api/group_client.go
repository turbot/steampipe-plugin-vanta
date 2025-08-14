package rest_api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/turbot/steampipe-plugin-vanta/v2/rest_api/model"
)

// ListGroups retrieves a paginated list of groups from Vanta
func (c *RestClient) ListGroups(ctx context.Context, options *model.ListGroupsOptions) (*model.ListGroupsOutput, error) {
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

	resp, err := c.makeRequest(ctx, "GET", "/v1/groups", params)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}

	respBodyBytes, err := c.readResponseBody(resp)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var result *model.ListGroupsOutput
	if err = json.Unmarshal(respBodyBytes, &result); err != nil {
		return nil, fmt.Errorf("failed to JSON-decode response body: %w", err)
	}

	return result, nil
}

// GetGroupByID retrieves a specific group by its ID
func (c *RestClient) GetGroupByID(ctx context.Context, id string) (*model.GroupItem, error) {
	if id == "" {
		return nil, fmt.Errorf("group ID cannot be empty")
	}

	resp, err := c.makeRequest(ctx, "GET", fmt.Sprintf("/v1/groups/%s", id), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}

	respBodyBytes, err := c.readResponseBody(resp)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var group *model.GroupItem
	if err = json.Unmarshal(respBodyBytes, &group); err != nil {
		return nil, fmt.Errorf("failed to JSON-decode response body: %w", err)
	}

	return group, nil
}
