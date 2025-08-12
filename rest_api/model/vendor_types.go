package model

import (
	"time"
)

// ListVendorsOptions represents options for listing vendors
type ListVendorsOptions struct {
	Limit  int    `json:"limit,omitempty"`
	Cursor string `json:"cursor,omitempty"`
}

// ListVendorsOutput represents the response from the list vendors API
type ListVendorsOutput struct {
	Results VendorResults `json:"results"`
}

// VendorResults contains the actual vendor data and pagination info
type VendorResults struct {
	PageInfo PageInfo  `json:"pageInfo"`
	Data     []*Vendor `json:"data"`
}

// Vendor represents a vendor/company in the Vanta system
type Vendor struct {
	ID                               string             `json:"id"`
	Name                             string             `json:"name"`
	WebsiteURL                       string             `json:"websiteUrl"`
	AccountManagerName               *string            `json:"accountManagerName"`
	AccountManagerEmail              *string            `json:"accountManagerEmail"`
	ServicesProvided                 *string            `json:"servicesProvided"`
	AdditionalNotes                  *string            `json:"additionalNotes"`
	SecurityOwnerUserID              string             `json:"securityOwnerUserId"`
	BusinessOwnerUserID              string             `json:"businessOwnerUserId"`
	ContractStartDate                *time.Time         `json:"contractStartDate"`
	ContractRenewalDate              *time.Time         `json:"contractRenewalDate"`
	ContractTerminationDate          *time.Time         `json:"contractTerminationDate"`
	NextSecurityReviewDueDate        *time.Time         `json:"nextSecurityReviewDueDate"`
	LastSecurityReviewCompletionDate *time.Time         `json:"lastSecurityReviewCompletionDate"`
	IsVisibleToAuditors              bool               `json:"isVisibleToAuditors"`
	IsRiskAutoScored                 bool               `json:"isRiskAutoScored"`
	Category                         *VendorCategory    `json:"category"`
	AuthDetails                      *VendorAuthDetails `json:"authDetails"`
	RiskAttributeIDs                 []string           `json:"riskAttributeIds"`
	Status                           string             `json:"status"`
	InherentRiskLevel                string             `json:"inherentRiskLevel"`
	ResidualRiskLevel                string             `json:"residualRiskLevel"`
	VendorHeadquarters               *string            `json:"vendorHeadquarters"`
	ContractAmount                   *float64           `json:"contractAmount"`
	CustomFields                     interface{}        `json:"customFields"`

	// Deprecated fields - not available in REST API but kept for backward compatibility
	Title             string      `json:"title,omitempty"`
	EvidenceRequestID string      `json:"evidenceRequestId,omitempty"`
	Description       string      `json:"description,omitempty"`
	UID               string      `json:"uid,omitempty"`
	AppUploadEnabled  *bool       `json:"appUploadEnabled,omitempty"`
	Restricted        *bool       `json:"restricted,omitempty"`
	DismissedStatus   interface{} `json:"dismissedStatus,omitempty"`
	RenewalMetadata   interface{} `json:"renewalMetadata,omitempty"`
	OrganizationName  string      `json:"organizationName,omitempty"`
}

// VendorCategory represents the category of a vendor
type VendorCategory struct {
	DisplayName string `json:"displayName"`
}

// VendorAuthDetails represents authentication details for a vendor
type VendorAuthDetails struct {
	Method                 string `json:"method"`
	PasswordMFA            *bool  `json:"passwordMFA"`
	PasswordRequiresNumber *bool  `json:"passwordRequiresNumber"`
	PasswordRequiresSymbol *bool  `json:"passwordRequiresSymbol"`
	PasswordMinimumLength  *int   `json:"passwordMinimumLength"`
}
