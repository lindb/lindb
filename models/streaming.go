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

package models

import (
	"fmt"
)

type Streaming struct {
	Name     string `json:"name"`
	Observer string `json:"observer"`
	Database string `json:"database"`
}

func (s *Streaming) String() string {
	return fmt.Sprintf(`create streaming "%s"(observer "%s", database "%s")`, s.Name, s.Observer, s.Database)
}

type ConsumeAssignment struct {
	ConsumerID NodeID
	Shards     []ShardID
}

type StreamingState struct {
	Config             Streaming                `json:"config"`
	Consumers          map[NodeID]StatelessNode `json:"consumers"`
	ConsumeAssignments []ConsumeAssignment      `json:"consumeAssignments"`
}

// ModifyStreamingJob represents modifying streaming job event.
type ModifyStreamingJob struct {
	Streaming string
	JobName   string
	Script    string
}

type DeleteStreamingJob struct {
	Streaming string
	JobName   string
}

type DeleteStreaming struct {
	Streaming string
}

type DeleteDatabase struct {
	Database string
}
