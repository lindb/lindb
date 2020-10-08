package series

import "github.com/lindb/lindb/pkg/collections"

//go:generate mockgen -source ./block.go -destination=./block_mock.go -package series

// Block represents series block which stores data points
type Block interface {
	// Append appends time slot and value into block
	Append(slot int, value float64) bool
	// Clear clears the values of block.
	Clear()
}

// block implements Block interface
type block struct {
	start, end int
	values     collections.FloatArray
}

// NewBlock creates a new block with start/end time slot
func NewBlock(start, end int) Block {
	return &block{
		start:  start,
		end:    end,
		values: collections.NewFloatArray(end - start + 1),
	}
}

// Append appends time slot and value into block
func (b *block) Append(slot int, value float64) bool {
	if slot > b.end {
		return true
	}
	if slot < b.start {
		return false
	}
	b.values.SetValue(slot, value)
	return false
}

// Clear clears the values of block.
func (b *block) Clear() {
	b.values.Reset()
}
