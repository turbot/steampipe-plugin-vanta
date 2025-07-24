package model

import (
	"time"
)

// ListComputersOptions represents options for listing computers
type ListComputersOptions struct {
	Limit  int    `json:"limit,omitempty"`
	Cursor string `json:"cursor,omitempty"`
}

// ListComputersOutput represents the response from the list computers API
type ListComputersOutput struct {
	Results ComputerResults `json:"results"`
}

// ComputerResults contains the actual computer data and pagination info
type ComputerResults struct {
	PageInfo PageInfo    `json:"pageInfo"`
	Data     []*Computer `json:"data"`
}

// Computer represents a computer/device in the Vanta system
type Computer struct {
	ID                    string           `json:"id"`
	IntegrationID         string           `json:"integrationId"`
	LastCheckDate         *time.Time       `json:"lastCheckDate,omitempty"`
	OperatingSystem       *OperatingSystem `json:"operatingSystem,omitempty"`
	Owner                 *ComputerOwner   `json:"owner,omitempty"`
	SerialNumber          string           `json:"serialNumber"`
	UDID                  string           `json:"udid"`
	Screenlock            *SecurityCheck   `json:"screenlock,omitempty"`
	DiskEncryption        *SecurityCheck   `json:"diskEncryption,omitempty"`
	PasswordManager       *SecurityCheck   `json:"passwordManager,omitempty"`
	AntivirusInstallation *SecurityCheck   `json:"antivirusInstallation,omitempty"`

	// Deprecated fields - not available in REST API but kept for backward compatibility
	OwnerName                  string        `json:"ownerName,omitempty"`
	AgentVersion               string        `json:"agentVersion,omitempty"`
	OsVersion                  string        `json:"osVersion,omitempty"`
	IsPasswordManagerInstalled *bool         `json:"isPasswordManagerInstalled,omitempty"`
	IsEncrypted                *bool         `json:"isEncrypted,omitempty"`
	HasScreenLock              *bool         `json:"hasScreenLock,omitempty"`
	Hostname                   string        `json:"hostname,omitempty"`
	HostIdentifier             string        `json:"hostIdentifier,omitempty"`
	LastPing                   *time.Time    `json:"lastPing,omitempty"`
	NumBrowserExtensions       *int          `json:"numBrowserExtensions,omitempty"`
	EndpointApplications       []interface{} `json:"endpointApplications,omitempty"`
	InstalledAvPrograms        []interface{} `json:"installedAvPrograms,omitempty"`
	InstalledPasswordManagers  []interface{} `json:"installedPasswordManagers,omitempty"`
	UnsupportedReasons         interface{}   `json:"unsupportedReasons,omitempty"`
	OwnerID                    string        `json:"ownerId,omitempty"`
	OrganizationName           string        `json:"organizationName,omitempty"`
}

// ComputerOwner represents the owner of a computer
type ComputerOwner struct {
	ID           string `json:"id"`
	EmailAddress string `json:"emailAddress"`
	DisplayName  string `json:"displayName"`
}

// OperatingSystem represents operating system information
type OperatingSystem struct {
	Type    string `json:"type"`
	Version string `json:"version"`
}

// SecurityCheck represents the outcome of a security check
type SecurityCheck struct {
	Outcome string `json:"outcome"`
}
