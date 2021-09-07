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
	"sync"

	flatbuffers "github.com/google/flatbuffers/go"

	"github.com/lindb/lindb/series/field"
)

// StorageRow represents a metric row with meta information and fields.
type StorageRow struct {
	MetricID  uint32
	SeriesID  uint32
	SlotIndex uint16
	FieldIDs  []field.ID

	readOnlyRow
}

// Unmarshal unmarshalls bytes slice into a metric-row without metric context
func (mr *StorageRow) Unmarshal(data []byte) {
	mr.m.Init(data, flatbuffers.GetUOffsetT(data))
	mr.MetricID = 0
	mr.SeriesID = 0
	mr.SlotIndex = 0
	mr.FieldIDs = mr.FieldIDs[:0]
}

// WriteCtx holds multi rows for inserting into memdb
// It is reused in sync.Pool
type WriteCtx struct {
	appendIndex int
	rows        []StorageRow
	itr         WriteCtxIterator
}

var writeCtxPool sync.Pool

// NewWriteCtx returns a fixed size context for batch writing.
func NewWriteCtx(size int) (ctx *WriteCtx, releaseFunc func(ctx *WriteCtx)) {
	releaseFunc = func(_ctx *WriteCtx) {
		_ctx.appendIndex = 0
		writeCtxPool.Put(_ctx)
	}
	item := writeCtxPool.Get()
	if item != nil {
		ctx = item.(*WriteCtx)
		if cap(ctx.rows) >= size {
			ctx.rows = ctx.rows[0:size]
			return ctx, releaseFunc
		}
	}
	return &WriteCtx{rows: make([]StorageRow, size)}, releaseFunc
}

func (ctx *WriteCtx) Append(data []byte) bool {
	if ctx.appendIndex >= len(ctx.rows) {
		return false
	}
	ctx.rows[ctx.appendIndex].Unmarshal(data)
	return true
}

func (ctx *WriteCtx) NewIterator() *WriteCtxIterator {
	ctx.itr.ctx = ctx
	return &ctx.itr
}

type WriteCtxIterator struct {
	idx int
	ctx *WriteCtx
}

func (itr *WriteCtxIterator) HasNext() bool {
	itr.idx++
	return itr.idx < len(itr.ctx.rows)
}

func (itr *WriteCtxIterator) NextRow() *StorageRow { return &itr.ctx.rows[itr.idx] }
func (itr *WriteCtxIterator) Reset()               { itr.idx = -1 }
