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

package storagequery

import (
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/lindb/roaring"
	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/flow"
	"github.com/lindb/lindb/series"
	"github.com/lindb/lindb/series/tag"
	"github.com/lindb/lindb/sql"
	"github.com/lindb/lindb/sql/stmt"
)

func TestSeriesSearch_Search(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockFilter := series.NewMockFilter(ctrl)
	seriesIDs := roaring.BitmapOf(10, 20, 30)

	// case 1: empty filter expr
	q, _ := sql.Parse("select f from cpu")
	query := q.(*stmt.Query)
	search := newSeriesSearch(mockFilter, nil, query.Condition)
	resultSet, err := search.Search()
	assert.NoError(t, err)
	assert.Equal(t, uint64(0), resultSet.GetCardinality())
	// case 2: equal tag filter
	q, _ = sql.Parse("select f from cpu where ip='1.1.1.1'")
	query = q.(*stmt.Query)
	mockFilter.EXPECT().GetSeriesIDsByTagValueIDs(tag.KeyID(1), gomock.Any()).Return(seriesIDs.Clone(), nil)
	search = newSeriesSearch(mockFilter, mockFilterResult(), query.Condition)
	resultSet, err = search.Search()
	assert.NoError(t, err)
	assert.Equal(t, seriesIDs, resultSet)
	// case 3: not expr
	q, _ = sql.Parse("select f from cpu where ip!='1.1.1.1'")
	query = q.(*stmt.Query)
	mockFilter.EXPECT().GetSeriesIDsByTagValueIDs(tag.KeyID(1), gomock.Any()).Return(seriesIDs.Clone(), nil)
	mockFilter.EXPECT().GetSeriesIDsForTag(tag.KeyID(1)).Return(roaring.BitmapOf(10, 20, 40, 50), nil)
	search = newSeriesSearch(mockFilter, mockFilterResult(), query.Condition)
	resultSet, err = search.Search()
	assert.NoError(t, err)
	assert.Equal(t, roaring.BitmapOf(40, 50), resultSet)
	// case 4: binary expr and
	q, _ = sql.Parse("select f from cpu " +
		"where ip='1.1.1.1' and path='/data' and time>'20190410 00:00:00' and time<'20190410 10:00:00'")
	query = q.(*stmt.Query)
	mockFilter.EXPECT().GetSeriesIDsByTagValueIDs(tag.KeyID(1), gomock.Any()).Return(seriesIDs.Clone(), nil)
	mockFilter.EXPECT().GetSeriesIDsByTagValueIDs(tag.KeyID(2), gomock.Any()).Return(roaring.BitmapOf(20), nil)
	search = newSeriesSearch(mockFilter, mockFilterResult(), query.Condition)
	resultSet, err = search.Search()
	assert.NoError(t, err)
	assert.Equal(t, roaring.BitmapOf(20), resultSet)
	// case 5: binary expr or
	q, _ = sql.Parse("select f from cpu " +
		"where ip='1.1.1.1' or path='/data' and time>'20190410 00:00:00' and time<'20190410 10:00:00'")
	query = q.(*stmt.Query)
	mockFilter.EXPECT().GetSeriesIDsByTagValueIDs(tag.KeyID(1), gomock.Any()).Return(seriesIDs.Clone(), nil)
	mockFilter.EXPECT().GetSeriesIDsByTagValueIDs(tag.KeyID(2), gomock.Any()).Return(roaring.BitmapOf(200), nil)
	search = newSeriesSearch(mockFilter, mockFilterResult(), query.Condition)
	resultSet, err = search.Search()
	assert.NoError(t, err)
	assert.Equal(t, roaring.BitmapOf(10, 20, 30, 200), resultSet)
	// case 6: paren expr
	q, _ = sql.Parse("select f from cpu where (ip='1.1.1.1')")
	query = q.(*stmt.Query)
	mockFilter.EXPECT().GetSeriesIDsByTagValueIDs(tag.KeyID(1), gomock.Any()).Return(seriesIDs.Clone(), nil)
	search = newSeriesSearch(mockFilter, mockFilterResult(), query.Condition)
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
	q, _ := sql.Parse("select f from cpu where ip='1.1.1.1'")
	query := q.(*stmt.Query)
	search := newSeriesSearch(mockFilter, make(map[string]*flow.TagFilterResult), query.Condition)
	resultSet, err := search.Search()
	assert.Error(t, err)
	assert.Nil(t, resultSet)
	// case 2: get series id err
	search = newSeriesSearch(mockFilter, mockFilterResult(), query.Condition)
	mockFilter.EXPECT().GetSeriesIDsByTagValueIDs(gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("err"))
	resultSet, err = search.Search()
	assert.Error(t, err)
	assert.Nil(t, resultSet)
	// case 3: not expr err
	q, _ = sql.Parse("select f from cpu where ip!='1.1.1.1'")
	query = q.(*stmt.Query)
	mockFilter.EXPECT().GetSeriesIDsByTagValueIDs(tag.KeyID(1), gomock.Any()).Return(seriesIDs, nil)
	mockFilter.EXPECT().GetSeriesIDsForTag(tag.KeyID(1)).Return(nil, fmt.Errorf("err"))
	search = newSeriesSearch(mockFilter, mockFilterResult(), query.Condition)
	resultSet, err = search.Search()
	assert.Error(t, err)
	assert.Nil(t, resultSet)
	// case 4: recursion err
	q, _ = sql.Parse("select f from cpu where ip='1.1.1.1' or ip='1.1.1.1'")
	query = q.(*stmt.Query)
	mockFilter.EXPECT().GetSeriesIDsByTagValueIDs(tag.KeyID(1), gomock.Any()).Return(nil, fmt.Errorf("err"))
	search = newSeriesSearch(mockFilter, mockFilterResult(), query.Condition)
	resultSet, err = search.Search()
	assert.Error(t, err)
	assert.Nil(t, resultSet)
}

func TestSeriesSearch_Search_expr_not_match(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockFilter := series.NewMockFilter(ctrl)

	q, _ := sql.Parse("select f from cpu where ip='1.1.1.1'")
	query := q.(*stmt.Query)
	query.Condition = &stmt.CallExpr{}
	search := newSeriesSearch(mockFilter, make(map[string]*flow.TagFilterResult), query.Condition)
	resultSet, err := search.Search()
	assert.NoError(t, err)
	assert.Equal(t, uint64(0), resultSet.GetCardinality())
}

func TestSeriesSearch_Search_complex(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockFilter := series.NewMockFilter(ctrl)

	q, _ := sql.Parse("select f from cpu" +
		" where (ip not in ('1.1.1.1','2.2.2.2') and region='sh') and (path='/data' or path='/home')")
	query := q.(*stmt.Query)
	mockFilter.EXPECT().GetSeriesIDsByTagValueIDs(tag.KeyID(1), roaring.BitmapOf(5)).Return(roaring.BitmapOf(1, 2), nil)
	mockFilter.EXPECT().GetSeriesIDsForTag(tag.KeyID(1)).Return(roaring.BitmapOf(1, 2, 3, 4, 5, 6, 7), nil)
	mockFilter.EXPECT().GetSeriesIDsByTagValueIDs(tag.KeyID(3), roaring.BitmapOf(4)).Return(roaring.BitmapOf(3, 5, 6, 7), nil)
	mockFilter.EXPECT().GetSeriesIDsByTagValueIDs(tag.KeyID(2), roaring.BitmapOf(2)).Return(roaring.BitmapOf(7), nil)
	mockFilter.EXPECT().GetSeriesIDsByTagValueIDs(tag.KeyID(2), roaring.BitmapOf(3)).Return(roaring.BitmapOf(5), nil)
	search := newSeriesSearch(mockFilter, mockFilterResult(), query.Condition)
	resultSet, err := search.Search()
	assert.NoError(t, err)
	assert.Equal(t, roaring.BitmapOf(5, 7), resultSet)
}

func mockFilterResult() map[string]*flow.TagFilterResult {
	result := make(map[string]*flow.TagFilterResult)
	result[(&stmt.EqualsExpr{Key: "ip", Value: "1.1.1.1"}).Rewrite()] = &flow.TagFilterResult{
		TagKeyID:    1,
		TagValueIDs: roaring.BitmapOf(1),
	}
	result[(&stmt.EqualsExpr{Key: "path", Value: "/data"}).Rewrite()] = &flow.TagFilterResult{
		TagKeyID:    2,
		TagValueIDs: roaring.BitmapOf(2),
	}
	result[(&stmt.EqualsExpr{Key: "path", Value: "/home"}).Rewrite()] = &flow.TagFilterResult{
		TagKeyID:    2,
		TagValueIDs: roaring.BitmapOf(3),
	}
	result[(&stmt.EqualsExpr{Key: "region", Value: "sh"}).Rewrite()] = &flow.TagFilterResult{
		TagKeyID:    3,
		TagValueIDs: roaring.BitmapOf(4),
	}
	result[(&stmt.InExpr{Key: "ip", Values: []string{"1.1.1.1", "2.2.2.2"}}).Rewrite()] = &flow.TagFilterResult{
		TagKeyID:    1,
		TagValueIDs: roaring.BitmapOf(5),
	}
	return result
}
