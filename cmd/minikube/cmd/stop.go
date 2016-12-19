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

	"github.com/docker/machine/libmachine"
	"github.com/minishift/minishift/pkg/minikube/cluster"
	"github.com/minishift/minishift/pkg/minikube/constants"
	"github.com/minishift/minishift/pkg/minishift/registration"
	"github.com/spf13/cobra"
)

// stopCmd represents the stop command
var stopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stops a running local OpenShift cluster.",
	Long: `Stops a local OpenShift cluster running in Virtualbox. This command stops the VM
itself, leaving all files intact. The cluster can be started again with the "start" command.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Stopping local OpenShift cluster...")
		api := libmachine.NewClient(constants.Minipath, constants.MakeMiniPath("certs"))
		defer api.Close()
		// Unregister VM when stop
		host, err := api.Load(constants.MachineName)
		if err != nil {
			fmt.Println("Error stopping machine: ", err)
			os.Exit(1)
		}
		registrator, err := registration.DetectRegistrator(host.Driver)
		if err != nil {
			fmt.Println("Distribution doesn't support unregister")
		} else {
			fmt.Println("Unregistering distribution...")
			if err := registrator.Unregister(); err != nil {
				fmt.Println("Error Unregistering distribution ", err)
				os.Exit(1)
			}
		}

		if err := cluster.StopHost(api); err != nil {
			fmt.Println("Error stopping machine: ", err)
			os.Exit(1)
		}
		fmt.Println("Machine stopped.")
	},
}

func init() {
	RootCmd.AddCommand(stopCmd)
}
