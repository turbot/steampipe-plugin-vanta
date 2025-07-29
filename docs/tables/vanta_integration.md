---
title: "Steampipe Table: vanta_integration - Query Vanta Integrations using SQL"
description: "Allows users to query Vanta Integrations, providing insights into the various connected integrations within the Vanta security system."
---

# Table: vanta_integration - Query Vanta Integrations using SQL

Vanta is a security and compliance automation platform. It simplifies the complex process of preparing for SOC 2, ISO 27001, and other security audits. Vanta Integrations are the various systems, applications, and services that Vanta connects to, in order to collect and analyze security-related data.

## Table Usage Guide

The `vanta_integration` table provides insights into the various connected integrations within the Vanta security system. As a security analyst, explore integration-specific details through this table, including connection status, resource kinds, and associated metadata. Utilize it to uncover information about integrations, such as their current connections, available resource types, and associated tests.

**Important Notes**

- There are various integrations available that can be integrated. The table `vanta_integration` only returns the integrations that are connected.

## Examples

### Basic info
Explore the specific details of your connected integrations, such as their names and unique identifiers.

```sql+postgres
select
  display_name,
  id,
  connections,
  scopable_resource
from
  vanta_integration;
```

```sql+sqlite
select
  display_name,
  id,
  connections,
  scopable_resource
from
  vanta_integration;
```

### List integrations with multiple connections
Identify integrations that have multiple connections configured, which may indicate complex setups or redundancy.

```sql+postgres
select
  display_name,
  id,
  jsonb_array_length(connections) as connection_count,
  connections
from
  vanta_integration
where
  jsonb_array_length(connections) > 1
order by
  connection_count desc;
```

```sql+sqlite
select
  display_name,
  id,
  json_array_length(connections) as connection_count,
  connections
from
  vanta_integration
where
  json_array_length(connections) > 1
order by
  connection_count desc;
```

### List integrations by connection status
Explore the status of integration connections to identify any that may need attention.

```sql+postgres
select
  i.display_name,
  i.id,
  c ->> 'status' as connection_status,
  c ->> 'displayName' as connection_name
from
  vanta_integration as i,
  jsonb_array_elements(i.connections) as c
order by
  connection_status, display_name;
```

```sql+sqlite
select
  i.display_name,
  i.id,
  json_extract(c.value, '$.status') as connection_status,
  json_extract(c.value, '$.displayName') as connection_name
from
  vanta_integration as i,
  json_each(i.connections) as c
order by
  connection_status, display_name;
```

### List integrations with available resource kinds
Explore which integrations provide specific types of scopable resources for monitoring.

```sql+postgres
select
  display_name,
  id,
  r as resource_kind
from
  vanta_integration,
  jsonb_array_elements_text(scopable_resource) as r
order by
  display_name, resource_kind;
```

```sql+sqlite
select
  display_name,
  id,
  json_extract(r.value, '$') as resource_kind
from
  vanta_integration,
  json_each(scopable_resource) as r
where
  json_valid(scopable_resource)
order by
  display_name, resource_kind;
```

### Count integrations by resource type availability
Analyze which resource types are most commonly available across your integrations.

```sql+postgres
select
  r as resource_type,
  count(*) as integration_count
from
  vanta_integration,
  jsonb_array_elements_text(scopable_resource) as r
group by
  r
order by
  integration_count desc;
```

```sql+sqlite
select
  json_extract(r.value, '$') as resource_type,
  count(*) as integration_count
from
  vanta_integration,
  json_each(scopable_resource) as r
where
  json_valid(scopable_resource)
group by
  json_extract(r.value, '$')
order by
  integration_count desc;
```

### List integrations with tests
Explore integrations that have associated tests for monitoring and compliance purposes.

```sql+postgres
select
  display_name,
  id,
  jsonb_array_length(tests) as test_count
from
  vanta_integration
where
  tests is not null
  and jsonb_array_length(tests) > 0
order by
  test_count desc;
```

```sql+sqlite
select
  display_name,
  id,
  json_array_length(tests) as test_count
from
  vanta_integration
where
  tests is not null
  and json_array_length(tests) > 0
order by
  test_count desc;
```

### Get integration test details
Explore the specific tests associated with each integration for monitoring purposes.

```sql+postgres
select
  i.display_name,
  i.id,
  t ->> 'testId' as test_id,
  t ->> 'displayName' as test_name
from
  vanta_integration as i,
  jsonb_array_elements(i.tests) as t
where
  i.tests is not null
order by
  i.display_name, test_name;
```

```sql+sqlite
select
  i.display_name,
  i.id,
  json_extract(t.value, '$.testId') as test_id,
  json_extract(t.value, '$.displayName') as test_name
from
  vanta_integration as i,
  json_each(i.tests) as t
where
  i.tests is not null
order by
  i.display_name, test_name;
```
