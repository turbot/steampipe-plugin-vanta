## v2.0.1 [2025-08-14]

_Bug fixes_

- Updated module path to v2 and adjusted import statements for consistency.

## v2.0.0 [2025-08-14]

_Breaking changes_

### Configuration Changes
- The plugin configuration format has changed to support the REST API.
- You must now configure either **OAuth client credentials** or an **access token** for authentication. Please refer the [Configuration](https://hub.steampipe.io/plugins/turbot/vanta#configuration) section for additional information.

### API Migration
- Migrated from Vanta's deprecated **GraphQL** API to the new **REST** API.
- Queries, dashboards, and benchmarks that reference removed columns (listed below) will fail until updated.

### Removed Columns

**`vanta_computer`**
- `agent_version`
- `hostname`
- `host_identifier`
- `last_ping`
- `num_browser_extensions`
- `endpoint_applications`
- `installed_av_programs`
- `installed_password_managers`
- `unsupported_reasons`
- `organization_name`

**`vanta_evidence`**
- `title`
- `evidence_request_id`
- `category`
- `uid`
- `app_upload_enabled`
- `restricted`
- `dismissed_status`
- `renewal_metadata`
- `organization_name`

**`vanta_group`**
- `checklist`
- `embedded_idp_group`
- `organization_name`

**`vanta_integration`**
- `description`
- `application_url`
- `installation_url`
- `logo_slug_id`
- `credentials`
- `integration_categories`
- `service_categories`
- `organization_name`

**`vanta_monitor`**
- `controls`
- `organization_name`

**`vanta_policy`**
- `policy_type`
- `created_at`
- `updated_at`
- `employee_acceptance_test_id`
- `num_users`
- `num_users_accepted`
- `source`
- `acceptance_controls`
- `approver`
- `standards`
- `uploaded_doc`
- `uploader`
- `organization_name`

**`vanta_user`**
- `is_from_scan`
- `needs_employee_digest_reminder`
- `is_not_human`

**`vanta_vendor`**
- `vendor_risk_locked`
- `owner`
- `risk_profile`
- `organization_name`

### Migration Notes
- **Reason for change:** Vanta has ended GraphQL API support; the REST API is now the only supported interface.
- **Action required:**
  - Update SQL queries to remove or replace references to removed columns.
  - Review table documentation for updated field availability.

## v1.1.0 [2025-04-18]

_Breaking changes_

The GraphQL API has deprecated the following fields from the tables, which have been removed from the plugin: ([#44](https://github.com/turbot/steampipe-plugin-vanta/pull/44))
- Removed the `permission_level` column from the `vanta_user` table.  
- Removed the following columns from the `vanta_policy` table: `num_users`, `num_users_accepted`, and `source`.  
- Removed the following columns from the `vanta_vendor` table: `assessment_documents`, `latest_security_review`, `severity`, `services_provided`, `shares_credit_card_data`, `submitted_vaqs` and `vendor_category`.

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
