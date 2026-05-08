# hull dependency

## Synopsis

`hull dependency` manages a package's declared layers and required co-deployed packages. Subcommands list the current state, update the lockfile, build (materialise) every layer into the local cache, and print a tree visualisation of the resolved dependency graph.

## When to use it

Use whenever `hull.yaml`'s `layers:` or `requires:` change, before committing the package. `hull dependency update` is the canonical "refresh my lock" command.

## Usage

```
hull dependency [command]
```

## Subcommands

- [`hull dependency list`](dependency-list.md) — List layers and required packages with their status
- [`hull dependency tree`](dependency-tree.md) — Display the layer composition chain
- [`hull dependency update`](dependency-update.md) — Re-resolve layer and dependency versions

## Flags

| Flag | Type | Default | Description |
|---|---|---|---|
| `-h, --help` | — | — | help for dependency |

## Persistent flags inherited from `hull`

| Flag | Type | Description |
|---|---|---|
| `--debug` | — | enable debug output |
| `--kube-context` | string | Kubernetes context to use |
| `--kubeconfig` | string | path to kubeconfig file |
| `-n, --namespace` | string | Kubernetes namespace |

## Examples

List the package's declared layers and their lock state:

```sh
hull dependency list ./my-app
```

Update the lockfile (resolve every layer's source, version, ref to a digest):

```sh
hull dependency update ./my-app
```

Materialise every locked layer into `./my-app/.hull/layers/`:

```sh
hull dependency build ./my-app
```

## See also

- [Layers guide](../guides/layers.md)
- [`hull.yaml` reference](../reference/hull-yaml.md)
