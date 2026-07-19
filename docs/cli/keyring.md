---
title: "hull keyring"
parent: "CLI"
---
{% raw %}
# hull keyring

## Synopsis

`hull keyring` manages the set of trusted PGP public keys hull uses to verify
package provenance. When you run a signed operation with `--verify` (for
example `hull package verify` or `hull install --verify`), hull checks the
package's `.prov` signature against the keys in this keyring.

The keyring is a per-user, per-machine directory, `~/.config/hull/keyring/`.
There is no shared cluster-wide keyring.

## Subcommands

| Command | What it does |
|---|---|
| [`add`](keyring-add.md) | Install a public key so its signer becomes trusted. |
| [`list`](keyring-list.md) | Show the installed keys and their fingerprints. |
| [`remove`](keyring-remove.md) | Delete a key, revoking trust in its signer. |

## Usage

```
hull keyring [command]
```

Add a signer, confirm it landed, and later remove it:

```sh
hull keyring add ./jane.pub
hull keyring list
hull keyring remove jane.pub
```

## See also

- [`keyring add`](keyring-add.md) — trust a signer
- [`keyring list`](keyring-list.md) — show trusted signers
- [`keyring remove`](keyring-remove.md) — revoke a signer
- [`package verify`](package-verify.md) — verify a package against the keyring
- [`login`](login.md) — store registry credentials
{% endraw %}
