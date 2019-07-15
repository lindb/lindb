package pipeline

// Pipeline contains a list of Stage
type Pipeline struct {
	stages []*Stage
}

// AddStage adds a Stage to the Pipeline
// establish the relation of the pre Stage and current Stage
func (pipeline *Pipeline) AddStage(stage *Stage) {
	pipeline.stages = append(pipeline.stages, stage)
	length := len(pipeline.stages)
	if length > 1 {
		preStage := pipeline.stages[length-2]
		currentStage := pipeline.stages[length-1]
		preStage.nextRouter = currentStage.currentRouter
		for i := 0; i < len(preStage.runs); i++ {
			preStage.runs[i].Task.SetRouter(preStage.nextRouter)
		}
	}
}

// Tell sends a message to the Pipeline
func (pipeline *Pipeline) Tell(ctx TaskContext, message Message) {
	pipeline.stages[0].currentRouter.Tell(ctx, message)
}

// Shutdown shutdown all stages
func (pipeline *Pipeline) Shutdown() {
	for i := 0; i < len(pipeline.stages); i++ {
		pipeline.stages[i].Shutdown()
	}
}
