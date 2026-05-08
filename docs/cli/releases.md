# hull releases

## Synopsis

`hull releases` orchestrates multiple separate releases declared in `hull-releases.yaml`. Subcommands install, upgrade, uninstall, plan, and report status across the whole graph using a topological order.

## When to use it

Use when you have a fleet of releases sourced from different places (local paths, OCI, HTTPS, git) with explicit dependency ordering between them. For releases that live in one repository, prefer `hull workspace`; for one-shot orchestration of disparate releases, this is the right tool.

## Usage

```
hull releases [command]
```

## Subcommands

- [`hull releases install`](releases-install.md) — install every release in topological order
- [`hull releases upgrade`](releases-upgrade.md) — upgrade every release; install if missing
- [`hull releases uninstall`](releases-uninstall.md) — uninstall every release in reverse topological order
- [`hull releases plan`](releases-plan.md) — print the topological order without applying
- [`hull releases status`](releases-status.md) — show current revision and status of every declared release

## Flags

| Flag | Type | Default | Description |
|---|---|---|---|
| `-h, --help` | — | — | help for releases |

## Persistent flags inherited from `hull`

| Flag | Type | Description |
|---|---|---|
| `--debug` | — | enable debug output |
| `--kube-context` | string | Kubernetes context to use |
| `--kubeconfig` | string | path to kubeconfig file |
| `-n, --namespace` | string | Kubernetes namespace |

## Examples

Plan an install of every release in `hull-releases.yaml`:

```sh
hull releases plan
```

Install every release in topological order:

```sh
hull releases install
```

Use a custom-named manifest:

```sh
hull releases install --file ./platform.releases.yaml
```

## See also

- [`hull-releases.yaml` reference](../reference/hull-releases-yaml.md)
- [Cross-release dependencies guide](../guides/releases.md)
- [`workspace`](workspace.md)
