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

package operator

import (
	"fmt"
	"strings"

	"github.com/lindb/roaring"

	"github.com/lindb/lindb/aggregation"
	"github.com/lindb/lindb/flow"
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/encoding"
	"github.com/lindb/lindb/pkg/timeutil"
)

// dataLoad represents load data operator by grouping context.
type dataLoad struct {
	executeCtx *flow.DataLoadContext
	segmentRS  *flow.TimeSegmentResultSet
	rs         flow.FilterResultSet

	foundSeries uint64
}

// NewDataLoad creates a dataLoad instance.
func NewDataLoad(executeCtx *flow.DataLoadContext,
	segmentRS *flow.TimeSegmentResultSet, rs flow.FilterResultSet,
) Operator {
	return &dataLoad{
		executeCtx: executeCtx,
		segmentRS:  segmentRS,
		rs:         rs,
	}
}

// Execute executes data load by low container(low series ids) by grouping context.
func (op *dataLoad) Execute() error {
	defer op.executeCtx.PendingDataLoadTasks.Dec()

	seriesIDs := op.executeCtx.ShardExecuteCtx.SeriesIDsAfterFiltering // after group result
	// double filtering, maybe some series ids be filtered out when do grouping.
	// filter logic: forward_reader.go -> GetGroupingScanner
	if roaring.FastAnd(seriesIDs, op.rs.SeriesIDs()).IsEmpty() {
		return nil
	}
	loader := op.rs.Load(op.executeCtx)
	if loader == nil {
		// maybe return nil loader
		return nil
	}

	familyTime := op.segmentRS.FamilyTime
	targetSlotRange := op.segmentRS.TargetRange
	queryIntervalRatio := op.segmentRS.IntervalRatio
	baseSlot := op.segmentRS.BaseTime

	// load field series data by series ids
	op.executeCtx.Decoder = encoding.GetTSDDecoder()
	op.executeCtx.DownSampling = func(slotRange timeutil.SlotRange, lowSeriesIdx uint16, fieldIdx int, getter encoding.TSDValueGetter) {
		var agg aggregation.FieldAggregator
		seriesAggregator := op.executeCtx.GetSeriesAggregator(lowSeriesIdx, fieldIdx)

		var ok bool
		agg, ok = seriesAggregator.GetAggregator(familyTime)
		if !ok {
			return
		}
		op.foundSeries++
		aggregation.DownSampling(
			slotRange, targetSlotRange, queryIntervalRatio, baseSlot,
			getter,
			agg.AggregateBySlot,
		)
	}

	// loads the metric data by given series id from load result.
	// if found data need to do down sampling aggregate.
	loader.Load(op.executeCtx)
	// release tsd decoder back to pool for re-use.
	encoding.ReleaseTSDDecoder(op.executeCtx.Decoder)
	return nil
}

// Identifier returns identifier value of data load operator.
func (op *dataLoad) Identifier() string {
	identifiers := strings.Split(op.rs.Identifier(), "segment")
	var identifier string
	if len(identifiers) > 1 {
		identifier = identifiers[1]
	} else {
		identifier = identifiers[0]
	}
	return fmt.Sprintf("Data Load[%s]", identifier)
}

// Stats returns the stats of data load operator.
func (op *dataLoad) Stats() interface{} {
	return &models.SeriesStats{
		NumOfSeries: op.foundSeries,
	}
}
