// Copyright The OpenTelemetry Authors
// Modified by flc1125 for github.com/flc1125/go-build-tools.
// SPDX-License-Identifier: Apache-2.0

package crosslink

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"golang.org/x/mod/modfile"
)

func TestCrosslink(t *testing.T) {
	lg, _ := zap.NewDevelopment()

	tests := []struct {
		testName string
		mockDir  string
		config   RunConfig
		expected map[string][]byte
	}{
		{
			testName: "testSimple",
			mockDir:  "testSimple",
			config:   DefaultRunConfig(),
			expected: map[string][]byte{
				"go.mod": []byte("module github.com/flc1125/go-build-tools/crosslink/testroot\n\n" +
					"go 1.20\n\n" +
					"require (\n\t" +
					"github.com/flc1125/go-build-tools/crosslink/testroot/testA v1.0.0\n" +
					")\n" +
					"replace github.com/flc1125/go-build-tools/crosslink/testroot/testA => ./testA\n\n" +
					"replace github.com/flc1125/go-build-tools/crosslink/testroot/testY => ./testY\n\n" +
					"replace github.com/flc1125/go-build-tools/crosslink/testroot/testZ => ./testZ\n\n" +
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
		{
			testName: "testCyclic",
			mockDir:  "testCyclic",
			config:   DefaultRunConfig(),
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
					"replace github.com/flc1125/go-build-tools/crosslink/testroot/testB => ../testB\n\n" +
					"replace github.com/flc1125/go-build-tools/crosslink/testroot => ../"),
				// b has req on root but not necessary to write out with current comparison logic
				filepath.Join("testB", "go.mod"): []byte("module github.com/flc1125/go-build-tools/crosslink/testroot/testB\n\n" +
					"go 1.20\n\n" +
					"replace github.com/flc1125/go-build-tools/crosslink/testroot/testA => ../testA\n\n" +
					"replace github.com/flc1125/go-build-tools/crosslink/testroot => ../\n\n"),
			},
		},
		{
			testName: "testSimpleWithPrune",
			mockDir:  "testSimple",
			config: RunConfig{
				Prune:  true,
				Logger: lg,
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

			err = Crosslink(test.config)
			require.NoError(t, err, "error executing crosslink")

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

func TestOverwrite(t *testing.T) {
	lg, _ := zap.NewDevelopment()

	tests := []struct {
		testName string
		config   RunConfig
		expected map[string][]byte
	}{
		{
			testName: "testOverwrite",
			config: RunConfig{
				Verbose:       true,
				Overwrite:     true,
				ExcludedPaths: map[string]struct{}{},
				Logger:        lg,
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
		{
			testName: "testNoOverwrite",
			config: RunConfig{
				ExcludedPaths: map[string]struct{}{},
				Verbose:       true,
				Logger:        lg,
			},
			expected: map[string][]byte{
				"go.mod": []byte("module github.com/flc1125/go-build-tools/crosslink/testroot\n\n" +
					"go 1.20\n\n" +
					"require (\n\t" +
					"github.com/flc1125/go-build-tools/crosslink/testroot/testA v1.0.0\n" +
					")\n" +
					"replace github.com/flc1125/go-build-tools/crosslink/testroot/testA => ../testA\n\n" +
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
			tmpRootDir := createTempTestDir(t, test.testName)
			err := renameGoMod(tmpRootDir)
			require.NoError(t, err, "error renaming gomod files")

			test.config.RootPath = tmpRootDir

			err = Crosslink(test.config)
			require.NoError(t, err, "error executing crosslink")

			// a mock_test_data_expected folder could be built instead of building expected files by hand.

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
	err := lg.Sync()
	if err != nil {
		fmt.Printf("failed to sync logger:  %v", err)
	}
}

// Testing exclude functionality for prune, overwrite, and no overwrite.
func TestExclude(t *testing.T) {
	testName := "testExclude"
	lg, _ := zap.NewDevelopment()
	tests := []struct {
		testCase string
		config   RunConfig
	}{
		{
			testCase: "Overwrite off",
			config: RunConfig{
				Prune: true,
				ExcludedPaths: map[string]struct{}{
					"github.com/flc1125/go-build-tools/crosslink/testroot/testB": {},
					"github.com/flc1125/go-build-tools/excludeme":                {},
					"github.com/flc1125/go-build-tools/crosslink/testroot/testA": {},
				},
				Verbose: true,
				Logger:  lg,
			},
		},
		{
			testCase: "Overwrite on",
			config: RunConfig{
				Overwrite: true,
				Prune:     true,
				ExcludedPaths: map[string]struct{}{
					"github.com/flc1125/go-build-tools/crosslink/testroot/testB": {},
					"github.com/flc1125/go-build-tools/excludeme":                {},
					"github.com/flc1125/go-build-tools/crosslink/testroot/testA": {},
				},
				Logger:  lg,
				Verbose: true,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.testCase, func(t *testing.T) {
			tmpRootDir := createTempTestDir(t, testName)
			err := renameGoMod(tmpRootDir)
			require.NoError(t, err, "error renaming gomod files")

			test.config.RootPath = tmpRootDir

			err = Crosslink(test.config)
			require.NoError(t, err, "error executing crosslink")

			// a mock_test_data_expected folder could be built instead of building expected files by hand.
			modFilesExpected := map[string][]byte{
				filepath.Join(tmpRootDir, "go.mod"): []byte("module github.com/flc1125/go-build-tools/crosslink/testroot\n\n" +
					"go 1.20\n\n" +
					"require (\n\t" +
					"github.com/flc1125/go-build-tools/crosslink/testroot/testA v1.0.0\n" +
					")\n" +
					"replace github.com/flc1125/go-build-tools/crosslink/testroot/testA => ../testA\n\n" +
					"replace github.com/flc1125/go-build-tools/excludeme => ../excludeme\n\n"),
				filepath.Join(tmpRootDir, "testA", "go.mod"): []byte("module github.com/flc1125/go-build-tools/crosslink/testroot/testA\n\n" +
					"go 1.20\n\n" +
					"require (\n\t" +
					"github.com/flc1125/go-build-tools/crosslink/testroot/testB v1.0.0\n" +
					")\n" +
					"replace github.com/flc1125/go-build-tools/crosslink/testroot/testB => ../testB"),
				filepath.Join(tmpRootDir, "testB", "go.mod"): []byte("module github.com/flc1125/go-build-tools/crosslink/testroot/testB\n\n" +
					"go 1.20\n\n"),
			}

			for modFilePath, modFilesExpected := range modFilesExpected {
				modFileActual, err := os.ReadFile(filepath.Clean(modFilePath))
				require.NoError(t, err, "test case %q: error reading actual mod file %q", test.testCase, modFilePath)

				actual, err := modfile.Parse("go.mod", modFileActual, nil)
				require.NoError(t, err, "test case %q: error decoding actual mod file %q", test.testCase, modFilePath)
				actual.Cleanup()

				expected, err := modfile.Parse("go.mod", modFilesExpected, nil)
				require.NoError(t, err, "test case %q: error decoding expected mod file %q", test.testCase, modFilePath)
				expected.Cleanup()

				// replace structs need to be assorted to avoid flaky fails in test
				replaceSortFunc := func(x, y *modfile.Replace) bool {
					return x.Old.Path < y.Old.Path
				}

				diff := cmp.Diff(expected, actual, cmpopts.IgnoreFields(modfile.Replace{}, "Syntax"),
					cmpopts.IgnoreFields(modfile.File{}, "Require", "Exclude", "Retract", "Syntax"),
					cmpopts.SortSlices(replaceSortFunc),
				)
				assert.Empty(t, diff, "test case %q, mod file %q: Replace{} mismatch (-want +got)", test.testCase, modFilePath)
			}
		})
	}
}

func TestBadRootPath(t *testing.T) {
	lg, _ := zap.NewDevelopment()
	tests := []struct {
		testName      string
		setConfigPath bool
		config        RunConfig
	}{
		{
			testName:      "noGoMod",
			setConfigPath: true,
			config: RunConfig{
				Logger: lg,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			tmpRootDir := t.TempDir()
			if test.setConfigPath {
				test.config.RootPath = tmpRootDir
			}

			err := Crosslink(test.config)
			assert.Error(t, err)
			err = Prune(test.config)
			assert.Error(t, err)
		})
	}
}

// Testing skipping specified go modules.
func TestSkip(t *testing.T) {
	testName := "testSkip"
	lg, _ := zap.NewDevelopment()
	tests := []struct {
		testCase string
		config   RunConfig
	}{
		{
			testCase: "No skipped go.mod",
			config: RunConfig{
				Prune:   true,
				Verbose: true,
				Logger:  lg,
			},
		},
		{
			testCase: "Include skipped go.mod",
			config: RunConfig{
				Prune: true,
				SkippedPaths: map[string]struct{}{
					"testA/go.mod": {},
				},
				Logger:  lg,
				Verbose: true,
			},
		},
		{
			testCase: "Include non-existent go.mod",
			config: RunConfig{
				Prune: true,
				SkippedPaths: map[string]struct{}{
					"non-existent/go.mod": {},
				},
				Logger:  lg,
				Verbose: true,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.testCase, func(t *testing.T) {
			tmpRootDir := createTempTestDir(t, testName)
			err := renameGoMod(tmpRootDir)
			require.NoError(t, err, "error renaming gomod files")

			test.config.RootPath = tmpRootDir

			err = Crosslink(test.config)
			require.NoError(t, err, "error message on execution %s")

			modFilesExpected := map[string][]byte{
				filepath.Join(tmpRootDir, "go.mod"): []byte("module github.com/flc1125/go-build-tools/crosslink/testroot\n\n" +
					"go 1.20\n\n" +
					"require (\n\t" +
					"github.com/flc1125/go-build-tools/crosslink/testroot/testA v1.0.0\n" +
					")\n" +
					"replace github.com/flc1125/go-build-tools/crosslink/testroot/testA => ./testA\n\n" +
					"replace github.com/flc1125/go-build-tools/excludeme => ../excludeme\n\n"),
				filepath.Join(tmpRootDir, "testA", "go.mod"): []byte("module github.com/flc1125/go-build-tools/crosslink/testroot/testA\n\n" +
					"go 1.20\n\n" +
					"require (\n\t" +
					"github.com/flc1125/go-build-tools/crosslink/testroot/testB v1.0.0\n" +
					")\n" +
					"replace github.com/flc1125/go-build-tools/crosslink/testroot/testB => ../testB"),
			}

			for modFilePath, modFilesExpected := range modFilesExpected {
				shouldDiffer := false
				for path := range test.config.SkippedPaths {
					if strings.HasSuffix(modFilePath, path) {
						shouldDiffer = true
					}
				}
				modFileActual, err := os.ReadFile(filepath.Clean(modFilePath))
				require.NoError(t, err, "test case %q: error reading actual mod file %q", test.testCase, modFilePath)

				actual, err := modfile.Parse("go.mod", modFileActual, nil)
				require.NoError(t, err, "test case %q: error decoding actual mod file %q", test.testCase, modFilePath)
				actual.Cleanup()

				expected, err := modfile.Parse("go.mod", modFilesExpected, nil)
				require.NoError(t, err, "test case %q: error decoding expected mod file %q", test.testCase, modFilePath)
				expected.Cleanup()

				// replace structs need to be assorted to avoid flaky fails in test
				replaceSortFunc := func(x, y *modfile.Replace) bool {
					return x.Old.Path < y.Old.Path
				}

				diff := cmp.Diff(expected, actual, cmpopts.IgnoreFields(modfile.Replace{}, "Syntax"),
					cmpopts.IgnoreFields(modfile.File{}, "Require", "Exclude", "Retract", "Syntax"),
					cmpopts.SortSlices(replaceSortFunc),
				)
				if shouldDiffer {
					assert.Empty(t, diff, "test case %q, mod file %q: Replace{} mismatch (-want +got)", test.testCase, modFilePath)
				}
			}
		})
	}
}
