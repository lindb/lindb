package query

import (
	"github.com/eleme/lindb/sql/stmt"
	"github.com/eleme/lindb/tsdb/index"
	"github.com/eleme/lindb/tsdb/series"
)

// seriesSearch represents a series search by condition expression,
// only do tag filter, return series ids.
// return multi-version series id set for condition
type seriesSearch struct {
	metricID uint32
	query    *stmt.Query

	index index.Index

	resultSet *series.MultiVerSeriesIDSet
	err       error
}

// newSeriesSearch creates a a series search using query condition
func newSeriesSearch(metricID uint32, index index.Index, query *stmt.Query) *seriesSearch {
	return &seriesSearch{
		metricID: metricID,
		index:    index,
		query:    query,
	}
}

// search searches series ids based on query condition
func (s *seriesSearch) search() {
	condition := s.query.Condition
	if condition == nil {
		return
	}
	seriesIDs, _ := s.findSeriesIDsByExpr(condition)
	s.resultSet = seriesIDs
}

// error returns error, if search fail
func (s *seriesSearch) error() error {
	return s.err
}

// getResultSet return series ids result set, if search success
func (s *seriesSearch) getResultSet() *series.MultiVerSeriesIDSet {
	return s.resultSet
}

// findSeriesIDsByExpr finds series ids by expr, recursion filter for expr
func (s *seriesSearch) findSeriesIDsByExpr(condition stmt.Expr) (series *series.MultiVerSeriesIDSet, tagKey string) {
	if condition == nil {
		return series, tagKey
	}
	if s.err != nil {
		return series, tagKey
	}
	switch expr := condition.(type) {
	case stmt.TagFilter:
		result, err := s.index.FindSeriesIDsByExpr(s.metricID, expr, s.query.TimeRange)
		if err != nil {
			s.err = err
			return
		}
		series = result
		tagKey = expr.TagKey()
	case *stmt.ParenExpr:
		series, tagKey = s.findSeriesIDsByExpr(expr.Expr)
	case *stmt.NotExpr:
		// find series ids by expr => a
		matchResult, tagKey := s.findSeriesIDsByExpr(expr.Expr)
		if len(tagKey) > 0 {
			// get all series ids for tag key
			all, err := s.index.GetSeriesIDsForTag(s.metricID, tagKey, s.query.TimeRange)
			if err != nil {
				s.err = err
				return nil, tagKey
			}
			// do and not got series ids not in 'a' list
			all.AndNot(matchResult)
			series = all
			return series, tagKey
		}
	case *stmt.BinaryExpr:
		if expr.Operator != stmt.AND && expr.Operator != stmt.OR {
			return series, tagKey
		}
		left, _ := s.findSeriesIDsByExpr(expr.Left)
		if left == nil {
			return series, tagKey
		}
		right, _ := s.findSeriesIDsByExpr(expr.Right)
		if right == nil {
			return series, tagKey
		}

		if expr.Operator == stmt.AND {
			left.And(right)
		} else {
			left.Or(right)
		}
		series = left
	}
	return series, tagKey
}
