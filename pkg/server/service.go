package server

// Service prepresents an operational state of server, lifecycle methods to transition between states.
type Service interface {
	// Run runs server
	Run() error
	// State returns current service state
	State() State
	// Stop shutdowns server, do some cleanup logic
	Stop() error
}
