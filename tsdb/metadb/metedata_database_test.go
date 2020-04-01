package metadb

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/constants"
	"github.com/lindb/lindb/pkg/fileutil"
	"github.com/lindb/lindb/pkg/timeutil"
	"github.com/lindb/lindb/series"
	"github.com/lindb/lindb/series/field"
	"github.com/lindb/lindb/series/tag"
)

func TestMetadataDatabase_New(t *testing.T) {
	defer func() {
		_ = fileutil.RemoveDir(testPath)
	}()

	// test: new success
	db, err := NewMetadataDatabase(context.TODO(), "test-db", testPath)
	assert.NoError(t, err)
	assert.NotNil(t, db)

	// test: can't re-open
	db1, err := NewMetadataDatabase(context.TODO(), "test-db", testPath)
	assert.Error(t, err)
	assert.Nil(t, db1)

	// close db
	err = db.Close()
	assert.NoError(t, err)
}

func TestMetadataDatabase_SuggestNamespace(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer func() {
		ctrl.Finish()
		createMetadataBackend = newMetadataBackend
	}()
	mockBackend := NewMockMetadataBackend(ctrl)
	createMetadataBackend = func(name, parent string) (backend MetadataBackend, err error) {
		return mockBackend, nil
	}
	db, err := NewMetadataDatabase(context.TODO(), "test-db", testPath)
	assert.NoError(t, err)
	mockBackend.EXPECT().suggestNamespace(gomock.Any(), gomock.Any()).Return([]string{"a"}, nil)
	values, err := db.SuggestNamespace("ns", 10)
	assert.NoError(t, err)
	assert.Equal(t, []string{"a"}, values)

	mockBackend.EXPECT().Close().Return(nil)
	_ = db.Close()
}

func TestMetadataDatabase_SuggestMetricName(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer func() {
		ctrl.Finish()
		createMetadataBackend = newMetadataBackend
	}()
	mockBackend := NewMockMetadataBackend(ctrl)
	createMetadataBackend = func(name, parent string) (backend MetadataBackend, err error) {
		return mockBackend, nil
	}
	db, err := NewMetadataDatabase(context.TODO(), "test-db", testPath)
	assert.NoError(t, err)
	mockBackend.EXPECT().suggestMetricName(gomock.Any(), gomock.Any(), gomock.Any()).Return([]string{"a"}, nil)
	values, err := db.SuggestMetrics("ns", "pp", 10)
	assert.NoError(t, err)
	assert.Equal(t, []string{"a"}, values)

	mockBackend.EXPECT().Close().Return(nil)
	_ = db.Close()
}

func TestMetadataDatabase_GetMetricID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer func() {
		ctrl.Finish()
		createMetadataBackend = newMetadataBackend
	}()
	mockBackend := NewMockMetadataBackend(ctrl)
	createMetadataBackend = func(name, parent string) (backend MetadataBackend, err error) {
		return mockBackend, nil
	}
	db, err := NewMetadataDatabase(context.TODO(), "test-db", testPath)
	assert.NoError(t, err)
	gomock.InOrder(
		mockBackend.EXPECT().loadMetricMetadata("ns-1", "name1").Return(nil, constants.ErrNotFound),
		mockBackend.EXPECT().genMetricID().Return(uint32(1)),
	)
	metricID, err := db.GenMetricID("ns-1", "name1")
	assert.NoError(t, err)
	assert.Equal(t, uint32(1), metricID)

	metricID, err = db.GetMetricID("ns-1", "name1")
	assert.NoError(t, err)
	assert.Equal(t, uint32(1), metricID)

	mockBackend.EXPECT().getMetricID("ns-1", "name2").Return(uint32(10), nil)
	metricID, err = db.GetMetricID("ns-1", "name2")
	assert.NoError(t, err)
	assert.Equal(t, uint32(10), metricID)

	mockBackend.EXPECT().saveMetadata(gomock.Any()).AnyTimes()
	mockBackend.EXPECT().Close().Return(nil)
	_ = db.Close()
}

func TestMetadataDatabase_GetTagKey(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer func() {
		ctrl.Finish()
		createMetadataBackend = newMetadataBackend
	}()
	mockBackend := NewMockMetadataBackend(ctrl)
	createMetadataBackend = func(name, parent string) (backend MetadataBackend, err error) {
		return mockBackend, nil
	}
	db, err := NewMetadataDatabase(context.TODO(), "test-db", testPath)
	assert.NoError(t, err)
	meta := NewMockMetricMetadata(ctrl)
	mockBackend.EXPECT().loadMetricMetadata("ns-1", "name1").Return(meta, nil)
	meta.EXPECT().getMetricID().Return(uint32(1))
	metricID, err := db.GenMetricID("ns-1", "name1")
	assert.NoError(t, err)
	assert.Equal(t, uint32(1), metricID)

	// case 1: from memory
	meta.EXPECT().getTagKeyID("tag-key").Return(uint32(10), true)
	tagKeyID, err := db.GetTagKeyID("ns-1", "name1", "tag-key")
	assert.NoError(t, err)
	assert.Equal(t, uint32(10), tagKeyID)

	// case 2: from memory not exist
	meta.EXPECT().getTagKeyID("tag-key").Return(uint32(10), false)
	tagKeyID, err = db.GetTagKeyID("ns-1", "name1", "tag-key")
	assert.Equal(t, constants.ErrNotFound, err)
	assert.Equal(t, uint32(0), tagKeyID)

	// case 4: backend, metric not exist
	mockBackend.EXPECT().getMetricID("ns-1", "name2").Return(uint32(0), constants.ErrNotFound)
	tagKeyID, err = db.GetTagKeyID("ns-1", "name2", "tag-key")
	assert.Equal(t, constants.ErrNotFound, err)
	assert.Equal(t, uint32(0), tagKeyID)

	// case 4: backend exist
	mockBackend.EXPECT().getMetricID("ns-1", "name2").Return(uint32(10), nil)
	mockBackend.EXPECT().getTagKeyID(uint32(10), "tag-key").Return(uint32(20), nil)
	tagKeyID, err = db.GetTagKeyID("ns-1", "name2", "tag-key")
	assert.NoError(t, err)
	assert.Equal(t, uint32(20), tagKeyID)

	// case 5: all tag keys from memory
	meta.EXPECT().getAllTagKeys().Return([]tag.Meta{{ID: 10, Key: "tag-key"}})
	tagKeys, err := db.GetAllTagKeys("ns-1", "name1")
	assert.NoError(t, err)
	assert.Equal(t, []tag.Meta{{ID: 10, Key: "tag-key"}}, tagKeys)

	// case 6: backend, metric not exist
	mockBackend.EXPECT().getMetricID("ns-1", "name2").Return(uint32(0), constants.ErrNotFound)
	tagKeys, err = db.GetAllTagKeys("ns-1", "name2")
	assert.Equal(t, constants.ErrNotFound, err)
	assert.Nil(t, tagKeys)

	// case 7: backend, tag keys exist
	mockBackend.EXPECT().getMetricID("ns-1", "name2").Return(uint32(10), nil)
	mockBackend.EXPECT().getAllTagKeys(uint32(10)).Return([]tag.Meta{{ID: 10, Key: "tag-key"}}, nil)
	tagKeys, err = db.GetAllTagKeys("ns-1", "name2")
	assert.NoError(t, err)
	assert.Equal(t, []tag.Meta{{ID: 10, Key: "tag-key"}}, tagKeys)

	mockBackend.EXPECT().saveMetadata(gomock.Any()).AnyTimes()
	mockBackend.EXPECT().Close().Return(nil)
	_ = db.Close()
}

func TestMetadataDatabase_SuggestTagKeys(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer func() {
		ctrl.Finish()
		createMetadataBackend = newMetadataBackend
	}()
	mockBackend := NewMockMetadataBackend(ctrl)
	createMetadataBackend = func(name, parent string) (backend MetadataBackend, err error) {
		return mockBackend, nil
	}
	db, err := NewMetadataDatabase(context.TODO(), "test-db", testPath)
	assert.NoError(t, err)
	mockBackend.EXPECT().getMetricID(gomock.Any(), gomock.Any()).Return(uint32(10), nil).AnyTimes()
	// case 1: suggest tag keys
	mockBackend.EXPECT().getAllTagKeys(gomock.Any()).Return([]tag.Meta{{ID: 10, Key: "tag-key"}}, nil)
	tagKeys, err := db.SuggestTagKeys("ns-1", "name1", "tag", 100)
	assert.NoError(t, err)
	assert.Equal(t, []string{"tag-key"}, tagKeys)
	// case 2: get tag values err
	mockBackend.EXPECT().getAllTagKeys(gomock.Any()).Return([]tag.Meta{{ID: 10, Key: "tag-key"}}, fmt.Errorf("err"))
	tagKeys, err = db.SuggestTagKeys("ns-1", "name1", "tag", 100)
	assert.Error(t, err)
	assert.Nil(t, tagKeys)
	// case 2: get tag values limit
	mockBackend.EXPECT().getAllTagKeys(gomock.Any()).Return([]tag.Meta{{ID: 10, Key: "tag-key"}, {ID: 10, Key: "tag-key1"}}, nil)
	tagKeys, err = db.SuggestTagKeys("ns-1", "name1", "tag", 1)
	assert.NoError(t, err)
	assert.Equal(t, []string{"tag-key"}, tagKeys)
}

func TestMetadataDatabase_GetField(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer func() {
		ctrl.Finish()
		createMetadataBackend = newMetadataBackend
	}()
	mockBackend := NewMockMetadataBackend(ctrl)
	createMetadataBackend = func(name, parent string) (backend MetadataBackend, err error) {
		return mockBackend, nil
	}
	db, err := NewMetadataDatabase(context.TODO(), "test-db", testPath)
	assert.NoError(t, err)
	meta := NewMockMetricMetadata(ctrl)
	mockBackend.EXPECT().loadMetricMetadata("ns-1", "name1").Return(meta, nil)
	meta.EXPECT().getMetricID().Return(uint32(1))
	metricID, err := db.GenMetricID("ns-1", "name1")
	assert.NoError(t, err)
	assert.Equal(t, uint32(1), metricID)

	// case 1: from memory
	meta.EXPECT().getField("f1").Return(field.Meta{ID: 19, Type: field.SumField}, true)
	f, err := db.GetField("ns-1", "name1", "f1")
	assert.NoError(t, err)
	assert.Equal(t, field.Meta{ID: 19, Type: field.SumField}, f)

	// case 2: from memory not exist
	meta.EXPECT().getField("f1").Return(field.Meta{}, false)
	f, err = db.GetField("ns-1", "name1", "f1")
	assert.Equal(t, constants.ErrNotFound, err)
	assert.Equal(t, field.Meta{}, f)

	// case 4: backend, metric not exist
	mockBackend.EXPECT().getMetricID("ns-1", "name2").Return(uint32(0), constants.ErrNotFound)
	f, err = db.GetField("ns-1", "name2", "f1")
	assert.Equal(t, constants.ErrNotFound, err)
	assert.Equal(t, field.Meta{}, f)

	// case 4: backend exist
	mockBackend.EXPECT().getMetricID("ns-1", "name2").Return(uint32(10), nil)
	mockBackend.EXPECT().getField(uint32(10), "f1").Return(field.Meta{ID: 19, Type: field.SumField}, nil)
	f, err = db.GetField("ns-1", "name2", "f1")
	assert.NoError(t, err)
	assert.Equal(t, field.Meta{ID: 19, Type: field.SumField}, f)

	// case 5: all tag keys from memory
	meta.EXPECT().getAllFields().Return([]field.Meta{{ID: 19, Type: field.SumField}})
	fields, err := db.GetAllFields("ns-1", "name1")
	assert.NoError(t, err)
	assert.Equal(t, []field.Meta{{ID: 19, Type: field.SumField}}, fields)

	// case 6: backend, metric not exist
	mockBackend.EXPECT().getMetricID("ns-1", "name2").Return(uint32(0), constants.ErrNotFound)
	fields, err = db.GetAllFields("ns-1", "name2")
	assert.Equal(t, constants.ErrNotFound, err)
	assert.Nil(t, fields)

	// case 7: backend, tag keys exist
	mockBackend.EXPECT().getMetricID("ns-1", "name2").Return(uint32(10), nil)
	mockBackend.EXPECT().getAllFields(uint32(10)).Return([]field.Meta{{ID: 19, Type: field.SumField}}, nil)
	fields, err = db.GetAllFields("ns-1", "name2")
	assert.NoError(t, err)
	assert.Equal(t, []field.Meta{{ID: 19, Type: field.SumField}}, fields)

	mockBackend.EXPECT().saveMetadata(gomock.Any()).AnyTimes()
	mockBackend.EXPECT().Close().Return(nil)
	_ = db.Close()
}

func TestMetadataDatabase_GenMetricID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer func() {
		ctrl.Finish()
		createMetadataBackend = newMetadataBackend
	}()
	mockBackend := NewMockMetadataBackend(ctrl)
	createMetadataBackend = func(name, parent string) (backend MetadataBackend, err error) {
		return mockBackend, nil
	}
	db, err := NewMetadataDatabase(context.TODO(), "test-db", testPath)
	assert.NoError(t, err)
	gomock.InOrder(
		mockBackend.EXPECT().loadMetricMetadata("ns-1", "name1").Return(nil, constants.ErrNotFound),
		mockBackend.EXPECT().genMetricID().Return(uint32(1)),
	)
	// case 1: gen new metric id
	metricID, err := db.GenMetricID("ns-1", "name1")
	assert.NoError(t, err)
	assert.Equal(t, uint32(1), metricID)

	// cast 2: get metric id from memory
	metricID, err = db.GenMetricID("ns-1", "name1")
	assert.NoError(t, err)
	assert.Equal(t, uint32(1), metricID)

	// case 3: load metric meta err
	mockBackend.EXPECT().loadMetricMetadata("ns-1", "name2").Return(nil, fmt.Errorf("err"))
	metricID, err = db.GenMetricID("ns-1", "name2")
	assert.Error(t, err)
	assert.Equal(t, uint32(0), metricID)

	// case 4: load metric meta ok
	meta := NewMockMetricMetadata(ctrl)
	meta.EXPECT().getMetricID().Return(uint32(100))
	mockBackend.EXPECT().loadMetricMetadata("ns-1", "name2").Return(meta, nil)
	metricID, err = db.GenMetricID("ns-1", "name2")
	assert.NoError(t, err)
	assert.Equal(t, uint32(100), metricID)

	mockBackend.EXPECT().saveMetadata(gomock.Any()).AnyTimes()
	mockBackend.EXPECT().Close().Return(nil)
	_ = db.Close()
}

func TestMetadataDatabase_GenFieldID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer func() {
		ctrl.Finish()
		createMetadataBackend = newMetadataBackend
	}()
	mockBackend := NewMockMetadataBackend(ctrl)
	createMetadataBackend = func(name, parent string) (backend MetadataBackend, err error) {
		return mockBackend, nil
	}
	meta := NewMockMetricMetadata(ctrl)
	db, err := NewMetadataDatabase(context.TODO(), "test-db", testPath)
	assert.NoError(t, err)
	gomock.InOrder(
		mockBackend.EXPECT().loadMetricMetadata("ns-1", "name1").Return(meta, nil),
		meta.EXPECT().getMetricID().Return(uint32(100)),
		meta.EXPECT().getField("f").Return(field.Meta{}, false),
		meta.EXPECT().createField(gomock.Any(), gomock.Any()).Return(field.ID(10), nil),
		meta.EXPECT().getMetricID().Return(uint32(1)),
	)
	// case 1: gen new field id
	_, err = db.GenMetricID("ns-1", "name1")
	assert.NoError(t, err)
	fieldID, err := db.GenFieldID("ns-1", "name1", "f", field.SumField)
	assert.NoError(t, err)
	assert.Equal(t, field.ID(10), fieldID)

	// case 2: get field id from memory
	meta.EXPECT().getField("f").Return(field.Meta{ID: 10, Type: field.SumField}, true)
	fieldID, err = db.GenFieldID("ns-1", "name1", "f", field.SumField)
	assert.NoError(t, err)
	assert.Equal(t, field.ID(10), fieldID)

	// case 3: get field id from memory, but type not match
	meta.EXPECT().getField("f").Return(field.Meta{ID: 10, Type: field.MinField}, true)
	fieldID, err = db.GenFieldID("ns-1", "name1", "f", field.SumField)
	assert.Equal(t, series.ErrWrongFieldType, err)
	assert.Equal(t, field.ID(0), fieldID)

	// case 4: create fail
	gomock.InOrder(
		meta.EXPECT().getField("f").Return(field.Meta{}, false),
		meta.EXPECT().createField(gomock.Any(), gomock.Any()).Return(field.ID(10), fmt.Errorf("err")),
	)
	fieldID, err = db.GenFieldID("ns-1", "name1", "f", field.SumField)
	assert.Error(t, err)
	assert.Equal(t, field.ID(0), fieldID)

	mockBackend.EXPECT().saveMetadata(gomock.Any()).AnyTimes()
	mockBackend.EXPECT().Close().Return(nil)
	_ = db.Close()
}

func TestMetadataDatabase_GenTagKeyID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer func() {
		ctrl.Finish()
		createMetadataBackend = newMetadataBackend
	}()
	mockBackend := NewMockMetadataBackend(ctrl)
	createMetadataBackend = func(name, parent string) (backend MetadataBackend, err error) {
		return mockBackend, nil
	}
	meta := NewMockMetricMetadata(ctrl)
	db, err := NewMetadataDatabase(context.TODO(), "test-db", testPath)
	assert.NoError(t, err)
	gomock.InOrder(
		mockBackend.EXPECT().loadMetricMetadata("ns-1", "name1").Return(meta, nil),
		meta.EXPECT().getMetricID().Return(uint32(100)),
		meta.EXPECT().getTagKeyID("tag-key").Return(uint32(0), false),
		meta.EXPECT().checkTagKeyCount().Return(nil),
		mockBackend.EXPECT().genTagKeyID().Return(uint32(10)),
		meta.EXPECT().createTagKey("tag-key", uint32(10)),
		meta.EXPECT().getMetricID().Return(uint32(1)),
	)
	// case 1: gen new tag key id
	_, err = db.GenMetricID("ns-1", "name1")
	assert.NoError(t, err)
	tagKeyID, err := db.GenTagKeyID("ns-1", "name1", "tag-key")
	assert.NoError(t, err)
	assert.Equal(t, uint32(10), tagKeyID)

	// case 2: get tag key id from memory
	meta.EXPECT().getTagKeyID("tag-key").Return(uint32(10), true)
	tagKeyID, err = db.GenTagKeyID("ns-1", "name1", "tag-key")
	assert.NoError(t, err)
	assert.Equal(t, uint32(10), tagKeyID)

	// case 3: too many tags
	gomock.InOrder(
		meta.EXPECT().getTagKeyID("tag-key").Return(uint32(0), false),
		meta.EXPECT().checkTagKeyCount().Return(series.ErrTooManyTagKeys),
	)
	tagKeyID, err = db.GenTagKeyID("ns-1", "name1", "tag-key")
	assert.Equal(t, series.ErrTooManyTagKeys, err)
	assert.Equal(t, uint32(0), tagKeyID)

	mockBackend.EXPECT().saveMetadata(gomock.Any()).AnyTimes()
	mockBackend.EXPECT().Close().Return(nil)
	_ = db.Close()
}

func TestMetadataDatabase_Close(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer func() {
		ctrl.Finish()
		createMetadataBackend = newMetadataBackend
	}()
	mockBackend := NewMockMetadataBackend(ctrl)
	createMetadataBackend = func(name, parent string) (backend MetadataBackend, err error) {
		return mockBackend, nil
	}
	db, err := NewMetadataDatabase(context.TODO(), "test-db", testPath)
	assert.NoError(t, err)
	db1 := db.(*metadataDatabase)
	db1.rwMux.Lock()
	db1.immutable = newMetadataUpdateEvent()
	db1.mutable.addMetric("ns1", "name", 1)
	db1.immutable.addMetric("ns1", "name22", 1)
	db1.rwMux.Unlock()

	mockBackend.EXPECT().saveMetadata(gomock.Any()).Return(fmt.Errorf("err"))
	err = db.Close()
	assert.Error(t, err)

	mockBackend.EXPECT().saveMetadata(gomock.Any()).Return(nil)
	mockBackend.EXPECT().saveMetadata(gomock.Any()).Return(fmt.Errorf("err"))
	err = db.Close()
	assert.Error(t, err)
}

func TestMetadataDatabase_checkSync(t *testing.T) {
	defer func() {
		_ = fileutil.RemoveDir(testPath)
		syncInterval = 2 * timeutil.OneSecond
	}()

	syncInterval = 100
	db, err := NewMetadataDatabase(context.TODO(), "test-db", testPath)
	assert.NoError(t, err)
	// mock one metric event
	metricID, err := db.GenMetricID("ns-1", "name1")
	assert.NoError(t, err)
	assert.Equal(t, uint32(1), metricID)
	time.Sleep(400 * time.Millisecond)

	metricID, err = db.GenMetricID("ns-1", "name2")
	assert.NoError(t, err)
	assert.Equal(t, uint32(2), metricID)
	_, err = db.GenTagKeyID("ns-1", "name1", "") // tag key cannot be empty
	assert.NoError(t, err)
	time.Sleep(400 * time.Millisecond)
	_ = db.Close()

	// reopen
	db, err = NewMetadataDatabase(context.TODO(), "test-db", testPath)
	assert.NoError(t, err)
	metricID, err = db.GenMetricID("ns-2", "name2")
	assert.NoError(t, err)
	assert.Equal(t, uint32(3), metricID)
	time.Sleep(400 * time.Millisecond)
	_ = db.Close()
}

func TestMetadataDatabase_notify_timeout(t *testing.T) {
	defer func() {
		syncInterval = 2 * timeutil.OneSecond
		_ = fileutil.RemoveDir(testPath)
	}()

	syncInterval = 100
	db, err := NewMetadataDatabase(context.TODO(), "test-db", testPath)
	assert.NoError(t, err)
	db1 := db.(*metadataDatabase)
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
	db1.rwMux.Lock()
	db1.mutable = newMetadataUpdateEvent()
	db1.mutable.addMetric("ns-1", "name-4", uint32(100))
	db1.rwMux.Unlock()
	// test notify timeout
	db1.notifySyncWithLock(true)
}

func TestMetadataDatabase_Sync(t *testing.T) {
	defer func() {
		_ = fileutil.RemoveDir(testPath)
	}()
	db := newMockMetadataDatabase(t)
	err := db.Sync()
	assert.NoError(t, err)
	err = db.Close()
	assert.NoError(t, err)

	_ = db.Close()
}

func newMockMetadataDatabase(t *testing.T) MetadataDatabase {
	db, err := NewMetadataDatabase(context.TODO(), "test-db", testPath)
	assert.NoError(t, err)
	assert.NotNil(t, db)

	return db
}
