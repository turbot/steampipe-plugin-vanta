# Table: vanta_group

A group is a collection of team members that helps in organizing them together based on criteria they have in common, such as employment status or department.

**NOTE:**

- To query the table; **you must set** `session_id` argument in the config file (`~/.steampipe/config/vanta.spc`).

## Examples

### Basic info

```sql
select
  name,
  id,
  checklist
from
  vanta_group;
```

### User details associated with each group

```sql
select
  g.name,
  u.display_name,
  u.email,
  u.permission_level
from
  vanta_group as g
  join vanta_user as u on g.id = u.role ->> 'id';
```

### List all users in each group having administrator access

```sql
select
  g.name,
  u.display_name,
  u.email
from
  vanta_group as g
  join vanta_user as u on g.id = u.role ->> 'id' and u.permission_level = 'Admin';
```

### Get the count of users in each group

```sql
select
  g.name,
  count(u.display_name)
from
  vanta_group as g
  join vanta_user as u on g.id = u.role ->> 'id'
group by
  g.name;
```
