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
	"errors"
	"path"
	"strconv"
	"sync"

	"github.com/lindb/lindb/config"
	"github.com/lindb/lindb/coordinator/storage"
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/queue"
	"github.com/lindb/lindb/rpc"
	"github.com/lindb/lindb/tsdb"
)

//go:generate mockgen -source=./wal.go -destination=./wal_mock.go -package=replica

// for testing
var (
	newFanOutQueue   = queue.NewFanOutQueue
	newWriteAheadLog = NewWriteAheadLog
)

// WriteAheadLogManager represents manage all writeTask ahead log.
type WriteAheadLogManager interface {
	// GetOrCreateLog returns writeTask ahead log for database,
	// if exist returns it, else creates a new log.
	GetOrCreateLog(database string) WriteAheadLog
}

// WriteAheadLog represents writeTask ahead log underlying fan out queue.
type WriteAheadLog interface {
	// GetOrCreatePartition returns a partition of writeTask ahead log.
	// if exist returns it, else create a new partition.
	GetOrCreatePartition(shardID models.ShardID) (Partition, error)
}

// writeAheadLogManager implements WriteAheadLogManager.
type writeAheadLogManager struct {
	ctx           context.Context
	cfg           config.WAL
	currentNodeID models.NodeID
	databaseLogs  map[string]WriteAheadLog
	engine        tsdb.Engine
	cliFct        rpc.ClientStreamFactory
	stateMgr      storage.StateManager

	mutex sync.Mutex
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
	return &writeAheadLogManager{
		ctx:           ctx,
		cfg:           cfg,
		currentNodeID: currentNodeID,
		engine:        engine,
		cliFct:        cliFct,
		stateMgr:      stateMgr,

		databaseLogs: make(map[string]WriteAheadLog),
	}
}

// GetOrCreateLog returns writeTask ahead log for database,
// if exist returns it, else creates a new.
func (w *writeAheadLogManager) GetOrCreateLog(database string) WriteAheadLog {
	w.mutex.Lock()
	defer w.mutex.Unlock()

	log, ok := w.databaseLogs[database]
	if ok {
		return log
	}
	// create new, then put in cache
	log = newWriteAheadLog(w.ctx, w.cfg, w.currentNodeID, database, w.engine, w.cliFct, w.stateMgr)
	w.databaseLogs[database] = log
	return log
}

// writeAheadLog implements WriteAheadLog.
type writeAheadLog struct {
	ctx           context.Context
	database      string
	cfg           config.WAL
	currentNodeID models.NodeID
	shardLogs     map[models.ShardID]Partition
	engine        tsdb.Engine
	cliFct        rpc.ClientStreamFactory
	stateMgr      storage.StateManager

	mutex sync.Mutex
}

// NewWriteAheadLog creates a WriteAheadLog instance.
func NewWriteAheadLog(ctx context.Context,
	cfg config.WAL,
	currentNodeID models.NodeID,
	database string,
	engine tsdb.Engine,
	cliFct rpc.ClientStreamFactory,
	stateMgr storage.StateManager,
) WriteAheadLog {
	return &writeAheadLog{
		ctx:           ctx,
		currentNodeID: currentNodeID,
		database:      database,
		cfg:           cfg,
		engine:        engine,
		cliFct:        cliFct,
		stateMgr:      stateMgr,
		shardLogs:     make(map[models.ShardID]Partition),
	}
}

// GetOrCreatePartition returns a partition of writeTask ahead log.
// if exist returns it, else create a new partition.
func (w *writeAheadLog) GetOrCreatePartition(shardID models.ShardID) (Partition, error) {
	w.mutex.Lock()
	defer w.mutex.Unlock()

	p, ok := w.shardLogs[shardID]
	if ok {
		return p, nil
	}
	shard, ok := w.engine.GetShard(w.database, shardID)
	if !ok {
		return nil, errors.New("shard not exist")
	}
	dirPath := path.Join(w.cfg.Dir, w.database, strconv.Itoa(int(shardID)))
	interval := w.cfg.RemoveTaskInterval.Duration()

	q, err := newFanOutQueue(dirPath, w.cfg.GetDataSizeLimit(), interval)
	if err != nil {
		return nil, err
	}
	p = NewPartition(w.ctx, shardID, shard, w.currentNodeID, q, w.cliFct, w.stateMgr)
	w.shardLogs[shardID] = p
	return p, nil
}
