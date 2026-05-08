# OCI

Hull packages can be pushed to and pulled from any OCI distribution-spec registry — GHCR, Docker Hub, Quay, Harbor, Artifactory, ECR, GAR, ACR, self-hosted Zot, self-hosted Distribution. Each version of a package becomes an OCI artifact; hull stores the `*.hull.tgz` archive as the artifact's blob and a thin manifest references it.

OCI is the recommended distribution mechanism when:

- You already have a container registry and want one place for both images and packages.
- You want IAM/auth integration with cloud providers.
- You want immutable, content-addressable artifacts.

For static HTTP repository hosting (GitHub Pages, S3 + index), see [Repositories](repositories.md). The two are complementary.

## Authentication

Hull uses the same credential layout as the registry's CLI tools. `hull registry login` writes to `~/.config/hull/credentials.json`:

```sh
hull registry login ghcr.io
hull registry login ghcr.io --username myuser --password ${GITHUB_TOKEN}
hull registry login registry.example.com --username u --password p
```

Subsequent push, pull, and OCI-source layer fetches authenticate transparently.

To remove credentials:

```sh
hull registry logout ghcr.io
```

Hull also reads `~/.docker/config.json` as a fallback. Once you've run a registry's standard login (`docker login`, `aws ecr get-login-password | docker login ...`, `gcloud auth configure-docker`), hull pushes and pulls just work.

## Pushing a package

Two equivalent commands push a packaged archive to an OCI registry:

```sh
# Build the archive
hull package ./my-app -d ./build

# Push via the dedicated subcommand
hull registry push ./build/my-app-1.2.3.hull.tgz oci://ghcr.io/example/charts/my-app

# OR via the multi-target publish command
hull publish ./build/my-app-1.2.3.hull.tgz --oci oci://ghcr.io/example/charts/my-app
```

Hull pushes the archive as a single blob with media type `application/vnd.hull.package.v1.tar+gzip` and creates a manifest tagged with the package's version (`1.2.3` from `hull.yaml`).

## Pulling

```sh
# Tag included in the URI
hull registry pull oci://ghcr.io/example/charts/my-app:1.2.3 -d ./pulled

# Or use the unified pull command, which accepts a SemVer constraint:
hull pull oci://ghcr.io/example/charts/my-app --version "^1.2.0" -d ./pulled

# Untar after download:
hull registry pull oci://ghcr.io/example/charts/my-app:1.2.3 -d ./pulled
hull pull   oci://ghcr.io/example/charts/my-app --version 1.2.3 -d ./pulled --untar
```

`hull pull --version` accepts a SemVer constraint and resolves it against the registry's tag list. `hull registry pull` always uses the literal tag in the URI.

## Installing directly from OCI

Pin the version via the URI tag:

```sh
hull install my-app oci://ghcr.io/example/charts/my-app:1.2.3 -n prod
```

Hull pulls, untars to a temp directory, and installs.

## Layers from OCI

In `hull.yaml`:

```yaml
layers:
  - name: my-layer
    source: oci://ghcr.io/example/charts/my-layer
    version: ^2.0.0
```

`hull dependency update .` resolves the constraint against the registry, locks the digest, and stores it in `hull.lock`:

```yaml
layers:
  - name: my-layer
    source: oci://ghcr.io/example/charts/my-layer
    resolvedVersion: 2.4.1
    digest: sha256:7e2a90d8c5...
```

## Insecure / plaintext registries

For local development against a registry without TLS, use `--plain-http`:

```sh
hull registry push  ./build/my-app-1.2.3.hull.tgz oci://localhost:5000/charts/my-app --plain-http
hull registry pull  oci://localhost:5000/charts/my-app:1.2.3 --plain-http
```

For self-signed TLS, use `--insecure-skip-tls-verify`:

```sh
hull registry push ./build/my-app-1.2.3.hull.tgz oci://registry.local/charts/my-app --insecure-skip-tls-verify
hull registry pull oci://registry.local/charts/my-app:1.2.3 --insecure-skip-tls-verify
```

The same effect can be set globally via environment variables (useful in CI):

| Env var | Effect |
|---|---|
| `HULL_OCI_PLAIN_HTTP=true` | use HTTP instead of HTTPS for OCI |
| `HULL_OCI_INSECURE_SKIP_TLS=true` | skip TLS verification for OCI |

## Provenance verification on install

Hull verifies the package's PGP `.prov` signature when `--verify` is passed:

```sh
hull install my-app oci://ghcr.io/example/charts/my-app:1.2.3 -n prod --verify
```

The PGP keyring lives at `~/.config/hull/keyring/` (configurable via `--keyring`). See [Signing](signing.md) for the full keyring management story.

For cosign-style OCI-attached signatures, run cosign verification *before* invoking hull (or as part of your CI's pull step). For example:

```sh
cosign verify --key cosign.pub ghcr.io/example/charts/my-app:1.2.3
hull install my-app oci://ghcr.io/example/charts/my-app:1.2.3 -n prod
```

This pattern keeps cosign as the externally-managed signature gate and hull as the package installer.

## Troubleshooting

- **`unauthorized: authentication required`** — log in: `hull registry login <host>`.
- **`x509: certificate signed by unknown authority`** — for an internal CA, set `HULL_OCI_CA_FILE=/path/to/ca.pem` (when supported) or use `--insecure-skip-tls-verify` for a quick test (do not leave that on in production).
- **`pull: no matching version for constraint X`** — list the tags via the registry's own API to see what's actually published.
- **Slow pulls from public registries** — hull caches blobs by digest in `~/.cache/hull/oci/` (or `${HULL_CACHE_HOME}/oci/`); a re-pull of the same digest is instant.
