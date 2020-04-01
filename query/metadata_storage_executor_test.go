package query

import (
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/pkg/encoding"
	"github.com/lindb/lindb/series/field"
	"github.com/lindb/lindb/sql/stmt"
	"github.com/lindb/lindb/tsdb"
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
	exec := newMetadataStorageExecutor(db, "ns", nil, &stmt.Metadata{
		Type: stmt.Namespace,
	})
	metadataIndex.EXPECT().SuggestNamespace(gomock.Any(), gomock.Any()).Return([]string{"a"}, nil)
	result, err := exec.Execute()
	assert.NoError(t, err)
	assert.Equal(t, []string{"a"}, result)

	// case 2: suggest metric name
	exec = newMetadataStorageExecutor(db, "ns", nil, &stmt.Metadata{
		Type: stmt.Metric,
	})
	metadataIndex.EXPECT().SuggestMetrics(gomock.Any(), gomock.Any(), gomock.Any()).Return([]string{"a"}, nil)
	result, err = exec.Execute()
	assert.NoError(t, err)
	assert.Equal(t, []string{"a"}, result)

	// case 3: suggest tag keys
	exec = newMetadataStorageExecutor(db, "ns", nil, &stmt.Metadata{
		Type: stmt.TagKey,
	})
	metadataIndex.EXPECT().SuggestTagKeys(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return([]string{"a"}, nil)
	result, err = exec.Execute()
	assert.NoError(t, err)
	assert.Equal(t, []string{"a"}, result)
	// case 4: get fields err
	exec = newMetadataStorageExecutor(db, "ns", nil, &stmt.Metadata{
		Type: stmt.Field,
	})
	metadataIndex.EXPECT().GetAllFields(gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("err"))
	result, err = exec.Execute()
	assert.Error(t, err)
	assert.Empty(t, result)

	// case 5: get fields
	exec = newMetadataStorageExecutor(db, "ns", nil, &stmt.Metadata{
		Type: stmt.Field,
	})
	metadataIndex.EXPECT().GetAllFields(gomock.Any(), gomock.Any()).Return([]field.Meta{{ID: 10}}, nil)
	result, err = exec.Execute()
	assert.NoError(t, err)
	assert.Equal(t, string(encoding.JSONMarshal([]field.Meta{{ID: 10}})), result[0])

	// case 6: suggest tag values
	exec = newMetadataStorageExecutor(db, "ns", []int32{1, 2}, &stmt.Metadata{
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
	exec = newMetadataStorageExecutor(db, "ns", []int32{1, 2}, &stmt.Metadata{
		Type: stmt.TagValue,
	})
	metadataIndex.EXPECT().GetTagKeyID(gomock.Any(), gomock.Any(), gomock.Any()).Return(uint32(0), fmt.Errorf("err"))

	result, err = exec.Execute()
	assert.Error(t, err)
	assert.Empty(t, result)
}
