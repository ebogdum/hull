# hull registry

## Synopsis

`hull registry` is the OCI counterpart of `hull repo`. Subcommands log in/out of registries, push and pull artifacts, and list available tags. Distinct from `hull repo` which targets HTTP repositories.

## When to use it

Use when distributing packages via OCI registries.

## Usage

```
hull registry [command]
```

## Subcommands

- [`hull registry push`](registry-push.md) — push a `.hull.tgz` archive to an OCI registry
- [`hull registry pull`](registry-pull.md) — pull a package from an OCI registry

## Flags

| Flag | Type | Default | Description |
|---|---|---|---|
| `-h, --help` | — | — | help for registry |

## Persistent flags inherited from `hull`

| Flag | Type | Description |
|---|---|---|
| `--debug` | — | enable debug output |
| `--kube-context` | string | Kubernetes context to use |
| `--kubeconfig` | string | path to kubeconfig file |
| `-n, --namespace` | string | Kubernetes namespace |

## Examples

Push a packaged archive to OCI:

```sh
hull registry push ./build/my-app-1.0.0.hull.tgz oci://ghcr.io/example/charts/my-app
```

Pull a package from OCI:

```sh
hull registry pull oci://ghcr.io/example/charts/my-app:1.0.0 -d ./pulled
```

## See also

- [`pull`](pull.md)
- [`registry push`](registry-push.md)
- [`login`](login.md)
- [OCI guide](../guides/oci.md)
