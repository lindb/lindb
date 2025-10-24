package log

import (
	"fmt"
	"path"
	"strconv"
	"sync"

	"github.com/lindb/common/pkg/fileutil"

	"github.com/lindb/lindb/kv"
	"github.com/lindb/lindb/pkg/timeutil"
	"github.com/lindb/lindb/storage/base"
	"github.com/lindb/lindb/storage/log/tblstore"
	"github.com/lindb/lindb/storage/store"
)

var (
	minuteInterval = timeutil.Interval(60_000)
	intervalCalc   = minuteInterval.Calculator()
)

type partition struct {
	base.Partition

	shard *shard

	kvStore        kv.Store
	timestampIndex kv.Family

	segments map[int]store.Segment

	mutex sync.Mutex
}

func NewPartition(timestamp int64, shard *shard) (store.Partition, error) {
	partitionName := intervalCalc.GetSegment(timestamp)
	dir := store.PartitionPath(shard.Database().Name(), shard.ShardID(), minuteInterval, partitionName)
	p := &partition{
		Partition: base.Partition{
			Timestamp: intervalCalc.CalcSegmentTime(timestamp),
			Dir:       dir,
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
	// TODO: close kv store if err

	kvFamily := p.kvStore.GetFamily("t")
	if kvFamily == nil {
		// create kv family
		var err error
		familyOption := kv.FamilyOption{
			CompactThreshold: 0,
			Merger:           string(tblstore.LogIndexMerger),
		}
		kvFamily, err = p.kvStore.CreateFamily("t", familyOption)
		if err != nil {
			return nil, err
		}
	}
	p.timestampIndex = kvFamily

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

	segment, err := NewSegment(timestamp, p)
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
		if (&timeRange).Overlap(segment.SegmentTimeRange()) {
			segments = append(segments, segment)
		}
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
