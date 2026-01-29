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

	"github.com/samber/lo"

	"github.com/lindb/lindb/coordinator/discovery"
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

	subChan   chan discovery.Subscriber
	unsubChan chan discovery.Subscriber
	events    chan discovery.MetaEvent
}

func NewCoordinator(engine Engine) *coordinator {
	ctx, cancel := context.WithCancel(context.Background())
	c := &coordinator{
		ctx:         ctx,
		cancal:      cancel,
		engine:      engine,
		events:      make(chan discovery.MetaEvent, 8),
		subChan:     make(chan discovery.Subscriber, 8),
		unsubChan:   make(chan discovery.Subscriber, 8),
		streamings:  make(map[string]*models.StreamingState),
		subscribers: make(map[string]*subscribers),
	}

	go c.run()

	return c
}

func (c *coordinator) run() {
	for {
		select {
		case sub := <-c.subChan:
			c.subscribe(sub)
		case unsub := <-c.unsubChan:
			c.unsubscribe(unsub)
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

func (c *coordinator) process(event discovery.MetaEvent) {
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

func (c *coordinator) OnEvent(event discovery.MetaEvent) {
	c.events <- event
}

func (c *coordinator) Subscribe(sub discovery.Subscriber) {
	c.subChan <- sub
}

func (c *coordinator) subscribe(sub discovery.Subscriber) {
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

func (c *coordinator) unsubscribe(sub discovery.Subscriber) {
	database := sub.Database()
	subs, ok := c.subscribers[database]
	if !ok {
		return
	}
	subs.unsubscribe(sub)
}

func (c *coordinator) Unsubscribe(sub discovery.Subscriber) {
	c.unsubChan <- sub
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
	subscribers map[models.ShardID][]discovery.Subscriber
}

func newSubscribers() *subscribers {
	return &subscribers{
		subscribers: make(map[models.ShardID][]discovery.Subscriber),
	}
}

func (s *subscribers) subscribe(sub discovery.Subscriber) {
	shardID := sub.Shard()
	subs, ok := s.subscribers[shardID]
	if !ok {
		subs = []discovery.Subscriber{}
	}
	subs = append(subs, sub)
	s.subscribers[shardID] = subs
}

func (s *subscribers) unsubscribe(sub discovery.Subscriber) {
	shardID := sub.Shard()
	subs, ok := s.subscribers[shardID]
	if !ok {
		return
	}
	// remove subscriber from list
	s.subscribers[shardID] = lo.Filter(subs, func(s discovery.Subscriber, _ int) bool {
		return s != sub
	})
}

func (s *subscribers) getSubscribers(shardID models.ShardID) []discovery.Subscriber {
	return s.subscribers[shardID]
}
