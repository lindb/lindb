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

	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/encoding"
	"github.com/lindb/lindb/pkg/fileutil"
	"github.com/lindb/lindb/pkg/logger"
	"github.com/lindb/lindb/pkg/state"
)

var (
	// memory
	systemMemScope      = monitorScope.Scope("system_mem_stat")
	memTotalGauge       = systemMemScope.NewGauge("total")
	memUsedGauge        = systemMemScope.NewGauge("used")
	memUsedPercentGauge = systemMemScope.NewGauge("used_percent")
	// cpu
	systemCPUScope  = monitorScope.Scope("system_cpu_stat")
	cpuIdleGauge    = systemCPUScope.NewGauge("idle")
	cpuNiceGauge    = systemCPUScope.NewGauge("nice")
	cpuSystemGauge  = systemCPUScope.NewGauge("system")
	cpuUserGauge    = systemCPUScope.NewGauge("user")
	cpuIrqGauge     = systemCPUScope.NewGauge("irq")
	cpuStealGauge   = systemCPUScope.NewGauge("steal")
	cpuSoftRiqGauge = systemCPUScope.NewGauge("softirq")
	cpuIOWaitGauge  = systemCPUScope.NewGauge("iowait")
	// disk usage
	systemDiskScope      = monitorScope.Scope("system_disk_usage_stats")
	diskTotalGauge       = systemDiskScope.NewGauge("total")
	diskUsedGauge        = systemDiskScope.NewGauge("used")
	diskFreeGauge        = systemDiskScope.NewGauge("free")
	diskUsedPercentGauge = systemDiskScope.NewGauge("used_percent")
	// disk inode
	inodesFreeGauge        = systemCPUScope.NewGauge("inodes_free")
	inodesUsedGauge        = systemCPUScope.NewGauge("inodes_used")
	inodesTotalGauge       = systemCPUScope.NewGauge("inodes_total")
	inodesUsedPercentGauge = systemCPUScope.NewGauge("inodes_used_percent")
	// net
	netScope              = monitorScope.Scope("system_net_stat")
	bytesSentCounterVec   = netScope.NewDeltaCounterVec("bytes_sent", "interface")
	bytesRecvCounterVec   = netScope.NewDeltaCounterVec("bytes_recv", "interface")
	packetsSentCounterVec = netScope.NewDeltaCounterVec("packets_sent", "interface")
	packetsRecvCounterVec = netScope.NewDeltaCounterVec("packets_recv", "interface")
	errInCounterVec       = netScope.NewDeltaCounterVec("errin", "interface")
	errOutCounterVec      = netScope.NewDeltaCounterVec("errout", "interface")
	dropInCounterVec      = netScope.NewDeltaCounterVec("dropin", "interface")
	dropOutCounterVec     = netScope.NewDeltaCounterVec("dropout", "interface")
)

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
		memTotalGauge.Update(float64(memStat.Total))
		memUsedGauge.Update(float64(memStat.Used))
		memUsedPercentGauge.Update(memStat.UsedPercent)
	}
}

func (r *SystemCollector) logCPUStat() {
	if r.systemStat.CPUStat != nil {
		cpuStat := r.systemStat.CPUStat
		cpuIdleGauge.Update(cpuStat.Idle)
		cpuNiceGauge.Update(cpuStat.Nice)
		cpuSystemGauge.Update(cpuStat.System)
		cpuUserGauge.Update(cpuStat.User)
		cpuIrqGauge.Update(cpuStat.Irq)
		cpuStealGauge.Update(cpuStat.Steal)
		cpuSoftRiqGauge.Update(cpuStat.Softirq)
		cpuIOWaitGauge.Update(cpuStat.Iowait)
	}
}

func (r *SystemCollector) logDiskUsageStat() {
	if r.systemStat.DiskUsageStat != nil {
		stat := r.systemStat.DiskUsageStat
		// usage
		diskTotalGauge.Update(float64(stat.Total))
		diskUsedGauge.Update(float64(stat.Used))
		diskFreeGauge.Update(float64(stat.Free))
		diskUsedPercentGauge.Update(stat.UsedPercent)
		// inode
		inodesFreeGauge.Update(float64(stat.InodesFree))
		inodesUsedGauge.Update(float64(stat.InodesUsed))
		inodesTotalGauge.Update(float64(stat.InodesTotal))
		inodesUsedPercentGauge.Update(stat.InodesUsedPercent)
	}
}
func (r *SystemCollector) logNetStat() {
	for _, stat := range r.netStats {
		lastStat, ok := r.netStats[stat.Name]
		// check time interval
		if ok && time.Since(r.netStatsUpdated[stat.Name]) <= 2*r.interval {
			bytesSentCounterVec.WithTagValues(stat.Name).Add(float64(stat.BytesSent - lastStat.BytesSent))
			bytesRecvCounterVec.WithTagValues(stat.Name).Add(float64(stat.BytesRecv - lastStat.BytesRecv))
			packetsSentCounterVec.WithTagValues(stat.Name).Add(float64(stat.PacketsSent - lastStat.PacketsSent))
			packetsRecvCounterVec.WithTagValues(stat.Name).Add(float64(stat.PacketsRecv - lastStat.PacketsRecv))
			errInCounterVec.WithTagValues(stat.Name).Add(float64(stat.Errin - lastStat.Errin))
			errOutCounterVec.WithTagValues(stat.Name).Add(float64(stat.Errout - lastStat.Errout))
			dropInCounterVec.WithTagValues(stat.Name).Add(float64(stat.Dropin - lastStat.Dropin))
			dropOutCounterVec.WithTagValues(stat.Name).Add(float64(stat.Dropout - lastStat.Dropout))
		}
		r.netStats[stat.Name] = stat
		r.netStatsUpdated[stat.Name] = time.Now()
	}
}
