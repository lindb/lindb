package models

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/pkg/option"
)

func TestNewShardAssignment(t *testing.T) {
	shardAssign := NewShardAssignment("test")
	shardAssign.AddReplica(1, 1)
	shardAssign.AddReplica(1, 2)
	shardAssign.AddReplica(2, 3)
	shardAssign.AddReplica(2, 5)
	assert.Equal(t, []int{1, 2}, shardAssign.Shards[1].Replicas)
	assert.Equal(t, []int{3, 5}, shardAssign.Shards[2].Replicas)
}

func TestDatabase_String(t *testing.T) {
	database := Database{
		Name:          "test",
		NumOfShard:    10,
		ReplicaFactor: 1,
		Option:        option.DatabaseOption{Interval: "10s"},
	}
	assert.Equal(t, "create database test with shard 10, replica 1, interval 10s", database.String())
}
