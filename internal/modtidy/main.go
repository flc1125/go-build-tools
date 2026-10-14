// Copyright 2026 flc1125
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"golang.org/x/mod/modfile"

	"github.com/flc1125/go-build-tools/internal/repo"
)

func main() {
	goBinary := flag.String("go", "go", "Go executable for dependency tidying.")
	flag.Parse()
	if err := run(*goBinary, flag.Args()); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(goBinary string, dirs []string) error {
	if len(dirs) == 0 {
		return errors.New("specify at least one module directory")
	}
	root, err := repo.FindRoot()
	if err != nil {
		return err
	}
	data, err := os.ReadFile(filepath.Join(root, "go.mod"))
	if err != nil {
		return err
	}
	rootPath := modfile.ModulePath(data)
	if rootPath == "" {
		return errors.New("root go.mod does not declare a module path")
	}
	for _, dir := range dirs {
		fmt.Printf("%s mod tidy in %s\n", goBinary, dir)
		if err = tidyModule(goBinary, rootPath, root, dir); err != nil {
			return fmt.Errorf("cannot tidy %s: %w", dir, err)
		}
	}
	return nil
}

func tidyModule(goBinary, rootPath, root, dir string) error {
	gomodPath := filepath.Join(dir, "go.mod")
	data, err := os.ReadFile(gomodPath)
	if err != nil {
		return err
	}
	f, err := modfile.Parse(gomodPath, data, nil)
	if err != nil {
		return err
	}
	temporaryReplace := needsLocalRoot(f, rootPath)
	if temporaryReplace {
		if err = f.AddReplace(rootPath, "", root, ""); err != nil {
			return err
		}
	}
	tmp, err := os.MkdirTemp("", "go-build-tools-tidy-")
	if err != nil {
		return err
	}
	defer func() { _ = os.RemoveAll(tmp) }()
	altMod := filepath.Join(tmp, "go.mod")
	data, err = f.Format()
	if err != nil {
		return err
	}
	if err = os.WriteFile(altMod, data, 0o600); err != nil {
		return err
	}
	sumPath := filepath.Join(dir, "go.sum")
	sum, err := os.ReadFile(sumPath)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return err
	}
	hadSum := err == nil
	altSum := filepath.Join(tmp, "go.sum")
	if hadSum {
		if err = os.WriteFile(altSum, sum, 0o600); err != nil {
			return err
		}
	}
	cmd := exec.Command(goBinary, "mod", "tidy", "-modfile", altMod)
	cmd.Dir = dir
	cmd.Env = append(os.Environ(), "GOWORK=off")
	if output, tidyErr := cmd.CombinedOutput(); tidyErr != nil {
		return fmt.Errorf("go mod tidy failed: %w\n%s", tidyErr, output)
	}
	data, err = os.ReadFile(altMod)
	if err != nil {
		return err
	}
	f, err = modfile.Parse(altMod, data, nil)
	if err != nil {
		return err
	}
	if temporaryReplace {
		if err = f.DropReplace(rootPath, ""); err != nil {
			return err
		}
		f.Cleanup()
	}
	data, err = f.Format()
	if err != nil {
		return err
	}
	sum, err = os.ReadFile(altSum)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return err
	}
	writeSum := hadSum || err == nil
	if err = os.WriteFile(gomodPath, data, 0o600); err != nil {
		return err
	}
	if writeSum {
		return os.WriteFile(sumPath, sum, 0o600)
	}
	return nil
}

func needsLocalRoot(f *modfile.File, rootPath string) bool {
	for _, replacement := range f.Replace {
		if replacement.Old.Path == rootPath {
			return false
		}
	}
	for _, dependency := range f.Require {
		if dependency.Mod.Path == rootPath {
			return true
		}
	}
	return false
}
