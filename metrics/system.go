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

package metrics

import (
	"github.com/lindb/lindb/internal/linmetric"
)

// SystemStatistics represents system statistics.
type SystemStatistics struct {
	// memory
	MemTotal *linmetric.BoundGauge
	MemUsed  *linmetric.BoundGauge
	MemFree  *linmetric.BoundGauge
	MemUsage *linmetric.BoundGauge
	// cpu
	Idle    *linmetric.BoundGauge
	Nice    *linmetric.BoundGauge
	System  *linmetric.BoundGauge
	User    *linmetric.BoundGauge
	Irq     *linmetric.BoundGauge
	Steal   *linmetric.BoundGauge
	SoftRiq *linmetric.BoundGauge
	IOWait  *linmetric.BoundGauge
	// disk usage
	DiskTotal *linmetric.BoundGauge
	DiskUsed  *linmetric.BoundGauge
	DiskFree  *linmetric.BoundGauge
	DiskUsage *linmetric.BoundGauge
	// disk inode
	INodesFree  *linmetric.BoundGauge
	INodesUsed  *linmetric.BoundGauge
	INodesTotal *linmetric.BoundGauge
	INodesUsage *linmetric.BoundGauge
	// net
	NetBytesSent   *linmetric.DeltaCounterVec
	NetBytesRecv   *linmetric.DeltaCounterVec
	NetPacketsSent *linmetric.DeltaCounterVec
	NetPacketsRecv *linmetric.DeltaCounterVec
	NetErrIn       *linmetric.DeltaCounterVec
	NetErrOut      *linmetric.DeltaCounterVec
	NetDropIn      *linmetric.DeltaCounterVec
	NetDropOut     *linmetric.DeltaCounterVec
}

// NewSystemStatistics creates a system statistics.
func NewSystemStatistics(registry *linmetric.Registry) *SystemStatistics {
	scope := registry.NewScope("lindb.monitor.system")
	cpuScope := scope.Scope("cpu_stat")
	memScope := scope.Scope("mem_stat")
	diskScope := scope.Scope("disk_usage_stats")
	inodesScope := scope.Scope("disk_inodes_stats")
	netScope := scope.Scope("net_stat")
	return &SystemStatistics{
		MemTotal: memScope.NewGauge("total"),
		MemUsed:  memScope.NewGauge("used"),
		MemFree:  memScope.NewGauge("free"),
		MemUsage: memScope.NewGauge("usage"),

		Idle:    cpuScope.NewGauge("idle"),
		Nice:    cpuScope.NewGauge("nice"),
		System:  cpuScope.NewGauge("system"),
		User:    cpuScope.NewGauge("user"),
		Irq:     cpuScope.NewGauge("irq"),
		Steal:   cpuScope.NewGauge("steal"),
		SoftRiq: cpuScope.NewGauge("softirq"),
		IOWait:  cpuScope.NewGauge("iowait"),

		DiskTotal: diskScope.NewGauge("total"),
		DiskUsed:  diskScope.NewGauge("used"),
		DiskFree:  diskScope.NewGauge("free"),
		DiskUsage: diskScope.NewGauge("usage"),

		INodesFree:     inodesScope.NewGauge("inodes_free"),
		INodesUsed:     inodesScope.NewGauge("inodes_used"),
		INodesTotal:    inodesScope.NewGauge("inodes_total"),
		INodesUsage:    inodesScope.NewGauge("inodes_usage"),
		NetBytesSent:   netScope.NewCounterVec("bytes_sent", "interface"),
		NetBytesRecv:   netScope.NewCounterVec("bytes_recv", "interface"),
		NetPacketsSent: netScope.NewCounterVec("packets_sent", "interface"),
		NetPacketsRecv: netScope.NewCounterVec("packets_recv", "interface"),
		NetErrIn:       netScope.NewCounterVec("errin", "interface"),
		NetErrOut:      netScope.NewCounterVec("errout", "interface"),
		NetDropIn:      netScope.NewCounterVec("dropin", "interface"),
		NetDropOut:     netScope.NewCounterVec("dropout", "interface"),
	}
}
