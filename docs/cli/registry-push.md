---
title: "hull registry push"
parent: "CLI"
---
{% raw %}
# hull registry push

Upload a packaged `.hull.tgz` archive to an OCI registry reference.

## When to use it

- To publish a release-ready package to a container registry so others can
  `hull registry pull` or install it by reference.
- To promote a build between registries (pull from one, push to another).
- In CI, right after [`hull package`](package.md) produces the archive.

## What happens

1. Reads the local archive at `<archive>` (the file
   [`hull package`](package.md) produced).
2. Uses the credentials you stored with [`hull login`](login.md) for the host
   in `<ref>`. If none are stored, the registry must allow anonymous push.
3. Uploads the archive as an OCI artifact to `<ref>`, tag included. An untagged
   reference is tagged with the package's own version.
4. Prints `Pushed <archive> to <ref>`. The package is now retrievable at that
   reference by anyone with pull access.

## Usage

```
hull registry push <archive> <ref>
```

`<archive>` is a path to a `.hull.tgz` file. `<ref>` is an `oci://…`
reference, optionally with a `:tag`.

## Flags

| Flag | Type | Default | Description |
|---|---|---|---|
| `--plain-http` | bool | false | Talk to the registry over plaintext HTTP instead of HTTPS — for a local or in-cluster registry without TLS. |
| `--insecure-skip-tls-verify` | bool | false | Keep HTTPS but skip certificate validation — for a registry with a self-signed certificate. |

Global flags `--oci-plain-http`, `--oci-insecure-skip-tls-verify`, and
`--allow-plaintext-auth` are inherited from `hull`.

## Worked example

You built an archive and want it in GitHub Container Registry.

**INPUT** — log in once, then push the local file to a tagged reference:

```sh
hull login ghcr.io -u alice --password-stdin < token.txt
hull registry push ./dist/my-app-1.0.0.hull.tgz \
  oci://ghcr.io/example/charts/my-app:1.0.0
```

**OUTPUT:**

```
Login succeeded for ghcr.io
Pushed ./dist/my-app-1.0.0.hull.tgz to oci://ghcr.io/example/charts/my-app:1.0.0
```

**RESULT:** the tag `ghcr.io/example/charts/my-app:1.0.0` now exists in the
registry. Anyone with pull access can fetch it:

```sh
hull registry pull oci://ghcr.io/example/charts/my-app:1.0.0
```

For a local registry without TLS, add `--plain-http`:

```sh
hull registry push ./dist/my-app-1.0.0.hull.tgz \
  oci://localhost:5000/charts/my-app:1.0.0 --plain-http
```

## See also

- [`login`](login.md) — store the credentials this command uses
- [`package`](package.md) — build the archive you push
- [`registry pull`](registry-pull.md) — fetch it back
- [`publish`](publish.md) — push to an HTTP API registry instead
- [`install`](install.md) — install a package from a reference
{% endraw %}
