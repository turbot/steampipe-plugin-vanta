package model

// ListIntegrationsOptions represents options for listing integrations
type ListIntegrationsOptions struct {
	Limit  int    `json:"limit,omitempty"`
	Cursor string `json:"cursor,omitempty"`
}

// ListIntegrationsOutput represents the response from the list connected integrations API
type ListIntegrationsOutput struct {
	Results IntegrationResults `json:"results"`
}

// IntegrationResults contains the actual integration data and pagination info
type IntegrationResults struct {
	PageInfo PageInfo       `json:"pageInfo"`
	Data     []*Integration `json:"data"`
}

// Integration represents a connected integration in the Vanta system
type Integration struct {
	IntegrationID string                   `json:"integrationId"`
	DisplayName   string                   `json:"displayName"`
	ResourceKinds []string                 `json:"resourceKinds"`
	Connections   []*IntegrationConnection `json:"connections"`

	// Deprecated fields - not available in REST API but kept for backward compatibility
	Description           string        `json:"description,omitempty"`
	ApplicationURL        string        `json:"applicationUrl,omitempty"`
	InstallationURL       string        `json:"installationUrl,omitempty"`
	LogoSlugID            string        `json:"logoSlugId,omitempty"`
	Credentials           interface{}   `json:"credentials,omitempty"`
	IntegrationCategories []string      `json:"integrationCategories,omitempty"`
	ServiceCategories     []string      `json:"serviceCategories,omitempty"`
	Tests                 []interface{} `json:"tests,omitempty"`
	OrganizationName      string        `json:"organizationName,omitempty"`
}

// IntegrationConnection represents a connection within an integration
type IntegrationConnection struct {
	ConnectionID           string  `json:"connectionId"`
	IsDisabled             bool    `json:"isDisabled"`
	ConnectionErrorMessage *string `json:"connectionErrorMessage"`
}
