package kv

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/kv/version"
	"github.com/lindb/lindb/pkg/fileutil"
)

func init() {
	RegisterMerger("mockMerger", newMockMerger)
}

func Test_Data_Write_Read(t *testing.T) {
	option := DefaultStoreOption(testKVPath)
	defer func() {
		_ = fileutil.RemoveDir(testKVPath)
	}()

	var kv, err = NewStore("test_kv", option)
	defer func() {
		_ = kv.Close()
	}()
	assert.Nil(t, err, "cannot create kv store")

	f, err := kv.CreateFamily("f", FamilyOption{Merger: "mockMerger"})
	assert.Nil(t, err, "cannot create family")
	flusher := f.NewFlusher()
	_ = flusher.Add(1, []byte("test"))
	_ = flusher.Add(10, []byte("test10"))
	commitErr := flusher.Commit()
	assert.Nil(t, commitErr)

	snapshot := f.GetSnapshot()
	readers, _ := snapshot.FindReaders(10)
	assert.Equal(t, 1, len(readers))
	value, _ := readers[0].Get(1)
	assert.Equal(t, []byte("test"), value)
	value, _ = readers[0].Get(10)
	assert.Equal(t, []byte("test10"), value)
	snapshot.Close()
}

func TestCommitEditLog(t *testing.T) {
	option := DefaultStoreOption(testKVPath)
	defer func() {
		_ = fileutil.RemoveDir(testKVPath)
	}()

	var kv, _ = NewStore("test_kv", option)
	defer func() {
		_ = kv.Close()
	}()

	f, _ := kv.CreateFamily("f", FamilyOption{Merger: "mockMerger"})

	editLog := version.NewEditLog(1)
	newFile := version.CreateNewFile(1, version.NewFileMeta(12, 1, 100, 2014))
	editLog.Add(newFile)
	editLog.Add(version.NewDeleteFile(1, 123))

	family, ok := f.(*family)
	if ok {
		flag := family.commitEditLog(editLog)
		assert.True(t, flag, "commit edit log error")
	}
	assert.True(t, ok)
}
