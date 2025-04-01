package join

import (
	"fmt"

	"github.com/samber/lo"

	"github.com/lindb/lindb/spi/types"
	"github.com/lindb/lindb/sql/execution/pipeline/operator"
	"github.com/lindb/lindb/sql/planner/plan"
)

type rows struct {
	left, right *types.Row
}

type HashJoinOperator struct {
	node        *plan.JoinNode
	left, right operator.Operator

	memTable              map[string]*rows
	leftKeys, rightKeys   []int
	leftScope, rightScope []*plan.Symbol
}

func NewHashJoinOperator(node *plan.JoinNode, left, right operator.Operator) operator.Operator {
	return &HashJoinOperator{
		node:       node,
		left:       left,
		right:      right,
		leftScope:  node.Left.GetOutputSymbols(),
		rightScope: node.Right.GetOutputSymbols(),
	}
}

// Children implements operator.Operator.
func (h *HashJoinOperator) Children() []operator.Operator {
	return []operator.Operator{h.left, h.right}
}

// GetLayout implements operator.Operator.
func (h *HashJoinOperator) GetLayout() []*plan.Symbol {
	return h.node.GetOutputSymbols()
}

func (h *HashJoinOperator) prepare() {
	h.memTable = make(map[string]*rows)

	joinCriteria := h.node.Criteria
	h.leftKeys = make([]int, len(joinCriteria))
	h.rightKeys = make([]int, len(joinCriteria))

	findIndex := func(keys []int, index int, key *plan.Symbol, scope []*plan.Symbol) {
		_, keySymbol, ok := lo.FindIndexOf(scope, func(item *plan.Symbol) bool {
			return item.Name == key.Name
		})
		if !ok {
			panic("not find left key")
		}
		keys[index] = keySymbol
	}

	for idx, expr := range joinCriteria {
		findIndex(h.leftKeys, idx, expr.Left, h.leftScope)
		findIndex(h.rightKeys, idx, expr.Right, h.rightScope)
	}
}

// Run implements operator.Operator.
func (h *HashJoinOperator) Run(output chan<- *types.Page) {
	h.prepare()

	leftInbound := make(chan *types.Page)
	go func() {
		defer func() {
			close(leftInbound)
		}()
		h.left.Run(leftInbound)
	}()

	rightInbound := make(chan *types.Page)
	go func() {
		defer func() {
			close(rightInbound)
		}()
		h.right.Run(rightInbound)
	}()

	for leftInbound != nil || rightInbound != nil {
		// consume page from left and right inbound
		select {
		case left, ok := <-leftInbound:
			if !ok {
				leftInbound = nil
			} else {
				h.process(left, h.leftKeys, true)
				fmt.Printf("left=>>>>>>>%v\n", left)
			}
		case right, ok := <-rightInbound:
			if !ok {
				rightInbound = nil
			} else {
				h.process(right, h.rightKeys, false)
				fmt.Printf("right=>>>>>>>%v\n", right)
			}
		}
	}
	newPage := types.NewPage()
	outputs := h.node.GetOutputSymbols()
	outputColumns := make([]*types.Column, len(outputs))
	for i, output := range outputs {
		outputColumns[i] = types.NewColumn()
		newPage.AppendColumn(types.NewColumnInfo(output.Name, output.DataType), outputColumns[i])
	}

	for key, memRows := range h.memTable {
		fmt.Printf("join key=%v,memRows=%v........\n", key, memRows)
		if memRows.left != nil && memRows.right != nil {
			for i := range len(h.leftScope) {
				outputColumns[i].Append(memRows.left.Get(i))
			}

			for i := range len(h.rightScope) {
				outputColumns[i+len(h.leftScope)].Append(memRows.right.Get(i))
			}
		}
	}
	output <- newPage
}

func (h *HashJoinOperator) process(page *types.Page, keys []int, isLeft bool) {
	it := page.Iterator()
	for row := it.Begin(); row != it.End(); row = it.Next() {
		column := row.GetString(keys[0])
		memRows, ok := h.memTable[string(*column)]
		if !ok {
			memRows = &rows{}
			h.memTable[string(*column)] = memRows
		}

		if isLeft {
			memRows.left = &row
		} else {
			memRows.right = &row
		}
	}
}
