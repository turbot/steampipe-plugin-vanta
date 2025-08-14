package rest_api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/turbot/steampipe-plugin-vanta/v2/rest_api/model"
)

// ListPolicies retrieves a paginated list of policies from Vanta
func (c *RestClient) ListPolicies(ctx context.Context, options *model.ListPoliciesOptions) (*model.ListPoliciesOutput, error) {
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

	resp, err := c.makeRequest(ctx, "GET", "/v1/policies", params)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}

	respBodyBytes, err := c.readResponseBody(resp)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var result *model.ListPoliciesOutput
	if err = json.Unmarshal(respBodyBytes, &result); err != nil {
		return nil, fmt.Errorf("failed to JSON-decode response body: %w", err)
	}

	return result, nil
}

// GetPolicyByID retrieves a specific policy by its ID
func (c *RestClient) GetPolicyByID(ctx context.Context, id string) (*model.PolicyItem, error) {
	if id == "" {
		return nil, fmt.Errorf("policy ID cannot be empty")
	}

	resp, err := c.makeRequest(ctx, "GET", fmt.Sprintf("/v1/policies/%s", id), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}

	respBodyBytes, err := c.readResponseBody(resp)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var policy *model.PolicyItem
	if err = json.Unmarshal(respBodyBytes, &policy); err != nil {
		return nil, fmt.Errorf("failed to JSON-decode response body: %w", err)
	}

	return policy, nil
}
