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

package metric

import (
	"bytes"
	"io"
	"strconv"
	"testing"

	"github.com/lindb/common/pkg/fasttime"
	commontimeutil "github.com/lindb/common/pkg/timeutil"
	"github.com/lindb/common/proto/gen/v1/flatMetricsV1"
	commonseries "github.com/lindb/common/series"
	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/pkg/timeutil"
)

func Test_BrokerBatchRows(t *testing.T) {
	for i := 0; i < 10; i++ {
		brokerRows := NewBrokerBatchRows()
		assertBrokerBatchRows(t, brokerRows)
		brokerRows.Release()
	}
}

func assertBrokerBatchRows(t *testing.T, brokerRows *BrokerBatchRows) {
	now := fasttime.UnixMilliseconds()

	assert.Zero(t, brokerRows.Len())
	for i := 0; i < 1000; i++ {
		i := i
		assert.NoError(t, brokerRows.TryAppend(func(row *BrokerRow) error {
			buildRow(row, now-int64(i)*1000*60)
			return nil
		}))
	}
	assert.Equal(t, 1000, brokerRows.Len())

	// only one shard
	itr := brokerRows.NewShardGroupIterator(1)
	assert.True(t, itr.HasRowsForNextShard())
	var interval timeutil.Interval
	_ = interval.ValueOf("10s")
	shardIdx, familyItr := itr.FamilyRowsForNextShard(interval)
	var (
		allRows  []BrokerRow
		families int
	)

	for familyItr.HasNextFamily() {
		_, rows := familyItr.NextFamily()
		families++
		allRows = append(allRows, rows...)
	}
	assert.Equal(t, 0, shardIdx)
	assert.Len(t, allRows, 1000)
	assert.False(t, itr.HasRowsForNextShard())

	itr = brokerRows.NewShardGroupIterator(10)
	for i := 0; i < 10; i++ {
		assert.True(t, itr.HasRowsForNextShard())
		shardIdx, familyItr = itr.FamilyRowsForNextShard(interval)
		assert.Equal(t, i, shardIdx)
		assert.True(t, familyItr.HasNextFamily())
		_, rows := familyItr.NextFamily()
		assert.True(t, len(rows) > 0)
		assert.Len(t, brokerRows.Rows(), 1000)
	}
	assert.False(t, itr.HasRowsForNextShard())

	// eviction
	assert.InDelta(t, 1000,
		brokerRows.EvictOutOfTimeRange(100, 100), 100)
}

func buildRow(row *BrokerRow, timestamp int64) {
	builder, releaseFunc := commonseries.NewRowBuilder()
	defer releaseFunc(builder)

	builder.AddMetricName([]byte("test"))
	_ = builder.AddTag([]byte("ts"), []byte(strconv.FormatInt(timestamp, 10)))
	_ = builder.AddSimpleField([]byte("f1"), flatMetricsV1.SimpleFieldTypeDeltaSum, 100)
	builder.AddTimestamp(timestamp)
	data, _ := builder.Build()
	row.FromBlock(data)
}

func Test_BrokerBatchRows_AppendError(t *testing.T) {
	batch := NewBrokerBatchRows()
	defer batch.Release()

	assert.Error(t, batch.TryAppend(func(row *BrokerRow) error {
		return io.ErrShortBuffer
	}))
	assert.Equal(t, 0, batch.Len())
}

func Test_BrokerRow_Writer(t *testing.T) {
	var row BrokerRow
	row.IsOutOfTimeRange = true
	row.buffer = append(row.buffer, []byte{1, 2, 3, 4}...)

	var buf bytes.Buffer
	assert.Equal(t, 0, row.Size())
	n, err := row.WriteTo(&buf)
	assert.Equal(t, 0, n)
	assert.NoError(t, err)

	_ = row.Metric()
	row.IsOutOfTimeRange = false
	assert.Equal(t, 4, row.Size())
	n, err = row.WriteTo(&buf)
	assert.Equal(t, 4, n)
	assert.NoError(t, err)
}

func Test_BrokerBatchRows_FamilyRowsForNextShard_SingleShard(t *testing.T) {
	now := fasttime.UnixMilliseconds()

	var brokerRows BrokerBatchRows
	for i := 0; i < 30; i++ {
		_ = brokerRows.TryAppend(func(row *BrokerRow) error {
			buildRow(row, now)
			return nil
		})
	}

	for i := 30; i < 50; i++ {
		_ = brokerRows.TryAppend(func(row *BrokerRow) error {
			buildRow(row, now-commontimeutil.OneHour)
			return nil
		})
	}

	for i := 50; i < 100; i++ {
		_ = brokerRows.TryAppend(func(row *BrokerRow) error {
			buildRow(row, now+commontimeutil.OneHour)
			return nil
		})
	}
	var interval timeutil.Interval
	_ = interval.ValueOf("10s")

	// one shard
	itr := brokerRows.NewShardGroupIterator(1)
	assert.True(t, itr.HasRowsForNextShard())

	shardIdx, familyItr := itr.FamilyRowsForNextShard(interval)
	assert.Equal(t, 0, shardIdx)

	// last family
	assert.True(t, familyItr.HasNextFamily())
	familyTime, rows := familyItr.NextFamily()
	t.Log(familyTime)
	assert.Len(t, rows, 20)

	// current family
	assert.True(t, familyItr.HasNextFamily())
	familyTime, rows = familyItr.NextFamily()
	t.Log(familyTime)
	assert.Len(t, rows, 30)

	// next family
	assert.True(t, familyItr.HasNextFamily())
	familyTime, rows = familyItr.NextFamily()
	t.Log(familyTime)
	assert.Len(t, rows, 50)

	assert.False(t, familyItr.HasNextFamily())
}

func Test_BrokerBatchRows_FamilyRowsForNextShard_SameFamily(t *testing.T) {
	now := fasttime.UnixMilliseconds()

	var brokerRows BrokerBatchRows
	for i := 0; i < 30; i++ {
		_ = brokerRows.TryAppend(func(row *BrokerRow) error {
			buildRow(row, now)
			return nil
		})
	}
	itr := brokerRows.NewShardGroupIterator(1)
	var interval timeutil.Interval
	_ = interval.ValueOf("10s")
	assert.True(t, itr.HasRowsForNextShard())
	_, familyItr := itr.FamilyRowsForNextShard(interval)
	assert.True(t, familyItr.HasNextFamily())
	assert.False(t, familyItr.HasNextFamily())
}
