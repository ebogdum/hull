# hull workspace

## Synopsis

`hull workspace` orchestrates multiple hull packages declared in `hull-workspace.yaml`. Subcommands install, upgrade, uninstall, plan, diff, and report status across the whole workspace using a topological order computed from `dependsOn` declarations.

## When to use it

Use when many sibling packages from one repository should roll out together with explicit dependency ordering between them. For releases sourced from disparate places (different registries, paths, repos), see `hull releases`.

## Usage

```
hull workspace [command]
```

## Subcommands

- [`hull workspace install`](workspace-install.md) — install every member in topological order
- [`hull workspace upgrade`](workspace-upgrade.md) — upgrade every member; install if missing
- [`hull workspace uninstall`](workspace-uninstall.md) — uninstall every member in reverse topological order
- [`hull workspace plan`](workspace-plan.md) — print the install plan with optional level grouping
- [`hull workspace status`](workspace-status.md) — show current revision and status of every declared member
- [`hull workspace diff`](workspace-diff.md) — show pending changes per member

## Flags

| Flag | Type | Default | Description |
|---|---|---|---|
| `-h, --help` | — | — | help for workspace |

## Persistent flags inherited from `hull`

| Flag | Type | Description |
|---|---|---|
| `--debug` | — | enable debug output |
| `--kube-context` | string | Kubernetes context to use |
| `--kubeconfig` | string | path to kubeconfig file |
| `-n, --namespace` | string | Kubernetes namespace |

## Examples

Plan a workspace install:

```sh
hull workspace plan .
```

Install with parallelism within levels and a health-gate between levels:

```sh
hull workspace install . --parallel 4 --health-gate
```

Diff every member's pending changes:

```sh
hull workspace diff .
```

## See also

- [`hull-workspace.yaml` reference](../reference/hull-workspace-yaml.md)
- [Workspaces guide](../guides/workspaces.md)
- [`releases`](releases.md)
