package table

import (
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/pkg/fileutil"
)

func TestMapCache_GetReader(t *testing.T) {
	_ = fileutil.MkDirIfNotExist(testKVPath)
	ctrl := gomock.NewController(t)
	defer func() {
		newMMapStoreReaderFunc = newMMapStoreReader
		_ = fileutil.RemoveDir(testKVPath)
		ctrl.Finish()
	}()
	cache := NewCache(testKVPath)
	// case 1: get reader err
	newMMapStoreReaderFunc = func(path string) (r Reader, err error) {
		return nil, fmt.Errorf("err")
	}
	r, err := cache.GetReader("f", "100000.sst")
	assert.Error(t, err)
	assert.Nil(t, r)
	// case 2: get reader success
	mockReader := NewMockReader(ctrl)
	newMMapStoreReaderFunc = func(path string) (reader Reader, err error) {
		return mockReader, nil
	}
	r, err = cache.GetReader("f", "100000.sst")
	assert.NoError(t, err)
	assert.Equal(t, mockReader, r)
	// case 3: get exist reader
	r, err = cache.GetReader("f", "100000.sst")
	assert.NoError(t, err)
	assert.Equal(t, mockReader, r)
	// case 4: evict not exist
	cache.Evict("f", "200000.sst")
	cache.Evict("f1", "100000.sst")
	// case 5: evict reader err
	mockReader.EXPECT().Close().Return(fmt.Errorf("err"))
	cache.Evict("f", "100000.sst")
	// case 6: close err
	mockReader.EXPECT().Close().Return(fmt.Errorf("err")).MaxTimes(2)
	_, _ = cache.GetReader("f", "100000.sst")
	_, _ = cache.GetReader("f", "200000.sst")
	err = cache.Close()
	assert.NoError(t, err)
}
