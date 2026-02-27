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
	"sync"

	larrow "github.com/lindb/arrow/pkg/arrow"

	"github.com/lindb/lindb/app/broker/write/writer"
	"github.com/lindb/lindb/models"
)

type shard[V larrow.EntryType] struct {
	ctx        context.Context
	id         models.ShardID
	database   writer.DatabaseAccessor[V]
	segments   sync.Map // segment time -> segment
	shardState models.ShardState
	liveNodes  map[models.NodeID]models.StatefulNode

	lock sync.Mutex
}

func NewShard[V larrow.EntryType](ctx context.Context,
	database writer.DatabaseAccessor[V],
	shardState models.ShardState, liveNodes map[models.NodeID]models.StatefulNode,
) writer.Shard[V] {
	return &shard[V]{
		ctx:        ctx,
		database:   database,
		id:         shardState.ID,
		shardState: shardState,
		liveNodes:  liveNodes,
	}
}

func (s *shard[V]) GetDatabase() writer.DatabaseAccessor[V] {
	return s.database
}

func (s *shard[V]) ID() models.ShardID {
	return s.id
}

func (s *shard[V]) GetOrCreateSegment(segmentTime int64,
	builder func() larrow.EntryBuilder[V],
) writer.Segment[V] {
	p, ok := s.segments.Load(segmentTime)
	if ok {
		return p.(writer.Segment[V])
	}

	s.lock.Lock()
	defer s.lock.Unlock()

	// double check
	p, ok = s.segments.Load(segmentTime)
	if ok {
		return p.(writer.Segment[V])
	}

	// create new segment
	segment := NewSegment[V](s.ctx, s, segmentTime, s.shardState, s.liveNodes, builder())
	s.segments.Store(segmentTime, segment)
	return segment
}

func (s *shard[V]) LeaderChanged(shardState models.ShardState, liveNodes map[models.NodeID]models.StatefulNode) {
	s.segments.Range(func(key, value any) bool {
		segment := value.(writer.Segment[V])
		segment.LeaderChanged(shardState, liveNodes)
		return true
	})
}
