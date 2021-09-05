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
	"sync"
	"time"

	"go.uber.org/atomic"

	"github.com/lindb/lindb/constants"
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/encoding"
	"github.com/lindb/lindb/pkg/logger"
	protoReplicaV1 "github.com/lindb/lindb/proto/gen/v1/replica"
	"github.com/lindb/lindb/rpc"
)

//go:generate mockgen -source=./sender.go -destination=./sender_mock.go -package=replica

type Sender interface {
	SyncShardState(shardState models.ShardState, liveNodes map[models.NodeID]models.StatefulNode)
	Send(data []byte) error
	Stop()
}

type sender struct {
	ctx      context.Context
	database string
	shardID  models.ShardID
	// factory to get writeTask streamClient
	fct rpc.ClientStreamFactory
	// current WriteStreamClient
	writeCli protoReplicaV1.ReplicaService_WriteClient
	// current WriteServiceClient
	writeService protoReplicaV1.ReplicaServiceClient

	shardState models.ShardState
	liveNodes  map[models.NodeID]models.StatefulNode

	// false -> notReady, true -> ready
	ready   *atomic.Bool
	running *atomic.Bool

	lock4Meta sync.Mutex
	logger    *logger.Logger
}

func newSender(
	ctx context.Context,
	database string,
	shardID models.ShardID,
	fct rpc.ClientStreamFactory,
) Sender {
	return &sender{
		ctx:      ctx,
		fct:      fct,
		database: database,
		shardID:  shardID,
		ready:    atomic.NewBool(false),
		running:  atomic.NewBool(true),
		logger:   logger.GetLogger("replica", "sender"),
	}
}

// recvLoop is a loop to receive message from rpc stream.
// The loop recovers from panic to prevent crash.
// The loop handles rpc re-connection issues.
// The loop only terminates when isStopped() return true.
func (c *sender) recvLoop() {
	defer func() {
		if rec := recover(); rec != nil {
			c.logger.Error("recover from panic, replicator.recvLoop",
				logger.Reflect("recover", rec),
				logger.Stack())

			c.logger.Info("restart recvLoop")
			//TODO modify sleep threshold for retry
			time.Sleep(500 * time.Millisecond)
			go c.recvLoop()
		}
	}()

	for {
		if !c.ready.Load() {
			if err := c.initClient(); err != nil {
				continue
			}
		}

		if !c.running.Load() {
			c.logger.Info("end recvLoop")
			return
		}
		// when connection is stopped, replicator.streamClient.Recv() returns error.
		resp, err := c.writeCli.Recv()
		if err != nil {
			//fixme if seq out of range need reset
			c.logger.Error("recvLoop receive error", logger.Error(err), logger.Stack())
			c.ready.Store(false)
			time.Sleep(time.Second)
			continue
		}
		if resp.Err != "" {
			c.ready.Store(false)
		}
	}
}
func (c *sender) initClient() error {
	// try to re-construct the streaming
	c.lock4Meta.Lock()
	defer c.lock4Meta.Unlock()

	if !c.running.Load() {
		return fmt.Errorf("not running")
	}

	storageNode, ok := c.liveNodes[c.shardState.Leader]
	if !ok {
		c.logger.Error("shard's leader node not live",
			logger.String("database", c.database),
			logger.Any("shard", c.shardID),
			logger.String("leader", (&storageNode).Indicator()))
		return fmt.Errorf("storage node not live")
	}

	writeService, err := c.fct.CreateReplicaServiceClient(&storageNode)
	if err != nil {
		c.logger.Error("recvLoop get service streamClient error", logger.Error(err))
		return err
	}
	c.writeService = writeService

	// pass metadata(database/shard state) when create rpc connection.
	shardState := encoding.JSONMarshal(&c.shardState)
	ctx := rpc.CreateOutgoingContextWithPairs(c.ctx,
		constants.RPCMetaKeyDatabase, c.database,
		constants.RPCMetaKeyShardState, string(shardState))
	writeCli, err := writeService.Write(ctx)

	if err != nil {
		c.logger.Error("recvLoop get clientStreaming error", logger.Error(err))
		return nil
	}

	//TODO need close old stream client?
	c.writeCli = writeCli
	c.ready.Store(true)

	c.logger.Info("initialize write client stream successfully",
		logger.String("database", c.database),
		logger.Any("shard", c.shardID),
		logger.String("leader", (&storageNode).Indicator()))
	return nil
}

// sendLoop is a loop to send message to rpc stream, it recovers from panic to prevent crash.
// The loop only terminates when isStopped() return true.
func (c *sender) Send(data []byte) error {
	if !c.ready.Load() {
		if err := c.initClient(); err != nil {
			return err
		}
	}

	if !c.ready.Load() {
		//TODO add log
		return nil
	}

	err := c.writeCli.Send(
		&protoReplicaV1.WriteRequest{
			Record: data,
		})

	if err != nil {
		c.ready.Store(false)
		return err
	}
	return nil
}
func (c *sender) SyncShardState(shardState models.ShardState, liveNodes map[models.NodeID]models.StatefulNode) {
	c.lock4Meta.Lock()
	defer c.lock4Meta.Unlock()

	// set new shard state and storage live nodes.
	c.shardState = shardState
	c.liveNodes = liveNodes

	// check if need start write task
	if c.running.CAS(false, true) {

		c.logger.Info("start shard data sender successfully", logger.String("db", c.database),
			logger.Any("shardID", c.shardID))
	}
}

func (c *sender) Stop() {
	c.running.Store(false)
}
