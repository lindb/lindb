package page

import (
	"fmt"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/pkg/fileutil"
)

var testPath = "test"
var fileName = "fileName"

func TestMappedPage_err(t *testing.T) {
	_ = fileutil.MkDirIfNotExist(testPath)

	defer func() {
		_ = fileutil.RemoveDir(testPath)
		mapFileFunc = fileutil.RWMap
	}()

	mapFileFunc = func(filePath string, size int) ([]byte, error) {
		return nil, fmt.Errorf("err")
	}
	mp, err := NewMappedPage(filepath.Join(testPath, fileName), 128)
	assert.Error(t, err)
	assert.Nil(t, mp)
}

func TestMappedPage(t *testing.T) {
	_ = fileutil.MkDirIfNotExist(testPath)

	defer func() {
		_ = fileutil.RemoveDir(testPath)
	}()

	bytes := []byte("12345")

	mp, err := NewMappedPage(filepath.Join(testPath, fileName), 128)
	assert.NoError(t, err)

	// copy data
	mp.WriteBytes(bytes, 0)

	assert.NoError(t, mp.Sync())
	assert.Equal(t, filepath.Join(testPath, fileName), mp.FilePath())
	assert.NotNil(t, 128, mp.Size())
	assert.Equal(t, bytes, mp.ReadBytes(0, 5))
	assert.False(t, mp.Closed())
	assert.NoError(t, mp.Close())
	assert.True(t, mp.Closed())
	assert.NoError(t, mp.Close())
}

func TestMappedPage_Write_number(t *testing.T) {
	_ = fileutil.MkDirIfNotExist(testPath)

	defer func() {
		_ = fileutil.RemoveDir(testPath)
	}()

	mp, err := NewMappedPage(filepath.Join(testPath, fileName), 128)
	assert.NoError(t, err)
	mp.PutUint32(10, 0)
	mp.PutUint64(999, 8)
	mp.PutUint8(50, 16)
	assert.Equal(t, uint32(999), mp.ReadUint32(8))
	assert.Equal(t, uint64(10), mp.ReadUint64(0))
	assert.Equal(t, uint8(50), mp.ReadUint8(16))

	err = mp.Close()
	assert.NoError(t, err)
}
