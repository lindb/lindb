package printer

import (
	"github.com/lindb/lindb/sql/planner/plan"
)

type PlanRepresentation struct {
	root     plan.PlanNode
	nodeInfo map[plan.PlanNodeID]*NodeRepresentation
}

func NewPlanRepresentation(root plan.PlanNode) *PlanRepresentation {
	return &PlanRepresentation{
		root:     root,
		nodeInfo: make(map[plan.PlanNodeID]*NodeRepresentation),
	}
}

func (pr *PlanRepresentation) getRoot() (nr *NodeRepresentation) {
	nr = pr.nodeInfo[pr.root.GetNodeID()]
	return nr
}

func (pr *PlanRepresentation) getNode(nodeID plan.PlanNodeID) (node *NodeRepresentation) {
	node = pr.nodeInfo[nodeID]
	return
}

func (pr *PlanRepresentation) addNode(node *NodeRepresentation) {
	pr.nodeInfo[node.getID()] = node
}

type NodeRepresentation struct {
	descriptor map[string]string
	name       string
	children   []plan.PlanNodeID
	details    []string
	outputs    []*plan.Symbol
	id         plan.PlanNodeID
}

func (nr *NodeRepresentation) appendDetails(detail string) {
	nr.details = append(nr.details, detail)
}

func (nr *NodeRepresentation) getID() plan.PlanNodeID {
	return nr.id
}

func (nr *NodeRepresentation) getName() string {
	return nr.name
}
