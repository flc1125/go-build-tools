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

Install the tools from any directory:

```sh
go install github.com/flc1125/go-build-tools/crosslink@latest
go install github.com/flc1125/go-build-tools/multimod@latest
go install github.com/flc1125/go-build-tools/gotmpl@latest
```

Use a specific version, such as `@v0.1.0`, instead of `@latest` to pin an
available release. Go installs the executables in `GOBIN`, or in the first
`GOPATH` entry's `bin` directory when `GOBIN` is empty.

Add the installation directory to your `PATH`.

### Source installation

Use a workspace when you want to install the local source, including local
changes to the shared internal packages.

1. Clone this repository:

   ```sh
   git clone https://github.com/flc1125/go-build-tools.git
   ```

2. Enter the repository directory:

   ```sh
   cd go-build-tools
   ```

3. Create the local workspace:

   ```sh
   make gowork
   ```

4. Install the tools from the repository root:

   ```sh
   go install github.com/flc1125/go-build-tools/crosslink
   go install github.com/flc1125/go-build-tools/multimod
   go install github.com/flc1125/go-build-tools/gotmpl
   ```

Do not add a version suffix for source installation. A versioned install uses
published source instead of your local changes.

## Development

The repository keeps the upstream multi-module layout. Run local checks from
the repository root:

```sh
make precommit
```

This command can change source files, documentation, and dependency files.
It also creates or updates the ignored `go.work` file for local dependencies.
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
