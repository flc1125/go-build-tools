# Upstream provenance

This repository's initial implementation is derived from:

- Project: [`open-telemetry/opentelemetry-go-build-tools`](https://github.com/open-telemetry/opentelemetry-go-build-tools)
- Release: [`v0.30.0`](https://github.com/open-telemetry/opentelemetry-go-build-tools/releases/tag/v0.30.0)
- Commit: `47adfe3ba15f8d3385ef12bcb846b98e8fa86ee5`
- License: Apache License 2.0

## Initial import

The initial import includes the upstream `crosslink`, `multimod`, and `gotmpl`
tools together with the root-module packages they require:

- `internal/repo`
- `internal/syncerror`

The initial import preserved the tools' command-line behavior and configuration
formats.

## Repository changes

Repository-specific changes include module paths, metadata, version
configuration, dependency packaging, and development automation.

The `multimod` dependency on `github.com/go-git/go-git/v5` is updated from
v5.18.0 to v5.19.2, together with the required transitive dependencies, to
address vulnerabilities reported after the upstream v0.30.0 release.

## Attribution

Original OpenTelemetry copyright and license headers are retained. Modified
source and example files carry modification notices.
