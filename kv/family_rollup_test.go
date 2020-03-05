package kv

import (
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/kv/table"
	"github.com/lindb/lindb/kv/version"
	"github.com/lindb/lindb/pkg/fileutil"
	"github.com/lindb/lindb/pkg/timeutil"
)

func TestFamily_needRollup(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer func() {
		_ = fileutil.RemoveDir(testKVPath)
		ctrl.Finish()
	}()
	store := NewMockStore(ctrl)
	store.EXPECT().Option().Return(DefaultStoreOption(testKVPath)).AnyTimes()
	fv := version.NewMockFamilyVersion(ctrl)
	snapshot := version.NewMockSnapshot(ctrl)
	v := version.NewMockVersion(ctrl)
	snapshot.EXPECT().Close().AnyTimes()
	snapshot.EXPECT().GetCurrent().Return(v).AnyTimes()
	fv.EXPECT().GetSnapshot().Return(snapshot).AnyTimes()
	store.EXPECT().createFamilyVersion(gomock.Any(), gomock.Any()).Return(fv)
	f, err := newFamily(store, FamilyOption{Merger: "mockMerger"})
	assert.NoError(t, err)
	f2 := f.(*family)
	// case 1: rollup job doing
	f2.rolluping.Store(true)
	assert.False(t, f2.needRollup())
	// case 2: rollup files nil
	f2.rolluping.Store(false)
	fv.EXPECT().GetLiveRollupFiles().Return(nil)
	assert.False(t, f2.needRollup())
	// case 3: rollup files < threshold
	fv.EXPECT().GetLiveRollupFiles().Return(map[table.FileNumber]timeutil.Interval{10: 10})
	assert.False(t, f2.needRollup())
	// case 4: need rollup
	fv.EXPECT().GetLiveRollupFiles().Return(map[table.FileNumber]timeutil.Interval{10: 10, 11: 10, 12: 10, 13: 10})
	assert.True(t, f2.needRollup())
}

func TestFamily_rollup(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer func() {
		_ = fileutil.RemoveDir(testKVPath)
		ctrl.Finish()
	}()
	store := NewMockStore(ctrl)
	store.EXPECT().Option().Return(DefaultStoreOption(testKVPath)).AnyTimes()
	fv := version.NewMockFamilyVersion(ctrl)
	snapshot := version.NewMockSnapshot(ctrl)
	v := version.NewMockVersion(ctrl)
	snapshot.EXPECT().Close().AnyTimes()
	snapshot.EXPECT().GetCurrent().Return(v).AnyTimes()
	fv.EXPECT().GetSnapshot().Return(snapshot).AnyTimes()
	fv.EXPECT().GetAllActiveFiles().Return(nil).AnyTimes()
	store.EXPECT().createFamilyVersion(gomock.Any(), gomock.Any()).Return(fv)
	f, err := newFamily(store, FamilyOption{Merger: "mockMerger"})
	assert.NoError(t, err)
	f2 := f.(*family)
	// case 1: rollup doing
	f2.rolluping.Store(true)
	f2.rollup()
	// case 2: get rollup files nil
	f2.rolluping.Store(false)
	fv.EXPECT().GetLiveRollupFiles().Return(nil).MaxTimes(2)
	f2.rollup()
	// case 3: rollup relation not found
	fv.EXPECT().GetLiveRollupFiles().Return(map[table.FileNumber]timeutil.Interval{10: 10}).MaxTimes(2)
	store.EXPECT().getRollup(timeutil.Interval(10)).Return(nil, false)
	f2.rollup()
	// case 4: do rollup err
	fv.EXPECT().GetLiveRollupFiles().Return(map[table.FileNumber]timeutil.Interval{10: 10}).MaxTimes(2)
	rollup := NewMockRollup(ctrl)
	tf := NewMockFamily(ctrl)
	rollup.EXPECT().GetTargetFamily(gomock.Any()).Return(tf).AnyTimes()
	store.EXPECT().getRollup(timeutil.Interval(10)).Return(rollup, true).AnyTimes()
	tf.EXPECT().doRollupWork(f2, rollup, []table.FileNumber{10}).Return(fmt.Errorf("err"))
	f2.rollup()
	// case 5: rollup success
	fv.EXPECT().GetLiveRollupFiles().Return(map[table.FileNumber]timeutil.Interval{10: 10}).MaxTimes(2)
	tf.EXPECT().doRollupWork(f2, rollup, []table.FileNumber{10}).Return(nil)
	store.EXPECT().commitFamilyEditLog(gomock.Any(), gomock.Any()).Return(nil)
	f2.rollup()
}

func TestFamily_doRollupWork(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer func() {
		_ = fileutil.RemoveDir(testKVPath)
		newCompactJobFunc = newCompactJob
		ctrl.Finish()
	}()
	store := NewMockStore(ctrl)
	store.EXPECT().Option().Return(DefaultStoreOption(testKVPath)).AnyTimes()
	fv := version.NewMockFamilyVersion(ctrl)
	store.EXPECT().createFamilyVersion(gomock.Any(), gomock.Any()).Return(fv)
	f, err := newFamily(store, FamilyOption{Merger: "mockMerger"})
	assert.NoError(t, err)
	f2 := f.(*family)
	// case 1: source file nil
	err = f2.doRollupWork(nil, nil, nil)
	assert.NoError(t, err)
	// case 2: source files already rollup
	fv.EXPECT().GetLiveReferenceFiles().Return(map[version.FamilyID][]table.FileNumber{10: {10, 20, 30}})
	sf := NewMockFamily(ctrl)
	sf.EXPECT().ID().Return(version.FamilyID(10)).AnyTimes()
	sf.EXPECT().familyInfo().Return("family").AnyTimes()
	err = f2.doRollupWork(sf, nil, []table.FileNumber{10, 20, 30})
	assert.NoError(t, err)
	// case 3: rollup source files
	fv.EXPECT().GetLiveReferenceFiles().Return(map[version.FamilyID][]table.FileNumber{10: {10, 30}}).AnyTimes()
	snapshot := version.NewMockSnapshot(ctrl)
	snapshot.EXPECT().Close().AnyTimes()
	sf.EXPECT().GetSnapshot().Return(snapshot).AnyTimes()
	err = f2.doRollupWork(sf, nil, []table.FileNumber{10, 20, 30})
	assert.NoError(t, err)
	// case 4: rollup job err
	compactJob := NewMockCompactJob(ctrl)
	newCompactJobFunc = func(family Family, state *compactionState, rollup Rollup) CompactJob {
		return compactJob
	}
	compactJob.EXPECT().Run().Return(fmt.Errorf("err"))
	err = f2.doRollupWork(sf, nil, []table.FileNumber{10, 20, 30})
	assert.Error(t, err)
}
