---
title: "Sign and verify packages"
nav_order: 11
parent: "Guides"
---
{% raw %}
# Sign and verify packages

Hull signs packages two ways: a detached PGP `.prov` provenance file that hull
produces and verifies natively, and cosign signatures attached to an OCI
artifact, which hull can **verify** on pull (signing is done with the cosign
CLI). Verification can gate installs so the cluster never receives a package
whose provenance was not checked.

## Why sign

- **Tamper detection** — a modified archive fails verification.
- **Authorship attestation** — a signature binds the archive to a key.
- **Supply-chain audit** — `hull audit` records who installed each revision,
  when, and with which flags.

## PGP provenance (`.prov`)

### Sign at package time

```sh
hull package ./my-app -d ./build --sign --key author@example.com \
  --keyring ~/.gnupg/secring.gpg
```

```
Successfully packaged to: ./build/my-app-1.2.3.hull.tgz
Signed: ./build/my-app-1.2.3.hull.tgz.prov
```

This writes the archive plus a detached `.prov`. The `.prov` embeds a manifest
(package name, version, archive SHA-256) signed by the named key; verification
re-derives the archive's SHA-256 and compares it. Point `--key` at a private
key file or a signer name, and `--keyring` at the PGP secret keyring holding
it. If the key is passphrase-protected, add `--passphrase-file`.

### Sign an already-built archive

To sign an archive built without `--sign` (or after a key rotation), use the
`package sign` subcommand — it takes a private-key **file**:

```sh
hull package sign ./build/my-app-1.2.3.hull.tgz --key ./signing-key.asc
```

```
Successfully signed: ./build/my-app-1.2.3.hull.tgz.prov
```

### Distribute the public key

Consumers install the signer's public key into their hull keyring:

```sh
hull keyring add ./author.pub
```

```
Installed key author.pub (3AA5C34371567BD2)
```

```sh
hull keyring list
```

```
FINGERPRINT                                  FILE
3AA5C34371567BD2                             author.pub
```

Remove a signer by fingerprint or filename with `hull keyring remove`. The
keyring is a directory of public keys at `~/.config/hull/keyring/` (override the
config root with `HULL_CONFIG_HOME`).

### Verify

`hull package verify` checks an archive's `.prov` against a public key or
keyring. It exits 0 and prints nothing on success, non-zero with a precise
reason on any mismatch — so it chains cleanly as an install gate:

```sh
hull package verify ./my-app-1.2.3.hull.tgz --keyring ./author.pub && \
  hull install my-app ./my-app-dir -n prod
```

The `--verify` flag on `hull pull` does the same inline: it fetches the `.prov`
alongside the archive and validates it against your hull keyring before the
archive is kept.

```sh
hull pull my-app --repo https://charts.example.com --version 1.2.3 \
  --prov --verify --destination ./pulled
hull install my-app ./pulled/my-app -n prod
```

`hull install --verify` validates a local package's signature before applying;
`--keyring` selects the keyring directory (default `~/.config/hull/keyring`).

## Cosign (OCI artifacts)

Cosign signatures live in the OCI registry as a sibling artifact. Hull
**verifies** them on pull but does not create them — sign with the cosign CLI.

### Sign

```sh
hull registry push ./build/my-app-1.2.3.hull.tgz oci://ghcr.io/example/charts/my-app:1.2.3
cosign sign --key cosign.key ghcr.io/example/charts/my-app:1.2.3
# or keyless, with OIDC:
cosign sign ghcr.io/example/charts/my-app:1.2.3
```

### Verify on pull

`hull registry pull` verifies the cosign signature before the artifact touches
disk — fail-closed. Supply a key, or a keyless identity + issuer:

```sh
hull registry pull oci://ghcr.io/example/charts/my-app:1.2.3 \
  --cosign-key cosign.pub -d ./pulled
```

```
cosign signature verified for oci://ghcr.io/example/charts/my-app:1.2.3
Pulled oci://ghcr.io/example/charts/my-app:1.2.3 to ./pulled/my-app-1.2.3.hull.tgz
```

```sh
hull registry pull oci://ghcr.io/example/charts/my-app:1.2.3 \
  --cosign-identity user@example.com \
  --cosign-issuer   https://accounts.google.com \
  -d ./pulled
```

Without the cosign flags the pull runs unverified. In CI, always pass them so
an unsigned artifact fails the pipeline before install.

## Audit trail

Every install, upgrade, and rollback is recorded in the release history.
`hull audit` prints who did what, when:

```sh
hull audit my-app -n prod
```

```
REVISION   ACTION     USER          STATUS       TIMESTAMP
1          install    alice@corp    superseded   2026-07-10 09:14:02
2          upgrade    bob@corp      superseded   2026-07-12 16:40:55
3          rollback   alice@corp    deployed     2026-07-15 11:22:07
```

`--output json` (or `yaml`) adds the full provenance of each revision — the
recorded flags, value files, kubeconfig context, and hostname:

```sh
hull audit my-app --revision 2 -n prod --output json
```

The record persists for the life of the release, giving a change-management and
SLSA/SOC 2 trail per release.

## Patterns

### Key-per-environment in CI

Sign in CI with an environment-specific key; give each operator only the
public keys they should trust:

```sh
# CI build
hull package ./my-app -d ./build --sign --key ci-${ENV}@example.com --keyring ${KEYRING}

# operator side — trust only the dev signer
hull keyring add /etc/hull/dev-pubkeys.asc
hull pull my-app --repo https://charts.internal --version 1.2.3 \
  --prov --verify --destination ./pulled     # rejects archives signed by other envs
hull install my-app ./pulled/my-app
```

### Supply-chain artifacts

Generate a CycloneDX 1.5 SBOM for a deployed release — its package plus every
container image it runs — for `cosign attest`, Grype, Trivy, or Dependency
Track:

```sh
hull sbom my-app -n prod > my-app.cdx.json
```

## Common errors

- **`unknown signer` / verification failed on the key** — add the signer's
  public key: `hull keyring add ./key.pub`.
- **archive digest does not match the signed digest** — the archive on disk
  differs from the one that was signed (usually re-packaged without re-signing);
  re-fetch the original, or re-sign with `hull package sign`.
- **`provenance file (.prov) not found`** — the source ships no `.prov`; either
  drop `--verify` or use a source that signs.
- **cosign `no matching signatures`** — no cosign signature is attached to that
  artifact; check with `cosign tree <ref>`.

## See also

- [`hull package`](../cli/package.md) ·
  [`hull package sign`](../cli/package-sign.md) ·
  [`hull package verify`](../cli/package-verify.md)
- [`hull keyring`](../cli/keyring.md) — manage trusted public keys
- [`hull registry pull`](../cli/registry-pull.md) — cosign verification on pull
- [`hull audit`](../cli/audit.md) · [`hull sbom`](../cli/sbom.md)
- [OCI](oci.md) · [Repositories](repositories.md)
{% endraw %}
