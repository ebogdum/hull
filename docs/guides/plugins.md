# Plugins

A plugin adds a new top-level command to hull. Once installed, `hull greet` or
`hull backup` runs your code as if it were built in.

## What a plugin actually is

A hull plugin is **a directory** containing:

- a `plugin.yaml` manifest, and
- an executable (or a shell one-liner) that the manifest points at.

That's the whole contract. Some specifics, because they're the questions every
plugin author asks first:

- **Does it have to be Go?** No. The plugin can be **any executable** — a bash
  or Python script, a compiled binary in any language, anything your OS can run.
  hull itself is Go, but it never loads your code into its own process.
- **Is it linked in / a Go plugin (`.so`)?** No. hull does **not** use Go's
  `plugin` package. Your plugin is a **separate program**.
- **Is it a pipe / RPC / IO protocol?** No. hull runs your program as a normal
  **child process** with `exec`, and wires its **stdin, stdout, and stderr
  straight to the terminal**. You read stdin and print to stdout exactly as you
  would in any CLI. There is no framing, no JSON-RPC, no socket.
- **Is it attached / long-running?** No. hull execs it, waits, and returns its
  **exit code**. When your program exits, the command is done.

This is the same model git, kubectl, and Helm use for their plugins — and hull
sets Helm's environment variables too, so many existing Helm plugins run
unchanged.

## Build your first plugin

A plugin is a directory. Create one:

```sh
mkdir greet && cd greet
```

Write the manifest, `plugin.yaml`:

```yaml
name: greet
version: 0.1.0
usage: "greet [name]"
description: Print a friendly greeting
command: greet.sh
```

Write the executable it points at, `greet.sh`:

```sh
#!/usr/bin/env bash
echo "Hello, ${1:-world}! (plugin dir: $HULL_PLUGIN_DIR)"
```

Install it from the local directory:

```sh
hull plugin install ./greet
```

```
Installed plugin: greet v0.1.0
```

Now run it — everything after the plugin name is passed to your program as
arguments:

```sh
hull greet Ada
```

```
Hello, Ada! (plugin dir: /home/you/.config/hull/plugins/greet)
```

```sh
hull greet
```

```
Hello, world! (plugin dir: /home/you/.config/hull/plugins/greet)
```

`hull greet Ada` ran `greet.sh` with `Ada` as `$1`; with no argument, the
script's own `${1:-world}` default produced `world`. hull added
`HULL_PLUGIN_DIR` to the environment (see below).

## The `plugin.yaml` manifest

| Field | Type | Required | Description |
|---|---|---|---|
| `name` | string | yes | The command users type: `hull <name>`. Must match `[a-zA-Z0-9][a-zA-Z0-9._-]*`. |
| `command` | string | yes | What to run. Either a filename in the plugin dir, or a shell expression (see below). |
| `version` | string | no | Shown in `hull plugin list`. |
| `description` | string | no | Shown in `hull plugin list`. |
| `usage` | string | no | Short usage line for your own help text. |
| `hooks` | map | no | `install` / `update` / `delete` shell commands run on lifecycle events. |
| `downloaders` | list | no | Custom protocol handlers (see [Custom downloaders](#custom-downloaders)). |

Unknown top-level keys are **rejected at install time** (strict parsing), so a
typo in a field name fails loudly instead of being silently ignored.

## The `command` field

`command` has two shapes, and hull picks automatically:

1. **A filename** — a single token naming a file **inside the plugin directory**,
   e.g. `command: greet.sh`. hull runs that file directly, appending the user's
   arguments. This is the common case. hull makes the file executable
   (`chmod 0755`) on install.

2. **A shell expression** — anything containing whitespace or shell characters
   (`| & ; < > $ ( ) * ?`), e.g. `command: "echo hello"` or
   `command: "python3 $HULL_PLUGIN_DIR/main.py"`. hull runs it through
   `/bin/sh -c` (`cmd /C` on Windows), so pipes, `$HULL_PLUGIN_DIR`, and
   multi-word commands all work.

For safety, a plain filename must resolve to a real file in the plugin dir; it
**cannot** contain a path separator or `..`, and a symlink is refused. A
filename that doesn't exist is an error (it catches typos) rather than silently
falling through to the shell.

## What hull passes your plugin

**Arguments.** Everything after `hull <name>` is passed to your program in
order. `hull greet Ada --loud` calls your program with `Ada` and `--loud`.

**Standard streams.** stdin, stdout, and stderr are the terminal's. Read input,
print output, and set your exit code normally.

**Environment variables.** hull adds these on top of the inherited environment:

| Variable | Value |
|---|---|
| `HULL_PLUGIN_DIR` | Absolute path to your plugin's installed directory. |
| `HULL_BIN` | Path to the hull executable — call back into hull with `"$HULL_BIN" template .`, etc. |
| `HULL_NAMESPACE` | The active namespace (`-n` / `HULL_NAMESPACE` / `HELM_NAMESPACE`). |
| `HULL_KUBECONFIG` | The active `KUBECONFIG` path. |
| `HELM_PLUGIN_DIR`, `HELM_BIN`, `HELM_NAMESPACE` | Helm-compatible aliases of the above, so Helm plugins work. |

Your program runs with its working directory set to the plugin directory.

## Lifecycle hooks

Declare shell commands that run on install, update, and remove — useful for
fetching a binary, compiling, or cleaning up:

```yaml
name: backup
version: 1.0.0
description: Back up a release
command: backup
hooks:
  install: "go build -o backup ./cmd"
  update:  "go build -o backup ./cmd"
  delete:  "rm -f backup"
```

Each hook runs through `/bin/sh -c` (`cmd /C` on Windows) with the working
directory set to the plugin directory and the same environment variables listed
above. The `install` hook runs **after** the files are in place; if it fails,
hull removes the half-installed plugin. hull logs the exact hook command before
running it, because a hook is the plugin author's code running with your
privileges.

## Custom downloaders

A plugin can teach hull to fetch packages over a custom protocol. Declare one or
more downloaders:

```yaml
downloaders:
  - command: myproto-downloader
    protocols: ["myproto"]
```

hull invokes the downloader as:

```
<command> <cert> <key> <ca> <full-url>
```

where the first three are TLS material (may be empty) and the last is the URL
the user requested. The downloader writes the package bytes to stdout.

## Installing and distributing

Install from a **local directory** (copied into place; symlinks are refused) or
from a **git repository** (`hull plugin install` runs `git clone --depth 1`):

```sh
hull plugin install ./greet                              # local
hull plugin install https://github.com/you/hull-greet.git   # git
```

Recognized git sources end in `.git` or start with `git@`, `git://`, `ssh://`,
or `file://`. Everything else is treated as a local path.

Plugins install under `~/.config/hull/plugins/<dir>`. Manage them with:

```sh
hull plugin list             # NAME / VERSION / DESCRIPTION
hull plugin update greet     # git plugins: git pull --ff-only, then update hook
hull plugin remove greet     # runs the delete hook, then deletes the directory
```

To publish for others, push the plugin directory (manifest + executable, or a
manifest plus an `install` hook that builds it) to a git repo and share the
clone URL. Signed plugins can be distributed through the
[marketplace](../cli/marketplace.md).

## Security

Installing a plugin means running its author's code on your machine — the
`command`, any shell-form command, and every lifecycle hook execute with **your**
privileges. hull reduces surprises (strict manifest parsing, no path traversal
in `command`, no symlinks copied in, hooks logged before they run), but it
cannot make untrusted code safe. Install plugins only from sources you trust,
and prefer [signed marketplace plugins](../cli/marketplace-verify.md).

## See also

- [`hull plugin`](../cli/plugin.md) — the plugin management commands
- [`hull plugin install`](../cli/plugin-install.md) · [`list`](../cli/plugin-list.md) · [`update`](../cli/plugin-update.md) · [`remove`](../cli/plugin-remove.md)
- [`hull marketplace`](../cli/marketplace.md) — discover and verify signed plugins
