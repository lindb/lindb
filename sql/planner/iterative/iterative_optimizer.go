package iterative

import (
	"fmt"
	"reflect"

	"github.com/lindb/common/pkg/logger"

	"github.com/lindb/lindb/sql/planner/plan"
	"github.com/lindb/lindb/sql/planner/printer"
)

type Context struct {
	IDAllocator *plan.PlanNodeIDAllocator
	memo        *Memo
}

type IterativeOptimizer struct {
	logger logger.Logger

	rules []Rule
}

func NewIterativeOptimizer(rules []Rule) *IterativeOptimizer {
	return &IterativeOptimizer{
		rules:  rules,
		logger: logger.GetLogger("Optimize", "Iterative"),
	}
}

func (opt *IterativeOptimizer) Optimize(node plan.PlanNode, idAllocator *plan.PlanNodeIDAllocator) plan.PlanNode {
	memo := NewMemo(idAllocator, node)
	context := &Context{
		memo:        memo,
		IDAllocator: idAllocator,
	}
	_ = opt.exploreGroup(context, memo.rootGroup)

	return memo.extract(memo.getNode(memo.rootGroup))
}

func (opt *IterativeOptimizer) exploreGroup(context *Context, group int) bool {
	progress := opt.exploreNode(context, group)
	for opt.exploreChildren(context, group) {
		progress = true

		if !opt.exploreNode(context, group) {
			break
		}
	}
	return progress
}

func (opt *IterativeOptimizer) exploreChildren(context *Context, group int) bool {
	progress := false
	expression := context.memo.getNode(group)

	for _, child := range expression.GetSources() {
		fmt.Printf("explore child: %v\n", child)
		if groupRef, ok := child.(*plan.GroupReference); ok {
			if opt.exploreGroup(context, groupRef.GroupID) {
				progress = true
			}
		} else {
			panic(fmt.Sprintf("expected child to be a group reference, found: %T", child))
		}
	}
	return progress
}

func (opt *IterativeOptimizer) exploreNode(context *Context, group int) bool {
	node := context.memo.getNode(group)
	done := false
	progress := false
	for !done {
		// TODO: match rules
		done = true
		for _, rule := range opt.rules {
			result := opt.transform(context, node, rule)
			if result != nil {
				node = context.memo.replace(group, result, reflect.TypeOf(rule).String())
				done = false
				progress = true
			}
		}
	}
	return progress
}

func (opt *IterativeOptimizer) transform(context *Context, node plan.PlanNode, rule Rule) plan.PlanNode {
	// TODO: iterator?
	result := rule.Apply(context, node)
	if result != nil && opt.logger.Enabled(logger.InfoLevel) {
		opt.logger.Info(fmt.Sprintf("rule:%T\nbefore:\n%s\nafter:\n%s",
			rule, printer.TextLogicalPlan(node), printer.TextLogicalPlan(result)))
	}
	return result
}
