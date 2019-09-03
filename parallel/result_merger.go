package parallel

import (
	"github.com/lindb/lindb/pkg/stream"
	pb "github.com/lindb/lindb/rpc/proto/common"
	"github.com/lindb/lindb/series"
	"github.com/lindb/lindb/series/field"
)

//go:generate mockgen -source=./result_merger.go -destination=./result_merger_mock.go -package=parallel

type ResultMerger interface {
	Merge(resp *pb.TaskResponse)
}

type resultMerger struct {
	resultSet chan *series.TimeSeriesEvent
}

func newResultMerger(resultSet chan *series.TimeSeriesEvent) ResultMerger {
	return &resultMerger{resultSet: resultSet}
}

func (m *resultMerger) Merge(resp *pb.TaskResponse) {
	data := resp.Payload
	ts := &pb.TimeSeries{}
	//todo handle err
	_ = ts.Unmarshal(data)
	m.resultSet <- &series.TimeSeriesEvent{
		Series: newGroupedIterator(ts.Fields),
	}
}

type groupedIterator struct {
	fields     map[string][]byte
	fieldNames []string

	idx int
}

func newGroupedIterator(fields map[string][]byte) series.GroupedIterator {
	it := &groupedIterator{fields: fields}
	for fieldName := range fields {
		it.fieldNames = append(it.fieldNames, fieldName)
	}
	return it
}

func (g *groupedIterator) Tags() map[string]string {
	return nil
}
func (g *groupedIterator) HasNext() bool {
	if g.idx >= len(g.fieldNames) {
		return false
	}
	g.idx++
	return true
}

func (g *groupedIterator) Next() series.FieldIterator {
	fieldName := g.fieldNames[g.idx-1]
	return newFieldIterator(fieldName, g.fields[fieldName])
}

func (g *groupedIterator) SeriesID() uint32 {
	return 10
}

type fieldIterator struct {
	fieldName string
	reader    *stream.Reader

	segmentStartTime int64
}

func newFieldIterator(fieldName string, data []byte) series.FieldIterator {
	it := &fieldIterator{
		fieldName: fieldName,
		reader:    stream.NewReader(data),
	}
	it.segmentStartTime = it.reader.ReadVarint64()
	return it
}

func (fsi *fieldIterator) FieldID() uint16       { return 0 }
func (fsi *fieldIterator) FieldName() string     { return fsi.fieldName }
func (fsi *fieldIterator) FieldType() field.Type { return field.Unknown }

func (fsi *fieldIterator) HasNext() bool {
	return !fsi.reader.Empty()
}

func (fsi *fieldIterator) Next() series.PrimitiveIterator {
	fieldID := fsi.reader.ReadUint16()
	length := fsi.reader.ReadVarint32()
	data := fsi.reader.ReadBytes(int(length))

	return series.NewPrimitiveIterator(fieldID, data)
}

func (fsi *fieldIterator) Bytes() ([]byte, error) {
	//FIXME stone1100
	return nil, nil
}

func (fsi *fieldIterator) SegmentStartTime() int64 {
	return fsi.segmentStartTime
}
