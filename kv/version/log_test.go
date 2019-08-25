package version

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRegisterLogType(t *testing.T) {
	defer func() {
		if err := recover(); err != nil {
			assert.True(t, true)
		} else {
			assert.Fail(t, "test panic fail")
		}
	}()

	RegisterLogType(1, func() Log {
		return &NewFile{}
	})
}
func TestNewFile(t *testing.T) {
	newFile := CreateNewFile(1, NewFileMeta(12, 1, 100, 2014))
	bytes, err := newFile.Encode()
	assert.Nil(t, err, "new file encode error")

	newFile2 := &NewFile{}
	err2 := newFile2.Decode(bytes)
	assert.Nil(t, err2, "new file decode error")

	assert.Equal(t, newFile, newFile2, "file1 not eqals files")
}

func TestDeleteFile(t *testing.T) {
	deleteFile := NewDeleteFile(1, 120)
	bytes, err := deleteFile.Encode()
	assert.Nil(t, err, "delete file encode error")

	deleteFile2 := &DeleteFile{}
	err2 := deleteFile2.Decode(bytes)

	assert.Nil(t, err2, "delete file decode error")
	assert.Equal(t, deleteFile, deleteFile2, "delete file1 not equals delete file2")
}

func TestNextFileNumber(t *testing.T) {
	nextFileNumber := NewNextFileNumber(12)
	bytes, err := nextFileNumber.Encode()
	assert.Nil(t, err, "next file nubmer encode error")

	nextFileNumber2 := &NextFileNumber{}
	err2 := nextFileNumber2.Decode(bytes)

	assert.Nil(t, err2, "next file nubmer decode error")
	assert.Equal(t, nextFileNumber, nextFileNumber2, "next file number 1 != next file number 2")
}
