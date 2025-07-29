---
title: "Steampipe Table: vanta_vendor - Query Vanta Vendor using SQL"
description: "Allows users to query Vanta Vendors, providing detailed information about the vendors' profiles, including their security and compliance status, risk levels, contracts, and review schedules."
---

# Table: vanta_vendor - Query Vanta Vendor using SQL

Vanta Vendor is a resource within the Vanta service that provides detailed information about the vendors used by an organization. It includes details about the vendors' profiles, such as their security and compliance status, risk assessments, contract information, and security review schedules. This information is crucial for organizations to understand and manage the security risks associated with their vendors.

## Table Usage Guide

The `vanta_vendor` table provides insights into the vendors used by an organization within the Vanta service. As a security analyst or vendor risk manager, you can explore vendor-specific details through this table, including their security and compliance status, risk levels, contract details, and review schedules. Utilize it to uncover information about your vendors, such as their risk scores, contract amounts, and security review status.

## Examples

### Basic info
Explore which vendors have been enlisted by Vanta, along with their respective risk levels and corresponding URLs.

```sql+postgres
select
  name,
  id,
  inherent_risk_level,
  residual_risk_level,
  status,
  website_url
from
  vanta_vendor;
```

```sql+sqlite
select
  name,
  id,
  inherent_risk_level,
  residual_risk_level,
  status,
  website_url
from
  vanta_vendor;
```

### List vendors with high inherent risk
Uncover the details of vendors categorized as high risk to prioritize vendor management tasks and focus on potential risk areas.

```sql+postgres
select
  name,
  id,
  inherent_risk_level,
  residual_risk_level,
  contract_amount,
  next_security_review_due_date
from
  vanta_vendor
where
  inherent_risk_level = 'HIGH';
```

```sql+sqlite
select
  name,
  id,
  inherent_risk_level,
  residual_risk_level,
  contract_amount,
  next_security_review_due_date
from
  vanta_vendor
where
  inherent_risk_level = 'HIGH';
```

### List vendors with security reviews overdue
Discover vendors whose security reviews are overdue to maintain security standards and ensure all vendors are regularly reviewed.

```sql+postgres
select
  name,
  id,
  inherent_risk_level,
  next_security_review_due_date,
  last_security_review_completion_date,
  extract(day from (current_timestamp - next_security_review_due_date)) as days_overdue
from
  vanta_vendor
where
  next_security_review_due_date < current_timestamp;
```

```sql+sqlite
select
  name,
  id,
  inherent_risk_level,
  next_security_review_due_date,
  last_security_review_completion_date,
  cast(julianday('now') - julianday(next_security_review_due_date) as integer) as days_overdue
from
  vanta_vendor
where
  next_security_review_due_date < datetime('now');
```

### List vendors by contract value
Explore vendors ordered by their contract value to understand the financial impact of your vendor relationships.

```sql+postgres
select
  name,
  id,
  contract_amount,
  inherent_risk_level,
  contract_start_date,
  contract_renewal_date
from
  vanta_vendor
where
  contract_amount is not null
order by
  contract_amount desc;
```

```sql+sqlite
select
  name,
  id,
  contract_amount,
  inherent_risk_level,
  contract_start_date,
  contract_renewal_date
from
  vanta_vendor
where
  contract_amount is not null
order by
  contract_amount desc;
```

### List vendors with expiring contracts
Identify vendors whose contracts are expiring within the next 90 days to plan for contract renewals.

```sql+postgres
select
  name,
  id,
  contract_renewal_date,
  contract_amount,
  account_manager_name,
  account_manager_email,
  extract(day from (contract_renewal_date - current_timestamp)) as days_until_renewal
from
  vanta_vendor
where
  contract_renewal_date between current_timestamp and (current_timestamp + interval '90 days');
```

```sql+sqlite
select
  name,
  id,
  contract_renewal_date,
  contract_amount,
  account_manager_name,
  account_manager_email,
  cast(julianday(contract_renewal_date) - julianday('now') as integer) as days_until_renewal
from
  vanta_vendor
where
  contract_renewal_date between datetime('now') and datetime('now', '+90 days');
```

### Count vendors by risk level
Analyze the distribution of vendors across different risk levels to understand your overall vendor risk profile.

```sql+postgres
select
  inherent_risk_level,
  count(*) as vendor_count,
  avg(contract_amount) as avg_contract_value
from
  vanta_vendor
where
  inherent_risk_level is not null
group by
  inherent_risk_level
order by
  vendor_count desc;
```

```sql+sqlite
select
  inherent_risk_level,
  count(*) as vendor_count,
  avg(contract_amount) as avg_contract_value
from
  vanta_vendor
where
  inherent_risk_level is not null
group by
  inherent_risk_level
order by
  vendor_count desc;
```

### List vendors by category
Explore vendors grouped by their categories to understand service distribution across your vendor portfolio.

```sql+postgres
select
  category_display_name,
  count(*) as vendor_count,
  sum(contract_amount) as total_contract_value,
  avg(contract_amount) as avg_contract_value
from
  vanta_vendor
where
  category_display_name is not null
group by
  category_display_name
order by
  vendor_count desc;
```

```sql+sqlite
select
  category_display_name,
  count(*) as vendor_count,
  sum(contract_amount) as total_contract_value,
  avg(contract_amount) as avg_contract_value
from
  vanta_vendor
where
  category_display_name is not null
group by
  category_display_name
order by
  vendor_count desc;
```

### List vendors visible to auditors
Identify which vendors are configured to be visible to auditors during compliance audits.

```sql+postgres
select
  name,
  id,
  inherent_risk_level,
  is_visible_to_auditors,
  last_security_review_completion_date,
  contract_amount
from
  vanta_vendor
where
  is_visible_to_auditors = true
order by
  inherent_risk_level, contract_amount desc;
```

```sql+sqlite
select
  name,
  id,
  inherent_risk_level,
  is_visible_to_auditors,
  last_security_review_completion_date,
  contract_amount
from
  vanta_vendor
where
  is_visible_to_auditors = 1
order by
  inherent_risk_level, contract_amount desc;
```

### List vendors with auto-scored risk
Explore vendors that have automated risk scoring enabled to understand how risk assessment is being automated.

```sql+postgres
select
  name,
  id,
  inherent_risk_level,
  residual_risk_level,
  is_risk_auto_scored,
  contract_amount
from
  vanta_vendor
where
  is_risk_auto_scored = true
order by
  inherent_risk_level, name;
```

```sql+sqlite
select
  name,
  id,
  inherent_risk_level,
  residual_risk_level,
  is_risk_auto_scored,
  contract_amount
from
  vanta_vendor
where
  is_risk_auto_scored = 1
order by
  inherent_risk_level, name;
```

### List vendors by headquarters location
Analyze vendor distribution by their headquarters location to understand geographic risk concentration.

```sql+postgres
select
  vendor_headquarters,
  count(*) as vendor_count,
  sum(contract_amount) as total_contract_value
from
  vanta_vendor
where
  vendor_headquarters is not null
group by
  vendor_headquarters
order by
  vendor_count desc;
```

```sql+sqlite
select
  vendor_headquarters,
  count(*) as vendor_count,
  sum(contract_amount) as total_contract_value
from
  vanta_vendor
where
  vendor_headquarters is not null
group by
  vendor_headquarters
order by
  vendor_count desc;
```
