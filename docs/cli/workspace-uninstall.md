---
title: "hull workspace uninstall"
parent: "CLI"
---
{% raw %}
# hull workspace uninstall

## Synopsis

`hull workspace uninstall` removes every member declared in
`hull-workspace.yaml`, in **reverse** dependency order. Dependents come down
before the things they depend on, so nothing is torn out from under a member
that still needs it.

## When to use it

- To tear a whole workspace down cleanly in one command.
- When you want the removal to stop at the first failure (the default) or, with
  `--continue-on-error`, to remove as much as possible and report the rest.

## What happens

1. Reads `hull-workspace.yaml` from `--dir` (default `.`) and sorts the members
   into dependency levels.
2. Reverses the order, so the most-dependent level is removed first and the
   deepest dependencies last.
3. Within a level, up to `--parallel` members are uninstalled at once; the
   default of `1` removes them one at a time.
4. A member that has no release is skipped rather than treated as an error.
5. On failure, stops unless `--continue-on-error` is set.

Each member is removed the same way [`hull uninstall`](uninstall.md) would remove
it, in the member's `namespace`.

## Usage

```
hull workspace uninstall [flags]
```

## Flags

| Flag | Type | Default | Description |
|---|---|---|---|
| `--dir` | string | `.` | Directory containing `hull-workspace.yaml`. Point it elsewhere to uninstall a workspace in another directory. |
| `--parallel` | int | `1` | Maximum members to uninstall at once within a level. Raise it to remove independent members concurrently; `1` keeps them sequential. |
| `--continue-on-error` | bool | `false` | Keep uninstalling the remaining members after one fails, then report all failures at the end. |
| `--atomic-workspace` | bool | `false` | If any member fails, roll back the successful ones. Mutually exclusive with `--continue-on-error`. |
| `--dry-run` | bool | `false` | Report the plan without removing anything from the cluster. |
| `--health-gate` | bool | `false` | No effect on uninstall; the health-gate applies only to install and upgrade. |
| `--health-gate-timeout` | duration | `5m0s` | Paired with `--health-gate`; has no effect on uninstall. |
| `--progress` | bool | `false` | Print live lines as each member starts and finishes, plus a final summary. |

Inherits the global flags.

## Worked example

**INPUT — `hull-workspace.yaml` with two members**, where `api` depends on
`postgres` (both currently deployed):

```yaml
apiVersion: hull/v1
defaults:
  namespace: apps
members:
  - name: postgres
    path: ./postgres
  - name: api
    path: ./api
    dependsOn: [postgres]
```

**Run it with live progress:**

```sh
hull workspace uninstall --progress
```

**OUTPUT:**

```
workspace: 2 members across 2 level(s), parallel=1, op=uninstall

[level 0/1] 1 member(s) starting concurrently
  → api (ns=apps) start
  ✓ api done in 1.8s

[level 1/1] 1 member(s) starting concurrently
  → postgres (ns=apps) start
  ✓ postgres done in 1.5s

All 2 member(s) succeeded.
```

The order is reversed compared to install: `api` is removed first because it
depends on `postgres`, and `postgres` is removed last. If a member were already
gone, it would be skipped rather than reported as a failure. Without
`--progress`, a fully successful run prints nothing; failures are always
reported.

## See also

- [`workspace`](workspace.md) — the workspace index
- [`workspace install`](workspace-install.md) — the forward-order counterpart
- [`uninstall`](uninstall.md) — the single-release analogue
{% endraw %}
