package rest_api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/turbot/steampipe-plugin-vanta/rest_api/model"
)

// ListEvidence retrieves a list of evidence for a specific audit
func (c *RestClient) ListEvidence(ctx context.Context, auditID string, options *model.ListEvidenceOptions) (*model.ListEvidenceOutput, error) {
	if auditID == "" {
		return nil, fmt.Errorf("audit ID is required")
	}

	path := fmt.Sprintf("/v1/audits/%s/evidence", auditID)

	// Build query parameters
	queryParams := url.Values{}
	if options != nil {
		if options.Limit > 0 {
			queryParams.Set("limit", fmt.Sprintf("%d", options.Limit))
		}
		if options.Cursor != "" {
			queryParams.Set("cursor", options.Cursor)
		}
	}

	resp, err := c.makeRequest(ctx, "GET", path, queryParams)
	if err != nil {
		return nil, err
	}

	body, err := c.readResponseBody(resp)
	if err != nil {
		return nil, err
	}

	var result model.ListEvidenceOutput
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal evidence list response: %v", err)
	}

	return &result, nil
}
