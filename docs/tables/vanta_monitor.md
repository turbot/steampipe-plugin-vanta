---
title: "Steampipe Table: vanta_monitor - Query Vanta Monitors using SQL"
description: "Allows users to query Vanta Monitors, providing insights into the status, type, and details of each monitor."
---

# Table: vanta_monitor - Query Vanta Monitors using SQL

Vanta is a security and compliance automation platform. It simplifies the complex and time-consuming process of preparing for SOC 2, HIPAA, and ISO 27001 audits. Vanta provides continuous monitoring of your applications, infrastructure, and cloud services to ensure they adhere to security best practices.

## Table Usage Guide

The `vanta_monitor` table provides insights into the monitors within Vanta's security and compliance automation platform. As a security analyst, explore monitor-specific details through this table, including status, type, and associated metadata. Utilize it to uncover information about monitors, such as those with alerts, the type of monitors, and the verification of monitor details.

**Important Notes**
- To query the table you must set `session_id` argument in the config file (`~/.steampipe/config/vanta.spc`).

## Examples

### Basic info
Explore the status and outcome of different categories within a monitoring system. This can help you understand the areas requiring attention and the effectiveness of remediation efforts.

```sql
select
  name,
  category,
  outcome,
  latest_flip_time,
  remediation_status ->> 'status' as status
from
  vanta_monitor;
```

### List all failed tests
Explore which monitors have failed tests to assess the areas of non-compliance and understand when the last status change occurred. This allows you to pinpoint specific issues and address them promptly.

```sql
select
  name,
  category,
  outcome,
  compliance_status,
  latest_flip_time
from
  vanta_monitor
where
  outcome = 'FAIL';
```

### Filter a specific test result by test ID
Explore the details of a specific test result by using its test ID, allowing you to gain insights into its category, outcome, and compliance status, as well as the time of the most recent status change. This is particularly useful for tracking and reviewing the performance of individual tests over time.

```sql
select
  name,
  category,
  outcome,
  compliance_status,
  latest_flip_time
from
  vanta_monitor
where
  test_id = 'test-03neqwol876pxg1iqjqib9';
```

### Count tests by remediation status
Analyze the status of remediation efforts by quantifying the number of tests associated with each status. This can help in prioritizing and tracking remediation tasks effectively.

```sql
select
  remediation_status ->> 'status' as status,
  count(name)
from
  vanta_monitor
group by
  remediation_status ->> 'status';
```

### List failed tests by owner
Explore which tests have failed and identify the owners responsible for these tests. This is useful for assessing the areas that need immediate attention or remediation.

```sql
select
  name,
  category,
  a ->> 'displayName' as owner,
  remediation_status ->> 'status' as status
from
  vanta_monitor
  left join jsonb_array_elements(assignees) as a on true
where
  outcome = 'FAIL';
```

### List failed tests by standard
Discover the segments that have failed tests according to different standards. This query can be used to assess the status of remediation and identify areas needing attention in order to meet various standards.

```sql
select
  name,
  category,
  s -> 'standardInfo' ->> 'standard' as standard,
  remediation_status ->> 'status' as status
from
  vanta_monitor,
  jsonb_array_elements(controls) as c,
  jsonb_array_elements(c -> 'standardSections') as s
where
  outcome = 'FAIL';
```

### List failed tests by integration
Explore which integrations have failed tests, allowing you to identify areas of concern and take necessary corrective actions. This is useful in maintaining system integrity and ensuring seamless integration performance.

```sql
select
  m.name,
  m.category,
  m.outcome,
  i.display_name as integration
from
  vanta_integration as i,
  jsonb_array_elements(i.tests) as t
  join vanta_monitor as m on m.test_id = t ->> 'testId' and m.outcome = 'FAIL';
```

### Count tests by outcome
Assess the distribution of test outcomes within your Vanta monitor system. This query is useful for understanding the frequency of different outcomes, helping to identify patterns or areas for improvement.

```sql
select
  outcome,
  count(name)
from
  vanta_monitor
group by
  outcome;
```

### Count active tests by category
Analyze the settings to understand the distribution of active tests across different categories. This can help in identifying the areas that are being frequently tested and those that require more attention.

```sql
select
  category,
  count(name)
from
  vanta_monitor
where
  disabled_status is null
group by
  category;
```