package encoding

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Write(t *testing.T) {
	e := NewFloatEncoder()
	e.Write(76.1)
	e.Write(50.0)
	e.Write(50.0)
	e.Write(999999999.0)
	e.Write(0.0099)
	e.Write(100.0099)

	data, err := e.Bytes()
	assert.Nil(t, err)
	d := NewFloatDecoder(data)
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
	for i := 0; i < loop; i++ {
		e.Write(76.1)
	}

	data, err := e.Bytes()
	assert.Nil(t, err)

	d := NewFloatDecoder(data)
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
}

func exceptValue(d *FloatDecoder, t *testing.T, except float64) {
	assert.True(t, d.Next())
	assert.Equal(t, except, d.Value())
}
