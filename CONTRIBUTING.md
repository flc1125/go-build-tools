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

1. Create or update the local Go workspace.
2. Format Go source files.
3. Run `go mod tidy` in every module, including `internal/tools`.
4. Run `go generate` in the public modules.
5. Run static checks.
6. Scan for reachable vulnerabilities.
7. Run unit tests.
8. Check the module sets in `versions.yaml`.
9. Check versioned installation of the CLI modules.

The command can change source files, documentation, and dependency files.
Review these changes before you commit them.

`go.work` and `go.work.sum` stay outside Git. The workspace uses local shared
packages without adding `replace` directives to the public CLI modules.

`make tidy` runs with `GOWORK=off` and temporary module files. CLI dependencies
on the root module resolve to local shared packages during tidying, even before
the new root-module version is published. Existing replacements are preserved.
The temporary root replacement is removed before dependency files are written
back. Local replacement paths and fixture checksums are not committed.

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
unit tests, license headers, and versioned CLI installation. It also fails if
tracked files differ from the Git index after the checks.

CI runs race tests and module-set checks separately:

```sh
make test-race
make multimod-verify
```

### Installation checks

To run only the versioned installation check:

```sh
make install-check
```

The check builds a temporary file proxy from the tracked root-module and CLI
files in your working tree. The root fixture uses the version required by the
CLI modules, so the check also works before release tags exist. The check runs
`go install` with `@latest` outside the repository.
It disables workspaces and uses a fresh module cache for the installations.

The check can reuse downloaded dependencies through the existing download
cache. Other dependencies use your configured Go proxy. Published dependencies
keep the configured checksum settings. Only the local fixture modules skip
checksum database checks.

The fixtures are local. This check does not publish modules, create
release tags, or change your installed tools.

## Pull requests

Describe the purpose of the change and its effect on command-line behavior,
configuration formats, or installation. Include the checks you ran.

Record notable changes in the `Unreleased` section of [CHANGELOG.md](./CHANGELOG.md).
Add the relevant pull request links to each entry after the pull request is
open. Use a commit link for a change made without a pull request.

Use [RELEASING.md](./RELEASING.md) when you prepare a release of this repository.
