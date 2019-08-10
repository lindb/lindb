package task

import (
	"encoding/json"
)

// State represents the task state.
type State uint8

const (
	// StateCreated is created
	StateCreated State = iota
	// StateRunning is running
	StateRunning
	// StateDoneOK is done
	StateDoneOK
	// StateDoneErr is done, but got error
	StateDoneErr
)

var statestrs = [...]string{
	"StateCreated",
	"StateRunning",
	"StateDoneOK",
	"StateDoneErr",
}

func (st State) String() string {
	return statestrs[int(st)]
}

type (
	// Kind is the task kind.
	Kind string
	// Task is the actual task with parameters and status.
	Task struct {
		Kind     Kind            `json:"kind"`
		Name     string          `json:"name"`
		Executor string          `json:"executor"`
		Params   json.RawMessage `json:"params"`
		State    State           `json:"state"`
		ErrMsg   string          `json:"err_msg,omitempty"`
	}
	groupedTasks struct {
		State State  `json:"state"`
		Tasks []Task `json:"tasks"`
	}
)
