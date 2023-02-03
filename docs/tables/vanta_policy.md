# Table: vanta_policy

Policy contains a set of rules related to information security for your organization. Policies touch on all business areas, so the creation process requires cross-team collaboration.

**NOTE:**

- To query the table; **you must set** `session_id` argument in the config file (`~/.steampipe/config/vanta.spc`).

## Examples

### Basic info

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
