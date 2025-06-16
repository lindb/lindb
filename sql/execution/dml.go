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

package execution

import (
	"fmt"

	"github.com/lindb/common/pkg/encoding"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/lindb/lindb/constants"
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/utils"
	protoCommandV1 "github.com/lindb/lindb/proto/gen/v1/command"
	"github.com/lindb/lindb/spi/types"
	"github.com/lindb/lindb/sql/execution/buffer"
	"github.com/lindb/lindb/sql/execution/model"
	"github.com/lindb/lindb/sql/execution/pipeline"
	"github.com/lindb/lindb/sql/interfaces"
	"github.com/lindb/lindb/sql/planner/plan"
	"github.com/lindb/lindb/sql/planner/printer"
	"github.com/lindb/lindb/sql/rewrite"
	"github.com/lindb/lindb/sql/tree"
)

type DMLExecutionFactory struct {
	deps *Deps
}

func NewDMLExecutionFactory(deps *Deps) ExecutionFactory {
	return &DMLExecutionFactory{
		deps: deps,
	}
}

func (f *DMLExecutionFactory) CreateExecution(session *Session, statement *tree.PreparedStatement) Execution {
	return NewDMLExecution(session, f.deps, statement)
}

type DMLContext struct {
	output  buffer.OutputBuffer
	rsBuild *buffer.ResultSetBuild

	err string
}

func NewDMLContext() *DMLContext {
	rsBuild := buffer.CreateResultSetBuild()
	return &DMLContext{
		rsBuild: rsBuild,
		output:  buffer.NewQueryOutputBuffer(rsBuild),
	}
}

func (ctx DMLContext) GetOutput() buffer.OutputBuffer {
	return ctx.output
}

func (ctx *DMLContext) Wait() {
	go ctx.rsBuild.Process()
	// FIXME: add timeout
}

func (ctx *DMLContext) SetError(err string) {
	ctx.err = err
}

func (ctx *DMLContext) Error() string {
	return ctx.err
}

func (ctx *DMLContext) ResultSet() *model.ResultSet {
	return ctx.rsBuild.ResultSet()
}

type DMLExecution struct {
	session *Session
	context *DMLContext
	planner *Planner

	preparedStatement *tree.PreparedStatement
	deps              *Deps
}

func NewDMLExecution(session *Session, deps *Deps, preparedStatement *tree.PreparedStatement) Execution {
	return &DMLExecution{
		session:           session,
		deps:              deps,
		planner:           NewPlanner(deps.AnalyzerFct),
		preparedStatement: preparedStatement,
	}
}

func (exec *DMLExecution) Start() any {
	defer func() {
		// cleanup execution context
		fmt.Println("cleanup execution context")
		pipeline.DriverManager.Cleanup(exec.session.RequestID)
	}()

	exec.context = NewDMLContext()

	// rewrite statement
	statement := exec.rewrite(exec.preparedStatement.Statement)
	// plan statement
	statementPlan := exec.planner.Plan(exec.session, statement)
	// distribute plan
	fragmentedPlan := exec.planner.PlanDistribution(statementPlan)
	// scheduler start
	exec.execute(fragmentedPlan, exec.context.GetOutput())

	// waiting query complete
	exec.context.Wait()
	if exec.context.Error() != "" {
		panic(exec.context.Error())
	}
	return exec.context.ResultSet()
}

func (exec *DMLExecution) rewrite(statement tree.Statement) tree.Statement {
	rewrites := rewrite.NewStatementRewrite([]interfaces.Rewrite{
		NewExplainRewrite(exec.session, NewQueryExplainer(exec.planner)),
		rewrite.NewShowQueriesRewrite(exec.session.Database, exec.session.NodeIDAllocator),
	})
	// rewrite
	return rewrites.Rewrite(statement)
}

func (exec *DMLExecution) execute(fragmentedPlan *plan.SubPlan, output buffer.OutputBuffer) {
	printer := printer.NewPlanPrinter(printer.NewTextRender(0))
	fmt.Println(printer.PrintDistributedPlan(fragmentedPlan))
	session := exec.session

	fragments := fragmentedPlan.GetAllFragments()

	rootFragment := fragments[0]
	currentTime, _ := utils.GetInt64FromContext(session.Context, constants.ContextKeyCurrentTime)

	// submit all task
	for i := range len(fragments) {
		fragment := fragments[i]
		taskID := model.TaskID{
			RequestID: session.RequestID,
			ID:        i,
		}

		fmt.Printf("remote parent node=====%v\n", fragment.ParentNode)
		go func() {
			// TODO: handle panic
			if fragment.ParentNode == nil {
				outputCh := make(chan *types.Page)
				// execute task under current node if it has no parent
				defer func() {
					close(outputCh)
					if err := recover(); err != nil {
						exec.context.SetError(fmt.Sprintf("%v", err))
					}
					// TODO::
					// close(exec.queryContext.completed)
				}()
				// run under current node
				taskFct := NewTaskExecutionFactory()
				taskExec := taskFct.Create(&SQLTask{
					currentTime: currentTime,
					id:          taskID,
					fragment:    rootFragment,
				})
				go func() {
					for page := range outputCh {
						output.AddPage(page)
					}
					output.Complete()
				}()
				taskExec.Execute(outputCh)
			} else {
				// execute task under remote node, send fragment to remote execution node
				fragment.Receivers = []models.InternalNode{*exec.deps.CurrentNode}
				data := encoding.JSONMarshal(fragment)

				for node, shards := range fragment.Partitions {
					exec.sendTask(node, taskID, shards, currentTime, data)
				}
			}
		}()
	}

	fmt.Println("done.......")
}

func (exec *DMLExecution) sendTask(node models.InternalNode, taskID model.TaskID,
	shards []int, currentTime int64, data []byte,
) {
	conn, err := grpc.Dial(node.Address(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	client := protoCommandV1.NewCommandServiceClient(conn)
	_, err = client.Command(exec.session.Context, &protoCommandV1.CommandRequest{
		Cmd: protoCommandV1.Command_SubmitTask,
		Payload: encoding.JSONMarshal(&model.TaskRequest{
			RequestContext: model.RequestContext{
				CurrentTime: currentTime,
			},
			TaskID:     taskID,
			Fragment:   data,
			Partitions: shards,
		}),
	})
	if err != nil {
		fmt.Println(exec.session.Context.Err())
		// TODO: check panic
		panic(err)
	}
}
