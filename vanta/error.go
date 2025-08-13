package vanta

import (
	"slices"

	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

// function which returns an ErrorPredicate for OCI API calls
func isNotFoundError(notFoundErrors []string) plugin.ErrorPredicate {
	return func(err error) bool {
		return slices.Contains(notFoundErrors, err.Error())
	}
}
