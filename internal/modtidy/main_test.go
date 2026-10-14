// Copyright 2026 flc1125
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
	"golang.org/x/mod/modfile"
)

func TestTidyModule(t *testing.T) {
	for _, existingReplace := range []bool{false, true} {
		name := "unpublished root"
		if existingReplace {
			name = "existing replacement"
		}
		t.Run(name, func(t *testing.T) {
			t.Setenv("GOPROXY", "off")
			root := t.TempDir()
			modulePath := "example.com/project"
			gomod := "module " + modulePath + "/cli\n\ngo 1.26.0\n\nrequire " + modulePath + " v0.1.0\n"
			files := map[string]string{
				"go.mod":                    "module " + modulePath + "\n\ngo 1.26.0\n",
				"internal/shared/shared.go": "package shared\n\nfunc Print() {}\n",
				"cli/main.go":               "package main\n\nimport \"example.com/project/internal/shared\"\n\nfunc main() { shared.Print() }\n",
			}
			if existingReplace {
				gomod += "\nreplace " + modulePath + " => ../replacement\n"
				files["replacement/go.mod"] = files["go.mod"]
				files["replacement/internal/shared/shared.go"] = files["internal/shared/shared.go"]
				// Only the existing replacement provides the imported package.
				delete(files, "internal/shared/shared.go")
			}
			files["cli/go.mod"] = gomod
			files["cli/go.sum"] = modulePath + " v0.0.1 h1:AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA=\n"
			writeFiles(t, root, files)
			cli := filepath.Join(root, "cli")
			require.NoError(t, tidyModule("go", modulePath, root, cli))
			data, err := os.ReadFile(filepath.Join(cli, "go.mod"))
			require.NoError(t, err)
			f, err := modfile.Parse("go.mod", data, nil)
			require.NoError(t, err)
			require.Len(t, f.Require, 1)
			require.Equal(t, "v0.1.0", f.Require[0].Mod.Version)
			if existingReplace {
				require.Len(t, f.Replace, 1)
				require.Equal(t, "../replacement", f.Replace[0].New.Path)
			} else {
				require.Empty(t, f.Replace)
			}
			sum, err := os.ReadFile(filepath.Join(cli, "go.sum"))
			require.NoError(t, err)
			require.NotContains(t, string(sum), modulePath)
			require.NotContains(t, string(data), root)
		})
	}
}

func TestTidyModuleFailurePreservesFiles(t *testing.T) {
	t.Setenv("GOPROXY", "off")
	root := t.TempDir()
	modulePath := "example.com/project"
	gomod := "module " + modulePath + "/cli\n\ngo 1.26.0\n\nrequire " + modulePath + " v0.1.0\n"
	sum := modulePath + " v0.0.1 h1:AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA=\n"
	writeFiles(t, root, map[string]string{
		"go.mod":      "module " + modulePath + "\n\ngo 1.26.0\n",
		"cli/go.mod":  gomod,
		"cli/go.sum":  sum,
		"cli/main.go": "package main\n\nimport _ \"example.invalid/missing\"\n\nfunc main() {}\n",
	})
	cli := filepath.Join(root, "cli")
	require.Error(t, tidyModule("go", modulePath, root, cli))
	data, err := os.ReadFile(filepath.Join(cli, "go.mod"))
	require.NoError(t, err)
	require.Equal(t, gomod, string(data))
	data, err = os.ReadFile(filepath.Join(cli, "go.sum"))
	require.NoError(t, err)
	require.Equal(t, sum, string(data))
}

func writeFiles(t *testing.T, root string, files map[string]string) {
	t.Helper()
	for name, data := range files {
		path := filepath.Join(root, name)
		require.NoError(t, os.MkdirAll(filepath.Dir(path), 0o755))
		require.NoError(t, os.WriteFile(path, []byte(data), 0o600))
	}
}
