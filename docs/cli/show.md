# hull show

`hull show` prints a package's metadata, values, README, or CRDs without
installing it.

## Subcommands

| Command | Displays |
|---|---|
| [`hull show chart`](show-chart.md) | the package metadata from `hull.yaml` |
| [`hull show values`](show-values.md) | the default `values.yaml` |
| [`hull show readme`](show-readme.md) | the package's README file |
| [`hull show crds`](show-crds.md) | the CRDs under the package's `crds/` directory |
| [`hull show all`](show-all.md) | chart, values, and README in one document |

Each subcommand takes a `<package-path>` that is a directory or a hull archive
(`.hull.tgz`, `.tgz`, `.tar.gz`) and reads the file straight from disk.

## Usage

```
hull show <chart|values|readme|crds|all> <package-path> [flags]
```

## See also

- [`values`](values.md) — resolve and trace the merged values
- [`template`](template.md) — render the package's manifests
- [`lint`](lint.md) — validate a package
