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
	"strings"
	"time"

	"github.com/lindb/common/pkg/logger"
	"github.com/lindb/common/pkg/ltoml"
)

// Engine represents the storage engine configuration shared by all data types.
// Engine-specific configs (e.g. MetricEngine) are nested inside to keep StorageBase clean.
type Engine struct {
	Dir                      string         `env:"DIR" toml:"dir"`
	MaxMemDBSize             ltoml.Size     `env:"MAX_MEMDB_SIZE" toml:"max-memdb-size"`
	MutableMemDBTTL          ltoml.Duration `env:"MUTABLE_MEMDB_TTL" toml:"mutable-memdb-ttl"`
	MaxMemUsageBeforeFlush   float64        `env:"MAX_MEM_USAGE_BEFORE_FLUSH" toml:"max-mem-usage-before-flush"`
	TargetMemUsageAfterFlush float64        `env:"TARGET_MEM_USAGE_AFTER_FLUSH" toml:"target-mem-usage-after-flush"`
	FlushConcurrency         int            `env:"FLUSH_CONCURRENCY" toml:"flush-concurrency"`
	FlushCheckInterval       ltoml.Duration `env:"FLUSH_CHECK_INTERVAL" toml:"flush-check-interval"`
	// Metric holds metric-engine-specific configuration, nested under [storage.engine.metric].
	Metric MetricEngine `envPrefix:"METRIC_" toml:"metric"`
}

func (t *Engine) TOML() string {
	return fmt.Sprintf(`
## The engine directory where the time series data and meta file stores.
## Default: %s
## Env: LINDB_STORAGE_ENGINE_DIR
dir = "%s"

## Flush configuration
##
## The amount of data to build up in each memdb,
## before it is queueing to the immutable list for flushing.
## larger memdb may improve query performance.
## Default: %s
## Env: LINDB_STORAGE_ENGINE_MAX_MEMDB_SIZE
max-memdb-size = "%s"
## Mutable memdb will switch to immutable this often,
## event if the configured memdb-size is not reached.
## Default: %s
## Env: LINDB_STORAGE_ENGINE_MUTABLE_MEMDB_TTL
mutable-memdb-ttl = "%s"
## Global flush operation will be triggered
## when system memory usage is higher than this ratio.
## Default: %.2f
## Env: LINDB_STORAGE_ENGINE_MAX_MEM_USAGE_BEFORE_FLUSH
max-mem-usage-before-flush = %.2f
## Global flush operation will be stopped
## when system memory usage is lower than this ration.
## Default: %.2f
## Env: LINDB_STORAGE_ENGINE_TARGET_MEM_USAGE_AFTER_FLUSH
target-mem-usage-after-flush = %.2f
## concurrency of goroutines for flushing.
## Default: %d
## Env: LINDB_STORAGE_ENGINE_FLUSH_CONCURRENCY
flush-concurrency = %d
## interval between flush checker cycles.
## Default: %s
## Env: LINDB_STORAGE_ENGINE_FLUSH_CHECK_INTERVAL
flush-check-interval = "%s"`,
		strings.ReplaceAll(t.Dir, "\\", "\\\\"),
		strings.ReplaceAll(t.Dir, "\\", "\\\\"),
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
		t.FlushCheckInterval.String(),
		t.FlushCheckInterval.String(),
	)
}

// MetricTOML returns the [storage.engine.metric] section as a TOML string.
func (t *Engine) MetricTOML() string {
	return t.Metric.TOML()
}

// MetricEngine represents metric-engine-specific configuration.
type MetricEngine struct {
	SeriesSequenceCache uint32 `env:"SERIES_SEQ_CACHE" toml:"series-sequence-cache"`
	MetaSequenceCache   uint32 `env:"META_SEQ_CACHE" toml:"meta-sequence-cache"`
}

func (m *MetricEngine) TOML() string {
	return fmt.Sprintf(`
## Cache size for series ID sequence allocation.
## Default: %d
## Env: LINDB_STORAGE_ENGINE_METRIC_SERIES_SEQ_CACHE
series-sequence-cache = %d
## Cache size for meta ID sequence allocation.
## Default: %d
## Env: LINDB_STORAGE_ENGINE_METRIC_META_SEQ_CACHE
meta-sequence-cache = %d`,
		m.SeriesSequenceCache,
		m.SeriesSequenceCache,
		m.MetaSequenceCache,
		m.MetaSequenceCache,
	)
}

// StorageBase represents a storage configuration
type StorageBase struct {
	BrokerEndpoint  string         `env:"BROKER_ENDPOINT" toml:"broker-endpoint"`
	WAL             WAL            `envPrefix:"WAL_" toml:"wal"`
	Engine          Engine         `envPrefix:"ENGINE_" toml:"engine"`
	HTTP            HTTP           `envPrefix:"HTTP_" toml:"http"`
	GRPC            GRPC           `envPrefix:"GRPC_" toml:"grpc"`
	TTLTaskInterval ltoml.Duration `env:"TTL_TASK_INTERVAL" toml:"ttl-task-interval"`
}

// TOML returns StorageBase's toml config string
func (s *StorageBase) TOML() string {
	return fmt.Sprintf(`
## Storage related configuration
[storage]
## interval for how often do ttl job
## Default: %s
## Env: LINDB_STORAGE_TTL_TASK_INTERVAL
ttl-task-interval = "%s"

## Storage HTTP related configuration.
[storage.http]%s

## Storage GRPC related configuration.
[storage.grpc]%s

## Write Ahead Log related configuration.
[storage.wal]%s

## Engine related configuration.
[storage.engine]%s

## Metric engine related configuration.
[storage.engine.metric]%s`,
		s.TTLTaskInterval,
		s.TTLTaskInterval,
		s.HTTP.TOML(),
		s.GRPC.TOML(),
		s.WAL.TOML(),
		s.Engine.TOML(),
		s.Engine.Metric.TOML(),
	)
}

// WAL represents config for write ahead log in storage.
type WAL struct {
	Dir                string         `env:"DIR" toml:"dir"`
	PageSize           ltoml.Size     `env:"PAGE_SIZE" toml:"page-size"`
	RemoveTaskInterval ltoml.Duration `env:"REMOVE_TASK_INTERVAL" toml:"remove-task-interval"`
}

func (rc *WAL) GetPageSize() int64 {
	if rc.PageSize <= 0 {
		return 128 * 1024 * 1024 // 128MB
	}
	if rc.PageSize >= 1024*1024*1024 {
		return 1024 * 1024 * 1024 // 1GB
	}
	return int64(rc.PageSize)
}

func (rc *WAL) TOML() string {
	return fmt.Sprintf(`
## WAL mmaped log directory
## Default: %s
## Env: LINDB_STORAGE_WAL_DIR
dir = "%s"
## page-size is the maximum page size in megabytes of the page file before a new
## file is created, available size is in [128MB, 1GB]
## Default: %s
## Env: LINDB_STORAGE_WAL_PAGE_SIZE
page-size = "%s"
## interval for how often remove expired write ahead log
## Default: %s
## Env: LINDB_STORAGE_WAL_REMOVE_TASK_INTERVAL
remove-task-interval = "%s"`,
		strings.ReplaceAll(rc.Dir, "\\", "\\\\"),
		strings.ReplaceAll(rc.Dir, "\\", "\\\\"),
		rc.PageSize.String(),
		rc.PageSize.String(),
		rc.RemoveTaskInterval.String(),
		rc.RemoveTaskInterval.String(),
	)
}

// Storage represents a storage configuration with common settings
type Storage struct {
	Coordinator RepoState      `envPrefix:"LINDB_COORDINATOR_" toml:"coordinator"`
	Monitor     Monitor        `envPrefix:"LINDB_MONITOR_" toml:"monitor"`
	Logging     logger.Setting `envPrefix:"LINDB_LOGGING_" toml:"logging"`
	StorageBase StorageBase    `envPrefix:"LINDB_STORAGE_" toml:"storage"`
	Query       Query          `envPrefix:"LINDB_QUERY_" toml:"query"`
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
		s.Logging.TOML("LINDB"),
	)
}

// NewDefaultStorageBase returns a new default StorageBase struct
func NewDefaultStorageBase() *StorageBase {
	return &StorageBase{
		TTLTaskInterval: ltoml.Duration(time.Hour * 24),
		HTTP: HTTP{
			Port:         2892,
			IdleTimeout:  ltoml.Duration(time.Minute * 2),
			ReadTimeout:  ltoml.Duration(time.Second * 5),
			WriteTimeout: ltoml.Duration(time.Second * 5),
		},
		GRPC: GRPC{
			Port:                 2891,
			MaxConcurrentStreams: 1024,
			ConnectTimeout:       ltoml.Duration(time.Second * 3),
		},
		WAL: WAL{
			Dir:                filepath.Join(defaultParentDir, "storage", "wal"),
			PageSize:           ltoml.Size(128 * 1024 * 1024),
			RemoveTaskInterval: ltoml.Duration(time.Minute),
		},
		Engine: Engine{
			Dir:                      filepath.Join(defaultParentDir, "storage", "data"),
			MaxMemDBSize:             ltoml.Size(500 * 1024 * 1024),
			MutableMemDBTTL:          ltoml.Duration(time.Minute * 30),
			MaxMemUsageBeforeFlush:   0.75,
			TargetMemUsageAfterFlush: 0.6,
			FlushConcurrency:         int(math.Ceil(float64(runtime.GOMAXPROCS(-1)) / 2)),
			FlushCheckInterval:       ltoml.Duration(time.Minute),
			// Metric-engine-specific defaults nested inside Engine.
			Metric: MetricEngine{
				SeriesSequenceCache: 1000,
				MetaSequenceCache:   100,
			},
		},
	}
}

// NewDefaultStorage creates storage's default config
func NewDefaultStorage() *Storage {
	return &Storage{
		Coordinator: *NewDefaultCoordinator(),
		Query:       *NewDefaultQuery(),
		Monitor:     *NewDefaultMonitor(),
		Logging:     *logger.NewDefaultSetting(),
		StorageBase: *NewDefaultStorageBase(),
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
		logger.NewDefaultSetting().TOML("LINDB"),
	)
}

func checkEngineCfg(engineCfg *Engine) error {
	defaultStorageCfg := NewDefaultStorageBase()
	if engineCfg.Dir == "" {
		return fmt.Errorf("engine dir cannot be empty")
	}
	if engineCfg.MaxMemDBSize <= 0 {
		engineCfg.MaxMemDBSize = defaultStorageCfg.Engine.MaxMemDBSize
	}
	if engineCfg.MutableMemDBTTL <= 0 {
		engineCfg.MutableMemDBTTL = defaultStorageCfg.Engine.MutableMemDBTTL
	}
	if engineCfg.MaxMemUsageBeforeFlush <= 0 {
		engineCfg.MaxMemUsageBeforeFlush = defaultStorageCfg.Engine.MaxMemUsageBeforeFlush
	}
	if engineCfg.TargetMemUsageAfterFlush <= 0 {
		engineCfg.TargetMemUsageAfterFlush = defaultStorageCfg.Engine.TargetMemUsageAfterFlush
	}
	if engineCfg.FlushConcurrency <= 0 {
		engineCfg.FlushConcurrency = defaultStorageCfg.Engine.FlushConcurrency
	}
	if engineCfg.FlushCheckInterval <= 0 {
		engineCfg.FlushCheckInterval = defaultStorageCfg.Engine.FlushCheckInterval
	}
	// Apply metric-engine-specific defaults (now nested inside Engine).
	checkMetricEngineCfg(&engineCfg.Metric)
	return nil
}

func checkMetricEngineCfg(metricCfg *MetricEngine) {
	defaultStorageCfg := NewDefaultStorageBase()
	if metricCfg.SeriesSequenceCache <= 0 {
		metricCfg.SeriesSequenceCache = defaultStorageCfg.Engine.Metric.SeriesSequenceCache
	}
	if metricCfg.MetaSequenceCache <= 0 {
		metricCfg.MetaSequenceCache = defaultStorageCfg.Engine.Metric.MetaSequenceCache
	}
}

// checkStorageBaseCfg checks storage config.
func checkStorageBaseCfg(storageBaseCfg *StorageBase) error {
	if err := checkGRPCCfg(&storageBaseCfg.GRPC); err != nil {
		return err
	}
	defaultStorageCfg := NewDefaultStorageBase()
	if storageBaseCfg.TTLTaskInterval <= 0 {
		storageBaseCfg.TTLTaskInterval = defaultStorageCfg.TTLTaskInterval
	}
	// checkEngineCfg also applies metric-engine defaults (Engine.Metric).
	if err := checkEngineCfg(&storageBaseCfg.Engine); err != nil {
		return err
	}
	return nil
}
