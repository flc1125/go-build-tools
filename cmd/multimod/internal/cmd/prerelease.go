// Copyright The OpenTelemetry Authors
// Modified by flc1125 for github.com/flc1125/go-build-tools.
// SPDX-License-Identifier: Apache-2.0

package cmd

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"

	"github.com/flc1125/go-build-tools/cmd/multimod/internal/prerelease"
)

var (
	allModuleSets           bool
	moduleSetNames          []string
	skipGoModTidy           bool
	commitToDifferentBranch bool
)

// prereleaseCmd represents the prerelease command
var prereleaseCmd = &cobra.Command{
	Use:   "prerelease",
	Short: "Prepares files for new version release",
	Long: `Updates version numbers and commits to a new branch for release:
- Checks that the working tree is clean.
- Checks that Git tags do not already exist for the new module set version.
- Switches to a new branch called prerelease_<module set name>_<new version>.
- Updates version.go files, if they exist.
- Updates module versions in all go.mod files.
- Attempts to call 'go mod tidy' in the directory of each modified go.mod file.
- Adds and commits changes to Git branch`,
	PreRunE: func(cmd *cobra.Command, _ []string) error {
		if allModuleSets {
			// do not require module set names if operating on all module sets
			if err := cmd.Flags().SetAnnotation(
				"module-set-names",
				cobra.BashCompOneRequiredFlag,
				[]string{"false"},
			); err != nil {
				return fmt.Errorf("could not set module-set-names flag as not required: %w", err)
			}
		}
		return nil
	},
	RunE: func(*cobra.Command, []string) error {
		log.Println("Using versioning file", versioningFile)

		return prerelease.Run(versioningFile, moduleSetNames, skipGoModTidy, commitToDifferentBranch)
	},
}

func init() {
	// Plain log output, no timestamps.
	log.SetFlags(0)

	rootCmd.AddCommand(prereleaseCmd)

	prereleaseCmd.Flags().BoolVarP(&allModuleSets, "all-module-sets", "a", false,
		"update versions of all modules sets",
	)
	prereleaseCmd.Flags().Lookup("all-module-sets").Deprecated = "default behavior"

	prereleaseCmd.Flags().StringSliceVarP(
		&moduleSetNames,
		"module-set-names",
		"m",
		nil,
		"allow-list of module sets to update",
	)
	prereleaseCmd.Flags().Lookup("module-set-names").DefValue = "all modules sets"

	prereleaseCmd.Flags().BoolVarP(
		&skipGoModTidy,
		"skip-go-mod-tidy",
		"s",
		false,
		"skip calling 'go mod tidy' after update",
	)
	prereleaseCmd.Flags().BoolVarP(&commitToDifferentBranch, "commit-to-different-branch", "b", true,
		"commit to a branch other than the current one",
	)
}
