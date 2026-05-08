# hull releases install

## Synopsis

Install every release declared in `hull-releases.yaml` in topological order. The graph is computed from each entry's `dependsOn` list; releases at the same topological level have no inter-dependencies among themselves.

## When to use it

Use to bring up a fresh fleet of related releases with one command — typically a platform bootstrap that combines locally-developed packages with upstream OCI releases.

## Usage

```
hull releases install [flags]
```

## Flags

| Flag | Type | Default | Description |
|---|---|---|---|
| `--file` | string | "hull-releases.yaml" | spec file path |
| `-h, --help` | — | — | help for install |

## Persistent flags inherited from `hull`

| Flag | Type | Description |
|---|---|---|
| `--debug` | — | enable debug output |
| `--kube-context` | string | Kubernetes context to use |
| `--kubeconfig` | string | path to kubeconfig file |
| `-n, --namespace` | string | Kubernetes namespace |

## Examples

Install every release in `./hull-releases.yaml`:

```sh
hull releases install
```

Install from a custom-named file:

```sh
hull releases install --file ./platform.releases.yaml
```

## See also

- [`releases`](releases.md)
