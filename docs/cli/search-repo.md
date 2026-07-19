# hull search repo

`hull search repo` searches the repositories you have added with
`hull repo add` for a package whose name or description contains a keyword.

## When to use it

Use it to find a package across your own trusted repositories. Reach for
[`search hub`](search-hub.md) instead when you want to discover packages you
have not added a repository for yet.

## What happens

1. You pass one keyword; hull lowercases it.
2. hull reads your configured repository list. If it is empty, it prints
   `No repositories configured. Use 'hull repo add' first.` and stops.
3. For each repository it fetches the index and keeps every chart whose name
   or whose description contains the keyword.
4. It prints one row per match — repository-qualified name, latest version,
   app version, and description — under a `NAME VERSION APP VERSION
   DESCRIPTION` header. With no matches it prints `No results found for
   "<keyword>"`.

No cluster is contacted.

## Usage

```
hull search repo <keyword> [flags]
```

## Flags

Inherits the global flags.

## Worked example

You have added the `bitnami` and `myrepo` repositories. Search for `redis`:

```sh
hull search repo redis
```

Output:

```
NAME                           VERSION         APP VERSION     DESCRIPTION
bitnami/redis                  20.1.3          7.4.0           Open source, in-memory data store
myrepo/cache-redis             10.2.0          7.2.4           Redis with clustering enabled
```

The `bitnami/redis` row matched because `redis` is in the chart name; the
`myrepo/cache-redis` row matched on either its name or its description. Each
row shows the latest version in that repository — feed the `NAME` straight
into `hull pull` or `hull install`.

## See also

- [`search`](search.md)
- [`search hub`](search-hub.md) — search Artifact Hub instead
- [`repo`](repo.md) — add the repositories searched here
- [`pull`](pull.md) — download a matched chart
- [`install`](install.md) — install a matched package
