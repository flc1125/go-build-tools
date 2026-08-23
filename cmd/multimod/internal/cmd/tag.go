// Copyright The OpenTelemetry Authors
// Modified by flc1125 for github.com/flc1125/go-build-tools.
// SPDX-License-Identifier: Apache-2.0

package cmd

import (
	"errors"
	"log"

	"github.com/spf13/cobra"

	"github.com/flc1125/go-build-tools/cmd/multimod/internal/tag"
)

var (
	commitHash          string
	deleteModuleSetTags bool
	moduleSetName       string
	printTags           bool
	previewTags         bool
)

// tagCmd represents the tag command
var tagCmd = &cobra.Command{
	Use:   "tag",
	Short: "Applies Git tags to specified commit",
	Long: `Tag script to add Git tags to a specified commit hash created by prerelease script:
- Creates new Git tags for all modules being updated.
- If tagging fails in the middle of the script, the recently created tags will be deleted.`,
	PreRunE: func(*cobra.Command, []string) error {
		if previewTags && deleteModuleSetTags {
			return errors.New("--preview-tags and --delete-module-set-tags cannot be used together")
		}
		if !previewTags && commitHash == "" {
			return errors.New("required flag(s) \"commit-hash\" not set")
		}
		return nil
	},
	RunE: func(*cobra.Command, []string) error {
		log.Println("Using versioning file", versioningFile)

		return tag.Run(versioningFile, moduleSetName, commitHash, deleteModuleSetTags, printTags, previewTags)
	},
}

func init() {
	// Plain log output, no timestamps.
	log.SetFlags(0)

	rootCmd.AddCommand(tagCmd)

	tagCmd.Flags().StringVarP(&commitHash, "commit-hash", "c", "",
		"Git commit hash to tag.",
	)
	tagCmd.Flags().StringVarP(&moduleSetName, "module-set-name", "m", "",
		"Name of module set being tagged. "+
			"Name must be listed in the module set versioning YAML. ",
	)
	if err := tagCmd.MarkFlagRequired("module-set-name"); err != nil {
		log.Fatalf("could not mark module-set-name flag as required: %v", err)
	}

	tagCmd.Flags().BoolVarP(&deleteModuleSetTags, "delete-module-set-tags", "d", false,
		"Specify this flag to delete all module tags associated with the version listed for the module set in the versioning file. Should only be used to undo recent tagging mistakes.",
	)

	tagCmd.Flags().BoolVarP(&printTags, "print-tags", "p", false,
		"Specify this flag to print all tags after tagging is complete. Printed tags are new-line delimited.",
	)
	tagCmd.Flags().BoolVar(&previewTags, "preview-tags", false,
		"Print the tags that would be created without creating or deleting any tags.",
	)
}
