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
	channel *ReplicatorChannel
}

func NewReplicator(channel *ReplicatorChannel) Replicator {
	return &replicator{
		channel: channel,
	}
}

func (r *replicator) Replica(_ int64, _ []byte) {
	// do nothing, need impl in child class
}

func (r *replicator) IsReady() bool {
	return true
}

func (r *replicator) Consume() int64 {
	return r.channel.Queue.Consume()
}

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

func (r *replicator) AppendIndex() int64 {
	return r.channel.Queue.Queue().HeadSeq()
}

func (r *replicator) ResetReplicaIndex(idx int64) error {
	return r.channel.Queue.SetHeadSeq(idx)
}

func (r *replicator) ResetAppendIndex(idx int64) {
	r.channel.Queue.Queue().SetAppendSeq(idx)
}

func (r *replicator) SetAckIndex(ackIdx int64) {
	r.channel.Queue.Ack(ackIdx)
}

func (r *replicator) String() string {
	return "[" +
		"database:" + r.channel.Database +
		",shard:" + strconv.Itoa(int(r.channel.ShardID)) +
		",from:" + strconv.Itoa(int(r.channel.From)) +
		",to:" + strconv.Itoa(int(r.channel.To)) +
		"]"
}
