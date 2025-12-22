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

package aggregation

import (
	"strings"
	"sync"

	"go.uber.org/atomic"
)

type StringMapping struct {
	strToID map[string]uint32
	idToStr map[uint32]string

	id atomic.Uint32

	mutex sync.Mutex
}

func NewStringMapping() *StringMapping {
	return &StringMapping{
		strToID: make(map[string]uint32),
		idToStr: make(map[uint32]string),
	}
}

func (m *StringMapping) GetID(s string) uint32 {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if id, ok := m.strToID[s]; ok {
		return id
	}

	/*
		make a new copy for s in order to remove references from bigger string slice
		for example:
		big := "a very large string..."
		s := big[10:20]  // s is a substring (slice) of big
	*/
	newStr := strings.Clone(s)

	id := m.id.Inc()
	m.strToID[newStr] = id
	m.idToStr[id] = newStr
	return id
}

func (m *StringMapping) GetString(id uint32) string {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	return m.idToStr[id]
}
