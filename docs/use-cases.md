# Hull use cases — for platform engineers, SREs, and GitOps teams

This page maps hull's features onto real workflows. Pick the role or pattern
that matches your team and follow the links.

## For platform engineers

You're building an internal developer platform on Kubernetes. Application teams
need a predictable way to ship workloads. Compliance wants signed artifacts and
audit trails. Operations wants drift detection.

**What hull gives you:**

- A versioned, signable package format (`*.hull.tgz`) capturing everything an
  app needs to run.
- A composition system (layers) so platform-provided base layers (logging,
  RBAC, network policies) flow into every app package.
- Per-release audit data (who, when, where, with what flags) baked into the
  cluster's release record.
- An OCI distribution model that reuses your existing registry and IAM.
- A workspace orchestrator for multi-component rollouts.
- Cluster-wide ownership queries via one label selector (`managedBy=hull`).

**Workflow recipes:**

- **Standard library of layers** — publish base layers (`org-base-rbac`,
  `org-base-monitoring`) to OCI; consumers pull them in via `layers:` in
  `hull.yaml`. Updates propagate by bumping the version constraint and running
  `hull dependency update`.
- **Multi-environment promotion** — define `environments:` (`dev`/`staging`/
  `prod`) in `hull.yaml` so the same package flows through environments. CI runs
  `hull plan --env staging -o staging.plan`, then `hull apply --plan
  staging.plan` after review.
- **Self-service signed packages** — teams run `hull package --sign` against
  CI's PGP key; each cluster's keyring admits only that key, so tampered
  packages fail `hull install --verify`.
- **Auditable rollouts** — `hull audit <release>` answers "who upgraded the
  auth service yesterday?" months later.

→ Start with: [Quickstart](guides/quickstart.md), [Layers](guides/layers.md),
[Workspaces](guides/workspaces.md), [Signing](guides/signing.md).

## For SRE / operations teams

You operate clusters. You need to know what's deployed, what's drifted, and
that rollback works.

**What hull gives you:**

- `hull list -A` — every release in every namespace.
- `hull drift ./pkg` — three-way compare of package, recorded state, and live
  cluster, with cluster-injected noise filtered out.
- `hull reconcile <release>` — converge cluster state back to the stored
  manifest by release name (no package source needed).
- `hull audit <release>` — full chronological history.
- `hull rollback <release> <rev>` — re-apply a previous revision and re-run its
  hooks.
- `hull metrics <release>` — sample CPU/memory and recommend requests/limits.
- `hull multi-install --to <ctx-list> --atomic-cross-cluster` — fleet-wide
  rollouts with rollback.
- `hull canary --stages 1,3,5 --bake 5m` — staged upgrades with bake periods.

**Workflow recipes:**

- **Drift-detection cron** — enumerate releases with `hull list -A -q`, then run
  `hull drift ./pkg -r <release>` per package (the check renders the package, so
  the source tree must be available); alert on any non-empty divergence.
- **Incident rollback drill** — practise `hull rollback <release> <prev-rev>`
  against staging before the next on-call. The behaviour is identical in prod.
- **Capacity right-sizing** — run `hull metrics <release>` after a release has
  baked for 24 hours; commit the recommended `requests` / `limits` to the
  package's `values.yaml`.
- **Force-cleanup after node failure** — `hull purge --yes --force` clears
  wedged releases and force-finalises stuck `Terminating` namespaces.
- **Ordered platform teardown** — a `hull-releases.yaml` with `dependsOn`
  brings the platform up (`hull releases install`) and tears it down in reverse
  (`hull releases uninstall`).

→ Start with: [`hull drift`](cli/drift.md), [`hull reconcile`](cli/reconcile.md),
[`hull audit`](cli/audit.md), [`hull canary`](cli/canary.md),
[`hull purge`](cli/purge.md).

## For GitOps teams (Argo CD, Flux)

You declare desired state in git and reconcile it into the cluster. You want
hull's packaging story without giving up your reconciler.

**What hull gives you:**

- `hull template` produces raw, deterministic manifest YAML you can commit for
  a reconciler to sync.
- `hull plan -o app.plan` produces an apply-able JSON artifact (rendered
  manifest plus a SHA-256 digest); `hull apply --plan app.plan` re-renders,
  verifies the digest, and applies — hull's own integrity-checked pipeline.
- `hull controller` is an in-cluster HullRelease CR reconciler if you want hull
  to reconcile directly.
- The `hull` binary can be registered as an Argo CD Config Management Plugin.
- Flux's OCI source can pull hull packages pushed via `hull registry push`.

**Workflow recipes (Argo CD):**

- **CMP plugin** — register `hull template <package>` as a CMP that emits the
  rendered manifest. Argo CD diffs and syncs as usual.
- **Pre-rendered manifest in git** — CI runs `hull template ./app --env prod >
  manifests.yaml` and commits it to a path Argo CD watches.
- **HullRelease CR pattern** — commit `HullRelease` CRs; Argo CD syncs the CRs
  and `hull controller` reconciles them.

**Workflow recipes (Flux):**

- **OCI source** — `hull registry push` to OCI; Flux's OCIRepository pulls; a
  Flux Kustomization applies the rendered manifests.
- **HullRelease CR + hull controller** — Flux syncs CRs; `hull controller`
  reconciles.

→ Start with: [`hull template`](cli/template.md), [`hull plan`](cli/plan.md),
[`hull apply`](cli/apply.md), [`hull controller`](cli/controller.md).

## For application developers

You write code. You want to ship to Kubernetes without becoming a YAML expert.

**What hull gives you:**

- `hull create my-app` — scaffold a working package in seconds.
- `hull init <template>` — start from a richer built-in template.
- `hull dev ./my-app` — watch the package and re-render on every save.
- `hull lint` — fast static validation before pushing to CI.
- `hull template` — render locally, no cluster contact.
- `hull diff` — show exactly what changes between two packages, value sets, or
  git revisions.
- `hull test` — run smoke tests against a deployed release.
- `hull config` — interactive walker that fills a values file from the schema.

**Workflow recipes:**

- **Local kind / k3d loop** — `hull dev ./my-app -n dev --interval 1s`
  re-renders on every save; apply the output from another terminal.
- **Pre-commit lint** — wire `hull lint .` into a pre-commit hook so malformed
  packages never land in `main`.
- **One-command smoke test** — `hull install my-app . -n test
  --create-namespace && hull test my-app -n test`.

→ Start with: [Quickstart](guides/quickstart.md), [`hull create`](cli/create.md),
[`hull dev`](cli/dev.md).

## For CI / release engineering

You automate deploys. You need predictability, integrity, and rollback paths.

**What hull gives you:**

- `hull plan` / `hull apply --plan` — separate rendering from cluster contact,
  with a digest checked before apply.
- `hull diff --from-ref v1.2.0 --to-ref HEAD ./chart` — compare two git
  revisions of a package.
- `hull rollback` — automated rollback on failure.
- `hull multi-install --atomic-cross-cluster` — deploy to many clusters with
  rollback.
- `hull canary --stages 10,50,100 --bake 5m` — staged rollouts with health
  gating.
- `hull sbom <release>` — emit a CycloneDX 1.5 SBOM for compliance.

**Workflow recipes:**

- **Plan-and-apply pipeline** — CI runs `hull plan -o app.plan` and posts the
  change preview as a PR comment; merging triggers `hull apply --plan app.plan`
  against prod.
- **Compare against a shipped revision** — dump a historical manifest with
  `hull get manifest <release> --revision N -o yaml > old.yaml`, render the
  candidate with `hull template ./chart > new.yaml`, then `hull diff old.yaml
  new.yaml`.
- **Canary in CI** — after merge, CI runs `hull canary <release> ./chart
  --stages 10,50,100 --bake 5m -n prod`; a failed stage auto-rolls back.

→ Start with: [`hull plan`](cli/plan.md), [`hull apply`](cli/apply.md),
[`hull canary`](cli/canary.md), [`hull sbom`](cli/sbom.md).

## For security and compliance teams

You audit. You sign. You verify. You want forensic trails.

**What hull gives you:**

- PGP-signed `.prov` provenance files.
- Cosign-attached OCI signatures (verify with the standard cosign workflow
  before install).
- A local PGP keyring (`hull keyring add/list/remove`) the operator owns.
- `hull install --verify` to fail closed on missing or wrong signatures.
- Audit data on every release record (who, when, where, what flags, what
  values).
- `hull sbom` — CycloneDX 1.5 SBOM per release.
- `managedBy=hull` on every applied resource and namespace.
- Render-time network policy: HTTP/Vault calls off by default; the render-time
  dial layer blocks loopback, link-local, RFC 1918, and metadata IPs.

**Workflow recipes:**

- **Signed-only install policy** — operators add only the keys you sign with;
  any package not signed by an admitted key fails `hull install --verify`.
- **Air-gapped install** — pull and mirror packages once; install offline.
  Render-time network knobs default off.
- **Forensic audit** — `hull audit <release>` reads the trail; `hull get
  manifest <release> --revision N` reproduces exactly what was applied.

→ Start with: [Signing](guides/signing.md), [`hull keyring`](cli/keyring.md),
[`hull audit`](cli/audit.md), [`hull sbom`](cli/sbom.md).

## For teams switching from Helm

You use Helm today. The pain points are typical: go-template surprises, brittle
umbrella charts, `values-<env>.yaml` proliferation, no native drift detection,
sparse audit trail, manual multi-cluster rollouts. You want a Helm alternative
without throwing away your chart investment.

**What hull gives you:**

- A **Helm-chart converter**: `hull migrate ./chart -o ./pkg` translates
  `Chart.yaml`, `templates/`, `_helpers.tpl`, and dependencies; go-template
  constructs become hull `${...}` expressions where the migrator can do it
  cleanly, and the rest is printed as a conversion report.
- A **Helm interop layer**: `hull helm-compat install my-app
  /path/to/upstream-chart` runs upstream charts under a hull release record —
  you get hull's reconcile, rollback, and audit around the unmigrated chart.
- **Side-by-side coexistence**: hull's `managedBy=hull` label and `hull.v1.*`
  Secret naming don't collide with Helm's, so one cluster hosts both during a
  phased migration.
- **Features Helm lacks natively**: drift detection (`hull drift`), reconcile
  (`hull reconcile`), per-revision audit (`hull audit`), multi-cluster atomic
  deploys, plan/apply with integrity hashing, workspace orchestration with
  health gating.

**Workflow recipes:**

- **Phased migration from leaf charts.** Migrate a small chart first (`hull
  migrate ./small-chart -o ./small-pkg`), deploy alongside existing Helm
  releases, build familiarity, then expand.
- **Wrap charts you can't fork.** For upstream charts (cert-manager,
  kube-prometheus-stack, ingress-nginx) where you want hull's release semantics
  without maintaining a fork, use `hull helm-compat install`, then `hull audit`
  / `hull reconcile` / `hull rollback` by release name.
- **Unified view of mixed releases.** `hull helm-compat install` records a hull
  release, so compat-managed releases appear in `hull list -A` alongside native
  hull releases.

→ Start with: [Migration guide](guides/migration.md),
[`hull migrate`](cli/migrate.md), [`hull helm-compat`](cli/helm-compat.md),
[Hull as a Helm alternative](comparison.md#hull-as-a-helm-alternative).

## Cross-cutting concerns

### Multi-tenant clusters

Operators with `create/get` on `HullRelease` CRs install via `hull controller`
while platform admins gate package paths with `--package-root`. Tenants can't
point the controller at host paths because the resolved package must live under
the allowlist.

### Air-gapped / sovereign environments

Pull once with `hull pull`, mirror archives to a local OCI registry, install
offline. Render-time network calls (HTTP/Vault/SOPS) are opt-in and off by
default.

### Compliance-driven environments (PCI, HIPAA, SOC 2)

Hull's audit trail, signed packages, and SBOM generation cover the
deployment-side controls these frameworks require. Combine with `hull policy
check` for in-package policy enforcement.

## Where next

- [Quickstart](guides/quickstart.md) — first install
- [FAQ](faq.md) — common questions
- [Glossary](glossary.md) — terminology
- [Comparison with other tools](comparison.md)
