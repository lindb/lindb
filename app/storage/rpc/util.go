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

package rpc

import (
	"errors"
	"fmt"

	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/storage"
	"github.com/lindb/lindb/storage/store"
)

// getOrCreateSegment returns write ahead log if it exists, else creates a new write ahead log.
func getOrCreateSegment(
	engine storage.Engine,
	database string,
	shardID models.ShardID,
	segmentTime int64,
	leader models.NodeID,
) (store.WriteAheadLog, error) {
	db, ok := engine.GetDatabase(database)
	if !ok {
		return nil, errors.New("database not exist")
	}

	shard, ok := db.GetShard(shardID)
	if !ok {
		fmt.Println(shardID)
		return nil, errors.New("shard not exist")
	}
	p, err := shard.GetOrCreatePartition(segmentTime)
	if err != nil {
		fmt.Printf("err0=%s\n", err)
		return nil, err
	}
	segment, err := p.GetOrCreateSegment(segmentTime)
	if err != nil {
		fmt.Printf("err1=%s\n", err)
		return nil, err
	}
	log, err := segment.GetOrCreateWAL(leader)
	if err != nil {
		fmt.Printf("err2=%s\n", err)
		return nil, err
	}
	return log, nil
}
