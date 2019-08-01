package handler

import (
	"context"
	"io"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/eleme/lindb/models"
	"github.com/eleme/lindb/rpc"
	"github.com/eleme/lindb/rpc/proto/storage"
)

/**
query server
register stream
get stream
deregister stream
*/
func TestQueryServerFactory_Register(t *testing.T) {
	ctl := gomock.NewController(t)

	mockServerStream := storage.NewMockQueryService_QueryServer(ctl)

	node := models.Node{
		IP:   "1.1.1.1",
		Port: 123,
	}

	done := make(chan struct{})

	gomock.InOrder(
		mockServerStream.EXPECT().Context().Return(buildContext(node)),
		mockServerStream.EXPECT().Recv().DoAndReturn(func() (*storage.QueryRequest, error) {
			<-done
			return nil, io.EOF
		}),
	)

	sfct := rpc.NewServerStreamFactory()

	query := NewQuery(sfct)

	go func() {
		if err := query.Query(mockServerStream); err != nil {
			t.Error(err)
		}
	}()

	// wait for query.Query exec
	time.Sleep(10 * time.Millisecond)

	assert.Equal(t, 1, len(sfct.Nodes()))
	assert.Equal(t, node, sfct.Nodes()[0])

	stream, ok := sfct.GetStream(node)
	assert.True(t, ok)

	_, ok = stream.(storage.QueryService_QueryServer)
	assert.True(t, ok)

	close(done)
}

func buildContext(logicNode models.Node) context.Context {
	return rpc.CreateIncomingContextWithNode(context.TODO(), logicNode)
}
