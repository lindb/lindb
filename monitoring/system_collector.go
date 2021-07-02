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

package monitoring

import (
	"context"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/mem"
	"github.com/shirou/gopsutil/net"

	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/encoding"
	"github.com/lindb/lindb/pkg/fileutil"
	"github.com/lindb/lindb/pkg/logger"
	"github.com/lindb/lindb/pkg/state"
)

var (
	memStatGauge = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "system_mem_stat",
		Help: "System mem stats",
	}, []string{"type"})
	cpuStatGauge = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "system_cpu_stat",
		Help: "System cpu stats",
	}, []string{"type"})
	diskUsageStatGauge = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "system_disk_usage_stat",
		Help: "System disk usage stats",
	}, []string{"type"})
	netStatCounter = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "net_stat",
		Help: "Network stat",
	}, []string{"type", "interface"})
)

func init() {
	BrokerRegistry.MustRegister(memStatGauge)
	BrokerRegistry.MustRegister(cpuStatGauge)
	BrokerRegistry.MustRegister(diskUsageStatGauge)
	BrokerRegistry.MustRegister(netStatCounter)
	BrokerRegistry.MustRegister(prometheus.NewProcessCollector(prometheus.ProcessCollectorOpts{}))
	BrokerRegistry.MustRegister(prometheus.NewGoCollector())

	StorageRegistry.MustRegister(memStatGauge)
	StorageRegistry.MustRegister(cpuStatGauge)
	StorageRegistry.MustRegister(diskUsageStatGauge)
	StorageRegistry.MustRegister(netStatCounter)
	StorageRegistry.MustRegister(prometheus.NewProcessCollector(prometheus.ProcessCollectorOpts{}))
	StorageRegistry.MustRegister(prometheus.NewGoCollector())
}

// SystemCollector collects the system stat
type SystemCollector struct {
	ctx             context.Context
	interval        time.Duration
	storage         string
	path            string                        // repository key
	netStats        map[string]net.IOCountersStat // interface-name as key
	netStatsUpdated map[string]time.Time          // last updated time
	systemStat      *models.SystemStat
	nodeStat        *models.NodeStat
	repository      state.Repository
	// used for mock
	MemoryStatGetter    MemoryStatGetter
	CPUStatGetter       CPUStatGetter
	DiskUsageStatGetter DiskUsageStatGetter
	NetStatGetter       NetStatGetter
}

// NewSystemCollector creates a new system stat collector
func NewSystemCollector(
	ctx context.Context,
	interval time.Duration,
	storage string,
	repository state.Repository,
	path string,
	node models.ActiveNode,
) *SystemCollector {
	r := &SystemCollector{
		interval:        interval,
		storage:         fileutil.GetExistPath(storage),
		repository:      repository,
		path:            path,
		netStats:        make(map[string]net.IOCountersStat),
		netStatsUpdated: make(map[string]time.Time),
		systemStat:      &models.SystemStat{},
		nodeStat: &models.NodeStat{
			Node: node,
		},
		ctx:                 ctx,
		MemoryStatGetter:    mem.VirtualMemory,
		CPUStatGetter:       GetCPUStat,
		DiskUsageStatGetter: disk.UsageWithContext,
		NetStatGetter:       GetNetStat,
	}
	return r
}

// Run starts a background goroutine that collects the monitoring stat
func (r *SystemCollector) Run() {
	ticker := time.NewTicker(r.interval)
	defer ticker.Stop()
	// collect system status
	r.collect()

	for {
		select {
		case <-ticker.C:
			// collect system status
			r.collect()
		case <-r.ctx.Done():
			return
		}
	}
}

// collect collects the monitoring stat
func (r *SystemCollector) collect() {
	var err error
	r.systemStat.CPUs = GetCPUs()

	if r.systemStat.MemoryStat, err = r.MemoryStatGetter(); err != nil {
		collectorLogger.Error("get memory stat", logger.Error(err))
	}
	if r.systemStat.CPUStat, err = r.CPUStatGetter(); err != nil {
		collectorLogger.Error("get cpu stat", logger.Error(err))
	}
	if r.systemStat.DiskUsageStat, err = r.DiskUsageStatGetter(r.ctx, r.storage); err != nil {
		collectorLogger.Error("get disk usage stat", logger.Error(err))
	}
	if stats, err := r.NetStatGetter(r.ctx); err != nil {
		collectorLogger.Error("get net stat", logger.Error(err))
	} else {
		for _, stat := range stats {
			r.netStats[stat.Name] = stat
			r.netStatsUpdated[stat.Name] = time.Now()
		}
	}

	r.nodeStat.System = *r.systemStat

	r.logMemStat()
	r.logDiskUsageStat()
	r.logCPUStat()
	r.logNetStat()

	if err := r.repository.Put(r.ctx, r.path, encoding.JSONMarshal(r.nodeStat)); err != nil {
		collectorLogger.Error("report stat error", logger.String("path", r.path))
	}
}

func (r *SystemCollector) logMemStat() {
	if r.systemStat.MemoryStat != nil {
		memStat := r.systemStat.MemoryStat
		memStatGauge.WithLabelValues("total").Set(float64(memStat.Total))
		memStatGauge.WithLabelValues("used").Set(float64(memStat.Used))
		memStatGauge.WithLabelValues("used_percent").Set(memStat.UsedPercent)
	}
}

func (r *SystemCollector) logCPUStat() {
	if r.systemStat.CPUStat != nil {
		cpuStat := r.systemStat.CPUStat
		cpuStatGauge.WithLabelValues("idle").Set(cpuStat.Idle)
		cpuStatGauge.WithLabelValues("nice").Set(cpuStat.Nice)
		cpuStatGauge.WithLabelValues("system").Set(cpuStat.System)
		cpuStatGauge.WithLabelValues("user").Set(cpuStat.User)
		cpuStatGauge.WithLabelValues("irq").Set(cpuStat.Irq)
		cpuStatGauge.WithLabelValues("steal").Set(cpuStat.Steal)
		cpuStatGauge.WithLabelValues("softirq").Set(cpuStat.Softirq)
		cpuStatGauge.WithLabelValues("iowait").Set(cpuStat.Iowait)
	}
}

func (r *SystemCollector) logDiskUsageStat() {
	if r.systemStat.DiskUsageStat != nil {
		stat := r.systemStat.DiskUsageStat
		// usage
		diskUsageStatGauge.WithLabelValues("total").Set(float64(stat.Total))
		diskUsageStatGauge.WithLabelValues("used").Set(float64(stat.Used))
		diskUsageStatGauge.WithLabelValues("free").Set(float64(stat.Free))
		diskUsageStatGauge.WithLabelValues("used_percent").Set(stat.UsedPercent)
		// inode
		diskUsageStatGauge.WithLabelValues("inodes_free").Set(float64(stat.InodesFree))
		diskUsageStatGauge.WithLabelValues("inodes_used").Set(float64(stat.InodesUsed))
		diskUsageStatGauge.WithLabelValues("inodes_total").Set(float64(stat.InodesTotal))
		diskUsageStatGauge.WithLabelValues("inodes_used_percent").Set(stat.InodesUsedPercent)
	}
}
func (r *SystemCollector) logNetStat() {
	for _, stat := range r.netStats {
		lastStat, ok := r.netStats[stat.Name]
		// check time interval
		if ok && time.Since(r.netStatsUpdated[stat.Name]) <= 2*r.interval {
			netStatCounter.WithLabelValues("bytes_sent", stat.Name).Add(float64(stat.BytesSent - lastStat.BytesSent))
			netStatCounter.WithLabelValues("bytes_recv", stat.Name).Add(float64(stat.BytesRecv - lastStat.BytesRecv))
			netStatCounter.WithLabelValues("packets_sent", stat.Name).Add(float64(stat.PacketsSent - lastStat.PacketsSent))
			netStatCounter.WithLabelValues("packets_recv", stat.Name).Add(float64(stat.PacketsRecv - lastStat.PacketsRecv))
			netStatCounter.WithLabelValues("errin", stat.Name).Add(float64(stat.Errin - lastStat.Errin))
			netStatCounter.WithLabelValues("errout", stat.Name).Add(float64(stat.Errout - lastStat.Errout))
			netStatCounter.WithLabelValues("dropin", stat.Name).Add(float64(stat.Dropin - lastStat.Dropin))
			netStatCounter.WithLabelValues("dropout", stat.Name).Add(float64(stat.Dropout - lastStat.Dropout))
		}
		r.netStats[stat.Name] = stat
		r.netStatsUpdated[stat.Name] = time.Now()
	}
}
