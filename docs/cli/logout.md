# hull logout

Remove the stored credentials for a registry host.

## When to use it

- When rotating or revoking credentials, or when you no longer need access to
  a host.
- On shared or CI machines, to clear a credential you set with `hull login`.

## What happens

1. hull looks up `<host>` in `~/.config/hull/credentials.json`.
2. If a credential is stored, it is deleted and the file is rewritten; hull
   prints `Logout succeeded for <host>`.
3. If nothing is stored for that host, hull prints `Not logged in to <host>`
   and changes nothing.

Only the saved credential is removed. Any repository registration and cached
packages stay in place.

## Usage

```
hull logout <host> [flags]
```

## Flags

Inherits the global flags.

## Worked example

Log out of a host you previously authenticated to:

```sh
hull logout registry.example.com
```

```
Logout succeeded for registry.example.com
```

Run it again — nothing is stored now, so hull says so and exits 0:

```sh
hull logout registry.example.com
```

```
Not logged in to registry.example.com
```

## See also

- [`login`](login.md) — store a credential for a host
- [`publish`](publish.md)
- [`registry`](registry.md)
