# hull plugin

## Synopsis

`hull plugin` manages the extra commands you bolt onto hull. A plugin is a small
program plus a `plugin.yaml` that names it; once installed, you run it as if it
were built in — `hull <plugin-name> [args]`.

Plugins live in `~/.config/hull/plugins/`, one directory per plugin. Everything
here is local to your machine and your user account — managing plugins never
touches a cluster.

## Subcommands

| Command | What it does |
|---|---|
| [`install`](plugin-install.md) | Install a plugin from a git URL or a local directory |
| [`list`](plugin-list.md) | Show the plugins you have installed |
| [`update`](plugin-update.md) | Pull the newest version of an installed plugin |
| [`remove`](plugin-remove.md) | Uninstall a plugin |

## Usage

```
hull plugin [command]
```

Once a plugin is installed you invoke it directly — hull treats an unknown
command as a plugin name:

```sh
hull plugin install https://github.com/acme/hull-backup.git
hull backup --release web        # "backup" now runs as a hull command
```

## See also

- [`marketplace`](marketplace.md) — find and verify signed plugins to install
- [`plugin install`](plugin-install.md)
- [`plugin list`](plugin-list.md)
