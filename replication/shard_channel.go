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

package replication

import (
	"context"
	"path"
	"strconv"
	"sync"
	"time"

	"github.com/lindb/lindb/config"
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/logger"
	"github.com/lindb/lindb/pkg/queue"
	"github.com/lindb/lindb/rpc"
	"github.com/lindb/lindb/rpc/proto/field"
)

//go:generate mockgen -source=./shard_channel.go -destination=./shard_channel_mock.go -package=replication

// for testing
var (
	newFanOutQueue = queue.NewFanOutQueue
)

// Channel represents a place to buffer the data for a specific cluster, database, shardID.
type Channel interface {
	// Database returns the database attribution.
	Database() string
	// ShardID returns the shardID attribution.
	ShardID() int32
	// Startup starts the channel internal goroutine worker which consumes chan data and writes wal
	Startup()
	// Write writes the data into the channel, ErrCanceled is returned when the channel is canceled before
	// data is wrote successfully.
	// Concurrent safe.
	Write(metric *field.Metric) error
	// GetOrCreateReplicator get a existed or creates a new replicator for target.
	// Concurrent safe.
	GetOrCreateReplicator(target models.Node) (Replicator, error)
	// Nodes returns all the target nodes for replication.
	Targets() []models.Node
}

// channel implements Channel.
type channel struct {
	// context to close channel
	ctx     context.Context
	dirPath string
	// factory to get WriteClient
	fct      rpc.ClientStreamFactory
	database string
	shardID  int32
	// underlying storage for written data
	q queue.FanOutQueue
	// chanel to convert multiple goroutine write to single goroutine write to FanOutQueue
	ch chan []byte

	chunk Chunk // buffer current write metric for compress

	// last flush time
	lastFlushTime time.Time
	// interval for check flush
	checkFlushInterval time.Duration
	// interval for flush
	flushInterval time.Duration
	//buffer size limit for batch bytes before append to queue
	bufferSizeLimit int

	// target -> replicator map
	replicatorMap sync.Map
	// lock to protect replicatorMap
	lock4map   sync.RWMutex
	lock4write sync.Mutex

	logger *logger.Logger
}

// newChannel returns a new channel with specific attribution.
func newChannel(
	cxt context.Context,
	cfg config.ReplicationChannel,
	database string,
	shardID int32,
	fct rpc.ClientStreamFactory,
) (Channel, error) {
	dirPath := path.Join(cfg.Dir, database, strconv.Itoa(int(shardID)))
	interval := cfg.RemoveTaskInterval.Duration()

	q, err := newFanOutQueue(dirPath, cfg.GetDataSizeLimit(), interval)
	if err != nil {
		return nil, err
	}
	bufferSize := defaultBufferSize
	if cfg.BufferSize > 0 {
		bufferSize = cfg.BufferSize
	}

	c := &channel{
		ctx:                cxt,
		dirPath:            dirPath,
		fct:                fct,
		database:           database,
		shardID:            shardID,
		q:                  q,
		ch:                 make(chan []byte, 2),
		chunk:              newChunk(bufferSize),
		lastFlushTime:      time.Now(),
		checkFlushInterval: cfg.CheckFlushInterval.Duration(),
		flushInterval:      cfg.FlushInterval.Duration(),
		bufferSizeLimit:    cfg.BufferSizeInBytes(),
		logger:             logger.GetLogger("replication", "Channel"),
	}

	return c, nil
}

// Database returns the database attribution.
func (c *channel) Database() string {
	return c.database
}

// ShardID returns the shardID attribution.
func (c *channel) ShardID() int32 {
	return c.shardID
}

// Startup starts the channel internal goroutine worker which consumes chan data and writes wal
func (c *channel) Startup() {
	c.initAppendTask()
	c.watchClose()
}

// GetOrCreateReplicator get a existed or creates a new replicator for target.
// Concurrent safe.
func (c *channel) GetOrCreateReplicator(target models.Node) (Replicator, error) {
	val, ok := c.replicatorMap.Load(target)
	if !ok {
		// double check
		c.lock4map.Lock()
		defer c.lock4map.Unlock()
		val, ok = c.replicatorMap.Load(target)
		if !ok {
			fo, err := c.q.GetOrCreateFanOut(target.Indicator())
			if err != nil {
				return nil, err
			}
			rep := newReplicator(target, c.database, c.shardID, fo, c.fct)

			c.replicatorMap.Store(target, rep)
			return rep, nil
		}
	}
	rep := val.(Replicator)
	return rep, nil
}

// Nodes returns all the nodes for replication.
func (c *channel) Targets() []models.Node {
	nodes := make([]models.Node, 0)
	c.replicatorMap.Range(func(key, value interface{}) bool {
		nd, _ := key.(models.Node)
		nodes = append(nodes, nd)
		return true
	})
	return nodes
}

// Write writes the data into the channel, ErrCanceled is returned when the ctx is canceled before
// data is wrote successfully.
// Concurrent safe.
func (c *channel) Write(metric *field.Metric) error {
	c.lock4write.Lock()
	defer c.lock4write.Unlock()

	c.chunk.Append(metric)

	if c.chunk.IsFull() {
		data, err := c.chunk.MarshalBinary()
		if err != nil {
			return err
		}
		if len(data) == 0 {
			return nil
		}
		select {
		case c.ch <- data:
			return nil
		case <-c.ctx.Done():
			return ErrCanceled
		}
	}
	return nil
}

// initAppendTask starts a goroutine to consume data from ch and batch append to q.
func (c *channel) initAppendTask() {
	go func() {
		// handle write wal loop
		c.writeWAL()
		c.writePendingBeforeClose()
		c.logger.Info("close channel append routine", logger.String("database", c.Database()), logger.Int32("shardID", c.ShardID()))
	}()
}

func (c *channel) writePendingBeforeClose() {
	var wait sync.WaitGroup
	wait.Add(1)
	go func() {
		// try to drain data from chan
		for data := range c.ch {
			err := c.q.Put(data)
			if err != nil {
				c.logger.Error("append to queue err", logger.Error(err))
			}
		}
		wait.Done()
	}()
	close(c.ch)
	wait.Wait()
	// flush chunk pending data if chunk not empty
	if !c.chunk.IsEmpty() {
		// flush chunk pending data if chunk not empty
		c.flushChunk()
	}
}

func (c *channel) checkFlush() {
	now := time.Now()
	if now.After(c.lastFlushTime.Add(c.flushInterval)) {
		c.lock4write.Lock()
		defer c.lock4write.Unlock()

		if !c.chunk.IsEmpty() {
			c.flushChunk()
		}
		c.lastFlushTime = now
	}
}

// flush flushes the chunk data and appends data into queue
func (c *channel) flushChunk() {
	data, err := c.chunk.MarshalBinary()
	if err != nil {
		c.logger.Error("chunk marshal err", logger.Error(err))
		return
	}
	if len(data) == 0 {
		return
	}
	err = c.q.Put(data)
	if err != nil {
		c.logger.Error("append to queue err", logger.Error(err))
	}
}

// writeWAL consumes data from chan, then appends the data into queue
func (c *channel) writeWAL() {
	// on avg 2 * limit could avoid buffer grow
	ticker := time.NewTicker(c.checkFlushInterval)
	defer ticker.Stop()
	for {
		select {
		case <-c.ctx.Done():
			return
		case data := <-c.ch:
			err := c.q.Put(data)
			if err != nil {
				c.logger.Error("append to queue err", logger.Error(err))
			}
		case <-ticker.C:
			// check
			c.checkFlush()
		}
	}
}

// watchClose waits on the context done then close the ch.
func (c *channel) watchClose() {
	go func() {
		<-c.ctx.Done()
		c.lock4map.RLock()
		defer c.lock4map.RUnlock()
		c.replicatorMap.Range(func(key, value interface{}) bool {
			rep, _ := value.(Replicator)
			rep.Stop()
			return true
		})
	}()
}
