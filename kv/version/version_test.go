package version

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestVersion_New(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	defer func() {
		if err := recover(); err != nil {
			assert.NotNil(t, err)
		} else {
			assert.True(t, false)
		}
	}()

	fv := NewMockFamilyVersion(ctrl)
	vs := NewMockStoreVersionSet(ctrl)
	fv.EXPECT().GetVersionSet().Return(vs).MaxTimes(2)
	vs.EXPECT().numberOfLevels().Return(2)
	v := newVersion(1, fv)
	assert.NotNil(t, v)
	assert.Equal(t, fv, v.GetFamilyVersion())
	// test new panic
	vs.EXPECT().numberOfLevels().Return(-1)
	_ = newVersion(1, fv)
}

func TestVersion_Release(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	fv := NewMockFamilyVersion(ctrl)
	vs := NewMockStoreVersionSet(ctrl)
	fv.EXPECT().GetVersionSet().Return(vs).MaxTimes(2)
	vs.EXPECT().numberOfLevels().Return(2)
	v := newVersion(1, fv)
	assert.Equal(t, int32(0), v.ref.Load())
	v.retain()
	assert.Equal(t, int32(1), v.ref.Load())
	fv.EXPECT().removeVersion(v)
	v.release()
	assert.Equal(t, int32(0), v.ref.Load())
}

func TestVersion_Files(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	fv := NewMockFamilyVersion(ctrl)
	vs := NewMockStoreVersionSet(ctrl)
	fv.EXPECT().GetVersionSet().Return(vs).AnyTimes()
	vs.EXPECT().numberOfLevels().Return(2).AnyTimes()
	v := newVersion(1, fv)
	v.addFile(0, &FileMeta{fileNumber: 1})
	v.addFile(-10, &FileMeta{fileNumber: 2})
	v.addFile(2, &FileMeta{fileNumber: 3})
	v.addFiles(1, []*FileMeta{{fileNumber: 4}})
	assert.Equal(t, 2, len(v.getAllFiles()))
	assert.Equal(t, 0, v.NumberOfFilesInLevel(-1))
	assert.Equal(t, 0, v.NumberOfFilesInLevel(10))
	assert.Equal(t, 1, v.NumberOfFilesInLevel(0))
	assert.Equal(t, 1, v.NumberOfFilesInLevel(1))

	vs.EXPECT().newVersionID().Return(int64(2))
	v2 := v.cloneVersion()
	assert.Equal(t, 1, v2.NumberOfFilesInLevel(0))
	assert.Equal(t, 1, v2.NumberOfFilesInLevel(1))

	assert.Nil(t, v.getFiles(-1))
	assert.Nil(t, v.getFiles(3))
	assert.Equal(t, 1, len(v.getFiles(0)))
	assert.Equal(t, 1, len(v.getFiles(1)))
	v.deleteFile(-1, int64(4))
	assert.Equal(t, 2, len(v.getAllFiles()))
	v.deleteFile(1, int64(4))
	assert.Equal(t, 1, len(v.getAllFiles()))
	assert.Nil(t, v.getFiles(1))
}

func TestVersion_Find_Files(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	fv := NewMockFamilyVersion(ctrl)
	vs := NewMockStoreVersionSet(ctrl)
	fv.EXPECT().GetVersionSet().Return(vs).AnyTimes()
	vs.EXPECT().numberOfLevels().Return(2).AnyTimes()
	v := newVersion(1, fv)
	f1 := FileMeta{fileNumber: 1, minKey: 10, maxKey: 200}
	f2 := FileMeta{fileNumber: 2, minKey: 50, maxKey: 400}
	v.addFile(0, &f1)
	v.addFile(1, &f2)
	files := v.findFiles(100)
	assert.Equal(t, 2, len(files))
	assert.Equal(t, f1, *files[0])
	assert.Equal(t, f2, *files[1])

	files = v.findFiles(20)
	assert.Equal(t, 1, len(files))
	assert.Equal(t, f1, *files[0])

	files = v.findFiles(300)
	assert.Equal(t, 1, len(files))
	assert.Equal(t, f2, *files[0])

	files = v.findFiles(3000)
	assert.Equal(t, 0, len(files))
	files = v.findFiles(5)
	assert.Equal(t, 0, len(files))
}

func TestVersion_PickL0Compaction(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	fv := NewMockFamilyVersion(ctrl)
	vs := NewMockStoreVersionSet(ctrl)
	fv.EXPECT().GetVersionSet().Return(vs).AnyTimes()
	fv.EXPECT().GetID().Return(1).AnyTimes()
	vs.EXPECT().numberOfLevels().Return(2).AnyTimes()
	v := newVersion(1, fv)
	/*
	* Level 0:
	* file 1: 1~10
	* file 2: 1000~1001
	 */
	f1 := FileMeta{fileNumber: 1, minKey: 10, maxKey: 100}
	f2 := FileMeta{fileNumber: 2, minKey: 1000, maxKey: 1001}
	v.addFiles(0, []*FileMeta{&f1, &f2})
	/*
	* Level 1:
	* file 3: 1~5
	* file 4: 100~200
	* file 5: 400~500
	 */
	f3 := FileMeta{fileNumber: 3, minKey: 1, maxKey: 5}
	f4 := FileMeta{fileNumber: 4, minKey: 100, maxKey: 200}
	f5 := FileMeta{fileNumber: 5, minKey: 400, maxKey: 500}
	v.addFiles(1, []*FileMeta{&f3, &f4, &f5})

	compaction := v.PickL0Compaction(5)
	assert.Nil(t, compaction)

	compaction = v.PickL0Compaction(1)
	assert.NotNil(t, compaction)
	assert.Equal(t, 2, len(compaction.levelInputs))
	assert.Equal(t, 1, len(compaction.levelUpInputs))
	assert.Equal(t, f4, *compaction.levelUpInputs[0])

	f6 := FileMeta{fileNumber: 6, minKey: 1, maxKey: 1000}
	v.addFiles(0, []*FileMeta{&f6})
	compaction = v.PickL0Compaction(1)
	assert.Equal(t, 3, len(compaction.levelUpInputs))
}
