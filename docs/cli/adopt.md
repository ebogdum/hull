# hull adopt

## Synopsis

`hull adopt` claims an existing in-cluster resource as part of a hull-managed release. Hull stamps the resource with `managedBy=hull`, generates a release record that references it, and from that point onwards the resource is upgrade-, rollback-, and drift-tracked the same as any other hull-managed resource.

## When to use it

Use when migrating an existing manually-applied workload (raw `kubectl apply`, Terraform, hand-written manifests) into hull's management without recreating it. Adoption is non-destructive — the resource keeps running while ownership changes; the only writes are the new `managedBy=hull` label and the release-storage Secret.

## What happens when you run it

1. Resolves each `<resource-ref>` to a live cluster object. Reference forms:
   - `apps/v1/Deployment/myns/myapp`
   - `v1/ConfigMap//cluster-scoped-cm`  (note the empty namespace segment)
   - `kind=Deployment,name=myapp,ns=myns`
2. Fetches each object, strips server-side metadata (resourceVersion, managedFields, status).
3. Composes a manifest, stamps `managedBy=hull` on each resource, and stores it as revision 1 of the named release.
4. The cluster resources themselves are not deleted or recreated — only the release record is new.

## Usage

```
hull adopt <release-name> <resource-ref>... [flags]
```

## Flags

| Flag | Type | Default | Description |
|---|---|---|---|
| `--create-namespace` | — | — | create the release namespace if it does not exist |
| `--description` | string | — | release description recorded in the audit trail |
| `-h, --help` | — | — | help for adopt |
| `--labels` | stringArray | — | label key=value to attach to the release (repeatable) |

## Persistent flags inherited from `hull`

| Flag | Type | Description |
|---|---|---|
| `--debug` | — | enable debug output |
| `--kube-context` | string | Kubernetes context to use |
| `--kubeconfig` | string | path to kubeconfig file |
| `-n, --namespace` | string | Kubernetes namespace |

## Examples

Adopt a Deployment by GVK/namespace/name:

```sh
hull adopt hello apps/v1/Deployment/prod/myapp -n prod
```

Adopt a cluster-scoped ConfigMap (note the empty namespace between the slashes):

```sh
hull adopt cluster-config v1/ConfigMap//global-config
```

Adopt several resources into one release:

```sh
hull adopt hello \
  apps/v1/Deployment/prod/myapp \
  v1/Service/prod/myapp \
  v1/ConfigMap/prod/myapp-config \
  -n prod
```

Adopt with a description for the audit trail:

```sh
hull adopt hello apps/v1/Deployment/prod/myapp -n prod \
  --description "imported from terraform module v3" \
  --labels source=terraform
```

## See also

- [`install`](install.md)
- [`drift`](drift.md)
