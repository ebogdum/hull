# hull migrate

Convert a Helm chart directory on disk into an equivalent hull package
directory.

## When to use it

- Adopting an upstream Helm chart as a hull-owned package you will maintain
  going forward.
- Moving your own charts off Helm and onto hull, one chart at a time.

The output is a starting point you then own and edit. `migrate` is a translation
tool, not a 1:1 Helm emulator — expect to review and finish some templates by
hand. To keep running an unmodified Helm chart as-is instead of converting it,
use [`hull helm-compat`](helm-compat.md).

## What happens

1. Reads the Helm chart at `<helm-chart-path>` — its `Chart.yaml`, `templates/`,
   `values.yaml`, and related files.
2. Converts the templates and metadata to hull's package format, rewriting Helm
   template constructs to their hull equivalents where it can.
3. Writes the result to a new package directory — `<chart-name>-hull/` by
   default, or wherever `-o/--output` points. With `--dry-run` nothing is
   written and the conversion is only reported.
4. Prints the output path, the list of converted files, any items that need
   manual review (with the file, line, and reason), and any warnings.

This works entirely on local files. It reads no cluster and installs nothing.

## Usage

```
hull migrate <helm-chart-path> [flags]
```

## Flags

| Flag | Type | Default | Description |
|---|---|---|---|
| `-o, --output` | string | `<chart-name>-hull/` | directory to write the converted package to |
| `--dry-run` | — | false | show what would be converted without writing anything |
| `--strict` | — | false | fail on any template that cannot be fully auto-converted |

## Persistent flags inherited from `hull`

| Flag | Type | Description |
|---|---|---|
| `--debug` | — | enable debug output |
| `--kube-context` | string | Kubernetes context to use |
| `--kubeconfig` | string | path to kubeconfig file |
| `-n, --namespace` | string | Kubernetes namespace |

## Worked example

**INPUT — an upstream Helm chart directory:**

```sh
ls ./redis
```

```
Chart.yaml  values.yaml  templates/
```

```
./redis/templates/
  deployment.yaml
  service.yaml
  _helpers.tpl
  NOTES.txt
```

**Convert it:**

```sh
hull migrate ./redis -o ./redis-hull
```

**OUTPUT — a hull package plus a conversion report:**

```
Output: ./redis-hull
Converted 4 files:
  - hull.yaml
  - values.yaml
  - templates/deployment.yaml
  - templates/service.yaml

Items requiring manual review (1):
  templates/deployment.yaml:24 — unsupported Helm function 'lookup'
    {{- $existing := lookup "v1" "Secret" .Release.Namespace "redis" }}

Migration complete.
```

```sh
ls ./redis-hull
```

```
hull.yaml  values.yaml  values.schema.json  templates/
```

Lint the result, then fix anything flagged for manual review:

```sh
hull lint ./redis-hull
```

## See also

- [`helm-compat`](helm-compat.md) — run an unmodified Helm chart without converting it
- [`adopt`](adopt.md) — bring existing in-cluster resources under hull management
- [`lint`](lint.md) — validate the converted package
