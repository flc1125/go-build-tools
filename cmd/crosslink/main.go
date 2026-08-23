// Copyright The OpenTelemetry Authors
// Modified by flc1125 for github.com/flc1125/go-build-tools.
// SPDX-License-Identifier: Apache-2.0

// Crosslink updates mono-repo go.mod. It ensures local versions of modules
// are used when building the code.
package main

import "github.com/flc1125/go-build-tools/cmd/crosslink/internal/cmd"

func main() {
	cmd.Execute()
}
