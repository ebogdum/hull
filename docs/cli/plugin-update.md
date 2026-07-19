# hull plugin update

## Synopsis

`hull plugin update` (aliases `up`, `upgrade`) refreshes an installed plugin. For
a plugin installed from git, it pulls the latest commit; then it re-reads the
plugin's metadata and runs its update hook.

## When to use it

Run it to pick up a newer version of a plugin — a bug fix, a new subcommand, a
changed flag. After updating, run `hull <name> --help` to see what changed.

## What happens

1. hull finds the plugin's directory under `~/.config/hull/plugins/`, matching
   `<name>` against either the directory name or the `name` in its
   `plugin.yaml`.
2. If that directory is a git checkout, hull runs `git pull --ff-only` there, so
   the update only applies when it fast-forwards cleanly. A plugin installed
   from a local directory has no git checkout, so this step is skipped.
3. hull re-reads `plugin.yaml` and re-checks the plugin's command.
4. hull runs the plugin's update hook if it declares one.
5. hull prints the plugin's name and version.

## Usage

```
hull plugin update <name> [flags]
```

## Flags

Inherits the global flags.

## Worked example

Update a git-installed plugin:

```sh
hull plugin update backup
```

```
Updated plugin: backup v0.5.0
```

You can use an alias:

```sh
hull plugin up backup
```

## See also

- [`plugin install`](plugin-install.md)
- [`plugin list`](plugin-list.md)
- [`plugin remove`](plugin-remove.md)
