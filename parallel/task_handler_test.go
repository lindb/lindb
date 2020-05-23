package parallel

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/metadata"

	"github.com/lindb/lindb/config"
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/ltoml"
	"github.com/lindb/lindb/rpc"
	commonmock "github.com/lindb/lindb/rpc/pbmock/common"
	pb "github.com/lindb/lindb/rpc/proto/common"
)

type mockTaskDispatcher struct {
}

func (d *mockTaskDispatcher) Dispatch(ctx context.Context, stream pb.TaskService_HandleServer, req *pb.TaskRequest) {
	panic("err")
}

var cfg = config.Query{
	MaxWorkers:  10,
	IdleTimeout: ltoml.Duration(time.Second * 5),
	Timeout:     ltoml.Duration(time.Second * 10),
}

func TestTaskHandler_Handle(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	dispatcher := NewMockTaskDispatcher(ctrl)
	taskServerFactory := rpc.NewMockTaskServerFactory(ctrl)
	taskServerFactory.EXPECT().Register(gomock.Any(), gomock.Any())
	taskServerFactory.EXPECT().Deregister(gomock.Any(), gomock.Any()).Return(true)
	handler := NewTaskHandler(cfg, taskServerFactory, dispatcher)

	server := commonmock.NewMockTaskService_HandleServer(ctrl)
	ctx := metadata.NewOutgoingContext(context.TODO(), metadata.Pairs())
	server.EXPECT().Context().Return(ctx)
	err := handler.Handle(server)
	assert.NotNil(t, err)

	ctx = rpc.CreateIncomingContextWithNode(context.TODO(), models.Node{IP: "1.1.1.1", Port: 9000})
	server.EXPECT().Context().Return(ctx)
	server.EXPECT().Recv().Return(nil, nil)
	server.EXPECT().Recv().Return(nil, fmt.Errorf("err"))
	dispatcher.EXPECT().Dispatch(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()
	_ = handler.Handle(server)
}

func TestTaskHandler_dispatch(t *testing.T) {
	handler := NewTaskHandler(cfg, nil, &mockTaskDispatcher{})
	// test dispatch panic
	handler.dispatch(nil, nil)
}
