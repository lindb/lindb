package parallel

import (
	"context"
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"

	"github.com/lindb/lindb/models"
	commonmock "github.com/lindb/lindb/rpc/pbmock/common"
	pb "github.com/lindb/lindb/rpc/proto/common"
)

func TestLeafTaskDispatcher_Dispatch(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	server := commonmock.NewMockTaskService_HandleServer(ctrl)
	server.EXPECT().Send(gomock.Any()).Return(fmt.Errorf("err"))
	leafTaskDispatcher := NewLeafTaskDispatcher(models.Node{IP: "1.1.1.1", Port: 9000}, nil, nil, nil)
	leafTaskDispatcher.Dispatch(context.TODO(), server, &pb.TaskRequest{PhysicalPlan: []byte{1, 1, 1}})
}

func TestIntermediateTaskDispatcher_Dispatch(t *testing.T) {
	dispatcher := NewIntermediateTaskDispatcher()
	dispatcher.Dispatch(context.TODO(), nil, &pb.TaskRequest{PhysicalPlan: []byte{1, 1, 1}})
}
