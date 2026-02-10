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

package writer

import (
	"context"

	larrow "github.com/lindb/arrow/pkg/arrow"

	"github.com/lindb/lindb/constants"
	"github.com/lindb/lindb/models"
)

type Writer interface {
	Write(ctx context.Context, data []byte, encoding constants.EncodingType) error
}

type Database interface {
	CreateShards(shards map[models.ShardID]models.ShardState, liveNodes map[models.NodeID]models.StatefulNode)
	LeaderChanged(shards map[models.ShardID]models.ShardState, liveNodes map[models.NodeID]models.StatefulNode)
}

type DatabaseAccessor[V any] interface {
	Database
	Name() string
	GetShards() []Shard[V]
}

type Shard[V any] interface {
	ID() models.ShardID
	GetDatabase() DatabaseAccessor[V]
	GetOrCreateSegment(segmentTime int64, builder func() larrow.EntryBuilder[V]) Segment[V]
	LeaderChanged(shardState models.ShardState, liveNodes map[models.NodeID]models.StatefulNode)
}

type Segment[V any] interface {
	Write(ctx context.Context, entry V)
	LeaderChanged(shardState models.ShardState, liveNodes map[models.NodeID]models.StatefulNode)
	Close()
}
