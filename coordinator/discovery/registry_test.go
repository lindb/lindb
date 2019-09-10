package discovery

import (
	"fmt"
	"testing"
	"time"

	"github.com/lindb/lindb/pkg/timeutil"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/state"
)

var testRegistryPath = "/test/registry"

func TestRegistry(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := state.NewMockRepository(ctrl)

	registry1 := NewRegistry(repo, testRegistryPath, 100)

	closedCh := make(chan state.Closed)

	node := models.Node{IP: "127.0.0.1", Port: 2080}
	nodeMap := models.ActiveNodeMap{
		OnlineTime: timeutil.Now(),
		NodeMap:    map[models.NodeType]*models.Node{models.NodeTypeRPC: &node},
	}
	gomock.InOrder(
		repo.EXPECT().Heartbeat(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Return(nil, fmt.Errorf("err")),
		repo.EXPECT().Heartbeat(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Return(closedCh, nil),
	)
	err := registry1.Register(node, &nodeMap)
	if err != nil {
		t.Fatal(err)
	}
	time.Sleep(600 * time.Millisecond)

	// maybe retry do heartbeat after close chan
	repo.EXPECT().Heartbeat(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, nil).AnyTimes()

	close(closedCh)

	nodePath := fmt.Sprintf("%s/%s", testRegistryPath, node.Indicator())
	repo.EXPECT().Delete(gomock.Any(), nodePath).Return(nil)
	err = registry1.Deregister(node)
	assert.Nil(t, err)

	err = registry1.Close()
	if err != nil {
		t.Fatal(err)
	}

	registry1 = NewRegistry(repo, testRegistryPath, 100)
	err = registry1.Close()
	if err != nil {
		t.Fatal(err)
	}
	r := registry1.(*registry)
	r.register("/data/pant", &nodeMap)

	registry1 = NewRegistry(repo, testRegistryPath, 100)
	r = registry1.(*registry)

	// cancel ctx in timer
	time.AfterFunc(100*time.Millisecond, func() {
		r.cancel()
	})
	r.register("/data/pant", &nodeMap)
}
