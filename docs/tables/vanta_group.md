---
title: "Steampipe Table: vanta_group - Query Vanta Groups using SQL"
description: "Allows users to query Groups in Vanta, specifically the group details and associated users, providing insights into group management and user assignments."
---

# Table: vanta_group - Query Vanta Groups using SQL

Vanta is a security monitoring platform that simplifies the complex process of security compliance. It provides comprehensive visibility into an organization's security posture, helping to identify and mitigate potential vulnerabilities. Vanta's Groups feature allows for the management of user permissions, providing a structured way to assign and control access rights.

## Table Usage Guide

The `vanta_group` table provides insights into Groups within Vanta's security monitoring platform. As a Security or Compliance Officer, explore group-specific details through this table, including group names, user assignments, and associated permissions. Utilize it to uncover information about groups, such as those with high-level permissions, the distribution of user assignments among groups, and the verification of access rights.

**Important Notes**
- To query the table you must set `session_id` argument in the config file (`~/.steampipe/config/vanta.spc`).

## Examples

### Basic info
Explore which Vanta groups are available by identifying their names and IDs, and assess the elements within each group's checklist. This can be useful to understand the composition and configuration of these groups for better management and organization.

```sql+postgres
select
  name,
  id,
  checklist
from
  vanta_group;
```

```sql+sqlite
select
  name,
  id,
  checklist
from
  vanta_group;
```

### User details associated with each group
Discover the segments that detail the relationship between user information and their respective groups. This can be beneficial in managing user permissions and understanding the distribution of roles within your organization.

```sql+postgres
select
  g.name,
  u.display_name,
  u.email,
  u.permission_level
from
  vanta_group as g
  join vanta_user as u on g.id = u.role ->> 'id';
```

```sql+sqlite
select
  g.name,
  u.display_name,
  u.email,
  u.permission_level
from
  vanta_group as g
  join vanta_user as u on g.id = json_extract(u.role, '$.id');
```

### List all users in each group having administrator access
Determine the areas in which users have been granted administrative access within different groups. This can help in understanding the distribution of administrative privileges across your organization, aiding in access control and security management.

```sql+postgres
select
  g.name,
  u.display_name,
  u.email
from
  vanta_group as g
  join vanta_user as u on g.id = u.role ->> 'id' and u.permission_level = 'Admin';
```

```sql+sqlite
select
  g.name,
  u.display_name,
  u.email
from
  vanta_group as g
  join vanta_user as u on g.id = json_extract(u.role, '$.id') and u.permission_level = 'Admin';
```

### Get the count of users in each group
Explore which user groups have the most members to better manage resources and permissions. This can help in identifying areas for optimization and ensuring balanced distribution of users across different groups.

```sql+postgres
select
  g.name,
  count(u.display_name)
from
  vanta_group as g
  join vanta_user as u on g.id = u.role ->> 'id'
group by
  g.name;
```

```sql+sqlite
select
  g.name,
  count(u.display_name)
from
  vanta_group as g
  join vanta_user as u on g.id = json_extract(u.role, '$.id')
group by
  g.name;
```