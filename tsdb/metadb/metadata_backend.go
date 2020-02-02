package metadb

import (
	"bytes"
	"encoding/binary"
	"io"
	"path"
	"time"

	"github.com/coreos/bbolt"
	"go.uber.org/atomic"

	"github.com/lindb/lindb/pkg/fileutil"
	"github.com/lindb/lindb/pkg/logger"
	"github.com/lindb/lindb/series"
	"github.com/lindb/lindb/series/field"
	"github.com/lindb/lindb/series/tag"
)

//go:generate mockgen -source ./metadata_backend.go -destination=./metadata_backend_mock.go -package=metadb

const MetaDB = "meta.db"

// for testing
var (
	mkDir            = fileutil.MkDirIfNotExist
	closeFunc        = closeDB
	setSequenceFunc  = setSequence
	createBucketFunc = createBucket
)

var (
	nsBucketName     = []byte("ns")
	metricBucketName = []byte("m")
	tagBucketName    = []byte("t")
	fieldBucketName  = []byte("f")
)

// MetadataBackend represents the metadata backend storage
type MetadataBackend interface {
	io.Closer

	// suggestNamespace suggests the namespace by namespace's prefix
	suggestNamespace(prefix string, limit int) (namespaces []string, err error)
	// suggestMetricName suggests the metric name by name's prefix
	suggestMetricName(namespace, prefix string, limit int) (metricNames []string, err error)

	// genMetricID generates the metric id in the memory
	genMetricID() uint32
	// genTagKeyID generates the tag key id in the memory
	genTagKeyID() uint32

	// loadMetricMetadata loads the metric metadata include all fields/tags by namespace and metric name,
	// if not exist return series.ErrNotFound
	loadMetricMetadata(namespace, metricName string) (MetricMetadata, error)
	// getMetricMetadata gets the metric metadata include all fields/tags by metric id, if not exist return series.ErrNotFound
	getMetricMetadata(metricID uint32) (metadata MetricMetadata, err error)

	// getMetricID gets the metric id by namespace and metric name, if not exist return series.ErrNotFound
	getMetricID(namespace string, metricName string) (metricID uint32, err error)
	// getTagKeyID gets the tag key id by metric id and tag key key, if not exist return series.ErrNotFound
	getTagKeyID(metricID uint32, tagKey string) (tagKeyID uint32, err error)
	// getAllTagKeys returns the all tag keys by metric id, if not exist return series.ErrNotFound
	getAllTagKeys(metricID uint32) (tags []tag.Meta, err error)
	// getField gets the field meta by metric id and field name, if not exist return series.ErrNotFound
	getField(metricID uint32, fieldName string) (f field.Meta, err error)
	// getAllFields returns the  all fields by metric id, if not exist return series.ErrNotFound
	getAllFields(metricID uint32) (fields []field.Meta, err error)

	// saveMetadata saves the pending metadata include namespace/metric metadata
	saveMetadata(event *metadataUpdateEvent) error

	// sync syncs bbolt.DB file data
	sync() error
}

// metadataBackend implements the MetadataBackend interface
type metadataBackend struct {
	db               *bbolt.DB
	metricIDSequence atomic.Uint32
	tagKeyIDSequence atomic.Uint32
}

// newMetadataBackend creates a new metadata backend storage
func newMetadataBackend(parent string) (MetadataBackend, error) {
	if err := mkDir(parent); err != nil {
		return nil, err
	}
	db, err := bbolt.Open(path.Join(parent, MetaDB), 0600, &bbolt.Options{Timeout: 1 * time.Second, NoSync: true})
	if err != nil {
		return nil, err
	}

	var metricIDSequence atomic.Uint32
	var tagKeyIDSequence atomic.Uint32
	err = db.Update(func(tx *bbolt.Tx) error {
		// create namespace bucket for save namespace/metric
		nsBucket, err := tx.CreateBucketIfNotExists(nsBucketName)
		if err != nil {
			return err
		}
		// load metric id sequence
		metricIDSequence.Store(uint32(nsBucket.Sequence()))
		// create metric bucket for save metric metadata
		metricBucket, err := tx.CreateBucketIfNotExists(metricBucketName)
		if err != nil {
			return err
		}
		// load tag key id sequence
		tagKeyIDSequence.Store(uint32(metricBucket.Sequence()))
		return nil
	})
	if err != nil {
		// close bbolt.DB if init metadata err
		if e := closeFunc(db); e != nil {
			metaLogger.Error("close bbolt.db err when create metadata backend fail", logger.Error(e))
		}
		return nil, err
	}
	return &metadataBackend{
		db:               db,
		metricIDSequence: metricIDSequence,
		tagKeyIDSequence: tagKeyIDSequence,
	}, err
}

// suggestNamespace suggests the namespace by namespace's prefix
func (mb *metadataBackend) suggestNamespace(prefix string, limit int) (namespaces []string, err error) {
	err = mb.db.View(func(tx *bbolt.Tx) error {
		cursor := tx.Bucket(nsBucketName).Cursor()
		prefix := []byte(prefix)
		for k, _ := cursor.Seek(prefix); k != nil && bytes.HasPrefix(k, prefix); k, _ = cursor.Next() {
			namespaces = append(namespaces, string(k))
			if len(namespaces) >= limit {
				return nil
			}
		}
		return nil
	})
	return
}

// suggestMetricName suggests the metric name by name's prefix
func (mb *metadataBackend) suggestMetricName(namespace, prefix string, limit int) (metricNames []string, err error) {
	err = mb.db.View(func(tx *bbolt.Tx) error {
		// 1. get namespace bucket
		nsBucket := tx.Bucket(nsBucketName).Bucket([]byte(namespace))
		if nsBucket == nil {
			return nil
		}

		// 2. scan metric name by prefix
		cursor := nsBucket.Cursor()
		prefix := []byte(prefix)
		for k, _ := cursor.Seek(prefix); k != nil && bytes.HasPrefix(k, prefix); k, _ = cursor.Next() {
			metricNames = append(metricNames, string(k))
			if len(metricNames) >= limit {
				return nil
			}
		}
		return nil
	})
	return
}

// genMetricID generates the metric id in the memory
func (mb *metadataBackend) genMetricID() uint32 {
	return mb.metricIDSequence.Inc()
}

// genTagKeyID generates the tag key id in the memory
func (mb *metadataBackend) genTagKeyID() uint32 {
	return mb.tagKeyIDSequence.Inc()
}

// loadMetricMetadata loads the metric metadata include all fields/tags by namespace and metric name,
// if not exist return series.ErrNotFound
func (mb *metadataBackend) loadMetricMetadata(namespace, metricName string) (metadata MetricMetadata, err error) {
	metricID, err := mb.getMetricID(namespace, metricName)
	if err != nil {
		return nil, err
	}
	return mb.getMetricMetadata(metricID)
}

// getMetricMetadata gets the metric metadata include all fields/tags by metric id, if not exist return series.ErrNotFound
func (mb *metadataBackend) getMetricMetadata(metricID uint32) (metadata MetricMetadata, err error) {
	var scratch [4]byte
	var fieldIDSeq int32
	var tags []tag.Meta
	var fields []field.Meta
	binary.LittleEndian.PutUint32(scratch[:], metricID)
	err = mb.db.View(func(tx *bbolt.Tx) error {
		metricBucket := tx.Bucket(metricBucketName).Bucket(scratch[:])
		if metricBucket == nil {
			return series.ErrNotFound
		}
		tags = loadTagKeys(metricBucket.Bucket(tagBucketName))
		fBucket := metricBucket.Bucket(fieldBucketName)
		fieldIDSeq = int32(fBucket.Sequence())
		fields = loadFields(fBucket)
		return nil
	})
	if err != nil {
		return
	}

	metadata = newMetricMetadata(metricID, fieldIDSeq)
	// initialize fields and tags
	metadata.initialize(fields, tags)
	return
}

// getMetricID gets the metric id by namespace and metric name, if not exist return series.ErrNotFound
func (mb *metadataBackend) getMetricID(namespace string, metricName string) (metricID uint32, err error) {
	err = mb.db.View(func(tx *bbolt.Tx) error {
		nsBucket := tx.Bucket(nsBucketName).Bucket([]byte(namespace))
		if nsBucket == nil {
			return series.ErrNotFound
		}
		value := nsBucket.Get([]byte(metricName))
		if len(value) == 0 {
			return series.ErrNotFound
		}
		metricID = binary.LittleEndian.Uint32(value)
		return nil
	})
	return
}

// getTagKeyID gets the tag key id by metric id and tag key key, if not exist return series.ErrNotFound
func (mb *metadataBackend) getTagKeyID(metricID uint32, tagKey string) (tagKeyID uint32, err error) {
	var scratch [4]byte
	binary.LittleEndian.PutUint32(scratch[:], metricID)
	err = mb.db.View(func(tx *bbolt.Tx) error {
		metricBucket := tx.Bucket(metricBucketName).Bucket(scratch[:])
		if metricBucket == nil {
			return series.ErrNotFound
		}
		value := metricBucket.Bucket(tagBucketName).Get([]byte(tagKey))
		if len(value) == 0 {
			return series.ErrNotFound
		}
		tagKeyID = binary.LittleEndian.Uint32(value)
		return nil
	})
	return
}

// getAllTagKeys returns the all tag keys by metric id, if not exist return series.ErrNotFound
func (mb *metadataBackend) getAllTagKeys(metricID uint32) (tags []tag.Meta, err error) {
	var scratch [4]byte
	binary.LittleEndian.PutUint32(scratch[:], metricID)
	err = mb.db.View(func(tx *bbolt.Tx) error {
		metricBucket := tx.Bucket(metricBucketName).Bucket(scratch[:])
		if metricBucket == nil {
			return series.ErrNotFound
		}
		tags = loadTagKeys(metricBucket.Bucket(tagBucketName))
		return nil
	})
	return
}

// getField gets the field meta by metric id and field name, if not exist return series.ErrNotFound
func (mb *metadataBackend) getField(metricID uint32, fieldName string) (f field.Meta, err error) {
	var scratch [4]byte
	binary.LittleEndian.PutUint32(scratch[:], metricID)
	err = mb.db.View(func(tx *bbolt.Tx) error {
		metricBucket := tx.Bucket(metricBucketName).Bucket(scratch[:])
		if metricBucket == nil {
			return series.ErrNotFound
		}
		value := metricBucket.Bucket(fieldBucketName).Get([]byte(fieldName))
		if len(value) == 0 {
			return series.ErrNotFound
		}
		f.Name = fieldName
		f.Type = field.Type(value[0])
		f.ID = binary.LittleEndian.Uint16(value[1:])
		return nil
	})
	return
}

// getAllFields returns the  all fields by metric id, if not exist return series.ErrNotFound
func (mb *metadataBackend) getAllFields(metricID uint32) (fields []field.Meta, err error) {
	var scratch [4]byte
	binary.LittleEndian.PutUint32(scratch[:], metricID)
	err = mb.db.View(func(tx *bbolt.Tx) error {
		metricBucket := tx.Bucket(metricBucketName).Bucket(scratch[:])
		if metricBucket == nil {
			return series.ErrNotFound
		}
		fields = loadFields(metricBucket.Bucket(fieldBucketName))
		return nil
	})
	return
}

// saveMetadata saves the pending metadata include namespace/metric metadata
func (mb *metadataBackend) saveMetadata(event *metadataUpdateEvent) (err error) {
	err = mb.db.Update(func(tx *bbolt.Tx) error {
		if err := saveNamespaceAndMetric(tx.Bucket(nsBucketName), event); err != nil {
			return err
		}
		if err := saveMetricMetadata(tx.Bucket(metricBucketName), event); err != nil {
			return err
		}
		return nil
	})
	return
}

// sync syncs the bbolt.DB file data
func (mb *metadataBackend) sync() error {
	return mb.db.Sync()
}

// Close closes the bbolt.DB
func (mb *metadataBackend) Close() error {
	return mb.db.Close()
}

// loadFields loads the fields from field bucket
func loadFields(fieldBucket *bbolt.Bucket) (fields []field.Meta) {
	cursor := fieldBucket.Cursor()
	for k, v := cursor.First(); k != nil; k, v = cursor.Next() {
		fields = append(fields, field.Meta{
			Name: string(k),
			Type: field.Type(v[0]),
			ID:   binary.LittleEndian.Uint16(v[1:]),
		})
	}
	return
}

// loadTagKeys loads the tag keys from tag key bucket
func loadTagKeys(tagKeyBucket *bbolt.Bucket) (tags []tag.Meta) {
	cursor := tagKeyBucket.Cursor()
	for k, v := cursor.First(); k != nil; k, v = cursor.Next() {
		tags = append(tags, tag.Meta{
			Key: string(k),
			ID:  binary.LittleEndian.Uint32(v),
		})
	}
	return
}

// saveNamespaceAndMetric saves namespaces and metric entry set
func saveNamespaceAndMetric(nsRootBucket *bbolt.Bucket, event *metadataUpdateEvent) (err error) {
	for ns, nsEvent := range event.namespaces {
		// save namespace name
		bucket, err := nsRootBucket.CreateBucketIfNotExists([]byte(ns))
		if err != nil {
			return err
		}
		// save metric entry
		for _, metric := range nsEvent.metrics {
			var scratch [4]byte
			binary.LittleEndian.PutUint32(scratch[:], metric.id)
			if err := bucket.Put([]byte(metric.name), scratch[:]); err != nil {
				return err
			}
		}
	}
	// final set metric id sequence
	if event.metricSeqID > 0 {
		if err = setSequenceFunc(nsRootBucket, uint64(event.metricSeqID)); err != nil {
			return err
		}
	}
	return
}

// saveMetricMetadata saves metric metadata include fields/tag keys if exist with metric root bucket
func saveMetricMetadata(metricRootBucket *bbolt.Bucket, event *metadataUpdateEvent) (err error) {
	for metricID, meta := range event.metrics {
		var scratch [4]byte
		binary.LittleEndian.PutUint32(scratch[:], metricID)
		mID := scratch[:]
		metricBucket := metricRootBucket.Bucket(mID)
		var fBucket *bbolt.Bucket
		var tBucket *bbolt.Bucket
		if metricBucket == nil {
			// if metric meta not exist, initialize metric bucket
			metricBucket, err = createBucketFunc(metricRootBucket, mID)
			if err != nil {
				return err
			}
			fBucket, err = metricBucket.CreateBucket(fieldBucketName)
			if err != nil {
				return err
			}
			tBucket, err = metricBucket.CreateBucket(tagBucketName)
			if err != nil {
				return err
			}
		}
		// save metric's fields
		if len(meta.fields) > 0 {
			if fBucket == nil {
				// for load field bucket for exist metric id
				fBucket = metricBucket.Bucket(fieldBucketName)
			}
			if err = saveFields(fBucket, meta.fieldIDSeq, meta.fields); err != nil {
				return err
			}
		}
		// save metric's tag keys
		if len(meta.tagKeys) > 0 {
			if tBucket == nil {
				// for tag key bucket for exist metric id
				tBucket = metricBucket.Bucket(tagBucketName)
			}
			if err = saveTagKeys(tBucket, meta.tagKeys); err != nil {
				return err
			}
		}
	}

	// final set tag key id sequence
	if event.tagKeySeqID > 0 {
		if err = setSequenceFunc(metricRootBucket, uint64(event.tagKeySeqID)); err != nil {
			return err
		}
	}
	return nil
}

// saveFields saves fields for metric with field bucket
func saveFields(fieldBucket *bbolt.Bucket, fieldIDSeq uint16, fields []field.Meta) (err error) {
	for _, f := range fields {
		var fieldValue [3]byte
		fieldValue[0] = byte(f.Type)
		binary.LittleEndian.PutUint16(fieldValue[1:], f.ID)
		if err = fieldBucket.Put([]byte(f.Name), fieldValue[:]); err != nil {
			return err
		}
	}
	// save field id sequence
	if err = setSequenceFunc(fieldBucket, uint64(fieldIDSeq)); err != nil {
		return err
	}
	return
}

// saveTagKeys saves tag keys for metric with tag key bucket
func saveTagKeys(tagKeyBucket *bbolt.Bucket, tagKeys []tag.Meta) (err error) {
	for _, t := range tagKeys {
		var scratch [4]byte
		binary.LittleEndian.PutUint32(scratch[:], t.ID)
		if err = tagKeyBucket.Put([]byte(t.Key), scratch[:]); err != nil {
			return err
		}
	}
	return
}

// closeDB closes the bbolt.DB
func closeDB(db *bbolt.DB) error {
	return db.Close()
}

// setSequence sets the bucket's sequence
func setSequence(bucket *bbolt.Bucket, seq uint64) error {
	return bucket.SetSequence(seq)
}

// createBucket creates the bucket with name
func createBucket(parentBucket *bbolt.Bucket, name []byte) (*bbolt.Bucket, error) {
	return parentBucket.CreateBucket(name)
}
