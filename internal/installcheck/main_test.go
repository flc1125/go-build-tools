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
	got, err := writeModule(root, "crosslink", proxy, fixtureVersion)
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
			_, err := writeModule(root, "crosslink", proxy, fixtureVersion)
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

func TestCreateProxyWithUnpublishedRoot(t *testing.T) {
	root := t.TempDir()
	modulePath := "example.com/BuildTools"
	files := map[string]string{
		"go.mod":                    "module " + modulePath + "\n\ngo 1.26.0\n",
		"internal/shared/shared.go": "package shared\n\nfunc Print() {}\n",
	}
	for _, dir := range []string{"crosslink", "gotmpl", "multimod"} {
		gomod := "module " + modulePath + "/" + dir + "\n\ngo 1.26.0\n"
		main := "package main\n\nfunc main() {}\n"
		if dir != "gotmpl" {
			gomod += "\nrequire " + modulePath + " v0.1.0\n"
			main = "package main\n\nimport \"example.com/BuildTools/internal/shared\"\n\nfunc main() { shared.Print() }\n"
		}
		files[dir+"/go.mod"] = gomod
		files[dir+"/main.go"] = main
	}
	for name, data := range files {
		path := filepath.Join(root, name)
		require.NoError(t, os.MkdirAll(filepath.Dir(path), 0o755))
		require.NoError(t, os.WriteFile(path, []byte(data), 0o600))
	}
	for _, args := range [][]string{{"init", "-q"}, {"add", "."}} {
		cmd := exec.CommandContext(t.Context(), "git", args...)
		cmd.Dir = root
		output, err := cmd.CombinedOutput()
		require.NoError(t, err, "%s", output)
	}
	proxy := t.TempDir()
	paths, err := createProxy(root, proxy)
	require.NoError(t, err)
	require.Equal(t, modulePath, paths[0])
	require.Len(t, paths, 4)
	escaped, err := module.EscapePath(modulePath)
	require.NoError(t, err)
	reader, err := zip.OpenReader(filepath.Join(proxy, escaped, "@v", "v0.1.0.zip"))
	require.NoError(t, err)
	defer func() { require.NoError(t, reader.Close()) }()
	names := make([]string, 0, len(reader.File))
	for _, file := range reader.File {
		names = append(names, strings.TrimPrefix(file.Name, modulePath+"@v0.1.0/"))
	}
	require.ElementsMatch(t, []string{"go.mod", "internal/shared/shared.go"}, names)
	tmp := t.TempDir()
	env, err := installEnvironment("go", tmp, proxy, paths)
	require.NoError(t, err)
	env = append(env, "GOPROXY="+fileURL(proxy))
	for _, path := range paths[1:] {
		cmd := exec.CommandContext(t.Context(), "go", "install", "-modcacherw", path+"@latest")
		cmd.Dir = tmp
		cmd.Env = env
		output, installErr := cmd.CombinedOutput()
		require.NoError(t, installErr, "%s", output)
	}
	path := filepath.Join(root, "multimod", "go.mod")
	data := strings.ReplaceAll(files["multimod/go.mod"], "v0.1.0", "v0.2.0")
	require.NoError(t, os.WriteFile(path, []byte(data), 0o600))
	_, err = createProxy(root, t.TempDir())
	require.ErrorContains(t, err, "root module versions do not match")
}
