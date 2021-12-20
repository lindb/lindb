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
	"path/filepath"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/kv/table"
	"github.com/lindb/lindb/kv/version"
	"github.com/lindb/lindb/pkg/timeutil"
)

func TestRollup(t *testing.T) {
	t.Run("10s->5min", func(t *testing.T) {
		sf, _ := timeutil.ParseTimestamp("2019-12-12 10:00:00")
		tf, _ := timeutil.ParseTimestamp("2019-12-12 00:00:00")
		in := newRollup(timeutil.Interval(10*1000), timeutil.Interval(5*60*1000), sf, tf)
		assert.Equal(t, uint16(30), in.IntervalRatio())
		timestamp := in.GetTimestamp(20)
		assert.Equal(t, sf+10*1000*20, timestamp)
		assert.Equal(t, uint16(10*60/5), in.CalcSlot(timestamp))
	})
	t.Run("10s->1hour", func(t *testing.T) {
		sf, _ := timeutil.ParseTimestamp("2019-12-12 10:00:00")
		tf, _ := timeutil.ParseTimestamp("2019-12-12 00:00:00")
		in := newRollup(timeutil.Interval(10*1000), timeutil.Interval(60*60*1000), sf, tf)
		assert.Equal(t, uint16(360), in.IntervalRatio())
		timestamp := in.GetTimestamp(20)
		assert.Equal(t, sf+10*1000*20, timestamp)
		assert.Equal(t, uint16(10), in.CalcSlot(timestamp))
	})
}

func TestFamily_needRollup(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	f, store := mockFamily(t, ctrl)
	fv := f.familyVersion.(*version.MockFamilyVersion)

	cases := []struct {
		name    string
		prepare func()
		need    bool
	}{
		{
			name: "rollup job doing",
			prepare: func() {
				f.rolluping.Store(true)
			},
			need: false,
		},
		{
			name: "no set rollup option",
			prepare: func() {
				store.EXPECT().Option().Return(StoreOption{})
			},
			need: false,
		},
		{
			name: "no live rollup file",
			prepare: func() {
				store.EXPECT().Option().Return(StoreOption{Rollup: []timeutil.Interval{10}})
				fv.EXPECT().GetLiveRollupFiles().Return(nil)
			},
			need: false,
		},
		{
			name: "rollup files < threshold",
			prepare: func() {
				store.EXPECT().Option().Return(StoreOption{Rollup: []timeutil.Interval{10}})
				fv.EXPECT().GetLiveRollupFiles().Return(map[table.FileNumber][]timeutil.Interval{10: {10}})
			},
			need: false,
		},
		{
			name: "need rollup",
			prepare: func() {
				store.EXPECT().Option().Return(StoreOption{Rollup: []timeutil.Interval{10}})
				fv.EXPECT().GetLiveRollupFiles().Return(
					map[table.FileNumber][]timeutil.Interval{
						10: {10}, 11: {10}, 12: {10}, 13: {10},
					})
			},
			need: true,
		},
	}

	for _, tt := range cases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			defer f.rolluping.Store(false)
			if tt.prepare != nil {
				tt.prepare()
			}
			need := f.needRollup()
			assert.True(t, need == tt.need)
		})
	}
}

func TestFamily_rollup(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer func() {
		Options.Store(&StoreOptions{})
		InitStoreManager(nil)
	}()

	Options.Store(&StoreOptions{Dir: t.TempDir()})
	targetStore := NewMockStore(ctrl)
	targetFamily := NewMockFamily(ctrl)
	storeMgr := NewMockStoreManager(ctrl)
	InitStoreManager(storeMgr)
	f, store := mockFamily(t, ctrl)
	fv := f.familyVersion.(*version.MockFamilyVersion)
	name := "13"
	cases := []struct {
		name    string
		prepare func()
	}{
		{
			name: "rollup job doing",
			prepare: func() {
				f.rolluping.Store(true)
			},
		},
		{
			name: "no live rollup file",
			prepare: func() {
				gomock.InOrder(
					fv.EXPECT().GetLiveRollupFiles().Return(nil),
					fv.EXPECT().GetAllActiveFiles().Return(nil),
					fv.EXPECT().GetLiveRollupFiles().Return(nil),
				)
			},
		},
		{
			name: "parse segment name failure",
			prepare: func() {
				gomock.InOrder(
					fv.EXPECT().GetLiveRollupFiles().
						Return(map[table.FileNumber][]timeutil.Interval{10: {10000}, 11: {10000}}),
					store.EXPECT().Option().Return(StoreOption{Source: timeutil.Interval(10000)}),
					store.EXPECT().Name().Return("xxx"),
					fv.EXPECT().GetAllActiveFiles().Return(nil),
					fv.EXPECT().GetLiveRollupFiles().Return(nil),
				)
			},
		},
		{
			name: "parse family name failure",
			prepare: func() {
				f.name = "aa"
				gomock.InOrder(
					fv.EXPECT().GetLiveRollupFiles().
						Return(map[table.FileNumber][]timeutil.Interval{10: {10000}, 11: {10000}}),
					store.EXPECT().Option().Return(StoreOption{Source: timeutil.Interval(10000)}),
					store.EXPECT().Name().Return("20190703"),
					fv.EXPECT().GetAllActiveFiles().Return(nil),
					fv.EXPECT().GetLiveRollupFiles().Return(nil),
				)
			},
		},
		{
			name: "target store not found",
			prepare: func() {
				f.name =
					name
				gomock.InOrder(
					fv.EXPECT().GetLiveRollupFiles().
						Return(map[table.FileNumber][]timeutil.Interval{10: {5 * 60 * 1000}, 11: {5 * 60 * 1000}}),
					store.EXPECT().Option().Return(StoreOption{Source: timeutil.Interval(10000)}),
					store.EXPECT().Name().Return("db/shard/1/segment/day/20190703"),
					storeMgr.EXPECT().GetStoreByName("db/shard/1/segment/month/201907").Return(nil, false),
					fv.EXPECT().GetAllActiveFiles().Return(nil),
					fv.EXPECT().GetLiveRollupFiles().Return(nil),
				)
			},
		},
		{
			name: "create target family failure",
			prepare: func() {
				f.name =
					name
				gomock.InOrder(
					fv.EXPECT().GetLiveRollupFiles().Return(map[table.FileNumber][]timeutil.Interval{
						10: {5 * 60 * 1000}, 11: {5 * 60 * 1000},
					}),
					store.EXPECT().Option().Return(StoreOption{Source: timeutil.Interval(10000)}),
					store.EXPECT().Name().Return("db/shard/1/segment/day/20190703"),
					storeMgr.EXPECT().GetStoreByName("db/shard/1/segment/month/201907").Return(targetStore, true),
					targetStore.EXPECT().CreateFamily("3", gomock.Any()).
						Return(nil, fmt.Errorf("err")),
					fv.EXPECT().GetAllActiveFiles().Return(nil),
					fv.EXPECT().GetLiveRollupFiles().Return(nil),
				)
			},
		},
		{
			name: "do rollup job err",
			prepare: func() {
				f.name =
					name
				gomock.InOrder(
					fv.EXPECT().GetLiveRollupFiles().Return(map[table.FileNumber][]timeutil.Interval{
						10: {5 * 60 * 1000}, 11: {5 * 60 * 1000},
					}),
					store.EXPECT().Option().Return(StoreOption{Source: timeutil.Interval(10000)}),
					store.EXPECT().Name().Return("db/shard/1/segment/day/20190703"),
					storeMgr.EXPECT().GetStoreByName("db/shard/1/segment/month/201907").Return(targetStore, true),
					targetStore.EXPECT().CreateFamily("3", gomock.Any()).Return(targetFamily, nil),
					targetFamily.EXPECT().doRollupWork(gomock.Any(), gomock.Any(), gomock.Any()).
						Return(fmt.Errorf("err")),
					fv.EXPECT().GetAllActiveFiles().Return(nil),
					fv.EXPECT().GetLiveRollupFiles().Return(nil),
				)
			},
		},
		{
			name: "do rollup job successfully",
			prepare: func() {
				f.name =
					name
				gomock.InOrder(
					fv.EXPECT().GetLiveRollupFiles().Return(map[table.FileNumber][]timeutil.Interval{
						10: {5 * 60 * 1000}, 11: {5 * 60 * 1000},
					}),
					store.EXPECT().Option().Return(StoreOption{Source: timeutil.Interval(10000)}),
					store.EXPECT().Name().Return("db/shard/1/segment/day/20190703"),
					storeMgr.EXPECT().GetStoreByName("db/shard/1/segment/month/201907").Return(targetStore, true),
					targetStore.EXPECT().CreateFamily("3", gomock.Any()).Return(targetFamily, nil),
					targetFamily.EXPECT().doRollupWork(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil),
					store.EXPECT().commitFamilyEditLog(gomock.Any(), gomock.Any()).Return(nil),
					fv.EXPECT().GetAllActiveFiles().Return(nil),
					fv.EXPECT().GetLiveRollupFiles().Return(nil),
				)
			},
		},
		{
			name: "do rollup job ok, but some targets failure",
			prepare: func() {
				f.name =
					name
				fv.EXPECT().GetLiveRollupFiles().Return(map[table.FileNumber][]timeutil.Interval{
					10: {5 * 60 * 1000}, 11: {5 * 60 * 60 * 1000}, // year target not found
				})
				store.EXPECT().Option().Return(StoreOption{Source: timeutil.Interval(10000)})
				store.EXPECT().Name().Return("db/shard/1/segment/day/20190703")
				storeMgr.EXPECT().GetStoreByName("db/shard/1/segment/month/201907").Return(targetStore, true)
				targetStore.EXPECT().CreateFamily("3", gomock.Any()).Return(targetFamily, nil)
				targetFamily.EXPECT().doRollupWork(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
				storeMgr.EXPECT().GetStoreByName("db/shard/1/segment/year/2019").Return(nil, false) // year store not found
				store.EXPECT().commitFamilyEditLog(gomock.Any(), gomock.Any()).Return(nil)
				fv.EXPECT().GetAllActiveFiles().Return(nil)
				fv.EXPECT().GetLiveRollupFiles().Return(nil)
			},
		},
	}

	for _, tt := range cases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			defer f.rolluping.Store(false)
			if tt.prepare != nil {
				tt.prepare()
			}
			f.rollup()
			time.Sleep(100 * time.Millisecond) // for waiting job completed
		})
	}
}

func TestFamily_doRollupWork(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer func() {
		ctrl.Finish()
	}()
	f, _ := mockFamily(t, ctrl)
	sourceFamily := NewMockFamily(ctrl)
	fv := f.familyVersion.(*version.MockFamilyVersion)
	rollup := NewMockRollup(ctrl)
	cases := []struct {
		name    string
		files   []table.FileNumber
		prepare func()
		wantErr bool
	}{
		{
			name:    "source file is empty",
			files:   nil,
			prepare: nil,
			wantErr: false,
		},
		{
			name:  "source files already rollup",
			files: []table.FileNumber{10, 20, 30},
			prepare: func() {
				gomock.InOrder(
					fv.EXPECT().GetLiveReferenceFiles().
						Return(map[version.FamilyID][]table.FileNumber{10: {10, 20, 30}}),
					sourceFamily.EXPECT().ID().Return(version.FamilyID(10)),
					sourceFamily.EXPECT().familyInfo().Return("familyInfo").MaxTimes(3),
				)
			},
			wantErr: false,
		},
		{
			name:  "rollup source files",
			files: []table.FileNumber{10, 20, 30},
			prepare: func() {
				snapshot := version.NewMockSnapshot(ctrl)
				v := version.NewMockVersion(ctrl)
				compactJob := NewMockCompactJob(ctrl)
				newCompactJobFunc = func(family Family, state *compactionState, rollup Rollup) CompactJob {
					return compactJob
				}
				gomock.InOrder(
					fv.EXPECT().GetLiveReferenceFiles().Return(map[version.FamilyID][]table.FileNumber{10: {10, 30}}),
					sourceFamily.EXPECT().ID().Return(version.FamilyID(10)),
					sourceFamily.EXPECT().familyInfo().Return("familyInfo").MaxTimes(2),
					sourceFamily.EXPECT().GetSnapshot().Return(snapshot),
					snapshot.EXPECT().GetCurrent().Return(v),
					v.EXPECT().GetFile(0, table.FileNumber(20)).Return(nil, true),
					compactJob.EXPECT().Run().Return(nil),
					snapshot.EXPECT().Close(),
				)
			},
			wantErr: false,
		},
		{
			name:  "rollup job failure",
			files: []table.FileNumber{10, 20, 30},
			prepare: func() {
				snapshot := version.NewMockSnapshot(ctrl)
				v := version.NewMockVersion(ctrl)
				compactJob := NewMockCompactJob(ctrl)
				newCompactJobFunc = func(family Family, state *compactionState, rollup Rollup) CompactJob {
					return compactJob
				}
				gomock.InOrder(
					fv.EXPECT().GetLiveReferenceFiles().Return(map[version.FamilyID][]table.FileNumber{10: {10, 30}}),
					sourceFamily.EXPECT().ID().Return(version.FamilyID(10)),
					sourceFamily.EXPECT().familyInfo().Return("familyInfo").MaxTimes(2),
					sourceFamily.EXPECT().GetSnapshot().Return(snapshot),
					snapshot.EXPECT().GetCurrent().Return(v),
					v.EXPECT().GetFile(0, table.FileNumber(20)).Return(nil, true),
					compactJob.EXPECT().Run().Return(fmt.Errorf("err")),
					snapshot.EXPECT().Close(),
				)
			},
			wantErr: true,
		},
	}

	for _, tt := range cases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				newCompactJobFunc = newCompactJob
			}()
			if tt.prepare != nil {
				tt.prepare()
			}
			if err := f.doRollupWork(sourceFamily, rollup, tt.files); (err != nil) != tt.wantErr {
				t.Errorf("doRollupWork() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func mockFamily(t *testing.T, ctrl *gomock.Controller) (*family, *MockStore) {
	path := filepath.Join(t.TempDir(), "need_rollup")
	store := NewMockStore(ctrl)
	store.EXPECT().Path().Return(path)
	fv := version.NewMockFamilyVersion(ctrl)
	store.EXPECT().createFamilyVersion(gomock.Any(), gomock.Any()).Return(fv)
	f, err := newFamily(store, FamilyOption{Merger: "mockMerger"})
	assert.NoError(t, err)
	assert.NotNil(t, f)
	f2 := f.(*family)
	return f2, store
}
