package indextbl

import (
	"sort"
	"sync"

	"github.com/RoaringBitmap/roaring"
)

// Implementation of an R-way Trie data structure.
// Ref: https://en.wikipedia.org/wiki/Trie

// A succinct data structure supporting rank/select efficiently is used here for querying and filtering data.
// Definition: https://en.wikipedia.org/wiki/Succinct_data_structure
// [1] SuRF: Practical Range Query Filtering with Fast Succinct Tries:
//     https://db.cs.cmu.edu/papers/2018/mod601-zhangA-hm.pdf
// [2] Fast, Small, Simple Rank/Select on Bitmaps:
//     https://users.dcc.uchile.cl/~gnavarro/ps/sea12.1.pdf

//go:generate mockgen -source ./series_trie_tree.go -destination=./series_trie_tree_mock.go -package indextbl

// trieTreeINTF abstract a trie tree in memory.
// All the descendants of a trieTreeNode have a common prefix of the string associated with that trieTreeNode,
// and the root trieTreeNode is associated with the empty string.
type trieTreeINTF interface {
	// Add adds the tagValue and bitmap to the Trie with the id
	Add(tagValue string, version int64, bitmap *roaring.Bitmap)
	// NodeNum returns the size of trieTreeNodes
	NodeNum() int
	// KeyNum returns the size of keys
	KeyNum() int
	// MarshalBinary marshals the trie tree to LOUDS encoded binary.
	MarshalBinary() *seriesBinData
	// Reset returns the trieTreeNodes to the pool for reuse.
	Reset()
}

// trieTree implements the trieTreeINTF interface.
// Not thread-safe
type trieTree struct {
	pool            sync.Pool
	root            *trieTreeNode
	trieTreeNodeNum int             // count of trieTreeNodes
	keyNum          int             // count of keys
	nodesBuf1       []*trieTreeNode // buffering trieTreeNodes for Breadth first search
	nodesBuf2       []*trieTreeNode // buffering trieTreeNodes for Breadth first search
	bin             *seriesBinData  // binary
}

// newTrieTree returns an empty trie tree.
func newTrieTree() trieTreeINTF {
	tt := &trieTree{
		pool:      sync.Pool{New: func() interface{} { return &trieTreeNode{} }},
		nodesBuf1: make([]*trieTreeNode, 0, 4096),
		nodesBuf2: make([]*trieTreeNode, 0, 4096),
	}
	tt.root = tt.newNode() // empty string
	tt.bin = new(seriesBinData)
	return tt
}

// trieTreeNodes implements sort.Interface
type trieTreeNodes []*trieTreeNode

func (n trieTreeNodes) Len() int           { return len(n) }
func (n trieTreeNodes) Swap(i, j int)      { n[i], n[j] = n[j], n[i] }
func (n trieTreeNodes) Less(i, j int) bool { return n[i].label < n[j].label }

// versionedBitmap is the trieTreeNode value
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

// trieTreeNode is the tree trieTreeNode of the trie
type trieTreeNode struct {
	label       byte
	values      versionedBitmaps // versioned bitmap slice
	isPrefixKey bool             // isPrefixKey indicates if this prefix a valid key
	children    trieTreeNodes    // sorted children trieTreeNodes list
}

// newNode returns a trieTreeNode with default value
func (tt *trieTree) newNode() *trieTreeNode {
	return tt.pool.Get().(*trieTreeNode)
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
			tt.trieTreeNodeNum++
		}
	}
	if !n.isPrefixKey {
		tt.keyNum++
	}
	// set the term trieTreeNode
	n.values.insert(version, bitmap)
	n.isPrefixKey = true
}

// NodeNum returns the size of trieTreeNodes
func (tt *trieTree) NodeNum() int { return tt.trieTreeNodeNum }

// KeyNum returns the size of keys
func (tt *trieTree) KeyNum() int { return tt.keyNum }

// seriesBinData is a LOUDS Encoded trie
type seriesBinData struct {
	LOUDS       RSINTF
	isPrefixKey RSINTF
	labels      []byte
	values      []versionedBitmaps
}

func (lb *seriesBinData) reset() {
	lb.LOUDS = nil
	lb.isPrefixKey = nil
	lb.labels = lb.labels[:0]
	for _, vb := range lb.values {
		vb.reset()
	}
	lb.values = lb.values[:0]
}

// MarshalBinary marshals the trie tree to LOUDS encoded binary.
func (tt *trieTree) MarshalBinary() *seriesBinData {
	tt.bin.LOUDS = NewRankSelect()
	tt.bin.isPrefixKey = NewRankSelect()
	tt.nodesBuf1 = []*trieTreeNode{tt.root}

	tt.bin.LOUDS.PushPseudoRoot()

	for len(tt.nodesBuf1) > 0 {
		for _, trieTreeNode := range tt.nodesBuf1 {
			if trieTreeNode.isPrefixKey {
				tt.bin.values = append(tt.bin.values, trieTreeNode.values)
			}
			tt.bin.isPrefixKey.PushBack(trieTreeNode.isPrefixKey)
			tt.bin.labels = append(tt.bin.labels, trieTreeNode.label)

			for _, trieTreeNode := range trieTreeNode.children {
				tt.nodesBuf2 = append(tt.nodesBuf2, trieTreeNode)
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

// Reset returns the trieTreeNodes to the underlying pool for reuse.
func (tt *trieTree) Reset() {
	tt.nodesBuf1 = []*trieTreeNode{tt.root}
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
	tt.trieTreeNodeNum = 0
	tt.nodesBuf1 = tt.nodesBuf1[:0]
	tt.nodesBuf2 = tt.nodesBuf2[:0]
}
