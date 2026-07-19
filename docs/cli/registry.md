# hull registry

Push and pull hull package archives to and from OCI-compliant registries.

`hull registry` groups the OCI transport commands. Use it to move a packaged
`.hull.tgz` archive into a container registry (GHCR, ECR, Docker Hub, Harbor,
Zot, …) and to fetch one back. Credentials come from `hull login`; the artifact
is stored as a standard OCI blob under the reference you name.

## Subcommands

| Command | What it does |
|---|---|
| [`registry push`](registry-push.md) | Upload a local `.hull.tgz` archive to an OCI reference. |
| [`registry pull`](registry-pull.md) | Download a package from an OCI reference, optionally verifying its cosign signature first. |

## Usage

```
hull registry [command]
```

Log in once per host, then push or pull:

```sh
hull login ghcr.io -u USER --password-stdin
hull registry push ./my-app-1.0.0.hull.tgz oci://ghcr.io/example/charts/my-app:1.0.0
hull registry pull  oci://ghcr.io/example/charts/my-app:1.0.0 -d ./pulled
```

## See also

- [`login`](login.md) — store registry credentials
- [`logout`](logout.md) — remove stored credentials
- [`package`](package.md) — build the `.hull.tgz` archive you push
- [`publish`](publish.md) — publish to an HTTP API registry instead of OCI
- [`pull`](pull.md) — fetch from OCI or an HTTP repository by chart name
