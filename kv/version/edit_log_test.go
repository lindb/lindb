package version

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEditLogCodec(t *testing.T) {
	editLog := NewEditLog(1)
	assert.True(t, editLog.IsEmpty(), "edit log not is empty")

	newFile := CreateNewFile(1, NewFileMeta(12, 1, 100, 2014))
	editLog.Add(newFile)
	editLog.Add(NewDeleteFile(1, 123))

	v, err := editLog.marshal()

	assert.Nil(t, err, "marshal error")
	assert.True(t, len(v) > 0, "encode edit log error")

	editLog2 := NewEditLog(1)
	err2 := editLog2.unmarshal(v)
	assert.Nil(t, err2, "unmarshal error")

	assert.Equal(t, editLog, editLog2, "edit log not eqauls")
}

func TestApply(t *testing.T) {
	initVersionSetTestData()
	defer destoryVersionTestData()

	var vs = NewStoreVersionSet(vsTestPath, 2)
	familyVersion := vs.CreateFamilyVersion("family", 1)
	editLog := NewEditLog(1)
	newFile := &NewFile{level: 1, file: NewFileMeta(12, 1, 100, 2014)}
	editLog.Add(newFile)
	version := newVersion(1, familyVersion)
	editLog.apply(version)

	assert.Equal(t, 1, len(version.getAllFiles()), "cannot add file into version")
	//delete file
	editLog2 := NewEditLog(1)
	editLog2.Add(NewDeleteFile(1, 12))
	editLog2.apply(version)
	assert.Equal(t, 0, len(version.getAllFiles()), "cannot delete file from version")
}
