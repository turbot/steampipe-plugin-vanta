---
title: "Steampipe Table: vanta_user - Query Vanta Users using SQL"
description: "Allows users to query Vanta Users, specifically providing information about each user's ID, email, name, employment status, and other employment details."
---

# Table: vanta_user - Query Vanta Users using SQL

Vanta is a security and compliance platform that simplifies the complex, time-consuming process of preparing for SOC 2, ISO 27001, and other security audits. It continuously monitors a company's technical infrastructure for security vulnerabilities and non-compliance. Vanta provides a user system where each user has an ID, email, name, employment status, and other employment-related details, which can be queried for identity and access management.

## Table Usage Guide

The `vanta_user` table provides insights into user identities within Vanta. As a security analyst or system administrator, explore user-specific details through this table, including user ID, email, name, employment status, and job details. Utilize it to manage user identities and access, monitor user activities, and maintain compliance with security standards.

## Examples

### Basic info
Explore user information including their display name, email, and employment status.

```sql+postgres
select
  display_name,
  id,
  email,
  employment_status,
  job_title
from
  vanta_user;
```

```sql+sqlite
select
  display_name,
  id,
  email,
  employment_status,
  job_title
from
  vanta_user;
```

### List current employees
Discover the segments that consist of currently employed individuals. This can be useful for understanding the active workforce within your organization.

```sql+postgres
select
  display_name,
  id,
  email,
  job_title,
  start_date
from
  vanta_user
where
  is_active = true;
```

```sql+sqlite
select
  display_name,
  id,
  email,
  job_title,
  start_date
from
  vanta_user
where
  is_active = 1;
```

### List inactive users
Explore which users are not currently active in your organization. This can be particularly useful for managing access controls and ensuring security compliance.

```sql+postgres
select
  display_name,
  id,
  email,
  employment_status,
  end_date
from
  vanta_user
where
  is_active = false;
```

```sql+sqlite
select
  display_name,
  id,
  email,
  employment_status,
  end_date
from
  vanta_user
where
  is_active = 0;
```

### List users with overdue security tasks
Discover users who have pending security tasks that need immediate attention.

```sql+postgres
select
  display_name,
  id,
  email,
  employment_status,
  task_status,
  job_title
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
  task_status,
  job_title
from
  vanta_user
where
  task_status = 'SECURITY_TASKS_OVERDUE';
```

### List current users by duration of employment
Analyze the duration of employment for your currently active users to gain insights into their tenure within your organization.

```sql+postgres
select
  display_name,
  employment_status,
  job_title,
  start_date::date,
  round(extract(day from (current_timestamp - start_date)) / 365.0, 1) as years_employed
from
  vanta_user
where
  is_active = true
  and start_date is not null
order by
  years_employed desc;
```

```sql+sqlite
select
  display_name,
  employment_status,
  job_title,
  date(start_date),
  round((julianday('now') - julianday(start_date)) / 365.0, 1) as years_employed
from
  vanta_user
where
  is_active = 1
  and start_date is not null
order by
  years_employed desc;
```

### Count users by employment status
Analyze the distribution of users across different employment statuses.

```sql+postgres
select
  employment_status,
  count(*) as user_count
from
  vanta_user
where
  employment_status is not null
group by
  employment_status
order by
  user_count desc;
```

```sql+sqlite
select
  employment_status,
  count(*) as user_count
from
  vanta_user
where
  employment_status is not null
group by
  employment_status
order by
  user_count desc;
```

### Count users by task status
Analyze the distribution of users based on their security task completion status.

```sql+postgres
select
  task_status,
  count(*) as user_count
from
  vanta_user
where
  task_status is not null
group by
  task_status
order by
  user_count desc;
```

```sql+sqlite
select
  task_status,
  count(*) as user_count
from
  vanta_user
where
  task_status is not null
group by
  task_status
order by
  user_count desc;
```

### List users by job title
Explore the organizational structure by listing users grouped by their job titles.

```sql+postgres
select
  job_title,
  count(*) as employee_count,
  array_agg(display_name order by display_name) as employees
from
  vanta_user
where
  job_title is not null
  and is_active = true
group by
  job_title
order by
  employee_count desc;
```

```sql+sqlite
select
  job_title,
  count(*) as employee_count,
  group_concat(display_name, ', ') as employees
from
  vanta_user
where
  job_title is not null
  and is_active = 1
group by
  job_title
order by
  employee_count desc;
```

### List recently hired employees
Identify employees who have joined the organization in the last 90 days.

```sql+postgres
select
  display_name,
  email,
  job_title,
  start_date,
  extract(day from (current_timestamp - start_date)) as days_since_start
from
  vanta_user
where
  start_date > (current_timestamp - interval '90 days')
  and is_active = true
order by
  start_date desc;
```

```sql+sqlite
select
  display_name,
  email,
  job_title,
  start_date,
  cast(julianday('now') - julianday(start_date) as integer) as days_since_start
from
  vanta_user
where
  start_date > datetime('now', '-90 days')
  and is_active = 1
order by
  start_date desc;
```
