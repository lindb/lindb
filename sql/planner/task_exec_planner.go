// Licensed to LinDB under one or more contributor
// license agreements. See the NOTICE file distributed with
// this work for additional information regarding copyright
// ownership. LinDB licenses this file to you under
// the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

package planner

import (
	"fmt"

	"github.com/samber/lo"

	"github.com/lindb/lindb/spi"
	"github.com/lindb/lindb/spi/types"
	"github.com/lindb/lindb/sql/context"
	"github.com/lindb/lindb/sql/execution/operator"
	"github.com/lindb/lindb/sql/execution/operator/exchange"
	"github.com/lindb/lindb/sql/execution/operator/join"
	"github.com/lindb/lindb/sql/execution/operator/output"
	"github.com/lindb/lindb/sql/execution/operator/scan"
	"github.com/lindb/lindb/sql/execution/operator/streaming"
	"github.com/lindb/lindb/sql/execution/operator/streaming/executor"
	"github.com/lindb/lindb/sql/execution/pipeline"
	planpkg "github.com/lindb/lindb/sql/planner/plan"
	"github.com/lindb/lindb/sql/tree"
)

type TaskExecutionPlanner struct{}

func NewTaskExecutionPlanner() *TaskExecutionPlanner {
	return &TaskExecutionPlanner{}
}

func (p *TaskExecutionPlanner) Plan(taskCtx *context.TaskContext, node planpkg.PlanNode) *TaskExecutionPlan {
	visitor := &TaskExecutionPlanVisitor{
		taskExecCtx: taskCtx,
		planner:     p,
	}
	taskExecPlanCtx := NewTaskExecutionPlanContext(taskCtx)
	var (
		op operator.Operator
		ok bool
	)
	if op, ok = node.Accept(taskExecPlanCtx, visitor).(operator.Operator); !ok {
		panic("cannot get physicalOperator")
	}
	fmt.Printf("task exec plan:%T,op:%T\n", node, op)
	return NewTaskExecutionPlan([]*pipeline.Pipeline{pipeline.NewPipeline(taskCtx, op)})
}

type TaskExecutionPlanVisitor struct {
	taskExecCtx *context.TaskContext
	planner     *TaskExecutionPlanner
}

// Visit visits all plan node and plans task execution physical operator.
func (v *TaskExecutionPlanVisitor) Visit(context any, n planpkg.PlanNode) (r any) {
	fmt.Printf("task exec plan visit: %T\n", n)
	switch node := n.(type) {
	case *planpkg.OutputNode:
		fmt.Println("output node...")
		if v.taskExecCtx.IsStreaming() {
			child := node.Source.Accept(context, v)
			return streaming.NewOutputOperator(v.taskExecCtx.Context, v.taskExecCtx.Database, v.taskExecCtx.StreamName,
				node, child.(operator.Operator))
		}
		child := node.Source.Accept(context, v).(operator.Operator)
		return output.NewRSOutputOperator(node, child)
	case *planpkg.InsertNode:
		return v.visitInsert(context, node)
	case *planpkg.AggregationNode:
		return v.visitAggregation(context, node)
	case *planpkg.RemoteSourceNode:
		return v.visitRemoteSource(context, node)
	case *planpkg.ExchangeNode:
		return v.visitExchange(context, node)
	case *planpkg.JoinNode:
		return v.visitJoin(context, node)
	case *planpkg.ProjectionNode:
		return v.visitProjection(context, node)
	case *planpkg.FilterNode:
		return v.visitFilter(context, node)
	case *planpkg.TableScanNode:
		return v.VisitTableScan(context, node)
	case *planpkg.ValuesNode:
		return v.visitValues(context, node)
	default:
		panic(fmt.Sprintf("umimplements task planner %T", n))
	}
}

func (v *TaskExecutionPlanVisitor) visitInsert(context any, node *planpkg.InsertNode) any {
	source := node.Source.Accept(context, v).(operator.Operator)
	return streaming.NewInsertOperator(v.taskExecCtx.Context, node, source)
}

func (v *TaskExecutionPlanVisitor) visitJoin(context any, node *planpkg.JoinNode) any {
	left := node.Left.Accept(context, v).(operator.Operator)
	right := node.Right.Accept(context, v).(operator.Operator)

	return join.NewHashJoinOperator(node, left, right)
}

func (v *TaskExecutionPlanVisitor) visitValues(_ any, node *planpkg.ValuesNode) (r any) {
	return operator.NewValuesOperator(node)
}

// visitFilter plans filter physical operator.
func (v *TaskExecutionPlanVisitor) visitFilter(context any, node *planpkg.FilterNode) (r any) {
	if tableScan, ok := node.Source.(*planpkg.TableScanNode); ok {
		return v.visitTableScan(context, tableScan, node.Predicate)
		// FIXME: source layout???
		// return NewPhysicalOperation(operatorFct, node.GetOutputSymbols(), nil)
	}
	panic("need impl visitFilter")
}

func (v *TaskExecutionPlanVisitor) visitExchange(context any, node *planpkg.ExchangeNode) (r any) {
	if node.Scope != planpkg.Local {
		panic("only local exchanges are supported in the local planner")
	}
	// FIXME: set child
	child := node.Sources[0].Accept(context, v).(operator.Operator)
	return exchange.NewLocalExchangeOperator(node, child)
}

func (v *TaskExecutionPlanVisitor) visitAggregation(context any, node *planpkg.AggregationNode) (r any) {
	source := node.Source.Accept(context, v).(operator.Operator)
	if v.taskExecCtx.IsStreaming() {
		var assignments []*planpkg.Assignment
		if projection, ok := node.Source.(*planpkg.ProjectionNode); ok {
			assignments = projection.Assignments
		}
		return streaming.NewTimeWindowOperator(node, executor.NewHashGrouping(node, assignments), source)
	}
	return v.planGroupByAggregation(node, source)
}

func (v *TaskExecutionPlanVisitor) planGroupByAggregation(
	node *planpkg.AggregationNode, source operator.Operator,
) operator.Operator {
	// TODO: need fixit
	return v.createHashAggregationOperatorFactory(node, source)
}

func (v *TaskExecutionPlanVisitor) createHashAggregationOperatorFactory(
	node *planpkg.AggregationNode, source operator.Operator,
) operator.Operator {
	return operator.NewHashAggregationOperator(node, source)
}

func (v *TaskExecutionPlanVisitor) visitProjection(context any, node *planpkg.ProjectionNode) (r any) {
	var source planpkg.PlanNode
	var predicate tree.Expression
	if filterNode, ok := node.Source.(*planpkg.FilterNode); ok {
		source = filterNode.Source
		predicate = filterNode.Predicate
	} else {
		source = node.Source
	}
	return v.visitScanFilterAndProjection(context, node, source, predicate)
}

func (v *TaskExecutionPlanVisitor) VisitTableScan(context any, node *planpkg.TableScanNode) (r any) {
	return v.visitTableScan(context, node, nil)
}

func (v *TaskExecutionPlanVisitor) visitRemoteSource(_ any, node *planpkg.RemoteSourceNode) (r any) {
	op := exchange.NewRemoteExchangeOperator(v.taskExecCtx.Context, node, len(node.SourceFragmentIDs))
	// register remote exchange source
	pipeline.DriverManager.RegisterSourceOperator(v.taskExecCtx.TaskID, op)
	return op
}

func (v *TaskExecutionPlanVisitor) visitScanFilterAndProjection(context any,
	project *planpkg.ProjectionNode, sourceNode planpkg.PlanNode, predicate tree.Expression,
) any {
	fmt.Printf("visitScanFilterAndProjection:%T,filter=%v\n", sourceNode, predicate)
	if tableScan, ok := sourceNode.(*planpkg.TableScanNode); ok {
		child := v.visitTableScan(context, tableScan, predicate)
		if v.taskExecCtx.IsStreaming() {
			return operator.NewProjectionOperator(v.taskExecCtx.Context, project, child)
		}
		return child
	}
	// plan source node
	child := sourceNode.Accept(context, v).(operator.Operator)
	return operator.NewProjectionOperator(v.taskExecCtx.Context, project, child)
}

func (v *TaskExecutionPlanVisitor) visitTableScan(_ any,
	node *planpkg.TableScanNode, predicate tree.Expression,
) operator.Operator {
	fmt.Printf("table node=%v\n", predicate)
	outputs := node.GetOutputSymbols()
	outputColumns := lo.Map(outputs, func(item *planpkg.Symbol, index int) types.ColumnMetadata {
		return types.ColumnMetadata{
			Name:     item.Name,
			DataType: item.DataType,
		}
	})
	provider := spi.GetSourceConnectorProvider(node.Table)
	connector := provider.CreateSourceConnector(v.taskExecCtx.Context,
		node.Table, v.taskExecCtx.Partitions, node.ColumnMapping, predicate, outputColumns, node.Assignments)

	return scan.NewTableScanOperator(connector, node)
}
