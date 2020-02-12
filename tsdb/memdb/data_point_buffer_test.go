package memdb

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/pkg/fileutil"
)

const testPath = "test_dp_buf"

func TestDataPointBuffer_AllocPage(t *testing.T) {
	defer func() {
		_ = fileutil.RemoveDir(testPath)
	}()
	_ = fileutil.MkDirIfNotExist(testPath)

	buf := newDataPointBuffer(testPath)
	for i := 0; i < 10000; i++ {
		b, err := buf.AllocPage()
		assert.NoError(t, err)
		assert.NotNil(t, b)
	}
	err := buf.Close()
	assert.NoError(t, err)
}

func TestDataPointBuffer_AllocPage_err(t *testing.T) {
	defer func() {
		mkdirFunc = fileutil.MkDirIfNotExist
		mapFunc = fileutil.RWMap
		_ = fileutil.RemoveDir(testPath)
	}()
	buf := newDataPointBuffer(testPath)
	mkdirFunc = func(path string) error {
		return fmt.Errorf("err")
	}
	// case 1: make file path err
	b, err := buf.AllocPage()
	assert.Error(t, err)
	assert.Nil(t, b)
	mkdirFunc = fileutil.MkDirIfNotExist

	// case 2: wrong region
	b, err = buf.AllocPage()
	assert.Error(t, err)
	assert.Nil(t, b)

	mapFunc = func(filePath string, size int) (bytes []byte, err error) {
		return nil, fmt.Errorf("err")
	}
	// case 3: map file err
	buf = newDataPointBuffer(testPath)
	b, err = buf.AllocPage()
	assert.Error(t, err)
	assert.Nil(t, b)
}

func TestDataPointBuffer_Close_err(t *testing.T) {
	defer func() {
		removeFunc = fileutil.RemoveDir
		unmapFunc = fileutil.Unmap
		_ = fileutil.RemoveDir(testPath)
	}()
	buf := newDataPointBuffer(testPath)
	b, err := buf.AllocPage()
	assert.NoError(t, err)
	assert.NotNil(t, b)
	// case 1: remove dir err
	removeFunc = func(path string) error {
		return fmt.Errorf("err")
	}
	err = buf.Close()
	assert.NoError(t, err)
	// case 2: unmap err
	removeFunc = fileutil.RemoveDir
	unmapFunc = func(data []byte) error {
		return fmt.Errorf("err")
	}
	err = buf.Close()
	assert.NoError(t, err)

}
