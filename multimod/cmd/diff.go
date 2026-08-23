// Copyright The OpenTelemetry Authors
// Modified by flc1125 for github.com/flc1125/go-build-tools.
// SPDX-License-Identifier: Apache-2.0

package cmd

import (
	"fmt"
	"log"
	"strings"

	"github.com/spf13/cobra"

	"github.com/flc1125/go-build-tools/internal/repo"
	"github.com/flc1125/go-build-tools/multimod/internal/diff"
)

var previousVersion string

var diffCmd = &cobra.Command{
	Use:   "diff",
	Short: "Determines if any files in a module have changed",
	RunE: func(*cobra.Command, []string) error {
		repoRoot, err := repo.FindRoot()
		if err != nil {
			return fmt.Errorf("could not find repo root: %w", err)
		}

		changedFiles, err := diff.HasChanged(repoRoot, versioningFile, previousVersion, moduleSetName)
		if err != nil {
			return fmt.Errorf("error running diff: %w", err)
		}

		if len(changedFiles) > 0 {
			return fmt.Errorf("the following files changed in %s modules since %s: \n%s\nrelease is required for %s modset", moduleSetName, previousVersion, strings.Join(changedFiles, "\n"), moduleSetName)
		}
		log.Printf("No %s modules have changed since %s", moduleSetName, previousVersion)
		return nil
	},
}

func init() {
	// Plain log output, no timestamps.
	log.SetFlags(0)

	rootCmd.AddCommand(diffCmd)
	diffCmd.Flags().StringVarP(&moduleSetName, "module-set-name", "m", "",
		"Name of module set being diff'd. "+
			"Name must be listed in the module set versioning YAML. ",
	)
	if err := diffCmd.MarkFlagRequired("module-set-name"); err != nil {
		log.Fatalf("could not mark module-set-name flag as required: %v", err)
	}

	diffCmd.Flags().StringVarP(&previousVersion, "previous-version", "p", "",
		"Previously released version."+
			"Version must be a tag in the repository. ",
	)
	if err := diffCmd.MarkFlagRequired("previous-version"); err != nil {
		log.Fatalf("could not mark previous-version flag as required: %v", err)
	}
}
