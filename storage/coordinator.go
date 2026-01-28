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
	"fmt"
	"sync"

	"github.com/lindb/lindb/meta"
	"github.com/lindb/lindb/models"
	protoReplicaV1 "github.com/lindb/lindb/proto/gen/v1/replica"
	"github.com/lindb/lindb/rpc"
	"github.com/lindb/lindb/storage/store"
)

type coordinator struct {
	ctx    context.Context
	cancal context.CancelFunc

	engine Engine

	streamings  map[string]*models.StreamingState
	subscribers map[string]*subscribers

	events chan meta.Event
	mutex  sync.Mutex
}

func NewCoordinator(engine Engine) *coordinator {
	ctx, cancel := context.WithCancel(context.Background())
	c := &coordinator{
		ctx:         ctx,
		cancal:      cancel,
		engine:      engine,
		events:      make(chan meta.Event, 8),
		streamings:  make(map[string]*models.StreamingState),
		subscribers: make(map[string]*subscribers),
	}

	go c.run()

	return c
}

func (c *coordinator) run() {
	for {
		select {
		case event, ok := <-c.events:
			if !ok {
				return
			}
			c.process(event)
		case <-c.ctx.Done():
			return
		}
	}
}

func (c *coordinator) process(event meta.Event) {
	switch stateEvent := event.(type) {
	case *models.DeleteStreaming:
		state, ok := c.streamings[stateEvent.Streaming]
		if !ok {
			return
		}
		defer delete(c.streamings, stateEvent.Streaming)

		subscribers, ok := c.subscribers[state.Config.Database]
		if !ok {
			return
		}
		for _, consumerAssign := range state.ConsumeAssignments {
			for _, shardID := range consumerAssign.Shards {
				subs := subscribers.getSubscribers(shardID)
				for _, sub := range subs {
					sub.Receive(&store.ConsumerStateChange{
						Streaming:  state.Config.Name,
						ConsumerID: consumerAssign.ConsumerID,
						IsDelete:   true,
					})
				}
			}
		}
	case *models.StreamingState:
		c.streamings[stateEvent.Config.Name] = stateEvent
		subscribers, ok := c.subscribers[stateEvent.Config.Database]
		if !ok {
			return
		}
		for _, consumerAssign := range stateEvent.ConsumeAssignments {
			for _, shardID := range consumerAssign.Shards {
				subs := subscribers.getSubscribers(shardID)
				for _, sub := range subs {
					sub.Receive(&store.ConsumerStateChange{
						Streaming:  stateEvent.Config.Name,
						ConsumerID: consumerAssign.ConsumerID,
					})
				}
			}
		}
	case *models.CreateShard:
		// TODO: add streaming coordinator create shard logic?
		c.engine.CreateShards(stateEvent.Database, stateEvent.Option, stateEvent.Shards)
	default:
		fmt.Println("TODO implement me")
	}
	// TODO implement me
	// panic("implement me")
}

func (c *coordinator) OnEvent(event meta.Event) {
	c.events <- event
}

func (c *coordinator) Subscribe(sub meta.Subscriber) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	database := sub.Database()
	subs, ok := c.subscribers[database]
	if !ok {
		subs = newSubscribers()
	}
	subs.subscribe(sub)
	c.subscribers[database] = subs

	for _, stateEvent := range c.streamings {
		if stateEvent.Config.Database != database {
			continue
		}

		for _, consumerAssign := range stateEvent.ConsumeAssignments {
			for _, shardID := range consumerAssign.Shards {
				if shardID != sub.Shard() {
					continue
				}

				sub.Receive(&store.ConsumerStateChange{
					Streaming:  stateEvent.Config.Name,
					ConsumerID: consumerAssign.ConsumerID,
				})
			}
		}
	}
}

func (c *coordinator) Unsubscribe(sub meta.Subscriber) {
	// no need to subscribe
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

func (c *coordinator) Close() {
	c.cancal()
	close(c.events)
}

type subscribers struct {
	subscribers map[models.ShardID][]meta.Subscriber
	mutex       sync.Mutex
}

func newSubscribers() *subscribers {
	return &subscribers{
		subscribers: make(map[models.ShardID][]meta.Subscriber),
	}
}

func (s *subscribers) subscribe(sub meta.Subscriber) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	shardID := sub.Shard()
	subs, ok := s.subscribers[shardID]
	if !ok {
		subs = []meta.Subscriber{}
	}
	subs = append(subs, sub)
	s.subscribers[shardID] = subs
}

func (s *subscribers) getSubscribers(shardID models.ShardID) []meta.Subscriber {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	return s.subscribers[shardID]
}
