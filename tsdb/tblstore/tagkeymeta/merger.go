package tagkeymeta

import (
	"github.com/lindb/lindb/kv"
	"github.com/lindb/lindb/pkg/encoding"
)

var MergerName kv.MergerType = "TagKeyMetaMerger"

// init registers tag meta merger create function
func init() {
	kv.RegisterMerger(MergerName, NewMerger)
}

// merger implements kv.Merger for merging tag trie meta data for each metric
type merger struct {
	flusher   Flusher
	kvFlusher *kv.NopFlusher
}

// NewTagMerger creates a merger for compacting tag-key-meta
func NewMerger() kv.Merger {
	kvFlusher := kv.NewNopFlusher()
	return &merger{
		kvFlusher: kvFlusher,
		flusher:   NewFlusher(kvFlusher),
	}
}

func (tm *merger) Init(params map[string]interface{}) {
	// do nothing
}

func cloneSlice(slice []byte) []byte {
	if len(slice) == 0 {
		return nil
	}
	cloned := make([]byte, len(slice))
	copy(cloned, slice)
	return cloned
}

// Merge merges the multi tag trie meta data into a trie for same metric
func (tm *merger) Merge(tagKeyID uint32, dataBlocks [][]byte) ([]byte, error) {
	maxSequenceID := uint32(0) // target sequence of tag value id
	// 1. prepare tagKeyMetas
	var tagKeyMetas []TagKeyMeta
	for _, dataBlock := range dataBlocks {
		tagKeyMeta, err := newTagKeyMeta(dataBlock)
		if err != nil {
			return nil, err
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
			return nil, err
		}
		for itr.Valid() {
			tm.flusher.FlushTagValue(cloneSlice(itr.Key()), encoding.ByteSlice2Uint32(itr.Value()))
			itr.Next()
		}
	}
	if err := tm.flusher.FlushTagKeyID(tagKeyID, maxSequenceID); err != nil {
		return nil, err
	}
	return tm.kvFlusher.Bytes(), nil
}
