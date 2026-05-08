# hull completion

## Synopsis

`hull completion` generates shell-completion scripts for `bash`, `zsh`, `fish`, or `powershell`. The generated script defines completion handlers; source it from your shell init or stage it into the shell's standard completion directory.

## When to use it

Use once per shell to enable tab completion of hull commands and flags.

## Usage

```
hull completion [bash|zsh|fish|powershell] [flags]
```

## Flags

| Flag | Type | Default | Description |
|---|---|---|---|
| `-h, --help` | — | — | help for completion |

## Persistent flags inherited from `hull`

| Flag | Type | Description |
|---|---|---|
| `--debug` | — | enable debug output |
| `--kube-context` | string | Kubernetes context to use |
| `--kubeconfig` | string | path to kubeconfig file |
| `-n, --namespace` | string | Kubernetes namespace |

## Examples

Bash:

```sh
hull completion bash > /etc/bash_completion.d/hull
```

Zsh:

```sh
hull completion zsh > "${fpath[1]}/_hull"
```

Fish:

```sh
hull completion fish > ~/.config/fish/completions/hull.fish
```
