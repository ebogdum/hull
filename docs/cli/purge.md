# hull purge

Find every release hull has installed anywhere in the cluster, uninstall each
one, and optionally delete the namespaces and CRDs hull created.

## When to use it

- Tearing down a test or CI cluster: remove everything hull put there in one
  command instead of uninstalling releases one at a time.
- Recovering after a node or cluster failure that left releases half-removed —
  `--force` deletes storage directly and unsticks Terminating namespaces.

`hull purge` is the cluster-wide broom that [`hull uninstall`](uninstall.md) is
not: `uninstall` removes one named release, `purge` sweeps them all.

## What happens

1. Connects to the cluster and lists every release-storage Secret labelled
   `managedBy=hull` (or the legacy `owner=hull`), reconstructing each release's
   name and namespace. Resources hull never installed are never touched.
2. Also picks up any namespace hull labelled at creation time, even if its
   releases are already gone.
3. Applies your scope filters: `--ns-prefix` narrows to matching namespaces,
   `--exclude-ns` protects specific ones. `kube-system`, `kube-public`,
   `kube-node-lease`, and `default` are always protected.
4. Without `--yes`, prints the plan and stops — nothing changes. With `--yes`,
   uninstalls every in-scope release (deleting its stored history), then
   optionally deletes namespaces (`--delete-namespaces`) and hull CRDs
   (`--delete-crds`).

`--yes` mutates the cluster and requires a reachable cluster. Running with
neither `--dry-run` nor `--yes` is refused.

## Usage

```
hull purge [flags]
```

## Flags

| Flag | Type | Default | Description |
|---|---|---|---|
| `--dry-run` | — | false | preview only; print what would be removed and exit |
| `--yes` | — | false | actually run; required when `--dry-run` is not set |
| `--ns-prefix` | string | — | restrict scope to namespaces beginning with this prefix |
| `--exclude-ns` | stringArray | — | namespaces to skip (repeatable, comma-separated) |
| `--delete-namespaces` | — | false | after uninstall, delete every namespace that held a hull release (never `kube-*`/`default`) |
| `--delete-crds` | — | false | remove hull-installed CRDs (currently `hullreleases.hull.dev`) |
| `--parallel` | int | 4 | number of namespaces to purge concurrently |
| `--ignore-failures` | — | false | keep going past failed uninstalls; report the count at the end |
| `--force` | — | false | skip graceful uninstall: delete storage Secrets directly, force-delete pods, force-finalize stuck namespaces (use after a node failure) |

## Persistent flags inherited from `hull`

| Flag | Type | Description |
|---|---|---|
| `--debug` | — | enable debug output |
| `--kube-context` | string | Kubernetes context to use |
| `--kubeconfig` | string | path to kubeconfig file |
| `-n, --namespace` | string | Kubernetes namespace |

## Worked example

**INPUT — three hull releases live in test namespaces:**

```sh
hull list --all-namespaces
```

```
NAME     NAMESPACE   REVISION   STATUS     UPDATED
api      hull-test   3          deployed   2026-07-18 09:10:22
cache    hull-test   1          deployed   2026-07-18 09:11:04
worker   hull-ci     2          deployed   2026-07-17 22:40:51
```

**Preview first — nothing is touched:**

```sh
hull purge --dry-run --ns-prefix hull-
```

```
found 3 release(s) across 2 namespace(s) in scope

--- DRY RUN ---
  ns/hull-ci
    └ release/worker
  ns/hull-test
    └ release/api
    └ release/cache

(re-run with --yes to actually purge)
```

**Execute — uninstall the releases and drop the namespaces:**

```sh
hull purge --yes --ns-prefix hull- --delete-namespaces
```

```
found 3 release(s) across 2 namespace(s) in scope
  ✓ hull-ci/worker [1/3]
  ✓ hull-test/api [2/3]
  ✓ hull-test/cache [3/3]
  ✓ ns/hull-ci deleted
  ✓ ns/hull-test deleted
purge complete
```

**OUTPUT — the releases and their history are gone:**

```sh
hull list --all-namespaces
```

```
NAME   NAMESPACE   REVISION   STATUS   UPDATED
```

Each release's stored revisions are deleted along with it, so
[`hull history`](history.md) and `hull rollback` have nothing left to act on.

## See also

- [`uninstall`](uninstall.md) — remove one named release instead of all of them
- [`prune`](prune.md) — trim old revisions of a release without removing it
- [`list`](list.md) — see what is in scope before you purge
