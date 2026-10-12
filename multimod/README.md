# multimod

> Derived from OpenTelemetry `opentelemetry-go-build-tools` v0.30.0.

`multimod` checks version configuration and prepares releases for repositories
with multiple Go modules. A module set groups modules that share one version.
Different sets can use different versions.

This document describes the commands for a target repository. For releases of
`go-build-tools` itself, use [RELEASING.md](../RELEASING.md).

## Installation

```sh
go install github.com/flc1125/go-build-tools/multimod@latest
```

Run this command from any directory. See the
[repository installation guide](../README.md#installation) for Go requirements
and source installation with a workspace.

## Version configuration

Create `versions.yaml` in the target Git repository's root directory. Run
`multimod` commands from that repository.

```yaml
module-sets:
  app:
    version: v0.2.0
    modules:
      - example.com/project
      - example.com/project/cli
excluded-modules:
  - example.com/project/internal/tools
```

Follow these rules:

- Give each set a name, a valid semantic version, and a list of module import paths.
- Use import paths, not filesystem paths, in `modules` and `excluded-modules`.
- Put each module in exactly one set unless you exclude it.
- Excluded modules do not receive release tags, but their dependencies can receive version updates.
- Update set versions or module groupings when you prepare a release.

The [extended example](./docs/versions-example.yaml) also shows custom version
file references. All commands accept `--versioning-file` to select another
configuration file. The default is `versions.yaml` in the target repository's
root directory.

## Check version configuration

```sh
multimod verify
```

The command checks these conditions:

1. Each module belongs to exactly one set, unless you exclude it.
2. Every configured module exists in the repository.
3. Each set has a valid semantic version.
4. No two sets share the same nonzero major version.

It also warns when a stable module depends on an unstable module in the
repository. These warnings do not cause the command to fail.

For this check, a major version of `v1` or higher is stable. A major version
of `v0` is unstable. Dependencies come from the `require` directives in
`go.mod` files.

## Prepare a release branch

`prerelease` requires a clean working tree. Commit your version configuration
changes before you run it.

1. Create a preparation branch from your intended release base:

   ```sh
   git switch -c release/prepare-<new-version>
   ```

2. Update the intended module set's version in `versions.yaml`.

3. Check the version configuration:

   ```sh
   multimod verify
   ```

4. Commit the configuration change:

   ```sh
   git add versions.yaml
   git commit -m "chore(release): prepare <new-version>"
   ```

5. Check that `git status --short` produces no output.

6. Prepare the module set:

   ```sh
   multimod prerelease --module-set-names <set-name>
   ```

   The command updates dependency versions and existing version reference
   files. It runs `go mod tidy` and commits the changes on a new branch.
   On success, it returns to the original branch.

7. Switch to the generated branch:

   ```sh
   git switch prerelease_<set-name>_<new-version>
   ```

   Include the `v` prefix in `<new-version>`.

8. Review the release changes against the intended release base.

9. Run the target repository's checks on the release branch.

10. Submit the release pull request.

11. Merge the pull request through the normal review process.

### Prerelease flags

| Flag | Behavior |
| --- | --- |
| `--module-set-names` | Optional comma-separated list of set names. Defaults to all sets. |
| `--commit-to-different-branch` | Defaults to `true`. Use `false` to commit on the current branch. |
| `--skip-go-mod-tidy` | Skips dependency tidying. Use only for debugging, not for releases. |

`prerelease` uses the plural `--module-set-names`. The `tag` command uses the
singular `--module-set-name`.

## Preview and create release tags

Update your local release base after the release pull request merges.
Check that `versions.yaml` contains the intended release version.

Preview the tags without creating or deleting them:

```sh
multimod tag --module-set-name <set-name> --preview-tags
```

Preview mode does not require `--commit-hash`. Tags for submodules include the
module directory prefix, such as `cli/v0.2.0`.

Create signed tags on the merged release commit:

```sh
multimod tag --module-set-name <set-name> --commit-hash <release-commit> --print-tags
```

The command runs `git tag -a -s`. Configure Git's signing backend and signing
key before you run this command. See the
[signing requirements](../RELEASING.md#requirements).

`--print-tags` writes one tag name per line to standard output. Progress
messages go to standard error. The command does not push tags.

1. Check each tag's signature with `git tag -v <tag>`.
2. Check its commit with `git rev-list -n 1 <tag>`.
3. Push each intended tag explicitly:

   ```sh
   git push origin <tag>
   ```

Do not reuse a published tag name for a different commit. Consumers and Go
module proxies can retain the original version.

## Delete local tags

**Caution:** This command deletes local tags for the configured set and
version. It does not check whether those tags exist on a remote.

Check the release remote before you delete tags. For example:

```sh
git ls-remote --tags origin
```

Do not delete or recreate tags that consumers can already access.

To undo an unpublished local tagging mistake:

```sh
multimod tag --module-set-name <set-name> --commit-hash <tagged-commit> --delete-module-set-tags
```

Provide `--commit-hash` in deletion mode. Every tag in the configured set
must exist as an annotated tag and point to that commit. The command only
deletes local tags. It does not delete remote tags.

## Release notes

Create the GitHub releases after you push the tags. Include the changes from
the release pull request in the release notes.
