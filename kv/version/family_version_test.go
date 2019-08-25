package version

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/kv/table"
)

func TestFamilyVersion(t *testing.T) {
	initVersionSetTestData()
	defer destroyVersionTestData()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cache := table.NewMockCache(ctrl)

	var vs = NewStoreVersionSet(vsTestPath, cache, 2)
	familyVersion1 := vs.CreateFamilyVersion("f", 1)
	assert.Equal(t, 1, familyVersion1.GetID())

	assert.Equal(t, 0, len(familyVersion1.GetAllActiveFiles()), "file list not empty")
	snapshot := familyVersion1.GetSnapshot()

	version1 := snapshot.GetCurrent()
	// add mock version ref
	file1 := NewFileMeta(12, 1, 50, 2014)
	version1.addFile(1, file1)
	file2 := NewFileMeta(13, 1, 10, 2014)
	file3 := NewFileMeta(14, 1, 100, 2014)
	var files []*FileMeta
	files = append(files, file2, file3)
	version1.addFiles(1, files)
	assert.Equal(t, 3, len(familyVersion1.GetAllActiveFiles()), "file list != 3")

	version2 := version1.cloneVersion()
	familyVersion1.appendVersion(version2)
	fv := familyVersion1.(*familyVersion)
	snapshot2 := familyVersion1.GetSnapshot()
	assert.Equal(t, 2, len(fv.activeVersions), "version list !=2")
	assert.Equal(t, version2, snapshot2.GetCurrent(), "get wrong current version")
	assert.Equal(t, 3, len(familyVersion1.GetAllActiveFiles()), "file list != 3")

	// delete file1
	version2.deleteFile(1, 12)
	// can get file from version1
	assert.Equal(t, 3, len(familyVersion1.GetAllActiveFiles()), "file list != 3")

	// release version1
	snapshot.Close()
	// cannot get file from version1
	assert.Equal(t, 2, len(familyVersion1.GetAllActiveFiles()), "file list != 2")

	file4 := NewFileMeta(40, 70, 100, 1024)
	// add invalid version
	version1.addFile(1, file4)
	assert.Equal(t, 2, len(familyVersion1.GetAllActiveFiles()), "file list != 2")

	reader := table.NewMockReader(ctrl)
	cache.EXPECT().GetReader(gomock.Any(), gomock.Any()).Return(reader, nil).MaxTimes(3)
	// add duplicate file
	version2.addFile(1, file3)
	assert.Equal(t, 2, len(familyVersion1.GetAllActiveFiles()), "file list != 2")
	reader1, _ := snapshot2.FindReaders(70)
	assert.Equal(t, 1, len(reader1))

	reader2, _ := snapshot2.FindReaders(5)
	assert.Equal(t, 2, len(reader2))

	assert.Equal(t, vs, familyVersion1.GetVersionSet())

	// test append new version then remove old version
	version1 = snapshot2.GetCurrent()
	version2 = version1.cloneVersion()
	snapshot2.Close()

	familyVersion1.appendVersion(version2)
	fv = familyVersion1.(*familyVersion)
	assert.Equal(t, 1, len(fv.activeVersions))
}
