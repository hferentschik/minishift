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
	"github.com/minishift/minishift/pkg/minishift/hostfolder/sshfs"

	"github.com/spf13/cobra"
)

// TODO - The whole setup would beed to live under a generic, overaching mechanism to configure host folders (see Windows shares approach) (HF)
var sshfsCmd = &cobra.Command{
	Use:   "sshfs",
	Short: "Snafu",
	Long:  "Snafu",
	Run:   runSshfs,
}

func init() {
	RootCmd.AddCommand(sshfsCmd)
}

func runSshfs(cmd *cobra.Command, args []string) {
	sshfs.CreateSftpDaemon()
}
