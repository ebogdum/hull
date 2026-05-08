# hull package

## Synopsis

`hull package` creates a `.hull.tgz` archive from a package directory. The archive is self-contained: every layer is materialised, `hull.lock` is included, and (with `--sign`) a detached PGP `.prov` file is produced alongside.

## When to use it

Use as the final step before publication. The output archive can be uploaded to an HTTP repository (`hull publish --repo`), pushed to OCI (`hull registry push` / `hull publish --oci`), or distributed by hand.

## What happens when you run it

1. Reads `<path>` and resolves every layer (uses cached materials from `hull dependency build` if present).
2. Validates `hull.yaml` and `values.yaml`.
3. Composes a tarball containing the package directory, layer cache, and `hull.lock`.
4. Names the output `<destination>/<name>-<version>.hull.tgz`.
5. With `--sign --key <path>`, also emits a `.prov` file (detached PGP signature) next to the archive.
6. With `--reproducible`, writes the archive deterministically (zero timestamps, canonical modes) for reproducible builds.

## Usage

```
hull package <path> [flags]
hull package [command]
```

## Subcommands

- [`hull package sign`](package-sign.md) — sign an existing archive with a PGP private key
- [`hull package verify`](package-verify.md) — verify a `.prov` signature against the local keyring

## Flags

| Flag | Type | Default | Description |
|---|---|---|---|
| `--app-version` | string | — | override the appVersion in hull.yaml |
| `-d, --destination` | string | "." | directory to write the archive to |
| `-h, --help` | — | — | help for package |
| `--key` | string | — | PGP private key file or signer name (used with --sign) |
| `--keyring` | string | — | PGP keyring file containing the signer (alternative to --key) |
| `--passphrase-file` | string | — | file containing the key's passphrase (- for stdin) |
| `--reproducible` | — | — | produce byte-identical output across machines (zero timestamps, canonical modes) |
| `--sign` | — | — | produce a .prov provenance file alongside the archive (requires --key or --keyring) |
| `--version` | string | — | override the version in hull.yaml |

## Persistent flags inherited from `hull`

| Flag | Type | Description |
|---|---|---|
| `--debug` | — | enable debug output |
| `--kube-context` | string | Kubernetes context to use |
| `--kubeconfig` | string | path to kubeconfig file |
| `-n, --namespace` | string | Kubernetes namespace |

## Examples

Package a directory:

```sh
hull package ./my-app -d ./build
```

Package and sign in one shot:

```sh
hull package ./my-app -d ./build --sign --key /path/to/secret-key.asc
```

Reproducible build (deterministic output across machines):

```sh
hull package ./my-app -d ./build --reproducible
```

Override the version captured in the archive (e.g. for a CI release-candidate tag):

```sh
hull package ./my-app -d ./build --version 1.2.3-rc.4 --app-version 1.5.0
```

Verify a previously-signed archive:

```sh
hull package verify ./build/my-app-1.0.0.hull.tgz
```

## See also

- [`publish`](publish.md)
- [`registry push`](registry-push.md)
- [Repositories guide](../guides/repositories.md)
- [Signing guide](../guides/signing.md)
