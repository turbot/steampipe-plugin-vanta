---
organization: Turbot
category: ["saas"]
icon_url: "/images/plugins/turbot/vanta.svg"
brand_color: "#5230d7"
display_name: "Vanta"
short_name: "vanta"
description: "Steampipe plugin to query users, policies, compliances, and more from your Vanta organization."
og_description: "Query Vanta with SQL! Open source CLI. No DB required."
og_image: "/images/plugins/turbot/vanta-social-graphic.png"
engines: ["steampipe", "sqlite", "postgres", "export"]
---

# Vanta + Steampipe

[Vanta](https://www.vanta.com) helps businesses get and stay compliant by continuously monitoring your people, systems and tools to improve the security posture.

[Steampipe](https://steampipe.io) is an open-source zero-ETL engine to instantly query cloud APIs using SQL.

List all active users in your organization:

```sql
select
  display_name,
  id,
  email,
  is_active
from
  vanta_user
where
  is_active;
```

```
+--------------+--------------------------+----------------+-----------+
| display_name | id                       | email          | is_active |
+--------------+--------------------------+----------------+-----------+
| Simba        | 5fb30b86a228f6b6f7024535 | simba@test.com | true      |
| Timon        | 5fb30b86a228f6b6f70245e7 | timon@test.com | true      |
+--------------+--------------------------+----------------+-----------+
```

## Documentation

- **[Table definitions & examples →](/plugins/turbot/vanta/tables)**

## Get started

### Install

Download and install the latest Vanta plugin:

```bash
steampipe plugin install vanta
```

### Credentials

| Item        | Description                                                                                                                                                                                                                      |
| ----------- | -------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| Credentials | The plugin supports OAuth-based authentication using either:<br/>1. OAuth client credentials (`client_id` and `client_secret`)<br/>2. Personal [access token]https://developer.vanta.com/docs/api-access-setup) (`access_token`) |
| Permissions | User requires appropriate access permissions to query Vanta resources. OAuth applications need to be configured with the scopes `auditor-api.audit:read` and `vanta-api.all:read` for the resources you want to access.          |
| Radius      | Each connection represents a single Vanta organization.                                                                                                                                                                          |
| Resolution  | Credentials explicitly set in a steampipe config file (`~/.steampipe/config/vanta.spc`).                                                                                                                                         |

The `vanta_evidence` table requires auditor-level access with specific scopes (`auditor-api.audit:read`, `auditor-api.audit:write`, `auditor-api.auditor:read`, `auditor-api.auditor:write`). You must be a registered Vanta Audit Partner to access evidence data.

For auditor setup instructions, visit the [Auditor Application Setup Guide](https://developer.vanta.com/docs/auditor-application-setup).

To generate standard credentials, visit the [Vanta Developer Documentation](https://developer.vanta.com/docs/api-access-setup) for detailed instructions.

### Configuration

Installing the latest vanta plugin will create a config file (`~/.steampipe/config/vanta.spc`) with a single connection named `vanta`:

```hcl
connection "vanta" {
  plugin = "vanta"

  # OAuth client credentials for authenticating with Vanta API
  # To generate OAuth credentials, refer: https://developer.vanta.com/docs/api-access-setup
  # client_id = "vci_jsur8ca2093fb6djsu847528d1629d424941ff545029urj"
  # client_secret = "vcs_jskaoer_kksjded84f8a40d5e64eedeaeolseru813710492300efee0dcff51208f093ujd"

  # Alternatively, you can use an access token instead of client credentials
  # For reference: https://developer.vanta.com/docs/api-access-setup#authentication-and-token-retrieval
  # access_token = "vat_9aa069_Bi3K7v9IoQPMIufU1w4GSJZIh2StgfC0"
}
```

### Authentication Methods

The Vanta plugin supports two authentication methods:

#### Method 1: OAuth Client Credentials (Recommended)

Use OAuth client credentials for production environments and automated workflows:

```hcl
connection "vanta" {
  plugin = "vanta"
  # OAuth client credentials for authenticating with Vanta API
  # To generate OAuth credentials, refer: https://developer.vanta.com/docs/api-access-setup
  client_id = "vci_jsur8ca2093fb6djsu847528d1629d424941ff545029urj"
  client_secret = "vcs_jskaoer_kksjded84f8a40d5e64eedeaeolseru813710492300efee0dcff51208f093ujd"
}
```

#### Method 2: Personal Access Token

Use a access token for development and testing:

```hcl
connection "vanta" {
  plugin = "vanta"
  # Alternatively, you can use a access token instead of client credentials
  access_token = "vat_9aa069_Bi3K7v9IoQPMIufU1w4GSJZIh2StgfC0"
}
```
