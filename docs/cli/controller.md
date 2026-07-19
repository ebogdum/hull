# hull controller

## Synopsis

`hull controller` runs hull as an in-cluster operator. Instead of you typing
`hull install` / `hull upgrade` from a workstation, you declare what you want
as `HullRelease` custom resources, and a long-running controller process
reconciles the cluster to match them.

The workflow has three parts, one per subcommand:

```
controller crd          — print the HullRelease CRD definition
controller install-crd  — register that CRD in the cluster
controller run          — start the reconcile loop that acts on HullReleases
```

Once the CRD is installed and the loop is running, you (or a GitOps engine)
create a `HullRelease` object like this:

```yaml
apiVersion: hull.dev/v1
kind: HullRelease
metadata:
  name: web
  namespace: apps
spec:
  package: web            # path under the controller's --package-root
  releaseName: web        # optional; defaults to metadata.name
  profile: prod           # optional
  values:                 # optional inline values
    replicas: 3
```

On its next tick the controller renders that package, installs or upgrades the
release, and writes the result back to the object's `status` (phase, message,
revision). Edit the `HullRelease` and the controller reconciles again; a CR it
has already applied and that has not changed is skipped.

## Subcommands

| Command | Description |
|---|---|
| [`crd`](controller-crd.md) | Print the HullRelease CRD YAML to stdout |
| [`install-crd`](controller-install-crd.md) | Apply the HullRelease CRD to the cluster |
| [`run`](controller-run.md) | Run the reconcile loop in the foreground |

## Usage

```
hull controller [command]
```

Bring the operator up in the usual order:

```sh
hull controller install-crd      # once per cluster
hull controller run              # long-running; deploy it as a Deployment
```

## See also

- [`install`](install.md) — the one-shot install the controller runs for you
- [`upgrade`](upgrade.md) — the one-shot upgrade the controller runs for you
- [`reconcile`](reconcile.md) — manually converge a single release
