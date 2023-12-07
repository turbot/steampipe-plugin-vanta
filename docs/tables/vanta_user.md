---
title: "Steampipe Table: vanta_user - Query Vanta Users using SQL"
description: "Allows users to query Vanta Users, specifically providing information about each user's ID, email, name, and role. This assists in managing user identities and access within the Vanta platform."
---

# Table: vanta_user - Query Vanta Users using SQL

Vanta is a security and compliance platform that simplifies the complex, time-consuming process of preparing for SOC 2, ISO 27001, and other security audits. It continuously monitors a company's technical infrastructure for security vulnerabilities and non-compliance. Vanta provides a user system where each user has an ID, email, name, and role, which can be queried for identity and access management.

## Table Usage Guide

The `vanta_user` table provides insights into user identities within Vanta. As a security analyst or system administrator, explore user-specific details through this table, including user ID, email, name, and role. Utilize it to manage user identities and access, monitor user activities, and maintain compliance with security standards.

**Important Notes**
- To query the table you must set `session_id` argument in the config file (`~/.steampipe/config/vanta.spc`).

## Examples

### Basic info
Analyze the employment status of users by using their display name and email. This can be useful for understanding the distribution of employment statuses within your user base.

```sql+postgres
select
  display_name,
  id,
  email,
  employment_status
from
  vanta_user;
```

```sql+sqlite
select
  display_name,
  id,
  email,
  employment_status
from
  vanta_user;
```

### List all admins
Identify instances where users have admin permissions. This could be useful for auditing purposes or to ensure that admin privileges are appropriately assigned.

```sql+postgres
select
  display_name,
  id,
  email,
  employment_status
from
  vanta_user
where
  permission_level = 'Admin';
```

```sql+sqlite
select
  display_name,
  id,
  email,
  employment_status
from
  vanta_user
where
  permission_level = 'Admin';
```

### List current employees
Discover the segments that consist of currently employed individuals. This can be useful for understanding the active workforce within your organization.

```sql+postgres
select
  display_name,
  id,
  email,
  employment_status
from
  vanta_user
where
  employment_status = 'CURRENTLY_EMPLOYED';
```

```sql+sqlite
select
  display_name,
  id,
  email,
  employment_status
from
  vanta_user
where
  employment_status = 'CURRENTLY_EMPLOYED';
```

### List inactive users
Explore which users are not currently active in your organization. This can be particularly useful for managing access controls and ensuring security compliance.

```sql+postgres
select
  display_name,
  id,
  email,
  employment_status
from
  vanta_user
where
  employment_status = 'INACTIVE_EMPLOYEE';
```

```sql+sqlite
select
  display_name,
  id,
  email,
  employment_status
from
  vanta_user
where
  employment_status = 'INACTIVE_EMPLOYEE';
```

### List users with security tasks overdue
Discover the segments that consist of users who have pending security tasks. This is crucial to identify potential security risks and ensure timely completion of these tasks.

```sql+postgres
select
  display_name,
  id,
  email,
  employment_status,
  'Due ' || extract(day from (current_timestamp - (task_status_info ->> 'dueDate')::timestamp)) || ' day(s) ago.' as security_task_status
from
  vanta_user
where
  task_status = 'SECURITY_TASKS_OVERDUE';
```

```sql+sqlite
select
  display_name,
  id,
  email,
  employment_status,
  'Due ' || julianday('now') - julianday(json_extract(task_status_info, '$.dueDate')) || ' day(s) ago.' as security_task_status
from
  vanta_user
where
  task_status = 'SECURITY_TASKS_OVERDUE';
```

### List current users by duration of employment
Analyze the duration of employment for your currently active users to gain insights into their tenure within your organization. This can be beneficial for HR planning, such as understanding workforce stability and planning for potential retirements or turnovers.

```sql+postgres
select
  display_name,
  employment_status,
  start_date::date,
  round(extract(day from (current_timestamp - start_date)) / 365, 1) as years
from
  vanta_user
where
  employment_status = 'CURRENTLY_EMPLOYED'
order by
  years desc;
```

```sql+sqlite
select
  display_name,
  employment_status,
  date(start_date),
  round(julianday('now') - julianday(start_date)) / 365.0 as years
from
  vanta_user
where
  employment_status = 'CURRENTLY_EMPLOYED'
order by
  years desc;
```

### Get the count of users by group
Analyze the distribution of users across different groups to understand the user composition in each group. This can be useful for managing user access and permissions, and for understanding the structure of your user base.

```sql+postgres
select
  role ->> 'name' as group_name,
  count(display_name)
from
  vanta_user
group by
  role ->> 'name';
```

```sql+sqlite
select
  json_extract(role, '$.name') as group_name,
  count(display_name)
from
  vanta_user
group by
  json_extract(role, '$.name');
```