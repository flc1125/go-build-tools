// Copyright The OpenTelemetry Authors
// Modified by flc1125 for github.com/flc1125/go-build-tools.
// SPDX-License-Identifier: Apache-2.0

package crosslink

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBuildDependencyGraph(t *testing.T) {
	tests := []struct {
		testName string
		mockDir  string
		config   RunConfig
		expected map[string][]string
	}{
		{
			testName: "testSimple",
			mockDir:  "testSimple",
			config:   DefaultRunConfig(),
			expected: map[string][]string{
				"github.com/flc1125/go-build-tools/crosslink/testroot": {
					"github.com/flc1125/go-build-tools/crosslink/testroot/testA",
					"github.com/flc1125/go-build-tools/crosslink/testroot/testB",
				},
				"github.com/flc1125/go-build-tools/crosslink/testroot/testA": {"github.com/flc1125/go-build-tools/crosslink/testroot/testB"},
				"github.com/flc1125/go-build-tools/crosslink/testroot/testB": {},
			},
		},
		{
			testName: "testCyclic",
			mockDir:  "testCyclic",
			config:   DefaultRunConfig(),
			expected: map[string][]string{
				"github.com/flc1125/go-build-tools/crosslink/testroot": {
					"github.com/flc1125/go-build-tools/crosslink/testroot/testA",
					"github.com/flc1125/go-build-tools/crosslink/testroot/testB",
				},
				"github.com/flc1125/go-build-tools/crosslink/testroot/testA": {
					"github.com/flc1125/go-build-tools/crosslink/testroot/testB",
					"github.com/flc1125/go-build-tools/crosslink/testroot",
				},
				// b has req on root but not necessary to write out with current comparison logic
				"github.com/flc1125/go-build-tools/crosslink/testroot/testB": {
					"github.com/flc1125/go-build-tools/crosslink/testroot/testA",
					"github.com/flc1125/go-build-tools/crosslink/testroot",
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			tmpRootDir := createTempTestDir(t, test.mockDir)
			err := renameGoMod(tmpRootDir)
			require.NoError(t, err, "error renaming gomod files")

			test.config.RootPath = tmpRootDir

			rootModulePath, err := identifyRootModule(test.config.RootPath)
			require.NoError(t, err, "error identifying root module")

			receivedMap, err := buildDependencyGraph(test.config, rootModulePath)
			require.NoError(t, err, "error building dependency graph")

			assert.Len(t, receivedMap, len(test.expected), "module count does not match")
			for modName, moduleInfoActual := range receivedMap {
				requiredReplaceStatementsActual := moduleInfoActual.requiredReplaceStatements
				expectedReplaceStatements := test.expected[modName]
				// verify that the amount of replace statements in module match the amount that are in module.
				assert.Len(t, requiredReplaceStatementsActual, len(expectedReplaceStatements),
					"module %q: unexpected replace statement count; expected: %v; actual: %v",
					modName, expectedReplaceStatements, requiredReplaceStatementsActual)
				// ensure that they contain the same values
				for _, expectedReplaceStatement := range expectedReplaceStatements {
					assert.Contains(t, requiredReplaceStatementsActual, expectedReplaceStatement,
						"module %q: expected replace statement is missing", modName)
				}
			}
		})
	}
}
