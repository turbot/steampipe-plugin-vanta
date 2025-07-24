package model

import (
	"time"
)

// ListPeopleOptions represents options for listing people
type ListPeopleOptions struct {
	Limit  int    `json:"limit,omitempty"`
	Cursor string `json:"cursor,omitempty"`
}

// ListPeopleOutput represents the response from the list people API
type ListPeopleOutput struct {
	Results PeopleResults `json:"results"`
}

// PeopleResults contains the actual people data and pagination info
type PeopleResults struct {
	PageInfo PageInfo  `json:"pageInfo"`
	Data     []*Person `json:"data"`
}

// Person represents a person in the Vanta system
type Person struct {
	ID           string        `json:"id"`
	EmailAddress string        `json:"emailAddress"`
	Employment   *Employment   `json:"employment,omitempty"`
	Name         *Name         `json:"name,omitempty"`
	GroupIDs     []string      `json:"groupIds,omitempty"`
	Sources      *Sources      `json:"sources,omitempty"`
	TasksSummary *TasksSummary `json:"tasksSummary,omitempty"`
}

// Employment contains employment information
type Employment struct {
	EndDate   *time.Time        `json:"endDate,omitempty"`
	JobTitle  *string           `json:"jobTitle,omitempty"`
	StartDate *time.Time        `json:"startDate,omitempty"`
	Status    *EmploymentStatus `json:"status,omitempty"`
}

// LeaveInfo contains leave information
type LeaveInfo struct {
	// Add fields as needed based on actual API response
}

// Name contains name information
type Name struct {
	Display string `json:"display"`
	Last    string `json:"last"`
	First   string `json:"first"`
}

// Sources contains data source information
type Sources struct {
	EmailAddress *SourceInfo `json:"emailAddress,omitempty"`
	Employment   *SourceInfo `json:"employment,omitempty"`
}

// SourceInfo contains information about a data source
type SourceInfo struct {
	IntegrationID string `json:"integrationId,omitempty"`
	ResourceID    string `json:"resourceId,omitempty"`
	Type          string `json:"type,omitempty"`
}

// TasksSummary contains task completion summary
type TasksSummary struct {
	CompletionDate *time.Time   `json:"completionDate,omitempty"`
	DueDate        *time.Time   `json:"dueDate,omitempty"`
	Status         string       `json:"status"`
	Details        TasksDetails `json:"details"`
}

// TasksDetails contains detailed task information
type TasksDetails struct {
	CompleteTrainings              *TaskDetail `json:"completeTrainings,omitempty"`
	CompleteCustomTasks            *TaskDetail `json:"completeCustomTasks,omitempty"`
	CompleteOffboardingCustomTasks *TaskDetail `json:"completeOffboardingCustomTasks,omitempty"`
	CompleteBackgroundChecks       *TaskDetail `json:"completeBackgroundChecks,omitempty"`
	AcceptPolicies                 *PolicyTask `json:"acceptPolicies,omitempty"`
	InstallDeviceMonitoring        *TaskDetail `json:"installDeviceMonitoring,omitempty"`
}

// TaskDetail represents a generic task detail
type TaskDetail struct {
	TaskType       string      `json:"taskType,omitempty"`
	Status         string      `json:"status,omitempty"`
	DueDate        *time.Time  `json:"dueDate,omitempty"`
	CompletionDate *time.Time  `json:"completionDate,omitempty"`
	Disabled       interface{} `json:"disabled,omitempty"`
}

// PolicyTask represents policy acceptance task details
type PolicyTask struct {
	TaskType           string      `json:"taskType,omitempty"`
	Status             string      `json:"status,omitempty"`
	DueDate            *time.Time  `json:"dueDate,omitempty"`
	CompletionDate     *time.Time  `json:"completionDate,omitempty"`
	Disabled           interface{} `json:"disabled,omitempty"`
	UnacceptedPolicies []Policy    `json:"unacceptedPolicies,omitempty"`
	AcceptedPolicies   []Policy    `json:"acceptedPolicies,omitempty"`
}

// Policy represents a policy item
type Policy struct {
	Name string `json:"name"`
}

type TaskStatus string

type EmploymentStatus string

const (
	TaskStatusComplete                      TaskStatus = "COMPLETE"
	TaskStatusDueSoon                       TaskStatus = "DUE_SOON"
	TaskStatusNone                          TaskStatus = "NONE"
	TaskStatusOverdue                       TaskStatus = "OVERDUE"
	TaskStatusPaused                        TaskStatus = "PAUSED"
	TaskStatusOffboardingComplete           TaskStatus = "OFFBOARDING_COMPLETE"
	TaskStatusOffboardingDueSoon            TaskStatus = "OFFBOARDING_DUE_SOON"
	TaskStatusOffboardingOffboardingOverdue TaskStatus = "OFFBOARDING_OVERDUE"

	EmploymentStatusUpcoming EmploymentStatus = "UPCOMING"
	EmploymentStatusCurrent  EmploymentStatus = "CURRENT"
	EmploymentStatusOnLeave  EmploymentStatus = "ON_LEAVE"
	EmploymentStatusInactive EmploymentStatus = "INACTIVE"
)

// Policy related types

// ListPoliciesOptions represents options for listing policies
type ListPoliciesOptions struct {
	Limit  int    `json:"limit,omitempty"`
	Cursor string `json:"cursor,omitempty"`
}

// ListPoliciesOutput represents the response from the list policies API
type ListPoliciesOutput struct {
	Results PolicyResults `json:"results"`
}

// PolicyResults contains the actual policy data and pagination info
type PolicyResults struct {
	PageInfo PageInfo      `json:"pageInfo"`
	Data     []*PolicyItem `json:"data"`
}

// PolicyItem represents a policy in the Vanta system
type PolicyItem struct {
	ID             string               `json:"id"`
	Name           string               `json:"name"`
	Description    string               `json:"description"`
	Status         string               `json:"status"`
	ApprovedAtDate *time.Time           `json:"approvedAtDate,omitempty"`
	LatestVersion  *PolicyLatestVersion `json:"latestVersion,omitempty"`
}

// PolicyLatestVersion represents the latest version information of a policy
type PolicyLatestVersion struct {
	Status string `json:"status"`
}

type PolicyStatus string

const (
	PolicyStatusNeedsRemediation PolicyStatus = "NEEDS_REMEDIATION"
	PolicyStatusCompliant        PolicyStatus = "COMPLIANT"
	PolicyStatusNotStarted       PolicyStatus = "NOT_STARTED"
)

