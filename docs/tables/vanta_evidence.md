---
title: "Steampipe Table: vanta_evidence - Query Vanta Evidence using SQL"
description: "Allows users to query Vanta Evidence, specifically the metadata and details of the evidence collected by Vanta for security and compliance monitoring."
---

# Table: vanta_evidence - Query Vanta Evidence using SQL

Vanta is a security and compliance platform that automates the collection of evidence for various security standards and regulations. It provides a centralized way to monitor and manage security controls across your infrastructure, applications, and services. Vanta Evidence is a key component of this platform, capturing and storing the necessary data to demonstrate compliance with these standards.

## Table Usage Guide

The `vanta_evidence` table offers insights into the evidence collected by Vanta for security and compliance monitoring. As a Security Analyst, you can use this table to explore specific details about each piece of evidence, including its metadata, associated controls, and status. By querying this table, you can effectively track and verify your organization's compliance status and identify potential security issues.

## Examples

### Basic info
Explore the various categories of evidence requests in the Vanta system, identifying instances where access to certain information might be restricted. This can help in understanding the nature of information requests and ensuring compliance with access control policies.

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
Explore which evidences have restricted document access to ensure compliance and maintain the integrity of sensitive information. This can be beneficial in situations where access needs to be limited due to confidentiality or security reasons.

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
Uncover the details of dismissed evidences in your Vanta security compliance data. This query is particularly useful in identifying and reviewing non-relevant evidences that have been marked as dismissed.

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
Explore which pieces of evidence are due for renewal within the next month. This is useful for staying on top of compliance requirements and ensuring that all evidence is updated in a timely manner.

```sql
select
  title,
  category,
  renewal_metadata ->> 'nextDate' as update_by
from
  vanta_evidence
where
  renewal_metadata ->> 'nextDate' != ''
  and current_timestamp < (renewal_metadata ->> 'nextDate')::timestamp
  and extract (day from ((renewal_metadata ->> 'nextDate')::timestamp - current_timestamp)) < 30
  and dismissed_status is null;
```

### Get the count of evidence by category
Explore which categories have the most evidence. This can be useful in identifying areas that may require additional scrutiny or attention.

```sql
select
  category,
  count(title)
from
  vanta_evidence
group by
  category;
```