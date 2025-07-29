package model

import "time"

// Monitor represents a monitor/test in Vanta
type Monitor struct {
	ID                     string                 `json:"id"`
	Name                   string                 `json:"name"`
	LastTestRunDate        *time.Time             `json:"lastTestRunDate"`
	LatestFlipDate         *time.Time             `json:"latestFlipDate"`
	Description            string                 `json:"description"`
	FailureDescription     string                 `json:"failureDescription"`
	RemediationDescription string                 `json:"remediationDescription"`
	Version                *MonitorVersion        `json:"version"`
	Category               string                 `json:"category"`
	Integrations           []string               `json:"integrations"`
	Status                 string                 `json:"status"`
	DeactivatedStatusInfo  *DeactivatedStatusInfo `json:"deactivatedStatusInfo"`
	RemediationStatusInfo  *RemediationStatusInfo `json:"remediationStatusInfo"`
	Owner                  *MonitorOwner          `json:"owner"`
}

// MonitorVersion represents version information for a monitor
type MonitorVersion struct {
	Major int    `json:"major"`
	Minor int    `json:"minor"`
	ID    string `json:"_id"`
}

// DeactivatedStatusInfo represents deactivation status information
type DeactivatedStatusInfo struct {
	IsDeactivated     bool       `json:"isDeactivated"`
	DeactivatedReason *string    `json:"deactivatedReason"`
	LastUpdatedDate   *time.Time `json:"lastUpdatedDate"`
}

// RemediationStatusInfo represents remediation status information
type RemediationStatusInfo struct {
	Status                 string     `json:"status"`
	SoonestRemediateByDate *time.Time `json:"soonestRemediateByDate"`
	ItemCount              int        `json:"itemCount"`
}

// MonitorOwner represents the owner of a monitor
type MonitorOwner struct {
	ID           string `json:"id"`
	EmailAddress string `json:"emailAddress"`
	DisplayName  string `json:"displayName"`
}

// MonitorResults represents the paginated response for monitor list
type MonitorResults struct {
	Results MonitorResultsData `json:"results"`
}

// MonitorResultsData represents the data portion of monitor results
type MonitorResultsData struct {
	PageInfo PageInfo   `json:"pageInfo"`
	Data     []*Monitor `json:"data"`
}

// ListMonitorsOptions represents options for listing monitors
type ListMonitorsOptions struct {
	Limit  int    `json:"limit,omitempty"`
	Cursor string `json:"cursor,omitempty"`
}

// TestEntity represents an entity associated with a test
type TestEntity struct {
	ID                string     `json:"id"`
	EntityStatus      string     `json:"entityStatus"`
	DisplayName       string     `json:"displayName"`
	ResponseType      string     `json:"responseType"`
	DeactivatedReason *string    `json:"deactivatedReason"`
	LastUpdatedDate   *time.Time `json:"lastUpdatedDate"`
	CreatedDate       *time.Time `json:"createdDate"`
}

// TestEntitiesResults represents the paginated response for test entities
type TestEntitiesResults struct {
	Results TestEntitiesResultsData `json:"results"`
}

// TestEntitiesResultsData represents the data portion of test entities results
type TestEntitiesResultsData struct {
	PageInfo PageInfo      `json:"pageInfo"`
	Data     []*TestEntity `json:"data"`
}

// ListTestEntitiesOptions represents options for listing test entities
type ListTestEntitiesOptions struct {
	Limit        int    `json:"limit,omitempty"`
	Cursor       string `json:"cursor,omitempty"`
	EntityStatus string `json:"entityStatus,omitempty"`
}
