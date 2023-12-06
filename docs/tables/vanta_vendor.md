---
title: "Steampipe Table: vanta_vendor - Query Vanta Vendor using SQL"
description: "Allows users to query Vanta Vendors, providing detailed information about the vendors' profiles, including their security and compliance status."
---

# Table: vanta_vendor - Query Vanta Vendor using SQL

Vanta Vendor is a resource within the Vanta service that provides detailed information about the vendors used by an organization. It includes details about the vendors' profiles, such as their security and compliance status. This information is crucial for organizations to understand and manage the security risks associated with their vendors.

## Table Usage Guide

The `vanta_vendor` table provides insights into the vendors used by an organization within the Vanta service. As a security analyst, you can explore vendor-specific details through this table, including their security and compliance status. Utilize it to uncover information about your vendors, such as their security scores, the number of employees, and the services they provide.

**Important Notes**
- To query the table you must set `session_id` argument in the config file (`~/.steampipe/config/vanta.spc`).

## Examples

### Basic info
Explore which vendors have been enlisted by Vanta, along with their respective severity levels and corresponding URLs. This can be useful for assessing the risk profile associated with each vendor and managing them efficiently.

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
Uncover the details of vendors categorized as high severity. This information can be useful for prioritizing vendor management tasks and focusing on potential risk areas.

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
Discover the vendors whose security checks are overdue by a year. This query is useful to maintain security standards and ensure all vendors are regularly reviewed.

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
Identify vendors who have not submitted any assessment documents. This query can be useful for compliance checks and to ensure all vendors are meeting documentation requirements.

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
Discover the segments that have specific vendor ownership. This query allows you to identify and understand the relationship between vendors and their owners, which is essential for managing vendor relationships and permissions effectively.

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