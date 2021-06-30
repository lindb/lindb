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

package metadb

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"path"
	"time"

	"go.etcd.io/bbolt"
	"go.uber.org/atomic"

	"github.com/lindb/lindb/constants"
	"github.com/lindb/lindb/pkg/fileutil"
	"github.com/lindb/lindb/pkg/logger"
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
	// rollbackMetricID rollbacks metric id
	rollbackMetricID(metricID uint32)
	// rollbackTagKeyID rollbacks tag key id
	rollbackTagKeyID(tagKeyID uint32)

	// loadMetricMetadata loads the metric metadata include all fields/tags by namespace and metric name,
	// if not exist return constants.ErrMetricIDNotFound, constants.ErrMetricBucketNotFound, constants.ErrMetricIDNotFound
	loadMetricMetadata(namespace, metricName string) (MetricMetadata, error)
	// getMetricMetadata gets the metric metadata include all fields/tags by metric id,
	// if not exist constants.ErrMetricBucketNotFound
	getMetricMetadata(metricID uint32) (metadata MetricMetadata, err error)

	// getMetricID gets the metric id by namespace and metric name,
	// if not exist return constants.ErrMetricIDNotFound
	getMetricID(namespace string, metricName string) (metricID uint32, err error)
	// getTagKeyID gets the tag key id by metric id and tag key key,
	// if not exist return constants.ErrTagKeyIDNotFound
	getTagKeyID(metricID uint32, tagKey string) (tagKeyID uint32, err error)
	// getAllTagKeys returns the all tag keys by metric id,
	// if not exist return constants.ErrMetricBucketNotFound
	getAllTagKeys(metricID uint32) (tags []tag.Meta, err error)
	// getField gets the field meta by metric id and field name,
	// if not exist return constants.ErrMetricBucketNotFound, constants.ErrFieldBucketNotFound
	getField(metricID uint32, fieldName field.Name) (f field.Meta, err error)
	// getAllFields returns the  all fields by metric id,
	// if not exist return constants.ErrMetricBucketNotFound
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
			metaLogger.Error("close bbolt.db err when create metadata backend fail",
				logger.String("db", parent), logger.Error(e))
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

// rollbackMetricID rollbacks metric id
func (mb *metadataBackend) rollbackMetricID(metricID uint32) {
	if metricID == mb.metricIDSequence.Load() {
		mb.metricIDSequence.Dec() // recycle metric id
	}
}

// rollbackMetricID rollbacks metric id
func (mb *metadataBackend) rollbackTagKeyID(tagKeyID uint32) {
	if tagKeyID == mb.tagKeyIDSequence.Load() {
		mb.tagKeyIDSequence.Dec() // recycle tag key id
	}
}

// genTagKeyID generates the tag key id in the memory
func (mb *metadataBackend) genTagKeyID() uint32 {
	return mb.tagKeyIDSequence.Inc()
}

// loadMetricMetadata loads the metric metadata include all fields/tags by namespace and metric name,
// if not exist return constants.ErrMetricIDNotFound, constants.ErrMetricBucketNotFound, constants.ErrMetricIDNotFound
func (mb *metadataBackend) loadMetricMetadata(namespace, metricName string) (metadata MetricMetadata, err error) {
	metricID, err := mb.getMetricID(namespace, metricName)
	if err != nil {
		return nil, err
	}
	return mb.getMetricMetadata(metricID)
}

// getMetricMetadata gets the metric metadata include all fields/tags by metric id,
// if not exist return constants.ErrMetricBucketNotFound
func (mb *metadataBackend) getMetricMetadata(metricID uint32) (metadata MetricMetadata, err error) {
	var scratch [4]byte
	var fieldIDSeq int32
	var tags []tag.Meta
	var fields []field.Meta
	binary.LittleEndian.PutUint32(scratch[:], metricID)
	err = mb.db.View(func(tx *bbolt.Tx) error {
		metricBucket := tx.Bucket(metricBucketName).Bucket(scratch[:])
		if metricBucket == nil {
			return fmt.Errorf("%w, metricID:%d", constants.ErrMetricBucketNotFound, metricID)
		}
		tags = loadTagKeys(metricBucket.Bucket(tagBucketName))
		fBucket := metricBucket.Bucket(fieldBucketName)
		fieldIDSeq = int32(fBucket.Sequence())
		fields = loadFields(fBucket)
		return nil
	})
	if err != nil {
		return metadata, fmt.Errorf("%w, metricID:%d with error: %s",
			constants.ErrMetricBucketNotFound, metricID, err)
	}

	metadata = newMetricMetadata(metricID, fieldIDSeq)
	// initialize fields and tags
	metadata.initialize(fields, tags)
	return
}

// getMetricID gets the metric id by namespace and metric name,
// if not exist return constants.ErrNameSpaceBucketNotFound, constants.ErrMetricIDNotFound
func (mb *metadataBackend) getMetricID(namespace string, metricName string) (metricID uint32, err error) {
	err = mb.db.View(func(tx *bbolt.Tx) error {
		nsBucket := tx.Bucket(nsBucketName).Bucket([]byte(namespace))
		if nsBucket == nil {
			return fmt.Errorf("%w, namepsace: %s, metricName: %s",
				constants.ErrNameSpaceBucketNotFound, namespace, metricName)
		}
		value := nsBucket.Get([]byte(metricName))
		if len(value) == 0 {
			return fmt.Errorf("%w, namepsace: %s, metricName: %s",
				constants.ErrMetricIDNotFound, namespace, metricName)
		}
		metricID = binary.LittleEndian.Uint32(value)
		return nil
	})
	return
}

// getTagKeyID gets the tag key id by metric id and tag key key, if not exist return constants.ErrTagKeyIDNotFound
func (mb *metadataBackend) getTagKeyID(metricID uint32, tagKey string) (tagKeyID uint32, err error) {
	var scratch [4]byte
	binary.LittleEndian.PutUint32(scratch[:], metricID)
	err = mb.db.View(func(tx *bbolt.Tx) error {
		metricBucket := tx.Bucket(metricBucketName).Bucket(scratch[:])
		if metricBucket == nil {
			return fmt.Errorf("%w, tagKey: %s", constants.ErrTagKeyIDNotFound, tagKey)
		}
		value := metricBucket.Bucket(tagBucketName).Get([]byte(tagKey))
		if len(value) == 0 {
			return fmt.Errorf("%w, tagKey: %s not in bucket", constants.ErrTagKeyIDNotFound, tagKey)
		}
		tagKeyID = binary.LittleEndian.Uint32(value)
		return nil
	})
	return
}

// getAllTagKeys returns the all tag keys by metric id, if not exist return constants.ErrMetricBucketNotFound
func (mb *metadataBackend) getAllTagKeys(metricID uint32) (tags []tag.Meta, err error) {
	var scratch [4]byte
	binary.LittleEndian.PutUint32(scratch[:], metricID)
	err = mb.db.View(func(tx *bbolt.Tx) error {
		metricBucket := tx.Bucket(metricBucketName).Bucket(scratch[:])
		if metricBucket != nil {
			tags = loadTagKeys(metricBucket.Bucket(tagBucketName))
			return nil
		}
		return fmt.Errorf("%w, metricID: %d", constants.ErrMetricBucketNotFound, metricID)
	})
	return
}

// getField gets the field meta by metric id and field name,
// if not exist return constants.ErrMetricBucketNotFound, constants.ErrFieldBucketNotFound
func (mb *metadataBackend) getField(metricID uint32, fieldName field.Name) (f field.Meta, err error) {
	var scratch [4]byte
	binary.LittleEndian.PutUint32(scratch[:], metricID)
	err = mb.db.View(func(tx *bbolt.Tx) error {
		metricBucket := tx.Bucket(metricBucketName).Bucket(scratch[:])
		if metricBucket == nil {
			return fmt.Errorf("%w during getField, metricID: %d", constants.ErrMetricBucketNotFound, metricID)
		}
		value := metricBucket.Bucket(fieldBucketName).Get([]byte(fieldName))
		if len(value) == 0 {
			return fmt.Errorf("%w during getField, fieldName: %s", constants.ErrFieldBucketNotFound, fieldName)
		}
		f.Name = fieldName
		f.ID = field.ID(value[0])
		f.Type = field.Type(value[1])
		return nil
	})
	return
}

// getAllFields returns the  all fields by metric id, if not exist return constants.ErrMetricBucketNotFound
func (mb *metadataBackend) getAllFields(metricID uint32) (fields []field.Meta, err error) {
	var scratch [4]byte
	binary.LittleEndian.PutUint32(scratch[:], metricID)
	err = mb.db.View(func(tx *bbolt.Tx) error {
		metricBucket := tx.Bucket(metricBucketName).Bucket(scratch[:])
		if metricBucket == nil {
			return fmt.Errorf("%w during getAllFields metricID: %d", constants.ErrMetricBucketNotFound, metricID)
		}
		fields = loadFields(metricBucket.Bucket(fieldBucketName))
		return nil
	})
	return
}

// saveMetadata saves the pending metadata include namespace/metric metadata
func (mb *metadataBackend) saveMetadata(event *metadataUpdateEvent) (err error) {
	err = mb.db.Update(func(tx *bbolt.Tx) error {
		if err := mb.saveNamespaceAndMetric(tx.Bucket(nsBucketName), event); err != nil {
			return err
		}
		if err := mb.saveMetricMetadata(tx.Bucket(metricBucketName), event); err != nil {
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
			Name: field.Name(k),
			ID:   field.ID(v[0]),
			Type: field.Type(v[1]),
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
func (mb *metadataBackend) saveNamespaceAndMetric(nsRootBucket *bbolt.Bucket, event *metadataUpdateEvent) (err error) {
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
		if mb.metricIDSequence.Load() < event.metricSeqID {
			mb.metricIDSequence.Store(event.metricSeqID)
		}
	}
	return
}

// saveMetricMetadata saves metric metadata include fields/tag keys if exist with metric root bucket
func (mb *metadataBackend) saveMetricMetadata(metricRootBucket *bbolt.Bucket, event *metadataUpdateEvent) (err error) {
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
		if mb.tagKeyIDSequence.Load() < event.tagKeySeqID {
			mb.tagKeyIDSequence.Store(event.tagKeySeqID)
		}
	}
	return nil
}

// saveFields saves fields for metric with field bucket
func saveFields(fieldBucket *bbolt.Bucket, fieldIDSeq uint16, fields []field.Meta) (err error) {
	for _, f := range fields {
		var fieldValue [2]byte
		fieldValue[0] = byte(f.ID)
		fieldValue[1] = byte(f.Type)
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
	if bucket.Sequence() < seq {
		return bucket.SetSequence(seq)
	}
	return nil
}

// createBucket creates the bucket with name
func createBucket(parentBucket *bbolt.Bucket, name []byte) (*bbolt.Bucket, error) {
	return parentBucket.CreateBucket(name)
}
