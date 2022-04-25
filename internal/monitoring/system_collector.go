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

	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/mem"
	"github.com/shirou/gopsutil/net"

	"github.com/lindb/lindb/metrics"
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/fileutil"
	"github.com/lindb/lindb/pkg/logger"
)

var sc *SystemCollector

// SystemCollector collects the system stat
type SystemCollector struct {
	ctx             context.Context
	interval        time.Duration
	storage         string
	netStats        map[string]net.IOCountersStat // interface-name as key
	netStatsUpdated map[string]time.Time          // last updated time
	systemStat      *models.SystemStat
	// used for mock
	MemoryStatGetter    MemoryStatGetter
	CPUStatGetter       CPUStatGetter
	DiskUsageStatGetter DiskUsageStatGetter
	NetStatGetter       NetStatGetter

	statistics *metrics.SystemStatistics
}

// NewSystemCollector creates a new system stat collector
func NewSystemCollector(
	ctx context.Context,
	storage string,
	statistics *metrics.SystemStatistics,
) *SystemCollector {
	sc = &SystemCollector{
		interval:            time.Second * 10,
		netStats:            make(map[string]net.IOCountersStat),
		netStatsUpdated:     make(map[string]time.Time),
		systemStat:          &models.SystemStat{},
		ctx:                 ctx,
		MemoryStatGetter:    mem.VirtualMemory,
		CPUStatGetter:       GetCPUStat,
		DiskUsageStatGetter: disk.UsageWithContext,
		NetStatGetter:       GetNetStat,
		statistics:          statistics,
	}
	if storage != "" {
		sc.storage = fileutil.GetExistPath(storage)
	}
	return sc
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
	if r.storage != "" {
		if r.systemStat.DiskUsageStat, err = r.DiskUsageStatGetter(r.ctx, r.storage); err != nil {
			collectorLogger.Error("get disk usage stat", logger.Error(err))
		}
	}
	if stats, err := r.NetStatGetter(r.ctx); err != nil {
		collectorLogger.Error("get net stat", logger.Error(err))
	} else {
		for _, stat := range stats {
			r.netStats[stat.Name] = stat
			r.netStatsUpdated[stat.Name] = time.Now()
		}
	}

	r.logMemStat()
	r.logDiskUsageStat()
	r.logCPUStat()
	r.logNetStat()
}

func (r *SystemCollector) logMemStat() {
	if r.systemStat.MemoryStat != nil {
		memStat := r.systemStat.MemoryStat
		r.statistics.MemTotal.Update(float64(memStat.Total))
		r.statistics.MemUsed.Update(float64(memStat.Used))
		r.statistics.MemFree.Update(float64(memStat.Free))
		r.statistics.MemUsage.Update(memStat.UsedPercent)
	}
}

func (r *SystemCollector) logCPUStat() {
	if r.systemStat.CPUStat != nil {
		cpuStat := r.systemStat.CPUStat
		r.statistics.Idle.Update(cpuStat.Idle)
		r.statistics.Nice.Update(cpuStat.Nice)
		r.statistics.System.Update(cpuStat.System)
		r.statistics.User.Update(cpuStat.User)
		r.statistics.Irq.Update(cpuStat.Irq)
		r.statistics.Steal.Update(cpuStat.Steal)
		r.statistics.SoftIrq.Update(cpuStat.Softirq)
		r.statistics.IOWait.Update(cpuStat.Iowait)
	}
}

func (r *SystemCollector) logDiskUsageStat() {
	if r.systemStat.DiskUsageStat != nil {
		stat := r.systemStat.DiskUsageStat
		// usage
		r.statistics.DiskTotal.Update(float64(stat.Total))
		r.statistics.DiskUsed.Update(float64(stat.Used))
		r.statistics.DiskFree.Update(float64(stat.Free))
		r.statistics.DiskUsed.Update(stat.UsedPercent)
		// inode
		r.statistics.INodesFree.Update(float64(stat.InodesFree))
		r.statistics.INodesUsed.Update(float64(stat.InodesUsed))
		r.statistics.INodesTotal.Update(float64(stat.InodesTotal))
		r.statistics.INodesUsage.Update(stat.InodesUsedPercent)
	}
}

func (r *SystemCollector) logNetStat() {
	for _, stat := range r.netStats {
		lastStat, ok := r.netStats[stat.Name]
		// check time interval
		if ok && time.Since(r.netStatsUpdated[stat.Name]) <= 2*r.interval {
			r.statistics.NetBytesSent.WithTagValues(stat.Name).Add(float64(stat.BytesSent - lastStat.BytesSent))
			r.statistics.NetBytesRecv.WithTagValues(stat.Name).Add(float64(stat.BytesRecv - lastStat.BytesRecv))
			r.statistics.NetPacketsSent.WithTagValues(stat.Name).Add(float64(stat.PacketsSent - lastStat.PacketsSent))
			r.statistics.NetPacketsRecv.WithTagValues(stat.Name).Add(float64(stat.PacketsRecv - lastStat.PacketsRecv))
			r.statistics.NetErrIn.WithTagValues(stat.Name).Add(float64(stat.Errin - lastStat.Errin))
			r.statistics.NetErrOut.WithTagValues(stat.Name).Add(float64(stat.Errout - lastStat.Errout))
			r.statistics.NetDropIn.WithTagValues(stat.Name).Add(float64(stat.Dropin - lastStat.Dropin))
			r.statistics.NetDropOut.WithTagValues(stat.Name).Add(float64(stat.Dropout - lastStat.Dropout))
		}
		r.netStats[stat.Name] = stat
		r.netStatsUpdated[stat.Name] = time.Now()
	}
}
