package rpc

import (
	"context"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/rpc/proto/common"
	"github.com/lindb/lindb/rpc/proto/storage"
)

var (
	node = models.Node{
		IP:   "127.0.0.1",
		Port: 123,
	}
	database = "database"
	shardID  = int32(0)
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

	fct := GetClientConnFactory()

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
	ctx := CreateIncomingContext(context.TODO(), database, shardID, node)

	n, err := GetLogicNodeFromContext(ctx)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, n, &node)

	db, err := GetDatabaseFromContext(ctx)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, db, database)

	sID, err := GetShardIDFromContext(ctx)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, shardID, sID)
}

func TestClientStreamFactory(t *testing.T) {
	target := models.Node{
		IP:   "127.0.0.1",
		Port: 1234,
	}
	fct := NewClientStreamFactory(node)
	_, err := fct.CreateWriteServiceClient(target)
	assert.Nil(t, err)

	assert.Equal(t, fct.LogicNode(), node)

	// stream client will dail the target address, it's no easy to test
}

func TestServerStreamFactory(t *testing.T) {
	fct := GetServerStreamFactory()

	_, ok := fct.GetStream(node)
	assert.False(t, ok)

	ctl := gomock.NewController(t)
	mockServerStream := storage.NewMockWriteService_WriteServer(ctl)

	fct.Register(node, mockServerStream)
	ss, ok := fct.GetStream(node)
	assert.True(t, ok)

	_, ok = ss.(storage.WriteService_WriteServer)
	assert.True(t, ok)
	assert.Equal(t, ss, mockServerStream)

	nodes := fct.Nodes()
	assert.Equal(t, 1, len(nodes))
	assert.Equal(t, node, nodes[0])

	fct.Deregister(node)
	nodes = fct.Nodes()
	assert.Equal(t, 0, len(nodes))
}

func TestClientStreamFactory_CreateTaskClient(t *testing.T) {
	ctrl := gomock.NewController(t)
	go ctrl.Finish()

	handler := common.NewMockTaskServiceServer(ctrl)

	factory := NewClientStreamFactory(models.Node{IP: "127.0.0.2", Port: 9000})
	target := models.Node{IP: "127.0.0.1", Port: 9000}

	client, err := factory.CreateTaskClient(target)
	assert.NotNil(t, err)
	assert.Nil(t, client)

	server := NewTCPServer(":9000")
	common.RegisterTaskServiceServer(server.GetServer(), handler)
	go func() {
		_ = server.Start()
	}()

	// wait server start finish
	time.Sleep(10 * time.Millisecond)

	_, _ = factory.CreateTaskClient(target)

	time.Sleep(10 * time.Millisecond)
	server.Stop()
}
