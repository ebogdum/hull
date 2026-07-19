# Migrate a Helm chart to a hull package

`hull migrate` converts an existing **Helm chart** directory into a **hull
package**: it walks the chart (`Chart.yaml`, `templates/`, `values.yaml`,
`crds/`, `_helpers.tpl`, `NOTES.txt`) and emits an equivalent hull package,
rewriting Go-template constructs to hull's `${...}` expressions where it can.
Anything it cannot translate cleanly is flagged for manual review.

The command reference is [`hull migrate`](../cli/migrate.md). The companion
[`hull helm-compat`](../cli/helm-compat.md) runs an unmodified Helm chart under
hull without converting it, and exports a hull package back into Helm's layout.

## When to migrate

Migrate when you want to **own** an upstream chart as a hull package long-term
and any of these apply:

- You want hull's `${...}` expressions instead of Go-templates with sprig.
- You want hull's ownership labels, drift detection, audit trail, and signing.
- You want to slot the chart into a hull workspace beside hull-native packages.

If you only need to run an upstream chart as-is, you do not need migration —
`hull helm-compat install` renders and installs the unmodified chart under a
hull release record.

## Size the job first

Before converting, `hull helm-compat report` counts the Go-template logic in a
chart so you can gauge the work:

```sh
hull helm-compat report ./redis
```

```json
{
  "chart": "redis",
  "templates": 4,
  "goTemplateBlocks": 36,
  "notes": [
    "_helpers.tpl: 7 Go-template blocks (run 'hull migrate' to translate)",
    "deployment.yaml: 21 Go-template blocks (run 'hull migrate' to translate)"
  ],
  "recommendations": [
    "Run 'hull migrate ./redis' to translate go-template blocks to hull's ${...} syntax"
  ]
}
```

A chart with few `{{ ... }}` blocks converts with little effort; one packed with
them needs more review afterward.

## What the migrator produces

| Helm input | hull output |
|---|---|
| `Chart.yaml` | `hull.yaml` (`apiVersion: hull/v1`, layers/dependencies translated) |
| `values.yaml` | `values.yaml` |
| `values.schema.json` | `values.schema.json` |
| `templates/*.yaml` | `templates/*.yaml` (template body rewritten where possible) |
| `templates/_helpers.tpl` | `templates/_helpers.yaml` (named-template partials) |
| `templates/NOTES.txt` | `templates/notes.yaml` |
| `templates/tests/*` | `tests/*` |
| `crds/*.yaml` | `crds/*.yaml` |

Inside templates it rewrites a curated set of Go-template constructs to hull
expressions — for example:

| Go-template | hull |
|---|---|
| `{{ .Values.x }}` | `${values.x}` |
| `{{ .Release.Name }}` | `${release.name}` |
| `{{ if .Values.enabled }}` … `{{ end }}` | `${if .Values.enabled}` … `${end}` |
| `{{ range .Values.items }}` … `{{ end }}` | `${range .Values.items}` … `${end}` |
| `{{ include "named" . }}` | `${include "named"}` |
| `{{ toYaml .Values.x \| nindent 4 }}` | `${values.x \| toYaml \| nindent 4}` |

Constructs it cannot translate cleanly — some multi-variable `with`/`range`
forms, heavily nested conditionals around YAML structure, or calls to functions
hull does not implement — are left unchanged and listed for manual review.

## Convert

Point `hull migrate` at the chart. The package is written into `-o/--output`
(default `<chart-name>-hull/`); the conversion report prints to stdout — there
is no separate report file:

```sh
hull migrate ./redis -o ./redis-hull
```

```
Output: ./redis-hull
Converted 8 files:
  - hull.yaml
  - values.yaml
  - templates/_helpers.yaml
  - templates/deployment.yaml
  - templates/service.yaml
  - tests/test-connection.yaml
  - templates/notes.yaml
  - .hullignore

Migration complete.
```

When a construct needs a human, the report names the file, line, and reason
before `Migration complete.`:

```
Items requiring manual review (1):
  templates/deployment.yaml:24 — unsupported Helm function 'lookup'
    {{- $existing := lookup "v1" "Secret" .Release.Namespace "redis" }}
```

Use `--dry-run` to see the report without writing anything, or `--strict` to
fail the command on any template that cannot be fully auto-converted (useful as
a CI gate).

## Review and finish

```sh
ls ./redis-hull
```

```
hull.yaml  values.yaml  templates/  tests/  .hullignore
```

Lint the result, then resolve any flagged items by hand:

```sh
hull lint ./redis-hull
```

The migrator is deterministic and idempotent — re-running on the same input
produces the same output — so the workflow is: migrate, hand-edit the review
items, commit, then re-migrate when the upstream chart releases a new version
and diff to apply the changes.

## The inverse: expose a hull package to Helm tooling

`hull helm-compat` also runs the other direction, for CI or GitOps tools that
only understand Helm:

```sh
# Render an unmodified Helm chart to manifests (like `helm template`):
hull helm-compat render ./redis

# Install an unmodified Helm chart under a hull release record:
hull helm-compat install redis ./redis -n data

# Export a hull package AS a Helm v3 chart:
hull helm-compat export ./my-pkg --out ./helm-export
```

```
exported helm-compat chart to ./helm-export
```

`export` writes a `Chart.yaml` (`apiVersion: v2`) plus the copied `values.yaml`
and `templates/` tree. It is best-effort: hull's `${...}` expressions are copied
verbatim and resolve only under hull, so the exported chart is suitable for
static analysis, not for `helm install`. To hand Helm fully rendered manifests,
pre-render with `hull template` and ship the output.

## Limits

`hull migrate` is **template translation**, not behaviour emulation:

- Release-name interpolation works the same way, so generated resource names
  match — but tooling that scrapes hull's release records reads a hull-specific
  schema, not Helm's `sh.helm.release.v1...` secrets.
- Helm's `helm.sh/hook: test` test pattern is rewritten to hull's lifecycle.

The goal is a working hull package you then own and maintain — not to run Helm
under hull forever.

## See also

- [`hull migrate`](../cli/migrate.md) — command reference
- [`hull helm-compat`](../cli/helm-compat.md) — render/install/export/report
- [`hull helm-compat export`](../cli/helm-compat-export.md) ·
  [`hull helm-compat report`](../cli/helm-compat-report.md)
- [`hull lint`](../cli/lint.md) — validate the converted package
- [Workspaces](workspaces.md) — slot the migrated package into a workspace
