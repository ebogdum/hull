# hull plugin remove

## Synopsis

`hull plugin remove` (aliases `rm`, `uninstall`) deletes an installed plugin.
Its directory under `~/.config/hull/plugins/` is removed, and `hull <name>` no
longer runs it.

## When to use it

Use it to retire a plugin you no longer need, clean up after developing one, or
clear a broken install before reinstalling.

## What happens

1. hull finds the plugin's directory under `~/.config/hull/plugins/`, matching
   `<name>` against either the directory name or the `name` in its
   `plugin.yaml`.
2. hull runs the plugin's delete hook if it declares one.
3. hull removes the directory.
4. hull prints a confirmation line.

## Usage

```
hull plugin remove <name> [flags]
```

## Flags

Inherits the global flags.

## Worked example

Remove a plugin by name:

```sh
hull plugin remove hello
```

```
Removed plugin: hello
```

You can use an alias:

```sh
hull plugin rm hello
```

## See also

- [`plugin install`](plugin-install.md)
- [`plugin list`](plugin-list.md)
- [`plugin update`](plugin-update.md)
