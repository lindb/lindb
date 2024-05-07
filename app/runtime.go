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

package app

import (
	"context"

	"github.com/lindb/common/pkg/logger"

	"github.com/lindb/lindb/config"
	"github.com/lindb/lindb/internal/linmetric"
	"github.com/lindb/lindb/internal/monitoring"
	"github.com/lindb/lindb/metrics"
	"github.com/lindb/lindb/series/tag"
)

var (
	newNativeProtoPusher = monitoring.NewNativeProtoPusher
	NewBaseRuntimeFn     = NewBaseRuntime
)

// BaseRuntime represents the common logic of runtime.
type BaseRuntime struct {
	ctx             context.Context
	pusher          monitoring.NativePusher
	logger          logger.Logger
	registry        *linmetric.Registry
	monitor         config.Monitor
	globalKeyValues tag.Tags
}

// NewBaseRuntime creates a base runtime instance.
func NewBaseRuntime(ctx context.Context, monitor config.Monitor, registry *linmetric.Registry, globalKeyValues tag.Tags) BaseRuntime {
	return BaseRuntime{
		ctx:             ctx,
		monitor:         monitor,
		registry:        registry,
		globalKeyValues: globalKeyValues,
		logger:          logger.GetLogger("Base", "Runtime"),
	}
}

// Shutdown stops the resource of base runtime.
func (r *BaseRuntime) Shutdown() {
	if r.pusher != nil {
		r.pusher.Stop()
		r.logger.Info("stopped native metric pusher successfully")
	}
}

// NativePusher pushes metric data into internal database.
func (r *BaseRuntime) NativePusher() {
	monitorEnabled := r.monitor.ReportInterval > 0
	if !monitorEnabled {
		r.logger.Info("pusher won't start because report-interval is 0")
		return
	}
	r.logger.Info("pusher is running",
		logger.String("interval", r.monitor.ReportInterval.String()))

	r.pusher = newNativeProtoPusher(
		r.ctx,
		r.monitor.URL,
		r.monitor.ReportInterval.Duration(),
		r.monitor.PushTimeout.Duration(),
		r.registry,
		r.globalKeyValues,
	)
	go r.pusher.Start()
}

// SystemCollector collects the system metric.
func (r *BaseRuntime) SystemCollector() {
	r.logger.Info("system collector is running")

	go monitoring.NewSystemCollector(
		r.ctx,
		"/",
		metrics.NewSystemStatistics(r.registry)).Run()
}
