package state

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/eleme/lindb/pkg/util"

	etcd "github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/integration"
	"github.com/stretchr/testify/assert"
	"gopkg.in/check.v1"
)

type address struct {
	Home string
}

type testEtcdRepoSuite struct {
}

var _ = check.Suite(&testEtcdRepoSuite{})

var cluster *integration.ClusterV3
var (
	test *testing.T
	cfg  etcd.Config
)

func Test(t *testing.T) {
	test = t
	check.TestingT(t)
}

func (ts *testEtcdRepoSuite) SetUpSuite(c *check.C) {
	cluster = integration.NewClusterV3(test, &integration.ClusterConfig{Size: 1})
	cfg = etcd.Config{
		Endpoints: []string{cluster.Members[0].GRPCAddr()},
	}
}

func (ts *testEtcdRepoSuite) TearDownSuite(c *check.C) {
	cluster.Terminate(test)
}

func (ts *testEtcdRepoSuite) Test_Write_Read(c *check.C) {
	var rep, err = newEtedRepository(cfg)
	if err != nil {
		c.Fatal(err)
	}

	home1 := &address{
		Home: "home1",
	}

	d, _ := json.Marshal(home1)
	err = rep.Put(context.TODO(), "/test/key1", d)
	if err != nil {
		c.Fatal(err)
	}
	d1, err1 := rep.Get(context.TODO(), "/test/key1")
	if err1 != nil {
		c.Fatal(err1)
	}
	home2 := &address{}
	_ = json.Unmarshal(d1, home2)
	c.Assert(*home1, check.Equals, *home2)

	_ = rep.Delete(context.TODO(), "/test/key1")

	_, err2 := rep.Get(context.TODO(), "/test/key1")
	c.Assert(err2, check.NotNil)

	_ = rep.Close()

}

func (ts *testEtcdRepoSuite) TestNew(c *check.C) {
	_, err := newEtedRepository("error type")
	c.Assert(err, check.NotNil)
}

func (ts *testEtcdRepoSuite) TestHeartBeat(c *check.C) {
	b, err := newEtedRepository(cfg)
	if err != nil {
		c.Fatal(err)
	}
	ip, _ := util.GetHostIP()
	heartbeat := fmt.Sprintf("/cluster1/storage/heartbeat/%s:%d", ip, 2918)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	var ch <-chan Closed
	ch, err = b.Heartbeat(ctx, heartbeat, []byte("test"), 1)
	if err != nil {
		c.Fatal(err)
	}
	_, err = b.Get(ctx, heartbeat)
	if err != nil {
		c.Fatal(err)
	}
	cancel()
	time.Sleep(time.Second)
	_, err = b.Get(ctx, heartbeat)
	if err == nil {
		c.Fatal("heartbeat should be deleted automatically")
	}
	select {
	case <-ch:
	case <-time.After(500 * time.Millisecond):
		c.Fatal("cancel heartbeat timeout")
	}
}

func (ts *testEtcdRepoSuite) TestWatch(c *check.C) {
	b, _ := newEtedRepository(cfg)
	ctx, cancel := context.WithCancel(context.Background())

	// test watch no exist path
	ch, err := b.Watch(ctx, "/cluster1/controller/1")
	if err != nil {
		c.Fatal(err)
	}
	c.Assert(ch, check.NotNil)
	var count int32
	var mutex sync.RWMutex
	val := make(map[string]string)
	go func() {
		for event := range ch {
			mutex.Lock()
			val[event.Key] = string(event.Value)
			atomic.AddInt32(&count, 1)
			mutex.Unlock()
		}
	}()
	_ = b.Put(ctx, "/cluster1/controller/1", []byte("1"))

	// test watch exist path
	_ = b.Put(ctx, "/cluster1/controller/2", []byte("2"))
	ch2, err2 := b.Watch(ctx, "/cluster1/controller/2")
	if err2 != nil {
		c.Fatal(err2)
	}
	go func() {
		for event := range ch2 {
			mutex.Lock()
			val[event.Key] = string(event.Value)
			atomic.AddInt32(&count, 1)
			mutex.Unlock()
		}
	}()
	// modify value of key2
	_ = b.Put(ctx, "/cluster1/controller/2", []byte("222"))
	time.Sleep(100 * time.Millisecond)
	cancel()

	// check watch trigger count
	c.Assert(atomic.AddInt32(&count, 0), check.Equals, int32(3))

	mutex.RLock()
	c.Assert(len(val), check.Equals, 2)
	c.Assert(val["/cluster1/controller/1"], check.Equals, "1")
	c.Assert(val["/cluster1/controller/2"], check.Equals, "222")
	mutex.RUnlock()
}

func (ts *testEtcdRepoSuite) TestGetWatchPrefix(c *check.C) {
	b, _ := newEtedRepository(cfg)
	ctx, cancel := context.WithCancel(context.Background())

	_ = b.Put(context.TODO(), "/lindb/broker/1", []byte("1"))

	_ = b.Put(context.TODO(), "/lindb/broker/2", []byte("2"))

	ch, err := b.Watch(ctx, "/lindb/broker")
	if err != nil {
		c.Fatal(err)
	}

	go func() {
		for event := range ch {
			println(event.Key + string(event.Value))
			c.Assert(event.Key, check.Equals, "/lindb/broker/3")
			c.Assert(string(event.Value), check.Equals, "3")
			switch event.Type {
			case EventTypeModify:

			case EventTypeDelete:
			}
		}
	}()

	_ = b.Put(context.TODO(), "/lindb/broker/3", []byte("3"))
	bytes1, _ := b.Get(context.TODO(), "/lindb/broker/3")
	c.Assert(string(bytes1), check.Equals, "3")
	time.Sleep(10 * time.Second)
	cancel()
}

func (ts *testEtcdRepoSuite) TestPutIfNotExitAndKeepLease(c *check.C) {
	b, _ := newEtedRepository(cfg)
	ctx, cancel := context.WithCancel(context.Background())
	// the key should not exist,it must be success
	success, ch, err := b.PutIfNotExist(ctx, "/lindb/breoker/master", []byte("test"), 1)
	if err != nil {
		c.Fatal(err)
	}
	c.Assert(success, check.Equals, true)
	time.Sleep(2 * time.Second)
	bytes, err := b.Get(context.TODO(), "/lindb/breoker/master")
	if err != nil {
		c.Fatal(err)
	}
	c.Assert(string(bytes), check.Equals, "test")

	ctx2, cancel2 := context.WithCancel(context.Background())
	shouldFalse, _, _ := b.PutIfNotExist(ctx2, "/lindb/breoker/master", []byte("test2"), 1)
	if cancel2 != nil {
		cancel2()
	}
	c.Assert(shouldFalse, check.Equals, false)
	cancel()
	select {
	case <-ch:
	case <-time.After(500 * time.Millisecond):
		c.Fatal("cancel heartbeat timeout")
	}

	time.Sleep(2 * time.Second)

	_, err = b.Get(context.TODO(), "/lindb/breoker/master")
	if err == nil {
		c.Fatal("the key should not exist")
	}

	ctx3, cancel3 := context.WithCancel(context.Background())
	shouldSuccess, _, _ := b.PutIfNotExist(ctx3, "/lindb/breoker/master", []byte("test3"), 1)
	c.Assert(shouldSuccess, check.Equals, true)

	bytes3, _ := b.Get(context.TODO(), "/lindb/breoker/master")

	c.Assert(string(bytes3), check.Equals, "test3")

	cancel3()
}
func (ts *testEtcdRepoSuite) TestDeleteWithTheValue(c *check.C) {
	b, _ := newEtedRepository(cfg)
	_ = b.Put(context.TODO(), "test", []byte("value1"))
	success, err := b.DeleteWithValue(context.TODO(), "test", []byte("value1"))
	if err != nil {
		c.Fatal("the operation should be success")
	}
	assert.True(c, success)

	_ = b.Put(context.TODO(), "test", []byte("value2"))
	success2, err := b.DeleteWithValue(context.TODO(), "test", []byte("value1"))
	if err != nil {
		c.Fatal("the operation should not have error")
	}
	assert.False(c, success2)

}
