package kv

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/pkg/fileutil"
)

var testKVPath = "../test_data"

func TestReOpen(t *testing.T) {
	option := DefaultStoreOption(testKVPath)
	defer fileutil.RemoveDir(testKVPath)

	var kv, _ = NewStore("test_kv", option)
	assert.NotNil(t, kv, "cannot create kv store")
	_, e := NewStore("test_kv", option)
	assert.NotNil(t, e, "store re-open not allow")

	kv, _ = kv.(*store)

	f1, _ := kv.CreateFamily("f", FamilyOption{})
	assert.NotNil(t, f1, "cannot create family")

	kvStore, ok := kv.(*store)
	if ok {
		assert.Equal(t, 1, kvStore.familyID, "store family id is wrong")
	}
	assert.True(t, ok)

	kv.Close()

	kv2, e := NewStore("test_kv", option)
	if e != nil {
		t.Error(e)
	}
	assert.NotNil(t, kv2, "cannot re-open kv store")
	f2 := kv.GetFamily("f")
	assert.Equal(t, f1.Name(), f2.Name(), "family diff when store reopen")
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
}

func TestCreateFamily(t *testing.T) {
	option := DefaultStoreOption("../test_data")
	defer fileutil.RemoveDir(testKVPath)

	var kv, err = NewStore("test_kv", option)
	defer kv.Close()
	assert.Nil(t, err, "cannot create kv store")

	f1, err2 := kv.CreateFamily("f", FamilyOption{})
	assert.Nil(t, err2, "cannot create family")

	var f2 = kv.GetFamily("f")
	assert.Equal(t, f1, f2, "family not same for same name")

	f11 := kv.GetFamily("f11")
	assert.Nil(t, f11)

	_, e := NewStore("test_kv", option)
	assert.NotNil(t, e, "store re-open not allow")
}
