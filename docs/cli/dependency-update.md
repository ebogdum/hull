# hull dependency update

`hull dependency update` re-resolves every layer and required package against
its source and rewrites `hull.lock` with the pinned versions and commits.

## When to use it

- After editing `layers:` or `requires:` in `hull.yaml`, to refresh the lock.
- To move a `git::` or registry source forward to the latest ref or version
  its constraint allows.

## What happens

hull reads `hull.yaml` and, for each layer and required package, resolves its
`source`:

- **git** sources are fetched and the exact commit is recorded as
  `resolvedCommit`.
- **registry** sources record the selected `resolvedVersion`.
- **local** sources are recorded by path.

hull then writes `hull.lock` with a fresh `generated` timestamp and one entry
per layer and require. This lock is what [`build`](dependency-build.md),
[`install`](install.md), and [`template`](template.md) read first, so every
later render uses the same pinned versions. On success hull prints
`Layers updated successfully.`

## Usage

```
hull dependency update <package-path> [name]
```

The optional `[name]` argument updates a single legacy `dependencies:` entry;
for `layers:`/`requires:` packages, all entries are re-resolved together.

## Flags

| Flag | Cause → effect |
|---|---|
| `--skip-refresh` | skip refreshing repository indexes first; resolve against the index already cached — faster, but may miss versions published since the last refresh |

## Worked example

**INPUT — `./web/hull.yaml`** with two layers and one required package, no
lock yet:

```yaml
apiVersion: hull/v1
name: web
version: 0.3.0
layers:
  - name: base-layer
    source: ../base-layer
  - name: common
    source: ../common-layer
requires:
  - name: redis
    source: ../redis-req
```

**Command:**

```sh
hull dependency update ./web
```

**OUTPUT:**

```
Layers updated successfully.
```

**`./web/hull.lock` that gets written:**

```yaml
apiVersion: hull/v1
generated: 2026-07-18T23:47:45.658955+02:00
layers:
    - name: base-layer
      source: ../base-layer
    - name: common
      source: ../common-layer
requires:
    - name: redis
      source: ../redis-req
```

**Tracing each input to its lock entry:**

| `hull.yaml` entry | `hull.lock` entry | Why |
|---|---|---|
| `layers[0]` `base-layer` | `layers: - name: base-layer` | resolved and pinned under `layers:` |
| `layers[1]` `common` | `layers: - name: common` | second layer, pinned |
| `requires[0]` `redis` | `requires: - name: redis` | requires are pinned in their own block |

A `git::` source would add `resolvedCommit:`, and a registry source
`resolvedVersion:`, so the lock records the exact bits, not just the name.
Run [`dependency list`](dependency-list.md) afterward and every `STATUS`
reads `locked`.

## See also

- [`dependency build`](dependency-build.md) — download what this lock pins
- [`dependency list`](dependency-list.md) — confirm entries are now `locked`
- [`install`](install.md) — install using the pinned versions
