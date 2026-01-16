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

	columns := tableScan.columns
	for _, column := range columns {
		column.createStream(len(familyLoaders))
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

	it := ds.split.lowSeriesIDs.PeekableIterator()
	// loop each low series ids
	for it.HasNext() {
		lowSeriesID := it.Next()
		aggregator := groupingAgg.GetAggregator(lowSeriesID)

		// load all columns data from data families
		for _, loader := range familyLoaders {
			loader.load(lowSeriesID)
		}

		// do rollup(down sampling)/aggregation for each column
		for _, column := range columns {
			// do downsampling
			column.downsampling(familyLoaders)

			// do aggregation
			column.aggregate(aggregator)

			// reset context for next series
			column.reset()
		}
	}

	if isGrouping {
		// emit tag value ids, reduce task need lookup tag value by id.
		ds.reduceCh <- tagsScanner.GetTagValueIDs()
	}

	ds.reduceCh <- ds.split
}

func (ds *dataScan) buildFamilyLoaders(split *DataSplit) []*loader {
	familyLoaderMap := make(map[int64]*loader)
	for i := range split.partition.fieldsData {
		rs := split.partition.fieldsData[i]
		// check series ids if match
		dataLoader := rs.Load(split.seriesIDHighKey, split.lowSeriesIDs)
		if dataLoader == nil {
			continue
		}

		familyTime := rs.FamilyTime()
		family, ok := familyLoaderMap[familyTime]
		fl := &familyLoader{filterResultSet: rs, loader: dataLoader, columns: split.partition.tableScan.columns}
		if ok {
			family.families = append(family.families, fl)
			family.timeRange = family.timeRange.Union(rs.SlotRange())
		} else {
			familyLoaderMap[familyTime] = &loader{
				familyTime: rs.FamilyTime(),
				interval:   rs.Interval(),
				families:   []*familyLoader{fl},
				timeRange:  rs.SlotRange(),
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
	// set family index
	for familyIndex, family := range familyLoaders {
		for _, fl := range family.families {
			fl.familyIndex = familyIndex
		}
	}
	return familyLoaders
}
