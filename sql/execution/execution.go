package execution

import (
	"context"
	"fmt"

	"github.com/lindb/common/pkg/encoding"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/lindb/lindb/constants"
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/utils"
	protoCommandV1 "github.com/lindb/lindb/proto/gen/v1/command"
	"github.com/lindb/lindb/sql/execution/buffer"
	"github.com/lindb/lindb/sql/execution/model"
	"github.com/lindb/lindb/sql/execution/pipeline"
	"github.com/lindb/lindb/sql/interfaces"
	"github.com/lindb/lindb/sql/planner/plan"
	"github.com/lindb/lindb/sql/planner/printer"
	"github.com/lindb/lindb/sql/rewrite"
	"github.com/lindb/lindb/sql/tree"
)

type Execution interface {
	Start() any
}

type DataDefinitionExecution struct {
	task DataDefinitionTask
}

func NewDataDefinitionExecution(task DataDefinitionTask) Execution {
	return &DataDefinitionExecution{
		task: task,
	}
}

func (exec *DataDefinitionExecution) Start() any {
	err := exec.task.Execute(context.TODO())
	// TODO: add log
	fmt.Println(err)

	fmt.Println("execution task")
	fmt.Println(exec.task.Name())
	return nil
}

type QueryContext struct {
	output  buffer.OutputBuffer
	rsBuild *buffer.ResultSetBuild

	completed chan struct{}
}

func NewQueryContext() *QueryContext {
	rsBuild := buffer.CreateResultSetBuild()
	return &QueryContext{
		rsBuild:   rsBuild,
		output:    buffer.NewQueryOutputBuffer(rsBuild),
		completed: make(chan struct{}),
	}
}

func (ctx QueryContext) GetOutput() buffer.OutputBuffer {
	return ctx.output
}

func (ctx *QueryContext) Wait() {
	go ctx.rsBuild.Process()
	// FIXME: add timeout
	<-ctx.completed
	ctx.rsBuild.Complete()
}

func (ctx *QueryContext) ResultSet() *model.ResultSet {
	return ctx.rsBuild.ResultSet()
}

type QueryExecution struct {
	session      *Session
	queryContext *QueryContext
	planner      *Planner

	preparedStatement *tree.PreparedStatement
	deps              *Deps
}

func NewQueryExecution(session *Session, deps *Deps, preparedStatement *tree.PreparedStatement) Execution {
	return &QueryExecution{
		session:           session,
		deps:              deps,
		planner:           NewPlanner(deps.AnalyzerFct),
		preparedStatement: preparedStatement,
	}
}

func (exec *QueryExecution) Start() any {
	defer func() {
		// cleanup execution context
		pipeline.DriverManager.Cleanup(exec.session.RequestID)
	}()

	exec.queryContext = NewQueryContext()

	// rewrite statement
	statement := exec.rewrite(exec.preparedStatement.Statement)
	// plan statement
	plan := exec.planner.Plan(exec.session, statement)
	// distribute plan
	fragmentedPlan := exec.planner.PlanDistribution(plan)
	// scheduler start
	exec.execute(fragmentedPlan, exec.queryContext.GetOutput())

	// waiting query complete
	exec.queryContext.Wait()
	return exec.queryContext.ResultSet()
}

func (exec *QueryExecution) rewrite(statement tree.Statement) tree.Statement {
	rewrites := rewrite.NewStatementRewrite([]interfaces.Rewrite{
		NewExplainRewrite(exec.session, NewQueryExplainer(exec.planner)),
		rewrite.NewShowQueriesRewrite(exec.session.Database),
	})
	// rewrite
	return rewrites.Rewrite(statement)
}

func (exec *QueryExecution) execute(fragmentedPlan *plan.SubPlan, output buffer.OutputBuffer) {
	printer := printer.NewPlanPrinter(printer.NewTextRender(0))
	fmt.Println(printer.PrintDistributedPlan(fragmentedPlan))
	session := exec.session

	fragments := fragmentedPlan.GetAllFragments()

	rootFragment := fragments[0]
	currentTime, _ := utils.GetInt64FromContext(session.Context, constants.ContextKeyCurrentTime)

	// submit all task
	for i := 0; i < len(fragments); i++ {
		taskID := model.TaskID{
			RequestID: session.RequestID,
			ID:        i,
		}
		fragment := fragments[i]

		fmt.Printf("remote parent node=====%v\n", fragment.RemoteParentNodeID)
		go func() {
			// TODO: handle panic
			if fragment.RemoteParentNodeID == nil {
				// run under current node
				taskFct := NewTaskExecutionFactory()
				taskExec := taskFct.Create(&SQLTask{
					currentTime: currentTime,
					id:          taskID,
					fragment:    rootFragment,
				}, output)
				taskExec.Execute()
				// TODO::
				close(exec.queryContext.completed)
			} else {
				fragment.Receivers = []models.InternalNode{*exec.deps.CurrentNode}
				for node, shards := range fragment.Partitions {
					// run under remote node, send fragment to remote execution node
					data := encoding.JSONMarshal(fragment)
					conn, err := grpc.Dial(node.Address(), grpc.WithTransportCredentials(insecure.NewCredentials()))
					if err != nil {
						panic(err)
					}
					defer conn.Close()

					client := protoCommandV1.NewCommandServiceClient(conn)
					_, err = client.Command(context.TODO(), &protoCommandV1.CommandRequest{
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
						// TODO: check panic
						panic(err)
					}
				}
			}
		}()
	}

	fmt.Println("done.......")
}
