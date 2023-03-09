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

func TestStorage_TOML(t *testing.T) {
	defaultCfg := NewDefaultStorageTOML()
	storageCfg := &Storage{}
	_, err := toml.Decode(defaultCfg, storageCfg)
	assert.NoError(t, err)
	assert.Equal(t, storageCfg.TOML(), defaultCfg)
}

func TestWAL_GetDataSizeLimit(t *testing.T) {
	wal := &WAL{}
	assert.Equal(t, int64(1024*1024), wal.GetDataSizeLimit())
	wal = &WAL{DataSizeLimit: 2 * 1024 * 1024 * 1024}
	assert.Equal(t, int64(1024*1024*1024), wal.GetDataSizeLimit())
	wal = &WAL{DataSizeLimit: 128 * 1024 * 1024}
	assert.Equal(t, int64(128*1024*1024), wal.GetDataSizeLimit())
}

func TestStorage_Env(t *testing.T) {
	cfg := Storage{}
	opts := env.Options{Environment: map[string]string{
		"COORDINATOR_NAMESPACE":                     "ns",
		"COORDINATOR_ENDPOINTS":                     "endpoint1,endpoint2",
		"COORDINATOR_LEASE_TTL":                     "60s",
		"COORDINATOR_TIMEOUT":                       "60s",
		"COORDINATOR_DIAL_TIMEOUT":                  "60s",
		"COORDINATOR_USERNAME":                      "LinDB",
		"COORDINATOR_PASSWORD":                      "pwd",
		"QUERY_CONCURRENCY":                         "100",
		"QUERY_IDLE_TIMEOUT":                        "100s",
		"QUERY_TIMEOUT":                             "120s",
		"STORAGE_BROKER_ENDPOINT":                   "broker_url",
		"STORAGE_TTL_TASK_INTERVAL":                 "2m",
		"STORAGE_HTTP_PORT":                         "3000",
		"STORAGE_HTTP_IDLE_TIMEOUT":                 "120s",
		"STORAGE_HTTP_WRITE_TIMEOUT":                "120s",
		"STORAGE_HTTP_READ_TIMEOUT":                 "2m",
		"STORAGE_GRPC_PORT":                         "2899",
		"STORAGE_GRPC_MAX_CONCURRENT_STREAMS":       "10000",
		"STORAGE_GRPC_CONNECT_TIMEOUT":              "2m",
		"STORAGE_WAL_REMOVE_TASK_INTERVAL":          "2m",
		"STORAGE_WAL_DIR":                           "wal_dir",
		"STORAGE_WAL_DATA_SIZE_LIMIT":               "1Mib",
		"STORAGE_TSDB_DIR":                          "tsdb_dir",
		"STORAGE_TSDB_MAX_MEMDB_SIZE":               "1Mib",
		"STORAGE_TSDB_MUTABLE_MEMDB_TTL":            "2m",
		"STORAGE_TSDB_MAX_MEM_USAGE_BEFORE_FLUSH":   "200.0",
		"STORAGE_TSDB_TARGET_MEM_USAGE_AFTER_FLUSH": "200.0",
		"STORAGE_TSDB_FLUSH_CONCURRENCY":            "2000",
		"STORAGE_TSDB_SERIES_SEQ_CACHE":             "1000",
		"STORAGE_TSDB_META_SEQ_CACHE":               "1000",
		"MONITOR_PUSH_TIMEOUT":                      "2m",
		"MONITOR_REPORT_INTERVAL":                   "2m",
		"MONITOR_URL":                               "monitor_url",
		"LOGGING_DIR":                               "log_dir",
		"LOGGING_LEVEL":                             "fatal",
		"LOGGING_MAX_SIZE":                          "1Mib",
		"LOGGING_MAX_BACKUPS":                       "10",
		"LOGGING_MAX_AGE":                           "20",
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

	assert.Equal(t, uint16(3000), cfg.StorageBase.HTTP.Port)
	assert.Equal(t, ltoml.Duration(time.Second*120), cfg.StorageBase.HTTP.WriteTimeout)
	assert.Equal(t, ltoml.Duration(time.Second*120), cfg.StorageBase.HTTP.ReadTimeout)
	assert.Equal(t, ltoml.Duration(time.Second*120), cfg.StorageBase.HTTP.IdleTimeout)
	assert.Equal(t, uint16(2899), cfg.StorageBase.GRPC.Port)
	assert.Equal(t, 10000, cfg.StorageBase.GRPC.MaxConcurrentStreams)
	assert.Equal(t, ltoml.Duration(time.Second*120), cfg.StorageBase.GRPC.ConnectTimeout)

	assert.Equal(t, "broker_url", cfg.StorageBase.BrokerEndpoint)
	assert.Equal(t, ltoml.Duration(time.Second*120), cfg.StorageBase.TTLTaskInterval)
	assert.Equal(t, ltoml.Duration(time.Second*120), cfg.StorageBase.WAL.RemoveTaskInterval)
	assert.Equal(t, "wal_dir", cfg.StorageBase.WAL.Dir)
	assert.Equal(t, ltoml.Size(1024*1024), cfg.StorageBase.WAL.DataSizeLimit)
	assert.Equal(t, "tsdb_dir", cfg.StorageBase.TSDB.Dir)
	assert.Equal(t, ltoml.Size(1024*1024), cfg.StorageBase.TSDB.MaxMemDBSize)
	assert.Equal(t, ltoml.Duration(time.Second*120), cfg.StorageBase.TSDB.MutableMemDBTTL)
	assert.Equal(t, float64(200.0), cfg.StorageBase.TSDB.MaxMemUsageBeforeFlush)
	assert.Equal(t, float64(200.0), cfg.StorageBase.TSDB.TargetMemUsageAfterFlush)
	assert.Equal(t, 2000, cfg.StorageBase.TSDB.FlushConcurrency)
	assert.Equal(t, uint32(1000), cfg.StorageBase.TSDB.SeriesSequenceCache)
	assert.Equal(t, uint32(1000), cfg.StorageBase.TSDB.MetaSequenceCache)

	assert.Equal(t, ltoml.Duration(time.Second*120), cfg.Monitor.PushTimeout)
	assert.Equal(t, ltoml.Duration(time.Second*120), cfg.Monitor.ReportInterval)
	assert.Equal(t, "monitor_url", cfg.Monitor.URL)
	assert.Equal(t, "log_dir", cfg.Logging.Dir)
	assert.Equal(t, "fatal", cfg.Logging.Level)
	assert.Equal(t, ltoml.Size(1024*1024), cfg.Logging.MaxSize)
	assert.Equal(t, uint16(10), cfg.Logging.MaxBackups)
	assert.Equal(t, uint16(20), cfg.Logging.MaxAge)
}
