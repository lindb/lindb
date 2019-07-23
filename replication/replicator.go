package replication

import (
	"sync"
	"sync/atomic"
	"time"

	"go.uber.org/zap"

	"github.com/eleme/lindb/models"
	"github.com/eleme/lindb/pkg/logger"
	"github.com/eleme/lindb/pkg/queue"
	"github.com/eleme/lindb/rpc"
	"github.com/eleme/lindb/rpc/proto/storage"
)

const (
	batchReplicaSize = 10
)

// Replicator represents a task to replicate data to target.
type Replicator interface {
	// Target returns the target target for replication.
	Target() models.Node
	// Cluster returns the cluster attribution.
	Cluster() string
	// Database returns the database attribution.
	Database() string
	// ShardID returns the shardID attribution.
	ShardID() uint32
	// Pending returns the num of messages remaining to replicate.
	Pending() int64
	// Stop stops the replication task.
	Stop()
}

// replicator implements Replicator.
type replicator struct {
	target   models.Node
	cluster  string
	database string
	shardID  uint32
	// underlying fanOut records the replication process.
	fo queue.FanOut
	// factory to get write client
	fct rpc.ClientStreamFactory
	// current WriteClient
	client storage.WriteService_WriteClient
	// lock to protect client
	lock4client sync.RWMutex
	// 0-> running, 1 -> stopped
	stopped int32
	logger  *logger.Logger
}

// newReplicator returns a Replicator with specific attributions.
func newReplicator(target models.Node, cluster, database string, shardID uint32,
	fo queue.FanOut, fct rpc.ClientStreamFactory) (Replicator, error) {
	r := &replicator{
		target:   target,
		cluster:  cluster,
		database: database,
		shardID:  shardID,
		fo:       fo,
		fct:      fct,
		logger:   logger.GetLogger("replication/replicator"),
	}

	client, err := fct.CreateWriteClient(database, shardID, target)
	if err != nil {
		return nil, err
	}

	r.client = client

	go r.recvLoop()
	go r.sendLoop()

	return r, nil
}

// Target returns the target target for replication.
func (r *replicator) Target() models.Node {
	return r.target
}

// Cluster returns the cluster attribution.
func (r *replicator) Cluster() string {
	return r.cluster
}

// Database returns the database attribution.
func (r *replicator) Database() string {
	return r.database
}

// ShardID returns the shardID attribution.
func (r *replicator) ShardID() uint32 {
	return r.shardID
}

// Pending returns the num of messages remaining to replicate.
func (r *replicator) Pending() int64 {
	return r.fo.Pending()
}

// Stop stops the replication task.
func (r *replicator) Stop() {
	atomic.StoreInt32(&r.stopped, 1)
}

// isStopped atomic check if is stopped.
func (r *replicator) isStopped() bool {
	return atomic.LoadInt32(&r.stopped) == 1
}

// recvLoop is a loop to receive message from rpc stream.
// The loop recovers from panic to prevent crash.
// The loop handles rpc re-connection issues.
// The loop only terminates when isStopped() return true.
func (r *replicator) recvLoop() {
	defer func() {
		if rec := recover(); rec != nil {
			r.logger.Error("recover from panic, replicator.recvLoop",
				zap.Reflect("rec", rec),
				zap.Stack("stack"))

			r.logger.Info("restart recvLoop")
			go r.recvLoop()
		}
	}()

	for {
		// when connection is stopped, replicator.client.Recv() returns error.
		resp, err := r.client.Recv()

		if err != nil {
			r.logger.Error("recvLoop receive error", zap.Error(err))
			for {
				// try to re-construct the streaming
				if r.isStopped() {
					r.logger.Info("end recvLoop")
					return
				}

				client, err := r.fct.CreateWriteClient(r.database, r.shardID, r.target)
				if err != nil {
					r.logger.Error("recvLoop re-construct the streaming error", zap.Error(err))
					time.Sleep(time.Second)
					continue
				}
				r.logger.Info("recvLoop re-construct the streaming success")
				r.lock4client.Lock()
				r.client = client
				r.lock4client.Unlock()
				break

			}
			continue
		}

		//logger.GetLogger("replication").Info("receive:", logger.Error(err), logger.Any("resp", resp))

		switch resp.Seq.(type) {
		case *storage.WriteResponse_AckSeq:
			r.fo.Ack(resp.GetAckSeq())
		case *storage.WriteResponse_ResetSeq:
			r.logger.Warn("reset head seq", zap.Int64("restSeq", resp.GetResetSeq()))
			if err := r.fo.SetHeadSeq(resp.GetResetSeq()); err != nil {
				r.logger.Error("reset head seq error", zap.Error(err))
			}
		}
	}

}

// sendLoop is a loop to send message to rpc stream, it recovers from panic to prevent crash.
// The loop only terminates when isStopped() return true.
func (r *replicator) sendLoop() {
	defer func() {
		if rec := recover(); rec != nil {
			r.logger.Error("recover from panic, replicator.sendLoop",
				zap.Reflect("rec", rec),
				zap.Stack("stack"))

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
		replicas := r.consumeBatch(&reusedReplicas)
		if len(replicas) == 0 {
			time.Sleep(10 * time.Millisecond)
			continue
		}
		wr := &storage.WriteRequest{
			Replicas: replicas,
		}
		seqs := make([]int64, 0, len(wr.Replicas))
		for _, rep := range wr.Replicas {
			seqs = append(seqs, rep.Seq)
		}
		r.logger.Info("send ", zap.Any("wr", seqs))

		r.lock4client.RLock()
		cli := r.client
		r.lock4client.RUnlock()
		if err := cli.Send(wr); err != nil {
			r.logger.Error("sendLoop write request error", zap.Error(err))
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
			r.logger.Error("get message from fanout queue error", zap.String("database", r.database),
				zap.Uint32("shardID", r.shardID))
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
