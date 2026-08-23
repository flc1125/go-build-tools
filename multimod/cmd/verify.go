// Copyright The OpenTelemetry Authors
// Modified by flc1125 for github.com/flc1125/go-build-tools.
// SPDX-License-Identifier: Apache-2.0

package cmd

import (
	"log"

	"github.com/spf13/cobra"

	"github.com/flc1125/go-build-tools/multimod/internal/verify"
)

// verifyCmd represents the verify command
var verifyCmd = &cobra.Command{
	Use:   "verify",
	Short: "Verifies that the versioning file is valid",
	Long: `verify checks that all modules listed in sets are valid by verifying the following properties:
- All modules are contained in exactly one module set.
- Versions conform to semver semantics.
- No more than one set of modules exists for any non-zero major version.
- Script warns if any stable modules depend on any unstable modules.
`,
	RunE: func(*cobra.Command, []string) error {
		log.Println("Using versioning file", versioningFile)

		return verify.Run(versioningFile)
	},
}

func init() {
	// Plain log output, no timestamps.
	log.SetFlags(0)

	rootCmd.AddCommand(verifyCmd)
}
