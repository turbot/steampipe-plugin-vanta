---
organization: Turbot
category: ["saas"]
icon_url: "/images/plugins/turbot/vanta.svg"
brand_color: "#5230D7"
display_name: "Vanta"
short_name: "vanta"
description: "Steampipe plugin to query users, policies, compliances, and more from your Vanta organization."
og_description: "Query Vanta with SQL! Open source CLI. No DB required."
og_image: "/images/plugins/turbot/vanta-social-graphic.png"
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
| Credentials | The plugin uses two different endpoints that uses different credential mechanism:<br/>1. Using a user's personal [API token](https://developer.vanta.com/docs/quick-start#1-make-an-api-token).<br/>2. Using the [cookie-based authentication](#getting-the-session-id-for-cookie-based-authentication) by passing a unique session ID for every request. |
| Permissions | User requires admin access to generate an API tokens to access the resources.                                                                                                                                                                                                                                                                             |
| Radius      | Each connection represents a single Vanta Installation.                                                                                                                                                                                                                                                                                                   |
| Resolution  | Credentials explicitly set in a steampipe config file (`~/.steampipe/config/vanta.spc`).                                                                                                                                                                                                                                                                  |

### Configuration

Installing the latest vanta plugin will create a config file (`~/.steampipe/config/vanta.spc`) with a single connection named `vanta`:

```hcl
connection "vanta" {
  plugin = "vanta"

  # A personal API token to access Vanta API
  # To generate an API token, refer: https://developer.vanta.com/docs/quick-start#1-make-an-api-token
  # api_token = "97GtVsdAPwowRToaWDtgZtILdXI_agszONwajQslZ1o"

  # Session Id of your current vanta session
  # Required to access tables that using the https://app.vanta.com/graphql endpoint
  # session_id = "s:3nZSteamPipe1fSu4iNV_1TB5UTesTToGK.zVANtaplugintest+GVxPvQffhnFY3skWlfkceZxXKSCjc"
}
```

### Getting the Session ID for cookie-based authentication

The Vanta APIs generally uses an user's personal API token to authenticate the requests. But some the tables in this plugin uses a different endpoint to which requires a unique session ID to access the query endpoint.

To locate your Session ID:

- Log into the Vanta console.
- Open your browser [developer tools](https://developer.mozilla.org/en-US/docs/Learn/Common_questions/What_are_browser_developer_tools).
- Open the `Network` view to see and analyze the network requests that make up each individual page load within a single user's session.
- Open up any `graphql` request from the list and check the `Cookies` section to get the list of request cookies.
- Get the session ID value from list named as `connect.sid`.

## Get involved

- Open source: https://github.com/turbot/steampipe-plugin-vanta
- Community: [Slack Channel](https://steampipe.io/community/join)
