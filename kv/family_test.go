package kv

import (
	"testing"

	"github.com/eleme/lindb/kv/version"
	"github.com/eleme/lindb/pkg/util"

	"github.com/stretchr/testify/assert"
)

func Test_Data_Write_Read(t *testing.T) {
	option := DefaultStoreOption(testKVPath)
	defer util.RemoveDir(testKVPath)

	var kv, err = NewStore("test_kv", option)
	defer kv.Close()
	assert.Nil(t, err, "cannot create kv store")

	f, err := kv.CreateFamily("f", FamilyOption{})
	assert.Nil(t, err, "cannot create family")
	flusher := f.NewFlusher()
	_ = flusher.Add(1, []byte("test"))
	_ = flusher.Add(10, []byte("test10"))
	commitErr := flusher.Commit()
	assert.Nil(t, commitErr)

	snapshot, _ := f.GetSnapshot(10)
	readers := snapshot.Readers()
	assert.Equal(t, 1, len(readers))
	assert.Equal(t, []byte("test"), readers[0].Get(1))
	assert.Equal(t, []byte("test10"), readers[0].Get(10))
	snapshot.Close()
}

func TestCommitEditLog(t *testing.T) {
	option := DefaultStoreOption(testKVPath)
	defer util.RemoveDir(testKVPath)

	var kv, _ = NewStore("test_kv", option)
	defer kv.Close()

	f, _ := kv.CreateFamily("f", FamilyOption{})

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
