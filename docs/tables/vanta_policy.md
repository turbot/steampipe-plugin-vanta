---
title: "Steampipe Table: vanta_policy - Query Vanta Policies using SQL"
description: "Allows users to query Vanta Policies, providing insights into the policy configurations and their associated details."
---

# Table: vanta_policy - Query Vanta Policies using SQL

Vanta is a security and compliance automation platform that simplifies the process of obtaining and maintaining compliance certifications. It automatically collects evidence of a company's security posture, tracks it over time, and streamlines workflows for certification renewals. A Vanta Policy is a set of rules and procedures that define how a company manages and secures its information.

## Table Usage Guide

The `vanta_policy` table provides insights into the policy configurations within Vanta. As a security analyst, explore policy-specific details through this table, including policy names, descriptions, and associated metadata. Utilize it to uncover information about policy configurations, such as policy status, approval dates, and the verification of policy details.

## Examples

### Basic info
Explore the various policies in your system by analyzing their title, status, and approval details.

```sql+postgres
select
  title,
  id,
  status,
  approved_at,
  latest_version_status
from
  vanta_policy;
```

```sql+sqlite
select
  title,
  id,
  status,
  approved_at,
  latest_version_status
from
  vanta_policy;
```

### List policies by status
Identify policies based on their current status to understand the state of your policy management.

```sql+postgres
select
  title,
  id,
  status,
  latest_version_status,
  approved_at
from
  vanta_policy
where
  status is not null
order by
  status, title;
```

```sql+sqlite
select
  title,
  id,
  status,
  latest_version_status,
  approved_at
from
  vanta_policy
where
  status is not null
order by
  status, title;
```

### List recently approved policies
Identify policies that have been approved in the last 90 days to track recent policy updates.

```sql+postgres
select
  title,
  id,
  status,
  approved_at,
  description
from
  vanta_policy
where
  approved_at > (current_timestamp - interval '90 days')
order by
  approved_at desc;
```

```sql+sqlite
select
  title,
  id,
  status,
  approved_at,
  description
from
  vanta_policy
where
  approved_at > datetime('now', '-90 days')
order by
  approved_at desc;
```

### List policies pending approval
Identify policies that need approval based on their latest version status.

```sql+postgres
select
  title,
  id,
  status,
  latest_version_status,
  description
from
  vanta_policy
where
  latest_version_status = 'PENDING_APPROVAL'
order by
  title;
```

```sql+sqlite
select
  title,
  id,
  status,
  latest_version_status,
  description
from
  vanta_policy
where
  latest_version_status = 'PENDING_APPROVAL'
order by
  title;
```

### Count policies by status
Analyze the distribution of policies across different statuses to understand your policy landscape.

```sql+postgres
select
  status,
  count(*) as policy_count
from
  vanta_policy
where
  status is not null
group by
  status
order by
  policy_count desc;
```

```sql+sqlite
select
  status,
  count(*) as policy_count
from
  vanta_policy
where
  status is not null
group by
  status
order by
  policy_count desc;
```

### Count policies by latest version status
Analyze the distribution of policies based on their latest version status.

```sql+postgres
select
  latest_version_status,
  count(*) as policy_count
from
  vanta_policy
where
  latest_version_status is not null
group by
  latest_version_status
order by
  policy_count desc;
```

```sql+sqlite
select
  latest_version_status,
  count(*) as policy_count
from
  vanta_policy
where
  latest_version_status is not null
group by
  latest_version_status
order by
  policy_count desc;
```

### List policies approved within the last year
Explore policies that have been approved within the last year to understand recent policy activity.

```sql+postgres
select
  title,
  id,
  status,
  approved_at,
  extract(days from (current_timestamp - approved_at)) as days_since_approval
from
  vanta_policy
where
  approved_at > (current_timestamp - interval '1 year')
order by
  approved_at desc;
```

```sql+sqlite
select
  title,
  id,
  status,
  approved_at,
  cast(julianday('now') - julianday(approved_at) as integer) as days_since_approval
from
  vanta_policy
where
  approved_at > datetime('now', '-1 year')
order by
  approved_at desc;
```

### List policies with descriptions
Explore policies that have detailed descriptions available.

```sql+postgres
select
  title,
  id,
  status,
  description,
  approved_at
from
  vanta_policy
where
  description is not null
  and length(trim(description)) > 0
order by
  title;
```

```sql+sqlite
select
  title,
  id,
  status,
  description,
  approved_at
from
  vanta_policy
where
  description is not null
  and length(trim(description)) > 0
order by
  title;
```
