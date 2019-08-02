package service

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/config"
	"github.com/lindb/lindb/pkg/fileutil"
	"github.com/lindb/lindb/pkg/interval"
	"github.com/lindb/lindb/pkg/option"
)

var testPath = "test_data"
var validOption = option.ShardOption{Interval: time.Second * 10, IntervalType: interval.Day}

func TestCreateShards(t *testing.T) {
	defer func() {
		_ = fileutil.RemoveDir(testPath)
	}()

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
