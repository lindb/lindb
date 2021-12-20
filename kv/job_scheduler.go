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

package kv

import (
	"context"
	"time"

	"github.com/lindb/lindb/pkg/logger"

	"go.uber.org/atomic"
)

// JobScheduler represents a background compaction job scheduler.
type JobScheduler interface {
	// Startup starts the job scheduler.
	Startup()
	// Shutdown stops the job scheduler.
	Shutdown()
	// IsRunning returns the scheduler if running.
	IsRunning() bool
}

// jobScheduler implements JobScheduler interface.
type jobScheduler struct {
	ctx    context.Context
	cancel context.CancelFunc

	option  StoreOptions
	running *atomic.Bool

	logger *logger.Logger
}

// NewJobScheduler creates a JobScheduler instance.
func NewJobScheduler(ctx context.Context, option StoreOptions) JobScheduler {
	ctx, cancel := context.WithCancel(ctx)
	return &jobScheduler{
		option:  option,
		ctx:     ctx,
		cancel:  cancel,
		running: atomic.NewBool(false),
		logger:  logger.GetLogger("KV", "JobScheduler"),
	}
}

// Startup starts the job scheduler.
func (js *jobScheduler) Startup() {
	if js.running.CAS(false, true) {
		js.schedule()
	}
}

// Shutdown stops the job scheduler.
func (js *jobScheduler) Shutdown() {
	if js.running.CAS(true, false) {
		js.cancel()
	}
}

// IsRunning returns the scheduler if running.
func (js *jobScheduler) IsRunning() bool {
	return js.running.Load()
}

// schedule a compaction background job.
// 1. check if it needs to do compact or rollup.
// 2. if it needs, start new goroutine does compact or rollup job.
func (js *jobScheduler) schedule() {
	interval := defaultCompactCheckInterval
	if js.option.CompactCheckInterval > 0 {
		interval = js.option.CompactCheckInterval
	}
	ticker := time.NewTicker(time.Duration(interval) * time.Second)
	go func() {
		for {
			select {
			case <-ticker.C:
				stores := GetStoreManager().GetStores()
				for idx := range stores {
					store := stores[idx]
					// schedule compact if it needs
					store.compact()
				}
			case <-js.ctx.Done():
				ticker.Stop()
				js.logger.Info("job scheduler exit......")
				return
			}
		}
	}()
}
