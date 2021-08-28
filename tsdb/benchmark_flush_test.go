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

package tsdb

import (
	"testing"

	"github.com/lindb/lindb/kv"
	"github.com/lindb/lindb/pkg/fileutil"
	protoMetricsV1 "github.com/lindb/lindb/proto/gen/v1/metrics"
	"github.com/lindb/lindb/series/field"
	"github.com/lindb/lindb/tsdb/memdb"
	"github.com/lindb/lindb/tsdb/tblstore/metricsdata"
)

const tmpBufDir = "./test_flush_dir/buf"
const tmpStoreDir = "./test_flush_dir/kv"

func prepareMemDB() memdb.MemoryDatabase {
	db, _ := memdb.NewMemoryDatabase(memdb.MemoryDatabaseCfg{TempPath: tmpBufDir})
	release := db.WithLock()
	for i := 0; i < 3200000; i++ {
		point := &memdb.MetricPoint{
			MetricID:  uint32(i % 1000),
			SeriesID:  uint32(i % 10000),
			SlotIndex: uint16(i % 360),
			FieldIDs:  []field.ID{1},
			Proto: &protoMetricsV1.Metric{
				Name:      "test",
				Namespace: "ns",
				SimpleFields: []*protoMetricsV1.SimpleField{
					{Name: "f1", Type: protoMetricsV1.SimpleFieldType_DELTA_SUM, Value: 10},
				},
			},
		}
		_ = db.WriteWithoutLock(point)
	}
	release()
	return db
}

func Benchmark_Memdb_Flush(b *testing.B) {
	db := prepareMemDB()
	defer func() {
		_ = fileutil.RemoveDir(tmpBufDir)
		_ = fileutil.RemoveDir(tmpStoreDir)
	}()
	b.Log("memdb size", db.MemSize())

	option := kv.DefaultStoreOption(tmpStoreDir)
	s, err := kv.NewStore("data", option)
	if err != nil {
		b.Fatal(err)
	}
	f, err := s.CreateFamily("f", kv.FamilyOption{
		CompactThreshold: 0,
		Merger:           string(metricsdata.MetricDataMerger)})
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		dataFlusher, _ := metricsdata.NewFlusher(f.NewFlusher())
		_ = db.FlushFamilyTo(dataFlusher)
	}
	b.ReportAllocs()
}
