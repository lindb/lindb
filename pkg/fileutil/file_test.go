package fileutil

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

type User struct {
	Name string
}

var testPath = "./file"

func TestFileUtil(t *testing.T) {
	_ = MkDirIfNotExist(testPath)

	defer func() {
		_ = RemoveDir(testPath)
	}()

	assert.True(t, Exist(testPath))

	files, _ := ListDir(testPath)
	assert.Len(t, files, 0)

	assert.Nil(t, MkDir(filepath.Join(os.TempDir(), "tmp/test.toml")))
}

func TestFileUtil_errors(t *testing.T) {
	// inexistent directory
	_, err := ListDir(filepath.Join(os.TempDir(), "/tmp/tmp/tmp/tmp"))

	// encode toml failure
	assert.NotNil(t, err)
}

func TestGetExistPath(t *testing.T) {
	assert.Equal(t, "/tmp", GetExistPath("/tmp/test1/test333"))
}
