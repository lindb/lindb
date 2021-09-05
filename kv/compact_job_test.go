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

package kv

import (
	"fmt"
	"sort"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/kv/table"
	"github.com/lindb/lindb/kv/version"
)

type mockAppendMerger struct {
	flusher Flusher
}

func newMockAppendMerger(flusher Flusher) (Merger, error) {
	return &mockAppendMerger{flusher: flusher}, nil
}

func (m *mockAppendMerger) Init(_ map[string]interface{}) {}

func (m *mockAppendMerger) Merge(key uint32, values [][]byte) error {
	var result []byte
	for _, v := range values {
		result = append(result, v...)
	}
	return m.flusher.Add(key, result)
}

func TestCompactJob_move_compact(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	snapshot := version.NewMockSnapshot(ctrl)
	merge := NewMockMerger(ctrl)
	family := generateMockFamily(ctrl, func(flusher Flusher) (Merger, error) {
		return merge, nil
	})
	family.EXPECT().familyInfo().Return("family").AnyTimes()
	f1 := version.NewFileMeta(1, 1, 100, 100)
	compaction := version.NewCompaction(1, 0, []*version.FileMeta{f1}, nil)
	state := newCompactionState(1000, snapshot, compaction)
	compact := newCompactJob(family, state, nil)
	err := compact.Run()
	assert.NoError(t, err)
	if err != nil {
		t.Fatal(err)
	}
	compact2 := compact.(*compactJob)
	logs := compact2.state.compaction.GetEditLog().GetLogs()
	assert.Equal(t, 2, len(logs))
	assert.Equal(t, version.NewDeleteFile(0, 1), logs[0])
	assert.Equal(t, version.CreateNewFile(1, f1), logs[1])
}

func TestCompactJob_merge_compact_get_read_fail(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	snapshot := version.NewMockSnapshot(ctrl)
	snapshot.EXPECT().GetReader(gomock.Any()).Return(nil, fmt.Errorf("err"))
	merge := NewMockMerger(ctrl)
	family := generateMockFamily(ctrl, func(flusher Flusher) (Merger, error) {
		return merge, nil
	})
	family.EXPECT().familyInfo().Return("family").AnyTimes()
	f1 := version.NewFileMeta(1, 1, 10, 100)
	f2 := version.NewFileMeta(2, 20, 50, 100)
	f3 := version.NewFileMeta(3, 1, 30, 100)
	f4 := version.NewFileMeta(4, 30, 100, 100)
	compaction := version.NewCompaction(1, 0, []*version.FileMeta{f1, f2}, []*version.FileMeta{f3, f4})
	state := newCompactionState(1000, snapshot, compaction)
	compactJob := newCompactJob(family, state, nil)
	err := compactJob.Run()
	assert.NotNil(t, err)
}

func TestCompactJob_merge_compact_merge_fail(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	snapshot := version.NewMockSnapshot(ctrl)
	reader := table.NewMockReader(ctrl)
	gomock.InOrder(
		reader.EXPECT().Iterator().Return(generateIterator(ctrl, map[uint32][]byte{
			1: []byte("value1"),
		})),
		reader.EXPECT().Iterator().Return(generateIterator(ctrl, map[uint32][]byte{})),
	)
	snapshot.EXPECT().GetReader(gomock.Any()).Return(reader, nil).MaxTimes(2)
	merge := NewMockMerger(ctrl)
	merge.EXPECT().Merge(gomock.Any(), gomock.Any()).Return(fmt.Errorf("err"))
	family := generateMockFamily(ctrl, func(flusher Flusher) (Merger, error) {
		return merge, nil
	})

	family.EXPECT().familyInfo().Return("family").AnyTimes()
	f1 := version.NewFileMeta(1, 1, 10, 100)
	f4 := version.NewFileMeta(4, 30, 100, 100)
	compaction := version.NewCompaction(1, 0, []*version.FileMeta{f1}, []*version.FileMeta{f4})
	state := newCompactionState(1000, snapshot, compaction)
	compactJob := newCompactJob(family, state, nil)
	err := compactJob.Run()
	assert.NotNil(t, err)

	gomock.InOrder(
		reader.EXPECT().Iterator().Return(generateIterator(ctrl, map[uint32][]byte{
			1: []byte("value1"),
		})),
		reader.EXPECT().Iterator().Return(generateIterator(ctrl, map[uint32][]byte{2: []byte("value2")})),
	)
	snapshot.EXPECT().GetReader(gomock.Any()).Return(reader, nil).MaxTimes(2)
	merge.EXPECT().Merge(uint32(1), gomock.Any()).Return(fmt.Errorf("err"))
	compactJob = newCompactJob(family, state, nil)
	err = compactJob.Run()
	assert.NotNil(t, err)
}

func TestCompactJob_merge_doMerge_fail(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	snapshot := version.NewMockSnapshot(ctrl)
	reader1 := table.NewMockReader(ctrl)
	reader2 := table.NewMockReader(ctrl)
	merge := NewMockMerger(ctrl)

	// test new store build fail
	gomock.InOrder(
		reader1.EXPECT().Iterator().Return(generateIterator(ctrl, map[uint32][]byte{
			1: []byte("value1"),
		})),
		reader2.EXPECT().Iterator().Return(generateIterator(ctrl, map[uint32][]byte{
			2: []byte("value2"),
		})),
	)
	gomock.InOrder(
		merge.EXPECT().Merge(gomock.Any(), gomock.Any()).Return(nil),
		merge.EXPECT().Merge(gomock.Any(), gomock.Any()).Return(nil),
	)
	snapshot.EXPECT().GetReader(table.FileNumber(1)).Return(reader1, nil)
	snapshot.EXPECT().GetReader(table.FileNumber(4)).Return(reader2, nil)
	family := generateMockFamily(ctrl, func(flusher Flusher) (Merger, error) {
		return merge, nil
	})
	family.EXPECT().familyInfo().Return("family").AnyTimes()

	f1 := version.NewFileMeta(1, 1, 10, 100)
	f4 := version.NewFileMeta(4, 30, 100, 100)
	compaction := version.NewCompaction(1, 0, []*version.FileMeta{f1}, []*version.FileMeta{f4})
	state := newCompactionState(10000, snapshot, compaction)
	compactJobIntf := newCompactJob(family, state, nil)
	err := compactJobIntf.Run()
	assert.NoError(t, err)
}

func TestCompactJob_output_fail(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	snapshot := version.NewMockSnapshot(ctrl)
	reader1 := table.NewMockReader(ctrl)
	reader2 := table.NewMockReader(ctrl)
	merge := NewMockMerger(ctrl)

	// test store build is empty
	gomock.InOrder(
		reader1.EXPECT().Iterator().Return(generateIterator(ctrl, map[uint32][]byte{
			1: []byte("value1"),
		})),
		reader2.EXPECT().Iterator().Return(generateIterator(ctrl, map[uint32][]byte{
			1: []byte("value2"),
		})),
	)
	snapshot.EXPECT().GetReader(table.FileNumber(1)).Return(reader1, nil)
	snapshot.EXPECT().GetReader(table.FileNumber(4)).Return(reader2, nil)
	family := generateMockFamily(ctrl, func(flusher Flusher) (Merger, error) {
		return merge, nil
	})
	family.EXPECT().familyInfo().Return("family").AnyTimes()
	builder := table.NewMockBuilder(ctrl)
	gomock.InOrder(
		family.EXPECT().newTableBuilder().Return(builder, nil),
		builder.EXPECT().FileNumber().Return(table.FileNumber(10)),
		family.EXPECT().addPendingOutput(table.FileNumber(10)),
		builder.EXPECT().Add(uint32(1), []byte{1, 2, 3}).Return(nil),
		builder.EXPECT().Size().Return(uint32(100)),
		// no output
		builder.EXPECT().Count().Return(uint64(0)),
	)
	f1 := version.NewFileMeta(1, 1, 10, 100)
	f4 := version.NewFileMeta(4, 30, 100, 100)
	compaction := version.NewCompaction(1, 0, []*version.FileMeta{f1}, []*version.FileMeta{f4})
	state := newCompactionState(1000, snapshot, compaction)
	compact := newCompactJob(family, state, nil)
	merge.EXPECT().Merge(uint32(1), gomock.Any()).DoAndReturn(func(key uint32, value [][]byte) error {
		_ = compact.(*compactJob).newCompactFlusher().Add(key, []byte{1, 2, 3})
		return nil
	})
	err := compact.Run()
	assert.NoError(t, err)
	assert.Equal(t, 2, len(state.compaction.GetEditLog().GetLogs()))
	assert.Equal(t, 0, len(state.outputs))

	// test finish output fail
	gomock.InOrder(
		reader1.EXPECT().Iterator().Return(generateIterator(ctrl, map[uint32][]byte{
			1: []byte("value1"),
		})),
		reader2.EXPECT().Iterator().Return(generateIterator(ctrl, map[uint32][]byte{
			1: []byte("value2"),
		})),
	)
	snapshot.EXPECT().GetReader(table.FileNumber(1)).Return(reader1, nil)
	snapshot.EXPECT().GetReader(table.FileNumber(4)).Return(reader2, nil)
	gomock.InOrder(
		family.EXPECT().newTableBuilder().Return(builder, nil),
		builder.EXPECT().FileNumber().Return(table.FileNumber(10)),
		family.EXPECT().addPendingOutput(table.FileNumber(10)),
		builder.EXPECT().Add(uint32(1), []byte{1, 2, 3}).Return(nil),
		builder.EXPECT().Size().Return(uint32(100)),
		builder.EXPECT().Count().Return(uint64(10)),
		builder.EXPECT().Close().Return(fmt.Errorf("err")),
		builder.EXPECT().Abandon().Return(fmt.Errorf("err")),
		family.EXPECT().removePendingOutput(table.FileNumber(10)),
	)
	state = newCompactionState(1000, snapshot, compaction)
	compact = newCompactJob(family, state, nil)
	merge.EXPECT().Merge(uint32(1), gomock.Any()).DoAndReturn(func(key uint32, value [][]byte) error {
		_ = compact.(*compactJob).newCompactFlusher().Add(key, []byte{1, 2, 3})
		return nil
	})
	err = compact.Run()
	assert.NotNil(t, err)

	// test build is nil, when finish output
	state = newCompactionState(1000, snapshot, compaction)
	compact = newCompactJob(family, state, nil)
	compact2 := compact.(*compactJob)
	err = compact2.finishCompactionOutputFile()
	assert.NotNil(t, err)
}

func TestCompactJob_finish_output_fail(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	snapshot := version.NewMockSnapshot(ctrl)
	family := NewMockFamily(ctrl)
	family.EXPECT().getNewMerger().Return(nil)
	compaction := version.NewCompaction(1, 0, nil, nil)
	state := newCompactionState(1000, snapshot, compaction)
	compact := newCompactJob(family, state, nil)
	compact2 := compact.(*compactJob)
	err := compact2.finishCompactionOutputFile()
	assert.NotNil(t, err)
}

func TestCompactJob_merge_compact(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	snapshot := version.NewMockSnapshot(ctrl)
	reader1 := table.NewMockReader(ctrl)
	reader2 := table.NewMockReader(ctrl)
	reader3 := table.NewMockReader(ctrl)
	reader4 := table.NewMockReader(ctrl)
	reader1.EXPECT().Iterator().Return(generateIterator(ctrl, map[uint32][]byte{
		1:  []byte("value1"),
		3:  []byte("value3"),
		10: []byte("value10"),
	}))
	reader2.EXPECT().Iterator().Return(generateIterator(ctrl, map[uint32][]byte{
		10: []byte("value10"),
		30: []byte("value30"),
		40: []byte("value40"),
	}))
	reader3.EXPECT().Iterator().Return(generateIterator(ctrl, map[uint32][]byte{
		1:  []byte("value1"),
		10: []byte("value10"),
	}))
	reader4.EXPECT().Iterator().Return(generateIterator(ctrl, map[uint32][]byte{
		10:  []byte("value10"),
		30:  []byte("value30"),
		100: []byte("value100"),
	}))
	snapshot.EXPECT().GetReader(table.FileNumber(1)).Return(reader1, nil)
	snapshot.EXPECT().GetReader(table.FileNumber(2)).Return(reader2, nil)
	snapshot.EXPECT().GetReader(table.FileNumber(3)).Return(reader3, nil)
	snapshot.EXPECT().GetReader(table.FileNumber(4)).Return(reader4, nil)
	family := generateMockFamily(ctrl, newMockAppendMerger)
	family.EXPECT().familyInfo().Return("family").AnyTimes()
	f1 := version.NewFileMeta(1, 1, 10, 100)
	f2 := version.NewFileMeta(2, 10, 50, 100)
	f3 := version.NewFileMeta(3, 1, 30, 100)
	f4 := version.NewFileMeta(4, 30, 100, 100)
	compaction := version.NewCompaction(1, 0, []*version.FileMeta{f1, f2}, []*version.FileMeta{f3, f4})
	state := newCompactionState(10000000, snapshot, compaction)
	compactJob := newCompactJob(family, state, NewMockRollup(ctrl))
	builder := table.NewMockBuilder(ctrl)
	gomock.InOrder(
		family.EXPECT().newTableBuilder().Return(builder, nil),
		builder.EXPECT().FileNumber().Return(table.FileNumber(5)),
		family.EXPECT().addPendingOutput(table.FileNumber(5)),
		builder.EXPECT().Add(uint32(1), []byte("value1value1")).Return(nil),
		builder.EXPECT().Size().Return(uint32(10)),
		builder.EXPECT().Add(uint32(3), []byte("value3")).Return(nil),
		builder.EXPECT().Size().Return(uint32(10)),
		builder.EXPECT().Add(uint32(10), []byte("value10value10value10value10")).Return(nil),
		builder.EXPECT().Size().Return(uint32(10)),
		builder.EXPECT().Add(uint32(30), []byte("value30value30")).Return(nil),
		builder.EXPECT().Size().Return(uint32(10)),
		builder.EXPECT().Add(uint32(40), []byte("value40")).Return(nil),
		builder.EXPECT().Size().Return(uint32(10)),
		builder.EXPECT().Add(uint32(100), []byte("value100")).Return(nil),
		builder.EXPECT().Size().Return(uint32(10)),
		builder.EXPECT().Count().Return(uint64(6)),
		builder.EXPECT().Close().Return(nil),
		builder.EXPECT().FileNumber().Return(table.FileNumber(5)),
		builder.EXPECT().MinKey().Return(uint32(1)),
		builder.EXPECT().MaxKey().Return(uint32(100)),
		builder.EXPECT().Size().Return(uint32(10)),
		family.EXPECT().removePendingOutput(table.FileNumber(5)),
	)
	err := compactJob.Run()
	assert.NoError(t, err)
	newFile := version.NewFileMeta(table.FileNumber(5), uint32(1), uint32(100), uint32(10))
	assert.Equal(t, 1, len(state.outputs))
	assert.Equal(t, *newFile, *(state.outputs[0]))
	editLog := state.compaction.GetEditLog()
	logs := editLog.GetLogs()
	assert.Equal(t, 5, len(logs))
	assert.Equal(t, version.NewDeleteFile(0, 1), logs[0])
	assert.Equal(t, version.NewDeleteFile(0, 2), logs[1])
	assert.Equal(t, version.NewDeleteFile(1, 3), logs[2])
	assert.Equal(t, version.NewDeleteFile(1, 4), logs[3])
	assert.Equal(t, version.CreateNewFile(1, newFile), logs[4])
}

func generateMockFamily(ctrl *gomock.Controller, merger NewMerger) *MockFamily {
	family := NewMockFamily(ctrl)
	family.EXPECT().getNewMerger().Return(merger).AnyTimes()
	family.EXPECT().Name().Return("test-family").AnyTimes()
	family.EXPECT().commitEditLog(gomock.Any()).Return(true).AnyTimes()
	return family
}

func generateIterator(ctrl *gomock.Controller, values map[uint32][]byte) table.Iterator {
	it1 := table.NewMockIterator(ctrl)
	var keys []uint32
	for key := range values {
		keys = append(keys, key)
	}
	sort.Slice(keys, func(i, j int) bool {
		return keys[i] < keys[j]
	})
	var calls []*gomock.Call
	for _, key := range keys {
		calls = append(calls,
			it1.EXPECT().HasNext().Return(true),
			it1.EXPECT().Key().Return(key),
			it1.EXPECT().Value().Return(values[key]))
	}
	calls = append(calls, it1.EXPECT().HasNext().Return(false))

	gomock.InOrder(calls...)
	return it1
}
