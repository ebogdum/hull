---
title: "hull dependency build"
parent: "CLI"
---
{% raw %}
# hull dependency build

`hull dependency build` resolves every layer and required package declared in
`hull.yaml` and downloads each one so the package is ready to render.

## When to use it

- After [`update`](dependency-update.md), to fetch the pinned dependencies onto
  this machine.
- Before rendering or installing offline, so all remote sources are already in
  hull's local cache.

## What happens

hull reads `hull.yaml` and resolves each layer and required package. If a
current `hull.lock` exists, it fetches the exact pinned versions and commits;
otherwise it resolves fresh and writes the lock. Remote sources (`git::`,
registry) are downloaded into hull's local cache; local paths need no
download. Once this succeeds, [`template`](template.md) and [`install`](install.md)
can render the package without reaching the network again. On success hull
prints `Dependencies resolved successfully.`

## Usage

```
hull dependency build <package-path>
```

## Flags

| Flag | Cause → effect |
|---|---|
| `--no-cache` | clear the cached repository index before resolving, so versions are re-read from the source instead of the last-fetched index |
| `--verify` | after downloading, check each installed dependency's digest against `hull.lock`; fail if any does not match |

## Worked example

**INPUT — `./web/hull.yaml`** with two layers and one required package, already
pinned in `hull.lock` by a prior `dependency update`:

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
hull dependency build ./web
```

**OUTPUT:**

```
Dependencies resolved successfully.
```

**What that line means, traced to the input:**

| `hull.yaml` entry | What build did |
|---|---|
| `layers[0]` `base-layer` | resolved `../base-layer` at its locked version, ready to merge |
| `layers[1]` `common` | resolved `../common-layer` at its locked version |
| `requires[0]` `redis` | resolved `../redis-req`, ready to install alongside `web` |

With `--verify`, hull additionally compares each downloaded dependency's digest
to the one recorded in `hull.lock` and stops with an error if they differ —
proof the fetched bytes match what was pinned.

## See also

- [`dependency update`](dependency-update.md) — pin versions into `hull.lock` first
- [`dependency tree`](dependency-tree.md) — see what will be downloaded
- [`install`](install.md) — install the package once dependencies are built
{% endraw %}
