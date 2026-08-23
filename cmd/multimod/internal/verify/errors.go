// Copyright The OpenTelemetry Authors
// Modified by flc1125 for github.com/flc1125/go-build-tools.
// SPDX-License-Identifier: Apache-2.0

package verify

import (
	"fmt"
	"strings"

	"github.com/flc1125/go-build-tools/cmd/multimod/internal/shared"
)

type errModuleNotInSet struct {
	modPath     shared.ModulePath
	modFilePath shared.ModuleFilePath
}

func (e *errModuleNotInSet) Error() string {
	return fmt.Sprintf("Module %v (defined in %v) is not listed in any module set.", e.modPath, e.modFilePath)
}

type errModuleNotInRepo struct {
	modPath    shared.ModulePath
	modSetName string
}

func (e *errModuleNotInRepo) Error() string {
	return fmt.Sprintf("Module %v in module set %v does not exist in the current repo.", e.modPath, e.modSetName)
}

type errInvalidVersion struct {
	modSetName    string
	modSetVersion string
}

func (e *errInvalidVersion) Error() string {
	return fmt.Sprintf("Module set %v has invalid version string: %v", e.modSetName, e.modSetVersion)
}

type errMultipleSetSameVersionSlice struct {
	errs []*errMultipleSetSameVersion
}

func (e *errMultipleSetSameVersionSlice) Error() string {
	errorStringSlice := make([]string, 0, len(e.errs))
	for _, err := range e.errs {
		errorStringSlice = append(errorStringSlice, err.Error())
	}

	return strings.Join(errorStringSlice, "\n")
}

type errMultipleSetSameVersion struct {
	modSetNames   []string
	modSetVersion string
}

func (e *errMultipleSetSameVersion) Error() string {
	return fmt.Sprintf("Multiple module sets have the same major version (%v): %v",
		e.modSetVersion, e.modSetNames)
}

// errDependency is logged upon discovery that a stable module depends on an unstable module.
type errDependency struct {
	modPath    shared.ModulePath
	modVersion string
	depPath    shared.ModulePath
	depVersion string
}

func (e *errDependency) Error() string {
	return fmt.Sprintf("WARNING: Stable module %v (%v) depends on unstable module %v (%v).\n",
		e.modPath, e.modVersion,
		e.depPath, e.depVersion)
}
