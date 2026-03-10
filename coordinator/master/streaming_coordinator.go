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

package master

import (
	"fmt"
	"strings"

	"github.com/google/go-cmp/cmp"
	"github.com/lindb/common/pkg/encoding"
	"github.com/lindb/common/pkg/logger"
	"github.com/samber/lo"

	"github.com/lindb/lindb/constants"
	"github.com/lindb/lindb/models"
)

// onStreamingCfgChange triggers when streaming config create/modify.
func (m *stateManager) onStreamingCfgChange(key string, data []byte) error {
	m.logger.Info("streaming config changed",
		logger.String("key", key),
		logger.String("data", string(data)))

	cfg := &models.Streaming{}
	if err := encoding.JSONUnmarshal(data, &cfg); err != nil {
		m.logger.Error("streaming config is changed, but unmarshal error",
			logger.Error(err))
		return err
	}

	m.streamings[cfg.Name] = cfg
	return nil
}

// onStreamingCfgDelete triggers when streaming config is deletion.
func (m *stateManager) onStreamingCfgDelete(key string) error {
	m.logger.Info("streaming config deleted",
		logger.String("key", key))
	name := strings.TrimPrefix(key, constants.GetStreamingConfigPath(""))
	_, ok := m.streamings[name]
	if !ok {
		return fmt.Errorf("streaming config not found: %s", name)
	}
	delete(m.streamings, name)
	delete(m.streamingStates, name)
	return nil
}

// onObserverNodeStartup triggers when observer node online.
func (m *stateManager) onObserverNodeStartup(key string, data []byte) error {
	m.logger.Info("new observer node online in streaming observer cluster",
		logger.String("key", key),
		logger.String("data", string(data)))

	consumerGroup, _, err := constants.ParseObserverLiveNode(key)
	if err != nil {
		m.logger.Error("parse observer live node path error",
			logger.String("key", key),
			logger.Error(err))
		return err
	}
	node := models.StatelessNode{}
	if err := encoding.JSONUnmarshal(data, &node); err != nil {
		m.logger.Error("new observer node online in streaming observer cluster but unmarshal error", logger.Error(err))
		return err
	}
	nodes, ok := m.observerNodes[consumerGroup]
	if !ok {
		m.observerNodes[consumerGroup] = []models.StatelessNode{node}
	} else {
		nodes = append(nodes, node)
		// remove duplicate node
		nodes = lo.UniqBy(nodes, func(n models.StatelessNode) string {
			return n.Indicator()
		})
		m.observerNodes[consumerGroup] = nodes
	}
	return nil
}

// onObserverNodeFailure triggers when observer node offline.
func (m *stateManager) onObserverNodeFailure(key string) error {
	m.logger.Info("observer node offline in streaming observer cluster",
		logger.String("key", key))

	consumerGroup, nodeKey, err := constants.ParseObserverLiveNode(key)
	if err != nil {
		m.logger.Error("parse observer live node path error",
			logger.String("key", key),
			logger.Error(err))
		return err
	}
	nodes, ok := m.observerNodes[consumerGroup]
	if !ok {
		m.logger.Warn("observer nodes not found for streaming",
			logger.String("streaming", consumerGroup))
		return nil
	}
	// delete node from observer nodes list
	nodes = lo.Filter(nodes, func(n models.StatelessNode, _ int) bool {
		return n.Indicator() != nodeKey
	})
	m.observerNodes[consumerGroup] = nodes
	return nil
}

// consumerReassign reassigns consumers when observer nodes change.
func (m *stateManager) consumerReassign() {
	for streamingName, cfg := range m.streamings {
		observerNodes, ok := m.observerNodes[cfg.Observer]
		if !ok {
			m.logger.Warn("observer nodes not found for streaming",
				logger.String("streaming", streamingName),
				logger.String("observer", cfg.Observer))
			continue
		}
		shardAssignment, ok := m.shardAssignments[cfg.Database]
		if !ok {
			m.logger.Warn("shard assignment not found for streaming",
				logger.String("streaming", streamingName),
				logger.String("database", cfg.Database))
			continue
		}

		shardIDs := lo.Keys(shardAssignment.Shards)
		nodeID := models.NodeID(0)
		consumers := lo.SliceToMap(observerNodes, func(node models.StatelessNode) (models.NodeID, models.StatelessNode) {
			nodeID++
			return nodeID, node
		})

		consumerAssignments := RangeConsumeAssign(lo.Keys(consumers), shardIDs)
		state := &models.StreamingState{
			Config:             *cfg,
			Consumers:          consumers,
			ConsumeAssignments: consumerAssignments,
		}

		oldState := m.streamingStates[streamingName]
		if !cmp.Equal(oldState, state) {
			data := encoding.JSONMarshal(state)
			if err := m.masterRepo.Put(m.ctx, constants.GetStreamingStatePath(streamingName), data); err != nil {
				m.logger.Error("update streaming state error", logger.Error(err))
				continue
			}
			m.streamingStates[streamingName] = state

			m.logger.Info("streaming consumer reassigned",
				logger.String("streaming", streamingName),
				logger.String("oldState", string(encoding.JSONMarshal(oldState))),
				logger.String("newState", string(data)))
		} else {
			m.logger.Info("streaming consumer assignment not changed")
		}

	}
}
