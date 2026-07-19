# hull workspace

## Synopsis

A **workspace** is a single file — `hull-workspace.yaml` — that lists several
hull packages and how they depend on each other. Instead of installing,
upgrading, or tearing down each release by hand, you point `hull workspace` at
the directory holding that file and it drives every member for you in
dependency order.

```yaml
# hull-workspace.yaml
apiVersion: hull/v1
defaults:
  namespace: apps            # used by any member that omits its own
members:
  - name: postgres
    path: ./postgres
  - name: api
    path: ./api
    dependsOn: [postgres]    # api is processed only after postgres
```

Each member names a release, a package `path`, and optionally a `namespace`,
`profile`, and a `dependsOn` list. `hull` reads the `dependsOn` edges, sorts the
members into dependency order, and processes them:

- **install** and **upgrade** run in dependency order — a member starts only
  after everything it depends on has finished.
- **uninstall** runs in **reverse** — dependents come down before the things
  they depend on.

Members that do not depend on each other share a level and can run concurrently
(see `--parallel`).

## Subcommands

| Command | What it does |
|---|---|
| [`plan`](workspace-plan.md) | Print the order the members will be processed in |
| [`install`](workspace-install.md) | Install every member in dependency order |
| [`upgrade`](workspace-upgrade.md) | Upgrade every member (installing any that are missing) |
| [`uninstall`](workspace-uninstall.md) | Uninstall every member in reverse order |
| [`status`](workspace-status.md) | Show each member's revision and status |
| [`diff`](workspace-diff.md) | Show what upgrading every member would change |

## Usage

```
hull workspace <command> [flags]
```

Every command reads `hull-workspace.yaml` from the directory given by `--dir`
(default `.`).

## See also

- [`install`](install.md), [`upgrade`](upgrade.md), [`uninstall`](uninstall.md)
  — the single-release commands the workspace runs once per member
- [`plan`](plan.md), [`status`](status.md), [`diff`](diff.md) — the
  single-release analogues of the workspace subcommands
