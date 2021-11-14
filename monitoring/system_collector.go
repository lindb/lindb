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

	"github.com/lindb/lindb/internal/linmetric"
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
	nodeStat        *models.NodeStat
	// used for mock
	MemoryStatGetter    MemoryStatGetter
	CPUStatGetter       CPUStatGetter
	DiskUsageStatGetter DiskUsageStatGetter
	NetStatGetter       NetStatGetter

	//  role symbols this collector is owned by storage or broker runtime
	role string
	// metrics
	memTotalGauge       *linmetric.BoundGauge
	memUsedGauge        *linmetric.BoundGauge
	memUsedPercentGauge *linmetric.BoundGauge
	// cpu
	cpuIdleGauge    *linmetric.BoundGauge
	cpuNiceGauge    *linmetric.BoundGauge
	cpuSystemGauge  *linmetric.BoundGauge
	cpuUserGauge    *linmetric.BoundGauge
	cpuIrqGauge     *linmetric.BoundGauge
	cpuStealGauge   *linmetric.BoundGauge
	cpuSoftRiqGauge *linmetric.BoundGauge
	cpuIOWaitGauge  *linmetric.BoundGauge
	// disk usage
	diskTotalGauge       *linmetric.BoundGauge
	diskUsedGauge        *linmetric.BoundGauge
	diskFreeGauge        *linmetric.BoundGauge
	diskUsedPercentGauge *linmetric.BoundGauge
	// disk inode
	inodesFreeGauge        *linmetric.BoundGauge
	inodesUsedGauge        *linmetric.BoundGauge
	inodesTotalGauge       *linmetric.BoundGauge
	inodesUsedPercentGauge *linmetric.BoundGauge
	// net
	bytesSentCounterVec   *linmetric.DeltaCounterVec
	bytesRecvCounterVec   *linmetric.DeltaCounterVec
	packetsSentCounterVec *linmetric.DeltaCounterVec
	packetsRecvCounterVec *linmetric.DeltaCounterVec
	errInCounterVec       *linmetric.DeltaCounterVec
	errOutCounterVec      *linmetric.DeltaCounterVec
	dropInCounterVec      *linmetric.DeltaCounterVec
	dropOutCounterVec     *linmetric.DeltaCounterVec
}

// NewSystemCollector creates a new system stat collector
func NewSystemCollector(
	ctx context.Context,
	storage string,
	node *models.StatelessNode,
	role string,
) *SystemCollector {
	sc = &SystemCollector{
		interval:        time.Second * 10,
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
		role:                role,
	}
	if storage != "" {
		sc.storage = fileutil.GetExistPath(storage)
	}
	sc.boundMetrics()
	return sc
}

func (r *SystemCollector) boundMetrics() {
	systemScope := linmetric.NewScope("lindb.monitor.system", "role", r.role)

	systemMemScope := systemScope.Scope("mem_stat")
	// memory
	r.memTotalGauge = systemMemScope.NewGauge("total")
	r.memUsedGauge = systemMemScope.NewGauge("used")
	r.memUsedPercentGauge = systemMemScope.NewGauge("used_percent")

	systemCPUScope := systemScope.Scope("cpu_stat")
	// cpu
	r.cpuIdleGauge = systemCPUScope.NewGauge("idle")
	r.cpuNiceGauge = systemCPUScope.NewGauge("nice")
	r.cpuSystemGauge = systemCPUScope.NewGauge("system")
	r.cpuUserGauge = systemCPUScope.NewGauge("user")
	r.cpuIrqGauge = systemCPUScope.NewGauge("irq")
	r.cpuStealGauge = systemCPUScope.NewGauge("steal")
	r.cpuSoftRiqGauge = systemCPUScope.NewGauge("softirq")
	r.cpuIOWaitGauge = systemCPUScope.NewGauge("iowait")

	systemDiskScope := systemScope.Scope("disk_usage_stats")
	// disk usage
	r.diskTotalGauge = systemDiskScope.NewGauge("total")
	r.diskUsedGauge = systemDiskScope.NewGauge("used")
	r.diskFreeGauge = systemDiskScope.NewGauge("free")
	r.diskUsedPercentGauge = systemDiskScope.NewGauge("used_percent")

	systemInodesScope := systemScope.Scope("disk_inodes_stats")
	// disk inode
	r.inodesFreeGauge = systemInodesScope.NewGauge("inodes_free")
	r.inodesUsedGauge = systemInodesScope.NewGauge("inodes_used")
	r.inodesTotalGauge = systemInodesScope.NewGauge("inodes_total")
	r.inodesUsedPercentGauge = systemInodesScope.NewGauge("inodes_used_percent")

	netScope := systemScope.Scope("net_stat")
	// net
	r.bytesSentCounterVec = netScope.NewCounterVec("bytes_sent", "interface")
	r.bytesRecvCounterVec = netScope.NewCounterVec("bytes_recv", "interface")
	r.packetsSentCounterVec = netScope.NewCounterVec("packets_sent", "interface")
	r.packetsRecvCounterVec = netScope.NewCounterVec("packets_recv", "interface")
	r.errInCounterVec = netScope.NewCounterVec("errin", "interface")
	r.errOutCounterVec = netScope.NewCounterVec("errout", "interface")
	r.dropInCounterVec = netScope.NewCounterVec("dropin", "interface")
	r.dropOutCounterVec = netScope.NewCounterVec("dropout", "interface")
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

	r.nodeStat.System = *r.systemStat

	r.logMemStat()
	r.logDiskUsageStat()
	r.logCPUStat()
	r.logNetStat()
}

func (r *SystemCollector) logMemStat() {
	if r.systemStat.MemoryStat != nil {
		memStat := r.systemStat.MemoryStat
		r.memTotalGauge.Update(float64(memStat.Total))
		r.memUsedGauge.Update(float64(memStat.Used))
		r.memUsedPercentGauge.Update(memStat.UsedPercent)
	}
}

func (r *SystemCollector) logCPUStat() {
	if r.systemStat.CPUStat != nil {
		cpuStat := r.systemStat.CPUStat
		r.cpuIdleGauge.Update(cpuStat.Idle)
		r.cpuNiceGauge.Update(cpuStat.Nice)
		r.cpuSystemGauge.Update(cpuStat.System)
		r.cpuUserGauge.Update(cpuStat.User)
		r.cpuIrqGauge.Update(cpuStat.Irq)
		r.cpuStealGauge.Update(cpuStat.Steal)
		r.cpuSoftRiqGauge.Update(cpuStat.Softirq)
		r.cpuIOWaitGauge.Update(cpuStat.Iowait)
	}
}

func (r *SystemCollector) logDiskUsageStat() {
	if r.systemStat.DiskUsageStat != nil {
		stat := r.systemStat.DiskUsageStat
		// usage
		r.diskTotalGauge.Update(float64(stat.Total))
		r.diskUsedGauge.Update(float64(stat.Used))
		r.diskFreeGauge.Update(float64(stat.Free))
		r.diskUsedPercentGauge.Update(stat.UsedPercent)
		// inode
		r.inodesFreeGauge.Update(float64(stat.InodesFree))
		r.inodesUsedGauge.Update(float64(stat.InodesUsed))
		r.inodesTotalGauge.Update(float64(stat.InodesTotal))
		r.inodesUsedPercentGauge.Update(stat.InodesUsedPercent)
	}
}
func (r *SystemCollector) logNetStat() {
	for _, stat := range r.netStats {
		lastStat, ok := r.netStats[stat.Name]
		// check time interval
		if ok && time.Since(r.netStatsUpdated[stat.Name]) <= 2*r.interval {
			r.bytesSentCounterVec.WithTagValues(stat.Name).Add(float64(stat.BytesSent - lastStat.BytesSent))
			r.bytesRecvCounterVec.WithTagValues(stat.Name).Add(float64(stat.BytesRecv - lastStat.BytesRecv))
			r.packetsSentCounterVec.WithTagValues(stat.Name).Add(float64(stat.PacketsSent - lastStat.PacketsSent))
			r.packetsRecvCounterVec.WithTagValues(stat.Name).Add(float64(stat.PacketsRecv - lastStat.PacketsRecv))
			r.errInCounterVec.WithTagValues(stat.Name).Add(float64(stat.Errin - lastStat.Errin))
			r.errOutCounterVec.WithTagValues(stat.Name).Add(float64(stat.Errout - lastStat.Errout))
			r.dropInCounterVec.WithTagValues(stat.Name).Add(float64(stat.Dropin - lastStat.Dropin))
			r.dropOutCounterVec.WithTagValues(stat.Name).Add(float64(stat.Dropout - lastStat.Dropout))
		}
		r.netStats[stat.Name] = stat
		r.netStatsUpdated[stat.Name] = time.Now()
	}
}
