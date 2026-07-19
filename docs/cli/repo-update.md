---
title: "hull repo update"
parent: "CLI"
---
{% raw %}
# hull repo update

Re-fetch the index of every registered repository.

## When to use it

- After [`hull repo add`](repo-add.md), so the new repository's charts become
  searchable and pullable.
- Before searching or pulling, to pick up newly published chart versions.

## What happens

Hull walks your repository list and, for each entry, fetches its `index.yaml`
into the local cache under `~/.cache/hull/indexes/`, using the credentials and
TLS material recorded for that repository. It prints one line per repository
and a final `Update complete.` Once the cache is fresh,
[`hull search`](search.md) and [`hull pull`](pull.md) see the current set of
charts.

A repository that cannot be reached is reported on its own line but does not
stop the others; pass `--fail-on-repo-update-fail` to make any failure exit
non-zero (useful in CI). With no repositories registered, hull prints
`No repositories configured.`

## Usage

```
hull repo update [flags]
```

## Flags

| Flag | Effect |
|---|---|
| `--fail-on-repo-update-fail` | Exit non-zero if any repository fails to update. |

## Worked example

```
$ hull repo update
...successfully got an update from "my-charts"
...successfully got an update from "private"
Update complete.
```

Refresh, then search the freshly updated catalogues:

```
$ hull repo update
...successfully got an update from "my-charts"
Update complete.

$ hull search repo redis
NAME              CHART VERSION   APP VERSION   DESCRIPTION
my-charts/redis   1.4.0           7.2.4         In-memory data store
```

## See also

- [`repo add`](repo-add.md) — register a repository first
- [`search`](search.md) — search the updated indexes
- [`pull`](pull.md) — download a chart
{% endraw %}
