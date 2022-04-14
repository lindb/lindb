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
	"io"
	"path"
	"strconv"
	"sync"

	"github.com/lindb/lindb/config"
	"github.com/lindb/lindb/coordinator/storage"
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/logger"
	"github.com/lindb/lindb/pkg/timeutil"
	"github.com/lindb/lindb/rpc"
	"github.com/lindb/lindb/tsdb"
)

//go:generate mockgen -source=./wal.go -destination=./wal_mock.go -package=replica

// WriteAheadLog represents write ahead log underlying fan out queue.
type WriteAheadLog interface {
	io.Closer

	// Name returns the name of write ahead log.
	Name() string
	// GetOrCreatePartition returns a partition of write ahead log.
	// if exist returns it, else create a new partition.
	GetOrCreatePartition(shardID models.ShardID, familyTime int64, leader models.NodeID) (Partition, error)
	// Stop stops all replicator channels.
	Stop()
	// getReplicaState returns the state of replica.
	getReplicaState() (rs []models.FamilyLogReplicaState)
	// recovery recoveries database write ahead log from local storage.
	recovery() error
	// destroy removes expired write ahead log.
	destroy()
}

// writeAheadLog implements WriteAheadLog.
type writeAheadLog struct {
	ctx           context.Context
	database      string
	dir           string
	cfg           config.WAL
	currentNodeID models.NodeID
	engine        tsdb.Engine
	cliFct        rpc.ClientStreamFactory
	stateMgr      storage.StateManager

	mutex sync.Mutex
	// family log = shard + family + leader
	familyLogs map[partitionKey]Partition

	logger *logger.Logger
}

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
		familyLogs:    make(map[partitionKey]Partition),
		logger:        logger.GetLogger("replica", "WriteAheadLog"),
	}
	return log
}

// Name returns the name of write ahead log.
func (w *writeAheadLog) Name() string {
	return w.database
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
	w.mutex.Lock()
	defer w.mutex.Unlock()

	p, ok := w.familyLogs[key]
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
	p = NewPartitionFn(w.ctx, shard, family, w.currentNodeID, q, w.cliFct, w.stateMgr)

	w.familyLogs[key] = p
	return p, nil
}

// getReplicaState returns the state of replica.
func (w *writeAheadLog) getReplicaState() (rs []models.FamilyLogReplicaState) {
	w.mutex.Lock()
	defer w.mutex.Unlock()

	for k, v := range w.familyLogs {
		state := v.getReplicaState()
		state.Leader = k.leader
		rs = append(rs, state)
	}
	return
}

// recovery recoveries database write ahead log from local storage.
func (w *writeAheadLog) recovery() error {
	shards, err := listDirFn(w.dir)
	if err != nil {
		return err
	}
	for _, shard := range shards {
		families, err := listDirFn(path.Join(w.dir, shard))
		if err != nil {
			return err
		}

		shardID := models.ParseShardID(shard)
		for _, family := range families {
			leaders, err := listDirFn(path.Join(w.dir, shard, family))
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
				err = partition.recovery(leaderID)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

// destroy removes expired write ahead log.
func (w *writeAheadLog) destroy() {
	w.mutex.Lock()

	newLogs := make(map[partitionKey]Partition)
	expireLogs := make(map[partitionKey]Partition)

	for key, log := range w.familyLogs {
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
	w.familyLogs = newLogs
	w.mutex.Unlock()

	for _, log := range expireLogs {
		w.logger.Info("write ahead log is expire, need destroy it", logger.String("path", log.Path()))
		log.Stop()
		if err := log.Close(); err != nil {
			w.logger.Warn("close write ahead log", logger.String("path", log.Path()), logger.Error(err))
		}
		if err := removeDirFn(log.Path()); err != nil {
			w.logger.Warn("remove write ahead log dir", logger.String("path", log.Path()), logger.Error(err))
		}
	}
}

// Close closes all log queues.
func (w *writeAheadLog) Close() error {
	w.mutex.Lock()
	defer w.mutex.Unlock()

	for _, log := range w.familyLogs {
		if err := log.Close(); err != nil {
			w.logger.Warn("close write ahead log err", logger.String("path", log.Path()))
		}
	}
	return nil
}

// Stop stops all replicator channels.
func (w *writeAheadLog) Stop() {
	w.mutex.Lock()
	defer w.mutex.Unlock()

	for _, log := range w.familyLogs {
		log.Stop()
	}
}
