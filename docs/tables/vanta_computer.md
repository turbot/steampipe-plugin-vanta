# Table: vanta_computer

The `vanta_computer` table can be used to query information about all computers within your organization to ensure that security-relevant settings are configured promptly.

**NOTE:**

- To query the table; **you must set** `session_id` argument in the config file (`~/.steampipe/config/vanta.spc`).

## Examples

### Basic info

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
