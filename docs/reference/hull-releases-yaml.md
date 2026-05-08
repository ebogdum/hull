# `hull-releases.yaml` Reference

`hull-releases.yaml` describes a **set of separate releases** with optional cross-release dependencies. It is consumed by the `hull releases` family of commands. Conceptually it sits between `hull-workspace.yaml` (which assumes one workspace = one repo of related packages) and free-form scripts that call `hull install` repeatedly: it lets you orchestrate releases of *unrelated* packages — possibly from different repositories, registries, or local paths — by declaring them in one file with a dependency graph.

Use this file when you have, say, a platform-engineering repo that rolls out infrastructure releases (cert-manager, external-dns, ingress, monitoring), each from a different upstream package, with strict ordering between them.

---

## Top-level structure

```yaml
releases:
  - name: ...
    package: ...
    ...
```

A single top-level key, `releases`, holding a list of `crossReleaseEntry` objects.

---

## `crossReleaseEntry` schema

| Field | Type | Required | Description |
|---|---|---|---|
| `name` | string | yes | The release name. Becomes the `<release>` argument hull would otherwise take from the CLI. Unique within the file. |
| `package` | string | yes | Where to find the package. Accepts the same forms as a layer `source` (local path, OCI URI, HTTPS chart URL, `git::...`). |
| `namespace` | string | no | Namespace to install into. Defaults to the active kubeconfig namespace if omitted. |
| `profile` | string | no | Profile to activate when rendering this release. |
| `values` | string list | no | Value files to merge. Paths are resolved relative to the directory containing `hull-releases.yaml`. |
| `set` | string list | no | `--set`-style overrides applied on top of `values` files. Same `key=value` syntax as the CLI flag. |
| `dependsOn` | string list | no | Names of other releases (in this same file) that must be installed/upgraded before this one. |

`dependsOn` is local to the file — entries can only reference other entries in the same file.

---

## Commands

```
hull releases plan       --file hull-releases.yaml   # show topological order
hull releases install    --file hull-releases.yaml   # install all releases in order
hull releases upgrade    --file hull-releases.yaml   # upgrade (or install if missing)
hull releases uninstall  --file hull-releases.yaml   # reverse-order uninstall
hull releases status     --file hull-releases.yaml   # per-release current status
```

The `--file` flag defaults to `hull-releases.yaml` in the working directory.

---

## Example: platform bootstrap

```yaml
releases:
  - name: cert-manager
    package: oci://quay.io/jetstack/cert-manager
    namespace: cert-manager
    values: [./values/cert-manager.yaml]
    set:
      - installCRDs=true

  - name: external-dns
    package: oci://ghcr.io/kubernetes-sigs/external-dns
    namespace: external-dns
    values: [./values/external-dns.yaml]
    dependsOn: [cert-manager]   # uses webhook from cert-manager

  - name: ingress
    package: ./packages/ingress-nginx
    namespace: ingress
    dependsOn: [cert-manager]

  - name: monitoring
    package: oci://ghcr.io/example/kube-prometheus-stack
    namespace: monitoring
    values:
      - ./values/monitoring-base.yaml
      - ./values/monitoring-overrides.yaml
    dependsOn: [ingress]
```

Plan output:

```
$ hull releases plan
Level 0:
  - cert-manager
Level 1:
  - external-dns
  - ingress
Level 2:
  - monitoring
```

Install:

```
$ hull releases install
[level 0] installing cert-manager...
[level 1] installing external-dns...
[level 1] installing ingress...
[level 2] installing monitoring...
```

Uninstall is the reverse:

```
$ hull releases uninstall
[level 2] uninstalling monitoring...
[level 1] uninstalling ingress...
[level 1] uninstalling external-dns...
[level 0] uninstalling cert-manager...
```

---

## Difference from `hull-workspace.yaml`

| Aspect | `hull-workspace.yaml` | `hull-releases.yaml` |
|---|---|---|
| Member packages | typically siblings in one repo | anywhere — local paths, OCI, HTTPS, git |
| Field name for package list | `members` (each with `path`) | `releases` (each with `package`) |
| Per-member overrides | `valueFiles` and inline `values` not on the entry; values come from the member's own `values.yaml` and the workspace's environment selection | `values:` and `set:` on the entry itself |
| Atomic / wait flags | per-member `atomic` / `wait` in the file | not exposed; releases use the standard `hull install` defaults |
| Parallel within a level | `--parallel` on `hull workspace install` | sequential within a level |
| Health-gate | `--health-gate` on `hull workspace install` | per-release `--wait` semantics (the default) |

You can use both in the same project: a workspace for tightly-coupled packages from one repo, and a `hull-releases.yaml` to roll out the workspace plus a few external releases together.
