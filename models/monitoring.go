package models

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
	Capacity           DiskStat         `json:"capacity,omitempty"`
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
	CPUs       int         `json:"cpus"`                 // number of cpu logic core
	CPUStat    *CPUStat    `json:"cpuStat,omitempty"`    // cpu stat
	MemoryStat *MemoryStat `json:"memoryStat,omitempty"` // memory stat
	DiskStat   *DiskStat   `json:"diskStat,omitempty"`   // disk stat
}

// DiskStat represents the disk usage statistics in system
type DiskStat struct {
	Total       uint64  `json:"total"`
	Used        uint64  `json:"used"`
	UsedPercent float64 `json:"usedPercent"`
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
