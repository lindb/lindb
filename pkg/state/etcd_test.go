// Licensed to LinDB under one or more contributor
// license agreements. See the NOTICE file distributed with
// this work for additional information regarding copyright
// ownership. LinDB licenses this file to you under
// the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

package state

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/lindb/lindb/config"
	"github.com/lindb/lindb/mock"
	"github.com/lindb/lindb/pkg/hostutil"

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
	var rep, err = newEtcdRepository(config.RepoState{
		Endpoints: ts.Cluster.Endpoints,
	}, "nobody")
	repo := rep.(*etcdRepository)
	repo.timeout = time.Second * 10
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
	var rep, err = newEtcdRepository(config.RepoState{
		Namespace: "/test/list",
		Endpoints: ts.Cluster.Endpoints,
	}, "nobody")
	repo := rep.(*etcdRepository)
	repo.timeout = time.Second * 10

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
	_, err := newEtcdRepository(config.RepoState{}, "nobody")
	c.Assert(err, check.NotNil)
}

func (ts *testEtcdRepoSuite) TestHeartBeat(c *check.C) {
	b, err := newEtcdRepository(config.RepoState{
		Endpoints: ts.Cluster.Endpoints,
	}, "nobody")
	repo := b.(*etcdRepository)
	repo.timeout = time.Second * 10
	if err != nil {
		c.Fatal(err)
	}
	ip, _ := hostutil.GetHostIP()
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
	b, _ := newEtcdRepository(config.RepoState{
		Endpoints: ts.Cluster.Endpoints,
	}, "nobody")
	ctx, cancel := context.WithCancel(context.Background())
	repo := b.(*etcdRepository)
	repo.timeout = time.Second * 10
	// test watch no exist path
	ch := b.Watch(ctx, "/cluster1/controller/1", true)
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
	ch2 := b.Watch(ctx, "/cluster1/controller/2", true)
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
	b, _ := newEtcdRepository(config.RepoState{
		Endpoints: ts.Cluster.Endpoints,
	}, "nobody")
	repo := b.(*etcdRepository)
	repo.timeout = time.Second * 10
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	_ = b.Put(context.TODO(), "/lindb/broker/1", []byte("1"))
	_ = b.Put(context.TODO(), "/lindb/broker/2", []byte("2"))
	ch := b.WatchPrefix(ctx, "/lindb/broker", true)
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

func (ts *testEtcdRepoSuite) TestElect(c *check.C) {
	b, _ := newEtcdRepository(config.RepoState{
		Endpoints: ts.Cluster.Endpoints,
	}, "nobody")
	repo := b.(*etcdRepository)
	repo.timeout = time.Second * 10

	ctx, cancel := context.WithCancel(context.Background())
	// the key should not exist,it must be success
	success, ch, err := b.Elect(ctx, "/lindb/breoker/master", []byte("test"), 1)
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
	shouldFalse, _, _ := b.Elect(ctx2, "/lindb/breoker/master", []byte("test2"), 1)
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
	shouldSuccess, _, _ := b.Elect(ctx3, "/lindb/breoker/master", []byte("test3"), 1)
	c.Assert(shouldSuccess, check.Equals, true)

	bytes3, _ := b.Get(context.TODO(), "/lindb/breoker/master")

	c.Assert(string(bytes3), check.Equals, "test3")

	cancel3()
}

func (ts *testEtcdRepoSuite) TestBatch(c *check.C) {
	b, _ := newEtcdRepository(config.RepoState{
		Namespace: "/test/batch",
		Endpoints: ts.Cluster.Endpoints,
	}, "nobody")
	repo := b.(*etcdRepository)
	repo.timeout = time.Second * 10
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
	b, _ := newEtcdRepository(config.RepoState{
		Namespace: "/test/batch",
		Endpoints: ts.Cluster.Endpoints,
	}, "nobody")
	repo := b.(*etcdRepository)
	repo.timeout = time.Second * 10

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

func (ts *testEtcdRepoSuite) TestNextSequence(c *check.C) {
	b, _ := newEtcdRepository(config.RepoState{
		Endpoints: ts.Cluster.Endpoints,
	}, "nobody")
	repo := b.(*etcdRepository)
	repo.timeout = time.Second * 10

	seq, err := repo.NextSequence(context.TODO(), "/test/seq")
	c.Assert(err, check.IsNil)
	c.Assert(int64(1), check.Equals, seq)
	for i := 2; i < 10; i++ {
		seq, err = repo.NextSequence(context.TODO(), "/test/seq")
		c.Assert(err, check.IsNil)
		c.Assert(int64(i), check.Equals, seq)
	}
}

func (ts *testEtcdRepoSuite) TestNextSequence_Concurrency(c *check.C) {
	b, _ := newEtcdRepository(config.RepoState{
		Endpoints: ts.Cluster.Endpoints,
	}, "nobody")
	repo := b.(*etcdRepository)
	repo.timeout = time.Minute

	var wait sync.WaitGroup
	wait.Add(10)
	var rs sync.Map
	for i := 0; i < 10; i++ {
		go func() {
			defer wait.Done()
			for j := 0; j < 10; j++ {
				seq, err := repo.NextSequence(context.TODO(), "/test/seq")
				c.Assert(err, check.IsNil)
				rs.Store(seq, seq)
			}
		}()
	}
	wait.Wait()
	for i := 1; i <= 100; i++ {
		_, ok := rs.Load(int64(i))
		c.Assert(ok, check.Equals, ok)
	}
}
