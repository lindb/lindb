package metricsmeta

import (
	"github.com/lindb/lindb/kv"
)

var TagMetaMerger kv.MergerType = "TagMetaMerger"

// init registers tag meta merger create function
func init() {
	kv.RegisterMerger(TagMetaMerger, NewTagMerger)
}

// tagMerger implements kv.Merger for merging tag trie meta data for each metric
type tagMerger struct {
	tagFlusher TagFlusher
	flusher    *kv.NopFlusher
}

// NewTagMerger creates a tag merger
func NewTagMerger() kv.Merger {
	flusher := kv.NewNopFlusher()
	return &tagMerger{
		flusher:    flusher,
		tagFlusher: NewTagFlusher(flusher),
	}
}

func (t *tagMerger) Init(params map[string]interface{}) {
	// do nothing
}

// Merge merges the multi tag trie meta data into a trie for same metric
func (t *tagMerger) Merge(key uint32, values [][]byte) ([]byte, error) {
	maxSequenceID := uint32(0) // target sequence of tag value id
	// 1. prepare tag trie readers
	var readers []TagKVEntrySetINTF
	for _, value := range values {
		reader, err := newTagKVEntrySetFunc(value)
		if err != nil {
			return nil, err
		}
		if maxSequenceID < reader.TagValueSeq() {
			maxSequenceID = reader.TagValueSeq()
		}
		readers = append(readers, reader)
	}
	// 2. iterator trie data, then merge the tag values
	for _, reader := range readers {
		q, err := reader.TrieTree()
		if err != nil {
			return nil, err
		}
		offsetsItr := q.Iterator("")
		for offsetsItr.HasNext() {
			tagValue, offset := offsetsItr.Next()
			tagValueID := reader.GetTagValueID(offset)
			t.tagFlusher.FlushTagValue(string(tagValue), tagValueID)
		}
	}
	if err := t.tagFlusher.FlushTagKeyID(key, maxSequenceID); err != nil {
		return nil, err
	}
	return t.flusher.Bytes(), nil
}
