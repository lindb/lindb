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

package tsdb

import (
	"context"
	"fmt"
	"math"
	"path/filepath"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/config"
	"github.com/lindb/lindb/constants"
	"github.com/lindb/lindb/kv"
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/fasttime"
	"github.com/lindb/lindb/pkg/fileutil"
	"github.com/lindb/lindb/pkg/option"
	"github.com/lindb/lindb/pkg/timeutil"
	protoMetricsV1 "github.com/lindb/lindb/proto/gen/v1/metrics"
	"github.com/lindb/lindb/series/field"
	"github.com/lindb/lindb/series/tag"
	"github.com/lindb/lindb/tsdb/indexdb"
	"github.com/lindb/lindb/tsdb/memdb"
	"github.com/lindb/lindb/tsdb/metadb"
)

var _testShard1Path = filepath.Join(testPath, shardDir, "1")

func TestShard_New(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer func() {
		_ = fileutil.RemoveDir(testPath)
		mkDirIfNotExist = fileutil.MkDirIfNotExist
		newReplicaSequenceFunc = newReplicaSequence
		newIntervalSegmentFunc = newIntervalSegment
		newKVStoreFunc = kv.NewStore
		newIndexDBFunc = indexdb.NewIndexDatabase
		newMemoryDBFunc = memdb.NewMemoryDatabase

		ctrl.Finish()
	}()

	db := NewMockDatabase(ctrl)
	meta := metadb.NewMockMetadata(ctrl)
	meta.EXPECT().DatabaseName().Return("test").AnyTimes()
	db.EXPECT().Name().Return("db").AnyTimes()
	db.EXPECT().Metadata().Return(meta).AnyTimes()
	// case 1: database option err
	thisShard, err := newShard(db, 1, _testShard1Path, option.DatabaseOption{})
	assert.Error(t, err)
	assert.Nil(t, thisShard)
	// case 2: interval err
	thisShard, err = newShard(db, 1, _testShard1Path, option.DatabaseOption{Interval: "as"})
	assert.Error(t, err)
	assert.Nil(t, thisShard)
	// case 3: create path err
	mkDirIfNotExist = func(path string) error {
		return fmt.Errorf("err")
	}
	thisShard, err = newShard(db, 1, _testShard1Path, option.DatabaseOption{Interval: "10s"})
	assert.Error(t, err)
	assert.Nil(t, thisShard)
	// case 4: new replica sequence err
	mkDirIfNotExist = fileutil.MkDirIfNotExist
	newReplicaSequenceFunc = func(dirPath string) (sequence ReplicaSequence, err error) {
		return nil, fmt.Errorf("err")
	}
	thisShard, err = newShard(db, 1, _testShard1Path, option.DatabaseOption{Interval: "10s"})
	assert.Error(t, err)
	assert.Nil(t, thisShard)
	// case 5: new interval segment err
	newReplicaSequenceFunc = newReplicaSequence
	newIntervalSegmentFunc = func(interval timeutil.Interval, path string) (segment IntervalSegment, err error) {
		return nil, fmt.Errorf("err")
	}
	thisShard, err = newShard(db, 1, _testShard1Path, option.DatabaseOption{Interval: "10s"})
	assert.Error(t, err)
	assert.Nil(t, thisShard)
	// case 6: new kv store err
	newIntervalSegmentFunc = newIntervalSegment
	newKVStoreFunc = func(name string, option kv.StoreOption) (store kv.Store, err error) {
		return nil, fmt.Errorf("err")
	}
	thisShard, err = newShard(db, 1, _testShard1Path, option.DatabaseOption{Interval: "10s"})
	assert.Error(t, err)
	assert.Nil(t, thisShard)
	// case 7: create forward family err
	kvStore := kv.NewMockStore(ctrl)
	kvStore.EXPECT().Close().Return(fmt.Errorf("err")).AnyTimes()
	newKVStoreFunc = func(name string, option kv.StoreOption) (store kv.Store, err error) {
		return kvStore, nil
	}
	kvStore.EXPECT().CreateFamily(forwardIndexDir, gomock.Any()).Return(nil, fmt.Errorf("err"))
	thisShard, err = newShard(db, 1, _testShard1Path, option.DatabaseOption{Interval: "10s"})
	assert.Error(t, err)
	assert.Nil(t, thisShard)
	// case 8: create forward family err
	family := kv.NewMockFamily(ctrl)
	kvStore.EXPECT().CreateFamily(forwardIndexDir, gomock.Any()).Return(family, nil)
	kvStore.EXPECT().CreateFamily(invertedIndexDir, gomock.Any()).Return(nil, fmt.Errorf("err"))
	thisShard, err = newShard(db, 1, _testShard1Path, option.DatabaseOption{Interval: "10s"})
	assert.Error(t, err)
	assert.Nil(t, thisShard)
	// case 9: create index db err
	kvStore.EXPECT().CreateFamily(gomock.Any(), gomock.Any()).Return(family, nil).AnyTimes()
	newIndexDBFunc = func(ctx context.Context, parent string,
		metadata metadb.Metadata, forward kv.Family, inverted kv.Family,
	) (indexDatabase indexdb.IndexDatabase, err error) {
		return nil, fmt.Errorf("err")
	}
	thisShard, err = newShard(db, 1, _testShard1Path, option.DatabaseOption{Interval: "10s"})
	assert.Error(t, err)
	assert.Nil(t, thisShard)
	newIndexDBFunc = indexdb.NewIndexDatabase

	// case 10: create shard success
	thisShard, err = newShard(db, 1, _testShard1Path, option.DatabaseOption{Interval: "10s"})
	assert.NoError(t, err)
	assert.NotNil(t, thisShard)
	assert.NotNil(t, thisShard.IndexDatabase())
	assert.Equal(t, "db", thisShard.DatabaseName())
	assert.Equal(t, models.ShardID(1), thisShard.ShardID())
	s, err := thisShard.GetOrCreateSequence("tes")
	assert.NoError(t, err)
	assert.NotNil(t, s)

	assert.True(t, fileutil.Exist(_testShard1Path))
	assert.False(t, thisShard.IsFlushing())
}

func TestShard_GetDataFamilies(t *testing.T) {
	defer func() {
		_ = fileutil.RemoveDir(testPath)
	}()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	db := NewMockDatabase(ctrl)
	meta := metadb.NewMockMetadata(ctrl)
	meta.EXPECT().DatabaseName().Return("test").AnyTimes()
	db.EXPECT().Name().Return("test-db").AnyTimes()
	db.EXPECT().Metadata().Return(meta).AnyTimes()
	s, _ := newShard(db, 1, _testShard1Path, option.DatabaseOption{Interval: "10s"})
	assert.Nil(t, s.GetDataFamilies(timeutil.Month, timeutil.TimeRange{}))
	assert.Nil(t, s.GetDataFamilies(timeutil.Day, timeutil.TimeRange{}))
	assert.Equal(t, 0, len(s.GetDataFamilies(timeutil.Day, timeutil.TimeRange{})))
}

func Test_Shard_validateMetric(t *testing.T) {
	s := &shard{metrics: *newShardMetrics("1", 1)}
	assert.Zero(t, s.CurrentInterval().Int64())
	// nil pb
	err := s.validateMetric(nil)
	assert.Error(t, err)
	// empty name
	err = s.validateMetric(&protoMetricsV1.Metric{Name: ""})
	assert.Error(t, err)
	// field empty
	err = s.validateMetric(&protoMetricsV1.Metric{Name: "1"})
	assert.Error(t, err)

	// bad tags, empty
	err = s.validateMetric(&protoMetricsV1.Metric{
		Name: "1",
		Tags: []*protoMetricsV1.KeyValue{
			{Key: "", Value: ""},
		},
		SimpleFields: []*protoMetricsV1.SimpleField{
			{Name: "1", Type: protoMetricsV1.SimpleFieldType_DELTA_SUM, Value: 1},
		},
		Timestamp: fasttime.UnixMilliseconds(),
	})
	assert.Error(t, err)

	// bad tags, nil
	err = s.validateMetric(&protoMetricsV1.Metric{
		Name: "1",
		Tags: []*protoMetricsV1.KeyValue{nil, nil},
		SimpleFields: []*protoMetricsV1.SimpleField{
			{Name: "1", Type: protoMetricsV1.SimpleFieldType_DELTA_SUM, Value: 1},
		},
		Timestamp: fasttime.UnixMilliseconds(),
	})
	assert.Error(t, err)
	// simple fields nil
	err = s.validateMetric(&protoMetricsV1.Metric{
		Name:         "1",
		SimpleFields: []*protoMetricsV1.SimpleField{nil, nil},
		Timestamp:    fasttime.UnixMilliseconds(),
	})
	assert.Error(t, err)
	// field name empty
	err = s.validateMetric(&protoMetricsV1.Metric{
		Name: "1",
		SimpleFields: []*protoMetricsV1.SimpleField{
			{Name: "", Type: protoMetricsV1.SimpleFieldType_DELTA_SUM, Value: 1},
		},
		Timestamp: fasttime.UnixMilliseconds(),
	})
	assert.Error(t, err)
	// sanitize field name
	err = s.validateMetric(&protoMetricsV1.Metric{
		Name: "1",
		SimpleFields: []*protoMetricsV1.SimpleField{
			{Name: "Histogram_2", Type: protoMetricsV1.SimpleFieldType_DELTA_SUM, Value: 1},
		},
		Timestamp: fasttime.UnixMilliseconds(),
	})
	assert.NoError(t, err)
	// sanitize field name, field type unspecified
	err = s.validateMetric(&protoMetricsV1.Metric{
		Name: "1",
		SimpleFields: []*protoMetricsV1.SimpleField{
			{Name: "xxx2", Type: protoMetricsV1.SimpleFieldType_SIMPLE_UNSPECIFIED, Value: 1},
		},
		Timestamp: fasttime.UnixMilliseconds(),
	})
	assert.Error(t, err)
	// Nan number
	err = s.validateMetric(&protoMetricsV1.Metric{
		Name: "1",
		SimpleFields: []*protoMetricsV1.SimpleField{
			{Value: math.Log(-1), Name: "222", Type: protoMetricsV1.SimpleFieldType_DELTA_SUM},
		},
		Timestamp: fasttime.UnixMilliseconds(),
	})
	assert.Error(t, err)

	// Inf number
	err = s.validateMetric(&protoMetricsV1.Metric{
		Name: "1",
		SimpleFields: []*protoMetricsV1.SimpleField{
			{Value: math.Inf(1) + 1, Name: "222", Type: protoMetricsV1.SimpleFieldType_DELTA_SUM},
		},
		Timestamp: fasttime.UnixMilliseconds(),
	})
	assert.Error(t, err)
	//
	// validate compound field
	//
	// length not match
	err = s.validateMetric(&protoMetricsV1.Metric{
		Name: "1",
		CompoundField: &protoMetricsV1.CompoundField{
			Values:         []float64{1, 2, 3},
			ExplicitBounds: []float64{1, 2, 3, 4},
		},
		Timestamp: fasttime.UnixMilliseconds(),
	})
	assert.Error(t, err)
	// length too short
	err = s.validateMetric(&protoMetricsV1.Metric{
		Name: "1",
		CompoundField: &protoMetricsV1.CompoundField{
			Values:         []float64{1, 2},
			ExplicitBounds: []float64{1, math.Inf(1) + 1},
		},
		Timestamp: fasttime.UnixMilliseconds(),
	})
	assert.Error(t, err)
	// min, max < 0
	err = s.validateMetric(&protoMetricsV1.Metric{
		Name: "1",
		CompoundField: &protoMetricsV1.CompoundField{
			Sum:            -1,
			Values:         []float64{1, 2, 3, 4},
			ExplicitBounds: []float64{1, 2, 3, math.Inf(1) + 1},
		},
		Timestamp: fasttime.UnixMilliseconds(),
	})
	assert.Error(t, err)
	// check value
	err = s.validateMetric(&protoMetricsV1.Metric{
		Name: "1",
		CompoundField: &protoMetricsV1.CompoundField{
			Sum:            11,
			Values:         []float64{-1, 2, 3, 4},
			ExplicitBounds: []float64{1, 2, 3, math.Inf(1) + 1},
		},
		Timestamp: fasttime.UnixMilliseconds(),
	})
	assert.Error(t, err)
	// check increase
	err = s.validateMetric(&protoMetricsV1.Metric{
		Name: "1",
		CompoundField: &protoMetricsV1.CompoundField{
			Sum:            11,
			Values:         []float64{1, 4, 3, 4},
			ExplicitBounds: []float64{1, 5, 3, math.Inf(1) + 1},
		},
		Timestamp: fasttime.UnixMilliseconds(),
	})
	assert.Error(t, err)
	// check last bound
	err = s.validateMetric(&protoMetricsV1.Metric{
		Name: "1",
		CompoundField: &protoMetricsV1.CompoundField{
			Sum:            11,
			Values:         []float64{1, 4, 3, 4},
			ExplicitBounds: []float64{1, 2, 3, 4},
		},
		Timestamp: fasttime.UnixMilliseconds(),
	})
	assert.Error(t, err)
	// ok
	err = s.validateMetric(&protoMetricsV1.Metric{
		Name: "1",
		CompoundField: &protoMetricsV1.CompoundField{
			Sum:            11,
			Values:         []float64{1, 4, 3, 4},
			ExplicitBounds: []float64{1, 2, 3, math.Inf(1) + 1},
		},
		Timestamp: fasttime.UnixMilliseconds(),
	})
	assert.NoError(t, err)
}

func TestShard_Write(t *testing.T) {
	defer func() {
		_ = fileutil.RemoveDir(testPath)
	}()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	db := NewMockDatabase(ctrl)
	metadata := metadb.NewMockMetadata(ctrl)
	metadata.EXPECT().DatabaseName().Return("test").AnyTimes()
	metadataDB := metadb.NewMockMetadataDatabase(ctrl)
	indexDB := indexdb.NewMockIndexDatabase(ctrl)
	metadata.EXPECT().MetadataDatabase().Return(metadataDB).AnyTimes()
	db.EXPECT().Name().Return("test-db").AnyTimes()
	db.EXPECT().Metadata().Return(metadata).AnyTimes()

	mockMemDB := memdb.NewMockMemoryDatabase(ctrl)
	mockMemDB.EXPECT().AcquireWrite().AnyTimes()
	mockMemDB.EXPECT().CompleteWrite().AnyTimes()
	mockMemDB.EXPECT().MemSize().Return(int64(100)).AnyTimes()
	mockMemDB.EXPECT().Write(gomock.Any()).Return(nil).AnyTimes()
	// calculate family start time and slot index
	shardINTF, _ := newShard(db, 1, _testShard1Path, option.DatabaseOption{Interval: "10s", Behind: "1m", Ahead: "1m"})
	timestamp := timeutil.Now()
	var interval timeutil.Interval
	_ = interval.ValueOf("10s")
	intervalCalc := interval.Calculator()
	segmentTime := intervalCalc.CalcSegmentTime(timestamp)              // day
	family := intervalCalc.CalcFamily(timestamp, segmentTime)           // hours
	familyTime := intervalCalc.CalcFamilyStartTime(segmentTime, family) // family timestamp
	shardIns := shardINTF.(*shard)
	shardIns.indexDB = indexDB
	shardIns.families.InsertFamily(familyTime, mockMemDB)

	// case 1: metric nil
	assert.Error(t, shardINTF.Write(nil))
	// case 2: metric name is empty
	assert.Error(t, shardINTF.Write(&protoMetricsV1.Metric{
		Timestamp: timestamp,
	}))
	// case 3: field is empty
	assert.Error(t, shardINTF.Write(&protoMetricsV1.Metric{
		Name:      "test",
		Timestamp: timestamp,
	}))
	// case 4: gen metric id err
	metadataDB.EXPECT().GenMetricID(constants.DefaultNamespace, "test").Return(uint32(0), fmt.Errorf("err"))
	assert.Error(t, shardINTF.Write(&protoMetricsV1.Metric{
		Name:      "test",
		Timestamp: timestamp,
		SimpleFields: []*protoMetricsV1.SimpleField{{
			Name:  "f1",
			Value: 1.0,
			Type:  protoMetricsV1.SimpleFieldType_DELTA_SUM,
		}},
	}))
	// case 5: gen series id err
	metadataDB.EXPECT().GenMetricID(constants.DefaultNamespace, "test").Return(uint32(10), nil).AnyTimes()
	indexDB.EXPECT().GetOrCreateSeriesID(uint32(10), uint64(9)).Return(uint32(0), false, fmt.Errorf("err"))
	assert.Error(t, shardINTF.Write(&protoMetricsV1.Metric{
		Name:      "test",
		Timestamp: timestamp,
		TagsHash:  9,
		Tags:      tag.KeyValuesFromMap(map[string]string{"ip": "1.1.1.1"}),
		SimpleFields: []*protoMetricsV1.SimpleField{{
			Name:  "f1",
			Value: 1.0,
			Type:  protoMetricsV1.SimpleFieldType_DELTA_SUM,
		}},
	}))
	// case 6: get old series id
	metadataDB.EXPECT().GenMetricID(constants.DefaultNamespace, "test").Return(uint32(10), nil).AnyTimes()
	metadataDB.EXPECT().GenFieldID(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(field.ID(1), nil)
	indexDB.EXPECT().GetOrCreateSeriesID(uint32(10), uint64(11)).Return(uint32(10), false, nil)
	assert.NoError(t, shardINTF.Write(&protoMetricsV1.Metric{
		Name:      "test",
		Timestamp: timestamp,
		TagsHash:  11,
		Tags:      tag.KeyValuesFromMap(map[string]string{"ip": "1.1.1.1"}),
		SimpleFields: []*protoMetricsV1.SimpleField{{
			Name:  "f1",
			Value: 1.0,
			Type:  protoMetricsV1.SimpleFieldType_DELTA_SUM,
		}},
	}))
	// case 7: create new series id
	indexDB.EXPECT().GetOrCreateSeriesID(uint32(10), uint64(10)).Return(uint32(10), true, nil)
	metadataDB.EXPECT().GenFieldID(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(field.ID(1), nil)
	indexDB.EXPECT().BuildInvertIndex(constants.DefaultNamespace, "test", tag.KeyValuesFromMap(map[string]string{"ip": "1.1.1.1"}), uint32(10))
	assert.NoError(t, shardINTF.Write(&protoMetricsV1.Metric{
		Name:      "test",
		Timestamp: timestamp,
		TagsHash:  10,
		Tags:      tag.KeyValuesFromMap(map[string]string{"ip": "1.1.1.1"}),
		SimpleFields: []*protoMetricsV1.SimpleField{{
			Name:  "f1",
			Value: 1.0,
			Type:  protoMetricsV1.SimpleFieldType_DELTA_SUM,
		}},
	}))
	// case 8: write metric without tags
	metadataDB.EXPECT().GenFieldID(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(field.ID(1), nil)
	assert.NoError(t, shardINTF.Write(&protoMetricsV1.Metric{
		Name:      "test",
		Timestamp: timestamp,
		TagsHash:  10,
		SimpleFields: []*protoMetricsV1.SimpleField{{
			Name:  "f1",
			Value: 1.0,
			Type:  protoMetricsV1.SimpleFieldType_DELTA_SUM,
		}},
	}))
}

func Test_Shard_howManyFieldsWillWrite(t *testing.T) {
	var s = &shard{}
	assert.Equal(t, s.howManyFieldsWillWrite(_testMetric), 26)
	assert.Equal(t, s.howManyFieldsWillWrite(&protoMetricsV1.Metric{
		Name: "xxxx",
		SimpleFields: []*protoMetricsV1.SimpleField{
			{Name: "111", Type: protoMetricsV1.SimpleFieldType_DELTA_SUM, Value: 111},
			{Name: "Histogram111", Type: protoMetricsV1.SimpleFieldType_DELTA_SUM, Value: 2222},
		}},
	), 2)
}

var _testMetric = &protoMetricsV1.Metric{
	Name: "xxxx",
	Tags: []*protoMetricsV1.KeyValue{
		{Key: "a", Value: "v"},
		{Key: "1", Value: "2"},
	},
	TagsHash: 11111,
	SimpleFields: []*protoMetricsV1.SimpleField{
		{Name: "111", Type: protoMetricsV1.SimpleFieldType_DELTA_SUM, Value: 111},
		{Name: "Histogram111", Type: protoMetricsV1.SimpleFieldType_DELTA_SUM, Value: 2222},
	},
	CompoundField: &protoMetricsV1.CompoundField{
		Sum:            1,
		Count:          2222,
		Min:            111,
		Max:            333343,
		Values:         []float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20},
		ExplicitBounds: []float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, math.Inf(1) + 1},
	},
}

func Benchmark_validate_metric(b *testing.B) {
	var s = &shard{metrics: *newShardMetrics("1", 1)}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = s.validateMetric(_testMetric)
	}
}

func Test_familyMemDBSet(t *testing.T) {
	set := newFamilyMemDBSet()
	for i := 1000; i >= 0; i -= 10 {
		set.InsertFamily(int64(i), nil)
	}

	for i := 0; i < 1000; i += 10 {
		_, exist := set.GetMutableFamily(int64(i))
		assert.True(t, exist)
		_, exist = set.GetMutableFamily(int64(i + 1))
		assert.False(t, exist)
		_, exist = set.GetMutableFamily(int64(i - 1))
		assert.False(t, exist)
	}
}

func TestShard_Close(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer func() {
		_ = fileutil.RemoveDir(testPath)
		newKVStoreFunc = kv.NewStore
		ctrl.Finish()
	}()
	kvStore := kv.NewMockStore(ctrl)
	family := kv.NewMockFamily(ctrl)
	kvStore.EXPECT().CreateFamily(gomock.Any(), gomock.Any()).Return(family, nil).AnyTimes()
	newKVStoreFunc = func(name string, option kv.StoreOption) (s kv.Store, err error) {
		return kvStore, nil
	}
	meta := metadb.NewMockMetadata(ctrl)
	meta.EXPECT().DatabaseName().Return("test").AnyTimes()
	db := NewMockDatabase(ctrl)
	db.EXPECT().Name().Return("test-db").AnyTimes()
	db.EXPECT().Metadata().Return(meta).AnyTimes()
	s, _ := newShard(db, 1, _testShard1Path, option.DatabaseOption{Interval: "10s"})
	index := indexdb.NewMockIndexDatabase(ctrl)
	s1 := s.(*shard)
	s1.indexDB = index

	// case 1: close index err
	index.EXPECT().Close().Return(fmt.Errorf("err"))
	err := s.Close()
	assert.Error(t, err)
	// case 2: close index store err
	index.EXPECT().Close().Return(nil).AnyTimes()
	kvStore.EXPECT().Close().Return(fmt.Errorf("exx"))
	err = s.Close()
	assert.Error(t, err)
	// case 3: flush immutable err
	kvStore.EXPECT().Close().Return(nil).AnyTimes()
	mutableDB := memdb.NewMockMemoryDatabase(ctrl)
	immutableDB := memdb.NewMockMemoryDatabase(ctrl)
	s1.families.InsertFamily(1, mutableDB)
	s1.families.InsertFamily(2, immutableDB)
	s1.families.SetFamilyImmutable(2)
	mutableDB.EXPECT().FamilyTime().Return(int64(1)).AnyTimes()
	immutableDB.EXPECT().FamilyTime().Return(int64(2)).AnyTimes()

	immutableDB.EXPECT().FlushFamilyTo(gomock.Any()).Return(fmt.Errorf("err"))
	err = s.Close()
	assert.Error(t, err)
	// case 4: close success
	mutableDB.EXPECT().FlushFamilyTo(gomock.Any()).Return(nil)
	mutableDB.EXPECT().Close().Return(nil)
	immutableDB.EXPECT().FlushFamilyTo(gomock.Any()).Return(nil)
	immutableDB.EXPECT().Close().Return(nil)
	err = s.Close()
	assert.NoError(t, err)
	// case 5: close memory database err
	immutableDB.EXPECT().FlushFamilyTo(gomock.Any()).Return(nil)
	immutableDB.EXPECT().Close().Return(fmt.Errorf("error"))
	err = s.Close()
	assert.Error(t, err)

	// case6, get segment error
	mockSegment := NewMockIntervalSegment(ctrl)
	s1.segment = mockSegment
	mockSegment.EXPECT().GetOrCreateSegment(gomock.Any()).Return(nil, fmt.Errorf("error"))
	err = s.Close()
	assert.Error(t, err)
}

func TestShard_Flush(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer func() {
		_ = fileutil.RemoveDir(testPath)
		newMemoryDBFunc = memdb.NewMemoryDatabase
		ctrl.Finish()
	}()

	meta := metadb.NewMockMetadata(ctrl)
	meta.EXPECT().DatabaseName().Return("test").AnyTimes()
	db := NewMockDatabase(ctrl)
	db.EXPECT().Name().Return("test-db").AnyTimes()
	db.EXPECT().Metadata().Return(meta).AnyTimes()
	s, _ := newShard(db, 2, _testShard1Path, option.DatabaseOption{Interval: "10s"})
	index := indexdb.NewMockIndexDatabase(ctrl)
	s2 := s.(*shard)
	s2.indexDB = index

	//// case 1: flush is doing
	s2.isFlushing.Store(true)
	assert.NoError(t, s2.Flush())
	// case 2: flush index error
	s2.isFlushing.Store(false)
	index.EXPECT().Flush().Return(fmt.Errorf("error"))
	assert.Error(t, s2.Flush())
	// case3: no memdb to flush
	index.EXPECT().Flush().Return(nil)
	assert.NoError(t, s2.Flush())

	index.EXPECT().Flush().Return(nil).AnyTimes()
	immutableDB1 := memdb.NewMockMemoryDatabase(ctrl)
	immutableDB1.EXPECT().FamilyTime().Return(int64(1)).AnyTimes()
	immutableDB1.EXPECT().MemSize().Return(int64(1000)).AnyTimes()
	immutableDB1.EXPECT().Close().Return(nil).AnyTimes()
	mutableDB2 := memdb.NewMockMemoryDatabase(ctrl)
	mutableDB2.EXPECT().FamilyTime().Return(int64(2)).AnyTimes()
	mutableDB2.EXPECT().MemSize().Return(int64(1000)).AnyTimes()
	mutableDB2.EXPECT().Close().Return(nil).AnyTimes()

	s2.families.InsertFamily(1, immutableDB1)
	s2.families.InsertFamily(2, mutableDB2)
	s2.families.SetFamilyImmutable(1)
	// case4: flush first immutable error
	immutableDB1.EXPECT().FlushFamilyTo(gomock.Any()).Return(fmt.Errorf("error"))
	assert.Error(t, s2.Flush())

	// case5, flush first immutable ok
	immutableDB1.EXPECT().FlushFamilyTo(gomock.Any()).Return(nil)
	assert.NoError(t, s2.Flush())

	// case6, move mutable to immutable
	mutableDB2.EXPECT().FlushFamilyTo(gomock.Any()).Return(nil)
	assert.NoError(t, s2.Flush())
}

func TestShard_NeedFlush(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	defer func() {
		_ = fileutil.RemoveDir(testPath)
	}()

	meta := metadb.NewMockMetadata(ctrl)
	meta.EXPECT().DatabaseName().Return("test").AnyTimes()
	db := NewMockDatabase(ctrl)
	db.EXPECT().Name().Return("test-db").AnyTimes()
	db.EXPECT().Metadata().Return(meta).AnyTimes()
	s, _ := newShard(db, 2, _testShard1Path, option.DatabaseOption{Interval: "10s"})
	index := indexdb.NewMockIndexDatabase(ctrl)
	s3 := s.(*shard)
	s3.indexDB = index
	// case 1: flush doing
	s3.isFlushing.Store(true)
	assert.False(t, s3.NeedFlush())
	s3.isFlushing.Store(false)

	db1 := memdb.NewMockMemoryDatabase(ctrl)
	db1.EXPECT().MemSize().Return(int64(config.GlobalStorageConfig().TSDB.MaxMemDBSize)).AnyTimes()
	db1.EXPECT().Uptime().Return(time.Second).AnyTimes()
	db2 := memdb.NewMockMemoryDatabase(ctrl)
	db2.EXPECT().MemSize().Return(int64(2)).AnyTimes()
	db2.EXPECT().Uptime().Return(time.Hour * 24).AnyTimes()
	db3 := memdb.NewMockMemoryDatabase(ctrl)
	db3.EXPECT().Uptime().Return(time.Second).AnyTimes()
	db3.EXPECT().MemSize().Return(int64(3)).AnyTimes()
	db4 := memdb.NewMockMemoryDatabase(ctrl)
	db4.EXPECT().Uptime().Return(time.Second).AnyTimes()
	db4.EXPECT().MemSize().Return(int64(4)).AnyTimes()
	db5 := memdb.NewMockMemoryDatabase(ctrl)
	db5.EXPECT().Uptime().Return(time.Second).AnyTimes()
	db5.EXPECT().MemSize().Return(int64(5)).AnyTimes()
	db6 := memdb.NewMockMemoryDatabase(ctrl)
	db6.EXPECT().Uptime().Return(time.Second).AnyTimes()
	db6.EXPECT().MemSize().Return(int64(config.GlobalStorageConfig().TSDB.MaxMemDBTotalSize) + 10000).AnyTimes()

	s3.families.InsertFamily(1, db1)
	s3.families.SetFamilyImmutable(1)
	s3.families.InsertFamily(2, db2)
	s3.families.SetFamilyImmutable(2)
	s3.families.InsertFamily(3, db3)
	s3.families.InsertFamily(4, db4)
	s3.families.InsertFamily(5, db5)
	s3.families.InsertFamily(6, db6)

	// case 2: too many memdbs
	assert.True(t, s3.NeedFlush())

	// case3, size too much
	s3.families.RemoveHeadImmutable()
	s3.families.RemoveHeadImmutable()
	assert.True(t, s3.NeedFlush())

	// case4, no need to flush
	s3.families.SetFamilyImmutable(3)
	s3.families.SetFamilyImmutable(4)
	s3.families.SetFamilyImmutable(5)
	s3.families.SetFamilyImmutable(6)
	s3.families.RemoveHeadImmutable()
	s3.families.RemoveHeadImmutable()
	s3.families.RemoveHeadImmutable()
	s3.families.RemoveHeadImmutable()
	assert.False(t, s3.NeedFlush())
}

func TestShard_GetOrCreateMemoryDatabase(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer func() {
		_ = fileutil.RemoveDir(testPath)
		newMemoryDBFunc = memdb.NewMemoryDatabase
		ctrl.Finish()
	}()
	newMemoryDBFunc = func(cfg memdb.MemoryDatabaseCfg) (memdb.MemoryDatabase, error) {
		return nil, nil
	}
	meta := metadb.NewMockMetadata(ctrl)
	meta.EXPECT().DatabaseName().Return("test").AnyTimes()
	db := NewMockDatabase(ctrl)
	db.EXPECT().Name().Return("test-db").AnyTimes()
	db.EXPECT().Metadata().Return(meta).AnyTimes()
	s, _ := newShard(db, 1, _testShard1Path, option.DatabaseOption{Interval: "10s"})
	s4 := s.(*shard)

	_, err := s4.GetOrCreateMemoryDatabase(1)
	assert.Nil(t, err)

	newMemoryDBFunc = func(cfg memdb.MemoryDatabaseCfg) (memdb.MemoryDatabase, error) {
		return nil, fmt.Errorf("err")
	}
	_, err = s4.GetOrCreateMemoryDatabase(2)
	assert.Error(t, err)
}
