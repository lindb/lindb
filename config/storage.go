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
	"time"

	"github.com/lindb/lindb/pkg/ltoml"
)

// TSDB represents the tsdb configuration
type TSDB struct {
	Dir                      string         `toml:"dir"`
	BatchWriteSize           int            `toml:"batch-write-size"`
	BatchPendingSize         int            `toml:"batch-pending-size"`
	BatchTimeout             ltoml.Duration `toml:"batch-timeout"`
	MaxMemDBSize             ltoml.Size     `toml:"max-memdb-size"`
	MaxMemDBTotalSize        ltoml.Size     `toml:"max-memdb-total-size"`
	MaxMemDBNumber           int            `toml:"max-memdb-number"`
	MutableMemDBTTL          ltoml.Duration `toml:"mutable-memdb-ttl"`
	MaxMemUsageBeforeFlush   float64        `toml:"max-mem-usage-before-flush"`
	TargetMemUsageAfterFlush float64        `toml:"target-mem-usage-after-flush"`
	FlushConcurrency         int            `toml:"flush-concurrency"`
	MaxSeriesIDsNumber       int            `toml:"max-seriesIDs"`
	MaxTagKeysNumber         int            `toml:"max-tagKeys"`
}

func (t *TSDB) TOML() string {
	return fmt.Sprintf(`
	## The TSDB directory where the time series data and meta file stores.
	dir = "%s"

	## Write configuration
	## 
	## Shard batch write to memdb with this many buffered metrics
	batch-write-size = %d
	## Shard pending with this many batched metrics in memory
	## if batch-write-size is 100, batch-pending-size is 10
	## at most 1000 metrics will be cached before write
	batch-pending-size = %d
	## Shard will write at least this often,
	## even if the configured batch-size is not reached.
	batch-timeout = "%s"
	
	## Flush configuration
	## 
	## The amount of data to build up in each memdb, 
	## before it is queueing to the immutable list for flushing.
	## Default: 500 MiB, larger memdb may improve query performance
	max-memdb-size = "%s"
	## The maximum size of mutable and immutable memdb of a shard
	## Flush operation will be triggered When this exceeds.
	## Default: 2 GiB
	max-memdb-total-size = "%s"
	## The maximum number mutable and immutable memdb stores in memory.
	## Default: 5. Notice that unlmitied time-range of metrics will make it uncontrollable。
	## If sets to 0, the memdb number is unlimited.
	max-memdb-number = %d
	## Mutable memdb will switch to immutable this often,
	## event if the configured memdb-size is not reached.
	## Default: 1h
	mutable-memdb-ttl = "%s"
	## Global flush operation will be triggered
	## when system memory usage is higher than this ratio.
	## Default: 0.85
	max-mem-usage-before-flush = %.2f
	## Global flush operation will be stopped 
	## when system memory usage is lower than this ration.
	## Defult: 0.60
	target-mem-usage-after-flush = %.2f
	## concurrency of goroutines for flushing. Default: 4
	flush-concurrency = %d

	## Time Series limitation
	## 
	## Limit for time series of metric.
	## Default: 200000
	max-seriesIDs = %d
	## Limit for tagKeys
	## Default: 32
	max-tagKeys = %d`,
		t.Dir,
		t.BatchWriteSize,
		t.BatchPendingSize,
		t.BatchTimeout.String(),
		t.MaxMemDBSize.String(),
		t.MaxMemDBTotalSize.String(),
		t.MaxMemDBNumber,
		t.MutableMemDBTTL.String(),
		t.MaxMemUsageBeforeFlush,
		t.TargetMemUsageAfterFlush,
		t.FlushConcurrency,
		t.MaxSeriesIDsNumber,
		t.MaxTagKeysNumber,
	)
}

// StorageBase represents a storage configuration
type StorageBase struct {
	Indicator int  `toml:"indicator"` // Indicator is unique id under current storage cluster.
	GRPC      GRPC `toml:"grpc"`
	TSDB      TSDB `toml:"tsdb"`
	WAL       WAL  `toml:"wal"`
}

// TOML returns StorageBase's toml config string
func (s *StorageBase) TOML() string {
	return fmt.Sprintf(`[storage]
  ## Indicator is a unique id for identifing each storage node
  ## Make sure indicator on each node is different
  indicator = %d

  [storage.grpc]%s

  [storage.wal]%s

  [storage.tsdb]%s`,
		s.Indicator,
		s.GRPC.TOML(),
		s.WAL.TOML(),
		s.TSDB.TOML(),
	)
}

// WAL represents config for write ahead log in storage.
type WAL struct {
	Dir                string         `toml:"dir"`
	DataSizeLimit      int64          `toml:"data-size-limit"`
	RemoveTaskInterval ltoml.Duration `toml:"remove-task-interval"`
	CheckFlushInterval ltoml.Duration `toml:"check-flush-interval"`
	FlushInterval      ltoml.Duration `toml:"flush-interval"`
	BufferSize         int            `toml:"buffer-size"`
}

func (rc *WAL) GetDataSizeLimit() int64 {
	if rc.DataSizeLimit <= 1 {
		return 1024 * 1024 // 1MB
	}
	if rc.DataSizeLimit >= 1024 {
		return 1024 * 1024 * 1024 // 1GB
	}
	return rc.DataSizeLimit * 1024 * 1024
}

func (rc *WAL) TOML() string {
	return fmt.Sprintf(`
	## WAL mmaped log directory
	dir = "%s"
	## data-size-limit is the maximum size in megabytes of the page file before a new
	## file is created. It defaults to 512 megabytes, available size is in [1MB, 1GB]
	data-size-limit = %d
	## interval for how often a new segment will be created
	remove-task-interval = "%s"
	## interval for how often buffer will be checked if it's available to flush
	check-flush-interval = "%s"
	## interval for how often data will be flushed if data not exceeds the buffer-size
	flush-interval = "%s"
	## will flush if this size of data in kegabytes get buffered
	buffer-size = %d`,
		rc.Dir,
		rc.DataSizeLimit,
		rc.RemoveTaskInterval.String(),
		rc.CheckFlushInterval.String(),
		rc.FlushInterval.String(),
		rc.BufferSize,
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

// NewDefaultStorageBase returns a new default StorageBase struct
func NewDefaultStorageBase() *StorageBase {
	return &StorageBase{
		Indicator: 1,
		GRPC: GRPC{
			Port:                 2891,
			TTL:                  ltoml.Duration(time.Second),
			MaxConcurrentStreams: 30,
			ConnectTimeout:       ltoml.Duration(time.Second * 3),
		},
		WAL: WAL{
			Dir:                filepath.Join(defaultParentDir, "storage/wal"),
			DataSizeLimit:      512,
			RemoveTaskInterval: ltoml.Duration(time.Minute),
			CheckFlushInterval: ltoml.Duration(time.Second),
			FlushInterval:      ltoml.Duration(5 * time.Second),
			BufferSize:         128,
		},
		TSDB: TSDB{
			Dir:                      filepath.Join(defaultParentDir, "storage/data"),
			BatchWriteSize:           100,
			BatchPendingSize:         10,
			BatchTimeout:             ltoml.Duration(time.Millisecond * 500),
			MaxMemDBSize:             ltoml.Size(500 * 1024 * 1024),
			MaxMemDBNumber:           5,
			MaxMemDBTotalSize:        ltoml.Size(2 * 1024 * 1024 * 1024),
			MutableMemDBTTL:          ltoml.Duration(time.Hour),
			MaxMemUsageBeforeFlush:   0.85,
			TargetMemUsageAfterFlush: 0.6,
			FlushConcurrency:         4,
			MaxSeriesIDsNumber:       200000,
			MaxTagKeysNumber:         32,
		},
	}
}

// NewDefaultStorageTOML creates storage's default toml config
func NewDefaultStorageTOML() string {
	return fmt.Sprintf(`%s

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
	if tsdbCfg.BatchWriteSize <= 0 {
		tsdbCfg.BatchWriteSize = defaultStorageCfg.TSDB.BatchWriteSize
	}
	if tsdbCfg.BatchPendingSize <= 0 {
		tsdbCfg.BatchPendingSize = defaultStorageCfg.TSDB.BatchPendingSize
	}
	if tsdbCfg.BatchTimeout <= 0 {
		tsdbCfg.BatchTimeout = defaultStorageCfg.TSDB.BatchTimeout
	}
	if tsdbCfg.MaxMemDBSize <= 0 {
		tsdbCfg.MaxMemDBSize = defaultStorageCfg.TSDB.MaxMemDBSize
	}
	if tsdbCfg.MaxMemDBNumber <= 0 {
		tsdbCfg.MaxMemDBNumber = defaultStorageCfg.TSDB.MaxMemDBNumber
	}
	if tsdbCfg.MaxMemDBTotalSize <= 0 {
		tsdbCfg.MaxMemDBTotalSize = defaultStorageCfg.TSDB.MaxMemDBTotalSize
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
	return checkTSDBCfg(&storageBaseCfg.TSDB)
}
