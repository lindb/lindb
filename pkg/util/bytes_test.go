package util

import (
	"fmt"
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIntToShort(t *testing.T) {
	var success = 0
	for i := 0; i < 10; i++ {
		for j := 0; j < math.MaxUint16; j++ {
			newValue := ShortToInt(uint16(i), uint16(j))
			high, low := IntToShort(newValue)
			assert.Equal(t, uint16(i), high)
			assert.Equal(t, uint16(j), low)
			success++
		}
	}
	fmt.Println("success:", success)
}

func TestUint32ToBytes(t *testing.T) {
	number := uint32(0)
	by := Uint32ToBytes(number)
	uInt := BytesToUint32(by)
	assert.Equal(t, number, uInt)
}
