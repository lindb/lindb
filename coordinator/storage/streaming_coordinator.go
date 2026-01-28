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

package storage

import (
	"path/filepath"

	"github.com/lindb/common/pkg/encoding"
	"github.com/lindb/common/pkg/logger"

	"github.com/lindb/lindb/models"
)

// onStreamingStateChange processes streaming state change event.
func (m *stateManager) onStreamingStateChange(key string, data []byte) error {
	m.logger.Info("streaming state changed", logger.String("key", key), logger.String("data", string(data)))

	var state models.StreamingState
	if err := encoding.JSONUnmarshal(data, &state); err != nil {
		m.logger.Error("unmarshal streaming state failed when do streaming state changed",
			logger.String("data", string(data)), logger.Error(err))
		return err
	}

	for _, watcher := range m.watchers {
		watcher.OnEvent(&state)
	}
	return nil
}

func (m *stateManager) onStreamingStateDeletion(key string) error {
	m.logger.Info("streaming state deleted", logger.String("key", key))

	_, streaming := filepath.Split(key)

	event := models.DeleteStreaming{
		Streaming: streaming,
	}
	for _, watcher := range m.watchers {
		watcher.OnEvent(&event)
	}
	return nil
}
