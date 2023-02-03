# Table: vanta_vendor

The `vanta_vendor` table can be used to query information about the vendors, and it's own security documentation to ensure compliance. Using Vanta, admin can conduct an assessment of vendors essential to your business' services and then take action to assign ownership and level of risk.

**NOTE:**

- To query the table; **you must set** `session_id` argument in the config file (`~/.steampipe/config/vanta.spc`).

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

### Get the owner information of each vendor

```sql
select
  v.name as vendor_name,
  v.severity as vendor_severity,
  u.display_name as owner_name,
  u.email as owner_email,
  u.permission_level
from
  vanta_vendor as v
  join vanta_user as u on v.owner ->> 'id' = u.id;
```
