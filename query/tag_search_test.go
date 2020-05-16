package query

import (
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/lindb/roaring"
	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/sql"
	"github.com/lindb/lindb/sql/stmt"
	"github.com/lindb/lindb/tsdb/metadb"
)

func TestTagSearch_Filter(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tagMeta := metadb.NewMockTagMetadata(ctrl)
	metadata := metadb.NewMockMetadata(ctrl)
	metadataDB := metadb.NewMockMetadataDatabase(ctrl)
	metadataDB.EXPECT().GetTagKeyID(gomock.Any(), gomock.Any(), gomock.Any()).Return(uint32(1), nil).AnyTimes()
	metadata.EXPECT().TagMetadata().Return(tagMeta).AnyTimes()
	metadata.EXPECT().MetadataDatabase().Return(metadataDB).AnyTimes()
	tagValueIDs := roaring.BitmapOf(1, 2, 3)

	// case 1: condition is empty
	q, _ := sql.Parse("select f from cpu")
	query := q.(*stmt.Query)
	search := newTagSearch("ns", query, metadata)
	resultSet, err := search.Filter()
	assert.NoError(t, err)
	assert.Empty(t, resultSet)
	// case 2: equal tag filter
	q, _ = sql.Parse("select f from cpu where ip='1.1.1.1'")
	query = q.(*stmt.Query)
	tagMeta.EXPECT().FindTagValueDsByExpr(gomock.Any(), &stmt.EqualsExpr{Key: "ip", Value: "1.1.1.1"}).Return(tagValueIDs, nil)
	search = newTagSearch("ns", query, metadata)
	resultSet, err = search.Filter()
	assert.NoError(t, err)
	assert.Len(t, resultSet, 1)
	assert.Equal(t, tagValueIDs, resultSet[(&stmt.EqualsExpr{Key: "ip", Value: "1.1.1.1"}).Rewrite()].tagValueIDs)
	// case 3: not tag filter
	q, _ = sql.Parse("select f from cpu where ip!='1.1.1.1'")
	query = q.(*stmt.Query)
	tagMeta.EXPECT().FindTagValueDsByExpr(gomock.Any(), &stmt.EqualsExpr{Key: "ip", Value: "1.1.1.1"}).Return(tagValueIDs, nil)
	search = newTagSearch("ns", query, metadata)
	resultSet, err = search.Filter()
	assert.NoError(t, err)
	assert.Len(t, resultSet, 1)
	assert.Equal(t, tagValueIDs, resultSet[(&stmt.EqualsExpr{Key: "ip", Value: "1.1.1.1"}).Rewrite()].tagValueIDs)
	// case 4: paren expr
	q, _ = sql.Parse("select f from cpu where (ip!='1.1.1.1')")
	query = q.(*stmt.Query)
	tagMeta.EXPECT().FindTagValueDsByExpr(gomock.Any(), &stmt.EqualsExpr{Key: "ip", Value: "1.1.1.1"}).Return(tagValueIDs, nil)
	search = newTagSearch("ns", query, metadata)
	resultSet, err = search.Filter()
	assert.NoError(t, err)
	assert.Len(t, resultSet, 1)
	assert.Equal(t, tagValueIDs, resultSet[(&stmt.EqualsExpr{Key: "ip", Value: "1.1.1.1"}).Rewrite()].tagValueIDs)
	// case 5: binary expr
	q, _ = sql.Parse("select f from cpu " +
		"where ip='1.1.1.1' and path='/data' and time>'20190410 00:00:00' and time<'20190410 10:00:00'")
	query = q.(*stmt.Query)
	tagMeta.EXPECT().FindTagValueDsByExpr(gomock.Any(), &stmt.EqualsExpr{Key: "ip", Value: "1.1.1.1"}).Return(tagValueIDs, nil)
	tagMeta.EXPECT().FindTagValueDsByExpr(gomock.Any(), &stmt.EqualsExpr{Key: "path", Value: "/data"}).Return(roaring.BitmapOf(10, 20), nil)
	search = newTagSearch("ns", query, metadata)
	resultSet, err = search.Filter()
	assert.NoError(t, err)
	assert.Len(t, resultSet, 2)
	assert.Equal(t, tagValueIDs, resultSet[(&stmt.EqualsExpr{Key: "ip", Value: "1.1.1.1"}).Rewrite()].tagValueIDs)
	assert.Equal(t, roaring.BitmapOf(10, 20), resultSet[(&stmt.EqualsExpr{Key: "path", Value: "/data"}).Rewrite()].tagValueIDs)
	// case 6: filter get empty
	q, _ = sql.Parse("select f from cpu where ip='1.1.1.1'")
	query = q.(*stmt.Query)
	tagMeta.EXPECT().FindTagValueDsByExpr(gomock.Any(), &stmt.EqualsExpr{Key: "ip", Value: "1.1.1.1"}).Return(nil, nil)
	search = newTagSearch("ns", query, metadata)
	resultSet, err = search.Filter()
	assert.NoError(t, err)
	assert.Len(t, resultSet, 0)
}

func TestTagSearch_Filter_err(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tagMeta := metadb.NewMockTagMetadata(ctrl)
	metadata := metadb.NewMockMetadata(ctrl)
	metadataDB := metadb.NewMockMetadataDatabase(ctrl)
	metadata.EXPECT().TagMetadata().Return(tagMeta).AnyTimes()
	metadata.EXPECT().MetadataDatabase().Return(metadataDB).AnyTimes()

	// case 1: get tag key err
	q, _ := sql.Parse("select f from cpu where ip='1.1.1.1'")
	query := q.(*stmt.Query)
	search := newTagSearch("ns", query, metadata)
	metadataDB.EXPECT().GetTagKeyID(gomock.Any(), gomock.Any(), gomock.Any()).Return(uint32(1), fmt.Errorf("err"))
	resultSet, err := search.Filter()
	assert.Error(t, err)
	assert.Nil(t, resultSet)
	// case 2: get tag value ids err
	search = newTagSearch("ns", query, metadata)
	metadataDB.EXPECT().GetTagKeyID(gomock.Any(), gomock.Any(), gomock.Any()).Return(uint32(1), nil).AnyTimes()
	tagMeta.EXPECT().FindTagValueDsByExpr(gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("err"))
	resultSet, err = search.Filter()
	assert.Error(t, err)
	assert.Nil(t, resultSet)
	// case 3: binary operator err
	q, _ = sql.Parse("select f from cpu " +
		"where ip='1.1.1.1' and path='/data'")
	query = q.(*stmt.Query)
	binary := query.Condition.(*stmt.BinaryExpr)
	binary.Operator = stmt.ADD
	search = newTagSearch("ns", query, metadata)
	resultSet, err = search.Filter()
	assert.Error(t, err)
	assert.Nil(t, resultSet)
	// case 4: recursion err
	q, _ = sql.Parse("select f from cpu where ip='1.1.1.1' or ip='1.1.1.1'")
	query = q.(*stmt.Query)
	search = newTagSearch("ns", query, metadata)
	metadataDB.EXPECT().GetTagKeyID(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(uint32(1), fmt.Errorf("err")).AnyTimes()
	tagMeta.EXPECT().FindTagValueDsByExpr(gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("err")).AnyTimes()
	resultSet, err = search.Filter()
	assert.Error(t, err)
	assert.Nil(t, resultSet)
}

func TestTagSearch_Filter_Complex(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tagMeta := metadb.NewMockTagMetadata(ctrl)
	metadata := metadb.NewMockMetadata(ctrl)
	metadataDB := metadb.NewMockMetadataDatabase(ctrl)
	metadataDB.EXPECT().GetTagKeyID(gomock.Any(), gomock.Any(), gomock.Any()).Return(uint32(1), nil).AnyTimes()
	metadata.EXPECT().TagMetadata().Return(tagMeta).AnyTimes()
	metadata.EXPECT().MetadataDatabase().Return(metadataDB).AnyTimes()
	tagValueIDs := roaring.BitmapOf(1, 2, 3)

	q, _ := sql.Parse("select f from cpu" +
		" where (ip not in ('1.1.1.1','2.2.2.2') and region='sh') and (path='/data' or path='/home')")
	query := q.(*stmt.Query)
	search := newTagSearch("ns", query, metadata)
	tagMeta.EXPECT().FindTagValueDsByExpr(gomock.Any(), &stmt.InExpr{Key: "ip", Values: []string{"1.1.1.1", "2.2.2.2"}}).
		Return(tagValueIDs, nil)
	tagMeta.EXPECT().FindTagValueDsByExpr(gomock.Any(), &stmt.EqualsExpr{Key: "region", Value: "sh"}).
		Return(tagValueIDs, nil)
	tagMeta.EXPECT().FindTagValueDsByExpr(gomock.Any(), &stmt.EqualsExpr{Key: "path", Value: "/data"}).
		Return(tagValueIDs, nil)
	tagMeta.EXPECT().FindTagValueDsByExpr(gomock.Any(), &stmt.EqualsExpr{Key: "path", Value: "/home"}).
		Return(tagValueIDs, nil)
	resultSet, err := search.Filter()
	assert.NoError(t, err)
	assert.NotNil(t, resultSet)
	assert.Len(t, resultSet, 4)
	assert.Equal(t, tagValueIDs, resultSet[(&stmt.InExpr{Key: "ip", Values: []string{"1.1.1.1", "2.2.2.2"}}).Rewrite()].tagValueIDs)
	assert.Equal(t, tagValueIDs, resultSet[(&stmt.EqualsExpr{Key: "region", Value: "sh"}).Rewrite()].tagValueIDs)
	assert.Equal(t, tagValueIDs, resultSet[(&stmt.EqualsExpr{Key: "path", Value: "/data"}).Rewrite()].tagValueIDs)
	assert.Equal(t, tagValueIDs, resultSet[(&stmt.EqualsExpr{Key: "path", Value: "/home"}).Rewrite()].tagValueIDs)
}
