// Copyright The OpenTelemetry Authors
// Modified by flc1125 for github.com/flc1125/go-build-tools.
// SPDX-License-Identifier: Apache-2.0

package crosslink

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"golang.org/x/mod/modfile"
)

func TestPrune(t *testing.T) {
	lg, _ := zap.NewDevelopment()
	tests := []struct {
		testName string
		mockDir  string
		config   RunConfig
		expected map[string][]byte
	}{
		{
			testName: "testSimplePrune",
			mockDir:  "testSimplePrune",
			config: RunConfig{
				Logger:  lg,
				Prune:   true,
				Verbose: true,
			},
			expected: map[string][]byte{
				"go.mod": []byte("module github.com/flc1125/go-build-tools/crosslink/testroot\n\n" +
					"go 1.20\n\n" +
					"require (\n\t" +
					"github.com/flc1125/go-build-tools/crosslink/testroot/testA v1.0.0\n" +
					")\n" +
					"replace github.com/flc1125/go-build-tools/crosslink/testroot/testA => ./testA\n\n" +
					"replace github.com/flc1125/go-build-tools/crosslink/testroot/testB => ./testB"),
				filepath.Join("testA", "go.mod"): []byte("module github.com/flc1125/go-build-tools/crosslink/testroot/testA\n\n" +
					"go 1.20\n\n" +
					"require (\n\t" +
					"github.com/flc1125/go-build-tools/crosslink/testroot/testB v1.0.0\n" +
					")\n" +
					"replace github.com/flc1125/go-build-tools/crosslink/testroot/testB => ../testB"),
				filepath.Join("testB", "go.mod"): []byte("module github.com/flc1125/go-build-tools/crosslink/testroot/testB\n\n" +
					"go 1.20\n\n"),
			},
		},
	}

	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			tmpRootDir := createTempTestDir(t, test.mockDir)
			err := renameGoMod(tmpRootDir)
			require.NoError(t, err, "error renaming gomod files")

			test.config.RootPath = tmpRootDir
			err = Prune(test.config)
			require.NoError(t, err, "error executing prune")

			for modFilePath, modFilesExpected := range test.expected {
				modFileActual, err := os.ReadFile(filepath.Clean(filepath.Join(tmpRootDir, modFilePath)))
				require.NoError(t, err, "error reading actual mod file %q", modFilePath)

				actual, err := modfile.Parse("go.mod", modFileActual, nil)
				require.NoError(t, err, "error decoding actual mod file %q", modFilePath)
				actual.Cleanup()

				expected, err := modfile.Parse("go.mod", modFilesExpected, nil)
				require.NoError(t, err, "error decoding expected mod file %q", modFilePath)
				expected.Cleanup()

				// replace structs need to be assorted to avoid flaky fails in test
				replaceSortFunc := func(x, y *modfile.Replace) bool {
					return x.Old.Path < y.Old.Path
				}

				diff := cmp.Diff(expected, actual, cmpopts.IgnoreFields(modfile.Replace{}, "Syntax"),
					cmpopts.IgnoreFields(modfile.File{}, "Require", "Exclude", "Retract", "Syntax"),
					cmpopts.SortSlices(replaceSortFunc),
				)
				assert.Empty(t, diff, "mod file %q: Replace{} mismatch (-want +got)", modFilePath)
			}
		})
	}
}

func TestPruneReplace(t *testing.T) {
	testName := "testPrune"

	tmpRootDir := createTempTestDir(t, testName)
	err := renameGoMod(tmpRootDir)
	require.NoError(t, err, "error renaming gomod files")

	modContents, err := os.ReadFile(filepath.Clean(filepath.Join(tmpRootDir, "go.mod")))
	require.NoError(t, err, "failed to read mock go.mod file")

	modFile, err := modfile.Parse("go.mod", modContents, nil)
	require.NoError(t, err, "failed to parse mock go.mod file")

	mockRequiredReplaceStatements := map[string]struct{}{
		"github.com/flc1125/go-build-tools/crosslink/testroot/testA": {},
		"github.com/flc1125/go-build-tools/crosslink/testroot/testB": {},
		"github.com/flc1125/go-build-tools/crosslink/testroot/testC": {},
		"github.com/flc1125/go-build-tools/crosslink/testroot/testD": {},
		"github.com/flc1125/go-build-tools/crosslink/testroot/testE": {},
		"github.com/flc1125/go-build-tools/crosslink/testroot/testF": {},
		"github.com/flc1125/go-build-tools/crosslink/testroot/testG": {},
		"github.com/flc1125/go-build-tools/crosslink/testroot/testH": {},
		"github.com/flc1125/go-build-tools/crosslink/testroot/testK": {},
	}

	mockModInfo := newModuleInfo(*modFile)
	mockModInfo.requiredReplaceStatements = mockRequiredReplaceStatements
	lg, _ := zap.NewDevelopment()
	pruneReplace("github.com/flc1125/go-build-tools/crosslink/testroot", mockModInfo, RunConfig{Prune: true, Verbose: true, Logger: lg})

	expectedModFile := []byte("module github.com/flc1125/go-build-tools/crosslink/testroot\n\n" +
		"go 1.20\n\n" +
		"require (\n\t" +
		"github.com/flc1125/go-build-tools/crosslink/testroot/testA v1.0.0\n" +
		"github.com/flc1125/go-build-tools/crosslink/testroot/testB v1.0.0\n" +
		"github.com/flc1125/go-build-tools/crosslink/testroot/testC v1.0.0\n" +
		"github.com/flc1125/go-build-tools/crosslink/testroot/testD v1.0.0\n" +
		"github.com/flc1125/go-build-tools/crosslink/testroot/testE v1.0.0\n" +
		"github.com/flc1125/go-build-tools/crosslink/testroot/testF v1.0.0\n" +
		")\n" +
		"replace github.com/flc1125/go-build-tools/crosslink/testroot/testA => ./testA\n\n" +
		"replace github.com/flc1125/go-build-tools/crosslink/testroot/testB => ./testB\n\n" +
		"replace github.com/flc1125/go-build-tools/crosslink/testroot/testC => ./testC\n\n" +
		"replace github.com/flc1125/go-build-tools/crosslink/testroot/testD => ./testD\n\n" +
		"replace github.com/flc1125/go-build-tools/crosslink/testroot/testE => ./testE\n\n" +
		"replace github.com/flc1125/go-build-tools/crosslink/testroot/testF => ./testF\n\n" +
		"replace github.com/flc1125/go-build-tools/crosslink/testroot/testG => ./testG\n\n" +
		"replace github.com/flc1125/go-build-tools/crosslink/testroot/testH => ./testH\n\n" +
		"replace go.opentelemetry.io/not-a-real-module/testFoo => ./testFoo\n\n" +
		"replace go.opentelemetry.io/fake-module/ => ./fake-module\n\n" +
		"replace github.com/flc1125/go-build-tools/multimod => ../multimod\n\n" +
		"replace foo.opentelemetery.io/bar => ../bar\n\n" +
		"replace github.com/flc1125/go-build-tools/crosslink/testroot/testK => ../crosslinkcopy/testK\n\n")

	expModParse, err := modfile.Parse("go.mod", expectedModFile, nil)
	require.NoError(t, err, "error parsing expected mod file")
	expModParse.Cleanup()

	actual := mockModInfo.moduleContents
	actual.Cleanup()

	// replace structs need to be assorted to avoid flaky fails in test
	replaceSortFunc := func(x, y *modfile.Replace) bool {
		return x.Old.Path < y.Old.Path
	}

	diff := cmp.Diff(*expModParse, actual, cmpopts.IgnoreFields(modfile.Replace{}, "Syntax"),
		cmpopts.IgnoreFields(modfile.File{}, "Require", "Exclude", "Retract", "Syntax", "Module"),
		cmpopts.SortSlices(replaceSortFunc),
	)
	assert.Empty(t, diff, "Replace{} mismatch (-want +got)")
}
