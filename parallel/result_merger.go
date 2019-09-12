package parallel

import (
	"context"

	pb "github.com/lindb/lindb/rpc/proto/common"
	"github.com/lindb/lindb/series"
)

//go:generate mockgen -source=./result_merger.go -destination=./result_merger_mock.go -package=parallel

// ResultMerger represents a merger which merges the task response and aggregates the result
type ResultMerger interface {
	// merge merges the task response and aggregates the result
	merge(resp *pb.TaskResponse)

	close()
}

type resultMerger struct {
	resultSet chan *series.TimeSeriesEvent

	events chan *pb.TaskResponse

	closed chan struct{}
	ctx    context.Context
}

// newResultMerger create a result merger
func newResultMerger(ctx context.Context, resultSet chan *series.TimeSeriesEvent) ResultMerger {
	merger := &resultMerger{
		resultSet: resultSet,
		events:    make(chan *pb.TaskResponse),
		closed:    make(chan struct{}),
		ctx:       ctx,
	}
	go func() {
		defer close(merger.closed)
		merger.process()
	}()
	return merger
}

// merge merges and aggregates the result
func (m *resultMerger) merge(resp *pb.TaskResponse) {
	m.events <- resp
}

func (m *resultMerger) close() {
	close(m.events)

	// waiting process completed
	<-m.closed
}

func (m *resultMerger) process() {
	for {
		select {
		case event, ok := <-m.events:
			if !ok {
				return
			}
			err := m.handleEvent(event)
			if err != nil {
				m.resultSet <- &series.TimeSeriesEvent{Err: err}
				return
			}
		case <-m.ctx.Done():
			return
		}
	}
}

func (m *resultMerger) handleEvent(resp *pb.TaskResponse) error {
	data := resp.Payload
	ts := &pb.TimeSeries{}
	err := ts.Unmarshal(data)
	if err != nil {
		return err
	}

	m.resultSet <- &series.TimeSeriesEvent{
		//Series: newGroupedIterator(ts.Fields),
	}
	return nil
}

//type groupedIterator struct {
//	fields     map[string][]byte
//	fieldNames []string
//
//	idx int
//}

//func newGroupedIterator(fields map[string][]byte) series.GroupedIterator {
//	it := &groupedIterator{fields: fields}
//	for fieldName := range fields {
//		it.fieldNames = append(it.fieldNames, fieldName)
//	}
//	return it
//}
//
//func (g *groupedIterator) Tags() map[string]string {
//	return nil
//}
//func (g *groupedIterator) HasNext() bool {
//	if g.idx >= len(g.fieldNames) {
//		return false
//	}
//	g.idx++
//	return true
//}
//
//func (g *groupedIterator) Next() series.FieldIterator {
//	fieldName := g.fieldNames[g.idx-1]
//	return newFieldIterator(fieldName, g.fields[fieldName])
//}
//
//func (g *groupedIterator) SeriesID() uint32 {
//	return 10
//}
//
//type fieldIterator struct {
//	fieldName string
//	reader    *stream.Reader
//
//	segmentStartTime int64
//}
//
//func newFieldIterator(fieldName string, data []byte) series.FieldIterator {
//	it := &fieldIterator{
//		fieldName: fieldName,
//		reader:    stream.NewReader(data),
//	}
//	it.segmentStartTime = it.reader.ReadVarint64()
//	return it
//}
//
//func (fsi *fieldIterator) FieldName() string     { return fsi.fieldName }
//func (fsi *fieldIterator) FieldType() field.Type { return field.Unknown }
//
//func (fsi *fieldIterator) HasNext() bool {
//	return !fsi.reader.Empty()
//}
//
//func (fsi *fieldIterator) Next() series.PrimitiveIterator {
//	fieldID := fsi.reader.ReadUint16()
//	length := fsi.reader.ReadVarint32()
//	data := fsi.reader.ReadBytes(int(length))
//
//	return series.NewPrimitiveIterator(fieldID, data)
//}
//
//func (fsi *fieldIterator) Bytes() ([]byte, error) {
//	//FIXME stone1100
//	return nil, nil
//}
//
//func (fsi *fieldIterator) SegmentStartTime() int64 {
//	return fsi.segmentStartTime
//}
