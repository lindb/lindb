package pipeline

// Stage is a part of Pipeline
// the config contains the number of the Runtime Tasks and can new a Task
// runs is a lot of Runtime Tasks for processing messages and generating new messages
// currentRouter is used for pre Stage routing new messages
// nextRouter is used for currentStage routing new messages
type Stage struct {
	config        Config
	runs          []*Runtime
	currentRouter Router
	nextRouter    Router
}

// NewStage returns a Stage only containing currentRouter currently
func NewStage(config Config) *Stage {
	runs := make([]*Runtime, config.GetTaskSize(), config.GetTaskSize())
	for i := 0; i < len(runs); i++ {
		runs[i], _ = NewTaskRuntime(config.NewTask())
	}
	return &Stage{
		config:        config,
		runs:          runs,
		currentRouter: NewRuntimeRouter(runs),
	}
}

// Shutdown shutdown all Runtime Tasks
func (stage *Stage) Shutdown() {
	message := new(ShutdownMessage)
	for i := 0; i < len(stage.runs); i++ {
		stage.runs[i].Tell(message)
	}
}
