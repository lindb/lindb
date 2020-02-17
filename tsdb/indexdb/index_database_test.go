package indexdb

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/pkg/fileutil"
	"github.com/lindb/lindb/pkg/timeutil"
	"github.com/lindb/lindb/series"
	"github.com/lindb/lindb/tsdb/metadb"
)

func TestNewIndexDatabase(t *testing.T) {
	defer func() {
		_ = fileutil.RemoveDir(testPath)
	}()
	db, err := NewIndexDatabase(context.TODO(), "test", testPath, nil, nil)
	assert.NoError(t, err)
	assert.NotNil(t, db)
	// can't new duplicate
	db2, err := NewIndexDatabase(context.TODO(), "test", testPath, nil, nil)
	assert.Error(t, err)
	assert.Nil(t, db2)
	err = db.Close()
	assert.NoError(t, err)
}

func TestIndexDatabase_BuildInvertIndex(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer func() {
		_ = fileutil.RemoveDir(testPath)
		ctrl.Finish()
	}()
	db, err := NewIndexDatabase(context.TODO(), "test", testPath, nil, nil)
	assert.NoError(t, err)
	assert.NotNil(t, db)
	db1 := db.(*indexDatabase)
	index := NewMockInvertedIndex(ctrl)
	db1.index = index
	index.EXPECT().buildInvertIndex(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any())
	db.BuildInvertIndex("ns", "cpu", map[string]string{"ip": "1.1.1.1"}, 10)
}

func TestIndexDatabase_GetOrCreateSeriesID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer func() {
		ctrl.Finish()
		_ = fileutil.RemoveDir(testPath)
	}()

	db, err := NewIndexDatabase(context.TODO(), "test", testPath, nil, nil)
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
	db, err = NewIndexDatabase(context.TODO(), "test", testPath, nil, nil)
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
	// close db
	err = db.Close()
	assert.NoError(t, err)
}

func TestIndexDatabase_GetOrCreateSeriesID_err(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer func() {
		ctrl.Finish()
		_ = fileutil.RemoveDir(testPath)
		createBackend = newIDMappingBackend
	}()

	backend := NewMockIDMappingBackend(ctrl)
	createBackend = func(name, parent string) (IDMappingBackend, error) {
		return backend, nil
	}
	metadata := metadb.NewMockMetadata(ctrl)
	metadataDB := metadb.NewMockMetadataDatabase(ctrl)
	metadata.EXPECT().MetadataDatabase().Return(metadataDB).AnyTimes()
	metadataDB.EXPECT().GenTagKeyID(gomock.Any(), gomock.Any(), gomock.Any()).Return(uint32(1), nil).AnyTimes()
	db, err := NewIndexDatabase(context.TODO(), "test", testPath, metadata, nil)
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

func TestIndexDatabase_FindSeriesIDsByExpr(t *testing.T) {
	defer func() {
		_ = fileutil.RemoveDir(testPath)
	}()

	//FIXME stone1100 need impl
	db, err := NewIndexDatabase(context.TODO(), "test", testPath, nil, nil)
	assert.NoError(t, err)
	assert.Panics(t, func() {
		_, _ = db.GetSeriesIDsForTag(1)
	})
	assert.Panics(t, func() {
		_, _ = db.GetSeriesIDsByTagValueIDs(1, nil)
	})
	assert.Panics(t, func() {
		_ = db.SuggestTagValues(1, "11", 100)
	})
	assert.Panics(t, func() {
		_, _ = db.GetGroupingContext(nil, series.NewVersion())
	})
}

func TestIndexDatabase_Close(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer func() {
		ctrl.Finish()
		_ = fileutil.RemoveDir(testPath)
		createBackend = newIDMappingBackend
	}()

	backend := NewMockIDMappingBackend(ctrl)
	createBackend = func(name, parent string) (IDMappingBackend, error) {
		return backend, nil
	}
	db, err := NewIndexDatabase(context.TODO(), "test", testPath, nil, nil)
	assert.NoError(t, err)

	// case 1: save mutable event err
	db2 := db.(*indexDatabase)
	db2.rwMutex.Lock()
	db2.mutable.addSeriesID(1, 1, 1)
	db2.rwMutex.Unlock()
	backend.EXPECT().saveMapping(gomock.Any()).Return(fmt.Errorf("err"))
	err = db.Close()
	assert.Error(t, err)

	// case 2: save immutable event err
	db, err = NewIndexDatabase(context.TODO(), "test", testPath, nil, nil)
	assert.NoError(t, err)
	db2 = db.(*indexDatabase)
	db2.rwMutex.Lock()
	db2.immutable = newMappingEvent()
	db2.immutable.addSeriesID(1, 1, 1)
	db2.rwMutex.Unlock()
	backend.EXPECT().saveMapping(gomock.Any()).Return(fmt.Errorf("err"))
	err = db.Close()
	assert.Error(t, err)
}

func TestIndexDatabase_checkSync(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer func() {
		ctrl.Finish()
		_ = fileutil.RemoveDir(testPath)
		syncInterval = 2 * timeutil.OneSecond
	}()

	syncInterval = 100
	db, err := NewIndexDatabase(context.TODO(), "test", testPath, nil, nil)
	assert.NoError(t, err)
	// mock one metric event
	seriesID, isCreated, err := db.GetOrCreateSeriesID(1, 30)
	assert.NoError(t, err)
	assert.True(t, isCreated)
	assert.Equal(t, uint32(1), seriesID)
	time.Sleep(400 * time.Millisecond)

	// mock one metric event, save event err
	backend := NewMockIDMappingBackend(ctrl)
	backend.EXPECT().saveMapping(gomock.Any()).Return(fmt.Errorf("err")).AnyTimes()
	db2 := db.(*indexDatabase)
	db2.rwMutex.Lock()
	db2.backend = backend
	db2.rwMutex.Unlock()

	seriesID, isCreated, err = db.GetOrCreateSeriesID(1, 40)
	assert.NoError(t, err)
	assert.True(t, isCreated)
	assert.Equal(t, uint32(2), seriesID)
	time.Sleep(400 * time.Millisecond)
	_ = db.Close()
}

func TestMetadataDatabase_notify_timeout(t *testing.T) {
	defer func() {
		syncInterval = 2 * timeutil.OneSecond
		_ = fileutil.RemoveDir(testPath)
	}()

	syncInterval = 100
	db, err := NewIndexDatabase(context.TODO(), "test", testPath, nil, nil)
	assert.NoError(t, err)
	db1 := db.(*indexDatabase)
	// mock notify
	db1.syncSignal <- struct{}{}
	time.Sleep(400 * time.Millisecond)

	// close it mock goroutine exit, no goroutine consume event
	_ = db.Close()

	// mock goroutine consume event
	go func() {
		time.Sleep(2 * time.Second)
		<-db1.syncSignal
	}()
	// add chan item
	db1.syncSignal <- struct{}{}
	// mock mutable isn't empty
	db1.rwMutex.Lock()
	db1.mutable = newMappingEvent()
	db1.mutable.addSeriesID(1, 1, 1)
	db1.rwMutex.Unlock()
	// test notify timeout
	db1.notifySyncWithLock(true)
}
