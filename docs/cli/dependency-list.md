# hull dependency list

## Synopsis

`hull dependency list` prints a tabular view of every layer and required package declared in the package's `hull.yaml`, alongside the lock state recorded in `hull.lock`. For each entry the output includes the source URL, the constraint, the resolved version (from the lock), and a brief status flag indicating whether the layer is currently materialised in the local cache.

## When to use it

Use to audit a package's composition before installing — confirm you're picking up the layer versions you expect, spot any layer whose lock entry is missing or out-of-date, and verify the cache is populated before going offline. Read-only; never mutates `hull.yaml` or `hull.lock`.

## What happens when you run it

1. Reads `hull.yaml` from `<package-path>`.
2. Reads `hull.lock` if present.
3. Walks each `layers:` and `requires:` entry, joining declared metadata with locked metadata.
4. Prints one row per entry to stdout. No cluster contact, no network access, no file writes.

## Usage

```
hull dependency list <package-path> [flags]
```

## Flags

| Flag | Type | Default | Description |
|---|---|---|---|
| `-h, --help` | bool | false | help for list |

## Persistent flags inherited from `hull`

| Flag | Type | Description |
|---|---|---|
| `--debug` | bool | enable debug output |
| `--kube-context` | string | Kubernetes context to use |
| `--kubeconfig` | string | path to kubeconfig file |
| `-n, --namespace` | string | Kubernetes namespace |

## Examples

List declared layers and their lock state:

```sh
hull dependency list ./my-app
```

Run from inside a package directory:

```sh
cd ./my-app && hull dependency list .
```

Combine with `dependency tree` to see the composition chain:

```sh
hull dependency list ./my-app
hull dependency tree ./my-app
```

## See also

- [`dependency`](dependency.md)
- [`dependency tree`](dependency-tree.md) — visualise nested composition
- [`dependency update`](dependency-update.md) — refresh the lock
- [`dependency build`](dependency-build.md) — materialise the cache
- [Layers guide](../guides/layers.md)
