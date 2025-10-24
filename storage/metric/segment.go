package metric

import (
	"errors"
	"fmt"
	"path/filepath"
	"sync"
	"time"

	"github.com/lindb/common/pkg/fasttime"
	"github.com/lindb/common/pkg/logger"
	"go.uber.org/atomic"

	"github.com/lindb/lindb/constants"
	"github.com/lindb/lindb/flow"
	"github.com/lindb/lindb/kv"
	"github.com/lindb/lindb/metrics"
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/compress"
	"github.com/lindb/lindb/pkg/timeutil"
	"github.com/lindb/lindb/series/metric"
	"github.com/lindb/lindb/storage/base"
	"github.com/lindb/lindb/storage/metric/memdb"
	"github.com/lindb/lindb/storage/metric/tblstore/metricsdata"
	"github.com/lindb/lindb/storage/store"
)

type Segment struct {
	base.Segment
	partition *partition

	kvFamily kv.Family

	intervalCalc timeutil.IntervalCalculator
	interval     timeutil.Interval

	immutableMemDB memdb.MemoryDatabase
	mutableMemDB   memdb.MemoryDatabase

	mutex sync.Mutex

	lastReadTime *atomic.Int64

	isFlushing     atomic.Bool
	flushCondition sync.WaitGroup
	lastFlushTime  int64

	statistics *metrics.FamilyStatistics
	logger     logger.Logger
}

func NewSegment(timestamp int64, partition *partition) (store.Segment, error) {
	interval := partition.PartitionInterval()
	intervalCalc := interval.Calculator()
	segmentTime := intervalCalc.CalcSegmentTime(timestamp)
	familySlot := intervalCalc.CalcFamily(timestamp, segmentTime)
	family := fmt.Sprintf("%d", familySlot)
	segmentPath := filepath.Join(partition.Path(), family)
	fmt.Printf("segmentPath:%s\n", segmentPath)
	kvFamily := partition.kvStore.GetFamily(family)
	if kvFamily == nil {
		// create kv family
		var err error
		familyOption := kv.FamilyOption{
			CompactThreshold: 0,
			Merger:           string(metricsdata.MetricDataMerger),
		}
		kvFamily, err = partition.kvStore.CreateFamily(family, familyOption)
		if err != nil {
			return nil, err
		}
	}
	segmentStartTime := intervalCalc.CalcFamilyStartTime(segmentTime, familySlot)
	shard := partition.shard
	db := shard.db
	seg := &Segment{
		Segment: base.Segment{
			TimeRange: timeutil.TimeRange{
				Start: segmentStartTime,
				End:   intervalCalc.CalcFamilyEndTime(segmentStartTime),
			},
			Path:              segmentPath,
			WALs:              make(map[models.NodeID]store.WriteAheadLog),
			Sequence:          make(map[models.NodeID]atomic.Int64),
			ImmutableSequence: make(map[models.NodeID]int64),
			PersistSequence:   make(map[models.NodeID]atomic.Int64),
		},
		partition:    partition,
		kvFamily:     kvFamily,
		interval:     interval,
		intervalCalc: intervalCalc,
		lastReadTime: atomic.NewInt64(fasttime.UnixMilliseconds()),
		statistics:   metrics.NewFamilyStatistics(db.Name(), shard.id.String()),
		logger:       logger.GetLogger("Metric", "Segment"),
	}

	seg.CreateWriteAheadLog = func(p string) (store.WriteAheadLog, error) {
		return store.CreateWriteAheadLog(p, seg)
	}

	return seg, nil
}

func (s *Segment) Partition() store.Partition {
	return s.partition
}

func (s *Segment) write(rows []*metric.StorageRow) error {
	if len(rows) == 0 {
		return nil
	}

	db, err := s.GetOrCreateMemoryDatabase(s.TimeRange.Start)
	if err != nil {
		// all rows are dropped
		s.statistics.WriteMetricFailures.Add(float64(len(rows)))
		return err
	}
	db.AcquireWrite()
	defer func() {
		s.statistics.WriteBatches.Incr()
		db.CompleteWrite()
	}()

	for idx := range rows {
		row := rows[idx]
		err := db.WriteRow(row)
		if err == nil {
			s.statistics.WriteMetrics.Incr()
			s.statistics.WriteFields.Add(row.WrittenFields)
		} else {
			s.statistics.WriteMetricFailures.Incr()
			s.logger.Error("failed writing row", logger.String("family", s.Path), logger.Error(err))
		}

		// waiting all operators done(write data/build meta and index)
		// TODO: add timeout??
		row.Wait()
	}

	return nil
}

// GetOrCreateMemoryDatabase returns memory database by given segment time.
func (s *Segment) GetOrCreateMemoryDatabase(segmentTime int64) (memdb.MemoryDatabase, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if s.mutableMemDB == nil {
		shard := s.partition.shard
		newDB, err := memdb.NewMemoryDatabase(&memdb.MemoryDatabaseCfg{
			SegmentTime:   segmentTime,
			IntervalCalc:  s.intervalCalc,
			Interval:      s.interval,
			Name:          shard.Database().Name(),
			IndexDatabase: shard.MemIndexDB(),
			BufferMgr:     shard.BufferManager(),
		})
		if err != nil {
			return nil, err
		}
		s.mutableMemDB = newDB
		s.statistics.ActiveMemDBs.Incr()
	}
	return s.mutableMemDB, nil
}

// Flush implements store.Segment.
func (s *Segment) Flush() error {
	if s.isFlushing.CompareAndSwap(false, true) {
		defer func() {
			// mark flush job complete, notify
			s.flushCondition.Done()
			s.isFlushing.Store(false)
		}()

		// 1. mark flush job doing
		s.flushCondition.Add(1)

		startTime := time.Now()

		// add lock when switch memory database
		s.mutex.Lock()
		if s.immutableMemDB != nil || s.mutableMemDB == nil || s.mutableMemDB.NumOfSeries() == 0 {
			// if immutable memory database not nil or no data need flush, return it
			s.mutex.Unlock()
			return nil
		}
		waitingFlushMemDB := s.mutableMemDB
		s.immutableMemDB = waitingFlushMemDB
		s.mutableMemDB = nil
		// mark mutable memory database nil, write data will be created
		waitingFlushMemDB.MarkReadOnly()
		immutableSeq := make(map[models.NodeID]int64)
		for leader, seq := range s.Sequence {
			immutableSeq[leader] = seq.Load()
		}
		s.ImmutableSequence = immutableSeq
		s.mutex.Unlock()

		if err := s.flushMemoryDatabase(immutableSeq, waitingFlushMemDB); err != nil {
			return err
		}

		// flush success, mark immutable memory database nil
		s.mutex.Lock()
		s.immutableMemDB = nil
		s.ImmutableSequence = nil
		// save persisted sequence, ack replica sequence in flushMemoryDatabase func
		for leader, seq := range immutableSeq {
			s.PersistSequence[leader] = *atomic.NewInt64(seq)
		}

		s.mutex.Unlock()

		endTime := time.Now()
		s.lastFlushTime = endTime.UnixMilli()
		s.logger.Info("flush memory database successfully",
			logger.String("segment", s.Path),
			logger.String("flush-duration", endTime.Sub(startTime).String()),
			logger.Int64("segmentTime", s.TimeRange.Start),
			logger.Int64("memDBSize", waitingFlushMemDB.MemSize()))
	}

	// another flush process is running
	return nil
}

// Write implements store.Segment.
func (s *Segment) Write(leader models.NodeID, seq int64, msg []byte) (rows int, err error) {
	reader := compress.NewSnappyReader()
	block, err := reader.Uncompress(msg)
	if err != nil {
		return
	}
	batchRows := metric.NewStorageBatchRows()
	batchRows.UnmarshalRows(block)
	rows = batchRows.Len()

	if rows == 0 {
		return
	}
	err = s.write(batchRows.Rows())

	return
}

// Close implements store.Segment.
func (s *Segment) Close() error {
	s.logger.Info("starting close data segment", logger.String("segment", s.Path))
	start := time.Now()

	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.flushCondition.Wait()

	if s.immutableMemDB != nil {
		if err := s.flushMemoryDatabase(s.ImmutableSequence, s.immutableMemDB); err != nil {
			return err
		}
	}
	if s.mutableMemDB != nil {
		sequences := make(map[models.NodeID]int64)
		for leader, seq := range s.Sequence {
			sequences[leader] = seq.Load()
		}
		if err := s.flushMemoryDatabase(sequences, s.mutableMemDB); err != nil {
			return err
		}
	}

	// FIXME:
	// GetFamilyManager().RemoveFamily(f)
	s.statistics.ActiveFamilies.Decr()

	s.logger.Info("close data segment complete", logger.String("segment", s.Path), logger.Any("cost", time.Since(start)))
	return nil
}

// Filter filters the data based on metric/version/seriesIDs,
// if it finds data then returns the FilterResultSet, else returns nil
func (s *Segment) Filter(ctx *flow.MetricScanContext) (resultSet []flow.FilterResultSet, err error) {
	s.lastReadTime.Store(fasttime.UnixMilliseconds())
	memRS, err := s.memoryFilter(ctx)
	if !errors.Is(err, constants.ErrNotFound) && err != nil {
		fmt.Printf("mem filter=%v\n", err)
		// FIXME: ignore not found??
		return nil, err
	}
	fileRS, err := s.fileFilter(ctx)
	if !errors.Is(err, constants.ErrNotFound) && err != nil {
		fmt.Printf("file filter=%v\n", err)
		return nil, err
	}
	resultSet = append(resultSet, memRS...)
	resultSet = append(resultSet, fileRS...)
	return
}

func (s *Segment) memoryFilter(ctx *flow.MetricScanContext) (resultSet []flow.FilterResultSet, err error) {
	memFilter := func(memDB memdb.MemoryDatabase) error {
		rs, err := memDB.Filter(ctx)
		if err != nil {
			fmt.Printf("mem db error=%v\n", err)
			return err
		}
		resultSet = append(resultSet, rs...)
		fmt.Println("found mem data", len(resultSet))
		return nil
	}

	s.mutex.Lock()
	defer s.mutex.Unlock()

	if s.mutableMemDB != nil {
		if err := memFilter(s.mutableMemDB); err != nil {
			return nil, err
		}
	}
	if s.immutableMemDB != nil {
		if err := memFilter(s.immutableMemDB); err != nil {
			return nil, err
		}
	}
	return
}

func (s *Segment) fileFilter(ctx *flow.MetricScanContext) (resultSet []flow.FilterResultSet, err error) {
	snapShot := s.kvFamily.GetSnapshot()
	defer func() {
		if err != nil || len(resultSet) == 0 {
			// if not find metrics data or has error, close snapshot directly
			snapShot.Close()
		}
	}()
	metricKey := uint32(ctx.MetricID)
	readers, err := snapShot.FindReaders(metricKey)
	if err != nil {
		s.logger.Error("filter data family error", logger.Error(err))
		return nil, err
	}
	querySlotRange := s.interval.CalcSlotRange(s.TimeRange.Start, ctx.TimeRange)
	fmt.Printf("find reader =%v,%v\n", readers, metricKey)
	var metricReaders []metricsdata.MetricReader
	for _, reader := range readers {
		value, err0 := reader.Get(metricKey)
		// metric data not found
		if err0 != nil {
			fmt.Println("metric not found from file")
			continue
		}
		r, err := metricsdata.NewReader(reader.Path(), value)
		if err != nil {
			fmt.Printf("new reader file=%v\n", err)
			return nil, err
		}
		storageTimeRange := r.GetTimeRange()
		if storageTimeRange.Overlap(querySlotRange) {
			metricReaders = append(metricReaders, r)
		} else {
			fmt.Printf("file time range out...,%v,%v\n", storageTimeRange, ctx.TimeRange)
		}
	}
	if len(metricReaders) == 0 {
		fmt.Println("no file found")
		return nil, nil
	}
	filter := metricsdata.NewFilter(s.TimeRange.Start, s.interval, querySlotRange, snapShot, metricReaders)
	return filter.Filter(ctx.SeriesIDs, ctx.Fields)
}

// flushMemoryDatabase flushes memory database to disk.
func (s *Segment) flushMemoryDatabase(sequences map[models.NodeID]int64, memDB memdb.MemoryDatabase) error {
	startTime := time.Now()
	flusher := s.kvFamily.NewFlusher()
	defer func() {
		flusher.Release()
		s.statistics.MemDBFlushDuration.UpdateSince(startTime)
	}()

	for leader, seq := range sequences {
		flusher.Sequence(int32(leader), seq)
	}

	dataFlusher, err := metricsdata.NewFlusher(flusher)
	if err != nil {
		return err
	}
	// flush family data
	if err := memDB.FlushFamilyTo(dataFlusher); err != nil {
		s.logger.Error("failed to flush memory database",
			logger.String("segment", s.Path),
			logger.Int64("memDBSize", memDB.MemSize()))
		s.statistics.MemDBFlushFailures.Incr()
		return err
	}

	// FIXME:
	// invoke sequence ack callback
	// for leader, seq := range sequences {
	// 	if callbacks, ok := f.callbacks[leader]; ok {
	// 		for _, fn := range callbacks {
	// 			fn(seq)
	// 		}
	// 	}
	// }

	s.statistics.ActiveMemDBs.Decr()

	if err := memDB.Close(); err != nil {
		// ignore close memory database err, if not maybe write duplicate data into file storage
		s.logger.Warn("failed to close memory database",
			logger.String("segment", s.Path),
			logger.Int64("memDBSize", memDB.MemSize()))
		return nil
	}

	return nil
}
