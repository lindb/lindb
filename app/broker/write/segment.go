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

package write

import (
	"context"
	"fmt"
	"runtime/pprof"
	"sync"

	"github.com/apache/arrow-go/v18/arrow"
	larrow "github.com/lindb/arrow/pkg/arrow"
	"github.com/lindb/common/pkg/logger"
	"github.com/lindb/common/pkg/timeutil"

	"github.com/lindb/lindb/app/broker/write/writer"
	"github.com/lindb/lindb/models"
	protoWriteV1 "github.com/lindb/lindb/proto/gen/v1/write"
)

type segment[V any] struct {
	ctx    context.Context
	cancel context.CancelFunc

	shard       writer.Shard[V]
	segmentTime int64
	builder     larrow.EntryBuilder[V]

	entries chan V
	records chan []byte

	shardState models.ShardState
	liveNodes  map[models.NodeID]models.StatefulNode
	state      *models.SegmentState

	serializers []*larrow.Serializer

	lock4state sync.Mutex

	logger logger.Logger
}

func NewSegment[V any](parent context.Context,
	shard writer.Shard[V], segmentTime int64,
	shardState models.ShardState, liveNodes map[models.NodeID]models.StatefulNode,
	builder larrow.EntryBuilder[V],
) writer.Segment[V] {
	ctx, cancel := context.WithCancel(parent)
	p := &segment[V]{
		ctx:         ctx,
		cancel:      cancel,
		shard:       shard,
		shardState:  shardState,
		liveNodes:   liveNodes,
		segmentTime: segmentTime,
		builder:     builder,
		entries:     make(chan V, 1024),
		records:     make(chan []byte, 256),
		logger:      logger.GetLogger("Write", "Segment"),
	}
	p.initialize()
	return p
}

func (p *segment[V]) initialize() {
	database := p.shard.GetDatabase()
	p.state = &models.SegmentState{
		Database:    database.Name(),
		SegmentTime: p.segmentTime,
	}

	// init serializers
	schemaIDs := p.builder.SchemaIDs()
	p.serializers = make([]*larrow.Serializer, len(schemaIDs))
	for i, schemaID := range schemaIDs {
		p.serializers[i] = larrow.NewSerializer(schemaID)
	}

	// start build and send goroutines
	go func() {
		channelFamilyLabels := pprof.Labels("type", "builder", "database", database.Name(),
			"shard", fmt.Sprintf("%d", p.shard.ID()),
			"segment", timeutil.FormatTimestamp(p.segmentTime, timeutil.DataTimeFormat2))
		pprof.Do(p.ctx, channelFamilyLabels, p.build)
	}()

	go func() {
		channelFamilyLabels := pprof.Labels("type", "send", "database", database.Name(),
			"shard", fmt.Sprintf("%d", p.shard.ID()),
			"segment", timeutil.FormatTimestamp(p.segmentTime, timeutil.DataTimeFormat2))
		pprof.Do(p.ctx, channelFamilyLabels, p.send)
	}()
}

func (p *segment[V]) Write(ctx context.Context, entry V) {
	select {
	case <-ctx.Done():
	case p.entries <- entry:
	}
}

func (p *segment[V]) LeaderChanged(shardState models.ShardState, liveNodes map[models.NodeID]models.StatefulNode) {
	p.lock4state.Lock()
	defer p.lock4state.Unlock()

	p.shardState = shardState
	p.liveNodes = liveNodes

	fmt.Println("leader change", liveNodes, shardState)
}

func (p *segment[V]) Close() {
	// release records builder
	p.builder.Release()
}

func (p *segment[V]) build(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case entry := <-p.entries:
			p.builder.Append(entry)
			if p.builder.NumOfRows() >= 10 {
				// emit record batches bytes to send channel
				p.records <- p.builder.Bytes()
			}
		}
	}
}

func (p *segment[V]) send(ctx context.Context) {
	var sender Sender
	var err error
	for {
		select {
		case <-ctx.Done():
			return
		case records := <-p.records:
			if sender == nil {
				// if sender is nil, build sender with current shard state and live nodes
				sender, err = p.buildSender()
				if err != nil {
					p.logger.Error("failed to build sender", logger.String("database", p.state.Database), logger.Error(err))
					continue
				}
			}

			if err := sender.Send(records); err != nil {
				p.logger.Error("failed to send data", logger.String("database", p.state.Database), logger.Error(err))
			}
			// TODO: add retry logic for send failure/send pending
		}
	}
}

func (p *segment[V]) serialze(records []arrow.RecordBatch) ([]*protoWriteV1.ArrowPayload, error) {
	defer func() {
		// need release record after serialize
		for _, record := range records {
			record.Release()
		}
	}()
	payloads := make([]*protoWriteV1.ArrowPayload, len(records))
	for i, record := range records {
		data, err := p.serializers[i].Serialize(record)
		if err != nil {
			return nil, err
		}
		payloads[i] = &protoWriteV1.ArrowPayload{
			SchemaIndex: int32(i),
			Record:      data,
		}
	}

	return payloads, nil
}

func (p *segment[V]) buildSender() (Sender, error) {
	p.lock4state.Lock()
	defer p.lock4state.Unlock()

	if len(p.liveNodes) == 0 {
		return nil, fmt.Errorf("no live node for shard %d", p.shard.ID())
	}
	leader, ok := p.liveNodes[p.shardState.Leader]
	if !ok {
		return nil, fmt.Errorf("leader node %d is not live for shard %d", p.shardState.Leader, p.shard.ID())
	}
	// set shard state for sender
	p.state.Shard = p.shardState
	return NewSender(p.ctx, &leader, p.state)
}
