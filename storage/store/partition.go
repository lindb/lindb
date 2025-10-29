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

package store

import (
	"io"

	"github.com/lindb/common/pkg/fileutil"
	"github.com/samber/lo"

	"github.com/lindb/lindb/pkg/timeutil"
)

type Partition interface {
	io.Closer

	PartitionTime() int64
	PartitionInterval() timeutil.Interval

	Path() string
	Shard() Shard

	GetOrCreateSegment(timestamp int64) (Segment, error)

	GetSegments(timeRange timeutil.TimeRange) []Segment
}

type Partitions struct {
	interval   timeutil.Interval
	partitions map[int64]*LazyPartition // partition timestamp -> partition
}

func NewPartitions(interval timeutil.Interval) *Partitions {
	return &Partitions{
		interval:   interval,
		partitions: make(map[int64]*LazyPartition),
	}
}

func (ps *Partitions) Load(path string, create func(timestamp int64) (*LazyPartition, error)) error {
	partitions, err := fileutil.ListDir(path)
	if err != nil {
		return err
	}
	intervalCalc := ps.interval.Calculator()
	for _, partition := range partitions {
		partitionTime, err := intervalCalc.ParseSegmentTime(partition)
		if err != nil {
			return err
		}
		lp, err := create(partitionTime)
		if err != nil {
			return err
		}
		ps.PutPartition(partitionTime, lp)
	}
	return nil
}

func (ps *Partitions) GetPartition(timestamp int64) (partition Partition, ok bool, err error) {
	var p *LazyPartition
	p, ok = ps.partitions[timestamp]
	if !ok {
		return
	}
	partition, err = p.Get()
	return
}

func (ps *Partitions) PutPartition(timestamp int64, partition *LazyPartition) {
	ps.partitions[timestamp] = partition
}

func (ps *Partitions) GetPartitions() []*LazyPartition {
	return lo.Values(ps.partitions)
}
