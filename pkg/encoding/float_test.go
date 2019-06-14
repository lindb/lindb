package encoding

import (
	"fmt"
	"testing"
	"time"
)

func Test_Write(t *testing.T) {
	e := NewFloatEncoder()
	e.Write(76.1)
	e.Write(50.0)
	e.Write(50.0)
	e.Write(999999999.0)
	e.Write(0.0099)
	e.Write(100.0099)

	e.Flush()

	d := NewFloatDecoder(e.Bytes())
	exceptValue(d, t, 76.1)
	exceptValue(d, t, 50.0)
	exceptValue(d, t, 50.0)
	exceptValue(d, t, 999999999.0)
	exceptValue(d, t, 0.0099)
	exceptValue(d, t, 100.0099)
}

func Benchmark_Write(t *testing.B) {
	e := NewFloatEncoder()
	loop := 1000000
	var now = time.Now().UnixNano() / 1000000
	for i := 0; i < loop; i++ {
		e.Write(76.1)
	}
	fmt.Printf("encode cost: %d \n", time.Now().UnixNano()/1000000-now)

	e.Flush()

	fmt.Printf("data size:%d \n", len(e.Bytes()))

	now = time.Now().UnixNano() / 1000000
	d := NewFloatDecoder(e.Bytes())
	for i := 0; i < loop; i++ {
		has := d.Next()
		f := d.Value()
		if !has {
			t.Errorf("haven't value")
		}
		if f != 76.1 {
			t.Errorf("get wrong value")
		}
	}
	fmt.Printf("decode cost: %d \n", time.Now().UnixNano()/1000000-now)
}

func exceptValue(d *FloatDecoder, t *testing.T, except float64) {
	if has := d.Next(); !has {
		t.Errorf("haven't value")
	}
	if f := d.Value(); f != except {
		t.Errorf("get wrong value, exception %f, actul: %f", except, f)
	}
}
