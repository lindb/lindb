package version

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFamilyVersion(t *testing.T) {
	initVersionSetTestData()
	defer destoryVersionTestData()
	var vs = NewStoreVersionSet(vsTestPath, 2)
	familyVersion := vs.CreateFamilyVersion("f", 1)

	assert.Equal(t, 0, len(familyVersion.GetAllFiles()), "file list not empty")

	version1 := familyVersion.GetCurrent()
	// add mock version ref
	version1.retain()
	file1 := NewFileMeta(12, 1, 50, 2014)
	version1.addFile(1, file1)
	file2 := NewFileMeta(13, 1, 10, 2014)
	file3 := NewFileMeta(14, 1, 100, 2014)
	var files []*FileMeta
	files = append(files, file2, file3)
	version1.addFiles(1, files)
	assert.Equal(t, 3, len(familyVersion.GetAllFiles()), "file list != 3")

	version2 := version1.cloneVersion()
	familyVersion.appendVersion(version2)
	assert.Equal(t, 2, len(familyVersion.activeVersions), "version list !=2")
	assert.Equal(t, version2, familyVersion.GetCurrent(), "get wrong current version")
	assert.Equal(t, 3, len(familyVersion.GetAllFiles()), "file list != 3")

	// delete file1
	version2.deleteFile(1, 12)
	// can get file from version1
	assert.Equal(t, 3, len(familyVersion.GetAllFiles()), "file list != 3")

	// release version1
	version1.Release()
	// cannot get file from version1
	assert.Equal(t, 2, len(familyVersion.GetAllFiles()), "file list != 2")

	file4 := NewFileMeta(40, 70, 100, 1024)
	// add invalid version
	version1.addFile(1, file4)
	assert.Equal(t, 2, len(familyVersion.GetAllFiles()), "file list != 2")

	// add dupliate file
	version2.addFile(1, file3)
	assert.Equal(t, 2, len(familyVersion.GetAllFiles()), "file list != 2")
	_, findFile1 := familyVersion.FindFiles(70)
	assert.Equal(t, 1, len(findFile1))

	_, findFile2 := familyVersion.FindFiles(5)
	assert.Equal(t, 2, len(findFile2))
}
