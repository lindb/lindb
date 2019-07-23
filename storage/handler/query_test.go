package handler

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/eleme/lindb/models"
	"github.com/eleme/lindb/pkg/logger"
	"github.com/eleme/lindb/rpc"
	"github.com/eleme/lindb/rpc/proto/storage"
)

var log = logger.GetLogger("handler/query_test")

/**
client -> connect, recv
server register stream
sever write
*/
func TestQueryStream(t *testing.T) {
	query := NewQuery(rpc.NewServerStreamFactory())

	sv := rpc.NewTCPServer(":9001")
	gs := sv.GetServer()

	storage.RegisterQueryServiceServer(gs, query)

	go func() {
		if err := sv.Start(); err != nil {
			t.Error(err)
		}
	}()

	logicNode := models.Node{
		IP:   "1.1.1.1",
		Port: 123,
	}
	cliFct := rpc.NewClientStreamFactory(logicNode)

	remoteNode := models.Node{
		IP:   "",
		Port: 9001,
	}

	queryCli, err := cliFct.CreateQueryClient(remoteNode)
	if err != nil {
		t.Fatal(err)
	}

	done := make(chan struct{})
	go func() {
		<-done
		s, ok := query.fct.GetStream(logicNode)
		if !ok {
			t.Error("should exists")
			return
		}

		sv := s.(storage.QueryService_QueryServer)

		err = sv.Send(&storage.QueryResponse{
			Msg: "server",
		})

		if err != nil {
			t.Error(err)
			return
		}
		log.Info("server send")
	}()

	err = queryCli.Send(&storage.QueryRequest{
		Msg: "client",
	})
	if err != nil {
		t.Fatal(err)
	}

	log.Info("client send")

	// wait for register connection
	time.Sleep(20 * time.Millisecond)
	close(done)

	resp, err := queryCli.Recv()
	if err != nil {
		t.Fatal(err)
	}

	log.Info("client recv:" + resp.Msg)

	err = queryCli.CloseSend()
	if err != nil {
		t.Fatal(err)
	}

	// wait for DisConn
	time.Sleep(100 * time.Millisecond)

	assert.Equal(t, 0, len(query.fct.Nodes()))

}

func TestMap(t *testing.T) {
	m := make(map[models.Node]string)

	node1 := models.Node{
		IP:   "123",
		Port: 123,
	}

	node2 := models.Node{
		IP:   "123",
		Port: 123,
	}

	m[node1] = "node1"
	m[node2] = "node2"

	assert.Equal(t, len(m), 1)
	assert.Equal(t, m[node1], "node2")

}
