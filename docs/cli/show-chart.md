---
title: "hull show chart"
parent: "CLI"
---
{% raw %}
# hull show chart

`hull show chart` prints a package's `hull.yaml` metadata unchanged.

## When to use it

- Verify a package's name, version, and apiVersion before installing.
- Inspect an unfamiliar or freshly pulled package's manifest from the terminal.

## What happens

1. Reads `hull.yaml` from `<package-path>` (a directory or a hull archive).
2. Prints it verbatim to stdout. No layer resolution, no value merging.

## Usage

```
hull show chart <package-path>
```

## Flags

Inherits the global flags.

## Worked example

**INPUT** — `test/fixtures/simple/hull.yaml` on disk:

```yaml
apiVersion: hull/v1
name: simple-app
version: 1.0.0
```

**OUTPUT** (`hull show chart test/fixtures/simple`) — the same manifest, so you
read the package's identity directly:

```yaml
apiVersion: hull/v1
name: simple-app
version: 1.0.0
```

Pipe it through `yq` to pull a single field, for example
`hull show chart test/fixtures/simple | yq '.version'` prints `1.0.0`.

## See also

- [`show`](show.md) — the show command index
- [`show values`](show-values.md) — the package's default values
- [`show all`](show-all.md) — chart, values, and README together
{% endraw %}
