# hull releases status

## Synopsis

`hull releases status` queries the cluster for the current state of every release declared in `hull-releases.yaml`. For each entry, it reports the current revision, status (deployed / pending / failed / missing), package version, and last-deployed time — the same information you'd get from `hull status <release>` for each, in one tabular view. Releases that are declared but not installed show as `not installed`.

## When to use it

Run as a CI gate after `hull releases install` / `upgrade` to confirm every declared release is at the expected status. Also useful for routine platform health checks: a single command tells you whether the platform graph is fully deployed and Ready.

## What happens when you run it

1. Reads `--file` (default `hull-releases.yaml`) from the current directory.
2. For each entry, queries the cluster for the corresponding release record (in the entry's namespace).
3. Composes a table with revision, status, package, namespace, and last-deployed timestamp.
4. Prints to stdout. Exits 0 if all are deployed; non-zero if any are missing or failed.

## Usage

```
hull releases status [flags]
```

## Flags

| Flag | Type | Default | Description |
|---|---|---|---|
| `--file` | string | hull-releases.yaml | spec file path |
| `-h, --help` | bool | false | help for status |

## Persistent flags inherited from `hull`

| Flag | Type | Description |
|---|---|---|
| `--debug` | bool | enable debug output |
| `--kube-context` | string | Kubernetes context to use |
| `--kubeconfig` | string | path to kubeconfig file |
| `-n, --namespace` | string | Kubernetes namespace |

## Examples

Default — status of every release declared in `./hull-releases.yaml`:

```sh
hull releases status
```

Status of a custom-named platform manifest:

```sh
hull releases status --file ./platform.releases.yaml
```

CI health check — fail if any release isn't deployed:

```sh
hull releases status || { echo "platform not fully up"; exit 1; }
```

## See also

- [`releases`](releases.md)
- [`releases install`](releases-install.md)
- [`releases upgrade`](releases-upgrade.md)
- [`status`](status.md) — single-release status
- [`hull-releases.yaml` reference](../reference/hull-releases-yaml.md)
