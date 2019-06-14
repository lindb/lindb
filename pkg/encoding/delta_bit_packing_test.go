package encoding

import (
	"fmt"
	"testing"
)

func TestAdd(t *testing.T) {
	p := NewDeltaBitPackingEncoder()

	p.Add(1)
	p.Add(10)
	p.Add(1)
	for i := 0; i < 100; i++ {
		p.Add(100)
	}

	p.Add(200)

	b, err := p.Bytes()
	if err != nil {
		t.Error(err)
	}
	fmt.Printf("xx,%p\n", &b)

	d := NewDeltaBitPackingDecoder(&b)

	for d.HasNext() {
		x := d.Next()
		fmt.Printf("first=%d\n", x)
	}

	fmt.Printf("xx,%p", &d)
}
