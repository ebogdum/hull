# hull releases upgrade

## Synopsis

Upgrade every release in `hull-releases.yaml`; install if missing. One invocation brings the platform graph up the first time and keeps it up-to-date thereafter.

## When to use it

Use as the canonical CI deploy command for the whole platform graph.

## Usage

```
hull releases upgrade [flags]
```

## Flags

| Flag | Type | Default | Description |
|---|---|---|---|
| `--file` | string | "hull-releases.yaml" | spec file path |
| `-h, --help` | — | — | help for upgrade |

## Persistent flags inherited from `hull`

| Flag | Type | Description |
|---|---|---|
| `--debug` | — | enable debug output |
| `--kube-context` | string | Kubernetes context to use |
| `--kubeconfig` | string | path to kubeconfig file |
| `-n, --namespace` | string | Kubernetes namespace |

## Examples

Bring the platform up to declared versions:

```sh
hull releases upgrade
```

Use a custom-named manifest file:

```sh
hull releases upgrade --file ./platform.releases.yaml
```

## See also

- [`releases`](releases.md)
- [`upgrade`](upgrade.md)
