package model

import "time"

// ListEvidenceOptions represents the options for listing evidence
type ListEvidenceOptions struct {
	AuditID string `json:"audit_id"`
	Limit   int    `json:"limit,omitempty"`
	Cursor  string `json:"cursor,omitempty"`
}

// ListEvidenceOutput represents the response from the list evidence API
type ListEvidenceOutput struct {
	Results EvidenceResults `json:"results"`
}

// EvidenceResults contains the actual evidence data and pagination info
type EvidenceResults struct {
	PageInfo PageInfo    `json:"pageInfo"`
	Data     []*Evidence `json:"data"`
}

// Evidence represents audit evidence in the Vanta system
type Evidence struct {
	ID                string            `json:"id"`
	ExternalID        string            `json:"externalId"`
	Status            string            `json:"status"`
	Name              string            `json:"name"`
	DeletionDate      *time.Time        `json:"deletionDate"`
	CreationDate      time.Time         `json:"creationDate"`
	StatusUpdatedDate time.Time         `json:"statusUpdatedDate"`
	TestStatus        *string           `json:"testStatus"`
	EvidenceType      string            `json:"evidenceType"`
	EvidenceID        string            `json:"evidenceId"`
	RelatedControls   []*RelatedControl `json:"relatedControls"`
	Description       *string           `json:"description"`
}

// RelatedControl represents a control associated to evidence
type RelatedControl struct {
	Name         string   `json:"name"`
	SectionNames []string `json:"sectionNames"`
}

// EvidenceStatus represents the possible statuses for evidence
const (
	EvidenceStatusAccepted    = "Accepted"
	EvidenceStatusFlagged     = "Flagged"
	EvidenceStatusInitialized = "Initialized"
	EvidenceStatusNA          = "NA"
	EvidenceStatusNotReady    = "Not ready for audit"
	EvidenceStatusReady       = "Ready for audit"
)

// EvidenceType represents the possible types of evidence
const (
	EvidenceTypeEvidence = "Evidence"
	EvidenceTypeRequest  = "Request"
	EvidenceTypePolicy   = "Policy"
	EvidenceTypeTest     = "Test"
)
