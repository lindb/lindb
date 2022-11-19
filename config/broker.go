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

// HTTP represents an HTTP level configuration of broker.
type HTTP struct {
	Port         uint16         `toml:"port"`
	IdleTimeout  ltoml.Duration `toml:"idle-timeout"`
	WriteTimeout ltoml.Duration `toml:"write-timeout"`
	ReadTimeout  ltoml.Duration `toml:"read-timeout"`
}

func (h *HTTP) TOML() string {
	return fmt.Sprintf(`
## port which the HTTP Server is listening on
## Default: %d
port = %d
## maximum duration the server should keep established connections alive.
## Default: %s
idle-timeout = "%s"
## maximum duration before timing out for server writes of the response
## Default: %s
write-timeout = "%s"
## maximum duration for reading the entire request, including the body.
## Default: %s
read-timeout = "%s"`,
		h.Port,
		h.Port,
		h.IdleTimeout.Duration().String(),
		h.IdleTimeout.Duration().String(),
		h.WriteTimeout.Duration().String(),
		h.WriteTimeout.Duration().String(),
		h.ReadTimeout.Duration().String(),
		h.ReadTimeout.Duration().String(),
	)
}

type Ingestion struct {
	MaxConcurrency int            `toml:"max-concurrency"`
	IngestTimeout  ltoml.Duration `toml:"ingest-timeout"`
}

func (i *Ingestion) TOML() string {
	return fmt.Sprintf(`
## How many goroutines can write metrics at the same time.
## If writes requests exceeds the concurrency,
## ingestion HTTP API will be throttled.
## Default: %d
max-concurrency = %d
## maximum duration before timeout for server ingesting metrics
## Default: %s
ingest-timeout = "%s"`,
		i.MaxConcurrency,
		i.MaxConcurrency,
		i.IngestTimeout.Duration().String(),
		i.IngestTimeout.Duration().String())
}

// User represents user model
type User struct {
	UserName string `toml:"username" json:"username" binding:"required"`
	Password string `toml:"password" json:"password" binding:"required"`
}

func (u *User) TOML() string {
	return fmt.Sprintf(`
## admin user setting
username = "%s"
password = "%s"`,
		u.UserName,
		u.Password)
}

// Write represents config for write replication in broker.
type Write struct {
	BatchTimeout   ltoml.Duration `toml:"batch-timeout"`
	BatchBlockSize ltoml.Size     `toml:"batch-block-size"`
	GCTaskInterval ltoml.Duration `toml:"gc-task-interval"`
}

func (rc *Write) TOML() string {
	return fmt.Sprintf(`
## Broker will write at least this often,
## even if the configured batch-size if not reached.
## Default: %s
batch-timeout = "%s"
## Broker will sending block to storage node in this size
## Default: %s
batch-block-size = "%s"
## interval for how often expired write write family garbage collect task execute
## Default: %s
gc-task-interval = "%s"`,
		rc.BatchTimeout.String(),
		rc.BatchTimeout.String(),
		rc.BatchBlockSize.String(),
		rc.BatchBlockSize.String(),
		rc.GCTaskInterval.String(),
		rc.GCTaskInterval.String(),
	)
}

// BrokerBase represents a broker configuration
type BrokerBase struct {
	HTTP      HTTP      `toml:"http"`
	Ingestion Ingestion `toml:"ingestion"`
	Write     Write     `toml:"write"`
	GRPC      GRPC      `toml:"grpc"`
}

// TOML returns broker's base configuration string as toml format.
func (bb *BrokerBase) TOML() string {
	return fmt.Sprintf(`
## Broker related configuration.
[broker]

## Controls how HTTP Server are configured.
[broker.http]%s

## Ingestion configuration for broker handle ingest request.
[broker.ingestion]%s

## Write configuration for writing replication block.
[broker.write]%s

## Controls how GRPC Server are configured.
[broker.grpc]%s`,
		bb.HTTP.TOML(),
		bb.Ingestion.TOML(),
		bb.Write.TOML(),
		bb.GRPC.TOML(),
	)
}

func NewDefaultBrokerBase() *BrokerBase {
	return &BrokerBase{
		HTTP: HTTP{
			Port:         9000,
			IdleTimeout:  ltoml.Duration(time.Minute * 2),
			ReadTimeout:  ltoml.Duration(time.Second * 5),
			WriteTimeout: ltoml.Duration(time.Second * 5),
		},
		Ingestion: Ingestion{
			MaxConcurrency: 256,
			IngestTimeout:  ltoml.Duration(time.Second * 5),
		},
		Write: Write{
			BatchTimeout:   ltoml.Duration(time.Second * 2),
			BatchBlockSize: ltoml.Size(256 * 1024),
			GCTaskInterval: ltoml.Duration(time.Minute),
		},
		GRPC: GRPC{
			Port:                 9001,
			MaxConcurrentStreams: 1024,
			ConnectTimeout:       ltoml.Duration(time.Second * 3),
		},
	}
}

// Broker represents a broker configuration with common settings
type Broker struct {
	Coordinator RepoState  `toml:"coordinator"`
	Query       Query      `toml:"query"`
	BrokerBase  BrokerBase `toml:"broker"`
	Monitor     Monitor    `toml:"monitor"`
	Logging     Logging    `toml:"logging"`
}

// TOML returns broker's configuration string as toml format.
func (b *Broker) TOML() string {
	return fmt.Sprintf(`## Coordinator related configuration.
%s

## Query related configuration.
%s
%s
%s
%s`,
		b.Coordinator.TOML(),
		b.Query.TOML(),
		b.BrokerBase.TOML(),
		b.Monitor.TOML(),
		b.Logging.TOML(),
	)
}

// NewDefaultBrokerTOML creates broker default toml config
func NewDefaultBrokerTOML() string {
	return fmt.Sprintf(`## Coordinator related configuration.
%s

## Query related configuration.
%s
%s
%s
%s`,
		NewDefaultCoordinator().TOML(),
		NewDefaultQuery().TOML(),
		NewDefaultBrokerBase().TOML(),
		NewDefaultMonitor().TOML(),
		NewDefaultLogging().TOML(),
	)
}

// checkBrokerBaseCfg checks broker base configuration, if not set using default value.
func checkBrokerBaseCfg(brokerBaseCfg *BrokerBase) error {
	if err := checkGRPCCfg(&brokerBaseCfg.GRPC); err != nil {
		return err
	}
	defaultBrokerCfg := NewDefaultBrokerBase()
	// http check
	if brokerBaseCfg.HTTP.Port <= 0 {
		return fmt.Errorf("http port cannot be empty")
	}
	if brokerBaseCfg.HTTP.ReadTimeout <= 0 {
		brokerBaseCfg.HTTP.ReadTimeout = defaultBrokerCfg.HTTP.ReadTimeout
	}
	if brokerBaseCfg.HTTP.WriteTimeout <= 0 {
		brokerBaseCfg.HTTP.WriteTimeout = defaultBrokerCfg.HTTP.WriteTimeout
	}
	if brokerBaseCfg.HTTP.IdleTimeout <= 0 {
		brokerBaseCfg.HTTP.IdleTimeout = defaultBrokerCfg.HTTP.IdleTimeout
	}

	// ingestion
	if brokerBaseCfg.Ingestion.IngestTimeout <= 0 {
		brokerBaseCfg.Ingestion.IngestTimeout = defaultBrokerCfg.Ingestion.IngestTimeout
	}
	if brokerBaseCfg.Ingestion.MaxConcurrency <= 0 {
		brokerBaseCfg.Ingestion.MaxConcurrency = defaultBrokerCfg.Ingestion.MaxConcurrency
	}
	// write check
	if brokerBaseCfg.Write.BatchTimeout <= 0 {
		brokerBaseCfg.Write.BatchTimeout = defaultBrokerCfg.Write.BatchTimeout
	}
	if brokerBaseCfg.Write.BatchBlockSize <= 0 {
		brokerBaseCfg.Write.BatchBlockSize = defaultBrokerCfg.Write.BatchBlockSize
	}
	if brokerBaseCfg.Write.GCTaskInterval <= 0 {
		brokerBaseCfg.Write.GCTaskInterval = defaultBrokerCfg.Write.GCTaskInterval
	}

	return nil
}
