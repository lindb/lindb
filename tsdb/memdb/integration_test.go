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

package memdb

import (
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"path/filepath"
	"runtime"
	"sync"
	"testing"

	"github.com/lindb/lindb/pkg/timeutil"
	protoMetricsV1 "github.com/lindb/lindb/proto/gen/v1/metrics"
	"github.com/lindb/lindb/series/field"
)

func BenchmarkMemoryDatabase_write(b *testing.B) {
	bufferMgr := NewBufferManager(filepath.Join(b.TempDir(), "data_temp"))
	cfg := MemoryDatabaseCfg{
		BufferMgr: bufferMgr,
	}
	db, err := NewMemoryDatabase(cfg)
	if err != nil {
		b.Fatal(err)
	}
	now := timeutil.Now()

	go func() {
		_ = http.ListenAndServe("0.0.0.0:6060", nil)
	}()
	// batch write
	release := db.WithLock()

	row := protoToStorageRow(&protoMetricsV1.Metric{
		Name:      "test",
		Namespace: "ns",
		SimpleFields: []*protoMetricsV1.SimpleField{
			{Name: "f1", Type: protoMetricsV1.SimpleFieldType_DELTA_SUM, Value: 10},
		},
	})
	row.MetricID = 1

	for i := 0; i < 3200000; i++ {
		row.MetricID = 1
		row.SeriesID = uint32(i)
		row.FieldIDs = []field.ID{1}
		_ = db.WriteRow(row)
	}
	release()

	runtime.GC()
	fmt.Printf("cost:=%d\n", timeutil.Now()-now)
	now = timeutil.Now()

	row = protoToStorageRow(&protoMetricsV1.Metric{
		Name:      "test",
		Namespace: "ns",
		SimpleFields: []*protoMetricsV1.SimpleField{
			{Name: "f1", Type: protoMetricsV1.SimpleFieldType_DELTA_SUM, Value: 10},
		},
	})

	for i := 0; i < 3200000; i++ {
		row.MetricID = 1
		row.SeriesID = uint32(i)
		row.SlotIndex = uint16(i % 1024)
		row.FieldIDs = []field.ID{1}
		release := db.WithLock()
		_ = db.WriteRow(row)
		release()
	}
	runtime.GC()
	fmt.Printf("cost:=%d\n", timeutil.Now()-now)
	select {}
}

func BenchmarkMemoryDatabase_write_sum(b *testing.B) {
	run := func(n int) {
		bufferMgr := NewBufferManager(filepath.Join(b.TempDir(), "data_temp", fmt.Sprintf("%d", n)))
		var cfg = MemoryDatabaseCfg{
			BufferMgr: bufferMgr,
		}
		db, err := NewMemoryDatabase(cfg)
		if err != nil {
			b.Fatal(err)
		}
		now := timeutil.Now()

		row := protoToStorageRow(&protoMetricsV1.Metric{
			Name:      "test",
			Namespace: "ns",
			SimpleFields: []*protoMetricsV1.SimpleField{
				{Name: "f1", Type: protoMetricsV1.SimpleFieldType_DELTA_SUM, Value: 10},
			},
		})
		for i := 0; i < 3200000; i++ {
			row.MetricID = 1
			row.SeriesID = uint32(i)
			row.SlotIndex = uint16(i % 1024)
			row.FieldIDs = []field.ID{1}

			_ = db.WriteRow(row)
		}
		fmt.Printf("n:=%d, cost:=%d\n", n, timeutil.Now()-now)
	}
	now := timeutil.Now()
	var wait sync.WaitGroup
	n := 4
	wait.Add(n)
	go func() {
		run(0)
		wait.Done()
	}()
	go func() {
		run(1)
		wait.Done()
	}()
	go func() {
		run(2)
		wait.Done()
	}()
	go func() {
		run(3)
		wait.Done()
	}()
	wait.Wait()
	fmt.Println(timeutil.Now() - now)
	run(0)
}
