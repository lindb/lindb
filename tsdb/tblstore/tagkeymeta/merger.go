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

package tagkeymeta

import (
	"github.com/lindb/lindb/kv"
)

var MergerName kv.MergerType = "TagKeyMetaMerger"

// init registers tag meta merger create function
func init() {
	kv.RegisterMerger(MergerName, NewMerger)
}

// merger implements kv.Merger for merging tag trie meta data for each metric
type merger struct {
	metaFlusher Flusher
	kvFlusher   kv.Flusher
}

// NewMerger creates a merger for compacting tag-key-meta
func NewMerger(kvFlusher kv.Flusher) (kv.Merger, error) {
	metaFlusher, err := NewFlusher(kvFlusher)
	if err != nil {
		return nil, err
	}
	return &merger{
		metaFlusher: metaFlusher,
		kvFlusher:   kvFlusher,
	}, nil
}

func (tm *merger) Init(_ map[string]interface{}) {}

func cloneSlice(slice []byte) []byte {
	if len(slice) == 0 {
		return nil
	}
	cloned := make([]byte, len(slice))
	copy(cloned, slice)
	return cloned
}

// Merge merges the multi tag trie meta data into a trie for same metric
func (tm *merger) Merge(tagKeyID uint32, dataBlocks [][]byte) error {
	maxSequenceID := uint32(0) // target sequence of tag value id
	// 1. prepare tagKeyMetas
	var tagKeyMetas []TagKeyMeta
	for _, dataBlock := range dataBlocks {
		tagKeyMeta, err := newTagKeyMeta(dataBlock)
		if err != nil {
			return err
		}
		if maxSequenceID < tagKeyMeta.TagValueIDSeq() {
			maxSequenceID = tagKeyMeta.TagValueIDSeq()
		}
		tagKeyMetas = append(tagKeyMetas, tagKeyMeta)
	}
	// 2. iterator trie data, then merge the tag values
	for _, tagKeyMeta := range tagKeyMetas {
		itr, err := tagKeyMeta.PrefixIterator(nil)
		if err != nil {
			return err
		}
		for itr.Valid() {
			tm.metaFlusher.FlushTagValue(cloneSlice(itr.Key()), itr.Value())
			itr.Next()
		}
	}
	if err := tm.metaFlusher.FlushTagKeyID(tagKeyID, maxSequenceID); err != nil {
		return err
	}
	return tm.metaFlusher.commitTagKeyID()
}
