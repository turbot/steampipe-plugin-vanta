package model

import "time"

// ListGroupsOptions represents options for listing groups
type ListGroupsOptions struct {
	Limit  int    `json:"limit,omitempty"`
	Cursor string `json:"cursor,omitempty"`
}

// ListGroupsOutput represents the response from the list groups API
type ListGroupsOutput struct {
	Results GroupResults `json:"results"`
}

// GroupResults contains the actual group data and pagination info
type GroupResults struct {
	PageInfo PageInfo     `json:"pageInfo"`
	Data     []*GroupItem `json:"data"`
}

// GroupItem represents a group in the Vanta system
type GroupItem struct {
	ID           string     `json:"id"`
	Name         string     `json:"name"`
	CreationDate *time.Time `json:"creationDate,omitempty"`
}
