# Table: vanta_user

The `vanta_user` table can be used to query information about all users in the organization.

**NOTE:**

- To query the table; **you must set** `session_id` argument in the config file (`~/.steampipe/config/vanta.spc`).

## Examples

### Basic info

```sql
select
  display_name,
  id,
  email,
  employment_status
from
  vanta_user;
```

### List all admins

```sql
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

```sql
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

```sql
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

```sql
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

### List current users by duration of employment

```sql
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

### Get the count of users by group

```sql
select
  role ->> 'name' as group_name,
  count(display_name)
from
  vanta_user
group by
  role ->> 'name';
```
