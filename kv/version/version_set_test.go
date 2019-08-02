package version

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/pkg/fileutil"
)

var vsTestPath = "test_data"

func TestVersionSetRecover(t *testing.T) {
	initVersionSetTestData()
	defer destoryVersionTestData()

	var vs = NewStoreVersionSet(vsTestPath, 2)
	err := vs.Recover()
	assert.Nil(t, err)
	vs.Destroy()

	vs = NewStoreVersionSet(vsTestPath, 2)
	err2 := vs.Recover()
	assert.Nil(t, err2)
	vs.Destroy()
}

func TestAssign_NextFileNumber(t *testing.T) {
	initVersionSetTestData()
	defer destoryVersionTestData()

	var vs = NewStoreVersionSet(vsTestPath, 2)
	assert.Equal(t, int64(2), vs.nextFileNumber, "wrong next file number")
	assert.Equal(t, int64(2), vs.NextFileNumber(), "assign wrong next file number")
	assert.Equal(t, int64(3), vs.nextFileNumber, "wrong next file number")
}

func TestVersionID(t *testing.T) {
	initVersionSetTestData()
	defer destoryVersionTestData()

	var vs = NewStoreVersionSet(vsTestPath, 2)
	assert.Equal(t, int64(0), vs.versionID, "wrong new version id")
	assert.Equal(t, int64(0), vs.newVersionID(), "assign wrong version id")
	assert.Equal(t, int64(1), vs.versionID, "wrong next version id")
}

func TestCommitFamilyEditLog(t *testing.T) {
	initVersionSetTestData()
	defer destoryVersionTestData()

	var vs = NewStoreVersionSet(vsTestPath, 2)
	assert.NotNil(t, vs, "cannot create store version")
	var err = vs.Recover()
	assert.Nil(t, err, "recover error")

	err = vs.CommitFamilyEditLog("f", nil)
	assert.NotNil(t, err, "commit not exist family version")

	familyID := 1
	vs.CreateFamilyVersion("f", familyID)
	editLog := NewEditLog(familyID)
	newFile := CreateNewFile(1, NewFileMeta(12, 1, 100, 2014))
	editLog.Add(newFile)
	editLog.Add(NewDeleteFile(1, 123))
	err = vs.CommitFamilyEditLog("f", editLog)
	assert.Nil(t, err, "commit family edit log error")

	vs.Destroy()

	// test recover many times
	for i := 0; i < 3; i++ {
		vs = NewStoreVersionSet(vsTestPath, 2)
		vs.CreateFamilyVersion("f", familyID)
		err = vs.Recover()
		assert.Nil(t, err, "recover error")

		familyVersion := vs.GetFamilyVersion("f")
		assert.Equal(t, newFile.file, familyVersion.GetCurrent().getAllFiles()[0], "cannot recover family version data")
		assert.Equal(t, int64(3+i), vs.nextFileNumber, "recover file number error")

		vs.Destroy()
	}
}

func TestCreateFamily(t *testing.T) {
	initVersionSetTestData()
	defer destoryVersionTestData()

	var vs = NewStoreVersionSet(vsTestPath, 2)

	familyVersion := vs.CreateFamilyVersion("family", 1)
	assert.NotNil(t, familyVersion, "get nil family version")

	familyVersion2 := vs.GetFamilyVersion("family")
	assert.NotNil(t, familyVersion2, "get nil family version2")

	assert.Equal(t, familyVersion, familyVersion2, "get diff family version")
}

func initVersionSetTestData() {
	if err := fileutil.MkDirIfNotExist(vsTestPath); err != nil {
		fmt.Println("create test path error")
	}
}

func destoryVersionTestData() {
	if err := fileutil.RemoveDir(vsTestPath); err != nil {
		fmt.Println("delete test path error")
	}
}
