// Copyright The OpenTelemetry Authors
// Modified by flc1125 for github.com/flc1125/go-build-tools.
// SPDX-License-Identifier: Apache-2.0

// Multimod enables the release of Go modules with flexible versioning.
package main

import "github.com/flc1125/go-build-tools/cmd/multimod/internal/cmd"

func main() {
	cmd.Execute()
}
