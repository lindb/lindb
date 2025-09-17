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
