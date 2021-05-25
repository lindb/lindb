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

// NodeStat represents the node monitoring stat
type NodeStat struct {
	Node     ActiveNode `json:"node,omitempty"`
	System   SystemStat `json:"system,omitempty"`
	Replicas int        `json:"replicas"` // the number of replica under the node
	IsDead   bool       `json:"isDead"`
}

// StorageClusterStat represents the storage cluster's stat
type StorageClusterStat struct {
	Name               string           `json:"name,omitempty"`
	Nodes              []*NodeStat      `json:"nodes,omitempty"`
	NodeStatus         NodeStatus       `json:"nodeStatus,omitempty"`
	ReplicaStatus      ReplicaStatus    `json:"replicaStatus,omitempty"`
	Capacity           disk.UsageStat   `json:"capacity,omitempty"`
	DatabaseStatusList []DatabaseStatus `json:"databaseStatusList,omitempty"`
}

// DatabaseStatus represents the database's status
type DatabaseStatus struct {
	Config        Database      `json:"config,omitempty"`
	ReplicaStatus ReplicaStatus `json:"replicaStatus,omitempty"`
}

// NodeStatus represents the status of cluster node
type NodeStatus struct {
	Total   int `json:"total"`
	Alive   int `json:"alive"`
	Suspect int `json:"suspect"`
	Dead    int `json:"dead"`
}

// ReplicaStatus represents the status of replica
type ReplicaStatus struct {
	Total           int `json:"total"`
	UnderReplicated int `json:"underReplicated"`
	Unavailable     int `json:"unavailable"`
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
