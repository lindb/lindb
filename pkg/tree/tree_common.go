package tree

import (
	"bytes"
)

const (
	HasParent       = byte('1')
	NoParent        = byte('0')
	HasChildrenNode = byte('1')
	NoChildrenNode  = byte('0')
	NotFound        = -1
)

//Iterator returns an iterator over elements
type Iterator interface {
	//If the iteration has more elements.
	Next() bool

	//Get current element's key
	GetKey() []byte

	//Get current element's value
	GetValue() int
}

//Filter represents iterator's validation
type Filter interface {
	//Start condition of Iterator
	beginMatch(key []byte) bool
	//End condition of Iterator
	endMatch(key []byte) bool
}

//SeekFilter represents filter for Seek queries
type SeekFilter struct {
	prefix []byte // Target prefix
}

//beginMatch represents match start key
func (s *SeekFilter) beginMatch(key []byte) bool {
	return bytes.HasPrefix(key, s.prefix)
}

//endMatch represents Match end key
func (s *SeekFilter) endMatch(key []byte) bool {
	return bytes.HasPrefix(key, s.prefix)
}

//RangeFilter represents filter for range queries
type RangeFilter struct {
	startKey []byte //startKey
	endKey   []byte //endKey
}

func (s *RangeFilter) beginMatch(key []byte) bool {
	return BytesCompare(key, s.startKey) >= 0
}

func (s *RangeFilter) endMatch(key []byte) bool {
	return BytesCompare(key, s.endKey) <= 0
}

//SkipFilter represents always matching Filter
type SkipFilter struct {
}

func (s *SkipFilter) beginMatch(key []byte) bool {
	return true
}

func (s *SkipFilter) endMatch(key []byte) bool {
	return true
}

//Compare the size of two byte arrays
// if a > b return 1
// else if a == b return 0
// if a < b return -1
func BytesCompare(a, b interface{}) int {
	return bytes.Compare(a.([]byte), b.([]byte))
}

// lcp finds the longest common prefix of the input bytes.
// It compares by bytes instead of runes (Unicode code points).
// It's up to the caller to do Unicode normalization if desired
// (e.g. see golang.org/x/text/unicode/norm).
func lcpByte(l [][]byte) []byte {
	// Special cases first
	switch len(l) {
	case 0:
		return nil
	case 1:
		return l[0]
	}
	// LCP of min and max
	// is the LCP of the whole set.
	min, max := l[0], l[0]
	for _, s := range l[1:] {
		switch {
		//bytes.Compare(a.([]byte), b.([]byte))
		case bytes.Compare(s, min) < 0:
			min = s
		case bytes.Compare(s, max) > 0:
			max = s
		}
	}
	for i := 0; i < len(min) && i < len(max); i++ {
		if min[i] != max[i] {
			return min[:i]
		}
	}
	// In the case where lengths are not equal but all bytes
	// are equal, min is the answer ("foo" < "foobar").
	return min
}

//bytesCombine represents byte arrays are merged in order
func bytesCombine(pBytes ...[]byte) []byte {
	return bytes.Join(pBytes, []byte(""))
}
