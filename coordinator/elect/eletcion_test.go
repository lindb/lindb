package elect

import (
	"bytes"
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/eleme/lindb/coordinator/discovery"
	"github.com/eleme/lindb/pkg/state"

	etcd "github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/integration"
	"github.com/coreos/pkg/capnslog"
	"github.com/stretchr/testify/assert"
)

func init() {
	capnslog.SetGlobalLogLevel(capnslog.CRITICAL)
}

func TestElect(t *testing.T) {
	clus := integration.NewClusterV3(t, &integration.ClusterConfig{Size: 1})
	defer clus.Terminate(t)

	_ = state.New("etcd", etcd.Config{
		Endpoints: []string{clus.Members[0].GRPCAddr()},
	})
	node := discovery.Node{IP: "127.0.0.1", Port: 2080}
	election := NewElection(node, "test", 1)
	ctx, cancel := context.WithCancel(context.Background())
	// first node election register must be success
	success, err := election.Elect(ctx)
	if err != nil {
		t.Errorf("Elect error :%s", err.Error())
	}
	assert.True(t, success)
	node2 := discovery.Node{IP: "127.0.0.2", Port: 2080}
	// second node election should be false

	election2 := NewElection(node2, "test", 1)
	ctx2, cancel2 := context.WithCancel(context.Background())
	success2, _ := election2.Elect(ctx2)
	assert.False(t, success2)
	isMaster, master, _ := election2.IsMaster()
	assert.False(t, isMaster)

	assert.NotNil(t, master)
	fmt.Println(master)
	assert.True(t, bytes.Equal([]byte(node.IP), []byte(master.IP)))
	cancel()
	// first node exist,the second node should be the master
	time.Sleep(2 * time.Second)
	isMaster, _, _ = election2.IsMaster()
	assert.True(t, isMaster)
	defer cancel2()
}

func TestResign(t *testing.T) {
	clus := integration.NewClusterV3(t, &integration.ClusterConfig{Size: 1})
	defer clus.Terminate(t)

	_ = state.New("etcd", etcd.Config{
		Endpoints: []string{clus.Members[0].GRPCAddr()},
	})
	node := discovery.Node{IP: "127.0.0.1", Port: 2080}
	election := NewElection(node, "test", 1)
	ctx, cancel := context.WithCancel(context.Background())
	success, _ := election.Elect(ctx)
	assert.True(t, success)
	_ = election.Resign(context.TODO())

	isMaster, master, err := election.IsMaster()
	assert.False(t, isMaster)
	assert.NotNil(t, err)
	assert.Nil(t, master)
	defer cancel()
}
