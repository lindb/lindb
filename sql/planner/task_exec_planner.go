package planner

import (
	"fmt"

	"github.com/samber/lo"

	"github.com/lindb/lindb/execution/buffer"
	"github.com/lindb/lindb/execution/pipeline"
	"github.com/lindb/lindb/execution/pipeline/operator"
	"github.com/lindb/lindb/execution/pipeline/operator/exchange"
	"github.com/lindb/lindb/execution/pipeline/operator/output"
	"github.com/lindb/lindb/execution/pipeline/operator/scan"
	"github.com/lindb/lindb/spi"
	"github.com/lindb/lindb/sql/context"
	planpkg "github.com/lindb/lindb/sql/planner/plan"
	"github.com/lindb/lindb/sql/tree"
)

type TaskExecutionPlanner struct{}

func NewTaskExecutionPlanner() *TaskExecutionPlanner {
	return &TaskExecutionPlanner{}
}

func (p *TaskExecutionPlanner) Plan(taskCtx *context.TaskContext, plan planpkg.PlanNode, outputBuffer buffer.OutputBuffer) *TaskExecutionPlan {
	visitor := &TaskExecutionPlanVisitor{
		taskExecCtx: taskCtx,
		planner:     p,
	}
	taskExecPlanCtx := NewTaskExecutionPlanContext(nil)
	var physicalOperator *PhysicalOperation
	if op, ok := plan.Accept(taskExecPlanCtx, visitor).(*PhysicalOperation); ok {
		physicalOperator = op
	}
	if physicalOperator == nil {
		panic("cannot get physicalOperator")
	}

	outputOperatorFct := output.NewRSOutputOperatorFactory(outputBuffer)

	// add output operator
	taskExecPlanCtx.AddDriverFactory(NewPhysicalOperation(outputOperatorFct, physicalOperator))

	var pipelines []*pipeline.Pipeline
	for i := range taskExecPlanCtx.driverFactories {
		var splitSource spi.SplitSource
		if taskExecPlanCtx.IsLocalStore() {
			// local data source
			splitSource = taskExecPlanCtx.splitSources[i]
		}

		pipelines = append(pipelines, pipeline.NewPipeline(taskCtx, splitSource, taskExecPlanCtx.driverFactories[0]))
	}

	return NewTaskExecutionPlan(pipelines)
}

type TaskExecutionPlanVisitor struct {
	taskExecCtx *context.TaskContext
	planner     *TaskExecutionPlanner
}

func (v *TaskExecutionPlanVisitor) Visit(context any, n planpkg.PlanNode) (r any) {
	switch node := n.(type) {
	case *planpkg.OutputNode:
		return node.Source.Accept(context, v)
	case *planpkg.ProjectionNode:
		return v.visitProjection(context, node)
	case *planpkg.TableScanNode:
		return v.VisitTableScan(context, node)
	case *planpkg.RemoteSourceNode:
		return v.visitRemoteSource(context, node)
	case *planpkg.AggregationNode:
		return v.visitAggregation(context, node)
	case *planpkg.ExchangeNode:
		return v.visitExchange(context, node)
	default:
		panic(fmt.Sprintf("imple task planner %v", n))
	}
}

func (v *TaskExecutionPlanVisitor) visitExchange(context any, node *planpkg.ExchangeNode) (r any) {
	if node.Scope != planpkg.Local {
		panic("only local exchanges are supported in the local planner")
	}
	source := node.Sources[0].Accept(context, v).(*PhysicalOperation)
	operatorFct := exchange.NewLocalExchangeOperatorFactory()
	return NewPhysicalOperation(operatorFct, source)
}

func (v *TaskExecutionPlanVisitor) visitAggregation(context any, node *planpkg.AggregationNode) (r any) {
	source := node.Source.Accept(context, v).(*PhysicalOperation)
	return v.planGroupByAggregation(node, source)
}

func (v *TaskExecutionPlanVisitor) planGroupByAggregation(node *planpkg.AggregationNode, source *PhysicalOperation) *PhysicalOperation {
	operatorFct := v.createHashAggregationOperatorFactory()
	return NewPhysicalOperation(operatorFct, source)
}

func (v *TaskExecutionPlanVisitor) createHashAggregationOperatorFactory() operator.OperatorFactory {
	return operator.NewHashAggregationOperatorFactory()
}

func (v *TaskExecutionPlanVisitor) visitProjection(context any, node *planpkg.ProjectionNode) (r any) {
	var source planpkg.PlanNode
	var filter tree.Expression
	if filterNode, ok := node.Source.(*planpkg.FilterNode); ok {
		source = filterNode.Source
		filter = filterNode.Predicate
	} else {
		source = node.Source
	}
	return v.visitScanFilterAndProjection(context, source, filter)
}

func (v *TaskExecutionPlanVisitor) VisitTableScan(context any, node *planpkg.TableScanNode) (r any) {
	operatorFct := v.visitTableScan(node, nil, context)
	return NewPhysicalOperation(operatorFct, nil)
}

func (v *TaskExecutionPlanVisitor) visitRemoteSource(context any, node *planpkg.RemoteSourceNode) (r any) {
	operatorFct := exchange.NewExchangeOperatorFactory(node.GetNodeID(), len(node.SourceFragmentIDs))
	return NewPhysicalOperation(operatorFct, nil)
}

func (v *TaskExecutionPlanVisitor) visitScanFilterAndProjection(context any, sourceNode planpkg.PlanNode, filter tree.Expression) any {
	var (
		source    *PhysicalOperation
		table     spi.TableHandle
		tableScan *planpkg.TableScanNode
		ok        bool
	)
	if tableScan, ok = sourceNode.(*planpkg.TableScanNode); ok {
		table = tableScan.Table
	} else {
		// plan source node
		source = sourceNode.Accept(context, v).(*PhysicalOperation)
	}

	if table != nil {
		operatorFct := v.visitTableScan(tableScan, filter, context)
		return NewPhysicalOperation(operatorFct, source)
	}

	return source
}

func (v *TaskExecutionPlanVisitor) visitTableScan(node *planpkg.TableScanNode, filter tree.Expression, context any) operator.OperatorFactory {
	outputs := node.GetOutputSymbols()
	columns := lo.Map(outputs, func(item *planpkg.Symbol, index int) spi.ColumnMetadata {
		return spi.ColumnMetadata{
			Name:     item.Name,
			DataType: item.DataType,
		}
	})
	splitSources := spi.GetSplitSourceProvider(node.Table).CreateSplitSources(node.Table, v.taskExecCtx.Partitions, columns, filter)
	// TODO: check source split
	planContext := context.(*TaskExecutionPlanContext)
	planContext.SetSplitSources(splitSources)
	planContext.SetLocalStore(true)

	return scan.NewTableScanOperatorFactory(node.GetNodeID(), node.Table, filter)
}
