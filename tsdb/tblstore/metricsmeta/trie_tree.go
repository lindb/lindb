package metricsmeta

import (
	"bytes"
	"regexp"
	"sort"
	"strings"
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

//go:generate mockgen -source ./trie_tree.go -destination=./trie_tree_mock.go -package metricsmeta

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

var trieTreeNodePool = sync.Pool{New: func() interface{} {
	return &trieTreeNode{}
}}

// trieTree implements the trieTreeBuilder interface.
// Not thread-safe
type trieTree struct {
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
	return trieTreeNodePool.Get().(*trieTreeNode)
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
			trieTreeNodePool.Put(n)
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
	// PrefixSearch returns keys by prefix-search
	PrefixSearch(value string, limit int) (founds []string)
	// Iterator returns the trie-tree iterator
	Iterator(prefixValue string) *TrieTreeIterator
	// GetValueByOffset picks the value at specified offset
	GetValueByOffset(offset int) (value string, ok bool)
}

// trieTreeBlock is the structured trie-tree-block of index-table
// trieTreeBlock implements trieTreeQuerier
type trieTreeBlock struct {
	labels      []byte
	isPrefixKey RSINTF
	LOUDS       RSINTF
}

func (block *trieTreeBlock) FindOffsetsByEqual(value string) (offsets []int) {
	exhausted, nodeNumber := block.walkTreeByValue([]byte(value))
	if exhausted && block.isPrefixKey.Bit(nodeNumber) {
		return []int{int(block.isPrefixKey.Rank1(nodeNumber) - 1)}
	}
	return
}

// walkTreeByValue walks on the tree and exhaust the char in the value
// if all chars are exhausted, it will return true
func (block *trieTreeBlock) walkTreeByValue(value []byte) (exhausted bool, nodeNumber uint64) {
	// first available node sequence is 1
	nodeNumber = 1
	if len(value) == 0 {
		return true, nodeNumber
	}
	var (
		lastWalkedIndex = 0
		maxNodeNumber   = block.LOUDS.MaxNodeNumber()
	)
	for idx, v := range value {
		firstChildNumber, ok := block.LOUDS.FirstChild(nodeNumber)
		if !ok {
			break
		}
		lastChildNumber, ok := block.LOUDS.LastChild(nodeNumber)
		if !ok {
			break
		}
		lastWalkedIndex = idx
		found := false
		for childNumber := firstChildNumber; childNumber <= lastChildNumber; childNumber++ {
			// validate labels length
			if int(childNumber) >= len(block.labels) {
				return false, maxNodeNumber + 1
			}
			if block.labels[int(childNumber)] == v {
				found = true
				nodeNumber = childNumber
				break
			}
		}
		if !found {
			return false, maxNodeNumber + 1
		}
	}
	// exhaust all
	if lastWalkedIndex == len(value)-1 {
		return true, nodeNumber
	}
	return false, maxNodeNumber + 1
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
	var prefix, middle, suffix []byte
	hashPrefix := strings.HasPrefix(value, "*")
	hasSuffix := strings.HasSuffix(value, "*")
	switch {
	case value == "":
		return nil
	case value == "*":
	// only endswith *
	case !hashPrefix && hasSuffix:
		prefix = []byte(value[:len(value)-1])
	// only startswith *
	case hashPrefix && !hasSuffix:
		suffix = []byte(value[1:])
	// startswith and endswith *
	case hashPrefix && hasSuffix:
		middle = []byte(value[1 : len(value)-1])
	default:
		return block.FindOffsetsByEqual(value)
	}

	itr := block.Iterator(string(prefix))
	for itr.HasNext() {
		key, offset := itr.Next()
		switch {
		case len(middle) > 0:
			if !bytes.Contains(key, middle) {
				continue
			}
		case len(suffix) > 0:
			if !bytes.HasSuffix(key, suffix) {
				continue
			}
		}
		offsets = append(offsets, offset)
	}
	return offsets
}

// _triePrefix represents a prefix on the trie tree
type _triePrefix struct {
	nodeNumber int
	payload    []byte
}

func (block *trieTreeBlock) FindOffsetsByRegex(pattern string) (offsets []int) {
	rp, err := regexp.Compile(pattern)
	if err != nil {
		return nil
	}
	literalPrefix, _ := rp.LiteralPrefix()

	itr := block.Iterator(literalPrefix)
	for itr.HasNext() {
		value, offset := itr.Next()
		if rp.Match(value) {
			offsets = append(offsets, offset)
		}
	}
	return offsets
}

func (block *trieTreeBlock) PrefixSearch(prefixValue string, limit int) (founds []string) {
	itr := block.Iterator(prefixValue)
	for itr.HasNext() {
		if len(founds) >= limit {
			break
		}
		value, _ := itr.Next()
		founds = append(founds, string(value))
	}
	return founds
}

// Iterator returns the trie-tree iterator
func (block *trieTreeBlock) Iterator(prefixValue string) *TrieTreeIterator {
	return &TrieTreeIterator{
		block:       block,
		prefixValue: []byte(prefixValue)}
}

func reverseBytes(s []byte) {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
}

func (block *trieTreeBlock) GetValueByOffset(offset int) (values string, ok bool) {
	var buf []byte // combine the values
	maxNodeNumber := block.LOUDS.MaxNodeNumber()
	// validate labels count and node-numbers
	if maxNodeNumber+1 != uint64(len(block.labels)) {
		return
	}
	buf = buf[:0]
	thisNodeNumber := block.isPrefixKey.Select1(uint64(offset) + 1)
	// combines the prefixKey path
	for 0 < thisNodeNumber && thisNodeNumber <= maxNodeNumber {
		buf = append(buf, block.labels[int(thisNodeNumber)])
		thisNodeNumber = block.LOUDS.Parent(thisNodeNumber)
		if thisNodeNumber >= maxNodeNumber {
			break
		}
	}
	if len(buf) == 0 {
		return
	}
	reverseBytes(buf)
	return string(buf[1:]), true
}

// TrieTreeIterator implements TrieTreeIterator
type TrieTreeIterator struct {
	block       *trieTreeBlock
	prefixValue []byte
	hasError    bool
	initialized bool
	prefixes    []_triePrefix
	value       []byte
	offset      int
}

func (itr *TrieTreeIterator) HasNext() bool {
	if itr.hasError {
		return false
	}
	if !itr.initialized {
		itr.initialized = true
		var startNodeNumber = 1
		if len(itr.prefixValue) != 0 {
			exhausted, nodeNumber := itr.block.walkTreeByValue(itr.prefixValue)
			if !exhausted {
				itr.hasError = true
				return false
			}
			// exhausted
			startNodeNumber = int(nodeNumber)
		}
		itr.prefixes = []_triePrefix{{nodeNumber: startNodeNumber, payload: itr.prefixValue}}
	}
	hasNext := false
	for len(itr.prefixes) > 0 {
		thisPrefix := itr.prefixes[len(itr.prefixes)-1]   // get the tail prefix
		itr.prefixes = itr.prefixes[:len(itr.prefixes)-1] // pop it
		// collect a new key and offset
		if itr.block.isPrefixKey.Bit(uint64(thisPrefix.nodeNumber)) {
			// gotcha
			hasNext = true
			itr.offset = int(itr.block.isPrefixKey.Rank1(uint64(thisPrefix.nodeNumber)) - 1)
			itr.value = thisPrefix.payload
		}
		var (
			firstChildNumber, lastChildNumber uint64
			ok                                bool
		)
		firstChildNumber, ok = itr.block.LOUDS.FirstChild(uint64(thisPrefix.nodeNumber))
		if !ok {
			goto CheckHasNext
		}
		lastChildNumber, ok = itr.block.LOUDS.LastChild(uint64(thisPrefix.nodeNumber))
		if !ok {
			goto CheckHasNext
		}
		for childNumber := firstChildNumber; childNumber <= lastChildNumber; childNumber++ {
			// validate labels length, error occurred
			if int(childNumber) >= len(itr.block.labels) {
				itr.hasError = true
				return false
			}
			newPrefix := _triePrefix{nodeNumber: int(childNumber), payload: make([]byte, 16)[:0]}
			newPrefix.payload = append(newPrefix.payload, thisPrefix.payload...)
			newPrefix.payload = append(newPrefix.payload, itr.block.labels[int(childNumber)])
			itr.prefixes = append(itr.prefixes, newPrefix)
		}
	CheckHasNext:
		{
			if !hasNext {
				continue
			}
			return hasNext
		}
	}
	return hasNext
}

func (itr *TrieTreeIterator) Next() (value []byte, offset int) {
	return itr.value, itr.offset
}
