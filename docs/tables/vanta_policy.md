---
title: "Steampipe Table: vanta_policy - Query Vanta Policies using SQL"
description: "Allows users to query Vanta Policies, providing insights into the policy configurations and their associated details."
---

# Table: vanta_policy - Query Vanta Policies using SQL

Vanta is a security and compliance automation platform that simplifies the process of obtaining and maintaining compliance certifications. It automatically collects evidence of a company's security posture, tracks it over time, and streamlines workflows for certification renewals. A Vanta Policy is a set of rules and procedures that define how a company manages and secures its information.

## Table Usage Guide

The `vanta_policy` table provides insights into the policy configurations within Vanta. As a security analyst, explore policy-specific details through this table, including policy names, descriptions, and associated metadata. Utilize it to uncover information about policy configurations, such as policy status, the type of policy, and the verification of policy details.

## Examples

### Basic info
Explore the various policies in your system by analyzing their title, type, status, and creation date. This can help you understand the range and scope of your current policies, as well as identify any gaps or inconsistencies.

```sql
select
  title,
  policy_type,
  status,
  created_at,
  standards
from
  vanta_policy;
```

### List unapproved policies
Identify policies that are pending approval to ensure timely review and validation for maintaining security compliance. This query is useful in managing organizational security by highlighting areas that need immediate attention.

```sql
select
  title,
  policy_type,
  status,
  created_at
from
  vanta_policy
where
  approver is null;
```

### List expired policies
Discover the segments that have policies which have expired. This is useful in understanding which areas need immediate attention for policy renewal, ensuring compliance and reducing risk.

```sql
select
  title,
  policy_type,
  status,
  created_at,
  approver ->> 'displayName' as approver
from
  vanta_policy
where
  (approved_at + interval '1 year') < current_timestamp;
```

### List policies expiring in the next 30 days
Determine the policies that are due to expire in the next 30 days. This can be useful for administrators to proactively manage policy renewals and avoid any lapses in coverage.

```sql
select
  title,
  policy_type,
  status,
  created_at,
  approver ->> 'displayName' as approver,
  'expires in ' || extract(day from ((approved_at + interval '1 year') - current_timestamp)) || ' day(s)' as status
from
  vanta_policy
where
  current_timestamp < (approved_at + interval '1 year')
  and extract(day from ((approved_at + interval '1 year') - current_timestamp)) <= '30';
```

### List users who have not accepted a specific policy
Determine the users who have not accepted a specific policy, such as a 'Code of Conduct'. This can be useful for ensuring all team members are in compliance with company policies.

```sql
with policy_summary as (
  select
    p.title as policy_name,
    p.status as policy_status,
    p.approved_at,
    p.approver ->> 'displayName' as approver,
    m.failing_resource_entities
  from
    vanta_policy as p
    join vanta_monitor as m on m.test_id = p.employee_acceptance_test_id
  where
    title = 'Code of Conduct'
  order by policy_type
)
select
  p.policy_name,
  f ->> 'displayName' as user_name,
  u.email
from
  policy_summary as p,
  jsonb_array_elements(p.failing_resource_entities) as f
  join vanta_user as u on u.display_name = f ->> 'displayName'
where
  f ->> '__typename' = 'User';
```