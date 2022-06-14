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

	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/timeutil"
)

//go:generate mockgen -source=./replicator.go -destination=./replicator_mock.go -package=replica

// state represents the state of replicator.
type state struct {
	state  models.ReplicatorState
	errMsg string
}

// Replicator represents write ahead log replicator.
type Replicator interface {
	fmt.Stringer
	// ReplicaState returns the replica state.
	ReplicaState() *models.ReplicaState
	// State returns the state of replicator.
	State() *state
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
	ResetReplicaIndex(idx int64) error
	// ResetAppendIndex resets append index.
	ResetAppendIndex(idx int64)
	// SetAckIndex sets ack index.
	SetAckIndex(ackIdx int64)
	// Pending returns lag of queue.
	Pending() int64
}

// replicator implements Replicator interface.
type replicator struct {
	channel *ReplicatorChannel
}

// ReplicaState returns the replica state.
func (r *replicator) ReplicaState() *models.ReplicaState {
	return r.channel.State
}

// Replica replicas message by replica index.
func (r *replicator) Replica(_ int64, _ []byte) {
	// do nothing, need impl in child class
}

// IsReady returns if replicator is ready.
func (r *replicator) IsReady() bool {
	return true
}

// Connect connects follower for sending replica message.
func (r *replicator) Connect() bool {
	return true
}

// Consume returns the index of message replica.
func (r *replicator) Consume() int64 {
	return r.channel.Queue.Consume()
}

// GetMessage returns message by replica index.
func (r *replicator) GetMessage(replicaIdx int64) ([]byte, error) {
	return r.channel.Queue.Get(replicaIdx)
}

// ReplicaIndex returns the index of message replica
func (r *replicator) ReplicaIndex() int64 {
	return r.channel.Queue.HeadSeq()
}

// AckIndex returns the index of message replica ack.
func (r *replicator) AckIndex() int64 {
	return r.channel.Queue.TailSeq()
}

// AppendIndex returns next append index.
func (r *replicator) AppendIndex() int64 {
	return r.channel.Queue.Queue().HeadSeq()
}

// ResetReplicaIndex resets replica index.
func (r *replicator) ResetReplicaIndex(idx int64) error {
	return r.channel.Queue.SetHeadSeq(idx)
}

// ResetAppendIndex resets append index.
func (r *replicator) ResetAppendIndex(idx int64) {
	r.channel.Queue.Queue().SetAppendSeq(idx)
}

// SetAckIndex sets ack index.
func (r *replicator) SetAckIndex(ackIdx int64) {
	r.channel.Queue.Ack(ackIdx)
}

// Pending returns lag of queue.
func (r *replicator) Pending() int64 {
	return r.channel.Queue.Pending()
}

// String returns string value of replicator.
func (r *replicator) String() string {
	return "[" +
		"database:" + r.channel.State.Database +
		",shard:" + r.channel.State.ShardID.String() +
		",family:" + timeutil.FormatTimestamp(r.channel.State.FamilyTime, timeutil.DataTimeFormat2) +
		",from(leader):" + r.channel.State.Leader.String() +
		",to(follower):" + r.channel.State.Follower.String() +
		"]"
}
