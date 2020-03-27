package indexdb

import (
	"context"
	"sync"
	"time"

	"github.com/lindb/roaring"

	"github.com/lindb/lindb/constants"
	"github.com/lindb/lindb/kv"
	"github.com/lindb/lindb/pkg/logger"
	"github.com/lindb/lindb/pkg/timeutil"
	"github.com/lindb/lindb/series"
	"github.com/lindb/lindb/tsdb/metadb"
)

// for testing
var (
	createBackend = newIDMappingBackend
)

var (
	syncInterval = 2 * timeutil.OneSecond
)

// indexDatabase implements IndexDatabase interface
type indexDatabase struct {
	name             string
	ctx              context.Context
	cancel           context.CancelFunc
	backend          IDMappingBackend           // id mapping backend storage
	metricID2Mapping map[uint32]MetricIDMapping // key: metric id, value: metric id mapping
	metadata         metadb.Metadata            // the metadata for generating ID of metric, field
	index            InvertedIndex

	mutable      *mappingEvent // pending update events
	immutable    *mappingEvent // syncing pending update events
	lastSyncTime int64
	syncSignal   chan struct{}

	syncInterval int64

	rwMutex sync.RWMutex // lock of create metric index
}

func (db *indexDatabase) SuggestTagValues(tagKeyID uint32, tagValuePrefix string, limit int) []string {
	//panic("implement me")
	return db.metadata.TagMetadata().SuggestTagValues(tagKeyID, tagValuePrefix, limit)
}

// GetGroupingContext returns the context of group by
func (db *indexDatabase) GetGroupingContext(tagKeyIDs []uint32, seriesIDs *roaring.Bitmap) (series.GroupingContext, error) {
	return db.index.GetGroupingContext(tagKeyIDs, seriesIDs)
}

// NewIndexDatabase creates a new index database
func NewIndexDatabase(ctx context.Context, name, parent string, metadata metadb.Metadata,
	forwardFamily kv.Family, invertedFamily kv.Family,
) (IndexDatabase, error) {
	backend, err := createBackend(name, parent)
	if err != nil {
		return nil, err
	}
	c, cancel := context.WithCancel(ctx)
	db := &indexDatabase{
		name:             name,
		ctx:              c,
		cancel:           cancel,
		backend:          backend,
		metadata:         metadata,
		metricID2Mapping: make(map[uint32]MetricIDMapping),
		index:            newInvertedIndex(metadata, invertedFamily, forwardFamily),
		mutable:          newMappingEvent(),
		lastSyncTime:     timeutil.Now(),
		syncSignal:       make(chan struct{}),
		syncInterval:     syncInterval,
	}
	go db.checkSync()
	go db.syncPendingEvent()
	return db, nil
}

// GetOrCreateSeriesID gets series by tags hash, if not exist generate new series id in memory,
// if generate a new series id returns isCreate is true
// if generate fail return err
func (db *indexDatabase) GetOrCreateSeriesID(metricID uint32, tagsHash uint64,
) (seriesID uint32, isCreated bool, err error) {
	db.rwMutex.Lock()
	defer db.rwMutex.Unlock()

	metricIDMapping, ok := db.metricID2Mapping[metricID]
	if ok {
		// get series id from memory cache
		seriesID, ok = metricIDMapping.GetSeriesID(tagsHash)
		if ok {
			return seriesID, false, nil
		}
	} else {
		// metric mapping not exist, need load from backend storage
		metricIDMapping, err = db.backend.loadMetricIDMapping(metricID)
		if err != nil && err != constants.ErrNotFound {
			return 0, false, err
		}
		// if metric id not exist in backend storage
		if err == constants.ErrNotFound {
			// create new metric id mapping with 0 sequence
			metricIDMapping = newMetricIDMapping(metricID, 0)
			// cache metric id mapping
			db.metricID2Mapping[metricID] = metricIDMapping
		} else {
			// cache metric id mapping
			db.metricID2Mapping[metricID] = metricIDMapping
			// metric id mapping exist, try get series id from backend storage
			seriesID, err = db.backend.getSeriesID(metricID, tagsHash)
			if err == nil {
				// cache load series id
				metricIDMapping.AddSeriesID(tagsHash, seriesID)
				return seriesID, false, nil
			}
		}
	}
	// throw err in backend storage
	if err != nil && err != constants.ErrNotFound {
		return 0, false, err
	}
	// generate new series id
	seriesID = metricIDMapping.GenSeriesID(tagsHash)

	// add pending event
	db.mutable.addSeriesID(metricID, tagsHash, seriesID)
	db.notifySyncWithoutLock(false)
	return seriesID, true, nil
}

// GetSeriesIDsByTagValueIDs gets series ids by tag value ids for spec metric's tag key
func (db *indexDatabase) GetSeriesIDsByTagValueIDs(tagKeyID uint32, tagValueIDs *roaring.Bitmap) (*roaring.Bitmap, error) {
	return db.index.GetSeriesIDsByTagValueIDs(tagKeyID, tagValueIDs)
}

// GetSeriesIDsForTag gets series ids for spec metric's tag key
func (db *indexDatabase) GetSeriesIDsForTag(tagKeyID uint32) (*roaring.Bitmap, error) {
	return db.index.GetSeriesIDsForTag(tagKeyID)
}

// GetSeriesIDsForMetric gets series ids for spec metric name
func (db *indexDatabase) GetSeriesIDsForMetric(namespace, metricName string) (*roaring.Bitmap, error) {
	// get all tags under metric
	tags, err := db.metadata.MetadataDatabase().GetAllTagKeys(namespace, metricName)
	if err != nil {
		return nil, err
	}
	tagLength := len(tags)
	if tagLength == 0 {
		// if metric hasn't any tags, returns default series id(0)
		return roaring.BitmapOf(constants.SeriesIDWithoutTags), nil
	}
	tagKeyIDs := make([]uint32, tagLength)
	for idx, tag := range tags {
		tagKeyIDs[idx] = tag.ID
	}
	// get series ids under all tag key ids
	return db.index.GetSeriesIDsForTags(tagKeyIDs)
}

// BuildInvertIndex builds the inverted index for tag value => series ids,
// the tags is considered as a empty key-value pair while tags is nil.
func (db *indexDatabase) BuildInvertIndex(namespace, metricName string, tags map[string]string, seriesID uint32) {
	db.index.buildInvertIndex(namespace, metricName, tags, seriesID)
}

// Flush flushes index data to disk
func (db *indexDatabase) Flush() error {
	//fixme flush mapping data
	//fixme inverted index need add wal???
	return db.index.Flush()
}

// Close closes the database, releases the resources
func (db *indexDatabase) Close() error {
	db.cancel()
	db.rwMutex.Lock()
	defer db.rwMutex.Unlock()
	saveMapping := func(event *mappingEvent) error {
		if event == nil {
			return nil
		}
		if !event.isEmpty() {
			if err := db.backend.saveMapping(event); err != nil {
				return err
			}
		}
		return nil
	}
	if err := saveMapping(db.mutable); err != nil {
		return err
	}
	if err := saveMapping(db.immutable); err != nil {
		return err
	}
	if err := db.backend.Close(); err != nil {
		return err
	}
	return db.index.Flush()
}

// checkSync checks if need sync pending series event in period
func (db *indexDatabase) checkSync() {
	ticker := time.NewTicker(time.Duration(db.syncInterval * 1000000))
	for {
		select {
		case <-ticker.C:
			db.notifySyncWithLock(false)
		case <-db.ctx.Done():
			ticker.Stop()
			indexLogger.Info("check series event update goroutine exit...", logger.String("db", db.name))
			return
		}
	}
}

// notifySyncWithoutLock notifies sync goroutine need save pending series events without lock
func (db *indexDatabase) notifySyncWithoutLock(force bool) {
	if (!db.mutable.isFull() || !force) && timeutil.Now()-db.lastSyncTime < db.syncInterval {
		return
	}

	if !db.mutable.isEmpty() && db.immutable == nil {
		db.immutable = db.mutable
		db.mutable = newMappingEvent()
		// notify with time out
		select {
		case <-time.After(time.Second):
			//FIXME add metric
			indexLogger.Error("notify sync series save timeout", logger.String("db", db.name))
		case db.syncSignal <- struct{}{}:
		}
	}
}

// notifySyncWithoutLock notifies sync goroutine need save pending series events with lock
func (db *indexDatabase) notifySyncWithLock(force bool) {
	db.rwMutex.Lock()
	defer db.rwMutex.Unlock()

	db.notifySyncWithoutLock(force)
}

// syncPendingEvent syncs the pending series event
func (db *indexDatabase) syncPendingEvent() {
	for {
		select {
		case <-db.ctx.Done():
			indexLogger.Info("sync update event goroutine exit...", logger.String("db", db.name))
			return
		case <-db.syncSignal:
			var event *mappingEvent
			db.rwMutex.RLock()
			event = db.immutable
			db.rwMutex.RUnlock()
			if event == nil {
				continue
			}
			if err := db.backend.saveMapping(event); err != nil {
				//FIXME stone1100 add metric
				indexLogger.Error("save mapping err", logger.String("db", db.name), logger.Error(err))
				continue
			}
			db.rwMutex.Lock()
			db.immutable = nil
			db.lastSyncTime = timeutil.Now()
			db.rwMutex.Unlock()
		}
	}
}
