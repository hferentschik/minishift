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

package cmd

import (
	"fmt"
	"os"
	"runtime"

	units "github.com/docker/go-units"
	"github.com/docker/machine/libmachine"
	"github.com/docker/machine/libmachine/provision"
	"github.com/docker/machine/libmachine/host"
	"github.com/golang/glog"
	ocfg "github.com/openshift/origin/pkg/cmd/cli/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"k8s.io/kubernetes/pkg/api"
	cfg "k8s.io/kubernetes/pkg/client/unversioned/clientcmd/api"

	"github.com/jimmidyson/minishift/pkg/minikube/cluster"
	"github.com/jimmidyson/minishift/pkg/minikube/constants"
	"github.com/jimmidyson/minishift/pkg/minikube/kubeconfig"
	"github.com/jimmidyson/minishift/pkg/util"
	"github.com/openshift/origin/pkg/bootstrap/docker"
	kcmdutil "k8s.io/kubernetes/pkg/kubectl/cmd/util"
	"github.com/openshift/origin/pkg/cmd/util/clientcmd"
	"github.com/docker/machine/libmachine/drivers"
)

const (
	isoURL                = "iso-url"
	memory                = "memory"
	cpus                  = "cpus"
	humanReadableDiskSize = "disk-size"
	vmDriver              = "vm-driver"
	openshiftVersion      = "openshift-version"
	hostOnlyCIDR          = "host-only-cidr"
	deployRegistry        = "deploy-registry"
	deployRouter          = "deploy-router"
)

var (
	dockerEnv        []string
	insecureRegistry []string
	registryMirror   []string
)

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Starts a local OpenShift cluster.",
	Long: `Starts a local OpenShift cluster using Virtualbox. This command
assumes you already have Virtualbox installed.`,
	Run: runStart,
}

type MinishiftProvisionerDetector struct {
}

func (detector *MinishiftProvisionerDetector) DetectProvisioner(d drivers.Driver) (provision.Provisioner, error) {
	fmt.Println("Hello world")
	return provision.NewBoot2DockerProvisioner(d), nil
}

func runStart(cmd *cobra.Command, args []string) {
	fmt.Println("Starting local OpenShift cluster...")
	provision.SetDetector(&MinishiftProvisionerDetector{})
	libMachineClient := libmachine.NewClient(constants.Minipath, constants.MakeMiniPath("certs"))
	defer libMachineClient.Close()

	config := cluster.MachineConfig{
		MinikubeISO:      viper.GetString(isoURL),
		Memory:           viper.GetInt(memory),
		CPUs:             viper.GetInt(cpus),
		DiskSize:         calculateDiskSizeInMB(viper.GetString(humanReadableDiskSize)),
		VMDriver:         viper.GetString(vmDriver),
		DockerEnv:        dockerEnv,
		InsecureRegistry: insecureRegistry,
		RegistryMirror:   registryMirror,
		HostOnlyCIDR:     viper.GetString(hostOnlyCIDR),
		DeployRouter:     viper.GetBool(deployRouter),
		DeployRegistry:   viper.GetBool(deployRegistry),
		OpenShiftVersion: viper.GetString(openshiftVersion),
	}

	var host *host.Host
	start := func() (err error) {
		host, err = cluster.StartHost(libMachineClient, config)
		if err != nil {
			glog.Errorf("Error starting host: %s. Retrying.\n", err)
		}
		return err
	}
	err := util.Retry(3, start)
	if err != nil {
		glog.Errorln("Error starting host: ", err)
		os.Exit(1)
	}

	envMap, err := cluster.GetHostDockerEnv(libMachineClient)
	for k, v := range envMap {
		os.Setenv(k, v)
	}

	clusterStartConfig := &docker.ClientStartConfig{
		Out:            os.Stdout,
		PortForwarding: defaultPortForwarding(),
	}

	f := clientcmd.New(cmd.PersistentFlags())
	kcmdutil.CheckErr(clusterStartConfig.Complete(f, cmd))
	kcmdutil.CheckErr(clusterStartConfig.Validate(os.Stdout))
	if err := clusterStartConfig.Start(os.Stdout); err != nil {
		os.Exit(1)
	}
}

func calculateDiskSizeInMB(humanReadableDiskSize string) int {
	diskSize, err := units.FromHumanSize(humanReadableDiskSize)
	if err != nil {
		glog.Errorf("Invalid disk size: %s", err)
	}
	return int(diskSize / units.MB)
}

// setupKubeconfig reads config from disk, adds the minikube settings, and writes it back.
// activeContext is true when minikube is the CurrentContext
// If no CurrentContext is set, the given name will be used.
func setupKubeconfig(server, certAuth string) error {
	configFile := constants.KubeconfigPath

	// read existing config or create new if does not exist
	config, err := kubeconfig.ReadConfigOrNew(configFile)
	if err != nil {
		return err
	}

	currentContextName := config.CurrentContext
	currentContext := config.Contexts[currentContextName]

	clusterName, err := ocfg.GetClusterNicknameFromURL(server)
	if err != nil {
		return err
	}
	cluster := cfg.NewCluster()
	cluster.Server = server
	cluster.CertificateAuthorityData = []byte(certAuth)
	config.Clusters[clusterName] = cluster

	// user
	userName := "admin/" + clusterName
	user := cfg.NewAuthInfo()
	if currentContext != nil && currentContext.AuthInfo == userName {
		currentUser := config.AuthInfos[userName]
		if currentUser != nil {
			user.Token = config.AuthInfos[userName].Token
		}
	}
	config.AuthInfos[userName] = user

	// context
	context := cfg.NewContext()
	context.Cluster = clusterName
	context.AuthInfo = userName
	context.Namespace = api.NamespaceDefault
	contextName := ocfg.GetContextNickname(api.NamespaceDefault, clusterName, userName)
	if currentContext != nil && currentContext.Cluster == clusterName && currentContext.AuthInfo == userName {
		contextName = currentContextName
		context.Namespace = currentContext.Namespace
	}
	config.Contexts[contextName] = context

	config.CurrentContext = contextName

	// write back to disk
	if err := kubeconfig.WriteConfig(config, configFile); err != nil {
		return err
	}

	fmt.Println("oc is now configured to use the cluster.")
	if len(user.Token) == 0 {
		fmt.Println("Run this command to use the cluster: ")
		fmt.Println("oc login --username=admin --password=admin")
	}

	return nil
}

func init() {
	startCmd.Flags().String(isoURL, constants.DefaultIsoUrl, "Location of the minishift iso")
	startCmd.Flags().String(vmDriver, constants.DefaultVMDriver, fmt.Sprintf("VM driver is one of: %v", constants.SupportedVMDrivers))
	startCmd.Flags().Int(memory, constants.DefaultMemory, "Amount of RAM allocated to the minishift VM")
	startCmd.Flags().Int(cpus, constants.DefaultCPUS, "Number of CPUs allocated to the minishift VM")
	startCmd.Flags().String(humanReadableDiskSize, constants.DefaultDiskSize, "Disk size allocated to the minishift VM (format: <number>[<unit>], where unit = b, k, m or g)")
	startCmd.Flags().String(hostOnlyCIDR, "192.168.99.1/24", "The CIDR to be used for the minishift VM (only supported with Virtualbox driver)")
	startCmd.Flags().StringSliceVar(&dockerEnv, "docker-env", nil, "Environment variables to pass to the Docker daemon. (format: key=value)")
	startCmd.Flags().StringSliceVar(&insecureRegistry, "insecure-registry", []string{"172.30.0.0/16"}, "Insecure Docker registries to pass to the Docker daemon")
	startCmd.Flags().StringSliceVar(&registryMirror, "registry-mirror", nil, "Registry mirrors to pass to the Docker daemon")
	startCmd.Flags().Bool(deployRegistry, true, "Should the OpenShift internal Docker registry be deployed?")
	startCmd.Flags().Bool(deployRouter, false, "Should the OpenShift router be deployed?")
	startCmd.Flags().String(openshiftVersion, "", "The OpenShift version that the minishift VM will run (ex: v1.2.3) OR a URI which contains an openshift binary (ex: file:///home/developer/go/src/github.com/openshift/origin/_output/local/bin/linux/amd64/openshift)")

	viper.BindPFlags(startCmd.Flags())

	RootCmd.AddCommand(startCmd)
}

func defaultPortForwarding() bool {
	// Defaults to true if running on Mac, with no DOCKER_HOST defined
	return runtime.GOOS == "darwin" && len(os.Getenv("DOCKER_HOST")) == 0
}
