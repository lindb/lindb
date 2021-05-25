// Licensed to LinDB under one or more contributor
// license agreements. See the NOTICE file distributed with
// this work for additional information regarding copyright
// ownership. LinDB licenses this file to you under
// the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

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
