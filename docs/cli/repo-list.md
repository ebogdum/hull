# hull repo list

## Synopsis

`hull repo list` prints every HTTP package repository currently registered with hull on this machine: name, URL, and a brief flag indicating whether credentials or TLS material are stored. The data comes from `~/.config/hull/repositories.yaml`.

## When to use it

Use to inventory configured repos, find a repo's URL for `hull pull --repo`, or verify a `hull repo add` succeeded.

## What happens when you run it

1. Reads `~/.config/hull/repositories.yaml`.
2. Prints in the requested output format (table by default).
3. No cluster contact, no network.

## Usage

```
hull repo list [flags]
```

## Flags

| Flag | Type | Default | Description |
|---|---|---|---|
| `-h, --help` | bool | false | help for list |
| `-o, --output` | string | table | output format: table, json, yaml |

## Persistent flags inherited from `hull`

| Flag | Type | Description |
|---|---|---|
| `--debug` | bool | enable debug output |
| `--kube-context` | string | Kubernetes context to use |
| `--kubeconfig` | string | path to kubeconfig file |
| `-n, --namespace` | string | Kubernetes namespace |

## Examples

Default tabular view:

```sh
hull repo list
```

JSON for scripting:

```sh
hull repo list -o json | jq '.[] | select(.name == "my-charts") | .url'
```

YAML for diffing across machines:

```sh
hull repo list -o yaml > /tmp/repos-machine-A.yaml
```

## See also

- [`repo`](repo.md)
- [`repo add`](repo-add.md)
- [`repo update`](repo-update.md)
- [`repo remove`](repo-remove.md)
