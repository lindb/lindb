package metadb

import (
	"context"
	"sync"
	"time"

	"github.com/lindb/lindb/constants"
	"github.com/lindb/lindb/pkg/logger"
	"github.com/lindb/lindb/pkg/timeutil"
	"github.com/lindb/lindb/series"
	"github.com/lindb/lindb/series/field"
	"github.com/lindb/lindb/series/tag"
)

// for testing
var (
	createMetadataBackend = newMetadataBackend
)

var (
	syncInterval = 2 * timeutil.OneSecond
)

// metadataDatabase implements the MetadataDatabase interface,
// !!!!NOTICE: need cache all tag keys/fields of metric
type metadataDatabase struct {
	name    string // tsdb's name
	ctx     context.Context
	cancel  context.CancelFunc
	backend MetadataBackend
	metrics map[string]MetricMetadata // metadata cache(key: namespace + metric-name, value: metric metadata)

	mutable      *metadataUpdateEvent // pending update events
	immutable    *metadataUpdateEvent // syncing pending update events
	lastSyncTime int64
	syncSignal   chan struct{}

	syncInterval int64

	rwMux sync.RWMutex
}

// NewMetadataDatabase creates new metadata database
func NewMetadataDatabase(ctx context.Context, name, parent string) (MetadataDatabase, error) {
	backend, err := createMetadataBackend(name, parent)
	if err != nil {
		return nil, err
	}
	c, cancel := context.WithCancel(ctx)
	mdb := &metadataDatabase{
		name:         name,
		ctx:          c,
		cancel:       cancel,
		backend:      backend,
		mutable:      newMetadataUpdateEvent(),
		metrics:      make(map[string]MetricMetadata),
		syncSignal:   make(chan struct{}),
		syncInterval: syncInterval,
		lastSyncTime: timeutil.Now(),
	}
	go mdb.checkSync()
	go mdb.syncPendingEvent()
	return mdb, nil
}

// SuggestNamespace suggests the namespace by namespace's prefix
func (mdb *metadataDatabase) SuggestNamespace(prefix string, limit int) (namespaces []string, err error) {
	return mdb.backend.suggestNamespace(prefix, limit)
}

// SuggestMetricName suggests the metric name by name's prefix
func (mdb *metadataDatabase) SuggestMetricName(namespace, prefix string, limit int) (namespaces []string, err error) {
	return mdb.backend.suggestMetricName(namespace, prefix, limit)
}

// GetMetricID gets the metric id by namespace and metric name, if not exist return constants.ErrNotFound
func (mdb *metadataDatabase) GetMetricID(namespace, metricName string) (metricID uint32, err error) {
	mdb.rwMux.RLock()
	// read from memory
	key := namespace + metricName
	metricMetadata, ok := mdb.metrics[key]
	if ok {
		defer mdb.rwMux.RUnlock()
		return metricMetadata.getMetricID(), nil
	}
	mdb.rwMux.RUnlock()

	// read from meta db
	return mdb.backend.getMetricID(namespace, metricName)
}

// GetTagKeyID gets the tag key id by namespace/metric name/tag key key, if not exist return constants.ErrNotFound
func (mdb *metadataDatabase) GetTagKeyID(namespace, metricName string, tagKey string) (tagKeyID uint32, err error) {
	key := namespace + metricName

	mdb.rwMux.RLock()
	metricMetadata, ok := mdb.metrics[key]
	if ok {
		defer mdb.rwMux.RUnlock()
		tagKeyID, ok = metricMetadata.getTagKeyID(tagKey)
		if ok {
			return
		}
		return 0, constants.ErrNotFound
	}
	mdb.rwMux.RUnlock()

	metricID, err := mdb.backend.getMetricID(namespace, metricName)
	if err != nil {
		return 0, err
	}

	return mdb.backend.getTagKeyID(metricID, tagKey)
}

// GetAllTagKeys returns the all tag keys by namespace/metric name, if not exist return constants.ErrNotFound
func (mdb *metadataDatabase) GetAllTagKeys(namespace, metricName string) (tags []tag.Meta, err error) {
	key := namespace + metricName
	mdb.rwMux.RLock()
	metricMetadata, ok := mdb.metrics[key]
	if ok {
		defer mdb.rwMux.RUnlock()
		return metricMetadata.getAllTagKeys(), nil
	}
	mdb.rwMux.RUnlock()

	metricID, err := mdb.backend.getMetricID(namespace, metricName)
	if err != nil {
		return
	}

	return mdb.backend.getAllTagKeys(metricID)
}

// GetField gets the field meta by namespace/metric name/field name, if not exist return constants.ErrNotFound
func (mdb *metadataDatabase) GetField(namespace, metricName, fieldName string) (f field.Meta, err error) {
	key := namespace + metricName
	mdb.rwMux.RLock()
	metricMetadata, ok := mdb.metrics[key]
	if ok {
		defer mdb.rwMux.RUnlock()
		f, ok = metricMetadata.getField(fieldName)
		if ok {
			return f, nil
		}
		return field.Meta{}, constants.ErrNotFound
	}
	mdb.rwMux.RUnlock()
	metricID, err := mdb.GetMetricID(namespace, metricName)
	if err != nil {
		return field.Meta{}, err
	}

	// read from db
	return mdb.backend.getField(metricID, fieldName)
}

func (mdb *metadataDatabase) GetAllFields(namespace, metricName string) (fields []field.Meta, err error) {
	key := namespace + metricName
	mdb.rwMux.RLock()
	metricMetadata, ok := mdb.metrics[key]
	if ok {
		defer mdb.rwMux.RUnlock()
		return metricMetadata.getAllFields(), nil
	}
	mdb.rwMux.RUnlock()
	metricID, err := mdb.GetMetricID(namespace, metricName)
	if err != nil {
		return nil, err
	}
	return mdb.backend.getAllFields(metricID)
}

// GenMetricID generates the metric id in the memory.
// 1) get metric id from memory if exist, if not exist goto 2
// 2) get metric metadata from backend storage, if not exist need create new metric metadata
func (mdb *metadataDatabase) GenMetricID(namespace, metricName string) (metricID uint32, err error) {
	key := namespace + metricName
	mdb.rwMux.RLock()
	// get metric id from memory, add read lock
	metricMetadata, ok := mdb.metrics[key]
	if ok {
		mdb.rwMux.RUnlock()
		return metricMetadata.getMetricID(), nil
	}
	mdb.rwMux.RUnlock()

	// assign metric id from memory, add write lock
	mdb.rwMux.Lock()
	defer mdb.rwMux.Unlock()
	// double check with memory
	metricMetadata, ok = mdb.metrics[key]
	if ok {
		return metricMetadata.getMetricID(), nil
	}

	// load metric metadata from backend storage
	metricMetadata, err = mdb.backend.loadMetricMetadata(namespace, metricName)
	if err == nil {
		// get metric metadata from backend
		mdb.metrics[key] = metricMetadata
		return metricMetadata.getMetricID(), nil
	}
	// isn't not found, return err
	if err != constants.ErrNotFound {
		return
	}
	// assign new metric id
	metricID = mdb.backend.genMetricID()
	mdb.metrics[key] = newMetricMetadata(metricID, 0)
	mdb.mutable.addMetric(namespace, metricName, metricID)

	mdb.notifySyncWithoutLock(false)
	return metricID, nil
}

// GenFieldID generates the field id in the memory,
// !!!!! NOTICE: metric metadata must be exist in memory, because gen metric has been saved
func (mdb *metadataDatabase) GenFieldID(namespace, metricName string,
	fieldName string, fieldType field.Type,
) (fieldID field.ID, err error) {
	key := namespace + metricName

	mdb.rwMux.Lock()
	defer mdb.rwMux.Unlock()
	// read from memory metric metadata
	metricMetadata := mdb.metrics[key]
	f, ok := metricMetadata.getField(fieldName)
	if ok {
		if f.Type == fieldType {
			return f.ID, nil
		}
		return 0, series.ErrWrongFieldType
	}
	// assign new field id
	fieldID, err = metricMetadata.createField(fieldName, fieldType)
	if err != nil {
		return 0, err
	}
	mdb.mutable.addField(metricMetadata.getMetricID(), field.Meta{
		ID:   fieldID,
		Type: fieldType,
		Name: fieldName,
	})
	mdb.notifySyncWithoutLock(false)
	return
}

// GenTagKeyID generates the tag key id in the memory
// !!!!! NOTICE: metric metadata must be exist in memory, because gen metric has been saved
func (mdb *metadataDatabase) GenTagKeyID(namespace, metricName, tagKey string) (tagKeyID uint32, err error) {
	key := namespace + metricName

	mdb.rwMux.Lock()
	defer mdb.rwMux.Unlock()
	// read from memory metric metadata
	metricMetadata := mdb.metrics[key]
	tagKeyID, ok := metricMetadata.getTagKeyID(tagKey)
	if ok {
		return tagKeyID, nil
	}
	// check tag keys count before create
	if err = metricMetadata.checkTagKeyCount(); err != nil {
		return 0, err
	}
	// assign new tag key id
	tagKeyID = mdb.backend.genTagKeyID()
	metricMetadata.createTagKey(tagKey, tagKeyID)
	mdb.mutable.addTagKey(metricMetadata.getMetricID(), tag.Meta{
		Key: tagKey,
		ID:  tagKeyID,
	})
	mdb.notifySyncWithoutLock(false)
	return
}

// Sync syncs the bbolt.DB's data file
func (mdb *metadataDatabase) Sync() error {
	//FIXME stone100 need impl sync force when flush metric data
	return mdb.backend.sync()
}

// Close closes the resources
func (mdb *metadataDatabase) Close() error {
	mdb.cancel()
	mdb.rwMux.Lock()
	defer mdb.rwMux.Unlock()

	if mdb.mutable != nil && !mdb.mutable.isEmpty() {
		if err := mdb.backend.saveMetadata(mdb.mutable); err != nil {
			return err
		}
	}
	if mdb.immutable != nil && !mdb.mutable.isEmpty() {
		if err := mdb.backend.saveMetadata(mdb.immutable); err != nil {
			return err
		}
	}
	return mdb.backend.Close()
}

// checkSync checks if need sync pending metadata event in period
func (mdb *metadataDatabase) checkSync() {
	ticker := time.NewTicker(time.Duration(mdb.syncInterval * 1000000))
	for {
		select {
		case <-ticker.C:
			mdb.notifySyncWithLock(false)
		case <-mdb.ctx.Done():
			ticker.Stop()
			metaLogger.Info("check metadata event update goroutine exit...", logger.String("db", mdb.name))
			return
		}
	}
}

// notifySyncWithoutLock notifies sync goroutine need save pending metadata events with lock
func (mdb *metadataDatabase) notifySyncWithLock(force bool) {
	mdb.rwMux.Lock()
	defer mdb.rwMux.Unlock()

	mdb.notifySyncWithoutLock(force)
}

// notifySyncWithoutLock notifies sync goroutine need save pending metadata events without lock
func (mdb *metadataDatabase) notifySyncWithoutLock(force bool) {
	if (!mdb.mutable.isFull() || !force) && timeutil.Now()-mdb.lastSyncTime < mdb.syncInterval {
		return
	}

	if !mdb.mutable.isEmpty() && mdb.immutable == nil {
		mdb.immutable = mdb.mutable
		mdb.mutable = newMetadataUpdateEvent()
		// notify with time out
		select {
		case mdb.syncSignal <- struct{}{}:
		case <-time.After(time.Second):
			//FIXME add metric
			metaLogger.Error("notify sync metadata save timeout", logger.String("db", mdb.name))
		}
	}
}

// syncPendingEvent syncs the pending metadata event
func (mdb *metadataDatabase) syncPendingEvent() {
	for {
		select {
		case <-mdb.syncSignal:
			var event *metadataUpdateEvent
			mdb.rwMux.RLock()
			event = mdb.immutable
			mdb.rwMux.RUnlock()
			if event == nil {
				continue
			}
			if err := mdb.backend.saveMetadata(event); err != nil {
				//FIXME stone1100 add metric
				metaLogger.Error("save metadata err", logger.String("db", mdb.name), logger.Error(err))
				continue
			}
			mdb.rwMux.Lock()
			mdb.immutable = nil
			mdb.lastSyncTime = timeutil.Now()
			mdb.rwMux.Unlock()
		case <-mdb.ctx.Done():
			metaLogger.Info("sync update event goroutine exit...", logger.String("db", mdb.name))
			return
		}
	}
}
