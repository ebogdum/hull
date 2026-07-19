# hull completion

`hull completion` prints a shell completion script for `bash`, `zsh`, `fish`,
or `powershell` to stdout, so your shell can tab-complete hull commands and
flags.

## When to use it

- Once per shell, to turn on tab completion for `hull`.
- After upgrading hull, to refresh the script if new commands or flags were
  added.

## What happens

1. You name a shell. hull writes that shell's completion script to stdout.
2. You either source the output in the current session or save it into the
   shell's completion directory so it loads automatically on every new shell.
3. Nothing is installed or contacted — the command just prints a script.

## Usage

```
hull completion [bash|zsh|fish|powershell]
```

## Flags

Inherits the global flags (`--debug`, `--kube-context`, `--kubeconfig`,
`-n/--namespace`); completion only prints a script, so they have no effect.

## Worked example

Pick the block for your shell.

**Bash** — load it for the current session, or install it permanently:

```sh
# current session only
source <(hull completion bash)

# every new shell (Linux)
hull completion bash | sudo tee /etc/bash_completion.d/hull > /dev/null
```

**Zsh** — ensure completion is enabled, then drop the script on `fpath`:

```sh
# once, if not already in ~/.zshrc
echo "autoload -U compinit; compinit" >> ~/.zshrc

hull completion zsh > "${fpath[1]}/_hull"
```

**Fish** — save it into the fish completions directory:

```sh
hull completion fish > ~/.config/fish/completions/hull.fish
```

**PowerShell** — load it for the session, or append it to your profile:

```powershell
# current session only
hull completion powershell | Out-String | Invoke-Expression

# every new session
hull completion powershell >> $PROFILE
```

Start a new shell (or re-source your profile) and `hull <Tab>` completes.

## See also

- [`env`](env.md) — the resolved paths and environment hull is using
- [`version`](version.md) — confirm the build you generated completions from
