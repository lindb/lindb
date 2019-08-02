package tree

import (
	"fmt"
	"sort"

	"github.com/lindb/lindb/pkg/logger"
	"github.com/lindb/lindb/pkg/stream"
)

const (
	rootHigh      = 1
	NotFoundValue = 0
)

//BTree represents for the key is bytes, the value is int B+Tree
type BTree struct {
	tree *Tree //Key is bytes, Value is int
}

//NewBTree create a key is bytes B+Tree
func NewBTree() *BTree {
	return &BTree{
		tree: New(BytesCompare),
	}
}

//Put sets the value associated with k.
func (b *BTree) Put(key []byte, value int) {
	b.tree.Set(key, value)
}

//Get returns the value associated with k and true if it exists. Otherwise Get
//returns (zero-value, false).
func (b *BTree) Get(key []byte) (int, bool) {
	value, ok := b.tree.Get(key)
	if ok {
		return value.(int), ok
	}
	return NotFoundValue, ok
}

//Clear removes all K/V pairs from the tree.
func (b *BTree) Clear() {
	b.tree.Clear()
}

//Len returns the number of items in the tree.
func (b *BTree) Len() int {
	return b.tree.Len()
}

//Writer represents encoding the B+tree into the encoder
type Writer struct {
	t       *Tree                  // B+Tree
	highBuf map[int]*stream.Binary //Key is the floor of B+Tree
}

//NewWriter returns new a B+Tree encoder
func NewWriter(tree *BTree) *Writer {
	return &Writer{
		t:       tree.tree,
		highBuf: make(map[int]*stream.Binary),
	}
}

//Encode returns serialize a B+Tree into bytes
func (w *Writer) Encode() ([]byte, error) {
	q := w.t.r
	//no nodes
	if q == nil {
		return nil, fmt.Errorf("empty data")
	}

	writer := stream.BinaryWriter()

	switch qType := q.(type) {
	case *x:
		//branch node
		writer.PutByte(HasChildrenNode)
		//Extract the longest common prefix for each node
		w.extractLcpByBranchNode(rootHigh, qType)
		w.serializeBranchNode(rootHigh, qType)
	case *d:
		//leaf node
		writer.PutByte(NoChildrenNode)
		//Extract the longest common prefix for leaf node
		commonPrefix, _ := w.extractLcpByLeafNode(qType)
		w.serializeLeafNode(rootHigh, qType, commonPrefix)
	}

	highs := make([]int, len(w.highBuf))
	var idx = 0
	for k := range w.highBuf {
		highs[idx] = k
		idx++
	}
	//sort
	sort.Ints(highs)
	//write the height of the tree
	writer.PutUvarint64(uint64(len(highs)))

	nodeWriter := stream.BinaryWriter()
	//write node information of each layer according to the height of the tree
	for _, high := range highs {
		writer.PutUvarint64(uint64(high))
		//starting position of each layer
		writer.PutUvarint64(uint64(nodeWriter.Len()))

		highWriter := w.highBuf[high]
		by, err := highWriter.Bytes()
		if nil != err {
			return nil, err
		}
		nodeWriter.PutBytes(by)
	}

	by, err := nodeWriter.Bytes()
	if nil != err {
		return nil, err
	}
	writer.PutBytes(by)
	return writer.Bytes()
}

//serializeBranchNode represents to write the longest common prefix first,
//and then write the node information including the starting position of the next layer.
func (w *Writer) serializeBranchNode(currentHigh int, parentNode *x) (startPos int) {
	branchWriter := w.getOrCreateByteBuf(currentHigh)
	startPos = branchWriter.Len()
	var nextHigh = currentHigh + 1
	//branch node lcp
	currentCommonPrefix := extractLcpByNodes(parentNode)

	var branchCount = 0
	bodyWriter := stream.BinaryWriter()
	for _, branchNode := range &parentNode.x {
		children := branchNode.ch
		if nil != children {
			branchCount++
			var startPos int
			switch cType := children.(type) {
			case *x:
				//branch node
				startPos = w.serializeBranchNode(nextHigh, cType)
			case *d:
				//leaf node
				startPos = w.serializeLeafNode(nextHigh, cType, branchNode.lcp)
			}

			if nil != branchNode.k {
				bodyWriter.PutByte(HasParent)
				//write suffix
				suffix := branchNode.k.([]byte)[len(currentCommonPrefix):]
				bodyWriter.PutLenBytes(suffix)
			} else {
				bodyWriter.PutByte(NoParent)
				//write suffix
				noParentKey := branchNode.lastKey
				noParentSuffix := noParentKey[len(currentCommonPrefix):]
				bodyWriter.PutLenBytes(noParentSuffix)
			}
			bodyWriter.PutUvarint64(uint64(startPos))
		}
	}

	//write branch node count
	branchWriter.PutUvarint64(uint64(branchCount))
	//write longest common prefix
	branchWriter.PutLenBytes(currentCommonPrefix)
	by, err := bodyWriter.Bytes()
	if nil != err {
		//TODO need handle error
		logger.GetLogger("pkg/tree").Error("serializeBranchNode get bytes error:", logger.Error(err))
	}
	branchWriter.PutBytes(by)
	return startPos
}

// serializeLeafNode represents to write the longest common prefix of the leaf node first,
// then write the value information.
func (w *Writer) serializeLeafNode(currentHigh int, dType *d, leafCommonPrefix []byte) (startPos int) {
	leafWriter := w.getOrCreateByteBuf(currentHigh)

	//startPosition
	startPos = leafWriter.Len()
	dataWriter := stream.BinaryWriter()

	var leafCount int
	for i := 0; i < len(dType.d); i++ {
		if nil != dType.d[i].k {
			leafCount++

			pair := dType.d[i]
			//write suffix
			dataWriter.PutLenBytes(pair.k.([]byte))
			//write value
			dataWriter.PutUvarint64(uint64(pair.v.(int)))
		}
	}
	//write leaf node count
	leafWriter.PutUvarint64(uint64(leafCount))
	//write leaf node lcp
	leafWriter.PutLenBytes(leafCommonPrefix)
	by, err := dataWriter.Bytes()
	if nil != err {
		//TODO need handle error
		logger.GetLogger("pkg/tree").Error("serializeLeafNode get bytes error:", logger.Error(err))
	}
	leafWriter.PutBytes(by)
	return
}

//=============================================Extract the longest common prefix of each node
//extractLcpByNodes represents extract the longest common prefix inside the node
func extractLcpByNodes(parentNode *x) (longestCommonPrefix []byte) {
	var keyArray [][]byte
	for i := 0; i <= parentNode.c; i++ {
		branchNode := parentNode.x[i]
		if nil != branchNode.k {
			keyArray = append(keyArray, branchNode.k.([]byte))
		} else {
			keyArray = append(keyArray, branchNode.lastKey)
		}
	}
	longestCommonPrefix = lcpByte(keyArray)
	return longestCommonPrefix
}

//extractLcpByBranchNode represents Extract the longest common prefix of the branch node.
//High indicates the height of the current tree, starting from 1.
func (w *Writer) extractLcpByBranchNode(high int, parentNode *x) (commonPrefix, lastKey []byte) {
	var branchCommonPrefix []byte //The longest common prefix with the same parent node
	var nextHigh = high + 1       //The next layer of tree height

	var branchKeyArray [][]byte
	for i := 0; i <= parentNode.c; i++ {
		var hasParent = false // Whether the current node has a parent node
		branchNode := &parentNode.x[i]
		if nil != branchNode.k {
			hasParent = true
			branchKeyArray = append(branchKeyArray, branchNode.k.([]byte))
		} else {
			hasParent = false
		}

		children := branchNode.ch
		var commonPrefix, endKey []byte
		switch cType := children.(type) {
		case *x:
			//branch node
			commonPrefix, endKey = w.extractLcpByBranchNode(nextHigh, cType)
		case *d:
			//leaf node
			commonPrefix, endKey = w.extractLcpByLeafNode(cType)
		}
		//no-parent
		if !hasParent {
			branchNode.lastKey = endKey
			lastKey = endKey
			branchKeyArray = append(branchKeyArray, lastKey)
		}
		branchNode.lcp = commonPrefix
	}
	branchCommonPrefix = lcpByte(branchKeyArray)
	return branchCommonPrefix, lastKey
}

//extractLcpByLeafNode represents extract the longest common prefix of the leaf node
func (w *Writer) extractLcpByLeafNode(leafNode *d) (commonPrefix, lastKey []byte) {
	var leafKeyArray [][]byte
	for i := 0; i <= leafNode.c; i++ {
		key := leafNode.d[i].k
		if nil != key {
			leafKeyArray = append(leafKeyArray, key.([]byte))
		}
	}
	//longest common prefix of leaf nodes
	commonPrefix = lcpByte(leafKeyArray)

	//Remove the longest common prefix
	for i := 0; i <= leafNode.c; i++ {
		pair := leafNode.d[i]
		if nil != pair.k {
			suffix := (pair.k.([]byte))[len(commonPrefix):]
			lastKey = pair.k.([]byte)
			leafNode.d[i].k = suffix
		}
	}
	return commonPrefix, lastKey
}

//getOrCreateByteBuf returns the current height of the byte buffer
func (w *Writer) getOrCreateByteBuf(high int) *stream.Binary {
	writer := w.highBuf[high]
	if nil == writer {
		writer = stream.BinaryWriter()
		w.highBuf[high] = writer
	}
	return writer
}
