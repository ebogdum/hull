# hull plugin

## Synopsis

`hull plugin` manages the local plugin set. Plugins are external commands hull invokes when an unknown top-level command is given (`hull foo` becomes `hull-plugin-foo args...` if `foo` is an installed plugin). Subcommands install, list, remove, and update plugins.

## When to use it

Use to extend hull with site-specific or organisation-wide commands without modifying hull itself.

## Usage

```
hull plugin [command]
```

## Subcommands

- [`hull plugin list`](plugin-list.md) — List installed plugins
- [`hull plugin remove`](plugin-remove.md) — Remove an installed plugin
- [`hull plugin update`](plugin-update.md) — Update an installed plugin

## Flags

| Flag | Type | Default | Description |
|---|---|---|---|
| `-h, --help` | — | — | help for plugin |

## Persistent flags inherited from `hull`

| Flag | Type | Description |
|---|---|---|
| `--debug` | — | enable debug output |
| `--kube-context` | string | Kubernetes context to use |
| `--kubeconfig` | string | path to kubeconfig file |
| `-n, --namespace` | string | Kubernetes namespace |

## Examples

Install a plugin from a local archive:

```sh
hull plugin install ./my-plugin-1.0.tgz
```

List installed plugins:

```sh
hull plugin list
```

Update every installed plugin:

```sh
hull plugin update
```

## See also

- [`marketplace`](marketplace.md)
