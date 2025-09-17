package base

import (
	"sync"
)

type Partition struct {
	Dir       string
	Timestamp int64

	mutex sync.Mutex
}

func (p *Partition) Path() string {
	return p.Dir
}

func (p *Partition) PartitionTime() int64 {
	return p.Timestamp
}
