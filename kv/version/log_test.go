package version

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/kv/table"
	"github.com/lindb/lindb/pkg/timeutil"
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
		return &newFile{}
	})
}

func TestNewFile(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	nf := CreateNewFile(1, NewFileMeta(12, 1, 100, 2014))
	bytes, err := nf.Encode()
	assert.NoError(t, err)

	newFile2 := &newFile{}
	err = newFile2.Decode(bytes)
	assert.NoError(t, err)
	assert.Equal(t, nf, newFile2)
	version := NewMockVersion(ctrl)
	version.EXPECT().AddFile(1, NewFileMeta(12, 1, 100, 2014))
	newFile2.apply(version)
}

func TestDeleteFile(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	dFile := NewDeleteFile(1, 120)
	bytes, err := dFile.Encode()
	assert.NoError(t, err)

	deleteFile2 := &deleteFile{}
	err = deleteFile2.Decode(bytes)
	assert.NoError(t, err)
	assert.Equal(t, dFile, deleteFile2)
	version := NewMockVersion(ctrl)
	version.EXPECT().DeleteFile(1, table.FileNumber(120))
	deleteFile2.apply(version)
}

func TestNextFileNumber(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	nextFile := NewNextFileNumber(12)
	bytes, err := nextFile.Encode()
	assert.NoError(t, err)

	nextFileNumber2 := &nextFileNumber{}
	err = nextFileNumber2.Decode(bytes)
	assert.NoError(t, err)
	assert.Equal(t, nextFile, nextFileNumber2)
	nextFileNumber2.apply(nil)
	sVersion := NewMockStoreVersionSet(ctrl)
	sVersion.EXPECT().setNextFileNumberWithoutLock(table.FileNumber(12))
	nextFileNumber2.applyVersionSet(sVersion)
}

func TestNewRollupFile(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	rollupFile := CreateNewRollupFile(12, 10)
	bytes, err := rollupFile.Encode()
	assert.NoError(t, err)

	rollupFile2 := &newRollupFile{}

	err = rollupFile2.Decode(bytes)
	assert.NoError(t, err)
	assert.Equal(t, rollupFile, rollupFile2)
	version := NewMockVersion(ctrl)
	version.EXPECT().AddRollupFile(table.FileNumber(12), timeutil.Interval(10))
	rollupFile2.apply(version)
}

func TestDeleteRollupFile(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	rollupFile := CreateDeleteRollupFile(12)
	bytes, err := rollupFile.Encode()
	assert.NoError(t, err)

	rollupFile2 := &deleteRollupFile{}

	err = rollupFile2.Decode(bytes)
	assert.NoError(t, err)
	assert.Equal(t, rollupFile, rollupFile2)
	version := NewMockVersion(ctrl)
	version.EXPECT().DeleteRollupFile(table.FileNumber(12))
	rollupFile2.apply(version)
}

func TestNewReferenceFile(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	referenceFile := CreateNewReferenceFile(FamilyID(10), 12)
	bytes, err := referenceFile.Encode()
	assert.NoError(t, err)

	referenceFile2 := &newReferenceFile{}

	err = referenceFile2.Decode(bytes)
	assert.NoError(t, err)
	assert.Equal(t, referenceFile, referenceFile2)
	version := NewMockVersion(ctrl)
	version.EXPECT().AddReferenceFile(FamilyID(10), table.FileNumber(12))
	referenceFile2.apply(version)
}

func TestDeleteReferenceFile(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	referenceFile := CreateDeleteReferenceFile(FamilyID(10), 12)
	bytes, err := referenceFile.Encode()
	assert.NoError(t, err)

	referenceFile2 := &deleteReferenceFile{}

	err = referenceFile2.Decode(bytes)
	assert.NoError(t, err)
	assert.Equal(t, referenceFile, referenceFile2)
	version := NewMockVersion(ctrl)
	version.EXPECT().DeleteReferenceFile(FamilyID(10), table.FileNumber(12))
	referenceFile2.apply(version)
}
