package rest_api

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/turbot/steampipe-plugin-vanta/v2/rest_api/model"
)

const (
	vantaAPIBaseURL = "https://api.vanta.com"
	ScopeAllRead    = "vanta-api.all:read"
)

// Vanta interface defines the methods available on the client
type Vanta interface {
	ListPeople(ctx context.Context, options *model.ListPeopleOptions) (*model.ListPeopleOutput, error)
	GetPersonByID(ctx context.Context, id string) (*model.Person, error)
	ListPolicies(ctx context.Context, options *model.ListPoliciesOptions) (*model.ListPoliciesOutput, error)
	GetPolicyByID(ctx context.Context, id string) (*model.PolicyItem, error)
	ListGroups(ctx context.Context, options *model.ListGroupsOptions) (*model.ListGroupsOutput, error)
	GetGroupByID(ctx context.Context, id string) (*model.GroupItem, error)
	ListConnectedIntegrations(ctx context.Context, options *model.ListIntegrationsOptions) (*model.ListIntegrationsOutput, error)
	GetIntegrationByID(ctx context.Context, id string) (*model.Integration, error)
	ListComputers(ctx context.Context, options *model.ListComputersOptions) (*model.ListComputersOutput, error)
	GetComputerByID(ctx context.Context, id string) (*model.Computer, error)
	ListVendors(ctx context.Context, options *model.ListVendorsOptions) (*model.ListVendorsOutput, error)
	GetVendorByID(ctx context.Context, id string) (*model.Vendor, error)
	ListMonitors(ctx context.Context, options *model.ListMonitorsOptions) (*model.MonitorResults, error)
	GetMonitorByID(ctx context.Context, id string) (*model.Monitor, error)
	ListTestEntities(ctx context.Context, testID string, options *model.ListTestEntitiesOptions) (*model.TestEntitiesResults, error)

	// Comprehensive Test API methods
	ListTests(ctx context.Context, options *model.ListTestsOptions) (*model.TestResults, error)
	GetTestByID(ctx context.Context, id string) (*model.Test, error)

	// Evidence API methods
	ListEvidence(ctx context.Context, auditID string, options *model.ListEvidenceOptions) (*model.ListEvidenceOutput, error)

	// Vulnerability API methods
	ListVulnerabilities(ctx context.Context, options *model.ListVulnerabilitiesOptions) (*model.ListVulnerabilitiesOutput, error)
	GetVulnerabilityByID(ctx context.Context, id string) (*model.Vulnerability, error)

	SetHTTPClient(client *http.Client)
}

// vanta is the internal implementation
type vanta struct {
	httpClient   *http.Client
	baseURL      string
	tokenStore   TokenStore
	clientID     string
	clientSecret string
	clientScopes []string
}

// Option represents a functional option for configuring the client
type Option func(*vanta)

// WithHTTPClient sets a custom HTTP client
func WithHTTPClient(httpClient *http.Client) Option {
	return func(v *vanta) { v.httpClient = httpClient }
}

// WithScopes sets the OAuth scopes for the client
func WithScopes(scopes ...string) Option {
	return func(v *vanta) { v.clientScopes = scopes }
}

// WithBaseURL sets the base URL for the Vanta API
func WithBaseURL(url string) Option {
	return func(v *vanta) { v.baseURL = url }
}

// WithOAuthCredentials sets OAuth client credentials for automatic token refresh
func WithOAuthCredentials(clientID, clientSecret string) Option {
	return func(v *vanta) {
		v.clientID = clientID
		v.clientSecret = clientSecret
	}
}

// WithToken sets a static Bearer token (bypasses OAuth)
func WithToken(token string) Option {
	return func(v *vanta) {
		v.tokenStore = NewBearerTokenStore(token)
	}
}

// GetOauthTokenInput represents the OAuth token request
type GetOauthTokenInput struct {
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	Scope        string `json:"scope"`
	GrantType    string `json:"grant_type"`
}

// GetOauthTokenOutput represents the OAuth token response
type GetOauthTokenOutput struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
	Scope       string `json:"scope"`
}

// New creates a new Vanta client with the given options
func New(ctx context.Context, opts ...Option) (Vanta, error) {
	v := &vanta{
		httpClient:   http.DefaultClient,
		baseURL:      vantaAPIBaseURL,
		clientScopes: []string{ScopeAllRead},
	}

	// Apply all options
	for _, opt := range opts {
		opt(v)
	}

	// If no token store is set and we have OAuth credentials, attempt to get a token
	if v.tokenStore == nil {
		if v.clientID != "" && v.clientSecret != "" {
			v.tokenStore = NewStaticTokenStore("", "")
			if err := v.refreshToken(ctx); err != nil {
				return nil, fmt.Errorf("failed to acquire auth token with oauth credentials: %v", err)
			}
		} else {
			return nil, errors.New("either provide a token with WithToken() or OAuth credentials with WithOAuthCredentials()")
		}
	}

	return v, nil
}

// refreshToken obtains a new access token using OAuth client credentials
func (v *vanta) refreshToken(ctx context.Context) error {
	if v.clientID == "" {
		return errors.New("empty oauth client id")
	}
	if v.clientSecret == "" {
		return errors.New("empty oauth client secret")
	}

	bodyBytes, err := json.Marshal(&GetOauthTokenInput{
		ClientID:     v.clientID,
		ClientSecret: v.clientSecret,
		Scope:        strings.Join(v.clientScopes, " "),
		GrantType:    "client_credentials",
	})
	if err != nil {
		return fmt.Errorf("failed to JSON-encode token request body: %v", err)
	}

	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/oauth/token", v.baseURL), bytes.NewBuffer(bodyBytes))
	if err != nil {
		return fmt.Errorf("failed to build http request: %v", err)
	}

	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/json")

	resp, err := v.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to execute http request: %v", err)
	}
	defer resp.Body.Close()

	respBodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read http response body: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("received non-200 http response status code (%d), body: %s", resp.StatusCode, string(respBodyBytes))
	}

	var oauthTokenOutput *GetOauthTokenOutput
	if err = json.Unmarshal(respBodyBytes, &oauthTokenOutput); err != nil {
		return fmt.Errorf("failed to JSON-decode token response body: %v", err)
	}

	// Update the token store with the new token
	if staticStore, ok := v.tokenStore.(*StaticTokenStore); ok {
		staticStore.SetToken(oauthTokenOutput.TokenType, oauthTokenOutput.AccessToken)
	}

	return nil
}

// Implement the Vanta interface methods by delegating to RestClient methods
func (v *vanta) ListPeople(ctx context.Context, options *model.ListPeopleOptions) (*model.ListPeopleOutput, error) {
	client := &RestClient{
		baseURL:    v.baseURL,
		httpClient: v.httpClient,
		tokenStore: v.tokenStore,
	}
	return client.ListPeople(ctx, options)
}

func (v *vanta) GetPersonByID(ctx context.Context, id string) (*model.Person, error) {
	client := &RestClient{
		baseURL:    v.baseURL,
		httpClient: v.httpClient,
		tokenStore: v.tokenStore,
	}
	return client.GetPersonByID(ctx, id)
}

func (v *vanta) ListPolicies(ctx context.Context, options *model.ListPoliciesOptions) (*model.ListPoliciesOutput, error) {
	client := &RestClient{
		baseURL:    v.baseURL,
		httpClient: v.httpClient,
		tokenStore: v.tokenStore,
	}
	return client.ListPolicies(ctx, options)
}

func (v *vanta) GetPolicyByID(ctx context.Context, id string) (*model.PolicyItem, error) {
	client := &RestClient{
		baseURL:    v.baseURL,
		httpClient: v.httpClient,
		tokenStore: v.tokenStore,
	}
	return client.GetPolicyByID(ctx, id)
}

func (v *vanta) ListGroups(ctx context.Context, options *model.ListGroupsOptions) (*model.ListGroupsOutput, error) {
	client := &RestClient{
		baseURL:    v.baseURL,
		httpClient: v.httpClient,
		tokenStore: v.tokenStore,
	}
	return client.ListGroups(ctx, options)
}

func (v *vanta) GetGroupByID(ctx context.Context, id string) (*model.GroupItem, error) {
	client := &RestClient{
		baseURL:    v.baseURL,
		httpClient: v.httpClient,
		tokenStore: v.tokenStore,
	}
	return client.GetGroupByID(ctx, id)
}

func (v *vanta) ListConnectedIntegrations(ctx context.Context, options *model.ListIntegrationsOptions) (*model.ListIntegrationsOutput, error) {
	client := &RestClient{
		baseURL:    v.baseURL,
		httpClient: v.httpClient,
		tokenStore: v.tokenStore,
	}
	return client.ListConnectedIntegrations(ctx, options)
}

func (v *vanta) GetIntegrationByID(ctx context.Context, id string) (*model.Integration, error) {
	client := &RestClient{
		baseURL:    v.baseURL,
		httpClient: v.httpClient,
		tokenStore: v.tokenStore,
	}
	return client.GetIntegrationByID(ctx, id)
}

func (v *vanta) ListComputers(ctx context.Context, options *model.ListComputersOptions) (*model.ListComputersOutput, error) {
	client := &RestClient{
		baseURL:    v.baseURL,
		httpClient: v.httpClient,
		tokenStore: v.tokenStore,
	}
	return client.ListComputers(ctx, options)
}

func (v *vanta) GetComputerByID(ctx context.Context, id string) (*model.Computer, error) {
	client := &RestClient{
		baseURL:    v.baseURL,
		httpClient: v.httpClient,
		tokenStore: v.tokenStore,
	}
	return client.GetComputerByID(ctx, id)
}

func (v *vanta) ListVendors(ctx context.Context, options *model.ListVendorsOptions) (*model.ListVendorsOutput, error) {
	client := &RestClient{
		baseURL:    v.baseURL,
		httpClient: v.httpClient,
		tokenStore: v.tokenStore,
	}
	return client.ListVendors(ctx, options)
}

func (v *vanta) GetVendorByID(ctx context.Context, id string) (*model.Vendor, error) {
	client := &RestClient{
		baseURL:    v.baseURL,
		httpClient: v.httpClient,
		tokenStore: v.tokenStore,
	}
	return client.GetVendorByID(ctx, id)
}

func (v *vanta) ListMonitors(ctx context.Context, options *model.ListMonitorsOptions) (*model.MonitorResults, error) {
	client := &RestClient{
		baseURL:    v.baseURL,
		httpClient: v.httpClient,
		tokenStore: v.tokenStore,
	}
	return client.ListMonitors(ctx, options)
}

func (v *vanta) GetMonitorByID(ctx context.Context, id string) (*model.Monitor, error) {
	client := &RestClient{
		baseURL:    v.baseURL,
		httpClient: v.httpClient,
		tokenStore: v.tokenStore,
	}
	return client.GetMonitorByID(ctx, id)
}

func (v *vanta) ListTestEntities(ctx context.Context, testID string, options *model.ListTestEntitiesOptions) (*model.TestEntitiesResults, error) {
	client := &RestClient{
		baseURL:    v.baseURL,
		httpClient: v.httpClient,
		tokenStore: v.tokenStore,
	}
	return client.ListTestEntities(ctx, testID, options)
}

// Comprehensive Test API method implementations
func (v *vanta) ListTests(ctx context.Context, options *model.ListTestsOptions) (*model.TestResults, error) {
	client := &RestClient{
		baseURL:    v.baseURL,
		httpClient: v.httpClient,
		tokenStore: v.tokenStore,
	}
	return client.ListTests(ctx, options)
}

func (v *vanta) GetTestByID(ctx context.Context, id string) (*model.Test, error) {
	client := &RestClient{
		baseURL:    v.baseURL,
		httpClient: v.httpClient,
		tokenStore: v.tokenStore,
	}
	return client.GetTestByID(ctx, id)
}

func (v *vanta) ListEvidence(ctx context.Context, auditID string, options *model.ListEvidenceOptions) (*model.ListEvidenceOutput, error) {
	client := &RestClient{
		baseURL:    v.baseURL,
		httpClient: v.httpClient,
		tokenStore: v.tokenStore,
	}
	return client.ListEvidence(ctx, auditID, options)
}

func (v *vanta) SetHTTPClient(client *http.Client) {
	v.httpClient = client
}

func (v *vanta) ListVulnerabilities(ctx context.Context, options *model.ListVulnerabilitiesOptions) (*model.ListVulnerabilitiesOutput, error) {
	client := &RestClient{
		baseURL:    v.baseURL,
		httpClient: v.httpClient,
		tokenStore: v.tokenStore,
	}
	return client.ListVulnerabilities(ctx, options)
}

func (v *vanta) GetVulnerabilityByID(ctx context.Context, id string) (*model.Vulnerability, error) {
	client := &RestClient{
		baseURL:    v.baseURL,
		httpClient: v.httpClient,
		tokenStore: v.tokenStore,
	}
	return client.GetVulnerabilityByID(ctx, id)
}
