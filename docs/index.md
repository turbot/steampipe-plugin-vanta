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

List all apps deployed in your organization:

```sql
select
  display_name,
  policy_type,
  url,
  created_at
from
  vanta_policy;
```

```
+----------------+----------------+----------------------------------+---------------------------+
| display_name   | policy_type    | url                              | created_at                |
+----------------+----------------+----------------------------------+---------------------------+
| Dummy Policy 1 | dummy-policy-1 | https://app.vanta.com/dummy_url1 | 2022-01-13T21:38:39+05:30 |
| Dummy Policy 2 | dummy-policy-2 | https://app.vanta.com/dummy_url2 | 2022-01-13T21:39:50+05:30 |
+----------------+----------------+----------------------------------+---------------------------+
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

| Item        | Description                                                                                                       |
| ----------- | ----------------------------------------------------------------------------------------------------------------- |
| Credentials | Vanta requires an [API token](https://developer.vanta.com/docs/quick-start#1-make-an-api-token) for all requests. |
| Permissions | User requires admin access to generate an API tokens to access the resources.                                     |
| Radius      | Each connection represents a single Vanta Installation.                                                           |
| Resolution  | Credentials explicitly set in a steampipe config file (`~/.steampipe/config/vanta.spc`).                          |

### Configuration

Installing the latest vanta plugin will create a config file (`~/.steampipe/config/vanta.spc`) with a single connection named `vanta`:

```hcl
connection "vanta" {
  plugin = "vanta"

  # A personal API token to access Vanta API
  # To generate an API token, refer: https://developer.vanta.com/docs/quick-start#1-make-an-api-token
  # api_token = "97GtVsdAPwowRToaWDtgZtILdXI_agszONwajQslZ1o"
}
```

## Get involved

- Open source: https://github.com/turbot/steampipe-plugin-vanta
- Community: [Slack Channel](https://steampipe.io/community/join)
