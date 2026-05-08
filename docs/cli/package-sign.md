# hull package sign

## Synopsis

`hull package sign` produces a detached PGP `.prov` (provenance) signature alongside an existing `.hull.tgz` archive. The `.prov` file is a cleartext-signed envelope that includes the archive's SHA-256 digest, package name, and version. Consumers verify it with `hull package verify` or with `--verify` on `hull pull` / `hull install`. The `--key` flag points at the **private key file** — typically an exported PGP secret key in armoured form.

## When to use it

Use when packaging happened without `--sign` (e.g. an archive built upstream that you want to re-sign before redistribution) or when re-signing after a key rotation. For new archives, prefer `hull package <pkg-dir> --sign --key <path>` which packages and signs in one shot.

## What happens when you run it

1. Reads the archive at `<archive.hull.tgz>` and computes its SHA-256.
2. Loads the PGP private key at `--key`.
3. Constructs the provenance manifest (name, version, digest, signer fingerprint).
4. PGP-signs the manifest and writes `<archive.hull.tgz>.prov` next to the archive.
5. No cluster contact, no network.

## Usage

```
hull package sign <archive.hull.tgz> [flags]
```

## Flags

| Flag | Type | Default | Description |
|---|---|---|---|
| `-h, --help` | bool | false | help for sign |
| `--key` | string | "" | path to PGP private key file (required) |

## Persistent flags inherited from `hull`

| Flag | Type | Description |
|---|---|---|
| `--debug` | bool | enable debug output |
| `--kube-context` | string | Kubernetes context to use |
| `--kubeconfig` | string | path to kubeconfig file |
| `-n, --namespace` | string | Kubernetes namespace |

## Examples

Sign an archive with an exported secret key file:

```sh
hull package sign ./build/my-app-1.0.0.hull.tgz --key /path/to/secret-key.asc
```

Result is `./build/my-app-1.0.0.hull.tgz.prov` next to the archive.

Re-sign an archive after key rotation (overwrites the existing `.prov`):

```sh
hull package sign ./build/my-app-1.0.0.hull.tgz --key /path/to/new-key.asc
```

Sign and immediately verify in one shell pipeline:

```sh
hull package sign   ./build/my-app-1.0.0.hull.tgz --key /path/to/key.asc
hull package verify ./build/my-app-1.0.0.hull.tgz --keyring /path/to/pubkey.asc
```

## See also

- [`package`](package.md) — package and sign in one command
- [`package verify`](package-verify.md)
- [`keyring add`](keyring-add.md)
- [Signing guide](../guides/signing.md)
