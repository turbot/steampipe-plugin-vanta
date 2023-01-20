package api

import (
	"context"
	"fmt"

	"github.com/machinebox/graphql"
	errorsHandler "github.com/turbot/steampipe-plugin-vanta/errors"
)

type WorkstationOwner struct {
	DisplayName string
	Email       string
}

type WorkstationData struct {
	AgentVersion               string   `json:"agentVersion"`
	OsVersion                  string   `json:"osVersion"`
	LastPing                   string   `json:"lastPing"`
	SerialNumber               string   `json:"serialNumber"`
	HostIdentifier             string   `json:"hostIdentifier"`
	Hostname                   string   `json:"hostname"`
	IsEncrypted                bool     `json:"isEncrypted"`
	IsPasswordManagerInstalled bool     `json:"isPasswordManagerInstalled"`
	NumBrowserExtensions       int      `json:"numBrowserExtensions"`
	HasScreenLock              bool     `json:"hasScreenlock"`
	InstalledAvPrograms        []string `json:"installedAvPrograms"`
	InstalledPasswordManagers  []string `json:"installedPasswordManagers"`
}

type Workstation struct {
	Id    string          `json:"id"`
	Data  WorkstationData `json:"data"`
	Owner WorkstationOwner
}

type UserQueryWorkstation struct {
	DisplayName  string        `json:"displayName"`
	Email        string        `json:"email"`
	Workstations []Workstation `json:"workstations"`
}

type DomainEndPointQueryUser struct {
	DisplayName  string        `json:"displayName"`
	Email        string        `json:"email"`
	Workstations []Workstation `json:"workstations"`
}

type DomainEndPointQueryOrganization struct {
	Users []DomainEndPointQueryUser `json:"users"`
}

type ListWorkstationsResponse struct {
	Organization DomainEndPointQueryOrganization `json:"organization"`
}

// Define the query
const (
	queryTest = `
query fetchDomainEndpoints {
	organization {
		id
		requiresLocationServices
		uiComponentStates {
			agentBannerIsCollapsed
		}
		users(includeRemovedUsers: true, includeNonHumanUsers: true) {
			id
			displayName
			email
			familyName
			givenName
			imageUrl
			isActive
			isNotHuman
			role {
				id
				name
			}
			agentTask {
				taskSatisfied
				completionDate
			}
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
)

// ListWorkstations returns the list all computers within your organization with their security-relevant settings information
func ListWorkstations(
	ctx context.Context,
	client *Client,
) (*ListWorkstationsResponse, error) {

	// Make a request
	req := graphql.NewRequest(queryTest)

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
