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

package openshift

import (
	"github.com/minishift/minishift/pkg/minishift/util"
	"testing"
)

func Test_config_commands_needs_existing_vm(t *testing.T) {
	setup(t)
	defer tearDown()

	util.RegisterExitHandler(createExitHandlerFunc(t, 1, nonExistentMachineError))

	target = "master"
	runConfig(nil, nil)
}

func Test_unknown_config_target_aborts_command(t *testing.T) {
	setup(t)
	defer tearDown()

	util.RegisterExitHandler(createExitHandlerFunc(t, 1, unknownPatchTargetError))

	target = "foo"
	runConfig(nil, nil)
}
