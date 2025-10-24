package base

import "github.com/lindb/lindb/pkg/timeutil"

type Partition struct {
	Dir       string
	Timestamp int64

	Interval timeutil.Interval
}

func (p *Partition) Path() string {
	return p.Dir
}

func (p *Partition) PartitionTime() int64 {
	return p.Timestamp
}

func (p *Partition) PartitionInterval() timeutil.Interval {
	return p.Interval
}
