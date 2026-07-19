---
title: "hull controller install-crd"
parent: "CLI"
---
{% raw %}
# hull controller install-crd

## Synopsis

`hull controller install-crd` registers the `HullRelease` CustomResourceDefinition
in your cluster. It applies the same YAML that [`controller crd`](controller-crd.md)
prints, then waits for the API server to start serving the new kind. The CRD
must exist before [`controller run`](controller-run.md) can see any
`HullRelease` objects.

## When to use it

- Once per cluster, before you deploy the controller.
- Again after upgrading hull, to pick up CRD schema changes. Re-running is
  safe: the apply is idempotent.

## What happens

1. Hull connects to the cluster using the active kubeconfig context.
2. It applies the embedded `HullRelease` CRD (group `hull.dev`).
3. It waits (up to two minutes) for the CRD to become `Established` — the point
   at which the API server serves the kind.
4. On success the command exits 0 with no output. On failure it prints the API
   server error and exits non-zero (for example, when your context lacks
   permission to create CRDs).

## Usage

```
hull controller install-crd
```

## Flags

Inherits the global flags.

| Flag | Type | Default | Description |
|---|---|---|---|
| `--debug` | — | — | print debug output while applying |
| `--kube-context` | string | (current) | which cluster to apply the CRD to |
| `--kubeconfig` | string | (default) | path to the kubeconfig file |
| `-n, --namespace` | string | — | Kubernetes namespace (the CRD itself is cluster-scoped) |

## Worked example

Register the CRD, then confirm it landed:

```sh
hull controller install-crd
kubectl get crd hullreleases.hull.dev
```

**Output:**

```
# hull controller install-crd prints nothing and exits 0

# kubectl now sees the registered CRD:
NAME                    CREATED AT
hullreleases.hull.dev   2026-07-18T14:02:11Z
```

Target a specific cluster:

```sh
hull controller install-crd --kube-context prod
```

## See also

- [`controller crd`](controller-crd.md) — print the CRD without applying it
- [`controller run`](controller-run.md) — start the reconciler once the CRD exists
- [`controller`](controller.md) — operator overview
{% endraw %}
