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
	"fmt"

	"github.com/lindb/lindb/meta"
	"github.com/lindb/lindb/models"
	protoReplicaV1 "github.com/lindb/lindb/proto/gen/v1/replica"
	"github.com/lindb/lindb/rpc"
)

type coordinator struct {
	engine Engine

	streamings map[string]*models.StreamingState
}

func NewCoordinator(engine Engine) *coordinator {
	return &coordinator{
		engine:     engine,
		streamings: make(map[string]*models.StreamingState),
	}
}

func (c *coordinator) OnEvent(event meta.Event) {
	switch stateEvent := event.(type) {
	case *models.StreamingState:
		db, ok := c.engine.GetDatabase(stateEvent.Config.Database)
		if !ok {
			return
		}

		c.streamings[stateEvent.Config.Name] = stateEvent

		// mark shards as consuming state
		for _, consumerAssign := range stateEvent.ConsumeAssignments {
			for _, shardID := range consumerAssign.Shards {
				shard, ok := db.GetShard(shardID)
				if !ok {
					continue
				}
				fmt.Println(shard)
				shard.Consume(stateEvent.Config.Name, consumerAssign.ConsumerID)
			}
		}
	default:
		fmt.Println("TODO implement me")
	}
	// TODO implement me
	// panic("implement me")
}

func (c *coordinator) GetLiveNode(name string, nodeID models.NodeID) (models.Node, bool) {
	state, ok := c.streamings[name]
	if !ok {
		return nil, false
	}
	node, ok := state.Consumers[nodeID]
	if !ok {
		return nil, false
	}
	return &node, true
}

func (c *coordinator) CreateReplicaServiceClient(target models.Node) (protoReplicaV1.ReplicaServiceClient, error) {
	conn, err := rpc.GetStorageClientConnFactory().GetClientConn(target)
	if err != nil {
		return nil, err
	}
	return protoReplicaV1.NewReplicaServiceClient(conn), nil
}
