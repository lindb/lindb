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

package models

import (
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/mem"
)

// FamilyLogReplicaState represents the family's log replica state.
type FamilyLogReplicaState struct {
	ShardID    ShardID `json:"shardId"`
	FamilyTime string  `json:"familyTime"`
	Leader     NodeID  `json:"leader"`

	Append int64 `json:"append"`

	Replicators []ReplicaPeerState `json:"replicators"`
}

// ReplicaPeerState represents current wal replica peer state.
type ReplicaPeerState struct {
	Replicator string `json:"replicator"`
	Consume    int64  `json:"consume"`
	ACK        int64  `json:"ack"`
	Pending    int64  `json:"pending"`
}

// SystemStat represents the system statistics
type SystemStat struct {
	CPUs          int                    `json:"cpus"`                    // number of cpu logic core
	CPUStat       *CPUStat               `json:"cpuStat,omitempty"`       // cpu stat
	MemoryStat    *mem.VirtualMemoryStat `json:"memoryStat,omitempty"`    // memory stat
	DiskUsageStat *disk.UsageStat        `json:"diskUsageStat,omitempty"` // disk usage stat
}

// MemoryStat represents the memory usage statistics in system
type MemoryStat struct {
	// Total amount of RAM on this system
	Total uint64 `json:"total"`
	// RAM used by programs
	// This value is computed from the kernel specific values.
	Used uint64 `json:"used"`
	// Percentage of RAM used by programs
	// This value is computed from the kernel specific values.
	UsedPercent float64 `json:"usedPercent"`
}

// CPUStat represents the amounts of time the CPU has spent performing different
// kinds of work.
type CPUStat struct {
	User    float64 `json:"user"`
	System  float64 `json:"system"`
	Idle    float64 `json:"idle"`
	Nice    float64 `json:"nice"`
	Iowait  float64 `json:"iowait"`
	Irq     float64 `json:"irq"`
	Softirq float64 `json:"softirq"`
	Steal   float64 `json:"steal"`
}
