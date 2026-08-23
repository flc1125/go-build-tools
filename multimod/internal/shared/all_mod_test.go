// Copyright The OpenTelemetry Authors
// Modified by flc1125 for github.com/flc1125/go-build-tools.
// SPDX-License-Identifier: Apache-2.0

package shared

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/flc1125/go-build-tools/multimod/internal/shared/sharedtest"
)

func TestNewAllModulePathMap(t *testing.T) {
	tmpRootDir := t.TempDir()
	modFiles := map[string][]byte{
		filepath.Join(tmpRootDir, "test", "test1", "go.mod"): []byte("module \"go.opentelemetry.io/test/test1\"\n\ngo 1.16\n\n" +
			"require (\n\t\"go.opentelemetry.io/testroot/v2\" v2.0.0\n)\n"),
		filepath.Join(tmpRootDir, "test", "go.mod"):          []byte("module go.opentelemetry.io/test3\n\ngo 1.16\n"),
		filepath.Join(tmpRootDir, "go.mod"):                  []byte("module go.opentelemetry.io/testroot/v2\n\ngo 1.16\n"),
		filepath.Join(tmpRootDir, "test", "test2", "go.mod"): []byte("module \"go.opentelemetry.io/test/testexcluded\"\n\ngo 1.16\n"),
	}

	require.NoError(t, sharedtest.WriteTempFiles(modFiles), "could not create go mod file tree")

	expected := ModulePathMap{
		"go.opentelemetry.io/test/test1":        ModuleFilePath(filepath.Join(tmpRootDir, "test", "test1", "go.mod")),
		"go.opentelemetry.io/test3":             ModuleFilePath(filepath.Join(tmpRootDir, "test", "go.mod")),
		"go.opentelemetry.io/testroot/v2":       ModuleFilePath(filepath.Join(tmpRootDir, "go.mod")),
		"go.opentelemetry.io/test/testexcluded": ModuleFilePath(filepath.Join(tmpRootDir, "test", "test2", "go.mod")),
	}

	result, err := newAllModulePathMap(tmpRootDir)

	require.NoError(t, err)
	assert.Equal(t, expected, result)
}
