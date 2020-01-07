package metricsdata

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/lindb/roaring"
	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/aggregation"
	"github.com/lindb/lindb/flow"
	"github.com/lindb/lindb/kv"
	"github.com/lindb/lindb/pkg/bit"
	"github.com/lindb/lindb/pkg/encoding"
	"github.com/lindb/lindb/pkg/stream"
	"github.com/lindb/lindb/series"
	"github.com/lindb/lindb/series/field"
	"github.com/lindb/lindb/tsdb/tblstore"
)

func TestMetricsDataFilter_Filter(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer func() {
		newVersionBlock = newMDTVersionBlock
		ctrl.Finish()
	}()

	mvBlockIt := tblstore.NewMockVersionBlockIterator(ctrl)
	dataFilter := NewFilter(0, nil, []tblstore.VersionBlockIterator{mvBlockIt})
	mvBlock := NewMockmetricVersionBlock(ctrl)
	newVersionBlock = func(familyTime int64, fieldIDs []uint16, block []byte) (block2 metricVersionBlock, e error) {
		return mvBlock, nil
	}
	gomock.InOrder(
		mvBlockIt.EXPECT().HasNext().Return(true),
		mvBlockIt.EXPECT().Peek().Return(series.Version(1), nil),
		mvBlockIt.EXPECT().Next(),
		mvBlockIt.EXPECT().HasNext().Return(true),
		mvBlockIt.EXPECT().Peek().Return(series.Version(10), nil),
		mvBlockIt.EXPECT().Next(),
		mvBlockIt.EXPECT().HasNext().Return(false),
	)
	rs, err := dataFilter.Filter(nil, series.Version(10), nil)
	assert.NoError(t, err)
	assert.Len(t, rs, 1)

	// test not found
	dataFilter = NewFilter(0, nil, []tblstore.VersionBlockIterator{mvBlockIt})
	gomock.InOrder(
		mvBlockIt.EXPECT().HasNext().Return(true),
		mvBlockIt.EXPECT().Peek().Return(series.Version(100), nil),
	)
	rs, err = dataFilter.Filter(nil, series.Version(10), nil)
	assert.NoError(t, err)
	assert.Nil(t, rs)

	// new metric version block err
	newVersionBlock = func(familyTime int64, fieldIDs []uint16, block []byte) (block2 metricVersionBlock, e error) {
		return nil, fmt.Errorf("err")
	}
	gomock.InOrder(
		mvBlockIt.EXPECT().HasNext().Return(true),
		mvBlockIt.EXPECT().Peek().Return(series.Version(10), nil),
	)
	rs, err = dataFilter.Filter(nil, series.Version(10), nil)
	assert.Error(t, err)
	assert.Nil(t, rs)
}

func Test_file_filterResultSet(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer func() {
		ctrl.Finish()
		newVersionBlock = newMDTVersionBlock
	}()

	newVersionBlock = func(familyTime int64, fieldIDs []uint16, block []byte) (block2 metricVersionBlock, e error) {
		return nil, fmt.Errorf("err")
	}
	rs, err := newFileFilterResultSet(0, nil, []byte{1, 2, 3})
	assert.Error(t, err)
	assert.Nil(t, rs)

	mvBlock := NewMockmetricVersionBlock(ctrl)
	newVersionBlock = func(familyTime int64, fieldIDs []uint16, block []byte) (block2 metricVersionBlock, e error) {
		return mvBlock, nil
	}
	mvBlock.EXPECT().load(gomock.Any(), gomock.Any(), gomock.Any())
	rs, err = newFileFilterResultSet(0, nil, []byte{1, 2, 3})
	assert.NoError(t, err)
	assert.NotNil(t, rs)
	rs.Load(nil, nil, 1, nil)
}

func Test_newMDTVersionBlock(t *testing.T) {
	// test block footer err
	vb, err := newMDTVersionBlock(1, nil, []byte{12, 3})
	assert.Error(t, err)
	assert.Nil(t, vb)

	// test read footer err
	var buf bytes.Buffer
	writer2 := stream.NewBufferWriter(&buf)
	writer2.PutUint32(2)
	writer2.PutUint32(999)
	writer2.PutUint32(10)
	writer2.PutUint32(10)
	data, _ := writer2.Bytes()
	vb, err = newMDTVersionBlock(1, nil, data)
	assert.Error(t, err)
	assert.Nil(t, vb)

	// mock metric version data
	data = buildMetricVersionBlockData()

	// test fields not match
	vb, err = newMDTVersionBlock(1, []uint16{9}, data)
	assert.NoError(t, err)
	assert.Nil(t, vb)

	// normal case
	vb, err = newMDTVersionBlock(1, []uint16{2}, data)
	assert.NoError(t, err)
	assert.NotNil(t, vb)

	// UnmarshalBinary bitmap err
	// set err bitmap data
	data[6]++
	vb, err = newMDTVersionBlock(1, []uint16{2}, data)
	assert.Error(t, err)
	assert.Nil(t, vb)
}

func Test_newMDTVersionBlock_load(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// mock metric version data
	data := buildMetricVersionBlockData()

	vb, err := newMDTVersionBlock(1, []uint16{2}, data)
	assert.NoError(t, err)
	assert.NotNil(t, vb)
	qf := flow.NewMockStorageQueryFlow(ctrl)
	// test high key not match
	vb.load(qf, 9, nil)

	// test load data
	sAgg := aggregation.NewMockSeriesAggregator(ctrl)
	qf.EXPECT().GetAggregator().Return(aggregation.FieldAggregates{sAgg})
	fAgg := aggregation.NewMockFieldAggregator(ctrl)
	sAgg.EXPECT().GetAggregator(gomock.Any()).Return(fAgg, false)
	qf.EXPECT().Reduce(gomock.Any(), gomock.Any())
	vb.load(qf, 0, map[string][]uint16{"test": {1, 2, 10}})
}

func Test_newMDTVersionBlock_readFiled(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// mock metric version data
	data := buildMetricVersionBlockData()

	vb, err := newMDTVersionBlock(1, []uint16{2}, data)
	assert.NoError(t, err)
	assert.NotNil(t, vb)
	vb1 := vb.(*mdtVersionBlock)

	sAgg := aggregation.NewMockSeriesAggregator(ctrl)
	fAgg := aggregation.NewMockFieldAggregator(ctrl)
	sAgg.EXPECT().GetAggregator(gomock.Any()).Return(fAgg, true)
	pAgg := aggregation.NewMockPrimitiveAggregator(ctrl)
	pAgg.EXPECT().Aggregate(gomock.Any(), gomock.Any())
	fAgg.EXPECT().GetAllAggregators().Return([]aggregation.PrimitiveAggregator{pAgg})

	tsd := encoding.GetTSDDecoder()
	encoder := encoding.NewTSDEncoder(10)
	encoder.AppendTime(bit.One)
	encoder.AppendValue(uint64(10))
	data, _ = encoder.Bytes()
	vb1.readData(field.SumField, sAgg, tsd, data)
}

func buildMetricVersionBlockData() []byte {
	nopKvFlusher := kv.NewNopFlusher()
	flusherImpl := NewFlusher(nopKvFlusher)
	flusherImpl.FlushFieldMetas([]field.Meta{
		{ID: 1, Type: field.SumField, Name: "sum"},
		{ID: 2, Type: field.MinField, Name: "min"},
		{ID: 3, Type: field.MaxField, Name: "max"},
	})
	flusherImpl.FlushField(1)
	flusherImpl.FlushField(2)
	flusherImpl.FlushSeries()
	flusherImpl.FlushSeriesBucket()
	flusherImpl.FlushVersion(series.NewVersion(), roaring.BitmapOf(1, 2, 3))
	f := flusherImpl.(*flusher)
	data, _ := f.writer.Bytes()
	return data
}
