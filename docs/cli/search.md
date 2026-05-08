# hull search

## Synopsis

`hull search` searches for packages. Subcommands search the public Artifact Hub catalogue (`hub`) or every registered repository (`repo`).

## When to use it

Use to discover available packages, both public and from your registered repos.

## Usage

```
hull search [command]
```

## Subcommands

- [`hull search repo`](search-repo.md) — Search configured repositories for packages

## Flags

| Flag | Type | Default | Description |
|---|---|---|---|
| `-h, --help` | — | — | help for search |

## Persistent flags inherited from `hull`

| Flag | Type | Description |
|---|---|---|
| `--debug` | — | enable debug output |
| `--kube-context` | string | Kubernetes context to use |
| `--kubeconfig` | string | path to kubeconfig file |
| `-n, --namespace` | string | Kubernetes namespace |

## Examples

Search the public hub:

```sh
hull search hub mqtt
```

Search registered repos:

```sh
hull search repo nginx
```

## See also

- [`repo`](repo.md)
- [`show`](show.md)
