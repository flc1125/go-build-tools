# go-build-tools

Reusable build tools for Go projects.

This repository provides `crosslink`, `multimod`, and `gotmpl` under the
`github.com/flc1125/go-build-tools` module path. Their initial implementation
is based on OpenTelemetry's
[`opentelemetry-go-build-tools`](https://github.com/open-telemetry/opentelemetry-go-build-tools)
[`v0.30.0`](https://github.com/open-telemetry/opentelemetry-go-build-tools/releases/tag/v0.30.0).

## Tools

### [`crosslink`](./crosslink/README.md)

Adds local `replace` directives for dependencies between Go modules in one
repository. It also removes unused replacements, manages `go.work`, and
generates a dependency-aware order for `go mod tidy`.

### [`multimod`](./multimod/README.md)

Checks version configuration and prepares releases for sets of Go modules in
one repository.

### [`gotmpl`](./gotmpl/README.md)

Generates files from Go templates and JSON data.

## Installation

The modules require Go 1.26.0 or newer. Use the latest patch release of Go 1.26
or 1.27 for development and security checks.

1. Clone this repository:

   ```sh
   git clone https://github.com/flc1125/go-build-tools.git
   ```

2. Enter the repository directory:

   ```sh
   cd go-build-tools
   ```

3. Install the tools from their module directories:

   ```sh
   (cd crosslink && go install .)
   (cd multimod && go install .)
   (cd gotmpl && go install .)
   ```

Go installs the executables in `GOBIN`, or in the first `GOPATH` entry's `bin`
directory when `GOBIN` is empty. Add the installation directory to your `PATH`.

`crosslink` and `multimod` use local `replace` directives to reference the root
module. Go rejects these directives when you install with a version suffix
such as `@latest`. Install these tools from a clone instead.

`gotmpl` has no local `replace` directives. You can also install it directly:

```sh
go install github.com/flc1125/go-build-tools/gotmpl@latest
```

## Development

The repository keeps the upstream multi-module layout. Run local checks from
the repository root:

```sh
make precommit
```

This command can change source files, documentation, and dependency files.
CI checks the latest patch releases of Go 1.26 and 1.27 on Linux and Windows.
[`internal/tools/go.mod`](./internal/tools/go.mod) specifies the development
tool versions.

See [CONTRIBUTING.md](./CONTRIBUTING.md) for development commands.
See [RELEASING.md](./RELEASING.md) for this repository's release process.
See [UPSTREAM.md](./UPSTREAM.md) for provenance and the compatibility baseline.

## License

This project is licensed under the Apache License 2.0. Portions of this
repository are derived from `open-telemetry/opentelemetry-go-build-tools`,
Copyright The OpenTelemetry Authors, and are used under the Apache License 2.0.

See [LICENSE](./LICENSE), [NOTICE](./NOTICE), and [UPSTREAM.md](./UPSTREAM.md)
for the license terms, attribution, and detailed provenance.
