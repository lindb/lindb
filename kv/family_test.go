package kv

import (
	"fmt"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/kv/table"
	"github.com/lindb/lindb/kv/version"
	"github.com/lindb/lindb/pkg/fileutil"
	"github.com/lindb/lindb/pkg/timeutil"
)

func init() {
	RegisterMerger("mockMerger", newMockMerger)
}

func TestFamily_New(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer func() {
		_ = fileutil.RemoveDir(testKVPath)
		mkDirFunc = fileutil.MkDir
		ctrl.Finish()
	}()
	store := NewMockStore(ctrl)
	store.EXPECT().Option().Return(DefaultStoreOption(testKVPath)).AnyTimes()
	// case 1: create family err, mkdir err
	mkDirFunc = func(path string) error {
		return fmt.Errorf("err")
	}
	f, err := newFamily(store, FamilyOption{Merger: "mockMerger"})
	assert.Error(t, err)
	assert.Nil(t, f)
	// case 2: create family err, merge not exist
	mkDirFunc = fileutil.MkDir
	f, err = newFamily(store, FamilyOption{Merger: "mockMerger_not_exist"})
	assert.Error(t, err)
	assert.Nil(t, f)
	// case 3: create family success
	vs := version.NewMockFamilyVersion(ctrl)
	store.EXPECT().createFamilyVersion(gomock.Any(), gomock.Any()).Return(vs)
	f, err = newFamily(store, FamilyOption{Merger: "mockMerger", ID: 10, Name: "f", MaxFileSize: 10})
	assert.NoError(t, err)
	assert.NotNil(t, f)
	assert.Equal(t, version.FamilyID(10), f.ID())
	assert.Equal(t, "f", f.Name())
	vs.EXPECT().GetSnapshot().Return(version.NewMockSnapshot(ctrl))
	assert.NotNil(t, f.GetSnapshot())
	assert.NotNil(t, f.NewFlusher())
	assert.NotNil(t, f.getFamilyVersion())
	assert.NotNil(t, f.getNewMerger())
}

func TestFamily_Data_Write_Read(t *testing.T) {
	option := DefaultStoreOption(testKVPath)
	defer func() {
		_ = fileutil.RemoveDir(testKVPath)
	}()

	var kv, err = NewStore("test_kv", option)
	defer func() {
		_ = kv.Close()
	}()
	assert.Nil(t, err, "cannot create kv store")

	f, err := kv.CreateFamily("f", FamilyOption{Merger: "mockMerger"})
	assert.Nil(t, err, "cannot create family")
	flusher := f.NewFlusher()
	_ = flusher.Add(1, []byte("test"))
	_ = flusher.Add(10, []byte("test10"))
	commitErr := flusher.Commit()
	assert.Nil(t, commitErr)

	snapshot := f.GetSnapshot()
	readers, _ := snapshot.FindReaders(10)
	assert.Equal(t, 1, len(readers))
	value, _ := readers[0].Get(1)
	assert.Equal(t, []byte("test"), value)
	value, _ = readers[0].Get(10)
	assert.Equal(t, []byte("test10"), value)
	snapshot.Close()
}

func TestFamily_commitEditLog(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer func() {
		_ = fileutil.RemoveDir(testKVPath)
		ctrl.Finish()
	}()
	store := NewMockStore(ctrl)
	store.EXPECT().Option().Return(DefaultStoreOption(testKVPath)).AnyTimes()
	fv := version.NewMockFamilyVersion(ctrl)
	store.EXPECT().createFamilyVersion(gomock.Any(), gomock.Any()).Return(fv)
	f, err := newFamily(store, FamilyOption{Merger: "mockMerger"})
	assert.NoError(t, err)
	// case 1: edit log is empty
	assert.True(t, f.commitEditLog(version.NewEditLog(10)))
	// case 2: edit log is nil
	assert.True(t, f.commitEditLog(nil))
	// case 3: commit edit log err
	editLog := version.NewEditLog(1)
	newFile := version.CreateNewFile(1, version.NewFileMeta(12, 1, 100, 2014))
	editLog.Add(newFile)
	editLog.Add(version.NewDeleteFile(1, 123))
	store.EXPECT().commitFamilyEditLog(gomock.Any(), gomock.Any()).Return(fmt.Errorf("err"))
	assert.False(t, f.commitEditLog(editLog))
	// case 4: commit edit log success
	store.EXPECT().commitFamilyEditLog(gomock.Any(), gomock.Any()).Return(nil)
	assert.True(t, f.commitEditLog(editLog))
}

func TestFamily_needCompact(t *testing.T) {
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
	// case 1: empty family
	v.EXPECT().NumberOfFilesInLevel(gomock.Any()).Return(0)
	assert.False(t, f.needCompat())
	// case 2: compacting
	f2 := f.(*family)
	f2.compacting.Store(true)
	assert.False(t, f.needCompat())
	f2.compacting.Store(false)
	// case 3: need compact
	v.EXPECT().NumberOfFilesInLevel(gomock.Any()).Return(10)
	assert.True(t, f.needCompat())
}

func TestFamily_compact(t *testing.T) {
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
	fv.EXPECT().GetAllActiveFiles().Return(nil).AnyTimes()
	fv.EXPECT().GetLiveRollupFiles().Return(nil).AnyTimes()
	// case 1: run compact job err
	v.EXPECT().PickL0Compaction(gomock.Any()).
		Return(version.NewCompaction(1, 0, nil, nil))
	compactJob := NewMockCompactJob(ctrl)
	f1 := f.(*family)
	f1.newCompactJobFunc = func(family Family, state *compactionState, rollup Rollup) CompactJob {
		return compactJob
	}
	compactJob.EXPECT().Run().Return(fmt.Errorf("err"))
	f.compact()
	time.Sleep(200 * time.Millisecond)
	// case 2: pick nil compaction
	v.EXPECT().PickL0Compaction(gomock.Any()).Return(nil)
	f.compact()

	time.Sleep(200 * time.Millisecond)
}

func TestFamily_compact_background(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer func() {
		newCompactJobFunc = newCompactJob
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
	fv.EXPECT().GetAllActiveFiles().Return(nil).AnyTimes()
	fv.EXPECT().GetLiveRollupFiles().Return(nil).AnyTimes()
	v.EXPECT().PickL0Compaction(gomock.Any()).
		Return(version.NewCompaction(1, 0, nil, nil)).AnyTimes()
	// case 2: compact job run err
	f2 := f.(*family)
	compactJob := NewMockCompactJob(ctrl)
	f2.newCompactJobFunc = func(family Family, state *compactionState, rollup Rollup) CompactJob {
		return compactJob
	}
	compactJob.EXPECT().Run().Return(fmt.Errorf("err"))
	err = f2.backgroundCompactionJob()
	assert.Error(t, err)
	// case 3: compact job run success
	compactJob.EXPECT().Run().Return(nil)
	err = f2.backgroundCompactionJob()
	assert.NoError(t, err)
}

func TestFamily_deleteObsoleteFiles(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer func() {
		listDirFunc = fileutil.ListDir
		removeDirFunc = fileutil.RemoveDir
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
	f1 := f.(*family)
	// case 1: list dir err
	listDirFunc = func(path string) (strings []string, err error) {
		return nil, fmt.Errorf("err")
	}
	f1.deleteObsoleteFiles()
	// case 2: delete obsolete file
	f.addPendingOutput(10)
	f.addPendingOutput(20)
	listDirFunc = func(path string) (strings []string, err error) {
		return []string{"000001.sst", "000002.sst", "000003.sst", "0000020.sst", "00010.sst", "000002.meta"}, nil
	}
	fv.EXPECT().GetAllActiveFiles().
		Return([]*version.FileMeta{version.NewFileMeta(2, 0, 0, 0)}).AnyTimes()
	fv.EXPECT().GetLiveRollupFiles().Return(map[table.FileNumber]timeutil.Interval{3: 10}).AnyTimes()
	store.EXPECT().evictFamilyFile(gomock.Any(), table.FileNumber(1))
	f1.deleteObsoleteFiles()
	// case 3: delete file err
	store.EXPECT().evictFamilyFile(gomock.Any(), table.FileNumber(1))
	removeDirFunc = func(name string) error {
		return fmt.Errorf("err")
	}
	f1.deleteObsoleteFiles()
}
