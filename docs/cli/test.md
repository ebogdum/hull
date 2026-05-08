# hull test

## Synopsis

`hull test` runs the test hooks for an installed release. Test hooks are Pods or Jobs declared with `$hook: test` (or under `tests/`); they are not run during install/upgrade. A failing test container exits the test as failed.

## When to use it

Use after install/upgrade to verify the release is functioning. Suitable for CI smoke tests.

## Usage

```
hull test <release-name> [flags]
```

## Flags

| Flag | Type | Default | Description |
|---|---|---|---|
| `--filter` | stringArray | — | only run tests whose filename contains the substring (repeatable) |
| `-h, --help` | — | — | help for test |
| `--logs` | — | — | show pod logs after test completes |
| `-o, --output` | string | "human" | output format: human, junit, json |
| `--parallel` | int | 1 | number of tests to run concurrently (1 = sequential) |
| `--retries` | int | — | number of retry attempts per test on failure |
| `--timeout` | duration | 5m0s | timeout waiting for test pods |

## Persistent flags inherited from `hull`

| Flag | Type | Description |
|---|---|---|
| `--debug` | — | enable debug output |
| `--kube-context` | string | Kubernetes context to use |
| `--kubeconfig` | string | path to kubeconfig file |
| `-n, --namespace` | string | Kubernetes namespace |

## Examples

Run tests for a release:

```sh
hull test my-app -n prod
```

Run with retries and parallel test workers:

```sh
hull test my-app -n prod --retries 2 --parallel 4
```

Show pod logs after the test completes (useful for debugging a failing test):

```sh
hull test my-app -n prod --logs
```

Run only one named test (filename substring match):

```sh
hull test my-app -n prod --filter connection
```

Emit JUnit XML for CI ingestion:

```sh
hull test my-app -n prod -o junit > test-results.xml
```

## See also

- [Hooks guide](../guides/hooks.md)
- [Hooks in templates](../templates/hooks.md)
