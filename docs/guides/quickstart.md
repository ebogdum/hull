# Quickstart

This guide takes you from zero to a deployed, upgraded, and rolled-back release in about ten minutes. It assumes you already have a working `kubectl` against some Kubernetes cluster — anything from `kind`, `k3d`, `k3s` and `minikube` to a managed service will work.

## Prerequisites

- A Kubernetes cluster reachable through your current `kubectl` context.
- The `hull` binary on your `$PATH`. See [Installing Hull](../../README.md#installing-hull) if not yet installed.
- Permission to create namespaces, Deployments, Services, ConfigMaps, and Secrets in the target cluster.

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
├── README.md
├── hull.yaml
├── values.yaml
└── templates/
    ├── _helpers.yaml
    ├── configmap.yaml
    ├── deployment.yaml
    └── service.yaml
```

The scaffolded package is a small "hello-world" web server. Open the files and look around — `hull.yaml` declares the package; `values.yaml` is the configurable surface; `templates/` is where YAML manifests with `${...}` expressions live. The `_helpers.yaml` file (note the leading underscore) is a *partial*: included from other templates, never rendered as a standalone manifest.

## 2. Lint it

```sh
hull lint .
```

Lint runs YAML parsing, schema validation (if `values.schema.json` is present), template rendering, and a static set of best-practice checks (e.g. resources without limits, hostPort usage, image tags pinned to `latest`). A clean lint exits 0.

## 3. Render templates locally

```sh
hull template .
```

This renders every template against the package's `values.yaml`, with the release name `hello`, and prints the resulting Kubernetes manifests to stdout. **Nothing has touched the cluster yet.** Use `hull template` whenever you want to inspect what hull would apply.

To override a value:

```sh
hull template . --set replicas=3
```

Or with a values file:

```sh
hull template . -f overrides.yaml
```

## 4. Install

```sh
hull install hello . -n hull-quickstart --create-namespace
```

What happens:

1. The package is rendered the same way `hull template` rendered it.
2. Hull stamps `managedBy=hull` on every resource and on the namespace it creates.
3. Pre-install hooks (if any) run.
4. Server-side apply pushes the rendered manifest into the cluster.
5. Post-install hooks run.
6. The release record is stored as a labelled Secret in the install namespace.
7. Hull waits for the rendered resources to become Ready (Deployment available, Pod Ready, etc.) and prints a summary.

Check what landed:

```sh
kubectl -n hull-quickstart get deploy,svc,cm
kubectl -n hull-quickstart get all -l managedBy=hull
```

## 5. List, status, manifest

```sh
hull list                          # releases in the current namespace
hull list -A                       # every release everywhere
hull status hello -n hull-quickstart
hull get manifest hello -n hull-quickstart
hull get values hello -n hull-quickstart
```

`hull status` shows the current revision, package, and per-resource readiness. `hull get manifest` prints the exact YAML hull stored for this revision (gzipped + base64 inside the release Secret; hull decompresses it on the fly).

## 6. Upgrade

Edit `values.yaml` (e.g. bump replicas) or any template, then:

```sh
hull upgrade hello . -n hull-quickstart
```

Each upgrade increments the release's revision counter and stores the new manifest. The cluster sees server-side apply with the new fields. Hooks tagged `pre-upgrade` and `post-upgrade` run before and after the apply.

## 7. Diff before applying

`hull diff` shows what would change without applying. It always uses server-side dry-run, so the cluster's defaulters and webhooks contribute to the comparison — you see the diff the cluster would actually compute.

```sh
hull diff hello . -n hull-quickstart
```

## 8. Plan and apply

For change-management workflows, separate "what hull would do" from "do it":

```sh
hull plan hello . -n hull-quickstart -o hello.plan
# review the plan...
hull apply hello.plan
```

`hull plan` produces a self-contained plan file (rendered manifest + parameters + hash). `hull apply` consumes one and executes the upgrade exactly as planned. The plan binds to the release name and namespace so you cannot accidentally apply it elsewhere.

## 9. History and rollback

```sh
hull history hello -n hull-quickstart
```

prints every revision with its timestamp, status, and audit data (who installed it, with what flags). To roll back:

```sh
hull rollback hello 1 -n hull-quickstart
```

Hull re-applies revision 1's stored manifest and re-runs revision 1's `pre-rollback` and `post-rollback` hooks (which are persisted alongside the manifest, so a rollback to an old revision uses the hooks that revision originally shipped, not the current ones).

## 10. Drift detection

After an install, the cluster might be edited out-of-band — somebody runs `kubectl edit deploy hello`, or another operator picks up a field. `hull drift` compares the live state against the release's stored manifest:

```sh
hull drift hello -n hull-quickstart
```

Drift output is per-resource, per-field. To re-converge to the stored manifest:

```sh
hull reconcile hello -n hull-quickstart
```

This re-applies the stored manifest, taking ownership of any drifted fields.

## 11. Audit

```sh
hull audit hello -n hull-quickstart
```

prints the full audit trail: every revision, the action (`install`, `upgrade`, `rollback`, `uninstall`), who initiated it, the kubeconfig context, the hull version, the flags as passed, and any value files supplied. This is signed metadata stored in the release record.

## 12. Uninstall

```sh
hull uninstall hello -n hull-quickstart
```

Pre-delete hooks run, the manifest's resources are deleted, post-delete hooks run, and the release record is removed. Pass `--keep-history` to keep the release record for forensic purposes; `hull list --filter status=uninstalled` shows kept-history releases.

To clean up the namespace too:

```sh
kubectl delete ns hull-quickstart
```

## Where next

- [Package anatomy](packages.md) — every file in a package, in detail.
- [Values](values.md) — how `values.yaml`, layers, environments, and CLI flags merge.
- [Layers](layers.md) — composing a package from reusable building blocks.
- [Hooks](hooks.md) — Job-based lifecycle hooks.
- [Workspaces](workspaces.md) — orchestrating many packages with one command.
- [Template expressions](../templates/expressions.md) — the `${...}` syntax.
- [Template functions](../templates/functions.md) — every built-in, with input/output examples.
