package meta

import (
	"sync/atomic"

	"github.com/lindb/lindb/models"
)

var currentNode atomic.Value

// SetCurrentNode sets current node.
func SetCurrentNode(node models.NodeID) {
	currentNode.Store(node)
}

// CurrentNode returns current node.
func CurrentNode() models.NodeID {
	return currentNode.Load().(models.NodeID)
}
