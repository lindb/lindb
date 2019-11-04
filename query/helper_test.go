package query

import (
	"github.com/golang/mock/gomock"

	"github.com/lindb/lindb/series/field"
	"github.com/lindb/lindb/tsdb"
	"github.com/lindb/lindb/tsdb/memdb"
	"github.com/lindb/lindb/tsdb/metadb"
)

///////////////////////////////////////////////////
//                mock interface				 //
///////////////////////////////////////////////////

func newMockDatabase(ctrl *gomock.Controller) *tsdb.MockDatabase {
	shard := tsdb.NewMockShard(ctrl)
	memDB := memdb.NewMockMemoryDatabase(ctrl)
	shard.EXPECT().MemoryDatabase().Return(memDB).AnyTimes()
	shard.EXPECT().GetDataFamilies(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()

	metadataIndex := metadb.NewMockIDGetter(ctrl)
	metadataIndex.EXPECT().GetMetricID(gomock.Any()).Return(uint32(10), nil).AnyTimes()
	metadataIndex.EXPECT().GetFieldID(gomock.Any(), gomock.Any()).Return(uint16(10), field.SumField, nil).AnyTimes()

	mockedDatabase := tsdb.NewMockDatabase(ctrl)
	mockedDatabase.EXPECT().GetShard(gomock.Any()).Return(shard, true).AnyTimes()
	mockedDatabase.EXPECT().IDGetter().Return(metadataIndex).AnyTimes()
	mockedDatabase.EXPECT().NumOfShards().Return(3).AnyTimes()
	return mockedDatabase
}
