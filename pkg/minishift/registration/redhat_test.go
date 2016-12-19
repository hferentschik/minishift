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

package registration

import (
	"github.com/docker/machine/libmachine/drivers"
	"github.com/minishift/minishift/pkg/minikube/tests"
	"testing"
)

func TestRedhatRegistratorCompatibleWithDistribution(t *testing.T) {
	s, _ := tests.NewSSHServer()
	s.CommandToOutput = make(map[string]string)
	s.CommandToOutput["sudo cat /etc/os-release"] = `NAME="Red Hat Enterprise Linux Server"
VERSION="7.3 (Maipo)"
ID="rhel"
ID_LIKE="fedora"
VERSION_ID="7.3"
PRETTY_NAME="Red Hat Enterprise Linux Server 7.3 (Maipo)"
ANSI_COLOR="0;31"
CPE_NAME="cpe:/o:redhat:enterprise_linux:7.3:GA:server"
HOME_URL="https://www.redhat.com/"
BUG_REPORT_URL="https://bugzilla.redhat.com/"

REDHAT_BUGZILLA_PRODUCT="Red Hat Enterprise Linux 7"
REDHAT_BUGZILLA_PRODUCT_VERSION=7.3
REDHAT_SUPPORT_PRODUCT="Red Hat Enterprise Linux"
REDHAT_SUPPORT_PRODUCT_VERSION="7.3"
VARIANT="minishift"
VARIANT_VERSION="1.0.0-alpha.1"
`
	s.CommandToOutput["sudo subscription-manager version"] = `server type: This system is currently not registered.
subscription management server: 0.9.51.11-1
subscription management rules: 5.15
subscription-manager: 1.17.15-1.el7
python-rhsm: 1.17.9-1.el7
`
	port, err := s.Start()
	if err != nil {
		t.Fatalf("Error starting ssh server: %s", err)
	}
	d := &tests.MockDriver{
		Port: port,
		BaseDriver: drivers.BaseDriver{
			IPAddress:  "127.0.0.1",
			SSHKeyPath: "",
		},
	}

	registrator := NewRedhatRegistrator(d)
	if !registrator.CompatibleWithDistribution() {
		t.Fatal("Registration capability should be in the Distribution")
	}
}

func TestRedhatRegistratorRegister(t *testing.T) {
	s, _ := tests.NewSSHServer()
	s.CommandToOutput = make(map[string]string)
	s.CommandToOutput["sudo subscription-manager version"] = `server type: This system is currently not registered.
subscription management server: 0.9.51.11-1
subscription management rules: 5.15
subscription-manager: 1.17.15-1.el7
python-rhsm: 1.17.9-1.el7
`
	s.CommandToOutput["sudo subscription-manager register --auto-attach --username foo --password foo"] = `
	The system has been registered with ID: 2af78473-n5th-ksa1-3748-983748c77da`
	port, err := s.Start()
	if err != nil {
		t.Fatalf("Error starting ssh server: %s", err)
	}
	d := &tests.MockDriver{
		Port: port,
		BaseDriver: drivers.BaseDriver{
			IPAddress:  "127.0.0.1",
			SSHKeyPath: "",
		},
	}

	registrator := NewRedhatRegistrator(d)
	m := make(map[string]string)
	m["username"] = "foo"
	m["password"] = "foo"
	if err := registrator.Register(m); err != nil {
		t.Fatal("Distribution should able to register")
	}
}

func TestRedhatRegistratorUnregister(t *testing.T) {
	s, _ := tests.NewSSHServer()
	s.CommandToOutput = make(map[string]string)
	s.CommandToOutput["sudo subscription-manager version"] = `server type: This system is currently not registered.
subscription management server: 0.9.51.11-1
subscription management rules: 5.15
subscription-manager: 1.17.15-1.el7
python-rhsm: 1.17.9-1.el7
`
	port, err := s.Start()
	if err != nil {
		t.Fatalf("Error starting ssh server: %s", err)
	}
	d := &tests.MockDriver{
		Port: port,
		BaseDriver: drivers.BaseDriver{
			IPAddress:  "127.0.0.1",
			SSHKeyPath: "",
		},
	}

	registrator := NewRedhatRegistrator(d)
	if err := registrator.Unregister(); err != nil {
		t.Fatal("Distribution shouldn't able to unregister")
	}
}
