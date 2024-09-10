package planner

import (
	"go.uber.org/atomic"

	"github.com/lindb/lindb/spi"
	"github.com/lindb/lindb/sql/execution/pipeline"
)

type TaskExecutionPlanContext struct {
	driverFactories []*pipeline.DriverFactory
	splitSources    []spi.SplitSource

	nextPipelineID atomic.Int32
	localStore     bool
}

func NewTaskExecutionPlanContext(driverFactories []*pipeline.DriverFactory) *TaskExecutionPlanContext {
	return &TaskExecutionPlanContext{
		driverFactories: driverFactories,
	}
}

func (ctx *TaskExecutionPlanContext) AddDriverFactory(physicalOperation *PhysicalOperation) {
	// FIXME: add lookup outer driver?
	driverFct := pipeline.NewDriverFactory(ctx.nextPipelineID.Inc(), physicalOperation.operatorFactories)
	ctx.driverFactories = append(ctx.driverFactories, driverFct)
}

func (ctx *TaskExecutionPlanContext) SetSplitSources(splitSources []spi.SplitSource) {
	ctx.splitSources = splitSources
}

func (ctx *TaskExecutionPlanContext) SetLocalStore(local bool) {
	ctx.localStore = local
}

func (ctx *TaskExecutionPlanContext) IsLocalStore() bool {
	return ctx.localStore
}

type TaskExecutionPlan struct {
	pipelines []*pipeline.Pipeline
}

func NewTaskExecutionPlan(pipelines []*pipeline.Pipeline) *TaskExecutionPlan {
	return &TaskExecutionPlan{
		pipelines: pipelines,
	}
}

func (p *TaskExecutionPlan) GetPipelines() []*pipeline.Pipeline {
	return p.pipelines
}
