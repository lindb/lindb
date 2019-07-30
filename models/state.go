package models

// StorageState represents storage cluster node state.
// NOTICE: it is not safe for concurrent use.
type StorageState struct {
	Name        string           `json:"name"`
	ActiveNodes map[string]*Node `json:"activeNodes"`
}

// NewStorageState creates storage cluster state
func NewStorageState() *StorageState {
	return &StorageState{
		ActiveNodes: make(map[string]*Node),
	}
}

// AddActiveNode adds a node into active node list
func (s *StorageState) AddActiveNode(node *Node) {
	key := node.Indicator()
	_, ok := s.ActiveNodes[key]
	if !ok {
		s.ActiveNodes[key] = node
	}
}

// RemoveActiveNode removes a node from active node list
func (s *StorageState) RemoveActiveNode(node string) {
	delete(s.ActiveNodes, node)
}

// GetActiveNodes returns all active nodes
func (s *StorageState) GetActiveNodes() []*Node {
	var nodes []*Node
	for _, node := range s.ActiveNodes {
		nodes = append(nodes, node)
	}
	return nodes
}
