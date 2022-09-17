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
	"context"
	"runtime/pprof"
	"time"

	"go.uber.org/atomic"

	"github.com/lindb/lindb/metrics"
	"github.com/lindb/lindb/pkg/logger"
)

//go:generate mockgen -source=./replicator_peer.go -destination=./replicator_peer_mock.go -package=replica

// ReplicatorPeer represents wal replica peer.
// local replicator: from == to.
// remote replicator: from != to.
type ReplicatorPeer interface {
	// Startup starts wal replicator channel,
	Startup()
	// Shutdown shutdowns gracefully.
	Shutdown()
	// ReplicatorState returns the state and type of the replicator.
	ReplicatorState() (string, *state)
}

// replicatorPeer implements ReplicatorPeer
type replicatorPeer struct {
	runner  *replicatorRunner
	running *atomic.Bool
}

// NewReplicatorPeer creates a ReplicatorPeer.
func NewReplicatorPeer(replicator Replicator) ReplicatorPeer {
	return &replicatorPeer{
		running: atomic.NewBool(false),
		runner:  newReplicatorRunner(replicator),
	}
}

// Startup starts wal replicator channel,
func (r *replicatorPeer) Startup() {
	if r.running.CAS(false, true) {
		go func() {
			replicatorLabels := pprof.Labels("type", r.runner.replicatorType,
				"replicator", r.runner.replicator.String())
			pprof.Do(context.Background(), replicatorLabels, r.runner.replicaLoop)
		}()
	}
}

// Shutdown shutdowns gracefully.
func (r *replicatorPeer) Shutdown() {
	if r.running.CAS(true, false) {
		r.runner.shutdown()
	}
}

// ReplicatorState returns the state and type of the replicator.
func (r *replicatorPeer) ReplicatorState() (string, *state) {
	return r.runner.replicatorType, r.runner.replicator.State()
}

type replicatorRunner struct {
	running        *atomic.Bool
	lastPending    *atomic.Int64
	replicatorType string
	replicator     Replicator

	closed          chan struct{}
	sleep, maxSleep int
	sleepFn         func(d time.Duration)

	statistics *metrics.StorageReplicatorRunnerStatistics
	logger     *logger.Logger
}

func newReplicatorRunner(replicator Replicator) *replicatorRunner {
	replicaType := "local"
	if _, ok := replicator.(*remoteReplicator); ok {
		replicaType = "remote"
	}
	state := replicator.ReplicaState()
	r := &replicatorRunner{
		replicator:     replicator,
		lastPending:    atomic.NewInt64(replicator.Pending()),
		replicatorType: replicaType,
		running:        atomic.NewBool(false),
		closed:         make(chan struct{}),
		sleep:          0,
		sleepFn:        time.Sleep,
		maxSleep:       2 * 10 * 1000, // 20 sec
		statistics:     metrics.NewStorageReplicatorRunnerStatistics(replicaType, state.Database, state.ShardID.String()),
		logger:         logger.GetLogger("Replica", "ReplicatorRunner"),
	}
	// set replica lag callback
	r.statistics.ReplicaLag.SetGetValueFn(func(val *atomic.Float64) {
		pending := replicator.Pending()
		val.Add(float64(r.lastPending.Load() - pending))
		r.lastPending.Store(pending)
	})
	return r
}

func (r *replicatorRunner) replicaLoop(ctx context.Context) {
	if r.running.CAS(false, true) {
		r.statistics.ActiveReplicators.Incr()
		r.loop(ctx)
	}
}

func (r *replicatorRunner) shutdown() {
	if r.running.CAS(true, false) {
		// wait for stop replica loop
		<-r.closed
	}
}

func (r *replicatorRunner) loop(ctx context.Context) {
	for r.running.Load() {
		r.replica(ctx)
	}

	// exit replica loop
	close(r.closed)

	r.statistics.ActiveReplicators.Decr()
}

func (r *replicatorRunner) replica(_ context.Context) {
	defer func() {
		if recovered := recover(); recovered != nil {
			r.statistics.ReplicaPanics.Incr()
			r.logger.Error("panic when replica data",
				logger.Any("err", recovered),
				logger.Stack(),
			)
		}
	}()

	hasData := false

	if r.replicator.IsReady() && r.replicator.Connect() {
		seq := r.replicator.Consume()
		if seq >= 0 {
			r.logger.Debug("replica write ahead log",
				logger.String("type", r.replicatorType),
				logger.String("replicator", r.replicator.String()),
				logger.Int64("index", seq))
			hasData = true
			r.sleep = 0
			data, err := r.replicator.GetMessage(seq)
			if err != nil {
				r.replicator.IgnoreMessage(seq)
				r.statistics.ConsumeMessageFailures.Incr()
				r.logger.Warn("cannot get replica message data, ignore replica message",
					logger.String("replicator", r.replicator.String()),
					logger.Int64("index", seq), logger.Error(err))
			} else {
				r.statistics.ConsumeMessage.Incr()
				r.replicator.Replica(seq, data)

				r.statistics.ReplicaBytes.Add(float64(len(data)))
			}
		}
	} else {
		r.logger.Warn("replica is not ready", logger.String("replicator", r.replicator.String()))
	}
	if !hasData {
		sleep := 2 << r.sleep
		if sleep < r.maxSleep {
			r.sleep++
		}
		if sleep > r.maxSleep {
			sleep = r.maxSleep
		}
		r.sleepFn(time.Duration(sleep) * time.Millisecond)
	}
}
