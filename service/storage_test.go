package service

import (
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/pkg/fileutil"
	"github.com/lindb/lindb/pkg/option"
	"github.com/lindb/lindb/tsdb"
)

var testPath = "test_data"
var validOption = option.EngineOption{Interval: "10s"}

func TestCreateShards(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer func() {
		ctrl.Finish()
		_ = fileutil.RemoveDir(testPath)
	}()

	factory := tsdb.NewMockEngineFactory(ctrl)
	engine := tsdb.NewMockEngine(ctrl)

	service := NewStorageService(factory)
	// 2 times for double check
	err := service.CreateShards("test_db", option.EngineOption{})
	assert.NotNil(t, err)

	factory.EXPECT().GetEngine(gomock.Any()).Return(nil).MaxTimes(2)
	factory.EXPECT().CreateEngine(gomock.Any()).Return(engine, nil)
	engine.EXPECT().CreateShards(gomock.Any(), gomock.Any()).Return(nil)
	err = service.CreateShards("test_db", validOption, 1, 2, 3)
	assert.Nil(t, err)

	shard := tsdb.NewMockShard(ctrl)
	factory.EXPECT().GetEngine(gomock.Any()).Return(engine)
	engine.EXPECT().GetShard(int32(1)).Return(shard)
	assert.NotNil(t, service.GetShard("test_db", 1))

	factory.EXPECT().GetEngine(gomock.Any()).Return(engine)
	engine.EXPECT().GetShard(int32(10)).Return(nil)
	assert.Nil(t, service.GetShard("test_db", 10))

	factory.EXPECT().GetEngine(gomock.Any()).Return(nil)
	assert.Nil(t, service.GetShard("not_exist_db", 10))

	// create engine error
	factory.EXPECT().GetEngine(gomock.Any()).Return(nil).MaxTimes(2)
	factory.EXPECT().CreateEngine("engine_err").Return(nil, fmt.Errorf("err"))
	err = service.CreateShards("engine_err", validOption, 1, 2, 3)
	assert.NotNil(t, err)

	// create shard error
	factory.EXPECT().GetEngine(gomock.Any()).Return(engine)
	engine.EXPECT().CreateShards(gomock.Any(), gomock.Any()).Return(fmt.Errorf("err"))
	err = service.CreateShards("test_db", validOption, 5)
	assert.NotNil(t, err)
}
