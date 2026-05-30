# Hull Use Cases — Kubernetes Package Management for Platform Engineers, SRE Teams, and GitOps Workflows

This page maps hull's features onto real workflows. Pick the role or pattern that matches your team and follow the links to the relevant guides.

## For platform engineers

You're building an internal developer platform on Kubernetes. Application teams need a predictable way to ship workloads without writing kubectl pipelines. Compliance wants signed artifacts and audit trails. Operations wants drift detection.

**What hull gives you:**

- A versioned, signable package format (`*.hull.tgz`) that captures everything an app needs to run on Kubernetes
- A composition system (layers) so platform-provided base layers (logging, RBAC, network policies) flow into every app package
- Per-release audit data (who, when, where, with what flags) baked into the cluster's release record
- An OCI distribution model that integrates with your existing container registry and IAM
- A workspace orchestrator for multi-component platform rollouts
- Cluster-wide ownership queries via a single label selector (`managedBy=hull`)

**Workflow recipes:**

- **Standard library of layers** — publish base layers (`org-base-rbac`, `org-base-monitoring`, `org-base-network`) to OCI; consumers pull them in via `layers:` in `hull.yaml`. Updates propagate by bumping the version constraint.
- **Multi-environment promotion** — define `environments:` (`dev`/`staging`/`prod`) inside `hull.yaml` so the *same* package metadata flows through environments. CI pipelines run `hull plan --env staging` then `hull apply` after review.
- **Self-service signed packages** — application teams `hull package --sign` against your CI's PGP key; platform's keyring on every cluster admits only that key. Tampered packages fail verification.
- **Auditable rollouts** — every install records flags, user, and source. `hull audit <release>` answers "who upgraded the auth service yesterday?" months later.

→ Start with: [Quickstart](guides/quickstart.md), [Layers](guides/layers.md), [Workspaces](guides/workspaces.md), [Signing](guides/signing.md).

## For SRE / operations teams

You operate clusters. You need to know what's deployed, what's drifted, what's at risk. You need rollback that works.

**What hull gives you:**

- `hull list -A` — every release in every namespace
- `hull drift <release>` — per-field comparison of stored manifest vs live state, with smart filtering for noise
- `hull reconcile <release>` — converge cluster state back to the stored manifest in one command
- `hull audit <release>` — full chronological history with audit data
- `hull rollback <release> <rev>` — re-apply a previous revision and re-run that revision's hooks
- `hull metrics <release>` — sample CPU/memory and recommend requests/limits
- `hull multi-install --to <ctx-list> --atomic-cross-cluster` — fleet-wide rollouts with rollback
- `hull canary` — staged upgrades with bake periods between replica counts

**Workflow recipes:**

- **Drift-detection cron** — schedule `hull drift -A` against every namespace; alert on non-empty output. The list of managed releases is enumerable via `hull list`, so the operator can drive its own scheduler.
- **Incident rollback drill** — practise `hull rollback <release> <prev-rev>` against staging before the next on-call incident. The behaviour is identical in prod.
- **Capacity right-sizing** — run `hull metrics <release>` after a release has baked for 24 hours; commit the recommended `requests` / `limits` to the package's `values.yaml`.
- **Force-cleanup after node failure** — `hull purge --yes --force` clears wedged releases and force-finalises stuck Terminating namespaces.
- **Multi-cluster atomic platform** — define a `hull-releases.yaml` with `dependsOn`-ordered releases; `hull releases install` brings up the platform; `hull releases uninstall` tears it down in reverse.

→ Start with: [`hull drift`](cli/drift.md), [`hull reconcile`](cli/reconcile.md), [`hull audit`](cli/audit.md), [`hull canary`](cli/canary.md), [`hull purge`](cli/purge.md).

## For GitOps teams (Argo CD, Flux)

You declare desired state in git and reconcile it into the cluster. You want hull's packaging story without giving up your GitOps reconciler.

**What hull gives you:**

- `hull plan` produces a deterministic rendered manifest you can commit to a git repo
- `hull apply <plan>` re-checks the plan's integrity (SHA-256 of the manifest) before applying — detects drift between plan and apply
- `hull controller` is an in-cluster HullRelease CR reconciler if you want hull to do reconciliation directly
- The `hull` binary can be invoked as a CMP (Config Management Plugin) for Argo CD
- Flux's OCI source can pull hull packages directly via `hull registry push`

**Workflow recipes (Argo CD):**

- **CMP plugin** — register `hull template <package>` as a CMP that emits the rendered manifest. Argo CD diffs and syncs as usual.
- **Pre-rendered manifest in git** — CI runs `hull plan ... -o manifest.yaml`, commits to a git path Argo CD watches. Plan integrity (SHA-256) ensures the rendered manifest hasn't been edited.
- **HullRelease CR pattern** — commit `HullRelease` CRs to git. Argo CD syncs the CRs; `hull controller` reconciles them.

**Workflow recipes (Flux):**

- **OCI source + Kustomize controller** — `hull registry push` to OCI; Flux's OCISource pulls; Kustomize-controller applies.
- **HullRelease CR + hull controller** — Flux syncs CRs; `hull controller` reconciles.

→ Start with: [`hull plan`](cli/plan.md), [`hull apply`](cli/apply.md), [`hull controller`](cli/controller.md), [Workspaces guide](guides/workspaces.md).

## For application developers

You write code. You want to ship to Kubernetes without becoming a YAML expert.

**What hull gives you:**

- `hull create my-app` — scaffold a complete, working package in seconds (Deployment + Service + ConfigMap + helpers)
- `hull init <template>` — start from richer templates (operator, batch, blank, web app)
- `hull dev` — watch the package and re-render on every save, like `webpack-dev-server` for Kubernetes
- `hull lint` — fast static validation before pushing to CI
- `hull template` — render locally, no cluster contact, see what would be applied
- `hull diff` — show exactly what would change on the next upgrade
- `hull test` — run smoke tests against a deployed release
- `hull config` — interactive walker that fills in `values.yaml` from the schema

**Workflow recipes:**

- **Local kind / k3d loop** — `hull dev ./my-app -n dev --interval 1s` re-renders on every file save; `kubectl apply` in another terminal applies the output.
- **Pre-commit lint** — wire `hull lint .` into a pre-commit hook so malformed packages never land in `main`.
- **One-command smoke test** — `hull install my-app . -n test --create-namespace && hull test my-app -n test`.

→ Start with: [Quickstart](guides/quickstart.md), [`hull create`](cli/create.md), [`hull dev`](cli/dev.md).

## For CI / release engineering

You automate deploys. You need predictability, integrity, and rollback paths.

**What hull gives you:**

- `hull plan` / `hull apply` — separates rendering from cluster contact; signed artifacts in git
- `hull diff --revision N` — compare the proposed render against any historical revision
- `hull rollback` — automated rollback on canary failure
- `hull multi-install --atomic-cross-cluster` — deploy to multiple clusters with rollback semantics
- `hull canary --stages 1,3,5 --bake 5m` — staged rollouts with health gating
- Plan integrity check — the SHA-256 of the rendered manifest is recorded; `hull apply` re-renders and verifies before applying
- `hull sbom <release>` — emit a CycloneDX 1.5 SBOM for compliance

**Workflow recipes:**

- **Plan-and-apply pipeline** — CI runs `hull plan` and posts the diff as a PR comment. Merging the PR triggers `hull apply` against prod.
- **Promotion through environments** — `hull plan --env staging`, review, `hull apply`. Then bump the package's appVersion and repeat for prod.
- **Canary in CI** — after merge to main, CI runs `hull canary --stages 10,50,100 --bake 5m -n prod`. Failed stages auto-roll back.

→ Start with: [`hull plan`](cli/plan.md), [`hull apply`](cli/apply.md), [`hull canary`](cli/canary.md), [`hull sbom`](cli/sbom.md).

## For security and compliance teams

You audit. You sign. You verify. You want forensic trails.

**What hull gives you:**

- PGP-signed `.prov` provenance files
- Cosign-attached OCI signatures (verify with the standard cosign workflow before install)
- A local PGP keyring (`hull keyring add/list/remove`) the operator owns
- `hull install --verify` to fail closed on missing or wrong signatures
- Audit data on every release record (who, when, where, what flags, what values)
- `hull sbom` — CycloneDX 1.5 SBOM emission per release
- `managedBy=hull` label on every applied resource and namespace
- Render-time network policy: HTTP/Vault calls disabled by default; SSRF-proof dial layer (loopback, link-local, RFC1918, metadata-IP all blocked)

**Workflow recipes:**

- **Signed-only install policy** — operators add only the keys you sign with; any package not signed by an admitted key fails `hull install --verify`.
- **Air-gapped install** — pull and mirror packages once; install offline. Render-time network knobs default off.
- **Forensic audit** — `hull audit <release>` reads the audit trail; `hull get manifest <release> --revision N` reproduces what was applied at any point.

→ Start with: [Signing](guides/signing.md), [`hull keyring`](cli/keyring.md), [`hull audit`](cli/audit.md), [`hull sbom`](cli/sbom.md).

## For teams switching from Helm

You're using Helm today. The pain points are typical: go-template surprises, brittle umbrella charts, the values-`<env>`.yaml proliferation, no native drift detection, sparse audit trail, manual multi-cluster rollouts. You want a **modern Helm alternative** without throwing away your existing chart investment.

**What hull gives you:**

- A built-in **Helm chart converter**: `hull migrate ./chart -d ./pkg` translates `Chart.yaml`, `templates/`, `_helpers.tpl`, and dependencies into a hull package; go-template constructs become hull `${...}` expressions where the migrator can do it cleanly.
- A **Helm interop layer**: `hull helm-compat install my-app /path/to/upstream-chart` runs upstream Helm charts under hull's release record without converting templates — you get hull's drift detection, audit trail, and reconcile model around the unmigrated chart.
- **Side-by-side coexistence**: hull's `managedBy=hull` label and `hull.v1.*` Secret naming don't collide with Helm's `app.kubernetes.io/managed-by=Helm` and `sh.helm.release.v1.*`, so a single cluster can host both during a phased migration.
- **The features Helm doesn't have natively**: drift detection (`hull drift`), reconcile (`hull reconcile`), per-revision audit (`hull audit`), multi-cluster atomic deploys (`hull multi-install --atomic-cross-cluster`), plan/apply with integrity hashing, workspace orchestration with health-gating.

**Workflow recipes:**

- **Phased migration starting with leaf charts.** Migrate a single small chart first (`hull migrate ./small-chart -d ./small-pkg`); deploy alongside existing Helm releases; build team familiarity; expand outward.
- **Wrap Helm charts you can't fork.** For upstream charts (cert-manager, kube-prometheus-stack, ingress-nginx) where you want hull's release semantics but don't want to maintain a fork, use `hull helm-compat install`. You keep upstream's update cadence and gain hull's audit / drift / reconcile loop.
- **Run drift detection across mixed Helm + hull releases.** `hull helm-compat list` enumerates Helm releases the compat layer manages; `hull drift` against each shows what drifted. Combine with `hull list -A` for a unified view of every managed release in the cluster.

→ Start with: [Helm-to-hull migration guide](guides/migration.md), [`hull migrate`](cli/migrate.md), [`hull helm-compat`](cli/helm-compat.md), [Hull as a Helm alternative](comparison.md#hull-as-a-helm-alternative).

## Cross-cutting concerns

### Multi-tenant clusters

Operators with `create/get` on `HullRelease` CRs can install via `hull controller` while platform admins gate package paths via `--package-root`. Tenants cannot point the controller at host paths because the resolved package must live under the allowlist.

### Air-gapped / sovereign environments

Pull once with `hull pull`, mirror archives to a local OCI registry, install offline. Render-time network calls (HTTP/Vault/SOPS) are opt-in and default off.

### Compliance-driven environments (PCI, HIPAA, SOC 2)

Hull's audit trail, signed packages, and SBOM generation cover the deployment-side controls these frameworks require. Combine with `hull policy` for in-package policy enforcement (hull-native declarative rules).

## Where next

- [Quickstart](guides/quickstart.md) — first install
- [FAQ](faq.md) — common questions
- [Glossary](glossary.md) — terminology
- [Comparison with other tools](comparison.md)
