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

[Steampipe](https://steampipe.io) is an open source CLI to instantly query cloud APIs using SQL.

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

The plugin uses two different endpoints that uses different credential mechanism

| Item        | Description                                                                                                                                                                                                                                                                                                                                               |
| ----------- | --------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| Credentials | The plugin uses two different endpoints that use different credential mechanisms:<br/>1. Using a user's personal [API token](https://developer.vanta.com/docs/quick-start#1-make-an-api-token).<br/>2. Using the [cookie-based authentication](#getting-the-session-id-for-cookie-based-authentication) by passing a unique session ID for every request. |
| Permissions | User requires admin access to generate an API token to access the resources.                                                                                                                                                                                                                                                                             |
| Radius      | Each connection represents a single Vanta installation.                                                                                                                                                                                                                                                                                                   |
| Resolution  | Credentials explicitly set in a steampipe config file (`~/.steampipe/config/vanta.spc`).                                                                                                                                                                                                                                                                  |

### Configuration

Installing the latest vanta plugin will create a config file (`~/.steampipe/config/vanta.spc`) with a single connection named `vanta`:

```hcl
connection "vanta" {
  plugin = "vanta"

  # A personal API token to access Vanta API
  # This is only required while querying `vanta_evidence` table.
  # To generate an API token, refer: https://developer.vanta.com/docs/quick-start#1-make-an-api-token
  # api_token = "97GtVsdAPwowRToaWDtgZtILdXI_agszONwajQslZ1o"

  # Session ID of your current vanta session
  # Set the value of `connect.sid` cookie from a logged in Vanta browser session
  # Required to access tables that are using the https://app.vanta.com/graphql endpoint
  # session_id = "s:3nZSteamPipe1fSu4iNV_1TB5UTesTToGK.zVANtaplugintest+GVxPvQffhnFY3skWlfkceZxXKSCjc"
}
```

### Getting the Session ID for cookie-based authentication

The Vanta APIs generally use a user's personal [API token](https://developer.vanta.com/docs/quick-start#1-make-an-api-token) to authenticate the requests. But some of the tables in this plugin use a different endpoint, which requires a unique session ID to access the query endpoint.

To retrieve your Session ID:

- Log into the Vanta console.
- Open your browser [developer tools](https://developer.mozilla.org/en-US/docs/Learn/Common_questions/What_are_browser_developer_tools).
- Open the `Network` view to see and analyze the network requests that make up each individual page load within a single user's session.
- Open any `graphql` request from the list and check the `Cookies` section to get the list of request cookies.
- Get the session ID value from the list named as `connect.sid`.


