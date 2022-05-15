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
	"math"
	"path/filepath"
	"runtime"
	"time"

	"github.com/lindb/lindb/pkg/ltoml"
)

// TSDB represents the tsdb configuration.
type TSDB struct {
	Dir                      string         `toml:"dir"`
	MaxMemDBSize             ltoml.Size     `toml:"max-memdb-size"`
	MutableMemDBTTL          ltoml.Duration `toml:"mutable-memdb-ttl"`
	MaxMemUsageBeforeFlush   float64        `toml:"max-mem-usage-before-flush"`
	TargetMemUsageAfterFlush float64        `toml:"target-mem-usage-after-flush"`
	FlushConcurrency         int            `toml:"flush-concurrency"`
	MaxSeriesIDsNumber       int            `toml:"max-seriesIDs"`
	SeriesSequenceCache      uint32         `toml:"series-sequence-cache"`
	MetaSequenceCache        uint32         `toml:"meta-sequence-cache"`
	MaxTagKeysNumber         int            `toml:"max-tagKeys"`
}

func (t *TSDB) TOML() string {
	return fmt.Sprintf(`
## The TSDB directory where the time series data and meta file stores.
## Default: %s
dir = "%s"

## Flush configuration
## 
## The amount of data to build up in each memdb, 
## before it is queueing to the immutable list for flushing.
## larger memdb may improve query performance.
## Default: %s
max-memdb-size = "%s"
## Mutable memdb will switch to immutable this often,
## event if the configured memdb-size is not reached.
## Default: %s
mutable-memdb-ttl = "%s"
## Global flush operation will be triggered
## when system memory usage is higher than this ratio.
## Default: %.2f
max-mem-usage-before-flush = %.2f
## Global flush operation will be stopped 
## when system memory usage is lower than this ration.
## Default: %.2f
target-mem-usage-after-flush = %.2f
## concurrency of goroutines for flushing.
## Default: %d
flush-concurrency = %d

## Time Series limitation
## 
## Limit for time series of metric.
## Default: %d
max-seriesIDs = %d
## Limit for tagKeys
## Default: %d
max-tagKeys = %d`,
		t.Dir,
		t.Dir,
		t.MaxMemDBSize.String(),
		t.MaxMemDBSize.String(),
		t.MutableMemDBTTL.String(),
		t.MutableMemDBTTL.String(),
		t.MaxMemUsageBeforeFlush,
		t.MaxMemUsageBeforeFlush,
		t.TargetMemUsageAfterFlush,
		t.TargetMemUsageAfterFlush,
		t.FlushConcurrency,
		t.FlushConcurrency,
		t.MaxSeriesIDsNumber,
		t.MaxSeriesIDsNumber,
		t.MaxTagKeysNumber,
		t.MaxTagKeysNumber,
	)
}

// StorageBase represents a storage configuration
type StorageBase struct {
	Indicator       int            `toml:"indicator"`       // Indicator is unique id under current storage cluster.
	BrokerEndpoint  string         `toml:"broker-endpoint"` // Broker http endpoint, auto register current storage cluster.
	TTLTaskInterval ltoml.Duration `toml:"ttl-task-interval"`
	HTTP            HTTP           `toml:"http"`
	GRPC            GRPC           `toml:"grpc"`
	TSDB            TSDB           `toml:"tsdb"`
	WAL             WAL            `toml:"wal"`
}

// TOML returns StorageBase's toml config string
func (s *StorageBase) TOML() string {
	return fmt.Sprintf(`
## Storage related configuration
[storage]
## Indicator is a unique id for identifing each storage node
## Make sure indicator on each node is different
## Default: %d
indicator = %d
## interval for how often do ttl job
## Default: %s
ttl-task-interval = "%s"
## Broker http endpoint which storage self register address
## Default: %s
broker-endpoint = "%s"

## Storage HTTP related configuration.
[storage.http]%s

## Storage GRPC related configuration.
[storage.grpc]%s

## Write Ahead Log related configuration.
[storage.wal]%s

## TSDB related configuration.
[storage.tsdb]%s`,
		s.Indicator,
		s.Indicator,
		s.TTLTaskInterval,
		s.TTLTaskInterval,
		s.BrokerEndpoint,
		s.BrokerEndpoint,
		s.HTTP.TOML(),
		s.GRPC.TOML(),
		s.WAL.TOML(),
		s.TSDB.TOML(),
	)
}

// WAL represents config for write ahead log in storage.
type WAL struct {
	Dir                string         `toml:"dir"`
	DataSizeLimit      ltoml.Size     `toml:"data-size-limit"`
	RemoveTaskInterval ltoml.Duration `toml:"remove-task-interval"`
}

func (rc *WAL) GetDataSizeLimit() int64 {
	if rc.DataSizeLimit <= 0 {
		return 1024 * 1024 // 1MB
	}
	if rc.DataSizeLimit >= 1024*1024*1024 {
		return 1024 * 1024 * 1024 // 1GB
	}
	return int64(rc.DataSizeLimit)
}

func (rc *WAL) TOML() string {
	return fmt.Sprintf(`
## WAL mmaped log directory
## Default: %s
dir = "%s"
## data-size-limit is the maximum size in megabytes of the page file before a new
## file is created. It defaults to 512 megabytes, available size is in [1MB, 1GB]
## Default: %s
data-size-limit = "%s"
## interval for how often remove expired write ahead log
## Default: %s
remove-task-interval = "%s"`,
		rc.Dir,
		rc.Dir,
		rc.DataSizeLimit.String(),
		rc.DataSizeLimit.String(),
		rc.RemoveTaskInterval.String(),
		rc.RemoveTaskInterval.String(),
	)
}

// Storage represents a storage configuration with common settings
type Storage struct {
	Coordinator RepoState   `toml:"coordinator"`
	Query       Query       `toml:"query"`
	StorageBase StorageBase `toml:"storage"`
	Monitor     Monitor     `toml:"monitor"`
	Logging     Logging     `toml:"logging"`
}

// TOML returns storage's configuration string as toml format.
func (s *Storage) TOML() string {
	return fmt.Sprintf(`## Coordinator related configuration.
%s

## Query related configuration.
%s
%s
%s
%s`,
		s.Coordinator.TOML(),
		s.Query.TOML(),
		s.StorageBase.TOML(),
		s.Monitor.TOML(),
		s.Logging.TOML(),
	)
}

// NewDefaultStorageBase returns a new default StorageBase struct
func NewDefaultStorageBase() *StorageBase {
	return &StorageBase{
		Indicator:       1,
		TTLTaskInterval: ltoml.Duration(time.Hour * 24),
		BrokerEndpoint:  "http://localhost:9000",
		HTTP: HTTP{
			Port:         2892,
			IdleTimeout:  ltoml.Duration(time.Minute * 2),
			ReadTimeout:  ltoml.Duration(time.Second * 5),
			WriteTimeout: ltoml.Duration(time.Second * 5),
		},
		GRPC: GRPC{
			Port:                 2891,
			MaxConcurrentStreams: runtime.GOMAXPROCS(-1) * 40,
			ConnectTimeout:       ltoml.Duration(time.Second * 3),
		},
		WAL: WAL{
			Dir:                filepath.Join(defaultParentDir, "storage", "wal"),
			DataSizeLimit:      ltoml.Size(128 * 1024 * 1024),
			RemoveTaskInterval: ltoml.Duration(time.Minute),
		},
		TSDB: TSDB{
			Dir:                      filepath.Join(defaultParentDir, "storage", "data"),
			MaxMemDBSize:             ltoml.Size(500 * 1024 * 1024),
			MutableMemDBTTL:          ltoml.Duration(time.Minute * 30),
			MaxMemUsageBeforeFlush:   0.75,
			TargetMemUsageAfterFlush: 0.6,
			FlushConcurrency:         int(math.Ceil(float64(runtime.GOMAXPROCS(-1)) / 2)),
			MaxSeriesIDsNumber:       200000,
			SeriesSequenceCache:      1000,
			MetaSequenceCache:        100,
			MaxTagKeysNumber:         32,
		},
	}
}

// NewDefaultStorageTOML creates storage's default toml config
func NewDefaultStorageTOML() string {
	return fmt.Sprintf(`## Coordinator related configuration.
%s

## Query related configuration.
%s
%s
%s
%s`,
		NewDefaultCoordinator().TOML(),
		NewDefaultQuery().TOML(),
		NewDefaultStorageBase().TOML(),
		NewDefaultMonitor().TOML(),
		NewDefaultLogging().TOML(),
	)
}

func checkTSDBCfg(tsdbCfg *TSDB) error {
	defaultStorageCfg := NewDefaultStorageBase()
	if tsdbCfg.Dir == "" {
		return fmt.Errorf("tsdb dir cannot be empty")
	}
	if tsdbCfg.MaxMemDBSize <= 0 {
		tsdbCfg.MaxMemDBSize = defaultStorageCfg.TSDB.MaxMemDBSize
	}
	if tsdbCfg.MutableMemDBTTL <= 0 {
		tsdbCfg.MutableMemDBTTL = defaultStorageCfg.TSDB.MutableMemDBTTL
	}
	if tsdbCfg.MaxMemUsageBeforeFlush <= 0 {
		tsdbCfg.MaxMemUsageBeforeFlush = defaultStorageCfg.TSDB.MaxMemUsageBeforeFlush
	}
	if tsdbCfg.TargetMemUsageAfterFlush <= 0 {
		tsdbCfg.TargetMemUsageAfterFlush = defaultStorageCfg.TSDB.TargetMemUsageAfterFlush
	}
	if tsdbCfg.FlushConcurrency <= 0 {
		tsdbCfg.FlushConcurrency = defaultStorageCfg.TSDB.FlushConcurrency
	}
	if tsdbCfg.MaxSeriesIDsNumber <= 0 {
		tsdbCfg.MaxSeriesIDsNumber = defaultStorageCfg.TSDB.MaxSeriesIDsNumber
	}
	if tsdbCfg.SeriesSequenceCache <= 0 {
		tsdbCfg.SeriesSequenceCache = defaultStorageCfg.TSDB.SeriesSequenceCache
	}
	if tsdbCfg.MetaSequenceCache <= 0 {
		tsdbCfg.MetaSequenceCache = defaultStorageCfg.TSDB.MetaSequenceCache
	}
	if tsdbCfg.MaxTagKeysNumber <= 0 {
		tsdbCfg.MaxTagKeysNumber = defaultStorageCfg.TSDB.MaxTagKeysNumber
	}
	return nil
}

func checkStorageBaseCfg(storageBaseCfg *StorageBase) error {
	if storageBaseCfg.Indicator <= 0 {
		return fmt.Errorf("indicator must > 0")
	}

	if err := checkGRPCCfg(&storageBaseCfg.GRPC); err != nil {
		return err
	}
	defaultStorageCfg := NewDefaultStorageBase()
	if storageBaseCfg.TTLTaskInterval <= 0 {
		storageBaseCfg.TTLTaskInterval = defaultStorageCfg.TTLTaskInterval
	}
	return checkTSDBCfg(&storageBaseCfg.TSDB)
}
