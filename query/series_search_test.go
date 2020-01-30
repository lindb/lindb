package query

import (
	"errors"
	"fmt"
	"testing"

	"github.com/lindb/lindb/series"
	"github.com/lindb/lindb/sql"
	"github.com/lindb/lindb/sql/stmt"

	"github.com/golang/mock/gomock"
	"github.com/lindb/roaring"
	"github.com/stretchr/testify/assert"
)

func TestSampleCondition(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockFilter := series.NewMockFilter(ctrl)
	s := mockSeriesIDSet(series.Version(1), roaring.BitmapOf(1, 2, 3, 4))

	query, _ := sql.Parse("select f from cpu")
	search := newSeriesSearch(1, mockFilter, query)
	resultSet, _ := search.Search()
	assert.Nil(t, resultSet)

	query, _ = sql.Parse("select f from cpu where ip='1.1.1.1'")
	mockFilter.EXPECT().
		FindSeriesIDsByExpr(uint32(1), &stmt.EqualsExpr{Key: "ip", Value: "1.1.1.1"}, query.TimeRange).
		Return(s, nil)
	search = newSeriesSearch(1, mockFilter, query)
	resultSet, _ = search.Search()
	assert.Equal(t, *s, *resultSet)

	query, _ = sql.Parse("select f from cpu where ip like '1.1.*.1'")
	mockFilter.EXPECT().
		FindSeriesIDsByExpr(uint32(1), &stmt.LikeExpr{Key: "ip", Value: "1.1.*.1"}, query.TimeRange).
		Return(s, nil)
	search = newSeriesSearch(1, mockFilter, query)
	resultSet, _ = search.Search()
	assert.Equal(t, *s, *resultSet)

	query, _ = sql.Parse("select f from cpu where ip =~ '1.1.*.1'")
	mockFilter.EXPECT().
		FindSeriesIDsByExpr(uint32(1), &stmt.RegexExpr{Key: "ip", Regexp: "1.1.*.1"}, query.TimeRange).
		Return(s, nil)
	search = newSeriesSearch(1, mockFilter, query)
	resultSet, _ = search.Search()
	assert.Equal(t, *s, *resultSet)

	query, _ = sql.Parse("select f from cpu where ip in ('1.1.1.1','1.1.3.3')")
	mockFilter.EXPECT().
		FindSeriesIDsByExpr(uint32(1), &stmt.InExpr{Key: "ip", Values: []string{"1.1.1.1", "1.1.3.3"}}, query.TimeRange).
		Return(s, nil)
	search = newSeriesSearch(1, mockFilter, query)
	resultSet, _ = search.Search()
	assert.Equal(t, *s, *resultSet)

	// search error
	query, _ = sql.Parse("select f from cpu where ip='1.1.1.1'")
	mockFilter.EXPECT().
		FindSeriesIDsByExpr(uint32(1), &stmt.EqualsExpr{Key: "ip", Value: "1.1.1.1"}, query.TimeRange).
		Return(nil, errors.New("search error"))
	search = newSeriesSearch(1, mockFilter, query)
	resultSet, err := search.Search()
	assert.Nil(t, resultSet)
	assert.NotNil(t, err)
}

func TestNotCondition(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockFilter := series.NewMockFilter(ctrl)

	query, _ := sql.Parse("select f from cpu where ip!='1.1.1.1'")
	mockFilter.EXPECT().
		FindSeriesIDsByExpr(gomock.Any(), &stmt.EqualsExpr{Key: "ip", Value: "1.1.1.1"}, query.TimeRange).
		Return(mockSeriesIDSet(series.Version(11), roaring.BitmapOf(3, 4)), nil)

	mockFilter.EXPECT().
		GetSeriesIDsForTag(gomock.Any(), query.TimeRange).
		Return(mockSeriesIDSet(series.Version(11), roaring.BitmapOf(1, 2, 3, 4)), nil)
	search := newSeriesSearch(1, mockFilter, query)
	resultSet, _ := search.Search()
	assert.Equal(t, *mockSeriesIDSet(series.Version(11), roaring.BitmapOf(1, 2)), *resultSet)

	// error
	query, _ = sql.Parse("select f from cpu where ip!='1.1.1.1'")
	mockFilter.EXPECT().
		FindSeriesIDsByExpr(gomock.Any(), &stmt.EqualsExpr{Key: "ip", Value: "1.1.1.1"}, query.TimeRange).
		Return(mockSeriesIDSet(series.Version(11), roaring.BitmapOf(3, 4)), nil)

	mockFilter.EXPECT().GetSeriesIDsForTag(uint32(1), query.TimeRange).
		Return(nil, errors.New("get series ids error"))
	search = newSeriesSearch(1, mockFilter, query)
	resultSet, err := search.Search()
	assert.Nil(t, resultSet)
	assert.NotNil(t, err)
}

func TestBinaryCondition(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockFilter := series.NewMockFilter(ctrl)

	// and
	query, _ := sql.Parse("select f from cpu " +
		"where ip='1.1.1.1' and path='/data' and time>'20190410 00:00:00' and time<'20190410 10:00:00'")
	mockFilter.EXPECT().
		FindSeriesIDsByExpr(uint32(1), &stmt.EqualsExpr{Key: "ip", Value: "1.1.1.1"}, query.TimeRange).
		Return(mockSeriesIDSet(series.Version(11), roaring.BitmapOf(1, 2, 3, 4)), nil)
	mockFilter.EXPECT().
		FindSeriesIDsByExpr(uint32(1), &stmt.EqualsExpr{Key: "path", Value: "/data"}, query.TimeRange).
		Return(mockSeriesIDSet(series.Version(11), roaring.BitmapOf(3, 5)), nil)
	search := newSeriesSearch(1, mockFilter, query)
	resultSet, _ := search.Search()
	assert.Equal(t, *mockSeriesIDSet(series.Version(11), roaring.BitmapOf(3)), *resultSet)

	// or
	mockFilter2 := series.NewMockFilter(ctrl)
	query, _ = sql.Parse("select f from cpu " +
		"where ip='1.1.1.1' or path='/data' and time>'20190410 00:00:00' and time<'20190410 10:00:00'")
	mockFilter2.EXPECT().
		FindSeriesIDsByExpr(uint32(1), &stmt.EqualsExpr{Key: "ip", Value: "1.1.1.1"}, query.TimeRange).
		Return(mockSeriesIDSet(series.Version(11), roaring.BitmapOf(1, 2, 3, 4)), nil)
	mockFilter2.EXPECT().
		FindSeriesIDsByExpr(uint32(1), &stmt.EqualsExpr{Key: "path", Value: "/data"}, query.TimeRange).
		Return(mockSeriesIDSet(series.Version(11), roaring.BitmapOf(3, 5)), nil)
	search = newSeriesSearch(1, mockFilter2, query)
	resultSet, _ = search.Search()
	assert.Equal(t, *mockSeriesIDSet(series.Version(11), roaring.BitmapOf(1, 2, 3, 4, 5)), *resultSet)

	// error
	mockFilter3 := series.NewMockFilter(ctrl)
	mockFilter3.EXPECT().
		FindSeriesIDsByExpr(uint32(1), &stmt.EqualsExpr{Key: "ip", Value: "1.1.1.1"}, query.TimeRange).
		Return(nil, errors.New("left error"))
	search = newSeriesSearch(1, mockFilter3, query)
	resultSet, err := search.Search()
	assert.Nil(t, resultSet)
	assert.NotNil(t, err)

	mockFilter4 := series.NewMockFilter(ctrl)
	query, _ = sql.Parse("select f from cpu " +
		"where ip='1.1.1.1' or path='/data' and time>'20190410 00:00:00' and time<'20190410 10:00:00'")
	mockFilter4.EXPECT().
		FindSeriesIDsByExpr(uint32(1), &stmt.EqualsExpr{Key: "ip", Value: "1.1.1.1"}, query.TimeRange).
		Return(mockSeriesIDSet(series.Version(11), roaring.BitmapOf(1, 2, 3, 4)), nil)
	mockFilter4.EXPECT().
		FindSeriesIDsByExpr(uint32(1), &stmt.EqualsExpr{Key: "path", Value: "/data"}, query.TimeRange).
		Return(nil, errors.New("right error"))
	search = newSeriesSearch(1, mockFilter4, query)
	resultSet, err = search.Search()
	assert.Nil(t, resultSet)
	assert.NotNil(t, err)
}

func TestComplexCondition(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockFilter := series.NewMockFilter(ctrl)

	query, _ := sql.Parse("select f from cpu" +
		" where (ip not in ('1.1.1.1','2.2.2.2') and region='sh') and (path='/data' or path='/home')")
	mockFilter.EXPECT().
		FindSeriesIDsByExpr(gomock.Any(), &stmt.InExpr{Key: "ip", Values: []string{"1.1.1.1", "2.2.2.2"}}, query.TimeRange).
		Return(mockSeriesIDSet(series.Version(11), roaring.BitmapOf(1, 2, 4)), nil)
	mockFilter.EXPECT().
		GetSeriesIDsForTag(gomock.Any(), query.TimeRange).
		Return(mockSeriesIDSet(series.Version(11), roaring.BitmapOf(1, 2, 3, 4, 6, 7, 8)), nil)
	mockFilter.EXPECT().
		FindSeriesIDsByExpr(gomock.Any(), &stmt.EqualsExpr{Key: "region", Value: "sh"}, query.TimeRange).
		Return(mockSeriesIDSet(series.Version(11), roaring.BitmapOf(2, 3, 4, 7)), nil)
	mockFilter.EXPECT().
		FindSeriesIDsByExpr(gomock.Any(), &stmt.EqualsExpr{Key: "path", Value: "/data"}, query.TimeRange).
		Return(mockSeriesIDSet(series.Version(11), roaring.BitmapOf(3, 5)), nil)
	mockFilter.EXPECT().
		FindSeriesIDsByExpr(gomock.Any(), &stmt.EqualsExpr{Key: "path", Value: "/home"}, query.TimeRange).
		Return(mockSeriesIDSet(series.Version(11), roaring.BitmapOf(1)), nil)
	search := newSeriesSearch(10, mockFilter, query)
	resultSet, _ := search.Search()
	// ip not in ('1.1.1.1','2.2.2.2') => 3,6,7,8
	// ip not in ('1.1.1.1','2.2.2.2') and region='sh' => 3,7
	// path='/data' or path='/home' => 1,3,5
	// final => 3
	assert.Equal(t, *mockSeriesIDSet(series.Version(11), roaring.BitmapOf(3)), *resultSet)

	// error
	mockFilter1 := series.NewMockFilter(ctrl)
	mockFilter1.EXPECT().
		FindSeriesIDsByExpr(gomock.Any(), &stmt.InExpr{Key: "ip", Values: []string{"1.1.1.1", "2.2.2.2"}}, query.TimeRange).
		Return(mockSeriesIDSet(series.Version(11), roaring.BitmapOf(1, 2, 4)), nil)
	mockFilter1.EXPECT().
		GetSeriesIDsForTag(gomock.Any(), query.TimeRange).
		Return(mockSeriesIDSet(series.Version(11), roaring.BitmapOf(1, 2, 3, 4, 6, 7, 8)), nil)
	mockFilter1.EXPECT().
		FindSeriesIDsByExpr(gomock.Any(), &stmt.EqualsExpr{Key: "region", Value: "sh"}, query.TimeRange).
		Return(nil, errors.New("complex error"))
	search = newSeriesSearch(10, mockFilter1, query)
	resultSet, err := search.Search()
	assert.Nil(t, resultSet)
	assert.NotNil(t, err)
}

func TestSeriesSearch_condition_fail(t *testing.T) {
	search := newSeriesSearch(10, nil, nil)
	result, _ := search.findSeriesIDsByExpr(nil)
	assert.Nil(t, result)

	search = newSeriesSearch(10, nil, nil)
	result, _ = search.findSeriesIDsByExpr(&stmt.BinaryExpr{Operator: stmt.ADD})
	assert.Nil(t, result)

	query, _ := sql.Parse("select f from disk " +
		"where (ip='1.1.1.1')")
	search = newSeriesSearch(10, nil, query)
	search.err = fmt.Errorf("err")
	resultSet, err := search.Search()
	assert.Nil(t, resultSet)
	assert.NotNil(t, err)
}

func mockSeriesIDSet(version series.Version, ids *roaring.Bitmap) *series.MultiVerSeriesIDSet {
	s := series.NewMultiVerSeriesIDSet()
	s.Add(version, ids)
	return s
}
