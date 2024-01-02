package iterative

import (
	"fmt"

	"github.com/samber/lo"

	"github.com/lindb/lindb/constants"
	"github.com/lindb/lindb/sql/planner/plan"
)

const RootGroupRef = 0

type Group struct {
	membership plan.PlanNode
}

func withMember(node plan.PlanNode) *Group {
	return &Group{
		membership: node,
	}
}

type GroupReference struct {
	outputs []*plan.Symbol
	groupID int

	plan.BaseNode
}

func (n *GroupReference) Accept(context any, visitor plan.Visitor) any {
	return visitor.Visit(context, n)
}

func (n *GroupReference) GetOutputSymbols() []*plan.Symbol {
	return n.outputs
}

func (n *GroupReference) GetSources() []plan.PlanNode {
	panic(constants.ErrNotSupportOperation)
}

func (n *GroupReference) ReplaceChildren(newChildren []plan.PlanNode) plan.PlanNode {
	panic(constants.ErrNotSupportOperation)
}

type Memo struct {
	idAllocator *plan.PlanNodeIDAllocator
	groups      map[int]*Group

	rootGroup   int
	nextGroupID int
}

func NewMemo(idAllocator *plan.PlanNodeIDAllocator, node plan.PlanNode) *Memo {
	memo := &Memo{
		idAllocator: idAllocator,
		nextGroupID: RootGroupRef + 1,
		groups:      make(map[int]*Group),
	}
	memo.rootGroup = memo.insertRecursive(node)
	// TODO: add ref????
	return memo
}

func (m *Memo) resolve(groupRef *GroupReference) plan.PlanNode {
	return m.getNode(groupRef.groupID)
}

func (m *Memo) replace(groupID int, node plan.PlanNode, _ string) plan.PlanNode {
	group := m.groups[groupID]
	// TODO:check old output???
	if groupRef, ok := node.(*GroupReference); ok {
		node = m.getNode(groupRef.groupID)
	} else {
		node = m.insertChildrenAndRewrite(node)
	}

	group.membership = node

	return node
}

func (m *Memo) extract(node plan.PlanNode) plan.PlanNode {
	return resolveGroupReferences(node, m.resolve)
}

func (m *Memo) getNode(group int) plan.PlanNode {
	return m.groups[group].membership
}

func (m *Memo) insertRecursive(node plan.PlanNode) int {
	if groupRef, ok := node.(*GroupReference); ok {
		return groupRef.groupID
	}
	group := m.genNextGroupID()
	rewritten := m.insertChildrenAndRewrite(node)
	m.groups[group] = withMember(rewritten)
	// TODO: inc ref

	return group
}

func (m *Memo) insertChildrenAndRewrite(node plan.PlanNode) plan.PlanNode {
	newChildren := lo.Map(node.GetSources(), func(item plan.PlanNode, index int) plan.PlanNode {
		fmt.Printf("kkk....%T=%v\n", item, item.GetOutputSymbols())
		return &GroupReference{
			BaseNode: plan.BaseNode{
				ID: m.idAllocator.Next(),
			},
			groupID: m.insertRecursive(item),
			outputs: item.GetOutputSymbols(),
		}
	})
	return node.ReplaceChildren(newChildren)
}

func (m *Memo) genNextGroupID() int {
	m.nextGroupID++
	return m.nextGroupID
}
