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

package replica

import (
	"context"
	"io"
	"sync"
	"time"

	"github.com/lindb/lindb/config"
	"github.com/lindb/lindb/coordinator/storage"
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/fileutil"
	"github.com/lindb/lindb/pkg/logger"
	"github.com/lindb/lindb/pkg/queue"
	"github.com/lindb/lindb/rpc"
	"github.com/lindb/lindb/tsdb"
)

//go:generate mockgen -source=./wal_manager.go -destination=./wal_manager_mock.go -package=replica

// for testing
var (
	newFanOutQueue   = queue.NewFanOutQueue
	NewPartitionFn   = NewPartition
	newWriteAheadLog = NewWriteAheadLog
	fileExistFn      = fileutil.Exist
	listDirFn        = fileutil.ListDir
	removeDirFn      = fileutil.RemoveDir
)

// partitionKey represents partition unique key.
type partitionKey struct {
	shardID    models.ShardID
	familyTime int64
	leader     models.NodeID
}

// WriteAheadLogManager represents manage all write ahead log.
type WriteAheadLogManager interface {
	io.Closer

	// GetOrCreateLog returns write ahead log for database,
	// if exist returns it, else creates a new log.
	GetOrCreateLog(database string) WriteAheadLog
	// GetReplicaState returns replica state for given database's name.
	GetReplicaState(database string) []models.FamilyLogReplicaState
	// DropDatabases drops write ahead log of databases, keep active databases.
	DropDatabases(activeDatabases map[string]struct{})
	// StopDatabases stop the replicator for write ahead log of databases, keep active databases.
	StopDatabases(activeDatabases map[string]struct{})
	// Recovery recoveries local history wal when server start.
	Recovery() error
	// Stop stops all replicator channel.
	Stop()
}

// writeAheadLogManager implements WriteAheadLogManager.
type writeAheadLogManager struct {
	ctx           context.Context
	cfg           config.WAL
	currentNodeID models.NodeID
	engine        tsdb.Engine
	cliFct        rpc.ClientStreamFactory
	stateMgr      storage.StateManager

	databaseLogs map[string]WriteAheadLog

	mutex  sync.Mutex
	logger *logger.Logger
}

// NewWriteAheadLogManager creates a WriteAheadLogManager instance.
func NewWriteAheadLogManager(
	ctx context.Context,
	cfg config.WAL,
	currentNodeID models.NodeID,
	engine tsdb.Engine,
	cliFct rpc.ClientStreamFactory,
	stateMgr storage.StateManager,
) WriteAheadLogManager {
	mgr := &writeAheadLogManager{
		ctx:           ctx,
		cfg:           cfg,
		currentNodeID: currentNodeID,
		engine:        engine,
		cliFct:        cliFct,
		databaseLogs:  make(map[string]WriteAheadLog),
		stateMgr:      stateMgr,
		logger:        logger.GetLogger("Replica", "WriteAheadLogManager"),
	}

	mgr.garbageCollectTask()

	return mgr
}

func (w *writeAheadLogManager) garbageCollect() {
	var logs []WriteAheadLog
	w.mutex.Lock()
	for _, log := range w.databaseLogs {
		logs = append(logs, log)
	}
	w.mutex.Unlock()

	for _, log := range logs {
		log.destroy()
	}
}

func (w *writeAheadLogManager) getDatabaseLogs() []WriteAheadLog {
	var logs []WriteAheadLog
	w.mutex.Lock()
	defer w.mutex.Unlock()

	for _, log := range w.databaseLogs {
		logs = append(logs, log)
	}
	return logs
}

func (w *writeAheadLogManager) garbageCollectTask() {
	go func() {
		ticker := time.NewTicker(w.cfg.RemoveTaskInterval.Duration())
		for {
			select {
			case <-ticker.C:
				w.garbageCollect()
			case <-w.ctx.Done():
				return
			}
		}
	}()
}

// GetOrCreateLog returns write ahead log for database,
// if exist returns it, else creates a new wal
func (w *writeAheadLogManager) GetOrCreateLog(database string) WriteAheadLog {
	w.mutex.Lock()
	defer w.mutex.Unlock()

	log, ok := w.databaseLogs[database]
	if ok {
		return log
	}

	log = newWriteAheadLog(w.ctx, w.cfg, w.currentNodeID, database, w.engine, w.cliFct, w.stateMgr)
	w.databaseLogs[database] = log
	return log
}

// Recovery recoveries local history wal when server start.
func (w *writeAheadLogManager) Recovery() error {
	if !fileExistFn(w.cfg.Dir) {
		return nil
	}
	databaseNames, err := listDirFn(w.cfg.Dir)
	if err != nil {
		return err
	}
	for _, databaseName := range databaseNames {
		log := w.GetOrCreateLog(databaseName)
		if err := log.recovery(); err != nil {
			return err
		}
	}
	return nil
}

// GetReplicaState returns replica state for given database's name.
func (w *writeAheadLogManager) GetReplicaState(database string) []models.FamilyLogReplicaState {
	w.mutex.Lock()
	defer w.mutex.Unlock()

	log, ok := w.databaseLogs[database]
	if ok {
		return log.getReplicaState()
	}
	return nil
}

// dropDatabase drops write ahead log.
func (w *writeAheadLogManager) dropDatabase(log WriteAheadLog) {
	if err := log.Close(); err != nil {
		w.logger.Error("close write ahead log err",
			logger.String("database", log.Name()), logger.Error(err))
		return
	}
	if err := log.Drop(); err != nil {
		w.logger.Warn("remove write ahead log dir failure", logger.String("path", log.Name()), logger.Error(err))
		return
	}

	// remove memory data
	w.mutex.Lock()
	delete(w.databaseLogs, log.Name())
	w.mutex.Unlock()
}

// DropDatabases drops write ahead log of databases, keep active databases.
func (w *writeAheadLogManager) DropDatabases(activeDatabases map[string]struct{}) {
	logs := w.getDatabaseLogs()
	for _, db := range logs {
		_, ok := activeDatabases[db.Name()]
		if ok {
			continue
		}
		w.dropDatabase(db)
		w.logger.Info("drop write ahead log successfully", logger.String("database", db.Name()))
	}
}

// StopDatabases stop the replicator for write ahead log of databases, keep active databases.
func (w *writeAheadLogManager) StopDatabases(activeDatabases map[string]struct{}) {
	logs := w.getDatabaseLogs()
	for _, db := range logs {
		_, ok := activeDatabases[db.Name()]
		if ok {
			continue
		}
		db.Stop()
		w.logger.Info("stop write ahead log replica successfully", logger.String("database", db.Name()))
	}
}

// Close closes all log queues.
func (w *writeAheadLogManager) Close() error {
	logs := w.getDatabaseLogs()
	for _, db := range logs {
		if err := db.Close(); err != nil {
			w.logger.Error("close write ahead log err",
				logger.String("database", db.Name()), logger.Error(err))
		}
	}
	return nil
}

// Stop stops all replicator channel.
func (w *writeAheadLogManager) Stop() {
	logs := w.getDatabaseLogs()
	for _, db := range logs {
		db.Stop()
	}
}
