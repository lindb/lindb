package rpc

import (
	"fmt"
	"io"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"

	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/rpc/proto/common"
)

type mockTaskHandle struct {
}

func (h *mockTaskHandle) Handle(stream common.TaskService_HandleServer) error {
	return nil
}

func TestTaskServerFactory(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()
	fct := NewTaskServerFactory()

	stream := fct.GetStream((&node).Indicator())
	assert.Nil(t, stream)

	mockServerStream := common.NewMockTaskService_HandleServer(ctl)

	fct.Register((&node).Indicator(), mockServerStream)
	stream = fct.GetStream((&node).Indicator())
	assert.NotNil(t, stream)

	nodes := fct.Nodes()
	assert.Equal(t, 1, len(nodes))
	assert.Equal(t, node, nodes[0])

	fct.Deregister((&node).Indicator())
	// parse node error
	fct.Register("node_err", mockServerStream)
	nodes = fct.Nodes()
	assert.Equal(t, 0, len(nodes))
}

func TestTaskClientFactory(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()
	oldClientConnFct := clientConnFct
	mockClientConnFct := NewMockClientConnFactory(ctl)
	clientConnFct = mockClientConnFct
	grpcServer := NewGRPCServer(":9000")
	common.RegisterTaskServiceServer(grpcServer.GetServer(), &mockTaskHandle{})

	go func() {
		err := grpcServer.Start()
		if err != nil {
			fmt.Print(err)
		}
	}()
	time.Sleep(100 * time.Millisecond)
	defer func() {
		clientConnFct = oldClientConnFct
	}()

	fct := NewTaskClientFactory(models.Node{IP: "127.0.0.1", Port: 123})
	receiver := NewMockTaskReceiver(ctl)
	fct.SetTaskReceiver(receiver)

	target := models.Node{IP: "127.0.0.1", Port: 122}
	mockClientConnFct.EXPECT().GetClientConn(target).Return(nil, fmt.Errorf("err"))
	err := fct.CreateTaskClient(target)
	assert.NotNil(t, err)

	conn, _ := grpc.Dial(target.Indicator(), grpc.WithInsecure())
	mockClientConnFct.EXPECT().GetClientConn(target).Return(conn, nil)
	err = fct.CreateTaskClient(target)
	assert.NotNil(t, err)

	target = models.Node{IP: "127.0.0.1", Port: 9000}
	conn, _ = grpc.Dial(target.Indicator(), grpc.WithInsecure())
	mockClientConnFct.EXPECT().GetClientConn(target).Return(conn, nil)
	err = fct.CreateTaskClient(target)
	assert.NoError(t, err)

	// not create new one if exist
	target = models.Node{IP: "127.0.0.1", Port: 9000}
	err = fct.CreateTaskClient(target)
	assert.NoError(t, err)

	cli := fct.GetTaskClient((&target).Indicator())
	assert.NotNil(t, cli)

	fct.CloseTaskClient((&target).Indicator())

	fct1 := fct.(*taskClientFactory)
	mockTaskClient := common.NewMockTaskService_HandleClient(ctl)
	fct1.taskStreams["mock_client"] = mockTaskClient

	mockTaskClient.EXPECT().CloseSend().Return(fmt.Errorf("err"))
	fct1.CloseTaskClient("mock_client")
}

func TestTaskClientFactory_handler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	receiver := NewMockTaskReceiver(ctrl)
	fct := NewTaskClientFactory(models.Node{IP: "127.0.0.1", Port: 123})
	fct.SetTaskReceiver(receiver)
	cli := common.NewMockTaskService_HandleClient(ctrl)

	factory := fct.(*taskClientFactory)
	gomock.InOrder(
		cli.EXPECT().Recv().Return(nil, fmt.Errorf("err")),
		cli.EXPECT().Recv().Return(nil, nil),
		receiver.EXPECT().Receive(gomock.Any()).Return(nil),
		cli.EXPECT().Recv().Return(nil, nil),
		receiver.EXPECT().Receive(gomock.Any()).Return(fmt.Errorf("err")),
		cli.EXPECT().Recv().Return(nil, io.EOF),
	)
	factory.handleTaskResponse(cli)
}
