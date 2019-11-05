package models

// NodeStat represents the node monitoring stat
type NodeStat struct {
	Node   ActiveNode `json:"node,omitempty"`
	System SystemStat `json:"system,omitempty"`
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
	User   float64 `json:"user"`
	System float64 `json:"system"`
	Idle   float64 `json:"idle"`
	Nice   float64 `json:"nice"`
}
