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

	"github.com/lindb/common/pkg/logger"
	"github.com/lindb/common/pkg/ltoml"
)

type Streaming struct {
	Coordinator   RepoState      `envPrefix:"LINDB_COORDINATOR_" toml:"coordinator"`
	StreamingBase StreamingBase  `envPrefix:"LINDB_STREAMING_" toml:"streaming"`
	Monitor       Monitor        `envPrefix:"LINDB_MONITOR_" toml:"monitor"`
	Logging       logger.Setting `envPrefix:"LINDB_LOGGING_" toml:"logging"`
}

// NewDefaultStreaming creates default config of streaming.
func NewDefaultStreaming() *Streaming {
	return &Streaming{
		Coordinator: *NewDefaultCoordinator(),
		Monitor:     *NewDefaultMonitor(),
		Logging:     *logger.NewDefaultSetting(),
	}
}

// NewDefaultStreamingTOML creates default toml config of streaming.
func NewDefaultStreamingTOML() string {
	return fmt.Sprintf(`## Coordinator related configuration.
%s

%s
%s
%s`,
		NewDefaultCoordinator().TOML(),
		NewDefaultStreamingBase().TOML(),
		NewDefaultMonitor().TOML(),
		logger.NewDefaultSetting().TOML("LINDB"),
	)
}

// StreamingBase represents a streaming configuration
type StreamingBase struct {
	Namespace string `env:"NAMESPACE" toml:"namespace"`
	GRPC      GRPC   `envPrefix:"GRPC_" toml:"grpc"`
	HTTP      HTTP   `envPrefix:"HTTP_" toml:"http"`
}

// NewDefaultStreamingBase returns a new default StreamingBase struct.
func NewDefaultStreamingBase() *StreamingBase {
	return &StreamingBase{
		Namespace: "observer-streaming",
		HTTP: HTTP{
			Port:         2992,
			IdleTimeout:  ltoml.Duration(time.Minute * 2),
			ReadTimeout:  ltoml.Duration(time.Second * 5),
			WriteTimeout: ltoml.Duration(time.Second * 5),
		},
		GRPC: GRPC{
			Port:                 2991,
			MaxConcurrentStreams: 1024,
			ConnectTimeout:       ltoml.Duration(time.Second * 3),
		},
	}
}

// TOML returns StreamingBase's toml config string
func (s *StreamingBase) TOML() string {
	return fmt.Sprintf(`[streaming]
namespace = "%s"

## Streaming HTTP related configuration.
[streaming.http]%s

## Streaming GRPC related configuration.
[streaming.grpc]%s`,
		s.Namespace,
		s.HTTP.TOML(),
		s.GRPC.TOML(),
	)
}

// TOML returns configuration string as toml format of streaming.
func (s *Streaming) TOML() string {
	return fmt.Sprintf(`## Coordinator related configuration.
%s

%s
%s`,
		s.Coordinator.TOML(),
		s.Monitor.TOML(),
		s.Logging.TOML("LINDB"),
	)
}
