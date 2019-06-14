package table

import "testing"

func Test_BuildStore(t *testing.T) {
	var builder, err = NewStoreBuilder("../../test_data/test_kv.sst")
	defer builder.Close()
	if err != nil {
		t.Error("new build error:", err)
		return
	}

	var result = builder.Add(1, []byte("test"))
	if !result {
		t.Error("write error")
		return
	}
	result = builder.Add(1, []byte("test"))
	if result {
		t.Error("write wrong data")
		return
	}

	NewReader("../../test_data/test_kv.sst")
}
