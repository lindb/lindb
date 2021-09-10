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
	"io"
	"sort"
	"sync"

	flatbuffers "github.com/google/flatbuffers/go"
	"github.com/lithammer/go-jump-consistent-hash"

	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/encoding"
	"github.com/lindb/lindb/pkg/fasttime"
	"github.com/lindb/lindb/proto/gen/v1/flatMetricsV1"
)

type BrokerRow struct {
	m      flatMetricsV1.Metric
	buffer []byte

	ShardID models.ShardID
	// IsOutOfTimeRange marks if this row is out-of time-range
	// data is not accessible when its set to true
	IsOutOfTimeRange bool
}

// FromBlock resets buffer, unmarshal from a new block,
// make sure that metric and shard id will be overwritten manually
func (row *BrokerRow) FromBlock(block []byte) {
	row.buffer = encoding.MustCopy(row.buffer, block)
	size := flatbuffers.GetSizePrefix(row.buffer, 0)
	partition := row.buffer[flatbuffers.SizeUOffsetT : flatbuffers.SizeUOffsetT+size]
	row.m.Init(partition, flatbuffers.GetUOffsetT(partition))
}

func (row *BrokerRow) Metric() flatMetricsV1.Metric { return row.m }

func (row *BrokerRow) Size() int {
	if row.IsOutOfTimeRange {
		return 0
	}
	return len(row.buffer)
}

func (row *BrokerRow) WriteTo(writer io.Writer) (int, error) {
	if row.IsOutOfTimeRange {
		return 0, nil
	}
	return writer.Write(row.buffer)
}

var brokerBatchRowsPool sync.Pool

// BrokerBatchRows holds rows from ingestion
// row will be putted into buffer after validation and re-building
type BrokerBatchRows struct {
	rows     []BrokerRow
	rowCount int

	shardGroupIterator BrokerBatchShardIterator
}

func newBrokerBatchRows() *BrokerBatchRows {
	return &BrokerBatchRows{}
}

// NewBrokerBatchRows returns a new batch for decoding flat metrics.
func NewBrokerBatchRows() (batch *BrokerBatchRows) {
	item := brokerBatchRowsPool.Get()
	if item != nil {
		builder := item.(*BrokerBatchRows)
		builder.reset()
		return builder
	}
	return newBrokerBatchRows()
}

// Release releases rows context into sync.Pool
func (br *BrokerBatchRows) Release() { brokerBatchRowsPool.Put(br) }

func (br *BrokerBatchRows) reset() { br.rowCount = 0 }

func (br *BrokerBatchRows) Len() int { return br.rowCount }
func (br *BrokerBatchRows) Less(i, j int) bool {
	return br.rows[i].ShardID < br.rows[j].ShardID
}
func (br *BrokerBatchRows) Swap(i, j int)     { br.rows[i], br.rows[j] = br.rows[j], br.rows[i] }
func (br *BrokerBatchRows) Rows() []BrokerRow { return br.rows[:br.rowCount] }

// EvictOutOfTimeRange evicts and marks out-of-range metrics invalid
func (br *BrokerBatchRows) EvictOutOfTimeRange(behind, ahead int64) (evicted int) {
	// check metric timestamp if in acceptable time range
	now := fasttime.UnixMilliseconds()
	for idx := 0; idx < br.Len(); idx++ {
		if (behind > 0 && br.rows[idx].m.Timestamp() < now-behind) ||
			(ahead > 0 && br.rows[idx].m.Timestamp() > now+ahead) {
			br.rows[idx].IsOutOfTimeRange = true
			evicted++
		}
	}
	return evicted
}

func (br *BrokerBatchRows) TryAppend(appendFunc func(row *BrokerRow) error) error {
	if len(br.rows) <= br.rowCount {
		br.rows = append(br.rows, BrokerRow{})
	}
	if err := appendFunc(&br.rows[br.rowCount]); err != nil {
		return err
	}
	// decoded successfully, move to next row index
	br.rowCount++
	return nil
}

func (br *BrokerBatchRows) NewShardGroupIterator(numOfShards int32) *BrokerBatchShardIterator {
	for i := 0; i < br.Len(); i++ {
		br.rows[i].ShardID = models.ShardID(jump.Hash(br.rows[i].m.Hash(), numOfShards))
	}
	br.shardGroupIterator.batch = br
	br.shardGroupIterator.Reset()
	return &br.shardGroupIterator
}

// BrokerBatchShardIterator grouping broker rows with shard-id,
// rows will be batched inserted into shard-channel for replication
type BrokerBatchShardIterator struct {
	groupEnd     int            // group end index
	groupStart   int            // group start index
	groupShardID models.ShardID // group shard id

	batch *BrokerBatchRows
}

// Reset re-sorts batch rows for batching inserting
func (itr *BrokerBatchShardIterator) Reset() {
	sort.Sort(itr.batch)
	itr.groupStart = 0
	itr.groupEnd = 0
	itr.groupShardID = models.ShardID(-1)
}

func (itr *BrokerBatchShardIterator) HasRowsForNextShard() bool {
	if itr.groupEnd >= itr.batch.Len() || itr.groupStart > itr.groupEnd {
		return false
	}
	itr.groupShardID = itr.batch.rows[itr.groupEnd].ShardID
	itr.groupStart = itr.groupEnd

	for itr.groupEnd < itr.batch.Len() {
		if !(itr.batch.rows[itr.groupEnd].ShardID == itr.groupShardID) {
			break
		}
		itr.groupEnd++
	}
	return itr.groupStart < itr.groupEnd
}

func (itr *BrokerBatchShardIterator) RowsForNextShard() (
	shardID models.ShardID,
	rows []BrokerRow,
) {
	return itr.groupShardID, itr.batch.rows[itr.groupStart:itr.groupEnd]
}
