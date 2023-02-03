# Table: vanta_evidence

The evidence request provides a list of documents that need to provide as a part of the audit for the chosen certificate, i.e., SOC2, ISO 27001, or HIPAA. Each request is a piece of evidence that is required to complete the audit.

**NOTE:**

- To query the table; **you must set** `api_token` argument in the config file (`~/.steampipe/config/vanta.spc`).

## Examples

### Basic info

```sql
select
  title,
  evidence_request_id,
  category,
  description,
  restricted
from
  vanta_evidence;
```

### List evidences with restricted document access

```sql
select
  title,
  evidence_request_id,
  category,
  description
from
  vanta_evidence
where
  restricted;
```

### List non-relevant evidences

```sql
select
  title,
  evidence_request_id,
  category,
  dismissed_status
from
  vanta_evidence
where
  dismissed_status -> 'isDismissed' = 'true';
```

### List evidences up for renewal within 30 days

```sql
select
  title,
  category,
  renewal_metadata ->> 'nextDate' as update_by
from
  vanta_evidence
where
  current_timestamp < (renewal_metadata ->> 'nextDate')::timestamp
  and extract (day from ((renewal_metadata ->> 'nextDate')::timestamp - current_timestamp)) < 30
  and dismissed_status is null;
```

### Get the count of evidence by category

```sql
select
  category,
  count(title)
from
  vanta_evidence
group by
  category;
```
