# crosslink

> Derived from OpenTelemetry `opentelemetry-go-build-tools` v0.30.0.

`crosslink` manages local dependencies in repositories with multiple Go
modules. It adds `replace` directives, removes unused replacements, manages
`go.work`, and generates an order for `go mod tidy`.

## Installation

```sh
go install github.com/flc1125/go-build-tools/crosslink@latest
```

Run this command from any directory. See the
[repository installation guide](../README.md#installation) for Go requirements
and source installation with a workspace.

## Before you run crosslink

The default command writes `go.mod` files. It adds missing replacements but
keeps existing replacements unless you specify `--overwrite`. It keeps unused
replacements unless you specify `--prune`.

Commit or save your existing changes before you run a command that changes
files. Review `git diff` after the command finishes.

The root module defines the namespace for dependency replacements. For example,
`example.com/project/library` belongs to the `example.com/project` namespace.
New replacements use local paths without module versions.

## Root directory

`--root` selects the directory to scan. All commands accept this flag.

Without `--root`, `crosslink` searches the current directory and its parents
for a Git repository. It returns an error if it cannot find one.

The default command, `prune`, and `tidylist` require a `go.mod` file in the root
directory. `work` scans the selected directory for Go modules.

## Add local replacements

Run the default command to add missing replacements for direct and transitive
dependencies within the root module namespace:

```sh
crosslink --root=.
```

Only dependencies with matching modules in the scanned repository receive
local replacements. The command does not add replacements for every module
in the repository.

### Overwrite existing replacements

**Caution:** `--overwrite` changes existing replacement paths and removes
their module versions.

```sh
crosslink --root=. --overwrite
```

## Remove unused replacements

**Caution:** Pruning removes unused `replace` directives from `go.mod` files.
It does not remove module directories or files.

Run `prune` to remove replacements without adding new ones:

```sh
crosslink prune --root=.
```

Use `--prune` or `-p` with the default command to add missing replacements and
remove unused replacements:

```sh
crosslink --root=. --prune
```

Pruning only removes replacements within the root module namespace. It keeps
replacements for direct and transitive dependencies. It also keeps
replacements outside that namespace.

## Exclude replacements or skip files

### `--exclude`

The default command and `prune` accept `--exclude`. Its values are module
import paths, not filesystem paths.

The flag protects replacement entries whose old module path matches a listed
value. It does not prevent changes inside the listed modules.

Use a comma-separated list without spaces:

```sh
crosslink --root=. --prune --exclude=example.com/project/a,example.com/project/b
```

You can repeat the flag instead:

```sh
crosslink --root=. --prune \
  --exclude=example.com/project/a \
  --exclude=example.com/project/b
```

### `--skip`

The default command and `tidylist` accept `--skip`. Its values are paths to
`go.mod` files, relative to `--root`.

The default command does not change these files. `tidylist` omits these files
from its dependency graph.

Use a comma-separated list without spaces:

```sh
crosslink --root=. --skip=cmd/example/go.mod,cmd/other/go.mod
```

You can repeat the flag instead:

```sh
crosslink --root=. \
  --skip=cmd/example/go.mod \
  --skip=cmd/other/go.mod
```

## Generate a workspace

`work` creates or updates `go.work`. It adds `use` directives for all Go
modules below the selected root directory.

It removes stale `use` paths that start with `./`. It keeps other existing
`use` paths and `replace` directives.

```sh
crosslink work --root=. --go=1.26.0
```

`--go` sets the Go version only when the command creates a new `go.work` file.
Its default is `1.23.0`. It does not change the version in an existing file.

Use a workspace version that supports all selected modules. The command
checks version syntax, not the modules' minimum Go versions.

## Generate a tidy order

`tidylist` writes module directories in dependency order. Each directory is
relative to `--root`, with one directory per line.

```sh
crosslink tidylist tidy-order.txt --root=.
```

The command overwrites the output file if it exists. It does not run
`go mod tidy` itself.

Run this shell example from the selected root directory to apply the order:

```sh
while IFS= read -r dir || [ -n "$dir" ]; do
  go -C "$dir" mod tidy || exit
done < tidy-order.txt
```

Modules with circular dependencies can appear more than once in the output.
Use `--allow-circular` to provide a list of directories with permitted cycles:

```sh
crosslink tidylist tidy-order.txt --root=. --allow-circular=circular.txt
```

The allowlist uses paths relative to `--root`, not module import paths.
For example:

```text
module-a
module-b
```

The file can contain empty lines and comments that start with `#`.
The command fails if the detected cycles do not match the allowlist.

Use `--validate` to check the generated order against the dependency graph:

```sh
crosslink tidylist tidy-order.txt --root=. --validate
```

## Verbose output

All commands accept `--verbose` or `-v`. Verbose output is off by default.
`--overwrite` enables it unless you explicitly set `--verbose`.

```sh
crosslink --root=. --verbose
```

To overwrite replacements without verbose output:

```sh
crosslink --root=. --overwrite --verbose=false
```

## Background

OpenTelemetry developed the original tool for its multi-module repositories.
See [opentelemetry-go](https://github.com/open-telemetry/opentelemetry-go) and
[opentelemetry-collector-contrib](https://github.com/open-telemetry/opentelemetry-collector-contrib)
for examples.
