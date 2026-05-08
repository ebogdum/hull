# hull registry push

## Synopsis

`hull registry push` uploads a `.hull.tgz` archive to an OCI distribution-spec registry. The destination reference identifies the artifact path (`oci://host/path/<name>`); hull tags the artifact with the package's `version` from `hull.yaml` (extracted from the archive). The result is an immutable, content-addressable artifact in the registry, suitable for `hull install oci://...:<version>` consumption.

## When to use it

Use as the publication step for OCI-based distribution. For multi-target publication (HTTP repo and OCI in one shot), use `hull publish`.

## What happens when you run it

1. Reads `<archive>` and extracts its `hull.yaml` to determine name and version.
2. Resolves credentials for `<ref>`'s host (from `~/.config/hull/credentials.json` or `~/.docker/config.json`).
3. Computes the blob digest, checks if it already exists in the registry (deduplicates).
4. Pushes the blob and a thin manifest tagged with the package's version.
5. Prints the resolved tag and digest on success.

## Usage

```
hull registry push <archive> <ref> [flags]
```

## Flags

| Flag | Type | Default | Description |
|---|---|---|---|
| `-h, --help` | bool | false | help for push |
| `--insecure-skip-tls-verify` | bool | false | skip TLS certificate verification |
| `--plain-http` | bool | false | use plaintext HTTP (no TLS) |

## Persistent flags inherited from `hull`

| Flag | Type | Description |
|---|---|---|
| `--debug` | bool | enable debug output |
| `--kube-context` | string | Kubernetes context to use |
| `--kubeconfig` | string | path to kubeconfig file |
| `-n, --namespace` | string | Kubernetes namespace |

## Examples

Push a packaged archive to a public registry:

```sh
hull registry push ./build/my-app-1.0.0.hull.tgz oci://ghcr.io/example/charts/my-app
```

Push to an internal registry over plain HTTP (development cluster):

```sh
hull registry push ./build/my-app-1.0.0.hull.tgz oci://localhost:5000/charts/my-app --plain-http
```

Push to a registry with a self-signed TLS certificate:

```sh
hull registry push ./build/my-app-1.0.0.hull.tgz oci://registry.local/charts/my-app --insecure-skip-tls-verify
```

Authenticate first, then push:

```sh
hull registry login ghcr.io --username "$GITHUB_USER" --password "$GITHUB_TOKEN"
hull registry push ./build/my-app-1.0.0.hull.tgz oci://ghcr.io/example/charts/my-app
```

## See also

- [`registry`](registry.md)
- [`registry pull`](registry-pull.md)
- [`publish`](publish.md) — multi-target publication
- [`package`](package.md) — produce the archive first
- [OCI guide](../guides/oci.md)
