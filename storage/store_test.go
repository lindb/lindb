package storage

import (
	"testing"
)

func Test_Create_Family(t *testing.T) {
	option := StoreOption{Path: "../test_data"}
	var kv, err = NewStore("test_kv", option)
	defer kv.Close()
	if nil != err {
		t.Error(err)
		return
	}
	_, err = kv.CreateFamily("f", FamilyOption{})
	if nil != err {
		t.Error(err)
		return
	}

	var _, ok = kv.GetFamily("f")
	if !ok {
		t.Errorf("can't get family")
		return
	}

	_, ok = kv.GetFamily("f1")
	if ok {
		t.Errorf("fail:get no exist family")
		return
	}

	_, e := NewStore("test_kv", option)
	if e == nil {
		t.Errorf("re-open not allow")
		return
	}
}
