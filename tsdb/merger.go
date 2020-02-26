package tsdb

import (
	"github.com/lindb/lindb/kv"
	"github.com/lindb/lindb/tsdb/tblstore/invertedindex"
	"github.com/lindb/lindb/tsdb/tblstore/metricsmeta"
)

const (
	forwardIndexMerger  = "forward_index_merger"
	invertedIndexMerger = "inverted_index_merger"
	tagMetaMerger       = "tag_meta_merger"
	//defaultTTLDuration  = time.Hour * 24 * 30
	nopMerger = "nop_merger"
)

func init() {
	// FIXME stone1100
	kv.RegisterMerger(
		invertedIndexMerger,
		invertedindex.NewInvertedMerger())
	kv.RegisterMerger(
		forwardIndexMerger,
		invertedindex.NewForwardMerge())
	kv.RegisterMerger(
		tagMetaMerger,
		metricsmeta.NewTagMerger())

	kv.RegisterMerger(nopMerger, &_nopMerger{})
}

// nopMerger does nothing
type _nopMerger struct{}

func (m *_nopMerger) Merge(key uint32, values [][]byte) ([]byte, error) {
	return nil, nil
}
