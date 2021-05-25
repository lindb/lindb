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

package config

import (
	"fmt"
	"path/filepath"
)

// Standalone represents the configuration of standalone mode
type Standalone struct {
	ETCD        ETCD        `toml:"etcd"`
	BrokerBase  BrokerBase  `toml:"broker"`
	StorageBase StorageBase `toml:"storage"`
	Logging     Logging     `toml:"logging"`
	Monitor     Monitor     `toml:"monitor"`
}

// ETCD represents embed etcd's configuration
type ETCD struct {
	Dir string `toml:"dir"`
	URL string `toml:"url"`
}

// TOML returns ETCD's toml config string
func (etcd *ETCD) TOML() string {
	return fmt.Sprintf(`## Config for embedding etcd server
[etcd]
  ## Where the ETCD data is stored
  dir = "%s"

  ## URL to listen on for client traffic 
  ## If 0.0.0.0 if specified as the IP, 
  ## etcd listens to the given port on all interfaces.
  ## If an IP address is given as well as a port, 
  ## etcd will listen on the given port and interface.
  ## example: http://10.0.0.1:2379
  url = "%s"
`,
		etcd.Dir,
		etcd.URL)
}

// NewDefaultETCD returns a default ETCD
func NewDefaultETCD() *ETCD {
	return &ETCD{
		Dir: filepath.Join(defaultParentDir, "coordinator"),
		URL: "http://localhost:2379",
	}
}

// NewDefaultStandaloneTOML creates default toml config for standalone
func NewDefaultStandaloneTOML() string {
	return fmt.Sprintf(`%s

%s

%s

%s

%s`,

		NewDefaultETCD().TOML(),
		NewDefaultBrokerBase().TOML(),
		NewDefaultStorageBase().TOML(),
		NewDefaultLogging().TOML(),
		NewDefaultMonitor().TOML(),
	)
}
