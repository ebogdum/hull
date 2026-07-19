# hull adopt

`hull adopt` claims resources that already exist in the cluster and records
them as revision 1 of a new hull release, so hull can manage them from then on.

## When to use it

- You have live resources created by `kubectl apply`, Terraform, or a
  hand-written manifest, and you want hull to track them.
- You are migrating a running deployment onto hull without deleting and
  recreating it.
- You want `hull diff`, `hull drift`, `hull upgrade`, and `hull uninstall` to
  work on resources hull never installed.

## What happens

1. Parses each `<resource-ref>` into a kind / namespace / name lookup.
2. Connects to the cluster in the current context.
3. If `--create-namespace` is set, creates the release namespace when missing.
4. Fetches each referenced object from the cluster and strips server-managed
   metadata (`status`, `managedFields`, resource version).
5. Stores the cleaned resources as **revision 1** of the release, recording
   any `--description` and `--labels` in the audit trail.
6. Prints how many resources were adopted.

Mutating: it writes new release state. The adopted resources themselves are
not changed. Requires a reachable cluster.

## Usage

```
hull adopt <release-name> <resource-ref>... [flags]
```

Each `<resource-ref>` takes one of these forms:

```
apps/v1/Deployment/myns/myapp        group/version/Kind/namespace/name
v1/ConfigMap//cluster-scoped-cm      empty namespace for cluster-scoped objects
kind=Deployment,name=myapp,ns=myns   key=value form
```

## Flags

| Flag | Type | Default | Description |
|---|---|---|---|
| `--create-namespace` | — | false | create the release namespace first, so adoption does not fail when it is missing |
| `--description` | string | — | text recorded against revision 1 in the audit trail |
| `--labels` | stringArray | — | `key=value` label attached to the release; repeat for more |

## Persistent flags inherited from `hull`

| Flag | Type | Description |
|---|---|---|
| `--debug` | — | enable debug output |
| `--kube-context` | string | Kubernetes context to use |
| `--kubeconfig` | string | path to kubeconfig file |
| `-n, --namespace` | string | Kubernetes namespace |

## Worked example — bring an unmanaged Deployment under hull

**INPUT — a Deployment someone created with raw `kubectl`,** so hull has no
record of it:

```sh
kubectl get deploy legacy-api -n apps
# NAME         READY   UP-TO-DATE   AVAILABLE
# legacy-api   3/3     3            3

hull list -n apps
# (legacy-api is absent — hull never installed it)
```

**Adopt it into a release named `legacy-api`:**

```sh
hull adopt legacy-api apps/v1/Deployment/apps/legacy-api \
  -n apps --description "onboarded from kubectl"
```

**OUTPUT:**

```
Adopted 1 resource(s) as release "legacy-api" (revision 1, namespace apps).
```

**Tracing the result:**

| Output | Cause |
|---|---|
| `Adopted 1 resource(s)` | one resource-ref was fetched and stored |
| `revision 1` | adoption always creates the release at revision 1 |
| now listed by `hull list -n apps` | the release state is recorded, so `hull diff`, `hull drift`, and `hull upgrade` work |

The Deployment keeps running untouched; only hull's stored state is new. From
here `hull drift legacy-api` compares that stored state against the live
object.

## See also

- [`drift`](drift.md) — compare a release's package, state, and cluster
- [`diff`](diff.md) — compare packages, manifests, or revisions
- [`upgrade`](upgrade.md) — roll a new package version onto the release
- [`uninstall`](uninstall.md) — remove an adopted release
