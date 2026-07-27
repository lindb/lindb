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

package spi

import (
	jsoniter "github.com/json-iterator/go"

	"github.com/lindb/lindb/constants"
	"github.com/lindb/lindb/pkg/encoding"
	"github.com/lindb/lindb/pkg/timeutil"
)

func init() {
	// register json encoder/decoder for table handle
	jsoniter.RegisterTypeEncoder("spi.TableHandle", &encoding.JSONEncoder[TableHandle]{})
	jsoniter.RegisterTypeDecoder("spi.TableHandle", &encoding.JSONDecoder[TableHandle]{})
}

type DatasourceKind int

const (
	InfoSchema DatasourceKind = iota + 1
	Metric
	Log
	Trace
	Streaming
)

func (kind DatasourceKind) String() string {
	switch kind {
	case InfoSchema:
		return constants.InformationSchema
	case Metric:
		return "metric"
	case Log:
		return "log"
	case Trace:
		return "trace"
	case Streaming:
		return "streaming"
	default:
		return "unknwon"
	}
}

// TableHandle represents a table handle that connect the storage engine.
type TableHandle interface {
	// SetTimeRange sets the time range of table handle.
	SetTimeRange(timeRange timeutil.TimeRange)
	// GetTimeRange returns the time range of table handle.
	GetTimeRange() timeutil.TimeRange
	SetInterval(interval timeutil.Interval)
	GetInterval() timeutil.Interval
	// Kind returns the kind of data source.
	Kind() DatasourceKind
	// String returns table info, format: ${database}:${namespace}:${tableName}
	String() string
}

// ShardCursor is a shard-local pagination position.
// (Timestamp, LogID) uniquely identifies a row within a single shard across leader switches
// because the per-row timestamp and per-segment logID are stable properties of the stored data.
type ShardCursor struct {
	Timestamp int64  `json:"timestamp"` // per-row log timestamp in nanoseconds
	LogID     uint32 `json:"logID"`     // per-segment log sequence number
}

// Paginator is an optional interface for TableHandle implementations that support
// cursor-based pagination. Use type assertion to check for support:
//
//	if p, ok := tableHandle.(spi.Paginator); ok { ... }
//
// Pagination uses a per-shard cursor map: each shard advances independently,
// so storage nodes only need to filter their own shard's cursor without
// cross-shard comparison.
type Paginator interface {
	// SetShardCursors sets the full per-shard cursor map.
	// Key: globally unique shard ID; Value: the last-seen position within that shard.
	SetShardCursors(cursors map[int64]ShardCursor)
	// GetShardCursor returns the cursor for the given shard, and whether one exists.
	GetShardCursor(shardID int64) (ShardCursor, bool)
	// HasAnyCursor returns true if at least one shard cursor has been set.
	HasAnyCursor() bool
	// SetLimit sets the maximum number of rows to return per page.
	SetLimit(limit int64)
	// GetLimit returns the page size limit (0 means use default).
	GetLimit() int64
}
