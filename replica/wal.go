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
	"fmt"
	"path"
	"strconv"
	"sync"
	"time"

	"go.uber.org/atomic"

	"github.com/lindb/lindb/config"
	"github.com/lindb/lindb/coordinator/storage"
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/fileutil"
	"github.com/lindb/lindb/pkg/logger"
	"github.com/lindb/lindb/pkg/queue"
	"github.com/lindb/lindb/pkg/timeutil"
	"github.com/lindb/lindb/rpc"
	"github.com/lindb/lindb/tsdb"
)

//go:generate mockgen -source=./wal.go -destination=./wal_mock.go -package=replica

// for testing
var (
	newFanOutQueue   = queue.NewFanOutQueue
	newWriteAheadLog = NewWriteAheadLog
)

type partitionKey struct {
	shardID    models.ShardID
	familyTime int64
	leader     models.NodeID
}

// WriteAheadLogManager represents manage all write ahead log.
type WriteAheadLogManager interface {
	// GetOrCreateLog returns write ahead log for database,
	// if exist returns it, else creates a new log.
	GetOrCreateLog(database string) WriteAheadLog
	// GetReplicaState returns replica state for given database's name.
	GetReplicaState(database string) []models.FamilyLogReplicaState
	// Recovery recoveries local history wal when server start.
	Recovery() error
}

// WriteAheadLog represents write ahead log underlying fan out queue.
type WriteAheadLog interface {
	// GetOrCreatePartition returns a partition of write ahead log.
	// if exist returns it, else create a new partition.
	GetOrCreatePartition(shardID models.ShardID, familyTime int64, leader models.NodeID) (Partition, error)
	getReplicaState() (rs []models.FamilyLogReplicaState)
	// recovery recoveries database write ahead log from local storage.
	recovery() error

	destroy()
}

// writeAheadLogManager implements WriteAheadLogManager.
type (
	databaseLogs         map[string]WriteAheadLog
	writeAheadLogManager struct {
		ctx           context.Context
		cfg           config.WAL
		currentNodeID models.NodeID
		engine        tsdb.Engine
		cliFct        rpc.ClientStreamFactory
		stateMgr      storage.StateManager
		// COW
		databaseLogs atomic.Value
		mutex        sync.Mutex
	}
)

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
		stateMgr:      stateMgr,
	}
	mgr.databaseLogs.Store(make(databaseLogs))

	mgr.garbageCollectTask()

	return mgr
}

func (w *writeAheadLogManager) garbageCollect() {
	logs := w.databaseLogs.Load().(databaseLogs)
	for _, log := range logs {
		log.destroy()
	}
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

func (w *writeAheadLogManager) getLog(database string) (WriteAheadLog, bool) {
	log, ok := w.databaseLogs.Load().(databaseLogs)[database]
	return log, ok
}

//TODO need remove log when database delete
func (w *writeAheadLogManager) insertLog(database string, newLog WriteAheadLog) {
	oldMap := w.databaseLogs.Load().(databaseLogs)
	newMap := make(databaseLogs)
	for database, log := range oldMap {
		newMap[database] = log
	}
	newMap[database] = newLog
	w.databaseLogs.Store(newMap)
}

// GetOrCreateLog returns write ahead log for database,
// if exist returns it, else creates a new wal
func (w *writeAheadLogManager) GetOrCreateLog(database string) WriteAheadLog {
	log, ok := w.getLog(database)
	if ok {
		return log
	}
	w.mutex.Lock()
	defer w.mutex.Unlock()

	if log, ok = w.getLog(database); ok {
		return log
	}

	log = newWriteAheadLog(w.ctx, w.cfg, w.currentNodeID, database, w.engine, w.cliFct, w.stateMgr)
	w.insertLog(database, log)
	return log
}

// Recovery recoveries local history wal when server start.
func (w *writeAheadLogManager) Recovery() error {
	if !fileutil.Exist(w.cfg.Dir) {
		return nil
	}
	databaseNames, err := fileutil.ListDir(w.cfg.Dir)
	if err != nil {
		return err
	}
	for _, databaseName := range databaseNames {
		log := w.GetOrCreateLog(databaseName)
		//
		if err := log.recovery(); err != nil {
			return err
		}
	}
	return nil
}

func (w *writeAheadLogManager) GetReplicaState(database string) []models.FamilyLogReplicaState {
	log, ok := w.getLog(database)
	if ok {
		return log.getReplicaState()
	}
	return nil
}

type (
	// family log = shard + family + leader
	familyLogs map[partitionKey]Partition
	// writeAheadLog implements WriteAheadLog.
	writeAheadLog struct {
		ctx           context.Context
		database      string
		dir           string
		cfg           config.WAL
		currentNodeID models.NodeID
		engine        tsdb.Engine
		cliFct        rpc.ClientStreamFactory
		stateMgr      storage.StateManager

		mutex      sync.Mutex
		familyLogs atomic.Value

		logger *logger.Logger
	}
)

// NewWriteAheadLog creates a WriteAheadLog instance.
func NewWriteAheadLog(
	ctx context.Context,
	cfg config.WAL,
	currentNodeID models.NodeID,
	database string,
	engine tsdb.Engine,
	cliFct rpc.ClientStreamFactory,
	stateMgr storage.StateManager,
) WriteAheadLog {
	log := &writeAheadLog{
		ctx:           ctx,
		currentNodeID: currentNodeID,
		database:      database,
		dir:           path.Join(cfg.Dir, database),
		cfg:           cfg,
		engine:        engine,
		cliFct:        cliFct,
		stateMgr:      stateMgr,
		logger:        logger.GetLogger("replica", "WriteAheadLogManager"),
	}
	log.familyLogs.Store(make(familyLogs))
	return log
}

// GetOrCreatePartition returns a partition of write ahead log.
// if exist returns it, else create a new partition.
func (w *writeAheadLog) GetOrCreatePartition(
	shardID models.ShardID,
	familyTime int64,
	leader models.NodeID,
) (Partition, error) {
	key := partitionKey{
		shardID:    shardID,
		familyTime: familyTime,
		leader:     leader,
	}
	p, ok := w.getPartition(key)
	if ok {
		return p, nil
	}
	w.mutex.Lock()
	defer w.mutex.Unlock()

	// double check
	p, ok = w.getPartition(key)
	if ok {
		return p, nil
	}
	shard, ok := w.engine.GetShard(w.database, shardID)
	if !ok {
		return nil, fmt.Errorf("shard: %d not exist", shardID.Int())
	}
	family, err := shard.GetOrCrateDataFamily(familyTime)
	if err != nil {
		return nil, err
	}
	// wal path: base dir + database + shard + family time + leader
	dir := path.Join(
		strconv.Itoa(int(shardID)),
		timeutil.FormatTimestamp(familyTime, timeutil.DataTimeFormat4),
		strconv.Itoa(int(leader)))
	dirPath := path.Join(w.dir, dir)

	interval := w.cfg.RemoveTaskInterval.Duration()

	q, err := newFanOutQueue(dirPath, w.cfg.GetDataSizeLimit(), interval)
	if err != nil {
		return nil, err
	}
	p = NewPartition(w.ctx, shard, family, w.currentNodeID, q, w.cliFct, w.stateMgr)

	w.insertPartition(key, p)
	return p, nil
}

func (w *writeAheadLog) getReplicaState() (rs []models.FamilyLogReplicaState) {
	logs := w.familyLogs.Load().(familyLogs)
	for k, v := range logs {
		state := v.getReplicaState()
		state.Leader = k.leader
		rs = append(rs, state)
	}
	return
}

func (w *writeAheadLog) getPartition(key partitionKey) (Partition, bool) {
	p, ok := w.familyLogs.Load().(familyLogs)[key]
	return p, ok
}

func (w *writeAheadLog) insertPartition(key partitionKey, p Partition) {
	oldMap := w.familyLogs.Load().(familyLogs)
	newMap := make(familyLogs)
	for key, partition := range oldMap {
		newMap[key] = partition
	}
	newMap[key] = p
	w.familyLogs.Store(newMap)
}

func (w *writeAheadLog) recovery() error {
	shards, err := fileutil.ListDir(w.dir)
	if err != nil {
		return err
	}
	for _, shard := range shards {
		families, err := fileutil.ListDir(path.Join(w.dir, shard))
		if err != nil {
			return err
		}

		shardID := models.ParseShardID(shard)
		for _, family := range families {
			leaders, err := fileutil.ListDir(path.Join(w.dir, shard, family))
			if err != nil {
				return err
			}
			familyTime, _ := timeutil.ParseTimestamp(family, timeutil.DataTimeFormat4)
			for _, leader := range leaders {
				leaderID := models.ParseNodeID(leader)
				partition, err := w.GetOrCreatePartition(shardID, familyTime, leaderID)
				if err != nil {
					return err
				}
				if err = partition.recovery(leaderID); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func (w *writeAheadLog) destroy() {
	logs := w.familyLogs.Load().(familyLogs)
	newLogs := make(familyLogs)
	expireLogs := make(familyLogs)

	for key, log := range logs {
		isExpire := log.IsExpire()
		w.logger.Debug("check write ahead log if expire", logger.String("path",
			log.Path()), logger.Any("expire", isExpire))
		if isExpire {
			expireLogs[key] = log
		} else {
			newLogs[key] = log
		}
	}

	// set new logs
	w.familyLogs.Store(newLogs)

	for _, log := range expireLogs {
		w.logger.Info("write ahead log is expire, need destroy it", logger.String("path", log.Path()))
		if err := log.Close(); err != nil {
			w.logger.Warn("close write ahead log", logger.String("path", log.Path()), logger.Error(err))
		}
		if err := fileutil.RemoveDir(log.Path()); err != nil {
			w.logger.Warn("remove write ahead log dir", logger.String("path", log.Path()), logger.Error(err))
		}
	}
	shards, err := fileutil.ListDir(w.dir)
	if err != nil {
		w.logger.Warn("list shard dir err")
	}
	for _, shard := range shards {
		families, err := fileutil.ListDir(path.Join(w.dir, shard))
		if err != nil {
			w.logger.Warn("list family dir err")
			continue
		}
		for _, family := range families {
			leaders, err := fileutil.ListDir(path.Join(w.dir, shard, family))
			if err != nil {
				continue
			}
			if len(leaders) == 0 {
				walPath := path.Join(w.dir, shard, family)
				if err := fileutil.RemoveDir(walPath); err != nil {
					w.logger.Warn("remove write ahead log err", logger.String("path", walPath))
				}
			}
		}

	}
}
