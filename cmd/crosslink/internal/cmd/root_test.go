// Copyright The OpenTelemetry Authors
// Modified by flc1125 for github.com/flc1125/go-build-tools.
// SPDX-License-Identifier: Apache-2.0

package cmd

import (
	"fmt"
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	cl "github.com/flc1125/go-build-tools/cmd/crosslink/internal"
)

func TestTransform(t *testing.T) {
	tests := []struct {
		testName   string
		inputSlice []string
	}{
		{
			testName: "with items",
			inputSlice: []string{
				"example.com/testA",
				"example.com/testB",
				"example.com/testC",
				"example.com/testD",
				"example.com/testE",
			},
		},
		{
			testName:   "with empty",
			inputSlice: []string{},
		},
	}

	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			actual := transformExclude(test.inputSlice)

			// len must match
			assert.Len(t, actual, len(test.inputSlice))

			// test for existence
			for _, val := range test.inputSlice {
				_, exists := actual[val]
				assert.True(t, exists)
			}
		})
	}
}

var configReset = func() {
	comCfg.runConfig = cl.DefaultRunConfig()
	comCfg.rootCommand.SetArgs([]string{})
}

// Validate run config is valid after pre run.
func TestPreRun(t *testing.T) {
	validRootPath, err := filepath.Abs("../../../..")
	require.NoError(t, err, "could not parse expected root path")
	validProdLogger, err := zap.NewProduction()
	require.NoError(t, err, "failed to create prod logger")

	tests := []struct {
		testName       string
		args           []string
		mockConfig     cl.RunConfig
		expectedConfig cl.RunConfig
	}{
		{
			testName:   "Default Config",
			args:       []string{},
			mockConfig: cl.DefaultRunConfig(),
			expectedConfig: cl.RunConfig{
				Overwrite:     false,
				RootPath:      validRootPath,
				Logger:        validProdLogger,
				ExcludedPaths: make(map[string]struct{}),
			},
		},
		{
			testName: "with overwrite",
			mockConfig: cl.RunConfig{
				Overwrite: true,
			},
			expectedConfig: cl.RunConfig{
				Overwrite: true,
				Verbose:   true,
				RootPath:  validRootPath,
			},
			args: []string{"--overwrite"},
		},
		{
			testName: "with overwrite and verbose=false",
			mockConfig: cl.RunConfig{
				Overwrite: true,
				Verbose:   false,
			},
			expectedConfig: cl.RunConfig{
				Overwrite: true,
				Verbose:   false,
				RootPath:  validRootPath,
			},
			args: []string{"--overwrite", "--verbose=false"},
		},
		{
			testName: "with prune exclusive",
			mockConfig: cl.RunConfig{
				Prune: true,
			},
			expectedConfig: cl.RunConfig{
				Prune:    true,
				RootPath: validRootPath,
			},
			args: []string{"--prune"},
		},
		{
			testName: "with prune exclusive short",
			mockConfig: cl.RunConfig{
				Prune: true,
			},
			expectedConfig: cl.RunConfig{
				Prune:    true,
				RootPath: validRootPath,
			},
			args: []string{"-p"},
		},
		{
			testName: "with verbose exclusive",
			mockConfig: cl.RunConfig{
				Verbose: true,
			},
			expectedConfig: cl.RunConfig{
				Verbose:  true,
				RootPath: validRootPath,
			},
			args: []string{"--verbose"},
		},
		{
			testName: "with verbose exclusive short",
			mockConfig: cl.RunConfig{
				Verbose: true,
			},
			expectedConfig: cl.RunConfig{
				Verbose:  true,
				RootPath: validRootPath,
			},
			args: []string{"-v"},
		},
		{
			testName:   "with good root path",
			mockConfig: cl.DefaultRunConfig(),
			expectedConfig: cl.RunConfig{
				RootPath: validRootPath,
				Logger:   validProdLogger,
			},
			args: []string{fmt.Sprintf("--root=%s", validRootPath)},
		},
	}

	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			t.Cleanup(configReset)
			comCfg.runConfig = test.mockConfig

			err := comCfg.rootCommand.ParseFlags(test.args)
			require.NoError(t, err, "failed to parse flags")

			testPreRun := comCfg.rootCommand.PersistentPreRunE
			err = testPreRun(&comCfg.rootCommand, nil)
			require.NoError(t, err, "pre-run returned error")

			diff := cmp.Diff(test.expectedConfig, comCfg.runConfig, cmpopts.IgnoreFields(cl.RunConfig{}, "Logger", "ExcludedPaths", "SkippedPaths"))
			assert.Empty(t, diff, "test case %q: Config{} mismatch (-want +got)", test.testName)
		})
	}
}

// isolated test because the working directory needs to changed
// and it will keep the happy path test above clean
func TestBadRootPath(t *testing.T) {
	t.Cleanup(configReset)
	mockConfig := cl.DefaultRunConfig()
	args := []string{}

	t.Chdir(t.TempDir())
	comCfg.runConfig = mockConfig

	err := comCfg.rootCommand.ParseFlags(args)
	require.NoError(t, err, "failed to parse flags")

	testPreRun := comCfg.rootCommand.PersistentPreRunE
	err = testPreRun(&comCfg.rootCommand, nil)
	assert.Error(t, err, "Pre Run did not return error")
}
