# hull show readme

## Synopsis

`hull show readme` prints the contents of a package's `README.md` to stdout. The output is the human-facing documentation the author shipped — typically install notes, configuration guidance, and links to upstream resources.

## When to use it

Use to read a package's documentation without opening the file in an editor — useful when reviewing pulled packages from the terminal, or piping the README through a renderer like `glow` or `bat`.

## What happens when you run it

1. Reads `<package-path>/README.md`.
2. Prints the content to stdout, unchanged.
3. No cluster contact, no network access.

## Usage

```
hull show readme <package-path> [flags]
```

## Flags

| Flag | Type | Default | Description |
|---|---|---|---|
| `-h, --help` | bool | false | help for readme |

## Persistent flags inherited from `hull`

| Flag | Type | Description |
|---|---|---|
| `--debug` | bool | enable debug output |
| `--kube-context` | string | Kubernetes context to use |
| `--kubeconfig` | string | path to kubeconfig file |
| `-n, --namespace` | string | Kubernetes namespace |

## Examples

Show a local package's README:

```sh
hull show readme ./my-app
```

After pulling a package from a repo:

```sh
hull pull my-app --repo https://charts.example.com --version 1.2.3 -d ./pulled --untar
hull show readme ./pulled/my-app
```

Pipe through a Markdown renderer for nicer terminal display:

```sh
hull show readme ./my-app | glow -
```

## See also

- [`show`](show.md)
- [`show all`](show-all.md)
- [`show chart`](show-chart.md)
