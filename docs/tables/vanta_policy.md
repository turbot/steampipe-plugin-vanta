# Table: vanta_policy

Policy contains a set of rules related to information security for your organization. Policies touch on all business areas, so the creation process requires cross-team collaboration.

**NOTE:**

- To query the table; **you must set** `api_token` argument in the config file (`~/.steampipe/config/vanta.spc`).

## Examples

### Basic info

```sql
select
  display_name,
  policy_type,
  url,
  created_at
from
  vanta_policy;
```

### List unapproved policies

```sql
select
  display_name,
  policy_type,
  url,
  created_at
from
  vanta_policy
where
  approver is null;
```

### List expired policies

```sql
select
  display_name,
  policy_type,
  url,
  created_at
from
  vanta_policy
where
  (approved_at + interval '1 year') < current_timestamp;
```

### List policies expires within 30 days

```sql
select
  display_name,
  url,
  approved_at,
  'expires in ' || extract(day from ((approved_at + interval '1 year') - current_timestamp)) || ' day(s)' as status
from
  vanta_policy
where
  current_timestamp < (approved_at + interval '1 year')
  and extract(day from ((approved_at + interval '1 year') - current_timestamp)) <= '30';
```
