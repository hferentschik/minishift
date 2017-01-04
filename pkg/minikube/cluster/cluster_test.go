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

package cluster

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/minishift/minishift/pkg/minikube/constants"
	"path/filepath"
)

func TestRemoteBoot2DockerURL(t *testing.T) {
	var machineConfig = MachineConfig{
		MinikubeISO: "http://github.com/fake/boot2docker.iso",
	}

	isoPath:= filepath.Join(constants.Minipath, "cache", "iso", filepath.Base(machineConfig.MinikubeISO))
	expectedURL := "file://" + filepath.ToSlash(isoPath)
	url := machineConfig.GetISOFileURI()
	assert.Equal(t, url, expectedURL, "they both should be equal")
}

func TestLocalBoot2DockerURL(t *testing.T) {
	isoPath := filepath.Join(constants.Minipath, "cache", "iso", "boot2docker.iso")
	localISOUrl := "file://" + filepath.ToSlash(isoPath)

	machineConfig := MachineConfig{
		MinikubeISO: localISOUrl,
	}

	url := machineConfig.GetISOFileURI()
	assert.Equal(t, url, localISOUrl, "they both should be equal")
}
