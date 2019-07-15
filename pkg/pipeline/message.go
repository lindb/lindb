package pipeline

// Message
type Message interface {
	// SetContext bind a message with a context
	SetContext(ctx TaskContext)

	// GetContext returns the context
	GetContext() TaskContext
}

// ShutdownMessage implements Message
type ShutdownMessage struct {
	Message
}
