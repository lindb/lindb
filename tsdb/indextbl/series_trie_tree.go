package indextbl

import (
	"sort"
	"sync"
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

// trieTreeBuilder abstract a trie tree in memory.
// All the descendants of a trieTreeNode have a common prefix of the string associated with that trieTreeNode,
// and the root trieTreeNode is associated with the empty string.
type trieTreeBuilder interface {
	// Add adds the tagValue and item to the tree
	Add(tagValue string, item interface{})
	// NodeNum returns the size of trieTreeNodes
	NodeNum() int
	// KeyNum returns the size of keys
	KeyNum() int
	// MarshalBinary marshals the trie tree to LOUDS encoded binary.
	MarshalBinary() *trieTreeData
	// Reset returns the trieTreeNodes to the pool for reuse.
	Reset()
}

// trieTree implements the trieTreeBuilder interface.
// Not thread-safe
type trieTree struct {
	pool      sync.Pool
	root      *trieTreeNode
	nodeNum   int             // count of trieTreeNodes
	keyNum    int             // count of keys
	nodesBuf1 []*trieTreeNode // buffering trieTreeNodes for Breadth first search
	nodesBuf2 []*trieTreeNode // buffering trieTreeNodes for Breadth first search
	treeData  *trieTreeData   // binary
}

// newTrieTree returns an empty trie tree.
func newTrieTree() trieTreeBuilder {
	tt := &trieTree{
		pool:      sync.Pool{New: func() interface{} { return &trieTreeNode{} }},
		nodesBuf1: make([]*trieTreeNode, 0, 4096),
		nodesBuf2: make([]*trieTreeNode, 0, 4096),
	}
	tt.root = tt.newNode() // empty string
	tt.treeData = new(trieTreeData)
	return tt
}

// trieTreeNodes implements sort.Interface
type trieTreeNodes []*trieTreeNode

func (n trieTreeNodes) Len() int           { return len(n) }
func (n trieTreeNodes) Swap(i, j int)      { n[i], n[j] = n[j], n[i] }
func (n trieTreeNodes) Less(i, j int) bool { return n[i].label < n[j].label }

// trieTreeNode is the tree trieTreeNode of the trie
type trieTreeNode struct {
	label       byte
	value       interface{}   // versioned bitmap slice
	isPrefixKey bool          // isPrefixKey indicates if this prefix a valid key
	children    trieTreeNodes // sorted children trieTreeNodes list
}

// newNode returns a trieTreeNode with default value
func (tt *trieTree) newNode() *trieTreeNode {
	return tt.pool.Get().(*trieTreeNode)
}

// Add adds a new key to the tree.
func (tt *trieTree) Add(tagValue string, item interface{}) {
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
	// set the term trieTreeNode
	n.value = item
	n.isPrefixKey = true
}

// NodeNum returns the size of trieTreeNodes
func (tt *trieTree) NodeNum() int { return tt.nodeNum }

// KeyNum returns the size of keys
func (tt *trieTree) KeyNum() int { return tt.keyNum }

// trieTreeData is a LOUDS Encoded trie tree with values
type trieTreeData struct {
	trieTreeBlock
	values []interface{}
}

func (lb *trieTreeData) reset() {
	lb.LOUDS = nil
	lb.isPrefixKey = nil
	lb.labels = lb.labels[:0]
	for idx := range lb.values {
		lb.values[idx] = nil
	}
	lb.values = lb.values[:0]
}

// MarshalBinary marshals the trie tree to LOUDS encoded binary.
func (tt *trieTree) MarshalBinary() *trieTreeData {
	tt.treeData.LOUDS = NewRankSelect()
	tt.treeData.isPrefixKey = NewRankSelect()
	tt.nodesBuf1 = []*trieTreeNode{tt.root}

	// for easier indexing
	tt.treeData.LOUDS.PushPseudoRoot()
	tt.treeData.isPrefixKey.PushBack(false)
	tt.treeData.labels = append(tt.treeData.labels, byte(0))

	for len(tt.nodesBuf1) > 0 {
		for _, trieTreeNode := range tt.nodesBuf1 {
			if trieTreeNode.isPrefixKey {
				tt.treeData.values = append(tt.treeData.values, trieTreeNode.value)
			}
			tt.treeData.isPrefixKey.PushBack(trieTreeNode.isPrefixKey)
			tt.treeData.labels = append(tt.treeData.labels, trieTreeNode.label)

			for _, trieTreeNode := range trieTreeNode.children {
				tt.nodesBuf2 = append(tt.nodesBuf2, trieTreeNode)
				tt.treeData.LOUDS.PushBack(true)
			}
			tt.treeData.LOUDS.PushBack(false)
		}
		// copy
		tt.nodesBuf1 = append(tt.nodesBuf1[:0], tt.nodesBuf2...)
		tt.nodesBuf2 = tt.nodesBuf2[:0]
	}
	return tt.treeData
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
			n.value = nil
			n.isPrefixKey = false
			n.label = byte(0)
			n.children = n.children[:0]
			tt.pool.Put(n)
		}
		tt.nodesBuf1 = append(tt.nodesBuf1[:0], tt.nodesBuf2...)
		tt.nodesBuf2 = tt.nodesBuf2[:0]
	}
	tt.treeData.reset()
	tt.root = tt.newNode()
	tt.keyNum = 0
	tt.nodeNum = 0
	tt.nodesBuf1 = tt.nodesBuf1[:0]
	tt.nodesBuf2 = tt.nodesBuf2[:0]
}

// trieTreeQuerier represents the ability for querying data-block offsets from the trie-tree
type trieTreeQuerier interface {
	// FindOffsetsByEqual find offsets of prefixKeys which equals to the value in the tree
	FindOffsetsByEqual(value string) (offsets []int)
	// FindOffsetsByIn find offsets of prefixKeys which in the value list in the tree
	FindOffsetsByIn(values []string) (offsets []int)
	// FindOffsetsByLike find offsets of prefixKeys which is like the value in the tree
	FindOffsetsByLike(value string) (offsets []int)
	// FindOffsetsByRegex find offsets of prefixKeys which regex matches the pattern in the tree
	FindOffsetsByRegex(pattern string) (offsets []int)
}

// trieTreeBlock is the structured trie-tree-block of series-index-table
// trieTreeBlock implements trieTreeQuerier
type trieTreeBlock struct {
	labels      []byte
	isPrefixKey RSINTF
	LOUDS       RSINTF
}

func (block *trieTreeBlock) FindOffsetsByEqual(value string) (offsets []int) {
	exhausted, nodeNumber := block.walkTreeByValue(value)
	if exhausted && block.isPrefixKey.Bit(nodeNumber) {
		return []int{int(block.isPrefixKey.Rank1(nodeNumber) - 1)}
	}
	return
}

// walkTreeByValue walks on the tree and exhaust the char in the value
// if all chars are exhausted, it will return true
func (block *trieTreeBlock) walkTreeByValue(value string) (exhausted bool, nodeNumber uint64) {
	// first available node sequence is 1
	nodeNumber = 1
	valueBytes := []byte(value)
	for idx, v := range valueBytes {
		firstChildNumber, ok := block.LOUDS.FirstChild(nodeNumber)
		if !ok {
			continue
		}
		lastChildNumber, ok := block.LOUDS.LastChild(nodeNumber)
		if !ok {
			continue
		}
		for childNumber := firstChildNumber; childNumber <= lastChildNumber; childNumber++ {
			// validate labels length
			if int(childNumber) >= len(block.labels) {
				return false, nodeNumber
			}
			if block.labels[int(childNumber)] == v {
				nodeNumber = childNumber
				break
			}
		}
		if idx == len(valueBytes)-1 {
			exhausted = true
		}
	}
	return
}

func (block *trieTreeBlock) FindOffsetsByIn(values []string) (offsets []int) {
	for _, value := range values {
		eachOffsets := block.FindOffsetsByEqual(value)
		if eachOffsets != nil {
			offsets = append(offsets, eachOffsets...)
		}
	}
	return
}

func (block *trieTreeBlock) FindOffsetsByLike(value string) (offsets []int) {
	exhausted, nodeNumber := block.walkTreeByValue(value)
	if !exhausted {
		return nil
	}
	// exhausted, check if it is the prefix-key
	if block.isPrefixKey.Bit(nodeNumber) {
		offsets = append(offsets, int(block.isPrefixKey.Rank1(nodeNumber)-1))
	}
	var (
		nextNodes = []uint64{nodeNumber}
		tmpNodes  []uint64
	)
	// BFS walk
	for len(nextNodes) > 0 {
		for _, node := range nextNodes {
			firstChildNumber, ok := block.LOUDS.FirstChild(node)
			if !ok {
				continue
			}
			lastChildNumber, ok := block.LOUDS.LastChild(node)
			if !ok {
				continue
			}
			for childNumber := firstChildNumber; childNumber <= lastChildNumber; childNumber++ {
				tmpNodes = append(tmpNodes, childNumber)
				if block.isPrefixKey.Bit(childNumber) {
					offsets = append(offsets, int(block.isPrefixKey.Rank1(childNumber)-1))
				}
			}
		}
		nextNodes = append(nextNodes[:0], tmpNodes...)
		tmpNodes = tmpNodes[:0]
	}
	return offsets
}

// todo: @codingcrush, implementation, traverse the values by forward index?
func (block *trieTreeBlock) FindOffsetsByRegex(pattern string) (offsets []int) { return nil }
