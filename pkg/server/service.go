package server

//go:generate mockgen -source=./service.go -destination=./service_mock.go -package=server

// Service represents an operational state of server, lifecycle methods to transition between states.
type Service interface {
	// Name returns the service's name
	Name() string
	// Run runs server
	Run() error
	// State returns current service state
	State() State
	// Stop shutdowns server, do some cleanup logic
	Stop() error
}
