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

After the release commit is merged, print the tags before pushing them:

```sh
make print-tags COMMIT=<release-commit>
```

When the commit and tags have been verified, push them explicitly:

```sh
make push-tags COMMIT=<release-commit> REMOTE=origin
```

Go module tags are effectively immutable after consumers observe them. Always
verify the commit and the complete tag list before pushing.
