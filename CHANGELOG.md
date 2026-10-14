# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project follows [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

The three CLI modules share each release version. Version links use the
`crosslink` release tag; `gotmpl` and `multimod` have matching release tags.
The shared root module uses published commit versions.

## [Unreleased]

## [0.1.0] - 2026-10-14

Initial release of `go-build-tools`.

### Added

- `crosslink` manages local dependencies in multi-module repositories with
  `replace` directives, Go workspaces, and dependency-aware tidy ordering.
- `multimod` checks module version configuration, prepares release branches,
  and creates release tags for module sets.
- `gotmpl` generates files from Go templates and JSON data.
- Workspace-based local development without local `replace` directives in
  published CLI modules.
- Automated dependency updates through Renovate.

### Changed

- Require Go 1.26.0 or newer. Use the latest patch release of Go 1.26 or 1.27.
  CI checks both versions on Linux and Windows.
- Update dependencies, including `golang.org/x`, `go-git` v5.19.3,
  `testify` v1.12.1, and `zap` v1.28.0.
- Update development linting to `golangci-lint` v2.14.0.
- Clarify installation and release workflows, consolidate upstream attribution,
  and use generic Go project examples.

### Fixed

- Support versioned CLI installation with `go install` from outside the
  repository. See [#17](https://github.com/flc1125/go-build-tools/pull/17).
- Skip the Unix file-mode assertion in `gotmpl` tests on Windows.
- Keep dependent module metadata tidy after automated dependency updates.

[unreleased]: https://github.com/flc1125/go-build-tools/compare/crosslink/v0.1.0...HEAD
[0.1.0]: https://github.com/flc1125/go-build-tools/releases/tag/crosslink/v0.1.0
