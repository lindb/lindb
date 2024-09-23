package execution

import (
	"context"
	"fmt"

	"github.com/lindb/common/pkg/encoding"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/lindb/lindb/models"
	protoCommandV1 "github.com/lindb/lindb/proto/gen/v1/command"
	sqlContext "github.com/lindb/lindb/sql/context"
	"github.com/lindb/lindb/sql/execution/buffer"
	"github.com/lindb/lindb/sql/execution/model"
	"github.com/lindb/lindb/sql/execution/pipeline"
	"github.com/lindb/lindb/sql/planner"
	"github.com/lindb/lindb/sql/planner/iterative"
	"github.com/lindb/lindb/sql/planner/iterative/rule"
	"github.com/lindb/lindb/sql/planner/optimization"
	"github.com/lindb/lindb/sql/planner/printer"
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
	session        *Session
	queryContext   *QueryContext
	plannerContext *sqlContext.PlannerContext

	preparedStatement *tree.PreparedStatement
	deps              *Deps
}

func NewQueryExecution(session *Session, deps *Deps, preparedStatement *tree.PreparedStatement) Execution {
	return &QueryExecution{
		session:           session,
		deps:              deps,
		preparedStatement: preparedStatement,
	}
}

func (exec *QueryExecution) Start() any {
	defer func() {
		// cleanup execution context
		pipeline.DriverManager.Cleanup(exec.session.RequestID)
	}()

	exec.queryContext = NewQueryContext()
	exec.plannerContext = sqlContext.NewPlannerContext(
		exec.session.Context,
		exec.session.Database,
		exec.session.NodeIDAllocator,
		exec.preparedStatement.Statement,
	)

	exec.analyze()

	plan := exec.planQuery(exec.queryContext.GetOutput())
	exec.planDistribution(plan)
	// scheduler start

	// waiting query complete
	exec.queryContext.Wait()
	return exec.queryContext.ResultSet()
}

func (exec *QueryExecution) analyze() {
	// rewrite
	rewrittenStatement := exec.deps.StatementRewrite.Rewrite(exec.preparedStatement.Statement)
	// create analyzer
	analyzer := exec.deps.AnalyzerFct.CreateAnalyzer(exec.plannerContext.AnalyzerContext)
	// do analyze
	analyzer.Analyze(rewrittenStatement)
}

func (exec *QueryExecution) planQuery(output buffer.OutputBuffer) *PlanRoot {
	// FIXME: fixme:

	// plan query
	planOptimizers := []optimization.PlanOptimizer{
		// optimization.NewPruneColumns(),
		iterative.NewIterativeOptimizer([]iterative.Rule{
			rule.NewRemoveRedundantIdentityProjections(),
		}),
		// column pruning optimizer
		iterative.NewIterativeOptimizer([]iterative.Rule{
			rule.NewPruneAggregationSourceColumns(),
			rule.NewPruneFilterColumns(),
			rule.NewPruneOutputSourceColumns(),
			rule.NewPruneProjectionColumns(),
			rule.NewPruneTableScanColumns(),
		}),
		iterative.NewIterativeOptimizer([]iterative.Rule{
			rule.NewRemoveRedundantIdentityProjections(),
		}),
		// push into table scan optimizer
		iterative.NewIterativeOptimizer([]iterative.Rule{
			rule.NewPushProjectionIntoTableScan(),
			rule.NewPushAggregationIntoTableScan(),
		}),
		optimization.NewAddExchanges(),
		optimization.NewAddLocalExchanges(),
		iterative.NewIterativeOptimizer([]iterative.Rule{
			rule.NewPushPartialAggregationThroughExchange(),
		}),
		iterative.NewIterativeOptimizer([]iterative.Rule{
			rule.NewRemoveRedundantIdentityProjections(),
		}),
	}
	logicalPlanner := planner.NewLogicalPlanner(exec.plannerContext, planOptimizers)
	plan := logicalPlanner.Plan()

	// fragment the plan
	fragmenter := planner.NewPlanFragmenter()
	fragmentedPlan := fragmenter.CreateSubPlans(plan)

	printer := printer.NewPlanPrinter(printer.NewTextRender(0))
	fmt.Println("******************")
	fmt.Println(printer.PrintLogicPlan(plan.Root))
	fmt.Println("******************")
	fmt.Println(printer.PrintDistributedPlan(fragmentedPlan))
	session := exec.session

	fragments := fragmentedPlan.GetAllFragments()

	rootFragment := fragments[0]

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
					id:       taskID,
					fragment: rootFragment,
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
	return &PlanRoot{
		root: fragmentedPlan,
	}
}

func (exec *QueryExecution) planDistribution(plan *PlanRoot) {
}
