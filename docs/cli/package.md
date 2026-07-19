# hull package

Package a hull package directory into a versioned `.hull.tgz` archive.

## When to use it

- To turn a working package directory into a single distributable file before
  you [`publish`](publish.md) it or hand it off.
- To sign at build time (`--sign`), or to produce byte-identical output for
  reproducible builds (`--reproducible`).

## What happens

1. hull reads `<path>`, validates the package, and writes an archive named
   `<name>-<version>.hull.tgz` into the `--destination` directory (default the
   current directory).
2. It prints `Successfully packaged to: <archive-path>`.
3. With `--sign`, it then signs the archive with the key from `--key` (or
   `--keyring`), writing a detached `<archive>.prov` provenance file, and
   prints `Signed: <prov-path>`. If the key is passphrase-protected, supply
   `--passphrase-file`.
4. `--version` / `--app-version` override the values recorded from `hull.yaml`;
   `--reproducible` zeroes timestamps and canonicalises file modes so the same
   inputs always yield the same bytes.

## Usage

```
hull package <path> [flags]
hull package [command]
```

## Subcommands

| Command | What it does |
|---|---|
| [`package sign`](package-sign.md) | sign an existing archive with a PGP private key |
| [`package verify`](package-verify.md) | verify an archive's `.prov` signature against a key |

## Flags

| Flag | Type | Default | Description |
|---|---|---|---|
| `-d, --destination` | string | "." | directory to write the archive to |
| `--version` | string | — | override the version recorded from hull.yaml |
| `--app-version` | string | — | override the appVersion recorded from hull.yaml |
| `--reproducible` | — | — | produce byte-identical output across machines (zero timestamps, canonical modes) |
| `--sign` | — | — | also produce a `.prov` provenance file (requires `--key` or `--keyring`) |
| `--key` | string | — | PGP private key file or signer name, used with `--sign` |
| `--keyring` | string | — | PGP keyring file containing the signer (alternative to `--key`) |
| `--passphrase-file` | string | — | file holding the key's passphrase (`-` for stdin) |

## Worked example

Package a directory and sign it in one step:

```sh
hull package ./my-app -d ./build --sign --key ./cosign.key
```

```
Successfully packaged to: ./build/my-app-1.0.0.hull.tgz
Signed: ./build/my-app-1.0.0.hull.tgz.prov
```

Now the archive and its provenance file are ready to
[`publish`](publish.md).

## See also

- [`package sign`](package-sign.md) — sign an already-built archive
- [`package verify`](package-verify.md) — check a signature
- [`publish`](publish.md) — upload the archive
- [`pull`](pull.md) · [`install`](install.md)
