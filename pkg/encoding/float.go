package encoding

import (
	"fmt"
	"math"
)

// FloatEncoder encodes float using xor compress
type FloatEncoder struct {
	encoder *XOREncoder
}

// NewFloatEncoder create float encoder
func NewFloatEncoder() *FloatEncoder {
	return &FloatEncoder{
		encoder: NewXOREncoder(),
	}
}

// Write writes a float value into xor stream, if fail return error
func (e *FloatEncoder) Write(v float64) error {
	if math.IsNaN(v) {
		return fmt.Errorf("unspoorted value for NaN")
	}

	val := math.Float64bits(v)
	if err := e.encoder.Write(val); err != nil {
		return err
	}
	return nil
}

// Bytes returns xor compress result, return error if fail
func (e *FloatEncoder) Bytes() ([]byte, error) {
	return e.encoder.Bytes()
}

// FloatDecoder decodes float using xor compress
type FloatDecoder struct {
	decoder *XORDecoder
}

// NewFloatDecoder creates float decoder using xor compress
func NewFloatDecoder(b []byte) *FloatDecoder {
	return &FloatDecoder{
		decoder: NewXORDecoder(b),
	}
}

// Next returns if has value in xor stream, return true has value
func (d *FloatDecoder) Next() bool {
	return d.decoder.Next()
}

// Value returns a float value in xor stream
func (d *FloatDecoder) Value() float64 {
	return math.Float64frombits(d.decoder.Value())
}
