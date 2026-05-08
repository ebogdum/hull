# hull controller install-crd

## Synopsis

`hull controller install-crd` applies the `HullRelease` CRD to the cluster. The CRD lets operators describe a release declaratively (`apiVersion: hull.dev/v1, kind: HullRelease`) and have the in-cluster reconciler converge it. The CRD must exist before `hull controller run` can watch and reconcile any `HullRelease` objects.

## When to use it

Run once per cluster as a prerequisite step before deploying the controller. Re-running is safe: the apply is idempotent and will leave an existing CRD untouched if the schema already matches, or upgrade it in place if it has changed.

## What happens when you run it

1. Hull connects to the cluster using the active kubeconfig context.
2. Server-side applies the embedded `HullRelease` CRD definition (group `hull.dev`).
3. Waits briefly for the CRD to reach `Established=true`.
4. Prints the applied object's name on success, or the API server error on failure.

## Usage

```
hull controller install-crd [flags]
```

## Flags

| Flag | Type | Default | Description |
|---|---|---|---|
| `-h, --help` | bool | false | help for install-crd |

## Persistent flags inherited from `hull`

| Flag | Type | Description |
|---|---|---|
| `--debug` | bool | enable debug output |
| `--kube-context` | string | Kubernetes context to use |
| `--kubeconfig` | string | path to kubeconfig file |
| `-n, --namespace` | string | Kubernetes namespace |

## Examples

Install the CRD using the current kubeconfig context:

```sh
hull controller install-crd
```

Install into an explicit cluster context:

```sh
hull controller install-crd --kube-context prod-cluster
```

Verify the CRD landed:

```sh
hull controller install-crd && kubectl get crd hullreleases.hull.dev
```

## See also

- [`controller`](controller.md)
- [`controller crd`](controller-crd.md) — print without applying
- [`controller run`](controller-run.md) — start the reconciler
