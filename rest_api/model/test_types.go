package model

import "time"

// Test represents a comprehensive test/monitor in Vanta with full API details
type Test struct {
	ID                     string                 `json:"id"`
	Name                   string                 `json:"name"`
	LastTestRunDate        *time.Time             `json:"lastTestRunDate"`
	LatestFlipDate         *time.Time             `json:"latestFlipDate"`
	Description            string                 `json:"description"`
	FailureDescription     string                 `json:"failureDescription"`
	RemediationDescription string                 `json:"remediationDescription"`
	Version                *TestVersion           `json:"version"`
	Category               string                 `json:"category"`
	Integrations           []string               `json:"integrations"`
	Status                 string                 `json:"status"`
	DeactivatedStatusInfo  *TestDeactivatedStatus `json:"deactivatedStatusInfo"`
	RemediationStatusInfo  *TestRemediationStatus `json:"remediationStatusInfo"`
	Owner                  *TestOwner             `json:"owner"`
}

// TestVersion represents version information for a test
type TestVersion struct {
	Major int    `json:"major"`
	Minor int    `json:"minor"`
	ID    string `json:"_id"`
}

// TestDeactivatedStatus represents deactivation status information
type TestDeactivatedStatus struct {
	IsDeactivated     bool       `json:"isDeactivated"`
	DeactivatedReason *string    `json:"deactivatedReason"`
	LastUpdatedDate   *time.Time `json:"lastUpdatedDate"`
}

// TestRemediationStatus represents remediation status information
type TestRemediationStatus struct {
	Status                 string     `json:"status"`
	SoonestRemediateByDate *time.Time `json:"soonestRemediateByDate"`
	ItemCount              int        `json:"itemCount"`
}

// TestOwner represents the owner of a test
type TestOwner struct {
	ID           string `json:"id"`
	EmailAddress string `json:"emailAddress"`
	DisplayName  string `json:"displayName"`
}

// TestResults represents the paginated response for comprehensive test list
type TestResults struct {
	Results TestResultsData `json:"results"`
}

// TestResultsData represents the data portion of test results
type TestResultsData struct {
	PageInfo PageInfo `json:"pageInfo"`
	Data     []*Test  `json:"data"`
}

// ListTestsOptions represents options for listing comprehensive tests
type ListTestsOptions struct {
	PageSize          int    `json:"pageSize,omitempty"`          // 1 to 100, defaults to 10
	PageCursor        string `json:"pageCursor,omitempty"`        // Pagination cursor
	StatusFilter      string `json:"statusFilter,omitempty"`      // OK, DEACTIVATED, NEEDS_ATTENTION, IN_PROGRESS, INVALID, NOT_APPLICABLE
	FrameworkFilter   string `json:"frameworkFilter,omitempty"`   // Filter by framework
	IntegrationFilter string `json:"integrationFilter,omitempty"` // Filter by integration (e.g., "aws")
	ControlFilter     string `json:"controlFilter,omitempty"`     // Filter by control ID
	OwnerFilter       string `json:"ownerFilter,omitempty"`       // Filter by owner ID
	CategoryFilter    string `json:"categoryFilter,omitempty"`    // Filter by category
	IsInRollout       *bool  `json:"isInRollout,omitempty"`       // Filter by rollout status
}
