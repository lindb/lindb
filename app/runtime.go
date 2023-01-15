package app

import (
	"context"

	"github.com/lindb/lindb/config"
	"github.com/lindb/lindb/internal/linmetric"
	"github.com/lindb/lindb/internal/monitoring"
	"github.com/lindb/lindb/metrics"
	"github.com/lindb/lindb/pkg/logger"
	"github.com/lindb/lindb/series/tag"
)

var (
	newNativeProtoPusher = monitoring.NewNativeProtoPusher
	NewBaseRuntimeFn     = NewBaseRuntime
)

// BaseRuntime represents the common logic of runtime.
type BaseRuntime struct {
	ctx             context.Context
	monitor         config.Monitor
	registry        *linmetric.Registry
	pusher          monitoring.NativePusher
	globalKeyValues tag.Tags

	logger *logger.Logger
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
