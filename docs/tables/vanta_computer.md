---
title: "Steampipe Table: vanta_computer - Query Vanta Computer Assets using SQL"
description: "Allows users to query Computer Assets in Vanta, specifically retrieving details about each computer including its ID, name, operating system, and other related information."
---

# Table: vanta_computer - Query Vanta Computer Assets using SQL

Vanta is a security and compliance automation platform. It simplifies the complex process of achieving and maintaining compliance with standards like SOC 2, HIPAA, and GDPR. Vanta's Computer Assets are individual computing devices that are part of your organization's network.

## Table Usage Guide

The `vanta_computer` table provides insights into computer assets within Vanta's security and compliance automation platform. As a security analyst or compliance officer, explore details about each computer, including its operating system, installed software, and other related information through this table. Utilize it to monitor the security status of each computer, track software installations, and maintain compliance with various standards.

## Examples

### Basic info
Explore which computers have a specific owner, serial number, and operating system version. This can help in identifying and managing the different machines within your network.

```sql+postgres
select
  owner_name,
  serial_number,
  os_version,
  udid,
  last_check_date
from
  vanta_computer;
```

```sql+sqlite
select
  owner_name,
  serial_number,
  os_version,
  udid,
  last_check_date
from
  vanta_computer;
```

### List computers with unencrypted hard drive
Discover the segments that consist of computers with unencrypted hard drives, allowing you to identify potential security vulnerabilities and take necessary actions to ensure data protection.

```sql+postgres
select
  owner_name,
  serial_number,
  os_version,
  udid,
  disk_encryption
from
  vanta_computer
where
  not is_encrypted;
```

```sql+sqlite
select
  owner_name,
  serial_number,
  os_version,
  udid,
  disk_encryption
from
  vanta_computer
where
  is_encrypted = 0;
```

### List computers with no screen lock configured
Discover the segments that include computers lacking screen lock configuration. This could be useful for identifying potential security risks within your network and implementing necessary security measures.

```sql+postgres
select
  owner_name,
  serial_number,
  os_version,
  udid,
  screenlock
from
  vanta_computer
where
  not has_screen_lock;
```

```sql+sqlite
select
  owner_name,
  serial_number,
  os_version,
  udid,
  screenlock
from
  vanta_computer
where
  has_screen_lock = 0;
```

### List computers with no password manager installed
Identify computers that may be vulnerable due to the absence of a password manager. This can help in enhancing system security by pinpointing machines that need password manager installations.

```sql+postgres
select
  owner_name,
  serial_number,
  os_version,
  udid,
  password_manager
from
  vanta_computer
where
  not is_password_manager_installed;
```

```sql+sqlite
select
  owner_name,
  serial_number,
  os_version,
  udid,
  password_manager
from
  vanta_computer
where
  is_password_manager_installed = 0;
```

### List computers not checked recently
Explore which computers haven't been checked in the last 30 days. This is useful to identify potential risks or issues that might have been overlooked due to lack of regular monitoring.

```sql+postgres
select
  owner_name,
  serial_number,
  os_version,
  udid,
  last_check_date
from
  vanta_computer
where
  last_check_date < (current_timestamp - interval '30 days');
```

```sql+sqlite
select
  owner_name,
  serial_number,
  os_version,
  udid,
  last_check_date
from
  vanta_computer
where
  last_check_date < datetime('now', '-30 day');
```

### List computers by operating system type
Discover the distribution of computers by operating system type. This can help in understanding your infrastructure composition and planning for OS-specific security measures.

```sql+postgres
select
  operating_system ->> 'type' as os_type,
  operating_system ->> 'version' as os_version,
  count(*) as computer_count
from
  vanta_computer
where
  operating_system is not null
group by
  operating_system ->> 'type',
  operating_system ->> 'version'
order by
  computer_count desc;
```

```sql+sqlite
select
  json_extract(operating_system, '$.type') as os_type,
  json_extract(operating_system, '$.version') as os_version,
  count(*) as computer_count
from
  vanta_computer
where
  operating_system is not null
group by
  json_extract(operating_system, '$.type'),
  json_extract(operating_system, '$.version')
order by
  computer_count desc;
```

### List computers with antivirus issues
Identify computers that have issues with antivirus installation or configuration.

```sql+postgres
select
  owner_name,
  serial_number,
  os_version,
  antivirus_installation ->> 'outcome' as antivirus_status,
  antivirus_installation ->> 'lastCheckDate' as last_antivirus_check
from
  vanta_computer
where
  antivirus_installation ->> 'outcome' != 'PASS';
```

```sql+sqlite
select
  owner_name,
  serial_number,
  os_version,
  json_extract(antivirus_installation, '$.outcome') as antivirus_status,
  json_extract(antivirus_installation, '$.lastCheckDate') as last_antivirus_check
from
  vanta_computer
where
  json_extract(antivirus_installation, '$.outcome') != 'PASS';
```

### List computers owned by inactive users
Explore which computers are owned by users who are no longer active. This can help in asset management and ensuring resources are effectively allocated.

```sql+postgres
select
  u.display_name as owner,
  c.serial_number,
  c.os_version,
  u.employment_status,
  c.last_check_date
from
  vanta_computer as c
  join vanta_user as u on c.owner_id = u.id and not u.is_active;
```

```sql+sqlite
select
  u.display_name as owner,
  c.serial_number,
  c.os_version,
  u.employment_status,
  c.last_check_date
from
  vanta_computer as c
  join vanta_user as u on c.owner_id = u.id and u.is_active = 0;
```
