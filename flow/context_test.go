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

package flow

import (
	"testing"

	"github.com/lindb/roaring"
	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/sql/stmt"
)

func TestDataLoadContext_IterateLowSeriesIDs(t *testing.T) {
	querySeriesIDs := roaring.BitmapOf(5, 11, 13)
	storageSeriesIDs := roaring.BitmapOf(1, 3, 5, 7, 9, 11, 13, 15)
	ctx := &DataLoadContext{
		LowSeriesIDsContainer: querySeriesIDs.GetContainer(0),
		ShardExecuteCtx: &ShardExecuteContext{
			StorageExecuteCtx: &StorageExecuteContext{
				Query: &stmt.Query{},
			},
		},
	}
	ctx.Grouping()
	findSeriesIDs := roaring.New()
	storageLowSeriesContainer := storageSeriesIDs.GetContainer(0)
	storageLowSeriesIDs := storageLowSeriesContainer.ToArray()
	ctx.IterateLowSeriesIDs(storageLowSeriesContainer, func(seriesIdxFromQuery uint16, seriesIdxFromStorage int) {
		findSeriesIDs.Add(uint32(storageLowSeriesIDs[seriesIdxFromStorage]))
	})
	assert.Equal(t, querySeriesIDs, findSeriesIDs)
}
