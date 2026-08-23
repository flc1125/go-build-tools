# Release process

All public modules are released together using `versions.yaml`.

## Prepare

1. Update the version for the `tools` module set in `versions.yaml`.
2. Verify the module set:

   ```sh
   make multimod-verify
   ```

3. Prepare the release changes:

   ```sh
   make multimod-prerelease
   ```

4. Review the generated branch and changes, run `make precommit`, and merge the
   release commit through the normal review process.

## Tag

After the release commit is merged, preview the tags. This command is read-only
and does not create tags:

```sh
make print-tags
```

Create the signed tags locally on the release commit:

```sh
make create-tags COMMIT=<release-commit>
```

Review the created tags and confirm that they point to the intended commit.
Then push the existing tags explicitly:

```sh
make push-tags COMMIT=<release-commit> REMOTE=origin
```

`push-tags` refuses to push if any expected tag is missing locally or points to
a different commit. Go module tags are effectively immutable after consumers
observe them, so always verify the commit and complete tag list before pushing.
