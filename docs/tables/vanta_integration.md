# Table: vanta_integration

An integration is a connection which can be integrated with Vanta to activate automated evidence collection and monitoring.

**Note:** There are various integrations available that can be integrated. The table `vanta_integration` only returns the integrations that are connected.

## Examples

### Basic info

```sql
select
  display_name,
  id,
  description,
  logo_slug_id
from
  vanta_integration;
```

### List integrations having disabled credentials

```sql
select
  display_name,
  id,
  (c ->> 'metadata')::jsonb as credential_metadata,
  c ->> 'service'
from
  vanta_integration,
  jsonb_array_elements(credentials) as c
where
  c ->> 'isDisabled' = 'true';
```

### List integrations with failed tests

```sql
select
  display_name,
  id,
  t ->> 'testId' as test_id,
  m.outcome as test_status
from
  vanta_integration as i,
  jsonb_array_elements(tests) as t
  join vanta_monitor as m on m.test_id = t ->> 'testId' and outcome = 'FAIL';
```
