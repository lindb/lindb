package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
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
