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
	"github.com/stretchr/testify/mock"
	"testing"
)

type MockDockerCommander struct{
	mock.Mock
}

func (m *MockDockerCommander) Ps() (string, error) {
	return "", nil
}

func (m *MockDockerCommander) Start(container string) (bool, error) {
	return true, nil
}

func (m *MockDockerCommander) Stop(container string) (bool, error) {
	return true, nil
}

func (m *MockDockerCommander) Cp(source string, container string, target string) error {
	return nil
}

func (m *MockDockerCommander) Exec(options string, container string, command string, args string) (string, error) {
	return "", nil
}

func (m *MockDockerCommander) LocalExec(cmd string) (string, error) {
	return "", nil
}

func (m *MockDockerCommander) Status(container string) (string, error) {
	return "", nil
}

func (m *MockDockerCommander) Restart(container string) (bool, error) {
	args := m.Called(container)
	return args.Bool(0), args.Error(1)
}

func TestMinishiftLogLevel(t *testing.T) {
	// create an instance of our test object
	commander := new(MockDockerCommander)

	// setup expectations
	commander.On("Restart", "origin").Return(true, nil)

	RestartOpenShift(commander)

	// assert that the expectations were met
	commander.AssertExpectations(t)
}

