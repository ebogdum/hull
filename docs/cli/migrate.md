# hull migrate

## Synopsis

`hull migrate` translates a Helm chart directory into a hull package. It walks the chart's structure (`Chart.yaml`, `templates/`, `values.yaml`, `crds/`, `_helpers.tpl`, `NOTES.txt`, `requirements.yaml`/`Chart.lock`) and emits an equivalent hull package, rewriting go-template constructs to hull `${...}` expressions where possible. Constructs that cannot be cleanly translated are flagged in a `hull-migration.md` report inside the output directory.

## When to use it

Use when adopting an upstream Helm chart as a hull-owned package. The output is a starting point that you then own and edit; the migrator is a translation tool, not a 1:1 emulator.

## Usage

```
hull migrate <helm-chart-path> [flags]
```

## Flags

| Flag | Type | Default | Description |
|---|---|---|---|
| `--dry-run` | — | — | show what would be converted without writing |
| `-h, --help` | — | — | help for migrate |
| `-o, --output` | string | — | output directory (default: <chart-name>-hull/) |
| `--strict` | — | — | fail on any template that cannot be fully auto-converted |

## Persistent flags inherited from `hull`

| Flag | Type | Description |
|---|---|---|
| `--debug` | — | enable debug output |
| `--kube-context` | string | Kubernetes context to use |
| `--kubeconfig` | string | path to kubeconfig file |
| `-n, --namespace` | string | Kubernetes namespace |

## Examples

Migrate an upstream Helm chart:

```sh
hull migrate ./upstream-chart -d ./migrated/
```

Lint the migrated package:

```sh
hull lint ./migrated/<chart-name>
```

## See also

- [Migration guide](../guides/migration.md)
- [`helm-compat`](helm-compat.md)
