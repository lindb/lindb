package encoding

import "testing"

func TestZigZag(t *testing.T) {
	var v = ZigZagEncode(1)
	if ZigZagDecode(v) != 1 {
		t.Errorf("error")
	}

	v = ZigZagEncode(-99999)
	if ZigZagDecode(v) != -99999 {
		t.Errorf("error")
	}
}
