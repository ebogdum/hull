---
title: "Host and consume an HTTP repository"
nav_order: 9
parent: "Guides"
---
{% raw %}
# Host and consume an HTTP repository

A hull repository is a directory served over HTTP(S) containing packaged
`*.hull.tgz` archives plus an `index.yaml` that catalogues them. It is the
simplest distribution mechanism — any static HTTP server hosts one, including
GitHub Pages, S3 with a static index, or an internal nginx.

For OCI-registry distribution, see [OCI](oci.md). The two are complementary; a
project often publishes both.

## The shape of a repo

```
https://charts.example.com/
├── index.yaml
├── my-app-1.0.0.hull.tgz
├── my-app-1.1.0.hull.tgz
├── my-app-1.2.0.hull.tgz
└── my-app-1.2.0.hull.tgz.prov      # optional detached PGP signature
```

`index.yaml` lists every package, its versions, each archive's URL, and each
archive's SHA-256 digest.

## Produce a repo

### Package each release

```sh
hull package ./my-app -d ./build
```

```
Successfully packaged to: ./build/my-app-1.2.3.hull.tgz
```

`hull package` reads the package's `hull.yaml`, resolves layers, and writes a
single self-contained `<name>-<version>.hull.tgz`. To sign at package time, add
`--sign` with a signer:

```sh
hull package ./my-app -d ./build --sign --key author@example.com \
  --keyring ~/.gnupg/secring.gpg
```

```
Successfully packaged to: ./build/my-app-1.2.3.hull.tgz
Signed: ./build/my-app-1.2.3.hull.tgz.prov
```

This writes the archive plus a detached `.prov` signature. See
[Signing](signing.md) for the full story.

### Generate the index

```sh
hull repo index ./build --url https://charts.example.com
```

```
Index generated at /home/you/build/index.yaml
```

`--url` is the base URL where the archives will be served; hull writes
per-version absolute download URLs into the index and records each archive's
SHA-256. To add a new version without regenerating from scratch, use `--merge`;
to sign the index, add `--sign <private-key>` (writes `index.yaml.prov`):

```sh
hull repo index ./build --url https://charts.example.com --merge --sign ./repo-key.asc
```

### Publish

The `./build/` directory is the entire repository — serve it from any static
HTTP host:

- **GitHub Pages** — commit `build/`, enable Pages, serve at
  `https://<user>.github.io/<repo>/`.
- **S3** — `aws s3 sync ./build s3://my-bucket/`, set `index.yaml`'s
  Content-Type to `application/yaml`.
- **nginx** — serve the directory; enable `autoindex` for a browsable view.

There is no special server; the repo is dumb static files.

## Consume a repo

Register the repo, refresh its index, then find and pull packages by name:

```sh
hull repo add my-charts https://charts.example.com
```

```
"my-charts" has been added to your repositories
```

```sh
hull repo update
```

```
...successfully got an update from "my-charts"
Update complete.
```

```sh
hull repo list
```

```
NAME                 URL
my-charts            https://charts.example.com
```

```sh
hull search repo my-app
```

`hull repo add` records the name and URL in
`~/.config/hull/repositories.yaml`; it does **not** fetch the index — that is
what `hull repo update` does, caching each repo's `index.yaml` under
`~/.cache/hull/indexes/` (30-minute TTL). Adding a name that already exists is
left untouched unless you pass `--force-update`.

### Pull a package

`hull pull` fetches a named chart from a repo (`--repo`), with SemVer version
selection, and can unpack it. Note the flag is `--destination` — `hull pull`
has no `-d` shorthand:

```sh
hull pull my-app --repo https://charts.example.com --version "^1.2.0" \
  --destination ./pulled --untar
```

```
Pulled and extracted: ./pulled/my-app
```

Then install from the unpacked **directory** (`hull install` takes a package
directory, not an archive or a URL):

```sh
hull install my-app ./pulled/my-app -n default --create-namespace
```

## Authenticate to a private repo

Supply credentials when you add the repo:

```sh
hull repo add private https://charts.example.com --username u --password p
```

For token-based auth (for example GitHub Pages behind a fine-grained token),
pass the token as the password — most providers accept it in place of a basic
password:

```sh
hull repo add private https://charts.example.com \
  --username "$GITHUB_USER" --password "$GITHUB_TOKEN"
```

Credentials are stored in `~/.config/hull/credentials.json`, keyed by host, and
reused automatically. To refresh a host's credential without re-adding the
repo, use `hull login` — it is non-interactive and needs a credential flag
(`-u/--username` with `-p/--password`, or `--token`, or `--api-key`):

```sh
echo "$GITHUB_TOKEN" | hull login charts.example.com -u "$GITHUB_USER" --password-stdin
```

```
Login succeeded for charts.example.com
```

`hull logout charts.example.com` removes only the stored credential, not the
repo registration. To unregister the repo itself, `hull repo remove private`.

## TLS options

`hull repo add`, `hull repo update`, and `hull pull` honour:

- `--ca-file <path>` — trust an extra CA bundle for a self-signed server.
- `--cert-file <path>` / `--key-file <path>` — client certificate for mTLS.
- `--insecure-skip-tls-verify` — skip server-cert validation (`hull repo add`
  only; do not use in production).

## Inspect a package

`hull show` operates on an unpacked package **directory**:

```sh
hull pull my-app --repo https://charts.example.com --version 1.2.3 \
  --destination ./pulled --untar
hull show chart   ./pulled/my-app     # hull.yaml metadata
hull show values  ./pulled/my-app     # default values.yaml
hull show readme  ./pulled/my-app     # README.md
hull show crds    ./pulled/my-app     # declared CRDs
hull show all     ./pulled/my-app     # everything in one document
```

To browse Artifact Hub without registering a repo:

```sh
hull search hub my-app
hull search hub my-app --endpoint https://artifacthub.io
```

`hull search hub` queries the public Artifact Hub API directly.

## Verify provenance on pull

When an archive ships a `.prov` sidecar, verify the signature as you pull:

```sh
hull pull my-app --repo https://charts.example.com --version 1.2.3 \
  --prov --verify --destination ./pulled
```

`--verify` fetches the `.prov` and checks it against your local keyring before
the archive is kept; a bad signature aborts the pull. See [Signing](signing.md).

## Troubleshooting

- **`digest mismatch`** — the archive on the server does not match the digest
  recorded in `index.yaml`. Common cause: an archive was re-published without
  regenerating the index (`hull repo index ./build --url ...`).
- **`401 Unauthorized`** — no credential was sent, or it is wrong. Re-run
  `hull login <host>`; check `~/.config/hull/credentials.json`.
- **index not found** — the URL does not serve `index.yaml` at its root. Verify
  with `curl <url>/index.yaml`.
- **TLS handshake error** — the server certificate is not trusted; pass
  `--ca-file` for a self-signed setup.

## See also

- [`hull repo add`](../cli/repo-add.md) · [`hull repo index`](../cli/repo-index.md)
  · [`hull repo update`](../cli/repo-update.md)
- [`hull pull`](../cli/pull.md) — download a package from a repo or OCI
- [`hull package`](../cli/package.md) — build the archives you host
- [OCI](oci.md) — registry-based distribution
- [Signing](signing.md) — sign and verify packages
{% endraw %}
