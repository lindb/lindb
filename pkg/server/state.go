package server

// State represents server current state
type State int

const (
	// New is inactive
	New State = iota
	// Running is running
	Running
	// Failed has encountered a problem
	Failed
	// Terminated is stopped
	Terminated
)
