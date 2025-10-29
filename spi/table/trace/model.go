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

package trace

import (
	"fmt"

	"github.com/lindb/lindb/pkg/encoding"
	"github.com/lindb/lindb/pkg/timeutil"
	"github.com/lindb/lindb/spi"
	"github.com/lindb/lindb/sql/tree"
	"github.com/lindb/lindb/storage/store"
)

func init() {
	// register table/column handle
	encoding.RegisterNodeType(TableHandle{})
	encoding.RegisterNodeType(ColumnHandle{})

	spi.RegisterCreateTableFn(spi.Trace, func(db, ns, name string) spi.TableHandle {
		return &TableHandle{
			Database:  db,
			Namespace: ns,
		}
	})
}

type TableHandle struct {
	Database  string `json:"database"`
	Namespace string `json:"namespace"`

	TimeRange timeutil.TimeRange `json:"timeRange"`
}

func (t *TableHandle) SetTimeRange(timeRange timeutil.TimeRange) {
	t.TimeRange = timeRange
}

func (t *TableHandle) GetTimeRange() timeutil.TimeRange {
	return t.TimeRange
}

func (t *TableHandle) SetInterval(interval timeutil.Interval) {
}

func (t *TableHandle) GetInterval() timeutil.Interval {
	// fake interval
	return timeutil.Interval(10_1000)
}

func (t *TableHandle) Kind() spi.DatasourceKind {
	return spi.Log
}

func (t *TableHandle) String() string {
	return fmt.Sprintf("%s:%s", t.Database, t.Namespace)
}

type ColumnHandle struct {
	Downsampling tree.FuncName `json:"downsampling"`
	Aggregation  tree.FuncName `json:"aggregation"`
}

type Partition struct {
	tableScan *TableScan
	shard     store.Shard
	segments  []store.Segment
}
