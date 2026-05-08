# CLI reference

Hull's command-line interface is rooted at `hull`. Every command has a `--help` output reachable via `hull <cmd> --help`; this directory holds the long-form reference for each, with description, every flag, and examples.

## Global flags

These apply to every command and can also be set as defaults via environment variables.

| Flag | Env | Description |
|---|---|---|
| `-n, --namespace <ns>` | — | Kubernetes namespace for the operation. |
| `--kubeconfig <path>` | `KUBECONFIG` | Path to kubeconfig file. |
| `--kube-context <name>` | — | Active kubeconfig context. |
| `--debug` | `HULL_DEBUG` | Enable debug logging. |

## Environment variables

| Variable | Purpose |
|---|---|
| `HULL_CONFIG_HOME` | Config dir (default: `~/.config/hull`). |
| `HULL_CACHE_HOME` | Cache dir (default: `~/.cache/hull`). |
| `HULL_DATA_HOME` | Data dir (default: `~/.local/share/hull`). |
| `HULL_DEBUG` | Enable debug logging when truthy. |
| `HULL_OCI_PLAIN_HTTP` | Use HTTP for OCI instead of HTTPS. |
| `HULL_OCI_INSECURE_SKIP_TLS` | Skip TLS verification for OCI. |
| `HULL_DRIVER` | Release storage driver (`secret`, `configmap`, `memory`, `sql`). Default `secret`. |
| `HULL_DRIVER_SQL_DSN` | SQL DSN when driver is `sql`. |
| `HULL_NAMESPACE` (or `HELM_NAMESPACE`) | Default namespace; legacy aliases recognised. |
| `KUBECONFIG` | Standard Kubernetes kubeconfig path. |

## Exit codes

| Code | Meaning |
|---|---|
| `0` | Success. |
| `1000-1999` | Input validation (CLI args, manifest format). |
| `2000-2999` | Business logic (release not found, version conflict). |
| `3000-3999` | Resource (cluster resource missing or unavailable). |
| `4000-4999` | Permission (RBAC, auth). |
| `5000-5999` | System (timeout, dependency failure, internal error). |

Hull prints both the textual error and the numeric code to stderr; scripts can branch on the code.

---

## Commands by category

### Release lifecycle

| Command | One-liner |
|---|---|
| [`install`](install.md) | Install a hull package as a new release. |
| [`upgrade`](upgrade.md) | Upgrade an existing release to a new revision. |
| [`rollback`](rollback.md) | Roll back a release to a previous revision. |
| [`uninstall`](uninstall.md) | Uninstall a release. |
| [`status`](status.md) | Show current release status and resource readiness. |
| [`list`](list.md) | List releases. |
| [`history`](history.md) | Show release history. |
| [`get`](get.md) | Get release details (manifest, values, hooks, notes). |
| [`audit`](audit.md) | Show the chronological audit trail for a release. |
| [`prune`](prune.md) | Drop superseded revisions, keeping the most recent N. |
| [`rename`](rename.md) | Rename a release in-place (preserve history). |
| [`releases`](releases.md) | Manage cross-release dependencies declared in `hull-releases.yaml`. |

### Diff, plan, drift

| Command | One-liner |
|---|---|
| [`diff`](diff.md) | Show what would change on upgrade. |
| [`plan`](plan.md) | Render and persist an apply-able plan. |
| [`apply`](apply.md) | Execute a previously-saved plan. |
| [`drift`](drift.md) | Detect drift between the stored manifest and live cluster state. |
| [`reconcile`](reconcile.md) | Re-apply the stored manifest to converge cluster state. |
| [`canary`](canary.md) | Staged upgrade through replica percentages with bake periods. |
| [`multi-install`](multi-install.md) | Install a release into multiple clusters. |

### Authoring

| Command | One-liner |
|---|---|
| [`create`](create.md) | Scaffold a new package. |
| [`init`](init.md) | Scaffold from a built-in template. |
| [`lint`](lint.md) | Validate a package for correctness. |
| [`template`](template.md) | Render templates locally. |
| [`debug`](debug.md) | Debug template rendering. |
| [`dev`](dev.md) | Watch a package and re-render on changes. |
| [`config`](config.md) | Interactively build a values file from `values.schema.json`. |
| [`values`](values.md) | Show effective values, optionally with per-key resolution trace. |
| [`scan`](scan.md) | Find common values across packages; extract a base layer. |
| [`policy`](policy.md) | Evaluate package policies against rendered manifests. |
| [`graph`](graph.md) | Render a dependency graph of a release. |
| [`metrics`](metrics.md) | Sample CPU/memory; recommend requests/limits. |

### Packaging

| Command | One-liner |
|---|---|
| [`package`](package.md) | Package as a `.hull.tgz` archive. |
| [`publish`](publish.md) | Publish a package to a registry. |
| [`sbom`](sbom.md) | Emit a CycloneDX 1.5 SBOM for a release. |
| [`adopt`](adopt.md) | Claim existing in-cluster resources as a hull-managed release. |
| [`dependency`](dependency.md) | Manage layers and required packages. |

### Repositories and OCI

| Command | One-liner |
|---|---|
| [`repo`](repo.md) | Manage hull package repositories. |
| [`registry`](registry.md) | Manage OCI registries. |
| [`login`](login.md) | Store credentials for a package registry. |
| [`logout`](logout.md) | Remove stored credentials. |
| [`pull`](pull.md) | Download a package from a repository or OCI. |
| [`registry push`](registry-push.md) | Push a package archive to an OCI registry. |
| [`search`](search.md) | Search for packages. |
| [`show`](show.md) | Show package metadata, values, README, etc. |

### Signing

| Command | One-liner |
|---|---|
| [`keyring`](keyring.md) | Manage the PGP keyring used for provenance verification. |
| [`package verify`](package-verify.md) | Verify a `.prov` signature against the local keyring. |

### Workspaces

| Command | One-liner |
|---|---|
| [`workspace`](workspace.md) | Orchestrate multiple packages declared in `hull-workspace.yaml`. |

### Plugins and marketplace

| Command | One-liner |
|---|---|
| [`plugin`](plugin.md) | Manage hull plugins. |
| [`marketplace`](marketplace.md) | Browse and install signed plugins. |

### Test

| Command | One-liner |
|---|---|
| [`test`](test.md) | Run tests for a release. |

### Operations

| Command | One-liner |
|---|---|
| [`controller`](controller.md) | Reconcile HullRelease CRs declared in the cluster. |
| [`purge`](purge.md) | Clean up everything hull has installed. |

### Compat and migration

| Command | One-liner |
|---|---|
| [`helm-compat`](helm-compat.md) | Helm chart compatibility helpers. |
| [`migrate`](migrate.md) | Convert a Helm chart to a hull package. |

### Misc

| Command | One-liner |
|---|---|
| [`version`](version.md) | Print the hull version. |
| [`env`](env.md) | Print hull's environment information. |
| [`completion`](completion.md) | Generate shell completion scripts. |
