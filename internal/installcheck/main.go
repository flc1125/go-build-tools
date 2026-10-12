// Copyright 2026 flc1125
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"golang.org/x/mod/modfile"
	"golang.org/x/mod/module"
	modzip "golang.org/x/mod/zip"

	"github.com/flc1125/go-build-tools/internal/repo"
)

const fixtureVersion = "v0.0.0-installcheck"

func main() {
	goBinary := flag.String("go", "go", "Go executable for installation checks.")
	flag.Parse()
	if err := run(*goBinary); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(goBinary string) error {
	root, err := repo.FindRoot()
	if err != nil {
		return err
	}
	tmp, err := os.MkdirTemp("", "go-build-tools-install-")
	if err != nil {
		return err
	}
	defer func() { _ = os.RemoveAll(tmp) }()

	proxy := filepath.Join(tmp, "proxy")
	paths, err := createProxy(root, proxy)
	if err != nil {
		return err
	}
	env, err := installEnvironment(goBinary, tmp, proxy, paths)
	if err != nil {
		return err
	}
	for _, path := range paths {
		fmt.Printf("Checking go install %s@latest\n", path)
		cmd := exec.Command(goBinary, "install", "-modcacherw", path+"@latest")
		cmd.Dir = tmp
		cmd.Env = env
		if output, installErr := cmd.CombinedOutput(); installErr != nil {
			return fmt.Errorf("cannot install %s@latest: %w\n%s", path, installErr, output)
		}
	}
	return nil
}

func installEnvironment(goBinary, tmp, proxy string, paths []string) ([]string, error) {
	cmd := exec.Command(goBinary, "env", "-json", "GOPROXY", "GOMODCACHE", "GONOSUMDB")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("cannot read Go settings: %w", err)
	}
	var settings map[string]string
	if err = json.Unmarshal(output, &settings); err != nil {
		return nil, err
	}
	proxies := []string{fileURL(proxy), fileURL(filepath.Join(settings["GOMODCACHE"], "cache", "download")), settings["GOPROXY"]}
	noSumDB := append([]string{}, paths...)
	if settings["GONOSUMDB"] != "" {
		noSumDB = append(noSumDB, settings["GONOSUMDB"])
	}
	return append(os.Environ(),
		"GOWORK=off",
		"GOBIN="+filepath.Join(tmp, "bin"),
		"GOMODCACHE="+filepath.Join(tmp, "modcache"),
		"GOPROXY="+strings.Join(proxies, ","),
		"GONOPROXY=none",
		"GONOSUMDB="+strings.Join(noSumDB, ","),
	), nil
}

func fileURL(path string) string {
	path = filepath.ToSlash(path)
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}
	return (&url.URL{Scheme: "file", Path: path}).String()
}

func createProxy(root, proxy string) ([]string, error) {
	var paths []string
	for _, dir := range []string{"crosslink", "gotmpl", "multimod"} {
		path, err := writeModule(root, dir, proxy)
		if err != nil {
			return nil, fmt.Errorf("cannot prepare %s installation: %w", dir, err)
		}
		paths = append(paths, path)
	}
	return paths, nil
}

func writeModule(root, dir, proxy string) (string, error) {
	gomod, err := os.ReadFile(filepath.Join(root, dir, "go.mod"))
	if err != nil {
		return "", err
	}
	path := modfile.ModulePath(gomod)
	escaped, err := module.EscapePath(path)
	if err != nil {
		return "", err
	}
	versionDir := filepath.Join(proxy, escaped, "@v")
	if err = os.MkdirAll(versionDir, 0o755); err != nil {
		return "", err
	}
	info, err := json.Marshal(struct {
		Version string
		Time    time.Time
	}{fixtureVersion, time.Unix(0, 0).UTC()})
	if err != nil {
		return "", err
	}
	files := map[string][]byte{
		filepath.Join(versionDir, fixtureVersion+".mod"):  gomod,
		filepath.Join(versionDir, fixtureVersion+".info"): info,
		filepath.Join(versionDir, "list"):                 []byte(fixtureVersion + "\n"),
		filepath.Join(proxy, escaped, "@latest"):          info,
	}
	for name, data := range files {
		if err = os.WriteFile(name, data, 0o600); err != nil {
			return "", err
		}
	}
	if err = writeZip(root, dir, filepath.Join(versionDir, fixtureVersion+".zip"), path); err != nil {
		return "", err
	}
	return path, nil
}

func writeZip(root, dir, archivePath, modulePath string) error {
	cmd := exec.Command("git", "ls-files", "-z", "--", dir)
	cmd.Dir = root
	output, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("cannot list module files: %w", err)
	}
	var files []modzip.File
	for path := range strings.SplitSeq(strings.TrimSuffix(string(output), "\x00"), "\x00") {
		if path != "" {
			files = append(files, diskFile{
				path:     filepath.Join(root, filepath.FromSlash(path)),
				relative: strings.TrimPrefix(path, dir+"/"),
			})
		}
	}
	archive, err := os.Create(archivePath)
	if err != nil {
		return err
	}
	err = modzip.Create(archive, module.Version{Path: modulePath, Version: fixtureVersion}, files)
	return errors.Join(err, archive.Close())
}

type diskFile struct {
	path     string
	relative string
}

func (f diskFile) Path() string { return f.relative }

func (f diskFile) Lstat() (os.FileInfo, error) { return os.Lstat(f.path) }

func (f diskFile) Open() (io.ReadCloser, error) { return os.Open(f.path) }
