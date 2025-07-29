---
title: "Steampipe Table: vanta_group - Query Vanta Groups using SQL"
description: "Allows users to query Groups in Vanta, specifically the group details including name, ID, and creation date."
---

# Table: vanta_group - Query Vanta Groups using SQL

Vanta is a security monitoring platform that simplifies the complex process of security compliance. It provides comprehensive visibility into an organization's security posture, helping to identify and mitigate potential vulnerabilities. Vanta's Groups feature allows for the management of user permissions, providing a structured way to assign and control access rights.

## Table Usage Guide

The `vanta_group` table provides insights into Groups within Vanta's security monitoring platform. As a Security or Compliance Officer, explore group-specific details through this table, including group names, IDs, and creation dates. This table provides basic group information available through the Vanta REST API.

## Examples

### Basic info
Explore which Vanta groups are available by identifying their names, IDs, and creation dates.

```sql+postgres
select
  name,
  id,
  creation_date
from
  vanta_group;
```

```sql+sqlite
select
  name,
  id,
  creation_date
from
  vanta_group;
```

### List groups created in the last 30 days
Discover recently created groups to track organizational changes and new team formations.

```sql+postgres
select
  name,
  id,
  creation_date
from
  vanta_group
where
  creation_date > (current_timestamp - interval '30 days')
order by
  creation_date desc;
```

```sql+sqlite
select
  name,
  id,
  creation_date
from
  vanta_group
where
  creation_date > datetime('now', '-30 days')
order by
  creation_date desc;
```

### List groups by creation date
Analyze the timeline of group creation to understand organizational growth and structure evolution.

```sql+postgres
select
  name,
  id,
  creation_date,
  extract(year from creation_date) as created_year,
  extract(month from creation_date) as created_month
from
  vanta_group
order by
  creation_date desc;
```

```sql+sqlite
select
  name,
  id,
  creation_date,
  strftime('%Y', creation_date) as created_year,
  strftime('%m', creation_date) as created_month
from
  vanta_group
order by
  creation_date desc;
```

### Count groups by creation year
Get insights into organizational expansion by analyzing group creation patterns over time.

```sql+postgres
select
  extract(year from creation_date) as creation_year,
  count(*) as groups_created
from
  vanta_group
where
  creation_date is not null
group by
  extract(year from creation_date)
order by
  creation_year desc;
```

```sql+sqlite
select
  strftime('%Y', creation_date) as creation_year,
  count(*) as groups_created
from
  vanta_group
where
  creation_date is not null
group by
  strftime('%Y', creation_date)
order by
  creation_year desc;
```
