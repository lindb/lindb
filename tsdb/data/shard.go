package data

import (
	"time"

	"github.com/eleme/lindb/models"
	"github.com/eleme/lindb/pkg/option"
	"github.com/eleme/lindb/tsdb/mem"
)

type Shard struct {
	id     int32
	option option.ShardOption
	mem    *mem.MemoryDatabase
}

func NewShard(shardID int32, option option.ShardOption) *Shard {
	shard := &Shard{
		id:     shardID,
		option: option,
		mem:    mem.NewMemDatabase(),
	}
	return shard
}

func (s *Shard) Write(point models.Point) {
	timestamp := point.Timestamp()
	now := time.Now().Unix()

	if timestamp < now-s.option.Behind {
		return
	}

	if timestamp > now+s.option.Ahead {
		return
	}

	timeSeriesStore := s.mem.GetTimeSeriesStore(point.Name(), point.Tags())

	//use family base time for memory store
	segmentTime := s.option.IntervalType.CalFamilyBaseTime(point.Timestamp())
	//slot time for ts data store
	slotTime := s.option.IntervalType.CalSlot(point.Timestamp())

	for k, v := range point.Fields() {
		fieldStore := timeSeriesStore.GetFieldStore(k)
		segmentStore := fieldStore.GetSegmentStore(segmentTime)

		segmentStore.Write(slotTime, v)
	}
}
