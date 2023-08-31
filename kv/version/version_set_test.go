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
	"fmt"
	"io/ioutil"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/lindb/common/pkg/fileutil"

	"github.com/lindb/lindb/kv/table"
	"github.com/lindb/lindb/pkg/bufioutil"
	"github.com/lindb/lindb/pkg/timeutil"
)

var vsTestPath = "test_data"

func TestVersionSetRecover(t *testing.T) {
	initVersionSetTestData()
	ctrl := gomock.NewController(t)
	defer func() {
		destroyVersionTestData()
		ctrl.Finish()
	}()
	cache := table.NewMockCache(ctrl)

	var vs = NewStoreVersionSet(vsTestPath, cache, 2)
	assert.Equal(t, table.FileNumber(1), vs.ManifestFileNumber())
	err := vs.Recover()
	assert.Nil(t, err)
	_ = vs.Destroy()

	vs = NewStoreVersionSet(vsTestPath, cache, 2)
	err2 := vs.Recover()
	assert.Nil(t, err2)
	assert.Equal(t, table.FileNumber(2), vs.ManifestFileNumber())
	_ = vs.Destroy()
}

func TestStoreVersionSet_Recover_err(t *testing.T) {
	initVersionSetTestData()
	ctrl := gomock.NewController(t)
	defer func() {
		writeFileFunc = ioutil.WriteFile
		renameFunc = os.Rename
		newBufferReaderFunc = bufioutil.NewBufioEntryReader
		newBufferWriterFunc = bufioutil.NewBufioEntryWriter
		readFileFunc = ioutil.ReadFile
		newEmptyEditLogFunc = newEmptyEditLog

		destroyVersionTestData()
		ctrl.Finish()
	}()
	cache := table.NewMockCache(ctrl)

	vs := NewStoreVersionSet(vsTestPath, cache, 2)
	// case 1: write manifest name err
	writeFileFunc = func(filename string, data []byte, perm os.FileMode) error {
		return fmt.Errorf("err")
	}
	err := vs.Recover()
	assert.Error(t, err)
	// case 2: remove current file name err
	writeFileFunc = ioutil.WriteFile
	renameFunc = func(oldpath, newpath string) error {
		return fmt.Errorf("err")
	}
	err = vs.Recover()
	assert.Error(t, err)
	// case 3: new buffer reader err
	renameFunc = os.Rename
	err = vs.Recover()
	assert.NoError(t, err)
	err = vs.Destroy()
	assert.NoError(t, err)
	vs = NewStoreVersionSet(vsTestPath, cache, 2) // reopen
	newBufferReaderFunc = func(fileName string) (reader bufioutil.BufioEntryReader, err error) {
		return nil, fmt.Errorf("ere")
	}
	err = vs.Recover()
	assert.Error(t, err)
	// case 4: read edit log err
	bufReader := bufioutil.NewMockBufioEntryReader(ctrl)
	newBufferReaderFunc = func(fileName string) (reader bufioutil.BufioEntryReader, err error) {
		return bufReader, nil
	}
	gomock.InOrder(
		bufReader.EXPECT().Next().Return(true),
		bufReader.EXPECT().Read().Return(nil, fmt.Errorf("err")),
		bufReader.EXPECT().Close().Return(fmt.Errorf("err")).AnyTimes(),
	)
	err = vs.Recover()
	assert.Error(t, err)
	// case 5: unmarshal edit log err
	editLog := NewMockEditLog(ctrl)
	newEmptyEditLogFunc = func() EditLog {
		return editLog
	}
	gomock.InOrder(
		bufReader.EXPECT().Next().Return(true),
		bufReader.EXPECT().Read().Return([]byte{1, 2}, nil),
		editLog.EXPECT().unmarshal(gomock.Any()).Return(fmt.Errorf("err")),
	)
	err = vs.Recover()
	assert.Error(t, err)
	// case 6: family id not found
	gomock.InOrder(
		bufReader.EXPECT().Next().Return(true),
		bufReader.EXPECT().Read().Return([]byte{1, 2}, nil),
		editLog.EXPECT().unmarshal(gomock.Any()).Return(nil),
		editLog.EXPECT().FamilyID().Return(FamilyID(9999)),
	)
	err = vs.Recover()
	assert.Error(t, err)
	// case 7: new buf writer err
	newBufferReaderFunc = bufioutil.NewBufioEntryReader
	newEmptyEditLogFunc = newEmptyEditLog
	newBufferWriterFunc = func(fileName string) (writer bufioutil.BufioWriter, err error) {
		return nil, fmt.Errorf("err")
	}
	err = vs.Recover()
	assert.Error(t, err)
	// case 8: write snapshot edit log err
	mockWriter := bufioutil.NewMockBufioWriter(ctrl)
	newBufferWriterFunc = func(fileName string) (writer bufioutil.BufioWriter, err error) {
		return mockWriter, nil
	}
	mockWriter.EXPECT().Write(gomock.Any()).Return(0, fmt.Errorf("err"))
	err = vs.Recover()
	assert.Error(t, err)
	// case 9: reader current file name err
	readFileFunc = func(filename string) (bytes []byte, err error) {
		return nil, fmt.Errorf("err")
	}
	err = vs.Recover()
	assert.Error(t, err)
}

func TestAssign_NextFileNumber(t *testing.T) {
	initVersionSetTestData()
	defer destroyVersionTestData()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	cache := table.NewMockCache(ctrl)

	var vs = NewStoreVersionSet(vsTestPath, cache, 2)
	vs1 := vs.(*storeVersionSet)
	assert.Equal(t, int64(2), vs1.nextFileNumber.Load(), "wrong next file number")
	assert.Equal(t, table.FileNumber(2), vs.NextFileNumber(), "assign wrong next file number")
	assert.Equal(t, int64(3), vs1.nextFileNumber.Load(), "wrong next file number")
}

func TestVersionID(t *testing.T) {
	initVersionSetTestData()
	defer destroyVersionTestData()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	cache := table.NewMockCache(ctrl)

	var vs = NewStoreVersionSet(vsTestPath, cache, 2)
	vs1 := vs.(*storeVersionSet)
	assert.Equal(t, int64(0), vs1.versionID.Load(), "wrong new version id")
	assert.Equal(t, int64(0), vs.newVersionID(), "assign wrong version id")
	assert.Equal(t, int64(1), vs1.versionID.Load(), "wrong next version id")
}

func TestCommitFamilyEditLog(t *testing.T) {
	initVersionSetTestData()
	ctrl := gomock.NewController(t)
	defer func() {
		destroyVersionTestData()
		ctrl.Finish()
	}()
	cache := table.NewMockCache(ctrl)
	cache.EXPECT().ReleaseReaders(gomock.Any()).AnyTimes()

	var vs = NewStoreVersionSet(vsTestPath, cache, 2)
	assert.NotNil(t, vs, "cannot create store version")
	var err = vs.Recover()
	assert.Nil(t, err, "recover error")

	err = vs.CommitFamilyEditLog("f", nil)
	assert.NotNil(t, err, "commit not exist family version")

	familyID := FamilyID(1)
	vs.CreateFamilyVersion("f", familyID)
	editLog := NewEditLog(familyID)
	nFile := CreateNewFile(1, NewFileMeta(12, 1, 100, 2014))
	nf := nFile.(*newFile)
	editLog.Add(nFile)
	editLog.Add(NewDeleteFile(1, 123))
	editLog.Add(CreateSequence(1, 10))
	editLog.Add(CreateNewRollupFile(1, 10000))
	editLog.Add(CreateNewReferenceFile("20230202", 1, 10))
	err = vs.CommitFamilyEditLog("f", editLog)
	assert.Nil(t, err, "commit family edit log error")

	_ = vs.Destroy()

	// test recover many times
	for i := 0; i < 3; i++ {
		vs = NewStoreVersionSet(vsTestPath, cache, 2)
		vs.CreateFamilyVersion("f", familyID)
		err = vs.Recover()
		assert.Nil(t, err, "recover error")

		familyVersion := vs.GetFamilyVersion("f")
		snapshot := familyVersion.GetSnapshot()
		vs1 := vs.(*storeVersionSet)
		current := snapshot.GetCurrent()
		assert.Equal(t, nf.file, current.GetAllFiles()[0], "cannot recover family version data")
		assert.Equal(t, int64(3+i), vs1.nextFileNumber.Load(), "recover file number error")
		assert.Equal(t, map[int32]int64{1: 10}, current.GetSequences())
		assert.Equal(t, map[FamilyID][]table.FileNumber{1: {10}}, current.GetReferenceFiles("20230202"))
		assert.Equal(t, map[table.FileNumber][]timeutil.Interval{1: {10000}}, current.GetRollupFiles())

		snapshot.Close()

		_ = vs.Destroy()
	}
}

func TestStoreVersionSet_CommitFamilyEditLog_err(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	cache := table.NewMockCache(ctrl)
	vs := NewStoreVersionSet(vsTestPath, cache, 2)
	familyID := FamilyID(1)
	vs.CreateFamilyVersion("f", familyID)
	manifest := bufioutil.NewMockBufioWriter(ctrl)
	versionSet := vs.(*storeVersionSet)
	versionSet.manifest = manifest

	// case 1: persist edit log encode err
	log := NewMockLog(ctrl)
	editLog := NewEditLog(familyID)
	editLog.Add(log)
	log.EXPECT().Encode().Return(nil, fmt.Errorf("err"))
	err := vs.CommitFamilyEditLog("f", editLog)
	assert.Error(t, err)
	// case 2: write manifest err
	log.EXPECT().Encode().Return([]byte{1, 2, 3}, nil).AnyTimes()
	manifest.EXPECT().Write(gomock.Any()).Return(0, fmt.Errorf("err"))
	err = vs.CommitFamilyEditLog("f", editLog)
	assert.Error(t, err)
	// case 4: sync manifest err
	manifest.EXPECT().Write(gomock.Any()).Return(10, nil)
	manifest.EXPECT().Sync().Return(fmt.Errorf("err"))
	err = vs.CommitFamilyEditLog("f", editLog)
	assert.Error(t, err)
}

func TestCreateFamily(t *testing.T) {
	initVersionSetTestData()
	ctrl := gomock.NewController(t)
	defer func() {
		destroyVersionTestData()
		ctrl.Finish()
	}()
	cache := table.NewMockCache(ctrl)

	var vs = NewStoreVersionSet(vsTestPath, cache, 2)

	familyVersion := vs.CreateFamilyVersion("family", 1)
	assert.NotNil(t, familyVersion, "get nil family version")
	familyVersion2 := vs.CreateFamilyVersion("family", 1)
	assert.Equal(t, familyVersion, familyVersion2)

	familyVersion2 = vs.GetFamilyVersion("family")
	assert.NotNil(t, familyVersion2, "get nil family version2")

	assert.Equal(t, familyVersion, familyVersion2, "get diff family version")
}

func TestStoreVersionSet_Destroy(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cache := table.NewMockCache(ctrl)
	vs := NewStoreVersionSet(vsTestPath, cache, 2)
	versionSet := vs.(*storeVersionSet)
	manifest := bufioutil.NewMockBufioWriter(ctrl)
	versionSet.manifest = manifest
	manifest.EXPECT().Close().Return(fmt.Errorf("err"))
	err := vs.Destroy()
	assert.Error(t, err)
}

func TestStoreVersionSet(t *testing.T) {
	path := t.TempDir()
	cache := table.NewCache(path, time.Minute)
	vs := NewStoreVersionSet(path, cache, 2)
	err := vs.Recover()
	assert.NoError(t, err)
	defer func() {
		assert.NoError(t, vs.Destroy())
	}()
	familyID := FamilyID(10)
	familyName := "10"
	fv := vs.CreateFamilyVersion(familyName, familyID)
	editLog := NewEditLog(familyID)
	editLog.Add(CreateNewRollupFile(table.FileNumber(10), timeutil.Interval(10000)))
	editLog.Add(CreateNewRollupFile(table.FileNumber(10), timeutil.Interval(30000)))
	err = fv.GetVersionSet().CommitFamilyEditLog(familyName, editLog)
	assert.NoError(t, err)

	editLog = NewEditLog(familyID)
	editLog.Add(CreateNewFile(0, &FileMeta{
		fileNumber: 1,
		minKey:     0,
		maxKey:     10,
		fileSize:   512,
	}))
	err = fv.GetVersionSet().CommitFamilyEditLog(familyName, editLog)
	assert.NoError(t, err)
	assert.Len(t, fv.GetLiveRollupFiles(), 1)
}

func TestStoreVersionSet_NextFileNumber(t *testing.T) {
	path := t.TempDir()
	cache := table.NewCache(path, time.Minute)
	vs := NewStoreVersionSet(path, cache, 2)
	err := vs.Recover()
	assert.NoError(t, err)
	defer func() {
		assert.NoError(t, vs.Destroy())
	}()
	familyID := FamilyID(10)
	familyName := "10"
	fv := vs.CreateFamilyVersion(familyName, familyID)
	assert.NotNil(t, fv)
	var wait sync.WaitGroup
	var rs sync.Map
	total := 24
	wait.Add(total)
	for i := 0; i < total; i++ {
		go func() {
			defer wait.Done()
			fn := vs.NextFileNumber()
			editLog := NewEditLog(familyID)
			editLog.Add(CreateNewFile(0, &FileMeta{
				fileNumber: fn,
				minKey:     0,
				maxKey:     10,
				fileSize:   10,
			}))
			time.Sleep(time.Millisecond)
			err := vs.CommitFamilyEditLog(familyName, editLog)
			assert.NoError(t, err)
			rs.Store(fn.Int64(), fn.Int64())
		}()
	}
	wait.Wait()
	c := 0
	rs.Range(func(_, _ any) bool {
		c++
		return true
	})
	assert.Equal(t, total, c)
}

func initVersionSetTestData() {
	if err := fileutil.MkDirIfNotExist(vsTestPath); err != nil {
		fmt.Println("create test path error")
	}
}

func destroyVersionTestData() {
	if err := fileutil.RemoveDir(vsTestPath); err != nil {
		fmt.Println("delete test path error")
	}
}
