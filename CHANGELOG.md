## v0.0.2 [2023-02-06]

_Bug fixes_

- Fixed the `vanta_integration` table to correctly return data instead of returning an error. ([#5](https://github.com/turbot/steampipe-plugin-vanta/pull/5))
- Renamed column `risk_attributes` to `risk_profile` in `vanta_vendor` table to stay in alignment with the GraphQL API attributes. ([#5](https://github.com/turbot/steampipe-plugin-vanta/pull/5))

## v0.0.1 [2023-02-04]

_What's new?_

- New tables added
  - [vanta_computer](https://hub.steampipe.io/plugins/turbot/vanta/tables/vanta_computer)
  - [vanta_evidence](https://hub.steampipe.io/plugins/turbot/vanta/tables/vanta_evidence)
  - [vanta_group](https://hub.steampipe.io/plugins/turbot/vanta/tables/vanta_group)
  - [vanta_integration](https://hub.steampipe.io/plugins/turbot/vanta/tables/vanta_integration)
  - [vanta_monitor](https://hub.steampipe.io/plugins/turbot/vanta/tables/vanta_monitor)
  - [vanta_policy](https://hub.steampipe.io/plugins/turbot/vanta/tables/vanta_policy)
  - [vanta_user](https://hub.steampipe.io/plugins/turbot/vanta/tables/vanta_user)
  - [vanta_vendor](https://hub.steampipe.io/plugins/turbot/vanta/tables/vanta_vendor)