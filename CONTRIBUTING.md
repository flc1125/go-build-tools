# Contributing

## Requirements

- Go 1.26.0 or newer. Use the latest patch release of Go 1.26 or 1.27.
- Git.
- GNU Make and a POSIX-compatible shell.
- Network access to download Go dependencies and vulnerability data.

An older Go patch release can fail the vulnerability scan because its standard
library contains known vulnerabilities.

## Development

This is a multi-module Go repository. Use the repository Makefile instead of
running checks against only the root module.

Run local checks from the repository root:

```sh
make precommit
```

The command performs these steps:

1. Format Go source files.
2. Run `go mod tidy` in every module, including `internal/tools`.
3. Run `go generate` in the public modules.
4. Run static checks.
5. Scan for reachable vulnerabilities.
6. Run unit tests.
7. Check the module sets in `versions.yaml`.

The command can change source files, documentation, and dependency files.
Review these changes before you commit them.

To run only tests:

```sh
make test
```

To run tests with the race detector:

```sh
make test-race
```

## CI checks

Run the CI checks after you commit your changes:

```sh
make ci
```

This command checks formatting, generation, static analysis, vulnerabilities,
unit tests, and license headers. It also fails if tracked files differ from
the Git index after the checks.

CI runs race tests and module-set checks separately:

```sh
make test-race
make multimod-verify
```

## Pull requests

Describe any intentional difference from the upstream compatibility baseline
in the pull request description. See [UPSTREAM.md](./UPSTREAM.md) for that
baseline.

Use [RELEASING.md](./RELEASING.md) when you prepare a release of this repository.
