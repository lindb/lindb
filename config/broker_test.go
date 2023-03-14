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
	"testing"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/caarlos0/env/v7"
	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/pkg/ltoml"
)

func TestBroker_TOML(t *testing.T) {
	defaultCfg := NewDefaultBrokerTOML()
	brokerCfg := &Broker{}
	_, err := toml.Decode(defaultCfg, brokerCfg)
	assert.NoError(t, err)
	assert.Equal(t, brokerCfg.TOML(), defaultCfg)

	assert.NotEmpty(t, (&User{}).TOML())
}

func TestDumpExampleCfg(t *testing.T) {
	assert.NoError(t, ltoml.WriteConfig("root.toml.example", NewDefaultRootTOML()))
	assert.NoError(t, ltoml.WriteConfig("broker.toml.example", NewDefaultBrokerTOML()))
	assert.NoError(t, ltoml.WriteConfig("storage.toml.example", NewDefaultStorageTOML()))
	assert.NoError(t, ltoml.WriteConfig("standalone.toml.example", NewDefaultStandaloneTOML()))
}

func TestBroker_Env(t *testing.T) {
	cfg := Broker{}
	opts := env.Options{Environment: map[string]string{
		"LINDB_COORDINATOR_NAMESPACE":              "ns",
		"LINDB_COORDINATOR_ENDPOINTS":              "endpoint1,endpoint2",
		"LINDB_COORDINATOR_LEASE_TTL":              "60s",
		"LINDB_COORDINATOR_TIMEOUT":                "60s",
		"LINDB_COORDINATOR_DIAL_TIMEOUT":           "60s",
		"LINDB_COORDINATOR_USERNAME":               "LinDB",
		"LINDB_COORDINATOR_PASSWORD":               "pwd",
		"LINDB_QUERY_CONCURRENCY":                  "100",
		"LINDB_QUERY_IDLE_TIMEOUT":                 "100s",
		"LINDB_QUERY_TIMEOUT":                      "120s",
		"LINDB_BROKER_SLOW_SQL":                    "120s",
		"LINDB_BROKER_HTTP_PORT":                   "3000",
		"LINDB_BROKER_HTTP_IDLE_TIMEOUT":           "120s",
		"LINDB_BROKER_HTTP_WRITE_TIMEOUT":          "120s",
		"LINDB_BROKER_HTTP_READ_TIMEOUT":           "2m",
		"LINDB_BROKER_INGESTION_CONCURRENCY":       "100",
		"LINDB_BROKER_INGESTION_TIMEOUT":           "2m",
		"LINDB_BROKER_WRITE_BATCH_TIMEOUT":         "2m",
		"LINDB_BROKER_WRITE_BLOCK_SIZE":            "1Mib",
		"LINDB_BROKER_WRITE_GC_INTERVAL":           "2m",
		"LINDB_BROKER_GRPC_PORT":                   "2899",
		"LINDB_BROKER_GRPC_MAX_CONCURRENT_STREAMS": "10000",
		"LINDB_BROKER_GRPC_CONNECT_TIMEOUT":        "2m",
		"LINDB_MONITOR_PUSH_TIMEOUT":               "2m",
		"LINDB_MONITOR_REPORT_INTERVAL":            "2m",
		"LINDB_MONITOR_URL":                        "monitor_url",
		"LINDB_LOGGING_DIR":                        "log_dir",
		"LINDB_LOGGING_LEVEL":                      "fatal",
		"LINDB_LOGGING_MAX_SIZE":                   "1Mib",
		"LINDB_LOGGING_MAX_BACKUPS":                "10",
		"LINDB_LOGGING_MAX_AGE":                    "20",
	}}
	err := env.Parse(&cfg, opts)
	assert.NoError(t, err)

	assert.Equal(t, "ns", cfg.Coordinator.Namespace)
	assert.Equal(t, []string{"endpoint1", "endpoint2"}, cfg.Coordinator.Endpoints)
	assert.Equal(t, ltoml.Duration(time.Second*60), cfg.Coordinator.LeaseTTL)
	assert.Equal(t, ltoml.Duration(time.Second*60), cfg.Coordinator.Timeout)
	assert.Equal(t, ltoml.Duration(time.Second*60), cfg.Coordinator.DialTimeout)
	assert.Equal(t, "LinDB", cfg.Coordinator.Username)
	assert.Equal(t, "pwd", cfg.Coordinator.Password)
	assert.Equal(t, 100, cfg.Query.QueryConcurrency)
	assert.Equal(t, ltoml.Duration(time.Second*100), cfg.Query.IdleTimeout)
	assert.Equal(t, ltoml.Duration(time.Second*120), cfg.Query.Timeout)
	assert.Equal(t, ltoml.Duration(time.Second*120), cfg.BrokerBase.SlowSQL)
	assert.Equal(t, uint16(3000), cfg.BrokerBase.HTTP.Port)
	assert.Equal(t, ltoml.Duration(time.Second*120), cfg.BrokerBase.HTTP.WriteTimeout)
	assert.Equal(t, ltoml.Duration(time.Second*120), cfg.BrokerBase.HTTP.ReadTimeout)
	assert.Equal(t, ltoml.Duration(time.Second*120), cfg.BrokerBase.HTTP.IdleTimeout)
	assert.Equal(t, 100, cfg.BrokerBase.Ingestion.MaxConcurrency)
	assert.Equal(t, ltoml.Duration(time.Second*120), cfg.BrokerBase.Ingestion.IngestTimeout)
	assert.Equal(t, ltoml.Duration(time.Second*120), cfg.BrokerBase.Write.BatchTimeout)
	assert.Equal(t, ltoml.Duration(time.Second*120), cfg.BrokerBase.Write.GCTaskInterval)
	assert.Equal(t, ltoml.Size(1024*1024), cfg.BrokerBase.Write.BatchBlockSize)
	assert.Equal(t, uint16(2899), cfg.BrokerBase.GRPC.Port)
	assert.Equal(t, 10000, cfg.BrokerBase.GRPC.MaxConcurrentStreams)
	assert.Equal(t, ltoml.Duration(time.Second*120), cfg.BrokerBase.GRPC.ConnectTimeout)
	assert.Equal(t, ltoml.Duration(time.Second*120), cfg.Monitor.PushTimeout)
	assert.Equal(t, ltoml.Duration(time.Second*120), cfg.Monitor.ReportInterval)
	assert.Equal(t, "monitor_url", cfg.Monitor.URL)
	assert.Equal(t, "log_dir", cfg.Logging.Dir)
	assert.Equal(t, "fatal", cfg.Logging.Level)
	assert.Equal(t, ltoml.Size(1024*1024), cfg.Logging.MaxSize)
	assert.Equal(t, uint16(10), cfg.Logging.MaxBackups)
	assert.Equal(t, uint16(20), cfg.Logging.MaxAge)
}
