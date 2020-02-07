package metadb

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/pkg/fileutil"
)

func TestNewMetadata(t *testing.T) {
	defer func() {
		_ = fileutil.RemoveDir(testPath)
	}()
	metadata, err := NewMetadata(context.TODO(), "test", testPath, nil)
	assert.NoError(t, err)
	assert.NotNil(t, metadata.TagMetadata())
	assert.NotNil(t, metadata.MetadataDatabase())

	metadata, err = NewMetadata(context.TODO(), "test", testPath, nil)
	assert.Error(t, err)
	assert.Nil(t, metadata)
}
