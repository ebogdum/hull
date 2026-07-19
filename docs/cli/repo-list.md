# hull repo list

Show the repositories you have registered.

## When to use it

- To confirm a [`hull repo add`](repo-add.md) took effect.
- To look up a repository's URL, for example to pass to `hull pull --repo`.

## What happens

Hull reads your repository list at `~/.config/hull/repositories.yaml` and
prints every entry's name and URL. Nothing is fetched — this reads only local
configuration, with no network or cluster access. If you have not registered
any repositories, hull prints `No repositories configured.`

## Usage

```
hull repo list [flags]
```

## Flags

| Flag | Effect |
|---|---|
| `-o, --output` | Choose the output format: `table` (default), `json`, or `yaml`. |

## Worked example

```
$ hull repo add my-charts https://charts.example.com
"my-charts" has been added to your repositories

$ hull repo list
NAME                 URL
my-charts            https://charts.example.com
```

Get one repository's URL as JSON for scripting:

```
$ hull repo list -o json | jq -r '.[] | select(.name=="my-charts") | .url'
https://charts.example.com
```

## See also

- [`repo add`](repo-add.md) — register a repository
- [`repo update`](repo-update.md) — refresh registered repositories
- [`search`](search.md) · [`pull`](pull.md)
