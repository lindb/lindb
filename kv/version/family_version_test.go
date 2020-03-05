package version

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/kv/table"
)

func TestFamilyID_Int(t *testing.T) {
	assert.Equal(t, 10, FamilyID(10).Int())
	assert.Equal(t, int32(10), FamilyID(10).Int32())
}

func TestFamilyVersion(t *testing.T) {
	initVersionSetTestData()
	ctrl := gomock.NewController(t)
	defer func() {
		destroyVersionTestData()
		ctrl.Finish()
	}()

	cache := table.NewMockCache(ctrl)

	var vs = NewStoreVersionSet(vsTestPath, cache, 2)
	familyVersion1 := vs.CreateFamilyVersion("f", 1)
	assert.Equal(t, FamilyID(1), familyVersion1.GetID())

	assert.Equal(t, 0, len(familyVersion1.GetAllActiveFiles()), "file list not empty")
	snapshot := familyVersion1.GetSnapshot()

	version1 := snapshot.GetCurrent()
	// add mock version ref
	file1 := NewFileMeta(12, 1, 50, 2014)
	version1.AddFile(1, file1)
	file2 := NewFileMeta(13, 1, 10, 2014)
	file3 := NewFileMeta(14, 1, 100, 2014)
	var files []*FileMeta
	files = append(files, file2, file3)
	version1.AddFiles(1, files)
	assert.Equal(t, 3, len(familyVersion1.GetAllActiveFiles()), "file list != 3")

	version2 := version1.Clone()
	familyVersion1.appendVersion(version2)
	fv := familyVersion1.(*familyVersion)
	snapshot2 := familyVersion1.GetSnapshot()
	assert.Equal(t, 2, len(fv.activeVersions), "version list !=2")
	assert.Equal(t, version2, snapshot2.GetCurrent(), "get wrong current version")
	assert.Equal(t, 3, len(familyVersion1.GetAllActiveFiles()), "file list != 3")

	// delete file1
	version2.DeleteFile(1, 12)
	// can get file from version1
	assert.Equal(t, 3, len(familyVersion1.GetAllActiveFiles()), "file list != 3")

	// release version1
	snapshot.Close()
	// cannot get file from version1
	assert.Equal(t, 2, len(familyVersion1.GetAllActiveFiles()), "file list != 2")

	file4 := NewFileMeta(40, 70, 100, 1024)
	// add invalid version
	version1.AddFile(1, file4)
	assert.Equal(t, 2, len(familyVersion1.GetAllActiveFiles()), "file list != 2")

	reader := table.NewMockReader(ctrl)
	cache.EXPECT().GetReader(gomock.Any(), gomock.Any()).Return(reader, nil).MaxTimes(3)
	// add duplicate file
	version2.AddFile(1, file3)
	assert.Equal(t, 2, len(familyVersion1.GetAllActiveFiles()), "file list != 2")
	reader1, _ := snapshot2.FindReaders(70)
	assert.Equal(t, 1, len(reader1))

	reader2, _ := snapshot2.FindReaders(5)
	assert.Equal(t, 2, len(reader2))

	assert.Equal(t, vs, familyVersion1.GetVersionSet())

	// test append new version then remove old version
	version1 = snapshot2.GetCurrent()
	version2 = version1.Clone()
	snapshot2.Close()

	familyVersion1.appendVersion(version2)
	fv = familyVersion1.(*familyVersion)
	assert.Equal(t, 1, len(fv.activeVersions))
}

func TestFamilyVersion_Rollup(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer func() {
		ctrl.Finish()
	}()

	vs := NewMockStoreVersionSet(ctrl)
	vs.EXPECT().newVersionID().Return(int64(100))
	vs.EXPECT().numberOfLevels().Return(2)
	fv := newFamilyVersion(1, "f", vs)
	version := NewMockVersion(ctrl)
	version.EXPECT().ID().Return(int64(10))
	version.EXPECT().GetFamilyVersion().Return(fv)
	fv.appendVersion(version)
	version.EXPECT().GetReferenceFiles().Return(nil)
	assert.Nil(t, fv.GetLiveReferenceFiles())
	version.EXPECT().GetRollupFiles()
	assert.Nil(t, fv.GetLiveRollupFiles())
}
