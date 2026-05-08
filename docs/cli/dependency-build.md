# hull dependency build

## Synopsis

`hull dependency build` resolves every layer and required package declared in `hull.yaml`, downloads each one to the package's local cache (`./.hull/layers/`), and verifies the lockfile is consistent. After a successful build, `hull install` and `hull template` can render the package fully offline. The command is the materialise-side counterpart to `hull dependency update` (which writes `hull.lock`).

## When to use it

Run after cloning a package fresh to a new machine, after `hull dependency update` to prefetch the new versions, or in CI to pre-populate the layer cache before running `hull lint` or `hull template`. Idempotent: re-running with no changes is a fast no-op.

## What happens when you run it

1. Hull reads `hull.yaml` and `hull.lock`.
2. For every layer / require, hull fetches the source (local copy, HTTPS archive, OCI artifact, git clone) into `./.hull/layers/<name>/`.
3. Each downloaded archive's digest is checked against `hull.lock`. With `--verify` set, the digest check fails the build instead of silently regenerating the lock entry.
4. Index caches are reused unless `--no-cache` is set.
5. On success, prints the resolved version of every layer; on failure, names the offending layer.

## Usage

```
hull dependency build <package-path> [flags]
```

## Flags

| Flag | Type | Default | Description |
|---|---|---|---|
| `-h, --help` | bool | false | help for build |
| `--no-cache` | bool | false | clear index cache before resolving |
| `--verify` | bool | false | verify digests of installed dependencies |

## Persistent flags inherited from `hull`

| Flag | Type | Description |
|---|---|---|
| `--debug` | bool | enable debug output |
| `--kube-context` | string | Kubernetes context to use |
| `--kubeconfig` | string | path to kubeconfig file |
| `-n, --namespace` | string | Kubernetes namespace |

## Examples

Materialise every layer for a package:

```sh
hull dependency build ./my-app
```

Strict mode for CI: any digest mismatch fails the build:

```sh
hull dependency build ./my-app --verify
```

Bust the index cache before resolving — useful when an upstream repository's `index.yaml` changed and you want a clean fetch:

```sh
hull dependency build ./my-app --no-cache
```

## See also

- [`dependency`](dependency.md)
- [`dependency update`](dependency-update.md) — refresh `hull.lock` before building
- [`dependency list`](dependency-list.md) — show what is declared
- [`dependency tree`](dependency-tree.md) — visualise the composition chain
- [Layers guide](../guides/layers.md)
