---
title: "Crypto and Secrets functions"
parent: "Functions"
grand_parent: "Templates"
---
{% raw %}
# Crypto and Secrets functions

> **Pipeline note.** `${value | f x}` = `f(value, x)`. Non-string inputs are stringified with `%v` for hashes; `rand*` length functions parse the value as a number.

## Crypto functions

### `sha1sum`
`sha1sum(value)` → string — SHA-1 hex (40 chars). Checksums only, not security.
```
${"hello" | sha1sum}   → "aaf4c61ddcc5e8a2dabede0f3b482cd9aea9434d"
```

### `sha256sum`
`sha256sum(value)` → string — SHA-256 hex (64 chars).
```
${"hello" | sha256sum} → "2cf24dba5fb0a30e26e83b2ac5b9e29e1b161e5c1fa7425e73043362938b9824"
```

### `sha512sum`
`sha512sum(value)` → string — SHA-512 hex (128 chars).

### `md5sum`
`md5sum(value)` → string — MD5 hex (32 chars). Broken; checksums only.
```
${"hello" | md5sum}    → "5d41402abc4b2a76b9719d911017c592"
```

### `adler32sum`
`adler32sum(value)` → string — Adler-32 checksum as a decimal string.
```
${"hello" | adler32sum} → "103547413"
```

### `hmacSha256`
`hmacSha256(value, key)` → string — HMAC-SHA256 hex. Errors `hmacSha256 requires a key argument` without a key.
```
${"message" | hmacSha256 "key"} → "6e9ef29b75fffc5b7abae527d58fdadb2fe42e7219011976917343065f58ed4a"
```

### `bcrypt`
`bcrypt(value, cost?)` → string — bcrypt hash (`$2a$…`, 60 chars). Default cost 10; an optional cost is applied only if it parses as an int in 4–31, else the default is kept. Salted → value varies.
```
${"password" | bcrypt}    → "$2a$10$N9qo8uLOickgx2ZMRZoMy..." (shape; varies)
${"password" | bcrypt 12} → "$2a$12$..." (shape; varies)
```

### `htpasswd`
`htpasswd(user, password)` → string — an Apache `user:bcryptHash` line (cost 10). The **value** is the username, the **arg** the password. Errors `htpasswd requires a password argument`.
```
${"admin" | htpasswd "s3cret"} → "admin:$2a$10$..." (shape; varies)
```

### `encryptAES`
`encryptAES(value, passphrase)` → string — AES-256-GCM; key via PBKDF2-HMAC-SHA256 (100k iters, random 16-byte salt). Output `base64(salt||nonce||ciphertext+tag)`. Non-deterministic. Errors `encryptAES requires a passphrase argument`.
```
${"topsecret" | encryptAES "hunter2"} → "b3NhbHQ...=" (base64; varies each call)
```

### `decryptAES`
`decryptAES(value, passphrase)` → string — inverse of `encryptAES`. Errors on bad base64, too-short blob, or `decryptAES: open failed` (wrong passphrase / tampering).
```
${cipher | decryptAES "hunter2"} → "topsecret"
```

### `genPrivateKey`
`genPrivateKey(kind, bits?)` → string (PEM) — key type from `value`:
- `""`/`"rsa"` → RSA (default 2048; optional `bits` must be **2048–8192**, else error; non-int ignored), PEM `RSA PRIVATE KEY`.
- `"ecdsa"` → ECDSA P-256, PEM `EC PRIVATE KEY`.
- `"ed25519"` → Ed25519, PKCS#8, PEM `PRIVATE KEY`.

Errors: `RSA key size N is below the 2048-bit minimum` / `exceeds the 8192-bit maximum` / `unsupported type "…"`.
```
${"rsa" | genPrivateKey}
# → (multi-line PEM; shape, value varies)
# -----BEGIN RSA PRIVATE KEY-----
# MIIEowIBAAKCAQEAn9xH4zxZRB7o+VGD0FHjIXoO2D3bAsgMt9FpYHjLJ1UEaG+x
# …
# -----END RSA PRIVATE KEY-----
${"ed25519" | genPrivateKey}
# → (multi-line PEM; shape, value varies)
# -----BEGIN PRIVATE KEY-----
# MC4CAQAwBQYDK2VwBCIEIEIciMvRpAMmWfkmaq9opPSfpNEnSA2DbU0nQzoqhct+
# -----END PRIVATE KEY-----
${"rsa" | genPrivateKey 1024} → error (RSA key size 1024 is below the 2048-bit minimum)
```

### `genCA`
`genCA(commonName, days?)` → map`{Cert, Key}` — self-signed CA cert + RSA-2048 key. CN defaults `hull-ca`; `days` default 365. Returns PEM `Cert`/`Key`.
```
${"my-ca" | genCA 3650 | get "Cert"} → the CA certificate PEM (shape; varies)
```

### `genSelfSignedCert`
`genSelfSignedCert(commonName, ...sans)` → map`{Cert, Key}` — self-signed leaf cert + RSA-2048 key, 1-year validity. CN defaults `localhost` and is added as a DNS SAN; each extra arg is an IP SAN if IP-parseable, else a DNS SAN.
```
${"svc.example.com" | genSelfSignedCert "svc" "10.0.0.5" | get "Key"} → RSA private key PEM (shape; varies)
```

### `randAlphaNum` / `randAlpha` / `randNumeric` / `randAscii`
`randAlphaNum(length)` → string, and siblings — cryptographically random strings over `[a-zA-Z0-9]` / `[a-zA-Z]` / `[0-9]` / printable-ASCII. Length from `value` (0–65536). Errors `rand: invalid length …` / `rand: length N exceeds 65536`.
```
${16 | randAlphaNum} → "aZ3kP0qR9sT1uV2w" (shape; varies, len 16)
${6  | randNumeric}  → "402913" (shape; varies)
```

### `randBytes`
`randBytes(n)` → string — `n` random bytes, base64-encoded. `n` in 0–65536.
```
${32 | randBytes} → "u8Qw…=" (44-char base64 of 32 bytes; varies)
```

### `uuidv4`
`uuidv4(value)` → string — random RFC-4122 v4 UUID (value/args ignored).
```
${"" | uuidv4} → "3f2504e0-4f89-41d3-9a0c-0305e82c3301" (shape; varies)
```

## Secrets functions

`sops`/`sopsKey` shell out to the host `sops` binary; a missing binary returns a structured error (not a panic). Paths must be relative, must not begin with `-`, and must not escape the package dir. `sopsKey`'s key path must match `^[A-Za-z0-9_][A-Za-z0-9_-]*(\.[A-Za-z0-9_][A-Za-z0-9_-]*)*$`.

### `sops`
`sops(path)` → string — runs `sops --decrypt <path>`, returns the plaintext (one trailing newline trimmed).
```
${"secrets/db.enc.yaml" | sops}
# → (decrypted plaintext, trailing newline trimmed)
# username: admin
# password: hunter2
```

### `sopsKey`
`sopsKey(path, keyPath)` → string — decrypts then extracts a single dotted key. Errors `sopsKey requires a key path argument`.
```
${"secrets/db.enc.yaml" | sopsKey "database.password"} → "hunter2"
```

### `externalSecret`
`externalSecret(name, secretStore, remoteKey, refreshInterval?)` → string (YAML) — emits an `external-secrets.io/v1beta1` ExternalSecret. Name from `value` (or first arg if empty). Default refresh `1h`.

Emits:
```yaml
apiVersion: external-secrets.io/v1beta1
kind: ExternalSecret
metadata:
  name: <name>
spec:
  refreshInterval: <refresh|1h>
  secretStoreRef:
    name: <secretStore>
    kind: SecretStore
  target:
    name: <name>
    creationPolicy: Owner
  dataFrom:
    - extract:
        key: <remoteKey>
```
```
${"db-creds" | externalSecret "vault-backend" "secret/data/db"}
  → the ExternalSecret above (name=db-creds, store=vault-backend, key=secret/data/db)
```

### `sealedSecret`
`sealedSecret(name, namespace, encryptedData)` → string (YAML) — wraps a **precomputed** ciphertext map (from `kubeseal`) into a `bitnami.com/v1alpha1` SealedSecret. It does NOT encrypt. `encryptedData` must be a map.

Emits:
```yaml
apiVersion: bitnami.com/v1alpha1
kind: SealedSecret
metadata:
  name: <name>
  namespace: <namespace>
spec:
  encryptedData:
    <k>: <ciphertext>
  template:
    metadata:
      name: <name>
      namespace: <namespace>
```
```
${"api-key" | sealedSecret "prod" (dict "token" "AgB2f…")}
  → the SealedSecret above with the sealed token under spec.encryptedData
```
{% endraw %}
