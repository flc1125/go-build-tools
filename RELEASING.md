# Release process

This repository releases the root module, `crosslink`, `gotmpl`, and `multimod`
together in the `tools` module set from `versions.yaml`.

The root module contains shared internal packages. It shares the CLI release
version. Each release has a root tag, such as `v0.1.0`, and matching CLI tags:
`crosslink/v0.1.0`, `gotmpl/v0.1.0`, and `multimod/v0.1.0`.

## Requirements

Use the development requirements in [CONTRIBUTING.md](./CONTRIBUTING.md#requirements).
You also need permission to push tags to the release remote.

`make create-tags` runs `git tag -a -s` to create signed tags. Configure Git's
signing backend and signing key before you create tags. Git uses settings such
as `gpg.format` and `user.signingkey`.

For OpenPGP signing, install GnuPG. Make your private key available to its
agent. Do not store private keys in this repository.

Run all commands below from the repository root. Replace placeholders such as
`<new-version>` and `<release-commit>` with the intended values.

## Shared dependencies

`crosslink` and `multimod` require the root module at the shared release version.
`multimod prerelease` updates these requirements together with other references
to modules in the `tools` set.

Public CLI modules must not contain local `replace` directives. Development
uses a workspace for local shared packages. `make tidy` uses temporary module
files to resolve the local root module before its release tag exists, without
adding replacements to the published CLI files.

`make install-check` validates the root module and CLI sources through a local
file proxy with workspaces disabled. No release tags need to exist for this
check. See [CONTRIBUTING.md](./CONTRIBUTING.md#installation-checks) for details.

## Prepare

1. Switch to `main`:

   ```sh
   git switch main
   ```

2. Update `main` from the release remote:

   ```sh
   git pull --ff-only origin main
   ```

3. Create a release preparation branch:

   ```sh
   git switch -c release/prepare-<new-version>
   ```

   This branch keeps version preparation commits separate from your local
   `main` branch.

4. Update the `tools` version in `versions.yaml`.

   Choose a version without existing release tags.

   Add a dated release entry to [CHANGELOG.md](./CHANGELOG.md). Move the relevant
   `Unreleased` entries into it and update the version links to use the root tag.

5. Check the module set:

   ```sh
   make multimod-verify
   ```

6. Commit the version configuration change:

   ```sh
   git add versions.yaml CHANGELOG.md
   git commit -m "chore(release): prepare <new-version>"
   ```

7. Check that the working tree is clean:

   ```sh
   git status --short
   ```

   This command must produce no output. `multimod prerelease` rejects
   uncommitted changes.

8. Prepare the release changes:

   ```sh
   make multimod-prerelease
   ```

    This command creates and commits the release changes on a new branch.
    On success, it returns to the original branch.

    Dependency tidying is deferred to `make precommit` on the generated branch.
    The release's root-module version does not need to be published first.

9. Switch to the generated release branch:

   ```sh
   git switch prerelease_tools_<new-version>
   ```

   Use the version from `versions.yaml`, including its `v` prefix.

10. Review the release changes:

    ```sh
    git diff main
    ```

11. Run local checks on the release branch:

    ```sh
    make precommit
    ```

12. Commit any additional fixes from the review or checks.

13. Open a pull request to merge the release branch into `main`.

14. Merge the pull request through the normal review process.

## Tag

1. Switch to `main` after the release pull request merges:

   ```sh
   git switch main
   ```

2. Update `main` from the release remote:

   ```sh
   git pull --ff-only origin main
   ```

3. Check that `versions.yaml` contains the version from the merged release commit.

4. Preview the complete tag list:

   ```sh
   make print-tags
   ```

    This command previews the root tag and all three CLI tags. It does not
    create or delete tags.

5. Create signed tags on the merged release commit:

   ```sh
   make create-tags COMMIT=<release-commit>
   ```

6. Check the signature of each tag:

   ```sh
   git tag -v <tag>
   ```

7. Check that each tag points to the intended release commit:

   ```sh
   git rev-list -n 1 <tag>
   ```

8. Push the existing tags to the release remote:

   ```sh
   make push-tags COMMIT=<release-commit> REMOTE=origin
   ```

`push-tags` refuses to push if any expected tag is missing locally or points to
a different commit. It checks commit targets, not tag signatures.

Do not reuse a published tag name for a different commit. Consumers and Go
module proxies can retain the original version.

For local tagging mistakes, see the [tag deletion cautions](./multimod/README.md#delete-local-tags).

## Release notes

Create the GitHub releases for the published tags. Use the corresponding version
entry in [CHANGELOG.md](./CHANGELOG.md) as the release notes.
