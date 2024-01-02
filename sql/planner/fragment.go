package planner

import (
	"github.com/lindb/lindb/models"
	planpkg "github.com/lindb/lindb/sql/planner/plan"
)

type PlanFragmenter struct{}

func NewPlanFragmenter() *PlanFragmenter {
	return &PlanFragmenter{}
}

func (f *PlanFragmenter) CreateSubPlans(plan *planpkg.Plan) *planpkg.SubPlan {
	fragmenter := NewFragmenter()
	rootProps := &FragmentProps{}
	root := plan.Root.Accept(rootProps, fragmenter).(planpkg.PlanNode)
	subPlan := fragmenter.buildRootFragment(root, rootProps)
	return subPlan
}

type FragmentProps struct {
	partitions map[models.InternalNode][]int
	children   []*planpkg.SubPlan
}

type Fragmenter struct {
	idAllocator planpkg.FragmentID
}

func NewFragmenter() *Fragmenter {
	return &Fragmenter{
		idAllocator: planpkg.RootFragmentID + 1,
	}
}

func (f *Fragmenter) Visit(context any, n planpkg.PlanNode) (r any) {
	switch node := n.(type) {
	case *planpkg.ExchangeNode:
		return f.visitExchange(context, node)
	case *planpkg.TableScanNode:
		if props, ok := context.(*FragmentProps); ok {
			props.partitions = node.Partitions
		}
		return f.defaultRewrite(context, n)
	default:
		return f.defaultRewrite(context, n)
	}
}

func (f *Fragmenter) visitExchange(context any, n *planpkg.ExchangeNode) (r any) {
	if n.Scope != planpkg.Remote {
		return f.defaultRewrite(context, n)
	}
	var children []*planpkg.SubPlan
	var childrenIDs []planpkg.FragmentID

	exchangeNodeID := n.GetNodeID()
	sources := n.GetSources()
	for i := range sources {
		childProps := &FragmentProps{}
		subPlan := f.buildSubPlan(sources[i], childProps)
		// set fragment's parent node
		subPlan.Fragment.RemoteParentNodeID = &exchangeNodeID

		children = append(children, subPlan)
		childrenIDs = append(childrenIDs, subPlan.Fragment.ID)
	}
	props := context.(*FragmentProps)
	props.children = append(props.children, children...)
	return &planpkg.RemoteSourceNode{
		BaseNode: planpkg.BaseNode{
			ID: n.GetNodeID(),
		},
		SourceFragmentIDs: childrenIDs,
	}
}

func (f *Fragmenter) defaultRewrite(context any, node planpkg.PlanNode) (r any) {
	var children []planpkg.PlanNode
	sources := node.GetSources()
	for i := range sources {
		children = append(children, sources[i].Accept(context, f).(planpkg.PlanNode))
	}
	return planpkg.ReplaceChildren(node, children)
}

func (f *Fragmenter) buildRootFragment(root planpkg.PlanNode, props *FragmentProps) *planpkg.SubPlan {
	return f.buildFragment(planpkg.RootFragmentID, root, props)
}

func (f *Fragmenter) buildFragment(id planpkg.FragmentID, root planpkg.PlanNode, props *FragmentProps) *planpkg.SubPlan {
	fragment := planpkg.NewPlanFragment(id, root)
	fragment.Partitions = props.partitions
	return &planpkg.SubPlan{
		Fragment: fragment,
		Children: props.children,
	}
}

func (f *Fragmenter) buildSubPlan(node planpkg.PlanNode, props *FragmentProps) *planpkg.SubPlan {
	child := node.Accept(props, f).(planpkg.PlanNode)
	id := f.nextID()
	return f.buildFragment(id, child, props)
}

func (f *Fragmenter) nextID() planpkg.FragmentID {
	id := f.idAllocator
	f.idAllocator++
	return id
}
