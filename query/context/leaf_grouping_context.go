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

package context

import (
	"encoding/binary"
	"sync"

	"go.uber.org/atomic"

	"github.com/lindb/lindb/series/tag"
)

const (
	tagValueNotFound = "tag_value_not_found"
)

// LeafGroupingContext represents collect grouping tags context under lead node.
type LeafGroupingContext struct {
	leafExecuteCtx *LeafExecuteContext

	collectGroupingTagsCompleted chan struct{}       // collect completed signal
	groupingRelatedTasks         atomic.Int32        // track how many tasks are pending
	collectRelatedTasks          atomic.Int32        // track if collect grouping tag value tasks completed
	tagsMap                      map[string]string   // tag value ids => tag values
	tagValuesMap                 []map[uint32]string // tag value id=> tag value for each group by tag key
	tagValues                    []string

	mutex sync.Mutex
}

// NewLeafGroupingContext creates a LeafGroupingContext instance.
func NewLeafGroupingContext(leafExecuteCtx *LeafExecuteContext) *LeafGroupingContext {
	// if not grouping, create empty context just for check grouping related task completed.
	ctx := &LeafGroupingContext{
		leafExecuteCtx: leafExecuteCtx,
	}
	storageExecuteCtx := leafExecuteCtx.StorageExecuteCtx
	groupByKenLen := len(storageExecuteCtx.Query.GroupBy)
	if groupByKenLen > 0 {
		ctx.tagValuesMap = make([]map[uint32]string, groupByKenLen)
		ctx.tagsMap = make(map[string]string)
		ctx.tagValues = make([]string, groupByKenLen) // temp cache
		ctx.collectGroupingTagsCompleted = make(chan struct{})
		ctx.collectRelatedTasks = *atomic.NewInt32(int32(groupByKenLen))
	}
	return ctx
}

// ForkGroupingTask forks a grouping task.
func (ctx *LeafGroupingContext) ForkGroupingTask() {
	ctx.groupingRelatedTasks.Inc()
}

// CompleteGroupingTask completes a grouping task, if all grouping tasks are completed, do collect grouping tag values.
func (ctx *LeafGroupingContext) CompleteGroupingTask() {
	ctx.groupingRelatedTasks.Dec()

	ctx.collectGroupByTagValues()
}

// collectGroupByTagValues collects group tag values
func (ctx *LeafGroupingContext) collectGroupByTagValues() {
	storageExecuteCtx := ctx.leafExecuteCtx.StorageExecuteCtx
	if ctx.groupingRelatedTasks.Load() != 0 || !storageExecuteCtx.Query.HasGroupBy() {
		// task not completed or isn't grouping query, return it.
		return
	}
	// all shard pending query tasks and grouping task completed, start collect tag values
	metadata := ctx.leafExecuteCtx.Database.Metadata()
	tagMetadata := metadata.TagMetadata()
	for idx, tagKeyID := range storageExecuteCtx.GroupByTags {
		tagKey := tagKeyID
		tagValueIDs := storageExecuteCtx.GroupingTagValueIDs[idx]
		tagIndex := idx
		if tagValueIDs == nil || tagValueIDs.IsEmpty() {
			ctx.reduceTagValues(tagIndex, nil)
			continue
		}

		tagValues := make(map[uint32]string) // tag value id => tag value
		err := tagMetadata.CollectTagValues(tagKey.ID, tagValueIDs, tagValues)
		if err != nil {
			ctx.leafExecuteCtx.SendResponse(err)
			return
		}
		ctx.reduceTagValues(tagIndex, tagValues)
	}
}

// reduceTagValues reduces the group by tag values
func (ctx *LeafGroupingContext) reduceTagValues(tagKeyIndex int, tagValues map[uint32]string) {
	ctx.mutex.Lock()
	defer ctx.mutex.Unlock()

	ctx.tagValuesMap[tagKeyIndex] = tagValues
	if ctx.collectRelatedTasks.Dec() == 0 {
		// notify all collect tag value tasks completed
		close(ctx.collectGroupingTagsCompleted)
	}
}

// getTagValues returns grouping tag string values by tag value ids.
func (ctx *LeafGroupingContext) getTagValues(tagValueIDs string) string {
	ctx.mutex.Lock()
	defer ctx.mutex.Unlock()

	if tagValues, ok := ctx.tagsMap[tagValueIDs]; ok {
		return tagValues
	}
	tagsData := []byte(tagValueIDs)
	for idx := range ctx.tagValues {
		tagValuesForKey := ctx.tagValuesMap[idx]
		offset := idx * 4
		tagValueID := binary.LittleEndian.Uint32(tagsData[offset:])
		if tagValue, ok := tagValuesForKey[tagValueID]; ok {
			ctx.tagValues[idx] = tagValue
		} else {
			ctx.tagValues[idx] = tagValueNotFound
		}
	}
	tagsOfStr := tag.ConcatTagValues(ctx.tagValues)
	ctx.tagsMap[tagValueIDs] = tagsOfStr
	return tagsOfStr
}
