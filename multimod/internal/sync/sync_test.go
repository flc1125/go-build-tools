// Copyright The OpenTelemetry Authors
// Modified by flc1125 for github.com/flc1125/go-build-tools.
// SPDX-License-Identifier: Apache-2.0

package sync

import (
	"bytes"
	"errors"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/flc1125/go-build-tools/multimod/internal/shared"
	"github.com/flc1125/go-build-tools/multimod/internal/shared/sharedtest"
)

var testDataDir, _ = filepath.Abs("./test_data")

// TestMain performs setup for the tests and suppress printing logs.
func TestMain(m *testing.M) {
	log.SetOutput(io.Discard)
	os.Exit(m.Run())
}

func TestNewSync(t *testing.T) {
	testName := "new_sync"
	versionsYamlDir := filepath.Join(testDataDir, testName)

	myVersioningFilename := filepath.Join(versionsYamlDir, "versions_valid.yaml")
	otherVersioningFilename := filepath.Join(versionsYamlDir, "other_versions_valid.yaml")

	tmpRootDir := t.TempDir()

	modFiles := map[string][]byte{
		filepath.Join(tmpRootDir, "my", "test", "test1", "go.mod"): []byte("module github.com/flc1125/go-build-tools/multimod/internal/sync/test/test1\n\n" +
			"go 1.16\n\n" +
			"require (\n\t" +
			"github.com/flc1125/go-build-tools/multimod/internal/sync/test/test2 v1.2.3-RC1+meta\n\t" +
			"go.opentelemetry.io/other/test/test1 v1.0.0-old\n\t" +
			"go.opentelemetry.io/other/testroot/v2 v2.2.2\n" +
			")"),
		filepath.Join(tmpRootDir, "my", "test", "test2", "go.mod"): []byte("module github.com/flc1125/go-build-tools/multimod/internal/sync/test/test2\n\n" +
			"go 1.16\n\n" +
			"require (\n\t" +
			"github.com/flc1125/go-build-tools/multimod/internal/sync/test/test1 v1.2.3-RC1+meta\n\t" +
			"go.opentelemetry.io/other/test/test1 v1.0.0-old\n\t" +
			"go.opentelemetry.io/other/testroot/v2 v2.2.2\n" +
			")"),
		filepath.Join(tmpRootDir, "my", "test", "go.mod"): []byte("module github.com/flc1125/go-build-tools/multimod/internal/sync/test3\n\n" +
			"go 1.16\n\n" +
			"require (\n\t" +
			"github.com/flc1125/go-build-tools/multimod/internal/sync/test/test1 v1.2.3-RC1+meta\n\t" +
			"github.com/flc1125/go-build-tools/multimod/internal/sync/test/test2 v1.2.3-RC1+meta\n\t" +
			"go.opentelemetry.io/other/test2 v0.1.0-old\n" +
			")"),
		filepath.Join(tmpRootDir, "my", "go.mod"): []byte("module github.com/flc1125/go-build-tools/multimod/internal/sync/testroot/v2\n\n" +
			"go 1.16\n\n" +
			"require (\n\t" +
			"github.com/flc1125/go-build-tools/multimod/internal/sync/test/test1 v1.2.3-RC1+meta\n\t" +
			"github.com/flc1125/go-build-tools/multimod/internal/sync/test/test2 v1.2.3-RC1+meta\n" +
			"go.opentelemetry.io/other/test/test1 v1.0.0-old\n\t" +
			")"),
	}

	require.NoError(t, sharedtest.WriteTempFiles(modFiles), "could not create go mod file tree")

	expectedMyModuleVersioning := shared.ModuleVersioning{
		ModSetMap: shared.ModuleSetMap{
			"my-mod-set-1": shared.ModuleSet{
				Version: "v1.2.3-RC1+meta",
				Modules: []shared.ModulePath{
					"github.com/flc1125/go-build-tools/multimod/internal/sync/test/test1",
					"github.com/flc1125/go-build-tools/multimod/internal/sync/test/test2",
				},
			},
			"my-mod-set-2": shared.ModuleSet{
				Version: "v0.1.0",
				Modules: []shared.ModulePath{
					"github.com/flc1125/go-build-tools/multimod/internal/sync/test3",
				},
			},
			"my-mod-set-3": shared.ModuleSet{
				Version: "v2.2.2",
				Modules: []shared.ModulePath{
					"github.com/flc1125/go-build-tools/multimod/internal/sync/testroot/v2",
				},
			},
		},
		ModPathMap: shared.ModulePathMap{
			"github.com/flc1125/go-build-tools/multimod/internal/sync/test/test1":  shared.ModuleFilePath(filepath.Join(tmpRootDir, "my", "test", "test1", "go.mod")),
			"github.com/flc1125/go-build-tools/multimod/internal/sync/test/test2":  shared.ModuleFilePath(filepath.Join(tmpRootDir, "my", "test", "test2", "go.mod")),
			"github.com/flc1125/go-build-tools/multimod/internal/sync/test3":       shared.ModuleFilePath(filepath.Join(tmpRootDir, "my", "test", "go.mod")),
			"github.com/flc1125/go-build-tools/multimod/internal/sync/testroot/v2": shared.ModuleFilePath(filepath.Join(tmpRootDir, "my", "go.mod")),
		},
		ModInfoMap: shared.ModuleInfoMap{
			"github.com/flc1125/go-build-tools/multimod/internal/sync/test/test1": shared.ModuleInfo{
				ModuleSetName: "my-mod-set-1",
				Version:       "v1.2.3-RC1+meta",
			},
			"github.com/flc1125/go-build-tools/multimod/internal/sync/test/test2": shared.ModuleInfo{
				ModuleSetName: "my-mod-set-1",
				Version:       "v1.2.3-RC1+meta",
			},
			"github.com/flc1125/go-build-tools/multimod/internal/sync/test3": shared.ModuleInfo{
				ModuleSetName: "my-mod-set-2",
				Version:       "v0.1.0",
			},
			"github.com/flc1125/go-build-tools/multimod/internal/sync/testroot/v2": shared.ModuleInfo{
				ModuleSetName: "my-mod-set-3",
				Version:       "v2.2.2",
			},
		},
	}

	testCases := []struct {
		modSetName          string
		expectedOtherModSet shared.ModuleSet
	}{
		{
			modSetName: "other-mod-set-1",
			expectedOtherModSet: shared.ModuleSet{
				Version: "v1.2.3-RC1+meta",
				Modules: []shared.ModulePath{"go.opentelemetry.io/other/test/test1"},
			},
		},
		{
			modSetName: "other-mod-set-2",
			expectedOtherModSet: shared.ModuleSet{
				Version: "v0.1.0",
				Modules: []shared.ModulePath{"go.opentelemetry.io/other/test2"},
			},
		},
		{
			modSetName: "other-mod-set-3",
			expectedOtherModSet: shared.ModuleSet{
				Version: "v2.2.2",
				Modules: []shared.ModulePath{"go.opentelemetry.io/other/testroot/v2"},
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.modSetName, func(t *testing.T) {
			actual, err := newSync(
				myVersioningFilename,
				otherVersioningFilename,
				tc.modSetName,
				tmpRootDir,
				"",
			)
			require.NoError(t, err)

			require.IsType(t, sync{}, actual)
			assert.Equal(t, tc.modSetName, actual.OtherModuleSetName)
			assert.Equal(t, tc.expectedOtherModSet, actual.OtherModuleSet)
			assert.Equal(t, expectedMyModuleVersioning, actual.MyModuleVersioning)
		})
	}
}

type roundTripFunc func(req *http.Request) *http.Response

func (f roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req), nil
}

type trackingBody struct {
	reader   io.Reader
	readErr  error
	closeErr error
	closed   bool
}

func (b *trackingBody) Read(p []byte) (int, error) {
	if b.readErr != nil {
		return 0, b.readErr
	}
	return b.reader.Read(p)
}

func (b *trackingBody) Close() error {
	b.closed = true
	return b.closeErr
}

func TestCheckRetryHandlesNilResponse(t *testing.T) {
	retry, err := checkRetry(t.Context(), nil, errors.New("network error"))

	require.NoError(t, err)
	assert.True(t, retry)
}

func TestParseVersionInfoClosesResponseBody(t *testing.T) {
	body := &trackingBody{reader: bytes.NewBufferString(`{"Version":"v1.2.3"}`)}
	s := sync{client: &http.Client{Transport: roundTripFunc(func(*http.Request) *http.Response {
		return &http.Response{StatusCode: http.StatusOK, Body: body, Header: make(http.Header)}
	})}}

	version, err := s.parseVersionInfo(t.Context(), "example.com/module", "main")

	require.NoError(t, err)
	assert.Equal(t, "v1.2.3", version)
	assert.True(t, body.closed)
}

func TestParseVersionInfoReturnsReadError(t *testing.T) {
	readErr := errors.New("read failed")
	body := &trackingBody{readErr: readErr}
	s := sync{client: &http.Client{Transport: roundTripFunc(func(*http.Request) *http.Response {
		return &http.Response{StatusCode: http.StatusOK, Body: body, Header: make(http.Header)}
	})}}

	_, err := s.parseVersionInfo(t.Context(), "example.com/module", "main")

	assert.ErrorIs(t, err, readErr)
	assert.ErrorContains(t, err, "failed to read response")
	assert.True(t, body.closed)
}

func TestParseVersionInfoReturnsCloseError(t *testing.T) {
	closeErr := errors.New("close failed")
	body := &trackingBody{reader: bytes.NewBufferString(`{"Version":"v1.2.3"}`), closeErr: closeErr}
	s := sync{client: &http.Client{Transport: roundTripFunc(func(*http.Request) *http.Response {
		return &http.Response{StatusCode: http.StatusOK, Body: body, Header: make(http.Header)}
	})}}

	_, err := s.parseVersionInfo(t.Context(), "example.com/module", "main")

	assert.ErrorIs(t, err, closeErr)
	assert.ErrorContains(t, err, "failed to close response")
	assert.True(t, body.closed)
}

func TestUpdateAllGoModFilesWithCommitHash(t *testing.T) {
	testName := "update_all_go_mod_files_with_commit_hash"
	versionsYamlDir := filepath.Join(testDataDir, testName)

	myVersioningFilename := filepath.Join(versionsYamlDir, "versions_valid.yaml")
	otherVersioningFilename := filepath.Join(versionsYamlDir, "other_versions_valid.yaml")

	testCases := []struct {
		modSetName             string
		expectedOutputModFiles map[string][]byte
		commit                 string
		client                 *http.Client
		expectedErr            string
	}{
		{
			modSetName: "other-mod-set-1",
			expectedOutputModFiles: map[string][]byte{
				filepath.Join("my", "test", "test1", "go.mod"): []byte("module github.com/flc1125/go-build-tools/multimod/internal/sync/test/test1\n\n" +
					"go 1.16\n\n" +
					"require (\n\t" +
					"github.com/flc1125/go-build-tools/multimod/internal/sync/test/test2 v1.2.3-RC1+meta\n\t" +
					"go.opentelemetry.io/other/test/test1 v1.2.4-RC1+meta\n\t" +
					"go.opentelemetry.io/other/testroot/v2 v2.2.2\n" +
					")"),
				filepath.Join("my", "test", "test2", "go.mod"): []byte("module github.com/flc1125/go-build-tools/multimod/internal/sync/test/test2\n\n" +
					"go 1.16\n\n" +
					"require (\n\t" +
					"github.com/flc1125/go-build-tools/multimod/internal/sync/test/test1 v1.2.3-RC1+meta\n\t" +
					"go.opentelemetry.io/other/test/test1 v1.2.4-RC1+meta\n\t" +
					"go.opentelemetry.io/other/testroot/v2 v2.2.2\n" +
					")"),
				filepath.Join("my", "test", "go.mod"): []byte("module github.com/flc1125/go-build-tools/multimod/internal/sync/test3\n\n" +
					"go 1.16\n\n" +
					"require (\n\t" +
					"github.com/flc1125/go-build-tools/multimod/internal/sync/test/test1 v1.2.3-RC1+meta\n\t" +
					"github.com/flc1125/go-build-tools/multimod/internal/sync/test/test2 v1.2.3-RC1+meta\n\t" +
					"go.opentelemetry.io/other/test2 v0.1.0-old\n" +
					")"),
				filepath.Join("my", "go.mod"): []byte("module github.com/flc1125/go-build-tools/multimod/internal/sync/testroot/v2\n\n" +
					"go 1.16\n\n" +
					"require (\n\t" +
					"github.com/flc1125/go-build-tools/multimod/internal/sync/test/test1 v1.2.3-RC1+meta\n\t" +
					"github.com/flc1125/go-build-tools/multimod/internal/sync/test/test2 v1.2.3-RC1+meta\n\t" +
					"go.opentelemetry.io/other/test/test1 v1.2.4-RC1+meta\n" +
					")"),
			},
			commit: "main",
			client: &http.Client{
				Transport: roundTripFunc(func(*http.Request) *http.Response {
					return &http.Response{
						StatusCode: 200,
						Body:       io.NopCloser(bytes.NewBufferString(`{"Version":"v1.2.4-RC1+meta","Time":"2023-10-12T21:04:47Z","Origin":{"VCS":"git","URL":"https://github.com/opentelemetry-go-build-tools/multimod/internal/sync/test/test2","Hash":"35cee309328ac126861ae6f554971aeb85a08bba"}}`)),
						Header:     make(http.Header),
					}
				}),
			},
		},
		{
			modSetName: "other-mod-set-1",
			commit:     "main",
			client: &http.Client{
				Transport: roundTripFunc(func(*http.Request) *http.Response {
					return &http.Response{
						StatusCode: 500,
						Body:       io.NopCloser(bytes.NewBufferString(`server error`)),
						Header:     make(http.Header),
					}
				}),
			},
			expectedErr: "request to proxy for package \"go.opentelemetry.io/other/test/test1\" at version \"main\" failed with status 500",
		},

		{
			modSetName: "other-mod-set-1",
			commit:     "main",
			client: &http.Client{
				Transport: roundTripFunc(func(*http.Request) *http.Response {
					return &http.Response{
						StatusCode: 200,
						Body:       io.NopCloser(bytes.NewBufferString(`invalid message`)),
						Header:     make(http.Header),
					}
				}),
			},
			expectedErr: "failed to unmarshal response",
		},
	}

	for _, tc := range testCases {
		tmpRootDir := t.TempDir()

		modFiles := map[string][]byte{
			filepath.Join(tmpRootDir, "my", "test", "test1", "go.mod"): []byte("module github.com/flc1125/go-build-tools/multimod/internal/sync/test/test1\n\n" +
				"go 1.16\n\n" +
				"require (\n\t" +
				"github.com/flc1125/go-build-tools/multimod/internal/sync/test/test2 v1.2.3-RC1+meta\n\t" +
				"go.opentelemetry.io/other/test/test1 v1.0.0-old\n\t" +
				"go.opentelemetry.io/other/testroot/v2 v2.2.2\n" +
				")"),
			filepath.Join(tmpRootDir, "my", "test", "test2", "go.mod"): []byte("module github.com/flc1125/go-build-tools/multimod/internal/sync/test/test2\n\n" +
				"go 1.16\n\n" +
				"require (\n\t" +
				"github.com/flc1125/go-build-tools/multimod/internal/sync/test/test1 v1.2.3-RC1+meta\n\t" +
				"go.opentelemetry.io/other/test/test1 v1.0.0-old\n\t" +
				"go.opentelemetry.io/other/testroot/v2 v2.2.2\n" +
				")"),
			filepath.Join(tmpRootDir, "my", "test", "go.mod"): []byte("module github.com/flc1125/go-build-tools/multimod/internal/sync/test3\n\n" +
				"go 1.16\n\n" +
				"require (\n\t" +
				"github.com/flc1125/go-build-tools/multimod/internal/sync/test/test1 v1.2.3-RC1+meta\n\t" +
				"github.com/flc1125/go-build-tools/multimod/internal/sync/test/test2 v1.2.3-RC1+meta\n\t" +
				"go.opentelemetry.io/other/test2 v0.1.0-old\n" +
				")"),
			filepath.Join(tmpRootDir, "my", "go.mod"): []byte("module github.com/flc1125/go-build-tools/multimod/internal/sync/testroot/v2\n\n" +
				"go 1.16\n\n" +
				"require (\n\t" +
				"github.com/flc1125/go-build-tools/multimod/internal/sync/test/test1 v1.2.3-RC1+meta\n\t" +
				"github.com/flc1125/go-build-tools/multimod/internal/sync/test/test2 v1.2.3-RC1+meta\n\t" +
				"go.opentelemetry.io/other/test/test1 v1.0.0-old\n" +
				")"),
		}

		t.Run(tc.modSetName, func(t *testing.T) {
			require.NoError(t, sharedtest.WriteTempFiles(modFiles), "could not create go mod file tree")

			s, err := newSync(
				myVersioningFilename,
				otherVersioningFilename,
				tc.modSetName,
				tmpRootDir,
				tc.commit,
			)
			s.client = tc.client
			require.NoError(t, err)

			err = s.updateAllGoModFiles(t.Context())
			if len(tc.expectedErr) == 0 {
				require.NoError(t, err)
			} else {
				require.ErrorContains(t, err, tc.expectedErr)
			}

			for modFilePathSuffix, expectedByteOutput := range tc.expectedOutputModFiles {
				actual, err := os.ReadFile(filepath.Clean(filepath.Join(tmpRootDir, modFilePathSuffix)))
				require.NoError(t, err)

				assert.Equal(t, expectedByteOutput, actual)
			}
		})
	}
}

func TestUpdateAllGoModFiles(t *testing.T) {
	testName := "update_all_go_mod_files"
	versionsYamlDir := filepath.Join(testDataDir, testName)

	myVersioningFilename := filepath.Join(versionsYamlDir, "versions_valid.yaml")
	otherVersioningFilename := filepath.Join(versionsYamlDir, "other_versions_valid.yaml")

	testCases := []struct {
		modSetName             string
		expectedOutputModFiles map[string][]byte
	}{
		{
			modSetName: "other-mod-set-1",
			expectedOutputModFiles: map[string][]byte{
				filepath.Join("my", "test", "test1", "go.mod"): []byte("module github.com/flc1125/go-build-tools/multimod/internal/sync/test/test1\n\n" +
					"go 1.16\n\n" +
					"require (\n\t" +
					"github.com/flc1125/go-build-tools/multimod/internal/sync/test/test2 v1.2.3-RC1+meta\n\t" +
					"go.opentelemetry.io/other/test/test1 v1.2.3-RC1+meta\n\t" +
					"go.opentelemetry.io/other/testroot/v2 v2.2.2\n" +
					")"),
				filepath.Join("my", "test", "test2", "go.mod"): []byte("module github.com/flc1125/go-build-tools/multimod/internal/sync/test/test2\n\n" +
					"go 1.16\n\n" +
					"require (\n\t" +
					"github.com/flc1125/go-build-tools/multimod/internal/sync/test/test1 v1.2.3-RC1+meta\n\t" +
					"go.opentelemetry.io/other/test/test1 v1.2.3-RC1+meta\n\t" +
					"go.opentelemetry.io/other/testroot/v2 v2.2.2\n" +
					")"),
				filepath.Join("my", "test", "go.mod"): []byte("module github.com/flc1125/go-build-tools/multimod/internal/sync/test3\n\n" +
					"go 1.16\n\n" +
					"require (\n\t" +
					"github.com/flc1125/go-build-tools/multimod/internal/sync/test/test1 v1.2.3-RC1+meta\n\t" +
					"github.com/flc1125/go-build-tools/multimod/internal/sync/test/test2 v1.2.3-RC1+meta\n\t" +
					"go.opentelemetry.io/other/test2 v0.1.0-old\n" +
					")"),
				filepath.Join("my", "go.mod"): []byte("module github.com/flc1125/go-build-tools/multimod/internal/sync/testroot/v2\n\n" +
					"go 1.16\n\n" +
					"require (\n\t" +
					"github.com/flc1125/go-build-tools/multimod/internal/sync/test/test1 v1.2.3-RC1+meta\n\t" +
					"github.com/flc1125/go-build-tools/multimod/internal/sync/test/test2 v1.2.3-RC1+meta\n\t" +
					"go.opentelemetry.io/other/test/test1 v1.2.3-RC1+meta\n" +
					")"),
				filepath.Join("my", "test", "testexcluded", "go.mod"): []byte("module go.opentelemetry.io/my/test/testexcluded\n\n" +
					"go 1.16\n\n" +
					"require (\n\t" +
					"github.com/flc1125/go-build-tools/multimod/internal/sync/test/test1 v1.0.0-old\n\t" +
					"go.opentelemetry.io/other/testroot/v2 v1.0.0\n\t" +
					"go.opentelemetry.io/other/test/test1 v1.2.3-RC1+meta\n" +
					")"),
			},
		},
		{
			modSetName: "other-mod-set-2",
			expectedOutputModFiles: map[string][]byte{
				filepath.Join("my", "test", "test1", "go.mod"): []byte("module github.com/flc1125/go-build-tools/multimod/internal/sync/test/test1\n\n" +
					"go 1.16\n\n" +
					"require (\n\t" +
					"github.com/flc1125/go-build-tools/multimod/internal/sync/test/test2 v1.2.3-RC1+meta\n\t" +
					"go.opentelemetry.io/other/test/test1 v1.0.0-old\n\t" +
					"go.opentelemetry.io/other/testroot/v2 v2.2.2\n" +
					")"),
				filepath.Join("my", "test", "test2", "go.mod"): []byte("module github.com/flc1125/go-build-tools/multimod/internal/sync/test/test2\n\n" +
					"go 1.16\n\n" +
					"require (\n\t" +
					"github.com/flc1125/go-build-tools/multimod/internal/sync/test/test1 v1.2.3-RC1+meta\n\t" +
					"go.opentelemetry.io/other/test/test1 v1.0.0-old\n\t" +
					"go.opentelemetry.io/other/testroot/v2 v2.2.2\n" +
					")"),
				filepath.Join("my", "test", "go.mod"): []byte("module github.com/flc1125/go-build-tools/multimod/internal/sync/test3\n\n" +
					"go 1.16\n\n" +
					"require (\n\t" +
					"github.com/flc1125/go-build-tools/multimod/internal/sync/test/test1 v1.2.3-RC1+meta\n\t" +
					"github.com/flc1125/go-build-tools/multimod/internal/sync/test/test2 v1.2.3-RC1+meta\n\t" +
					"go.opentelemetry.io/other/test2 v0.1.0\n" +
					")"),
				filepath.Join("my", "go.mod"): []byte("module github.com/flc1125/go-build-tools/multimod/internal/sync/testroot/v2\n\n" +
					"go 1.16\n\n" +
					"require (\n\t" +
					"github.com/flc1125/go-build-tools/multimod/internal/sync/test/test1 v1.2.3-RC1+meta\n\t" +
					"github.com/flc1125/go-build-tools/multimod/internal/sync/test/test2 v1.2.3-RC1+meta\n\t" +
					"go.opentelemetry.io/other/test/test1 v1.0.0-old\n" +
					")"),
				filepath.Join("my", "test", "testexcluded", "go.mod"): []byte("module go.opentelemetry.io/my/test/testexcluded\n\n" +
					"go 1.16\n\n" +
					"require (\n\t" +
					"github.com/flc1125/go-build-tools/multimod/internal/sync/test/test1 v1.0.0-old\n\t" +
					"go.opentelemetry.io/other/testroot/v2 v1.0.0\n\t" +
					"go.opentelemetry.io/other/test/test1 v1.0.0-old\n" +
					")"),
			},
		},
		{
			modSetName: "other-mod-set-3",
			expectedOutputModFiles: map[string][]byte{
				filepath.Join("my", "test", "test1", "go.mod"): []byte("module github.com/flc1125/go-build-tools/multimod/internal/sync/test/test1\n\n" +
					"go 1.16\n\n" +
					"require (\n\t" +
					"github.com/flc1125/go-build-tools/multimod/internal/sync/test/test2 v1.2.3-RC1+meta\n\t" +
					"go.opentelemetry.io/other/test/test1 v1.0.0-old\n\t" +
					"go.opentelemetry.io/other/testroot/v2 v2.2.2\n" +
					")"),
				filepath.Join("my", "test", "test2", "go.mod"): []byte("module github.com/flc1125/go-build-tools/multimod/internal/sync/test/test2\n\n" +
					"go 1.16\n\n" +
					"require (\n\t" +
					"github.com/flc1125/go-build-tools/multimod/internal/sync/test/test1 v1.2.3-RC1+meta\n\t" +
					"go.opentelemetry.io/other/test/test1 v1.0.0-old\n\t" +
					"go.opentelemetry.io/other/testroot/v2 v2.2.2\n" +
					")"),
				filepath.Join("my", "test", "go.mod"): []byte("module github.com/flc1125/go-build-tools/multimod/internal/sync/test3\n\n" +
					"go 1.16\n\n" +
					"require (\n\t" +
					"github.com/flc1125/go-build-tools/multimod/internal/sync/test/test1 v1.2.3-RC1+meta\n\t" +
					"github.com/flc1125/go-build-tools/multimod/internal/sync/test/test2 v1.2.3-RC1+meta\n\t" +
					"go.opentelemetry.io/other/test2 v0.1.0-old\n" +
					")"),
				filepath.Join("my", "go.mod"): []byte("module github.com/flc1125/go-build-tools/multimod/internal/sync/testroot/v2\n\n" +
					"go 1.16\n\n" +
					"require (\n\t" +
					"github.com/flc1125/go-build-tools/multimod/internal/sync/test/test1 v1.2.3-RC1+meta\n\t" +
					"github.com/flc1125/go-build-tools/multimod/internal/sync/test/test2 v1.2.3-RC1+meta\n\t" +
					"go.opentelemetry.io/other/test/test1 v1.0.0-old\n" +
					")"),
				filepath.Join("my", "test", "testexcluded", "go.mod"): []byte("module go.opentelemetry.io/my/test/testexcluded\n\n" +
					"go 1.16\n\n" +
					"require (\n\t" +
					"github.com/flc1125/go-build-tools/multimod/internal/sync/test/test1 v1.0.0-old\n\t" +
					"go.opentelemetry.io/other/testroot/v2 v2.2.2\n\t" +
					"go.opentelemetry.io/other/test/test1 v1.0.0-old\n" +
					")"),
			},
		},
	}

	for _, tc := range testCases {
		tmpRootDir := t.TempDir()

		modFiles := map[string][]byte{
			filepath.Join(tmpRootDir, "my", "test", "test1", "go.mod"): []byte("module github.com/flc1125/go-build-tools/multimod/internal/sync/test/test1\n\n" +
				"go 1.16\n\n" +
				"require (\n\t" +
				"github.com/flc1125/go-build-tools/multimod/internal/sync/test/test2 v1.2.3-RC1+meta\n\t" +
				"go.opentelemetry.io/other/test/test1 v1.0.0-old\n\t" +
				"go.opentelemetry.io/other/testroot/v2 v2.2.2\n" +
				")"),
			filepath.Join(tmpRootDir, "my", "test", "test2", "go.mod"): []byte("module github.com/flc1125/go-build-tools/multimod/internal/sync/test/test2\n\n" +
				"go 1.16\n\n" +
				"require (\n\t" +
				"github.com/flc1125/go-build-tools/multimod/internal/sync/test/test1 v1.2.3-RC1+meta\n\t" +
				"go.opentelemetry.io/other/test/test1 v1.0.0-old\n\t" +
				"go.opentelemetry.io/other/testroot/v2 v2.2.2\n" +
				")"),
			filepath.Join(tmpRootDir, "my", "test", "go.mod"): []byte("module github.com/flc1125/go-build-tools/multimod/internal/sync/test3\n\n" +
				"go 1.16\n\n" +
				"require (\n\t" +
				"github.com/flc1125/go-build-tools/multimod/internal/sync/test/test1 v1.2.3-RC1+meta\n\t" +
				"github.com/flc1125/go-build-tools/multimod/internal/sync/test/test2 v1.2.3-RC1+meta\n\t" +
				"go.opentelemetry.io/other/test2 v0.1.0-old\n" +
				")"),
			filepath.Join(tmpRootDir, "my", "go.mod"): []byte("module github.com/flc1125/go-build-tools/multimod/internal/sync/testroot/v2\n\n" +
				"go 1.16\n\n" +
				"require (\n\t" +
				"github.com/flc1125/go-build-tools/multimod/internal/sync/test/test1 v1.2.3-RC1+meta\n\t" +
				"github.com/flc1125/go-build-tools/multimod/internal/sync/test/test2 v1.2.3-RC1+meta\n\t" +
				"go.opentelemetry.io/other/test/test1 v1.0.0-old\n" +
				")"),
			filepath.Join(tmpRootDir, "my", "test", "testexcluded", "go.mod"): []byte("module go.opentelemetry.io/my/test/testexcluded\n\n" +
				"go 1.16\n\n" +
				"require (\n\t" +
				"github.com/flc1125/go-build-tools/multimod/internal/sync/test/test1 v1.0.0-old\n\t" +
				"go.opentelemetry.io/other/testroot/v2 v1.0.0\n\t" +
				"go.opentelemetry.io/other/test/test1 v1.0.0-old\n" +
				")"),
		}

		t.Run(tc.modSetName, func(t *testing.T) {
			require.NoError(t, sharedtest.WriteTempFiles(modFiles), "could not create go mod file tree")

			s, err := newSync(
				myVersioningFilename,
				otherVersioningFilename,
				tc.modSetName,
				tmpRootDir,
				"",
			)
			require.NoError(t, err)

			err = s.updateAllGoModFiles(t.Context())
			require.NoError(t, err)

			for modFilePathSuffix, expectedByteOutput := range tc.expectedOutputModFiles {
				actual, err := os.ReadFile(filepath.Clean(filepath.Join(tmpRootDir, modFilePathSuffix)))
				require.NoError(t, err)

				assert.Equal(t, expectedByteOutput, actual)
			}
		})
	}
}
