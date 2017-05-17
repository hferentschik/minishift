/*
Copyright (C) 2017 Red Hat, Inc.

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

package clusterup

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/blang/semver"
	"github.com/docker/machine/libmachine/host"
	"github.com/docker/machine/libmachine/provision"
	"github.com/minishift/minishift/pkg/minikube/constants"
	"github.com/minishift/minishift/pkg/minikube/kubeconfig"
	"github.com/minishift/minishift/pkg/minishift/addon/command"
	"github.com/minishift/minishift/pkg/minishift/addon/manager"
	"github.com/minishift/minishift/pkg/minishift/oc"
	"strings"
	"time"
)

const (
	ip_key             = "ip"
	routing_suffix_key = "routing-suffix"
)

type ClusterUpConfig struct {
	OpenShiftVersion string
	MachineName      string
	Ip               string
	Port             int
	RoutingSuffix    string
	HostPvDir        string
	User             string
	Project          string
}

// postClusterUp runs the Minishift specific provisioning after cluster up has run
func PostClusterUp(config *ClusterUpConfig, ocRunner *oc.OcRunner, sshCommander provision.SSHCommander, addOnManager *manager.AddOnManager) error {
	err := kubeconfig.CacheSystemAdminEntries(ocRunner.KubeConfigPath, getConfigClusterName(config.Ip, config.Port))
	if err != nil {
		return err
	}

	err = ocRunner.AddSudoerRoleForUser(config.User)
	if err != nil {
		return err
	}

	err = ocRunner.AddCliContext(config.MachineName, config.Ip, config.User, config.Project)
	if err != nil {
		return err
	}

	err = configurePersistentVolumes(config.HostPvDir, addOnManager, sshCommander, ocRunner)
	if err != nil {
		return err
	}

	err = applyAddOns(addOnManager, config.Ip, config.RoutingSuffix, ocRunner, sshCommander)
	if err != nil {
		return err
	}

	return nil
}

//Create the host directories in the virtual machine if not present
func EnsureHostDirectoriesExist(host *host.Host, dirs []string) error {
	cmd := fmt.Sprintf("sudo mkdir -p %s", strings.Join(dirs, " "))
	_, err := host.RunSSHCommand(cmd)
	if err != nil {
		return err
	}
	return nil
}

// TODO - persistent volume creation should really be fixed upstream, aka 'cluster up'. See https://github.com/openshift/origin/issues/14076 (HF)
// configurePersistentVolumes makes sure that the default persistent volumes created by 'cluster up' have the right permissions - see https://github.com/minishift/minishift/issues/856
func configurePersistentVolumes(hostPvDir string, addOnManager *manager.AddOnManager, sshCommander provision.SSHCommander, ocRunner *oc.OcRunner) error {
	// don't apply this if anyuid is not enabled
	anyuid := addOnManager.Get("anyuid")
	if anyuid == nil || !anyuid.IsEnabled() {
		return nil
	}

	fmt.Print("-- Waiting for persistent volumes to be created ... ")

	var out, err *bytes.Buffer

	// poll the status of the persistent-volume-setup job to determine when the persitent volume creates is completed
	for {
		out = new(bytes.Buffer)
		err = new(bytes.Buffer)
		exitStatus := ocRunner.Run("get job persistent-volume-setup -n default -o 'jsonpath={ .status.active }'", out, err)

		if exitStatus != 0 || len(err.String()) > 0 {
			return errors.New("Unable to monitor persistent volume creation")
		}

		if out.String() != "1" {
			break
		}

		time.Sleep(1 * time.Second)
	}

	// verify the job succeeded
	out = new(bytes.Buffer)
	err = new(bytes.Buffer)
	exitStatus := ocRunner.Run("get job persistent-volume-setup -n default -o 'jsonpath={ .status.succeeded }'", out, err)

	if exitStatus != 0 || len(err.String()) > 0 || out.String() != "1" {
		return errors.New("Persistent volume creation failed")
	}

	cmd := fmt.Sprintf("sudo chmod -R 777 %s/pv*", hostPvDir)
	sshCommander.SSHCommand(cmd)

	// if we have SELinux enabled we need to sort things out there as well
	// 'cluster up' does this as well, but we do it here as well to have all required actions collected in one
	// place, instead of relying on some implicit knowledge on what 'cluster up does (HF)
	cmd = fmt.Sprintf("sudo which chcon; if [ $? -eq 0 ]; then chcon -R -t svirt_sandbox_file_t %s/*; fi", hostPvDir)
	sshCommander.SSHCommand(cmd)

	cmd = fmt.Sprintf("sudo which restorecon; if [ $? -eq 0 ]; then restorecon -R %s; fi", hostPvDir)
	sshCommander.SSHCommand(cmd)

	fmt.Println("OK")
	fmt.Println()

	return nil
}

func applyAddOns(addOnManager *manager.AddOnManager, ip string, routingSuffix string, ocRunner *oc.OcRunner, sshCommander provision.SSHCommander) error {
	context, err := getExecutionContext(ip, routingSuffix, ocRunner, sshCommander)
	err = addOnManager.Apply(context)
	if err != nil {
		return err
	}

	return nil
}

func getConfigClusterName(ip string, port int) string {
	return fmt.Sprintf("%s:%d", strings.Replace(ip, ".", "-", -1), port)
}

func getExecutionContext(ip string, routingSuffix string, ocRunner *oc.OcRunner, sshCommander provision.SSHCommander) (*command.ExecutionContext, error) {
	context, err := command.NewExecutionContext(ocRunner, sshCommander)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Unable to initialise execution context: %s", err.Error()))
	}

	context.AddToContext(ip_key, ip)
	context.AddToContext(routing_suffix_key, routingSuffix)

	return context, nil
}

func ValidateOpenshiftMinVersion(ver string, minVersion string) bool {
	v, _ := semver.Parse(strings.TrimPrefix(ver, constants.VersionPrefix))
	minSupportedVersion := strings.TrimPrefix(minVersion, constants.VersionPrefix)
	versionRange, err := semver.ParseRange(fmt.Sprintf(">=%s", minSupportedVersion))
	if err != nil {
		fmt.Println("Not able to parse version info", err)
		return false
	}
	if versionRange(v) {
		return true
	}
	return false
}
