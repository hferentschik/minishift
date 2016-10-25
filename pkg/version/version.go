/*
Copyright (C) 2016 Red Hat, Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package version

import (
	"strings"

	"github.com/blang/semver"
)

const VersionPrefix = "v"

// The current version of minishift
// This is a private field and should be set when compiling with --ldflags="-X github.com/jimmidyson/minishift/pkg/version.version=vX.Y.Z"
var version = "0.9.0"

func GetVersion() string {
	return version
}

func GetSemverVersion() (semver.Version, error) {
	return semver.Make(strings.TrimPrefix(GetVersion(), VersionPrefix))
}

// The default version of OpenShift
// This is a private field and should be set when compiling with --ldflags="-X github.com/jimmidyson/minishift/pkg/version.openshiftVersion=vX.Y.Z"
var openshiftVersion = "v0.0.0-unset"

func GetOpenShiftVersion() string {
	return openshiftVersion
}

func GetOpenShiftSemverVersion() (semver.Version, error) {
	return semver.Make(strings.TrimPrefix(GetOpenShiftVersion(), VersionPrefix))
}
