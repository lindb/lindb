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
	"fmt"
	"strconv"

	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/queue"
)

//go:generate mockgen -source=./replicator.go -destination=./replicator_mock.go -package=replica

type ReplicatorState int

const (
	ReplicatorInitState ReplicatorState = iota
	ReplicatorReadyState
	ReplicatorFailureState
)

type Replicator interface {
	fmt.Stringer

	// Consume returns the index of message replica.
	Consume() int64
	GetMessage(replicaIdx int64) ([]byte, error)
	Replica(idx int64, msg []byte)
	IsReady() bool
	// ReplicaIndex returns the index of message replica
	ReplicaIndex() int64
	// AckIndex returns the index of message replica ack
	AckIndex() int64
	AppendIndex() int64
	ResetReplicaIndex(idx int64) error
	ResetAppendIndex(idx int64)
	SetAckIndex(ackIdx int64)
}

type replicator struct {
	database string
	shardID  models.ShardID
	// underlying fanOut records the replication process.
	queue    queue.FanOut
	from, to models.NodeID // replicator node peer
}

func NewReplicator(database string, shardID models.ShardID, from, to models.NodeID, wal queue.FanOut) Replicator {
	return &replicator{
		database: database,
		shardID:  shardID,
		queue:    wal,
		from:     from,
		to:       to,
	}
}

func (r *replicator) Replica(_ int64, _ []byte) {
	// do nothing, need impl in child class
}

func (r *replicator) IsReady() bool {
	return true
}

func (r *replicator) Consume() int64 {
	return r.queue.Consume()
}

func (r *replicator) GetMessage(replicaIdx int64) ([]byte, error) {
	return r.queue.Get(replicaIdx)
}

// ReplicaIndex returns the index of message replica
func (r *replicator) ReplicaIndex() int64 {
	return r.queue.HeadSeq()
}

// AckIndex returns the index of message replica ack
func (r *replicator) AckIndex() int64 {
	return r.queue.TailSeq()
}

func (r *replicator) AppendIndex() int64 {
	return r.queue.Queue().HeadSeq()
}

func (r *replicator) ResetReplicaIndex(idx int64) error {
	return r.queue.SetHeadSeq(idx)
}

func (r *replicator) ResetAppendIndex(idx int64) {
	r.queue.Queue().SetAppendSeq(idx)
}

func (r *replicator) SetAckIndex(ackIdx int64) {
	r.queue.Ack(ackIdx)
}

func (r *replicator) String() string {
	return "[" +
		"database:" + r.database +
		",shard:" + strconv.Itoa(int(r.shardID)) +
		",from:" + string(r.from) +
		",to:" + string(r.to) +
		"]"
}
