package service

import (
	"fmt"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/config"
	"github.com/lindb/lindb/pkg/fileutil"
	"github.com/lindb/lindb/pkg/interval"
	"github.com/lindb/lindb/pkg/option"
	"github.com/lindb/lindb/tsdb"
)

var testPath = "test_data"
var validOption = option.ShardOption{Interval: time.Second * 10, IntervalType: interval.Day}

func TestCreateShards(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer func() {
		ctrl.Finish()
		_ = fileutil.RemoveDir(testPath)
	}()

	factory := tsdb.NewMockEngineFactory(ctrl)
	engine := tsdb.NewMockEngine(ctrl)

	cfg := config.Engine{
		Path: testPath,
	}
	service := NewStorageService(cfg, factory)

	factory.EXPECT().CreateEngine(gomock.Any(), gomock.Any()).Return(engine, nil)
	err := service.CreateShards("test_db", option.ShardOption{})
	assert.NotNil(t, err)

	engine.EXPECT().CreateShards(gomock.Any(), gomock.Any()).Return(nil)
	err = service.CreateShards("test_db", validOption, 1, 2, 3)
	assert.Nil(t, err)

	shard := tsdb.NewMockShard(ctrl)
	engine.EXPECT().GetShard(int32(1)).Return(shard)
	assert.NotNil(t, service.GetShard("test_db", 1))

	engine.EXPECT().GetShard(int32(10)).Return(nil)
	assert.Nil(t, service.GetShard("test_db", 10))
	assert.Nil(t, service.GetShard("not_exist_db", 10))

	// create engine error
	factory.EXPECT().CreateEngine("engine_err", gomock.Any()).Return(nil, fmt.Errorf("err"))
	err = service.CreateShards("engine_err", validOption, 1, 2, 3)
	assert.NotNil(t, err)

	// create shard error
	engine.EXPECT().CreateShards(gomock.Any(), gomock.Any()).Return(fmt.Errorf("err"))
	err = service.CreateShards("test_db", validOption, 5)
	assert.NotNil(t, err)
}
