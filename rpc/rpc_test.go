package rpc

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/eleme/lindb/models"
)

func TestClientConnFactory(t *testing.T) {
	node1 := models.Node{
		IP:   "1.1.1.1",
		Port: 123,
	}

	node2 := models.Node{
		IP:   "1.1.1.1",
		Port: 456,
	}

	fct := NewClientConnFactory()

	conn1, err := fct.GetClientConn(node1)
	if err != nil {
		t.Fatal(err)
	}

	conn11, err := fct.GetClientConn(node1)
	if err != nil {
		t.Fatal(err)
	}

	conn2, err := fct.GetClientConn(node2)
	if err != nil {
		t.Fatal(err)
	}

	assert.True(t, conn1 == conn11)
	assert.False(t, conn1 == conn2)

}

func TestContext(t *testing.T) {
	node := models.Node{
		IP:   "1.1.1.1",
		Port: 123,
	}
	ctx := CreateIngoingContext(context.TODO(), "db", 0, node)

	n, err := GetLogicNodeFromContext(ctx)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, n, &node)

	db, err := GetDatabaseFromContext(ctx)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, db, "db")

	shardID, err := GetShardIDFromContext(ctx)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, shardID, int32(0))
}
