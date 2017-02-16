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

package config

import (
	"fmt"
	"net"
	"net/url"
	"os"
	"strconv"

	units "github.com/docker/go-units"

	"github.com/minishift/minishift/pkg/minikube/constants"
)

func IsValidDriver(string, driver string) error {
	for _, d := range constants.SupportedVMDrivers {
		if driver == d {
			return nil
		}
	}
	return fmt.Errorf("Driver %s is not supported", driver)
}

func RequiresRestartMsg(string, string) error {
	fmt.Fprintln(os.Stdout, "To apply the changes, you must delete the current VM with `minishift delete` and start a new VM with `minishift start`.")
	return nil
}

func IsValidDiskSize(name string, disksize string) error {
	_, err := units.FromHumanSize(disksize)
	if err != nil {
		return fmt.Errorf("Disk size is not valid: %v", err)
	}
	return nil
}

func IsPositive(name string, val string) error {
	i, err := strconv.Atoi(val)
	if err != nil {
		return fmt.Errorf("%s:%v", name, err)
	}
	if i <= 0 {
		return fmt.Errorf("%s must be > 0", name)
	}
	return nil
}

func IsValidCIDR(name string, cidr string) error {
	_, _, err := net.ParseCIDR(cidr)
	if err != nil {
		return fmt.Errorf("Error parsing CIDR: %v", err)
	}
	return nil
}

func IsValidPath(name string, path string) error {
	_, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("%s path is not valid: %v", name, err)
	}
	return nil
}

func IsValidUrl(_ string, isoURL string) error {
	_, err := url.ParseRequestURI(isoURL)
	if err != nil {
		return fmt.Errorf("%s url is not valid: %v", isoURL, err)
	}
	return nil
}
