package indexdb

import (
	"sync"

	"github.com/cespare/xxhash"

	"github.com/lindb/lindb/tsdb/metadb"
	"github.com/lindb/lindb/tsdb/tblstore/invertedindex"
)

type memoryIndexDatabase struct {
	metricHash2Index sync.Map           // key: xxhash.Sum64String(metric-name), value: metric index
	generator        metadb.IDGenerator // the generator for generating ID of metric, field
	index            InvertedIndex

	lock sync.RWMutex // lock of create metric index
}

func NewMemoryIndexDatabase(generator metadb.IDGenerator) MemoryIndexDatabase {
	return &memoryIndexDatabase{
		generator: generator,
		index:     newInvertedIndex(generator),
	}
}

func (db *memoryIndexDatabase) GetTimeSeriesID(metricName string,
	tags map[string]string, tagsHash uint64,
) (metricID, seriesID uint32) {
	hash := xxhash.Sum64String(metricName)
	metricIDMappingINTF, ok := db.metricHash2Index.Load(hash)
	if !ok {
		// not found need create new metric id mapping
		db.lock.Lock()
		defer db.lock.Unlock()

		// double check metric id mapping if exist
		metricIDMappingINTF, ok = db.metricHash2Index.Load(hash)
		if !ok {
			// creates new metric id mapping
			// gen new metric id for new metric id mapping
			metricID := db.generator.GenMetricID(metricName)

			metricIDMappingINTF = newMetricIDMapping(metricID)
			db.metricHash2Index.Store(hash, metricIDMappingINTF)
		}
	}
	idMapping := metricIDMappingINTF.(MetricIDMapping)
	seriesID, created := idMapping.GetOrCreateSeriesID(tagsHash)
	if created {
		//FIXME stone100 need add goroutine
		db.index.buildInvertIndex(metricID, tags, seriesID)
	}
	metricID = idMapping.GetMetricID()
	return
}

// FlushInvertedIndexTo flushes the series data to a inverted-index file.
func (db *memoryIndexDatabase) FlushInvertedIndexTo(flusher invertedindex.Flusher) (err error) {
	//db.metricHash2Index.Range(func(key, value interface{}) bool {
	//	metricIDMapping := value.(MetricIDMapping)
	//	if err = metricIDMapping.FlushInvertedIndexTo(flusher, db.generator); err != nil {
	//		//FIXME need add log
	//		return true
	//	}
	//	return true
	//})
	//return flusher.Commit()
	return nil
}
