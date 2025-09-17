package log

import (
	"fmt"
	"path"
	"path/filepath"
	"strconv"
	"sync"

	"github.com/lindb/common/pkg/fileutil"

	"github.com/lindb/lindb/kv"
	"github.com/lindb/lindb/pkg/timeutil"
	"github.com/lindb/lindb/storage/base"
	"github.com/lindb/lindb/storage/store"
	"github.com/lindb/lindb/storage/utils"
)

var (
	minuteInterval = timeutil.Interval(60_000)
	intervalCalc   = minuteInterval.Calculator()
)

type partition struct {
	base.Partition

	shard *shard

	kvStore kv.Store

	segments map[int]store.Segment

	mutex sync.Mutex
}

func NewPartition(timestamp int64, shard *shard) (store.Partition, error) {
	p := &partition{
		Partition: base.Partition{
			Timestamp: intervalCalc.CalcSegmentTime(timestamp),
			Dir:       filepath.Join(utils.ShardPath(shard.Database().Name(), shard.ShardID()), intervalCalc.GetSegment(timestamp)),
		},

		shard:    shard,
		segments: make(map[int]store.Segment),
	}
	storeOption := kv.DefaultStoreOption()
	kvStore, err := kv.GetStoreManager().CreateStore(path.Join(p.Dir, "secondary"), storeOption)
	if err != nil {
		return nil, fmt.Errorf("create kv store for segment error:%s", err)
	}
	p.kvStore = kvStore

	segments, err := fileutil.ListDir(p.Dir)
	if err != nil {
		return nil, err
	}
	for _, segment := range segments {
		if segment == "secondary" {
			continue
		}
		segmentSlot, err := strconv.Atoi(segment)
		if err != nil {
			// TODO: add metric
			continue
		}
		// create data family
		segmentTime := intervalCalc.CalcFamilyStartTime(timestamp, segmentSlot)
		fmt.Printf("create segment: %d\n", segmentTime)
		p.GetOrCreateSegment(segmentTime)
	}

	return p, nil
}

func (p *partition) GetOrCreateSegment(timestamp int64) (store.Segment, error) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	segmentKey := intervalCalc.CalcFamily(timestamp, intervalCalc.CalcSegmentTime(timestamp))

	if segment, ok := p.segments[segmentKey]; ok {
		return segment, nil
	}

	segment, err := NewSegment(p)
	if err != nil {
		return nil, err
	}
	p.segments[segmentKey] = segment
	return segment, nil
}

func (p *partition) GetSegments(timeRange timeutil.TimeRange) (segments []store.Segment) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	for _, segment := range p.segments {
		segments = append(segments, segment)
	}

	return
}

func (p *partition) Shard() store.Shard {
	return p.shard
}

// Close implements store.Partition.
func (p *partition) Close() error {
	for _, segment := range p.segments {
		segment.Close()
	}
	return nil
}
