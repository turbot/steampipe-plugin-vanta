---
title: "Steampipe Table: vanta_evidence - Query Vanta Evidence using SQL"
description: "Allows users to query Vanta Evidence, specifically the metadata and details of the evidence collected by Vanta for security and compliance monitoring during audits."
---

# Table: vanta_evidence - Query Vanta Evidence using SQL

Vanta is a security and compliance platform that automates the collection of evidence for various security standards and regulations. The `vanta_evidence` table provides access to audit evidence data, allowing you to explore evidence collected during specific audits.

## Table Usage Guide

The `vanta_evidence` table offers insights into evidence collected during Vanta audits. As a Security Analyst, you can use this table to explore specific details about each piece of evidence, including its status, type, associated controls, and metadata. By querying this table, you can effectively track audit progress and verify evidence collection status.

**Important Notes**

- You must provide the `audit_id` in the query parameter in order to query this table.
- The access token must have the scope `auditor-api.audit:read`.

## Examples

### Basic info

Explore the various types of evidence collected during an audit, including their status and creation details.

```sql+postgres
select
  id,
  name,
  evidence_type,
  status,
  creation_date,
  description
from
  vanta_evidence
where
  audit_id = 'your_audit_id';
```

```sql+sqlite
select
  id,
  name,
  evidence_type,
  status,
  creation_date,
  description
from
  vanta_evidence
where
  audit_id = 'your_audit_id';
```

### List evidence by status

Get all evidence that is ready for audit review.

```sql+postgres
select
  id,
  name,
  evidence_type,
  status,
  status_updated_date,
  related_control_names
from
  vanta_evidence
where
  audit_id = 'your_audit_id'
  and status = 'Ready for audit';
```

```sql+sqlite
select
  id,
  name,
  evidence_type,
  status,
  status_updated_date,
  related_control_names
from
  vanta_evidence
where
  audit_id = 'your_audit_id'
  and status = 'Ready for audit';
```

### List evidence by type

Get all test-type evidence and their test results.

```sql+postgres
select
  id,
  name,
  evidence_type,
  test_status,
  status,
  creation_date
from
  vanta_evidence
where
  audit_id = 'your_audit_id'
  and evidence_type = 'Test';
```

```sql+sqlite
select
  id,
  name,
  evidence_type,
  test_status,
  status,
  creation_date
from
  vanta_evidence
where
  audit_id = 'your_audit_id'
  and evidence_type = 'Test';
```

### List flagged evidence

Identify evidence that has been flagged for review.

```sql+postgres
select
  id,
  name,
  evidence_type,
  status,
  status_updated_date,
  description,
  related_control_names
from
  vanta_evidence
where
  audit_id = 'your_audit_id'
  and status = 'Flagged';
```

```sql+sqlite
select
  id,
  name,
  evidence_type,
  status,
  status_updated_date,
  description,
  related_control_names
from
  vanta_evidence
where
  audit_id = 'your_audit_id'
  and status = 'Flagged';
```

### Count evidence by status

Get a summary of evidence status distribution for an audit.

```sql+postgres
select
  status,
  count(*) as evidence_count
from
  vanta_evidence
where
  audit_id = 'your_audit_id'
group by
  status
order by
  evidence_count desc;
```

```sql+sqlite
select
  status,
  count(*) as evidence_count
from
  vanta_evidence
where
  audit_id = 'your_audit_id'
group by
  status
order by
  evidence_count desc;
```

### List recently created evidence

Get evidence created in the last 30 days.

```sql+postgres
select
  id,
  name,
  evidence_type,
  status,
  creation_date,
  related_control_names
from
  vanta_evidence
where
  audit_id = 'your_audit_id'
  and creation_date >= now() - interval '30 days'
order by
  creation_date desc;
```

```sql+sqlite
select
  id,
  name,
  evidence_type,
  status,
  creation_date,
  related_control_names
from
  vanta_evidence
where
  audit_id = 'your_audit_id'
  and creation_date >= datetime('now', '-30 days')
order by
  creation_date desc;
```

### Get specific evidence details

Retrieve detailed information about a specific piece of evidence.

```sql+postgres
select
  id,
  external_id,
  name,
  evidence_type,
  status,
  test_status,
  creation_date,
  status_updated_date,
  deletion_date,
  description,
  related_controls
from
  vanta_evidence
where
  audit_id = 'your_audit_id'
  and id = 'evidence_id';
```

```sql+sqlite
select
  id,
  external_id,
  name,
  evidence_type,
  status,
  test_status,
  creation_date,
  status_updated_date,
  deletion_date,
  description,
  related_controls
from
  vanta_evidence
where
  audit_id = 'your_audit_id'
  and id = 'evidence_id';
```
