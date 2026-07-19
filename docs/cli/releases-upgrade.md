---
title: "hull releases upgrade"
parent: "CLI"
---
{% raw %}
# hull releases upgrade

## Synopsis

`hull releases upgrade` upgrades every release in your `hull-releases.yaml`, in
dependency order, and installs any release that is not there yet. It is the
one command you re-run to keep the whole set current: it works the first time
(everything installs) and every time after (everything upgrades).

## When to use it

- As your repeatable deploy command for the set — safe to run whether or not
  the releases already exist.
- After editing a package, its values, or its `set` overrides, to roll the
  change out across the set in the right order.

## What happens

1. Reads the spec file (`--file`, default `hull-releases.yaml`) and sorts the
   releases into dependency order. A cycle or unknown `dependsOn` name stops
   here, before anything is applied.
2. For each release in order, upgrades its package under the release `name`, in
   the entry's `namespace` (falling back to `-n/--namespace`), applying its
   `profile`, `values`, and `set`. A release that does not exist yet is
   installed instead.
3. Each release is upgraded atomically: if one fails, it rolls back to its
   previous revision, the command stops, and the releases already processed are
   left as they are.
4. Prints one line per release, with its new revision number.

Like `install`, this does not wait for pods to become ready between releases.
For that, use [`hull workspace`](workspace.md) with `--health-gate`.

## Usage

```
hull releases upgrade [flags]
```

## Flags

| Flag | Type | Default | Effect |
|---|---|---|---|
| `--file` | string | `hull-releases.yaml` | read the spec from this path instead of the default |

Inherits the global flags (`--kube-context`, `--kubeconfig`, `-n/--namespace`,
`--debug`). `-n` sets the namespace for any release that does not name its own.

## Worked example

**INPUT** — `hull-releases.yaml`. `postgres` and `redis` are already installed
at revision 1; you bumped the image in `./charts/api` and now roll it out:

```yaml
releases:
  - name: postgres
    package: ./charts/postgres
    namespace: data

  - name: redis
    package: ./charts/redis
    namespace: data

  - name: api
    package: ./charts/api
    namespace: apps
    dependsOn:
      - postgres
      - redis
```

**COMMAND:**

```sh
hull releases upgrade
```

**OUTPUT:**

```
[postgres] upgraded (revision 2)
[redis] upgraded (revision 2)
[api] upgraded (revision 2)
```

Every release moves up one revision, in dependency order — `postgres` and
`redis` before `api`. Had `api` not existed yet, its line would instead read
`[api] installed (revision 1, ns apps)`.

## See also

- [`releases`](releases.md) — the parent command and the spec-file format
- [`plan`](releases-plan.md) — preview the order first
- [`install`](releases-install.md) — first-time install of the set
- [`upgrade`](upgrade.md) — upgrade a single release
{% endraw %}
