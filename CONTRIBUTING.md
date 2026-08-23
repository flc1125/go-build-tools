# Contributing

## Development

This is a multi-module Go repository. Use the repository Makefile instead of
running commands against only the root module.

```sh
make precommit
```

That command formats source files, tidies every module, runs static checks, and
runs all unit tests.

To run only tests:

```sh
make test
```

Changes that intentionally diverge from the upstream compatibility baseline
should document the behavior difference in the pull request description.
