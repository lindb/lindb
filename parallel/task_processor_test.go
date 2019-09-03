package parallel

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"

	"github.com/lindb/lindb/models"
	pb "github.com/lindb/lindb/rpc/proto/common"
)

func TestLeafTaskDispatcher_Dispatch(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	leafTaskDispatcher := NewLeafTaskDispatcher(models.Node{IP: "1.1.1.1", Port: 9000}, nil, nil, nil)
	leafTaskDispatcher.Dispatch(context.TODO(), &pb.TaskRequest{PhysicalPlan: []byte{1, 1, 1}})
}

func TestIntermediateTaskDispatcher_Dispatch(t *testing.T) {
	dispatcher := NewIntermediateTaskDispatcher()
	dispatcher.Dispatch(context.TODO(), &pb.TaskRequest{PhysicalPlan: []byte{1, 1, 1}})
}
