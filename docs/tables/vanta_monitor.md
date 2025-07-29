---
title: "Steampipe Table: vanta_monitor - Query Vanta Monitors using SQL"
description: "Allows users to query Vanta Monitors, providing insights into the status, type, and details of each monitor."
---

# Table: vanta_monitor - Query Vanta Monitors using SQL

Vanta is a security and compliance automation platform. It simplifies the complex and time-consuming process of preparing for SOC 2, HIPAA, and ISO 27001 audits. Vanta provides continuous monitoring of your applications, infrastructure, and cloud services to ensure they adhere to security best practices.

## Table Usage Guide

The `vanta_monitor` table provides insights into the monitors within Vanta's security and compliance automation platform. As a security analyst, explore monitor-specific details through this table, including status, type, and associated metadata. Utilize it to uncover information about monitors, such as those with alerts, the type of monitors, and the verification of monitor details.

## Examples

### Basic info
Explore the status and details of different monitors within your monitoring system.

```sql+postgres
select
  name,
  category,
  status,
  latest_flip_date,
  remediation_status
from
  vanta_monitor;
```

```sql+sqlite
select
  name,
  category,
  status,
  latest_flip_date,
  remediation_status
from
  vanta_monitor;
```

### List all failed tests
Explore which monitors are in a failing state to assess areas of non-compliance and understand when the last status change occurred.

```sql+postgres
select
  name,
  category,
  status,
  latest_flip_date,
  owner_display_name,
  remediation_item_count
from
  vanta_monitor
where
  status = 'NEEDS_ATTENTION';
```

```sql+sqlite
select
  name,
  category,
  status,
  latest_flip_date,
  owner_display_name,
  remediation_item_count
from
  vanta_monitor
where
  status = 'NEEDS_ATTENTION';
```

### Filter monitors by specific test ID
Explore the details of a specific monitor by using its test ID.

```sql+postgres
select
  name,
  category,
  status,
  latest_flip_date,
  description,
  failure_description
from
  vanta_monitor
where
  id = 'test-03neqwol876pxg1iqjqib9';
```

```sql+sqlite
select
  name,
  category,
  status,
  latest_flip_date,
  description,
  failure_description
from
  vanta_monitor
where
  id = 'test-03neqwol876pxg1iqjqib9';
```

### Count monitors by status
Analyze the distribution of monitor statuses to understand your overall compliance health.

```sql+postgres
select
  status,
  count(*) as monitor_count
from
  vanta_monitor
group by
  status
order by
  monitor_count desc;
```

```sql+sqlite
select
  status,
  count(*) as monitor_count
from
  vanta_monitor
group by
  status
order by
  monitor_count desc;
```

### List failed monitors by owner
Explore which monitors have failed and identify the owners responsible for these tests.

```sql+postgres
select
  name,
  category,
  status,
  owner_display_name,
  owner_email,
  remediation_item_count
from
  vanta_monitor
where
  status = 'NEEDS_ATTENTION'
  and owner_display_name is not null
order by
  owner_display_name;
```

```sql+sqlite
select
  name,
  category,
  status,
  owner_display_name,
  owner_email,
  remediation_item_count
from
  vanta_monitor
where
  status = 'NEEDS_ATTENTION'
  and owner_display_name is not null
order by
  owner_display_name;
```

### List monitors with remediation in progress
Identify monitors that currently have active remediation efforts.

```sql+postgres
select
  name,
  category,
  status,
  remediation_status,
  remediation_item_count,
  owner_display_name
from
  vanta_monitor
where
  remediation_status = 'IN_PROGRESS'
order by
  remediation_item_count desc;
```

```sql+sqlite
select
  name,
  category,
  status,
  remediation_status,
  remediation_item_count,
  owner_display_name
from
  vanta_monitor
where
  remediation_status = 'IN_PROGRESS'
order by
  remediation_item_count desc;
```

### List deactivated monitors
Discover monitors that have been deactivated and understand the reasons.

```sql+postgres
select
  name,
  category,
  status,
  is_deactivated,
  deactivated_reason,
  latest_flip_date
from
  vanta_monitor
where
  is_deactivated = true;
```

```sql+sqlite
select
  name,
  category,
  status,
  is_deactivated,
  deactivated_reason,
  latest_flip_date
from
  vanta_monitor
where
  is_deactivated = 1;
```

### Count active monitors by category
Analyze the distribution of active monitors across different categories.

```sql+postgres
select
  category,
  count(*) as active_monitor_count
from
  vanta_monitor
where
  is_deactivated = false
group by
  category
order by
  active_monitor_count desc;
```

```sql+sqlite
select
  category,
  count(*) as active_monitor_count
from
  vanta_monitor
where
  is_deactivated = 0
group by
  category
order by
  active_monitor_count desc;
```

### List monitors with failing resource entities
Explore monitors that have failing resources requiring attention.

```sql+postgres
select
  name,
  category,
  status,
  jsonb_array_length(failing_resource_entities) as failing_entities_count,
  owner_display_name
from
  vanta_monitor
where
  failing_resource_entities is not null
  and jsonb_array_length(failing_resource_entities) > 0
order by
  failing_entities_count desc;
```

```sql+sqlite
select
  name,
  category,
  status,
  json_array_length(failing_resource_entities) as failing_entities_count,
  owner_display_name
from
  vanta_monitor
where
  failing_resource_entities is not null
  and json_array_length(failing_resource_entities) > 0
order by
  failing_entities_count desc;
```
