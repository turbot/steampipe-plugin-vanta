# Table: vanta_user

The `vanta_user` table can be used to query information about the currently active users.

## Examples

### Basic info

```sql
select
  display_name,
  uid,
  email,
  created_at
from
  vanta_user;
```
