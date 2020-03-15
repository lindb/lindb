package query

import (
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/lindb/roaring"
	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/series"
	"github.com/lindb/lindb/sql"
	"github.com/lindb/lindb/sql/stmt"
)

func TestSeriesSearch_Search(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockFilter := series.NewMockFilter(ctrl)
	seriesIDs := roaring.BitmapOf(10, 20, 30)

	// case 1: empty filter expr
	query, _ := sql.Parse("select f from cpu")
	search := newSeriesSearch(mockFilter, nil, query)
	resultSet, err := search.Search()
	assert.NoError(t, err)
	assert.Equal(t, uint64(0), resultSet.GetCardinality())
	// case 2: equal tag filter
	query, _ = sql.Parse("select f from cpu where ip='1.1.1.1'")
	mockFilter.EXPECT().GetSeriesIDsByTagValueIDs(uint32(1), gomock.Any()).Return(seriesIDs.Clone(), nil)
	search = newSeriesSearch(mockFilter, mockFilterResult(), query)
	resultSet, err = search.Search()
	assert.NoError(t, err)
	assert.Equal(t, seriesIDs, resultSet)
	// case 3: not expr
	query, _ = sql.Parse("select f from cpu where ip!='1.1.1.1'")
	mockFilter.EXPECT().GetSeriesIDsByTagValueIDs(uint32(1), gomock.Any()).Return(seriesIDs.Clone(), nil)
	mockFilter.EXPECT().GetSeriesIDsForTag(uint32(1)).Return(roaring.BitmapOf(10, 20, 40, 50), nil)
	search = newSeriesSearch(mockFilter, mockFilterResult(), query)
	resultSet, err = search.Search()
	assert.NoError(t, err)
	assert.Equal(t, roaring.BitmapOf(40, 50), resultSet)
	// case 4: binary expr and
	query, _ = sql.Parse("select f from cpu " +
		"where ip='1.1.1.1' and path='/data' and time>'20190410 00:00:00' and time<'20190410 10:00:00'")
	mockFilter.EXPECT().GetSeriesIDsByTagValueIDs(uint32(1), gomock.Any()).Return(seriesIDs.Clone(), nil)
	mockFilter.EXPECT().GetSeriesIDsByTagValueIDs(uint32(2), gomock.Any()).Return(roaring.BitmapOf(20), nil)
	search = newSeriesSearch(mockFilter, mockFilterResult(), query)
	resultSet, err = search.Search()
	assert.NoError(t, err)
	assert.Equal(t, roaring.BitmapOf(20), resultSet)
	// case 5: binary expr or
	query, _ = sql.Parse("select f from cpu " +
		"where ip='1.1.1.1' or path='/data' and time>'20190410 00:00:00' and time<'20190410 10:00:00'")
	mockFilter.EXPECT().GetSeriesIDsByTagValueIDs(uint32(1), gomock.Any()).Return(seriesIDs.Clone(), nil)
	mockFilter.EXPECT().GetSeriesIDsByTagValueIDs(uint32(2), gomock.Any()).Return(roaring.BitmapOf(200), nil)
	search = newSeriesSearch(mockFilter, mockFilterResult(), query)
	resultSet, err = search.Search()
	assert.NoError(t, err)
	assert.Equal(t, roaring.BitmapOf(10, 20, 30, 200), resultSet)
	// case 6: paren expr
	query, _ = sql.Parse("select f from cpu where (ip='1.1.1.1')")
	mockFilter.EXPECT().GetSeriesIDsByTagValueIDs(uint32(1), gomock.Any()).Return(seriesIDs.Clone(), nil)
	search = newSeriesSearch(mockFilter, mockFilterResult(), query)
	resultSet, err = search.Search()
	assert.NoError(t, err)
	assert.Equal(t, seriesIDs, resultSet)
}

func TestSeriesSearch_Search_err(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockFilter := series.NewMockFilter(ctrl)
	seriesIDs := roaring.BitmapOf(10, 20, 30)

	// case 1: expr not exist
	query, _ := sql.Parse("select f from cpu where ip='1.1.1.1'")
	search := newSeriesSearch(mockFilter, make(map[string]*tagFilterResult), query)
	resultSet, err := search.Search()
	assert.Error(t, err)
	assert.Nil(t, resultSet)
	// case 2: get series id err
	search = newSeriesSearch(mockFilter, mockFilterResult(), query)
	mockFilter.EXPECT().GetSeriesIDsByTagValueIDs(gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("err"))
	resultSet, err = search.Search()
	assert.Error(t, err)
	assert.Nil(t, resultSet)
	// case 3: not expr err
	query, _ = sql.Parse("select f from cpu where ip!='1.1.1.1'")
	mockFilter.EXPECT().GetSeriesIDsByTagValueIDs(uint32(1), gomock.Any()).Return(seriesIDs, nil)
	mockFilter.EXPECT().GetSeriesIDsForTag(uint32(1)).Return(nil, fmt.Errorf("err"))
	search = newSeriesSearch(mockFilter, mockFilterResult(), query)
	resultSet, err = search.Search()
	assert.Error(t, err)
	assert.Nil(t, resultSet)
	// case 4: recursion err
	query, _ = sql.Parse("select f from cpu where ip='1.1.1.1' or ip='1.1.1.1'")
	mockFilter.EXPECT().GetSeriesIDsByTagValueIDs(uint32(1), gomock.Any()).Return(nil, fmt.Errorf("err"))
	search = newSeriesSearch(mockFilter, mockFilterResult(), query)
	resultSet, err = search.Search()
	assert.Error(t, err)
	assert.Nil(t, resultSet)
}

func TestSeriesSearch_Search_expr_not_match(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockFilter := series.NewMockFilter(ctrl)

	query, _ := sql.Parse("select f from cpu where ip='1.1.1.1'")
	query.Condition = &stmt.CallExpr{}
	search := newSeriesSearch(mockFilter, make(map[string]*tagFilterResult), query)
	resultSet, err := search.Search()
	assert.NoError(t, err)
	assert.Equal(t, uint64(0), resultSet.GetCardinality())
}

func TestSeriesSearch_Search_complex(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockFilter := series.NewMockFilter(ctrl)

	query, _ := sql.Parse("select f from cpu" +
		" where (ip not in ('1.1.1.1','2.2.2.2') and region='sh') and (path='/data' or path='/home')")
	mockFilter.EXPECT().GetSeriesIDsByTagValueIDs(uint32(1), roaring.BitmapOf(5)).Return(roaring.BitmapOf(1, 2), nil)
	mockFilter.EXPECT().GetSeriesIDsForTag(uint32(1)).Return(roaring.BitmapOf(1, 2, 3, 4, 5, 6, 7), nil)
	mockFilter.EXPECT().GetSeriesIDsByTagValueIDs(uint32(3), roaring.BitmapOf(4)).Return(roaring.BitmapOf(3, 5, 6, 7), nil)
	mockFilter.EXPECT().GetSeriesIDsByTagValueIDs(uint32(2), roaring.BitmapOf(2)).Return(roaring.BitmapOf(7), nil)
	mockFilter.EXPECT().GetSeriesIDsByTagValueIDs(uint32(2), roaring.BitmapOf(3)).Return(roaring.BitmapOf(5), nil)
	search := newSeriesSearch(mockFilter, mockFilterResult(), query)
	resultSet, err := search.Search()
	assert.NoError(t, err)
	assert.Equal(t, roaring.BitmapOf(5, 7), resultSet)
}

func mockFilterResult() map[string]*tagFilterResult {
	result := make(map[string]*tagFilterResult)
	result[(&stmt.EqualsExpr{Key: "ip", Value: "1.1.1.1"}).Rewrite()] = &tagFilterResult{
		tagKey:      1,
		tagValueIDs: roaring.BitmapOf(1),
	}
	result[(&stmt.EqualsExpr{Key: "path", Value: "/data"}).Rewrite()] = &tagFilterResult{
		tagKey:      2,
		tagValueIDs: roaring.BitmapOf(2),
	}
	result[(&stmt.EqualsExpr{Key: "path", Value: "/home"}).Rewrite()] = &tagFilterResult{
		tagKey:      2,
		tagValueIDs: roaring.BitmapOf(3),
	}
	result[(&stmt.EqualsExpr{Key: "region", Value: "sh"}).Rewrite()] = &tagFilterResult{
		tagKey:      3,
		tagValueIDs: roaring.BitmapOf(4),
	}
	result[(&stmt.InExpr{Key: "ip", Values: []string{"1.1.1.1", "2.2.2.2"}}).Rewrite()] = &tagFilterResult{
		tagKey:      1,
		tagValueIDs: roaring.BitmapOf(5),
	}
	return result
}

//
//func TestComplexCondition(t *testing.T) {
//	ctrl := gomock.NewController(t)
//	defer ctrl.Finish()
//	mockFilter := series.NewMockFilter(ctrl)
//
//	query, _ := sql.Parse("select f from cpu" +
//		" where (ip not in ('1.1.1.1','2.2.2.2') and region='sh') and (path='/data' or path='/home')")
//	mockFilter.EXPECT().
//		FindSeriesIDsByExpr(gomock.Any(), &stmt.InExpr{Key: "ip", Values: []string{"1.1.1.1", "2.2.2.2"}}, query.TimeRange).
//		Return(mockSeriesIDSet(series.Version(11), roaring.BitmapOf(1, 2, 4)), nil)
//	mockFilter.EXPECT().
//		GetSeriesIDsForTag(gomock.Any(), query.TimeRange).
//		Return(mockSeriesIDSet(series.Version(11), roaring.BitmapOf(1, 2, 3, 4, 6, 7, 8)), nil)
//	mockFilter.EXPECT().
//		FindSeriesIDsByExpr(gomock.Any(), &stmt.EqualsExpr{Key: "region", Value: "sh"}, query.TimeRange).
//		Return(mockSeriesIDSet(series.Version(11), roaring.BitmapOf(2, 3, 4, 7)), nil)
//	mockFilter.EXPECT().
//		FindSeriesIDsByExpr(gomock.Any(), &stmt.EqualsExpr{Key: "path", Value: "/data"}, query.TimeRange).
//		Return(mockSeriesIDSet(series.Version(11), roaring.BitmapOf(3, 5)), nil)
//	mockFilter.EXPECT().
//		FindSeriesIDsByExpr(gomock.Any(), &stmt.EqualsExpr{Key: "path", Value: "/home"}, query.TimeRange).
//		Return(mockSeriesIDSet(series.Version(11), roaring.BitmapOf(1)), nil)
//	search := newSeriesSearch(10, mockFilter, query)
//	resultSet, _ := search.Search()
//	// ip not in ('1.1.1.1','2.2.2.2') => 3,6,7,8
//	// ip not in ('1.1.1.1','2.2.2.2') and region='sh' => 3,7
//	// path='/data' or path='/home' => 1,3,5
//	// final => 3
//	assert.Equal(t, *mockSeriesIDSet(series.Version(11), roaring.BitmapOf(3)), *resultSet)
//
//	// error
//	mockFilter1 := series.NewMockFilter(ctrl)
//	mockFilter1.EXPECT().
//		FindSeriesIDsByExpr(gomock.Any(), &stmt.InExpr{Key: "ip", Values: []string{"1.1.1.1", "2.2.2.2"}}, query.TimeRange).
//		Return(mockSeriesIDSet(series.Version(11), roaring.BitmapOf(1, 2, 4)), nil)
//	mockFilter1.EXPECT().
//		GetSeriesIDsForTag(gomock.Any(), query.TimeRange).
//		Return(mockSeriesIDSet(series.Version(11), roaring.BitmapOf(1, 2, 3, 4, 6, 7, 8)), nil)
//	mockFilter1.EXPECT().
//		FindSeriesIDsByExpr(gomock.Any(), &stmt.EqualsExpr{Key: "region", Value: "sh"}, query.TimeRange).
//		Return(nil, errors.New("complex error"))
//	search = newSeriesSearch(10, mockFilter1, query)
//	resultSet, err := search.Search()
//	assert.Nil(t, resultSet)
//	assert.NotNil(t, err)
//}
//
//func TestSeriesSearch_condition_fail(t *testing.T) {
//	search := newSeriesSearch(10, nil, nil)
//	result, _ := search.findSeriesIDsByExpr(nil)
//	assert.Nil(t, result)
//
//	search = newSeriesSearch(10, nil, nil)
//	result, _ = search.findSeriesIDsByExpr(&stmt.BinaryExpr{Operator: stmt.ADD})
//	assert.Nil(t, result)
//
//	query, _ := sql.Parse("select f from disk " +
//		"where (ip='1.1.1.1')")
//	search = newSeriesSearch(10, nil, query)
//	search.err = fmt.Errorf("err")
//	resultSet, err := search.Search()
//	assert.Nil(t, resultSet)
//	assert.NotNil(t, err)
//}
//
