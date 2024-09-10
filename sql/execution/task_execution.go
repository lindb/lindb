package execution

import (
	"github.com/lindb/lindb/sql/context"
	"github.com/lindb/lindb/sql/execution/buffer"
	"github.com/lindb/lindb/sql/planner"
)

type TaskExecutionFactory struct{}

func NewTaskExecutionFactory() *TaskExecutionFactory {
	return &TaskExecutionFactory{}
}

func (fct *TaskExecutionFactory) Create(task *SQLTask, output buffer.OutputBuffer) *TaskExecution {
	planner := planner.NewTaskExecutionPlanner()

	ctx := &context.TaskContext{
		TaskID:     task.id,
		Fragment:   task.fragment,
		Partitions: task.partitions,
	}
	plan := planner.Plan(ctx, task.fragment.Root, output)

	return &TaskExecution{
		taskCtx: ctx,
		plan:    plan,
	}
}

type TaskExecution struct {
	taskCtx *context.TaskContext
	plan    *planner.TaskExecutionPlan
}

func (exe *TaskExecution) Execute() {
	pipelines := exe.plan.GetPipelines()
	for i := range pipelines {
		pipeline := pipelines[i]
		pipeline.Run()
	}
}
