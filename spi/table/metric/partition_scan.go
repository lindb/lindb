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
	"errors"
	"fmt"

	"github.com/lindb/roaring"

	"github.com/lindb/lindb/constants"
	"github.com/lindb/lindb/flow"
	"github.com/lindb/lindb/series/tag"
	"github.com/lindb/lindb/sql/tree"
)

type partitionScan struct {
	ctx       *ExecutionContext
	partition *Partition
	reduceCh  chan<- any
}

func NewPartitionScan(ctx *ExecutionContext, partition *Partition, reduceCh chan<- any) *partitionScan {
	return &partitionScan{
		ctx:       ctx,
		partition: partition,
		reduceCh:  reduceCh,
	}
}

func (ps *partitionScan) Run() {
	// TODO: handle panic
	seriesIDs := ps.findSeriesIDs(ps.partition)
	var groupingContext flow.GroupingContext

	if seriesIDs.IsEmpty() {
		panic(constants.ErrSeriesIDNotFound)
	}

	var err error
	var seriesIDsAfterGrouping *roaring.Bitmap

	tableScan := ps.partition.tableScan
	if tableScan.isGrouping() {
		// if it has grouping, do group by tag keys, else just split series ids as batch first.
		seriesIDsAfterGrouping, groupingContext, err = ps.partition.shard.IndexDB().
			GetGroupingContext(tableScan.grouping.tags, seriesIDs)
		if err != nil && !errors.Is(err, constants.ErrNotFound) {
			// TODO: add not found check
			panic(err)
		}
		// maybe filtering some series ids after grouping that is result of filtering.
		// if not found, return empty series ids.
		seriesIDs = seriesIDsAfterGrouping
		fmt.Println("groupiung.....")
	}

	highKeys := seriesIDs.GetHighKeys()
	for index, highKey := range highKeys {
		data := &DataSplit{
			partition:       ps.partition,
			groupingContext: groupingContext,

			seriesIDHighKey: highKey,
			lowSeriesIDs:    seriesIDs.GetContainerAtIndex(index),
		}
		if tableScan.fields.Len() == 0 {
			fmt.Printf("sereis meta grouping=%v\n", data.groupingContext)
			// if fields is empty, then do series query(send reducer task).
			ps.reduceCh <- data
		} else {
			execute(ps.ctx, tableScan.db.ExecutorPool().DataFetcher, func() {
				dataScan := NewDataScan(data, ps.reduceCh)
				dataScan.Run()
			})
		}
	}
}

func (ps *partitionScan) findSeriesIDs(partition *Partition) *roaring.Bitmap {
	tableScan := partition.tableScan
	fmt.Printf("families=%v,fields=%v\n", partition.segments, tableScan.fields)
	if tableScan.fields.Len() == 0 && len(partition.segments) == 0 {
		// no data family return empty series ids
		return roaring.New()
	}

	seriesIDs := ps.lookupSeriesIDs(partition)
	result := roaring.New()

	for i := range partition.segments {
		family := partition.segments[i]
		// check family data if matches condition(series ids)
		resultSet, err := family.Filter(&flow.MetricScanContext{
			MetricID:  tableScan.metricID,
			SeriesIDs: seriesIDs,
			Fields:    tableScan.fields, // set fields when search the data of time series
			TimeRange: tableScan.timeRange,
		})

		if err != nil && !errors.Is(err, constants.ErrNotFound) {
			panic(err)
		}

		for i := range resultSet {
			rs := resultSet[i]

			// check double, maybe some series ids be filtered out when do grouping.
			finalSeriesIDs := roaring.FastAnd(seriesIDs, rs.SeriesIDs())
			if finalSeriesIDs.IsEmpty() {
				continue
			}

			result.Or(finalSeriesIDs)

			if tableScan.fields.Len() > 0 {
				partition.fieldsData = append(partition.fieldsData, rs)
			}
		}
	}
	return result
}

func (ps *partitionScan) lookupSeriesIDs(partition *Partition) *roaring.Bitmap {
	var (
		seriesIDs *roaring.Bitmap
		err       error
		ok        bool
	)
	tableScan := partition.tableScan
	predicate := tableScan.predicate

	if predicate == nil {
		// if predicate nil, find all series ids under metric
		seriesIDs, err = partition.shard.IndexDB().GetSeriesIDsForMetric(tableScan.metricID)
		if err != nil {
			panic(err)
		}
	} else {
		// find series ids based on where condition
		lookup := NewRowLookupVisitor(partition)
		if seriesIDs, ok = predicate.Accept(nil, lookup).(*roaring.Bitmap); !ok {
			panic(constants.ErrSeriesIDNotFound)
		}
	}

	if seriesIDs == nil || seriesIDs.IsEmpty() {
		panic(constants.ErrSeriesIDNotFound)
	}
	return seriesIDs
}

type RowsLookupVisitor struct {
	partition *Partition
}

func NewRowLookupVisitor(partition *Partition) *RowsLookupVisitor {
	return &RowsLookupVisitor{
		partition: partition,
	}
}

func (v *RowsLookupVisitor) Visit(context any, n tree.Node) any {
	fmt.Printf("row lookup visitor: %v\n", v.partition.tableScan.filterResult)
	var seriesIDs *roaring.Bitmap
	var tagKey tag.KeyID
	indexDB := v.partition.shard.IndexDB()

	switch node := n.(type) {
	case *tree.ComparisonExpression:
		tagKey, seriesIDs = v.visitPredicate(node)
		if node.Operator == tree.ComparisonNEQ {
			// get all series ids for tag key
			all, err := indexDB.GetSeriesIDsForTag(tagKey)
			if err != nil {
				panic(err)
			}
			// do and not got series ids not in 'a' list
			all.AndNot(seriesIDs)
			return all
		}
	// TODO: add null predicate
	case *tree.InPredicate, *tree.RegexPredicate, *tree.LikePredicate:
		_, seriesIDs = v.visitPredicate(node)
	case *tree.NotExpression:
		// get filter series ids
		tagKey, seriesIDs = v.visitPredicate(node.Value)
		// TODO: cache if dup
		// get all series ids for tag key
		all, err := indexDB.GetSeriesIDsForTag(tagKey)
		if err != nil {
			panic(err)
		}
		// do and not got series ids not in 'a' list
		all.AndNot(seriesIDs)
		return all
	case *tree.LogicalExpression:
		for _, term := range node.Terms {
			matchResult := term.Accept(context, v).(*roaring.Bitmap)
			if seriesIDs == nil {
				seriesIDs = matchResult
			} else {
				if node.Operator == tree.LogicalAND {
					seriesIDs.And(matchResult)
				} else {
					seriesIDs.Or(matchResult)
				}
			}
		}
		return seriesIDs
	case *tree.Cast:
		return node.Expression.Accept(context, v)
	}
	return seriesIDs
}

func (v *RowsLookupVisitor) visitPredicate(node tree.Node) (tag.KeyID, *roaring.Bitmap) {
	columnResult, ok := v.partition.tableScan.filterResult[node.GetID()]
	if !ok {
		panic(constants.ErrSeriesIDNotFound)
	}
	fmt.Printf("tag value ids=%v\n", columnResult.TagValueIDs)
	indexDB := v.partition.shard.IndexDB()
	seriesIDs, err := indexDB.GetSeriesIDsByTagValueIDs(columnResult.TagKeyID, columnResult.TagValueIDs)
	if err != nil {
		panic(err)
	}
	return columnResult.TagKeyID, seriesIDs
}
