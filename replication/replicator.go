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

package replication

import (
	"context"
	"sync"
	"time"

	"go.uber.org/atomic"

	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/logger"
	"github.com/lindb/lindb/pkg/queue"
	"github.com/lindb/lindb/rpc"
	"github.com/lindb/lindb/rpc/proto/storage"
)

//go:generate mockgen -source=./replicator.go -destination=./replicator_mock.go -package=replication

const (
	batchReplicaSize = 10
	//maxPendingSeqSize = 100
	unaryRPCTimeout = time.Second * 3
)

// Replicator represents a task to replicate data to target.
type Replicator interface {
	// Target returns the target target for replication.
	Target() models.Node
	// Database returns the database attribution.
	Database() string
	// ShardID returns the shardID attribution.
	ShardID() int32
	// Pending returns the num of messages remaining to replicate.
	Pending() int64
	// ReplicaIndex returns the index of message replica
	ReplicaIndex() int64
	// AckIndex returns the index of message replica ack
	AckIndex() int64
	// Stop stops the replication task.
	Stop()
}

// replicator implements Replicator.
type replicator struct {
	target   models.Node
	database string
	shardID  int32
	// underlying fanOut records the replication process.
	fo queue.FanOut
	// factory to get write streamClient
	fct rpc.ClientStreamFactory
	// current WriteStreamClient
	streamClient storage.WriteService_WriteClient
	// current WriteServiceClient
	serviceClient storage.WriteServiceClient
	// lock to protect clients
	lock4client sync.RWMutex
	// false -> running, true -> stopped
	stopped atomic.Bool
	// false -> notReady, true -> ready
	ready atomic.Bool
	//storage received cur sequence num
	//storageCurSeq int64
	logger *logger.Logger
}

// newReplicator returns a Replicator with specific attributions.
func newReplicator(target models.Node, database string, shardID int32,
	fo queue.FanOut, fct rpc.ClientStreamFactory) Replicator {
	r := &replicator{
		target:   target,
		database: database,
		shardID:  shardID,
		fo:       fo,
		fct:      fct,
		logger:   logger.GetLogger("replication", "Replicator"),
	}

	go r.recvLoop()
	go r.sendLoop()

	return r
}

// Target returns the target target for replication.
func (r *replicator) Target() models.Node {
	return r.target
}

// Database returns the database attribution.
func (r *replicator) Database() string {
	return r.database
}

// ShardID returns the shardID attribution.
func (r *replicator) ShardID() int32 {
	return r.shardID
}

// Pending returns the num of messages remaining to replicate.
func (r *replicator) Pending() int64 {
	return r.fo.Pending()
}

// ReplicaIndex returns the index of message replica
func (r *replicator) ReplicaIndex() int64 {
	return r.fo.HeadSeq()
}

// AckIndex returns the index of message replica ack
func (r *replicator) AckIndex() int64 {
	return r.fo.TailSeq()
}

// Stop stops the replication task.
func (r *replicator) Stop() {
	r.stopped.Store(true)
}

// isStopped atomic check if is stopped.
func (r *replicator) isStopped() bool {
	return r.stopped.Load()
}

func (r *replicator) isReady() bool {
	return r.ready.Load()
}

func (r *replicator) setReady(ready bool) {
	r.ready.Store(ready)
}

// recvLoop is a loop to receive message from rpc stream.
// The loop recovers from panic to prevent crash.
// The loop handles rpc re-connection issues.
// The loop only terminates when isStopped() return true.
func (r *replicator) recvLoop() {
	defer func() {
		if rec := recover(); rec != nil {
			r.logger.Error("recover from panic, replicator.recvLoop",
				logger.Reflect("recover", rec),
				logger.Stack())

			r.logger.Info("restart recvLoop")
			//TODO modify sleep threshold for retry
			time.Sleep(500 * time.Millisecond)
			go r.recvLoop()
		}
	}()

	for {
		if !r.isReady() {
			r.initClient()
		}

		if r.isStopped() {
			r.logger.Info("end recvLoop")
			return
		}
		// when connection is stopped, replicator.streamClient.Recv() returns error.
		resp, err := r.streamClient.Recv()
		if err != nil {
			//fixme if seq out of range need reset
			r.logger.Error("recvLoop receive error", logger.Error(err))
			r.setReady(false)
			time.Sleep(time.Second)
			continue
		}

		// todo@TianliangXia use resp.curSeq for sliding window control
		// ackSeq could be nil, means no ack signal
		ack, ok := resp.Ack.(*storage.WriteResponse_AckSeq)
		if ok {
			r.fo.Ack(ack.AckSeq)
		}
	}
}

func (r *replicator) initClient() {
	// try to re-construct the streaming
	for {
		if r.isStopped() {
			return
		}

		serviceClient, err := r.fct.CreateWriteServiceClient(r.target)
		if err != nil {
			r.logger.Error("recvLoop get service streamClient error", logger.Error(err))
			time.Sleep(time.Second)
			continue
		}
		r.serviceClient = serviceClient

		// get storage head seq, reset fanOut headSeq or reset storage headSeq.
		nextSeq, err := r.remoteNextSeq()
		if err != nil {
			r.logger.Error("recvLoop get remote next seq error", logger.Error(err),
				logger.String("target", r.target.Indicator()))
			// typically CreateWriteServiceClient won't return err if remote target is unavailable(async dial), the real rpc call will.
			// sleep to avoid dead for loop
			time.Sleep(time.Second)
			continue
		}

		// try to reset fanOut headSeq, if success, consume from new headSeq,
		// if fail, try to reset remote headSeq.
		r.logger.Info("recvLoop try to set fanOut head seq", logger.Int64("headSeq", nextSeq))
		if err := r.fo.SetHeadSeq(nextSeq); err != nil {
			r.logger.Error("recvLoop reset fanOut head seq error", logger.Error(err))

			foHeadSeq := r.fo.HeadSeq()
			r.logger.Info("recvLoop try to set remote storage head seq", logger.Int64("headSeq", foHeadSeq))
			if err := r.resetRemoteSeq(foHeadSeq); err != nil {
				r.logger.Error("recvLoop reset remote head seq error", logger.Error(err))
				continue
			}
		}

		streamClient, err := r.fct.CreateWriteClient(r.database, r.shardID, r.target)
		if err != nil {
			r.logger.Error("recvLoop get clientStreaming error", logger.Error(err))
			continue
		}

		r.logger.Info("recvLoop get clientStreaming success")
		r.lock4client.Lock()
		r.streamClient = streamClient
		r.lock4client.Unlock()
		break
	}
	r.setReady(true)
}

func (r *replicator) remoteNextSeq() (int64, error) {
	nextReq := &storage.NextSeqRequest{
		Database: r.database,
		ShardID:  r.shardID,
	}

	ctx, cancel := context.WithTimeout(context.TODO(), unaryRPCTimeout)
	ctx = rpc.CreateOutgoingContextWithNode(ctx, r.fct.LogicNode())
	nextResp, err := r.serviceClient.Next(ctx, nextReq)
	cancel()
	if err != nil {
		r.logger.Debug("failed to call remoteNextSeq",
			logger.String("database", r.database),
			logger.Int32("shardID", r.shardID),
		)
		return -1, err
	}
	return nextResp.Seq, nil
}

func (r *replicator) resetRemoteSeq(resetSeq int64) error {
	// reset storage headSeq
	nextReq := &storage.ResetSeqRequest{
		Database: r.database,
		ShardID:  r.shardID,
		Seq:      resetSeq,
	}
	ctx, cancel := context.WithTimeout(context.TODO(), unaryRPCTimeout)
	ctx = rpc.CreateOutgoingContextWithNode(ctx, r.fct.LogicNode())
	// response body is empty, if no error return, reset seq success
	_, err := r.serviceClient.Reset(ctx, nextReq)
	cancel()
	return err
}

// sendLoop is a loop to send message to rpc stream, it recovers from panic to prevent crash.
// The loop only terminates when isStopped() return true.
func (r *replicator) sendLoop() {
	defer func() {
		if rec := recover(); rec != nil {
			r.logger.Error("recover from panic, replicator.sendLoop",
				logger.Reflect("recover", rec),
				logger.Stack())

			r.logger.Info("restart sendLoop")
			go r.sendLoop()
		}
	}()

	// reuse the fix size slice
	reusedReplicas := make([]*storage.Replica, 0, batchReplicaSize)

	for {
		if r.isStopped() {
			r.logger.Info("end sendLoop")
			return
		}

		// conn not ready
		if !r.isReady() {
			time.Sleep(time.Second)
			continue
		}

		replicas := r.consumeBatch(&reusedReplicas)
		// no more replicas
		if len(replicas) == 0 {
			time.Sleep(10 * time.Millisecond)
			continue
		}
		wr := &storage.WriteRequest{
			Replicas: replicas,
		}

		r.logger.Debug("send replicas",
			logger.Int64("begin", replicas[0].Seq),
			logger.Int64("end", replicas[len(replicas)-1].Seq))

		// recvLoop may change streamClient
		r.lock4client.RLock()
		cli := r.streamClient
		r.lock4client.RUnlock()
		if err := cli.Send(wr); err != nil {
			r.logger.Error("sendLoop write request error", logger.Error(err))
			r.setReady(false)
		}
	}
}

// consumeBatch consumes a batch of Replicas(limited by batchReplicaSize), the input slice is reused.
func (r *replicator) consumeBatch(repPointer *[]*storage.Replica) []*storage.Replica {
	replicas := *repPointer
	replicas = replicas[:0]
	var i int
	for i = 0; i < batchReplicaSize; i++ {
		seq := r.fo.Consume()
		if seq == queue.SeqNoNewMessageAvailable {
			break
		}

		data, err := r.fo.Get(seq)
		if err != nil {
			r.logger.Error("get message from fanout queue error", logger.String("database", r.database),
				logger.Int32("shardID", r.shardID))
			break
		}

		replica := &storage.Replica{
			Seq:  seq,
			Data: data,
		}
		replicas = append(replicas, replica)
	}
	return replicas[:i]
}
