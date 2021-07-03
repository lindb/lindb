package broker

import (
	"io"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/rpc"
)

func Test_ConnectionManager(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTaskClientFactory := rpc.NewMockTaskClientFactory(ctrl)
	cm := &connectionManager{
		RoleFrom:          "broker",
		RoleTo:            "broker",
		connections:       make(map[string]struct{}),
		taskClientFactory: mockTaskClientFactory,
	}
	mockTaskClientFactory.EXPECT().CreateTaskClient(gomock.Any()).Return(nil).Times(3)
	cm.createConnection(models.Node{IP: "192.168.1.1", Port: 1000})
	cm.createConnection(models.Node{IP: "192.168.1.2", Port: 2000})
	cm.createConnection(models.Node{IP: "192.168.1.3", Port: 3000})
	assert.Len(t, cm.connections, 3)
	mockTaskClientFactory.EXPECT().CreateTaskClient(gomock.Any()).Return(io.ErrClosedPipe).Times(2)
	cm.createConnection(models.Node{IP: "192.168.1.3", Port: 4000})
	cm.createConnection(models.Node{IP: "192.168.1.3", Port: 3000})
	assert.Len(t, cm.connections, 2)

	mockTaskClientFactory.EXPECT().CloseTaskClient(gomock.Any()).
		Return(true, nil).AnyTimes()
	cm.closeInactiveNodeConnections([]string{
		"192.168.1.1:9000",
		"192.168.1.1:1000",
		"192.168.1.2:2000",
	})
	assert.Len(t, cm.connections, 2)

	cm.closeInactiveNodeConnections([]string{"192.168.1.1:1000"})
	assert.Len(t, cm.connections, 1)

	cm.closeInactiveNodeConnections([]string{})
	assert.Len(t, cm.connections, 0)

}
