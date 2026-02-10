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

package write

import (
	"context"
	"errors"
	"reflect"
	"sync"

	"github.com/lindb/arrow/pkg/logs"
	"github.com/lindb/common/pkg/logger"

	logswriter "github.com/lindb/lindb/app/broker/write/logs"
	"github.com/lindb/lindb/app/broker/write/writer"
	"github.com/lindb/lindb/coordinator/discovery"
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/option"
)

type Manager interface {
	discovery.Watcher

	GetWriter(database string) (writer.Writer, bool)
}

type manager struct {
	ctx       context.Context
	writers   map[string]writer.Writer
	databases map[string]writer.Database

	lock   sync.RWMutex
	logger logger.Logger
}

func NewManager(ctx context.Context) Manager {
	return &manager{
		ctx:       ctx,
		writers:   make(map[string]writer.Writer),
		databases: make(map[string]writer.Database),
		logger:    logger.GetLogger("Write", "Manager"),
	}
}

func (m *manager) GetWriter(database string) (writer.Writer, bool) {
	m.lock.RLock()
	defer m.lock.RUnlock()

	db, ok := m.writers[database]
	return db, ok
}

// OnEvent implements [Manager].
func (m *manager) OnEvent(event discovery.MetaEvent) {
	switch e := event.(type) {
	case *models.ChangeShardStateEvent:
		if err := m.createShard(e); err != nil {
			m.logger.Error("create shards failed", logger.String("db", e.DatabaseCfg.Name), logger.Error(err))
		}

		// notify shard state change
		m.notifyLeaderChange(e)
	default:
		m.logger.Warn("unsupported event type", logger.Any("eventType", reflect.TypeOf(event)))
	}
}

func (m *manager) notifyLeaderChange(e *models.ChangeShardStateEvent) {
	m.lock.Lock()
	defer m.lock.Unlock()

	database, ok := m.databases[e.DatabaseCfg.Name]
	if !ok {
		return
	}
	database.LeaderChanged(e.Shards, e.LiveNodes)
}

func (m *manager) createShard(e *models.ChangeShardStateEvent) error {
	if len(e.Shards) == 0 {
		return errors.New("numOfShards must be greater than zero")
	}

	m.lock.Lock()
	defer m.lock.Unlock()

	cfg := e.DatabaseCfg
	database, ok := m.databases[cfg.Name]
	if !ok {
		switch cfg.Option.Engine {
		case option.Log:
			databaseAccessor := NewDatabase[*logs.Log](m.ctx, cfg)
			m.writers[cfg.Name] = logswriter.NewWriter(m.ctx, cfg, databaseAccessor)
			database = databaseAccessor
		default:
			return errors.New("unsupported engine type: " + string(cfg.Option.Engine))
		}
		m.databases[cfg.Name] = database
	}
	database.CreateShards(e.Shards, e.LiveNodes)
	return nil
}

// Subscribe implements [Manager].
func (m *manager) Subscribe(sub discovery.Subscriber) {
	panic("unimplemented")
}

// Unsubscribe implements [Manager].
func (m *manager) Unsubscribe(sub discovery.Subscriber) {
	panic("unimplemented")
}

// Close implements [Manager].
func (m *manager) Close() {
	panic("unimplemented")
}
