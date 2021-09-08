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
	"sort"

	"github.com/lindb/lindb/pkg/timeutil"
	"github.com/lindb/lindb/series/field"

	flatbuffers "github.com/google/flatbuffers/go"
)

// StorageRow represents a metric row with meta information and fields.
type StorageRow struct {
	MetricID  uint32
	SeriesID  uint32
	SlotIndex uint16
	FieldIDs  []field.ID

	Writable bool // Writable symbols if all meta information is set
	readOnlyRow
}

// Unmarshal unmarshalls bytes slice into a metric-row without metric context
func (mr *StorageRow) Unmarshal(data []byte) {
	mr.m.Init(data, flatbuffers.GetUOffsetT(data))
	mr.MetricID = 0
	mr.SeriesID = 0
	mr.SlotIndex = 0
	mr.FieldIDs = mr.FieldIDs[:0]
	mr.Writable = false
}

// BatchRows holds multi rows for inserting into memdb
// It is reused in sync.Pool
type BatchRows struct {
	appendIndex    int
	rows           []StorageRow
	familyIterator StorageRowFamilyIterator
}

// NewBatchRows returns write-context for batch writing.
func NewBatchRows() (ctx *BatchRows) {
	return &BatchRows{}
}
func (br *BatchRows) reset() { br.appendIndex = 0 }

func (br *BatchRows) UnmarshalRows(rowsBlock []byte) {
	br.reset()
	// uint32 length + block encoding
	for len(rowsBlock) > 0 {
		size := flatbuffers.GetSizePrefix(rowsBlock, 0)
		br.append(rowsBlock[flatbuffers.SizeUOffsetT : flatbuffers.SizeUOffsetT+size])
		rowsBlock = rowsBlock[flatbuffers.SizeUOffsetT+size:]
	}
}

func (br *BatchRows) append(data []byte) {
	defer func() { br.appendIndex++ }()
	if br.appendIndex < len(br.rows) {
		br.rows[br.appendIndex].Unmarshal(data)
		return
	}
	var sr StorageRow
	sr.Unmarshal(data)
	br.rows = append(br.rows, sr)
}

func (br *BatchRows) Len() int           { return br.appendIndex }
func (br *BatchRows) Less(i, j int) bool { return br.rows[i].Timestamp() < br.rows[j].Timestamp() }
func (br *BatchRows) Swap(i, j int)      { br.rows[i], br.rows[j] = br.rows[j], br.rows[i] }
func (br *BatchRows) Rows() []StorageRow { return br.rows[:br.Len()] }

// NewFamilyIterator provides a method for iterating data with family
func (br *BatchRows) NewFamilyIterator(interval timeutil.Interval) *StorageRowFamilyIterator {
	br.familyIterator.batch = br
	br.familyIterator.Reset(interval)
	return &br.familyIterator
}

type StorageRowFamilyIterator struct {
	groupEnd        int   // group end index
	groupStart      int   // group start index
	groupFamilyTime int64 // group family time

	sameFamily bool

	batch        *BatchRows
	intervalCalc timeutil.IntervalCalculator
}

func (itr *StorageRowFamilyIterator) HasNextFamily() bool {
	if itr.groupEnd >= itr.batch.Len() || itr.groupStart > itr.groupEnd {
		return false
	}
	if itr.sameFamily {
		itr.groupEnd = itr.batch.Len()
		itr.groupStart = 0
		return true
	}

	firstTimestamp := itr.batch.rows[itr.groupEnd].Timestamp()
	timeRange := itr.timeRangeOfTimestamp(firstTimestamp)
	itr.groupStart = itr.groupEnd
	itr.groupFamilyTime = itr.familyTimeOfTimestamp(firstTimestamp)

	for itr.groupEnd < itr.batch.Len() {
		if !timeRange.Contains(itr.batch.rows[itr.groupEnd].Timestamp()) {
			break
		}
		itr.groupEnd++
	}
	return itr.groupStart < itr.groupEnd
}

func (itr *StorageRowFamilyIterator) NextFamily() (familyTime int64, rows []StorageRow) {
	return itr.groupFamilyTime, itr.batch.rows[itr.groupStart:itr.groupEnd]
}

func (itr *StorageRowFamilyIterator) Reset(interval timeutil.Interval) {
	itr.intervalCalc = interval.Calculator()
	itr.groupStart = 0
	itr.groupEnd = 0
	itr.groupFamilyTime = 0
	// fast path, all rows are same family
	if itr.sameFamily = itr.isSameFamily(); itr.sameFamily {
		return
	}
	// slow path, re-sort it with family time
	if !sort.IsSorted(itr.batch) {
		sort.Sort(itr.batch)
	}
}

func (itr *StorageRowFamilyIterator) isSameFamily() bool {
	if itr.batch.appendIndex <= 0 {
		return true
	}
	firstTimestamp := itr.batch.rows[0].Timestamp()
	itr.groupFamilyTime = itr.familyTimeOfTimestamp(firstTimestamp)
	timeRange := itr.timeRangeOfTimestamp(firstTimestamp)
	for i := 1; i < len(itr.batch.rows); i++ {
		if !timeRange.Contains(itr.batch.rows[i].Timestamp()) {
			return false
		}
	}
	return true
}

func (itr *StorageRowFamilyIterator) familyTimeOfTimestamp(timestamp int64) int64 {
	segmentTime := itr.intervalCalc.CalcSegmentTime(timestamp)
	family := itr.intervalCalc.CalcFamily(timestamp, segmentTime)
	return itr.intervalCalc.CalcFamilyStartTime(segmentTime, family)
}

func (itr *StorageRowFamilyIterator) timeRangeOfTimestamp(timestamp int64) timeutil.TimeRange {
	segmentTime := itr.intervalCalc.CalcSegmentTime(timestamp)
	family := itr.intervalCalc.CalcFamily(timestamp, segmentTime)

	familyStartTime := itr.intervalCalc.CalcFamilyStartTime(segmentTime, family)
	return timeutil.TimeRange{
		Start: familyStartTime,
		End:   itr.intervalCalc.CalcFamilyEndTime(familyStartTime),
	}
}
