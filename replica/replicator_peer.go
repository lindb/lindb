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
	"time"

	"go.uber.org/atomic"

	"github.com/lindb/lindb/pkg/logger"
)

// ReplicatorPeer represents wal replica peer.
// local replicator: from == to.
// remote replicator: from != to.
type ReplicatorPeer interface {
	// Startup starts wal replicator channel,
	Startup()
	// Shutdown shutdown gracefully.
	Shutdown()
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
func (r replicatorPeer) Startup() {
	if r.running.CAS(false, true) {
		go r.runner.replicaLoop()
	}
}

// Shutdown shutdown gracefully.
func (r replicatorPeer) Shutdown() {
	if r.running.CAS(true, false) {
		r.runner.shutdown()
	}
}

type replicatorRunner struct {
	running    *atomic.Bool
	replicator Replicator

	closed chan struct{}

	logger *logger.Logger
}

func newReplicatorRunner(replicator Replicator) *replicatorRunner {
	return &replicatorRunner{
		replicator: replicator,
		running:    atomic.NewBool(false),
		closed:     make(chan struct{}),
		logger:     logger.GetLogger("replica", "replicatorRunner"),
	}
}

func (r *replicatorRunner) replicaLoop() {
	if r.running.CAS(false, true) {
		r.loop()
	}
}

func (r *replicatorRunner) shutdown() {
	if r.running.CAS(true, false) {
		// wait for stop replica loop
		<-r.closed
	}
}

func (r *replicatorRunner) loop() {
	for r.running.Load() {
		//TODO need handle panic
		hasData := false

		if r.replicator.IsReady() {
			seq := r.replicator.Consume()
			if seq > 0 {
				hasData = true
				data, err := r.replicator.GetMessage(seq)
				if err != nil {
					//TODO add metric
					r.logger.Warn("cannot get replica message data",
						logger.Any("replicator", r.replicator),
						logger.Int64("index", seq))
				} else {
					r.replicator.Replica(seq, data)
				}
			}
		} else {
			r.logger.Warn("replica is not ready", logger.Any("replicator", r.replicator))
		}
		if !hasData {
			//TODO add config?
			time.Sleep(10 * time.Millisecond)
		}
	}

	// exit replica loop
	close(r.closed)
}
