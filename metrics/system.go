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
	MemTotal *linmetric.BoundGauge // Total amount of RAM on this system
	MemUsed  *linmetric.BoundGauge // RAM used by programs
	MemFree  *linmetric.BoundGauge // Free RAM
	MemUsage *linmetric.BoundGauge // Percentage of RAM used by programs
	// cpu
	Idle    *linmetric.BoundGauge // CPU time that's not actively being used
	Nice    *linmetric.BoundGauge // CPU time used by processes that have a positive niceness
	System  *linmetric.BoundGauge // CPU time used by the kernel
	User    *linmetric.BoundGauge // CPU time used by user space processes
	Irq     *linmetric.BoundGauge // Interrupt Requests
	Steal   *linmetric.BoundGauge // The percentage of time a virtual CPU waits for a real CPU
	SoftIrq *linmetric.BoundGauge // The kernel is servicing interrupt requests (IRQs)
	IOWait  *linmetric.BoundGauge // It marks time spent waiting for input or output operations
	// disk usage
	DiskTotal *linmetric.BoundGauge // Total amount of disk
	DiskUsed  *linmetric.BoundGauge // Disk used by programs
	DiskFree  *linmetric.BoundGauge // Free disk
	DiskUsage *linmetric.BoundGauge // Percentage of disk used by programs
	// disk inode
	INodesTotal *linmetric.BoundGauge // Total amount of inode
	INodesUsed  *linmetric.BoundGauge // INode used by programs
	INodesFree  *linmetric.BoundGauge // Free inode
	INodesUsage *linmetric.BoundGauge // Percentage of inode used by programs
	// net
	NetBytesSent   *linmetric.DeltaCounterVec // number of bytes sent
	NetBytesRecv   *linmetric.DeltaCounterVec // number of bytes received
	NetPacketsSent *linmetric.DeltaCounterVec // number of packets sent
	NetPacketsRecv *linmetric.DeltaCounterVec // number of packets received
	NetErrIn       *linmetric.DeltaCounterVec // total number of errors while receiving
	NetErrOut      *linmetric.DeltaCounterVec // total number of errors while sending
	NetDropIn      *linmetric.DeltaCounterVec // total number of incoming packets which were dropped
	NetDropOut     *linmetric.DeltaCounterVec // total number of outgoing packets which were dropped (always 0 on OSX and BSD)
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
		SoftIrq: cpuScope.NewGauge("softirq"),
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
