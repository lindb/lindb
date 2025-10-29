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

package store

import (
	"fmt"
	"io"

	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/queue"
)

type ReplicatorType string

var (
	CreateWriteAheadLog  func(path string, segment Segment) (WriteAheadLog, error)
	CreateReplicatorPeer func(replicator Replicator) ReplicatorPeer
)

const (
	ReplicatorTypeLocal  ReplicatorType = "local"
	ReplicatorTypeRemote ReplicatorType = "remote"
)

type WriteAheadLog interface {
	io.Closer

	Get(index int64) ([]byte, error)
	Write(msg []byte) error
	Replica(replicaIndex int64, msg []byte) (int64, error)

	// ReplicaAckIndex returns the index which replica appended index.
	ReplicaAckIndex() int64
	// ResetReplicaIndex resets replica index.
	ResetReplicaIndex(idx int64)

	// BuildReplicaForLeader builds replica relation when handle writeTask connection.
	BuildReplicaForLeader(leader models.NodeID, replicas []models.NodeID) error
	// BuildReplicaForFollower builds replica relation when handle replica connection.
	BuildReplicaForFollower(leader models.NodeID, replica models.NodeID) error
}

// ReplicatorState represents the state of replicator.
type ReplicatorState struct {
	ErrMsg string
	State  models.ReplicatorState
}

// Replicator represents write ahead log replicator.
type Replicator interface {
	fmt.Stringer
	// ReplicaState returns the replica state.
	ReplicaState() *models.ReplicaState
	// State returns the state of replicator.
	State() *ReplicatorState
	// Pause paused replica data.
	Pause()
	// Consume returns the index of message replica.
	Consume() int64
	// GetMessage returns message by replica index.
	GetMessage(replicaIdx int64) ([]byte, error)
	// Replica replicas message by replica index.
	Replica(idx int64, msg []byte)
	// IsReady returns if replicator is ready.
	IsReady() bool
	// Connect connects follower for sending replica message.
	Connect() bool
	// ReplicaIndex returns the index of message replica
	ReplicaIndex() int64
	// AckIndex returns the index of message replica ack
	AckIndex() int64
	// AppendIndex returns next append index.
	AppendIndex() int64
	// ResetReplicaIndex resets replica index.
	ResetReplicaIndex(idx int64)
	// ResetAppendIndex resets append index.
	ResetAppendIndex(idx int64)
	// SetAckIndex sets ack index.
	SetAckIndex(ackIdx int64)
	// Pending returns lag of queue.
	Pending() int64
	// IgnoreMessage ignores invalid message.
	IgnoreMessage(replicaIdx int64)
	// Close closes replicator, releases resource.
	Close()
}

// ReplicatorPeer represents wal replica peer.
// local replicator: from == to.
// remote replicator: from != to.
type ReplicatorPeer interface {
	// Startup starts wal replicator channel,
	Startup()
	// Shutdown shutdowns gracefully.
	Shutdown()
	// ReplicatorState returns the state and type of the replicator.
	ReplicatorState() (string, *ReplicatorState)
}

// ReplicatorChannel represents channel peer[from,to] for the shard of database.
type ReplicatorChannel struct {
	State *models.ReplicaState

	// underlying ConsumerGroup records the replication process.
	ConsumerGroup queue.ConsumerGroup
}
