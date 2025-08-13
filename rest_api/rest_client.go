package rest_api

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

// TokenStore interface for managing authentication tokens
type TokenStore interface {
	GetToken() (tokenType, token string)
}

// RestClient wraps the Vanta API client with additional functionality
type RestClient struct {
	baseURL    string
	httpClient *http.Client
	tokenStore TokenStore
}

// NewRestClient creates a new REST client instance
func NewRestClient(baseURL string, tokenStore TokenStore) *RestClient {
	return &RestClient{
		baseURL:    baseURL,
		httpClient: http.DefaultClient,
		tokenStore: tokenStore,
	}
}

// SetHTTPClient allows setting a custom HTTP client
func (c *RestClient) SetHTTPClient(client *http.Client) {
	c.httpClient = client
}

// makeRequest performs HTTP requests with proper authentication
func (c *RestClient) makeRequest(ctx context.Context, method, path string, queryParams url.Values) (*http.Response, error) {
	tokenType, token := c.tokenStore.GetToken()
	if token == "" {
		return nil, errors.New("no auth token present")
	}

	u, err := url.Parse(fmt.Sprintf("%s%s", c.baseURL, path))
	if err != nil {
		return nil, fmt.Errorf("failed to parse URL: %v", err)
	}

	if queryParams != nil {
		u.RawQuery = queryParams.Encode()
	}

	req, err := http.NewRequest(method, u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to build http request: %v", err)
	}

	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("%s %s", tokenType, token))

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute http request: %v", err)
	}

	return resp, nil
}

// readResponseBody reads and validates HTTP response
func (c *RestClient) readResponseBody(resp *http.Response) ([]byte, error) {
	defer resp.Body.Close()

	respBodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read http response body: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("received non-200 http response status code (%d), body: %s", resp.StatusCode, string(respBodyBytes))
	}

	return respBodyBytes, nil
}
