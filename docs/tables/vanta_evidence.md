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

### List requests with restricted document access

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

### List non-relevant requests

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

### Get the count of request by document category

```sql
select
  category,
  count(title)
from
  vanta_evidence
group by
  category;
```
