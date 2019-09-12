package encoding

import (
	"math"
	"math/rand"
	"sort"
	"testing"
	"time"

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

	d := NewDeltaBitPackingDecoder(b)

	count := 0
	for d.HasNext() {
		_ = d.Next()
		count++
	}
	assert.Equal(t, 104, count)
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

func Test_DeltaBitPackingEncoder_Decoder(t *testing.T) {
	p := NewDeltaBitPackingEncoder()
	d := NewDeltaBitPackingDecoder(nil)

	for range [10]struct{}{} {
		p.Reset()
		list := getRandomList()
		for _, v := range list {
			p.Add(v)
		}
		b := p.Bytes()

		d.Reset(b)
		var count = 0
		for d.HasNext() {
			value := d.Next()
			assert.Equal(t, list[count], value)
			count++
		}
	}
}

func getRandomList() []int32 {
	var list []int32

	rand.Seed(time.Now().UnixNano())
	for i := 0; i < 1000; i++ {
		list = append(list, rand.Int31n(math.MaxInt32))
	}
	sort.Slice(list, func(i, j int) bool {
		return list[i] < list[j]
	})
	return list
}
