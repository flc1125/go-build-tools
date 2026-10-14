# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project follows [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

The root module and the three CLI modules share each release version. Version
links use the root tag; each CLI module has a matching module-prefixed tag.

## [Unreleased]

## [0.1.0] - 2026-10-14

Initial release of `go-build-tools`.

### Added

- `crosslink` manages local dependencies in multi-module repositories with
  `replace` directives, Go workspaces, and dependency-aware tidy ordering.
  ([4567c5c](https://github.com/flc1125/go-build-tools/commit/4567c5c238660ad56d6058da9da0cd3562120d2d))
- `multimod` checks module version configuration, prepares release branches,
  and creates release tags for module sets.
  ([4567c5c](https://github.com/flc1125/go-build-tools/commit/4567c5c238660ad56d6058da9da0cd3562120d2d))
- `gotmpl` generates files from Go templates and JSON data.
  ([4567c5c](https://github.com/flc1125/go-build-tools/commit/4567c5c238660ad56d6058da9da0cd3562120d2d))
- Shared internal packages released with the CLI modules under the same version.
  ([#22](https://github.com/flc1125/go-build-tools/pull/22))
- Workspace-based local development without local `replace` directives in
  published CLI modules.
  ([#17](https://github.com/flc1125/go-build-tools/pull/17))
- Automated dependency updates through Renovate.
  ([#3](https://github.com/flc1125/go-build-tools/pull/3))

### Changed

- Require Go 1.26.0 or newer. Use the latest patch release of Go 1.26 or 1.27.
  CI checks both versions on Linux and Windows.
  ([#4](https://github.com/flc1125/go-build-tools/pull/4),
  [#7](https://github.com/flc1125/go-build-tools/pull/7))
- Update `golang.org/x` dependencies.
  ([#8](https://github.com/flc1125/go-build-tools/pull/8),
  [#13](https://github.com/flc1125/go-build-tools/pull/13),
  [#16](https://github.com/flc1125/go-build-tools/pull/16))
- Update `go-git` to v5.19.3.
  ([#15](https://github.com/flc1125/go-build-tools/pull/15))
- Update `testify` to v1.12.1.
  ([#10](https://github.com/flc1125/go-build-tools/pull/10))
- Update `zap` to v1.28.0.
  ([#11](https://github.com/flc1125/go-build-tools/pull/11))
- Configure development linting and update `golangci-lint` to v2.14.0.
  ([#5](https://github.com/flc1125/go-build-tools/pull/5),
  [#12](https://github.com/flc1125/go-build-tools/pull/12),
  [#14](https://github.com/flc1125/go-build-tools/pull/14))
- Clarify installation and release workflows, consolidate upstream attribution,
  and use generic Go project examples.
  ([#17](https://github.com/flc1125/go-build-tools/pull/17),
  [#20](https://github.com/flc1125/go-build-tools/pull/20),
  [#22](https://github.com/flc1125/go-build-tools/pull/22))

### Fixed

- Support versioned CLI installation with `go install` from outside the
  repository. ([#17](https://github.com/flc1125/go-build-tools/pull/17))
- Skip the Unix file-mode assertion in `gotmpl` tests on Windows.
  ([#2](https://github.com/flc1125/go-build-tools/pull/2))
- Keep dependent module metadata tidy after automated dependency updates.
  ([#19](https://github.com/flc1125/go-build-tools/pull/19))

[unreleased]: https://github.com/flc1125/go-build-tools/compare/v0.1.0...HEAD
[0.1.0]: https://github.com/flc1125/go-build-tools/releases/tag/v0.1.0
