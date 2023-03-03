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

	multierror "github.com/hashicorp/go-multierror"

	"github.com/lindb/lindb/config"
	"github.com/lindb/lindb/constants"
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/fileutil"
	"github.com/lindb/lindb/pkg/logger"
	"github.com/lindb/lindb/pkg/unique"
	"github.com/lindb/lindb/series/field"
	"github.com/lindb/lindb/series/metric"
	"github.com/lindb/lindb/series/tag"
)

//go:generate mockgen -source ./metadata_backend.go -destination=./metadata_backend_mock.go -package=metadb

const (
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

// sequenceItem represents sequence metadata.
type sequenceItem struct {
	sequence unique.Sequence
	store    unique.IDStore
	key      []byte
}

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
	saveTagKey(metricID metric.ID, tagKey string) (tag.KeyID, error)
	// getAllTagKeys returns the all tag keys by metric id,
	// if not exist return empty.
	getAllTagKeys(metricID metric.ID) (tags tag.Metas, err error)
	// saveField saves the field meta for given metric id.
	saveField(metricID metric.ID, field field.Meta) error
	// getAllFields returns the  all fields by metric id,
	// if not exist return empty.
	getAllFields(metricID metric.ID) (fields field.Metas, max field.ID, err error)

	// getOrCreateMetricMetadata creates metric metadata if not exist, else load metric metadata from backend storage.
	getOrCreateMetricMetadata(namespace, metricName string, limits *models.Limits) (MetricMetadata, error)
	// getMetricMetadata gets the metric metadata include all fields/tags by metric id,
	// if not exist constants.ErrMetricBucketNotFound
	getMetricMetadata(metricID metric.ID) (metadata MetricMetadata, err error)

	// sync the backend memory data into persist storage.
	sync() error
}

// metadataBackend implements the MetadataBackend interface.
type metadataBackend struct {
	namespace, metric, tagKey, field                        unique.IDStore
	namespaceIDSequence, metricIDSequence, tagKeyIDSequence unique.Sequence

	dbs       map[string]unique.IDStore
	sequences []sequenceItem
}

// newMetadataBackend creates a new metadata backend storage.
func newMetadataBackend(parent string) (MetadataBackend, error) {
	var storageDBs map[string]unique.IDStore
	var err error
	defer func() {
		if err != nil {
			metaLogger.Error("new metadata backend fail, need close backend storage")
			// if got err, need close storage db if not nil
			for name, db := range storageDBs {
				if err0 := db.Close(); err0 != nil {
					metaLogger.Error("close storage db err when create metadata backend fail",
						logger.String("db", name), logger.Error(err0))
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
	}
	// init seq function
	initSeq := func(db unique.IDStore, key []byte, init func(seq unique.Sequence)) error {
		val, exist, err0 := db.Get(key)
		if err0 != nil {
			return err0
		}
		sequenceInitValue := uint32(0)
		if exist {
			sequenceInitValue = binary.LittleEndian.Uint32(val)
		}
		cacheSize := config.GlobalStorageConfig().TSDB.MetaSequenceCache

		// persist cached sequence value
		err = unique.SaveSequence(db, key, sequenceInitValue+cacheSize)
		if err != nil {
			return err
		}

		sequence := unique.NewSequence(sequenceInitValue, cacheSize)
		init(sequence)
		// cache all sequences
		backend.sequences = append(backend.sequences,
			sequenceItem{
				sequence: sequence,
				store:    db,
				key:      key,
			})
		return nil
	}
	var sequences = []struct {
		key  []byte
		db   unique.IDStore
		init func(seq unique.Sequence)
	}{
		{
			key: namespaceIDSequenceKey,
			db:  backend.namespace,
			init: func(seq unique.Sequence) {
				backend.namespaceIDSequence = seq
			},
		},
		{
			key: metricIDSequenceKey,
			db:  backend.metric,
			init: func(seq unique.Sequence) {
				backend.metricIDSequence = seq
			},
		},
		{
			key: tagKeyIDSequenceKey,
			db:  backend.tagKey,
			init: func(seq unique.Sequence) {
				backend.tagKeyIDSequence = seq
			},
		},
	}

	// init sequence with value
	for _, arg := range sequences {
		err = initSeq(arg.db, arg.key, arg.init)
		if err != nil {
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
		db, err := newIDStoreFn(path.Join(parent, name))
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
		if bytes.Equal(val, namespaceIDSequenceKey) {
			continue
		}
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
	nsLen := len(namespaceVal)
	for _, val := range values {
		metricNames = append(metricNames, string(val[nsLen:]))
	}
	return
}

// getMetricID gets the metric id by namespace and metric name,
// if not exist return constants.ErrMetricIDNotFound.
func (mb *metadataBackend) getMetricID(namespace, metricName string) (metricID metric.ID, err error) {
	// 1. get namespace id
	namespaceVal, exist, err := mb.namespace.Get([]byte(namespace))
	if err != nil {
		return
	}
	if !exist {
		err = fmt.Errorf("%w, metric: %s", constants.ErrMetricIDNotFound, metricName)
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
		err = fmt.Errorf("%w, metric: %s", constants.ErrMetricIDNotFound, metricName)
		return
	}
	metricID = metric.ID(binary.LittleEndian.Uint32(metricVal))
	// TODO too many query for one query????
	return
}

// saveTagKey saves the tag meta for given metric id.
func (mb *metadataBackend) saveTagKey(metricID metric.ID, tagKey string) (tag.KeyID, error) {
	tagKeyID, err := nextSequence(mb.tagKeyIDSequence, mb.tagKey, tagKeyIDSequenceKey, false, 0)
	if err != nil {
		return tag.EmptyTagKeyID, err
	}
	id := tag.KeyID(tagKeyID)
	tagMeta := &tag.Meta{Key: tagKey, ID: id}

	val, err := tagMeta.MarshalBinary()
	if err != nil {
		return tag.EmptyTagKeyID, err
	}
	var scratch [4]byte
	binary.LittleEndian.PutUint32(scratch[:], uint32(metricID))

	if err := mb.tagKey.Merge(scratch[:], val); err != nil {
		return tag.EmptyTagKeyID, err
	}
	return id, nil
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
func (mb *metadataBackend) saveField(metricID metric.ID, f field.Meta) error {
	val, err := f.MarshalBinary()
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
func (mb *metadataBackend) getOrCreateMetricMetadata(namespace, metricName string, limits *models.Limits) (MetricMetadata, error) {
	nsKey := []byte(namespace)
	nsIDVal, exist, err := mb.namespace.Get(nsKey)
	if err != nil {
		return nil, err
	}
	if !exist {
		// gen namespace id
		var nsID uint32
		nsID, err = nextSequence(mb.namespaceIDSequence, mb.namespace, namespaceIDSequenceKey,
			limits.EnableNamespacesCheck(), limits.MaxNamespaces)
		if err != nil {
			return nil, err
		}
		var scratch [4]byte
		binary.LittleEndian.PutUint32(scratch[:], nsID)
		nsIDVal = scratch[:]
		err = mb.namespace.Put(nsKey, nsIDVal)
		if err != nil {
			return nil, err
		}
	}

	var key []byte
	key = append(key, nsIDVal...)
	key = append(key, metricName...)
	metricIDVal, exist, err := mb.metric.Get(key)
	if err != nil {
		return nil, err
	}
	if !exist {
		// gen metric id
		var metricID uint32
		metricID, err = nextSequence(mb.metricIDSequence, mb.metric, metricIDSequenceKey,
			limits.EnableMetricsCheck(), limits.MaxMetrics)
		if err != nil {
			return nil, err
		}
		var scratch [4]byte
		binary.LittleEndian.PutUint32(scratch[:], metricID)
		metricIDVal = scratch[:]
		err = mb.metric.Put(key, metricIDVal)
		if err != nil {
			return nil, err
		}
		return newMetricMetadata(metric.ID(metricID)), nil
	}

	// if metric exist, need load metric from storage
	metricID := metric.ID(binary.LittleEndian.Uint32(metricIDVal))
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
	if err := mb.saveSequences(); err != nil {
		result = multierror.Append(result, err)
	}

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
	if err := mb.sync(); err != nil {
		result = multierror.Append(result, err)
	}

	for _, db := range mb.dbs {
		if err := db.Close(); err != nil {
			result = multierror.Append(result, err)
		}
	}
	return result
}

// saveSequences persists the current value of sequence for all sequences.
func (mb *metadataBackend) saveSequences() error {
	var result error
	for _, sequence := range mb.sequences {
		current := sequence.sequence.Current()
		if err := unique.SaveSequence(sequence.store, sequence.key, current); err != nil {
			result = multierror.Append(result, err)
		}
	}
	return result
}

// nextSequence returns next value from sequence,
// if no data in cache, need to cache next back from storage.
func nextSequence(seq unique.Sequence, store unique.IDStore, key []byte, enable bool, limit uint32) (uint32, error) {
	if !seq.HasNext() {
		cur := seq.Current()
		if enable && cur > limit {
			return 0, constants.ErrTooManyMetadata
		}
		nextBatchSeriesSeq := cur + config.GlobalStorageConfig().TSDB.SeriesSequenceCache
		if err := unique.SaveSequence(store, key, nextBatchSeriesSeq); err != nil {
			return 0, err
		}
		seq.Limit(seq.Current() + config.GlobalStorageConfig().TSDB.SeriesSequenceCache)
	}
	return seq.Next(), nil
}
