package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReplicaState(t *testing.T) {
	replica := ReplicaState{
		Database:     "db",
		ShardID:      int32(1),
		Pending:      100,
		ReplicaIndex: 50,
	}

	assert.Equal(t, "db/1", replica.ShardIndicator())
}
