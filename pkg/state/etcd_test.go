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
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/lindb/lindb/config"
	"github.com/lindb/lindb/internal/mock"
	"github.com/lindb/lindb/pkg/encoding"
	"github.com/lindb/lindb/pkg/hostutil"

	"github.com/stretchr/testify/assert"
)

type address struct {
	Home string
}

func Test_Write_Read(t *testing.T) {
	cluster := mock.StartEtcdCluster(t, "http://localhost:8700")
	defer cluster.Terminate(t)

	var rep, err = newEtcdRepository(&config.RepoState{
		Endpoints: cluster.Endpoints,
	}, "nobody")
	assert.NoError(t, err)

	repo := rep.(*etcdRepository)
	repo.timeout = time.Second * 10

	home1 := &address{
		Home: "home1",
	}

	d := encoding.JSONMarshal(home1)
	err = rep.Put(context.TODO(), "/test/key1", d)
	assert.NoError(t, err)

	d1, err := rep.Get(context.TODO(), "/test/key1")
	assert.NoError(t, err)

	home2 := &address{}
	_ = encoding.JSONUnmarshal(d1, home2)
	assert.Equal(t, *home1, *home2)

	_ = rep.Delete(context.TODO(), "/test/key1")

	_, err2 := rep.Get(context.TODO(), "/test/key1")
	assert.Error(t, err2)

	_ = rep.Close()
}

func TestList(t *testing.T) {
	cluster := mock.StartEtcdCluster(t, "http://localhost:8701")
	defer cluster.Terminate(t)

	var rep, err = newEtcdRepository(&config.RepoState{
		Namespace: "/test/list",
		Endpoints: cluster.Endpoints,
	}, "nobody")
	assert.NoError(t, err)

	repo := rep.(*etcdRepository)
	repo.timeout = time.Second * 10

	home1 := &address{
		Home: "home1",
	}

	d := encoding.JSONMarshal(home1)
	_ = rep.Put(context.TODO(), "/test/key1", d)
	_ = rep.Put(context.TODO(), "/test/key2", d)
	// value is empty, will ignore
	_ = rep.Put(context.TODO(), "/test/key3", []byte{})
	list, err := rep.List(context.TODO(), "/test")

	assert.NoError(t, err)
	assert.Len(t, list, 2)
}

func TestWalkEntry(t *testing.T) {
	cluster := mock.StartEtcdCluster(t, "http://localhost:8701")
	defer cluster.Terminate(t)
	var rep, err = newEtcdRepository(&config.RepoState{
		Namespace: "/test/list",
		Endpoints: cluster.Endpoints,
	}, "nobody")
	assert.NoError(t, err)

	repo := rep.(*etcdRepository)
	repo.timeout = time.Second * 10

	home1 := &address{
		Home: "home1",
	}

	d := encoding.JSONMarshal(home1)
	_ = rep.Put(context.TODO(), "/test/key1", d)
	_ = rep.Put(context.TODO(), "/test/key2", d)
	// value is empty, will ignore
	_ = rep.Put(context.TODO(), "/test/key3", []byte{})
	count := 0
	err = rep.WalkEntry(context.TODO(), "/test", func(key, value []byte) {
		count++
	})
	assert.NoError(t, err)
	assert.Equal(t, 2, count)
}

func TestNew(t *testing.T) {
	cfg := &config.RepoState{}
	_, err := newEtcdRepository(cfg, "nobody")
	assert.Error(t, err)
}

func TestHeartBeat(t *testing.T) {
	cluster := mock.StartEtcdCluster(t, "http://localhost:8702")
	defer cluster.Terminate(t)

	cfg := &config.RepoState{
		Endpoints: cluster.Endpoints,
	}
	b, err := newEtcdRepository(cfg, "nobody")
	assert.NoError(t, err)

	repo := b.(*etcdRepository)
	repo.timeout = time.Second * 10

	ip, _ := hostutil.GetHostIP()
	heartbeat := fmt.Sprintf("/cluster1/storage/heartbeat/%s:%d", ip, 2918)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	var ch <-chan Closed
	ch, err = b.Heartbeat(ctx, heartbeat, []byte("test"), 1)
	assert.NoError(t, err)

	_, err = b.Get(ctx, heartbeat)
	assert.NoError(t, err)

	cancel()
	time.Sleep(time.Second)
	_, err = b.Get(ctx, heartbeat)
	assert.Error(t, err)
	select {
	case <-ch:
	case <-time.After(500 * time.Millisecond):
		t.Fatal("cancel heartbeat timeout")
	}
}

func TestWatch(t *testing.T) {
	cluster := mock.StartEtcdCluster(t, "http://localhost:8703")
	defer cluster.Terminate(t)

	cfg := &config.RepoState{
		Endpoints: cluster.Endpoints,
	}
	b, _ := newEtcdRepository(cfg, "nobody")
	ctx, cancel := context.WithCancel(context.Background())
	repo := b.(*etcdRepository)
	repo.timeout = time.Second * 10
	// test watch no exist path
	ch := b.Watch(ctx, "/cluster1/controller/1", true)
	assert.NotNil(t, ch)

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
	assert.Len(t, val, 2)
	assert.Equal(t, "1", val["/cluster1/controller/1"])
	assert.Equal(t, "222", val["/cluster1/controller/2"])
}

func TestGetWatchPrefix(t *testing.T) {
	cluster := mock.StartEtcdCluster(t, "http://localhost:8704")
	defer cluster.Terminate(t)

	cfg := &config.RepoState{
		Endpoints: cluster.Endpoints,
	}
	b, _ := newEtcdRepository(cfg, "nobody")
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
	assert.Equal(t, "3", string(bytes1))
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
			assert.False(t, allEvt)
			assert.False(t, modifyEvt)
			assert.False(t, deleteEvt)
			assert.Len(t, kvs, 2)
			assert.Equal(t, "1", kvs["/lindb/broker/1"])
			assert.Equal(t, "2", kvs["/lindb/broker/2"])
			allEvt = true
		case EventTypeModify:
			assert.True(t, allEvt)
			assert.False(t, modifyEvt)
			assert.False(t, deleteEvt)
			assert.Len(t, kvs, 1)
			assert.Equal(t, "3", kvs["/lindb/broker/3"])
			modifyEvt = true
		case EventTypeDelete:
			assert.True(t, allEvt)
			assert.True(t, modifyEvt)
			assert.False(t, deleteEvt)
			assert.Len(t, kvs, 1)
			_, ok := kvs["/lindb/broker/3"]
			assert.True(t, ok)
			deleteEvt = true
			cancel()
		}
	}
	assert.True(t, deleteEvt)
}

func TestElect(t *testing.T) {
	cluster := mock.StartEtcdCluster(t, "http://localhost:8705")
	defer cluster.Terminate(t)

	cfg := &config.RepoState{
		Endpoints: cluster.Endpoints,
	}
	b, _ := newEtcdRepository(cfg, "nobody")
	repo := b.(*etcdRepository)
	repo.timeout = time.Second * 10

	ctx, cancel := context.WithCancel(context.Background())
	// the key should not exist,it must be success
	success, ch, err := b.Elect(ctx, "/lindb/broker/master", []byte("test"), 1)
	assert.NoError(t, err)
	assert.True(t, success)
	time.Sleep(2 * time.Second)
	bytes, err := b.Get(context.TODO(), "/lindb/broker/master")
	assert.NoError(t, err)
	assert.Equal(t, "test", string(bytes))

	ctx2, cancel2 := context.WithCancel(context.Background())
	shouldFalse, _, _ := b.Elect(ctx2, "/lindb/broker/master", []byte("test2"), 1)
	if cancel2 != nil {
		cancel2()
	}
	assert.False(t, shouldFalse)
	cancel()
	select {
	case <-ch:
	case <-time.After(500 * time.Millisecond):
		t.Fatal("cancel heartbeat timeout")
	}

	time.Sleep(3 * time.Second)

	_, err = b.Get(context.TODO(), "/lindb/broker/master")
	assert.Error(t, err)

	ctx3, cancel3 := context.WithCancel(context.Background())
	shouldSuccess, _, _ := b.Elect(ctx3, "/lindb/broker/master", []byte("test3"), 1)
	assert.True(t, shouldSuccess)

	bytes3, _ := b.Get(context.TODO(), "/lindb/broker/master")

	assert.Equal(t, "test3", string(bytes3))

	cancel3()
}

func TestBatch(t *testing.T) {
	cluster := mock.StartEtcdCluster(t, "http://localhost:8706")
	defer cluster.Terminate(t)

	cfg := &config.RepoState{
		Namespace: "/test/batch",
		Endpoints: cluster.Endpoints,
	}
	b, _ := newEtcdRepository(cfg, "nobody")
	repo := b.(*etcdRepository)
	repo.timeout = time.Second * 10
	batch := Batch{
		KVs: []KeyValue{
			{"key1", []byte("value1")},
			{"key2", []byte("value2")},
			{"key3", []byte("value3")},
		}}
	success, _ := b.Batch(context.TODO(), batch)
	assert.True(t, success)

	list, _ := b.List(context.TODO(), "key")
	assert.Len(t, list, 3)
}

func TestTransaction(t *testing.T) {
	cluster := mock.StartEtcdCluster(t, "http://localhost:8707")
	defer cluster.Terminate(t)

	cfg := &config.RepoState{
		Namespace: "/test/batch",
		Endpoints: cluster.Endpoints,
	}
	b, _ := newEtcdRepository(cfg, "nobody")
	repo := b.(*etcdRepository)
	repo.timeout = time.Second * 10

	txn := b.NewTransaction()
	txn.Put("test", []byte("value"))
	err := b.Commit(context.TODO(), txn)
	assert.NoError(t, err)

	v, _ := b.Get(context.TODO(), "test")
	assert.Equal(t, []byte("value"), v)

	txn = b.NewTransaction()
	txn.ModRevisionCmp("key", "=", 0)
	txn.Put("test", []byte("value2"))
	err = b.Commit(context.TODO(), txn)
	assert.NoError(t, err)
	v, _ = b.Get(context.TODO(), "test")
	assert.Equal(t, []byte("value2"), v)

	txn = b.NewTransaction()
	txn.ModRevisionCmp("key", "=", 33)
	txn.Delete("test")
	err = b.Commit(context.TODO(), txn)
	assert.Error(t, err)

	v, _ = b.Get(context.TODO(), "test")
	assert.Equal(t, []byte("value2"), v)

	txn = b.NewTransaction()
	txn.ModRevisionCmp("key", "=", 0)
	txn.Delete("test")
	err = b.Commit(context.TODO(), txn)
	assert.NoError(t, err)
	_, err = b.Get(context.TODO(), "test")
	assert.Error(t, err)

	assert.Error(t, TxnErr(nil, fmt.Errorf("err")))
}

func TestNextSequence(t *testing.T) {
	cluster := mock.StartEtcdCluster(t, "http://localhost:8708")
	defer cluster.Terminate(t)

	cfg := &config.RepoState{
		Endpoints: cluster.Endpoints,
	}
	b, _ := newEtcdRepository(cfg, "nobody")
	repo := b.(*etcdRepository)
	repo.timeout = time.Second * 10

	seq, err := repo.NextSequence(context.TODO(), "/test/seq")
	assert.NoError(t, err)
	assert.Equal(t, int64(1), seq)
	for i := 2; i < 10; i++ {
		seq, err = repo.NextSequence(context.TODO(), "/test/seq")
		assert.NoError(t, err)
		assert.Equal(t, int64(i), seq)
	}
}

func TestNextSequence_Concurrency(t *testing.T) {
	cluster := mock.StartEtcdCluster(t, "http://localhost:8709")
	defer cluster.Terminate(t)

	cfg := &config.RepoState{
		Endpoints: cluster.Endpoints,
	}
	b, _ := newEtcdRepository(cfg, "nobody")
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
				assert.NoError(t, err)
				rs.Store(seq, seq)
			}
		}()
	}
	wait.Wait()
	for i := 1; i <= 100; i++ {
		_, ok := rs.Load(int64(i))
		assert.True(t, ok)
	}
}
