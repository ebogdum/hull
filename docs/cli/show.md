# hull show

## Synopsis

`hull show` prints information about a package without installing it: the `hull.yaml` (`chart`), the `values.yaml` (`values`), the README (`readme`), the CRDs in `crds/` (`crds`), or all of the above (`all`).

## When to use it

Use to inspect a package's structure before installing or before vendoring it. Works against local directories, registered repos, and OCI references.

## Usage

```
hull show [command]
```

## Subcommands

- [`hull show chart`](show-chart.md) — Show package metadata (hull.yaml)
- [`hull show crds`](show-crds.md) — Show CRDs declared by the package
- [`hull show readme`](show-readme.md) — Show package README
- [`hull show values`](show-values.md) — Show default values (values.yaml)

## Flags

| Flag | Type | Default | Description |
|---|---|---|---|
| `-h, --help` | — | — | help for show |

## Persistent flags inherited from `hull`

| Flag | Type | Description |
|---|---|---|
| `--debug` | — | enable debug output |
| `--kube-context` | string | Kubernetes context to use |
| `--kubeconfig` | string | path to kubeconfig file |
| `-n, --namespace` | string | Kubernetes namespace |

## Examples

Show the package metadata:

```sh
hull show chart my-charts/my-app
```

Show the default values:

```sh
hull show values my-charts/my-app
```

Show all package data:

```sh
hull show all my-charts/my-app -o yaml
```

## See also

- [`pull`](pull.md)
- [`search`](search.md)
