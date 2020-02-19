package models

// PhysicalPlan represents the distribution query's physical plan
type PhysicalPlan struct {
	Database      string         `json:"database"`            // database name
	Namespace     string         `json:"namespace,omitempty"` // namespace
	Root          Root           `json:"root"`                // root node
	Intermediates []Intermediate `json:"intermediates"`       // intermediate node if need
	Leafs         []Leaf         `json:"leafs"`               // leaf nodes(storage nodes of query database)
}

// NewPhysicalPlan creates the physical plan with root node
func NewPhysicalPlan(root Root) *PhysicalPlan {
	return &PhysicalPlan{Root: root}
}

// AddIntermediate adds an intermediate node into the intermediate node list
func (t *PhysicalPlan) AddIntermediate(intermediate Intermediate) {
	t.Intermediates = append(t.Intermediates, intermediate)
}

// AddLeaf adds a leaf node into the leaf node list
func (t *PhysicalPlan) AddLeaf(leaf Leaf) {
	t.Leafs = append(t.Leafs, leaf)
}

// Root represents the root node info
type Root struct {
	Indicator string
	NumOfTask int32
}

type BaseNode struct {
	Parent    string // parent node's indicator
	Indicator string // current node's indicator
}

// Intermediate represents the intermediate node info
type Intermediate struct {
	BaseNode

	NumOfTask int32
}

// Leaf represents the leaf node info
type Leaf struct {
	BaseNode

	Receivers []Node
	ShardIDs  []int32
}
