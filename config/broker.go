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
	"runtime"
	"time"

	"github.com/lindb/lindb/pkg/ltoml"
)

// HTTP represents a HTTP level configuration of broker.
type HTTP struct {
	Port         uint16         `toml:"port"`
	IdleTimeout  ltoml.Duration `toml:"idle-timeout"`
	WriteTimeout ltoml.Duration `toml:"write-timeout"`
	ReadTimeout  ltoml.Duration `toml:"read-timeout"`
}

func (h *HTTP) TOML() string {
	return fmt.Sprintf(`
## Controls how HTTP Server are configured.
##
## which port broker's HTTP Server is listening on 
port = %d
## maximum duration the server should keep established connections alive.
## Default: 2m
idle-timeout = "%s"
## maximum duration before timing out for server writes of the response
## Default: 5s
write-timeout = "%s"	
## maximum duration for reading the entire request, including the body.
## Default: 5s
read-timeout = "%s"`,
		h.Port,
		h.IdleTimeout.Duration().String(),
		h.WriteTimeout.Duration().String(),
		h.ReadTimeout.Duration().String(),
	)
}

type Ingestion struct {
	IngestTimeout ltoml.Duration `toml:"ingest-timeout"`
}

func (i *Ingestion) TOML() string {
	return fmt.Sprintf(`
## maximum duration before timeout for server ingesting metrics
## Default: 5s
ingest-timeout = "%s"`,
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
	MaxConcurrency int            `toml:"max-write-concurrency"`
	BatchTimeout   ltoml.Duration `toml:"batch-timout"`
	BatchBlockSize ltoml.Size     `toml:"batch-block-size"`
}

func (rc *Write) TOML() string {
	return fmt.Sprintf(`
## Write Configuration for writing replication block
## 
## How many goroutines can write metrics at the same time.
## If writes requests exceeds the concurrency, 
## ingestion HTTP API will be throttled.
## Default: runtime.GOMAXPROCS(-1) * 2
max-concurrency = %d
## Broker will write at least this often,
## even if the configured batch-size if not reached.
batch-timeout = "%s"
## Broker will sending block to storage node in this size
batch-size = "%s"`,
		rc.MaxConcurrency,
		rc.BatchTimeout.String(),
		rc.BatchBlockSize.String(),
	)
}

// BrokerBase represents a broker configuration
type BrokerBase struct {
	HTTP      HTTP      `toml:"http"`
	Ingestion Ingestion `toml:"ingestion"`
	Write     Write     `toml:"write"`
	User      User      `toml:"user"`
	GRPC      GRPC      `toml:"grpc"`
}

func (bb *BrokerBase) TOML() string {
	return fmt.Sprintf(`
[broker]

[broker.http]%s

[broker.ingestion]%s

[broker.write]%s

[broker.user]%s

[broker.grpc]%s`,
		bb.HTTP.TOML(),
		bb.Ingestion.TOML(),
		bb.Write.TOML(),
		bb.User.TOML(),
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
			IngestTimeout: ltoml.Duration(time.Second * 5),
		},
		Write: Write{
			MaxConcurrency: runtime.GOMAXPROCS(-1) * 2,
			BatchTimeout:   ltoml.Duration(time.Second * 2),
			BatchBlockSize: ltoml.Size(256 * 1024),
		},
		GRPC: GRPC{
			Port:                 9001,
			MaxConcurrentStreams: runtime.GOMAXPROCS(-1) * 2,
			ConnectTimeout:       ltoml.Duration(time.Second * 3),
		},
		User: User{
			UserName: "admin",
			Password: "admin123",
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

// NewDefaultBrokerTOML creates broker default toml config
func NewDefaultBrokerTOML() string {
	return fmt.Sprintf(`%s

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
	// write check
	if brokerBaseCfg.Write.BatchTimeout <= 0 {
		brokerBaseCfg.Write.BatchTimeout = defaultBrokerCfg.Write.BatchTimeout
	}
	if brokerBaseCfg.Write.MaxConcurrency <= 0 {
		brokerBaseCfg.Write.MaxConcurrency = defaultBrokerCfg.Write.MaxConcurrency
	}
	if brokerBaseCfg.Write.BatchBlockSize <= 0 {
		brokerBaseCfg.Write.BatchBlockSize = defaultBrokerCfg.Write.BatchBlockSize
	}

	return nil
}
