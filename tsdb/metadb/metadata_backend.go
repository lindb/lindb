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
	"encoding/binary"
	"fmt"
	"io"
	"path"

	"github.com/hashicorp/go-multierror"
	"go.uber.org/atomic"

	"github.com/lindb/lindb/constants"
	"github.com/lindb/lindb/pkg/fileutil"
	"github.com/lindb/lindb/pkg/logger"
	"github.com/lindb/lindb/pkg/unique"
	"github.com/lindb/lindb/series/field"
	"github.com/lindb/lindb/series/metric"
	"github.com/lindb/lindb/series/tag"
)

//go:generate mockgen -source ./metadata_backend.go -destination=./metadata_backend_mock.go -package=metadb

const (
	MetaDB      = "meta"
	NamespaceDB = "namespace"
	MetricDB    = "metric"
	TagKeyDB    = "tagkey"
	FieldDB     = "field"
)

// for testing
var (
	mkDirFn      = fileutil.MkDirIfNotExist
	newIDStoreFn = unique.NewIDStore
)

var (
	namespaceIDSequenceKey = []byte("__$$ns_seq$$__")
	metricIDSequenceKey    = []byte("__$$metric_seq$$__")
	tagKeyIDSequenceKey    = []byte("__$$key_key_seq$$__")

	storageDBNames = []string{NamespaceDB, MetricDB, TagKeyDB, FieldDB}
)

// MetadataBackend represents the metadata backend storage.
type MetadataBackend interface {
	io.Closer

	// suggestNamespace suggests the namespace by namespace's prefix.
	suggestNamespace(prefix string, limit int) (namespaces []string, err error)
	// suggestMetricName suggests the metric name by namespace and name's prefix.
	suggestMetricName(namespace, prefix string, limit int) (metricNames []string, err error)
	// getMetricID gets the metric id by namespace and metric name,
	// if not exist return constants.ErrMetricIDNotFound.
	getMetricID(namespace string, metricName string) (metricID metric.ID, err error)
	// saveTagKey saves the tag meta for given metric id.
	saveTagKey(metricID metric.ID, tagKey string) (uint32, error)
	// getAllTagKeys returns the all tag keys by metric id,
	// if not exist return empty.
	getAllTagKeys(metricID metric.ID) (tags tag.Metas, err error)
	// saveField saves the field meta for given metric id.
	saveField(metricID metric.ID, field field.Meta) error
	// getAllFields returns the  all fields by metric id,
	// if not exist return empty.
	getAllFields(metricID metric.ID) (fields field.Metas, max field.ID, err error)

	// getOrCreateMetricMetadata creates metric metadata if not exist, else load metric metadata from backend storage.
	getOrCreateMetricMetadata(namespace, metricName string) (MetricMetadata, error)
	// getMetricMetadata gets the metric metadata include all fields/tags by metric id,
	// if not exist constants.ErrMetricBucketNotFound
	getMetricMetadata(metricID metric.ID) (metadata MetricMetadata, err error)

	// sync the backend memory data into persist storage.
	sync() error
}

// metadataBackend implements the MetadataBackend interface.
type metadataBackend struct {
	namespace, metric, tagKey, field                        unique.IDStore
	namespaceIDSequence, metricIDSequence, tagKeyIDSequence *atomic.Uint32

	dbs map[string]unique.IDStore
}

// newMetadataBackend creates a new metadata backend storage
func newMetadataBackend(parent string) (MetadataBackend, error) {
	var storageDBs map[string]unique.IDStore
	var err error
	defer func() {
		if err != nil {
			metaLogger.Error("new metadata backend fail, need close backend storage")
			// if got err, need close storage db if not nil
			for name, db := range storageDBs {
				if err := db.Close(); err != nil {
					metaLogger.Error("close storage db err when create metadata backend fail",
						logger.String("db", name), logger.Error(err))
				}
			}
		}
	}()
	storageDBs, err = newStorageDB(parent)
	if err != nil {
		return nil, err
	}
	backend := &metadataBackend{
		namespace: storageDBs[NamespaceDB],
		metric:    storageDBs[MetricDB],
		tagKey:    storageDBs[TagKeyDB],
		field:     storageDBs[FieldDB],

		dbs: storageDBs,

		namespaceIDSequence: atomic.NewUint32(0),
		metricIDSequence:    atomic.NewUint32(0),
		tagKeyIDSequence:    atomic.NewUint32(0),
	}
	// init seq function
	initSeq := func(db unique.IDStore, key []byte, seq *atomic.Uint32) error {
		val, exist, err := db.Get(key)
		if err != nil {
			return err
		}
		if exist {
			seq.Store(binary.LittleEndian.Uint32(val))
		}
		return nil
	}
	var sequences = []struct {
		key []byte
		db  unique.IDStore
		seq *atomic.Uint32
	}{
		{
			key: namespaceIDSequenceKey,
			db:  backend.namespace,
			seq: backend.namespaceIDSequence,
		},
		{
			key: metricIDSequenceKey,
			db:  backend.metric,
			seq: backend.metricIDSequence,
		},
		{
			key: tagKeyIDSequenceKey,
			db:  backend.tagKey,
			seq: backend.tagKeyIDSequence,
		},
	}

	for _, arg := range sequences {
		if err = initSeq(arg.db, arg.key, arg.seq); err != nil {
			return nil, err
		}
	}

	return backend, err
}

// newStorageDB creates backend id store.
func newStorageDB(parent string) (map[string]unique.IDStore, error) {
	if err := mkDirFn(parent); err != nil {
		return nil, err
	}
	dbs := make(map[string]unique.IDStore)
	for _, name := range storageDBNames {
		db, err := newIDStoreFn(path.Join(parent, MetaDB, name))
		if err != nil {
			return dbs, err
		}
		dbs[name] = db
	}
	return dbs, nil
}

// suggestNamespace suggests the namespace by namespace's prefix.
func (mb *metadataBackend) suggestNamespace(prefix string, limit int) (namespaces []string, err error) {
	values, err := mb.namespace.IterKeys([]byte(prefix), limit)
	if err != nil {
		return nil, err
	}
	for _, val := range values {
		namespaces = append(namespaces, string(val))
	}
	return
}

// suggestMetricName suggests the metric name by namespace and name's prefix.
func (mb *metadataBackend) suggestMetricName(namespace, prefix string, limit int) (metricNames []string, err error) {
	// 1. get namespace id
	namespaceVal, exist, err := mb.namespace.Get([]byte(namespace))
	if err != nil {
		return
	}
	if !exist {
		return
	}
	// 2. scan metric name by prefix
	var key []byte
	key = append(key, namespaceVal...)
	key = append(key, prefix...)
	values, err := mb.metric.IterKeys(key, limit)
	if err != nil {
		return
	}
	for _, val := range values {
		metricNames = append(metricNames, string(val))
	}
	return
}

// getMetricID gets the metric id by namespace and metric name,
// if not exist return constants.ErrMetricIDNotFound.
func (mb *metadataBackend) getMetricID(namespace string, metricName string) (metricID metric.ID, err error) {
	// 1. get namespace id
	namespaceVal, exist, err := mb.namespace.Get([]byte(namespace))
	if err != nil {
		return
	}
	if !exist {
		err = constants.ErrMetricIDNotFound
		return
	}
	// 2. get metric id by namespace id and name
	var key []byte
	key = append(key, namespaceVal...)
	key = append(key, metricName...)
	metricVal, exist, err := mb.metric.Get(key)
	if err != nil {
		return
	}
	if !exist {
		err = constants.ErrMetricIDNotFound
		return
	}
	metricID = metric.ID(binary.LittleEndian.Uint32(metricVal))
	fmt.Printf("metric: %s, %s, %d", namespace, metricName, metricID)
	return
}

// saveTagKey saves the tag meta for given metric id.
func (mb *metadataBackend) saveTagKey(metricID metric.ID, tagKey string) (uint32, error) {
	tagKeyID := mb.tagKeyIDSequence.Inc()
	tagMeta := &tag.Meta{Key: tagKey, ID: tagKeyID}

	val, err := tagMeta.MarshalBinary()
	fmt.Printf("tag key: %d, %s, %d\n", metricID, tagMeta.Key, tagMeta.ID)
	if err != nil {
		return tag.EmptyTagKeyID, err
	}
	var scratch [4]byte
	binary.LittleEndian.PutUint32(scratch[:], uint32(metricID))

	if err := mb.tagKey.Merge(scratch[:], val); err != nil {
		return tag.EmptyTagKeyID, err
	}
	return tagKeyID, nil
}

// getAllTagKeys returns the all tag keys by metric id, if not exist returns empty.
func (mb *metadataBackend) getAllTagKeys(metricID metric.ID) (tags tag.Metas, err error) {
	val, exist, err := mb.tagKey.Get(metricID.MarshalBinary())
	if err != nil {
		return nil, err
	}
	if !exist {
		return nil, nil
	}
	tags, err = tag.UnmarshalBinary(val)
	if err != nil {
		return nil, err
	}
	return
}

// saveField saves the field meta for given metric id.
func (mb *metadataBackend) saveField(metricID metric.ID, field field.Meta) error {
	val, err := field.MarshalBinary()
	if err != nil {
		return err
	}
	var scratch [4]byte
	binary.LittleEndian.PutUint32(scratch[:], uint32(metricID))
	return mb.field.Merge(scratch[:], val)
}

// getAllFields returns the  all fields by metric id, if not exist returns empty.
func (mb *metadataBackend) getAllFields(metricID metric.ID) (fields field.Metas, max field.ID, err error) {
	val, exist, err := mb.field.Get(metricID.MarshalBinary())
	if err != nil {
		return nil, 0, err
	}
	if !exist {
		return nil, 0, nil
	}
	fields, max, err = field.UnmarshalBinary(val)
	if err != nil {
		return nil, 0, err
	}
	return
}

// getOrCreateMetricMetadata creates metric metadata if not exist, else load metric metadata from backend storage.
func (mb *metadataBackend) getOrCreateMetricMetadata(namespace, metricName string) (MetricMetadata, error) {
	nsKey := []byte(namespace)
	nsVal, exist, err := mb.namespace.Get(nsKey)
	if err != nil {
		return nil, err
	}
	if !exist {
		// gen namespace id
		nsID := mb.namespaceIDSequence.Inc()
		var scratch [4]byte
		binary.LittleEndian.PutUint32(scratch[:], nsID)
		nsVal = scratch[:]
		if err := mb.namespace.Put(namespaceIDSequenceKey, nsVal); err != nil {
			return nil, err
		}
	}

	var key []byte
	key = append(key, nsVal...)
	key = append(key, metricName...)
	metricVal, exist, err := mb.metric.Get(key)
	if err != nil {
		return nil, err
	}
	if !exist {
		// gen metric id
		metricID := mb.metricIDSequence.Inc()
		var scratch [4]byte
		binary.LittleEndian.PutUint32(scratch[:], metricID)
		metricVal = scratch[:]
		if err := mb.metric.Put(metricIDSequenceKey, metricVal); err != nil {
			return nil, err
		}
		return newMetricMetadata(metric.ID(metricID)), nil
	}

	// if metric exist, need load metric from storage
	metricID := metric.ID(binary.LittleEndian.Uint32(metricVal))
	return mb.getMetricMetadata(metricID)
}

// getMetricMetadata gets the metric metadata include all fields/tags by metric id,
// if not exist return constants.ErrMetricBucketNotFound
func (mb *metadataBackend) getMetricMetadata(metricID metric.ID) (metadata MetricMetadata, err error) {
	fields, fieldMaxID, err := mb.getAllFields(metricID)
	if err != nil {
		return
	}
	tags, err := mb.getAllTagKeys(metricID)
	if err != nil {
		return
	}

	metadata = newMetricMetadata(metricID)
	// initialize fields and tags
	metadata.initialize(fields, int32(fieldMaxID), tags)
	return
}

// sync the backend memory data into persist storage.
func (mb *metadataBackend) sync() error {
	var result error
	for _, db := range mb.dbs {
		if err := db.Flush(); err != nil {
			result = multierror.Append(result, err)
		}
	}
	return result
}

// Close closes the backend storage.
func (mb *metadataBackend) Close() error {
	var result error
	for _, db := range mb.dbs {
		if err := db.Close(); err != nil {
			result = multierror.Append(result, err)
		}
	}
	return result
}
