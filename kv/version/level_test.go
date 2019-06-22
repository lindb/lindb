package version

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_File_Level(t *testing.T) {
	level := newLevel()

	level.addFile(NewFileMeta(1, 1, 10, 1024))
	level.addFile(NewFileMeta(1, 1, 10, 1024))
	level.addFile(NewFileMeta(1, 1, 10, 1024))

	var files = level.getFiles()

	assert.Equal(t, 1, len(files), "add file wrong")

	//add file
	level.addFile(NewFileMeta(2, 1, 10, 1024))
	level.addFile(NewFileMeta(20, 1, 10, 1024))

	//delete file
	level.deleteFile(2)

	files = level.getFiles()
	assert.Equal(t, 2, len(files), "dlete file wrong")
}

func Test_Add_Files(t *testing.T) {
	level := newLevel()

	level.addFiles(NewFileMeta(1, 1, 10, 1024), NewFileMeta(2, 1, 10, 1024), NewFileMeta(3, 1, 10, 1024))

	var files = level.getFiles()

	assert.Equal(t, 3, len(files), "add files wrong")
}
