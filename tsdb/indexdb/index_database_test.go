package indexdb

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/lindb/roaring"
	"github.com/stretchr/testify/assert"
	"go.uber.org/atomic"

	"github.com/lindb/lindb/pkg/fileutil"
	"github.com/lindb/lindb/pkg/timeutil"
	"github.com/lindb/lindb/series/tag"
	"github.com/lindb/lindb/tsdb/metadb"
	"github.com/lindb/lindb/tsdb/wal"
)

func TestNewIndexDatabase(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer func() {
		_ = fileutil.RemoveDir(testPath)

		ctrl.Finish()
	}()
	mockMetadata := metadb.NewMockMetadata(ctrl)
	mockMetadata.EXPECT().DatabaseName().Return("test").AnyTimes()
	db, err := NewIndexDatabase(context.TODO(), testPath, mockMetadata, nil, nil)
	assert.NoError(t, err)
	assert.NotNil(t, db)
	// can't new duplicate
	db2, err := NewIndexDatabase(context.TODO(), testPath, nil, nil, nil)
	assert.Error(t, err)
	assert.Nil(t, db2)

	err = db.Close()
	assert.NoError(t, err)
}

func TestNewIndexDatabase_err(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer func() {
		createBackend = newIDMappingBackend
		createSeriesWAL = wal.NewSeriesWAL
		_ = fileutil.RemoveDir(testPath)

		ctrl.Finish()
	}()
	mockMetadata := metadb.NewMockMetadata(ctrl)
	mockMetadata.EXPECT().DatabaseName().Return("test").AnyTimes()

	backend := NewMockIDMappingBackend(ctrl)
	createBackend = func(parent string) (IDMappingBackend, error) {
		return backend, nil
	}

	// case 1: create series wal err
	backend.EXPECT().Close().Return(fmt.Errorf("err"))
	createSeriesWAL = func(path string) (wal.SeriesWAL, error) {
		return nil, fmt.Errorf("err")
	}

	db, err := NewIndexDatabase(context.TODO(), testPath, mockMetadata, nil, nil)
	assert.Error(t, err)
	assert.Nil(t, db)
	// case 2: series wal recovery err
	mockSeriesWAl := wal.NewMockSeriesWAL(ctrl)
	createSeriesWAL = func(path string) (wal.SeriesWAL, error) {
		return mockSeriesWAl, nil
	}
	backend.EXPECT().Close().Return(fmt.Errorf("err"))
	mockSeriesWAl.EXPECT().Recovery(gomock.Any(), gomock.Any())
	mockSeriesWAl.EXPECT().NeedRecovery().Return(true)
	db, err = NewIndexDatabase(context.TODO(), testPath, mockMetadata, nil, nil)
	assert.Error(t, err)
	assert.Nil(t, db)
}

func TestIndexDatabase_SuggestTagValues(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer func() {
		_ = fileutil.RemoveDir(testPath)

		ctrl.Finish()
	}()
	metaDB := metadb.NewMockMetadata(ctrl)
	metaDB.EXPECT().DatabaseName().Return("test")
	tagMeta := metadb.NewMockTagMetadata(ctrl)
	metaDB.EXPECT().TagMetadata().Return(tagMeta)
	db, err := NewIndexDatabase(context.TODO(), testPath, metaDB, nil, nil)
	assert.NoError(t, err)
	tagMeta.EXPECT().SuggestTagValues(gomock.Any(), gomock.Any(), gomock.Any()).Return([]string{"a", "b"})
	tagValues := db.SuggestTagValues(10, "test", 100)
	assert.Equal(t, []string{"a", "b"}, tagValues)

	err = db.Close()
	assert.NoError(t, err)
}

func TestIndexDatabase_BuildInvertIndex(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer func() {
		_ = fileutil.RemoveDir(testPath)
		ctrl.Finish()
	}()
	meta := metadb.NewMockMetadata(ctrl)
	meta.EXPECT().DatabaseName().Return("test").AnyTimes()
	db, err := NewIndexDatabase(context.TODO(), testPath, meta, nil, nil)
	assert.NoError(t, err)
	assert.NotNil(t, db)
	db1 := db.(*indexDatabase)
	index := NewMockInvertedIndex(ctrl)
	db1.index = index
	index.EXPECT().buildInvertIndex(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any())
	db.BuildInvertIndex("ns", "cpu", map[string]string{"ip": "1.1.1.1"}, 10)

	index.EXPECT().Flush().Return(nil)
	err = db.Close()
	assert.NoError(t, err)
}

func TestIndexDatabase_series_Recovery_err(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer func() {
		createBackend = newIDMappingBackend
		_ = fileutil.RemoveDir(testPath)

		ctrl.Finish()
	}()

	meta := metadb.NewMockMetadata(ctrl)
	meta.EXPECT().DatabaseName().Return("test").AnyTimes()
	db, err := NewIndexDatabase(context.TODO(), testPath, meta, nil, nil)
	assert.NoError(t, err)
	assert.NotNil(t, db)
	for i := 0; i < 11000; i++ {
		_, isCreated, err := db.GetOrCreateSeriesID(1, uint64(i))
		assert.NoError(t, err)
		assert.True(t, isCreated)
	}
	err = db.Close()
	assert.NoError(t, err)

	backend := NewMockIDMappingBackend(ctrl)
	backend.EXPECT().Close().Return(nil).AnyTimes()
	createBackend = func(parent string) (IDMappingBackend, error) {
		return backend, nil
	}
	backend.EXPECT().saveMapping(gomock.Any()).Return(fmt.Errorf("err"))
	db, err = NewIndexDatabase(context.TODO(), testPath, meta, nil, nil)
	assert.Error(t, err)
	assert.Nil(t, db)

	createBackend = newIDMappingBackend
	// recovery success
	db, err = NewIndexDatabase(context.TODO(), testPath, meta, nil, nil)
	assert.NoError(t, err)
	assert.NotNil(t, db)

	for i := 0; i < 100; i++ {
		_, isCreated, err := db.GetOrCreateSeriesID(1, uint64(1000000+i))
		assert.NoError(t, err)
		assert.True(t, isCreated)
	}
	err = db.Close()
	assert.NoError(t, err)

	createBackend = func(parent string) (IDMappingBackend, error) {
		return backend, nil
	}
	backend.EXPECT().saveMapping(gomock.Any()).Return(fmt.Errorf("err"))
	db, err = NewIndexDatabase(context.TODO(), testPath, meta, nil, nil)
	assert.Error(t, err)
	assert.Nil(t, db)
}

func TestIndexDatabase_GetOrCreateSeriesID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer func() {
		_ = fileutil.RemoveDir(testPath)

		ctrl.Finish()
	}()

	meta := metadb.NewMockMetadata(ctrl)
	meta.EXPECT().DatabaseName().Return("test").AnyTimes()
	db, err := NewIndexDatabase(context.TODO(), testPath, meta, nil, nil)
	assert.NoError(t, err)
	// case 1: generate new series id and create new metric id mapping
	seriesID, isCreated, err := db.GetOrCreateSeriesID(1, 10)
	assert.NoError(t, err)
	assert.True(t, isCreated)
	assert.Equal(t, uint32(1), seriesID)
	// case 2: get series id from memory
	seriesID, isCreated, err = db.GetOrCreateSeriesID(1, 10)
	assert.NoError(t, err)
	assert.False(t, isCreated)
	assert.Equal(t, uint32(1), seriesID)
	// case 3: generate new series id from memory
	seriesID, isCreated, err = db.GetOrCreateSeriesID(1, 20)
	assert.NoError(t, err)
	assert.True(t, isCreated)
	assert.Equal(t, uint32(2), seriesID)
	// close db
	err = db.Close()
	assert.NoError(t, err)

	// reopen
	db, err = NewIndexDatabase(context.TODO(), testPath, meta, nil, nil)
	assert.NoError(t, err)
	// case 4: get series id from backend
	seriesID, isCreated, err = db.GetOrCreateSeriesID(1, 20)
	assert.NoError(t, err)
	assert.False(t, isCreated)
	assert.Equal(t, uint32(2), seriesID)
	// case 5: gen series id, id sequence reset from backend
	seriesID, isCreated, err = db.GetOrCreateSeriesID(1, 30)
	assert.NoError(t, err)
	assert.True(t, isCreated)
	assert.Equal(t, uint32(3), seriesID)
	// case 6: append series wal err, need rollback new series id
	mockSeriesWAl := wal.NewMockSeriesWAL(ctrl)
	db1 := db.(*indexDatabase)
	oldWAL := db1.seriesWAL
	db1.seriesWAL = mockSeriesWAl
	mockSeriesWAl.EXPECT().Append(uint32(1), uint64(50), uint32(4)).Return(fmt.Errorf("err"))
	seriesID, isCreated, err = db.GetOrCreateSeriesID(1, 50)
	assert.Error(t, err)
	assert.False(t, isCreated)
	assert.Equal(t, uint32(0), seriesID)
	// add use series id => 4
	db1.seriesWAL = oldWAL
	seriesID, isCreated, err = db.GetOrCreateSeriesID(1, 50)
	assert.NoError(t, err)
	assert.True(t, isCreated)
	assert.Equal(t, uint32(4), seriesID)

	// close db
	err = db.Close()
	assert.NoError(t, err)
}

func TestIndexDatabase_GetOrCreateSeriesID_err(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer func() {
		_ = fileutil.RemoveDir(testPath)
		createBackend = newIDMappingBackend

		ctrl.Finish()
	}()

	backend := NewMockIDMappingBackend(ctrl)
	createBackend = func(parent string) (IDMappingBackend, error) {
		return backend, nil
	}
	metadata := metadb.NewMockMetadata(ctrl)
	metadata.EXPECT().DatabaseName().Return("test").AnyTimes()
	metadataDB := metadb.NewMockMetadataDatabase(ctrl)
	metadata.EXPECT().MetadataDatabase().Return(metadataDB).AnyTimes()
	metadataDB.EXPECT().GenTagKeyID(gomock.Any(), gomock.Any(), gomock.Any()).Return(uint32(1), nil).AnyTimes()
	db, err := NewIndexDatabase(context.TODO(), testPath, metadata, nil, nil)
	assert.NoError(t, err)
	// case 1: load metric mapping err
	backend.EXPECT().loadMetricIDMapping(uint32(1)).Return(nil, fmt.Errorf("err"))
	seriesID, isCreated, err := db.GetOrCreateSeriesID(1, 30)
	assert.Error(t, err)
	assert.False(t, isCreated)
	assert.Equal(t, uint32(0), seriesID)

	// case 2: load series err
	backend.EXPECT().loadMetricIDMapping(uint32(1)).Return(newMetricIDMapping(1, 0), nil)
	backend.EXPECT().getSeriesID(uint32(1), uint64(30)).Return(uint32(0), fmt.Errorf("err"))
	seriesID, isCreated, err = db.GetOrCreateSeriesID(1, 30)
	assert.Error(t, err)
	assert.False(t, isCreated)
	assert.Equal(t, uint32(0), seriesID)

	backend.EXPECT().Close().Return(nil)
	err = db.Close()
	assert.NoError(t, err)
}

func TestIndexDatabase_GetGroupingContext(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer func() {
		_ = fileutil.RemoveDir(testPath)
		ctrl.Finish()
	}()
	meta := metadb.NewMockMetadata(ctrl)
	meta.EXPECT().DatabaseName().Return("test").AnyTimes()
	db, err := NewIndexDatabase(context.TODO(), testPath, meta, nil, nil)
	assert.NoError(t, err)
	assert.NotNil(t, db)
	index := NewMockInvertedIndex(ctrl)
	db1 := db.(*indexDatabase)
	db1.index = index
	index.EXPECT().GetGroupingContext(gomock.Any(), gomock.Any()).Return(nil, nil)
	ctx, err := db.GetGroupingContext([]uint32{1, 2}, roaring.BitmapOf(1, 2, 3))
	assert.NoError(t, err)
	assert.Nil(t, ctx)

	index.EXPECT().Flush().Return(nil)
	err = db.Close()
	assert.NoError(t, err)
}

func TestIndexDatabase_GetSeriesIDs(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer func() {
		_ = fileutil.RemoveDir(testPath)
		ctrl.Finish()
	}()

	index := NewMockInvertedIndex(ctrl)
	metaDB := metadb.NewMockMetadataDatabase(ctrl)
	meta := metadb.NewMockMetadata(ctrl)
	meta.EXPECT().DatabaseName().Return("test")
	meta.EXPECT().MetadataDatabase().Return(metaDB).AnyTimes()
	db, err := NewIndexDatabase(context.TODO(), testPath, meta, nil, nil)
	db2 := db.(*indexDatabase)
	db2.index = index
	db2.metadata = meta
	assert.NoError(t, err)

	// case 1: get series ids by tag key
	index.EXPECT().GetSeriesIDsForTag(uint32(1)).Return(roaring.BitmapOf(1, 2), nil)
	seriesIDs, err := db.GetSeriesIDsForTag(1)
	assert.NoError(t, err)
	assert.NotNil(t, seriesIDs)
	// case 2: get series ids by tag value ids
	index.EXPECT().GetSeriesIDsByTagValueIDs(uint32(1), roaring.BitmapOf(1, 2, 3)).Return(roaring.BitmapOf(1, 2), nil)
	seriesIDs, err = db.GetSeriesIDsByTagValueIDs(1, roaring.BitmapOf(1, 2, 3))
	assert.NoError(t, err)
	assert.NotNil(t, seriesIDs)
	// case 3: get tags err
	metaDB.EXPECT().GetAllTagKeys(gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("err"))
	seriesIDs, err = db.GetSeriesIDsForMetric("ns", "name")
	assert.Equal(t, fmt.Errorf("err"), err)
	assert.Nil(t, seriesIDs)
	// case 4: get empty tags
	metaDB.EXPECT().GetAllTagKeys(gomock.Any(), gomock.Any()).Return(nil, nil)
	seriesIDs, err = db.GetSeriesIDsForMetric("ns", "name")
	assert.NoError(t, err)
	assert.Equal(t, roaring.BitmapOf(0), seriesIDs)
	// case 5: get series ids for metric
	metaDB.EXPECT().GetAllTagKeys(gomock.Any(), gomock.Any()).Return([]tag.Meta{{ID: 1}}, nil)
	index.EXPECT().GetSeriesIDsForTags([]uint32{1}).Return(roaring.BitmapOf(1, 2, 3), nil)
	seriesIDs, err = db.GetSeriesIDsForMetric("ns", "name")
	assert.NoError(t, err)
	assert.NotNil(t, seriesIDs)

	index.EXPECT().Flush().Return(nil)
	err = db.Close()
	assert.NoError(t, err)
}

func TestIndexDatabase_Close(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer func() {
		_ = fileutil.RemoveDir(testPath)
		createBackend = newIDMappingBackend
		createSeriesWAL = wal.NewSeriesWAL
		ctrl.Finish()
	}()

	backend := NewMockIDMappingBackend(ctrl)
	createBackend = func(parent string) (IDMappingBackend, error) {
		return backend, nil
	}
	mockSeriesWAL := wal.NewMockSeriesWAL(ctrl)
	mockSeriesWAL.EXPECT().Close().Return(fmt.Errorf("err"))

	meta := metadb.NewMockMetadata(ctrl)
	meta.EXPECT().DatabaseName().Return("test")
	db, err := NewIndexDatabase(context.TODO(), testPath, meta, nil, nil)
	db1 := db.(*indexDatabase)
	db1.seriesWAL = mockSeriesWAL

	assert.NoError(t, err)
	backend.EXPECT().Close().Return(fmt.Errorf("err"))
	err = db.Close()
	assert.Error(t, err)
}

func TestIndexDatabase_Flush(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer func() {
		createSeriesWAL = wal.NewSeriesWAL
		_ = fileutil.RemoveDir(testPath)

		ctrl.Finish()
	}()
	mockSeriesWAL := wal.NewMockSeriesWAL(ctrl)
	mockSeriesWAL.EXPECT().Close().Return(nil)
	mockSeriesWAL.EXPECT().Recovery(gomock.Any(), gomock.Any())
	mockSeriesWAL.EXPECT().NeedRecovery().Return(false).AnyTimes()
	createSeriesWAL = func(path string) (wal.SeriesWAL, error) {
		return mockSeriesWAL, nil
	}

	meta := metadb.NewMockMetadata(ctrl)
	meta.EXPECT().DatabaseName().Return("test")
	db, err := NewIndexDatabase(context.TODO(), testPath, meta, nil, nil)
	assert.NoError(t, err)
	mockSeriesWAL.EXPECT().Sync().Return(fmt.Errorf("err"))
	err = db.Flush()
	assert.NoError(t, err)
	err = db.Close()
	assert.NoError(t, err)
}

func TestIndexDatabase_checkSync(t *testing.T) {
	syncInterval = 100
	ctrl := gomock.NewController(t)
	defer func() {
		syncInterval = 2 * timeutil.OneSecond
		_ = fileutil.RemoveDir(testPath)
		createSeriesWAL = wal.NewSeriesWAL

		ctrl.Finish()
	}()

	var count atomic.Int32
	mockSeriesWAL := wal.NewMockSeriesWAL(ctrl)
	mockSeriesWAL.EXPECT().NeedRecovery().DoAndReturn(func() bool {
		count.Inc()
		return count.Load() != 1
	}).AnyTimes()
	mockSeriesWAL.EXPECT().Recovery(gomock.Any(), gomock.Any()).AnyTimes()
	createSeriesWAL = func(path string) (wal.SeriesWAL, error) {
		return mockSeriesWAL, nil
	}

	meta := metadb.NewMockMetadata(ctrl)
	meta.EXPECT().DatabaseName().Return("test").AnyTimes()
	db, err := NewIndexDatabase(context.TODO(), testPath, meta, nil, nil)
	assert.NoError(t, err)
	assert.NotNil(t, db)

	time.Sleep(time.Second)

	mockSeriesWAL.EXPECT().Close().Return(nil)
	err = db.Close()
	assert.NoError(t, err)
}
