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

//go:build integration
// +build integration

package tsdb

import (
	"bytes"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	commontimeutil "github.com/lindb/common/pkg/timeutil"
	protoMetricsV1 "github.com/lindb/common/proto/gen/v1/linmetrics"

	"github.com/lindb/lindb/config"
	"github.com/lindb/lindb/kv"
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/option"
	"github.com/lindb/lindb/pkg/timeutil"
	"github.com/lindb/lindb/series/metric"
	"github.com/lindb/lindb/tsdb"
)

func TestDatabase_Write_And_Rollup(t *testing.T) {
	dir := t.TempDir()
	config.SetGlobalStorageConfig(&config.StorageBase{
		TSDB: config.TSDB{Dir: dir},
	})

	engine, err := tsdb.NewEngine()
	assert.NoError(t, err)
	assert.NotNil(t, engine)
	defer func() {
		engine.Close()
	}()

	interval := timeutil.Interval(10 * 1000)
	rollupInterval := timeutil.Interval(5 * 60 * 1000)
	opt := &option.DatabaseOption{
		Intervals:    option.Intervals{{Interval: interval}, {Interval: rollupInterval}},
		AutoCreateNS: true,
	}
	err = engine.CreateShards("write-db", opt, models.ShardID(1))
	assert.NoError(t, err)
	db, ok := engine.GetDatabase("write-db")
	assert.True(t, ok)
	assert.NotNil(t, db)
	shard, ok := db.GetShard(models.ShardID(1))
	assert.True(t, ok)
	assert.NotNil(t, shard)

	now, _ := commontimeutil.ParseTimestamp("20190702 19:10:00", "20060102 15:04:05")
	familyTime := interval.Calculator().CalcFamilyTime(now)
	f, err := shard.GetOrCrateDataFamily(familyTime)
	assert.NoError(t, err)
	assert.NotNil(t, f)

	for i := 0; i < 5; i++ {
		rows := mockBatchRows(&protoMetricsV1.Metric{
			Name:      "test",
			Timestamp: now,
			SimpleFields: []*protoMetricsV1.SimpleField{{
				Name:  "f1",
				Value: 1.0,
				Type:  protoMetricsV1.SimpleFieldType_DELTA_SUM,
			}},
		})

		err = f.WriteRows(rows)
		assert.NoError(t, err)

		err = f.Flush()
		assert.NoError(t, err)
	}

	storeName := tsdb.ShardSegmentPath("write-db", models.ShardID(1), interval, "20190702")
	store, ok := kv.GetStoreManager().GetStoreByName(storeName)
	assert.True(t, ok)
	assert.NotNil(t, store)
	storeName = tsdb.ShardSegmentPath("write-db", models.ShardID(1), rollupInterval, "201907")
	rollupTargetStore, ok := kv.GetStoreManager().GetStoreByName(storeName)
	assert.True(t, ok)
	assert.NotNil(t, rollupTargetStore)

	store.ForceRollup()
	time.Sleep(200 * time.Millisecond)
}

func mockBatchRows(m *protoMetricsV1.Metric) []metric.StorageRow {
	var ml = protoMetricsV1.MetricList{Metrics: []*protoMetricsV1.Metric{m}}
	var buf bytes.Buffer
	converter := metric.NewProtoConverter(models.NewDefaultLimits())
	_, _ = converter.MarshalProtoMetricListV1To(ml, &buf)

	var br metric.StorageBatchRows
	br.UnmarshalRows(buf.Bytes())
	return br.Rows()
}
