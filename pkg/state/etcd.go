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
	"path/filepath"
	"strconv"
	"strings"
	"time"

	etcdcliv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/concurrency"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc"

	"github.com/lindb/lindb/config"
	"github.com/lindb/lindb/pkg/logger"
)

// etcdRepository is repository based on etcd storage
type etcdRepository struct {
	namespace string
	client    *etcdcliv3.Client
	logger    *logger.Logger
	timeout   time.Duration
}

// newEtcdRepository creates a new repository based on etcd storage
func newEtcdRepository(repoState *config.RepoState, owner string) (Repository, error) {
	zapCfg := zap.NewProductionConfig()
	zapCfg.Level = zap.NewAtomicLevelAt(zapcore.ErrorLevel)
	cfg := etcdcliv3.Config{
		Endpoints:            repoState.Endpoints,
		DialTimeout:          repoState.DialTimeout.Duration(),
		DialKeepAliveTime:    repoState.DialTimeout.Duration(),
		DialKeepAliveTimeout: repoState.DialTimeout.Duration(),
		DialOptions:          []grpc.DialOption{grpc.WithBlock()},
		Username:             repoState.Username,
		Password:             repoState.Password,
		LogConfig:            &zapCfg,
	}
	cli, err := etcdcliv3.New(cfg)
	if err != nil {
		return nil, fmt.Errorf("create etcd client error:%s", err)
	}

	repo := etcdRepository{
		namespace: repoState.Namespace,
		client:    cli,
		timeout:   repoState.Timeout.Duration(),
		logger:    logger.GetLogger(owner, "ETCD")}

	repo.logger.Info("new etcd client successfully",
		logger.Any("endpoints", repoState.Endpoints))
	return &repo, nil
}

// Get retrieves value for given key from etcd
func (r *etcdRepository) Get(ctx context.Context, key string) ([]byte, error) {
	thisCtx, cancelFunc := context.WithTimeout(ctx, r.timeout)
	defer cancelFunc()
	resp, err := r.get(thisCtx, key)
	if err != nil {
		return nil, err
	}
	return r.getValue(key, resp)
}

// List retrieves list for given prefix from etcd
func (r *etcdRepository) List(ctx context.Context, prefix string) ([]KeyValue, error) {
	thisCtx, cancelFunc := context.WithTimeout(ctx, r.timeout)
	defer cancelFunc()
	resp, err := r.client.Get(thisCtx, r.keyPath(prefix), etcdcliv3.WithPrefix())
	if err != nil {
		return nil, err
	}
	var result []KeyValue

	if len(resp.Kvs) > 0 {
		for _, kv := range resp.Kvs {
			if len(kv.Value) > 0 {
				result = append(result, KeyValue{Key: r.parseKey(string(kv.Key)), Value: kv.Value})
			}
		}
	}
	return result, nil
}

// WalkEntry walks each kv entry via fn for given prefix from repository.
func (r *etcdRepository) WalkEntry(ctx context.Context, prefix string, fn func(key, value []byte)) error {
	thisCtx, cancelFunc := context.WithTimeout(ctx, r.timeout)
	defer cancelFunc()
	resp, err := r.client.Get(thisCtx, r.keyPath(prefix), etcdcliv3.WithPrefix())
	if err != nil {
		return err
	}
	if len(resp.Kvs) > 0 {
		for _, kv := range resp.Kvs {
			if len(kv.Value) > 0 {
				fn(kv.Key, kv.Value)
			}
		}
	}
	return nil
}

// Put puts a key-value pair into etcd
func (r *etcdRepository) Put(ctx context.Context, key string, val []byte) error {
	thisCtx, cancelFunc := context.WithTimeout(ctx, r.timeout)
	defer cancelFunc()
	_, err := r.client.Put(thisCtx, r.keyPath(key), string(val))
	if err != nil {
		r.logger.Error("put error", logger.String("path", key),
			logger.String("namespace", r.namespace),
			logger.Error(err))
	}
	return err
}

func (r *etcdRepository) PutWithTX(ctx context.Context, key string, val []byte, check func(oldVal []byte) error) (bool, error) {
	thisCtx, cancelFunc := context.WithTimeout(ctx, r.timeout)
	defer cancelFunc()
	keyPath := r.keyPath(key)
	resp, err := r.client.Get(thisCtx, keyPath)
	if err != nil {
		return false, err
	}
	var txn etcdcliv3.Txn
	if len(resp.Kvs) > 0 {
		old := resp.Kvs[0]
		if check != nil {
			err = check(old.Value)
			if err != nil {
				return false, err
			}
		}
		txn = r.client.Txn(thisCtx).If(etcdcliv3.Compare(etcdcliv3.ModRevision(keyPath), "=", old.ModRevision))
	} else {
		txn = r.client.Txn(thisCtx).If(etcdcliv3.Compare(etcdcliv3.CreateRevision(keyPath), "=", 0))
	}
	putResp, err := txn.Then(etcdcliv3.OpPut(keyPath, string(val))).Commit()
	if err != nil {
		return false, err
	}
	return putResp.Succeeded, nil
}

// Delete deletes value for given key from etcd
func (r *etcdRepository) Delete(ctx context.Context, key string) error {
	thisCtx, cancelFunc := context.WithTimeout(ctx, r.timeout)
	defer cancelFunc()
	_, err := r.client.Delete(thisCtx, r.keyPath(key))
	return err
}

// Close closes etcd client
func (r *etcdRepository) Close() error {
	return r.client.Close()
}

// Heartbeat does heartbeat on the key with a value and ttl based on etcd
func (r *etcdRepository) Heartbeat(ctx context.Context, key string, value []byte, ttl int64) (<-chan Closed, error) {
	h := newHeartbeat(r.client, r.keyPath(key), value, ttl, false)
	h.withLogger(r.logger)
	_, err := h.grantKeepAliveLease(ctx)
	if err != nil {
		return nil, err
	}
	ch := make(chan Closed)
	// do keepalive/retry background
	go func() {
		// closed channel, if keep alive stopped
		defer close(ch)
		h.keepAlive(ctx)
	}()
	return ch, nil
}

// Elect puts a key with a value.it will be success
// if the key does not exist,otherwise it will be failed.When this
// operation success,it will do keepalive background
func (r *etcdRepository) Elect(
	ctx context.Context, key string,
	value []byte, ttl int64,
) (ok bool, closed <-chan Closed, err error) {
	h := newHeartbeat(r.client, r.keyPath(key), value, ttl, true)
	h.withLogger(r.logger)
	success, err := h.grantKeepAliveLease(ctx)
	if err != nil {
		return false, nil, err
	}
	// when put success,do keep alive
	if success {
		ch := make(chan Closed)
		// do keepalive/retry background
		go func() {
			// closed channel, if keep alive stopped
			defer func() {
				close(ch)
			}()
			h.keepAlive(ctx)
		}()
		return success, ch, nil
	}
	return success, nil, nil
}

// get returns response of get operator
func (r *etcdRepository) get(ctx context.Context, key string) (*etcdcliv3.GetResponse, error) {
	thisCtx, cancelFunc := context.WithTimeout(ctx, r.timeout)
	defer cancelFunc()
	resp, err := r.client.Get(thisCtx, r.keyPath(key))
	if err != nil {
		return nil, fmt.Errorf("get value failure for key[%s], error:%s", key, err)
	}
	return resp, nil
}

// getValue returns value of response which got.
func (r *etcdRepository) getValue(key string, resp *etcdcliv3.GetResponse) ([]byte, error) {
	if len(resp.Kvs) == 0 {
		return nil, ErrNotExist
	}

	firstKV := resp.Kvs[0]
	if len(firstKV.Value) == 0 {
		return nil, fmt.Errorf("key[%s]'s value is empty", key)
	}
	return firstKV.Value, nil
}

// Watch watches on a key. The watched events will be returned through the returned channel.
//
// NOTE: when caller meets EventTypeAll, it must clean all previous values, since it may contain
// deleted values we do not know.
func (r *etcdRepository) Watch(ctx context.Context, key string, fetchVal bool) WatchEventChan {
	watcher := newWatcher(ctx, r, r.keyPath(key), fetchVal)
	return watcher.EventC
}

// WatchPrefix watches on a prefix. All the changes who have the prefix
// will be notified through the WatchEventChan channel.
//
// NOTE: when caller meets EventTypeAll, it must clean all previous values, since it may contain
// deleted values we do not know.
func (r *etcdRepository) WatchPrefix(ctx context.Context, prefixKey string, fetchVal bool) WatchEventChan {
	watcher := newWatcher(ctx, r, r.keyPath(prefixKey), fetchVal, etcdcliv3.WithPrefix())
	return watcher.EventC
}

// Batch puts k/v list, this operation is atomic
func (r *etcdRepository) Batch(ctx context.Context, batch Batch) (bool, error) {
	var ops []etcdcliv3.Op
	for _, kv := range batch.KVs {
		ops = append(ops, etcdcliv3.OpPut(
			r.keyPath(kv.Key),
			string(kv.Value),
		))
	}

	resp, err := r.client.Txn(ctx).Then(ops...).Commit()
	if err != nil {
		return false, err
	}
	return resp.Succeeded, nil
}

// NewTransaction creates a new transaction
func (r *etcdRepository) NewTransaction() Transaction {
	return newTransaction(r)
}

// Commit commits the transaction, if fail return err
func (r *etcdRepository) Commit(ctx context.Context, txn Transaction) error {
	t, ok := txn.(*transaction)
	if !ok {
		return ErrTxnConvert
	}
	resp, err := r.client.Txn(ctx).If(t.cmps...).Then(t.ops...).Commit()
	return TxnErr(resp, err)
}

// NextSequence returns next sequence number.
func (r *etcdRepository) NextSequence(ctx context.Context, key string) (int64, error) {
	s, err := concurrency.NewSession(r.client) // explore options to pass
	if err != nil {
		return 0, err
	}

	m := concurrency.NewMutex(s, key)

	err = m.Lock(ctx)
	if err != nil {
		return 0, err
	}
	defer func() {
		_ = m.Unlock(ctx)
	}()

	resp, err := r.client.Get(ctx, key)
	if err != nil {
		return 0, err
	}
	var seq int64
	if resp.Count > 0 {
		seq, err = strconv.ParseInt(string(resp.OpResponse().Get().Kvs[0].Value), 10, 64)
		if err != nil {
			return 0, err
		}
		seq++
	} else {
		seq = 1 // init value
	}

	_, err = r.client.Put(ctx, key, strconv.FormatInt(seq, 10))
	if err != nil {
		return 0, err
	}
	return seq, nil
}

// keyPath return new key path with namespace prefix
func (r *etcdRepository) keyPath(key string) string {
	if len(r.namespace) > 0 {
		return filepath.Join(r.namespace, key)
	}
	return key
}

// parseKey parses the key, removes the namespace
func (r *etcdRepository) parseKey(key string) string {
	if r.namespace == "" {
		return key
	}
	return strings.Replace(key, r.namespace, "", 1)
}

type transaction struct {
	ops  []etcdcliv3.Op
	cmps []etcdcliv3.Cmp
	repo *etcdRepository
}

func newTransaction(repo *etcdRepository) Transaction {
	return &transaction{repo: repo}
}

func (t *transaction) ModRevisionCmp(key, op string, v interface{}) {
	t.cmps = append(t.cmps, etcdcliv3.Compare(etcdcliv3.ModRevision(t.repo.keyPath(key)), op, v))
}

func (t *transaction) Put(key string, value []byte) {
	t.ops = append(t.ops, etcdcliv3.OpPut(t.repo.keyPath(key), string(value)))
}

func (t *transaction) Delete(key string) {
	t.ops = append(t.ops, etcdcliv3.OpDelete(t.repo.keyPath(key)))
}
