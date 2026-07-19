---
title: "Quickstart"
parent: "Guides"
---
{% raw %}
# Quickstart

By the end of this guide you will have scaffolded a package, installed it as a
release, upgraded it, previewed and reverted changes, inspected drift, and
uninstalled it — all with the real commands and their real output. It takes
about ten minutes and assumes a working `kubectl` against any cluster (`kind`,
`k3d`, `k3s`, `minikube`, or a managed service).

## Prerequisites

- A Kubernetes cluster reachable through your current `kubectl` context.
- The `hull` binary on your `$PATH`. See
  [Quick install](../../README.md#quick-install) if you do not have it yet.
- Permission to create namespaces, Deployments, and Services in the cluster.

Sanity check:

```sh
hull version
kubectl get nodes
```

## 1. Scaffold a package

```sh
hull create hello
cd hello
```

`hull create` writes:

```
hello/
├── .hullignore
├── hull.yaml
├── values.yaml
└── templates/
    ├── _helpers.yaml
    ├── deployment.yaml
    ├── notes.yaml
    └── service.yaml
```

The scaffold is a minimal nginx Deployment plus a Service. `hull.yaml` declares
the package; `values.yaml` is the configurable surface; `templates/` holds the
manifests, written with `${...}` expressions. `_helpers.yaml` (leading
underscore) is a *partial* — a bag of reusable snippets, never emitted as a
standalone manifest. `templates/notes.yaml` is a document with a single
`message:` key; hull treats any such document as the release notes rather than a
manifest.

The default `values.yaml`:

```yaml
name: hello
replicaCount: 1
image:
  repository: nginx
  tag: latest
service:
  port: 80
```

## 2. Lint it

```sh
hull lint .
```

```
lint passed
```

`hull lint` checks that `hull.yaml` and `values.yaml` parse, that
`values.schema.json` (if present) is valid JSON, and that every template
renders. It does **not** validate values against the schema — that happens at
render time (step 3). See [`hull lint`](../cli/lint.md).

## 3. Render templates locally

```sh
hull template .
```

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  labels:
    app: hello
  name: hello
spec:
  replicas: 1
  selector:
    matchLabels:
      app: hello
  template:
    metadata:
      labels:
        app: hello
    spec:
      containers:
        - image: nginx:latest
          name: hello
          ports:
            - containerPort: 80
---
apiVersion: v1
kind: Service
metadata:
  labels:
    app: hello
  name: hello
spec:
  ports:
    - port: 80
      protocol: TCP
      targetPort: 80
  selector:
    app: hello
  type: ClusterIP
```

`hull template` renders every template against `values.yaml` and prints the
manifests to stdout. **Nothing has touched the cluster.** Override a value with
`--set`:

```sh
hull template . --set replicaCount=3
```

or with a values file:

```sh
hull template . -f overrides.yaml
```

See [`hull template`](../cli/template.md).

## 4. Install

```sh
hull install hello . -n hull-quickstart --create-namespace
```

```
NOTES:
hello has been installed successfully.
Namespace: hull-quickstart
Run "kubectl get deployments" to verify.
```

Hull renders the package, stamps `managedBy=hull` on every resource,
server-side-applies the manifest, waits for the resources to become ready, and
stores the release record as a labelled Secret in the namespace. Check what
landed:

```sh
kubectl -n hull-quickstart get all -l managedBy=hull
```

```
NAME                        READY   STATUS    RESTARTS   AGE
pod/hello-bc94584c5-lfgvr   1/1     Running   0          24s

NAME            TYPE        CLUSTER-IP     EXTERNAL-IP   PORT(S)   AGE
service/hello   ClusterIP   10.43.65.214   <none>        80/TCP    24s

NAME                    READY   UP-TO-DATE   AVAILABLE   AGE
deployment.apps/hello   1/1     1            1           25s
```

See [`hull install`](../cli/install.md).

## 5. List, status, manifest

```sh
hull list -n hull-quickstart
```

```
NAME     NAMESPACE          REVISION    STATUS      PACKAGE    VERSION    UPDATED
hello    hull-quickstart    1           deployed    hello      0.1.0      2026-07-18 22:06:02
```

```sh
hull status hello -n hull-quickstart
```

```
NAME:       hello
NAMESPACE:  hull-quickstart
STATUS:     deployed
REVISION:   1
PACKAGE:    hello-0.1.0
UPDATED:    2026-07-18 22:06:02

NOTES:
hello has been installed successfully.
...
```

`hull get manifest hello -n hull-quickstart` prints the exact YAML hull stored
for the revision; `hull get values hello -n hull-quickstart` prints the merged
values it used. Use `hull list -A` to see releases across all namespaces. See
[`hull list`](../cli/list.md), [`hull status`](../cli/status.md), and
[`hull get`](../cli/get.md).

## 6. Upgrade

Edit `values.yaml` or a template, then re-apply. Here, bump the replica count:

```sh
hull upgrade hello . -n hull-quickstart --set replicaCount=3
```

Each upgrade increments the revision counter, stores the new manifest, and
server-side-applies it. `hull history` now shows two revisions:

```sh
hull history hello -n hull-quickstart
```

```
REVISION    STATUS        PACKAGE        UPDATED                DESCRIPTION
1           superseded    hello-0.1.0    2026-07-18 22:06:02
2           deployed      hello-0.1.0    2026-07-18 22:07:32
```

See [`hull upgrade`](../cli/upgrade.md) and [`hull history`](../cli/history.md).

## 7. Preview changes before applying

`hull diff` compares local inputs only — it never reads the cluster. Render the
package two ways and diff them:

```sh
hull diff . --to-set replicaCount=5
```

```
diff: from → to

~ update  Deployment/hello
      ~ spec.replicas
          - 1
          + 5

Summary: 0 added, 1 changed, 0 removed.
```

To compare the package against what hull last recorded for the release, use
`hull plan` instead:

```sh
hull plan . -r hello -n hull-quickstart --action upgrade
```

```
hull plan: update  hello / hull-quickstart  (package .)

~ update  Deployment/hello
      from: deployment.yaml
      ~ spec.replicas
          - 3   (state)
          + 1   ← package-default (values.yaml)

Plan: 0 to add, 1 to change, 0 to destroy.
```

See [`hull diff`](../cli/diff.md) and [`hull plan`](../cli/plan.md).

## 8. Plan and apply

For change-management workflows, separate "what hull would do" from "do it".
`hull plan --out` writes a self-contained JSON artifact (rendered manifest plus
a sha256 integrity digest, bound to the release name and namespace):

```sh
hull plan . -r hello -n hull-quickstart --action upgrade --set replicaCount=5 --out plan.json
```

```
plan written to plan.json
```

```sh
hull apply --plan plan.json -n hull-quickstart
```

```
applied upgrade for hello revision 3
```

`hull apply` executes exactly what the plan captured. See
[`hull apply`](../cli/apply.md).

## 9. Roll back

```sh
hull rollback hello 1 -n hull-quickstart
```

Hull re-applies revision 1's stored manifest and records a new revision. The
audit trail records every action:

```sh
hull audit hello -n hull-quickstart
```

```
REVISION    ACTION     USER      STATUS        TIMESTAMP
1           install    bogdan    superseded    2026-07-18 22:06:02
2           upgrade    bogdan    deployed      2026-07-18 22:07:32
```

See [`hull rollback`](../cli/rollback.md) and [`hull audit`](../cli/audit.md).

## 10. Detect and reconcile drift

`hull drift` compares three views — the package as it renders now, the recorded
state, and the live cluster. It locates each live object by name **and
namespace**, so the resource templates must carry their namespace. Add one line
under `metadata` in `templates/deployment.yaml` and `templates/service.yaml`:

```yaml
metadata:
  name: "${values.name}"
  namespace: ${release.namespace}
```

Apply the edit, then change the cluster out of band and compare:

```sh
hull upgrade hello . -n hull-quickstart
kubectl -n hull-quickstart scale deploy hello --replicas=7
hull drift . -r hello -n hull-quickstart
```

```
drift: package ↔ state ↔ running   (release hello)

~ differs                Deployment/hello  (namespace hull-quickstart)
      spec.replicas  ⚠ cluster drift
          package: 1
          state:   1
          running: 7

1 cluster-drift, 0 pending-apply, 0 orphan, 0 missing, 0 to-create.
```

Push the recorded state back onto the cluster:

```sh
hull reconcile hello -n hull-quickstart
```

```
Reconciled 2 resource(s):
  - Deployment/hello
  - Service/hello
```

See [`hull drift`](../cli/drift.md) and [`hull reconcile`](../cli/reconcile.md).

## 11. Uninstall

```sh
hull uninstall hello -n hull-quickstart
```

Hull deletes the release's resources. History is **kept** by default (so
`hull audit` and `hull rollback` still work); pass `--purge` to delete the
release record too. List kept-history releases with:

```sh
hull list --uninstalled -n hull-quickstart
```

Remove the namespace when you are done:

```sh
kubectl delete ns hull-quickstart
```

See [`hull uninstall`](../cli/uninstall.md).

## Next steps

- [Package anatomy](packages.md) — every file in a package.
- [Values](values.md) — how `values.yaml`, layers, environments, profiles, and
  CLI flags merge.
- [Layers](layers.md) — composing a package from reusable building blocks.
- [Schema validation](schema-validation.md) — `values.schema.json` patterns.
- [Hooks](hooks.md) — lifecycle Jobs and Pods.
- [Workspaces](workspaces.md) — orchestrating many packages at once.
- [Template expressions](../templates/expressions.md) and
  [function reference](../templates/functions.md) — the `${...}` language.
- [CLI reference](../cli/README.md) — every command and flag.
{% endraw %}
