![image](https://hub.steampipe.io/images/plugins/turbot/vanta-social-graphic.png)

# Vanta Plugin for Steampipe

Use SQL to query users, policies, compliances, and more from your Vanta organization.

- **[Get started →](https://hub.steampipe.io/plugins/turbot/vanta)**
- Documentation: [Table definitions & examples](https://hub.steampipe.io/plugins/turbot/vanta/tables)
- Community: [Slack Channel](https://steampipe.io/community/join)
- Get involved: [Issues](https://github.com/turbot/steampipe-plugin-vanta/issues)

## Quick start

Install the plugin with [Steampipe](https://steampipe.io):

```shell
steampipe plugin install vanta
```

Configure your [credentials](https://hub.steampipe.io/plugins/turbot/vanta#credentials) and [config file](https://hub.steampipe.io/plugins/turbot/vanta#configuration).

Run a query:

```sql
select
  display_name,
  id,
  email,
  is_active
from
  vanta_user
where
  is_active;
```

## Developing

Prerequisites:

- [Steampipe](https://steampipe.io/downloads)
- [Golang](https://golang.org/doc/install)

Clone:

```sh
git clone https://github.com/turbot/steampipe-plugin-vanta.git
cd steampipe-plugin-vanta
```

Build, which automatically installs the new version to your `~/.steampipe/plugins` directory:

```sh
make
```

Configure the plugin:

```sh
cp config/* ~/.steampipe/config
vi ~/.steampipe/config/vanta.spc
```

Try it!

```shell
steampipe query
> .inspect vanta
```

Further reading:

- [Writing plugins](https://steampipe.io/docs/develop/writing-plugins)
- [Writing your first table](https://steampipe.io/docs/develop/writing-your-first-table)

## Contributing

Please see the [contribution guidelines](https://github.com/turbot/steampipe/blob/main/CONTRIBUTING.md) and our [code of conduct](https://github.com/turbot/steampipe/blob/main/CODE_OF_CONDUCT.md). All contributions are subject to the [Apache 2.0 open source license](https://github.com/turbot/steampipe-plugin-vanta/blob/main/LICENSE).

`help wanted` issues:

- [Steampipe](https://github.com/turbot/steampipe/labels/help%20wanted)
- [Vanta Plugin](https://github.com/turbot/steampipe-plugin-vanta/labels/help%20wanted)
