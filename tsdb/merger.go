package tsdb

import (
	"github.com/lindb/lindb/kv"
	"github.com/lindb/lindb/tsdb/tblstore/invertedindex"
	"github.com/lindb/lindb/tsdb/tblstore/metricsdata"
	"github.com/lindb/lindb/tsdb/tblstore/metricsmeta"
)

const (
	forwardIndexMerger  = "forward_index_merger"
	invertedIndexMerger = "inverted_index_merger"
	tagMetaMerger       = "tag_meta_merger"
	dataMerger          = "data_merge"
)

func init() {
	// FIXME stone1100
	kv.RegisterMerger(
		invertedIndexMerger,
		invertedindex.NewInvertedMerger())
	kv.RegisterMerger(
		forwardIndexMerger,
		invertedindex.NewForwardMerger())
	kv.RegisterMerger(
		tagMetaMerger,
		metricsmeta.NewTagMerger())
	kv.RegisterMerger(
		dataMerger,
		metricsdata.NewMerger())
}
