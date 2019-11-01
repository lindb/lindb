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

func MockTSDBEngine(ctrl *gomock.Controller) *tsdb.MockEngine {
	shard := tsdb.NewMockShard(ctrl)
	memDB := memdb.NewMockMemoryDatabase(ctrl)
	shard.EXPECT().GetMemoryDatabase().Return(memDB).AnyTimes()
	shard.EXPECT().GetDataFamilies(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()

	metadataIndex := metadb.NewMockIDGetter(ctrl)
	metadataIndex.EXPECT().GetMetricID(gomock.Any()).Return(uint32(10), nil).AnyTimes()
	metadataIndex.EXPECT().GetFieldID(gomock.Any(), gomock.Any()).Return(uint16(10), field.SumField, nil).AnyTimes()

	engine := tsdb.NewMockEngine(ctrl)
	engine.EXPECT().GetShard(gomock.Any()).Return(shard).AnyTimes()
	engine.EXPECT().GetIDGetter().Return(metadataIndex).AnyTimes()
	engine.EXPECT().NumOfShards().Return(3).AnyTimes()
	return engine
}
