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

package streaming

import (
	"context"
	"path/filepath"
	"strings"
	"sync"

	"github.com/lindb/common/pkg/encoding"
	"github.com/lindb/common/pkg/logger"

	"github.com/lindb/lindb/constants"
	"github.com/lindb/lindb/coordinator/discovery"
	"github.com/lindb/lindb/internal/linmetric"
	"github.com/lindb/lindb/metrics"
	"github.com/lindb/lindb/models"
)

type StateManager interface {
	discovery.StateMachineEventHandle

	// RegisterWatcher registers state manager watcher.
	RegisterWatcher(watcher discovery.Watcher)
}

type stateManager struct {
	discovery.StateMachineEventHandle

	ctx    context.Context
	cancel context.CancelFunc

	events   chan *discovery.Event
	watchers []discovery.Watcher

	mutex sync.RWMutex

	statistics *metrics.StateManagerStatistics
	logger     logger.Logger
}

func NewStateManager(ctx context.Context) StateManager {
	c, cancel := context.WithCancel(ctx)
	mgr := &stateManager{
		ctx:    c,
		cancel: cancel,
		events: make(chan *discovery.Event),

		statistics: metrics.NewStateManagerStatistics(linmetric.StreamingRegistry),
		logger:     logger.GetLogger("Streaming", "StateManager"),
	}

	// start discovery event consumer task
	go mgr.consumeEvents()

	return mgr
}

func (s *stateManager) RegisterWatcher(watcher discovery.Watcher) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.watchers = append(s.watchers, watcher)
}

// EmitEvent emits the discovery event to state manager.
func (s *stateManager) EmitEvent(event *discovery.Event) {
	s.events <- event
}

// consumeEvents consumes discovery events, processes events based on event type.
func (s *stateManager) consumeEvents() {
	for {
		select {
		case <-s.ctx.Done():
			s.logger.Info("state manager event consumer exiting")
			return
		case event := <-s.events:
			s.processEvent(event)
		}
	}
}

// processEvent processes the discovery event.
func (s *stateManager) processEvent(event *discovery.Event) {
	eventType := event.Type.String()
	defer func() {
		if err := recover(); err != nil {
			s.statistics.Panics.WithTagValues(eventType, constants.BrokerRole).Incr()
			s.logger.Error("panic when process discovery event, lost the state",
				logger.Any("err", err), logger.Stack())
		}
	}()

	s.mutex.Lock()
	defer s.mutex.Unlock()

	var err error
	switch event.Type {
	case discovery.DatabaseConfigChanged:
		err = s.onDatabaseCfgChange(event.Key, event.Value)
	case discovery.DatabaseConfigDeletion:
		err = s.onDatabaseCfgDelete(event.Key)
	case discovery.StreamingConfigChanged:
		err = s.onStreamingCfgChange(event.Key, event.Value)
	case discovery.StreamingConfigDeletion:
		err = s.onStreamingCfgDelete(event.Key)
	case discovery.StreamingJobChanged:
		err = s.onStreamingJobCfgChange(event.Key, event.Value)
	case discovery.StreamingJobDeletion:
		err = s.onStreamingJobCfgDelete(event.Key)
	default:
		s.logger.Warn("unknown event type", logger.String("type", event.Type.String()))
	}

	if err != nil {
		s.statistics.HandleEventFailure.WithTagValues(eventType, constants.StreamingRole).Incr()
	} else {
		s.statistics.HandleEvents.WithTagValues(eventType, constants.StreamingRole).Incr()
	}
}

// onDatabaseCfgChange triggers when database create/modify.
func (s *stateManager) onDatabaseCfgChange(key string, data []byte) error {
	s.logger.Info("database config is modified",
		logger.String("key", key),
		logger.String("data", string(data)))

	cfg := models.Database{}
	if err := encoding.JSONUnmarshal(data, &cfg); err != nil {
		s.logger.Error("database config modified but unmarshal error", logger.Error(err))
		return err
	}

	if cfg.Name == "" {
		s.logger.Error("database name cannot be empty")
		return constants.ErrNameEmpty
	}

	for _, watcher := range s.watchers {
		watcher.OnEvent(&cfg)
	}

	return nil
}

// onDatabaseCfgDelete triggers when database delete.
func (s *stateManager) onDatabaseCfgDelete(key string) error {
	s.logger.Info("database config is deleted",
		logger.String("key", key))

	_, dbName := filepath.Split(key)
	event := &models.DeleteDatabase{
		Database: dbName,
	}

	for _, watcher := range s.watchers {
		watcher.OnEvent(event)
	}

	return nil
}

// onStreamingCfgDelete triggers when streaming config delete.
func (s *stateManager) onStreamingCfgDelete(key string) error {
	s.logger.Info("streaming config is deleted",
		logger.String("key", key))

	name := strings.TrimPrefix(key, constants.GetStreamingConfigPath(""))

	event := &models.DeleteStreaming{
		Streaming: name,
	}

	for _, watcher := range s.watchers {
		watcher.OnEvent(event)
	}

	return nil
}

// onStreamingJobCfgDelete triggers when streaming job delete.
func (s *stateManager) onStreamingJobCfgDelete(key string) error {
	s.logger.Info("streaming job config is deleted",
		logger.String("key", key))

	stream, job, err := constants.ParseStreamingJob(key)
	if err != nil {
		s.logger.Error("parse streaming job key error", logger.String("key", key), logger.Error(err))
		return err
	}

	event := &models.DeleteStreamingJob{
		Streaming: stream,
		JobName:   job,
	}

	for _, watcher := range s.watchers {
		watcher.OnEvent(event)
	}

	return nil
}

// onStreamingCfgChange triggers when streaming config create/modify.
func (s *stateManager) onStreamingCfgChange(key string, data []byte) error {
	s.logger.Info("streaming config is modified",
		logger.String("key", key),
		logger.String("data", string(data)))

	cfg := models.Streaming{}
	if err := encoding.JSONUnmarshal(data, &cfg); err != nil {
		s.logger.Error("streaming config modified but unmarshal error", logger.Error(err))
		return err
	}

	if cfg.Name == "" {
		s.logger.Error("streaming config name cannot be empty")
		return constants.ErrNameEmpty
	}

	for _, watcher := range s.watchers {
		watcher.OnEvent(&cfg)
	}

	return nil
}

// onStreamingJobCfgChange triggers when streaming job create/modify.
func (s *stateManager) onStreamingJobCfgChange(key string, data []byte) error {
	s.logger.Info("streaming job config is modified",
		logger.String("key", key),
		logger.String("data", string(data)))

	stream, job, err := constants.ParseStreamingJob(key)
	if err != nil {
		s.logger.Error("parse streaming job key error", logger.String("key", key), logger.Error(err))
		return err
	}

	event := &models.ModifyStreamingJob{
		Streaming: stream,
		JobName:   job,
		Script:    string(data),
	}

	for _, watcher := range s.watchers {
		watcher.OnEvent(event)
	}

	return nil
}

// Close implements [StateManager].
// Subtle: this method shadows the method (StateMachineEventHandle).Close of stateManager.StateMachineEventHandle.
func (s *stateManager) Close() {
	s.cancel()
}
