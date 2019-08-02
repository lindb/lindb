package query

import (
	"github.com/lindb/lindb/sql/stmt"
	"github.com/lindb/lindb/tsdb/index"
	"github.com/lindb/lindb/tsdb/series"
)

//go:generate mockgen -source=./series_search.go -destination=./series_search_mock.go -package=query

// SeriesSearch represents a series search by condition expression
type SeriesSearch interface {
	// Search searches series ids base on condition, if search fail return nil, else return multi-version series ids
	Search() (*series.MultiVerSeriesIDSet, error)
}

// seriesSearch represents a series search by condition expression,
// only do tag filter, return series ids.
// return multi-version series id set for condition
type seriesSearch struct {
	metricID uint32
	query    *stmt.Query

	filter index.SeriesIDsFilter

	err error
}

// newSeriesSearch creates a a series search using query condition
func newSeriesSearch(metricID uint32, filter index.SeriesIDsFilter, query *stmt.Query) *seriesSearch {
	return &seriesSearch{
		metricID: metricID,
		filter:   filter,
		query:    query,
	}
}

// Search searches series ids base on condition, if search fail return nil, else return multi-version series ids
func (s *seriesSearch) Search() (*series.MultiVerSeriesIDSet, error) {
	condition := s.query.Condition
	if condition == nil {
		return nil, nil
	}
	seriesIDs, _ := s.findSeriesIDsByExpr(condition)
	if s.err != nil {
		return nil, s.err
	}
	return seriesIDs, nil
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
		result, err := s.filter.FindSeriesIDsByExpr(s.metricID, expr, s.query.TimeRange)
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
			all, err := s.filter.GetSeriesIDsForTag(s.metricID, tagKey, s.query.TimeRange)
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
