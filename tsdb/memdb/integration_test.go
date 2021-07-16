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
	"context"
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"path/filepath"
	"runtime"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/kv"
	"github.com/lindb/lindb/pkg/fileutil"
	"github.com/lindb/lindb/pkg/timeutil"
	protoMetricsV1 "github.com/lindb/lindb/proto/gen/v1/metrics"
	"github.com/lindb/lindb/tsdb/metadb"
)

func BenchmarkMemoryDatabase_write(b *testing.B) {
	defer func() {
		_ = fileutil.RemoveDir(testPath)
	}()
	kvStore, err := kv.NewStore("test", kv.StoreOption{
		Path:                 filepath.Join(testPath, "kv"),
		Levels:               2,
		CompactCheckInterval: 5 * 60,
	})
	assert.NoError(b, err)

	metadata, err := metadb.NewMetadata(context.TODO(), "test", filepath.Join(testPath, "meta"),
		kvStore.GetFamily("meta"))
	assert.NoError(b, err)

	metricID, err := metadata.MetadataDatabase().GenMetricID("ns", "test")
	assert.NoError(b, err)
	var cfg = MemoryDatabaseCfg{
		Metadata: metadata,
		TempPath: filepath.Join(testPath, "data_temp"),
	}
	db, err := NewMemoryDatabase(cfg)
	if err != nil {
		b.Fatal(err)
	}
	now := timeutil.Now()
	for i := 0; i < 3200000; i++ {
		_ = db.Write("ns", "test", metricID, uint32(i), uint16(now%1024), []*protoMetricsV1.SimpleField{{
			Name:  "f1",
			Type:  protoMetricsV1.SimpleFieldType_DELTA_SUM,
			Value: 10.0,
		}}, nil)
	}
	runtime.GC()
	fmt.Printf("cost:=%d\n", timeutil.Now()-now)
	now = timeutil.Now()
	for i := 0; i < 3200000; i++ {
		_ = db.Write("ns", "test", metricID, uint32(i), uint16(now%1024), []*protoMetricsV1.SimpleField{{
			Name:  "f1",
			Type:  protoMetricsV1.SimpleFieldType_DELTA_SUM,
			Value: 10.0,
		}}, nil)
	}
	runtime.GC()
	fmt.Printf("cost:=%d\n", timeutil.Now()-now)
	_ = http.ListenAndServe("0.0.0.0:6060", nil)
}

func BenchmarkMemoryDatabase_write_sum(b *testing.B) {
	defer func() {
		_ = fileutil.RemoveDir(testPath)
	}()
	kvStore, err := kv.NewStore("test", kv.StoreOption{
		Path:                 filepath.Join(testPath, "kv"),
		Levels:               2,
		CompactCheckInterval: 5 * 60,
	})
	assert.NoError(b, err)

	metadata, err := metadb.NewMetadata(context.TODO(), "test",
		filepath.Join(testPath, "meta"), kvStore.GetFamily("meta"))
	assert.NoError(b, err)

	metricID, err := metadata.MetadataDatabase().GenMetricID("ns", "test")
	assert.NoError(b, err)
	run := func(n int) {
		var cfg = MemoryDatabaseCfg{
			Metadata: metadata,
			TempPath: filepath.Join(testPath, "data_temp", fmt.Sprintf("%d", n)),
		}
		db, err := NewMemoryDatabase(cfg)
		if err != nil {
			b.Fatal(err)
		}
		now := timeutil.Now()
		for i := 0; i < 400000; i++ {
			_ = db.Write("ns", "test", metricID, uint32(i), uint16(now%1024), []*protoMetricsV1.SimpleField{{
				Name:  "f1",
				Type:  protoMetricsV1.SimpleFieldType_DELTA_SUM,
				Value: 10.0,
			}}, nil)
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
