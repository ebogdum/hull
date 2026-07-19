# hull env

`hull env` prints the resolved paths and settings hull uses, one
`KEY="value"` per line.

## When to use it

- Confirm which cache, config, and data directories hull reads before you
  debug a path or plugin problem.
- Check the namespace and kubeconfig hull resolves from your flags and
  environment.

## What happens

1. Resolves each value from its environment variable, falling back to the OS
   default when unset (for example `HULL_CACHE_HOME` defaults to the user cache
   directory, and derived paths like `HULL_REPOSITORY_CACHE` sit under it).
2. Resolves the namespace in order: `--namespace` flag → `HULL_NAMESPACE` →
   `HELM_NAMESPACE` → `default`.
3. Sorts the keys and prints them as quoted `KEY="value"` lines.

No cluster is contacted.

## Usage

```
hull env [flags]
```

## Flags

Inherits the global flags. `-n, --namespace` changes the reported
`HULL_NAMESPACE`.

## Worked example

**INPUT** — no overrides set. Run `hull env`:

```
HULL_BIN="/usr/local/bin/hull"
HULL_CACHE_HOME="/Users/you/Library/Caches/hull"
HULL_CONFIG_HOME="/Users/you/Library/Application Support/hull"
HULL_DATA_HOME="/Users/you/.local/share/hull"
HULL_KUBECONFIG=""
HULL_KUBECONTEXT=""
HULL_NAMESPACE="default"
HULL_PLUGINS="/Users/you/.local/share/hull/plugins"
HULL_REGISTRY_CONFIG="/Users/you/Library/Application Support/hull/registry.json"
HULL_REPOSITORY_CACHE="/Users/you/Library/Caches/hull/repository"
HULL_REPOSITORY_CONFIG="/Users/you/Library/Application Support/hull/repositories.yaml"
```

**Now set values and watch the output follow.** With `HULL_CACHE_HOME` set, the
derived repository cache moves with it (`HULL_CACHE_HOME=/data/cache hull env`):

```
HULL_CACHE_HOME="/data/cache"
HULL_REPOSITORY_CACHE="/data/cache/repository"
```

Pass `-n staging` and the resolved namespace changes
(`hull -n staging env`):

```
HULL_NAMESPACE="staging"
```

Each printed line reflects exactly what you set, so you can verify hull reads
the values you expect.

## See also

- [`config`](config.md) — build a values file interactively
- [`version`](version.md) — print the hull build version
