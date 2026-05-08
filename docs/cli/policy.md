# hull policy

## Synopsis

`hull policy` evaluates package-defined policy rules (under `policies/`) against the rendered manifest. Rules can be Hull policy YAML (declarative match-and-require) or Rego (full Open Policy Agent expressiveness).

## When to use it

Use as a CI gate to enforce organisation-wide rules: "every Pod sets runAsNonRoot", "every Service has a selector", "no hostNetwork", etc. Policies live with the package so they ship together.

## Usage

```
hull policy [command]
```

## Subcommands

- [`hull policy list`](policy-list.md) — List policy rules declared in the package

## Flags

| Flag | Type | Default | Description |
|---|---|---|---|
| `-h, --help` | — | — | help for policy |

## Persistent flags inherited from `hull`

| Flag | Type | Description |
|---|---|---|
| `--debug` | — | enable debug output |
| `--kube-context` | string | Kubernetes context to use |
| `--kubeconfig` | string | path to kubeconfig file |
| `-n, --namespace` | string | Kubernetes namespace |

## Examples

Evaluate policies:

```sh
hull policy check ./my-app
```

List declared policies:

```sh
hull policy list ./my-app
```

## See also

- [`lint`](lint.md)
