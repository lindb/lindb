package state

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/coreos/etcd/integration"
	"github.com/stretchr/testify/assert"
	etcd "go.etcd.io/etcd/clientv3"
)

type address struct {
	Home string
}

func Test_Write_Read(t *testing.T) {
	clus := integration.NewClusterV3(t, &integration.ClusterConfig{Size: 1})
	defer clus.Terminate(t)
	cfg := etcd.Config{
		Endpoints: []string{clus.Members[0].GRPCAddr()},
	}

	var rep, err = newETCDRepository(cfg)
	if err != nil {
		t.Fatal(err)
	}

	home1 := &address{
		Home: "home1",
	}

	d, _ := json.Marshal(home1)
	err = rep.Put(context.TODO(), "/test/key1", d)
	if err != nil {
		t.Fatal(err)
	}
	d1, err1 := rep.Get(context.TODO(), "/test/key1")
	if err1 != nil {
		t.Fatal(err1)
	}
	home2 := &address{}
	json.Unmarshal(d1, home2)
	assert.Equal(t, *home1, *home2)

	rep.Delete(context.TODO(), "/test/key1")

	_, err2 := rep.Get(context.TODO(), "/test/key1")
	assert.NotNil(t, err2)

	rep.Close()
}

func TestNew(t *testing.T) {
	_, err := newETCDRepository("error type")
	assert.NotNil(t, err)
}
