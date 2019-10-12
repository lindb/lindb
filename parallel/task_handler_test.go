package parallel

import (
	"context"
	"fmt"
	"io"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/metadata"

	"github.com/lindb/lindb/config"
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/rpc"
	pb "github.com/lindb/lindb/rpc/proto/common"
)

var cfg = config.Query{
	NumOfTasks:    10,
	QueueCapacity: 10,
	Timeout:       10,
}

func TestTaskHandler_Handle(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	dispatcher := NewMockTaskDispatcher(ctrl)
	taskServerFactory := rpc.NewMockTaskServerFactory(ctrl)
	taskServerFactory.EXPECT().Register(gomock.Any(), gomock.Any())
	taskServerFactory.EXPECT().Deregister(gomock.Any())
	handler := NewTaskHandler(cfg, taskServerFactory, dispatcher)

	server := pb.NewMockTaskService_HandleServer(ctrl)
	ctx := metadata.NewOutgoingContext(context.TODO(), metadata.Pairs())
	server.EXPECT().Context().Return(ctx)
	err := handler.Handle(server)
	assert.NotNil(t, err)

	ctx = rpc.CreateIncomingContextWithNode(context.TODO(), models.Node{IP: "1.1.1.1", Port: 9000})
	server.EXPECT().Context().Return(ctx)
	server.EXPECT().Recv().Return(nil, fmt.Errorf("err"))
	server.EXPECT().Recv().Return(nil, nil)
	server.EXPECT().Recv().Return(nil, io.EOF)
	dispatcher.EXPECT().Dispatch(gomock.Any(), gomock.Any()).AnyTimes()
	_ = handler.Handle(server)
}
