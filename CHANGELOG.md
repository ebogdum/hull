# Changelog

All notable changes to hull are documented in this file.
The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.0.0] — first stable release

The first stable release of hull — a Kubernetes package manager, YAML
templating engine, and drift-detection CLI written in Go.

### Template engine

- YAML-native expression engine with `${...}` substitution and
  pipeline syntax (`${values.x | upper | quote}`).
- Lowercase root namespaces: `values`, `release`, `package`,
  `capabilities`, `files`. Capitalised sub-keys on capabilities
  (`kubeVersion.Major`/`.Minor`/`.GitVersion`).
- YAML-level control-flow directives: `$if`/`$then`/`$else`,
  `$each`/`$as`/`$yield`, `$switch`/`$cases`/`$default`.
- `$include` partials with depth-bounded recursion and a
  `_helpers.yaml` convention.
- Hooks declared via `$hook` directive with `$weight`,
  `$delete-policy`, and per-hook `$timeout`.
- Files API: `files.Get`, `files.Glob`, `files.AsConfig`,
  `files.AsSecrets`, `files.Lines`.
- Per-render function registry — context-bound `tpl`, `lookup`,
  `include`, and `files.*` closures wrapped over a stateless base
  registry, safe under parallel rendering.
- ~180 built-in functions covering string, math, regex, date, crypto,
  encoding, type, collection, logic, path, secret, and Sprig-extra
  surfaces. RSA key generation clamped to a sane bit range; range,
  repeat, regex, and rand-* helpers carry allocation caps.
- `tpl` recursion guard and `$include` depth limit.
- `lookup` with per-render cache and automatic mapper-refresh on
  miss; live-cluster reads gated by an explicit context flag.
- Strftime → Go layout converter for `date` / `dateInZone`.
- Truthy/falsy rules cover string forms (`""`, `"false"`, `"0"`,
  `"no"`, case variants) so values round-tripped through env vars
  and `--set` keep their semantics.

### Package model

- `hull.yaml` schema (`apiVersion`, `name`, `version`, `appVersion`,
  `kubeVersion`, `layers`, `requires`, `environments`, `profiles`,
  `immutables`, `dependencies`).
- `values.yaml` with deep-merge semantics across layers and `--set`.
- `values.schema.json` validation: `pattern`, `format`, `enum`,
  `const`, `multipleOf`, `min`/`max`/`exclusive*`, `minItems` /
  `maxItems` / `uniqueItems`, `minProperties` / `maxProperties`,
  `dependentRequired`, `patternProperties`, `allOf`/`anyOf`/`oneOf`/
  `not`, `$ref`/`$defs`/`definitions` with cycle detection.
- `_helpers.yaml` partials.
- `crds/` directory applied as a separate phase, waiting on
  `Established=true` before the regular manifest applies.
- `tests/` directory for `hull test`.
- `files/` directory accessible from templates via the Files API.
- `notes.yaml` post-install message.
- `profiles/` for named overlays.
- `policies/` for in-package policy enforcement (Rego or
  hull-native rules).
- `.hullignore` exclusion file.

### Layers and composition

- `layers:` for in-release composition; `requires:` for separate
  releases.
- Layer sources: local path, HTTPS, git, OCI.
- Tag-based and condition-based layer enablement, evaluated against
  fully-merged values (parent overrides applied).
- Parent override of nested layer keys via `layers.<name>.<key>`.
- Forced precedence with the `!` suffix.
- `hull.lock` auto-generation for reproducible resolution.

### Environments and profiles

- `environments:` block in `hull.yaml` with `inherits:` chains,
  per-environment values, value-files, profile, namespace, and
  cluster context.
- Profiles independent of environments, applicable to any release.

### Release lifecycle

- `hull install` / `hull upgrade` (with `--install`,
  `--reset-values`, `--reuse-values`, `--only`, `--wait`,
  `--wait-for-jobs`, `--atomic`, `--cleanup-on-fail`, `--force`).
- `hull rollback` — re-applies a previous revision and re-runs that
  revision's stored hooks.
- `hull uninstall` (with `--keep-history`, `--no-hooks`).
- `hull list` (`-A`, `--selector`, `--filter`, `--date`).
- `hull status`, `hull history`, `hull get`
  (`manifest`/`values`/`hooks`/`notes`/`metadata`/`all`,
  `--revision`, `--template`).
- `hull audit` — full chronological release trail (who, when, where,
  flags, values).
- `hull diff` — smart, server-side, per-resource diffing against any
  historical revision.
- `hull plan` — deterministic rendered manifest with SHA-256
  integrity hash.
- `hull apply` — re-renders, verifies the plan integrity hash, then
  applies.
- `hull drift` — per-field comparison of stored manifest vs live
  state with smart filtering.
- `hull reconcile` — converges cluster state back to the stored
  manifest.
- `hull prune` — drops superseded revisions.
- `hull rename` — copies revisions to a new release name.
- `hull canary` — staged rollouts with bake periods between replica
  counts and automatic rollback.
- `hull multi-install` — fleet-wide rollouts with
  `--atomic-cross-cluster` rollback semantics.

### Storage drivers

- Secret driver (default; `hull.v1.<release>.v<rev>` naming).
- ConfigMap driver.
- Memory driver (in-process, with deep-clone isolation between
  reads and writes).
- SQL driver (PostgreSQL, SQLite, MySQL) with
  schema-versioned migrations executed under transactions.
- Connection-string env-var configuration for the SQL driver.

### Workspaces and cross-source orchestration

- `hull-workspace.yaml` with members, defaults, and `dependsOn`.
- Kahn-level scheduling with parallel execution within a level.
- Health-gating between levels.
- Atomic-workspace mode with cross-cluster rollback.
- Continue-on-error mode for best-effort rollouts.
- `hull-releases.yaml` for cross-source release groups.

### Distribution and registry

- `hull package` — builds `*.hull.tgz` archives with optional
  `--sign`, `--app-version`, and `--version` overrides.
- `hull repo add` / `list` / `remove` / `update` / `index` for
  HTTPS chart-style repositories.
- `hull pull` — HTTPS or OCI, with `--version` constraint, `--prov`
  provenance download, and TLS material flags
  (`--ca-file`/`--cert-file`/`--key-file`).
- `hull push`.
- `hull publish` — HTTP API or OCI, with bounded response bodies
  and per-call timeouts.
- `hull registry login` / `logout` / `push` / `pull`, with
  `--insecure-skip-tls-verify` and `--plain-http`.
- `hull search` — Artifact Hub plus repo search, with
  `--endpoint`, `--kind`, and repository-URL output.
- `hull show` — chart, values, readme, crds, all (with bounded
  archive entry sizes).
- OCI distribution via ORAS.
- Authenticated client with retries, redirects, mTLS, cross-host
  redirect blocking, and refusal of plaintext-HTTP Basic Auth
  unless explicitly opted in.

### Signing and verification

- `hull package --sign` — PGP `.prov` provenance file emission.
- `hull keyring add` / `list` / `remove` — operator-controlled
  trust store with PGP public-key validation on add.
- `hull verify` — standalone verification.
- `hull install --verify` — fail-closed install gated on signature.
- Cosign integration — verify external cosign signatures before
  install.

### Cluster operations

- Server-side apply with `FieldManager=hull`.
- Deterministic install order (`CustomResourceDefinition` first,
  then namespaces, RBAC, and the rest).
- `WaitForReady` for `Deployment`, `StatefulSet`, `DaemonSet`,
  `Job`, `CronJob`, `Pod`, `ReplicaSet`, `ReplicationController`.
- Permanent-failure detection for `ProgressDeadlineExceeded` and
  `ReplicaFailure` so waits don't hang on broken rollouts.
- `managedBy=hull` label stamped on every applied resource and on
  the namespace itself.
- Stamped pod templates inside `Deployment`, `StatefulSet`,
  `DaemonSet`, `ReplicaSet`, `ReplicationController`, `Job`,
  `CronJob` so embedded controllers are queryable cluster-wide.
- Capabilities API (`kubeVersion`, `apiVersions.Has`).
- `helm.sh/resource-policy: keep` honoured across all delete paths
  (uninstall, upgrade `--force`, atomic cleanup).

### Controller and CRDs

- `HullRelease` Custom Resource Definition.
- `hull controller` — in-cluster reconciler with package-root
  allowlist, per-cycle eviction of deleted CRs, and structured
  status logging.

### Tooling

- `hull create` — scaffold a working package.
- `hull init` — template-based scaffold (operator, batch, blank,
  web app).
- `hull lint` — schema, template, manifest, and best-practice
  validation.
- `hull template` — local rendering, no cluster contact.
- `hull dev` — watch and re-render on every save.
- `hull config` — interactive walker over `values.schema.json`.
- `hull values --trace` — show the resolved value source per key.
- `hull policy` — run policies against rendered manifests.
- `hull scan` — find common values across packages and extract a
  base layer.
- `hull sbom` — CycloneDX 1.5 SBOM emission per release.
- `hull metrics` — sample CPU and memory of running workloads and
  recommend `requests` / `limits`.
- `hull graph` — dependency graph rendering.
- `hull adopt` — claim existing resources into a release record.
- `hull migrate` — Helm-chart-to-hull-package converter.
- `hull helm-compat` — interop layer that runs unmodified upstream
  Helm charts under hull's release record.

### Plugins and marketplace

- `hull plugin install` / `list` / `remove` / `update` from git URL
  or local path, with split install / update / delete hooks.
- Plugin-as-extension subcommands.
- `hull marketplace search` / `verify` against a signed plugin
  index.
- `plugin.yaml` schema validation (known-fields-only decode).
- Cross-platform plugin runner (POSIX shell on Unix, `cmd.exe` on
  Windows).
- Compound URL-scheme matching (`s3+https`, `git+ssh`).
- TLS-material plumbing for downloaders.

### Operations and forensics

- `hull purge` — cluster-wide cleanup driven by the
  `managedBy=hull` selector (with legacy-label fallback).
- `hull purge --force` — drains stuck namespaces, force-deletes
  pods (two-pass with `grace=0`), force-finalises terminating
  namespaces, and sweeps orphan pods left behind by node failure.
- `hull purge --delete-namespaces` for full namespace teardown.
- Status redaction for secret-shaped strings in CLI output.
- `hull test` with `--parallel`, `--retries`, and per-test
  defer-based cleanup.

### Security

- SSRF-proof dial layer — blocks loopback, link-local, RFC1918,
  CGNAT, and cloud metadata IPs.
- TLS 1.2 minimum across all outbound connections.
- Tar extraction protections: path-traversal rejection, entry-count
  cap, per-entry size cap, file-mode sanitisation, symlink refusal.
- Migration safety: symlink rejection on input, refusal to
  overwrite a non-empty output directory unless explicitly
  forced.
- Plan-time path/action validation.
- PGP keyring with operator-controlled trust and public-key
  validation on add.
- Render-time HTTP / Vault calls opt-in via env var, default off.
- Per-call timeouts on OCI push/pull, marketplace, and publish
  paths; bounded error-response body reads.
- Resource caps on `repeat`, `indent`, `nindent`, `until`,
  `untilStep`, `randAlphaNum`/`randAscii`/`randNumeric`/
  `randAlpha`, and regex find-all/split.
- Function-argument literal-only model (no path injection from the
  argument side of a pipeline).
- Plugin-install argument hardening (dash-prefix rejection,
  `--` separator).
- Cross-host redirect blocking and refusal of plaintext-HTTP Basic
  Auth without an explicit opt-in.

### Documentation

- MIT licence.
- Comprehensive README with quickstart, feature overview, and
  Helm-alternative comparison.
- ~100 per-command CLI reference pages.
- Per-template-field reference (expressions, control flow,
  functions, capabilities, layers, hooks).
- Format references for `hull.yaml`, `values.yaml`,
  `values.schema.json`, `hull-workspace.yaml`, `hull-releases.yaml`.
- Guides: quickstart, layers, workspaces, signing, migration from
  Helm, GitOps integration.
- FAQ, glossary, troubleshooting, comparison, and use-case pages.

[1.0.0]: https://github.com/ebogdum/hull/releases/tag/v1.0.0
