# hull workspace plan

## Synopsis

`hull workspace plan` reads `hull-workspace.yaml` and prints the topological order in which `hull workspace install` / `upgrade` would process the declared members. With `--levels`, output is grouped by Kahn topological depth — members in the same level have no inter-dependencies among themselves and are eligible to run in parallel (under `--parallel N` on install/upgrade). The command is read-only; no cluster contact, no file writes.

## When to use it

Run before `install` / `upgrade` to confirm dependency resolution, particularly after editing member `dependsOn` declarations. The `--levels` view tells you exactly where parallelism kicks in: a level with three members can run all three in parallel.

## What happens when you run it

1. Reads `<dir>/hull-workspace.yaml` (default: current directory).
2. Builds the dependency graph from each member's `dependsOn` list.
3. Computes the topological order (with cycle detection — a cycle aborts with a clear error naming the involved members).
4. Prints the plan: flat list by default; level-grouped with `--levels`.

## Usage

```
hull workspace plan [flags]
```

## Flags

| Flag | Type | Default | Description |
|---|---|---|---|
| `--dir` | string | . | directory containing `hull-workspace.yaml` |
| `-h, --help` | bool | false | help for plan |
| `--levels` | bool | false | group output by topological depth (parallelisable groups) |

## Persistent flags inherited from `hull`

| Flag | Type | Description |
|---|---|---|
| `--debug` | bool | enable debug output |
| `--kube-context` | string | Kubernetes context to use |
| `--kubeconfig` | string | path to kubeconfig file |
| `-n, --namespace` | string | Kubernetes namespace |

## Examples

Flat list of members in install order:

```sh
hull workspace plan
```

Level-grouped view (each group can run in parallel):

```sh
hull workspace plan --levels
```

Plan a workspace at a non-default path:

```sh
hull workspace plan --dir ./platform --levels
```

## See also

- [`workspace`](workspace.md)
- [`workspace install`](workspace-install.md)
- [`workspace upgrade`](workspace-upgrade.md)
- [`workspace status`](workspace-status.md)
- [`hull-workspace.yaml` reference](../reference/hull-workspace-yaml.md)
- [Workspaces guide](../guides/workspaces.md)
