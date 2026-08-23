// Copyright 2026 flc1125
// SPDX-License-Identifier: Apache-2.0

package cmd

import (
	"io"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHelpOutsideRepository(t *testing.T) {
	t.Chdir(t.TempDir())
	rootCmd.SetArgs([]string{"--help"})
	rootCmd.SetOut(io.Discard)
	rootCmd.SetErr(io.Discard)
	t.Cleanup(func() {
		rootCmd.SetArgs(nil)
		rootCmd.SetOut(nil)
		rootCmd.SetErr(nil)
	})

	_, err := rootCmd.ExecuteC()

	require.NoError(t, err)
}

func TestResolveVersioningFileExplicitPathOutsideRepository(t *testing.T) {
	t.Chdir(t.TempDir())
	original := versioningFile
	versioningFile = filepath.Join(t.TempDir(), "versions.yaml")
	t.Cleanup(func() { versioningFile = original })

	require.NoError(t, resolveVersioningFile())
}
