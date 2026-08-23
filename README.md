# go-build-tools

Reusable build tools for Go projects.

The first release is a focused derivative of OpenTelemetry's
[`opentelemetry-go-build-tools`](https://github.com/open-telemetry/opentelemetry-go-build-tools)
at [`v0.30.0`](https://github.com/open-telemetry/opentelemetry-go-build-tools/releases/tag/v0.30.0).
It preserves the behavior of three tools while publishing them under this
repository's module path.

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

The modules require Go 1.25 or newer. Static analysis is pinned to Go 1.25.14
and golangci-lint v2.11.4 to reproduce the upstream compatibility baseline.

See [UPSTREAM.md](./UPSTREAM.md) for provenance and the compatibility baseline.

## License

Apache License 2.0. See [LICENSE](./LICENSE) and [NOTICE](./NOTICE).
