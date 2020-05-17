package query

import (
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/lindb/roaring"
	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/pkg/encoding"
	"github.com/lindb/lindb/series"
	"github.com/lindb/lindb/series/field"
	"github.com/lindb/lindb/sql/stmt"
	"github.com/lindb/lindb/tsdb"
	"github.com/lindb/lindb/tsdb/indexdb"
	"github.com/lindb/lindb/tsdb/metadb"
)

func TestMetadataStorageExecutor_Execute(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	db := tsdb.NewMockDatabase(ctrl)

	metadata := metadb.NewMockMetadata(ctrl)
	db.EXPECT().Metadata().Return(metadata).AnyTimes()
	metadataIndex := metadb.NewMockMetadataDatabase(ctrl)
	metadata.EXPECT().MetadataDatabase().Return(metadataIndex).AnyTimes()

	// case 1: suggest namespace
	exec := newMetadataStorageExecutor(db, nil, &stmt.Metadata{
		Type: stmt.Namespace,
	})
	metadataIndex.EXPECT().SuggestNamespace(gomock.Any(), gomock.Any()).Return([]string{"a"}, nil)
	result, err := exec.Execute()
	assert.NoError(t, err)
	assert.Equal(t, []string{"a"}, result)

	// case 2: suggest metric name
	exec = newMetadataStorageExecutor(db, nil, &stmt.Metadata{
		Type: stmt.Metric,
	})
	metadataIndex.EXPECT().SuggestMetrics(gomock.Any(), gomock.Any(), gomock.Any()).Return([]string{"a"}, nil)
	result, err = exec.Execute()
	assert.NoError(t, err)
	assert.Equal(t, []string{"a"}, result)

	// case 3: suggest tag keys
	exec = newMetadataStorageExecutor(db, nil, &stmt.Metadata{
		Type: stmt.TagKey,
	})
	metadataIndex.EXPECT().SuggestTagKeys(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return([]string{"a"}, nil)
	result, err = exec.Execute()
	assert.NoError(t, err)
	assert.Equal(t, []string{"a"}, result)
	// case 4: get fields err
	exec = newMetadataStorageExecutor(db, nil, &stmt.Metadata{
		Type: stmt.Field,
	})
	metadataIndex.EXPECT().GetAllFields(gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("err"))
	result, err = exec.Execute()
	assert.Error(t, err)
	assert.Empty(t, result)

	// case 5: get fields
	exec = newMetadataStorageExecutor(db, nil, &stmt.Metadata{
		Type: stmt.Field,
	})
	metadataIndex.EXPECT().GetAllFields(gomock.Any(), gomock.Any()).Return([]field.Meta{{ID: 10}}, nil)
	result, err = exec.Execute()
	assert.NoError(t, err)
	assert.Equal(t, string(encoding.JSONMarshal([]field.Meta{{ID: 10}})), result[0])

	// case 6: suggest tag values
	exec = newMetadataStorageExecutor(db, []int32{1, 2}, &stmt.Metadata{
		Type: stmt.TagValue,
	})
	metadataIndex.EXPECT().GetTagKeyID(gomock.Any(), gomock.Any(), gomock.Any()).Return(uint32(2), nil)

	tagMeta := metadb.NewMockTagMetadata(ctrl)
	metadata.EXPECT().TagMetadata().Return(tagMeta).AnyTimes()

	tagMeta.EXPECT().SuggestTagValues(gomock.Any(), gomock.Any(), gomock.Any()).Return([]string{"a"})
	result, err = exec.Execute()
	assert.NoError(t, err)
	assert.Equal(t, []string{"a"}, result)

	// case 7: suggest tag values err
	exec = newMetadataStorageExecutor(db, []int32{1, 2}, &stmt.Metadata{
		Type: stmt.TagValue,
	})
	metadataIndex.EXPECT().GetTagKeyID(gomock.Any(), gomock.Any(), gomock.Any()).Return(uint32(0), fmt.Errorf("err"))

	result, err = exec.Execute()
	assert.Error(t, err)
	assert.Empty(t, result)
}

func TestMetadataStorageExecutor_Execute_With_Tag_Condition(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer func() {
		newTagSearchFunc = newTagSearch
		newSeriesSearchFunc = newSeriesSearch

		ctrl.Finish()
	}()

	db := tsdb.NewMockDatabase(ctrl)

	metadata := metadb.NewMockMetadata(ctrl)
	db.EXPECT().Metadata().Return(metadata).AnyTimes()
	metadataIndex := metadb.NewMockMetadataDatabase(ctrl)
	metadata.EXPECT().MetadataDatabase().Return(metadataIndex).AnyTimes()
	metadataIndex.EXPECT().GetTagKeyID(gomock.Any(), gomock.Any(), gomock.Any()).Return(uint32(2), nil).AnyTimes()

	// case 1: tag search err
	tagSearch := NewMockTagSearch(ctrl)
	newTagSearchFunc = func(namespace, metricName string, condition stmt.Expr, metadata metadb.Metadata) TagSearch {
		return tagSearch
	}
	exec := newMetadataStorageExecutor(db, []int32{1}, &stmt.Metadata{
		Type:      stmt.TagValue,
		Condition: &stmt.EqualsExpr{},
		Limit:     2,
	})
	tagSearch.EXPECT().Filter().Return(nil, fmt.Errorf("err"))
	_, err := exec.Execute()
	assert.Error(t, err)
	// case 2: tag not found
	tagSearch.EXPECT().Filter().Return(nil, nil)
	_, err = exec.Execute()
	assert.Error(t, err)

	shard := tsdb.NewMockShard(ctrl)
	db.EXPECT().GetShard(gomock.Any()).Return(shard, true).AnyTimes()
	indexDB := indexdb.NewMockIndexDatabase(ctrl)
	shard.EXPECT().IndexDatabase().Return(indexDB).AnyTimes()

	tagSearch.EXPECT().Filter().Return(map[string]*tagFilterResult{"key": {}}, nil).AnyTimes()
	// case 3: series search err
	seriesSearch := NewMockSeriesSearch(ctrl)
	newSeriesSearchFunc = func(filter series.Filter, filterResult map[string]*tagFilterResult, condition stmt.Expr) SeriesSearch {
		return seriesSearch
	}
	seriesSearch.EXPECT().Search().Return(nil, fmt.Errorf("err"))
	_, err = exec.Execute()
	assert.Error(t, err)

	seriesSearch.EXPECT().Search().Return(roaring.BitmapOf(1, 2, 3), nil).AnyTimes()
	// case 4: get grouping err
	indexDB.EXPECT().GetGroupingContext(gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("err"))
	_, err = exec.Execute()
	assert.Error(t, err)

	gCtx := series.NewMockGroupingContext(ctrl)
	indexDB.EXPECT().GetGroupingContext(gomock.Any(), gomock.Any()).Return(gCtx, nil).AnyTimes()
	gCtx.EXPECT().ScanTagValueIDs(gomock.Any(), gomock.Any()).
		Return([]*roaring.Bitmap{roaring.BitmapOf(1, 2, 3)}).AnyTimes()
	tagMeta := metadb.NewMockTagMetadata(ctrl)
	metadata.EXPECT().TagMetadata().Return(tagMeta).AnyTimes()

	// case 5: collect tag value err
	tagMeta.EXPECT().CollectTagValues(gomock.Any(), gomock.Any(), gomock.Any()).Return(fmt.Errorf("err"))
	_, err = exec.Execute()
	assert.Error(t, err)

	// case 6: collect tag values
	tagMeta.EXPECT().CollectTagValues(gomock.Any(), gomock.Any(), gomock.Any()).
		DoAndReturn(func(tagKeyID uint32,
			tagValueIDs *roaring.Bitmap,
			tagValues map[uint32]string,
		) error {
			tagValues[12] = "a"
			tagValues[13] = "b"
			tagValues[14] = "c"
			tagValues[15] = "d"
			return nil
		})
	result, err := exec.Execute()
	assert.NoError(t, err)
	assert.Len(t, result, 2)
}
