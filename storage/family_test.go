package storage

import "testing"

func Test_NewTableBuilder(t *testing.T) {
	option := StoreOption{Path: "../test_data"}
	var kv, err = NewStore("test_kv", option)
	defer kv.Close()
	if nil != err {
		t.Error(err)
		return
	}
	f, err := kv.CreateFamily("f", FamilyOption{})
	if nil != err {
		t.Error(err)
		return
	}

	f.NewTableBuilder()
}
