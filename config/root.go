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
	"time"

	"github.com/lindb/lindb/pkg/ltoml"
)

// Root represents a root configuration with common settings.
type Root struct {
	Coordinator RepoState `toml:"coordinator"`
	Query       Query     `toml:"query"`
	HTTP        HTTP      `toml:"http"`
	GRPC        GRPC      `toml:"grpc"`
	Monitor     Monitor   `toml:"monitor"`
	Logging     Logging   `toml:"logging"`
}

// TOML returns root's configuration string as toml format.
func (r *Root) TOML() string {
	return fmt.Sprintf(`## Coordinator related configuration.
%s

## Query related configuration.
%s

## Controls how HTTP Server are configured.
[http]%s

## Controls how GRPC Server are configured.
[grpc]%s
%s
%s`,
		r.Coordinator.TOML(),
		r.Query.TOML(),
		r.HTTP.TOML(),
		r.GRPC.TOML(),
		r.Monitor.TOML(),
		r.Logging.TOML(),
	)
}

// NewDefaultBrokerTOML creates root default toml config.
func NewDefaultRootTOML() string {
	return NewDefaultRoot().TOML()
}

func NewDefaultRoot() *Root {
	return &Root{
		Coordinator: *NewDefaultCoordinator(),
		Query:       *NewDefaultQuery(),
		HTTP: HTTP{
			Port:         3000,
			IdleTimeout:  ltoml.Duration(time.Minute * 2),
			ReadTimeout:  ltoml.Duration(time.Second * 5),
			WriteTimeout: ltoml.Duration(time.Second * 5),
		},
		GRPC: GRPC{
			Port:                 3001,
			MaxConcurrentStreams: 1024,
			ConnectTimeout:       ltoml.Duration(time.Second * 3),
		},
		Monitor: *NewDefaultMonitor(),
		Logging: *NewDefaultLogging(),
	}
}
