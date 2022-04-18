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
	"context"
	"path"
	"time"

	"github.com/lindb/lindb/config"
	"github.com/lindb/lindb/constants"
	"github.com/lindb/lindb/pkg/logger"
	"github.com/lindb/lindb/pkg/state"
	"github.com/lindb/lindb/replica"
	"github.com/lindb/lindb/tsdb"
)

//go:generate mockgen -source=./database_lifecycle.go -destination=./database_lifecycle_mock.go -package=storage

// DatabaseLifecycle represents database's lifecycle manager include data and write ahead log.
type DatabaseLifecycle interface {
	// Startup startups database's lifecycle, includes background task(ttl etc.)
	Startup()
	// Shutdown shutdowns database's lifecycle.
	Shutdown()
}

// databaseLifecycle implements DatabaseLifecycle interface.
type databaseLifecycle struct {
	ctx    context.Context
	cancel context.CancelFunc

	walMgr replica.WriteAheadLogManager
	engine tsdb.Engine
	repo   state.Repository

	logger *logger.Logger
}

// NewDatabaseLifecycle creates a DatabaseLifecycle instance.
func NewDatabaseLifecycle(
	ctx context.Context,
	repo state.Repository,
	walMgr replica.WriteAheadLogManager,
	engine tsdb.Engine,
) DatabaseLifecycle {
	c, cancel := context.WithCancel(ctx)
	return &databaseLifecycle{
		ctx:    c,
		cancel: cancel,
		repo:   repo,
		walMgr: walMgr,
		engine: engine,
		logger: logger.GetLogger("lifecycle", "Database"),
	}
}

// Startup startups database's lifecycle, includes background task(ttl etc.)
func (l *databaseLifecycle) Startup() {
	l.ttlTask()
}

// Shutdown shutdowns database's lifecycle.
func (l *databaseLifecycle) Shutdown() {
	l.cancel()

	if l.walMgr != nil {
		l.logger.Info("stopping write ahead log replicator...")
		l.walMgr.Stop()
		l.logger.Info("stopped write ahead log replicator...")
	}

	// close the storage engine
	if l.engine != nil {
		l.logger.Info("stopping tsdb engine...")
		l.engine.Close()
		l.logger.Info("stopped tsdb engine")
	}

	if l.walMgr != nil {
		l.logger.Info("Closing write ahead log ...")
		if err := l.walMgr.Close(); err != nil {
			l.logger.Error("stopped write ahead log replicator with error", logger.Error(err))
		} else {
			l.logger.Info("write ahead log closed...")
		}
	}
}

// ttlTask runs ttl task in background goroutine.
func (l *databaseLifecycle) ttlTask() {
	go func() {
		ticker := time.NewTicker(config.GlobalStorageConfig().TTLTaskInterval.Duration())
		for {
			select {
			case <-ticker.C:
				l.tryDropDatabases()

				// support dynamic modify config
				ticker.Reset(config.GlobalStorageConfig().TTLTaskInterval.Duration())
			case <-l.ctx.Done():
				return
			}
		}
	}()
}

// tryDropDatabases tries drop database's resource(data/write ahead log), keeps active databases.
func (l *databaseLifecycle) tryDropDatabases() {
	activeDatabases := make(map[string]struct{})
	if err := l.repo.WalkEntry(l.ctx, constants.ShardAssigmentPath, func(key, _ []byte) {
		_, name := path.Split(string(key))
		activeDatabases[name] = struct{}{}
	}); err != nil {
		l.logger.Error("list active database list failure", logger.Error(err))
		return
	}
	if len(activeDatabases) == 0 {
		// if active database is empty, do not drop database operation.
		return
	}
	l.walMgr.StopDatabases(activeDatabases)
	l.engine.DropDatabases(activeDatabases)
	l.walMgr.DropDatabases(activeDatabases)
}
