package parallel

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	pb "github.com/lindb/lindb/rpc/proto/common"
)

func TestTaskSenderManager_ClientStream(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	taskSender := NewTaskSenderManager()
	taskSender.AddClientStream("nil", nil)
	assert.Nil(t, taskSender.GetClientStream("nil"))

	taskHandler := pb.NewMockTaskService_HandleClient(ctrl)
	assert.Nil(t, taskSender.GetClientStream("key1"))
	taskSender.AddClientStream("key1", taskHandler)
	assert.Equal(t, taskHandler, taskSender.GetClientStream("key1"))
	taskSender.RemoveClientStream("key1")
	assert.Nil(t, taskSender.GetClientStream("key1"))
}

func TestTaskSenderManager_ServerStream(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	taskSender := NewTaskSenderManager()
	taskSender.AddServerStream("nil", nil)
	assert.Nil(t, taskSender.GetServerStream("nil"))

	taskHandler := pb.NewMockTaskService_HandleServer(ctrl)
	assert.Nil(t, taskSender.GetServerStream("key1"))
	taskSender.AddServerStream("key1", taskHandler)
	assert.Equal(t, taskHandler, taskSender.GetServerStream("key1"))
	taskSender.RemoveServerStream("key1")
	assert.Nil(t, taskSender.GetServerStream("key1"))
}
