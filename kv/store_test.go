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

package kv

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/kv/table"
	"github.com/lindb/lindb/kv/version"
	"github.com/lindb/lindb/pkg/fileutil"
	"github.com/lindb/lindb/pkg/lockers"
	"github.com/lindb/lindb/pkg/ltoml"
)

var mergerStr = "mockMergerAppend"

func newMockMerger(flusher Flusher) (Merger, error) {
	return &mockAppendMerger{flusher: flusher}, nil
}

func init() {
	RegisterMerger(MergerType(mergerStr), newMockMerger)
}

func TestRegisterMerger(t *testing.T) {
	assert.Panics(t, func() {
		RegisterMerger("test", newMockMerger)
		RegisterMerger("test", newMockMerger)
	})
}

func TestStore_New(t *testing.T) {
	ctrl := gomock.NewController(t)
	tmpDir := filepath.Join(t.TempDir(), "test_data")
	option := DefaultStoreOption(tmpDir)
	defer func() {
		encodeTomlFunc = ltoml.EncodeToml
		mkDirFunc = fileutil.MkDir
		newVersionSetFunc = version.NewStoreVersionSet
		newFileLockFunc = lockers.NewFileLock
		ctrl.Finish()
	}()
	// case 1: create store dir err
	mkDirFunc = func(path string) error {
		return fmt.Errorf("err")
	}
	kv, err := NewStore("test_kv", option)
	assert.Error(t, err)
	assert.Nil(t, kv)
	// case 2: dump store option err
	lock := lockers.NewMockFileLock(ctrl)
	newFileLockFunc = func(fileName string) lockers.FileLock {
		return lock
	}
	lock.EXPECT().Lock().Return(nil)
	lock.EXPECT().Unlock().Return(fmt.Errorf("err"))
	mkDirFunc = fileutil.MkDir
	encodeTomlFunc = func(fileName string, v interface{}) error {
		return fmt.Errorf("err")
	}
	kv, err = NewStore("test_kv", option)
	assert.Error(t, err)
	assert.Nil(t, kv)
	_ = fileutil.RemoveDir(tmpDir)
	encodeTomlFunc = ltoml.EncodeToml
	newFileLockFunc = lockers.NewFileLock
	// case 3: new store success
	kv, err = NewStore("test_kv", option)
	assert.NoError(t, err)
	assert.NotNil(t, kv, "cannot create kv store")
	// case 4: new store fail, because try lock file err
	_, err = NewStore("test_kv", option)
	assert.Error(t, err)
	kv, _ = kv.(*store)
	kvStore, ok := kv.(*store)
	_, err = kv.CreateFamily("f", FamilyOption{Merger: mergerStr})
	assert.NoError(t, err, "cannot create family")
	if ok {
		assert.Equal(t, 1, kvStore.familySeq, "store family id is wrong")
	}
	assert.True(t, ok)
	_ = kv.Close()
	// case 5: reopen store
	kv2, e := NewStore("test_kv", option)
	assert.NoError(t, e)
	assert.NotNil(t, kv2, "cannot re-open kv store")

	kvStore, ok = kv2.(*store)
	if ok {
		assert.Equal(t, 1, kvStore.familySeq, "store family id is wrong")
	}
	assert.True(t, ok)
	_ = kv2.Close()
	delete(mergers, MergerType(mergerStr))
	// case 6: decode option err
	_, e = NewStore("test_kv", option)
	assert.NotNil(t, e)
	assert.Nil(t, nil)
	RegisterMerger(MergerType(mergerStr), newMockMerger)

	_ = ioutil.WriteFile(filepath.Join(tmpDir, version.Options), []byte("err"), 0644)
	kv, e = NewStore("test_kv", option)
	assert.Error(t, e)
	assert.Nil(t, kv)
	// case 7: recover version err
	_ = fileutil.RemoveDir(tmpDir)
	vs := version.NewMockStoreVersionSet(ctrl)
	newVersionSetFunc = func(storePath string, storeCache table.Cache, numOfLevels int) version.StoreVersionSet {
		return vs
	}
	vs.EXPECT().Recover().Return(fmt.Errorf("err"))
	vs.EXPECT().Destroy().Return(nil) // close store
	vs.EXPECT().ManifestFileNumber().Return(table.FileNumber(10))
	_, e = NewStore("test_kv", option)
	assert.Error(t, e)
	assert.Nil(t, kv)
}

func TestStore_CreateFamily(t *testing.T) {
	option := DefaultStoreOption(filepath.Join(t.TempDir(), "test_data"))
	defer func() {
		encodeTomlFunc = ltoml.EncodeToml
		newFamilyFunc = newFamily
	}()

	kv, err := NewStore("test_kv", option)
	defer func() {
		_ = kv.Close()
	}()
	assert.NoError(t, err, "cannot create kv store")

	// case 1: create family success
	f1, err := kv.CreateFamily("f", FamilyOption{Merger: mergerStr})
	assert.NoError(t, err, "cannot create family")
	// case 2: create family, but exist
	f2, err := kv.CreateFamily("f", FamilyOption{Merger: mergerStr})
	assert.NoError(t, err)
	assert.Equal(t, f1, f2)
	// case 3: get family
	f2 = kv.GetFamily("f")
	assert.Equal(t, f1, f2, "family not same for same name")
	// case 4: get family exist
	f11 := kv.GetFamily("f11")
	assert.Nil(t, f11)
	// case 5: toml dump err
	encodeTomlFunc = func(fileName string, v interface{}) error {
		return fmt.Errorf("err")
	}
	f2, err = kv.CreateFamily("f1_err", FamilyOption{Merger: mergerStr})
	assert.Error(t, err)
	assert.Nil(t, f2)
	// case 6: new exist family err
	encodeTomlFunc = ltoml.EncodeToml
	newFamilyFunc = func(store Store, option FamilyOption) (f Family, err error) {
		return nil, fmt.Errorf("err")
	}
	f2, err = kv.CreateFamily("f1_err", FamilyOption{Merger: mergerStr})
	assert.Error(t, err)
	assert.Nil(t, f2)
	// case 7: list family name
	names := kv.ListFamilyNames()
	assert.Len(t, names, 1)
	assert.Equal(t, "f", names[0])
}

func TestStore_deleteObsoleteFiles(t *testing.T) {
	option := DefaultStoreOption(filepath.Join(t.TempDir(), "test_data"))
	defer func() {
		listDirFunc = fileutil.ListDir
		removeFunc = os.Remove
	}()

	listDirFunc = func(path string) (strings []string, err error) {
		return nil, fmt.Errorf("err")
	}
	// case 1: list dir err
	kv, err := NewStore("test_kv", option)
	assert.NoError(t, err)
	err = kv.Close()
	assert.NoError(t, err)

	// case 2: remove file err
	listDirFunc = fileutil.ListDir
	removeFunc = func(name string) error {
		return fmt.Errorf("err")
	}
	kv, err = NewStore("test_kv", option)
	assert.NoError(t, err)
	err = kv.Close()
	assert.NoError(t, err)

}

func TestStore_Compact(t *testing.T) {
	option := DefaultStoreOption(filepath.Join(t.TempDir(), "test_data"))
	option.CompactCheckInterval = 1

	kv, err := NewStore("test_kv", option)
	assert.NoError(t, err)
	defer func() {
		_ = kv.Close()
	}()
	f1, err2 := kv.CreateFamily("f", FamilyOption{
		CompactThreshold: 2,
		Merger:           mergerStr,
		MaxFileSize:      1 * 1024 * 1024,
	})
	assert.Nil(t, err2, "cannot create family")

	for i := 0; i < 2; i++ {
		flusher := f1.NewFlusher()
		_ = flusher.Add(1, []byte("test"))
		_ = flusher.Add(10, []byte("test10"))
		commitErr := flusher.Commit()
		assert.Nil(t, commitErr)
	}
	time.Sleep(2 * time.Second)

	snapshot := f1.GetSnapshot()
	readers, err := snapshot.FindReaders(10)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(readers))
	value, _ := readers[0].Get(1)
	assert.Equal(t, []byte("testtest"), value)
	value, _ = readers[0].Get(10)
	assert.Equal(t, []byte("test10test10"), value)
	snapshot.Close()
}

func TestStore_Close(t *testing.T) {
	option := DefaultStoreOption(filepath.Join(t.TempDir(), "test_data"))
	option.CompactCheckInterval = 1
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	kv, err := NewStore("test_kv", option)
	assert.NoError(t, err)
	kv1 := kv.(*store)
	cache := table.NewMockCache(ctrl)
	cache.EXPECT().Close().Return(fmt.Errorf("err"))
	kv1.cache = cache
	vs := version.NewMockStoreVersionSet(ctrl)
	vs.EXPECT().Destroy().Return(fmt.Errorf("err"))
	kv1.versions = vs
	err = kv.Close()
	assert.NoError(t, err)
}

func TestStore_RegisterRollup(t *testing.T) {
	option := DefaultStoreOption(filepath.Join(t.TempDir(), "test_data"))
	option.CompactCheckInterval = 1
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	kv, err := NewStore("test_kv", option)
	assert.NoError(t, err)
	rollup := NewMockRollup(ctrl)
	kv.RegisterRollup(10, rollup)
	kv.RegisterRollup(10, rollup) // reject
	rollup2, ok := kv.getRollup(10)
	assert.True(t, ok)
	assert.Equal(t, rollup, rollup2)
}
