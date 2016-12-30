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
	"reflect"
	"runtime"
	"path/filepath"
	"github.com/minishift/minishift/pkg/minikube/constants"
	"github.com/docker/machine/drivers/hyperv"
)

var machineConfig = MachineConfig{
	MinikubeISO: "https://github.com/fake/boot2docker.iso",
	Memory: 2048,
	CPUs: 2,
	DiskSize: 10000,
	VMDriver: "hyperv",
}

func TestCreateHypervHost(t *testing.T) {
	if runtime.GOOS != "windows" {
		t.Skip("Skip TestCreateHypervHost for non-windows test.")
	}
	isoPath := filepath.Join(constants.Minipath, "cache", "iso", "boot2docker.iso")
	expectedURL := "file://" + filepath.ToSlash(isoPath)

	d := createHypervHost(machineConfig)
	expectedDriver := "*hyperv.Driver"

	if reflect.TypeOf(d).String() != expectedDriver {
		t.Fatalf("Unexpected driver type. Expected '%s' but got '%s'", expectedDriver, reflect.TypeOf(d).String())
	}

	driver := d.(*hyperv.Driver)

	assert.Equal(t, driver.Boot2DockerURL, expectedURL, "they should be equal")
	assert.Equal(t, driver.DiskSize, machineConfig.DiskSize, "they should be equal")
	assert.Equal(t, driver.MemSize, machineConfig.Memory, "they should be equal")
	assert.Equal(t, driver.CPU, machineConfig.CPUs, "they should be equal")
	assert.Equal(t, driver.SSHUser, "docker", "they should be equal")
}
