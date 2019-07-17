package util

import (
	"path"
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

	files, _ := ListDir(testPath)
	assert.Equal(t, "toml", files[0])
}
