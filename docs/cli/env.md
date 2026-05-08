# hull env

## Synopsis

`hull env` prints hull's environment information: the resolved config dir, cache dir, data dir, namespace, kubeconfig path, kubeconfig context, and any overrides supplied via environment variables.

## When to use it

Use when debugging path or auth issues — when you're not sure whether `~/.config/hull` is being used, when `KUBECONFIG` isn't being read as expected, etc.

## Usage

```
hull env [flags]
```

## Flags

| Flag | Type | Default | Description |
|---|---|---|---|
| `-h, --help` | — | — | help for env |

## Persistent flags inherited from `hull`

| Flag | Type | Description |
|---|---|---|
| `--debug` | — | enable debug output |
| `--kube-context` | string | Kubernetes context to use |
| `--kubeconfig` | string | path to kubeconfig file |
| `-n, --namespace` | string | Kubernetes namespace |

## Examples

Print hull's environment:

```sh
hull env
```

## See also

- [CLI reference index](README.md)
