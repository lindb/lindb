// Licensed to LinDB under one or more contributor
// license agreements. See the NOTICE file distributed with
// this work for additional information regarding copyright
// ownership. LinDB licenses this file to you under
// the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

package version

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/kv/table"
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
	assert.Len(t, v.Levels(), 2)
	assert.Equal(t, int64(1), v.ID())
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
	assert.Equal(t, int32(0), v.NumOfRef())
	v.Retain()
	assert.Equal(t, int32(1), v.NumOfRef())
	fv.EXPECT().removeVersion(v)
	v.Release()
	assert.Equal(t, int32(0), v.NumOfRef())
}

func TestVersion_Files(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	fv := NewMockFamilyVersion(ctrl)
	vs := NewMockStoreVersionSet(ctrl)
	fv.EXPECT().GetVersionSet().Return(vs).AnyTimes()
	vs.EXPECT().numberOfLevels().Return(2).AnyTimes()
	v := newVersion(1, fv)
	v.AddFile(0, &FileMeta{fileNumber: 1})
	v.AddFile(-10, &FileMeta{fileNumber: 2})
	v.AddFile(2, &FileMeta{fileNumber: 3})
	v.AddFiles(1, []*FileMeta{{fileNumber: 4}})
	assert.Equal(t, 2, len(v.GetAllFiles()))
	assert.Equal(t, 0, v.NumberOfFilesInLevel(-1))
	assert.Equal(t, 0, v.NumberOfFilesInLevel(10))
	assert.Equal(t, 1, v.NumberOfFilesInLevel(0))
	assert.Equal(t, 1, v.NumberOfFilesInLevel(1))

	vs.EXPECT().newVersionID().Return(int64(2))
	v2 := v.Clone()
	assert.Equal(t, 1, v2.NumberOfFilesInLevel(0))
	assert.Equal(t, 1, v2.NumberOfFilesInLevel(1))

	assert.Nil(t, v.GetFiles(-1))
	assert.Nil(t, v.GetFiles(3))
	assert.Equal(t, 1, len(v.GetFiles(0)))
	assert.Equal(t, 1, len(v.GetFiles(1)))
	v.DeleteFile(-1, table.FileNumber(4))
	assert.Equal(t, 2, len(v.GetAllFiles()))
	v.DeleteFile(1, table.FileNumber(4))
	assert.Equal(t, 1, len(v.GetAllFiles()))
	assert.Nil(t, v.GetFiles(1))
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
	v.AddFile(0, &f1)
	v.AddFile(1, &f2)
	files := v.FindFiles(100)
	assert.Equal(t, 2, len(files))
	assert.Equal(t, f1, *files[0])
	assert.Equal(t, f2, *files[1])

	files = v.FindFiles(20)
	assert.Equal(t, 1, len(files))
	assert.Equal(t, f1, *files[0])

	files = v.FindFiles(300)
	assert.Equal(t, 1, len(files))
	assert.Equal(t, f2, *files[0])

	files = v.FindFiles(3000)
	assert.Equal(t, 0, len(files))
	files = v.FindFiles(5)
	assert.Equal(t, 0, len(files))
}

func TestVersion_PickL0Compaction(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	fv := NewMockFamilyVersion(ctrl)
	vs := NewMockStoreVersionSet(ctrl)
	fv.EXPECT().GetVersionSet().Return(vs).AnyTimes()
	fv.EXPECT().GetID().Return(FamilyID(1)).AnyTimes()
	vs.EXPECT().numberOfLevels().Return(2).AnyTimes()
	v := newVersion(1, fv)
	/*
	* Level 0:
	* file 1: 1~10
	* file 2: 1000~1001
	 */
	f1 := FileMeta{fileNumber: 1, minKey: 10, maxKey: 100}
	f2 := FileMeta{fileNumber: 2, minKey: 1000, maxKey: 1001}
	v.AddFiles(0, []*FileMeta{&f1, &f2})
	/*
	* Level 1:
	* file 3: 1~5
	* file 4: 100~200
	* file 5: 400~500
	 */
	f3 := FileMeta{fileNumber: 3, minKey: 1, maxKey: 5}
	f4 := FileMeta{fileNumber: 4, minKey: 100, maxKey: 200}
	f5 := FileMeta{fileNumber: 5, minKey: 400, maxKey: 500}
	v.AddFiles(1, []*FileMeta{&f3, &f4, &f5})

	compaction := v.PickL0Compaction(5)
	assert.Nil(t, compaction)

	compaction = v.PickL0Compaction(1)
	assert.NotNil(t, compaction)
	assert.Equal(t, 2, len(compaction.levelInputs))
	assert.Equal(t, 1, len(compaction.levelUpInputs))
	assert.Equal(t, f4, *compaction.levelUpInputs[0])

	f6 := FileMeta{fileNumber: 6, minKey: 1, maxKey: 1000}
	v.AddFiles(0, []*FileMeta{&f6})
	compaction = v.PickL0Compaction(1)
	assert.Equal(t, 3, len(compaction.levelUpInputs))
}

func TestVersion_RollupJob(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	fv := NewMockFamilyVersion(ctrl)
	vs := NewMockStoreVersionSet(ctrl)
	fv.EXPECT().GetVersionSet().Return(vs).AnyTimes()
	fv.EXPECT().GetID().Return(FamilyID(1)).AnyTimes()
	vs.EXPECT().numberOfLevels().Return(2).AnyTimes()
	v := newVersion(1, fv)
	v.AddRollupFile(10, 3)
	v.DeleteRollupFile(10)
	assert.Empty(t, v.GetRollupFiles())
	v.AddReferenceFile(10, 100)
	v.AddReferenceFile(10, 10)
	v.DeleteReferenceFile(10, 10)
	assert.Equal(t, map[FamilyID][]table.FileNumber{10: {100}}, v.GetReferenceFiles())
}
