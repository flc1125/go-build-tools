# Upstream provenance

This repository contains a focused derivative of:

- Project: `open-telemetry/opentelemetry-go-build-tools`
- Release: `v0.30.0`
- Commit: `47adfe3ba15f8d3385ef12bcb846b98e8fa86ee5`
- License: Apache License 2.0

The initial import includes the upstream `crosslink`, `multimod`, and `gotmpl`
tools together with the root-module packages they require:

- `internal/repo`
- `internal/syncerror`

Initial modifications are intentionally limited to module paths, repository
metadata, version configuration, and repository-level development automation.
The command-line behavior and configuration formats remain based on the
upstream release.

The `multimod` dependency on `github.com/go-git/go-git/v5` is updated from
v5.18.0 to v5.19.2, together with the required transitive dependencies, to
address vulnerabilities reported after the upstream v0.30.0 release.

## Installation and local development

The public CLI modules do not contain local `replace` directives. This permits
versioned installation with `go install`. `crosslink` and `multimod` pin a
published root-module commit for their shared internal packages.

The CLI modules share the `tools` release version. The root module stays
outside that release set. A local workspace replaces the root module during
development without changing the published CLI dependencies.

This changes dependency packaging and development automation, not the tools'
command-line behavior or configuration formats.

Original OpenTelemetry copyright and license headers are retained. Files whose
Go import paths were adapted carry an additional modification notice.
