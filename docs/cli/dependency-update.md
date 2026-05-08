# hull dependency update

## Synopsis

`hull dependency update` re-resolves every layer and required package against its source (HTTP repo index, OCI registry tag list, git ref, local path) and rewrites `hull.lock` with the freshest pinned digests. The lockfile is the source of truth for `hull install`, `hull template`, and `hull dependency build` — they all consult `hull.lock` first and only fall back to `hull.yaml`'s constraint if the lock is missing or stale.

## When to use it

Run whenever you edit `hull.yaml`'s `layers:` or `requires:`, or when you want to bump a layer to the highest version still satisfying its constraint. Pass an optional `[name]` argument to update a single layer in place; without it, every layer is re-resolved. Always commit the resulting `hull.lock` to source control — without it, two builds of the same package can pick up different layer versions if the constraint allows it.

## What happens when you run it

1. Reads `hull.yaml` from `<package-path>`.
2. Refreshes upstream indexes (HTTP repo `index.yaml`, OCI registry tag list) unless `--skip-refresh` is set.
3. For each layer / require, resolves the constraint (`version:`, `ref:`) to a specific tag/commit/digest.
4. Writes `hull.lock` next to `hull.yaml`.
5. Prints a summary of changed layers (added, upgraded, removed) on stdout.

## Usage

```
hull dependency update <package-path> [name] [flags]
```

## Flags

| Flag | Type | Default | Description |
|---|---|---|---|
| `-h, --help` | bool | false | help for update |
| `--skip-refresh` | bool | false | skip repository index refresh before resolving |

## Persistent flags inherited from `hull`

| Flag | Type | Description |
|---|---|---|
| `--debug` | bool | enable debug output |
| `--kube-context` | string | Kubernetes context to use |
| `--kubeconfig` | string | path to kubeconfig file |
| `-n, --namespace` | string | Kubernetes namespace |

## Examples

Re-resolve every layer and rewrite `hull.lock`:

```sh
hull dependency update ./my-app
```

Update a single layer by name (other layers remain at their locked versions):

```sh
hull dependency update ./my-app shared-base
```

Skip the index refresh — useful in CI when an earlier step has already populated the cache:

```sh
hull dependency update ./my-app --skip-refresh
```

Pair with `dependency build` to actually fetch the resolved versions:

```sh
hull dependency update ./my-app
hull dependency build  ./my-app
```

## See also

- [`dependency`](dependency.md)
- [`dependency build`](dependency-build.md) — materialise the cache
- [`dependency list`](dependency-list.md) — inspect the lock
- [`dependency tree`](dependency-tree.md) — visualise the composition chain
- [Layers guide](../guides/layers.md)
