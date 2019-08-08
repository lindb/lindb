package context

import (
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"

	"github.com/lindb/lindb/coordinator/database"
	"github.com/lindb/lindb/coordinator/storage"
)

func TestMasterContext_Close(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	stateMachine := &StateMachine{}
	ctx := NewMasterContext(stateMachine)
	ctx.Close()

	admin := database.NewMockAdminStateMachine(ctrl)
	cluster := storage.NewMockClusterStateMachine(ctrl)

	stateMachine.DatabaseAdmin = admin
	stateMachine.StorageCluster = cluster

	admin.EXPECT().Close().Return(nil)
	cluster.EXPECT().Close().Return(nil)
	ctx.Close()

	admin.EXPECT().Close().Return(fmt.Errorf("err"))
	cluster.EXPECT().Close().Return(fmt.Errorf("err"))
	ctx.Close()
}
