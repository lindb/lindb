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

package buffer

import (
	"fmt"

	"github.com/apache/arrow-go/v18/arrow"
	"github.com/apache/arrow-go/v18/arrow/array"
	"github.com/apache/arrow-go/v18/arrow/memory"
)

// ResultSetBuild accumulates Arrow RecordBatches produced by the execution pipeline
// and merges them into a single RecordBatch when the query completes.
//
// Ownership model:
//   - Each record sent via AddRecord is retained by the buffer (Retain called on receipt).
//   - ResultSet() releases all buffered records after merging and transfers ownership of
//     the merged record to the caller. The caller is responsible for calling Release().
type ResultSetBuild struct {
	inbound   chan arrow.RecordBatch
	completed chan struct{}
	buffer    []arrow.RecordBatch
}

func CreateResultSetBuild() *ResultSetBuild {
	return &ResultSetBuild{
		inbound:   make(chan arrow.RecordBatch),
		completed: make(chan struct{}),
	}
}

func (rsb *ResultSetBuild) AddRecord(record arrow.RecordBatch) {
	if record != nil && record.NumRows() > 0 {
		rsb.inbound <- record
	}
}

func (rsb *ResultSetBuild) Process() {
	defer close(rsb.completed)
	for record := range rsb.inbound {
		// Retain so the buffer holds its own reference independent of the sender.
		record.Retain()
		rsb.buffer = append(rsb.buffer, record)
	}
}

func (rsb *ResultSetBuild) Complete() {
	close(rsb.inbound)
}

// ResultSet waits for processing to finish and returns a single merged RecordBatch.
// The caller owns the returned record and must call Release() on it.
// Returns nil if no records were received.
func (rsb *ResultSetBuild) ResultSet() arrow.RecordBatch {
	// Wait for Process() to drain the inbound channel.
	<-rsb.completed

	records := rsb.buffer
	// Release all buffered records once we are done with them.
	defer func() {
		for _, rec := range records {
			rec.Release()
		}
	}()

	switch len(records) {
	case 0:
		return nil
	case 1:
		// Single record: retain once more so the caller gets a clean ref-count,
		// then the deferred release above drops the buffer's ref.
		records[0].Retain()
		return records[0]
	}

	schema := records[0].Schema()
	numCols := schema.NumFields()
	numRows := int64(0)

	// Validate schema consistency and total row count.
	for _, rec := range records {
		if !rec.Schema().Equal(schema) {
			panic("schema mismatch during concatenation")
		}
		numRows += rec.NumRows()
	}

	// Concatenate each column across all records.
	newCols := make([]arrow.Array, numCols)
	for i := range numCols {
		colArrays := make([]arrow.Array, len(records))
		for j, rec := range records {
			colArrays[j] = rec.Column(i)
		}
		var err error
		newCols[i], err = array.Concatenate(colArrays, memory.DefaultAllocator)
		if err != nil {
			// Release columns already allocated before panicking.
			for k := range i {
				newCols[k].Release()
			}
			panic(fmt.Sprintf("failed to concatenate column %d: %v", i, err))
		}
	}

	// array.NewRecord takes ownership of newCols (ref-count already incremented
	// by Concatenate), so no additional Retain/Release needed here.
	return array.NewRecordBatch(schema, newCols, numRows)
}
