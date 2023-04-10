## v0.2.0 [2023-04-11]

_Dependencies_

- Recompiled plugin with [steampipe-plugin-sdk v5.3.0](https://github.com/turbot/steampipe-plugin-sdk/blob/main/CHANGELOG.md#v530-2023-03-16) which includes fixes for query cache pending item mechanism and aggregator connections not working for dynamic tables. ([#12](https://github.com/turbot/steampipe-plugin-vanta/pull/12))

## v0.1.0 [2023-02-10]

_Enhancements_

- Added column `endpoint_applications` to `vanta_computer` table. ([#9](https://github.com/turbot/steampipe-plugin-vanta/pull/9))

_Dependencies_

- Recompiled plugin with [steampipe-plugin-sdk v5.1.3](https://github.com/turbot/steampipe-plugin-sdk/blob/main/CHANGELOG.md#v513-2023-02-09) which fixes the issue pertaining to query caching functionality. ([#10](https://github.com/turbot/steampipe-plugin-vanta/pull/10))

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