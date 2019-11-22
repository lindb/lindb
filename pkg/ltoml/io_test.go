package ltoml

import (
	"os"
	"path"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/pkg/fileutil"
)

type User struct {
	Name string
}

var testPath = "./file"

func Test_Encode(t *testing.T) {
	_ = fileutil.MkDirIfNotExist(testPath)
	defer func() {
		_ = fileutil.RemoveDir(testPath)
	}()
	user := User{Name: "LinDB"}
	file := path.Join(testPath, "toml")
	err := EncodeToml(file, &user)
	if err != nil {
		t.Fatal(err)
	}
	user2 := User{}
	err = DecodeToml(file, &user2)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, user, user2)

	files, _ := fileutil.ListDir(testPath)
	assert.Equal(t, "toml", files[0])

	assert.NotNil(t, EncodeToml(filepath.Join(os.TempDir(), "/tmp/test.toml"), []byte{}))
}

func Test_WriteConfig(t *testing.T) {
	_ = fileutil.MkDirIfNotExist(testPath)
	defer func() {
		_ = fileutil.RemoveDir(testPath)
	}()
	assert.Nil(t, WriteConfig(path.Join(testPath, "toml"), ""))
}
