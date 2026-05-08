# hull purge

## Synopsis

`hull purge` removes every release hull has installed across the cluster. By default it operates only on resources stamped with `managedBy=hull`; namespaces hull created are also discoverable by label. Always preview with `--dry-run` before running with `--yes`.

## When to use it

Use to clean up after a test sprawl, after a node failure that left wreckage, or to reset a development cluster. NOT a routine operation. Pair `--force` with `--delete-namespaces` for the most aggressive cleanup; `--force` includes pod- and namespace-finalizer escape hatches.

## Usage

```
hull purge [flags]
```

## Flags

| Flag | Type | Default | Description |
|---|---|---|---|
| `--delete-crds` | — | — | remove hull-installed CRDs (currently: hullreleases.hull.dev) |
| `--delete-namespaces` | — | — | after uninstall, delete every namespace that contained a hull release (excludes kube-*/default) |
| `--dry-run` | — | — | preview only; print what would be removed |
| `--exclude-ns` | strings | — | namespaces to skip (repeatable, comma-separated) |
| `--force` | — | — | skip graceful uninstall: delete storage Secrets directly, force-delete pods, force-finalize stuck namespaces (use after a node failure) |
| `-h, --help` | — | — | help for purge |
| `--ignore-failures` | — | — | keep going past failed uninstalls; report at the end |
| `--ns-prefix` | string | — | restrict scope to namespaces beginning with this prefix |
| `--parallel` | int | 4 | number of namespaces to purge concurrently |
| `--yes` | — | — | actually run; required when --dry-run is not set |

## Persistent flags inherited from `hull`

| Flag | Type | Description |
|---|---|---|
| `--debug` | — | enable debug output |
| `--kube-context` | string | Kubernetes context to use |
| `--kubeconfig` | string | path to kubeconfig file |
| `-n, --namespace` | string | Kubernetes namespace |

## Examples

Preview what would be removed:

```sh
hull purge --dry-run
```

Delete every hull-installed release:

```sh
hull purge --yes
```

Aggressive recovery after a kubelet failure (force-finalizes stuck namespaces, sweeps orphan pods):

```sh
hull purge --yes --force --delete-namespaces
```

Limit scope to namespaces with a given prefix:

```sh
hull purge --yes --ns-prefix hull-test --delete-namespaces
```

## See also

- [`uninstall`](uninstall.md)
