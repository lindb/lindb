package fileutil

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

var testPath = "./file"

func TestMkDirIfNotExist(t *testing.T) {
	defer func() {
		mkdirAllFunc = os.MkdirAll
		_ = RemoveDir(testPath)
	}()

	mkdirAllFunc = func(path string, perm os.FileMode) error {
		return fmt.Errorf("err")
	}
	err := MkDirIfNotExist(testPath)
	assert.Error(t, err)

	err = MkDir(testPath)
	assert.Error(t, err)
	mkdirAllFunc = os.MkdirAll
	err = MkDir(testPath)
	assert.NoError(t, err)
}

func TestRemoveDir(t *testing.T) {
	_ = MkDirIfNotExist(testPath)

	defer func() {
		removeAllFunc = os.RemoveAll
		_ = RemoveDir(testPath)
	}()
	removeAllFunc = func(path string) error {
		return fmt.Errorf("err")
	}
	err := RemoveDir(testPath)
	assert.Error(t, err)
}

func TestFileUtil(t *testing.T) {
	_ = MkDirIfNotExist(testPath)

	defer func() {
		_ = RemoveDir(testPath)
	}()

	assert.True(t, Exist(testPath))
}

func TestFileUtil_errors(t *testing.T) {
	// not existent directory
	_, err := ListDir(filepath.Join(os.TempDir(), "/tmp/tmp/tmp/tmp"))

	// encode toml failure
	assert.NotNil(t, err)
}

func TestGetExistPath(t *testing.T) {
	assert.Equal(t, "/tmp", GetExistPath("/tmp/test1/test333"))
}

func TestListDir(t *testing.T) {
	_ = MkDirIfNotExist(testPath)

	defer func() {
		_ = RemoveDir(testPath)
	}()
	_, _ = os.Create(testPath + "/file1")
	files, err := ListDir(testPath)
	assert.NoError(t, err)
	assert.Len(t, files, 1)
}

func TestRemoveFile(t *testing.T) {
	_ = MkDirIfNotExist(testPath)

	defer func() {
		_ = RemoveDir(testPath)
		removeFunc = os.Remove
	}()
	_, _ = os.Create(testPath + "/file1")
	err := RemoveFile(testPath + "/file1")
	assert.NoError(t, err)
	files, err := ListDir(testPath)
	assert.NoError(t, err)
	assert.Len(t, files, 0)

	_, _ = os.Create(testPath + "/file1")
	removeFunc = func(name string) error {
		return fmt.Errorf("err")
	}
	err = RemoveFile(testPath + "/file1")
	assert.Error(t, err)
	err = RemoveFile(testPath + "/file2")
	assert.NoError(t, err)
	files, err = ListDir(testPath)
	assert.NoError(t, err)
	assert.Len(t, files, 1)
}
