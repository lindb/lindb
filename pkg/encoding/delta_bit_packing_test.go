package encoding

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_DeltaBitPackingEncoder_Add(t *testing.T) {
	p := NewDeltaBitPackingEncoder()

	p.Add(1)
	p.Add(10)
	p.Add(1)
	for i := 0; i < 100; i++ {
		p.Add(100)
	}

	p.Add(200)

	b := p.Bytes()

	t.Logf("xx,%p\n", &b)

	d := NewDeltaBitPackingDecoder(b)

	count := 0
	for d.HasNext() {
		_ = d.Next()
		count++
	}
	assert.Equal(t, 104, count)

	t.Logf("xx,%p", &d)
}

func Test_DeltaBitPackingEncoder_Reset(t *testing.T) {
	p := NewDeltaBitPackingEncoder()
	for i := 0; i < 100; i++ {
		p.Add(100)
	}
	b1 := p.Bytes()
	p.Reset()
	for i := 0; i < 100; i++ {
		p.Add(100)
	}
	b2 := p.Bytes()
	assert.Equal(t, b1, b2)

}
