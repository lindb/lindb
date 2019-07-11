package service

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/eleme/lindb/config"
	"github.com/eleme/lindb/pkg/interval"
	"github.com/eleme/lindb/pkg/option"
	"github.com/eleme/lindb/pkg/util"
)

var testPath = "test_data"
var validOption = option.ShardOption{Interval: time.Second * 10, IntervalType: interval.Day}

func TestCreateShards(t *testing.T) {
	defer util.RemoveDir(testPath)

	cfg := config.Engine{
		Path: testPath,
	}
	service := NewStorageService(cfg)
	err := service.CreateShards("test_db", option.ShardOption{})
	assert.NotNil(t, err)

	err = service.CreateShards("test_db", validOption, 1, 2, 3)
	assert.Nil(t, err)

	assert.NotNil(t, service.GetShard("test_db", 1))
	assert.NotNil(t, service.GetShard("test_db", 2))
	assert.NotNil(t, service.GetShard("test_db", 3))
	assert.Nil(t, service.GetShard("test_db", 10))
	assert.Nil(t, service.GetShard("test_db2", 2))
}
