// Copyright 2026 flc1125
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"archive/zip"
	"encoding/json"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"golang.org/x/mod/module"
)

func TestFileURL(t *testing.T) {
	path := filepath.Join(t.TempDir(), "module proxy")
	u, err := url.Parse(fileURL(path))
	require.NoError(t, err)
	require.Equal(t, "file", u.Scheme)
	require.Empty(t, u.Host)
	require.Equal(t, "/"+strings.TrimPrefix(filepath.ToSlash(path), "/"), u.Path)
	require.Contains(t, u.String(), "%20")
}

func TestWriteModule(t *testing.T) {
	root := t.TempDir()
	dir := filepath.Join(root, "crosslink")
	require.NoError(t, os.MkdirAll(filepath.Join(dir, "nested"), 0o755))
	path := "example.com/BuildTools/crosslink"
	files := map[string]string{
		"go.mod":         "module " + path + "\n\ngo 1.26.0\n",
		"main.go":        "package main\n\nfunc main() {}\n",
		"nested/go.mod":  "module example.com/nested\n\ngo 1.26.0\n",
		"nested/main.go": "package main\n\nfunc main() {}\n",
	}
	for name, content := range files {
		require.NoError(t, os.WriteFile(filepath.Join(dir, name), []byte(content), 0o600))
	}
	for _, args := range [][]string{{"init", "-q"}, {"add", "crosslink"}} {
		cmd := exec.Command("git", args...)
		cmd.Dir = root
		output, err := cmd.CombinedOutput()
		require.NoError(t, err, "%s", output)
	}
	require.NoError(t, os.WriteFile(filepath.Join(dir, "untracked.txt"), []byte("not a release file"), 0o600))
	proxy := t.TempDir()
	got, err := writeModule(root, "crosslink", proxy)
	require.NoError(t, err)
	require.Equal(t, path, got)
	escaped, err := module.EscapePath(path)
	require.NoError(t, err)
	info, err := os.ReadFile(filepath.Join(proxy, escaped, "@latest"))
	require.NoError(t, err)
	var metadata map[string]any
	require.NoError(t, json.Unmarshal(info, &metadata))
	require.Equal(t, fixtureVersion, metadata["Version"])
	reader, err := zip.OpenReader(filepath.Join(proxy, escaped, "@v", fixtureVersion+".zip"))
	require.NoError(t, err)
	defer func() { require.NoError(t, reader.Close()) }()
	names := make([]string, 0, len(reader.File))
	for _, file := range reader.File {
		names = append(names, strings.TrimPrefix(file.Name, path+"@"+fixtureVersion+"/"))
	}
	require.ElementsMatch(t, []string{"go.mod", "main.go"}, names)

	for _, tt := range []struct {
		name    string
		gomod   string
		wantErr bool
	}{
		{name: "installable", gomod: files["go.mod"]},
		{name: "local replacement", gomod: files["go.mod"] + "\nreplace example.com/shared => ../shared\n", wantErr: true},
	} {
		t.Run(tt.name, func(t *testing.T) {
			require.NoError(t, os.WriteFile(filepath.Join(dir, "go.mod"), []byte(tt.gomod), 0o600))
			proxy := t.TempDir()
			_, err := writeModule(root, "crosslink", proxy)
			require.NoError(t, err)
			tmp := t.TempDir()
			env, err := installEnvironment("go", tmp, proxy, []string{path})
			require.NoError(t, err)
			cmd := exec.CommandContext(t.Context(), "go", "install", "-modcacherw", path+"@latest")
			cmd.Dir = tmp
			cmd.Env = env
			output, err := cmd.CombinedOutput()
			if tt.wantErr {
				require.Error(t, err)
				require.Contains(t, string(output), "replace directives")
				return
			}
			require.NoError(t, err, "%s", output)
		})
	}
}
