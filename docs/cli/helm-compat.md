# hull helm-compat

## Synopsis

`hull helm-compat` bridges hull and Helm in both directions. It renders and
installs unmodified upstream Helm charts under a hull release record, exports a
hull package into Helm's chart layout so Helm-only tooling can read it, and
reports how much of a Helm chart hull would have to translate to adopt it.

## Subcommands

| Command | What it does |
|---|---|
| [`export`](helm-compat-export.md) | write a hull package out as a Helm v3 chart |
| [`report`](helm-compat-report.md) | analyse a Helm chart and report its Go-template usage |
| `render` | render an unmodified Helm chart to manifests (like `helm template`) |
| `install` | render an unmodified Helm chart and apply it under a hull release |

## Usage

```
hull helm-compat <command> <path>
```

## See also

- [`migrate`](migrate.md) — convert a Helm chart into a native hull package
- [`template`](template.md) — render a hull package to manifests
