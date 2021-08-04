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

	"github.com/lindb/lindb/pkg/fileutil"
	"github.com/lindb/lindb/pkg/timeutil"
	protoMetricsV1 "github.com/lindb/lindb/proto/gen/v1/metrics"
	"github.com/lindb/lindb/series/field"
)

func BenchmarkMemoryDatabase_write(b *testing.B) {
	defer func() {
		_ = fileutil.RemoveDir(testPath)
	}()
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
	for i := 0; i < 3200000; i++ {
		point := &MetricPoint{
			MetricID:  1,
			SeriesID:  uint32(i),
			SlotIndex: uint16(now % 1024),
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

	runtime.GC()
	fmt.Printf("cost:=%d\n", timeutil.Now()-now)
	now = timeutil.Now()

	for i := 0; i < 3200000; i++ {
		point := &MetricPoint{
			MetricID:  1,
			SeriesID:  uint32(i),
			SlotIndex: uint16(now % 1024),
			FieldIDs:  []field.ID{1},
			Proto: &protoMetricsV1.Metric{
				Name:      "test",
				Namespace: "ns",
				SimpleFields: []*protoMetricsV1.SimpleField{
					{Name: "f1", Type: protoMetricsV1.SimpleFieldType_DELTA_SUM, Value: 10},
				},
			},
		}
		_ = db.Write(point)
	}
	runtime.GC()
	fmt.Printf("cost:=%d\n", timeutil.Now()-now)
	select {}
}

func BenchmarkMemoryDatabase_write_sum(b *testing.B) {
	defer func() {
		_ = fileutil.RemoveDir(testPath)
	}()

	run := func(n int) {
		var cfg = MemoryDatabaseCfg{
			TempPath: filepath.Join(testPath, "data_temp", fmt.Sprintf("%d", n)),
		}
		db, err := NewMemoryDatabase(cfg)
		if err != nil {
			b.Fatal(err)
		}
		now := timeutil.Now()
		for i := 0; i < 3200000; i++ {
			point := &MetricPoint{
				MetricID:  1,
				SeriesID:  uint32(i),
				SlotIndex: uint16(now % 1024),
				FieldIDs:  []field.ID{1},
				Proto: &protoMetricsV1.Metric{
					Name:      "test",
					Namespace: "ns",
					SimpleFields: []*protoMetricsV1.SimpleField{
						{Name: "f1", Type: protoMetricsV1.SimpleFieldType_DELTA_SUM, Value: 10},
					},
				},
			}
			_ = db.Write(point)
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
