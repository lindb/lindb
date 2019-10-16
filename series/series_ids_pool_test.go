package series

import "testing"

func Test_Uint32Pool(t *testing.T) {
	item := Uint32Pool.Get()
	Uint32Pool.Put(item)
}
