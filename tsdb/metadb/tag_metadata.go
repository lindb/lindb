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
	"errors"
	"strings"
	"sync"

	"github.com/lindb/roaring"
	"github.com/prometheus/client_golang/prometheus"

	"github.com/lindb/lindb/constants"
	"github.com/lindb/lindb/kv"
	"github.com/lindb/lindb/monitoring"
	"github.com/lindb/lindb/pkg/strutil"
	"github.com/lindb/lindb/sql/stmt"
	"github.com/lindb/lindb/tsdb/tblstore/tagkeymeta"
)

//go:generate mockgen -source ./tag_metadata.go -destination=./tag_metadata_mock.go -package metadb

// for testing
var (
	newTagReaderFunc  = tagkeymeta.NewReader
	newTagFlusherFunc = tagkeymeta.NewFlusher
)

var (
	genTagValueIDCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "meta_gen_tag_value_id",
			Help: "Generate tag value id counter.",
		},
		[]string{"db"},
	)
)

func init() {
	monitoring.StorageRegistry.MustRegister(genTagValueIDCounter)
}

// TagMetadata represents the tag metadata, stores all tag values under spec tag key
type TagMetadata interface {
	// GenTagValueID generates the tag value id for spec tag key
	GenTagValueID(tagKeyID uint32, tagValue string) (uint32, error)
	// SuggestTagValues returns suggestions from given tag key id and prefix of tag value
	SuggestTagValues(tagKeyID uint32, tagValuePrefix string, limit int) []string
	// FindTagValueDsByExpr finds tag value ids by tag filter expr for spec tag key,
	// if not exist, return nil, constants.ErrNotFound, else returns tag value ids
	FindTagValueDsByExpr(tagKeyID uint32, expr stmt.TagFilter) (*roaring.Bitmap, error)
	// GetTagValueIDsForTag get tag value ids for spec metric's tag key,
	// if not exist, return nil, constants.ErrNotFound, else returns tag value ids
	GetTagValueIDsForTag(tagKeyID uint32) (*roaring.Bitmap, error)
	// CollectTagValues collects the tag values by tag value ids,
	CollectTagValues(tagKeyID uint32,
		tagValueIDs *roaring.Bitmap,
		tagValues map[uint32]string,
	) error
	// Flush flushes the memory tag metadata into kv store
	Flush() error
}

// tagMetadata implements TagMetadata interface
type tagMetadata struct {
	databaseName string
	family       kv.Family // store tag key/value data using common kv store
	mutable      *TagStore // mutable store current writeable memory store
	immutable    *TagStore // immutable need to flush into kv store

	rwMutex sync.RWMutex
}

// NewTagMetadata creates a tag metadata
func NewTagMetadata(databaseName string, family kv.Family) TagMetadata {
	m := &tagMetadata{
		databaseName: databaseName,
		family:       family,
		mutable:      NewTagStore(),
	}
	return m
}

// GenTagValueID generates the tag value id for spec tag key
func (m *tagMetadata) GenTagValueID(tagKeyID uint32, tagValue string) (tagValueID uint32, err error) {
	// get tag value id from memory with read lock
	m.rwMutex.RLock()
	tagValueID, ok := m.getTagValueIDInMem(tagKeyID, tagValue)
	if ok {
		m.rwMutex.RUnlock()
		return tagValueID, nil
	}
	m.rwMutex.RUnlock()

	// try load tag value id from kv store
	snapshot := m.family.GetSnapshot()
	defer snapshot.Close()

	readers, err := snapshot.FindReaders(tagKeyID)
	if err != nil {
		// find table.Reader err, return it
		return
	}
	var reader tagkeymeta.Reader
	if len(readers) > 0 {
		// found tag data in kv store, try load tag value data
		reader = newTagReaderFunc(readers)
		tagValueID, err = reader.GetTagValueID(tagKeyID, tagValue)
		if err == nil {
			// got tag value id from kv store
			return tagValueID, nil
		}
		if !errors.Is(err, constants.ErrNotFound) {
			// if load tag value id err, return it
			return
		}
	}

	// tag value not exist, need assign new tag value id for tag value with write lock
	m.rwMutex.Lock()
	defer m.rwMutex.Unlock()
	// double check, memory if exist tag value
	tagValueID, ok = m.getTagValueIDInMem(tagKeyID, tagValue)
	if ok {
		return tagValueID, nil
	}

	// assign new tag value id
	tag, ok := m.mutable.Get(tagKeyID)
	if !ok {
		if reader != nil {
			// if tag data exist in kv store, need load tag value id auto sequence
			seq, err := reader.GetTagValueSeq(tagKeyID)
			if err != nil {
				return 0, err
			}
			tag = newTagEntry(seq)
		} else {
			// for new tag, auto sequence start with 1
			tag = newTagEntry(0)
		}
		// cache tag entry
		m.mutable.Put(tagKeyID, tag)
	}

	// assign new id
	tagValueID = tag.genTagValueID()
	tag.addTagValue(tagValue, tagValueID)
	//TODO add wal???

	genTagValueIDCounter.WithLabelValues(m.databaseName).Inc()

	return tagValueID, nil
}

// SuggestTagValues returns suggestions from given tag key id and prefix of tag value
func (m *tagMetadata) SuggestTagValues(tagKeyID uint32, tagValuePrefix string, limit int) []string {
	result := make([]string, 0)
	m.loadTagValueIDsInMem(tagKeyID, func(tagEntry TagEntry) {
		for value := range tagEntry.getTagValues() {
			if strings.HasPrefix(value, tagValuePrefix) {
				result = append(result, value)
			}
		}
	})

	snapshot := m.family.GetSnapshot()
	defer snapshot.Close()

	readers, err := snapshot.FindReaders(tagKeyID)
	if err != nil {
		// find table.Reader err, return it
		return nil
	}
	var reader tagkeymeta.Reader
	if len(readers) > 0 {
		// found tag data in kv store, try load tag value data
		reader = newTagReaderFunc(readers)
		readerValues := reader.SuggestTagValues(tagKeyID, tagValuePrefix, limit)
		result = append(result, readerValues...)
	}
	return result
}

// FindTagValueDsByExpr finds tag value ids by tag filter expr for spec tag key,
// if not exist, return nil, constants.ErrNotFound, else returns tag value ids
func (m *tagMetadata) FindTagValueDsByExpr(tagKeyID uint32, expr stmt.TagFilter) (*roaring.Bitmap, error) {
	result := roaring.New()
	m.loadTagValueIDsInMem(tagKeyID, func(tagEntry TagEntry) {
		ids := tagEntry.findSeriesIDsByExpr(expr)
		if ids != nil {
			result.Or(ids)
		}
	})

	err := m.loadTagValueIDsInKV(tagKeyID, func(reader tagkeymeta.Reader) error {
		tagValueIDs, err := reader.FindValueIDsByExprForTagKeyID(tagKeyID, expr)
		if err != nil {
			return err
		}
		result.Or(tagValueIDs)
		return nil
	})
	if err != nil {
		return nil, err
	}
	return result, nil
}

// GetTagValueIDsForTag get tag value ids for spec metric's tag key,
// if not exist, return nil, constants.ErrNotFound, else returns tag value ids
func (m *tagMetadata) GetTagValueIDsForTag(tagKeyID uint32) (*roaring.Bitmap, error) {
	result := roaring.New()
	m.loadTagValueIDsInMem(tagKeyID, func(tagEntry TagEntry) {
		ids := tagEntry.getTagValueIDs()
		if ids != nil {
			result.Or(ids)
		}
	})

	err := m.loadTagValueIDsInKV(tagKeyID, func(reader tagkeymeta.Reader) error {
		tagValueIDs, err := reader.GetTagValueIDsForTagKeyID(tagKeyID)
		if err != nil {
			return err
		}
		result.Or(tagValueIDs)
		return nil
	})
	if err != nil {
		return nil, err
	}
	return result, nil
}

// CollectTagValues collects the tag values by tag value ids for spec tag key,
func (m *tagMetadata) CollectTagValues(tagKeyID uint32,
	tagValueIDs *roaring.Bitmap,
	tagValues map[uint32]string,
) error {
	m.loadTagValueIDsInMem(tagKeyID, func(tagEntry TagEntry) {
		if tagValueIDs.IsEmpty() {
			// if found all tag value in memory, return it(maybe load immutable memory)
			return
		}
		tagEntry.collectTagValues(tagValueIDs, tagValues)
	})

	if tagValueIDs.IsEmpty() {
		// no need collect tag value ids, returns it
		return nil
	}
	err := m.loadTagValueIDsInKV(tagKeyID, func(reader tagkeymeta.Reader) error {
		return reader.CollectTagValues(tagKeyID, tagValueIDs, tagValues)
	})
	if err != nil {
		return err
	}
	return nil
}

// Flush flushes the memory tag metadata into kv store
func (m *tagMetadata) Flush() error {
	if !m.checkFlush() {
		return nil
	}

	// flush immutable data into kv store
	fluster := m.family.NewFlusher()
	tagFluster := newTagFlusherFunc(fluster)
	if err := m.immutable.WalkEntry(func(key uint32, value TagEntry) error {
		tagValues := value.getTagValues()
		for tagValue, tagValueID := range tagValues {
			tagFluster.FlushTagValue(strutil.String2ByteSlice(tagValue), tagValueID)
		}
		if err := tagFluster.FlushTagKeyID(key, value.getTagValueIDSeq()); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return err
	}
	if err := tagFluster.Commit(); err != nil {
		return err
	}
	// finally clear immutable
	m.rwMutex.Lock()
	m.immutable = nil
	m.rwMutex.Unlock()
	return nil
}

// checkFlush checks if need do flush job, if need, do switch mutable/immutable
func (m *tagMetadata) checkFlush() bool {
	m.rwMutex.Lock()
	defer m.rwMutex.Unlock()

	if m.mutable.Size() == 0 && m.immutable == nil {
		// no new data or immutable is not nil
		return false
	}
	if m.mutable.Size() > 0 && m.immutable == nil {
		// reset mutable, if flush fail immutable is not nil
		m.immutable = m.mutable
		m.mutable = NewTagStore()
	}
	return true
}

// loadTagValueIDsInKV loads tag value ids in kv store
func (m *tagMetadata) loadTagValueIDsInKV(tagKeyID uint32, fn func(reader tagkeymeta.Reader) error) error {
	// try load tag value id from kv store
	snapshot := m.family.GetSnapshot()
	defer snapshot.Close()

	readers, err := snapshot.FindReaders(tagKeyID)
	if err != nil {
		// find table.Reader err, return it
		return err
	}
	var reader tagkeymeta.Reader
	if len(readers) > 0 {
		// found tag data in kv store, try load tag value data
		reader = newTagReaderFunc(readers)
		if err := fn(reader); err != nil {
			return err
		}
	}
	return nil
}

// loadTagValueIDsInMem loads tag value ids from mutable/immutable store
func (m *tagMetadata) loadTagValueIDsInMem(tagKeyID uint32, fn func(tagEntry TagEntry)) {
	// define get tag value ids func
	getTagValueIDs := func(tagStore *TagStore) {
		tag, ok := tagStore.Get(tagKeyID)
		if ok {
			fn(tag)
		}
	}

	m.rwMutex.RLock()
	defer m.rwMutex.RUnlock()

	getTagValueIDs(m.mutable)
	if m.immutable != nil {
		getTagValueIDs(m.immutable)
	}
}

// getTagValueIDInMem gets tag value id from mutable/immutable store
func (m *tagMetadata) getTagValueIDInMem(tagKeyID uint32, tagValue string) (tagValueID uint32, ok bool) {
	tagValueID, ok = getTagValueID(m.mutable, tagKeyID, tagValue)
	if ok {
		return
	}
	if m.immutable != nil {
		tagValueID, ok = getTagValueID(m.immutable, tagKeyID, tagValue)
		if ok {
			return
		}
	}
	return
}

// getTagValueID gets tag value id from tag store based on tag key id and tag value
func getTagValueID(tags *TagStore, tagKeyID uint32, tagValue string) (tagValueID uint32, ok bool) {
	tag, ok := tags.Get(tagKeyID)
	if ok {
		tagValueID, ok = tag.getTagValueID(tagValue)
		if ok {
			return tagValueID, true
		}
	}
	return
}
