# go-build-tools

Reusable build tools for Go projects.

This repository publishes `crosslink`, `multimod`, and `gotmpl` under the
`github.com/flc1125/go-build-tools` module path. Their initial implementation
is based on OpenTelemetry's
[`opentelemetry-go-build-tools`](https://github.com/open-telemetry/opentelemetry-go-build-tools)
[`v0.30.0`](https://github.com/open-telemetry/opentelemetry-go-build-tools/releases/tag/v0.30.0).

## Tools

### [`crosslink`](./crosslink/README.md)

Manages repositories containing multiple `go.mod` files, including local
`replace` directives, `go.work` files, pruning, and dependency-aware tidy
ordering.

```sh
go install github.com/flc1125/go-build-tools/crosslink@latest
```

### [`multimod`](./multimod/README.md)

Verifies and releases sets of Go modules from a single repository.

```sh
go install github.com/flc1125/go-build-tools/multimod@latest
```

### [`gotmpl`](./gotmpl/README.md)

Generates files from Go templates and JSON data.

```sh
go install github.com/flc1125/go-build-tools/gotmpl@latest
```

## Development

The repository intentionally keeps the upstream multi-module layout. Run the
full local validation with:

```sh
make precommit
```

The modules require Go 1.26.0 or newer. CI validates the latest patch releases
of Go 1.26 and 1.27 on Linux and Windows. Development tool versions, including
golangci-lint v2.11.4, are pinned in `internal/tools`.

See [UPSTREAM.md](./UPSTREAM.md) for provenance and the compatibility baseline.

## License

This project is licensed under the Apache License 2.0. Portions of this
repository are derived from `open-telemetry/opentelemetry-go-build-tools`,
Copyright The OpenTelemetry Authors, and are used under the Apache License 2.0.

See [LICENSE](./LICENSE), [NOTICE](./NOTICE), and [UPSTREAM.md](./UPSTREAM.md)
for the license terms, attribution, and detailed provenance.
