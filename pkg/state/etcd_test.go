package state

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/lindb/lindb/mock"
	"github.com/lindb/lindb/pkg/util"

	"gopkg.in/check.v1"
)

type address struct {
	Home string
}

type testEtcdRepoSuite struct {
	mock.RepoTestSuite
}

func TestETCDRepo(t *testing.T) {
	check.Suite(&testEtcdRepoSuite{})
	check.TestingT(t)
}

func (ts *testEtcdRepoSuite) Test_Write_Read(c *check.C) {
	var rep, err = newEtedRepository(Config{
		Endpoints: ts.Cluster.Endpoints,
	})
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

func (ts *testEtcdRepoSuite) TestList(c *check.C) {
	var rep, err = newEtedRepository(Config{
		Namespace: "/test/list",
		Endpoints: ts.Cluster.Endpoints,
	})
	if err != nil {
		c.Fatal(err)
	}

	home1 := &address{
		Home: "home1",
	}

	d, _ := json.Marshal(home1)
	_ = rep.Put(context.TODO(), "/test/key1", d)
	_ = rep.Put(context.TODO(), "/test/key2", d)
	// value is empty, will ignore
	_ = rep.Put(context.TODO(), "/test/key3", []byte{})
	list, _ := rep.List(context.TODO(), "/test")

	c.Assert(2, check.Equals, len(list))
}

func (ts *testEtcdRepoSuite) TestNew(c *check.C) {
	_, err := newEtedRepository(Config{})
	c.Assert(err, check.NotNil)
}

func (ts *testEtcdRepoSuite) TestHeartBeat(c *check.C) {
	b, err := newEtedRepository(Config{
		Endpoints: ts.Cluster.Endpoints,
	})
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
	b, _ := newEtedRepository(Config{
		Endpoints: ts.Cluster.Endpoints,
	})
	ctx, cancel := context.WithCancel(context.Background())

	// test watch no exist path
	ch := b.Watch(ctx, "/cluster1/controller/1")
	c.Assert(ch, check.NotNil)
	var wg sync.WaitGroup
	var mutex sync.RWMutex
	val := make(map[string]string)
	syncKVs := func(ch WatchEventChan) {
		for event := range ch {
			if event.Err != nil {
				continue
			}
			mutex.Lock()
			for _, kv := range event.KeyValues {
				val[kv.Key] = string(kv.Value)
			}
			mutex.Unlock()
		}
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		syncKVs(ch)
	}()
	_ = b.Put(ctx, "/cluster1/controller/1", []byte("1"))

	// test watch exist path
	_ = b.Put(ctx, "/cluster1/controller/2", []byte("2"))
	ch2 := b.Watch(ctx, "/cluster1/controller/2")
	wg.Add(1)
	go func() {
		defer wg.Done()
		syncKVs(ch2)
	}()

	// modify value of key2
	_ = b.Put(ctx, "/cluster1/controller/2", []byte("222"))
	time.Sleep(200 * time.Millisecond)
	cancel()
	wg.Wait()

	// check watch trigger count
	c.Assert(len(val), check.Equals, 2)
	c.Assert(val["/cluster1/controller/1"], check.Equals, "1")
	c.Assert(val["/cluster1/controller/2"], check.Equals, "222")
}

func (ts *testEtcdRepoSuite) TestGetWatchPrefix(c *check.C) {
	b, _ := newEtedRepository(Config{
		Endpoints: ts.Cluster.Endpoints,
	})
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	_ = b.Put(context.TODO(), "/lindb/broker/1", []byte("1"))
	_ = b.Put(context.TODO(), "/lindb/broker/2", []byte("2"))
	ch := b.WatchPrefix(ctx, "/lindb/broker")
	time.Sleep(100 * time.Millisecond)

	_ = b.Put(context.TODO(), "/lindb/broker/3", []byte("3"))
	bytes1, _ := b.Get(context.TODO(), "/lindb/broker/3")
	c.Assert(string(bytes1), check.Equals, "3")
	_ = b.Delete(context.TODO(), "/lindb/broker/3")
	time.Sleep(time.Second)

	var allEvt, modifyEvt, deleteEvt bool
	for event := range ch {
		if event.Err != nil {
			continue
		}
		kvs := map[string]string{}
		for _, kv := range event.KeyValues {
			kvs[kv.Key] = string(kv.Value)
		}
		switch event.Type {
		case EventTypeAll:
			c.Assert(allEvt, check.Equals, false)
			c.Assert(modifyEvt, check.Equals, false)
			c.Assert(deleteEvt, check.Equals, false)
			c.Assert(len(kvs), check.Equals, 2)
			c.Assert(kvs["/lindb/broker/1"], check.Equals, "1")
			c.Assert(kvs["/lindb/broker/2"], check.Equals, "2")
			allEvt = true
		case EventTypeModify:
			c.Assert(allEvt, check.Equals, true)
			c.Assert(modifyEvt, check.Equals, false)
			c.Assert(deleteEvt, check.Equals, false)
			c.Assert(len(kvs), check.Equals, 1)
			c.Assert(kvs["/lindb/broker/3"], check.Equals, "3")
			modifyEvt = true
		case EventTypeDelete:
			c.Assert(allEvt, check.Equals, true)
			c.Assert(modifyEvt, check.Equals, true)
			c.Assert(deleteEvt, check.Equals, false)
			c.Assert(len(kvs), check.Equals, 1)
			_, ok := kvs["/lindb/broker/3"]
			c.Assert(ok, check.Equals, true)
			deleteEvt = true
			cancel()
		}
	}
	c.Assert(deleteEvt, check.Equals, true)
}

func (ts *testEtcdRepoSuite) TestPutIfNotExitAndKeepLease(c *check.C) {
	b, _ := newEtedRepository(Config{
		Endpoints: ts.Cluster.Endpoints,
	})
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

func (ts *testEtcdRepoSuite) TestBatch(c *check.C) {
	b, _ := newEtedRepository(Config{
		Namespace: "/test/batch",
		Endpoints: ts.Cluster.Endpoints,
	})
	batch := Batch{
		KVs: []KeyValue{
			{"key1", []byte("value1")},
			{"key2", []byte("value2")},
			{"key3", []byte("value3")},
		}}
	success, _ := b.Batch(context.TODO(), batch)
	c.Assert(true, check.Equals, success)

	list, _ := b.List(context.TODO(), "key")
	c.Assert(3, check.Equals, len(list))
}

func (ts *testEtcdRepoSuite) TestTransaction(c *check.C) {
	b, _ := newEtedRepository(Config{
		Namespace: "/test/batch",
		Endpoints: ts.Cluster.Endpoints,
	})

	txn := b.NewTransaction()
	txn.Put("test", []byte("value"))
	err := b.Commit(context.TODO(), txn)
	if err != nil {
		c.Fatal(err)
	}

	v, _ := b.Get(context.TODO(), "test")
	c.Assert([]byte("value"), check.DeepEquals, v)

	txn = b.NewTransaction()
	txn.ModRevisionCmp("key", "=", 0)
	txn.Put("test", []byte("value2"))
	err = b.Commit(context.TODO(), txn)
	if err != nil {
		c.Fatal(err)
	}
	v, _ = b.Get(context.TODO(), "test")
	c.Assert([]byte("value2"), check.DeepEquals, v)

	txn = b.NewTransaction()
	txn.ModRevisionCmp("key", "=", 33)
	txn.Delete("test")
	err = b.Commit(context.TODO(), txn)
	c.Assert(err, check.NotNil)

	v, _ = b.Get(context.TODO(), "test")
	c.Assert([]byte("value2"), check.DeepEquals, v)

	txn = b.NewTransaction()
	txn.ModRevisionCmp("key", "=", 0)
	txn.Delete("test")
	err = b.Commit(context.TODO(), txn)
	if err != nil {
		c.Fatal(err)
	}
	_, err = b.Get(context.TODO(), "test")
	c.Assert(err, check.NotNil)

	c.Assert(TxnErr(nil, fmt.Errorf("err")), check.NotNil)
}
