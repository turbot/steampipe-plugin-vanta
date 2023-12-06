---
title: "Steampipe Table: vanta_computer - Query Vanta Computer Assets using SQL"
description: "Allows users to query Computer Assets in Vanta, specifically retrieving details about each computer including its ID, name, operating system, and other related information."
---

# Table: vanta_computer - Query Vanta Computer Assets using SQL

Vanta is a security and compliance automation platform. It simplifies the complex process of achieving and maintaining compliance with standards like SOC 2, HIPAA, and GDPR. Vanta's Computer Assets are individual computing devices that are part of your organization's network.

## Table Usage Guide

The `vanta_computer` table provides insights into computer assets within Vanta's security and compliance automation platform. As a security analyst or compliance officer, explore details about each computer, including its operating system, installed software, and other related information through this table. Utilize it to monitor the security status of each computer, track software installations, and maintain compliance with various standards.

**Important Notes**
- To query the table you must set `session_id` argument in the config file (`~/.steampipe/config/vanta.spc`).

## Examples

### Basic info
Explore which computers have a specific owner, serial number, and operating system version. This can help in identifying and managing the different machines within your network.

```sql
select
  owner_name,
  serial_number,
  agent_version,
  os_version,
  hostname
from
  vanta_computer;
```

### List computers with unencrypted hard drive
Discover the segments that consist of computers with unencrypted hard drives, allowing you to identify potential security vulnerabilities and take necessary actions to ensure data protection.

```sql
select
  owner_name,
  serial_number,
  agent_version,
  os_version,
  hostname
from
  vanta_computer
where
  not is_encrypted;
```

### List computers with no screen lock configured
Discover the segments that include computers lacking screen lock configuration. This could be useful for identifying potential security risks within your network and implementing necessary security measures.

```sql
select
  owner_name,
  serial_number,
  agent_version,
  os_version,
  hostname
from
  vanta_computer
where
  not has_screen_lock;
```

### List computers with no password manager installed
Identify computers that may be vulnerable due to the absence of a password manager. This can help in enhancing system security by pinpointing machines that need password manager installations.

```sql
select
  owner_name,
  serial_number,
  agent_version,
  os_version,
  hostname
from
  vanta_computer
where
  installed_password_managers is null;
```

### List computers not checked over the last 90 days
Explore which computers haven't been checked in the last 90 days. This is useful to identify potential risks or issues that might have been overlooked due to lack of regular monitoring.

```sql
select
  owner_name,
  serial_number,
  agent_version,
  os_version,
  hostname,
  last_ping
from
  vanta_computer
where
  last_ping < (current_timestamp - interval '90 days');
```

### List unmonitored computers
Discover the segments that contain computers that are not being monitored. This is useful in identifying potential gaps in your IT infrastructure, allowing you to address any unsupported operating systems or versions.

```sql
select
  owner_name,
  serial_number,
  agent_version,
  os_version,
  hostname,
  case
    when (unsupported_reasons -> 'unsupportedOsVersion')::boolean then 'OS version not supported'
    when (unsupported_reasons -> 'unsupportedOsType')::boolean then 'OS not supported'
  end as status
from
  vanta_computer
where
  unsupported_reasons is not null;
```

### List computers with Tailscale app installed
Determine the areas in which computers have the Tailscale app installed. This query is useful for gaining insights into the distribution and usage of this specific application within your network.

```sql
select
  owner_name,
  serial_number,
  last_ping,
  app as application
from
  vanta_computer,
  jsonb_array_elements_text(endpoint_applications) as app
where
  app like 'Tailscale %';
```

### List computers with no Slack app installed
Determine the computers that do not have the Slack app installed. This can be useful for IT administrators to identify and rectify gaps in software deployment across the organization.

```sql
with device_with_slack_installed as (
  select
    distinct id
  from
    vanta_computer,
    jsonb_array_elements_text(endpoint_applications) as app
  where
    app like 'Slack %'
)
select
  owner_name,
  serial_number,
  last_ping
from
  vanta_computer
where
  endpoint_applications is not null
  and id not in (
    select
      id
    from
      device_with_slack_installed
  );
```

### List computers with an older version of Zoom app (< 5.12)
Determine the areas in which computers are running outdated versions of the Zoom app for potential software updates. This helps in maintaining system security and ensuring all devices are up-to-date with the latest software versions.

```sql
select
  owner_name,
  serial_number,
  last_ping,
  app as application
from
  vanta_computer,
  jsonb_array_elements_text(endpoint_applications) as app
where
  app like 'zoom.us %'
  and string_to_array(split_part(app, ' ', 2), '.')::int[] < string_to_array('5.12', '.')::int[];
```

### List computers owned by inactive users
Explore which computers are owned by users who are no longer active. This can help in asset management and ensuring resources are effectively allocated.

```sql
select
  u.display_name as owner,
  c.serial_number,
  u.end_date,
  c.last_ping
from
  vanta_computer as c
  join vanta_user as u on c.owner_id = u.id and not u.is_active;
```