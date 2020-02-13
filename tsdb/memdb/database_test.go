package memdb

import (
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/pkg/fileutil"
	"github.com/lindb/lindb/pkg/timeutil"
	pb "github.com/lindb/lindb/rpc/proto/field"
	"github.com/lindb/lindb/series/field"
	"github.com/lindb/lindb/tsdb/indexdb"
	"github.com/lindb/lindb/tsdb/metadb"
	"github.com/lindb/lindb/tsdb/tblstore/metricsdata"
)

const testDBPath = "test_db"

var cfg = MemoryDatabaseCfg{
	Interval: timeutil.Interval(10 * timeutil.OneSecond),
}

func TestNewMemoryDatabase(t *testing.T) {
	mdINTF := NewMemoryDatabase(cfg)
	assert.NotNil(t, mdINTF)
	assert.Equal(t, 10*timeutil.OneSecond, mdINTF.Interval())
}

func TestMemoryDatabase_Write(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer func() {
		_ = fileutil.RemoveDir(testDBPath)
		defer ctrl.Finish()
	}()
	_ = fileutil.MkDirIfNotExist(testDBPath)
	cfg.TempPath = testDBPath

	// mock
	mockMetadata := metadb.NewMockMetadata(ctrl)
	mockMetadataDatabase := metadb.NewMockMetadataDatabase(ctrl)
	mockMetadata.EXPECT().MetadataDatabase().Return(mockMetadataDatabase).AnyTimes()
	mockIndex := indexdb.NewMockIndexDatabase(ctrl)
	mockMStore := NewMockmStoreINTF(ctrl)
	tStore := NewMocktStoreINTF(ctrl)
	fStore := NewMockfStoreINTF(ctrl)
	mockMStore.EXPECT().GetOrCreateTStore(uint32(10)).Return(tStore, 10).AnyTimes()
	// build memory-database
	cfg.Index = mockIndex
	cfg.Metadata = mockMetadata
	mdINTF := NewMemoryDatabase(cfg)
	md := mdINTF.(*memoryDatabase)
	assert.Zero(t, md.MemSize())

	// load mock
	md.mStores.Put(uint32(1), mockMStore)
	// case 1: write ok
	gomock.InOrder(
		mockMetadataDatabase.EXPECT().GenFieldID("ns", "test1", "f1", field.SumField).Return(field.ID(1), nil),
		tStore.EXPECT().GetFStore(gomock.Any(), gomock.Any(), gomock.Any()).Return(fStore, true),
		fStore.EXPECT().Write(gomock.Any(), gomock.Any(), gomock.Any()).Return(10),
		mockMStore.EXPECT().AddField(gomock.Any(), gomock.Any()),
		mockMStore.EXPECT().SetTimestamp(gomock.Any(), gomock.Any()),
	)
	err := md.Write("ns", "test1", uint32(1), uint32(10), 1564300800000, []*pb.Field{{
		Name:  "f1",
		Field: &pb.Field_Sum{Sum: &pb.Sum{Value: 10.0}},
	}})
	assert.NoError(t, err)
	assert.Len(t, md.Families(), 1)
	// case 2: field type unknown
	err = md.Write("ns", "test1", uint32(1), uint32(10), 1564300800000, []*pb.Field{{
		Name:  "f1",
		Field: nil,
	}})
	assert.NoError(t, err)
	// case 3: generate field err
	mockMetadataDatabase.EXPECT().GenFieldID("ns", "test1", "f1-err", field.SumField).Return(field.ID(0), fmt.Errorf("err"))
	err = md.Write("ns", "test1", uint32(1), uint32(10), 1564300800000, []*pb.Field{{
		Name:  "f1-err",
		Field: &pb.Field_Sum{Sum: &pb.Sum{Value: 10.0}},
	}})
	assert.NoError(t, err)
	// case 4: new family times
	gomock.InOrder(
		mockMetadataDatabase.EXPECT().GenFieldID("ns", "test1", "f1", field.SumField).Return(field.ID(1), nil),
		tStore.EXPECT().GetFStore(gomock.Any(), gomock.Any(), gomock.Any()).Return(fStore, true),
		fStore.EXPECT().Write(gomock.Any(), gomock.Any(), gomock.Any()).Return(10),
		mockMStore.EXPECT().AddField(gomock.Any(), gomock.Any()),
		mockMStore.EXPECT().SetTimestamp(gomock.Any(), gomock.Any()),
	)
	err = md.Write("ns", "test1", uint32(1), uint32(10), 1564300800000+timeutil.OneHour, []*pb.Field{{
		Name:  "f1",
		Field: &pb.Field_Sum{Sum: &pb.Sum{Value: 10.0}},
	}})
	assert.NoError(t, err)
	assert.Len(t, md.Families(), 2)
	assert.True(t, md.MemSize() > 0)
	// case 5: new metric store
	err = md.Write("ns", "test1", uint32(20), uint32(20), 1564300800000, []*pb.Field{{
		Name: "f1",
	}})
	assert.NoError(t, err)
	// case 6: create new field store
	gomock.InOrder(
		mockMetadataDatabase.EXPECT().GenFieldID("ns", "test1", "f4", field.SumField).Return(field.ID(1), nil),
		tStore.EXPECT().GetFStore(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, false),
		tStore.EXPECT().InsertFStore(gomock.Any()),
		mockMStore.EXPECT().AddField(gomock.Any(), gomock.Any()),
		mockMStore.EXPECT().SetTimestamp(gomock.Any(), gomock.Any()),
	)
	err = md.Write("ns", "test1", uint32(1), uint32(10), 1564300800000+timeutil.OneHour, []*pb.Field{{
		Name:  "f4",
		Field: &pb.Field_Sum{Sum: &pb.Sum{Value: 10.0}},
	}})
	assert.NoError(t, err)
	assert.Len(t, md.Families(), 2)
	assert.True(t, md.MemSize() > 0)
}

func TestMemoryDatabase_Write_err(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer func() {
		_ = fileutil.RemoveDir(testDBPath)
		defer ctrl.Finish()
	}()
	cfg.TempPath = testDBPath

	// mock
	mockMetadata := metadb.NewMockMetadata(ctrl)
	mockMetadataDatabase := metadb.NewMockMetadataDatabase(ctrl)
	mockMetadata.EXPECT().MetadataDatabase().Return(mockMetadataDatabase).AnyTimes()
	mockIndex := indexdb.NewMockIndexDatabase(ctrl)
	mockMStore := NewMockmStoreINTF(ctrl)
	tStore := NewMocktStoreINTF(ctrl)
	mockMStore.EXPECT().GetOrCreateTStore(uint32(10)).Return(tStore, 10).AnyTimes()
	// build memory-database
	cfg.Index = mockIndex
	cfg.Metadata = mockMetadata
	mdINTF := NewMemoryDatabase(cfg)
	buf := NewMockDataPointBuffer(ctrl)
	buf.EXPECT().AllocPage().Return(nil, fmt.Errorf("err"))
	md := mdINTF.(*memoryDatabase)
	md.buf = buf

	// load mock
	md.mStores.Put(uint32(1), mockMStore)
	// case 1: write ok
	gomock.InOrder(
		mockMetadataDatabase.EXPECT().GenFieldID("ns", "test1", "f1", field.SumField).Return(field.ID(1), nil),
		tStore.EXPECT().GetFStore(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, false),
	)
	err := md.Write("ns", "test1", uint32(1), uint32(10), 1564300800000, []*pb.Field{{
		Name:  "f1",
		Field: &pb.Field_Sum{Sum: &pb.Sum{Value: 10.0}},
	}})
	assert.Error(t, err)
}

func TestMemoryDatabase_FlushFamilyTo(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mdINTF := NewMemoryDatabase(cfg)
	md := mdINTF.(*memoryDatabase)
	flusher := metricsdata.NewMockFlusher(ctrl)
	flusher.EXPECT().Commit().Return(nil).AnyTimes()
	// mock mStore
	mockMStore := NewMockmStoreINTF(ctrl)
	md.mStores.Put(uint32(3333), mockMStore)

	// case 1: flusher ok
	mockMStore.EXPECT().FlushMetricsDataTo(gomock.Any(), gomock.Any()).Return(nil)
	err := md.FlushFamilyTo(flusher, 10)
	assert.NoError(t, err)
	// case 2: flusher err
	mockMStore.EXPECT().FlushMetricsDataTo(gomock.Any(), gomock.Any()).Return(fmt.Errorf("err"))
	err = md.FlushFamilyTo(flusher, 10)
	assert.Error(t, err)
}

func TestMemoryDatabase_Filter(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mdINTF := NewMemoryDatabase(cfg)
	md := mdINTF.(*memoryDatabase)

	// not found
	_, _ = md.Filter(0, []uint16{1}, 1, nil)

	// mock mStore
	mockMStore := NewMockmStoreINTF(ctrl)
	mockMStore.EXPECT().Filter(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, nil)
	md.mStores.Put(uint32(3333), mockMStore)

	_, _ = md.Filter(uint32(3333), []uint16{1}, 1, nil)
}

func TestMemoryDatabase_getFieldValue(t *testing.T) {
	mdINTF := NewMemoryDatabase(cfg)
	md := mdINTF.(*memoryDatabase)
	assert.Equal(t, 10.1, md.getFieldValue(field.SumField, &pb.Field{
		Field: &pb.Field_Sum{Sum: &pb.Sum{
			Value: 10.1,
		}},
	}))
	assert.Equal(t, 10.1, md.getFieldValue(field.MinField, &pb.Field{
		Field: &pb.Field_Min{Min: &pb.Min{
			Value: 10.1,
		}},
	}))
	assert.Equal(t, 10.1, md.getFieldValue(field.MaxField, &pb.Field{
		Field: &pb.Field_Max{Max: &pb.Max{
			Value: 10.1,
		}},
	}))
	assert.Equal(t, 10.1, md.getFieldValue(field.GaugeField, &pb.Field{
		Field: &pb.Field_Gauge{Gauge: &pb.Gauge{
			Value: 10.1,
		}},
	}))
	assert.Equal(t, 0.0, md.getFieldValue(field.Unknown, &pb.Field{
		Field: &pb.Field_Gauge{Gauge: &pb.Gauge{
			Value: 10.1,
		}},
	}))
}
