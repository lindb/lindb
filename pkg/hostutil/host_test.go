// Licensed to LinDB under one or more contributor
// license agreements. See the NOTICE file distributed with
// this work for additional information regarding copyright
// ownership. LinDB licenses this file to you under
// the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

package hostutil

import (
	"fmt"
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetHostIP(t *testing.T) {
	ip, err := GetHostIP()
	assert.NoError(t, err)
	assert.NotEmpty(t, ip)
}

func Test_getHostInfo(t *testing.T) {
	defer func() {
		netInterfaces = net.Interfaces
	}()

	// mock err
	netInterfaces = func() (interfaces []net.Interface, e error) {
		return nil, fmt.Errorf("err")
	}
	host := getHostInfo()
	assert.Empty(t, host.hostIP)
	assert.Error(t, host.err)

	// mock empty
	netInterfaces = func() (interfaces []net.Interface, e error) {
		return nil, nil
	}
	host = getHostInfo()
	assert.Empty(t, host.hostIP)
	assert.Error(t, host.err)

	netInterfaces = func() (interfaces []net.Interface, e error) {
		return []net.Interface{
			{
				Name:  "mock_test",
				Flags: 1,
			},
		}, nil
	}
	_ = getHostInfo()
}
