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
	"strings"
)

// Standalone represents the configuration of standalone mode
type Standalone struct {
	ETCD        ETCD        `envPrefix:"LINDB_ETCD_" toml:"etcd"`
	Coordinator RepoState   `envPrefix:"LINDB_COORDINATOR_" toml:"coordinator"`
	Query       Query       `envPrefix:"LINDB_QUERY_" toml:"query"`
	BrokerBase  BrokerBase  `envPrefix:"LINDB_BROKER_" toml:"broker"`
	StorageBase StorageBase `envPrefix:"LINDB_STORAGE_" toml:"storage"`
	Logging     Logging     `envPrefix:"LINDB_LOGGING_" toml:"logging"`
	Monitor     Monitor     `envPrefix:"LINDB_MONITOR_" toml:"monitor"`
}

// ETCD represents embed etcd's configuration
type ETCD struct {
	Dir string `env:"DIR" toml:"dir"`
	URL string `env:"URL" toml:"url"`
}

// TOML returns ETCD's toml config string
func (etcd *ETCD) TOML() string {
	return fmt.Sprintf(`[etcd]
## Where the ETCD data is stored
## Default: %s
## Env: LINDB_ETCD_DIR
dir = "%s"
## URL to listen on for client traffic 
## If 0.0.0.0 if specified as the IP, 
## etcd listens to the given port on all interfaces.
## If an IP address is given as well as a port, 
## etcd will listen on the given port and interface.
## Default: %s
## Env: LINDB_ETCD_URL
url = "%s"`,
		strings.ReplaceAll(etcd.Dir, "\\", "\\\\"),
		strings.ReplaceAll(etcd.Dir, "\\", "\\\\"),
		etcd.URL,
		etcd.URL,
	)
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
	return fmt.Sprintf(`## Embed ETCD related configuration.
%s

## Coordinator related configuration.
%s

## Query related configuration.
%s
%s
%s
%s
%s`,

		NewDefaultETCD().TOML(),
		NewDefaultCoordinator().TOML(),
		NewDefaultQuery().TOML(),
		NewDefaultBrokerBase().TOML(),
		NewDefaultStorageBase().TOML(),
		NewDefaultLogging().TOML(),
		NewDefaultMonitor().TOML(),
	)
}

// NewDefaultStandalone creates standalone default configuration.
func NewDefaultStandalone() Standalone {
	return Standalone{
		ETCD:        *NewDefaultETCD(),
		Coordinator: *NewDefaultCoordinator(),
		Query:       *NewDefaultQuery(),
		BrokerBase:  *NewDefaultBrokerBase(),
		StorageBase: *NewDefaultStorageBase(),
		Logging:     *NewDefaultLogging(),
		Monitor:     *NewDefaultMonitor(),
	}
}
