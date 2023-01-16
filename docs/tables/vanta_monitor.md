# Table: vanta_monitor

Vanta helps businesses get and stay compliant by continuously monitoring people, systems and tools to improve security posture.

The table `vanta_monitor` provides information about all the monitors and their status.

## Examples

### Basic info

```sql
select
  name,
  category,
  outcome,
  latest_flip,
  timestamp
from
  vanta_monitor;
```

### List all failed tests

```sql
select
  name,
  category,
  outcome,
  latest_flip,
  timestamp
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
  latest_flip,
  timestamp,
  fail_message,
  remediation
from
  vanta_monitor
where
  test_id = 'test-03neqwol876pxg1iqjqib9';
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
