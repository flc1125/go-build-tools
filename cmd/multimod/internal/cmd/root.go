// Copyright The OpenTelemetry Authors
// Modified by flc1125 for github.com/flc1125/go-build-tools.
// SPDX-License-Identifier: Apache-2.0

// Package cmd contains the command line interface for multimod.
package cmd

import (
	"fmt"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/flc1125/go-build-tools/internal/repo"
)

var versioningFile string

const (
	defaultVersionsConfigName = "versions"
	defaultVersionsConfigType = "yaml"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "multimod",
	Short: "Enables the release of Go modules with flexible versioning",
	Long: `A Golang release versioning and tagging tool that simplifies and
automates versioning for repos with multiple Go modules.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func resolveVersioningFile() error {
	if versioningFile != "" {
		return nil
	}

	repoRoot, err := repo.FindRoot()
	if err != nil {
		return fmt.Errorf("could not find repo root: %w", err)
	}

	versioningFile = filepath.Join(repoRoot,
		fmt.Sprintf("%v.%v", defaultVersionsConfigName, defaultVersionsConfigType))
	return nil
}

func init() {
	cobra.OnInitialize()

	rootCmd.PersistentPreRunE = func(*cobra.Command, []string) error {
		return resolveVersioningFile()
	}
	rootCmd.PersistentFlags().StringVarP(&versioningFile, "versioning-file", "v", "",
		"Path to versioning file that contains definitions of all module sets. "+
			"If unspecified, defaults to versions.yaml in the Git repo root.")
}
