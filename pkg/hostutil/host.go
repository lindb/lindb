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
	"errors"
	"net"
	"sync"
)

var (
	once sync.Once
	host hostInfo
)

// hostInfo defines host basic info, if cannot get host info returns error
type hostInfo struct {
	err    error
	hostIP string
}

// just for testing
var netInterfaces = net.Interfaces

// extractHostInfo extracts host info, just do it once
func extractHostInfo() {
	once.Do(func() {
		host = getHostInfo()
	})
}

// getHostInfo returns host info like ip
func getHostInfo() (host hostInfo) {
	ifaces, err := netInterfaces()
	if err != nil {
		host.err = err
		return host
	}
	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp == 0 || iface.Flags&net.FlagLoopback != 0 {
			// interface is down or loopback
			continue
		}
		addrs, err := iface.Addrs()
		if err != nil {
			host.err = err
			return hostInfo{}
		}
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			if ip == nil || ip.IsLoopback() {
				continue
			}
			ip = ip.To4()
			if ip == nil {
				// not an ipv4 address
				continue
			}
			host.hostIP = ip.String()
			return host
		}
	}
	host.err = errors.New("cannot extract host info")
	return host
}

// GetHostIP returns current host ip address
func GetHostIP() (string, error) {
	extractHostInfo()
	return host.hostIP, host.err
}
