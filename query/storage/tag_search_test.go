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
	"github.com/lindb/lindb/series/tag"
	"github.com/lindb/lindb/sql"
	"github.com/lindb/lindb/sql/stmt"
	"github.com/lindb/lindb/tsdb"
	"github.com/lindb/lindb/tsdb/metadb"
)

func TestTagSearch_Filter(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tagMeta := metadb.NewMockTagMetadata(ctrl)
	metadata := metadb.NewMockMetadata(ctrl)
	db := tsdb.NewMockDatabase(ctrl)
	db.EXPECT().Metadata().Return(metadata).AnyTimes()
	metadataDB := metadb.NewMockMetadataDatabase(ctrl)
	metadataDB.EXPECT().GetTagKeyID(gomock.Any(), gomock.Any(), gomock.Any()).Return(tag.KeyID(1), nil).AnyTimes()
	metadata.EXPECT().TagMetadata().Return(tagMeta).AnyTimes()
	metadata.EXPECT().MetadataDatabase().Return(metadataDB).AnyTimes()
	tagValueIDs := roaring.BitmapOf(1, 2, 3)

	// case 1: condition is empty
	q, _ := sql.Parse("select f from cpu")
	query := q.(*stmt.Query)
	ctx := &executeContext{
		database: db,
		storageExecuteCtx: &flow.StorageExecuteContext{
			Query: query,
		},
	}
	search := newTagSearch(ctx)
	err := search.Filter()
	assert.NoError(t, err)
	assert.Empty(t, ctx.storageExecuteCtx.TagFilterResult)
	// case 2: equal tag filter
	q, _ = sql.Parse("select f from cpu where ip='1.1.1.1'")
	query = q.(*stmt.Query)
	tagMeta.EXPECT().FindTagValueDsByExpr(gomock.Any(), &stmt.EqualsExpr{Key: "ip", Value: "1.1.1.1"}).Return(tagValueIDs, nil)
	ctx.storageExecuteCtx.Query = query
	search = newTagSearch(ctx)
	err = search.Filter()
	assert.NoError(t, err)
	assert.Len(t, ctx.storageExecuteCtx.TagFilterResult, 1)
	assert.Equal(t, tagValueIDs,
		ctx.storageExecuteCtx.TagFilterResult[(&stmt.EqualsExpr{Key: "ip", Value: "1.1.1.1"}).Rewrite()].TagValueIDs)
	// case 3: not tag filter
	q, _ = sql.Parse("select f from cpu where ip!='1.1.1.1'")
	query = q.(*stmt.Query)
	tagMeta.EXPECT().FindTagValueDsByExpr(gomock.Any(), &stmt.EqualsExpr{Key: "ip", Value: "1.1.1.1"}).Return(tagValueIDs, nil)
	ctx.storageExecuteCtx.Query = query
	search = newTagSearch(ctx)
	err = search.Filter()
	assert.NoError(t, err)
	assert.Len(t, ctx.storageExecuteCtx.TagFilterResult, 1)
	assert.Equal(t, tagValueIDs,
		ctx.storageExecuteCtx.TagFilterResult[(&stmt.EqualsExpr{Key: "ip", Value: "1.1.1.1"}).Rewrite()].TagValueIDs)
	// case 4: paren expr
	q, _ = sql.Parse("select f from cpu where (ip!='1.1.1.1')")
	query = q.(*stmt.Query)
	tagMeta.EXPECT().FindTagValueDsByExpr(gomock.Any(), &stmt.EqualsExpr{Key: "ip", Value: "1.1.1.1"}).Return(tagValueIDs, nil)
	ctx.storageExecuteCtx.Query = query
	search = newTagSearch(ctx)
	err = search.Filter()
	assert.NoError(t, err)
	assert.Len(t, ctx.storageExecuteCtx.TagFilterResult, 1)
	assert.Equal(t, tagValueIDs,
		ctx.storageExecuteCtx.TagFilterResult[(&stmt.EqualsExpr{Key: "ip", Value: "1.1.1.1"}).Rewrite()].TagValueIDs)
	// case 5: binary expr
	q, _ = sql.Parse("select f from cpu " +
		"where ip='1.1.1.1' and path='/data' and time>'20190410 00:00:00' and time<'20190410 10:00:00'")
	query = q.(*stmt.Query)
	tagMeta.EXPECT().FindTagValueDsByExpr(gomock.Any(), &stmt.EqualsExpr{Key: "ip", Value: "1.1.1.1"}).Return(tagValueIDs, nil)
	tagMeta.EXPECT().FindTagValueDsByExpr(gomock.Any(), &stmt.EqualsExpr{Key: "path", Value: "/data"}).Return(roaring.BitmapOf(10, 20), nil)
	ctx.storageExecuteCtx.Query = query
	search = newTagSearch(ctx)
	err = search.Filter()
	assert.NoError(t, err)
	assert.Len(t, ctx.storageExecuteCtx.TagFilterResult, 2)
	assert.Equal(t, tagValueIDs,
		ctx.storageExecuteCtx.TagFilterResult[(&stmt.EqualsExpr{Key: "ip", Value: "1.1.1.1"}).Rewrite()].TagValueIDs)
	assert.Equal(t, roaring.BitmapOf(10, 20),
		ctx.storageExecuteCtx.TagFilterResult[(&stmt.EqualsExpr{Key: "path", Value: "/data"}).Rewrite()].TagValueIDs)
	// case 6: filter get empty
	q, _ = sql.Parse("select f from cpu where ip='1.1.1.1'")
	query = q.(*stmt.Query)
	tagMeta.EXPECT().FindTagValueDsByExpr(gomock.Any(), &stmt.EqualsExpr{Key: "ip", Value: "1.1.1.1"}).Return(nil, nil)
	ctx.storageExecuteCtx.Query = query
	search = newTagSearch(ctx)
	err = search.Filter()
	assert.NoError(t, err)
	assert.Len(t, ctx.storageExecuteCtx.TagFilterResult, 0)
}

func TestTagSearch_Filter_err(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tagMeta := metadb.NewMockTagMetadata(ctrl)
	metadata := metadb.NewMockMetadata(ctrl)
	db := tsdb.NewMockDatabase(ctrl)
	db.EXPECT().Metadata().Return(metadata).AnyTimes()
	metadataDB := metadb.NewMockMetadataDatabase(ctrl)
	metadata.EXPECT().TagMetadata().Return(tagMeta).AnyTimes()
	metadata.EXPECT().MetadataDatabase().Return(metadataDB).AnyTimes()

	// case 1: get tag key err
	q, _ := sql.Parse("select f from cpu where ip='1.1.1.1'")
	query := q.(*stmt.Query)
	ctx := &executeContext{
		database: db,
		storageExecuteCtx: &flow.StorageExecuteContext{
			Query: query,
		},
	}
	search := newTagSearch(ctx)
	metadataDB.EXPECT().GetTagKeyID(gomock.Any(), gomock.Any(), gomock.Any()).Return(tag.KeyID(1), fmt.Errorf("err"))
	err := search.Filter()
	resultSet := ctx.storageExecuteCtx.TagFilterResult
	assert.Error(t, err)
	assert.Empty(t, resultSet)
	// case 2: get tag value ids err
	search = newTagSearch(ctx)
	metadataDB.EXPECT().GetTagKeyID(gomock.Any(), gomock.Any(), gomock.Any()).Return(tag.KeyID(1), nil).AnyTimes()
	tagMeta.EXPECT().FindTagValueDsByExpr(gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("err"))
	err = search.Filter()
	resultSet = ctx.storageExecuteCtx.TagFilterResult
	assert.Error(t, err)
	assert.Empty(t, resultSet)
	// case 3: binary operator err
	q, _ = sql.Parse("select f from cpu " +
		"where ip='1.1.1.1' and path='/data'")
	query = q.(*stmt.Query)
	ctx.storageExecuteCtx.Query = query
	binary := query.Condition.(*stmt.BinaryExpr)
	binary.Operator = stmt.ADD
	search = newTagSearch(ctx)
	err = search.Filter()
	resultSet = ctx.storageExecuteCtx.TagFilterResult
	assert.Error(t, err)
	assert.Empty(t, resultSet)
	// case 4: recursion err
	q, _ = sql.Parse("select f from cpu where ip='1.1.1.1' or ip='1.1.1.1'")
	query = q.(*stmt.Query)
	ctx.storageExecuteCtx.Query = query
	search = newTagSearch(ctx)
	metadataDB.EXPECT().GetTagKeyID(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(tag.KeyID(1), fmt.Errorf("err")).AnyTimes()
	tagMeta.EXPECT().FindTagValueDsByExpr(gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("err")).AnyTimes()
	err = search.Filter()
	resultSet = ctx.storageExecuteCtx.TagFilterResult
	assert.Error(t, err)
	assert.Empty(t, resultSet)
}

func TestTagSearch_Filter_Complex(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tagMeta := metadb.NewMockTagMetadata(ctrl)
	metadata := metadb.NewMockMetadata(ctrl)
	db := tsdb.NewMockDatabase(ctrl)
	db.EXPECT().Metadata().Return(metadata).AnyTimes()
	metadataDB := metadb.NewMockMetadataDatabase(ctrl)
	metadataDB.EXPECT().GetTagKeyID(gomock.Any(), gomock.Any(), gomock.Any()).Return(tag.KeyID(1), nil).AnyTimes()
	metadata.EXPECT().TagMetadata().Return(tagMeta).AnyTimes()
	metadata.EXPECT().MetadataDatabase().Return(metadataDB).AnyTimes()
	tagValueIDs := roaring.BitmapOf(1, 2, 3)

	q, _ := sql.Parse("select f from cpu" +
		" where (ip not in ('1.1.1.1','2.2.2.2') and region='sh') and (path='/data' or path='/home')")
	query := q.(*stmt.Query)
	ctx := &executeContext{
		database: db,
		storageExecuteCtx: &flow.StorageExecuteContext{
			Query: query,
		},
	}
	search := newTagSearch(ctx)
	tagMeta.EXPECT().FindTagValueDsByExpr(gomock.Any(), &stmt.InExpr{Key: "ip", Values: []string{"1.1.1.1", "2.2.2.2"}}).
		Return(tagValueIDs, nil)
	tagMeta.EXPECT().FindTagValueDsByExpr(gomock.Any(), &stmt.EqualsExpr{Key: "region", Value: "sh"}).
		Return(tagValueIDs, nil)
	tagMeta.EXPECT().FindTagValueDsByExpr(gomock.Any(), &stmt.EqualsExpr{Key: "path", Value: "/data"}).
		Return(tagValueIDs, nil)
	tagMeta.EXPECT().FindTagValueDsByExpr(gomock.Any(), &stmt.EqualsExpr{Key: "path", Value: "/home"}).
		Return(tagValueIDs, nil)
	err := search.Filter()
	resultSet := ctx.storageExecuteCtx.TagFilterResult
	assert.NoError(t, err)
	assert.NotNil(t, resultSet)
	assert.Len(t, resultSet, 4)
	assert.Equal(t, tagValueIDs, resultSet[(&stmt.InExpr{Key: "ip", Values: []string{"1.1.1.1", "2.2.2.2"}}).Rewrite()].TagValueIDs)
	assert.Equal(t, tagValueIDs, resultSet[(&stmt.EqualsExpr{Key: "region", Value: "sh"}).Rewrite()].TagValueIDs)
	assert.Equal(t, tagValueIDs, resultSet[(&stmt.EqualsExpr{Key: "path", Value: "/data"}).Rewrite()].TagValueIDs)
	assert.Equal(t, tagValueIDs, resultSet[(&stmt.EqualsExpr{Key: "path", Value: "/home"}).Rewrite()].TagValueIDs)
}
