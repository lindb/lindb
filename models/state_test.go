package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStorageState(t *testing.T) {
	storageState := NewStorageState()
	storageState.AddActiveNode(&ActiveNode{Node: Node{IP: "1.1.1.1", Port: 9000}})
	storageState.AddActiveNode(&ActiveNode{Node: Node{IP: "1.1.1.2", Port: 9000}})
	storageState.AddActiveNode(&ActiveNode{Node: Node{IP: "1.1.1.3", Port: 9000}})
	assert.Equal(t, 3, len(storageState.GetActiveNodes()))
	storageState.RemoveActiveNode("1.1.1.2:9000")
	assert.Equal(t, 2, len(storageState.GetActiveNodes()))
}
