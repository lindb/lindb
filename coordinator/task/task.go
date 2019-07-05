package task

import "encoding/json"

// State represents the task state.
type State uint8

const (
	StateCreated State = iota
	StateRunning
	StateDoneOK
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

// UnsafeMarshal marshals itself by JSON encoder, it will panic if error occurs.
func (t Task) UnsafeMarshal() []byte {
	data, err := json.Marshal(t)
	if err != nil {
		panic(err)
	}
	return data
}

// UnsafeUnmarshal unmarshals itself by JSON decoder, it will panic if error occurs.
func (t *Task) UnsafeUnmarshal(data []byte) {
	if err := json.Unmarshal(data, t); err != nil {
		panic(err)
	}
}

func (gt groupedTasks) UnsafeMarshal() []byte {
	data, err := json.Marshal(gt)
	if err != nil {
		panic(err)
	}
	return data
}

func (gt *groupedTasks) UnsafeUnmarshal(data []byte) {
	if err := json.Unmarshal(data, gt); err != nil {
		panic(err)
	}
}
