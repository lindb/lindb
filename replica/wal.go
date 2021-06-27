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
	"errors"
	"path"
	"strconv"
	"sync"

	"github.com/lindb/lindb/config"
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/queue"
	"github.com/lindb/lindb/rpc"
	"github.com/lindb/lindb/service"
)

//go:generate mockgen -source=./wal.go -destination=./wal_mock.go -package=replica

// for testing
var (
	newFanOutQueue   = queue.NewFanOutQueue
	newWriteAheadLog = NewWriteAheadLog
)

// WriteAheadLogManager represents manage all write ahead log.
type WriteAheadLogManager interface {
	// GetOrCreateLog returns write ahead log for database,
	// if exist returns it, else creates a new log.
	GetOrCreateLog(database string) WriteAheadLog
}

// WriteAheadLog represents write ahead log underlying fan out queue.
type WriteAheadLog interface {
	// GetOrCreatePartition returns a partition of write ahead log.
	// if exist returns it, else create a new partition.
	GetOrCreatePartition(shardID models.ShardID) (Partition, error)
}

// writeAheadLogManager implements WriteAheadLogManager.
type writeAheadLogManager struct {
	cfg           config.Replica
	currentNodeID models.NodeID
	databaseLogs  map[string]WriteAheadLog
	storageSrv    service.StorageService
	cliFct        rpc.ClientStreamFactory

	mutex sync.Mutex
}

// NewWriteAheadLogManager creates a WriteAheadLogManager instance.
func NewWriteAheadLogManager(cfg config.Replica,
	currentNodeID models.NodeID,
	storageSrv service.StorageService,
	cliFct rpc.ClientStreamFactory,
) WriteAheadLogManager {
	return &writeAheadLogManager{
		cfg:           cfg,
		currentNodeID: currentNodeID,
		storageSrv:    storageSrv,
		cliFct:        cliFct,

		databaseLogs: make(map[string]WriteAheadLog),
	}
}

// GetOrCreateLog returns write ahead log for database,
// if exist returns it, else creates a new.
func (w *writeAheadLogManager) GetOrCreateLog(database string) WriteAheadLog {
	w.mutex.Lock()
	defer w.mutex.Unlock()

	log, ok := w.databaseLogs[database]
	if ok {
		return log
	}
	// create new, then put in cache
	log = newWriteAheadLog(w.cfg, w.currentNodeID, database, w.storageSrv, w.cliFct)
	w.databaseLogs[database] = log
	return log
}

// writeAheadLog implements WriteAheadLog.
type writeAheadLog struct {
	database      string
	cfg           config.Replica
	currentNodeID models.NodeID
	shardLogs     map[models.ShardID]Partition
	storageSrv    service.StorageService
	cliFct        rpc.ClientStreamFactory

	mutex sync.Mutex
}

// NewWriteAheadLog creates a WriteAheadLog instance.
func NewWriteAheadLog(cfg config.Replica,
	currentNodeID models.NodeID,
	database string,
	storageSrv service.StorageService,
	cliFct rpc.ClientStreamFactory,
) WriteAheadLog {
	return &writeAheadLog{
		currentNodeID: currentNodeID,
		database:      database,
		cfg:           cfg,
		storageSrv:    storageSrv,
		cliFct:        cliFct,
		shardLogs:     make(map[models.ShardID]Partition),
	}
}

// GetOrCreatePartition returns a partition of write ahead log.
// if exist returns it, else create a new partition.
func (w *writeAheadLog) GetOrCreatePartition(shardID models.ShardID) (Partition, error) {
	w.mutex.Lock()
	defer w.mutex.Unlock()

	p, ok := w.shardLogs[shardID]
	if ok {
		return p, nil
	}
	shard, ok := w.storageSrv.GetShard(w.database, int32(shardID))
	if !ok {
		return nil, errors.New("shard not exist")
	}
	dirPath := path.Join(w.cfg.Dir, w.database, strconv.Itoa(int(shardID)))
	interval := w.cfg.RemoveTaskInterval.Duration()

	q, err := newFanOutQueue(dirPath, w.cfg.GetDataSizeLimit(), interval)
	if err != nil {
		return nil, err
	}
	p = NewPartition(shardID, shard, w.currentNodeID, q, w.cliFct)
	w.shardLogs[shardID] = p
	return p, nil
}
