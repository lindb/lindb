package indextbl

import (
	"sort"
	"sync"

	"github.com/RoaringBitmap/roaring"
	"github.com/hillbig/rsdic"
)

// Implementation of an R-way Trie data structure.
// Ref: https://en.wikipedia.org/wiki/Trie

// A succinct data structure supporting rank/select efficiently is used here for querying and filtering data.
// Definition: https://en.wikipedia.org/wiki/Succinct_data_structure
// [1] SuRF: Practical Range Query Filtering with Fast Succinct Tries:
//     https://db.cs.cmu.edu/papers/2018/mod601-zhangA-hm.pdf
// [2] Fast, Small, Simple Rank/Select on Bitmaps:
//     https://users.dcc.uchile.cl/~gnavarro/ps/sea12.1.pdf

//go:generate mockgen -source ./trie_tree.go -destination=./trie_tree_mock.go -package indextbl

// trieTreeINTF abstract a trie tree in memory.
// All the descendants of a node have a common prefix of the string associated with that node,
// and the root node is associated with the empty string.
type trieTreeINTF interface {
	// Add adds the tagValue and bitmap to the Trie with the id
	Add(tagValue string, version int64, bitmap *roaring.Bitmap)
	// NodeNum returns the size of nodes
	NodeNum() int
	// KeyNum returns the size of keys
	KeyNum() int
	// MarshalBinary marshals the trie tree to LOUDS encoded binary.
	MarshalBinary() *loudsBinary
	// Reset returns the nodes to the pool for reuse.
	Reset()
}

// rankSelectINTF abstracts the rank select implementation(RSDIC).
type rankSelectINTF interface {
	// Num returns the number of bits
	Num() uint64
	// OneNum returns the number of ones in bits
	OneNum() uint64
	// ZeroNum returns the number of zeros in bits
	ZeroNum() uint64
	// PushBack appends the bit to the end of B
	PushBack(bit bool)
	// Bit returns the (pos+1)-th bit in bits, i.e. bits[pos]
	Bit(pos uint64) bool
	// Rank returns the number of bit's in B[0...pos)
	Rank(pos uint64, bit bool) uint64
	// Select returns the position of (rank+1)-th occurrence of bit in B
	// Select returns num if rank+1 is larger than the possible range.
	// (i.e. Select(oneNum, true) = num, Select(zeroNum, false) = num)
	Select(rank uint64, bit bool) uint64
	Select1(rank uint64) uint64
	Select0(rank uint64) uint64
	// BitAndRank returns the (pos+1)-th bit (=b) and Rank(pos, b)
	// Although this is equivalent to b := Bit(pos), r := Rank(pos, b),
	// BitAndRank is faster.
	BitAndRank(pos uint64) (bool, uint64)
	// AllocSize returns the allocated size in bytes.
	AllocSize() int
	// MarshalBinary encodes the RSDic into a binary form and returns the result.
	MarshalBinary() (out []byte, err error)
	// UnmarshalBinary decodes the RSDic from a binary from generated MarshalBinary
	UnmarshalBinary(in []byte) (err error)
}

// trieTree implements the trieTreeINTF interface.
// Not thread-safe
type trieTree struct {
	pool      sync.Pool
	root      *node
	nodeNum   int          // count of nodes
	keyNum    int          // count of keys
	nodesBuf1 []*node      // buffering nodes for Breadth first search
	nodesBuf2 []*node      // buffering nodes for Breadth first search
	bin       *loudsBinary // binary
}

// newTrieTree returns an empty trie tree.
func newTrieTree() trieTreeINTF {
	tt := &trieTree{
		pool:      sync.Pool{New: func() interface{} { return &node{} }},
		nodesBuf1: make([]*node, 0, 4096),
		nodesBuf2: make([]*node, 0, 4096),
	}
	tt.root = tt.newNode() // empty string
	tt.bin = new(loudsBinary)
	return tt
}

// nodes implements sort.Interface
type nodes []*node

func (n nodes) Len() int           { return len(n) }
func (n nodes) Swap(i, j int)      { n[i], n[j] = n[j], n[i] }
func (n nodes) Less(i, j int) bool { return n[i].label < n[j].label }

// versionedBitmap is the node value
type versionedBitmap struct {
	version int64
	bitmap  *roaring.Bitmap
}

// versionedBitmaps implements sort.Interface
type versionedBitmaps []versionedBitmap

func (n versionedBitmaps) Len() int           { return len(n) }
func (n versionedBitmaps) Swap(i, j int)      { n[i], n[j] = n[j], n[i] }
func (n versionedBitmaps) Less(i, j int) bool { return n[i].version < n[j].version }

func (n *versionedBitmaps) insert(version int64, bitmap *roaring.Bitmap) {
	idx := sort.Search(len(*n), func(i int) bool { return (*n)[i].version >= version })
	// existed before
	if idx < len(*n) && (*n)[idx].version == version {
		(*n)[idx].bitmap = bitmap
		return
	}
	vb := versionedBitmap{version: version, bitmap: bitmap}
	*n = append(*n, vb)
	sort.Sort(*n)
}

// reset set the bitmap pointer and slice length
func (n *versionedBitmaps) reset() {
	for _, item := range *n {
		item.bitmap = nil
	}
	*n = (*n)[:0]
}

// node is the tree node of the trie
type node struct {
	label       byte
	values      versionedBitmaps // versioned bitmap slice
	isPrefixKey bool             // isPrefixKey indicates if this prefix a valid key
	children    nodes            // sorted children nodes list
}

// newNode returns a node with default value
func (tt *trieTree) newNode() *node {
	return tt.pool.Get().(*node)
}

// Add adds a new key to the tree with id.
func (tt *trieTree) Add(tagValue string, version int64, bitmap *roaring.Bitmap) {
	if tagValue == "" {
		return
	}
	n := tt.root
	for _, k := range []byte(tagValue) {
		k := k
		childIdx := sort.Search(len(n.children), func(i int) bool { return n.children[i].label >= k })
		// exist before
		if childIdx < len(n.children) && n.children[childIdx].label == k {
			n = n.children[childIdx]
			continue
		} else {
			theNewNode := tt.newNode()
			theNewNode.label = k
			n.children = append(n.children, theNewNode)
			sort.Sort(n.children)
			n = theNewNode
			tt.nodeNum++
		}
	}
	if !n.isPrefixKey {
		tt.keyNum++
	}
	// set the term node
	n.values.insert(version, bitmap)
	n.isPrefixKey = true
}

// NodeNum returns the size of nodes
func (tt *trieTree) NodeNum() int { return tt.nodeNum }

// KeyNum returns the size of keys
func (tt *trieTree) KeyNum() int { return tt.keyNum }

// loudsBinary is a LOUDS Encoded trie
type loudsBinary struct {
	LOUDS       rankSelectINTF
	isPrefixKey rankSelectINTF
	labels      []byte
	values      []versionedBitmaps
}

func (lb *loudsBinary) reset() {
	lb.LOUDS = nil
	lb.isPrefixKey = nil
	lb.labels = lb.labels[:0]
	for _, vb := range lb.values {
		vb.reset()
	}
	lb.values = lb.values[:0]
}

// MarshalBinary marshals the trie tree to LOUDS encoded binary.
func (tt *trieTree) MarshalBinary() *loudsBinary {
	tt.bin.LOUDS = rsdic.New()
	tt.bin.isPrefixKey = rsdic.New()
	tt.nodesBuf1 = []*node{tt.root}

	for len(tt.nodesBuf1) > 0 {
		for _, node := range tt.nodesBuf1 {
			if node.isPrefixKey {
				tt.bin.values = append(tt.bin.values, node.values)
			}
			tt.bin.isPrefixKey.PushBack(node.isPrefixKey)
			tt.bin.labels = append(tt.bin.labels, node.label)

			for _, node := range node.children {
				tt.nodesBuf2 = append(tt.nodesBuf2, node)
				tt.bin.LOUDS.PushBack(true)
			}
			tt.bin.LOUDS.PushBack(false)
		}
		// copy
		tt.nodesBuf1 = append(tt.nodesBuf1[:0], tt.nodesBuf2...)
		tt.nodesBuf2 = tt.nodesBuf2[:0]
	}
	return tt.bin
}

// Reset returns the nodes to the underlying pool for reuse.
func (tt *trieTree) Reset() {
	tt.nodesBuf1 = []*node{tt.root}
	tt.nodesBuf2 = tt.nodesBuf2[:0]

	for len(tt.nodesBuf1) > 0 {
		for _, n := range tt.nodesBuf1 {
			for _, child := range n.children {
				tt.nodesBuf2 = append(tt.nodesBuf2, child)
			}
			n.values.reset()
			n.isPrefixKey = false
			n.label = byte(0)
			n.children = n.children[:0]
			tt.pool.Put(n)
		}
		tt.nodesBuf1 = append(tt.nodesBuf1[:0], tt.nodesBuf2...)
		tt.nodesBuf2 = tt.nodesBuf2[:0]
	}
	tt.bin.reset()
	tt.root = tt.newNode()
	tt.keyNum = 0
	tt.nodeNum = 0
	tt.nodesBuf1 = tt.nodesBuf1[:0]
	tt.nodesBuf2 = tt.nodesBuf2[:0]
}
