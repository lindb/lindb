package metadb

import (
	"context"
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/pkg/fileutil"
)

func TestNewMetadata(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer func() {
		_ = fileutil.RemoveDir(testPath)
		ctrl.Finish()
	}()
	metadata1, err := NewMetadata(context.TODO(), "test", testPath, nil)
	assert.NoError(t, err)
	assert.NotNil(t, metadata1.TagMetadata())
	assert.NotNil(t, metadata1.MetadataDatabase())

	metadata2, err := NewMetadata(context.TODO(), "test", testPath, nil)
	assert.Error(t, err)
	assert.Nil(t, metadata2)

	err = metadata1.Close()
	assert.NoError(t, err)

	db := NewMockMetadataDatabase(ctrl)
	m := metadata1.(*metadata)
	m.metadataDatabase = db
	db.EXPECT().Close().Return(fmt.Errorf("err"))
	err = m.Close()
	assert.Error(t, err)
}

func TestMetadata_Flush(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer func() {
		_ = fileutil.RemoveDir(testPath)
		ctrl.Finish()
	}()
	metadata1, err := NewMetadata(context.TODO(), "test", testPath, nil)
	assert.NoError(t, err)
	db := NewMockMetadataDatabase(ctrl)
	m := metadata1.(*metadata)
	m.metadataDatabase = db
	db.EXPECT().Sync().Return(fmt.Errorf("err"))
	err = metadata1.Flush()
	assert.Error(t, err)

	db.EXPECT().Sync().Return(nil)
	err = metadata1.Flush()
	assert.NoError(t, err)
}
