# Table: vanta_monitor

Vanta helps businesses get and stay compliant by continuously monitoring people, systems and tools to improve security posture.

The table `vanta_monitor` provides information about all the monitors and their status.

**NOTE:**

- To query the table; **you must set** `session_id` argument in the config file (`~/.steampipe/config/vanta.spc`).

## Examples

### Basic info

```sql
select
  name,
  category,
  outcome,
  latest_flip_time,
  remediation_status ->> 'status' as status
from
  vanta_monitor;
```

### List all failed tests

```sql
select
  name,
  category,
  outcome,
  compliance_status,
  latest_flip_time
from
  vanta_monitor
where
  outcome = 'FAIL';
```

### Filter a specific test result by test ID

```sql
select
  name,
  category,
  outcome,
  compliance_status,
  latest_flip_time
from
  vanta_monitor
where
  test_id = 'test-03neqwol876pxg1iqjqib9';
```

### Count tests by remediation status

```sql
select
  remediation_status ->> 'status' as status,
  count(name)
from
  vanta_monitor
group by
  remediation_status ->> 'status';
```

### List failed tests by owner

```sql
select
  name,
  category,
  a ->> 'displayName' as owner,
  remediation_status ->> 'status' as status
from
  vanta_monitor
  left join jsonb_array_elements(assignees) as a on true
where
  outcome = 'FAIL';
```

### List failed tests by standard

```sql
select
  name,
  category,
  s -> 'standardInfo' ->> 'standard' as standard,
  remediation_status ->> 'status' as status
from
  vanta_monitor,
  jsonb_array_elements(controls) as c,
  jsonb_array_elements(c -> 'standardSections') as s
where
  outcome = 'FAIL';
```

### List failed tests by integration

```sql
select
  m.name,
  m.category,
  m.outcome,
  i.display_name as integration
from
  vanta_integration as i,
  jsonb_array_elements(i.tests) as t
  join vanta_monitor as m on m.test_id = t ->> 'testId' and m.outcome = 'FAIL'
```

### Count tests by outcome

```sql
select
  outcome,
  count(name)
from
  vanta_monitor
group by
  outcome;
```

### Count active tests by category

```sql
select
  category,
  count(name)
from
  vanta_monitor
where
  disabled_status is null
group by
  category;
```
