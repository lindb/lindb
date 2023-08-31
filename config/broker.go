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

// HTTP represents an HTTP level configuration of broker.
type HTTP struct {
	Port         uint16         `env:"PORT" toml:"port"`
	IdleTimeout  ltoml.Duration `env:"IDLE_TIMEOUT" toml:"idle-timeout"`
	WriteTimeout ltoml.Duration `env:"WRITE_TIMEOUT" toml:"write-timeout"`
	ReadTimeout  ltoml.Duration `env:"READ_TIMEOUT" toml:"read-timeout"`
}

func (h *HTTP) TOML() string {
	return fmt.Sprintf(`
## port which the HTTP Server is listening on
## Default: %d
## Env: LINDB_BROKER_HTTP_PORT
## Env: LINDB_STORAGE_HTTP_PORT
## Env: LINDB_ROOT_HTTP_PORT
port = %d
## maximum duration the server should keep established connections alive.
## Default: %s
## Env: LINDB_BROKER_HTTP_IDLE_TIMEOUT
## Env: LINDB_STORAGE_HTTP_IDLE_TIMEOUT
## Env: LINDB_ROOT_HTTP_IDLE_TIMEOUT
idle-timeout = "%s"
## maximum duration before timing out for server writes of the response
## Default: %s
## Env: LINDB_BROKER_HTTP_WRITE_TIMEOUT
## Env: LINDB_STORAGE_HTTP_WRITE_TIMEOUT
## Env: LINDB_ROOT_HTTP_WRITE_TIMEOUT
write-timeout = "%s"
## maximum duration for reading the entire request, including the body.
## Default: %s
## Env: LINDB_BROKER_HTTP_READ_TIMEOUT
## Env: LINDB_STORAGE_HTTP_READ_TIMEOUT
## Env: LINDB_ROOT_HTTP_READ_TIMEOUT
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
	MaxConcurrency int            `env:"CONCURRENCY" toml:"max-concurrency"`
	IngestTimeout  ltoml.Duration `env:"TIMEOUT" toml:"ingest-timeout"`
}

func (i *Ingestion) TOML() string {
	return fmt.Sprintf(`
## How many goroutines can write metrics at the same time.
## If writes requests exceeds the concurrency,
## ingestion HTTP API will be throttled.
## Default: %d
## Env: LINDB_BROKER_INGESTION_CONCURRENCY
max-concurrency = %d
## maximum duration before timeout for server ingesting metrics
## Default: %s
## Env: LINDB_BROKER_INGESTION_TIMEOUT
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
	BatchTimeout   ltoml.Duration `env:"BATCH_TIMEOUT" toml:"batch-timeout"`
	BatchBlockSize ltoml.Size     `env:"BLOCK_SIZE" toml:"batch-block-size"`
	GCTaskInterval ltoml.Duration `env:"GC_INTERVAL" toml:"gc-task-interval"`
}

func (rc *Write) TOML() string {
	return fmt.Sprintf(`
## Broker will write at least this often,
## even if the configured batch-size if not reached.
## Default: %s
## Env: LINDB_BROKER_WRITE_BATCH_TIMEOUT
batch-timeout = "%s"
## Broker will sending block to storage node in this size
## Default: %s
## Env: LINDB_BROKER_WRITE_BLOCK_SIZE
batch-block-size = "%s"
## interval for how often expired write write family garbage collect task execute
## Default: %s
## Env: LINDB_BROKER_WRITE_GC_INTERVAL
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
	SlowSQL   ltoml.Duration `env:"SLOW_SQL" toml:"slow-sql"`
	HTTP      HTTP           `envPrefix:"HTTP_" toml:"http"`
	Ingestion Ingestion      `envPrefix:"INGESTION_" toml:"ingestion"`
	Write     Write          `envPrefix:"WRITE_" toml:"write"`
	GRPC      GRPC           `envPrefix:"GRPC_" toml:"grpc"`
}

// TOML returns broker's base configuration string as toml format.
func (bb *BrokerBase) TOML() string {
	return fmt.Sprintf(`
## Broker related configuration.
[broker]

## Throttle duration for slow sql.
## Default: %s
## Env: LINDB_BROKER_SLOW_SQL
slow-sql = "%s"

## Controls how HTTP Server are configured.
[broker.http]%s

## Ingestion configuration for broker handle ingest request.
[broker.ingestion]%s

## Write configuration for writing replication block.
[broker.write]%s

## Controls how GRPC Server are configured.
[broker.grpc]%s`,
		bb.SlowSQL.String(),
		bb.SlowSQL.String(),
		bb.HTTP.TOML(),
		bb.Ingestion.TOML(),
		bb.Write.TOML(),
		bb.GRPC.TOML(),
	)
}

func NewDefaultBrokerBase() *BrokerBase {
	return &BrokerBase{
		SlowSQL: ltoml.Duration(time.Second * 30),
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
	Coordinator RepoState      `envPrefix:"LINDB_COORDINATOR_" toml:"coordinator"`
	Query       Query          `envPrefix:"LINDB_QUERY_" toml:"query"`
	BrokerBase  BrokerBase     `envPrefix:"LINDB_BROKER_" toml:"broker"`
	Monitor     Monitor        `envPrefix:"LINDB_MONITOR_" toml:"monitor"`
	Logging     logger.Setting `envPrefix:"LINDB_LOGGING_" toml:"logging"`
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
		b.Logging.TOML("LINDB"),
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
		logger.NewDefaultSetting().TOML("LINDB"),
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
