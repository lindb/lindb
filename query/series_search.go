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

package query

import (
	"fmt"

	"github.com/lindb/roaring"

	"github.com/lindb/lindb/constants"
	"github.com/lindb/lindb/series"
	"github.com/lindb/lindb/sql/stmt"
)

//go:generate mockgen -source ./series_search.go -destination=./series_search_mock.go -package=query

// SeriesSearch represents a series search by condition expression
type SeriesSearch interface {
	// Search searches series ids base on condition, if search fail return nil, else return series ids
	Search() (*roaring.Bitmap, error)
}

// seriesSearch represents a series search by condition expression,
// only do tag filter, return series ids.
// return series id set for condition
type seriesSearch struct {
	condition    stmt.Expr
	filterResult map[string]*tagFilterResult

	filter series.Filter

	err error
}

// newSeriesSearch creates a a series search using query condition
func newSeriesSearch(filter series.Filter, filterResult map[string]*tagFilterResult, condition stmt.Expr) SeriesSearch {
	return &seriesSearch{
		filterResult: filterResult,
		filter:       filter,
		condition:    condition,
	}
}

// Search searches series ids base on condition, if search fail return nil, else return series ids
func (s *seriesSearch) Search() (*roaring.Bitmap, error) {
	_, seriesIDs := s.findSeriesIDsByExpr(s.condition)
	if s.err != nil {
		return nil, s.err
	}
	return seriesIDs, nil
}

// findSeriesIDsByExpr finds series ids by expr, recursion filter for expr
func (s *seriesSearch) findSeriesIDsByExpr(condition stmt.Expr) (uint32, *roaring.Bitmap) {
	if condition == nil {
		return 0, roaring.New() // create a empty series ids for parent expr
	}
	if s.err != nil {
		return 0, roaring.New() // create a empty series ids for parent expr
	}
	switch expr := condition.(type) {
	case stmt.TagFilter:
		tagKey, seriesIDs, err := s.getSeriesIDsByExpr(expr)
		if err != nil {
			s.err = err
			return tagKey, roaring.New() // create a empty series ids for parent expr
		}
		return tagKey, seriesIDs
	case *stmt.ParenExpr:
		return s.findSeriesIDsByExpr(expr.Expr)
	case *stmt.NotExpr:
		// get filter series ids
		tagKey, matchResult := s.findSeriesIDsByExpr(expr.Expr)
		// get all series ids for tag key
		all, err := s.filter.GetSeriesIDsForTag(tagKey)
		if err != nil {
			s.err = err
			return tagKey, roaring.New() // create a empty series ids for parent expr
		}
		// do and not got series ids not in 'a' list
		all.AndNot(matchResult)
		return 0, all
	case *stmt.BinaryExpr:
		_, left := s.findSeriesIDsByExpr(expr.Left)
		_, right := s.findSeriesIDsByExpr(expr.Right)
		if expr.Operator == stmt.AND {
			left.And(right)
		} else {
			left.Or(right)
		}
		return 0, left
	}
	return 0, roaring.New() // create a empty series ids for parent expr
}

// getTagKeyID returns the tag key id by tag key
func (s *seriesSearch) getSeriesIDsByExpr(expr stmt.Expr) (uint32, *roaring.Bitmap, error) {
	tagValues, ok := s.filterResult[expr.Rewrite()]
	if !ok {
		return 0, nil, fmt.Errorf("%w, expr: %s", constants.ErrTagValueFilterResultNotFound, expr.Rewrite())
	}
	seriesIDs, err := s.filter.GetSeriesIDsByTagValueIDs(tagValues.tagKey, tagValues.tagValueIDs)
	if err != nil {
		return 0, nil, err
	}
	return tagValues.tagKey, seriesIDs, nil
}
