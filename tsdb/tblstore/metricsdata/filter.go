package metricsdata

import (
	"errors"
	"math"

	"github.com/lindb/roaring"

	"github.com/lindb/lindb/aggregation"
	"github.com/lindb/lindb/flow"
	"github.com/lindb/lindb/kv/version"
	"github.com/lindb/lindb/pkg/collections"
	"github.com/lindb/lindb/pkg/encoding"
	"github.com/lindb/lindb/pkg/stream"
	"github.com/lindb/lindb/series"
	"github.com/lindb/lindb/series/field"
	f "github.com/lindb/lindb/tsdb/field"
	"github.com/lindb/lindb/tsdb/tblstore"
)

//go:generate mockgen -source ./filter.go -destination=./filter_mock.go -package metricsdata

// for testing
var (
	newVersionBlock = newMDTVersionBlock
)

var (
	errBlockLength  = errors.New("failed validating version-block length")
	errReadPosition = errors.New("failed validating position")
)

const (
	mdtLevel3FooterSize = 4 + // Series offset position
		4 + // series bitmap position
		4 //  field-meta position
)

// fieldIndex represents need read field index based on metric field metas
type fieldIndex struct {
	fieldType field.Type
	need      bool // if need read
	idx       int  // search field's index
}

// Filter implements filtering metrics from sst files.
type Filter interface {
	// Filter filters data under each sst file based on query condition
	Filter(fieldIDs []uint16, version series.Version, seriesIDs *roaring.Bitmap) ([]flow.FilterResultSet, error)
}

// metricVersionBlock represents the metric version data block for loading data
type metricVersionBlock interface {
	load(flow flow.StorageQueryFlow, highKey uint16, groupedSeries map[string][]uint16)
}

// metricsDataFilter represents the sst file data filter
type metricsDataFilter struct {
	familyTime int64
	snapshot   version.Snapshot //FIXME stone1100, need close version snapshot
	blockIts   []tblstore.VersionBlockIterator
}

// NewFilter creates the sst file data filter
func NewFilter(familyTime int64, snapshot version.Snapshot, blockIts []tblstore.VersionBlockIterator) Filter {
	return &metricsDataFilter{
		familyTime: familyTime,
		snapshot:   snapshot,
		blockIts:   blockIts,
	}
}

// Filter filters the data under each sst file based on metric/version/seriesIDs,
// if finds data then returns the FilterResultSet, else returns nil
func (f *metricsDataFilter) Filter(fieldIDs []uint16,
	version series.Version, seriesIDs *roaring.Bitmap,
) (rs []flow.FilterResultSet, err error) {
	for _, it := range f.blockIts {
		for it.HasNext() {
			v, block := it.Peek()
			switch {
			case v == version:
				// build file filter result set
				r, err := newFileFilterResultSet(f.familyTime, fieldIDs, block)
				if err != nil {
					return nil, err
				}
				// maybe result set is nil, because field ids not found
				if r != nil {
					rs = append(rs, r)
				}
				// skip to next version
				it.Next()
			case v < version:
				// skip to next version
				it.Next()
			default:
				// if version > v, query version not found
				return
			}
		}
	}
	return
}

// mdt is short for metric-data-table

// mdtVersionBlock implements ScanEvent
type mdtVersionBlock struct {
	familyTime         int64
	block              []byte
	seriesBucketOffset *encoding.FixedOffsetDecoder
	seriesBitmap       *roaring.Bitmap
	fieldMetas         field.Metas
	fieldIndexes       []fieldIndex
	// position
	seriesBucketOffsetPos int
	seriesBitmapPos       int
	fieldMetaPos          int

	bitArrayLen int
}

func newMDTVersionBlock(familyTime int64, fieldIDs []uint16, block []byte) (metricVersionBlock, error) {
	if len(block) <= mdtLevel3FooterSize {
		return nil, errBlockLength
	}
	vb := &mdtVersionBlock{
		familyTime: familyTime,
		block:      block,
	}
	// read footer
	if err := vb.readFooter(); err != nil {
		return nil, err
	}
	// read field-meta
	vb.readFieldMetas()
	// check query fields is match
	needFieldCount := len(fieldIDs)
	vb.fieldIndexes = make([]fieldIndex, len(vb.fieldMetas))
	found := 0
	for idx, fieldMeta := range vb.fieldMetas {
		fieldIdx := fieldIndex{}
		for needIdx, fieldID := range fieldIDs {
			if fieldMeta.ID == fieldID {
				fieldIdx.fieldType = fieldMeta.Type
				fieldIdx.need = true
				fieldIdx.idx = needIdx
				found++
			}
		}
		vb.fieldIndexes[idx] = fieldIdx
	}
	// check query fields if found
	if found != needFieldCount {
		return nil, nil
	}

	vb.bitArrayLen = int(math.Ceil(float64(vb.fieldMetas.Len()+1) / float64(8)))
	// read offsets
	if err := vb.readOffsetsAndBitmap(); err != nil {
		return nil, err
	}
	return vb, nil
}

// initialize step1
func (vb *mdtVersionBlock) readFooter() error {
	// read footer
	footerStartPos := len(vb.block) - mdtLevel3FooterSize
	vb.seriesBitmapPos = int(stream.ReadUint32(vb.block, footerStartPos))
	vb.seriesBucketOffsetPos = int(stream.ReadUint32(vb.block, footerStartPos+4))
	vb.fieldMetaPos = int(stream.ReadUint32(vb.block, footerStartPos+8))

	if 0 < vb.seriesBitmapPos &&
		vb.seriesBitmapPos < vb.seriesBucketOffsetPos &&
		vb.seriesBucketOffsetPos < vb.fieldMetaPos &&
		vb.fieldMetaPos < len(vb.block) {
		return nil
	}
	return errReadPosition
}

// initialize step2
func (vb *mdtVersionBlock) readFieldMetas() {
	offset := vb.fieldMetaPos
	fieldCount := stream.ReadUint16(vb.block, offset)
	offset += 2
	for i := 0; i < int(fieldCount); i++ {
		fieldID := stream.ReadUint16(vb.block, offset)
		offset += 2
		fieldType := field.Type(vb.block[offset])
		offset++
		vb.fieldMetas = append(vb.fieldMetas, field.Meta{
			ID:   fieldID,
			Type: fieldType,
		})
	}
}

// initialize step3
func (vb *mdtVersionBlock) readOffsetsAndBitmap() error {
	// read bitmap
	vb.seriesBitmap = roaring.New()
	if err := vb.seriesBitmap.UnmarshalBinary(vb.block[vb.seriesBitmapPos:vb.seriesBucketOffsetPos]); err != nil {
		return err
	}
	// read series bucket offsets
	vb.seriesBucketOffset = encoding.NewFixedOffsetDecoder(vb.block[vb.seriesBucketOffsetPos:vb.fieldMetaPos])
	return nil
}

func (vb *mdtVersionBlock) load(flow flow.StorageQueryFlow, highKey uint16, groupedSeries map[string][]uint16) {
	// 1. get high container index by the high key of series ID
	highContainerIdx := vb.seriesBitmap.GetContainerIndex(highKey)
	if highContainerIdx < 0 {
		// if high container index < 0(series IDs not exist) return it
		return
	}
	// 2. get low container include all low keys by the high container index, delete op will clean empty low container
	lowContainer := vb.seriesBitmap.GetContainerAtIndex(highContainerIdx)
	seriesOffsets := encoding.NewFixedOffsetDecoder(vb.block[vb.seriesBucketOffset.Get(highContainerIdx):])

	//var aggregators aggregation.FieldAggregates
	tsd := encoding.GetTSDDecoder()
	defer encoding.ReleaseTSDDecoder(tsd)

	for groupByTags, lowSeriesIDs := range groupedSeries {
		aggregator := flow.GetAggregator()
		for _, lowSeriesID := range lowSeriesIDs {
			// check low series id if exist
			if !lowContainer.Contains(lowSeriesID) {
				continue
			}
			// get the index of low series id in container
			idx := lowContainer.Rank(lowSeriesID)
			// scan the data and aggregate the values
			seriesPos := seriesOffsets.Get(idx - 1)
			// read fields data and agg it
			vb.readFieldsData(tsd, aggregator, seriesPos)
		}
		flow.Reduce(groupByTags, aggregator)
	}
}

func (vb *mdtVersionBlock) readFieldsData(tsd *encoding.TSDDecoder, agg aggregation.FieldAggregates, position int) {
	// read bit-array
	offset := position + vb.bitArrayLen
	bitArray := collections.NewBitArray(vb.block[position:offset])

	lens := make([]int, len(vb.fieldIndexes))
	// preparing 2 stream readers
	for idx := range vb.fieldIndexes {
		if bitArray.GetBit(uint16(idx)) {
			l, length, _ := stream.ReadUvarint(vb.block, offset)
			offset += length
			lens[idx] = int(l)
		}
	}

	// jump to fields-data
	for idx, fieldIdx := range vb.fieldIndexes {
		if !bitArray.GetBit(uint16(idx)) {
			continue
		}
		dataLength := lens[idx]
		end := offset + dataLength
		// if field need, read the field data and aggregate the values
		if fieldIdx.need {
			data := vb.block[offset:end]
			a := agg[fieldIdx.idx]
			vb.readData(fieldIdx.fieldType, a, tsd, data)
		}
		// goto next field's offset
		offset = end
	}
}

func (vb *mdtVersionBlock) readData(fieldType field.Type, agg aggregation.SeriesAggregator,
	tsd *encoding.TSDDecoder, data []byte,
) {
	segmentAgg, ok := agg.GetAggregator(vb.familyTime)
	if !ok {
		return
	}
	f.Aggregate(fieldType, segmentAgg, tsd, data)
}

type fileFilterResultSet struct {
	block metricVersionBlock
}

// newFileFilterResultSet creates the file filter result set
func newFileFilterResultSet(familyTime int64, fieldIDs []uint16, block []byte) (flow.FilterResultSet, error) {
	b, err := newVersionBlock(familyTime, fieldIDs, block)
	if err != nil {
		return nil, err
	}
	return &fileFilterResultSet{
		block: b,
	}, nil
}

// Load reads data from sst files, finds series data based on grouped series IDs and does down sampling,
// finally reduces the down sampling result set.
func (f *fileFilterResultSet) Load(flow flow.StorageQueryFlow, fieldIDs []uint16,
	highKey uint16, groupedSeries map[string][]uint16,
) {
	f.block.load(flow, highKey, groupedSeries)
}
