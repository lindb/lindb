package metricsdata

import (
	"fmt"
	"math"

	"github.com/lindb/roaring"

	"github.com/lindb/lindb/aggregation"
	"github.com/lindb/lindb/pkg/collections"
	"github.com/lindb/lindb/pkg/encoding"
	"github.com/lindb/lindb/pkg/stream"
	"github.com/lindb/lindb/series"
	"github.com/lindb/lindb/series/field"
)

//go:generate mockgen -source ./scanner.go -destination=./scanner_mock.go -package metricsdata

const (
	mdtLevel3FooterSize = 4 + // Series offset position
		4 + // series bitmap position
		4 //  field-meta position
)

// Scanner implements metrics from sstable.
type Scanner interface {
	//flow.Scanner
}

/*type metricsDataScanner struct {
	readers []table.Reader
	sr      *stream.Reader
}
*/
//func NewScanner(readers []table.Reader) flow.Scanner {
//	return &metricsDataScanner{
//		readers: readers,
//		sr:      stream.NewReader(nil)}
//}
//
//func (r *metricsDataScanner) Scan(qCtx *flow.StorageQueryContext) {
//version2Blocks := r.pickVersion2Blocks(qCtx)
//for _, mdtVersionBlocks := range version2Blocks {
//	for _, mdt := range mdtVersionBlocks {
//		//qCtx.Worker.Emit(mdt)
//	}
//}
//}
//
//func (r *metricsDataScanner) pickVersion2Blocks(
//	qCtx *flow.StorageQueryContext,
//) (
//	version2Blocks map[series.Version][]*mdtVersionBlock,
//) {
//version2Blocks = make(map[series.Version][]*mdtVersionBlock)
//for _, reader := range r.readers {
//	itr, err := tblstore.NewVersionBlockIterator(reader.Get(sCtx.MetricID))
//	if err != nil {
//		continue
//	}
//	for itr.HasNext() {
//		version, block := itr.Next()
//		if !sCtx.SeriesIDSet.Contains(version) {
//			continue
//		}
//		structuredBlock, err := newMDTVersionBlock(version, block, sCtx)
//		if err != nil {
//			continue
//		}
//		blockList, ok := version2Blocks[version]
//		if ok {
//			blockList = append(blockList, structuredBlock)
//		} else {
//			blockList = []*mdtVersionBlock{structuredBlock}
//		}
//		version2Blocks[version] = blockList
//	}
//}
//return version2Blocks
//	return
//}

// mdt is short for metric-data-table
// mdtVersionBlock implements ScanEvent
type mdtVersionBlock struct {
	version       series.Version
	block         []byte
	sr1           *stream.Reader
	sr2           *stream.Reader
	seriesOffsets *encoding.DeltaBitPackingDecoder
	seriesBitmap  *roaring.Bitmap
	fieldMetas    field.Metas
	//qCtx          *flow.StorageQueryContext
	// position
	seriesOffsetPos int
	seriesBitmapPos int
	fieldMetaPos    int

	bitArray    *collections.BitArray
	aggregators aggregation.FieldAggregates
}

//func newMDTVersionBlock(
//	version series.Version,
//	block []byte,
//) (
//	*mdtVersionBlock,
//	error,
//) {
//	if len(block) <= mdtLevel3FooterSize {
//		return nil, fmt.Errorf("failed validating version-block length")
//	}
//	vb := &mdtVersionBlock{
//		version:  version,
//		block:    block,
//		sr1:      stream.NewReader(block),
//		sr2:      stream.NewReader(block),
//		bitArray: collections.NewBitArray(nil),
//	}
//	// read footer
//	if err := vb.readFooter(); err != nil {
//		return nil, err
//	}
//	// read field-meta and time-range
//	if err := vb.readFieldMetas(); err != nil {
//		return nil, err
//	}
//	// read offsets
//	if err := vb.readOffsetsAndBitmap(); err != nil {
//		return nil, err
//	}
//	return vb, nil
//}

// initialize step1
func (vb *mdtVersionBlock) readFooter() error {
	vb.sr1.SeekStart()
	// read footer
	_ = vb.sr1.ReadSlice(len(vb.block) - mdtLevel3FooterSize)
	vb.seriesOffsetPos = int(vb.sr1.ReadUint32())
	vb.seriesBitmapPos = int(vb.sr1.ReadUint32())
	vb.fieldMetaPos = int(vb.sr1.ReadUint32())

	if 0 < vb.seriesOffsetPos &&
		vb.seriesOffsetPos < vb.seriesBitmapPos &&
		vb.seriesBitmapPos < vb.fieldMetaPos &&
		vb.fieldMetaPos < len(vb.block) {
		return nil
	}
	return fmt.Errorf("failed validating position")
}

// initialize step2
func (vb *mdtVersionBlock) readFieldMetas() error {
	vb.sr1.SeekStart()
	_ = vb.sr1.ReadSlice(vb.fieldMetaPos)
	// validate timeRange
	fieldCount := vb.sr1.ReadUvarint64()
	for i := 0; i < int(fieldCount); i++ {
		fieldID := vb.sr1.ReadUint16()
		fieldType := field.Type(vb.sr1.ReadByte())
		fieldName := vb.sr1.ReadSlice(int(vb.sr1.ReadUvarint64()))
		vb.fieldMetas = append(vb.fieldMetas, field.Meta{
			ID:   fieldID,
			Type: fieldType,
			Name: string(fieldName)})
		if vb.sr1.Error() != nil {
			return vb.sr1.Error()
		}
	}
	return nil
}

// initialize step3
func (vb *mdtVersionBlock) readOffsetsAndBitmap() error {
	// read offsets
	vb.seriesOffsets = encoding.NewDeltaBitPackingDecoder(vb.block[vb.seriesOffsetPos:vb.seriesBitmapPos])
	// read bitmap
	vb.seriesBitmap = roaring.New()
	if err := vb.seriesBitmap.UnmarshalBinary(vb.block[vb.seriesBitmapPos:vb.fieldMetaPos]); err != nil {
		return err
	}
	return nil
}

func (vb *mdtVersionBlock) SeriesIDs() *roaring.Bitmap {
	return vb.seriesBitmap
}

func (vb *mdtVersionBlock) TotalSeriesIDs() *roaring.Bitmap {
	return nil
}

func (vb *mdtVersionBlock) Version() series.Version {
	return vb.version
}

func (vb *mdtVersionBlock) Release() {
	// todo
	if vb.aggregators == nil {
		return
	}
	vb.aggregators.Reset()
}

func (vb *mdtVersionBlock) ResultSet() interface{} {
	// todo

	return nil
}

func (vb *mdtVersionBlock) SetGroupedSeries(highKey uint16, groupedSeries map[string][]uint16) {

}

/*
1. Scan version block -> Worker pool
2. Scan series entry
3. Scan field data

*/

func (vb *mdtVersionBlock) Scan(highKey uint16, groupedSeries map[string][]uint16) {
	//expectedSeriesIDs := vb.qCtx.SeriesIDSet.Versions()[vb.version]
	//var (
	//	currentPosition int32
	//	currentSeriesID uint32
	//)
	//itr := vb.seriesBitmap.Iterator()
	//for itr.HasNext() {
	//	currentSeriesID = itr.Next()
	//	if vb.seriesOffsets.HasNext() {
	//		currentPosition = vb.seriesOffsets.Next()
	//	} else {
	//		//FIXME
	//		return
	//	}
	//	//if !expectedSeriesIDs.Contains(currentSeriesID) {
	//	//	continue
	//	//}
	//	if vb.readFieldsData(currentPosition) != nil {
	//		//FIXME
	//		return
	//	}
	//}
	//return
}

func (vb *mdtVersionBlock) readFieldsData(position int32) error {
	vb.sr1.SeekStart()
	vb.sr1.ReadSlice(int(position))
	// read series entry
	// read fields-info
	_ = vb.sr1.ReadVarint64()
	_ = vb.sr1.ReadVarint64()
	// read bit-array
	bitArrayLen := int(math.Ceil(float64(vb.fieldMetas.Len()+1) / float64(8)))
	vb.bitArray.Reset(vb.sr1.ReadSlice(bitArrayLen))

	// preparing 2 stream readers
	endPosOfBitArray := vb.sr1.Position()
	for idx := range vb.fieldMetas {
		if vb.bitArray.GetBit(uint16(idx)) {
			_ = vb.sr1.ReadVarint64()
		}
	}
	startPosOfFieldsData := vb.sr1.Position()
	// sr2 points to fields-data
	vb.sr2.SeekStart()
	_ = vb.sr2.ReadSlice(startPosOfFieldsData)
	// sr1 points to data length list
	vb.sr1.SeekStart()
	_ = vb.sr1.ReadSlice(endPosOfBitArray)

	// jump to fields-data
	// TODO???
	//for idx, fm := range vb.fieldMetas {
	//	dataLength := vb.sr1.ReadUvarint64()
	//	if vb.sCtx.ContainsFieldID(fm.ID) {
	//		if !vb.bitArray.GetBit(uint16(idx)) {
	//			continue
	//		}
	//		data := vb.sr2.ReadSlice(int(dataLength))
	//		if vb.sr2.Error() != nil {
	//			return vb.sr2.Error()
	//		}
	//		if err := vb.readData(data); err != nil {
	//			return err
	//		}
	//	}
	//}
	return nil
}

func (vb *mdtVersionBlock) readData(data []byte) error {
	// todo

	return nil
}
