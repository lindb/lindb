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

package query

import "sync"

var (
	mgr      PipelineManager
	once4Mgr sync.Once
)

// PipelineManager represents the manager which store current exuecting Pipeline.
type PipelineManager interface {
	// AddPipeline adds a Pipeline when it starts.
	AddPipeline(requestID string, pipeline Pipeline)
	// RemovePipeline removes a Pipeline when it completed.
	RemovePipeline(requestID string)
	// GetPipeline returns a Pipeline by given request id, if not exist then return nil.
	GetPipeline(requestID string) Pipeline
	// GetAllAlivePipelines returns all alive request ids for pipelines.
	GetAllAlivePipelines() []string
}

// GetPipelineManager returns a singleton PipelineManager instance.
func GetPipelineManager() PipelineManager {
	if mgr != nil {
		return mgr
	}
	once4Mgr.Do(func() {
		mgr = newPipelineManager()
	})
	return mgr
}

// pipelineManager implements PipelineManager interface.
type pipelineManager struct {
	pipelines map[string]Pipeline

	mutex sync.RWMutex
}

// newPipelineManager creates a PipelineManager instance.
func newPipelineManager() PipelineManager {
	return &pipelineManager{
		pipelines: make(map[string]Pipeline),
	}
}

// AddPipeline adds a Pipeline when it starts.
func (m *pipelineManager) AddPipeline(requestID string, pipeline Pipeline) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	m.pipelines[requestID] = pipeline
}

// GetPipeline returns a Pipeline by given request id, if not exist then return nil.
func (m *pipelineManager) GetPipeline(requestID string) Pipeline {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	p := m.pipelines[requestID]
	return p
}

// RemovePipeline removes a Pipeline when it completed.
func (m *pipelineManager) RemovePipeline(requestID string) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	delete(m.pipelines, requestID)
}

// GetAllAlivePipelines returns all alive request ids for pipelines.
func (m *pipelineManager) GetAllAlivePipelines() (rs []string) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	for reqID := range m.pipelines {
		rs = append(rs, reqID)
	}
	return
}
