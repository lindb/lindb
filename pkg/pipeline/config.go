package pipeline

// Config represents the config of a Stage
type Config interface {
	// GetTaskSize returns the number of Tasks
	GetTaskSize() int
	// NewTask returns a Task instance
	NewTask() Task
}
