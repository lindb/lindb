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

	"github.com/lindb/common/models"
	"github.com/lindb/roaring"

	"github.com/lindb/lindb/constants"
	"github.com/lindb/lindb/flow"
	"github.com/lindb/lindb/index"
	"github.com/lindb/lindb/series/tag"
	"github.com/lindb/lindb/sql/stmt"
	"github.com/lindb/lindb/tsdb"
)

// seriesFiltering represents series filtering operator.
type seriesFiltering struct {
	executeCtx   *flow.ShardExecuteContext
	indexSegment index.MetricIndexSegment

	err error
}

// NewSeriesFiltering creates a seriesFiltering instance.
func NewSeriesFiltering(executeCtx *flow.ShardExecuteContext, shard tsdb.Shard) Operator {
	return &seriesFiltering{
		executeCtx:   executeCtx,
		indexSegment: shard.IndexSegment(),
	}
}

// Execute executes filtering series ids based on tag values result set.
func (op *seriesFiltering) Execute() error {
	queryStmt := op.executeCtx.StorageExecuteCtx.Query
	// if it gets tag filter result do series ids searching
	_, seriesIDs := op.findSeriesIDsByExpr(queryStmt.Condition)
	if op.err != nil {
		return op.err
	}
	op.executeCtx.SeriesIDsAfterFiltering.Or(seriesIDs)
	return nil
}

// findSeriesIDsByExpr finds series ids by expr, recursion filter for expr
func (op *seriesFiltering) findSeriesIDsByExpr(condition stmt.Expr) (tag.KeyID, *roaring.Bitmap) {
	if condition == nil {
		return 0, roaring.New() // create an empty series ids for parent expr
	}
	if op.err != nil {
		return 0, roaring.New() // create an empty series ids for parent expr
	}
	switch expr := condition.(type) {
	case stmt.TagFilter:
		tagKey, seriesIDs, err := op.getSeriesIDsByExpr(expr)
		if err != nil {
			op.err = err
			return tagKey, roaring.New() // create an empty series ids for parent expr
		}
		return tagKey, seriesIDs
	case *stmt.ParenExpr:
		return op.findSeriesIDsByExpr(expr.Expr)
	case *stmt.NotExpr:
		// get filter series ids
		tagKey, matchResult := op.findSeriesIDsByExpr(expr.Expr)
		// get all series ids for tag key
		all, err := op.indexSegment.GetSeriesIDsForTag(tagKey, op.executeCtx.StorageExecuteCtx.Query.TimeRange)
		if err != nil {
			op.err = err
			return tagKey, roaring.New() // create an empty series ids for parent expr
		}
		// do and not got series ids not in 'a' list
		all.AndNot(matchResult)
		return 0, all
	case *stmt.BinaryExpr:
		_, left := op.findSeriesIDsByExpr(expr.Left)
		_, right := op.findSeriesIDsByExpr(expr.Right)
		if expr.Operator == stmt.AND {
			left.And(right)
		} else {
			left.Or(right)
		}
		return 0, left
	}
	return 0, roaring.New() // create an empty series ids for parent expr
}

// getTagKeyID returns the tag key id by tag key
func (op *seriesFiltering) getSeriesIDsByExpr(expr stmt.Expr) (tag.KeyID, *roaring.Bitmap, error) {
	tagValues, ok := op.executeCtx.StorageExecuteCtx.TagFilterResult[expr.Rewrite()]
	if !ok {
		return 0, nil, fmt.Errorf("%w, expr: %s", constants.ErrTagValueFilterResultNotFound, expr.Rewrite())
	}
	seriesIDs, err := op.indexSegment.GetSeriesIDsByTagValueIDs(tagValues.TagKeyID,
		tagValues.TagValueIDs, op.executeCtx.StorageExecuteCtx.Query.TimeRange)
	if err != nil {
		return 0, nil, err
	}
	return tagValues.TagKeyID, seriesIDs, nil
}

// Identifier returns identifier value of series filtering operator.
func (op *seriesFiltering) Identifier() string {
	return "Series Filtering"
}

// Stats returns the stats of series filtering operator.
func (op *seriesFiltering) Stats() interface{} {
	return &models.SeriesStats{
		NumOfSeries: op.executeCtx.SeriesIDsAfterFiltering.GetCardinality(),
	}
}
