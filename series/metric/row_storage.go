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

// StorageBatchRows holds multi rows for inserting into memdb
// It is reused in sync.Pool
type StorageBatchRows struct {
	appendIndex int
	rows        []StorageRow
}

// NewStorageBatchRows returns write-context for batch writing.
func NewStorageBatchRows() (ctx *StorageBatchRows) {
	return &StorageBatchRows{}
}
func (br *StorageBatchRows) reset() { br.appendIndex = 0 }

func (br *StorageBatchRows) UnmarshalRows(rowsBlock []byte) {
	br.reset()
	// uint32 length + block encoding
	for len(rowsBlock) > 0 {
		size := flatbuffers.GetSizePrefix(rowsBlock, 0)
		br.append(rowsBlock[flatbuffers.SizeUOffsetT : flatbuffers.SizeUOffsetT+size])
		rowsBlock = rowsBlock[flatbuffers.SizeUOffsetT+size:]
	}
}

func (br *StorageBatchRows) append(data []byte) {
	defer func() { br.appendIndex++ }()
	if br.appendIndex < len(br.rows) {
		br.rows[br.appendIndex].Unmarshal(data)
		return
	}
	var sr StorageRow
	sr.Unmarshal(data)
	br.rows = append(br.rows, sr)
}

func (br *StorageBatchRows) Len() int { return br.appendIndex }
func (br *StorageBatchRows) Less(i, j int) bool {
	return br.rows[i].Timestamp() < br.rows[j].Timestamp()
}
func (br *StorageBatchRows) Swap(i, j int)      { br.rows[i], br.rows[j] = br.rows[j], br.rows[i] }
func (br *StorageBatchRows) Rows() []StorageRow { return br.rows[:br.Len()] }
