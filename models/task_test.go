package models

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/pkg/option"
)

func TestCreateShardTask_Bytes(t *testing.T) {
	task := CreateShardTask{
		Database:    "test",
		ShardIDs:    []int32{1, 4, 6},
		ShardOption: option.ShardOption{TimeWindow: 100},
	}
	data := task.Bytes()
	task1 := CreateShardTask{}
	_ = json.Unmarshal(data, &task1)
	assert.Equal(t, task, task1)
}
