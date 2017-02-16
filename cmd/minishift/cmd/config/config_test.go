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
	"bytes"
	"reflect"
	"testing"
)

type configTestCase struct {
	data   string
	config map[string]interface{}
}

var configTestCases = []configTestCase{
	{
		data: `{
    "memory": 2
}`,
		config: map[string]interface{}{
			"memory": 2,
		},
	},
	{
		data: `{
    "ReminderWaitPeriodInHours": 99,
    "cpus": 4,
    "disk-size": "20g",
    "iso-url": "http://foo.bar/minishift-centos.iso",
    "log_dir": "/etc/hosts",
    "show-libmachine-logs": true,
    "v": 5,
    "vm-driver": "kvm"
}`,
		config: map[string]interface{}{
			"iso-url":   "http://foo.bar/minishift-centos.iso",
			"vm-driver": "kvm",
			"cpus":      4,
			"disk-size": "20g",
			"v":         5,
			"show-libmachine-logs":      true,
			"log_dir":                   "/etc/hosts",
			"ReminderWaitPeriodInHours": 99,
		},
	},
}

func TestReadConfig(t *testing.T) {
	for _, tt := range configTestCases {
		r := bytes.NewBufferString(tt.data)
		config, err := decode(r)
		if reflect.DeepEqual(config, tt.config) || err != nil {
			t.Errorf("Cannot decode config. \n\n expected %+v, \n\n received %+v", tt.config, config)
		}
	}
}

func TestWriteConfig(t *testing.T) {
	var b bytes.Buffer
	for _, tt := range configTestCases {
		err := encode(&b, tt.config)
		if err != nil {
			t.Errorf("Error encoding: %s", err)
		}
		if b.String() != tt.data {
			t.Errorf("Cannot encode config. \n\n expected \n %+v, \n\n received \n %+v", tt.data, b.String())
		}
		b.Reset()
	}
}
