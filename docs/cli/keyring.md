# hull keyring

## Synopsis

`hull keyring` manages the PGP keyring used to verify package provenance signatures. Subcommands add, list, and remove armoured public keys.

## When to use it

Use to maintain the set of trusted signers for `--verify` operations.

## Usage

```
hull keyring [command]
```

## Subcommands

- [`hull keyring list`](keyring-list.md) — List keys in the hull keyring
- [`hull keyring remove`](keyring-remove.md) — Remove a key from the hull keyring

## Flags

| Flag | Type | Default | Description |
|---|---|---|---|
| `-h, --help` | — | — | help for keyring |

## Persistent flags inherited from `hull`

| Flag | Type | Description |
|---|---|---|
| `--debug` | — | enable debug output |
| `--kube-context` | string | Kubernetes context to use |
| `--kubeconfig` | string | path to kubeconfig file |
| `-n, --namespace` | string | Kubernetes namespace |

## Examples

Add a signer's public key:

```sh
hull keyring add /path/to/signer.pub
```

List trusted signers:

```sh
hull keyring list
```

Remove a signer by fingerprint:

```sh
hull keyring remove ABCDEF1234567890
```

## See also

- [Signing guide](../guides/signing.md)
- [`pull`](pull.md)
