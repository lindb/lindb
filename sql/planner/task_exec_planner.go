package planner

import (
	"fmt"

	"github.com/samber/lo"

	"github.com/lindb/lindb/spi"
	"github.com/lindb/lindb/sql/context"
	"github.com/lindb/lindb/sql/execution/buffer"
	"github.com/lindb/lindb/sql/execution/pipeline"
	"github.com/lindb/lindb/sql/execution/pipeline/operator"
	"github.com/lindb/lindb/sql/execution/pipeline/operator/exchange"
	"github.com/lindb/lindb/sql/execution/pipeline/operator/output"
	"github.com/lindb/lindb/sql/execution/pipeline/operator/scan"
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
			// TODO: need check split source if nil
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

// Visit visits all plan node and plans task execution physical operator.
func (v *TaskExecutionPlanVisitor) Visit(context any, n planpkg.PlanNode) (r any) {
	switch node := n.(type) {
	case *planpkg.OutputNode:
		return node.Source.Accept(context, v)
	case *planpkg.AggregationNode:
		return v.visitAggregation(context, node)
	case *planpkg.RemoteSourceNode:
		return v.visitRemoteSource(context, node)
	case *planpkg.ExchangeNode:
		return v.visitExchange(context, node)
	case *planpkg.ProjectionNode:
		return v.visitProjection(context, node)
	case *planpkg.FilterNode:
		return v.visitFilter(context, node)
	case *planpkg.TableScanNode:
		return v.VisitTableScan(context, node)
	default:
		panic(fmt.Sprintf("umimplements task planner %T", n))
	}
}

// visitFilter plans filter physical operator.
func (v *TaskExecutionPlanVisitor) visitFilter(context any, node *planpkg.FilterNode) (r any) {
	if tableScan, ok := node.Source.(*planpkg.TableScanNode); ok {
		operatorFct := v.visitTableScan(context, tableScan, node.Predicate)
		return NewPhysicalOperation(operatorFct, nil)
	}
	panic("need impl visitFilter")
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
	return v.visitScanFilterAndProjection(context, node, source, filter)
}

func (v *TaskExecutionPlanVisitor) VisitTableScan(context any, node *planpkg.TableScanNode) (r any) {
	operatorFct := v.visitTableScan(context, node, nil)
	return NewPhysicalOperation(operatorFct, nil)
}

func (v *TaskExecutionPlanVisitor) visitRemoteSource(context any, node *planpkg.RemoteSourceNode) (r any) {
	operatorFct := exchange.NewExchangeOperatorFactory(node.GetNodeID(), len(node.SourceFragmentIDs))
	return NewPhysicalOperation(operatorFct, nil)
}

func (v *TaskExecutionPlanVisitor) visitScanFilterAndProjection(context any, project *planpkg.ProjectionNode, sourceNode planpkg.PlanNode, filter tree.Expression) any {
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
		operatorFct := v.visitTableScan(context, tableScan, filter)
		return NewPhysicalOperation(operatorFct, source)
	}
	projectOpFct := operator.NewProjectionOperatorFactory(project, sourceNode.GetOutputSymbols())
	return NewPhysicalOperation(projectOpFct, source)
}

func (v *TaskExecutionPlanVisitor) visitTableScan(context any, node *planpkg.TableScanNode, filter tree.Expression) operator.OperatorFactory {
	outputs := node.GetOutputSymbols()
	outputColumns := lo.Map(outputs, func(item *planpkg.Symbol, index int) spi.ColumnMetadata {
		return spi.ColumnMetadata{
			Name:     item.Name,
			DataType: item.DataType,
		}
	})
	splitSources := spi.GetSplitSourceProvider(node.Table).CreateSplitSources(node.Table, v.taskExecCtx.Partitions, outputColumns, filter)
	// TODO: check source split
	planContext := context.(*TaskExecutionPlanContext)
	planContext.SetSplitSources(splitSources)
	planContext.SetLocalStore(true)

	return scan.NewTableScanOperatorFactory(node.GetNodeID(), node.Table, outputColumns, filter)
}
