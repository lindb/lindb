package table

// Iterator iterates over a store's key/value pairs in key order.
type Iterator interface {
	// Next moves the iterator to the next key/value pair.
	// It returns false if the iterator is exhausted.
	Next() bool
	// Key returns the key of the current key/value pair
	Key() uint32
	// Value returns the value of the current key/value pair
	Value() []byte
}

// mergedIterator iteratores over some iterator in key order
// type mergedIterator struct {
// 	iters []Iterator

// 	key   uint32
// 	value [][]byte

// 	err error
// }

// newMergedIterator returns an iterator that merges its input.
// func newMergedIterator(iters []Iterator) *mergedIterator {
// 	return &mergedIterator{
// 		iters: iters,
// 	}
// }

// func (it *mergedIterator) Next() bool {
// 	return true
// }

// func (it *mergedIterator) Key() uint32 {
// 	return 1
// }
// func (it *mergedIterator) Value() []byte {
// 	return nil
// }
