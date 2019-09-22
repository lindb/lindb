package art

type iteratorLevel struct {
	node     *artNode
	childIdx int
}

type iterator struct {
	version int // tree version

	tree       *tree
	nextNode   *artNode
	depthLevel int
	depth      []*iteratorLevel
}

type bufferedIterator struct {
	options  int
	nextNode Node
	err      error
	it       *iterator
}

func traverseOptions(opts ...int) int {
	options := 0
	for _, opt := range opts {
		options |= opt
	}
	options &= TraverseAll
	if options == 0 {
		// By default filter only leafs
		options = TraverseLeaf
	}

	return options
}

func traverseFilter(options int, callback Callback) Callback {
	if options == TraverseAll {
		return callback
	}

	return func(node Node) bool {
		if options&TraverseLeaf == TraverseLeaf && node.Kind() == Leaf {
			return callback(node)
		} else if options&TraverseNode == TraverseNode && node.Kind() != Leaf {
			return callback(node)
		}

		return true
	}
}

func (t *tree) ForEach(callback Callback, opts ...int) {
	options := traverseOptions(opts...)
	t.recursiveForEach(t.root, traverseFilter(options, callback))
}

func (t *tree) recursiveForEach(current *artNode, callback Callback) {
	if current == nil {
		return
	}

	if !callback(current) {
		return
	}

	switch current.kind {
	case Node4:
		t.forEachChildren(current.node4().children[:], callback)

	case Node16:
		t.forEachChildren(current.node16().children[:], callback)

	case Node48:
		node := current.node48()
		child := node.children[current.maxChildren()]
		if child != nil {
			t.recursiveForEach(child, callback)
		}

		for _, idx := range node.keys {
			if idx.Absent() {
				continue
			}

			child := node.children[idx.Get()]
			if child != nil {
				t.recursiveForEach(child, callback)
			}
		}

	case Node256:
		t.forEachChildren(current.node256().children[:], callback)
	}
}

func (t *tree) forEachChildren(children []*artNode, callback Callback) {
	nullChild := children[len(children) - 1]
	if  nullChild != nil {
		t.recursiveForEach(nullChild, callback)
	}

	for _, child := range children {
		if child != nil && child != nullChild {
			t.recursiveForEach(child, callback)
		}
	}
}

func (t *tree) ForEachPrefix(key Key, callback Callback) {
	t.forEachPrefix(t.root, key, callback)
}

func (t *tree) forEachPrefix(current *artNode, key Key, callback Callback) {
	if current == nil {
		return
	}

	depth := 0
	for current != nil {
		if current.isLeaf() {
			leaf := current.leaf()
			if leaf.prefixMatch(key) {
				callback(current)
			}
			return
		}

		if depth == len(key) {
			leaf := current.minimum()
			if leaf.prefixMatch(key) {
				t.recursiveForEach(current, callback)
			}
			return
		}

		node := current.node()
		if node.prefixLen > 0 {
			prefixLen := current.matchDeep(key, depth)
			if prefixLen > node.prefixLen {
				prefixLen = node.prefixLen
			}

			if prefixLen == 0 {
				return
			} else if depth+prefixLen == len(key) {
				t.recursiveForEach(current, callback)
				return
			}
			depth += node.prefixLen
		}

		// Find a child to recursive to
		next := current.findChild(key.charAt(depth))
		if *next == nil {
			return
		}
		current = *next
		depth++
	}
}

// Iterator pattern
func (t *tree) Iterator(opts ...int) Iterator {
	options := traverseOptions(opts...)

	it := &iterator{
		version:    t.version,
		tree:       t,
		nextNode:   t.root,
		depthLevel: 0,
		depth:      []*iteratorLevel{{t.root, nullIdx}}}

	if options&TraverseAll == TraverseAll {
		return it
	}

	bti := &bufferedIterator{
		options: options,
		it:      it,
	}
	return bti
}

func (ti *iterator) checkConcurrentModification() error {
	if ti.version == ti.tree.version {
		return nil
	}

	return ErrConcurrentModification
}

func (ti *iterator) HasNext() bool {
	return ti != nil && ti.nextNode != nil
}

func (ti *iterator) Next() (Node, error) {
	if !ti.HasNext() {
		return nil, ErrNoMoreNodes
	}

	err := ti.checkConcurrentModification()
	if err != nil {
		return nil, err
	}

	cur := ti.nextNode
	ti.next()

	return cur, nil
}

const nullIdx = -1

func nextChild(childIdx int, children []*artNode) ( /*nextChildIdx*/ int /*nextNode*/, *artNode) {
	nullChild := children[len(children) - 1]
	if childIdx == nullIdx {
		if nullChild != nil {
			return 0, nullChild
		}

		childIdx = 0
	}

	for i := childIdx; i < len(children); i++ {
		child := children[i]
		if child != nil && child != nullChild {
			return i + 1, child
		}
	}

	return 0, nil
}

func (ti *iterator) next() {
	for {
		var nextNode *artNode
		nextChildIdx := nullIdx

		curNode := ti.depth[ti.depthLevel].node
		curChildIdx := ti.depth[ti.depthLevel].childIdx

		switch curNode.kind {
		case Node4:
			nextChildIdx, nextNode = nextChild(curChildIdx, curNode.node4().children[:])

		case Node16:
			nextChildIdx, nextNode = nextChild(curChildIdx, curNode.node16().children[:])

		case Node48:
			node := curNode.node48()
			nullChild := node.children[curNode.maxChildren()]
			if curChildIdx == nullIdx {
				if nullChild == nil {
					curChildIdx = 0  // try from 0 based child
				} else {
					nextChildIdx = 0 // we have a child with null suffix
					nextNode = nullChild
					break
				}
			}

			for i := curChildIdx; i < len(node.keys); i++ {
				idx := node.keys[byte(i)]
				if idx.Absent() {
					continue
				}
				child := node.children[idx.Get()]
				if child != nil && child != nullChild {
					nextChildIdx = i + 1
					nextNode = child
					break
				}
			}

		case Node256:
			nextChildIdx, nextNode = nextChild(curChildIdx, curNode.node256().children[:])
		}

		if nextNode == nil {
			if ti.depthLevel > 0 {
				// return to previous level
				ti.depthLevel--
			} else {
				ti.nextNode = nil // done!
				return
			}
		} else {
			// star from the next when we come back from the child node
			ti.depth[ti.depthLevel].childIdx = nextChildIdx
			ti.nextNode = nextNode

			// make sure that we have enough space for levels
			if ti.depthLevel+1 >= cap(ti.depth) {
				newDepthLevel := make([]*iteratorLevel, ti.depthLevel+2)
				copy(newDepthLevel, ti.depth)
				ti.depth = newDepthLevel
			}

			ti.depthLevel++
			ti.depth[ti.depthLevel] = &iteratorLevel{nextNode, nullIdx}
			return
		}
	}
}

func (bti *bufferedIterator) HasNext() bool {
	for bti.it.HasNext() {
		bti.nextNode, bti.err = bti.it.Next()
		if bti.err != nil {
			return true
		}
		if bti.options&TraverseLeaf == TraverseLeaf && bti.nextNode.Kind() == Leaf {
			return true
		} else if bti.options&TraverseNode == TraverseNode && bti.nextNode.Kind() != Leaf {
			return true
		}
	}
	bti.nextNode = nil
	bti.err = nil

	return false
}

func (bti *bufferedIterator) Next() (Node, error) {
	return bti.nextNode, bti.err
}
