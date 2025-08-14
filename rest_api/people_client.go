package rest_api

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/url"

	"github.com/turbot/steampipe-plugin-vanta/v2/rest_api/model"
)

// ListPeople retrieves a paginated list of people from Vanta
func (c *RestClient) ListPeople(ctx context.Context, options *model.ListPeopleOptions) (*model.ListPeopleOutput, error) {
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

	resp, err := c.makeRequest(ctx, "GET", "/v1/people", params)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}

	respBodyBytes, err := c.readResponseBody(resp)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var result *model.ListPeopleOutput
	if err = json.Unmarshal(respBodyBytes, &result); err != nil {
		return nil, fmt.Errorf("failed to JSON-decode response body: %w", err)
	}

	// Add debug logging for the response
	log.Printf("DEBUG: ListPeople response: people_count=%d, has_next_page=%t, end_cursor=%s",
		len(result.Results.Data),
		result.Results.PageInfo.HasNextPage,
		result.Results.PageInfo.EndCursor)

	return result, nil
}

// GetPersonByID retrieves a specific person by their ID
func (c *RestClient) GetPersonByID(ctx context.Context, id string) (*model.Person, error) {
	if id == "" {
		return nil, fmt.Errorf("person ID cannot be empty")
	}

	resp, err := c.makeRequest(ctx, "GET", fmt.Sprintf("/v1/people/%s", id), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}

	respBodyBytes, err := c.readResponseBody(resp)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var person *model.Person
	if err = json.Unmarshal(respBodyBytes, &person); err != nil {
		return nil, fmt.Errorf("failed to JSON-decode response body: %w", err)
	}

	return person, nil
}
