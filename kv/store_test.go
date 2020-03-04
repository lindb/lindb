package kv

import (
	"io/ioutil"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/kv/version"
	"github.com/lindb/lindb/pkg/fileutil"
)

var testKVPath = "./test_data"
var mergerStr = "mockMergerAppend"

func newMockMerger() Merger {
	return &mockAppendMerger{}
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

func TestReOpen(t *testing.T) {
	option := DefaultStoreOption(testKVPath)
	defer func() {
		_ = fileutil.RemoveDir(testKVPath)
	}()

	var kv, _ = NewStore("test_kv", option)
	assert.NotNil(t, kv, "cannot create kv store")
	_, e := NewStore("test_kv", option)
	assert.NotNil(t, e, "store re-open not allow")

	kv, _ = kv.(*store)

	f1, _ := kv.CreateFamily("f", FamilyOption{Merger: mergerStr})
	assert.NotNil(t, f1, "cannot create family")
	names := kv.ListFamilyNames()
	assert.Len(t, names, 1)
	assert.Equal(t, "f", names[0])

	kvStore, ok := kv.(*store)
	if ok {
		assert.Equal(t, 1, kvStore.familyID, "store family id is wrong")
	}
	assert.True(t, ok)

	_ = kv.Close()

	kv2, e := NewStore("test_kv", option)
	assert.NoError(t, e)
	assert.NotNil(t, kv2, "cannot re-open kv store")
	f2 := kv.GetFamily("f")
	assert.Equal(t, f1.Name(), f2.Name(), "family diff when store reopen")
	names = kv.ListFamilyNames()
	assert.Len(t, names, 1)
	assert.Equal(t, "f", names[0])

	family, flag := f1.(*family)
	if flag {
		assert.Equal(t, family.option.ID, family.option.ID, "family id diff")
	}
	assert.True(t, flag)
	kvStore, ok = kv2.(*store)
	if ok {
		assert.Equal(t, 1, kvStore.familyID, "store family id is wrong")
	}
	assert.True(t, ok)
	_ = kv2.Close()
	delete(mergers, MergerType(mergerStr))
	_, e = NewStore("test_kv", option)
	assert.NotNil(t, e)
	assert.Nil(t, nil)
	RegisterMerger(MergerType(mergerStr), newMockMerger)

	_ = ioutil.WriteFile(filepath.Join(testKVPath, version.Options), []byte("err"), 0644)
	_, e = NewStore("test_kv", option)
	assert.Error(t, e)
	assert.Nil(t, nil)
}

func TestCreateFamily(t *testing.T) {
	option := DefaultStoreOption(testKVPath)
	defer func() {
		_ = fileutil.RemoveDir(testKVPath)
	}()

	var kv, err = NewStore("test_kv", option)
	defer func() {
		_ = kv.Close()
	}()
	assert.Nil(t, err, "cannot create kv store")

	f1, err2 := kv.CreateFamily("f", FamilyOption{Merger: mergerStr})
	assert.Nil(t, err2, "cannot create family")

	var f2 = kv.GetFamily("f")
	assert.Equal(t, f1, f2, "family not same for same name")

	f11 := kv.GetFamily("f11")
	assert.Nil(t, f11)

	// cannot re-open
	_, e := NewStore("test_kv", option)
	assert.NotNil(t, e)
}

func TestStore_Compact(t *testing.T) {
	option := DefaultStoreOption(testKVPath)
	option.CompactCheckInterval = 1
	defer func() {
		_ = fileutil.RemoveDir(testKVPath)
	}()

	kv, err := NewStore("test_kv", option)
	if err != nil {
		t.Fatal(err)
	}
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
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, 1, len(readers))
	value, _ := readers[0].Get(1)
	assert.Equal(t, []byte("testtest"), value)
	value, _ = readers[0].Get(10)
	assert.Equal(t, []byte("test10test10"), value)
	snapshot.Close()
}
