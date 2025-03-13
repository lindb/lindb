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
	"fmt"
	"slices"

	"github.com/samber/lo"

	"github.com/lindb/lindb/pkg/collections"
	"github.com/lindb/lindb/pkg/encoding"
	"github.com/lindb/lindb/series/field"
)

type dataScan struct {
	split *DataSplit

	reduceCh chan<- any
}

func NewDataScan(split *DataSplit, reduceCh chan<- any) *dataScan {
	return &dataScan{
		split:    split,
		reduceCh: reduceCh,
	}
}

func (ds *dataScan) Run() {
	// find time series data of filed
	tableScan := ds.split.partition.tableScan
	familyLoaders := ds.buildFamilyLoaders(ds.split)
	if len(familyLoaders) == 0 {
		// family not match
		return
	}

	isGrouping := tableScan.isGrouping()
	var (
		groupingAgg grouping
		tagsScanner *TagsScanner
	)

	if isGrouping {
		tagScanners := ds.split.groupingContext.BuildGroup(ds.split.seriesIDHighKey, ds.split.lowSeriesIDs)
		tagsScanner = NewTagsScanner(tagScanners)
		groupingAgg = newGroupingWithTags(tagsScanner, tableScan)
	} else {
		groupingAgg = newGroupingWithoutTags(tableScan)
	}

	ds.split.groupingAgg = groupingAgg
	fmt.Println("run data scan")

	columnRollups := tableScan.createRollups()
	var numOfPoints int
	step := int64(0)
	start := tableScan.timeRange.Start
	if tableScan.isTimestampSelected {
		numOfPoints = tableScan.timeRange.NumOfPoints(tableScan.interval)
		step = tableScan.interval.Int64()
	} else {
		numOfPoints = 1 // timestamp not in select item list
	}

	it := ds.split.lowSeriesIDs.PeekableIterator()
	// loop each low series ids
	for it.HasNext() {
		lowSeriesID := it.Next()
		aggregator := groupingAgg.GetAggregator(lowSeriesID)

		// load all fields data from families
		for _, loader := range familyLoaders {
			// fmt.Printf("start load family time=%v\n", loader.familyTime)
			// TODO: go? reset stream
			for _, family := range loader.families {
				slotRange := family.filterResultSet.SlotRange()

				family.loader.Load(lowSeriesID, func(field field.Meta, getter encoding.TSDValueGetter) {
					columnStream := loader.streams.GetStreamByIndex(field.Index)
					fn := field.Type.AggType().Aggregate
					for movingSourceSlot := slotRange.Start; movingSourceSlot <= slotRange.End; movingSourceSlot++ {
						value, ok := getter.GetValue(movingSourceSlot)
						if !ok {
							// no data, goto next loop
							continue
						}
						columnStream.SetAtStep(int(movingSourceSlot), value, fn)
					}
				})
			}
		}

		// do rollup(down sampling)/aggregation for each column
		for index, column := range tableScan.columns {
			field := column.meta

			// do rollup(down sampling)
			for _, loader := range familyLoaders {
				familyTime := loader.familyTime
				slotRange := loader.timeRange
				interval := loader.interval.Int64()
				for movingSourceSlot := slotRange.Start; movingSourceSlot <= slotRange.End; movingSourceSlot++ {
					timestamp := familyTime + int64(movingSourceSlot)*interval
					value := loader.streams.GetStreamByIndex(field.Index).GetAtStep(int(movingSourceSlot))

					// rollup
					for idx, r := range column.rollups {
						columnRollups[idx].doRollup(r.aggType, timestamp, value)
					}
				}
			}

			// do aggregation
			for aggIdx, agg := range column.aggs {
				idx := index + aggIdx
				if aggregator[idx] == nil {
					aggregator[idx] = collections.NewFloatArray(numOfPoints)
				}
				// aggregate down sampled data
				aggregate(agg.aggType, start, step, aggregator[idx], columnRollups[agg.target].getTimeSeries())
			}

			// reset rollup context for next column
			for idx := range column.rollups {
				columnRollups[idx].reset()
			}
		}
	}

	if isGrouping {
		// emit tag value ids, reduce task need lookup tag value by id.
		ds.reduceCh <- tagsScanner.GetTagValueIDs()
	}

	ds.reduceCh <- ds.split
}

func (ds *dataScan) buildFamilyLoaders(split *DataSplit) []*loader {
	tableScan := split.partition.tableScan
	familyLoaderMap := make(map[int64]*loader)
	for i := range split.partition.fieldsData {
		rs := split.partition.fieldsData[i]
		// check series ids if match
		dataLoader := rs.Load(split.seriesIDHighKey, split.lowSeriesIDs)
		if dataLoader != nil {
			family, ok := familyLoaderMap[rs.FamilyTime()]
			if ok {
				family.families = append(family.families, &familyLoader{filterResultSet: rs, loader: dataLoader})
				family.timeRange = family.timeRange.Union(rs.SlotRange())
			} else {
				familyLoaderMap[rs.FamilyTime()] = &loader{
					familyTime: rs.FamilyTime(),
					interval:   rs.Interval(),
					families:   []*familyLoader{{filterResultSet: rs, loader: dataLoader}},
					timeRange:  rs.SlotRange(),
					streams:    newStreams(tableScan.fields.Len(), tableScan.db.GetStream),
				}
			}
		}
	}
	if len(familyLoaderMap) == 0 {
		// family not match
		return nil
	}

	familyLoaders := lo.Values(familyLoaderMap)
	// order by family time
	slices.SortFunc(familyLoaders, func(a, b *loader) int {
		return int(a.familyTime - b.familyTime)
	})
	return familyLoaders
}
