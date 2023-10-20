package api

import (
	"context"
	"fmt"

	"github.com/machinebox/graphql"
	errorsHandler "github.com/turbot/steampipe-plugin-vanta/errors"
)

type WorkstationData struct {
	AgentVersion               string   `json:"agentVersion"`
	HasScreenLock              bool     `json:"hasScreenlock"`
	HostIdentifier             string   `json:"hostIdentifier"`
	Hostname                   string   `json:"hostname"`
	Id                         string   `json:"id"`
	InstalledAvPrograms        []string `json:"installedAvPrograms"`
	InstalledPasswordManagers  []string `json:"installedPasswordManagers"`
	IsEncrypted                bool     `json:"isEncrypted"`
	IsPasswordManagerInstalled bool     `json:"isPasswordManagerInstalled"`
	LastPing                   string   `json:"lastPing"`
	NumBrowserExtensions       int      `json:"numBrowserExtensions"`
	OsVersion                  string   `json:"osVersion"`
	SerialNumber               string   `json:"serialNumber"`
}

type WorkstationOwner struct {
	DisplayName string
	Id          string
}

type WorkstationUnsupportedReasons struct {
	UnsupportedOsVersion bool `json:"unsupportedOsVersion"`
	UnsupportedOsType    bool `json:"unsupportedOsType"`
}

type Workstation struct {
	Data               WorkstationData               `json:"data"`
	Id                 string                        `json:"id"`
	OrganizationName   string                        `json:"-"`
	Owner              WorkstationOwner              `json:"-"`
	UnsupportedReasons WorkstationUnsupportedReasons `json:"unsupportedReasons"`
}

type DomainEndPointQueryUser struct {
	DisplayName  string        `json:"displayName"`
	Id           string        `json:"id"`
	Workstations []Workstation `json:"workstations"`
}

type DomainEndPointQueryOrganization struct {
	Id    string                    `json:"id"`
	Name  string                    `json:"name"`
	Users []DomainEndPointQueryUser `json:"users"`
}

// ListWorkstationsResponse is returned by ListWorkstations on success
type ListWorkstationsResponse struct {
	Organization DomainEndPointQueryOrganization `json:"organization"`
}

// Define the query
const (
	queryWorkstationList = `
query fetchDomainEndpoints {
  organization {
    id
    name
    uiComponentStates {
      agentBannerIsCollapsed
    }
    users(includeRemovedUsers: true, includeNonHumanUsers: true) {
      id
      displayName
      ...UserComputerFields
    }
  }
}

fragment UserComputerFields on User {
  id
  workstations(includeUnsupported: true) {
    id
    unsupportedReasons {
      unsupportedOsVersion
      unsupportedOsType
    }
    data {
      id
      agentVersion
      osVersion
      lastPing
      serialNumber
      hostIdentifier
      hostname
      ... on macosWorkstationData {
        installedAvPrograms
        installedPasswordManagers
        isEncrypted
        isPasswordManagerInstalled
        numBrowserExtensions
        hasScreenlock
      }
      ... on windowsWorkstationData {
        installedAvPrograms
        installedPasswordManagers
        isEncrypted
        isPasswordManagerInstalled
        numBrowserExtensions
        hasScreenlock
      }
      ... on linuxWorkstationData {
      installedAvPrograms
        isEncrypted
      }
    }
  }
  managedComputers {
    id
    uniqueId
    udid
    updatedAt
    hasScreenlock
    name
    isEncrypted
    operatingSystem {
      name
      version
    }
    hardware {
      serialNumber
    }
    passwordManagers {
      name
    }
    antivirusNames
    vantaAttributes {
      key
      value
      managedExternally
    }
    ... on SpecificMicrosoftEndpointManagerManagedComputerResource {
      avPolicies {
        id
        name
      }
    }
    ... on ReceivedMacosUserComputerResource {
      installedApp {
        id
        app {
          id
          name
        }
      }
    }
    ... on ReceivedWindowsUserComputerResource {
      installedApp {
        id
        app {
          id
          name
        }
      }
    }
  }
}	
`

	queryEndpointApplicationsList = `
query fetchEndpointApplications($endpointId: String!) {
  organization {
    id
    osqueryEndpointApplicationData(endpointId: $endpointId)
  }
}
`
)

// ListWorkstations returns the list all computers within your organization with their security-relevant settings information
//
// @param ctx context for configuration
//
// @param client the API client
func ListWorkstations(
	ctx context.Context,
	client *Client,
) (*ListWorkstationsResponse, error) {

	// Make a request
	req := graphql.NewRequest(queryWorkstationList)

	// set header fields
	req.Header.Set("user-agent", "steampipe")
	req.Header.Set("Cookie", fmt.Sprintf("connect.sid=%s", *client.SessionId))
	req.Header.Set("x-csrf-token", "this_csrf_header_is_constant")
	req.Header.Set("Cache-Control", "no-cache")

	var err error
	var data ListWorkstationsResponse

	if err := client.Graphql.Run(ctx, req, &data); err != nil {
		err = errorsHandler.BuildErrorMessage(err)
		return nil, err
	}

	return &data, err
}

type EndpointApplicationQueryOrganization struct {
	EndpointApplications []string `json:"osqueryEndpointApplicationData"`
}

// ListEndpointApplicationsResponse is returned by ListEndpointApplications on success
type ListEndpointApplicationsResponse struct {
	Organization EndpointApplicationQueryOrganization `json:"organization"`
}

// ListEndpointApplications returns the list of applications installed on a device
//
// @param ctx context for configuration
//
// @param client the API client
//
// @param id the endpoint ID
func ListEndpointApplications(
	ctx context.Context,
	client *Client,
	id string,
) (*ListEndpointApplicationsResponse, error) {

	// Make a request
	req := graphql.NewRequest(queryEndpointApplicationsList)

	req.Var("endpointId", id)

	// set header fields
	req.Header.Set("user-agent", "steampipe")
	req.Header.Set("Cookie", fmt.Sprintf("connect.sid=%s", *client.SessionId))
	req.Header.Set("x-csrf-token", "this_csrf_header_is_constant")
	req.Header.Set("Cache-Control", "no-cache")

	var err error
	var data ListEndpointApplicationsResponse

	if err := client.Graphql.Run(ctx, req, &data); err != nil {
		err = errorsHandler.BuildErrorMessage(err)
		return nil, err
	}

	return &data, err
}
