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

package cmd

import (
	"fmt"

	"github.com/containers/image/copy"
	"github.com/containers/image/signature"
	"github.com/containers/image/transports/alltransports"
	"github.com/containers/image/types"
	"github.com/minishift/minishift/pkg/util/os/atexit"
	"github.com/spf13/cobra"
	"os"
)

var imageCmd = &cobra.Command{
	Use:   "image",
	Short: "foo",
	Long:  `bar`,
	Run:   copyImage,
}

func init() {
	RootCmd.AddCommand(imageCmd)
}

func copyImage(cmd *cobra.Command, args []string) {
	fmt.Println("Trying to copy an image")

	policy := &signature.Policy{Default: []signature.PolicyRequirement{signature.NewPRInsecureAcceptAnything()}}

	policyContext, err := signature.NewPolicyContext(policy)
	if err != nil {
		atexit.ExitWithMessage(1, "Error creating security context")
	}

	srcRef, err := alltransports.ParseImageName("docker://busybox:latest")
	if err != nil {
		atexit.ExitWithMessage(1, fmt.Sprintf("Invalid source name %s: %v", "docker://busybox:latest", err))
	}
	destRef, err := alltransports.ParseImageName("ostree:/Users/hardy/tmp-2/skopeo")
	if err != nil {
		atexit.ExitWithMessage(1, fmt.Sprintf("Invalid destination name %s: %v", "dir:/Users/hardy/tmp-2/skopeo", err))
	}

	copy.Image(policyContext, destRef, srcRef, &copy.Options{
		RemoveSignatures: false,
		SignBy:           "",
		ReportWriter:     os.Stdout,
		SourceCtx:        &types.SystemContext{},
		DestinationCtx:   &types.SystemContext{},
	})

}
