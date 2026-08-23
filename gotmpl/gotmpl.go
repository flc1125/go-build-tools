// Copyright The OpenTelemetry Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"text/template"
)

// gotmpl uses [text/template] package to generate a file using
// bodyPath as template body filepath,
// jsonData as data in JSON format,
// outPath as output filepath.
func gotmpl(bodyPath, jsonData, outPath string) error {
	if bodyPath == "" {
		return errors.New("gotmpl: template body filepath must be set")
	}
	if outPath == "" {
		return errors.New("gotmpl: output filepath must be set")
	}

	tmpl, err := template.ParseFiles(bodyPath)
	if err != nil {
		return fmt.Errorf("gotmpl: cannot parse template body file: %w", err)
	}

	var data any
	if err = json.Unmarshal(([]byte)(jsonData), &data); err != nil {
		return fmt.Errorf("gotmpl: data must be in JSON format: %w", err)
	}

	var rendered bytes.Buffer
	err = tmpl.Option("missingkey=error").Execute(&rendered, data)
	if err != nil {
		return fmt.Errorf("gotmpl: execution failed: %w", err)
	}

	outDir := filepath.Dir(outPath)
	outFile, err := os.CreateTemp(outDir, "."+filepath.Base(outPath)+".*")
	if err != nil {
		return fmt.Errorf("gotmpl: cannot create temporary output file: %w", err)
	}
	tmpPath := outFile.Name()
	defer func() {
		_ = outFile.Close()
		_ = os.Remove(tmpPath)
	}()

	perm := os.FileMode(0o644)
	if info, statErr := os.Stat(outPath); statErr == nil {
		perm = info.Mode().Perm()
	} else if !errors.Is(statErr, os.ErrNotExist) {
		return fmt.Errorf("gotmpl: cannot inspect output file: %w", statErr)
	}
	if err := outFile.Chmod(perm); err != nil {
		return fmt.Errorf("gotmpl: cannot set temporary output permissions: %w", err)
	}
	if _, err := outFile.Write(rendered.Bytes()); err != nil {
		return fmt.Errorf("gotmpl: cannot write temporary output file: %w", err)
	}
	if err := outFile.Sync(); err != nil {
		return fmt.Errorf("gotmpl: cannot sync temporary output file: %w", err)
	}
	if err := outFile.Close(); err != nil {
		return fmt.Errorf("gotmpl: cannot close temporary output file: %w", err)
	}
	if err := os.Rename(tmpPath, outPath); err != nil {
		return fmt.Errorf("gotmpl: cannot replace output file: %w", err)
	}
	return nil
}
