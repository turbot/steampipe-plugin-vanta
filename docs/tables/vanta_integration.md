---
title: "Steampipe Table: vanta_integration - Query Vanta Integrations using SQL"
description: "Allows users to query Vanta Integrations, providing insights into the various integrations within the Vanta security system."
---

# Table: vanta_integration - Query Vanta Integrations using SQL

Vanta is a security and compliance automation platform. It simplifies the complex process of preparing for SOC 2, ISO 27001, and other security audits. Vanta Integrations are the various systems, applications, and services that Vanta connects to, in order to collect and analyze security-related data.

## Table Usage Guide

The `vanta_integration` table provides insights into the various integrations within the Vanta security system. As a security analyst, explore integration-specific details through this table, including status, type, and associated metadata. Utilize it to uncover information about integrations, such as their current status, the type of integration, and other critical data.

**Important Notes**
- There are various integrations available that can be integrated. The table `vanta_integration` only returns the integrations that are connected.
- To query the table you must set `session_id` argument in the config file (`~/.steampipe/config/vanta.spc`).

## Examples

### Basic info
Explore the specific details of your integrated applications, such as their names, unique identifiers, and associated logos. This can assist in managing and tracking your integrations more effectively.

```sql+postgres
select
  display_name,
  id,
  description,
  logo_slug_id
from
  vanta_integration;
```

```sql+sqlite
select
  display_name,
  id,
  description,
  logo_slug_id
from
  vanta_integration;
```

### List integrations having disabled credentials
Identify instances where certain integrations have been disabled. This is useful in maintaining system security and functionality by quickly pinpointing any inactive credentials.

```sql+postgres
select
  display_name,
  id,
  (c ->> 'metadata')::jsonb as credential_metadata,
  c ->> 'service'
from
  vanta_integration,
  jsonb_array_elements(credentials) as c
where
  c ->> 'isDisabled' = 'true';
```

```sql+sqlite
select
  display_name,
  id,
  json_extract(c.value, '$.metadata') as credential_metadata,
  json_extract(c.value, '$.service')
from
  vanta_integration,
  json_each(credentials) as c
where
  json_extract(c.value, '$.isDisabled') = 'true';
```

### List integrations with failed tests
Assess the elements within your system integrations to identify instances where tests have failed. This can be beneficial in pinpointing specific areas of concern and taking corrective actions to improve system performance.

```sql+postgres
select
  i.display_name,
  i.id,
  t ->> 'testId' as test_id,
  m.outcome as test_status
from
  vanta_integration as i,
  jsonb_array_elements(tests) as t
  join vanta_monitor as m on m.test_id = t ->> 'testId' and outcome = 'FAIL';
```

```sql+sqlite
select
  i.display_name,
  i.id,
  json_extract(t.value, '$.testId') as test_id,
  m.outcome as test_status
from
  vanta_integration as i,
  json_each(tests) as t
  join vanta_monitor as m on m.test_id = json_extract(t.value, '$.testId') and outcome = 'FAIL';
```