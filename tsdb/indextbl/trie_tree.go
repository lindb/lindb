package indextbl

import (
	"math"
	"sort"
	"sync"

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
	// Add adds the key to the Trie with the id
	Add(key string, id uint32)
	// NodeNum returns the size of nodes
	NodeNum() int
	// KeyNum returns the size of keys
	KeyNum() int
	// MarshalBinary marshals the trie tree to LOUDS encoded binary.
	MarshalBinary() *loudsBinary
	// Reset returns the nodes to the pool for reuse.
	Reset()
}

// trieTree implements the trieTreeINTF interface.
// Not thread-safe
type trieTree struct {
	pool      sync.Pool
	root      *node
	nodeNum   int     // count of nodes
	keyNum    int     // count of keys
	nodesBuf1 []*node // buffering nodes for Breadth first search
	nodesBuf2 []*node // buffering nodes for Breadth first search
}

// newTrieTree returns an empty trie tree.
func newTrieTree() trieTreeINTF {
	tt := &trieTree{
		pool: sync.Pool{
			New: func() interface{} {
				return &node{value: math.MaxUint32}
			},
		},
		nodesBuf1: make([]*node, 0, 4096),
		nodesBuf2: make([]*node, 0, 4096),
	}
	tt.root = tt.newNode() // empty string
	return tt
}

// nodes implements sort.Interface
type nodes []*node

func (n nodes) Len() int           { return len(n) }
func (n nodes) Swap(i, j int)      { n[i], n[j] = n[j], n[i] }
func (n nodes) Less(i, j int) bool { return n[i].label < n[j].label }

// node is the tree node of the trie
type node struct {
	label       byte
	value       uint32 // value for seriesID and metricID
	isPrefixKey bool   // isPrefixKey indicates if this prefix a valid key
	children    nodes  // sorted children nodes list
}

// newNode returns a node with default value
func (tt *trieTree) newNode() *node {
	return tt.pool.Get().(*node)
}

// Add adds a new key to the tree with id.
func (tt *trieTree) Add(key string, id uint32) {
	if key == "" {
		return
	}
	n := tt.root
	for _, k := range []byte(key) {
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
	n.value = id
	n.isPrefixKey = true
}

// NodeNum returns the size of nodes
func (tt *trieTree) NodeNum() int { return tt.nodeNum }

// KeyNum returns the size of keys
func (tt *trieTree) KeyNum() int { return tt.keyNum }

// loudsBinary is a LOUDS Encoded trie
type loudsBinary struct {
	LOUDS       *rsdic.RSDic
	isPrefixKey *rsdic.RSDic
	labels      []byte
	values      []uint32
}

// MarshalBinary marshals the trie tree to LOUDS encoded binary.
func (tt *trieTree) MarshalBinary() *loudsBinary {
	LOUDS := rsdic.New()
	isPrefixKey := rsdic.New()

	var (
		labels []byte
		values []uint32
	)
	tt.nodesBuf1 = []*node{tt.root}

	for len(tt.nodesBuf1) > 0 {
		for _, node := range tt.nodesBuf1 {
			if node.isPrefixKey {
				values = append(values, node.value)
			}
			isPrefixKey.PushBack(node.isPrefixKey)
			labels = append(labels, node.label)

			for _, node := range node.children {
				tt.nodesBuf2 = append(tt.nodesBuf2, node)
				LOUDS.PushBack(true)
			}
			LOUDS.PushBack(false)
		}
		// copy
		tt.nodesBuf1 = append(tt.nodesBuf1[:0], tt.nodesBuf2...)
		tt.nodesBuf2 = tt.nodesBuf2[:0]
	}

	return &loudsBinary{
		LOUDS:       LOUDS,
		isPrefixKey: isPrefixKey,
		labels:      labels,
		values:      values}
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
			n.value = math.MaxUint32
			n.isPrefixKey = false
			n.label = byte(0)
			n.children = n.children[:0]
			tt.pool.Put(n)
		}
		tt.nodesBuf1 = append(tt.nodesBuf1[:0], tt.nodesBuf2...)
		tt.nodesBuf2 = tt.nodesBuf2[:0]
	}

	tt.root = tt.newNode()
	tt.keyNum = 0
	tt.nodeNum = 0
	tt.nodesBuf1 = tt.nodesBuf1[:0]
	tt.nodesBuf2 = tt.nodesBuf2[:0]
}
