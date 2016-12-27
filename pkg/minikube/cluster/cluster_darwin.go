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
	//"os"
	"github.com/docker/machine/drivers/vmwarefusion"
	"github.com/docker/machine/libmachine/drivers"
	//"github.com/docker/machine/libmachine/drivers/plugin"
	//"github.com/docker/machine/libmachine/drivers/plugin/localbinary"
	"github.com/minishift/minishift/pkg/minikube/constants"
	//"github.com/zchee/docker-machine-driver-xhyve/xhyve"
)

func createVMwareFusionHost(config MachineConfig) drivers.Driver {
	d := vmwarefusion.NewDriver(constants.MachineName, constants.Minipath).(*vmwarefusion.Driver)
	d.Boot2DockerURL = config.GetISOFileURI()
	d.Memory = config.Memory
	d.CPU = config.CPUs

	// TODO(philips): push these defaults upstream to fixup this driver
	d.SSHPort = 22
	d.ISO = d.ResolveStorePath("boot2docker.iso")
	return d
}

type xhyveDriver struct {
	*drivers.BaseDriver
	Boot2DockerURL string
	BootCmd        string
	CPU            int
	CaCertPath     string
	DiskSize       int64
	MacAddr        string
	Memory         int
	PrivateKeyPath string
	UUID           string
	NFSShare       bool
	DiskNumber     int
	Virtio9p       bool
	Virtio9pFolder string
}

func createXhyveHost(config MachineConfig) *xhyveDriver {
	// Warning, this is a hack. CoreDrivers is a fixed sized array, so we cannot just add an element.
	// Good thing is we know we want the xhyve driver which we bundle, so we can just set the first element
//	localbinary.CoreDrivers[0] = "xhyve"

//	os.Setenv(localbinary.PluginEnvKey, localbinary.PluginEnvVal)
	// Register the xhyve driver
	//plugin.RegisterDriver(xhyve.NewDriver("", ""))

	return &xhyveDriver{
		BaseDriver: &drivers.BaseDriver{
			MachineName: constants.MachineName,
			StorePath:   constants.Minipath,
		},
		Memory:         config.Memory,
		CPU:            config.CPUs,
		Boot2DockerURL: config.GetISOFileURI(),
		DiskSize:       int64(config.DiskSize),
		Virtio9p:       true,
		Virtio9pFolder: "/Users",
		UUID:           "F4BB3F79-AB4E-4708-95CA-E32FBFCDEFD3",
	}
}
