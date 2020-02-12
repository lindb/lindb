package models

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/pkg/option"
)

func TestCreateShardTask_Bytes(t *testing.T) {
	task := CreateShardTask{
		DatabaseName:   "test",
		ShardIDs:       []int32{1, 4, 6},
		DatabaseOption: option.DatabaseOption{},
	}
	data := task.Bytes()
	task1 := CreateShardTask{}
	_ = json.Unmarshal(data, &task1)
	assert.Equal(t, task, task1)
}

func TestDatabaseFlushTask_Bytes(t *testing.T) {
	task := DatabaseFlushTask{
		DatabaseName: "test",
	}
	data := task.Bytes()
	task1 := DatabaseFlushTask{}
	_ = json.Unmarshal(data, &task1)
	assert.Equal(t, task, task1)
}
