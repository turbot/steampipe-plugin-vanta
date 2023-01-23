# Table: vanta_vendor

The `vanta_vendor` table can be used to query information about the vendors, and it's own security documentation to ensure compliance. Using Vanta, admin can conduct an assessment of vendors essential to your business’ services and then take action to assign ownership and level of risk

## Examples

### Basic info

```sql
select
  name,
  id,
  severity,
  url
from
  vanta_vendor;
```

### List vendors with high severity

```sql
select
  name,
  id,
  severity,
  url
from
  vanta_vendor
where
  severity = 'high';
```

### List vendors with security checks overdue

```sql
select
  name,
  id,
  severity,
  url
from
  vanta_vendor
where
  current_timestamp > (latest_security_review_completed_at + interval '1 year');
```

### List vendors with no documents provided

```sql
select
  name,
  id,
  severity,
  url
from
  vanta_vendor
where
  assessment_documents is null;
```
