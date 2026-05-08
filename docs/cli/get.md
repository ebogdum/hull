# hull get

## Synopsis

`hull get` retrieves the contents of a release record. Subcommands fetch specific subresources: the rendered manifest, the merged values used at install time, the rendered hooks (with last-run results), the rendered post-install notes, the package metadata, and a combined `all` view for one-shot extraction.

## When to use it

Use when you need to inspect what hull stored for a release — what it actually rendered, what values it used, what hooks it ran. Combine with `--revision N` to look at a historical revision.

## Usage

```
hull get [command]
```

## Subcommands

- [`hull get all`](get-all.md) — full release record
- [`hull get hooks`](get-hooks.md) — release hooks (manifests + last-run results)
- [`hull get manifest`](get-manifest.md) — rendered Kubernetes manifest
- [`hull get metadata`](get-metadata.md) — release metadata (name, package, revision, status)
- [`hull get notes`](get-notes.md) — post-install notes
- [`hull get values`](get-values.md) — merged values used at install

## Flags

| Flag | Type | Default | Description |
|---|---|---|---|
| `-h, --help` | — | — | help for get |

## Persistent flags inherited from `hull`

| Flag | Type | Description |
|---|---|---|
| `--debug` | — | enable debug output |
| `--kube-context` | string | Kubernetes context to use |
| `--kubeconfig` | string | path to kubeconfig file |
| `-n, --namespace` | string | Kubernetes namespace |

## Examples

The rendered manifest of the current revision:

```sh
hull get manifest my-app -n prod
```

Merged values for revision 3:

```sh
hull get values my-app --revision 3 -n prod
```

Everything (manifest, values, hooks, notes, metadata):

```sh
hull get all my-app -n prod -o yaml
```

## See also

- [`status`](status.md)
- [`history`](history.md)
- [`audit`](audit.md)
