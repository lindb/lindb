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
	"github.com/lindb/common/models"
	"github.com/lindb/roaring"
	"github.com/samber/lo"

	"github.com/lindb/lindb/series/field"
	"github.com/lindb/lindb/spi/types"
)

type reducer struct {
	ctx       *ExecutionContext
	tableScan *TableScan
	inbound   <-chan any // *DataSplit or []*roaring.Bitmap
	outbound  chan<- *types.Page

	result map[*GroupingKey][]Result // tags => series data of fields(aggregators)
}

func NewReducer(ctx *ExecutionContext, tableScan *TableScan, inbound <-chan any, outbound chan<- *types.Page) *reducer {
	return &reducer{
		ctx:       ctx,
		tableScan: tableScan,
		inbound:   inbound,
		outbound:  outbound,
		result:    make(map[*GroupingKey][]Result),
	}
}

func (r *reducer) Run() {
	for {
		select {
		case <-r.ctx.ctx.Done():
			err := r.ctx.ctx.Err()
			if err != nil {
				panic(err)
			}
			goto END
		case event := <-r.inbound:
			if event == nil {
				goto END
			}
			switch e := event.(type) {
			case *DataSplit:
				r.process(e)
			case []*roaring.Bitmap:
				r.tableScan.grouping.CollectTagValueIDs(e)
			}
		}
	}

END:
	if r.tableScan.isGrouping() {
		r.tableScan.grouping.CollectTagValues()
	}

	r.outbound <- r.buildOutputPage()
}

func (r *reducer) process(split *DataSplit) {
	if r.tableScan.fields.Len() == 0 {
		// find series data(columns)
		r.findSeries(split)
		return
	}

	// merge the data of time series
	split.groupingAgg.ForEach(func(tags *GroupingKey, rs []Result) {
		if _, ok := r.result[tags]; !ok {
			r.result[tags] = rs
		} else {
			panic("need impl merge")
		}
	})
}

func (r *reducer) findSeries(split *DataSplit) {
	tagScanners := split.groupingContext.BuildGroup(split.seriesIDHighKey, split.lowSeriesIDs)
	tagsScanner := NewTagsScanner(tagScanners)
	it := split.lowSeriesIDs.PeekableIterator()
	for it.HasNext() {
		// loop each low series ids
		lowSeriesID := it.Next()
		// build grouping keys
		key := tagsScanner.FindTagValues(lowSeriesID)
		if _, ok := r.result[key]; !ok {
			r.result[key.Clone()] = nil
		}
	}
	tableScan := split.partition.tableScan

	tableScan.grouping.CollectTagValueIDs(tagsScanner.GetTagValueIDs())
}

func (r *reducer) buildOutputPage() *types.Page {
	page := types.NewPage()
	var (
		fields          []*types.Column
		grouping        []*types.Column
		groupingIndexes []int
	)
	for idx, output := range r.tableScan.outputs {
		column := types.NewColumn()
		page.AppendColumn(output, column)
		if lo.ContainsBy(r.tableScan.fields, func(item field.Meta) bool {
			return item.Name.String() == output.Name
		}) {
			fields = append(fields, column)
		} else if output.DataType == types.DTString {
			grouping = append(grouping, column)
			groupingIndexes = append(groupingIndexes, idx)
		}
	}
	// set grouping index of the columns
	page.SetGrouping(groupingIndexes)

	hasGrouping := r.tableScan.isGrouping()
	// set tag values
	for tags, seriesData := range r.result {
		if hasGrouping {
			tags := r.tableScan.grouping.GetTagValues(*tags)
			for idx, tag := range tags {
				grouping[idx].Append(tag)
			}
		}
		for fieldIdx, stream := range seriesData {
			if stream == nil {
				fields[fieldIdx].Append(nil)
				continue
			}
			switch dst := stream.(type) {
			case *result[float64]:
				timeSeries := types.NewTimeSeriesWithValues(r.tableScan.timeRange, r.tableScan.interval, dst.array.Values())
				fields[fieldIdx].Append(timeSeries)
			case *result[*models.Exemplar]:
				fields[fieldIdx].Append(dst.array.Values())
			}
		}
	}
	return page
}
