---
title: "hull marketplace"
parent: "CLI"
---
{% raw %}
# hull marketplace

## Synopsis

`hull marketplace` helps you find signed plugins from a marketplace index and
check one before you trust it. A marketplace is a JSON index listing plugins,
their download URLs, and per-plugin signatures. The default index is
`https://plugins.hull.dev/index.json`; point `--index` at another one to use a
private marketplace.

The marketplace commands don't install anything themselves — once you've found
and verified a plugin, you install it with
[`hull plugin install`](plugin-install.md).

## Subcommands

| Command | What it does |
|---|---|
| [`search`](marketplace-search.md) | List plugins in a marketplace index, optionally filtered by keyword |
| [`verify`](marketplace-verify.md) | Check a downloaded plugin archive against the index's digest and signature |

## Usage

```
hull marketplace [command]
```

A typical flow: find a plugin, verify the archive you downloaded, then install
it.

```sh
hull marketplace search backup
hull marketplace verify --archive ./backup.tar.gz --name backup
hull plugin install ./backup
```

## See also

- [`plugin`](plugin.md) — install and manage plugins
- [`marketplace search`](marketplace-search.md)
- [`marketplace verify`](marketplace-verify.md)
{% endraw %}
