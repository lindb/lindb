package cli

import (
	"os"
	"os/signal"
	"syscall"
)

type CommandLine struct {
	osSignal chan os.Signal
}

// New an instance CommandLine
func New() *CommandLine {
	return &CommandLine{
		osSignal: make(chan os.Signal, 1),
	}
}

// Run CommandLine
func (c *CommandLine) Run() error {

	// register OS signals for graceful termination
	signal.Notify(c.osSignal, os.Interrupt, os.Kill, syscall.SIGINT, syscall.SIGTERM)

	return c.mainLoop()
}

// mainLoop runs the main prompt loop for CommandLine
func (c *CommandLine) mainLoop() error {
	for {
		select {
		case <-c.osSignal:
			c.exit()
			return nil
		}
	}
}

// exit CommandLine
func (c *CommandLine) exit() {

}
