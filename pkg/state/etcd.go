package state

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/lindb/lindb/config"
	"github.com/lindb/lindb/pkg/logger"

	etcdcliv3 "github.com/coreos/etcd/clientv3"
)

// etcdRepository is repository based on etcd storage
type etcdRepository struct {
	namespace string
	client    *etcdcliv3.Client
}

// newEtedRepository creates a new repository based on etcd storage
func newEtedRepository(repoState config.RepoState) (Repository, error) {
	cfg := etcdcliv3.Config{
		Endpoints: repoState.Endpoints,
		// DialTimeout: config.DialTimeout * time.Second,
	}
	cli, err := etcdcliv3.New(cfg)
	if err != nil {
		return nil, fmt.Errorf("create etc client error:%s", err)
	}
	logger.GetLogger("pkg/state", "ETCDRepository").Info("new etcd client successfully",
		logger.Any("endpoints", repoState.Endpoints))
	return &etcdRepository{
		namespace: repoState.Namespace,
		client:    cli,
	}, nil
}

// Get retrieves value for given key from etcd
func (r *etcdRepository) Get(ctx context.Context, key string) ([]byte, error) {
	resp, err := r.get(ctx, key)
	if err != nil {
		return nil, err
	}
	return r.getValue(key, resp)
}

// List retrieves list for given prefix from etcd
func (r *etcdRepository) List(ctx context.Context, prefix string) ([]KeyValue, error) {
	resp, err := r.client.Get(ctx, r.keyPath(prefix), etcdcliv3.WithPrefix())
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

// Put puts a key-value pair into etcd
func (r *etcdRepository) Put(ctx context.Context, key string, val []byte) error {
	_, err := r.client.Put(ctx, r.keyPath(key), string(val))
	return err
}

// Delete deletes value for given key from etcd
func (r *etcdRepository) Delete(ctx context.Context, key string) error {
	_, err := r.client.Delete(ctx, r.keyPath(key))
	return err
}

// Close closes etcd client
func (r *etcdRepository) Close() error {
	return r.client.Close()
}

// Heartbeat does heartbeat on the key with a value and ttl based on etcd
func (r *etcdRepository) Heartbeat(ctx context.Context, key string, value []byte, ttl int64) (<-chan Closed, error) {
	h := newHeartbeat(r.client, r.keyPath(key), value, ttl, false)
	_, err := h.grantKeepAliveLease(ctx)
	if err != nil {
		return nil, err
	}
	ch := make(chan Closed)
	// do keepalive/retry background
	go func() {
		// close closed channel, if keep alive stopped
		defer close(ch)
		h.keepAlive(ctx)
	}()
	return ch, nil
}

// Elect puts a key with a value.it will be success
// if the key does not exist,otherwise it will be failed.When this
// operation success,it will do keepalive background
func (r *etcdRepository) Elect(ctx context.Context, key string,
	value []byte, ttl int64) (bool, <-chan Closed, error) {
	h := newHeartbeat(r.client, r.keyPath(key), value, ttl, true)
	success, err := h.grantKeepAliveLease(ctx)
	if err != nil {
		return false, nil, err
	}
	// when put success,do keep alive
	if success {
		ch := make(chan Closed)
		// do keepalive/retry background
		go func() {
			// close closed channel, if keep alive stopped
			defer close(ch)
			h.keepAlive(ctx)
		}()
		return success, ch, nil
	}
	return success, nil, nil
}

// get returns response of get operator
func (r *etcdRepository) get(ctx context.Context, key string) (*etcdcliv3.GetResponse, error) {
	resp, err := r.client.Get(ctx, r.keyPath(key))
	if err != nil {
		return nil, fmt.Errorf("get value failure for key[%s], error:%s", key, err)
	}
	return resp, nil
}

// getValue returns value of get's response
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
// NOTE: when caller meets EventTypeAll, it must clean all previous values, since it may contains
// deleted values we do not know.
func (r *etcdRepository) Watch(ctx context.Context, key string, fetchVal bool) WatchEventChan {
	watcher := newWatcher(ctx, r, r.keyPath(key), fetchVal)
	return watcher.EventC
}

// WatchPrefix watches on a prefix.All of the changes who has the prefix
// will be notified through the WatchEventChan channel.
//
// NOTE: when caller meets EventTypeAll, it must clean all previous values, since it may contains
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

// keyPath return new key path with namespace prefix
func (r *etcdRepository) keyPath(key string) string {
	if len(r.namespace) > 0 {
		return filepath.Join(r.namespace, key)
	}
	return key
}

// parseKey parses the key, removes the namespace
func (r *etcdRepository) parseKey(key string) string {
	if len(r.namespace) == 0 {
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
