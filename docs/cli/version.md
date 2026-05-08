# hull version

## Synopsis

`hull version` prints the hull binary version, the Git commit it was built from, and the build date.

## When to use it

Use to confirm which version is installed.

## Usage

```
hull version [flags]
```

## Flags

| Flag | Type | Default | Description |
|---|---|---|---|
| `-h, --help` | — | — | help for version |

## Persistent flags inherited from `hull`

| Flag | Type | Description |
|---|---|---|
| `--debug` | — | enable debug output |
| `--kube-context` | string | Kubernetes context to use |
| `--kubeconfig` | string | path to kubeconfig file |
| `-n, --namespace` | string | Kubernetes namespace |

## Examples

Print version:

```sh
hull version
```
