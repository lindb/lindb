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

package index

import (
	"context"
	"fmt"
	"math"
	"path"
	"time"

	commonfileutil "github.com/lindb/common/pkg/fileutil"
	"github.com/lindb/common/pkg/logger"
	"github.com/lindb/roaring"
	"go.uber.org/atomic"

	"github.com/lindb/lindb/constants"
	v1 "github.com/lindb/lindb/index/v1"
	"github.com/lindb/lindb/kv"
	"github.com/lindb/lindb/pkg/strutil"
	"github.com/lindb/lindb/series/field"
	"github.com/lindb/lindb/series/metric"
	"github.com/lindb/lindb/series/tag"
	"github.com/lindb/lindb/sql/stmt"
)

// for testing
var (
	mkdir       = commonfileutil.MkDirIfNotExist
	newSequence = NewSequence
)

const (
	metaPath           = "kv"
	sequencePath       = "sequence"
	nsFamilyName       = "ns"
	metricFamilyName   = "metric"
	schemaFamilyName   = "schema"
	tagValueFamilyName = "tv"
)

// metricMetaDatabase implements MetricMetaDatabase interface.
type metricMetaDatabase struct {
	ctx      context.Context
	cancel   context.CancelFunc
	kvStore  kv.Store
	ns       IndexKVStore // namespace -> namespace id
	metric   IndexKVStore // metric name -> metric id
	tagValue IndexKVStore // tag key id -> tag values

	schemaStore MetricSchemaStore // metric id -> metric metadata like table schema(tag keys/fields/shard sequences etc.)

	sequence *Sequence // sequence(namespace/name/tag key)
	worker   *NotifyWorker
	flushing atomic.Bool

	logger logger.Logger
}

// NewMetricMetaDatabase creates a metric meta store.
func NewMetricMetaDatabase(dir string) (MetricMetaDatabase, error) {
	if err := mkdir(dir); err != nil {
		return nil, err
	}
	sequence, err := newSequence(path.Join(dir, sequencePath))
	if err != nil {
		return nil, err
	}
	kvStore, err := kv.GetStoreManager().CreateStore(path.Join(dir, metaPath), kv.DefaultStoreOption())
	if err != nil {
		return nil, err
	}
	nsFamily, err := kvStore.CreateFamily(nsFamilyName, kv.FamilyOption{
		Merger: string(v1.IndexKVMerger),
	})
	if err != nil {
		return nil, err
	}
	metricFamily, err := kvStore.CreateFamily(metricFamilyName, kv.FamilyOption{
		Merger: string(v1.IndexKVMerger),
	})
	if err != nil {
		return nil, err
	}
	tagValueFamily, err := kvStore.CreateFamily(tagValueFamilyName, kv.FamilyOption{
		Merger: string(v1.IndexKVMerger),
	})
	if err != nil {
		return nil, err
	}
	schemaFamily, err := kvStore.CreateFamily(schemaFamilyName, kv.FamilyOption{
		Merger: string(v1.MetricSchemaMerger),
	})
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithCancel(context.Background())
	mm := &metricMetaDatabase{
		ctx:         ctx,
		cancel:      cancel,
		kvStore:     kvStore,
		ns:          NewIndexKVStore(nsFamily, 1000, 20*time.Minute),
		metric:      NewIndexKVStore(metricFamily, 1000, 10*time.Minute),
		tagValue:    NewIndexKVStore(tagValueFamily, 10000, 10*time.Minute),
		schemaStore: NewMetricSchemaStore(schemaFamily),
		sequence:    sequence,
		logger:      logger.GetLogger("Index", "MetricMetaDatabase"),
	}

	mm.worker = NewWorker(ctx, dir, 100*time.Millisecond, mm.handle)
	return mm, nil
}

func (mm *metricMetaDatabase) Notify(n Notifier) {
	mm.worker.Notify(n)
}

func (mm *metricMetaDatabase) handle(n Notifier) {
	switch notifier := n.(type) {
	case *MetaNotifier:
		metricID, err := mm.genMetricID(strutil.String2ByteSlice(notifier.Namespace), strutil.String2ByteSlice(notifier.MetricName))
		notifier.Callback(uint32(metricID), err)
		PutMetaNotifier(notifier)
	case *FieldNotifier:
		metricID, err := mm.genMetricID(strutil.String2ByteSlice(notifier.Namespace), strutil.String2ByteSlice(notifier.MetricName))
		var fieldID field.ID
		if err == nil {
			fieldID, err = mm.genFieldID(metricID, notifier.Field)
		}
		notifier.Callback(fieldID, err)
		PutFieldNotifier(notifier)
	case *TagNotifier:
		metricID := notifier.metricID
		for _, tag := range notifier.tags {
			tagKeyID, err := mm.genTagKeyID(metricID, tag.Key)
			if err != nil {
				mm.logger.Warn("gen tag key id failure, ignore build this tag index",
					logger.String("key", string(tag.Key)), logger.String("value", string(tag.Value)), logger.Error(err))
				continue
			}
			tagValueID, err := mm.genTagValueID(tagKeyID, tag.Value)
			if err != nil {
				mm.logger.Warn("gen tag value id failure, ignore build this tag index",
					logger.String("key", string(tag.Key)), logger.String("value", string(tag.Value)), logger.Error(err))
				continue
			}
			notifier.buildIndex(uint32(tagKeyID), tagValueID)
		}
		PutTagNotifier(notifier)
	case *FlushNotifier:
		if mm.flushing.CompareAndSwap(false, true) {
			mm.PrepareFlush() // do under worker goroutine
			// start flush goroutine, do io block background
			go func() {
				notifier.Callback(mm.Flush())
			}()
		} else {
			notifier.Callback(nil)
		}
	}
}

func (mm *metricMetaDatabase) genMetricID(namespace, metricName []byte) (metric.ID, error) {
	nsID, err := mm.ns.GetOrCreateValue(uint32(namespace[0]), namespace, mm.sequence.GetNamespaceSeq)
	if err != nil {
		return 0, err
	}
	metricID, err := mm.metric.GetOrCreateValue(nsID, metricName, mm.sequence.GetMetricNameSeq)
	if err != nil {
		return 0, err
	}
	return metric.ID(metricID), nil
}

func (mm *metricMetaDatabase) genFieldID(metricID metric.ID, f field.Meta) (field.ID, error) {
	return mm.schemaStore.genFieldID(metricID, f)
}

func (mm *metricMetaDatabase) genTagKeyID(metricID metric.ID, tagKey []byte) (tag.KeyID, error) {
	return mm.schemaStore.genTagKeyID(metricID, tagKey, mm.sequence.GetTagKeySeq)
}

func (mm *metricMetaDatabase) genTagValueID(tagKeyID tag.KeyID, tagValue []byte) (uint32, error) {
	return mm.tagValue.GetOrCreateValue(uint32(tagKeyID), tagValue, mm.sequence.GetTagValueSeq)
}

func (mm *metricMetaDatabase) GetSchema(metricID metric.ID) (*metric.Schema, error) {
	schema, err := mm.schemaStore.GetSchema(metricID)
	return schema, err
}

func (mm *metricMetaDatabase) SuggestNamespace(prefix string, limit int) (namespaces []string, err error) {
	for i := 0; i < math.MaxUint8; i++ {
		rs, err := mm.ns.Suggest(uint32(i), prefix, limit)
		if err != nil {
			return nil, err
		}
		namespaces = append(namespaces, rs...)
		if len(namespaces) >= limit {
			return namespaces, nil
		}
	}
	return
}

func (mm *metricMetaDatabase) SuggestMetrics(namespace, metricPrefix string, limit int) ([]string, error) {
	ns := strutil.String2ByteSlice(namespace)
	nsID, ok, err := mm.ns.GetValue(uint32(ns[0]), ns)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, nil
	}
	return mm.metric.Suggest(nsID, metricPrefix, limit)
}

func (mm *metricMetaDatabase) SuggestTagValues(tagKeyID tag.KeyID, tagValuePrefix string, limit int) ([]string, error) {
	return mm.tagValue.Suggest(uint32(tagKeyID), tagValuePrefix, limit)
}

func (mm *metricMetaDatabase) GetMetricID(namespace, metricName string) (metric.ID, error) {
	ns := strutil.String2ByteSlice(namespace)
	nsID, ok, err := mm.ns.GetValue(uint32(ns[0]), ns)
	if err != nil {
		return 0, err
	}
	if !ok {
		return 0, fmt.Errorf("%w, metric: %s", constants.ErrMetricIDNotFound, metricName)
	}
	metricID, ok, err := mm.metric.GetValue(nsID, strutil.String2ByteSlice(metricName))
	if err != nil {
		return 0, err
	}
	if !ok {
		return 0, fmt.Errorf("%w, metric: %s", constants.ErrMetricIDNotFound, metricName)
	}
	return metric.ID(metricID), nil
}

// FindTagValueDsByExpr finds tag value ids by tag filter expr for spec tag key,
// if not exist, return nil, constants.ErrNotFound, else returns tag value ids
func (mm *metricMetaDatabase) FindTagValueDsByExpr(tagKeyID tag.KeyID, expr stmt.TagFilter) (*roaring.Bitmap, error) {
	ids, err := mm.tagValue.FindValuesByExpr(uint32(tagKeyID), expr)
	if err != nil {
		return nil, err
	}
	result := roaring.New()
	result.AddMany(ids)
	return result, nil
}

// FindTagValueIDsForTag get tag value ids for spec tag key of metric,
// if not exist, return nil, constants.ErrNotFound, else returns tag value ids
func (mm *metricMetaDatabase) FindTagValueIDsForTag(tagKeyID tag.KeyID) (tagValueIDs *roaring.Bitmap, err error) {
	ids, err := mm.tagValue.GetValues(uint32(tagKeyID))
	if err != nil {
		return nil, err
	}
	result := roaring.New()
	result.AddMany(ids)
	return result, nil
}

// CollectTagValues collects the tag values by tag value ids,
func (mm *metricMetaDatabase) CollectTagValues(
	tagKeyID tag.KeyID,
	tagValueIDs *roaring.Bitmap,
	tagValues map[uint32]string,
) error {
	return mm.tagValue.CollectKVs(uint32(tagKeyID), tagValueIDs, tagValues)
}

func (mm *metricMetaDatabase) PrepareFlush() {
	mm.ns.PrepareFlush()
	mm.metric.PrepareFlush()
	mm.tagValue.PrepareFlush()
	mm.schemaStore.PrepareFlush()
}

func (mm *metricMetaDatabase) Flush() error {
	defer func() {
		mm.flushing.Store(false)
	}()
	if err := mm.sequence.Sync(); err != nil {
		return err
	}
	if err := mm.ns.Flush(); err != nil {
		return err
	}
	if err := mm.metric.Flush(); err != nil {
		return err
	}
	if err := mm.schemaStore.Flush(); err != nil {
		return err
	}
	if err := mm.tagValue.Flush(); err != nil {
		return err
	}
	return nil
}

// Close closes metric meta database.
func (mm *metricMetaDatabase) Close() error {
	mm.cancel()
	mm.worker.Shutdown()
	if err := mm.sequence.Close(); err != nil {
		return err
	}
	return kv.GetStoreManager().CloseStore(mm.kvStore.Name())
}
