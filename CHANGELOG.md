## v1.1.0 [2025-04-17]

_Dependencies_

- Recompiled plugin with Go version `1.23.1`. ([#45](https://github.com/turbot/steampipe-plugin-vanta/pull/45))
- Recompiled plugin with [steampipe-plugin-sdk v5.11.5](https://github.com/turbot/steampipe-plugin-sdk/blob/v5.11.5/CHANGELOG.md#v5115-2025-03-31) that addresses critical and high vulnerabilities in dependent packages. ([#45](https://github.com/turbot/steampipe-plugin-vanta/pull/45))

## v1.0.0 [2024-10-22]

There are no significant changes in this plugin version; it has been released to align with [Steampipe's v1.0.0](https://steampipe.io/changelog/steampipe-cli-v1-0-0) release. This plugin adheres to [semantic versioning](https://semver.org/#semantic-versioning-specification-semver), ensuring backward compatibility within each major version.

_Dependencies_

- Recompiled plugin with Go version `1.22`. ([#40](https://github.com/turbot/steampipe-plugin-vanta/pull/40))
- Recompiled plugin with [steampipe-plugin-sdk v5.10.4](https://github.com/turbot/steampipe-plugin-sdk/blob/develop/CHANGELOG.md#v5104-2024-08-29) that fixes logging in the plugin export tool. ([#40](https://github.com/turbot/steampipe-plugin-vanta/pull/40))

## v0.4.0 [2023-12-12]

_What's new?_

- The plugin can now be downloaded and used with the [Steampipe CLI](https://steampipe.io/docs), as a [Postgres FDW](https://steampipe.io/docs/steampipe_postgres/overview), as a [SQLite extension](https://steampipe.io/docs//steampipe_sqlite/overview) and as a standalone [exporter](https://steampipe.io/docs/steampipe_export/overview). ([#35](https://github.com/turbot/steampipe-plugin-vanta/pull/35))
- The table docs have been updated to provide corresponding example queries for Postgres FDW and SQLite extension. ([#35](https://github.com/turbot/steampipe-plugin-vanta/pull/35))
- Docs license updated to match Steampipe [CC BY-NC-ND license](https://github.com/turbot/steampipe-plugin-vanta/blob/main/docs/LICENSE). ([#35](https://github.com/turbot/steampipe-plugin-vanta/pull/35))

_Dependencies_

- Recompiled plugin with [steampipe-plugin-sdk v5.8.0](https://github.com/turbot/steampipe-plugin-sdk/blob/main/CHANGELOG.md#v580-2023-12-11) that includes plugin server encapsulation for in-process and GRPC usage, adding Steampipe Plugin SDK version to `_ctx` column, and fixing connection and potential divide-by-zero bugs. ([#34](https://github.com/turbot/steampipe-plugin-vanta/pull/34))

## v0.3.2 [2023-10-20]

_Bug fixes_

- Fixed `vanta_computer` table queries failing due to inclusion of deprecated API field `requiresLocationServices` in `fetchDomainEndpoints` query. ([#19](https://github.com/turbot/steampipe-plugin-vanta/pull/19)) (Thanks [@eric-glb](https://github.com/eric-glb) for the contribution!)

## v0.3.1 [2023-10-05]

_Dependencies_

- Recompiled plugin with [steampipe-plugin-sdk v5.6.2](https://github.com/turbot/steampipe-plugin-sdk/blob/main/CHANGELOG.md#v562-2023-10-03) which prevents nil pointer reference errors for implicit hydrate configs. ([#23](https://github.com/turbot/steampipe-plugin-vanta/pull/23))

## v0.3.0 [2023-10-02]

_Dependencies_

- Upgraded to [steampipe-plugin-sdk v5.6.1](https://github.com/turbot/steampipe-plugin-sdk/blob/main/CHANGELOG.md#v561-2023-09-29) with support for rate limiters. ([#20](https://github.com/turbot/steampipe-plugin-vanta/pull/20))
- Recompiled plugin with Go version `1.21`. ([#20](https://github.com/turbot/steampipe-plugin-vanta/pull/20))

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
