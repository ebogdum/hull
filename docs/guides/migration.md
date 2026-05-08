# Migrate a Helm chart to a hull package — Helm to hull conversion guide

This guide walks through converting (migrating) an existing **Helm chart** into a **hull package** using the `hull migrate` command. If you're searching for a **Helm chart converter**, **Helm chart migration tool**, or **how to import a Helm chart into hull**, you're in the right place.

The `hull migrate` command translates a Helm chart directory into a hull package: it walks the Helm chart structure (`Chart.yaml`, `templates/`, `values.yaml`, `crds/`, `_helpers.tpl`, `NOTES.txt`, `requirements.yaml`/`Chart.lock`) and emits an equivalent hull package, rewriting go-template constructs to hull's `${...}` expressions where possible. Constructs the migrator can't translate cleanly are flagged in a `hull-migration.md` review report inside the output directory.

The companion `hull helm-compat` command provides the inverse direction: rendering a hull package as a Helm-compatible artifact for downstream tooling that consumes Helm output (e.g. `helm template`-driven CI gates, Helm-aware OCI scanners, GitOps tools that only know Helm).

> **Glossary search hooks:** "Helm to hull migration", "convert Helm chart", "Helm chart to hull package", "Helm migrator", "Helm chart converter", "import Helm chart into hull", "Helm-compat hull", "Helm chart hull replacement".

## When to migrate

You have an upstream Helm chart you want to install through hull, **and** any of:

- You want hull's expression syntax instead of go-templates with sprig.
- You want hull's ownership labels, drift detection, audit trail, and signing.
- You want to slot the upstream chart into a hull workspace alongside hull-native packages.

If you only need a one-shot install of an upstream chart, you don't need migration — `hull install` accepts a Helm chart's tarball or directory directly via the compatibility layer (`hull helm-compat install`).

Migration is for **owning** the package long-term.

## The migrator's job

`hull migrate` walks a Helm chart directory and produces a hull package directory:

| Helm input | hull output |
|---|---|
| `Chart.yaml` | `hull.yaml` (with `apiVersion: hull/v1`, layers translated, dependencies translated) |
| `values.yaml` | `values.yaml` (unchanged) |
| `values.schema.json` | `values.schema.json` (unchanged) |
| `templates/*.yaml` | `templates/*.yaml` (template body rewritten where possible) |
| `templates/_helpers.tpl` | `templates/_helpers.yaml` (named templates → hull `${define}` partials) |
| `templates/NOTES.txt` | `notes.yaml` |
| `crds/*.yaml` | `crds/*.yaml` (unchanged) |
| `Chart.lock` | `hull.lock` |
| `requirements.yaml` (Helm v2) | layers entries in `hull.yaml` |

Inside templates, the migrator translates a curated set of go-template constructs to hull expressions:

| Go-template | hull |
|---|---|
| `{{ .Values.x }}` | `${values.x}` |
| `{{ .Release.Name }}` | `${release.name}` |
| `{{ if .Values.enabled }}` ... `{{ end }}` | `${if .Values.enabled}` ... `${end}` |
| `{{ range .Values.items }}` ... `{{ end }}` | `${range .Values.items}` ... `${end}` |
| `{{ toYaml .Values.x | nindent 4 }}` | `${values.x | toYaml | nindent 4}` |
| `{{ printf "%s-%s" $a $b }}` | `${printf "%s-%s" $a $b}` |
| `{{ tpl .Values.foo . }}` | `${tpl .Values.foo}` |
| `{{ lookup "v1" "Secret" "default" "x" }}` | `${lookup "v1" "Secret" "default" "x"}` |
| `{{ include "named" . }}` | `${include "named"}` |
| `{{ index .Values "foo" "bar" }}` | `${get .Values "foo" "bar"}` |

Conditional blocks, range blocks, and sprig functions (math, string, regex, date, crypto, etc.) are passed through with hull's equivalents — hull's expression engine implements every sprig function the migrator can't otherwise translate.

## Things the migrator can't translate

When the migrator finds a construct it can't translate cleanly, it emits the original token unchanged AND adds an entry to the migration's review list:

- `{{ with $foo := ... }}` over multiple variables.
- `{{ range $i, $e := ... }}` with explicit index naming.
- Heavily nested conditionals around YAML structure (which sometimes break a 1:1 line translation).
- Calls to functions hull doesn't implement (rare; the migrator names them).

## Workflow

```sh
hull migrate ./upstream-chart -d ./migrated/
# walks upstream-chart, writes ./migrated/<chart-name>/...
hull lint ./migrated/<chart-name>
# review hull-migration.md inside the output:
cat ./migrated/<chart-name>/hull-migration.md
```

The `hull-migration.md` report lists:

- What the migrator did automatically.
- What it left for manual review (with file + line references).
- Warnings (deprecated fields, ambiguous translations).

## Iteration

The migrator is deterministic and idempotent — re-running on the same input produces the same output. It's safe to:

1. Run `hull migrate` to get a starting point.
2. Hand-edit the output to clean up review items.
3. Commit.
4. Re-run later when the upstream chart releases a new version, diff the output, and apply the upstream changes selectively.

## Compatibility layer (the inverse)

`hull helm-compat` exposes hull packages to Helm-aware tooling:

```sh
hull helm-compat template my-app ./my-pkg                # Helm-shaped manifest output
hull helm-compat export ./my-pkg -d ./helm-export/       # writes a Helm chart skeleton
```

This is useful when CI runs `helm-diff`, `helm-conftest`, or `helm-secrets`-style tools and you want hull to look like the Helm chart it would have been. The export is best-effort — hull's `${...}` expressions may end up as inert literal strings in the exported chart, so the export is suitable for static analysis but not for actual Helm install.

## Limits

The migrator is **template translation**, not behaviour translation. If the upstream chart relies on:

- A specific Helm release-name format used by clients later (`my-release-redis-master`), the migrated hull package generates the same names (release name interpolation works the same way), but workflow tooling that scrapes the release record may need updates.
- The Helm release Secret schema (`sh.helm.release.v1.<release>.v<rev>`), hull's Secret schema differs (`hull.v1.<release>.v<rev>` with hull-specific labels). Tools reading Helm releases directly need updating; `hull helm-compat list` returns the hull-style data.
- A specific Helm-only test pattern (`helm test` with `helm.sh/hook: test`), `hull migrate` rewrites it to `$hook: test` with hull's lifecycle.

## Summary

`hull migrate` exists because most Kubernetes packages today are Helm charts and there's no point making users rewrite from scratch. It's a translation tool, not a 1:1 emulator — the goal is to land a working hull package that you then own, not to run Helm under hull forever.
